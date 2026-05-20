import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/generated/manpasik.pb.dart';
import 'package:manpasik/shared/providers/admin_settings_provider.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// AdminSettingsScreen — Sanggam Orbit 시스템 설정
//
// [Rule 4] +sanggam_theme.dart + sanggam_container.dart
// [Rule 4] AppBar+TabBar → body 내 커스텀 헤더+탭
// [Rule 4] Theme.of(context) 4x + ThemeData 파라미터 6x 제거
// [Rule 4] theme.textTheme.* ~20x → 직접 TextStyle
// [Rule 4] theme.colorScheme.* ~20x → SanggamTheme 상수
// [Rule 4] Colors.blue/orange/green/red/purple/grey → SanggamTheme
// [Rule 4] Card → SanggamContainer
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] borderRadius:14→16, 10→16, 6→8
// ───────────────────────────────────────────────────

/// 관리자 시스템 설정 화면 (AS-6)
class AdminSettingsScreen
    extends ConsumerStatefulWidget {
  const AdminSettingsScreen({super.key});

  @override
  ConsumerState<AdminSettingsScreen> createState() =>
      _AdminSettingsScreenState();
}

class _AdminSettingsScreenState
    extends ConsumerState<AdminSettingsScreen>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  final _searchController = TextEditingController();

  @override
  void initState() {
    super.initState();
    _tabController = TabController(
      length: adminConfigCategories.length,
      vsync: this,
    );
    _tabController.addListener(_onTabChanged);

    WidgetsBinding.instance
        .addPostFrameCallback((_) {
      ref
          .read(adminSettingsProvider.notifier)
          .loadConfigs();
    });
  }

  void _onTabChanged() {
    if (!_tabController.indexIsChanging) {
      final category = adminConfigCategories[
          _tabController.index];
      ref
          .read(adminSettingsProvider.notifier)
          .changeCategory(category);
    }
  }

  @override
  void dispose() {
    _tabController.removeListener(_onTabChanged);
    _tabController.dispose();
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final settingsState =
        ref.watch(adminSettingsProvider);

    ref.listen<AdminSettingsState>(
        adminSettingsProvider, (prev, next) {
      if (next.errorMessage != null &&
          next.errorMessage !=
              prev?.errorMessage) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(next.errorMessage!),
            backgroundColor: SanggamTheme.error,
            behavior: SnackBarBehavior.floating,
            action: SnackBarAction(
              label: '닫기',
              textColor: Colors.white,
              onPressed: () {
                ref
                    .read(adminSettingsProvider
                        .notifier)
                    .clearError();
              },
            ),
          ),
        );
      }
    });

    return Scaffold(
      backgroundColor: SanggamTheme.background,
      body: SafeArea(
        child: Column(
          children: [
            // 헤더
            Padding(
              padding: const EdgeInsets.symmetric(
                  horizontal: 8, vertical: 8),
              child: Row(
                children: [
                  IconButton(
                    icon: const Icon(
                        Icons.arrow_back,
                        color: Colors.white),
                    tooltip: '뒤로 가기',
                    onPressed: () => context.pop(),
                  ),
                  const SizedBox(width: 8),
                  const Expanded(
                    child: Text(
                      '시스템 설정 관리',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                ],
              ),
            ),

            // 탭바
            TabBar(
              controller: _tabController,
              isScrollable: true,
              tabAlignment: TabAlignment.start,
              indicatorColor: SanggamTheme.primary,
              labelColor: SanggamTheme.primary,
              unselectedLabelColor:
                  SanggamTheme.onSurfaceDim,
              dividerColor:
                  SanggamTheme.surfaceVariant,
              tabs: adminConfigCategories
                  .map((cat) {
                final label =
                    categoryLabels[cat] ?? cat;
                final count = settingsState
                    .categoryCounts[cat];
                return Tab(
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Icon(
                          _getCategoryIcon(cat),
                          size: 18),
                      const SizedBox(width: 8),
                      Text(label),
                      if (count != null &&
                          count > 0) ...[
                        const SizedBox(width: 8),
                        Container(
                          padding:
                              const EdgeInsets
                                  .symmetric(
                                  horizontal: 8,
                                  vertical: 2),
                          decoration:
                              BoxDecoration(
                            color: SanggamTheme
                                .primary
                                .withValues(
                                    alpha: 0.15),
                            borderRadius:
                                BorderRadius
                                    .circular(
                                        16),
                          ),
                          child: Text(
                            '$count',
                            style:
                                const TextStyle(
                              color: SanggamTheme
                                  .primary,
                              fontSize: 10,
                            ),
                          ),
                        ),
                      ],
                    ],
                  ),
                );
              }).toList(),
            ),

            // 검색 바
            Padding(
              padding:
                  const EdgeInsets.fromLTRB(
                      16, 16, 16, 8),
              child: TextField(
                controller: _searchController,
                style: const TextStyle(
                    color: Colors.white),
                decoration: InputDecoration(
                  hintText: '설정 검색...',
                  hintStyle: const TextStyle(
                      color: SanggamTheme
                          .onSurfaceDim),
                  prefixIcon: const Icon(
                      Icons.search,
                      color: SanggamTheme
                          .onSurfaceDim),
                  suffixIcon: _searchController
                          .text.isNotEmpty
                      ? IconButton(
                          icon: const Icon(
                              Icons.clear,
                              color: SanggamTheme
                                  .onSurfaceDim),
                          tooltip: '검색 초기화',
                          onPressed: () {
                            _searchController
                                .clear();
                            ref
                                .read(
                                    adminSettingsProvider
                                        .notifier)
                                .setSearchQuery(
                                    '');
                          },
                        )
                      : null,
                  filled: true,
                  fillColor: SanggamTheme.surface,
                  border: OutlineInputBorder(
                    borderRadius:
                        BorderRadius.circular(16),
                    borderSide: BorderSide.none,
                  ),
                  enabledBorder:
                      OutlineInputBorder(
                    borderRadius:
                        BorderRadius.circular(16),
                    borderSide: const BorderSide(
                        color: SanggamTheme
                            .surfaceVariant),
                  ),
                  focusedBorder:
                      OutlineInputBorder(
                    borderRadius:
                        BorderRadius.circular(16),
                    borderSide: const BorderSide(
                        color:
                            SanggamTheme.primary),
                  ),
                ),
                onChanged: (value) {
                  ref
                      .read(adminSettingsProvider
                          .notifier)
                      .setSearchQuery(value);
                  setState(() {});
                },
              ),
            ),

            // 설정 카드 목록
            Expanded(
              child: settingsState.isLoading
                  ? const Center(
                      child:
                          CircularProgressIndicator(
                              color: SanggamTheme
                                  .primary))
                  : settingsState
                          .filteredConfigs.isEmpty
                      ? _buildEmptyState()
                      : _buildConfigList(
                          settingsState
                              .filteredConfigs),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildEmptyState() {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(48),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              width: 100,
              height: 100,
              decoration: BoxDecoration(
                color: SanggamTheme.surfaceVariant,
                shape: BoxShape.circle,
              ),
              child: const Icon(
                Icons.settings_suggest_rounded,
                size: 48,
                color: SanggamTheme.onSurfaceDim,
              ),
            ),
            const SizedBox(height: 24),
            const Text(
              '설정 항목이 없습니다',
              style: TextStyle(
                color: Colors.white,
                fontSize: 16,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            const Text(
              '해당 카테고리에 등록된 설정이 없거나\n검색 결과가 없습니다',
              textAlign: TextAlign.center,
              style: TextStyle(
                color: SanggamTheme.onSurfaceDim,
                fontSize: 14,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildConfigList(
      List<ConfigWithMeta> configs) {
    return ListView.builder(
      padding: const EdgeInsets.fromLTRB(
          16, 0, 16, 16),
      itemCount: configs.length,
      itemBuilder: (context, index) {
        final config = configs[index];
        return _ConfigCard(
          config: config,
          onTap: () => _showEditDialog(config),
        );
      },
    );
  }

  Future<void> _showEditDialog(
      ConfigWithMeta config) async {
    final result = await showDialog<String>(
      context: context,
      builder: (ctx) =>
          _ConfigEditDialog(config: config),
    );

    if (result != null && mounted) {
      final success = await ref
          .read(adminSettingsProvider.notifier)
          .saveConfig(config.key, result);
      if (success && mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(
                '${config.displayName.isNotEmpty ? config.displayName : config.key} 저장 완료'),
            behavior: SnackBarBehavior.floating,
          ),
        );
      }
    }
  }

  IconData _getCategoryIcon(String category) {
    switch (category) {
      case 'general':
        return Icons.settings_rounded;
      case 'security':
        return Icons.shield_rounded;
      case 'ai':
        return Icons.psychology_rounded;
      case 'integration':
        return Icons.hub_rounded;
      case 'notification':
        return Icons.notifications_rounded;
      case 'measurement':
        return Icons.sensors_rounded;
      case 'payment':
        return Icons.payment_rounded;
      case 'ui':
        return Icons.palette_rounded;
      default:
        return Icons.tune_rounded;
    }
  }
}

// ── 설정 카드 위젯 ──

class _ConfigCard extends StatelessWidget {
  const _ConfigCard({
    required this.config,
    required this.onTap,
  });

  final ConfigWithMeta config;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final displayName =
        config.displayName.isNotEmpty
            ? config.displayName
            : config.key;
    final isSecret =
        config.securityLevel == 'secret' ||
            config.valueType == 'secret';
    final displayValue =
        isSecret ? '••••••••' : config.value;

    return SanggamContainer(
      borderRadius: 16,
      padding: EdgeInsets.zero,
      margin: const EdgeInsets.only(bottom: 16),
      onTap: onTap,
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment:
              CrossAxisAlignment.start,
          children: [
            // 헤더: 키 이름 + 타입 배지
            Row(
              children: [
                Expanded(
                  child: Text(
                    displayName,
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 14,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
                _ValueTypeBadge(
                    valueType: config.valueType),
                if (config.isRequired) ...[
                  const SizedBox(width: 8),
                  Container(
                    padding:
                        const EdgeInsets.symmetric(
                            horizontal: 8,
                            vertical: 2),
                    decoration: BoxDecoration(
                      color: SanggamTheme.error
                          .withValues(alpha: 0.15),
                      borderRadius:
                          BorderRadius.circular(8),
                    ),
                    child: const Text(
                      '필수',
                      style: TextStyle(
                        color: SanggamTheme.error,
                        fontSize: 10,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                ],
              ],
            ),

            const SizedBox(height: 8),

            // 현재 값
            Container(
              width: double.infinity,
              padding: const EdgeInsets.symmetric(
                  horizontal: 16, vertical: 8),
              decoration: BoxDecoration(
                color: SanggamTheme.surfaceVariant
                    .withValues(alpha: 0.5),
                borderRadius:
                    BorderRadius.circular(8),
              ),
              child: Text(
                displayValue.isNotEmpty
                    ? displayValue
                    : '(미설정)',
                style: TextStyle(
                  fontFamily: 'monospace',
                  fontSize: 14,
                  color: displayValue.isNotEmpty
                      ? Colors.white
                      : SanggamTheme.onSurfaceDim,
                ),
                maxLines: 2,
                overflow: TextOverflow.ellipsis,
              ),
            ),

            // 설명
            if (config
                .description.isNotEmpty) ...[
              const SizedBox(height: 8),
              Text(
                config.description,
                style: const TextStyle(
                  color:
                      SanggamTheme.onSurfaceDim,
                  fontSize: 12,
                ),
                maxLines: 2,
                overflow: TextOverflow.ellipsis,
              ),
            ],

            // 메타 정보
            const SizedBox(height: 8),
            Row(
              children: [
                Icon(
                  Icons.key_rounded,
                  size: 14,
                  color: SanggamTheme.onSurfaceDim
                      .withValues(alpha: 0.6),
                ),
                const SizedBox(width: 4),
                Expanded(
                  child: Text(
                    config.key,
                    style: TextStyle(
                      color: SanggamTheme
                          .onSurfaceDim
                          .withValues(alpha: 0.6),
                      fontSize: 10,
                      fontFamily: 'monospace',
                    ),
                    overflow:
                        TextOverflow.ellipsis,
                  ),
                ),
                if (config
                    .restartRequired) ...[
                  const SizedBox(width: 8),
                  Icon(
                    Icons.restart_alt_rounded,
                    size: 14,
                    color: SanggamTheme.error
                        .withValues(alpha: 0.7),
                  ),
                  const SizedBox(width: 2),
                  Text(
                    '재시작 필요',
                    style: TextStyle(
                      color: SanggamTheme.error
                          .withValues(alpha: 0.7),
                      fontSize: 10,
                    ),
                  ),
                ],
              ],
            ),
          ],
        ),
      ),
    );
  }
}

// ── 값 타입 배지 ──

class _ValueTypeBadge extends StatelessWidget {
  const _ValueTypeBadge(
      {required this.valueType});

  final String valueType;

  @override
  Widget build(BuildContext context) {
    final (label, color) = _typeInfo;
    return Container(
      padding: const EdgeInsets.symmetric(
          horizontal: 8, vertical: 2),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.12),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Text(
        label,
        style: TextStyle(
          color: color,
          fontSize: 10,
          fontWeight: FontWeight.w600,
        ),
      ),
    );
  }

  (String, Color) get _typeInfo {
    switch (valueType) {
      case 'string':
        return ('문자열', SanggamTheme.jagaeCyan);
      case 'number':
        return ('숫자', SanggamTheme.primary);
      case 'boolean':
        return ('불리언', SanggamTheme.jagaeCyan);
      case 'secret':
        return ('비밀', SanggamTheme.error);
      case 'select':
        return ('선택', SanggamTheme.jagaeMagenta);
      default:
        return (
          valueType.isNotEmpty
              ? valueType
              : '기타',
          SanggamTheme.onSurfaceDim
        );
    }
  }
}

// ── 편집 다이얼로그 ──

class _ConfigEditDialog extends StatefulWidget {
  const _ConfigEditDialog(
      {required this.config});

  final ConfigWithMeta config;

  @override
  State<_ConfigEditDialog> createState() =>
      _ConfigEditDialogState();
}

class _ConfigEditDialogState
    extends State<_ConfigEditDialog> {
  late TextEditingController _valueController;
  late bool _boolValue;
  late String _selectValue;
  String? _validationError;

  ConfigWithMeta get config => widget.config;

  @override
  void initState() {
    super.initState();
    _valueController =
        TextEditingController(text: config.value);
    _boolValue =
        config.value.toLowerCase() == 'true';
    _selectValue = config.value;
  }

  @override
  void dispose() {
    _valueController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final displayName =
        config.displayName.isNotEmpty
            ? config.displayName
            : config.key;

    return AlertDialog(
      backgroundColor: SanggamTheme.surface,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(16),
      ),
      title: Text(displayName,
          style: const TextStyle(
              color: Colors.white)),
      content: SizedBox(
        width:
            MediaQuery.of(context).size.width *
                0.8,
        child: SingleChildScrollView(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment:
                CrossAxisAlignment.start,
            children: [
              if (config
                  .description.isNotEmpty) ...[
                Text(
                  config.description,
                  style: const TextStyle(
                    color:
                        SanggamTheme.onSurfaceDim,
                    fontSize: 12,
                  ),
                ),
                const SizedBox(height: 16),
              ],

              if (config
                  .helpText.isNotEmpty) ...[
                Container(
                  padding:
                      const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: SanggamTheme.primary
                        .withValues(alpha: 0.1),
                    borderRadius:
                        BorderRadius.circular(8),
                  ),
                  child: Row(
                    crossAxisAlignment:
                        CrossAxisAlignment.start,
                    children: [
                      const Icon(
                        Icons
                            .info_outline_rounded,
                        size: 18,
                        color:
                            SanggamTheme.primary,
                      ),
                      const SizedBox(width: 8),
                      Expanded(
                        child: Text(
                          config.helpText,
                          style: const TextStyle(
                            color: SanggamTheme
                                .primary,
                            fontSize: 12,
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 16),
              ],

              _buildInputWidget(),

              if (_validationError !=
                  null) ...[
                const SizedBox(height: 8),
                Text(
                  _validationError!,
                  style: const TextStyle(
                    color: SanggamTheme.error,
                    fontSize: 12,
                  ),
                ),
              ],

              const SizedBox(height: 16),
              _buildMetaInfo(),
            ],
          ),
        ),
      ),
      actions: [
        TextButton(
          onPressed: () =>
              Navigator.of(context).pop(),
          child: const Text('취소',
              style: TextStyle(
                  color: SanggamTheme
                      .onSurfaceDim)),
        ),
        FilledButton(
          onPressed: _onSave,
          style: FilledButton.styleFrom(
            backgroundColor: SanggamTheme.primary,
            foregroundColor:
                SanggamTheme.background,
            shape: RoundedRectangleBorder(
              borderRadius:
                  BorderRadius.circular(16),
            ),
          ),
          child: const Text('저장'),
        ),
      ],
    );
  }

  Widget _buildInputWidget() {
    switch (config.valueType) {
      case 'boolean':
        return SwitchListTile(
          title: Text(
            _boolValue ? '활성화' : '비활성화',
            style: const TextStyle(
              color: Colors.white,
              fontSize: 16,
            ),
          ),
          value: _boolValue,
          activeColor: SanggamTheme.primary,
          onChanged: (v) =>
              setState(() => _boolValue = v),
          contentPadding: EdgeInsets.zero,
        );

      case 'select':
        final allowed = config.allowedValues;
        if (allowed.isEmpty) {
          return _buildTextField();
        }
        return DropdownButtonFormField<String>(
          value:
              allowed.contains(_selectValue)
                  ? _selectValue
                  : null,
          dropdownColor: SanggamTheme.surface,
          style: const TextStyle(
              color: Colors.white),
          decoration: InputDecoration(
            labelText: '값 선택',
            labelStyle: const TextStyle(
                color:
                    SanggamTheme.onSurfaceDim),
            hintText:
                config.placeholder.isNotEmpty
                    ? config.placeholder
                    : '선택하세요',
            hintStyle: const TextStyle(
                color:
                    SanggamTheme.onSurfaceDim),
          ),
          items: allowed.map((v) {
            return DropdownMenuItem(
                value: v, child: Text(v));
          }).toList(),
          onChanged: (v) {
            if (v != null) {
              setState(() {
                _selectValue = v;
                _validationError = null;
              });
            }
          },
        );

      case 'number':
        return TextField(
          controller: _valueController,
          style: const TextStyle(
              color: Colors.white),
          keyboardType:
              const TextInputType
                  .numberWithOptions(
                  decimal: true),
          inputFormatters: [
            FilteringTextInputFormatter.allow(
                RegExp(r'[\d.\-]')),
          ],
          decoration: InputDecoration(
            labelText: '값',
            labelStyle: const TextStyle(
                color:
                    SanggamTheme.onSurfaceDim),
            hintText:
                config.placeholder.isNotEmpty
                    ? config.placeholder
                    : '숫자를 입력하세요',
            hintStyle: const TextStyle(
                color:
                    SanggamTheme.onSurfaceDim),
            helperText: _buildRangeHelper(),
            helperStyle: const TextStyle(
                color:
                    SanggamTheme.onSurfaceDim),
          ),
          onChanged: (_) => setState(
              () => _validationError = null),
        );

      case 'secret':
        return TextField(
          controller: _valueController,
          style: const TextStyle(
              color: Colors.white),
          obscureText: true,
          decoration: InputDecoration(
            labelText: '값',
            labelStyle: const TextStyle(
                color:
                    SanggamTheme.onSurfaceDim),
            hintText:
                config.placeholder.isNotEmpty
                    ? config.placeholder
                    : '비밀 값을 입력하세요',
            hintStyle: const TextStyle(
                color:
                    SanggamTheme.onSurfaceDim),
          ),
          onChanged: (_) => setState(
              () => _validationError = null),
        );

      default:
        return _buildTextField();
    }
  }

  Widget _buildTextField() {
    return TextField(
      controller: _valueController,
      style:
          const TextStyle(color: Colors.white),
      maxLines:
          config.value.contains('\n') ? 5 : 1,
      decoration: InputDecoration(
        labelText: '값',
        labelStyle: const TextStyle(
            color: SanggamTheme.onSurfaceDim),
        hintText: config.placeholder.isNotEmpty
            ? config.placeholder
            : '값을 입력하세요',
        hintStyle: const TextStyle(
            color: SanggamTheme.onSurfaceDim),
      ),
      onChanged: (_) => setState(
          () => _validationError = null),
    );
  }

  String? _buildRangeHelper() {
    final min = config.validationMin;
    final max = config.validationMax;
    if (min != 0 || max != 0) {
      return '범위: $min ~ $max';
    }
    return null;
  }

  Widget _buildMetaInfo() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: SanggamTheme.surfaceVariant
            .withValues(alpha: 0.5),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Column(
        crossAxisAlignment:
            CrossAxisAlignment.start,
        children: [
          _metaRow('키', config.key),
          if (config.category.isNotEmpty)
            _metaRow(
                '카테고리',
                categoryLabels[
                        config.category] ??
                    config.category),
          if (config.defaultValue.isNotEmpty)
            _metaRow('기본값', config.defaultValue),
          if (config.serviceName.isNotEmpty)
            _metaRow('서비스', config.serviceName),
          if (config.updatedBy.isNotEmpty)
            _metaRow('수정자', config.updatedBy),
          if (config.restartRequired)
            _metaRow('재시작', '변경 적용 시 재시작 필요'),
        ],
      ),
    );
  }

  Widget _metaRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Row(
        crossAxisAlignment:
            CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 64,
            child: Text(
              label,
              style: const TextStyle(
                color: SanggamTheme.onSurfaceDim,
                fontSize: 10,
                fontWeight: FontWeight.bold,
              ),
            ),
          ),
          Expanded(
            child: Text(
              value,
              style: const TextStyle(
                color: Colors.white,
                fontSize: 10,
                fontFamily: 'monospace',
              ),
            ),
          ),
        ],
      ),
    );
  }

  void _onSave() {
    String finalValue;
    switch (config.valueType) {
      case 'boolean':
        finalValue = _boolValue.toString();
        break;
      case 'select':
        finalValue = _selectValue;
        break;
      default:
        finalValue =
            _valueController.text.trim();
    }

    if (config.isRequired &&
        finalValue.isEmpty) {
      setState(() => _validationError =
          '필수 항목입니다. 값을 입력하세요.');
      return;
    }

    if (config.valueType == 'number' &&
        finalValue.isNotEmpty) {
      final parsed =
          double.tryParse(finalValue);
      if (parsed == null) {
        setState(() =>
            _validationError = '유효한 숫자를 입력하세요.');
        return;
      }
      final min = config.validationMin;
      final max = config.validationMax;
      if ((min != 0 || max != 0) &&
          (parsed < min || parsed > max)) {
        setState(() => _validationError =
            '범위를 벗어났습니다 ($min ~ $max)');
        return;
      }
    }

    Navigator.of(context).pop(finalValue);
  }
}
