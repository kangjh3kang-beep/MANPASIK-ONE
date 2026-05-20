import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/theme/app_theme.dart';

class DataScreen extends ConsumerWidget {
  const DataScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      extendBodyBehindAppBar: true,
      appBar: AppBar(
        title: const Text('Data Vault'),
        centerTitle: true,
        flexibleSpace: ClipRect(
          child: BackdropFilter(
            filter: ImageFilter.blur(sigmaX: 10, sigmaY: 10),
            child: Container(color: AppTheme.backgroundDark.withValues(alpha: 0.6)),
          ),
        ),
      ),
      body: Stack(
        children: [
          Positioned(
            bottom: 100, left: -50,
            child: Container(
              width: 300, height: 300,
              decoration: BoxDecoration(shape: BoxShape.circle, color: Colors.blueAccent.withValues(alpha: 0.1)),
              child: BackdropFilter(filter: ImageFilter.blur(sigmaX: 100, sigmaY: 100), child: Container()),
            ),
          ),
          SafeArea(
            child: ListView(
              padding: const EdgeInsets.all(20),
              physics: const BouncingScrollPhysics(),
              children: [
                _buildDataCard('16-Channel Raw Diff', 'CRDT 동기화 대기 중: 0건', Icons.graphic_eq, '324 MB'),
                _buildDataCard('896-Dim AI Fingerprint', 'PostgreSQL 영구 아카이빙 완료', Icons.fingerprint, '1.2 GB'),
                _buildDataCard('Self-Healing DB', '시스템 복구 이벤트: 42건', Icons.memory, '15 MB'),
                _buildDataCard('Digital Twin 3D 모델', '업데이트 버전: v2.4.3', Icons.view_in_ar, '450 MB'),
              ],
            ),
          )
        ],
      ),
    );
  }

  Widget _buildDataCard(String title, String subtitle, IconData icon, String size) {
    return Container(
      margin: const EdgeInsets.only(bottom: 16),
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: AppTheme.surfaceDark.withValues(alpha: 0.5),
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: Colors.white.withValues(alpha: 0.05)),
      ),
      child: Row(
        children: [
          Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(color: Colors.blueAccent.withValues(alpha: 0.1), shape: BoxShape.circle),
            child: Icon(icon, color: Colors.blueAccent),
          ),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(title, style: const TextStyle(color: Colors.white, fontSize: 16, fontWeight: FontWeight.bold)),
                const SizedBox(height: 4),
                Text(subtitle, style: TextStyle(color: Colors.white.withValues(alpha: 0.5), fontSize: 13)),
              ],
            ),
          ),
          Text(size, style: const TextStyle(color: AppTheme.primaryNeonTeal, fontWeight: FontWeight.bold)),
        ],
      ),
    );
  }
}
