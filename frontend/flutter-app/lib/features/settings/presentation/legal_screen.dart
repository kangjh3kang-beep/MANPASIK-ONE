import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

/// 이용약관 / 개인정보처리방침 화면
///
/// [type] 파라미터로 'terms' 또는 'privacy'를 받아 분기합니다.
class LegalScreen extends StatelessWidget {
  const LegalScreen({super.key, required this.type});

  final String type;

  bool get _isTerms => type == 'terms';

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: Text(_isTerms ? '이용약관' : '개인정보처리방침'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(24),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              _isTerms ? 'ManPaSik 서비스 이용약관' : 'ManPaSik 개인정보처리방침',
              style: theme.textTheme.headlineSmall?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              '시행일: 2026-01-01 | 버전: 1.0',
              style: theme.textTheme.bodySmall?.copyWith(
                color: theme.colorScheme.onSurfaceVariant,
              ),
            ),
            const Divider(height: 32),
            Text(
              _isTerms ? _termsContent : _privacyContent,
              style: theme.textTheme.bodyMedium?.copyWith(height: 1.8),
            ),
          ],
        ),
      ),
    );
  }
}

const _termsContent = '''
제1조 (목적)
본 약관은 만파식(ManPaSik) 서비스(이하 "서비스")를 제공하는 주식회사 만파식(이하 "회사")과 서비스를 이용하는 고객(이하 "회원") 간의 권리, 의무 및 책임사항을 규정함을 목적으로 합니다.

제2조 (정의)
① "서비스"란 회사가 제공하는 AI 기반 건강관리 플랫폼 및 관련 부가 서비스를 의미합니다.
② "리더기"란 ManPaSik 전용 차등측정 디바이스를 의미합니다.
③ "카트리지"란 리더기에 장착하여 특정 바이오마커를 측정하는 소모품을 의미합니다.
④ "측정 데이터"란 리더기를 통해 수집된 차등 신호 및 AI 분석 결과를 의미합니다.

제3조 (서비스의 제공)
① 회사는 다음과 같은 서비스를 제공합니다:
  1. 차등측정 기반 건강 바이오마커 분석
  2. AI 건강 코칭 및 맞춤형 추천
  3. 카트리지 마켓플레이스
  4. 비대면 진료 연결 (의료기관 연계)
  5. 가족 건강 관리
  6. 커뮤니티 서비스

제4조 (회원가입)
① 서비스 이용을 위해 회원가입이 필요합니다.
② 회원은 실명 및 실제 정보를 제공해야 합니다.
③ 14세 미만의 아동은 법정대리인의 동의가 필요합니다.

제5조 (구독 서비스)
① 무료(Free), 기본(Basic), 프로(Pro), 클리니컬(Clinical) 4개 등급으로 구분됩니다.
② 유료 구독은 월 단위로 자동 갱신되며, 해지 시 잔여 기간까지 이용 가능합니다.
③ 구독 등급별 제공 기능은 마켓 내 구독 관리 페이지에서 확인할 수 있습니다.

제6조 (측정 데이터의 처리)
① 측정 데이터는 회원 본인의 건강관리 목적으로만 사용됩니다.
② 데이터는 AES-256 암호화로 저장되며, TLS 1.3으로 전송됩니다.
③ 회원은 언제든 데이터 삭제를 요청할 수 있습니다.

제7조 (면책 조항)
① 본 서비스의 측정 결과는 의학적 진단을 대체하지 않습니다.
② 건강 이상 소견 시 반드시 전문 의료기관을 방문하시기 바랍니다.
③ 카트리지의 올바른 사용법을 준수하지 않아 발생한 부정확한 결과에 대해 회사는 책임지지 않습니다.

제8조 (분쟁 해결)
본 약관과 관련된 분쟁은 대한민국 법률에 따라 해결하며, 관할 법원은 서울중앙지방법원으로 합니다.
''';

const _privacyContent = '''
1. 개인정보의 수집 및 이용 목적
회사는 다음 목적을 위해 개인정보를 수집·이용합니다:
  - 서비스 제공 및 회원 관리
  - 건강 측정 데이터 분석 및 맞춤형 건강 코칭
  - AI 모델 개선 (익명화 데이터)
  - 고객 지원 및 불만 처리
  - 법적 의무 이행

2. 수집하는 개인정보 항목
① 필수 항목: 이메일, 비밀번호, 이름
② 선택 항목: 생년월일, 성별, 혈액형, 키, 체중, 기저질환, 알레르기
③ 자동 수집: 기기 정보, 앱 사용 기록, 측정 데이터, IP 주소

3. 개인정보의 보유 및 이용 기간
① 회원 탈퇴 시 즉시 파기 (단, 법령에 따른 보존 의무 항목 제외)
② 건강 측정 데이터: 탈퇴 후 90일 이내 파기 (GDPR Art.17 준수)
③ 감사 로그: 10년 보관 (IEC 62304, FDA 21 CFR Part 11)
④ 결제 기록: 5년 보관 (전자상거래법)

4. 개인정보의 제3자 제공
① 원칙적으로 회원의 동의 없이 제3자에게 제공하지 않습니다.
② 예외: 법령에 따른 요청, 회원의 명시적 동의, 응급 상황 (긴급 연락망)
③ 가족 건강 공유: 회원의 명시적 동의 하에 가족 구성원에게만 제한적 공유

5. 개인정보의 안전성 확보 조치
① 암호화: AES-256-GCM (저장), TLS 1.3 (전송)
② 접근 통제: RBAC 기반 역할별 접근 권한
③ PHI 접근 감사: 모든 건강정보 접근 기록 10년 보관
④ 개인정보 영향평가: 연 1회 실시
⑤ 규정 준수: PIPA(한국), HIPAA(미국), GDPR(EU), PIPL(중국), APPI(일본)

6. 정보주체의 권리
① 열람권: 수집된 개인정보의 열람을 요청할 수 있습니다.
② 정정권: 부정확한 정보의 정정을 요청할 수 있습니다.
③ 삭제권: 개인정보의 삭제를 요청할 수 있습니다.
④ 처리정지권: 개인정보 처리의 정지를 요청할 수 있습니다.
⑤ 이동권: 개인정보를 FHIR R4 형식으로 내보낼 수 있습니다.

7. 개인정보 보호책임자
  - 이름: 개인정보보호팀
  - 이메일: privacy@manpasik.com
  - 전화: 1588-0000
''';
