import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/app_theme.dart';

/// 규제 준수 체크리스트 화면
class AdminComplianceScreen extends ConsumerWidget {
  const AdminComplianceScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('규제 준수 현황'),
        leading: IconButton(icon: const Icon(Icons.arrow_back), onPressed: () => context.pop()),
      ),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          // 종합 준수율
          Card(
            color: Colors.green.withOpacity(0.05),
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Row(
                children: [
                  Stack(
                    alignment: Alignment.center,
                    children: [
                      SizedBox(width: 60, height: 60, child: CircularProgressIndicator(value: 0.92, strokeWidth: 6, color: Colors.green, backgroundColor: Colors.green.withOpacity(0.1))),
                      const Text('92%', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
                    ],
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text('종합 준수율', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                        Text('46/50 항목 충족', style: theme.textTheme.bodySmall),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),

          // 규제별 상세
          ..._regulations.map((reg) => _buildRegulationCard(theme, reg)),
        ],
      ),
    );
  }

  Widget _buildRegulationCard(ThemeData theme, _RegulationData reg) {
    final progress = reg.passed / reg.total;
    final color = progress >= 0.9 ? Colors.green : progress >= 0.7 ? Colors.orange : Colors.red;

    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: ExpansionTile(
        leading: Icon(reg.icon, color: AppTheme.sanggamGold),
        title: Text(reg.name, style: const TextStyle(fontWeight: FontWeight.w600)),
        subtitle: Row(
          children: [
            Expanded(
              child: LinearProgressIndicator(value: progress, color: color, backgroundColor: color.withOpacity(0.1)),
            ),
            const SizedBox(width: 8),
            Text('${reg.passed}/${reg.total}', style: theme.textTheme.bodySmall),
          ],
        ),
        children: reg.items.map((item) => ListTile(
          dense: true,
          leading: Icon(
            item.passed ? Icons.check_circle : Icons.cancel,
            size: 18,
            color: item.passed ? Colors.green : Colors.red,
          ),
          title: Text(item.label, style: theme.textTheme.bodySmall),
        )).toList(),
      ),
    );
  }

  static final _regulations = [
    _RegulationData('GDPR (EU)', Icons.public, 12, 12, [
      _CheckItem('데이터 수집 동의', true), _CheckItem('잊힐 권리 구현', true),
      _CheckItem('데이터 이동권', true), _CheckItem('DPO 지정', true),
      _CheckItem('데이터 침해 통지 절차', true), _CheckItem('프라이버시 영향 평가', true),
    ]),
    _RegulationData('PIPA (한국)', Icons.flag, 10, 10, [
      _CheckItem('개인정보처리방침 공개', true), _CheckItem('수집·이용 동의', true),
      _CheckItem('제3자 제공 동의', true), _CheckItem('파기 절차', true),
      _CheckItem('안전성 확보 조치', true),
    ]),
    _RegulationData('HIPAA (미국)', Icons.health_and_safety, 14, 12, [
      _CheckItem('PHI 암호화 (전송)', true), _CheckItem('PHI 암호화 (저장)', true),
      _CheckItem('접근 제어', true), _CheckItem('감사 로그', true),
      _CheckItem('비상 접근 절차', false), _CheckItem('비즈니스 연관 계약', false),
    ]),
    _RegulationData('ISO 13485', Icons.verified, 10, 8, [
      _CheckItem('품질 관리 시스템', true), _CheckItem('설계 및 개발 제어', true),
      _CheckItem('위험 관리 (ISO 14971)', true), _CheckItem('추적성', true),
      _CheckItem('CAPA 프로세스', false), _CheckItem('내부 감사 일정', false),
    ]),
    _RegulationData('IEC 62304', Icons.code, 8, 8, [
      _CheckItem('소프트웨어 안전 분류', true), _CheckItem('개발 계획', true),
      _CheckItem('요구사항 분석', true), _CheckItem('검증/확인(V&V)', true),
    ]),
  ];
}

class _RegulationData {
  final String name;
  final IconData icon;
  final int total, passed;
  final List<_CheckItem> items;
  const _RegulationData(this.name, this.icon, this.total, this.passed, this.items);
}

class _CheckItem {
  final String label;
  final bool passed;
  const _CheckItem(this.label, this.passed);
}
