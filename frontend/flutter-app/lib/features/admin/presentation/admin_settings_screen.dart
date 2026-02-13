import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/generated/manpasik.pb.dart';
import 'package:manpasik/shared/providers/admin_settings_provider.dart';

/// 관리자 시스템 설정 화면 (AS-6)
///
/// 카테고리별 설정 조회, 검색, 편집 기능 제공.
/// AdminService gRPC (ListSystemConfigs, SetSystemConfig, ValidateConfigValue) 연동.
class AdminSettingsScreen extends ConsumerStatefulWidget {
  const AdminSettingsScreen({super.key});

  @override
  ConsumerState<AdminSettingsScreen> createState() =>
      _AdminSettingsScreenState();
}

class _AdminSettingsScreenState extends ConsumerState<AdminSettingsScreen>
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

    // 초기 데이터 로드
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(adminSettingsProvider.notifier).loadConfigs();
    });
  }

  void _onTabChanged() {
    if (!_tabController.indexIsChanging) {
      final category = adminConfigCategories[_tabController.index];
      ref.read(adminSettingsProvider.notifier).changeCategory(category);
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
    final theme = Theme.of(context);
    final settingsState = ref.watch(adminSettingsProvider);

    // 에러 스낵바
    ref.listen<AdminSettingsState>(adminSettingsProvider, (prev, next) {
      if (next.errorMessage != null && next.errorMessage != prev?.errorMessage) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(next.errorMessage!),
            backgroundColor: theme.colorScheme.error,
            behavior: SnackBarBehavior.floating,
            action: SnackBarAction(
              label: '닫기',
              textColor: Colors.white,
              onPressed: () {
                ref.read(adminSettingsProvider.notifier).clearError();
              },
            ),
          ),
        );
      }
    });

    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
        title: const Text('시스템 설정 관리'),
        bottom: TabBar(
          controller: _tabController,
          isScrollable: true,
          tabAlignment: TabAlignment.start,
          tabs: adminConfigCategories.map((cat) {
            final label = categoryLabels[cat] ?? cat;
            final count = settingsState.categoryCounts[cat];
            return Tab(
              child: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Icon(_getCategoryIcon(cat), size: 18),
                  const SizedBox(width: 6),
                  Text(label),
                  if (count != null && count > 0) ...[
                    const SizedBox(width: 4),
                    Container(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 6,
                        vertical: 1,
                      ),
                      decoration: BoxDecoration(
                        color: theme.colorScheme.primaryContainer,
                        borderRadius: BorderRadius.circular(10),
                      ),
                      child: Text(
                        '$count',
                        style: theme.textTheme.labelSmall?.copyWith(
                          color: theme.colorScheme.onPrimaryContainer,
                        ),
                      ),
                    ),
                  ],
                ],
              ),
            );
          }).toList(),
        ),
      ),
      body: Column(
        children: [
          // 검색 바
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
            child: TextField(
              controller: _searchController,
              decoration: InputDecoration(
                hintText: '설정 검색...',
                prefixIcon: const Icon(Icons.search),
                suffixIcon: _searchController.text.isNotEmpty
                    ? IconButton(
                        icon: const Icon(Icons.clear),
                        onPressed: () {
                          _searchController.clear();
                          ref
                              .read(adminSettingsProvider.notifier)
                              .setSearchQuery('');
                        },
                      )
                    : null,
              ),
              onChanged: (value) {
                ref
                    .read(adminSettingsProvider.notifier)
                    .setSearchQuery(value);
                setState(() {}); // suffixIcon 갱신
              },
            ),
          ),

          // 설정 카드 목록
          Expanded(
            child: settingsState.isLoading
                ? const Center(child: CircularProgressIndicator())
                : settingsState.filteredConfigs.isEmpty
                    ? _buildEmptyState(theme)
                    : _buildConfigList(theme, settingsState.filteredConfigs),
          ),
        ],
      ),
    );
  }

  /// 빈 상태 UI
  Widget _buildEmptyState(ThemeData theme) {
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
                color: theme.colorScheme.surfaceContainerHighest,
                shape: BoxShape.circle,
              ),
              child: Icon(
                Icons.settings_suggest_rounded,
                size: 48,
                color: theme.colorScheme.onSurfaceVariant,
              ),
            ),
            const SizedBox(height: 20),
            Text(
              '설정 항목이 없습니다',
              style: theme.textTheme.titleMedium?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              '해당 카테고리에 등록된 설정이 없거나\n검색 결과가 없습니다',
              textAlign: TextAlign.center,
              style: theme.textTheme.bodyMedium?.copyWith(
                color: theme.colorScheme.onSurfaceVariant,
              ),
            ),
          ],
        ),
      ),
    );
  }

  /// 설정 카드 목록
  Widget _buildConfigList(ThemeData theme, List<ConfigWithMeta> configs) {
    return ListView.builder(
      padding: const EdgeInsets.fromLTRB(16, 0, 16, 16),
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

  /// 편집 다이얼로그
  Future<void> _showEditDialog(ConfigWithMeta config) async {
    final result = await showDialog<String>(
      context: context,
      builder: (ctx) => _ConfigEditDialog(config: config),
    );

    if (result != null && mounted) {
      final success = await ref
          .read(adminSettingsProvider.notifier)
          .saveConfig(config.key, result);
      if (success && mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('${config.displayName.isNotEmpty ? config.displayName : config.key} 저장 완료'),
            behavior: SnackBarBehavior.floating,
          ),
        );
      }
    }
  }

  /// 카테고리 아이콘
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
    final theme = Theme.of(context);
    final displayName = config.displayName.isNotEmpty
        ? config.displayName
        : config.key;
    final isSecret = config.securityLevel == 'secret' ||
        config.valueType == 'secret';
    final displayValue = isSecret ? '••••••••' : config.value;

    return Card(
      margin: const EdgeInsets.only(bottom: 10),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(14),
      ),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(14),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // 헤더: 키 이름 + 타입 배지
              Row(
                children: [
                  Expanded(
                    child: Text(
                      displayName,
                      style: theme.textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
                  _ValueTypeBadge(valueType: config.valueType),
                  if (config.isRequired) ...[
                    const SizedBox(width: 6),
                    Container(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 6,
                        vertical: 2,
                      ),
                      decoration: BoxDecoration(
                        color: theme.colorScheme.errorContainer,
                        borderRadius: BorderRadius.circular(6),
                      ),
                      child: Text(
                        '필수',
                        style: theme.textTheme.labelSmall?.copyWith(
                          color: theme.colorScheme.onErrorContainer,
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
                  horizontal: 12,
                  vertical: 8,
                ),
                decoration: BoxDecoration(
                  color: theme.colorScheme.surfaceContainerHighest
                      .withValues(alpha: 0.5),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Text(
                  displayValue.isNotEmpty ? displayValue : '(미설정)',
                  style: theme.textTheme.bodyMedium?.copyWith(
                    fontFamily: 'monospace',
                    color: displayValue.isNotEmpty
                        ? theme.colorScheme.onSurface
                        : theme.colorScheme.onSurfaceVariant,
                  ),
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),
              ),

              // 설명
              if (config.description.isNotEmpty) ...[
                const SizedBox(height: 8),
                Text(
                  config.description,
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: theme.colorScheme.onSurfaceVariant,
                  ),
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),
              ],

              // 메타 정보 (키, 서비스명)
              const SizedBox(height: 8),
              Row(
                children: [
                  Icon(
                    Icons.key_rounded,
                    size: 14,
                    color: theme.colorScheme.onSurfaceVariant
                        .withValues(alpha: 0.6),
                  ),
                  const SizedBox(width: 4),
                  Expanded(
                    child: Text(
                      config.key,
                      style: theme.textTheme.labelSmall?.copyWith(
                        color: theme.colorScheme.onSurfaceVariant
                            .withValues(alpha: 0.6),
                        fontFamily: 'monospace',
                      ),
                      overflow: TextOverflow.ellipsis,
                    ),
                  ),
                  if (config.restartRequired) ...[
                    const SizedBox(width: 8),
                    Icon(
                      Icons.restart_alt_rounded,
                      size: 14,
                      color: theme.colorScheme.error.withValues(alpha: 0.7),
                    ),
                    const SizedBox(width: 2),
                    Text(
                      '재시작 필요',
                      style: theme.textTheme.labelSmall?.copyWith(
                        color: theme.colorScheme.error.withValues(alpha: 0.7),
                      ),
                    ),
                  ],
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}

// ── 값 타입 배지 ──

class _ValueTypeBadge extends StatelessWidget {
  const _ValueTypeBadge({required this.valueType});

  final String valueType;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final (label, color) = _typeInfo;
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.12),
        borderRadius: BorderRadius.circular(6),
      ),
      child: Text(
        label,
        style: theme.textTheme.labelSmall?.copyWith(
          color: color,
          fontWeight: FontWeight.w600,
        ),
      ),
    );
  }

  (String, Color) get _typeInfo {
    switch (valueType) {
      case 'string':
        return ('문자열', Colors.blue);
      case 'number':
        return ('숫자', Colors.orange);
      case 'boolean':
        return ('불리언', Colors.green);
      case 'secret':
        return ('비밀', Colors.red);
      case 'select':
        return ('선택', Colors.purple);
      default:
        return (valueType.isNotEmpty ? valueType : '기타', Colors.grey);
    }
  }
}

// ── 편집 다이얼로그 ──

class _ConfigEditDialog extends StatefulWidget {
  const _ConfigEditDialog({required this.config});

  final ConfigWithMeta config;

  @override
  State<_ConfigEditDialog> createState() => _ConfigEditDialogState();
}

class _ConfigEditDialogState extends State<_ConfigEditDialog> {
  late TextEditingController _valueController;
  late bool _boolValue;
  late String _selectValue;
  String? _validationError;

  ConfigWithMeta get config => widget.config;

  @override
  void initState() {
    super.initState();
    _valueController = TextEditingController(text: config.value);
    _boolValue = config.value.toLowerCase() == 'true';
    _selectValue = config.value;
  }

  @override
  void dispose() {
    _valueController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final displayName = config.displayName.isNotEmpty
        ? config.displayName
        : config.key;

    return AlertDialog(
      title: Text(displayName),
      content: SizedBox(
        width: MediaQuery.of(context).size.width * 0.8,
        child: SingleChildScrollView(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // 설명
              if (config.description.isNotEmpty) ...[
                Text(
                  config.description,
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: theme.colorScheme.onSurfaceVariant,
                  ),
                ),
                const SizedBox(height: 16),
              ],

              // 도움말
              if (config.helpText.isNotEmpty) ...[
                Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: theme.colorScheme.primaryContainer
                        .withValues(alpha: 0.3),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Icon(
                        Icons.info_outline_rounded,
                        size: 18,
                        color: theme.colorScheme.primary,
                      ),
                      const SizedBox(width: 8),
                      Expanded(
                        child: Text(
                          config.helpText,
                          style: theme.textTheme.bodySmall?.copyWith(
                            color: theme.colorScheme.primary,
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 16),
              ],

              // 편집 입력
              _buildInputWidget(theme),

              // 유효성 오류
              if (_validationError != null) ...[
                const SizedBox(height: 8),
                Text(
                  _validationError!,
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: theme.colorScheme.error,
                  ),
                ),
              ],

              // 메타 정보
              const SizedBox(height: 16),
              _buildMetaInfo(theme),
            ],
          ),
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          child: const Text('취소'),
        ),
        FilledButton(
          onPressed: _onSave,
          child: const Text('저장'),
        ),
      ],
    );
  }

  /// 값 타입별 입력 위젯
  Widget _buildInputWidget(ThemeData theme) {
    switch (config.valueType) {
      case 'boolean':
        return SwitchListTile(
          title: Text(
            _boolValue ? '활성화' : '비활성화',
            style: theme.textTheme.bodyLarge,
          ),
          value: _boolValue,
          onChanged: (v) => setState(() => _boolValue = v),
          contentPadding: EdgeInsets.zero,
        );

      case 'select':
        final allowed = config.allowedValues;
        if (allowed.isEmpty) {
          return _buildTextField(theme);
        }
        return DropdownButtonFormField<String>(
          value: allowed.contains(_selectValue) ? _selectValue : null,
          decoration: InputDecoration(
            labelText: '값 선택',
            hintText: config.placeholder.isNotEmpty
                ? config.placeholder
                : '선택하세요',
          ),
          items: allowed.map((v) {
            return DropdownMenuItem(value: v, child: Text(v));
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
          keyboardType: const TextInputType.numberWithOptions(decimal: true),
          inputFormatters: [
            FilteringTextInputFormatter.allow(RegExp(r'[\d.\-]')),
          ],
          decoration: InputDecoration(
            labelText: '값',
            hintText: config.placeholder.isNotEmpty
                ? config.placeholder
                : '숫자를 입력하세요',
            helperText: _buildRangeHelper(),
          ),
          onChanged: (_) => setState(() => _validationError = null),
        );

      case 'secret':
        return TextField(
          controller: _valueController,
          obscureText: true,
          decoration: InputDecoration(
            labelText: '값',
            hintText: config.placeholder.isNotEmpty
                ? config.placeholder
                : '비밀 값을 입력하세요',
          ),
          onChanged: (_) => setState(() => _validationError = null),
        );

      default: // string 등
        return _buildTextField(theme);
    }
  }

  /// 기본 텍스트 입력
  Widget _buildTextField(ThemeData theme) {
    return TextField(
      controller: _valueController,
      maxLines: config.value.contains('\n') ? 5 : 1,
      decoration: InputDecoration(
        labelText: '값',
        hintText: config.placeholder.isNotEmpty
            ? config.placeholder
            : '값을 입력하세요',
      ),
      onChanged: (_) => setState(() => _validationError = null),
    );
  }

  /// 범위 도움말 (number 타입)
  String? _buildRangeHelper() {
    final min = config.validationMin;
    final max = config.validationMax;
    if (min != 0 || max != 0) {
      return '범위: $min ~ $max';
    }
    return null;
  }

  /// 메타 정보 표시
  Widget _buildMetaInfo(ThemeData theme) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: theme.colorScheme.surfaceContainerHighest.withValues(alpha: 0.5),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _metaRow(theme, '키', config.key),
          if (config.category.isNotEmpty)
            _metaRow(theme, '카테고리',
                categoryLabels[config.category] ?? config.category),
          if (config.defaultValue.isNotEmpty)
            _metaRow(theme, '기본값', config.defaultValue),
          if (config.serviceName.isNotEmpty)
            _metaRow(theme, '서비스', config.serviceName),
          if (config.updatedBy.isNotEmpty)
            _metaRow(theme, '수정자', config.updatedBy),
          if (config.restartRequired)
            _metaRow(theme, '재시작', '변경 적용 시 재시작 필요'),
        ],
      ),
    );
  }

  Widget _metaRow(ThemeData theme, String label, String value) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 4),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 64,
            child: Text(
              label,
              style: theme.textTheme.labelSmall?.copyWith(
                color: theme.colorScheme.onSurfaceVariant,
                fontWeight: FontWeight.bold,
              ),
            ),
          ),
          Expanded(
            child: Text(
              value,
              style: theme.textTheme.labelSmall?.copyWith(
                color: theme.colorScheme.onSurface,
                fontFamily: 'monospace',
              ),
            ),
          ),
        ],
      ),
    );
  }

  /// 저장
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
        finalValue = _valueController.text.trim();
    }

    // 클라이언트 측 기본 유효성 검증
    if (config.isRequired && finalValue.isEmpty) {
      setState(() => _validationError = '필수 항목입니다. 값을 입력하세요.');
      return;
    }

    if (config.valueType == 'number' && finalValue.isNotEmpty) {
      final parsed = double.tryParse(finalValue);
      if (parsed == null) {
        setState(() => _validationError = '유효한 숫자를 입력하세요.');
        return;
      }
      final min = config.validationMin;
      final max = config.validationMax;
      if ((min != 0 || max != 0) && (parsed < min || parsed > max)) {
        setState(
            () => _validationError = '범위를 벗어났습니다 ($min ~ $max)');
        return;
      }
    }

    Navigator.of(context).pop(finalValue);
  }
}
