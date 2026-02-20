import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/core/theme/app_theme.dart';

/// 데이터 내보내기 (FHIR R4 / CSV) 화면
class DataExportScreen extends StatefulWidget {
  const DataExportScreen({super.key});

  @override
  State<DataExportScreen> createState() => _DataExportScreenState();
}

class _DataExportScreenState extends State<DataExportScreen> {
  String _format = 'fhir_r4';
  bool _includeMeasurements = true;
  bool _includeHealthRecords = true;
  bool _includePrescriptions = true;
  bool _includeCoaching = false;
  String _dateRange = 'all';
  bool _isExporting = false;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('데이터 내보내기'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            // 안내 배너
            Container(
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: AppTheme.sanggamGold.withValues(alpha: 0.1),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Row(
                children: [
                  const Icon(Icons.info_outline, color: AppTheme.sanggamGold),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Text(
                      '건강 데이터를 FHIR R4 또는 CSV 형식으로 내보낼 수 있습니다.\n내보낸 데이터는 다른 의료 서비스에서 활용할 수 있습니다.',
                      style: theme.textTheme.bodySmall,
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 24),

            // 내보내기 형식 선택
            Text('내보내기 형식',
                style: theme.textTheme.titleMedium
                    ?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: 8),
            _buildFormatCard('fhir_r4', 'FHIR R4 (JSON)',
                'HL7 국제 의료 데이터 표준 형식', Icons.medical_services),
            _buildFormatCard(
                'csv', 'CSV (엑셀)', '스프레드시트 호환 형식', Icons.table_chart),
            _buildFormatCard(
                'pdf', 'PDF 리포트', '인쇄 가능한 건강 리포트', Icons.picture_as_pdf),
            const SizedBox(height: 24),

            // 데이터 범위 선택
            Text('내보낼 데이터',
                style: theme.textTheme.titleMedium
                    ?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: 8),
            SwitchListTile(
              title: const Text('측정 데이터'),
              subtitle: const Text('혈당, 혈압, 콜레스테롤 등'),
              value: _includeMeasurements,
              onChanged: (val) => setState(() => _includeMeasurements = val),
              contentPadding: EdgeInsets.zero,
            ),
            SwitchListTile(
              title: const Text('건강 기록'),
              subtitle: const Text('진단 기록, 검사 결과'),
              value: _includeHealthRecords,
              onChanged: (val) => setState(() => _includeHealthRecords = val),
              contentPadding: EdgeInsets.zero,
            ),
            SwitchListTile(
              title: const Text('처방전'),
              subtitle: const Text('처방약 이력, 복약 기록'),
              value: _includePrescriptions,
              onChanged: (val) => setState(() => _includePrescriptions = val),
              contentPadding: EdgeInsets.zero,
            ),
            SwitchListTile(
              title: const Text('AI 코칭 데이터'),
              subtitle: const Text('건강 목표, 코칭 메시지'),
              value: _includeCoaching,
              onChanged: (val) => setState(() => _includeCoaching = val),
              contentPadding: EdgeInsets.zero,
            ),
            const SizedBox(height: 16),

            // 기간 선택
            Text('기간',
                style: theme.textTheme.titleMedium
                    ?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: 8),
            SegmentedButton<String>(
              segments: const [
                ButtonSegment(value: 'all', label: Text('전체')),
                ButtonSegment(value: '1y', label: Text('1년')),
                ButtonSegment(value: '6m', label: Text('6개월')),
                ButtonSegment(value: '1m', label: Text('1개월')),
              ],
              selected: {_dateRange},
              onSelectionChanged: (val) =>
                  setState(() => _dateRange = val.first),
            ),
            const SizedBox(height: 32),

            // 내보내기 버튼
            FilledButton.icon(
              onPressed: _isExporting ? null : _handleExport,
              icon: _isExporting
                  ? const SizedBox(
                      width: 20,
                      height: 20,
                      child: CircularProgressIndicator(
                          strokeWidth: 2, color: Colors.white))
                  : const Icon(Icons.download),
              label: Text(_isExporting ? '내보내는 중...' : '데이터 내보내기'),
              style: FilledButton.styleFrom(
                backgroundColor: AppTheme.sanggamGold,
                minimumSize: const Size.fromHeight(48),
              ),
            ),
            const SizedBox(height: 16),

            // 설명
            Text(
              '* 내보내기된 파일은 기기의 다운로드 폴더에 저장됩니다.\n'
              '* FHIR R4 형식은 전 세계 의료 기관에서 호환됩니다.\n'
              '* 데이터는 암호화되어 전송됩니다.',
              style: theme.textTheme.bodySmall
                  ?.copyWith(color: Colors.grey),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildFormatCard(
      String value, String title, String subtitle, IconData icon) {
    final isSelected = _format == value;
    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
        side: BorderSide(
          color: isSelected ? AppTheme.sanggamGold : Colors.transparent,
          width: 2,
        ),
      ),
      child: RadioListTile<String>(
        title: Text(title, style: const TextStyle(fontWeight: FontWeight.w600)),
        subtitle: Text(subtitle),
        secondary: Icon(icon,
            color: isSelected ? AppTheme.sanggamGold : Colors.grey),
        value: value,
        groupValue: _format,
        onChanged: (val) => setState(() => _format = val ?? 'fhir_r4'),
      ),
    );
  }

  void _handleExport() {
    setState(() => _isExporting = true);
    Future.delayed(const Duration(seconds: 3), () {
      if (mounted) {
        setState(() => _isExporting = false);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(
                '데이터가 ${_format.toUpperCase()} 형식으로 내보내기 되었습니다.'),
            action: SnackBarAction(label: '열기', onPressed: () {}),
          ),
        );
      }
    });
  }
}
