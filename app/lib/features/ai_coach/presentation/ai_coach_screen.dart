import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/theme/app_theme.dart';

class AiCoachScreen extends ConsumerWidget {
  const AiCoachScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      extendBodyBehindAppBar: true,
      appBar: AppBar(
        title: const Text('AI Personal Advisor'),
        centerTitle: true,
        flexibleSpace: ClipRect(
          child: BackdropFilter(
            filter: ImageFilter.blur(sigmaX: 10, sigmaY: 10),
            child: Container(color: AppTheme.backgroundDark.withValues(alpha: 0.5)),
          ),
        ),
      ),
      body: Stack(
        children: [
          // 배경 네온 글로우
          Positioned(
            top: 100, left: -50,
            child: Container(
              width: 300, height: 300,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                color: AppTheme.primaryNeonPurple.withValues(alpha: 0.15),
              ),
              child: BackdropFilter(filter: ImageFilter.blur(sigmaX: 100, sigmaY: 100), child: Container()),
            ),
          ),
          
          SafeArea(
            child: Column(
              children: [
                _buildGlowAvatar(),
                const SizedBox(height: 10),
                const Text('ManPaSik AI', style: TextStyle(fontSize: 22, fontWeight: FontWeight.bold, color: Colors.white, letterSpacing: 1.2)),
                const Text('항상 당신을 지켜보고 있습니다.', style: TextStyle(fontSize: 14, color: AppTheme.primaryNeonTeal)),
                const SizedBox(height: 20),
                Expanded(
                  child: ListView(
                    physics: const BouncingScrollPhysics(),
                    padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 10),
                    children: [
                      _buildChatBubble('현재 심박 변이도(HRV)가 평소보다 15% 낮아졌습니다. 스트레스가 누적된 상태인가요?', isBot: true),
                      _buildChatBubble('네, 어제 야근을 좀 했어요.', isBot: false),
                      _buildChatBubble('그렇군요. 오늘은 격렬한 웨이트 트레이닝보다 30분 가벼운 요가를 권장합니다. 하단 마켓에서 마그네슘 보충제를 살펴보시겠어요?', isBot: true, isInteractive: true),
                      _buildChatBubble('마그네슘이 수면에 도움이 될까요?', isBot: false),
                      _buildChatBubble('네, 16채널 혈류 피처를 분석한 결과 칼슘 채널 활성도가 떨어져 있습니다. 마그네슘이 근육 이완과 숙면에 직접적인 도움을 줍니다.', isBot: true),
                      const SizedBox(height: 20),
                    ],
                  ),
                ),
                _buildInputArea(),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildGlowAvatar() {
    return Container(
      margin: const EdgeInsets.only(top: 20),
      width: 100, height: 100,
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        gradient: AppTheme.premiumGradient,
        boxShadow: [
          BoxShadow(color: AppTheme.primaryNeonPurple.withValues(alpha: 0.5), blurRadius: 30, spreadRadius: 5)
        ],
      ),
      child: const Center(
        child: Icon(Icons.auto_awesome, color: Colors.white, size: 50),
      ),
    );
  }

  Widget _buildChatBubble(String text, {required bool isBot, bool isInteractive = false}) {
    return Align(
      alignment: isBot ? Alignment.centerLeft : Alignment.centerRight,
      child: Container(
        margin: const EdgeInsets.only(bottom: 16),
        constraints: const BoxConstraints(maxWidth: 300),
        decoration: BoxDecoration(
          color: isBot ? AppTheme.surfaceDark : AppTheme.primaryNeonTeal.withValues(alpha: 0.15),
          borderRadius: BorderRadius.only(
            topLeft: const Radius.circular(20),
            topRight: const Radius.circular(20),
            bottomLeft: Radius.circular(isBot ? 0 : 20),
            bottomRight: Radius.circular(isBot ? 20 : 0),
          ),
          border: Border.all(
            color: isBot ? Colors.white.withValues(alpha: 0.05) : AppTheme.primaryNeonTeal.withValues(alpha: 0.5),
          ),
        ),
        padding: const EdgeInsets.symmetric(horizontal: 18, vertical: 14),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              text,
              style: TextStyle(
                color: isBot ? Colors.white.withValues(alpha: 0.9) : AppTheme.primaryNeonTeal,
                fontSize: 15,
                height: 1.5,
              ),
            ),
            if (isInteractive) ...[
              const SizedBox(height: 12),
              ElevatedButton.icon(
                onPressed: () {},
                icon: const Icon(Icons.shopping_bag, size: 16),
                label: const Text('보충제 추천 보기'),
                style: ElevatedButton.styleFrom(
                  backgroundColor: AppTheme.primaryNeonPurple,
                  foregroundColor: Colors.white,
                  minimumSize: const Size(double.infinity, 36),
                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                ),
              )
            ]
          ],
        ),
      ),
    );
  }

  Widget _buildInputArea() {
    return Container(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 30),
      decoration: BoxDecoration(
        color: AppTheme.surfaceDark,
        borderRadius: const BorderRadius.only(topLeft: Radius.circular(24), topRight: Radius.circular(24)),
      ),
      child: Row(
        children: [
          Expanded(
            child: Container(
              height: 50,
              decoration: BoxDecoration(
                color: AppTheme.backgroundDark,
                borderRadius: BorderRadius.circular(25),
                border: Border.all(color: Colors.white.withValues(alpha: 0.1)),
              ),
              padding: const EdgeInsets.symmetric(horizontal: 20),
              child: const TextField(
                style: TextStyle(color: Colors.white),
                decoration: InputDecoration(
                  border: InputBorder.none,
                  hintText: 'AI 코치에게 증상이나 목표를 물어보세요...',
                  hintStyle: TextStyle(color: Colors.white30, fontSize: 13),
                ),
              ),
            ),
          ),
          const SizedBox(width: 12),
          Container(
            width: 50, height: 50,
            decoration: const BoxDecoration(
              shape: BoxShape.circle,
              gradient: AppTheme.premiumGradient,
            ),
            child: const Icon(Icons.send, color: Colors.white, size: 20),
          )
        ],
      ),
    );
  }
}
