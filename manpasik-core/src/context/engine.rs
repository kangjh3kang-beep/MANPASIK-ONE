// src/context/engine.rs
// ManPaSik (萬波息) v2.4.3 - §6.11 ContextEngine

#[derive(Debug, Clone, PartialEq, Eq, PartialOrd, Ord)]
pub enum ContextCardType {
    NutritionRecommendation = 1, // 영양 추천
    ProductSuggestion = 2,       // 제품 제안
    HabitNudge = 3,              // 습관 넛지
    CheckupReminder = 4,         // 검진 알림
    EnvironmentAlert = 5,        // 환경 연계
}

pub struct ContextCard {
    pub card_id: String,
    pub card_type: ContextCardType,
    pub priority: u8, // 1~100 (90+ 는 무조건 강제 노출)
    pub title: String,
    pub body: String,
    pub dismissed: bool,
}

pub struct UserContext {
    pub recent_glucose: Option<f64>,
    pub habit_streak_days: u16,
    pub has_danger_risks: bool,
}

pub struct ContextEngine;

impl ContextEngine {
    pub fn new() -> Self {
        Self
    }

    /// 사용자 컨텍스트를 입력받아 우선순위 정렬된 컨텍스트 카드 발급
    pub fn generate_cards(&self, context: &UserContext) -> Vec<ContextCard> {
        let mut cards = Vec::new();

        // 1. 긴급 위험 카드 (Danger Risk)
        if context.has_danger_risks {
            cards.push(ContextCard {
                card_id: "risk-1".into(),
                card_type: ContextCardType::CheckupReminder,
                priority: 95,
                title: "긴급 의료 지원 필요".into(),
                body: "최근 측정치에서 위험 수준의 징후가 포착되었습니다. 전문의를 찾아주세요."
                    .into(),
                dismissed: false,
            });
        }

        // 2. 혈당 컨텍스트(영양 추천)
        if let Some(glucose) = context.recent_glucose {
            if glucose > 120.0 {
                cards.push(ContextCard {
                    card_id: "nutri-1".into(),
                    card_type: ContextCardType::NutritionRecommendation,
                    priority: 80,
                    title: "크롬 및 식이섬유 보충 추천".into(),
                    body: "비정상 혈당 변동을 잡아줄 영양소를 추천합니다.".into(),
                    dismissed: false,
                });
            }
        }

        // 3. 습관 달성(넛지)
        if context.habit_streak_days < 3 {
            cards.push(ContextCard {
                card_id: "habit-1".into(),
                card_type: ContextCardType::HabitNudge,
                priority: 60,
                title: "규칙적 수면 패턴 깨짐".into(),
                body: "최근 3일간 수면 패턴이 틀어졌습니다. 오늘 저녁 스트레칭을 추천합니다."
                    .into(),
                dismissed: false,
            });
        }

        // 우선순위(Priority) 내림차순 정렬
        cards.sort_by(|a, b| b.priority.cmp(&a.priority));

        // 최대 10장 노출 제약
        cards.into_iter().take(10).collect()
    }
}
