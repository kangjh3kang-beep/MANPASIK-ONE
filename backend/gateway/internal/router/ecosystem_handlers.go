package router

import (
	"net/http"
)

// ============================================================================
// Rewards Service Handlers (Mock — gRPC 미연결)
// ============================================================================

// GET /api/v1/rewards/{userId}/balance
func (r *Router) handleGetRewardBalance(w http.ResponseWriter, req *http.Request) {
	userID := req.PathValue("userId")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "userId가 필요합니다")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"userId":       userID,
		"balance":      2847,
		"currency":     "MPS",
		"totalEarned":  12480,
		"weeklyChange": 650,
	})
}

// GET /api/v1/rewards/{userId}/contributions
func (r *Router) handleGetRewardContributions(w http.ResponseWriter, req *http.Request) {
	userID := req.PathValue("userId")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "userId가 필요합니다")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"userId": userID,
		"contributions": []map[string]interface{}{
			{"id": "1", "date": "2026-04-26", "type": "심박 데이터", "points": 150, "quality": "A+"},
			{"id": "2", "date": "2026-04-25", "type": "혈압 데이터", "points": 120, "quality": "A"},
			{"id": "3", "date": "2026-04-25", "type": "수면 패턴", "points": 200, "quality": "A+"},
			{"id": "4", "date": "2026-04-24", "type": "활동량 데이터", "points": 80, "quality": "B+"},
			{"id": "5", "date": "2026-04-23", "type": "식이 기록", "points": 100, "quality": "A"},
		},
	})
}

// ============================================================================
// GxP Compliance Service Handlers (Mock — gRPC 미연결)
// ============================================================================

// GET /api/v1/gxp/audit-logs
func (r *Router) handleGetGxpAuditLogs(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"audit_logs": []map[string]interface{}{
			{"id": "1", "user": "김품질", "action": "배치 검사 승인", "batchId": "B-2026-0421", "hasSignature": true, "timestamp": "14:32:07"},
			{"id": "2", "user": "이감사", "action": "원자재 입고 확인", "batchId": "B-2026-0420", "hasSignature": true, "timestamp": "11:15:33"},
			{"id": "3", "user": "박규정", "action": "출하 승인", "batchId": "B-2026-0419", "hasSignature": false, "timestamp": "09:45:12"},
		},
	})
}

// GET /api/v1/gxp/compliance
func (r *Router) handleGetGxpCompliance(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"rules": []map[string]interface{}{
			{"ruleId": "KGMP 의약품 제조관리", "status": "pass", "itemsCount": 41, "passedCount": 41},
			{"ruleId": "GLP 비임상시험관리", "status": "pass", "itemsCount": 28, "passedCount": 27},
			{"ruleId": "GCP 임상시험관리", "status": "pass", "itemsCount": 35, "passedCount": 35},
			{"ruleId": "GDP 유통관리", "status": "warning", "itemsCount": 22, "passedCount": 20},
		},
	})
}

// ============================================================================
// Developer Portal Service Handlers (Mock — gRPC 미연결)
// ============================================================================

// GET /api/v1/dev-portal/keys
func (r *Router) handleGetDevPortalKeys(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"keys": []map[string]interface{}{
			{"id": "1", "key": "mmup_live_sk_x7k9...d6e8", "name": "Production Key", "createdAt": "2026-03-15", "status": "active"},
			{"id": "2", "key": "mmup_test_sk_a1b2...o5p6", "name": "Test Key", "createdAt": "2026-04-01", "status": "active"},
			{"id": "3", "key": "mmup_dev_sk_q1w2...g5h6", "name": "Dev Key", "createdAt": "2026-01-10", "status": "revoked"},
		},
	})
}

// GET /api/v1/dev-portal/endpoints
func (r *Router) handleGetDevPortalEndpoints(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"endpoints": []map[string]interface{}{
			{"method": "GET", "path": "/api/v1/patients", "description": "환자 목록 조회", "auth": "Bearer", "status": "stable"},
			{"method": "POST", "path": "/api/v1/measurements", "description": "측정 데이터 전송", "auth": "Bearer", "status": "stable"},
			{"method": "GET", "path": "/api/v1/biomarkers/:id", "description": "바이오마커 조회", "auth": "Bearer", "status": "stable"},
			{"method": "POST", "path": "/api/v1/predictions/run", "description": "AI 예측 실행", "auth": "API Key", "status": "beta"},
			{"method": "GET", "path": "/api/v1/rewards/balance", "description": "토큰 잔액 조회", "auth": "Bearer", "status": "stable"},
			{"method": "POST", "path": "/api/v1/fhir/bundle", "description": "FHIR 번들 전송", "auth": "mTLS", "status": "beta"},
		},
	})
}

// ============================================================================
// Telemedicine Service Handlers (Mock — gRPC 미연결)
// ============================================================================

// GET /api/v1/telemedicine/doctors
func (r *Router) handleListTelemedicineDoctors(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"doctors": []map[string]interface{}{
			{"id": "1", "name": "김내과", "specialty": "내과", "rating": 4.8, "reviews": 234, "available": true, "nextSlot": "오늘 15:00", "fee": "₩25,000"},
			{"id": "2", "name": "이피부", "specialty": "피부과", "rating": 4.9, "reviews": 189, "available": true, "nextSlot": "오늘 16:30", "fee": "₩30,000"},
			{"id": "3", "name": "박가정", "specialty": "가정의학과", "rating": 4.7, "reviews": 312, "available": false, "nextSlot": "내일 09:00", "fee": "₩20,000"},
			{"id": "4", "name": "최소아", "specialty": "소아과", "rating": 4.9, "reviews": 156, "available": true, "nextSlot": "오늘 17:00", "fee": "₩25,000"},
		},
	})
}

// ============================================================================
// Admin Extra Handlers (Mock — gRPC 미연결)
// ============================================================================

// GET /api/v1/admin/kpis
func (r *Router) handleGetAdminKPIs(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"kpis": []map[string]interface{}{
			{"label": "총 사용자", "value": "12,847", "change": "+3.2%", "icon": "Users", "color": "text-sky-600 bg-sky-50"},
			{"label": "오늘 측정", "value": "3,421", "change": "+12%", "icon": "Activity", "color": "text-emerald-600 bg-emerald-50"},
			{"label": "활성 리더기", "value": "8,923", "change": "+1.5%", "icon": "Package", "color": "text-purple-600 bg-purple-50"},
			{"label": "이상 알림", "value": "7", "change": "-2건", "icon": "AlertTriangle", "color": "text-amber-600 bg-amber-50"},
		},
	})
}

// GET /api/v1/admin/events
func (r *Router) handleGetAdminEvents(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"events": []map[string]interface{}{
			{"time": "14:32", "type": "측정", "desc": "서울 강남구 — 혈당 측정 이상치 감지", "severity": "warning"},
			{"time": "14:28", "type": "리더기", "desc": "MPK-A1B2C3 펌웨어 업데이트 완료", "severity": "info"},
			{"time": "14:15", "type": "사용자", "desc": "신규 가입 +23명 (일일 목표 120%)", "severity": "success"},
			{"time": "13:50", "type": "재고", "desc": "혈당 카트리지 재고 부족 경고 (창고 B)", "severity": "warning"},
			{"time": "13:30", "type": "결제", "desc": "정기배송 결제 처리 완료 342건", "severity": "info"},
			{"time": "12:45", "type": "품질", "desc": "Lot-2026-A 보정 검증 통과", "severity": "success"},
		},
	})
}
