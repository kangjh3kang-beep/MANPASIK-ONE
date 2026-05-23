# Mapping Registry

패킷 ↔ OMOP CDM ↔ FHIR R4 ↔ LOINC/SNOMED 매핑 (버전관리)

## 파일 목록
- `loinc_mapping.json` — 바이오마커 ↔ LOINC/UCUM 매핑 (15개 항목)

## 매핑 원칙
- 모든 측정 단위는 UCUM 코드 사용
- FHIR R4 Observation 카테고리 (laboratory / vital-signs) 구분
- 매핑 변경 시 버전 업데이트 필수
- mg/dL ↔ mmol/L 등 단위 변환은 이 레지스트리 기준으로 수행

## 향후 확장
- OMOP CDM concept_id 매핑 추가 예정
- SNOMED CT 코드 매핑 추가 예정
