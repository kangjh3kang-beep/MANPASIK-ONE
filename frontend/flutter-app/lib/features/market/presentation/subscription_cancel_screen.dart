import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/core/theme/app_theme.dart';

/// 구독 해지 화면
class SubscriptionCancelScreen extends StatefulWidget {
  const SubscriptionCancelScreen({super.key});

  @override
  State<SubscriptionCancelScreen> createState() =>
      _SubscriptionCancelScreenState();
}

class _SubscriptionCancelScreenState extends State<SubscriptionCancelScreen> {
  String? _selectedReason;
  final _detailController = TextEditingController();
  bool _isCancelling = false;

  static const _reasons = [
    '서비스를 더 이상 사용하지 않음',
    '가격이 너무 비쌈',
    '다른 서비스로 전환',
    '필요한 기능이 부족',
    '기기를 분실/교체함',
    '기타',
  ];

  @override
  void dispose() {
    _detailController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('구독 해지'),
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
            // 현재 구독 정보
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text('현재 구독 정보',
                        style: theme.textTheme.titleMedium
                            ?.copyWith(fontWeight: FontWeight.bold)),
                    const Divider(),
                    _buildInfoRow('플랜', 'Pro 구독'),
                    _buildInfoRow('월 요금', '₩29,900'),
                    _buildInfoRow('다음 결제일', '2026-03-18'),
                    _buildInfoRow('구독 시작일', '2025-12-18'),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 16),

            // 혜택 안내
            Container(
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: AppTheme.sanggamGold.withValues(alpha: 0.1),
                borderRadius: BorderRadius.circular(12),
                border: Border.all(color: AppTheme.sanggamGold),
              ),
              child: Column(
                children: [
                  const Icon(Icons.info_outline,
                      color: AppTheme.sanggamGold, size: 32),
                  const SizedBox(height: 8),
                  Text('해지 시 잃게 되는 혜택',
                      style: theme.textTheme.titleSmall
                          ?.copyWith(fontWeight: FontWeight.bold)),
                  const SizedBox(height: 8),
                  ...[
                    'AI 건강 코칭 이용',
                    '원격 진료 서비스',
                    '가족 그룹 (최대 6명)',
                    'FHIR 데이터 내보내기',
                    '고급 카트리지 접근',
                  ].map((feature) => Padding(
                        padding: const EdgeInsets.symmetric(vertical: 2),
                        child: Row(
                          children: [
                            const Icon(Icons.remove_circle_outline,
                                size: 16, color: Colors.red),
                            const SizedBox(width: 8),
                            Text(feature),
                          ],
                        ),
                      )),
                ],
              ),
            ),
            const SizedBox(height: 24),

            // 해지 사유 선택
            Text('해지 사유를 선택해주세요',
                style: theme.textTheme.titleMedium
                    ?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: 8),
            ...(_reasons.map((reason) => RadioListTile<String>(
                  title: Text(reason),
                  value: reason,
                  groupValue: _selectedReason,
                  onChanged: (val) => setState(() => _selectedReason = val),
                  contentPadding: EdgeInsets.zero,
                ))),
            const SizedBox(height: 8),

            // 상세 의견
            TextField(
              controller: _detailController,
              maxLines: 3,
              decoration: const InputDecoration(
                labelText: '추가 의견 (선택)',
                hintText: '서비스 개선을 위해 의견을 남겨주세요',
                border: OutlineInputBorder(),
              ),
            ),
            const SizedBox(height: 24),

            // 해지 버튼
            FilledButton(
              onPressed:
                  _selectedReason != null && !_isCancelling ? _handleCancel : null,
              style: FilledButton.styleFrom(
                backgroundColor: theme.colorScheme.error,
                minimumSize: const Size.fromHeight(48),
              ),
              child: _isCancelling
                  ? const SizedBox(
                      width: 20,
                      height: 20,
                      child: CircularProgressIndicator(
                          strokeWidth: 2, color: Colors.white))
                  : const Text('구독 해지하기'),
            ),
            const SizedBox(height: 8),

            // 유지 버튼
            OutlinedButton(
              onPressed: () => context.pop(),
              style: OutlinedButton.styleFrom(
                minimumSize: const Size.fromHeight(48),
                foregroundColor: AppTheme.sanggamGold,
              ),
              child: const Text('구독 유지하기'),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildInfoRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label, style: const TextStyle(color: Colors.grey)),
          Text(value, style: const TextStyle(fontWeight: FontWeight.w600)),
        ],
      ),
    );
  }

  void _handleCancel() {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('정말 해지하시겠습니까?'),
        content: const Text(
            '현재 결제 기간이 끝난 후 Pro 혜택이 중단됩니다.\n무료 플랜으로 전환됩니다.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text('취소'),
          ),
          FilledButton(
            onPressed: () {
              Navigator.pop(ctx);
              setState(() => _isCancelling = true);
              Future.delayed(const Duration(seconds: 2), () {
                if (mounted) {
                  setState(() => _isCancelling = false);
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(
                        content: Text('구독이 해지되었습니다. 현재 결제 기간까지 이용 가능합니다.')),
                  );
                  context.pop();
                }
              });
            },
            style: FilledButton.styleFrom(
                backgroundColor: Theme.of(ctx).colorScheme.error),
            child: const Text('해지 확인'),
          ),
        ],
      ),
    );
  }
}
