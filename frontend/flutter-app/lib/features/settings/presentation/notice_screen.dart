import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

/// 공지사항 화면
class NoticeScreen extends StatelessWidget {
  const NoticeScreen({super.key});

  static final _notices = [
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
          '주요 변경 사항은 건강 데이터 처리 범위와 제3자 제공 관련 내용입니다. '
          '자세한 내용은 설정 > 개인정보처리방침에서 확인하실 수 있습니다.',
      isImportant: true,
    ),
    _NoticeItem(
      title: '서버 정기 점검 안내 (2/20 03:00~05:00)',
      category: '점검',
      date: '2026-02-08',
      content: '서비스 안정성 향상을 위한 정기 점검이 진행됩니다.\n'
          '점검 시간: 2026년 2월 20일 03:00 ~ 05:00 (KST)\n'
          '점검 중에는 측정 및 데이터 동기화가 일시 중단됩니다.',
      isImportant: false,
    ),
    _NoticeItem(
      title: '건강 챌린지 시즌 2 오픈!',
      category: '이벤트',
      date: '2026-02-01',
      content: '커뮤니티 건강 챌린지 시즌 2가 시작되었습니다. '
          '30일 혈당 관리 챌린지에 참여하고 특별 뱃지를 획득하세요!',
      isImportant: false,
    ),
    _NoticeItem(
      title: '신규 카트리지 출시 - 비타민D, 코르티솔',
      category: '상품',
      date: '2026-01-25',
      content: 'Premium/Professional 등급의 새로운 카트리지가 출시되었습니다:\n'
          '- 비타민D (25(OH)D): 골밀도 및 면역 관련 지표\n'
          '- 코르티솔: 스트레스 호르몬 수치 측정',
      isImportant: false,
    ),
  ];

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('공지사항'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => context.pop(),
        ),
      ),
      body: ListView.builder(
        padding: const EdgeInsets.all(16),
        itemCount: _notices.length,
        itemBuilder: (context, index) => _buildNoticeCard(theme, _notices[index]),
      ),
    );
  }

  Widget _buildNoticeCard(ThemeData theme, _NoticeItem notice) {
    final catColor = switch (notice.category) {
      '업데이트' => Colors.blue,
      '정책' => Colors.orange,
      '점검' => Colors.red,
      '이벤트' => Colors.green,
      '상품' => Colors.purple,
      _ => Colors.grey,
    };

    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: ExpansionTile(
        leading: notice.isImportant
            ? Icon(Icons.campaign, color: theme.colorScheme.error)
            : Icon(Icons.article_outlined, color: theme.colorScheme.outline),
        title: Text(
          notice.title,
          style: theme.textTheme.titleSmall?.copyWith(
            fontWeight: notice.isImportant ? FontWeight.bold : FontWeight.normal,
          ),
        ),
        subtitle: Row(
          children: [
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 1),
              decoration: BoxDecoration(
                color: catColor.withOpacity(0.1),
                borderRadius: BorderRadius.circular(8),
              ),
              child: Text(notice.category, style: TextStyle(fontSize: 10, color: catColor)),
            ),
            const SizedBox(width: 8),
            Text(notice.date, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.outline)),
          ],
        ),
        children: [
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 0, 16, 16),
            child: Text(notice.content, style: theme.textTheme.bodyMedium),
          ),
        ],
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
}
