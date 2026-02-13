import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/shared/widgets/measurement_card.dart';
import 'package:manpasik/shared/widgets/sanggam_decoration.dart';
import 'package:manpasik/shared/widgets/primary_button.dart';
import 'package:manpasik/l10n/app_localizations.dart';


class HomeScreen extends ConsumerWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final authState = ref.watch(authProvider);
    final historyAsync = ref.watch(measurementHistoryProvider);

    return Scaffold(
      body: SafeArea(
        child: Column(
          children: [
            // 커스텀 앱바
            Padding(
              padding: const EdgeInsets.fromLTRB(24, 24, 24, 0),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        AppLocalizations.of(context)!.greeting,
                        style: theme.textTheme.bodyLarge?.copyWith(
                          color: theme.colorScheme.onSurfaceVariant,
                        ),
                      ),
                      Text(
                        AppLocalizations.of(context)!.greetingWithName(authState.displayName ?? AppLocalizations.of(context)!.user),
                        style: theme.textTheme.headlineMedium?.copyWith(
                          fontWeight: FontWeight.bold,
                          color: theme.colorScheme.primary, // Sanggam Gold from theme
                        ),
                      ),
                    ],
                  ),
                  Row(
                    children: [
                      IconButton(
                        icon: const Icon(Icons.smart_toy_rounded),
                        tooltip: 'AI 건강 어시스턴트',
                        onPressed: () => context.push('/chat'),
                      ),
                      IconButton(
                        icon: const Icon(Icons.devices_rounded),
                        tooltip: AppLocalizations.of(context)!.devices,
                        onPressed: () => context.push('/devices'),
                      ),
                      IconButton(
                        icon: const Icon(Icons.settings_rounded),
                        tooltip: AppLocalizations.of(context)!.settings,
                        onPressed: () => context.push('/settings'),
                      ),
                    ],
                  ),
                ],
              ),
            ),

            const SizedBox(height: 32),

            // 측정 시작 카드 (Sanggam Premium)
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 24),
              child: SanggamContainer(
                borderRadius: 24,
                child: Padding(
                  padding: const EdgeInsets.all(24),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          Icon(
                            Icons.query_stats_rounded,
                            color: theme.colorScheme.secondary.withOpacity(0.9),
                          ),
                          const SizedBox(width: 8),
                          Text(
                            AppLocalizations.of(context)!.newMeasurement,
                            style: theme.textTheme.titleMedium?.copyWith(
                              color: theme.colorScheme.secondary,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 16),
                      Text(
                        AppLocalizations.of(context)!.checkHealth,
                        style: theme.textTheme.headlineSmall?.copyWith(
                          color: Colors.white,
                          fontWeight: FontWeight.bold,
                          height: 1.2,
                        ),
                      ),
                      const SizedBox(height: 24),
                      PrimaryButton(
                        text: AppLocalizations.of(context)!.startMeasurementAction,
                        icon: Icons.auto_awesome_rounded,
                        onPressed: () => context.push('/measurement'),
                      ),
                    ],
                  ),
                ),
              ),
            ),

            const SizedBox(height: 32),

            // 최근 기록 헤더
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 24),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    AppLocalizations.of(context)!.recentHistory,
                    style: theme.textTheme.titleLarge?.copyWith(
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  TextButton(
                    onPressed: () {},
                    child: Text(AppLocalizations.of(context)!.viewAll),
                  ),
                ],
              ),
            ),

            const SizedBox(height: 8),

            // 측정 기록 목록 (gRPC GetMeasurementHistory)
            Expanded(
              child: historyAsync.when(
                data: (result) {
                  if (result.items.isEmpty) {
                    return Center(
                      child: Text(
                        AppLocalizations.of(context)!.noDevicesRegistered,
                        style: theme.textTheme.bodyLarge?.copyWith(
                          color: theme.colorScheme.onSurfaceVariant,
                        ),
                      ),
                    );
                  }
                  return ListView.builder(
                    padding: const EdgeInsets.symmetric(horizontal: 24),
                    itemCount: result.items.length,
                    itemBuilder: (context, index) {
                      final item = result.items[index];
                      final date = item.measuredAt ?? DateTime.now();
                      final type = item.primaryValue <= 100
                          ? 'normal'
                          : (item.primaryValue <= 125 ? 'warning' : 'high');
                      return MeasurementCard(
                        date: date,
                        value: item.primaryValue,
                        unit: item.unit.isNotEmpty ? item.unit : 'mg/dL',
                        resultType: type,
                        onTap: () {},
                      );
                    },
                  );
                },
                loading: () => const Center(child: CircularProgressIndicator()),
                error: (err, _) => Center(
                  child: Padding(
                    padding: const EdgeInsets.all(24),
                    child: Text(
                      '기록을 불러올 수 없습니다.\n$err',
                      textAlign: TextAlign.center,
                      style: theme.textTheme.bodyMedium?.copyWith(
                        color: theme.colorScheme.error,
                      ),
                    ),
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
