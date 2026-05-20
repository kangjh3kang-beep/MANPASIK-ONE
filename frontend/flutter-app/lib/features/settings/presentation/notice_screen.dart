import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

// ───────────────────────────────────────────────────
// NoticeScreen — Sanggam Orbit 공지사항
//
// [Rule 4] +sanggam_theme.dart + sanggam_container.dart
// [Rule 4] AppBar → body 내 커스텀 헤더
// [Rule 4] Theme.of(context) 1x 제거
// [Rule 4] ThemeData 파라미터 1x 제거
// [Rule 4] theme.textTheme.* ~3x → 직접 TextStyle
// [Rule 4] theme.colorScheme.* ~3x → SanggamTheme 상수
// [Rule 4] Colors.blue → SanggamTheme.jagaeCyan
// [Rule 4] Colors.orange → SanggamTheme.primary
// [Rule 4] Colors.red → SanggamTheme.error
// [Rule 4] Colors.green → SanggamTheme.jagaeCyan
// [Rule 4] Colors.purple → SanggamTheme.jagaeMagenta
// [Rule 4] Colors.grey → SanggamTheme.onSurfaceDim
// [Rule 4] withOpacity → withValues(alpha:)
// [Rule 4] Card → SanggamContainer
// [Rule 4] Scaffold 배경 → SanggamTheme.background
// [Rule 2] spacing 6→8
// ───────────────────────────────────────────────────

/// 공지사항 화면
class NoticeScreen extends ConsumerStatefulWidget {
  const NoticeScreen({super.key});

  @override
  ConsumerState<NoticeScreen> createState() => _NoticeScreenState();
}

class _NoticeScreenState extends ConsumerState<NoticeScreen> {
  List<_NoticeItem> _notices = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _loadNotices();
  }

  Future<void> _loadNotices() async {
    try {
      final rest = ref.read(restClientProvider);
      final data = await rest.listNotices();
      final rawItems = data['notices'] as List? ?? [];
      setState(() {
        _notices = rawItems
            .map((e) => _NoticeItem.fromMap(e as Map<String, dynamic>))
            .toList();
        _isLoading = false;
      });
    } on DioException {
      // Fallback to default notices on API failure
      setState(() {
        _notices = _defaultNotices;
        _isLoading = false;
      });
    }
  }

  static final _defaultNotices = [
    _NoticeItem(
      title: 'ManPaSik v1.2.0 업데이트 안내',
      category: '업데이트',
      date: '2026-02-15',
      content: '새로운 기능이 추가되었습니다:\n'
          '- AI 건강 코칭 고도화\n'
          '- 카트리지 백과사전 추가\n'
          '- 가족 건강 리포트 기능\n'
          '- 성능 및 안정성 개선',
      isImportant: true,
    ),
    _NoticeItem(
      title: '개인정보 처리방침 변경 안내',
      category: '정책',
      date: '2026-02-10',
      content: '개인정보 처리방침이 일부 변경되었습니다. '
          '주요 변경 사항은 건강 데이터 처리 범위와 제3자 제공 관련 내용입니다.',
      isImportant: true,
    ),
  ];

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
                      '공지사항',
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
              child: _isLoading
                  ? const Center(
                      child: CircularProgressIndicator(
                          color: SanggamTheme.primary))
                  : _notices.isEmpty
                      ? const Center(
                          child: Text('공지사항이 없습니다.',
                              style: TextStyle(
                                  color: SanggamTheme.onSurfaceDim)))
                      : ListView.builder(
                          padding: const EdgeInsets.all(16),
                          itemCount: _notices.length,
                          itemBuilder: (_, index) =>
                              _buildNoticeCard(_notices[index]),
                        ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildNoticeCard(
      _NoticeItem notice) {
    final catColor = switch (notice.category) {
      '업데이트' => SanggamTheme.jagaeCyan,
      '정책' => SanggamTheme.primary,
      '점검' => SanggamTheme.error,
      '이벤트' => SanggamTheme.jagaeCyan,
      '상품' => SanggamTheme.jagaeMagenta,
      _ => SanggamTheme.onSurfaceDim,
    };

    return Padding(
      padding:
          const EdgeInsets.only(bottom: 8),
      child: SanggamContainer(
        borderRadius: 16,
        padding: EdgeInsets.zero,
        child: Theme(
          data: ThemeData.dark().copyWith(
            dividerColor:
                SanggamTheme.surfaceVariant,
          ),
          child: ExpansionTile(
            leading: notice.isImportant
                ? const Icon(Icons.campaign,
                    color: SanggamTheme.error)
                : const Icon(
                    Icons.article_outlined,
                    color: SanggamTheme
                        .onSurfaceDim),
            title: Text(
              notice.title,
              style: TextStyle(
                color: Colors.white,
                fontSize: 14,
                fontWeight: notice.isImportant
                    ? FontWeight.bold
                    : FontWeight.normal,
              ),
            ),
            subtitle: Row(
              children: [
                Container(
                  padding: const EdgeInsets
                      .symmetric(
                      horizontal: 8,
                      vertical: 2),
                  decoration: BoxDecoration(
                    color: catColor.withValues(
                        alpha: 0.1),
                    borderRadius:
                        BorderRadius.circular(
                            8),
                  ),
                  child: Text(
                      notice.category,
                      style: TextStyle(
                        fontSize: 10,
                        color: catColor,
                      )),
                ),
                const SizedBox(width: 8),
                Text(notice.date,
                    style: const TextStyle(
                      color: SanggamTheme
                          .onSurfaceDim,
                      fontSize: 12,
                    )),
              ],
            ),
            children: [
              Padding(
                padding:
                    const EdgeInsets.fromLTRB(
                        16, 0, 16, 16),
                child: Text(notice.content,
                    style: const TextStyle(
                      color: Colors.white,
                      fontSize: 14,
                    )),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _NoticeItem {
  final String title, category, date, content;
  final bool isImportant;
  const _NoticeItem({
    required this.title,
    required this.category,
    required this.date,
    required this.content,
    required this.isImportant,
  });

  factory _NoticeItem.fromMap(Map<String, dynamic> m) {
    return _NoticeItem(
      title: (m['title'] ?? '') as String,
      category: (m['category'] ?? '') as String,
      date: (m['date'] ?? m['created_at'] ?? '') as String,
      content: (m['content'] ?? '') as String,
      isImportant: (m['is_important'] as bool?) ?? false,
    );
  }
}
