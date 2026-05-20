import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:isar/isar.dart';

import '../../../core/database/isar_provider.dart';
import '../../../core/database/local_measurement.dart';

/// Flutter 앱 내부의 Riverpod 기반 로컬 데이터베이스 레포지토리
/// 이 레포지토리 덕분에 폰이 오프라인이어도 측정이 저장/차트표시 가능합나다.
final measurementLocalRepoProvider = FutureProvider<MeasurementLocalRepository>((ref) async {
  final isar = await ref.watch(isarDatabaseProvider.future);
  return MeasurementLocalRepository(isar);
});

class MeasurementLocalRepository {
  final Isar isar;

  MeasurementLocalRepository(this.isar);

  /// 16채널 바이너리와 896차원 핑거프린트 결괏값을 로컬 폰 용량에 초고속 적재 (CRDT 준비)
  Future<int> saveMeasurement(LocalMeasurement measurement) async {
    return await isar.writeTxn(() async {
      return await isar.localMeasurements.put(measurement);
    });
  }

  /// 메인 홈 화면 차트에서 보여주기 위해 최근 측정 데이터를 불러옵니다.
  Future<List<LocalMeasurement>> getRecentMeasurements({int limit = 50}) async {
    return await isar.localMeasurements
        .where()
        .sortByMeasuredAtDesc()
        .limit(limit)
        .findAll();
  }

  /// 서버 통신 복구 시, 백엔드(Go)로 쏘기 위해 안넘어간(isSynced=false) 데이터 색출
  Future<List<LocalMeasurement>> getUnsyncedMeasurements() async {
    return await isar.localMeasurements
        .filter()
        .isSyncedEqualTo(false)
        .findAll();
  }

  /// Go 서버 적재 완료 시 Sync 상태 동기화
  Future<void> markAsSynced(int localId, String backendServerId) async {
    await isar.writeTxn(() async {
      final item = await isar.localMeasurements.get(localId);
      if (item != null) {
        item.isSynced = true;
        item.serverId = backendServerId;
        item.syncedAt = DateTime.now();
        await isar.localMeasurements.put(item);
      }
    });
  }
}
