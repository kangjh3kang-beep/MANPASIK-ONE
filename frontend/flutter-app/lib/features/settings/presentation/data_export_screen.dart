import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// DataExportScreen — Sanggam Orbit 데이터 내보내기
//
// [Rule 4] app_theme.dart → sanggam_theme.dart
// [Rule 4] +sanggam_container.dart
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) 1x 제거
// [Rule 4] theme.textTheme.* ~5x → 직접 TextStyle
// [Rule 4] AppTheme.sanggamGold 4x → SanggamTheme.primary
// [Rule 4] Colors.grey → SanggamTheme.onSurfaceDim
// [Rule 4] Card → SanggamContainer
// [Rule 4] FilledButton foregroundColor + borderRadius:16
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] borderRadius:12→16, spacing:12→16
// ───────────────────────────────────────────────────

/// 데이터 내보내기 (FHIR R4 / CSV) 화면
class DataExportScreen extends ConsumerStatefulWidget {
  const DataExportScreen({super.key});

  @override
  ConsumerState<DataExportScreen> createState() =>
      _DataExportScreenState();
}

class _DataExportScreenState
    extends ConsumerState<DataExportScreen> {
  String _format = 'fhir_r4';
  bool _includeMeasurements = true;
  bool _includeHealthRecords = true;
  bool _includePrescriptions = true;
  bool _includeCoaching = false;
  String _dateRange = 'all';
  bool _isExporting = false;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor:
          SanggamTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            // 헤더
            Padding(
              padding:
                  const EdgeInsets.symmetric(
                      horizontal: 8,
                      vertical: 8),
              child: Row(
                children: [
                  IconButton(
                    icon: const Icon(
                        Icons.arrow_back,
                        color: Colors.white),
                    tooltip: '뒤로 가기',
                    onPressed: () =>
                        context.pop(),
                  ),
                  const SizedBox(width: 8),
                  const Expanded(
                    child: Text(
                      '데이터 내보내기',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight:
                            FontWeight.bold,
                      ),
                    ),
                  ),
                ],
              ),
            ),
            // 본문
            Expanded(
              child: SingleChildScrollView(
                padding:
                    const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment:
                      CrossAxisAlignment
                          .stretch,
                  children: [
                    // 안내 배너
                    Container(
                      padding:
                          const EdgeInsets.all(
                              16),
                      decoration:
                          BoxDecoration(
                        color: SanggamTheme
                            .primary
                            .withValues(
                                alpha: 0.1),
                        borderRadius:
                            BorderRadius
                                .circular(16),
                      ),
                      child: Row(
                        children: [
                          const Icon(
                              Icons
                                  .info_outline,
                              color:
                                  SanggamTheme
                                      .primary),
                          const SizedBox(
                              width: 16),
                          const Expanded(
                            child: Text(
                              '건강 데이터를 FHIR R4 또는 CSV 형식으로 내보낼 수 있습니다.\n내보낸 데이터는 다른 의료 서비스에서 활용할 수 있습니다.',
                              style: TextStyle(
                                color: SanggamTheme
                                    .onSurfaceDim,
                                fontSize: 12,
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: 24),

                    // 내보내기 형식 선택
                    const Text('내보내기 형식',
                        style: TextStyle(
                          color: Colors.white,
                          fontSize: 16,
                          fontWeight:
                              FontWeight.bold,
                        )),
                    const SizedBox(height: 8),
                    _buildFormatCard(
                        'fhir_r4',
                        'FHIR R4 (JSON)',
                        'HL7 국제 의료 데이터 표준 형식',
                        Icons
                            .medical_services),
                    _buildFormatCard(
                        'csv',
                        'CSV (엑셀)',
                        '스프레드시트 호환 형식',
                        Icons.table_chart),
                    _buildFormatCard(
                        'pdf',
                        'PDF 리포트',
                        '인쇄 가능한 건강 리포트',
                        Icons
                            .picture_as_pdf),
                    const SizedBox(height: 24),

                    // 데이터 범위 선택
                    const Text('내보낼 데이터',
                        style: TextStyle(
                          color: Colors.white,
                          fontSize: 16,
                          fontWeight:
                              FontWeight.bold,
                        )),
                    const SizedBox(height: 8),
                    SwitchListTile(
                      title: const Text(
                          '측정 데이터',
                          style: TextStyle(
                              color: Colors
                                  .white)),
                      subtitle: const Text(
                          '혈당, 혈압, 콜레스테롤 등',
                          style: TextStyle(
                            color: SanggamTheme
                                .onSurfaceDim,
                            fontSize: 12,
                          )),
                      value:
                          _includeMeasurements,
                      activeColor:
                          SanggamTheme.primary,
                      onChanged: (val) =>
                          setState(() =>
                              _includeMeasurements =
                                  val),
                      contentPadding:
                          EdgeInsets.zero,
                    ),
                    SwitchListTile(
                      title: const Text(
                          '건강 기록',
                          style: TextStyle(
                              color: Colors
                                  .white)),
                      subtitle: const Text(
                          '진단 기록, 검사 결과',
                          style: TextStyle(
                            color: SanggamTheme
                                .onSurfaceDim,
                            fontSize: 12,
                          )),
                      value:
                          _includeHealthRecords,
                      activeColor:
                          SanggamTheme.primary,
                      onChanged: (val) =>
                          setState(() =>
                              _includeHealthRecords =
                                  val),
                      contentPadding:
                          EdgeInsets.zero,
                    ),
                    SwitchListTile(
                      title: const Text('처방전',
                          style: TextStyle(
                              color: Colors
                                  .white)),
                      subtitle: const Text(
                          '처방약 이력, 복약 기록',
                          style: TextStyle(
                            color: SanggamTheme
                                .onSurfaceDim,
                            fontSize: 12,
                          )),
                      value:
                          _includePrescriptions,
                      activeColor:
                          SanggamTheme.primary,
                      onChanged: (val) =>
                          setState(() =>
                              _includePrescriptions =
                                  val),
                      contentPadding:
                          EdgeInsets.zero,
                    ),
                    SwitchListTile(
                      title: const Text(
                          'AI 코칭 데이터',
                          style: TextStyle(
                              color: Colors
                                  .white)),
                      subtitle: const Text(
                          '건강 목표, 코칭 메시지',
                          style: TextStyle(
                            color: SanggamTheme
                                .onSurfaceDim,
                            fontSize: 12,
                          )),
                      value:
                          _includeCoaching,
                      activeColor:
                          SanggamTheme.primary,
                      onChanged: (val) =>
                          setState(() =>
                              _includeCoaching =
                                  val),
                      contentPadding:
                          EdgeInsets.zero,
                    ),
                    const SizedBox(height: 16),

                    // 기간 선택
                    const Text('기간',
                        style: TextStyle(
                          color: Colors.white,
                          fontSize: 16,
                          fontWeight:
                              FontWeight.bold,
                        )),
                    const SizedBox(height: 8),
                    SegmentedButton<String>(
                      style: ButtonStyle(
                        foregroundColor:
                            WidgetStateProperty
                                .resolveWith(
                                    (states) =>
                                        states.contains(
                                                WidgetState
                                                    .selected)
                                            ? SanggamTheme
                                                .background
                                            : Colors
                                                .white),
                        backgroundColor:
                            WidgetStateProperty
                                .resolveWith(
                                    (states) =>
                                        states.contains(
                                                WidgetState
                                                    .selected)
                                            ? SanggamTheme
                                                .primary
                                            : SanggamTheme
                                                .surface),
                      ),
                      segments: const [
                        ButtonSegment(
                            value: 'all',
                            label:
                                Text('전체')),
                        ButtonSegment(
                            value: '1y',
                            label:
                                Text('1년')),
                        ButtonSegment(
                            value: '6m',
                            label:
                                Text('6개월')),
                        ButtonSegment(
                            value: '1m',
                            label:
                                Text('1개월')),
                      ],
                      selected: {_dateRange},
                      onSelectionChanged:
                          (val) => setState(
                              () => _dateRange =
                                  val.first),
                    ),
                    const SizedBox(height: 32),

                    // 내보내기 버튼
                    FilledButton.icon(
                      onPressed: _isExporting
                          ? null
                          : _handleExport,
                      icon: _isExporting
                          ? const SizedBox(
                              width: 20,
                              height: 20,
                              child:
                                  CircularProgressIndicator(
                                      strokeWidth:
                                          2,
                                      color: Colors
                                          .white))
                          : const Icon(
                              Icons.download),
                      label: Text(
                          _isExporting
                              ? '내보내는 중...'
                              : '데이터 내보내기'),
                      style: FilledButton
                          .styleFrom(
                        backgroundColor:
                            SanggamTheme
                                .primary,
                        foregroundColor:
                            SanggamTheme
                                .background,
                        minimumSize:
                            const Size
                                .fromHeight(
                                48),
                        shape:
                            RoundedRectangleBorder(
                          borderRadius:
                              BorderRadius
                                  .circular(
                                      16),
                        ),
                      ),
                    ),
                    const SizedBox(height: 16),

                    // 설명
                    const Text(
                      '* 내보내기된 파일은 기기의 다운로드 폴더에 저장됩니다.\n'
                      '* FHIR R4 형식은 전 세계 의료 기관에서 호환됩니다.\n'
                      '* 데이터는 암호화되어 전송됩니다.',
                      style: TextStyle(
                        color: SanggamTheme
                            .onSurfaceDim,
                        fontSize: 12,
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildFormatCard(String value,
      String title, String subtitle,
      IconData icon) {
    final isSelected = _format == value;
    return Padding(
      padding:
          const EdgeInsets.only(bottom: 8),
      child: SanggamContainer(
        borderRadius: 16,
        padding: EdgeInsets.zero,
        borderColor: isSelected
            ? SanggamTheme.primary
            : SanggamTheme.surfaceVariant,
        borderWidth: isSelected ? 2 : 1.5,
        child: RadioListTile<String>(
          title: Text(title,
              style: const TextStyle(
                color: Colors.white,
                fontWeight: FontWeight.w600,
              )),
          subtitle: Text(subtitle,
              style: const TextStyle(
                color:
                    SanggamTheme.onSurfaceDim,
                fontSize: 12,
              )),
          secondary: Icon(icon,
              color: isSelected
                  ? SanggamTheme.primary
                  : SanggamTheme
                      .onSurfaceDim),
          activeColor: SanggamTheme.primary,
          value: value,
          groupValue: _format,
          onChanged: (val) => setState(
              () => _format = val ?? 'fhir_r4'),
        ),
      ),
    );
  }

  Future<void> _handleExport() async {
    setState(() => _isExporting = true);
    try {
      final rest = ref.read(restClientProvider);
      final userId =
          ref.read(authProvider).userId ?? '';
      final categories = <String>[
        if (_includeMeasurements) 'measurements',
        if (_includeHealthRecords)
          'health_records',
        if (_includePrescriptions)
          'prescriptions',
        if (_includeCoaching) 'coaching',
      ];
      await rest.requestDataExport(
        userId: userId,
        format: _format,
        categories: categories,
      );
      if (mounted) {
        setState(
            () => _isExporting = false);
        ScaffoldMessenger.of(context)
            .showSnackBar(
          SnackBar(
            content: Text(
                '데이터가 ${_format.toUpperCase()} 형식으로 내보내기 되었습니다.'),
            action: SnackBarAction(
                label: '열기',
                onPressed: () {}),
          ),
        );
      }
    } on DioException {
      if (mounted) {
        setState(
            () => _isExporting = false);
        ScaffoldMessenger.of(context)
            .showSnackBar(
          const SnackBar(
            content: Text(
                '데이터 내보내기에 실패했습니다. 다시 시도해주세요.'),
          ),
        );
      }
    }
  }
}
