import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:grpc/grpc.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/services/auth_interceptor.dart';
import 'package:manpasik/generated/manpasik.pb.dart';
import 'package:manpasik/generated/manpasik.pbgrpc.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';

// ── 카테고리 상수 ──

const List<String> adminConfigCategories = [
  'general',
  'security',
  'ai',
  'integration',
  'notification',
  'measurement',
  'payment',
  'ui',
];

const Map<String, String> categoryLabels = {
  'general': '일반',
  'security': '보안',
  'ai': 'AI',
  'integration': '연동',
  'notification': '알림',
  'measurement': '측정',
  'payment': '결제',
  'ui': 'UI',
};

const Map<String, String> categoryIcons = {
  'general': 'settings',
  'security': 'shield',
  'ai': 'psychology',
  'integration': 'hub',
  'notification': 'notifications',
  'measurement': 'sensors',
  'payment': 'payment',
  'ui': 'palette',
};

// ── 상태 모델 ──

class AdminSettingsState {
  final List<ConfigWithMeta> configs;
  final String selectedCategory;
  final String searchQuery;
  final bool isLoading;
  final String? errorMessage;
  final Map<String, int> categoryCounts;

  const AdminSettingsState({
    this.configs = const [],
    this.selectedCategory = 'general',
    this.searchQuery = '',
    this.isLoading = false,
    this.errorMessage,
    this.categoryCounts = const {},
  });

  AdminSettingsState copyWith({
    List<ConfigWithMeta>? configs,
    String? selectedCategory,
    String? searchQuery,
    bool? isLoading,
    String? errorMessage,
    Map<String, int>? categoryCounts,
  }) {
    return AdminSettingsState(
      configs: configs ?? this.configs,
      selectedCategory: selectedCategory ?? this.selectedCategory,
      searchQuery: searchQuery ?? this.searchQuery,
      isLoading: isLoading ?? this.isLoading,
      errorMessage: errorMessage,
      categoryCounts: categoryCounts ?? this.categoryCounts,
    );
  }

  /// 현재 검색어로 필터링된 설정 목록
  List<ConfigWithMeta> get filteredConfigs {
    if (searchQuery.isEmpty) return configs;
    final query = searchQuery.toLowerCase();
    return configs.where((c) {
      return c.key.toLowerCase().contains(query) ||
          c.displayName.toLowerCase().contains(query) ||
          c.description.toLowerCase().contains(query);
    }).toList();
  }
}

// ── Notifier ──

class AdminSettingsNotifier extends StateNotifier<AdminSettingsState> {
  AdminSettingsNotifier(this._client) : super(const AdminSettingsState());

  final AdminServiceClient _client;

  /// 설정 목록 로드
  Future<void> loadConfigs({String? category}) async {
    final cat = category ?? state.selectedCategory;
    state = state.copyWith(isLoading: true, errorMessage: null, selectedCategory: cat);
    try {
      final response = await _client.listSystemConfigs(
        ListSystemConfigsRequest(
          languageCode: 'ko',
          category: cat,
        ),
      );
      state = state.copyWith(
        configs: response.configs,
        categoryCounts: response.categoryCounts,
        isLoading: false,
      );
    } on GrpcError catch (e) {
      state = state.copyWith(
        isLoading: false,
        errorMessage: 'gRPC 오류: ${e.message ?? e.codeName}',
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        errorMessage: '설정을 불러올 수 없습니다: $e',
      );
    }
  }

  /// 카테고리 변경
  Future<void> changeCategory(String category) async {
    if (category == state.selectedCategory) return;
    await loadConfigs(category: category);
  }

  /// 검색어 변경
  void setSearchQuery(String query) {
    state = state.copyWith(searchQuery: query);
  }

  /// 설정 값 유효성 검증
  Future<ValidateConfigValueResponse> validateValue(String key, String value) async {
    return _client.validateConfigValue(
      ValidateConfigValueRequest(key: key, value: value),
    );
  }

  /// 설정 값 저장
  Future<bool> saveConfig(String key, String value, {String? description}) async {
    try {
      // 먼저 유효성 검증
      final validation = await validateValue(key, value);
      if (!validation.valid) {
        state = state.copyWith(
          errorMessage: validation.errorMessage.isNotEmpty
              ? validation.errorMessage
              : '유효하지 않은 값입니다',
        );
        return false;
      }

      // 저장
      await _client.setSystemConfig(
        SetSystemConfigRequest(
          key: key,
          value: value,
          description: description,
        ),
      );

      // 설정 목록 다시 로드
      await loadConfigs();
      return true;
    } on GrpcError catch (e) {
      state = state.copyWith(
        errorMessage: '저장 실패: ${e.message ?? e.codeName}',
      );
      return false;
    } catch (e) {
      state = state.copyWith(errorMessage: '저장 실패: $e');
      return false;
    }
  }

  /// 에러 메시지 제거
  void clearError() {
    state = state.copyWith(errorMessage: null);
  }
}

// ── Provider ──

final adminSettingsProvider =
    StateNotifierProvider<AdminSettingsNotifier, AdminSettingsState>((ref) {
  final manager = ref.watch(grpcClientManagerProvider);
  final accessToken = ref.watch(authProvider).accessToken;

  final interceptors = <ClientInterceptor>[];
  if (accessToken != null && accessToken.isNotEmpty) {
    interceptors.add(AuthInterceptor(() => accessToken));
  }

  final client = AdminServiceClient(
    manager.adminChannel,
    interceptors: interceptors,
  );
  return AdminSettingsNotifier(client);
});
