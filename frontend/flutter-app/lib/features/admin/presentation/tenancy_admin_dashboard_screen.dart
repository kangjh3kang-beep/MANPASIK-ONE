import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/widgets/sanggam_container.dart';

/// Tenancy 운영자 대시보드 (Phase AM-3).
///
/// admin/ops 사용자가 통계를 한눈에 보는 화면. /ops/tenancy/stats 와
/// /ops/tenancy/audit/stats 를 조회하여 카드로 표시.
class TenancyAdminDashboardScreen extends ConsumerStatefulWidget {
  const TenancyAdminDashboardScreen({
    super.key,
    this.opsBaseUrl,
    this.targetTenantId,
  });

  /// admin-service:9100 등 ops 엔드포인트 baseURL. 미설정 시 기본 gateway 사용.
  final String? opsBaseUrl;

  /// audit 통계 조회 대상 tenant. 미설정 시 화면에서 입력받음.
  final String? targetTenantId;

  @override
  ConsumerState<TenancyAdminDashboardScreen> createState() =>
      _TenancyAdminDashboardScreenState();
}

class _TenancyAdminDashboardScreenState
    extends ConsumerState<TenancyAdminDashboardScreen> {
  Map<String, dynamic>? _tenancyStats;
  Map<String, dynamic>? _auditStats;
  List<Map<String, dynamic>> _auditFailures = const [];
  List<Map<String, dynamic>> _dlqEntries = const [];
  String? _error;
  bool _loading = false;
  late final TextEditingController _tenantController;

  @override
  void initState() {
    super.initState();
    _tenantController =
        TextEditingController(text: widget.targetTenantId ?? '');
    _refresh();
  }

  @override
  void dispose() {
    _tenantController.dispose();
    super.dispose();
  }

  Future<void> _refresh() async {
    setState(() {
      _loading = true;
      _error = null;
    });
    try {
      final client = ref.read(restClientProvider);
      final stats =
          await client.getTenancyStats(opsBaseUrl: widget.opsBaseUrl);
      Map<String, dynamic>? auditStats;
      var failures = <Map<String, dynamic>>[];
      final tid = _tenantController.text.trim();
      if (tid.isNotEmpty) {
        auditStats = await client.getAuditStats(
          tenantId: tid,
          opsBaseUrl: widget.opsBaseUrl,
        );
        // 실패 호출 조회는 best-effort — Audit 통계가 성공해도 실패 호출 엔드포인트가
        // 미설정/비활성일 수 있음 (DB pool 미설정 환경 등).
        try {
          failures = await client.getAuditFailures(
            tenantId: tid,
            limit: 10,
            opsBaseUrl: widget.opsBaseUrl,
          );
        } catch (_) {
          failures = const [];
        }
      }
      // DLQ 는 tenant 무관 — 항상 시도. 실패해도 다른 카드 표시는 유지.
      var dlq = <Map<String, dynamic>>[];
      try {
        dlq = await client.listWebhookDLQ(opsBaseUrl: widget.opsBaseUrl);
      } catch (_) {
        dlq = const [];
      }
      setState(() {
        _tenancyStats = stats;
        _auditStats = auditStats;
        _auditFailures = failures;
        _dlqEntries = dlq;
        _loading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _loading = false;
      });
    }
  }

  Future<void> _replayDLQ(int id) async {
    try {
      final client = ref.read(restClientProvider);
      await client.replayWebhookDLQ(id: id, opsBaseUrl: widget.opsBaseUrl);
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('DLQ #$id 재발송 성공')),
      );
      await _refresh();
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('재발송 실패: $e')),
      );
    }
  }

  Future<void> _dropDLQ(int id) async {
    try {
      final client = ref.read(restClientProvider);
      await client.dropWebhookDLQ(id: id, opsBaseUrl: widget.opsBaseUrl);
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('DLQ #$id 폐기 완료')),
      );
      await _refresh();
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('폐기 실패: $e')),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: SanggamTheme.background,
      body: SafeArea(
        child: ListView(
          padding: const EdgeInsets.all(16),
          children: [
            _buildHeader(context),
            const SizedBox(height: 16),
            _buildTenantSelector(),
            const SizedBox(height: 16),
            if (_loading)
              const Center(
                child: CircularProgressIndicator(color: SanggamTheme.primary),
              )
            else if (_error != null)
              _buildError(_error!)
            else ...[
              if (_tenancyStats != null) _buildTenancyCard(_tenancyStats!),
              const SizedBox(height: 12),
              if (_auditStats != null) _buildAuditCard(_auditStats!),
              if (_auditFailures.isNotEmpty) ...[
                const SizedBox(height: 12),
                _buildAuditFailuresCard(_auditFailures),
              ],
              const SizedBox(height: 12),
              _buildDLQCard(_dlqEntries),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildHeader(BuildContext context) {
    return Row(
      children: [
        IconButton(
          icon: const Icon(Icons.arrow_back, color: SanggamTheme.primary),
          onPressed: () => Navigator.of(context).maybePop(),
        ),
        const SizedBox(width: 8),
        const Expanded(
          child: Text(
            'Tenancy 운영 대시보드',
            style: TextStyle(
              color: SanggamTheme.primary,
              fontSize: 20,
              fontWeight: FontWeight.w600,
            ),
          ),
        ),
        IconButton(
          key: const Key('refresh-button'),
          icon: const Icon(Icons.refresh, color: SanggamTheme.primary),
          tooltip: '새로고침',
          onPressed: _loading ? null : _refresh,
        ),
      ],
    );
  }

  Widget _buildTenantSelector() {
    return SanggamContainer(
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              'Audit 조회 대상 Tenant',
              style: TextStyle(color: SanggamTheme.onSurfaceDim, fontSize: 12),
            ),
            const SizedBox(height: 6),
            Row(
              children: [
                Expanded(
                  child: TextField(
                    key: const Key('tenant-input'),
                    controller: _tenantController,
                    style: const TextStyle(color: Colors.white),
                    decoration: const InputDecoration(
                      hintText: 'tenant_id 입력 (예: hospA)',
                      hintStyle: TextStyle(color: SanggamTheme.onSurfaceDim),
                      filled: true,
                      fillColor: SanggamTheme.surfaceVariant,
                      border:
                          OutlineInputBorder(borderSide: BorderSide.none),
                    ),
                  ),
                ),
                const SizedBox(width: 8),
                ElevatedButton(
                  key: const Key('apply-tenant'),
                  onPressed: _loading ? null : _refresh,
                  style: ElevatedButton.styleFrom(
                    backgroundColor: SanggamTheme.primary,
                    foregroundColor: SanggamTheme.onPrimary,
                  ),
                  child: const Text('조회'),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildTenancyCard(Map<String, dynamic> stats) {
    final invitations = (stats['invitations'] is Map<String, dynamic>)
        ? stats['invitations'] as Map<String, dynamic>
        : <String, dynamic>{};
    return SanggamContainer(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              'Tenancy 전체',
              style: TextStyle(
                color: SanggamTheme.primary,
                fontSize: 14,
                fontWeight: FontWeight.w600,
              ),
            ),
            const SizedBox(height: 12),
            _statRow('활성 조직', '${stats['tenant_count'] ?? 0}'),
            _statRow('활성 멤버', '${stats['active_member_count'] ?? 0}'),
            _statRow('비활성 멤버', '${stats['inactive_member_count'] ?? 0}'),
            const Divider(color: SanggamTheme.surfaceVariant, height: 24),
            const Text(
              '초대 상태',
              style: TextStyle(
                color: SanggamTheme.onSurfaceDim,
                fontSize: 11,
              ),
            ),
            const SizedBox(height: 6),
            _statRow('대기 중', '${invitations['pending'] ?? 0}'),
            _statRow('수락', '${invitations['accepted'] ?? 0}'),
            _statRow('만료', '${invitations['expired'] ?? 0}'),
            _statRow('취소', '${invitations['revoked'] ?? 0}'),
          ],
        ),
      ),
    );
  }

  Widget _buildAuditCard(Map<String, dynamic> stats) {
    return SanggamContainer(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Icon(Icons.bar_chart,
                    color: SanggamTheme.primary, size: 18),
                const SizedBox(width: 6),
                Text(
                  'LLM Audit — ${stats['TenantID'] ?? ''}',
                  style: const TextStyle(
                    color: SanggamTheme.primary,
                    fontSize: 14,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 12),
            _statRow('총 호출', '${stats['TotalCalls'] ?? 0}'),
            _statRow('성공', '${stats['SuccessfulCalls'] ?? 0}'),
            _statRow('실패', '${stats['FailedCalls'] ?? 0}'),
            _statRow('누적 토큰', '${stats['TotalTokens'] ?? 0}'),
            _statRow(
              '${stats['WindowHours'] ?? 24}h 토큰',
              '${stats['TokensInWindow'] ?? 0}',
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildAuditFailuresCard(List<Map<String, dynamic>> failures) {
    return SanggamContainer(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Icon(Icons.error_outline,
                    color: SanggamTheme.error, size: 18),
                const SizedBox(width: 6),
                Text(
                  '최근 실패 호출 (${failures.length})',
                  style: const TextStyle(
                    color: SanggamTheme.error,
                    fontSize: 14,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 8),
            for (final f in failures.take(10))
              Padding(
                padding: const EdgeInsets.symmetric(vertical: 4),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      '${f['Provider'] ?? '?'} / ${f['Model'] ?? '?'} — ${_truncated(f['ErrorMessage'])}',
                      style: const TextStyle(
                        color: Colors.white,
                        fontSize: 12,
                      ),
                    ),
                    Text(
                      '${f['UserID'] ?? ''}  ${f['RecordedAt'] ?? ''}',
                      style: const TextStyle(
                        color: SanggamTheme.onSurfaceDim,
                        fontSize: 10,
                      ),
                    ),
                  ],
                ),
              ),
          ],
        ),
      ),
    );
  }

  Widget _buildDLQCard(List<Map<String, dynamic>> entries) {
    return SanggamContainer(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                const Icon(Icons.report_problem,
                    color: SanggamTheme.primary, size: 18),
                const SizedBox(width: 6),
                Text(
                  'Webhook DLQ (${entries.length})',
                  style: const TextStyle(
                    color: SanggamTheme.primary,
                    fontSize: 14,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 8),
            if (entries.isEmpty)
              const Text(
                '대기 중인 항목 없음',
                style: TextStyle(
                  color: SanggamTheme.onSurfaceDim,
                  fontSize: 12,
                ),
              )
            else
              for (final e in entries)
                _buildDLQRow(e),
          ],
        ),
      ),
    );
  }

  Widget _buildDLQRow(Map<String, dynamic> entry) {
    final id = (entry['id'] is num) ? (entry['id'] as num).toInt() : 0;
    final ev = (entry['event'] is Map<String, dynamic>)
        ? entry['event'] as Map<String, dynamic>
        : <String, dynamic>{};
    final lastErr = entry['last_error']?.toString() ?? '';
    final attempts = entry['attempts'];
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 6),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  '#$id  ${ev['type'] ?? '?'}  · tenant=${ev['tenant_id'] ?? ''}',
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 12,
                    fontWeight: FontWeight.w600,
                  ),
                ),
                Text(
                  '시도 $attempts회 · ${_truncated(lastErr)}',
                  style: const TextStyle(
                    color: SanggamTheme.onSurfaceDim,
                    fontSize: 10,
                  ),
                ),
              ],
            ),
          ),
          IconButton(
            key: Key('dlq-replay-$id'),
            icon: const Icon(Icons.refresh,
                color: SanggamTheme.primary, size: 18),
            tooltip: '재발송',
            onPressed: _loading ? null : () => _replayDLQ(id),
          ),
          IconButton(
            key: Key('dlq-drop-$id'),
            icon:
                const Icon(Icons.delete, color: SanggamTheme.error, size: 18),
            tooltip: '폐기',
            onPressed: _loading ? null : () => _dropDLQ(id),
          ),
        ],
      ),
    );
  }

  String _truncated(Object? s) {
    final str = s?.toString() ?? '';
    if (str.length <= 80) return str;
    return '${str.substring(0, 80)}...';
  }

  Widget _statRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        children: [
          Expanded(
            child: Text(
              label,
              style: const TextStyle(color: Colors.white, fontSize: 13),
            ),
          ),
          Text(
            value,
            style: const TextStyle(
              color: SanggamTheme.primary,
              fontSize: 14,
              fontWeight: FontWeight.w600,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildError(String message) {
    return SanggamContainer(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            const Icon(Icons.error_outline,
                color: SanggamTheme.error, size: 32),
            const SizedBox(height: 8),
            Text(
              message,
              textAlign: TextAlign.center,
              style: const TextStyle(color: SanggamTheme.error, fontSize: 12),
            ),
            const SizedBox(height: 8),
            TextButton(
              onPressed: _refresh,
              child: const Text('다시 시도',
                  style: TextStyle(color: SanggamTheme.primary)),
            ),
          ],
        ),
      ),
    );
  }
}
