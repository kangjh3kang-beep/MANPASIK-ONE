// Package service??measurement-service??鍮꾩쫰?덉뒪 濡쒖쭅??援ы쁽?⑸땲??
//
// Measurement business logic.
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/manpasik/backend/shared/assay"
	apperrors "github.com/manpasik/backend/shared/errors"
	"go.uber.org/zap"
)

// MeasurementService??痢≪젙 ?쒕퉬?ㅼ쓽 鍮꾩쫰?덉뒪 濡쒖쭅?낅땲??
type MeasurementService struct {
	logger         *zap.Logger
	sessionRepo    SessionRepository
	measureRepo    MeasurementRepository
	vectorRepo     VectorRepository
	eventPublisher EventPublisher
	searchIndexer  SearchIndexer // optional: nil?대㈃ ?몃뜳??鍮꾪솢?깊솕
}

// SessionRepository??痢≪젙 ?몄뀡 ??μ냼 ?명꽣?섏씠?ㅼ엯?덈떎 (PostgreSQL).
type SessionRepository interface {
	CreateSession(ctx context.Context, session *MeasurementSession) error
	GetSession(ctx context.Context, sessionID string) (*MeasurementSession, error)
	EndSession(ctx context.Context, sessionID string, totalMeasurements int, endedAt time.Time) error
}

// MeasurementRepository??痢≪젙 ?곗씠????μ냼 ?명꽣?섏씠?ㅼ엯?덈떎 (TimescaleDB).
type MeasurementRepository interface {
	Store(ctx context.Context, data *MeasurementData) error
	GetHistory(ctx context.Context, userID string, start, end time.Time, limit, offset int) ([]*MeasurementSummary, int, error)
}

// VectorRepository???묎굅?꾨┛??踰≫꽣 ??μ냼 ?명꽣?섏씠?ㅼ엯?덈떎 (Milvus).
type VectorRepository interface {
	StoreFingerprint(ctx context.Context, sessionID string, vector []float32) error
	SearchSimilar(ctx context.Context, vector []float32, topK int) ([]SimilarResult, error)
}

// EventPublisher???대깽??諛쒗뻾 ?명꽣?섏씠?ㅼ엯?덈떎 (Kafka).
type EventPublisher interface {
	PublishMeasurementCompleted(ctx context.Context, event *MeasurementCompletedEvent) error
}

// SearchIndexer??痢≪젙 ?곗씠??寃???몃뜳???명꽣?섏씠?ㅼ엯?덈떎 (Elasticsearch).
type SearchIndexer interface {
	IndexMeasurement(ctx context.Context, sessionID string, data *MeasurementData) error
}

// MeasurementSession? 痢≪젙 ?몄뀡 ?뷀떚?곗엯?덈떎.
type MeasurementSession struct {
	ID                string
	DeviceID          string
	CartridgeID       string
	UserID            string
	StartedAt         time.Time
	EndedAt           *time.Time
	TotalMeasurements int
	Status            string // "active", "completed", "error"
}

// MeasurementData??媛쒕퀎 痢≪젙 ?곗씠?곗엯?덈떎.
type MeasurementData struct {
	Time              time.Time
	SessionID         string
	DeviceID          string
	UserID            string
	CartridgeType     string
	RawChannels       []float64
	SDet              float64
	SRef              float64
	Alpha             float64
	SCorrected        float64
	PrimaryValue      float64
	Unit              string
	Confidence        float64
	EvidenceStatus    assay.EvidenceStatus
	DiagnosticReady   bool
	EvidenceGaps      []string
	FingerprintVector []float32
	TempC             float32
	HumidityPct       float32
	BatteryPct        int
}

// MeasurementSummary??痢≪젙 寃곌낵 ?붿빟?낅땲??
type MeasurementSummary struct {
	SessionID       string
	CartridgeType   string
	PrimaryValue    float64
	Unit            string
	EvidenceStatus  assay.EvidenceStatus
	DiagnosticReady bool
	EvidenceGaps    []string
	MeasuredAt      time.Time
}

// SimilarResult???좎궗 踰≫꽣 寃??寃곌낵?낅땲??
type SimilarResult struct {
	SessionID string
	Score     float32
	Distance  float32
}

// MeasurementCompletedEvent??痢≪젙 ?꾨즺 ?대깽?몄엯?덈떎.
type MeasurementCompletedEvent struct {
	SessionID    string
	UserID       string
	DeviceID     string
	PrimaryValue float64
	Unit         string
	CompletedAt  time.Time
}

// NewMeasurementService????MeasurementService瑜??앹꽦?⑸땲??
func NewMeasurementService(
	logger *zap.Logger,
	sessionRepo SessionRepository,
	measureRepo MeasurementRepository,
	vectorRepo VectorRepository,
	eventPublisher EventPublisher,
) *MeasurementService {
	return &MeasurementService{
		logger:         logger,
		sessionRepo:    sessionRepo,
		measureRepo:    measureRepo,
		vectorRepo:     vectorRepo,
		eventPublisher: eventPublisher,
	}
}

// SetSearchIndexer??寃???몃뜳?쒕? ?ㅼ젙?⑸땲??(optional).
func (s *MeasurementService) SetSearchIndexer(indexer SearchIndexer) {
	s.searchIndexer = indexer
}

// StartSession? ??痢≪젙 ?몄뀡???쒖옉?⑸땲??
func (s *MeasurementService) StartSession(
	ctx context.Context,
	deviceID, cartridgeID, userID string,
) (*MeasurementSession, error) {
	if deviceID == "" || cartridgeID == "" || userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "device_id, cartridge_id, user_id are required")
	}

	session := &MeasurementSession{
		ID:          uuid.New().String(),
		DeviceID:    deviceID,
		CartridgeID: cartridgeID,
		UserID:      userID,
		StartedAt:   time.Now().UTC(),
		Status:      "active",
	}

	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
		s.logger.Error("?몄뀡 ?앹꽦 ?ㅽ뙣",
			zap.String("device_id", deviceID),
			zap.Error(err),
		)
		return nil, apperrors.New(apperrors.ErrInternal, "痢≪젙 ?몄뀡 ?쒖옉???ㅽ뙣?덉뒿?덈떎")
	}

	s.logger.Info("痢≪젙 ?몄뀡 ?쒖옉",
		zap.String("session_id", session.ID),
		zap.String("device_id", deviceID),
		zap.String("cartridge_id", cartridgeID),
	)

	return session, nil
}

// ProcessMeasurement???ㅽ듃由щ컢 痢≪젙 ?곗씠?곕? 泥섎━?⑸땲??
func (s *MeasurementService) ProcessMeasurement(
	ctx context.Context,
	data *MeasurementData,
) (*ProcessedResult, error) {
	if data.SessionID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "session_id is required")
	}

	session, err := s.sessionRepo.GetSession(ctx, data.SessionID)
	if err != nil || session == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "measurement session not found")
	}
	if session.Status != "active" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "measurement session is not active")
	}
	if data.DeviceID == "" {
		data.DeviceID = session.DeviceID
	}
	if data.UserID == "" {
		data.UserID = session.UserID
	}
	if data.CartridgeType == "" {
		data.CartridgeType = session.CartridgeID
	}
	if data.Time.IsZero() {
		data.Time = time.Now().UTC()
	}
	if err := applyAssaySemantics(data); err != nil {
		return nil, err
	}

	if err := s.measureRepo.Store(ctx, data); err != nil {
		s.logger.Error("measurement data store failed",
			zap.String("session_id", data.SessionID),
			zap.Error(err),
		)
		return nil, apperrors.New(apperrors.ErrInternal, "measurement data store failed")
	}

	if len(data.FingerprintVector) > 0 {
		if err := s.vectorRepo.StoreFingerprint(ctx, data.SessionID, data.FingerprintVector); err != nil {
			s.logger.Warn("fingerprint vector store failed",
				zap.String("session_id", data.SessionID),
				zap.Error(err),
			)
		}
	}

	return &ProcessedResult{
		SessionID:       data.SessionID,
		PrimaryValue:    data.PrimaryValue,
		Unit:            data.Unit,
		Confidence:      data.Confidence,
		EvidenceStatus:  data.EvidenceStatus,
		DiagnosticReady: data.DiagnosticReady,
		EvidenceGaps:    append([]string(nil), data.EvidenceGaps...),
		ProcessedAt:     time.Now().UTC(),
	}, nil
}

func applyAssaySemantics(data *MeasurementData) error {
	result, err := assay.Evaluate(data.CartridgeType, assay.Signal{
		SCorrected:      data.SCorrected,
		RawChannelCount: len(data.RawChannels),
	})
	if err != nil {
		return apperrors.New(apperrors.ErrInvalidInput, err.Error())
	}
	if data.PrimaryValue == 0 {
		data.PrimaryValue = result.PrimaryValue
	}
	if data.Unit == "" {
		data.Unit = result.Unit
	}
	if data.Confidence == 0 {
		data.Confidence = result.Confidence
	}
	data.EvidenceStatus = result.Definition.Evidence.Status
	data.DiagnosticReady = result.Definition.IsDiagnosticReady()
	data.EvidenceGaps = result.Definition.EvidenceGaps()
	return nil
}

// EndSession? 痢≪젙 ?몄뀡??醫낅즺?⑸땲??
func (s *MeasurementService) EndSession(
	ctx context.Context,
	sessionID string,
) (*SessionEndResult, error) {
	if sessionID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "session_id is required")
	}

	session, err := s.sessionRepo.GetSession(ctx, sessionID)
	if err != nil || session == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "痢≪젙 ?몄뀡??李얠쓣 ???놁뒿?덈떎")
	}

	endedAt := time.Now().UTC()
	if err := s.sessionRepo.EndSession(ctx, sessionID, session.TotalMeasurements, endedAt); err != nil {
		return nil, apperrors.New(apperrors.ErrInternal, "?몄뀡 醫낅즺???ㅽ뙣?덉뒿?덈떎")
	}

	// 痢≪젙 ?꾨즺 ?대깽??諛쒗뻾 (Kafka)
	event := &MeasurementCompletedEvent{
		SessionID:   sessionID,
		UserID:      session.UserID,
		DeviceID:    session.DeviceID,
		CompletedAt: endedAt,
	}
	if err := s.eventPublisher.PublishMeasurementCompleted(ctx, event); err != nil {
		s.logger.Warn("?대깽??諛쒗뻾 ?ㅽ뙣 (鍮꾩튂紐낆쟻)",
			zap.String("session_id", sessionID),
			zap.Error(err),
		)
	}

	s.logger.Info("痢≪젙 ?몄뀡 醫낅즺",
		zap.String("session_id", sessionID),
		zap.Int("total_measurements", session.TotalMeasurements),
	)

	return &SessionEndResult{
		SessionID:         sessionID,
		TotalMeasurements: session.TotalMeasurements,
		EndedAt:           endedAt,
	}, nil
}

// GetHistory??痢≪젙 湲곕줉??議고쉶?⑸땲??
func (s *MeasurementService) GetHistory(
	ctx context.Context,
	userID string,
	start, end time.Time,
	limit, offset int,
) ([]*MeasurementSummary, int, error) {
	if userID == "" {
		return nil, 0, apperrors.New(apperrors.ErrInvalidInput, "user_id is required")
	}

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.measureRepo.GetHistory(ctx, userID, start, end, limit, offset)
}

// defaultLOINCMap? 諛붿씠?ㅻ쭏而??대쫫?????LOINC 肄붾뱶 留ㅽ븨?낅땲??
var defaultLOINCMap = map[string]string{
	"blood_glucose":     "15074-8",
	"blood_pressure":    "85354-9",
	"cholesterol_total": "2093-3",
	"hemoglobin_a1c":    "4548-4",
	"heart_rate":        "8867-4",
	"body_temperature":  "8310-5",
	"oxygen_saturation": "2708-6",
}

// loincCodeFor??諛붿씠?ㅻ쭏而??대쫫?????LOINC 肄붾뱶瑜?諛섑솚?⑸땲?? 留ㅽ븨???놁쑝硫?湲곕낯媛믪쓣 諛섑솚?⑸땲??
func loincCodeFor(cartridgeType string) string {
	if code, ok := defaultLOINCMap[cartridgeType]; ok {
		return code
	}
	return "29463-7" // 湲곕낯 LOINC 肄붾뱶
}

// ExportSingleMeasurement???⑥씪 ?몄뀡??痢≪젙 寃곌낵瑜?FHIR R4 Observation Bundle JSON?쇰줈 ?대낫?낅땲??
func (s *MeasurementService) ExportSingleMeasurement(ctx context.Context, sessionID string) (string, error) {
	if sessionID == "" {
		return "", apperrors.New(apperrors.ErrInvalidInput, "session_id is required")
	}

	// 1. ?몄뀡 議고쉶
	session, err := s.sessionRepo.GetSession(ctx, sessionID)
	if err != nil || session == nil {
		return "", apperrors.New(apperrors.ErrNotFound, "痢≪젙 ?몄뀡??李얠쓣 ???놁뒿?덈떎")
	}

	// 2. ?대떦 ?몄뀡??痢≪젙 ?곗씠??議고쉶 (?좎? 湲곗??쇰줈 媛?몄삩 ???몄뀡 ?꾪꽣留?
	summaries, _, err := s.measureRepo.GetHistory(ctx, session.UserID, time.Time{}, time.Now().UTC(), 500, 0)
	if err != nil {
		return "", fmt.Errorf("痢≪젙 ?곗씠??議고쉶 ?ㅽ뙣: %w", err)
	}

	// 3. Filter by session ID and convert to FHIR Observation resources.
	observations := make([]map[string]interface{}, 0)
	for _, m := range summaries {
		if m.SessionID != sessionID {
			continue
		}
		obs := map[string]interface{}{
			"resourceType": "Observation",
			"status":       "final",
			"code": map[string]interface{}{
				"coding": []map[string]interface{}{{
					"system":  "http://loinc.org",
					"code":    loincCodeFor(m.CartridgeType),
					"display": m.CartridgeType,
				}},
			},
			"subject": map[string]interface{}{
				"reference": "Patient/" + session.UserID,
			},
			"effectiveDateTime": m.MeasuredAt.Format(time.RFC3339),
			"valueQuantity": map[string]interface{}{
				"value":  m.PrimaryValue,
				"unit":   m.Unit,
				"system": "http://unitsofmeasure.org",
				"code":   m.Unit,
			},
			"identifier": []map[string]interface{}{{
				"system": "https://manpasik.com/session",
				"value":  m.SessionID,
			}},
		}
		observations = append(observations, obs)
	}

	// 4. FHIR Bundle 援ъ꽦
	entries := make([]map[string]interface{}, 0, len(observations))
	for _, o := range observations {
		entries = append(entries, map[string]interface{}{"resource": o})
	}
	bundle := map[string]interface{}{
		"resourceType": "Bundle",
		"type":         "collection",
		"entry":        entries,
	}

	jsonBytes, err := json.Marshal(bundle)
	if err != nil {
		return "", fmt.Errorf("FHIR Bundle 吏곷젹???ㅽ뙣: %w", err)
	}
	return string(jsonBytes), nil
}

// ExportToFHIRObservations??痢≪젙 寃곌낵瑜?FHIR R4 Observation Bundle濡??대낫?낅땲??(Agent D).
// biomarkerNames媛 鍮꾩뼱 ?덉쑝硫??꾩껜 諛붿씠?ㅻ쭏而??ы븿.
func (s *MeasurementService) ExportToFHIRObservations(
	ctx context.Context,
	userID string,
	fromDate, toDate *time.Time,
	biomarkerNames []string,
) (fhirBundleJSON string, observationCount int, err error) {
	if userID == "" {
		return "", 0, apperrors.New(apperrors.ErrInvalidInput, "user_id is required")
	}
	var start, end time.Time
	if fromDate != nil {
		start = *fromDate
	}
	if toDate != nil {
		end = *toDate
	} else {
		end = time.Now().UTC()
	}
	summaries, _, err := s.measureRepo.GetHistory(ctx, userID, start, end, 500, 0)
	if err != nil {
		return "", 0, err
	}
	biomarkerSet := make(map[string]struct{})
	for _, n := range biomarkerNames {
		biomarkerSet[n] = struct{}{}
	}
	observations := make([]map[string]interface{}, 0, len(summaries))
	for _, m := range summaries {
		if len(biomarkerSet) > 0 {
			if _, ok := biomarkerSet[m.CartridgeType]; !ok {
				continue
			}
		}
		obs := map[string]interface{}{
			"resourceType": "Observation",
			"status":       "final",
			"code": map[string]interface{}{
				"coding": []map[string]interface{}{{
					"system":  "http://loinc.org",
					"code":    "29463-7",
					"display": m.CartridgeType,
				}},
			},
			"subject": map[string]interface{}{
				"reference": "Patient/" + userID,
			},
			"effectiveDateTime": m.MeasuredAt.Format(time.RFC3339),
			"valueQuantity": map[string]interface{}{
				"value":  m.PrimaryValue,
				"unit":   m.Unit,
				"system": "http://unitsofmeasure.org",
				"code":   m.Unit,
			},
			"identifier": []map[string]interface{}{{
				"system": "https://manpasik.com/session",
				"value":  m.SessionID,
			}},
		}
		observations = append(observations, obs)
	}
	entries := make([]map[string]interface{}, 0, len(observations))
	for _, o := range observations {
		entries = append(entries, map[string]interface{}{"resource": o})
	}
	bundle := map[string]interface{}{
		"resourceType": "Bundle",
		"type":         "collection",
		"entry":        entries,
	}
	jsonBytes, err := json.Marshal(bundle)
	if err != nil {
		return "", 0, fmt.Errorf("FHIR Bundle 吏곷젹???ㅽ뙣: %w", err)
	}
	return string(jsonBytes), len(observations), nil
}

// DigitalTwinState???붿????몄쐢 ?숆린???곹깭?낅땲??
type DigitalTwinState struct {
	TwinID                string
	SessionID             string
	UserID                string
	DeviceID              string
	Residuals             []float64
	EWMAValue             float64
	CUSUMPos              float64
	CUSUMNeg              float64
	HealthState           string // "nominal", "warning", "drift", "recalibrate"
	DriftScore            float64
	RemainingMeasurements int
	FingerprintVector     []float32
	FingerprintDim        int
	SyncedAt              time.Time
}

// CalibrationStatus??蹂댁젙 ?곹깭?낅땲??
type CalibrationStatus struct {
	CalibrationState         string // "valid", "expiring", "expired", "needs_recalibration"
	DriftScore               float64
	MeasurementsSinceCal     int
	MaxMeasurementsBeforeCal int
	LastCalibratedAt         time.Time
	NextCalibrationAt        time.Time
}

// SyncDigitalTwin? ?붿????몄쐢 ?곹깭瑜??숆린?뷀빀?덈떎.
func (s *MeasurementService) SyncDigitalTwin(
	ctx context.Context,
	sessionID, userID, deviceID string,
	residuals []float64,
	ewmaValue, cusumPos, cusumNeg float64,
	healthState string,
	driftScore float64,
	remainingMeasurements int,
	fingerprintVector []float32,
	fingerprintDim int,
) (*DigitalTwinState, error) {
	if sessionID == "" || userID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "session_id and user_id are required")
	}

	// ?몄뀡 議댁옱 ?뺤씤
	session, err := s.sessionRepo.GetSession(ctx, sessionID)
	if err != nil || session == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "痢≪젙 ?몄뀡??李얠쓣 ???놁뒿?덈떎")
	}

	// ?쒕━?꾪듃 ?먯닔 湲곕컲 沅뚯옣 議곗튂 寃곗젙
	recommendedAction := "continue"
	if driftScore > 0.8 {
		recommendedAction = "replace_cartridge"
	} else if driftScore > 0.5 {
		recommendedAction = "recalibrate"
	}

	state := &DigitalTwinState{
		TwinID:                uuid.New().String(),
		SessionID:             sessionID,
		UserID:                userID,
		DeviceID:              deviceID,
		Residuals:             residuals,
		EWMAValue:             ewmaValue,
		CUSUMPos:              cusumPos,
		CUSUMNeg:              cusumNeg,
		HealthState:           healthState,
		DriftScore:            driftScore,
		RemainingMeasurements: remainingMeasurements,
		FingerprintVector:     fingerprintVector,
		FingerprintDim:        fingerprintDim,
		SyncedAt:              time.Now().UTC(),
	}

	if len(fingerprintVector) > 0 {
		if err := s.vectorRepo.StoreFingerprint(ctx, sessionID, fingerprintVector); err != nil {
			s.logger.Warn("?붿????몄쐢 踰≫꽣 ????ㅽ뙣 (鍮꾩튂紐낆쟻)", zap.Error(err))
		}
	}

	// ?쒕━?꾪듃 寃쎄퀬 ?대깽??諛쒗뻾
	if driftScore > 0.5 {
		s.logger.Warn("移댄듃由ъ? ?쒕━?꾪듃 媛먯?",
			zap.String("session_id", sessionID),
			zap.Float64("drift_score", driftScore),
			zap.String("action", recommendedAction),
		)
	}

	s.logger.Info("?붿????몄쐢 ?숆린???꾨즺",
		zap.String("twin_id", state.TwinID),
		zap.String("session_id", sessionID),
		zap.String("health_state", healthState),
		zap.Float64("drift_score", driftScore),
	)

	return state, nil
}

// GetCalibrationStatus??蹂댁젙 ?곹깭瑜?議고쉶?⑸땲??
func (s *MeasurementService) GetCalibrationStatus(
	ctx context.Context,
	sessionID, deviceID string,
) (*CalibrationStatus, error) {
	if sessionID == "" && deviceID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "session_id or device_id is required")
	}

	// 湲곕낯 蹂댁젙 ?곹깭 諛섑솚 (硫붾え由???μ냼 ??Phase A?먯꽌??DB ?놁씠 湲곕낯媛?
	now := time.Now().UTC()
	maxMeasurements := 500 // 移댄듃由ъ???理쒕? 痢≪젙 ?잛닔

	status := &CalibrationStatus{
		CalibrationState:         "valid",
		DriftScore:               0.0,
		MeasurementsSinceCal:     0,
		MaxMeasurementsBeforeCal: maxMeasurements,
		LastCalibratedAt:         now.Add(-24 * time.Hour),
		NextCalibrationAt:        now.Add(72 * time.Hour),
	}

	return status, nil
}

// ProcessedResult??泥섎━??痢≪젙 寃곌낵?낅땲??
type ProcessedResult struct {
	SessionID       string
	PrimaryValue    float64
	Unit            string
	Confidence      float64
	EvidenceStatus  assay.EvidenceStatus
	DiagnosticReady bool
	EvidenceGaps    []string
	ProcessedAt     time.Time
}

// SessionEndResult???몄뀡 醫낅즺 寃곌낵?낅땲??
type SessionEndResult struct {
	SessionID         string
	TotalMeasurements int
	EndedAt           time.Time
}

// SearchSimilarFingerprints??踰≫꽣 ?좎궗??湲곕컲?쇰줈 ?좎궗??痢≪젙??寃?됲빀?덈떎 (Milvus).
func (s *MeasurementService) SearchSimilarFingerprints(
	ctx context.Context,
	vector []float32,
	topK int,
) ([]SimilarResult, error) {
	if len(vector) == 0 {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "寃??踰≫꽣媛 鍮꾩뼱 ?덉뒿?덈떎")
	}
	if topK <= 0 {
		topK = 10
	}

	results, err := s.vectorRepo.SearchSimilar(ctx, vector, topK)
	if err != nil {
		s.logger.Warn("踰≫꽣 ?좎궗??寃???ㅽ뙣", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "?좎궗??寃?됱뿉 ?ㅽ뙣?덉뒿?덈떎")
	}

	return results, nil
}

// ?ъ슜?섏? ?딅뒗 import 諛⑹?
var _ = fmt.Sprintf
