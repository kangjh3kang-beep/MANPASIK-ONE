#!/bin/bash
# Phase 2: AppBar leading tooltip 일괄 추가
# 인라인 패턴: leading: IconButton(icon: const Icon(Icons.arrow_back), onPressed: ...)
# 멀티라인 패턴: icon: const Icon(Icons.arrow_back), 다음 줄에 tooltip 삽입

BASE="/home/kangjh3kang/Manpasik/frontend/flutter-app/lib"
FEAT="$BASE/features"
SHARED="$BASE/shared"

# -- community --
# challenge_screen: 인라인 arrow_back x2
sed -i "s|leading: IconButton(icon: const Icon(Icons.arrow_back), onPressed: () => context.pop())|leading: IconButton(icon: const Icon(Icons.arrow_back), tooltip: '뒤로 가기', onPressed: () => context.pop())|g" "$FEAT/community/presentation/challenge_screen.dart"

# create_post_screen: 멀티라인
sed -i '/leading: IconButton(/{n;/icon: const Icon(Icons.close),/{/tooltip/!{s|icon: const Icon(Icons.close),|icon: const Icon(Icons.close),\n          tooltip: '"'"'닫기'"'"',|}}}' "$FEAT/community/presentation/create_post_screen.dart"

# research_post_screen: 멀티라인 arrow_back
sed -i '/leading: IconButton(/{n;/icon: const Icon(Icons.arrow_back),/{/tooltip/!{s|icon: const Icon(Icons.arrow_back),|icon: const Icon(Icons.arrow_back),\n          tooltip: '"'"'뒤로 가기'"'"',|}}}' "$FEAT/community/presentation/research_post_screen.dart"

# qna_screen: 인라인 arrow_back + close
sed -i "s|leading: IconButton(icon: const Icon(Icons.arrow_back), onPressed: () => context.pop())|leading: IconButton(icon: const Icon(Icons.arrow_back), tooltip: '뒤로 가기', onPressed: () => context.pop())|g" "$FEAT/community/presentation/qna_screen.dart"
sed -i "s|leading: IconButton(icon: const Icon(Icons.close), onPressed: () => context.pop())|leading: IconButton(icon: const Icon(Icons.close), tooltip: '닫기', onPressed: () => context.pop())|g" "$FEAT/community/presentation/qna_screen.dart"

# post_detail_screen: 멀티라인
sed -i '/leading: IconButton(/{n;/icon: const Icon(Icons.arrow_back),/{/tooltip/!{s|icon: const Icon(Icons.arrow_back),|icon: const Icon(Icons.arrow_back),\n          tooltip: '"'"'뒤로 가기'"'"',|}}}' "$FEAT/community/presentation/post_detail_screen.dart"

# -- measurement --
sed -i '/leading: IconButton(/{n;/icon: const Icon(Icons.arrow_back),/{/tooltip/!{s|icon: const Icon(Icons.arrow_back),|icon: const Icon(Icons.arrow_back),\n          tooltip: '"'"'뒤로 가기'"'"',|}}}' "$FEAT/measurement/presentation/measurement_screen.dart"

sed -i '/leading: IconButton(/{n;/icon: const Icon(Icons.arrow_back),/{/tooltip/!{s|icon: const Icon(Icons.arrow_back),|icon: const Icon(Icons.arrow_back),\n          tooltip: '"'"'뒤로 가기'"'"',|}}}' "$FEAT/measurement/presentation/measurement_result_screen.dart"

# result_screen: close
sed -i '/leading: IconButton(/{n;/icon: const Icon(Icons.close),/{/tooltip/!{s|icon: const Icon(Icons.close),|icon: const Icon(Icons.close),\n          tooltip: '"'"'닫기'"'"',|}}}' "$FEAT/measurement/presentation/result_screen.dart"

# vision_analyzer
sed -i '/leading: IconButton(/{n;/icon: const Icon(Icons.arrow_back),/{/tooltip/!{s|icon: const Icon(Icons.arrow_back),|icon: const Icon(Icons.arrow_back),\n          tooltip: '"'"'뒤로 가기'"'"',|}}}' "$FEAT/measurement/presentation/vision_analyzer.dart"

# -- ai_coach --
sed -i '/leading: IconButton(/{n;/icon: const Icon(Icons.arrow_back),/{/tooltip/!{s|icon: const Icon(Icons.arrow_back),|icon: const Icon(Icons.arrow_back),\n          tooltip: '"'"'뒤로 가기'"'"',|}}}' "$FEAT/ai_coach/presentation/exercise_video_screen.dart"
sed -i '/leading: IconButton(/{n;/icon: const Icon(Icons.arrow_back),/{/tooltip/!{s|icon: const Icon(Icons.arrow_back),|icon: const Icon(Icons.arrow_back),\n          tooltip: '"'"'뒤로 가기'"'"',|}}}' "$FEAT/ai_coach/presentation/food_analysis_screen.dart"

# -- medical --
for f in telemedicine_screen facility_search_screen prescription_detail_screen consultation_result_screen video_call_screen; do
  sed -i '/leading: IconButton(/{n;/icon: const Icon(Icons.arrow_back),/{/tooltip/!{s|icon: const Icon(Icons.arrow_back),|icon: const Icon(Icons.arrow_back),\n          tooltip: '"'"'뒤로 가기'"'"',|}}}' "$FEAT/medical/presentation/${f}.dart"
  sed -i "s|leading: IconButton(icon: const Icon(Icons.arrow_back), onPressed: () => context.pop())|leading: IconButton(icon: const Icon(Icons.arrow_back), tooltip: '뒤로 가기', onPressed: () => context.pop())|g" "$FEAT/medical/presentation/${f}.dart"
done

# -- admin --
for f in admin_compliance_screen admin_hierarchy_screen admin_monitor_screen admin_dashboard_screen admin_audit_screen admin_users_screen admin_settings_screen admin_ecosystem_screen admin_revenue_screen admin_inventory_table; do
  sed -i '/leading: IconButton(/{n;/icon: const Icon(Icons.arrow_back),/{/tooltip/!{s|icon: const Icon(Icons.arrow_back),|icon: const Icon(Icons.arrow_back),\n          tooltip: '"'"'뒤로 가기'"'"',|}}}' "$FEAT/admin/presentation/${f}.dart"
  sed -i "s|leading: IconButton(icon: const Icon(Icons.arrow_back), onPressed: () => context.pop())|leading: IconButton(icon: const Icon(Icons.arrow_back), tooltip: '뒤로 가기', onPressed: () => context.pop())|g" "$FEAT/admin/presentation/${f}.dart"
  # close 패턴도 처리
  sed -i '/leading: IconButton(/{n;/icon: const Icon(Icons.close),/{/tooltip/!{s|icon: const Icon(Icons.close),|icon: const Icon(Icons.close),\n          tooltip: '"'"'닫기'"'"',|}}}' "$FEAT/admin/presentation/${f}.dart"
  sed -i "s|leading: IconButton(icon: const Icon(Icons.close), onPressed: () => context.pop())|leading: IconButton(icon: const Icon(Icons.close), tooltip: '닫기', onPressed: () => context.pop())|g" "$FEAT/admin/presentation/${f}.dart"
done

# -- settings --
for f in accessibility_screen notice_screen settings_screen data_export_screen notification_settings_screen emergency_settings_screen profile_edit_screen support_screen escalation_screen security_screen legal_screen inquiry_create_screen; do
  sed -i '/leading: IconButton(/{n;/icon: const Icon(Icons.arrow_back),/{/tooltip/!{s|icon: const Icon(Icons.arrow_back),|icon: const Icon(Icons.arrow_back),\n          tooltip: '"'"'뒤로 가기'"'"',|}}}' "$FEAT/settings/presentation/${f}.dart"
  sed -i "s|leading: IconButton(icon: const Icon(Icons.arrow_back), onPressed: () => context.pop())|leading: IconButton(icon: const Icon(Icons.arrow_back), tooltip: '뒤로 가기', onPressed: () => context.pop())|g" "$FEAT/settings/presentation/${f}.dart"
  sed -i '/leading: IconButton(/{n;/icon: const Icon(Icons.close),/{/tooltip/!{s|icon: const Icon(Icons.close),|icon: const Icon(Icons.close),\n          tooltip: '"'"'닫기'"'"',|}}}' "$FEAT/settings/presentation/${f}.dart"
  sed -i "s|leading: IconButton(icon: const Icon(Icons.close), onPressed: () => context.pop())|leading: IconButton(icon: const Icon(Icons.close), tooltip: '닫기', onPressed: () => context.pop())|g" "$FEAT/settings/presentation/${f}.dart"
done

# -- market --
for f in order_history_screen cartridge_detail_screen product_detail_screen checkout_screen plan_comparison_screen order_detail_screen subscription_cancel_screen subscription_screen cart_screen encyclopedia_screen market_screen; do
  sed -i '/leading: IconButton(/{n;/icon: const Icon(Icons.arrow_back),/{/tooltip/!{s|icon: const Icon(Icons.arrow_back),|icon: const Icon(Icons.arrow_back),\n          tooltip: '"'"'뒤로 가기'"'"',|}}}' "$FEAT/market/presentation/${f}.dart"
  sed -i "s|leading: IconButton(icon: const Icon(Icons.arrow_back), onPressed: () => context.pop())|leading: IconButton(icon: const Icon(Icons.arrow_back), tooltip: '뒤로 가기', onPressed: () => context.pop())|g" "$FEAT/market/presentation/${f}.dart"
  sed -i '/leading: IconButton(/{n;/icon: const Icon(Icons.close),/{/tooltip/!{s|icon: const Icon(Icons.close),|icon: const Icon(Icons.close),\n          tooltip: '"'"'닫기'"'"',|}}}' "$FEAT/market/presentation/${f}.dart"
  sed -i "s|leading: IconButton(icon: const Icon(Icons.close), onPressed: () => context.pop())|leading: IconButton(icon: const Icon(Icons.close), tooltip: '닫기', onPressed: () => context.pop())|g" "$FEAT/market/presentation/${f}.dart"
done

# -- shared widgets --
for f in pass_verification_webview toss_payment_webview; do
  sed -i '/leading: IconButton(/{n;/icon: const Icon(Icons.close),/{/tooltip/!{s|icon: const Icon(Icons.close),|icon: const Icon(Icons.close),\n          tooltip: '"'"'닫기'"'"',|}}}' "$SHARED/widgets/${f}.dart"
  sed -i "s|leading: IconButton(icon: const Icon(Icons.close), onPressed: () => context.pop())|leading: IconButton(icon: const Icon(Icons.close), tooltip: '닫기', onPressed: () => context.pop())|g" "$SHARED/widgets/${f}.dart"
  sed -i '/leading: IconButton(/{n;/icon: const Icon(Icons.arrow_back),/{/tooltip/!{s|icon: const Icon(Icons.arrow_back),|icon: const Icon(Icons.arrow_back),\n          tooltip: '"'"'뒤로 가기'"'"',|}}}' "$SHARED/widgets/${f}.dart"
done

echo "Phase 2 complete"

# 검증: tooltip 개수
TOTAL=$(grep -r "tooltip:" "$FEAT" "$SHARED" --include="*.dart" | wc -l)
echo "Total tooltip occurrences: $TOTAL"
