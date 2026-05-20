import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:isar/isar.dart';
import 'package:path_provider/path_provider.dart';

import 'local_measurement.dart';

/// 앱 전체에서 사용할 수 있는 전역 Isar DB 인스턴스 프로바이더
/// 비동기 초기화를 위해 FutureProvider로 구성 (앱 진입 지점인 main.dart 에서 대기)
final isarDatabaseProvider = FutureProvider<Isar>((ref) async {
  // 모바일/데스크탑 환경의 안전한 앱 전용 문서 디렉토리 확보
  final dir = await getApplicationDocumentsDirectory();
  
  // 앱 내 모든 Isar Collection(테이블)을 등록하여 DB 오픈
  final isar = await Isar.open(
    [
      LocalMeasurementSchema,
      // 차후 User 스키마, 기기 진단 스키마 등 50여 개 UI 관련 스키마 추가 예정
    ],
    directory: dir.path,
    name: 'manpasik_db',
  );

  // 캐싱 처리를 위한 리소스 해제(Dispose) 훅 등록
  ref.onDispose(() {
    if (isar.isOpen) {
      isar.close();
    }
  });

  return isar;
});
