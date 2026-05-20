import 'package:isar/isar.dart';

part 'local_measurement.g.dart';

/// Flutter 앱 로컬 스토리지를 위한 오프라인 우선(CRDT) 측정 데이터 모델
/// Go 백엔드의 `models.go (Measurement)` 구조체와 1:1로 매핑됩니다.
@collection
class LocalMeasurement {
  Id id = Isar.autoIncrement; // Isar 내부 정수 키

  @Index(unique: true)
  late String clientLocalId; // 앱 자체 발급 UUID (서버 동기화 시 고유 식별자)

  String? serverId; // 서버 적재 완료 시 할당받는 백엔드 DB 측 UUID (매핑용)
  
  late String deviceMac;
  late DateTime measuredAt;

  // Rust Core Engine 출력값 (IEEE-754 배열)
  late List<double> diffSignal;   // 16채널 차분신호 
  late List<double> fingerprint;  // 896차원 확장 피처

  // AI Inference 출력값
  late int healthScore;
  late String riskLabel;

  // CRDT 및 오프라인-퍼스트 제어 변수
  @Index()
  bool isSynced = false; // 서버로의 전송 성공 여부 플래그
  DateTime? syncedAt;    // 최종 동기화(Ping) 시간
  
  // Soft Delete (휴지통 기능 및 서버 단 삭제 전파)
  bool isDeleted = false;
  late DateTime updatedAt;
}
