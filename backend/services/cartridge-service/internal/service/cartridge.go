// Package service는 cartridge-service의 비즈니스 로직을 구현합니다.
package service

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	apperrors "github.com/manpasik/backend/shared/errors"
	"go.uber.org/zap"
)

// ============================================================================
// 카트리지 카테고리 정의
// ============================================================================

// CartridgeCategory는 카트리지 카테고리 코드입니다.
type CartridgeCategory int32

const (
	CategoryUnknown        CartridgeCategory = 0
	CategoryHealthBiomarker CartridgeCategory = 1
	CategoryEnvironmental  CartridgeCategory = 2
	CategoryFoodSafety     CartridgeCategory = 3
	CategoryElectronicSensor CartridgeCategory = 4
	CategoryAdvancedAnalysis CartridgeCategory = 5
	CategoryIndustrial     CartridgeCategory = 6
	CategoryVeterinary     CartridgeCategory = 7
	CategoryPharmaceutical CartridgeCategory = 8
	CategoryAgricultural   CartridgeCategory = 9
	CategoryCosmetic       CartridgeCategory = 10
	CategoryForensic       CartridgeCategory = 11
	CategoryMarine         CartridgeCategory = 12
	CategoryBeta           CartridgeCategory = 254
	CategoryCustomResearch CartridgeCategory = 255
)

// CategoryInfo는 카테고리 정보입니다.
type CategoryInfo struct {
	Code        int32
	NameEN      string
	NameKO      string
	Description string
	TypeCount   int32
	IsActive    bool
}

// ============================================================================
// 카트리지 타입 정보 (레지스트리 항목)
// ============================================================================

// CartridgeTypeInfo는 카트리지 타입 상세 정보입니다.
type CartridgeTypeInfo struct {
	CategoryCode     int32
	TypeIndex        int32
	LegacyCode       int32
	NameEN           string
	NameKO           string
	Description      string
	RequiredChannels int32
	MeasurementSecs  int32
	Unit             string
	ReferenceRange   string
	IsActive         bool
	IsBeta           bool
	Manufacturer     string
}

// ============================================================================
// 카트리지 상세 (NFC 태그 파싱 결과)
// ============================================================================

// CartridgeDetail은 NFC 태그에서 읽은 카트리지 상세 정보입니다.
type CartridgeDetail struct {
	CartridgeUID        string
	CategoryCode        int32
	TypeIndex           int32
	LegacyCode          int32
	NameKO              string
	NameEN              string
	LotID               string
	ExpiryDate          string // YYYYMMDD
	RemainingUses       int32
	MaxUses             int32
	AlphaCoefficient    float64
	TempCoefficient     float64
	HumidityCoefficient float64
	RequiredChannels    int32
	MeasurementSecs     int32
	Unit                string
	ReferenceRange      string
	IsValid             bool
	ValidationMessage   string
}

// ============================================================================
// 카트리지 사용 기록
// ============================================================================

// CartridgeUsageRecord는 카트리지 사용 이력 항목입니다.
type CartridgeUsageRecord struct {
	RecordID     string
	UserID       string
	SessionID    string
	CartridgeUID string
	CategoryCode int32
	TypeIndex    int32
	TypeNameKO   string
	UsedAt       time.Time
}

// ============================================================================
// 카트리지 잔여 사용 정보
// ============================================================================

// CartridgeRemainingInfo는 카트리지 잔여 사용 정보입니다.
type CartridgeRemainingInfo struct {
	CartridgeUID  string
	RemainingUses int32
	MaxUses       int32
	ExpiryDate    string
	IsExpired     bool
}

// ============================================================================
// 카트리지 유효성 검증 결과
// ============================================================================

// CartridgeAccessLevel은 카트리지 접근 레벨입니다.
type CartridgeAccessLevel int32

const (
	AccessUnknown    CartridgeAccessLevel = 0
	AccessIncluded   CartridgeAccessLevel = 1
	AccessLimited    CartridgeAccessLevel = 2
	AccessAddOn      CartridgeAccessLevel = 3
	AccessRestricted CartridgeAccessLevel = 4
	AccessBeta       CartridgeAccessLevel = 5
)

// ValidateResult는 카트리지 유효성 검증 결과입니다.
type ValidateResult struct {
	IsValid       bool
	Reason        string
	RemainingUses int32
	AccessLevel   CartridgeAccessLevel
	Detail        *CartridgeDetail
}

// ============================================================================
// 저장소 인터페이스
// ============================================================================

// CartridgeUsageRepository는 카트리지 사용 기록 저장소입니다.
type CartridgeUsageRepository interface {
	Create(ctx context.Context, record *CartridgeUsageRecord) error
	ListByUserID(ctx context.Context, userID string, limit, offset int32) ([]*CartridgeUsageRecord, int32, error)
}

// CartridgeStateRepository는 카트리지별 잔여 사용 상태 저장소입니다.
type CartridgeStateRepository interface {
	GetByUID(ctx context.Context, uid string) (*CartridgeRemainingInfo, error)
	Upsert(ctx context.Context, info *CartridgeRemainingInfo) error
	DecrementUses(ctx context.Context, uid string) (int32, error)
}

// ============================================================================
// 카트리지 레지스트리 (30종 내장)
// ============================================================================

// 기본 카테고리 정의 (15개)
var defaultCategories = []*CategoryInfo{
	{Code: 0, NameEN: "Unknown", NameKO: "알 수 없음", Description: "알 수 없는 카테고리", IsActive: false},
	{Code: 1, NameEN: "HealthBiomarker", NameKO: "건강 바이오마커", Description: "혈액/타액/체액 기반 건강 바이오마커", IsActive: true},
	{Code: 2, NameEN: "Environmental", NameKO: "환경 모니터링", Description: "수질, 공기질, 방사선 등 환경 모니터링", IsActive: true},
	{Code: 3, NameEN: "FoodSafety", NameKO: "식품 안전", Description: "농약, 신선도, 알레르겐 등 식품 안전 검사", IsActive: true},
	{Code: 4, NameEN: "ElectronicSensor", NameKO: "전자코/전자혀", Description: "전자코, 전자혀, EHD 기체 분석", IsActive: true},
	{Code: 5, NameEN: "AdvancedAnalysis", NameKO: "고급 분석", Description: "비표적/다중패널 고급 분석", IsActive: true},
	{Code: 6, NameEN: "Industrial", NameKO: "산업용", Description: "산업용 분석 (Phase 3-4)", IsActive: false},
	{Code: 7, NameEN: "Veterinary", NameKO: "수의학", Description: "수의학 분석 (Phase 3-4)", IsActive: false},
	{Code: 8, NameEN: "Pharmaceutical", NameKO: "제약", Description: "제약 분석 (Phase 3-4)", IsActive: false},
	{Code: 9, NameEN: "Agricultural", NameKO: "농업", Description: "농업 분석 (Phase 3-4)", IsActive: false},
	{Code: 10, NameEN: "Cosmetic", NameKO: "화장품", Description: "화장품 분석 (Phase 3-4)", IsActive: false},
	{Code: 11, NameEN: "Forensic", NameKO: "법의학", Description: "법의학 분석 (Phase 3-4)", IsActive: false},
	{Code: 12, NameEN: "Marine", NameKO: "해양", Description: "해양 분석 (Phase 3-4)", IsActive: false},
	{Code: 254, NameEN: "Beta", NameKO: "베타/실험용", Description: "베타 및 실험용 카트리지", IsActive: true},
	{Code: 255, NameEN: "CustomResearch", NameKO: "맞춤형 연구", Description: "맞춤형 연구용 카트리지", IsActive: true},
}

// 기본 카트리지 타입 레지스트리 (30종, Rust nfc 모듈과 동일)
var defaultTypes = []*CartridgeTypeInfo{
	// HealthBiomarker 14종
	{CategoryCode: 1, TypeIndex: 1, LegacyCode: 0x01, NameEN: "Glucose", NameKO: "혈당", RequiredChannels: 88, MeasurementSecs: 15, Unit: "mg/dL", ReferenceRange: "70-100", IsActive: true, Manufacturer: "ManPaSik"},
	{CategoryCode: 1, TypeIndex: 2, LegacyCode: 0x02, NameEN: "LipidPanel", NameKO: "지질 패널", RequiredChannels: 88, MeasurementSecs: 15, Unit: "mg/dL", ReferenceRange: "", IsActive: true, Manufacturer: "ManPaSik"},
	{CategoryCode: 1, TypeIndex: 3, LegacyCode: 0x03, NameEN: "HbA1c", NameKO: "당화혈색소", RequiredChannels: 88, MeasurementSecs: 15, Unit: "%", ReferenceRange: "4.0-5.6", IsActive: true, Manufacturer: "ManPaSik"},
	{CategoryCode: 1, TypeIndex: 4, LegacyCode: 0x04, NameEN: "UricAcid", NameKO: "요산", RequiredChannels: 88, MeasurementSecs: 15, Unit: "mg/dL", ReferenceRange: "3.5-7.2", IsActive: true, Manufacturer: "ManPaSik"},
	{CategoryCode: 1, TypeIndex: 5, LegacyCode: 0x05, NameEN: "Creatinine", NameKO: "크레아티닌", RequiredChannels: 88, MeasurementSecs: 15, Unit: "mg/dL", ReferenceRange: "0.7-1.3", IsActive: true, Manufacturer: "ManPaSik"},
	{CategoryCode: 1, TypeIndex: 6, LegacyCode: 0x06, NameEN: "VitaminD", NameKO: "비타민 D", RequiredChannels: 88, MeasurementSecs: 15, Unit: "ng/mL", ReferenceRange: "30-100", IsActive: true, Manufacturer: "ManPaSik"},
	{CategoryCode: 1, TypeIndex: 7, LegacyCode: 0x07, NameEN: "VitaminB12", NameKO: "비타민 B12", RequiredChannels: 88, MeasurementSecs: 15, Unit: "pg/mL", ReferenceRange: "200-900", IsActive: true, Manufacturer: "ManPaSik"},
	{CategoryCode: 1, TypeIndex: 8, LegacyCode: 0x08, NameEN: "Ferritin", NameKO: "철분(페리틴)", RequiredChannels: 88, MeasurementSecs: 15, Unit: "ng/mL", ReferenceRange: "12-300", IsActive: true, Manufacturer: "ManPaSik"},
	{CategoryCode: 1, TypeIndex: 9, LegacyCode: 0x09, NameEN: "Tsh", NameKO: "갑상선(TSH)", RequiredChannels: 88, MeasurementSecs: 15, Unit: "mIU/L", ReferenceRange: "0.4-4.0", IsActive: true, Manufacturer: "ManPaSik"},
	{CategoryCode: 1, TypeIndex: 10, LegacyCode: 0x0A, NameEN: "Cortisol", NameKO: "코르티솔", RequiredChannels: 88, MeasurementSecs: 15, Unit: "μg/dL", ReferenceRange: "6-23", IsActive: true, Manufacturer: "ManPaSik"},
	{CategoryCode: 1, TypeIndex: 11, LegacyCode: 0x0B, NameEN: "Testosterone", NameKO: "테스토스테론", RequiredChannels: 88, MeasurementSecs: 15, Unit: "ng/dL", ReferenceRange: "300-1000", IsActive: true, Manufacturer: "ManPaSik"},
	{CategoryCode: 1, TypeIndex: 12, LegacyCode: 0x0C, NameEN: "Estrogen", NameKO: "에스트로겐", RequiredChannels: 88, MeasurementSecs: 15, Unit: "pg/mL", ReferenceRange: "15-350", IsActive: true, Manufacturer: "ManPaSik"},
	{CategoryCode: 1, TypeIndex: 13, LegacyCode: 0x0D, NameEN: "Crp", NameKO: "C-반응성단백", RequiredChannels: 88, MeasurementSecs: 15, Unit: "mg/L", ReferenceRange: "0-3", IsActive: true, Manufacturer: "ManPaSik"},
	{CategoryCode: 1, TypeIndex: 14, LegacyCode: 0x0E, NameEN: "Insulin", NameKO: "인슐린", RequiredChannels: 88, MeasurementSecs: 15, Unit: "μIU/mL", ReferenceRange: "2.6-24.9", IsActive: true, Manufacturer: "ManPaSik"},
	// Environmental 4종
	{CategoryCode: 2, TypeIndex: 1, LegacyCode: 0x20, NameEN: "WaterQuality", NameKO: "수질 검사", RequiredChannels: 88, MeasurementSecs: 15, Unit: "", ReferenceRange: "", IsActive: true, Manufacturer: "ManPaSik"},
	{CategoryCode: 2, TypeIndex: 2, LegacyCode: 0x21, NameEN: "IndoorAirQuality", NameKO: "실내 공기질", RequiredChannels: 88, MeasurementSecs: 15, Unit: "", ReferenceRange: "", IsActive: true, Manufacturer: "ManPaSik"},
	{CategoryCode: 2, TypeIndex: 3, LegacyCode: 0x22, NameEN: "Radon", NameKO: "라돈", RequiredChannels: 88, MeasurementSecs: 15, Unit: "", ReferenceRange: "", IsActive: true, Manufacturer: "ManPaSik"},
	{CategoryCode: 2, TypeIndex: 4, LegacyCode: 0x23, NameEN: "Radiation", NameKO: "방사능", RequiredChannels: 88, MeasurementSecs: 15, Unit: "", ReferenceRange: "", IsActive: true, Manufacturer: "ManPaSik"},
	// FoodSafety 4종
	{CategoryCode: 3, TypeIndex: 1, LegacyCode: 0x30, NameEN: "PesticideResidue", NameKO: "농약 잔류", RequiredChannels: 88, MeasurementSecs: 15, Unit: "", ReferenceRange: "", IsActive: true, Manufacturer: "ManPaSik"},
	{CategoryCode: 3, TypeIndex: 2, LegacyCode: 0x31, NameEN: "FoodFreshness", NameKO: "식품 신선도", RequiredChannels: 88, MeasurementSecs: 15, Unit: "", ReferenceRange: "", IsActive: true, Manufacturer: "ManPaSik"},
	{CategoryCode: 3, TypeIndex: 3, LegacyCode: 0x32, NameEN: "Allergen", NameKO: "알레르겐", RequiredChannels: 88, MeasurementSecs: 15, Unit: "", ReferenceRange: "", IsActive: true, Manufacturer: "ManPaSik"},
	{CategoryCode: 3, TypeIndex: 4, LegacyCode: 0x33, NameEN: "DateDrug", NameKO: "데이트약물", RequiredChannels: 88, MeasurementSecs: 15, Unit: "", ReferenceRange: "", IsActive: true, Manufacturer: "ManPaSik"},
	// ElectronicSensor 3종
	{CategoryCode: 4, TypeIndex: 1, LegacyCode: 0x40, NameEN: "ENose", NameKO: "전자코", RequiredChannels: 8, MeasurementSecs: 30, Unit: "", ReferenceRange: "", IsActive: true, Manufacturer: "ManPaSik"},
	{CategoryCode: 4, TypeIndex: 2, LegacyCode: 0x41, NameEN: "ETongue", NameKO: "전자혀", RequiredChannels: 8, MeasurementSecs: 30, Unit: "", ReferenceRange: "", IsActive: true, Manufacturer: "ManPaSik"},
	{CategoryCode: 4, TypeIndex: 3, LegacyCode: 0x42, NameEN: "EhdGas", NameKO: "EHD 기체", RequiredChannels: 8, MeasurementSecs: 30, Unit: "", ReferenceRange: "", IsActive: true, Manufacturer: "ManPaSik"},
	// AdvancedAnalysis 4종 (NonTarget1792 포함)
	{CategoryCode: 5, TypeIndex: 1, LegacyCode: 0x50, NameEN: "NonTarget448", NameKO: "비표적 448차원", RequiredChannels: 448, MeasurementSecs: 45, Unit: "", ReferenceRange: "", IsActive: true, Manufacturer: "ManPaSik"},
	{CategoryCode: 5, TypeIndex: 2, LegacyCode: 0x51, NameEN: "NonTarget896", NameKO: "비표적 896차원", RequiredChannels: 896, MeasurementSecs: 90, Unit: "", ReferenceRange: "", IsActive: true, Manufacturer: "ManPaSik"},
	{CategoryCode: 5, TypeIndex: 3, LegacyCode: 0x52, NameEN: "NonTarget1792", NameKO: "비표적 1792차원(궁극)", RequiredChannels: 1792, MeasurementSecs: 180, Unit: "", ReferenceRange: "", IsActive: true, Manufacturer: "ManPaSik"},
	{CategoryCode: 5, TypeIndex: 4, LegacyCode: 0x53, NameEN: "MultiBiomarker", NameKO: "다중 바이오마커", RequiredChannels: 88, MeasurementSecs: 90, Unit: "", ReferenceRange: "", IsActive: true, Manufacturer: "ManPaSik"},
	// CustomResearch 1종
	{CategoryCode: 255, TypeIndex: 1, LegacyCode: 0xFF, NameEN: "CustomResearch", NameKO: "맞춤형 연구용", RequiredChannels: 896, MeasurementSecs: 90, Unit: "", ReferenceRange: "", IsActive: true, Manufacturer: "ManPaSik"},
}

// registryKey는 카테고리+타입 인덱스로 레지스트리 키를 생성합니다.
func registryKey(categoryCode, typeIndex int32) string {
	return fmt.Sprintf("%d:%d", categoryCode, typeIndex)
}

// ============================================================================
// CartridgeService
// ============================================================================

// CartridgeService는 카트리지 비즈니스 로직입니다.
type CartridgeService struct {
	logger    *zap.Logger
	usageRepo CartridgeUsageRepository
	stateRepo CartridgeStateRepository

	// 내장 레지스트리 (30종)
	typeRegistry     map[string]*CartridgeTypeInfo // key: "category:typeIndex"
	categoryRegistry map[int32]*CategoryInfo
}

// NewCartridgeService는 새 CartridgeService를 생성합니다.
func NewCartridgeService(
	logger *zap.Logger,
	usageRepo CartridgeUsageRepository,
	stateRepo CartridgeStateRepository,
) *CartridgeService {
	svc := &CartridgeService{
		logger:           logger,
		usageRepo:        usageRepo,
		stateRepo:        stateRepo,
		typeRegistry:     make(map[string]*CartridgeTypeInfo),
		categoryRegistry: make(map[int32]*CategoryInfo),
	}

	// 카테고리 레지스트리 초기화
	for _, cat := range defaultCategories {
		cp := *cat
		svc.categoryRegistry[cat.Code] = &cp
	}

	// 타입 레지스트리 초기화 (30종)
	for _, t := range defaultTypes {
		cp := *t
		key := registryKey(t.CategoryCode, t.TypeIndex)
		svc.typeRegistry[key] = &cp
	}

	// 카테고리별 타입 수 계산
	for _, t := range svc.typeRegistry {
		if cat, ok := svc.categoryRegistry[t.CategoryCode]; ok {
			cat.TypeCount++
		}
	}

	return svc
}

// ============================================================================
// ReadCartridge — NFC 태그 데이터 파싱
// ============================================================================

// ReadCartridge는 NFC 태그 원시 데이터를 파싱하여 카트리지 정보를 반환합니다.
func (s *CartridgeService) ReadCartridge(_ context.Context, tagData []byte, tagVersion int32) (*CartridgeDetail, error) {
	if len(tagData) == 0 {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "NFC 태그 데이터가 비어 있습니다")
	}

	switch tagVersion {
	case 1:
		return s.parseTagV1(tagData)
	case 2:
		return s.parseTagV2(tagData)
	default:
		return nil, apperrors.New(apperrors.ErrInvalidInput, fmt.Sprintf("지원하지 않는 태그 버전: %d", tagVersion))
	}
}

// parseTagV1은 v1.0 NFC 태그를 파싱합니다.
// v1.0 레이아웃 (53+ 바이트):
//
//	[0-7]:   카트리지 UID (8바이트)
//	[8]:     레거시 타입 코드 (1-byte)
//	[9-16]:  로트 ID (8바이트)
//	[17-24]: 유효 기간 (YYYYMMDD)
//	[25-26]: 잔여 사용 횟수 (u16 LE)
//	[27-28]: 최대 사용 횟수 (u16 LE)
//	[29-36]: α 계수 (f64 LE)
//	[37-44]: 온도 보정 계수 (f64 LE)
//	[45-52]: 습도 보정 계수 (f64 LE)
func (s *CartridgeService) parseTagV1(data []byte) (*CartridgeDetail, error) {
	if len(data) < 53 {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "v1.0 태그: 최소 53바이트 필요")
	}

	uid := hex.EncodeToString(data[0:8])
	legacyCode := int32(data[8])
	categoryCode, typeIndex := legacyToFullCode(legacyCode)

	lotID := strings.TrimRight(string(data[9:17]), "\x00")
	expiryDate := string(data[17:25])

	remainingUses := int32(binary.LittleEndian.Uint16(data[25:27]))
	maxUses := int32(binary.LittleEndian.Uint16(data[27:29]))

	alpha := math.Float64frombits(binary.LittleEndian.Uint64(data[29:37]))
	tempCoeff := math.Float64frombits(binary.LittleEndian.Uint64(data[37:45]))
	humidityCoeff := math.Float64frombits(binary.LittleEndian.Uint64(data[45:53]))

	detail := &CartridgeDetail{
		CartridgeUID:        uid,
		CategoryCode:        categoryCode,
		TypeIndex:           typeIndex,
		LegacyCode:          legacyCode,
		LotID:               lotID,
		ExpiryDate:          expiryDate,
		RemainingUses:       remainingUses,
		MaxUses:             maxUses,
		AlphaCoefficient:    alpha,
		TempCoefficient:     tempCoeff,
		HumidityCoefficient: humidityCoeff,
		IsValid:             true,
	}

	// 레지스트리에서 타입 정보 보충
	s.enrichDetailFromRegistry(detail)

	return detail, nil
}

// parseTagV2는 v2.0 NFC 태그를 파싱합니다.
// v2.0 레이아웃 (80+ 바이트):
//
//	[0-7]:   카트리지 UID
//	[8]:     카테고리 코드
//	[9]:     타입 인덱스
//	[10]:    레거시 호환 코드 (0x00이면 없음)
//	[11]:    태그 포맷 버전 (0x02)
//	[12-19]: 로트 ID
//	[20-27]: 유효 기간
//	[28-29]: 잔여 사용 횟수
//	[30-31]: 최대 사용 횟수
//	[32-33]: 필요 채널 수 (u16 BE)
//	[34]:    측정 시간 (초)
//	[35]:    플래그
//	[36-43]: α 계수
//	[44-51]: 온도 보정 계수
//	[52-59]: 습도 보정 계수
//	[60-63]: CRC-32
//	[64+]:   확장 보정 데이터
func (s *CartridgeService) parseTagV2(data []byte) (*CartridgeDetail, error) {
	if len(data) < 80 {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "v2.0 태그: 최소 80바이트 필요")
	}

	uid := hex.EncodeToString(data[0:8])
	categoryCode := int32(data[8])
	typeIndex := int32(data[9])
	legacyCode := int32(data[10])

	lotID := strings.TrimRight(string(data[12:20]), "\x00")
	expiryDate := string(data[20:28])

	remainingUses := int32(binary.LittleEndian.Uint16(data[28:30]))
	maxUses := int32(binary.LittleEndian.Uint16(data[30:32]))

	requiredChannels := int32(binary.BigEndian.Uint16(data[32:34]))
	measurementSecs := int32(data[34])

	alpha := math.Float64frombits(binary.LittleEndian.Uint64(data[36:44]))
	tempCoeff := math.Float64frombits(binary.LittleEndian.Uint64(data[44:52]))
	humidityCoeff := math.Float64frombits(binary.LittleEndian.Uint64(data[52:60]))

	detail := &CartridgeDetail{
		CartridgeUID:        uid,
		CategoryCode:        categoryCode,
		TypeIndex:           typeIndex,
		LegacyCode:          legacyCode,
		LotID:               lotID,
		ExpiryDate:          expiryDate,
		RemainingUses:       remainingUses,
		MaxUses:             maxUses,
		AlphaCoefficient:    alpha,
		TempCoefficient:     tempCoeff,
		HumidityCoefficient: humidityCoeff,
		RequiredChannels:    requiredChannels,
		MeasurementSecs:     measurementSecs,
		IsValid:             true,
	}

	// 레지스트리에서 타입 정보 보충
	s.enrichDetailFromRegistry(detail)

	return detail, nil
}

// enrichDetailFromRegistry는 레지스트리에서 타입 정보를 보충합니다.
func (s *CartridgeService) enrichDetailFromRegistry(detail *CartridgeDetail) {
	key := registryKey(detail.CategoryCode, detail.TypeIndex)
	if typeInfo, ok := s.typeRegistry[key]; ok {
		detail.NameKO = typeInfo.NameKO
		detail.NameEN = typeInfo.NameEN
		detail.Unit = typeInfo.Unit
		detail.ReferenceRange = typeInfo.ReferenceRange
		// v1에서는 태그에 채널/측정시간이 없으므로 레지스트리에서 보충
		if detail.RequiredChannels == 0 {
			detail.RequiredChannels = typeInfo.RequiredChannels
		}
		if detail.MeasurementSecs == 0 {
			detail.MeasurementSecs = typeInfo.MeasurementSecs
		}
	}
}

// legacyToFullCode는 레거시 1-byte 코드를 (category, typeIndex)로 변환합니다.
func legacyToFullCode(legacy int32) (categoryCode, typeIndex int32) {
	switch {
	case legacy >= 0x01 && legacy <= 0x0E:
		return 1, legacy // HealthBiomarker
	case legacy >= 0x20 && legacy <= 0x23:
		return 2, legacy - 0x1F // Environmental
	case legacy >= 0x30 && legacy <= 0x33:
		return 3, legacy - 0x2F // FoodSafety
	case legacy >= 0x40 && legacy <= 0x42:
		return 4, legacy - 0x3F // ElectronicSensor
	case legacy >= 0x50 && legacy <= 0x53:
		return 5, legacy - 0x4F // AdvancedAnalysis
	case legacy == 0xFF:
		return 255, 1 // CustomResearch
	default:
		return 0, 0
	}
}

// ============================================================================
// RecordUsage — 카트리지 사용 기록
// ============================================================================

// RecordUsage는 카트리지 사용을 기록하고 잔여 횟수를 감소시킵니다.
func (s *CartridgeService) RecordUsage(ctx context.Context, userID, sessionID, cartridgeUID string, categoryCode, typeIndex int32) (int32, error) {
	if userID == "" || cartridgeUID == "" {
		return 0, apperrors.New(apperrors.ErrInvalidInput, "user_id와 cartridge_uid는 필수입니다")
	}

	// 잔여 횟수 감소
	remaining, err := s.stateRepo.DecrementUses(ctx, cartridgeUID)
	if err != nil {
		s.logger.Error("카트리지 사용 횟수 감소 실패", zap.Error(err))
		return 0, apperrors.New(apperrors.ErrInternal, "사용 횟수 감소에 실패했습니다")
	}

	// 타입 이름 조회
	typeNameKO := ""
	key := registryKey(categoryCode, typeIndex)
	if typeInfo, ok := s.typeRegistry[key]; ok {
		typeNameKO = typeInfo.NameKO
	}

	// 사용 기록 저장
	record := &CartridgeUsageRecord{
		RecordID:     uuid.New().String(),
		UserID:       userID,
		SessionID:    sessionID,
		CartridgeUID: cartridgeUID,
		CategoryCode: categoryCode,
		TypeIndex:    typeIndex,
		TypeNameKO:   typeNameKO,
		UsedAt:       time.Now().UTC(),
	}

	if err := s.usageRepo.Create(ctx, record); err != nil {
		s.logger.Error("카트리지 사용 기록 저장 실패", zap.Error(err))
		return remaining, apperrors.New(apperrors.ErrInternal, "사용 기록 저장에 실패했습니다")
	}

	s.logger.Info("카트리지 사용 기록",
		zap.String("user_id", userID),
		zap.String("cartridge_uid", cartridgeUID),
		zap.Int32("remaining", remaining),
	)

	return remaining, nil
}

// ============================================================================
// GetUsageHistory — 사용 이력 조회
// ============================================================================

// GetUsageHistory는 사용자의 카트리지 사용 이력을 조회합니다.
func (s *CartridgeService) GetUsageHistory(ctx context.Context, userID string, limit, offset int32) ([]*CartridgeUsageRecord, int32, error) {
	if userID == "" {
		return nil, 0, apperrors.New(apperrors.ErrInvalidInput, "user_id는 필수입니다")
	}

	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	records, totalCount, err := s.usageRepo.ListByUserID(ctx, userID, limit, offset)
	if err != nil {
		s.logger.Error("사용 이력 조회 실패", zap.Error(err))
		return nil, 0, apperrors.New(apperrors.ErrInternal, "사용 이력 조회에 실패했습니다")
	}

	return records, totalCount, nil
}

// ============================================================================
// GetCartridgeType — 타입 정보 조회
// ============================================================================

// GetCartridgeType은 카트리지 타입 정보를 레지스트리에서 조회합니다.
func (s *CartridgeService) GetCartridgeType(_ context.Context, categoryCode, typeIndex int32) (*CartridgeTypeInfo, error) {
	key := registryKey(categoryCode, typeIndex)
	typeInfo, ok := s.typeRegistry[key]
	if !ok {
		return nil, apperrors.New(apperrors.ErrNotFound, fmt.Sprintf("카트리지 타입을 찾을 수 없습니다: category=%d, type=%d", categoryCode, typeIndex))
	}

	cp := *typeInfo
	return &cp, nil
}

// ============================================================================
// ListCategories — 카테고리 목록 조회
// ============================================================================

// ListCategories는 모든 카테고리 목록을 반환합니다.
func (s *CartridgeService) ListCategories() []*CategoryInfo {
	result := make([]*CategoryInfo, 0, len(s.categoryRegistry))
	// 정렬된 순서로 반환
	orderedCodes := []int32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 254, 255}
	for _, code := range orderedCodes {
		if cat, ok := s.categoryRegistry[code]; ok {
			cp := *cat
			result = append(result, &cp)
		}
	}
	return result
}

// ============================================================================
// ListTypesByCategory — 카테고리별 타입 목록
// ============================================================================

// ListTypesByCategory는 지정 카테고리에 속하는 타입 목록을 반환합니다.
func (s *CartridgeService) ListTypesByCategory(_ context.Context, categoryCode int32) ([]*CartridgeTypeInfo, error) {
	var result []*CartridgeTypeInfo
	for _, t := range s.typeRegistry {
		if t.CategoryCode == categoryCode {
			cp := *t
			result = append(result, &cp)
		}
	}
	if len(result) == 0 {
		return nil, apperrors.New(apperrors.ErrNotFound, fmt.Sprintf("카테고리 %d에 등록된 타입이 없습니다", categoryCode))
	}
	return result, nil
}

// ============================================================================
// GetRemainingUses — 잔여 사용 횟수 조회
// ============================================================================

// GetRemainingUses는 카트리지의 잔여 사용 정보를 조회합니다.
func (s *CartridgeService) GetRemainingUses(ctx context.Context, cartridgeUID string) (*CartridgeRemainingInfo, error) {
	if cartridgeUID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "cartridge_uid는 필수입니다")
	}

	info, err := s.stateRepo.GetByUID(ctx, cartridgeUID)
	if err != nil {
		s.logger.Error("잔여 사용 정보 조회 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "잔여 사용 정보 조회에 실패했습니다")
	}
	if info == nil {
		return nil, apperrors.New(apperrors.ErrNotFound, "카트리지 상태 정보가 없습니다")
	}

	// 만료 확인
	today := time.Now().UTC().Format("20060102")
	info.IsExpired = info.ExpiryDate != "" && info.ExpiryDate < today

	return info, nil
}

// ============================================================================
// ValidateCartridge — 카트리지 유효성 검증
// ============================================================================

// ValidateCartridge는 카트리지의 유효성을 검증합니다.
func (s *CartridgeService) ValidateCartridge(ctx context.Context, cartridgeUID string, categoryCode, typeIndex int32, _ string) (*ValidateResult, error) {
	if cartridgeUID == "" {
		return nil, apperrors.New(apperrors.ErrInvalidInput, "cartridge_uid는 필수입니다")
	}

	result := &ValidateResult{
		IsValid:     true,
		Reason:      "ok",
		AccessLevel: AccessIncluded,
	}

	// 타입 존재 여부 확인
	key := registryKey(categoryCode, typeIndex)
	_, typeExists := s.typeRegistry[key]
	if !typeExists {
		result.IsValid = false
		result.Reason = "unknown_type"
		result.AccessLevel = AccessRestricted
		return result, nil
	}

	// 카트리지 상태 조회
	state, err := s.stateRepo.GetByUID(ctx, cartridgeUID)
	if err != nil {
		s.logger.Error("카트리지 상태 조회 실패", zap.Error(err))
		return nil, apperrors.New(apperrors.ErrInternal, "카트리지 상태 조회에 실패했습니다")
	}

	if state == nil {
		// 상태 없으면 유효 (아직 사용되지 않은 카트리지)
		return result, nil
	}

	result.RemainingUses = state.RemainingUses

	// 만료 확인
	if state.ExpiryDate != "" {
		today := time.Now().UTC().Format("20060102")
		if state.ExpiryDate < today {
			result.IsValid = false
			result.Reason = "expired"
			return result, nil
		}
	}

	// 잔여 사용 횟수 확인
	if state.RemainingUses <= 0 {
		result.IsValid = false
		result.Reason = "no_uses"
		return result, nil
	}

	return result, nil
}

// ============================================================================
// InitCartridgeState — 카트리지 상태 초기화 (ReadCartridge 후 호출)
// ============================================================================

// InitCartridgeState는 ReadCartridge로 읽은 카트리지의 상태를 저장소에 기록합니다.
func (s *CartridgeService) InitCartridgeState(ctx context.Context, detail *CartridgeDetail) error {
	if detail == nil {
		return nil
	}

	info := &CartridgeRemainingInfo{
		CartridgeUID:  detail.CartridgeUID,
		RemainingUses: detail.RemainingUses,
		MaxUses:       detail.MaxUses,
		ExpiryDate:    detail.ExpiryDate,
	}

	return s.stateRepo.Upsert(ctx, info)
}
