import { NextRequest, NextResponse } from 'next/server';

/**
 * Catch-all API Route — 실제 백엔드 연결 전까지 데모 데이터 서빙
 * Go Gateway가 연결되면 이 파일을 제거하고 next.config.mjs에서 rewrites로 프록시
 */

const MOCK_DATA: Record<string, () => unknown> = {
  'patients': () => [
    { id: 'pt-001', name: '김환자', age: 45, gender: 'M' },
    { id: 'pt-002', name: '이건강', age: 32, gender: 'F' },
    { id: 'pt-003', name: '박진료', age: 58, gender: 'M' },
    { id: 'pt-004', name: '최예방', age: 27, gender: 'F' },
    { id: 'pt-005', name: '정관리', age: 63, gender: 'M' },
  ],
  'patients/*/vitals': () => [
    { timestamp: '08:00', heartRate: 72, bloodPressureSys: 120, spo2: 98 },
    { timestamp: '09:00', heartRate: 75, bloodPressureSys: 122, spo2: 97 },
    { timestamp: '10:00', heartRate: 68, bloodPressureSys: 118, spo2: 99 },
    { timestamp: '11:00', heartRate: 80, bloodPressureSys: 125, spo2: 97 },
  ],
  'patients/*/biomarkers': () => [
    { biomarkerId: 'bm-001', name: 'HbA1c', value: 6.2, unit: '%', timestamp: '2026-05-22T08:00:00Z' },
    { biomarkerId: 'bm-002', name: 'LDL-C', value: 110, unit: 'mg/dL', timestamp: '2026-05-22T08:00:00Z' },
    { biomarkerId: 'bm-003', name: 'CRP', value: 0.8, unit: 'mg/L', timestamp: '2026-05-22T08:00:00Z' },
  ],
  'sensors': () => [
    { id: 'sns-001', type: 'electrochemical', status: 'active', battery: 95 },
    { id: 'sns-002', type: 'optical', status: 'active', battery: 87 },
    { id: 'sns-003', type: 'impedance', status: 'active', battery: 72 },
    { id: 'sns-004', type: 'temperature', status: 'active', battery: 99 },
    { id: 'sns-005', type: 'electrochemical', status: 'calibrating', battery: 63 },
    { id: 'sns-006', type: 'optical', status: 'inactive', battery: 12 },
    { id: 'sns-007', type: 'impedance', status: 'active', battery: 81 },
    { id: 'sns-008', type: 'temperature', status: 'active', battery: 90 },
  ],
  'agents': () => [
    { id: 'diag-nlp', name: '진단 NLP', model: 'BioBERT-KR v3.2', status: 'running', accuracy: 98.7, latency: 42, totalRequests: 12847 },
    { id: 'ecg-cls', name: 'ECG 분류', model: 'ResNet-ECG v2.1', status: 'running', accuracy: 99.2, latency: 18, totalRequests: 8923 },
    { id: 'xray-det', name: 'X-Ray 탐지', model: 'YOLO-Med v4.0', status: 'running', accuracy: 97.8, latency: 85, totalRequests: 3421 },
    { id: 'drug-int', name: '약물 상호작용', model: 'GNN-Drug v1.8', status: 'idle', accuracy: 96.1, latency: 120, totalRequests: 1876 },
    { id: 'risk-pred', name: '위험도 예측', model: 'XGBoost-Risk v5.0', status: 'running', accuracy: 94.5, latency: 35, totalRequests: 6234 },
    { id: 'ocr-rx', name: '처방전 OCR', model: 'TrOCR-Med v2.4', status: 'error', accuracy: 91.3, latency: 0, totalRequests: 425 },
  ],
  'agents/pipeline': () => ({
    steps: [
      { id: 1, name: '데이터 수집' },
      { id: 2, name: '전처리' },
      { id: 3, name: '특성 추출' },
      { id: 4, name: '모델 추론' },
      { id: 5, name: '후처리' },
      { id: 6, name: '결과 저장' },
    ],
    activeStep: 4,
    totalRequests: 33726,
  }),
  'predictor/risk-scores': () => [
    { disease: '제2형 당뇨', riskScore: 0.82, contributingFactors: ['BMI', 'HbA1c'], trend: 'up' },
    { disease: '관상동맥 질환', riskScore: 0.15, contributingFactors: ['Cholesterol'], trend: 'stable' },
    { disease: '고혈압', riskScore: 0.45, contributingFactors: ['Na+', 'Stress'], trend: 'down' },
    { disease: '만성 신장질환', riskScore: 0.08, contributingFactors: ['Creatinine'], trend: 'stable' },
    { disease: '간경변', riskScore: 0.12, contributingFactors: ['ALT', 'AST'], trend: 'stable' },
  ],
  'rewards/*/balance': () => ({
    userId: 'user-1', balance: 2847, currency: 'MPS', totalEarned: 12480, weeklyChange: 650,
  }),
  'rewards/*/contributions': () => [
    { id: '1', date: '2026-04-26', type: '심박 데이터', points: 150, quality: 'A+' },
    { id: '2', date: '2026-04-25', type: '혈압 데이터', points: 120, quality: 'A' },
    { id: '3', date: '2026-04-25', type: '수면 패턴', points: 200, quality: 'A+' },
    { id: '4', date: '2026-04-24', type: '활동량 데이터', points: 80, quality: 'B+' },
    { id: '5', date: '2026-04-23', type: '식이 기록', points: 100, quality: 'A' },
  ],
  'partners': () => [
    { id: '1', name: '서울대학교병원', protocol: 'FHIR R4', status: 'active', totalExchanges: 12847, lastSync: '2분 전' },
    { id: '2', name: '연세 세브란스', protocol: 'HL7 v2.5', status: 'active', totalExchanges: 8923, lastSync: '5분 전' },
    { id: '3', name: '삼성서울병원', protocol: 'FHIR R4', status: 'active', totalExchanges: 7456, lastSync: '8분 전' },
    { id: '4', name: '아산병원', protocol: 'FHIR R4', status: 'syncing', totalExchanges: 3421, lastSync: '동기화 중...' },
    { id: '5', name: '국립암센터', protocol: 'HL7 v2.5', status: 'inactive', totalExchanges: 1205, lastSync: '3시간 전' },
  ],
  'gxp/audit-logs': () => [
    { id: '1', user: '김품질', action: '배치 검사 승인', batchId: 'B-2026-0421', hasSignature: true, timestamp: '14:32:07' },
    { id: '2', user: '이감사', action: '원자재 입고 확인', batchId: 'B-2026-0420', hasSignature: true, timestamp: '11:15:33' },
    { id: '3', user: '박규정', action: '출하 승인', batchId: 'B-2026-0419', hasSignature: false, timestamp: '09:45:12' },
  ],
  'gxp/compliance': () => [
    { ruleId: 'KGMP 의약품 제조관리', status: 'pass', itemsCount: 41, passedCount: 41 },
    { ruleId: 'GLP 비임상시험관리', status: 'pass', itemsCount: 28, passedCount: 27 },
    { ruleId: 'GCP 임상시험관리', status: 'pass', itemsCount: 35, passedCount: 35 },
    { ruleId: 'GDP 유통관리', status: 'warning', itemsCount: 22, passedCount: 20 },
  ],
  'dev-portal/keys': () => [
    { id: '1', key: 'mmup_live_sk_x7k9...d6e8', name: 'Production Key', createdAt: '2026-03-15', status: 'active' },
    { id: '2', key: 'mmup_test_sk_a1b2...o5p6', name: 'Test Key', createdAt: '2026-04-01', status: 'active' },
    { id: '3', key: 'mmup_dev_sk_q1w2...g5h6', name: 'Dev Key', createdAt: '2026-01-10', status: 'revoked' },
  ],
  'dev-portal/endpoints': () => [
    { method: 'GET', path: '/api/v1/patients', description: '환자 목록 조회', auth: 'Bearer', status: 'stable' },
    { method: 'POST', path: '/api/v1/measurements', description: '측정 데이터 전송', auth: 'Bearer', status: 'stable' },
    { method: 'GET', path: '/api/v1/biomarkers/:id', description: '바이오마커 조회', auth: 'Bearer', status: 'stable' },
    { method: 'POST', path: '/api/v1/predictions/run', description: 'AI 예측 실행', auth: 'API Key', status: 'beta' },
    { method: 'GET', path: '/api/v1/rewards/balance', description: '토큰 잔액 조회', auth: 'Bearer', status: 'stable' },
    { method: 'POST', path: '/api/v1/fhir/bundle', description: 'FHIR 번들 전송', auth: 'mTLS', status: 'beta' },
  ],
  'hardware/devices': () => [
    { name: 'BLE GATT Service', status: 'Connected', latency: '12ms' },
    { name: 'UART Serial Bus', status: 'Ready', latency: '2ms' },
    { name: 'Edge TPU Accelerator', status: 'Active', latency: '45ms' },
  ],
  'hardware/diagnostics': () => ({
    engineVersion: 'v2.4.3', cpuTemp: 42.3, memoryMB: 15.4, securityStatus: 'IEC 62304 OK',
  }),
  'app/metrics': () => [
    { label: '심박수', value: '72', unit: 'bpm', trend: '정상', iconName: 'Heart', color: 'text-red-500 bg-red-50' },
    { label: '산소포화도', value: '98', unit: '%', trend: '정상', iconName: 'Droplets', color: 'text-sky-500 bg-sky-50' },
    { label: '걸음 수', value: '8,247', unit: '보', trend: '목표 82%', iconName: 'Footprints', color: 'text-emerald-500 bg-emerald-50' },
    { label: '수면', value: '7h 23m', unit: '', trend: '양호', iconName: 'Moon', color: 'text-violet-500 bg-violet-50' },
  ],
  'predictions/run': () => ({
    predictionId: 'pred-001', status: 'completed', riskScore: 0.73, confidence: 0.91,
    modelVersion: 'GNN-Risk v2.1', recommendations: ['혈당 모니터링 강화', '식이 조절 권장'],
  }),
  'measurements/history': () => [
    { id: '1', date: '2026-05-24', time: '14:32', biomarker: '혈당', value: 95, unit: 'mg/dL', status: 'normal', confidence: 0.96, trend: 'stable' },
    { id: '2', date: '2026-05-23', time: '08:15', biomarker: '혈당', value: 112, unit: 'mg/dL', status: 'normal', confidence: 0.94, trend: 'up' },
    { id: '3', date: '2026-05-22', time: '14:00', biomarker: '당화혈색소', value: 6.8, unit: '%', status: 'caution', confidence: 0.93, trend: 'stable' },
    { id: '4', date: '2026-05-21', time: '09:30', biomarker: '총 콜레스테롤', value: 195, unit: 'mg/dL', status: 'normal', confidence: 0.91, trend: 'down' },
    { id: '5', date: '2026-05-20', time: '08:00', biomarker: '혈당', value: 88, unit: 'mg/dL', status: 'normal', confidence: 0.95, trend: 'down' },
    { id: '6', date: '2026-05-19', time: '14:45', biomarker: '요산', value: 7.2, unit: 'mg/dL', status: 'caution', confidence: 0.92, trend: 'up' },
  ],
  'telemedicine/doctors': () => [
    { id: '1', name: '김내과', specialty: '내과', rating: 4.8, reviews: 234, available: true, nextSlot: '오늘 15:00', photo: '👨‍⚕️', fee: '₩25,000' },
    { id: '2', name: '이피부', specialty: '피부과', rating: 4.9, reviews: 189, available: true, nextSlot: '오늘 16:30', photo: '👩‍⚕️', fee: '₩30,000' },
    { id: '3', name: '박가정', specialty: '가정의학과', rating: 4.7, reviews: 312, available: false, nextSlot: '내일 09:00', photo: '👨‍⚕️', fee: '₩20,000' },
    { id: '4', name: '최소아', specialty: '소아과', rating: 4.9, reviews: 156, available: true, nextSlot: '오늘 17:00', photo: '👩‍⚕️', fee: '₩25,000' },
  ],
  'store/products': () => [
    { id: '1', name: '혈당 카트리지 (10개입)', category: '카트리지', price: 29000, rating: 4.8, reviews: 1234, image: '🩸', badge: '인기', subscription: true },
    { id: '2', name: '콜레스테롤 카트리지 (10개입)', category: '카트리지', price: 35000, rating: 4.7, reviews: 892, image: '💉', badge: null, subscription: true },
    { id: '3', name: '당화혈색소 카트리지 (5개입)', category: '카트리지', price: 45000, rating: 4.9, reviews: 567, image: '🔬', badge: '신규', subscription: true },
    { id: '4', name: '종합검사 카트리지 세트', category: '카트리지', price: 89000, rating: 4.8, reviews: 345, image: '📦', badge: '추천', subscription: false },
    { id: '5', name: 'MPS 리더기 Pro', category: '리더기', price: 290000, rating: 4.9, reviews: 2341, image: '📱', badge: '베스트', subscription: false },
    { id: '6', name: '채혈침 세트 (100개입)', category: '소모품', price: 12000, rating: 4.6, reviews: 3456, image: '💊', badge: null, subscription: false },
    { id: '7', name: '알코올 솜 (200매)', category: '소모품', price: 8000, rating: 4.5, reviews: 5678, image: '🧴', badge: null, subscription: false },
    { id: '8', name: '휴대용 파우치', category: '액세서리', price: 25000, rating: 4.7, reviews: 890, image: '👜', badge: null, subscription: false },
  ],
  'admin/kpis': () => [
    { label: '총 사용자', value: '12,847', change: '+3.2%', icon: 'Users', color: 'text-sky-600 bg-sky-50' },
    { label: '오늘 측정', value: '3,421', change: '+12%', icon: 'Activity', color: 'text-emerald-600 bg-emerald-50' },
    { label: '활성 리더기', value: '8,923', change: '+1.5%', icon: 'Package', color: 'text-purple-600 bg-purple-50' },
    { label: '이상 알림', value: '7', change: '-2건', icon: 'AlertTriangle', color: 'text-amber-600 bg-amber-50' },
  ],
  'admin/events': () => [
    { time: '14:32', type: '측정', desc: '서울 강남구 — 혈당 측정 이상치 감지', severity: 'warning' },
    { time: '14:28', type: '리더기', desc: 'MPK-A1B2C3 펌웨어 업데이트 완료', severity: 'info' },
    { time: '14:15', type: '사용자', desc: '신규 가입 +23명 (일일 목표 120%)', severity: 'success' },
    { time: '13:50', type: '재고', desc: '혈당 카트리지 재고 부족 경고 (창고 B)', severity: 'warning' },
    { time: '13:30', type: '결제', desc: '정기배송 결제 처리 완료 342건', severity: 'info' },
    { time: '12:45', type: '품질', desc: 'Lot-2026-A 보정 검증 통과', severity: 'success' },
  ],
};

function matchPath(pattern: string, path: string): boolean {
  const regex = new RegExp('^' + pattern.replace(/\*/g, '[^/]+') + '$');
  return regex.test(path);
}

function findMockData(path: string): unknown | null {
  // Exact match first
  if (MOCK_DATA[path]) return MOCK_DATA[path]();
  // Wildcard match
  for (const [pattern, factory] of Object.entries(MOCK_DATA)) {
    if (pattern.includes('*') && matchPath(pattern, path)) {
      return factory();
    }
  }
  return null;
}

export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  const { path } = await params;
  const apiPath = path.join('/');
  const data = findMockData(apiPath);

  if (data !== null) {
    return NextResponse.json({ data });
  }

  return NextResponse.json(
    { error: 'Not found', path: apiPath },
    { status: 404 }
  );
}

export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  const { path } = await params;
  const apiPath = path.join('/');
  const data = findMockData(apiPath);

  if (data !== null) {
    return NextResponse.json({ data });
  }

  return NextResponse.json(
    { error: 'Not found', path: apiPath },
    { status: 404 }
  );
}
