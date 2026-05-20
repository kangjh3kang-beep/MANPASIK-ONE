import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:grpc/grpc.dart';
import 'dart:io';

import 'pb/measurement.pbgrpc.dart';
import '../../features/measure/data/measurement_local_repository.dart';

// gRPC 싱글톤 클라이언트 프로바이더
final grpcClientProvider = Provider<MeasurementServiceClient>((ref) {
  // 에뮬레이터 환경과 실기기/로컬 호스트를 모두 지원하는 IP 설정
  final host = Platform.isAndroid ? '10.0.2.2' : '127.0.0.1';
  
  final channel = ClientChannel(
    host,
    port: 50051,
    options: const ChannelOptions(credentials: ChannelCredentials.insecure()),
  );

  ref.onDispose(() {
    channel.shutdown();
  });

  return MeasurementServiceClient(channel);
});

// CRDT 핵심: 백그라운드 동기화 엔진 프로바이더
final syncEngineProvider = Provider<SyncEngine>((ref) {
  final grpcClient = ref.watch(grpcClientProvider);
  return SyncEngine(ref, grpcClient);
});

class SyncEngine {
  final Ref ref;
  final MeasurementServiceClient client;

  SyncEngine(this.ref, this.client);

  /// 로컬 Isar DB에서 isSynced == false 인 항목들을 찾아 Go 백엔드로 전송
  Future<void> runDataSync() async {
    final localRepo = await ref.read(measurementLocalRepoProvider.future);
    
    // 1. 미전송 데이터 수집
    final pendingRecords = await localRepo.getUnsyncedMeasurements();
    if (pendingRecords.isEmpty) {
      print('[SyncEngine] 동기화할 누락 데이터가 없습니다.');
      return;
    }

    print('[SyncEngine] ${pendingRecords.length}건의 데이터를 Go 백엔드로 Bulk 전송 시작...');

    // 2. gRPC Protobuf 변환
    final grpcRecords = pendingRecords.map((item) {
      final r = MeasurementRecord()
        ..id = item.clientLocalId
        ..deviceMac = item.deviceMac
        ..timestamp = Int64(item.measuredAt.millisecondsSinceEpoch)
        ..healthScore = item.healthScore
        ..riskLabel = item.riskLabel;
        
      r.diffSignal.addAll(item.diffSignal);
      r.fingerprint.addAll(item.fingerprint);
      return r;
    }).toList();

    final request = SyncMeasurementsRequest()..records.addAll(grpcRecords);

    // 3. gRPC 전송 및 응답 확인
    try {
      final response = await client.syncMeasurements(request);
      print('[SyncEngine] 동기화 성공: ${response.syncedCount}건 완료 (실패: ${response.failedIds.length}건)');

      // 4. 로컬 DB 상태 업데이트 (CRDT 동기화 확정)
      for (var record in pendingRecords) {
        if (!response.failedIds.contains(record.clientLocalId)) {
          // 백엔드의 실제 발급 UUID가 넘어온다면 여기서 매핑할 수 있음. 임시로 SUCCESS 문자열 할당.
          await localRepo.markAsSynced(record.id, "REMOTE_SUCCESS_ACK");
        }
      }
    } catch (e) {
      print('[SyncEngine] gRPC 서버 전송 오류. 네트워크가 복구되면 다시 시도합니다. 에러: $e');
    }
  }
}
