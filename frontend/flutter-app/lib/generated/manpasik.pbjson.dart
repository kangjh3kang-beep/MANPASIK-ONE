// This is a generated file - do not edit.
//
// Generated from manpasik.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports
// ignore_for_file: unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use socialProviderDescriptor instead')
const SocialProvider$json = {
  '1': 'SocialProvider',
  '2': [
    {'1': 'SOCIAL_PROVIDER_UNSPECIFIED', '2': 0},
    {'1': 'SOCIAL_PROVIDER_GOOGLE', '2': 1},
    {'1': 'SOCIAL_PROVIDER_APPLE', '2': 2},
    {'1': 'SOCIAL_PROVIDER_KAKAO', '2': 3},
    {'1': 'SOCIAL_PROVIDER_NAVER', '2': 4},
  ],
};

/// Descriptor for `SocialProvider`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List socialProviderDescriptor = $convert.base64Decode(
    'Cg5Tb2NpYWxQcm92aWRlchIfChtTT0NJQUxfUFJPVklERVJfVU5TUEVDSUZJRUQQABIaChZTT0'
    'NJQUxfUFJPVklERVJfR09PR0xFEAESGQoVU09DSUFMX1BST1ZJREVSX0FQUExFEAISGQoVU09D'
    'SUFMX1BST1ZJREVSX0tBS0FPEAMSGQoVU09DSUFMX1BST1ZJREVSX05BVkVSEAQ=');

@$core.Deprecated('Use genderDescriptor instead')
const Gender$json = {
  '1': 'Gender',
  '2': [
    {'1': 'GENDER_UNSPECIFIED', '2': 0},
    {'1': 'GENDER_MALE', '2': 1},
    {'1': 'GENDER_FEMALE', '2': 2},
    {'1': 'GENDER_OTHER', '2': 3},
    {'1': 'GENDER_PREFER_NOT_TO_SAY', '2': 4},
  ],
};

/// Descriptor for `Gender`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List genderDescriptor = $convert.base64Decode(
    'CgZHZW5kZXISFgoSR0VOREVSX1VOU1BFQ0lGSUVEEAASDwoLR0VOREVSX01BTEUQARIRCg1HRU'
    '5ERVJfRkVNQUxFEAISEAoMR0VOREVSX09USEVSEAMSHAoYR0VOREVSX1BSRUZFUl9OT1RfVE9f'
    'U0FZEAQ=');

@$core.Deprecated('Use deviceStatusDescriptor instead')
const DeviceStatus$json = {
  '1': 'DeviceStatus',
  '2': [
    {'1': 'DEVICE_STATUS_UNKNOWN', '2': 0},
    {'1': 'DEVICE_STATUS_ONLINE', '2': 1},
    {'1': 'DEVICE_STATUS_OFFLINE', '2': 2},
    {'1': 'DEVICE_STATUS_MEASURING', '2': 3},
    {'1': 'DEVICE_STATUS_UPDATING', '2': 4},
    {'1': 'DEVICE_STATUS_ERROR', '2': 5},
  ],
};

/// Descriptor for `DeviceStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List deviceStatusDescriptor = $convert.base64Decode(
    'CgxEZXZpY2VTdGF0dXMSGQoVREVWSUNFX1NUQVRVU19VTktOT1dOEAASGAoUREVWSUNFX1NUQV'
    'RVU19PTkxJTkUQARIZChVERVZJQ0VfU1RBVFVTX09GRkxJTkUQAhIbChdERVZJQ0VfU1RBVFVT'
    'X01FQVNVUklORxADEhoKFkRFVklDRV9TVEFUVVNfVVBEQVRJTkcQBBIXChNERVZJQ0VfU1RBVF'
    'VTX0VSUk9SEAU=');

@$core.Deprecated('Use commandTypeDescriptor instead')
const CommandType$json = {
  '1': 'CommandType',
  '2': [
    {'1': 'COMMAND_TYPE_UNKNOWN', '2': 0},
    {'1': 'COMMAND_TYPE_START_MEASUREMENT', '2': 1},
    {'1': 'COMMAND_TYPE_STOP_MEASUREMENT', '2': 2},
    {'1': 'COMMAND_TYPE_CALIBRATE', '2': 3},
    {'1': 'COMMAND_TYPE_REBOOT', '2': 4},
    {'1': 'COMMAND_TYPE_OTA_UPDATE', '2': 5},
  ],
};

/// Descriptor for `CommandType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List commandTypeDescriptor = $convert.base64Decode(
    'CgtDb21tYW5kVHlwZRIYChRDT01NQU5EX1RZUEVfVU5LTk9XThAAEiIKHkNPTU1BTkRfVFlQRV'
    '9TVEFSVF9NRUFTVVJFTUVOVBABEiEKHUNPTU1BTkRfVFlQRV9TVE9QX01FQVNVUkVNRU5UEAIS'
    'GgoWQ09NTUFORF9UWVBFX0NBTElCUkFURRADEhcKE0NPTU1BTkRfVFlQRV9SRUJPT1QQBBIbCh'
    'dDT01NQU5EX1RZUEVfT1RBX1VQREFURRAF');

@$core.Deprecated('Use subscriptionTierDescriptor instead')
const SubscriptionTier$json = {
  '1': 'SubscriptionTier',
  '2': [
    {'1': 'SUBSCRIPTION_TIER_FREE', '2': 0},
    {'1': 'SUBSCRIPTION_TIER_BASIC', '2': 1},
    {'1': 'SUBSCRIPTION_TIER_PRO', '2': 2},
    {'1': 'SUBSCRIPTION_TIER_CLINICAL', '2': 3},
  ],
};

/// Descriptor for `SubscriptionTier`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List subscriptionTierDescriptor = $convert.base64Decode(
    'ChBTdWJzY3JpcHRpb25UaWVyEhoKFlNVQlNDUklQVElPTl9USUVSX0ZSRUUQABIbChdTVUJTQ1'
    'JJUFRJT05fVElFUl9CQVNJQxABEhkKFVNVQlNDUklQVElPTl9USUVSX1BSTxACEh4KGlNVQlND'
    'UklQVElPTl9USUVSX0NMSU5JQ0FMEAM=');

@$core.Deprecated('Use subscriptionStatusDescriptor instead')
const SubscriptionStatus$json = {
  '1': 'SubscriptionStatus',
  '2': [
    {'1': 'SUBSCRIPTION_STATUS_UNKNOWN', '2': 0},
    {'1': 'SUBSCRIPTION_STATUS_ACTIVE', '2': 1},
    {'1': 'SUBSCRIPTION_STATUS_CANCELLED', '2': 2},
    {'1': 'SUBSCRIPTION_STATUS_EXPIRED', '2': 3},
    {'1': 'SUBSCRIPTION_STATUS_SUSPENDED', '2': 4},
    {'1': 'SUBSCRIPTION_STATUS_TRIAL', '2': 5},
  ],
};

/// Descriptor for `SubscriptionStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List subscriptionStatusDescriptor = $convert.base64Decode(
    'ChJTdWJzY3JpcHRpb25TdGF0dXMSHwobU1VCU0NSSVBUSU9OX1NUQVRVU19VTktOT1dOEAASHg'
    'oaU1VCU0NSSVBUSU9OX1NUQVRVU19BQ1RJVkUQARIhCh1TVUJTQ1JJUFRJT05fU1RBVFVTX0NB'
    'TkNFTExFRBACEh8KG1NVQlNDUklQVElPTl9TVEFUVVNfRVhQSVJFRBADEiEKHVNVQlNDUklQVE'
    'lPTl9TVEFUVVNfU1VTUEVOREVEEAQSHQoZU1VCU0NSSVBUSU9OX1NUQVRVU19UUklBTBAF');

@$core.Deprecated('Use productCategoryDescriptor instead')
const ProductCategory$json = {
  '1': 'ProductCategory',
  '2': [
    {'1': 'PRODUCT_CATEGORY_UNKNOWN', '2': 0},
    {'1': 'PRODUCT_CATEGORY_CARTRIDGE', '2': 1},
    {'1': 'PRODUCT_CATEGORY_READER', '2': 2},
    {'1': 'PRODUCT_CATEGORY_ACCESSORY', '2': 3},
    {'1': 'PRODUCT_CATEGORY_BUNDLE', '2': 4},
  ],
};

/// Descriptor for `ProductCategory`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List productCategoryDescriptor = $convert.base64Decode(
    'Cg9Qcm9kdWN0Q2F0ZWdvcnkSHAoYUFJPRFVDVF9DQVRFR09SWV9VTktOT1dOEAASHgoaUFJPRF'
    'VDVF9DQVRFR09SWV9DQVJUUklER0UQARIbChdQUk9EVUNUX0NBVEVHT1JZX1JFQURFUhACEh4K'
    'GlBST0RVQ1RfQ0FURUdPUllfQUNDRVNTT1JZEAMSGwoXUFJPRFVDVF9DQVRFR09SWV9CVU5ETE'
    'UQBA==');

@$core.Deprecated('Use orderStatusDescriptor instead')
const OrderStatus$json = {
  '1': 'OrderStatus',
  '2': [
    {'1': 'ORDER_STATUS_UNKNOWN', '2': 0},
    {'1': 'ORDER_STATUS_PENDING', '2': 1},
    {'1': 'ORDER_STATUS_PAID', '2': 2},
    {'1': 'ORDER_STATUS_SHIPPED', '2': 3},
    {'1': 'ORDER_STATUS_DELIVERED', '2': 4},
    {'1': 'ORDER_STATUS_CANCELLED', '2': 5},
    {'1': 'ORDER_STATUS_REFUNDED', '2': 6},
  ],
};

/// Descriptor for `OrderStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List orderStatusDescriptor = $convert.base64Decode(
    'CgtPcmRlclN0YXR1cxIYChRPUkRFUl9TVEFUVVNfVU5LTk9XThAAEhgKFE9SREVSX1NUQVRVU1'
    '9QRU5ESU5HEAESFQoRT1JERVJfU1RBVFVTX1BBSUQQAhIYChRPUkRFUl9TVEFUVVNfU0hJUFBF'
    'RBADEhoKFk9SREVSX1NUQVRVU19ERUxJVkVSRUQQBBIaChZPUkRFUl9TVEFUVVNfQ0FOQ0VMTE'
    'VEEAUSGQoVT1JERVJfU1RBVFVTX1JFRlVOREVEEAY=');

@$core.Deprecated('Use paymentTypeDescriptor instead')
const PaymentType$json = {
  '1': 'PaymentType',
  '2': [
    {'1': 'PAYMENT_TYPE_UNKNOWN', '2': 0},
    {'1': 'PAYMENT_TYPE_ONE_TIME', '2': 1},
    {'1': 'PAYMENT_TYPE_SUBSCRIPTION', '2': 2},
  ],
};

/// Descriptor for `PaymentType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List paymentTypeDescriptor = $convert.base64Decode(
    'CgtQYXltZW50VHlwZRIYChRQQVlNRU5UX1RZUEVfVU5LTk9XThAAEhkKFVBBWU1FTlRfVFlQRV'
    '9PTkVfVElNRRABEh0KGVBBWU1FTlRfVFlQRV9TVUJTQ1JJUFRJT04QAg==');

@$core.Deprecated('Use paymentStatusDescriptor instead')
const PaymentStatus$json = {
  '1': 'PaymentStatus',
  '2': [
    {'1': 'PAYMENT_STATUS_UNKNOWN', '2': 0},
    {'1': 'PAYMENT_STATUS_PENDING', '2': 1},
    {'1': 'PAYMENT_STATUS_COMPLETED', '2': 2},
    {'1': 'PAYMENT_STATUS_FAILED', '2': 3},
    {'1': 'PAYMENT_STATUS_CANCELLED', '2': 4},
    {'1': 'PAYMENT_STATUS_REFUNDED', '2': 5},
    {'1': 'PAYMENT_STATUS_PARTIAL_REFUND', '2': 6},
  ],
};

/// Descriptor for `PaymentStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List paymentStatusDescriptor = $convert.base64Decode(
    'Cg1QYXltZW50U3RhdHVzEhoKFlBBWU1FTlRfU1RBVFVTX1VOS05PV04QABIaChZQQVlNRU5UX1'
    'NUQVRVU19QRU5ESU5HEAESHAoYUEFZTUVOVF9TVEFUVVNfQ09NUExFVEVEEAISGQoVUEFZTUVO'
    'VF9TVEFUVVNfRkFJTEVEEAMSHAoYUEFZTUVOVF9TVEFUVVNfQ0FOQ0VMTEVEEAQSGwoXUEFZTU'
    'VOVF9TVEFUVVNfUkVGVU5ERUQQBRIhCh1QQVlNRU5UX1NUQVRVU19QQVJUSUFMX1JFRlVORBAG');

@$core.Deprecated('Use aiModelTypeDescriptor instead')
const AiModelType$json = {
  '1': 'AiModelType',
  '2': [
    {'1': 'AI_MODEL_TYPE_UNSPECIFIED', '2': 0},
    {'1': 'AI_MODEL_TYPE_BIOMARKER_CLASSIFIER', '2': 1},
    {'1': 'AI_MODEL_TYPE_ANOMALY_DETECTOR', '2': 2},
    {'1': 'AI_MODEL_TYPE_TREND_PREDICTOR', '2': 3},
    {'1': 'AI_MODEL_TYPE_HEALTH_SCORER', '2': 4},
    {'1': 'AI_MODEL_TYPE_FOOD_CALORIE_ESTIMATOR', '2': 5},
  ],
};

/// Descriptor for `AiModelType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List aiModelTypeDescriptor = $convert.base64Decode(
    'CgtBaU1vZGVsVHlwZRIdChlBSV9NT0RFTF9UWVBFX1VOU1BFQ0lGSUVEEAASJgoiQUlfTU9ERU'
    'xfVFlQRV9CSU9NQVJLRVJfQ0xBU1NJRklFUhABEiIKHkFJX01PREVMX1RZUEVfQU5PTUFMWV9E'
    'RVRFQ1RPUhACEiEKHUFJX01PREVMX1RZUEVfVFJFTkRfUFJFRElDVE9SEAMSHwobQUlfTU9ERU'
    'xfVFlQRV9IRUFMVEhfU0NPUkVSEAQSKAokQUlfTU9ERUxfVFlQRV9GT09EX0NBTE9SSUVfRVNU'
    'SU1BVE9SEAU=');

@$core.Deprecated('Use riskLevelDescriptor instead')
const RiskLevel$json = {
  '1': 'RiskLevel',
  '2': [
    {'1': 'RISK_LEVEL_UNSPECIFIED', '2': 0},
    {'1': 'RISK_LEVEL_LOW', '2': 1},
    {'1': 'RISK_LEVEL_MODERATE', '2': 2},
    {'1': 'RISK_LEVEL_HIGH', '2': 3},
    {'1': 'RISK_LEVEL_CRITICAL', '2': 4},
  ],
};

/// Descriptor for `RiskLevel`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List riskLevelDescriptor = $convert.base64Decode(
    'CglSaXNrTGV2ZWwSGgoWUklTS19MRVZFTF9VTlNQRUNJRklFRBAAEhIKDlJJU0tfTEVWRUxfTE'
    '9XEAESFwoTUklTS19MRVZFTF9NT0RFUkFURRACEhMKD1JJU0tfTEVWRUxfSElHSBADEhcKE1JJ'
    'U0tfTEVWRUxfQ1JJVElDQUwQBA==');

@$core.Deprecated('Use calibrationTypeDescriptor instead')
const CalibrationType$json = {
  '1': 'CalibrationType',
  '2': [
    {'1': 'CALIBRATION_TYPE_UNKNOWN', '2': 0},
    {'1': 'CALIBRATION_TYPE_FACTORY', '2': 1},
    {'1': 'CALIBRATION_TYPE_FIELD', '2': 2},
    {'1': 'CALIBRATION_TYPE_AUTO', '2': 3},
  ],
};

/// Descriptor for `CalibrationType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List calibrationTypeDescriptor = $convert.base64Decode(
    'Cg9DYWxpYnJhdGlvblR5cGUSHAoYQ0FMSUJSQVRJT05fVFlQRV9VTktOT1dOEAASHAoYQ0FMSU'
    'JSQVRJT05fVFlQRV9GQUNUT1JZEAESGgoWQ0FMSUJSQVRJT05fVFlQRV9GSUVMRBACEhkKFUNB'
    'TElCUkFUSU9OX1RZUEVfQVVUTxAD');

@$core.Deprecated('Use calibrationStatusDescriptor instead')
const CalibrationStatus$json = {
  '1': 'CalibrationStatus',
  '2': [
    {'1': 'CALIBRATION_STATUS_UNKNOWN', '2': 0},
    {'1': 'CALIBRATION_STATUS_VALID', '2': 1},
    {'1': 'CALIBRATION_STATUS_EXPIRING', '2': 2},
    {'1': 'CALIBRATION_STATUS_EXPIRED', '2': 3},
    {'1': 'CALIBRATION_STATUS_NEEDED', '2': 4},
  ],
};

/// Descriptor for `CalibrationStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List calibrationStatusDescriptor = $convert.base64Decode(
    'ChFDYWxpYnJhdGlvblN0YXR1cxIeChpDQUxJQlJBVElPTl9TVEFUVVNfVU5LTk9XThAAEhwKGE'
    'NBTElCUkFUSU9OX1NUQVRVU19WQUxJRBABEh8KG0NBTElCUkFUSU9OX1NUQVRVU19FWFBJUklO'
    'RxACEh4KGkNBTElCUkFUSU9OX1NUQVRVU19FWFBJUkVEEAMSHQoZQ0FMSUJSQVRJT05fU1RBVF'
    'VTX05FRURFRBAE');

@$core.Deprecated('Use goalCategoryDescriptor instead')
const GoalCategory$json = {
  '1': 'GoalCategory',
  '2': [
    {'1': 'GOAL_CATEGORY_UNKNOWN', '2': 0},
    {'1': 'GOAL_CATEGORY_BLOOD_GLUCOSE', '2': 1},
    {'1': 'GOAL_CATEGORY_BLOOD_PRESSURE', '2': 2},
    {'1': 'GOAL_CATEGORY_CHOLESTEROL', '2': 3},
    {'1': 'GOAL_CATEGORY_WEIGHT', '2': 4},
    {'1': 'GOAL_CATEGORY_EXERCISE', '2': 5},
    {'1': 'GOAL_CATEGORY_NUTRITION', '2': 6},
    {'1': 'GOAL_CATEGORY_SLEEP', '2': 7},
    {'1': 'GOAL_CATEGORY_STRESS', '2': 8},
    {'1': 'GOAL_CATEGORY_CUSTOM', '2': 9},
  ],
};

/// Descriptor for `GoalCategory`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List goalCategoryDescriptor = $convert.base64Decode(
    'CgxHb2FsQ2F0ZWdvcnkSGQoVR09BTF9DQVRFR09SWV9VTktOT1dOEAASHwobR09BTF9DQVRFR0'
    '9SWV9CTE9PRF9HTFVDT1NFEAESIAocR09BTF9DQVRFR09SWV9CTE9PRF9QUkVTU1VSRRACEh0K'
    'GUdPQUxfQ0FURUdPUllfQ0hPTEVTVEVST0wQAxIYChRHT0FMX0NBVEVHT1JZX1dFSUdIVBAEEh'
    'oKFkdPQUxfQ0FURUdPUllfRVhFUkNJU0UQBRIbChdHT0FMX0NBVEVHT1JZX05VVFJJVElPThAG'
    'EhcKE0dPQUxfQ0FURUdPUllfU0xFRVAQBxIYChRHT0FMX0NBVEVHT1JZX1NUUkVTUxAIEhgKFE'
    'dPQUxfQ0FURUdPUllfQ1VTVE9NEAk=');

@$core.Deprecated('Use goalStatusDescriptor instead')
const GoalStatus$json = {
  '1': 'GoalStatus',
  '2': [
    {'1': 'GOAL_STATUS_UNKNOWN', '2': 0},
    {'1': 'GOAL_STATUS_ACTIVE', '2': 1},
    {'1': 'GOAL_STATUS_ACHIEVED', '2': 2},
    {'1': 'GOAL_STATUS_PAUSED', '2': 3},
    {'1': 'GOAL_STATUS_CANCELLED', '2': 4},
  ],
};

/// Descriptor for `GoalStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List goalStatusDescriptor = $convert.base64Decode(
    'CgpHb2FsU3RhdHVzEhcKE0dPQUxfU1RBVFVTX1VOS05PV04QABIWChJHT0FMX1NUQVRVU19BQ1'
    'RJVkUQARIYChRHT0FMX1NUQVRVU19BQ0hJRVZFRBACEhYKEkdPQUxfU1RBVFVTX1BBVVNFRBAD'
    'EhkKFUdPQUxfU1RBVFVTX0NBTkNFTExFRBAE');

@$core.Deprecated('Use coachingTypeDescriptor instead')
const CoachingType$json = {
  '1': 'CoachingType',
  '2': [
    {'1': 'COACHING_TYPE_UNKNOWN', '2': 0},
    {'1': 'COACHING_TYPE_MEASUREMENT_FEEDBACK', '2': 1},
    {'1': 'COACHING_TYPE_DAILY_TIP', '2': 2},
    {'1': 'COACHING_TYPE_GOAL_PROGRESS', '2': 3},
    {'1': 'COACHING_TYPE_ALERT', '2': 4},
    {'1': 'COACHING_TYPE_MOTIVATION', '2': 5},
    {'1': 'COACHING_TYPE_RECOMMENDATION', '2': 6},
  ],
};

/// Descriptor for `CoachingType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List coachingTypeDescriptor = $convert.base64Decode(
    'CgxDb2FjaGluZ1R5cGUSGQoVQ09BQ0hJTkdfVFlQRV9VTktOT1dOEAASJgoiQ09BQ0hJTkdfVF'
    'lQRV9NRUFTVVJFTUVOVF9GRUVEQkFDSxABEhsKF0NPQUNISU5HX1RZUEVfREFJTFlfVElQEAIS'
    'HwobQ09BQ0hJTkdfVFlQRV9HT0FMX1BST0dSRVNTEAMSFwoTQ09BQ0hJTkdfVFlQRV9BTEVSVB'
    'AEEhwKGENPQUNISU5HX1RZUEVfTU9USVZBVElPThAFEiAKHENPQUNISU5HX1RZUEVfUkVDT01N'
    'RU5EQVRJT04QBg==');

@$core.Deprecated('Use recommendationTypeDescriptor instead')
const RecommendationType$json = {
  '1': 'RecommendationType',
  '2': [
    {'1': 'RECOMMENDATION_TYPE_UNKNOWN', '2': 0},
    {'1': 'RECOMMENDATION_TYPE_FOOD', '2': 1},
    {'1': 'RECOMMENDATION_TYPE_EXERCISE', '2': 2},
    {'1': 'RECOMMENDATION_TYPE_SUPPLEMENT', '2': 3},
    {'1': 'RECOMMENDATION_TYPE_LIFESTYLE', '2': 4},
    {'1': 'RECOMMENDATION_TYPE_CHECKUP', '2': 5},
  ],
};

/// Descriptor for `RecommendationType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List recommendationTypeDescriptor = $convert.base64Decode(
    'ChJSZWNvbW1lbmRhdGlvblR5cGUSHwobUkVDT01NRU5EQVRJT05fVFlQRV9VTktOT1dOEAASHA'
    'oYUkVDT01NRU5EQVRJT05fVFlQRV9GT09EEAESIAocUkVDT01NRU5EQVRJT05fVFlQRV9FWEVS'
    'Q0lTRRACEiIKHlJFQ09NTUVOREFUSU9OX1RZUEVfU1VQUExFTUVOVBADEiEKHVJFQ09NTUVORE'
    'FUSU9OX1RZUEVfTElGRVNUWUxFEAQSHwobUkVDT01NRU5EQVRJT05fVFlQRV9DSEVDS1VQEAU=');

@$core.Deprecated('Use cartridgeAccessLevelDescriptor instead')
const CartridgeAccessLevel$json = {
  '1': 'CartridgeAccessLevel',
  '2': [
    {'1': 'CARTRIDGE_ACCESS_UNKNOWN', '2': 0},
    {'1': 'CARTRIDGE_ACCESS_INCLUDED', '2': 1},
    {'1': 'CARTRIDGE_ACCESS_LIMITED', '2': 2},
    {'1': 'CARTRIDGE_ACCESS_ADD_ON', '2': 3},
    {'1': 'CARTRIDGE_ACCESS_RESTRICTED', '2': 4},
    {'1': 'CARTRIDGE_ACCESS_BETA', '2': 5},
  ],
};

/// Descriptor for `CartridgeAccessLevel`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List cartridgeAccessLevelDescriptor = $convert.base64Decode(
    'ChRDYXJ0cmlkZ2VBY2Nlc3NMZXZlbBIcChhDQVJUUklER0VfQUNDRVNTX1VOS05PV04QABIdCh'
    'lDQVJUUklER0VfQUNDRVNTX0lOQ0xVREVEEAESHAoYQ0FSVFJJREdFX0FDQ0VTU19MSU1JVEVE'
    'EAISGwoXQ0FSVFJJREdFX0FDQ0VTU19BRERfT04QAxIfChtDQVJUUklER0VfQUNDRVNTX1JFU1'
    'RSSUNURUQQBBIZChVDQVJUUklER0VfQUNDRVNTX0JFVEEQBQ==');

@$core.Deprecated('Use facilityTypeDescriptor instead')
const FacilityType$json = {
  '1': 'FacilityType',
  '2': [
    {'1': 'FACILITY_TYPE_UNKNOWN', '2': 0},
    {'1': 'FACILITY_TYPE_HOSPITAL', '2': 1},
    {'1': 'FACILITY_TYPE_CLINIC', '2': 2},
    {'1': 'FACILITY_TYPE_PHARMACY', '2': 3},
    {'1': 'FACILITY_TYPE_DENTAL', '2': 4},
    {'1': 'FACILITY_TYPE_ORIENTAL', '2': 5},
  ],
};

/// Descriptor for `FacilityType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List facilityTypeDescriptor = $convert.base64Decode(
    'CgxGYWNpbGl0eVR5cGUSGQoVRkFDSUxJVFlfVFlQRV9VTktOT1dOEAASGgoWRkFDSUxJVFlfVF'
    'lQRV9IT1NQSVRBTBABEhgKFEZBQ0lMSVRZX1RZUEVfQ0xJTklDEAISGgoWRkFDSUxJVFlfVFlQ'
    'RV9QSEFSTUFDWRADEhgKFEZBQ0lMSVRZX1RZUEVfREVOVEFMEAQSGgoWRkFDSUxJVFlfVFlQRV'
    '9PUklFTlRBTBAF');

@$core.Deprecated('Use doctorSpecialtyDescriptor instead')
const DoctorSpecialty$json = {
  '1': 'DoctorSpecialty',
  '2': [
    {'1': 'DOCTOR_SPECIALTY_UNKNOWN', '2': 0},
    {'1': 'DOCTOR_SPECIALTY_GENERAL', '2': 1},
    {'1': 'DOCTOR_SPECIALTY_INTERNAL', '2': 2},
    {'1': 'DOCTOR_SPECIALTY_CARDIOLOGY', '2': 3},
    {'1': 'DOCTOR_SPECIALTY_ENDOCRINOLOGY', '2': 4},
    {'1': 'DOCTOR_SPECIALTY_DERMATOLOGY', '2': 5},
    {'1': 'DOCTOR_SPECIALTY_PEDIATRICS', '2': 6},
    {'1': 'DOCTOR_SPECIALTY_PSYCHIATRY', '2': 7},
    {'1': 'DOCTOR_SPECIALTY_ORTHOPEDICS', '2': 8},
    {'1': 'DOCTOR_SPECIALTY_OPHTHALMOLOGY', '2': 9},
    {'1': 'DOCTOR_SPECIALTY_ENT', '2': 10},
  ],
};

/// Descriptor for `DoctorSpecialty`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List doctorSpecialtyDescriptor = $convert.base64Decode(
    'Cg9Eb2N0b3JTcGVjaWFsdHkSHAoYRE9DVE9SX1NQRUNJQUxUWV9VTktOT1dOEAASHAoYRE9DVE'
    '9SX1NQRUNJQUxUWV9HRU5FUkFMEAESHQoZRE9DVE9SX1NQRUNJQUxUWV9JTlRFUk5BTBACEh8K'
    'G0RPQ1RPUl9TUEVDSUFMVFlfQ0FSRElPTE9HWRADEiIKHkRPQ1RPUl9TUEVDSUFMVFlfRU5ET0'
    'NSSU5PTE9HWRAEEiAKHERPQ1RPUl9TUEVDSUFMVFlfREVSTUFUT0xPR1kQBRIfChtET0NUT1Jf'
    'U1BFQ0lBTFRZX1BFRElBVFJJQ1MQBhIfChtET0NUT1JfU1BFQ0lBTFRZX1BTWUNISUFUUlkQBx'
    'IgChxET0NUT1JfU1BFQ0lBTFRZX09SVEhPUEVESUNTEAgSIgoeRE9DVE9SX1NQRUNJQUxUWV9P'
    'UEhUSEFMTU9MT0dZEAkSGAoURE9DVE9SX1NQRUNJQUxUWV9FTlQQCg==');

@$core.Deprecated('Use reservationStatusDescriptor instead')
const ReservationStatus$json = {
  '1': 'ReservationStatus',
  '2': [
    {'1': 'RESERVATION_STATUS_UNKNOWN', '2': 0},
    {'1': 'RESERVATION_STATUS_PENDING', '2': 1},
    {'1': 'RESERVATION_STATUS_CONFIRMED', '2': 2},
    {'1': 'RESERVATION_STATUS_COMPLETED', '2': 3},
    {'1': 'RESERVATION_STATUS_CANCELLED', '2': 4},
    {'1': 'RESERVATION_STATUS_NO_SHOW', '2': 5},
  ],
};

/// Descriptor for `ReservationStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List reservationStatusDescriptor = $convert.base64Decode(
    'ChFSZXNlcnZhdGlvblN0YXR1cxIeChpSRVNFUlZBVElPTl9TVEFUVVNfVU5LTk9XThAAEh4KGl'
    'JFU0VSVkFUSU9OX1NUQVRVU19QRU5ESU5HEAESIAocUkVTRVJWQVRJT05fU1RBVFVTX0NPTkZJ'
    'Uk1FRBACEiAKHFJFU0VSVkFUSU9OX1NUQVRVU19DT01QTEVURUQQAxIgChxSRVNFUlZBVElPTl'
    '9TVEFUVVNfQ0FOQ0VMTEVEEAQSHgoaUkVTRVJWQVRJT05fU1RBVFVTX05PX1NIT1cQBQ==');

@$core.Deprecated('Use adminRoleDescriptor instead')
const AdminRole$json = {
  '1': 'AdminRole',
  '2': [
    {'1': 'ADMIN_ROLE_UNKNOWN', '2': 0},
    {'1': 'ADMIN_ROLE_SUPER', '2': 1},
    {'1': 'ADMIN_ROLE_NATIONAL', '2': 2},
    {'1': 'ADMIN_ROLE_REGIONAL', '2': 3},
    {'1': 'ADMIN_ROLE_BRANCH', '2': 4},
    {'1': 'ADMIN_ROLE_STORE', '2': 5},
  ],
};

/// Descriptor for `AdminRole`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List adminRoleDescriptor = $convert.base64Decode(
    'CglBZG1pblJvbGUSFgoSQURNSU5fUk9MRV9VTktOT1dOEAASFAoQQURNSU5fUk9MRV9TVVBFUh'
    'ABEhcKE0FETUlOX1JPTEVfTkFUSU9OQUwQAhIXChNBRE1JTl9ST0xFX1JFR0lPTkFMEAMSFQoR'
    'QURNSU5fUk9MRV9CUkFOQ0gQBBIUChBBRE1JTl9ST0xFX1NUT1JFEAU=');

@$core.Deprecated('Use auditActionDescriptor instead')
const AuditAction$json = {
  '1': 'AuditAction',
  '2': [
    {'1': 'AUDIT_ACTION_UNKNOWN', '2': 0},
    {'1': 'AUDIT_ACTION_LOGIN', '2': 1},
    {'1': 'AUDIT_ACTION_LOGOUT', '2': 2},
    {'1': 'AUDIT_ACTION_CREATE', '2': 3},
    {'1': 'AUDIT_ACTION_UPDATE', '2': 4},
    {'1': 'AUDIT_ACTION_DELETE', '2': 5},
    {'1': 'AUDIT_ACTION_CONFIG_CHANGE', '2': 6},
    {'1': 'AUDIT_ACTION_ROLE_CHANGE', '2': 7},
  ],
};

/// Descriptor for `AuditAction`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List auditActionDescriptor = $convert.base64Decode(
    'CgtBdWRpdEFjdGlvbhIYChRBVURJVF9BQ1RJT05fVU5LTk9XThAAEhYKEkFVRElUX0FDVElPTl'
    '9MT0dJThABEhcKE0FVRElUX0FDVElPTl9MT0dPVVQQAhIXChNBVURJVF9BQ1RJT05fQ1JFQVRF'
    'EAMSFwoTQVVESVRfQUNUSU9OX1VQREFURRAEEhcKE0FVRElUX0FDVElPTl9ERUxFVEUQBRIeCh'
    'pBVURJVF9BQ1RJT05fQ09ORklHX0NIQU5HRRAGEhwKGEFVRElUX0FDVElPTl9ST0xFX0NIQU5H'
    'RRAH');

@$core.Deprecated('Use familyRoleDescriptor instead')
const FamilyRole$json = {
  '1': 'FamilyRole',
  '2': [
    {'1': 'FAMILY_ROLE_UNKNOWN', '2': 0},
    {'1': 'FAMILY_ROLE_OWNER', '2': 1},
    {'1': 'FAMILY_ROLE_GUARDIAN', '2': 2},
    {'1': 'FAMILY_ROLE_MEMBER', '2': 3},
    {'1': 'FAMILY_ROLE_CHILD', '2': 4},
    {'1': 'FAMILY_ROLE_ELDERLY', '2': 5},
  ],
};

/// Descriptor for `FamilyRole`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List familyRoleDescriptor = $convert.base64Decode(
    'CgpGYW1pbHlSb2xlEhcKE0ZBTUlMWV9ST0xFX1VOS05PV04QABIVChFGQU1JTFlfUk9MRV9PV0'
    '5FUhABEhgKFEZBTUlMWV9ST0xFX0dVQVJESUFOEAISFgoSRkFNSUxZX1JPTEVfTUVNQkVSEAMS'
    'FQoRRkFNSUxZX1JPTEVfQ0hJTEQQBBIXChNGQU1JTFlfUk9MRV9FTERFUkxZEAU=');

@$core.Deprecated('Use invitationStatusDescriptor instead')
const InvitationStatus$json = {
  '1': 'InvitationStatus',
  '2': [
    {'1': 'INVITATION_STATUS_UNKNOWN', '2': 0},
    {'1': 'INVITATION_STATUS_PENDING', '2': 1},
    {'1': 'INVITATION_STATUS_ACCEPTED', '2': 2},
    {'1': 'INVITATION_STATUS_REJECTED', '2': 3},
    {'1': 'INVITATION_STATUS_EXPIRED', '2': 4},
  ],
};

/// Descriptor for `InvitationStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List invitationStatusDescriptor = $convert.base64Decode(
    'ChBJbnZpdGF0aW9uU3RhdHVzEh0KGUlOVklUQVRJT05fU1RBVFVTX1VOS05PV04QABIdChlJTl'
    'ZJVEFUSU9OX1NUQVRVU19QRU5ESU5HEAESHgoaSU5WSVRBVElPTl9TVEFUVVNfQUNDRVBURUQQ'
    'AhIeChpJTlZJVEFUSU9OX1NUQVRVU19SRUpFQ1RFRBADEh0KGUlOVklUQVRJT05fU1RBVFVTX0'
    'VYUElSRUQQBA==');

@$core.Deprecated('Use healthRecordTypeDescriptor instead')
const HealthRecordType$json = {
  '1': 'HealthRecordType',
  '2': [
    {'1': 'HEALTH_RECORD_TYPE_UNKNOWN', '2': 0},
    {'1': 'HEALTH_RECORD_TYPE_LAB_RESULT', '2': 1},
    {'1': 'HEALTH_RECORD_TYPE_IMAGING', '2': 2},
    {'1': 'HEALTH_RECORD_TYPE_VITAL_SIGN', '2': 3},
    {'1': 'HEALTH_RECORD_TYPE_ALLERGY', '2': 4},
    {'1': 'HEALTH_RECORD_TYPE_CONDITION', '2': 5},
    {'1': 'HEALTH_RECORD_TYPE_IMMUNIZATION', '2': 6},
    {'1': 'HEALTH_RECORD_TYPE_PROCEDURE', '2': 7},
  ],
};

/// Descriptor for `HealthRecordType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List healthRecordTypeDescriptor = $convert.base64Decode(
    'ChBIZWFsdGhSZWNvcmRUeXBlEh4KGkhFQUxUSF9SRUNPUkRfVFlQRV9VTktOT1dOEAASIQodSE'
    'VBTFRIX1JFQ09SRF9UWVBFX0xBQl9SRVNVTFQQARIeChpIRUFMVEhfUkVDT1JEX1RZUEVfSU1B'
    'R0lORxACEiEKHUhFQUxUSF9SRUNPUkRfVFlQRV9WSVRBTF9TSUdOEAMSHgoaSEVBTFRIX1JFQ0'
    '9SRF9UWVBFX0FMTEVSR1kQBBIgChxIRUFMVEhfUkVDT1JEX1RZUEVfQ09ORElUSU9OEAUSIwof'
    'SEVBTFRIX1JFQ09SRF9UWVBFX0lNTVVOSVpBVElPThAGEiAKHEhFQUxUSF9SRUNPUkRfVFlQRV'
    '9QUk9DRURVUkUQBw==');

@$core.Deprecated('Use fHIRResourceTypeDescriptor instead')
const FHIRResourceType$json = {
  '1': 'FHIRResourceType',
  '2': [
    {'1': 'FHIR_RESOURCE_TYPE_UNKNOWN', '2': 0},
    {'1': 'FHIR_RESOURCE_TYPE_OBSERVATION', '2': 1},
    {'1': 'FHIR_RESOURCE_TYPE_CONDITION', '2': 2},
    {'1': 'FHIR_RESOURCE_TYPE_ALLERGY_INTOLERANCE', '2': 3},
    {'1': 'FHIR_RESOURCE_TYPE_IMMUNIZATION', '2': 4},
    {'1': 'FHIR_RESOURCE_TYPE_PROCEDURE', '2': 5},
    {'1': 'FHIR_RESOURCE_TYPE_DIAGNOSTIC_REPORT', '2': 6},
    {'1': 'FHIR_RESOURCE_TYPE_BUNDLE', '2': 7},
  ],
};

/// Descriptor for `FHIRResourceType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List fHIRResourceTypeDescriptor = $convert.base64Decode(
    'ChBGSElSUmVzb3VyY2VUeXBlEh4KGkZISVJfUkVTT1VSQ0VfVFlQRV9VTktOT1dOEAASIgoeRk'
    'hJUl9SRVNPVVJDRV9UWVBFX09CU0VSVkFUSU9OEAESIAocRkhJUl9SRVNPVVJDRV9UWVBFX0NP'
    'TkRJVElPThACEioKJkZISVJfUkVTT1VSQ0VfVFlQRV9BTExFUkdZX0lOVE9MRVJBTkNFEAMSIw'
    'ofRkhJUl9SRVNPVVJDRV9UWVBFX0lNTVVOSVpBVElPThAEEiAKHEZISVJfUkVTT1VSQ0VfVFlQ'
    'RV9QUk9DRURVUkUQBRIoCiRGSElSX1JFU09VUkNFX1RZUEVfRElBR05PU1RJQ19SRVBPUlQQBh'
    'IdChlGSElSX1JFU09VUkNFX1RZUEVfQlVORExFEAc=');

@$core.Deprecated('Use prescriptionStatusDescriptor instead')
const PrescriptionStatus$json = {
  '1': 'PrescriptionStatus',
  '2': [
    {'1': 'PRESCRIPTION_STATUS_UNKNOWN', '2': 0},
    {'1': 'PRESCRIPTION_STATUS_DRAFT', '2': 1},
    {'1': 'PRESCRIPTION_STATUS_ACTIVE', '2': 2},
    {'1': 'PRESCRIPTION_STATUS_DISPENSED', '2': 3},
    {'1': 'PRESCRIPTION_STATUS_COMPLETED', '2': 4},
    {'1': 'PRESCRIPTION_STATUS_CANCELLED', '2': 5},
    {'1': 'PRESCRIPTION_STATUS_EXPIRED', '2': 6},
  ],
};

/// Descriptor for `PrescriptionStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List prescriptionStatusDescriptor = $convert.base64Decode(
    'ChJQcmVzY3JpcHRpb25TdGF0dXMSHwobUFJFU0NSSVBUSU9OX1NUQVRVU19VTktOT1dOEAASHQ'
    'oZUFJFU0NSSVBUSU9OX1NUQVRVU19EUkFGVBABEh4KGlBSRVNDUklQVElPTl9TVEFUVVNfQUNU'
    'SVZFEAISIQodUFJFU0NSSVBUSU9OX1NUQVRVU19ESVNQRU5TRUQQAxIhCh1QUkVTQ1JJUFRJT0'
    '5fU1RBVFVTX0NPTVBMRVRFRBAEEiEKHVBSRVNDUklQVElPTl9TVEFUVVNfQ0FOQ0VMTEVEEAUS'
    'HwobUFJFU0NSSVBUSU9OX1NUQVRVU19FWFBJUkVEEAY=');

@$core.Deprecated('Use drugInteractionSeverityDescriptor instead')
const DrugInteractionSeverity$json = {
  '1': 'DrugInteractionSeverity',
  '2': [
    {'1': 'DRUG_INTERACTION_SEVERITY_UNKNOWN', '2': 0},
    {'1': 'DRUG_INTERACTION_SEVERITY_NONE', '2': 1},
    {'1': 'DRUG_INTERACTION_SEVERITY_MINOR', '2': 2},
    {'1': 'DRUG_INTERACTION_SEVERITY_MODERATE', '2': 3},
    {'1': 'DRUG_INTERACTION_SEVERITY_MAJOR', '2': 4},
    {'1': 'DRUG_INTERACTION_SEVERITY_CONTRAINDICATED', '2': 5},
  ],
};

/// Descriptor for `DrugInteractionSeverity`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List drugInteractionSeverityDescriptor = $convert.base64Decode(
    'ChdEcnVnSW50ZXJhY3Rpb25TZXZlcml0eRIlCiFEUlVHX0lOVEVSQUNUSU9OX1NFVkVSSVRZX1'
    'VOS05PV04QABIiCh5EUlVHX0lOVEVSQUNUSU9OX1NFVkVSSVRZX05PTkUQARIjCh9EUlVHX0lO'
    'VEVSQUNUSU9OX1NFVkVSSVRZX01JTk9SEAISJgoiRFJVR19JTlRFUkFDVElPTl9TRVZFUklUWV'
    '9NT0RFUkFURRADEiMKH0RSVUdfSU5URVJBQ1RJT05fU0VWRVJJVFlfTUFKT1IQBBItCilEUlVH'
    'X0lOVEVSQUNUSU9OX1NFVkVSSVRZX0NPTlRSQUlORElDQVRFRBAF');

@$core.Deprecated('Use postCategoryDescriptor instead')
const PostCategory$json = {
  '1': 'PostCategory',
  '2': [
    {'1': 'POST_CATEGORY_UNKNOWN', '2': 0},
    {'1': 'POST_CATEGORY_GENERAL', '2': 1},
    {'1': 'POST_CATEGORY_HEALTH_TIP', '2': 2},
    {'1': 'POST_CATEGORY_QNA', '2': 3},
    {'1': 'POST_CATEGORY_EXPERIENCE', '2': 4},
    {'1': 'POST_CATEGORY_RECIPE', '2': 5},
    {'1': 'POST_CATEGORY_EXERCISE', '2': 6},
  ],
};

/// Descriptor for `PostCategory`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List postCategoryDescriptor = $convert.base64Decode(
    'CgxQb3N0Q2F0ZWdvcnkSGQoVUE9TVF9DQVRFR09SWV9VTktOT1dOEAASGQoVUE9TVF9DQVRFR0'
    '9SWV9HRU5FUkFMEAESHAoYUE9TVF9DQVRFR09SWV9IRUFMVEhfVElQEAISFQoRUE9TVF9DQVRF'
    'R09SWV9RTkEQAxIcChhQT1NUX0NBVEVHT1JZX0VYUEVSSUVOQ0UQBBIYChRQT1NUX0NBVEVHT1'
    'JZX1JFQ0lQRRAFEhoKFlBPU1RfQ0FURUdPUllfRVhFUkNJU0UQBg==');

@$core.Deprecated('Use challengeStatusDescriptor instead')
const ChallengeStatus$json = {
  '1': 'ChallengeStatus',
  '2': [
    {'1': 'CHALLENGE_STATUS_UNKNOWN', '2': 0},
    {'1': 'CHALLENGE_STATUS_UPCOMING', '2': 1},
    {'1': 'CHALLENGE_STATUS_ACTIVE', '2': 2},
    {'1': 'CHALLENGE_STATUS_COMPLETED', '2': 3},
    {'1': 'CHALLENGE_STATUS_CANCELLED', '2': 4},
  ],
};

/// Descriptor for `ChallengeStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List challengeStatusDescriptor = $convert.base64Decode(
    'Cg9DaGFsbGVuZ2VTdGF0dXMSHAoYQ0hBTExFTkdFX1NUQVRVU19VTktOT1dOEAASHQoZQ0hBTE'
    'xFTkdFX1NUQVRVU19VUENPTUlORxABEhsKF0NIQUxMRU5HRV9TVEFUVVNfQUNUSVZFEAISHgoa'
    'Q0hBTExFTkdFX1NUQVRVU19DT01QTEVURUQQAxIeChpDSEFMTEVOR0VfU1RBVFVTX0NBTkNFTE'
    'xFRBAE');

@$core.Deprecated('Use challengeTypeDescriptor instead')
const ChallengeType$json = {
  '1': 'ChallengeType',
  '2': [
    {'1': 'CHALLENGE_TYPE_UNKNOWN', '2': 0},
    {'1': 'CHALLENGE_TYPE_STEPS', '2': 1},
    {'1': 'CHALLENGE_TYPE_MEASUREMENT', '2': 2},
    {'1': 'CHALLENGE_TYPE_DIET', '2': 3},
    {'1': 'CHALLENGE_TYPE_EXERCISE', '2': 4},
    {'1': 'CHALLENGE_TYPE_SLEEP', '2': 5},
  ],
};

/// Descriptor for `ChallengeType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List challengeTypeDescriptor = $convert.base64Decode(
    'Cg1DaGFsbGVuZ2VUeXBlEhoKFkNIQUxMRU5HRV9UWVBFX1VOS05PV04QABIYChRDSEFMTEVOR0'
    'VfVFlQRV9TVEVQUxABEh4KGkNIQUxMRU5HRV9UWVBFX01FQVNVUkVNRU5UEAISFwoTQ0hBTExF'
    'TkdFX1RZUEVfRElFVBADEhsKF0NIQUxMRU5HRV9UWVBFX0VYRVJDSVNFEAQSGAoUQ0hBTExFTk'
    'dFX1RZUEVfU0xFRVAQBQ==');

@$core.Deprecated('Use roomTypeDescriptor instead')
const RoomType$json = {
  '1': 'RoomType',
  '2': [
    {'1': 'ROOM_TYPE_UNKNOWN', '2': 0},
    {'1': 'ROOM_TYPE_ONE_TO_ONE', '2': 1},
    {'1': 'ROOM_TYPE_GROUP', '2': 2},
    {'1': 'ROOM_TYPE_WEBINAR', '2': 3},
    {'1': 'ROOM_TYPE_CONSULTATION', '2': 4},
  ],
};

/// Descriptor for `RoomType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List roomTypeDescriptor = $convert.base64Decode(
    'CghSb29tVHlwZRIVChFST09NX1RZUEVfVU5LTk9XThAAEhgKFFJPT01fVFlQRV9PTkVfVE9fT0'
    '5FEAESEwoPUk9PTV9UWVBFX0dST1VQEAISFQoRUk9PTV9UWVBFX1dFQklOQVIQAxIaChZST09N'
    'X1RZUEVfQ09OU1VMVEFUSU9OEAQ=');

@$core.Deprecated('Use roomStatusDescriptor instead')
const RoomStatus$json = {
  '1': 'RoomStatus',
  '2': [
    {'1': 'ROOM_STATUS_UNKNOWN', '2': 0},
    {'1': 'ROOM_STATUS_WAITING', '2': 1},
    {'1': 'ROOM_STATUS_ACTIVE', '2': 2},
    {'1': 'ROOM_STATUS_ENDED', '2': 3},
    {'1': 'ROOM_STATUS_FAILED', '2': 4},
  ],
};

/// Descriptor for `RoomStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List roomStatusDescriptor = $convert.base64Decode(
    'CgpSb29tU3RhdHVzEhcKE1JPT01fU1RBVFVTX1VOS05PV04QABIXChNST09NX1NUQVRVU19XQU'
    'lUSU5HEAESFgoSUk9PTV9TVEFUVVNfQUNUSVZFEAISFQoRUk9PTV9TVEFUVVNfRU5ERUQQAxIW'
    'ChJST09NX1NUQVRVU19GQUlMRUQQBA==');

@$core.Deprecated('Use signalTypeDescriptor instead')
const SignalType$json = {
  '1': 'SignalType',
  '2': [
    {'1': 'SIGNAL_TYPE_UNKNOWN', '2': 0},
    {'1': 'SIGNAL_TYPE_OFFER', '2': 1},
    {'1': 'SIGNAL_TYPE_ANSWER', '2': 2},
    {'1': 'SIGNAL_TYPE_ICE_CANDIDATE', '2': 3},
    {'1': 'SIGNAL_TYPE_RENEGOTIATE', '2': 4},
    {'1': 'SIGNAL_TYPE_MUTE', '2': 5},
    {'1': 'SIGNAL_TYPE_UNMUTE', '2': 6},
  ],
};

/// Descriptor for `SignalType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List signalTypeDescriptor = $convert.base64Decode(
    'CgpTaWduYWxUeXBlEhcKE1NJR05BTF9UWVBFX1VOS05PV04QABIVChFTSUdOQUxfVFlQRV9PRk'
    'ZFUhABEhYKElNJR05BTF9UWVBFX0FOU1dFUhACEh0KGVNJR05BTF9UWVBFX0lDRV9DQU5ESURB'
    'VEUQAxIbChdTSUdOQUxfVFlQRV9SRU5FR09USUFURRAEEhQKEFNJR05BTF9UWVBFX01VVEUQBR'
    'IWChJTSUdOQUxfVFlQRV9VTk1VVEUQBg==');

@$core.Deprecated('Use notificationTypeDescriptor instead')
const NotificationType$json = {
  '1': 'NotificationType',
  '2': [
    {'1': 'NOTIFICATION_TYPE_UNKNOWN', '2': 0},
    {'1': 'NOTIFICATION_TYPE_MEASUREMENT', '2': 1},
    {'1': 'NOTIFICATION_TYPE_HEALTH_ALERT', '2': 2},
    {'1': 'NOTIFICATION_TYPE_APPOINTMENT', '2': 3},
    {'1': 'NOTIFICATION_TYPE_PRESCRIPTION', '2': 4},
    {'1': 'NOTIFICATION_TYPE_COMMUNITY', '2': 5},
    {'1': 'NOTIFICATION_TYPE_SYSTEM', '2': 6},
    {'1': 'NOTIFICATION_TYPE_PROMOTION', '2': 7},
  ],
};

/// Descriptor for `NotificationType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List notificationTypeDescriptor = $convert.base64Decode(
    'ChBOb3RpZmljYXRpb25UeXBlEh0KGU5PVElGSUNBVElPTl9UWVBFX1VOS05PV04QABIhCh1OT1'
    'RJRklDQVRJT05fVFlQRV9NRUFTVVJFTUVOVBABEiIKHk5PVElGSUNBVElPTl9UWVBFX0hFQUxU'
    'SF9BTEVSVBACEiEKHU5PVElGSUNBVElPTl9UWVBFX0FQUE9JTlRNRU5UEAMSIgoeTk9USUZJQ0'
    'FUSU9OX1RZUEVfUFJFU0NSSVBUSU9OEAQSHwobTk9USUZJQ0FUSU9OX1RZUEVfQ09NTVVOSVRZ'
    'EAUSHAoYTk9USUZJQ0FUSU9OX1RZUEVfU1lTVEVNEAYSHwobTk9USUZJQ0FUSU9OX1RZUEVfUF'
    'JPTU9USU9OEAc=');

@$core.Deprecated('Use notificationChannelDescriptor instead')
const NotificationChannel$json = {
  '1': 'NotificationChannel',
  '2': [
    {'1': 'NOTIFICATION_CHANNEL_UNKNOWN', '2': 0},
    {'1': 'NOTIFICATION_CHANNEL_PUSH', '2': 1},
    {'1': 'NOTIFICATION_CHANNEL_EMAIL', '2': 2},
    {'1': 'NOTIFICATION_CHANNEL_SMS', '2': 3},
    {'1': 'NOTIFICATION_CHANNEL_IN_APP', '2': 4},
  ],
};

/// Descriptor for `NotificationChannel`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List notificationChannelDescriptor = $convert.base64Decode(
    'ChNOb3RpZmljYXRpb25DaGFubmVsEiAKHE5PVElGSUNBVElPTl9DSEFOTkVMX1VOS05PV04QAB'
    'IdChlOT1RJRklDQVRJT05fQ0hBTk5FTF9QVVNIEAESHgoaTk9USUZJQ0FUSU9OX0NIQU5ORUxf'
    'RU1BSUwQAhIcChhOT1RJRklDQVRJT05fQ0hBTk5FTF9TTVMQAxIfChtOT1RJRklDQVRJT05fQ0'
    'hBTk5FTF9JTl9BUFAQBA==');

@$core.Deprecated('Use notificationPriorityDescriptor instead')
const NotificationPriority$json = {
  '1': 'NotificationPriority',
  '2': [
    {'1': 'NOTIFICATION_PRIORITY_UNKNOWN', '2': 0},
    {'1': 'NOTIFICATION_PRIORITY_LOW', '2': 1},
    {'1': 'NOTIFICATION_PRIORITY_NORMAL', '2': 2},
    {'1': 'NOTIFICATION_PRIORITY_HIGH', '2': 3},
    {'1': 'NOTIFICATION_PRIORITY_URGENT', '2': 4},
  ],
};

/// Descriptor for `NotificationPriority`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List notificationPriorityDescriptor = $convert.base64Decode(
    'ChROb3RpZmljYXRpb25Qcmlvcml0eRIhCh1OT1RJRklDQVRJT05fUFJJT1JJVFlfVU5LTk9XTh'
    'AAEh0KGU5PVElGSUNBVElPTl9QUklPUklUWV9MT1cQARIgChxOT1RJRklDQVRJT05fUFJJT1JJ'
    'VFlfTk9STUFMEAISHgoaTk9USUZJQ0FUSU9OX1BSSU9SSVRZX0hJR0gQAxIgChxOT1RJRklDQV'
    'RJT05fUFJJT1JJVFlfVVJHRU5UEAQ=');

@$core.Deprecated('Use consultationStatusDescriptor instead')
const ConsultationStatus$json = {
  '1': 'ConsultationStatus',
  '2': [
    {'1': 'CONSULTATION_STATUS_UNKNOWN', '2': 0},
    {'1': 'CONSULTATION_STATUS_REQUESTED', '2': 1},
    {'1': 'CONSULTATION_STATUS_MATCHED', '2': 2},
    {'1': 'CONSULTATION_STATUS_SCHEDULED', '2': 3},
    {'1': 'CONSULTATION_STATUS_IN_PROGRESS', '2': 4},
    {'1': 'CONSULTATION_STATUS_COMPLETED', '2': 5},
    {'1': 'CONSULTATION_STATUS_CANCELLED', '2': 6},
    {'1': 'CONSULTATION_STATUS_NO_SHOW', '2': 7},
  ],
};

/// Descriptor for `ConsultationStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List consultationStatusDescriptor = $convert.base64Decode(
    'ChJDb25zdWx0YXRpb25TdGF0dXMSHwobQ09OU1VMVEFUSU9OX1NUQVRVU19VTktOT1dOEAASIQ'
    'odQ09OU1VMVEFUSU9OX1NUQVRVU19SRVFVRVNURUQQARIfChtDT05TVUxUQVRJT05fU1RBVFVT'
    'X01BVENIRUQQAhIhCh1DT05TVUxUQVRJT05fU1RBVFVTX1NDSEVEVUxFRBADEiMKH0NPTlNVTF'
    'RBVElPTl9TVEFUVVNfSU5fUFJPR1JFU1MQBBIhCh1DT05TVUxUQVRJT05fU1RBVFVTX0NPTVBM'
    'RVRFRBAFEiEKHUNPTlNVTFRBVElPTl9TVEFUVVNfQ0FOQ0VMTEVEEAYSHwobQ09OU1VMVEFUSU'
    '9OX1NUQVRVU19OT19TSE9XEAc=');

@$core.Deprecated('Use videoSessionStatusDescriptor instead')
const VideoSessionStatus$json = {
  '1': 'VideoSessionStatus',
  '2': [
    {'1': 'VIDEO_SESSION_STATUS_UNKNOWN', '2': 0},
    {'1': 'VIDEO_SESSION_STATUS_WAITING', '2': 1},
    {'1': 'VIDEO_SESSION_STATUS_CONNECTED', '2': 2},
    {'1': 'VIDEO_SESSION_STATUS_ENDED', '2': 3},
    {'1': 'VIDEO_SESSION_STATUS_FAILED', '2': 4},
  ],
};

/// Descriptor for `VideoSessionStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List videoSessionStatusDescriptor = $convert.base64Decode(
    'ChJWaWRlb1Nlc3Npb25TdGF0dXMSIAocVklERU9fU0VTU0lPTl9TVEFUVVNfVU5LTk9XThAAEi'
    'AKHFZJREVPX1NFU1NJT05fU1RBVFVTX1dBSVRJTkcQARIiCh5WSURFT19TRVNTSU9OX1NUQVRV'
    'U19DT05ORUNURUQQAhIeChpWSURFT19TRVNTSU9OX1NUQVRVU19FTkRFRBADEh8KG1ZJREVPX1'
    'NFU1NJT05fU1RBVFVTX0ZBSUxFRBAE');

@$core.Deprecated('Use registerRequestDescriptor instead')
const RegisterRequest$json = {
  '1': 'RegisterRequest',
  '2': [
    {'1': 'email', '3': 1, '4': 1, '5': 9, '10': 'email'},
    {'1': 'password', '3': 2, '4': 1, '5': 9, '10': 'password'},
    {'1': 'display_name', '3': 3, '4': 1, '5': 9, '10': 'displayName'},
    {'1': 'birth_date', '3': 4, '4': 1, '5': 9, '10': 'birthDate'},
    {
      '1': 'gender',
      '3': 5,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.Gender',
      '10': 'gender'
    },
    {'1': 'blood_type', '3': 6, '4': 1, '5': 9, '10': 'bloodType'},
    {'1': 'height_cm', '3': 7, '4': 1, '5': 1, '10': 'heightCm'},
    {'1': 'weight_kg', '3': 8, '4': 1, '5': 1, '10': 'weightKg'},
    {
      '1': 'medical_conditions',
      '3': 9,
      '4': 3,
      '5': 9,
      '10': 'medicalConditions'
    },
    {'1': 'allergies', '3': 10, '4': 3, '5': 9, '10': 'allergies'},
    {
      '1': 'emergency_contact_name',
      '3': 11,
      '4': 1,
      '5': 9,
      '10': 'emergencyContactName'
    },
    {
      '1': 'emergency_contact_phone',
      '3': 12,
      '4': 1,
      '5': 9,
      '10': 'emergencyContactPhone'
    },
    {'1': 'terms_agreed', '3': 13, '4': 1, '5': 8, '10': 'termsAgreed'},
    {'1': 'privacy_agreed', '3': 14, '4': 1, '5': 8, '10': 'privacyAgreed'},
    {
      '1': 'health_data_agreed',
      '3': 15,
      '4': 1,
      '5': 8,
      '10': 'healthDataAgreed'
    },
    {'1': 'marketing_agreed', '3': 16, '4': 1, '5': 8, '10': 'marketingAgreed'},
  ],
};

/// Descriptor for `RegisterRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List registerRequestDescriptor = $convert.base64Decode(
    'Cg9SZWdpc3RlclJlcXVlc3QSFAoFZW1haWwYASABKAlSBWVtYWlsEhoKCHBhc3N3b3JkGAIgAS'
    'gJUghwYXNzd29yZBIhCgxkaXNwbGF5X25hbWUYAyABKAlSC2Rpc3BsYXlOYW1lEh0KCmJpcnRo'
    'X2RhdGUYBCABKAlSCWJpcnRoRGF0ZRIrCgZnZW5kZXIYBSABKA4yEy5tYW5wYXNpay52MS5HZW'
    '5kZXJSBmdlbmRlchIdCgpibG9vZF90eXBlGAYgASgJUglibG9vZFR5cGUSGwoJaGVpZ2h0X2Nt'
    'GAcgASgBUghoZWlnaHRDbRIbCgl3ZWlnaHRfa2cYCCABKAFSCHdlaWdodEtnEi0KEm1lZGljYW'
    'xfY29uZGl0aW9ucxgJIAMoCVIRbWVkaWNhbENvbmRpdGlvbnMSHAoJYWxsZXJnaWVzGAogAygJ'
    'UglhbGxlcmdpZXMSNAoWZW1lcmdlbmN5X2NvbnRhY3RfbmFtZRgLIAEoCVIUZW1lcmdlbmN5Q2'
    '9udGFjdE5hbWUSNgoXZW1lcmdlbmN5X2NvbnRhY3RfcGhvbmUYDCABKAlSFWVtZXJnZW5jeUNv'
    'bnRhY3RQaG9uZRIhCgx0ZXJtc19hZ3JlZWQYDSABKAhSC3Rlcm1zQWdyZWVkEiUKDnByaXZhY3'
    'lfYWdyZWVkGA4gASgIUg1wcml2YWN5QWdyZWVkEiwKEmhlYWx0aF9kYXRhX2FncmVlZBgPIAEo'
    'CFIQaGVhbHRoRGF0YUFncmVlZBIpChBtYXJrZXRpbmdfYWdyZWVkGBAgASgIUg9tYXJrZXRpbm'
    'dBZ3JlZWQ=');

@$core.Deprecated('Use registerResponseDescriptor instead')
const RegisterResponse$json = {
  '1': 'RegisterResponse',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'email', '3': 2, '4': 1, '5': 9, '10': 'email'},
    {'1': 'display_name', '3': 3, '4': 1, '5': 9, '10': 'displayName'},
    {'1': 'role', '3': 4, '4': 1, '5': 9, '10': 'role'},
  ],
};

/// Descriptor for `RegisterResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List registerResponseDescriptor = $convert.base64Decode(
    'ChBSZWdpc3RlclJlc3BvbnNlEhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZBIUCgVlbWFpbBgCIA'
    'EoCVIFZW1haWwSIQoMZGlzcGxheV9uYW1lGAMgASgJUgtkaXNwbGF5TmFtZRISCgRyb2xlGAQg'
    'ASgJUgRyb2xl');

@$core.Deprecated('Use loginRequestDescriptor instead')
const LoginRequest$json = {
  '1': 'LoginRequest',
  '2': [
    {'1': 'email', '3': 1, '4': 1, '5': 9, '10': 'email'},
    {'1': 'password', '3': 2, '4': 1, '5': 9, '10': 'password'},
  ],
};

/// Descriptor for `LoginRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List loginRequestDescriptor = $convert.base64Decode(
    'CgxMb2dpblJlcXVlc3QSFAoFZW1haWwYASABKAlSBWVtYWlsEhoKCHBhc3N3b3JkGAIgASgJUg'
    'hwYXNzd29yZA==');

@$core.Deprecated('Use loginResponseDescriptor instead')
const LoginResponse$json = {
  '1': 'LoginResponse',
  '2': [
    {'1': 'access_token', '3': 1, '4': 1, '5': 9, '10': 'accessToken'},
    {'1': 'refresh_token', '3': 2, '4': 1, '5': 9, '10': 'refreshToken'},
    {'1': 'expires_in', '3': 3, '4': 1, '5': 3, '10': 'expiresIn'},
    {'1': 'token_type', '3': 4, '4': 1, '5': 9, '10': 'tokenType'},
    {'1': 'user_id', '3': 5, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `LoginResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List loginResponseDescriptor = $convert.base64Decode(
    'Cg1Mb2dpblJlc3BvbnNlEiEKDGFjY2Vzc190b2tlbhgBIAEoCVILYWNjZXNzVG9rZW4SIwoNcm'
    'VmcmVzaF90b2tlbhgCIAEoCVIMcmVmcmVzaFRva2VuEh0KCmV4cGlyZXNfaW4YAyABKANSCWV4'
    'cGlyZXNJbhIdCgp0b2tlbl90eXBlGAQgASgJUgl0b2tlblR5cGUSFwoHdXNlcl9pZBgFIAEoCV'
    'IGdXNlcklk');

@$core.Deprecated('Use refreshTokenRequestDescriptor instead')
const RefreshTokenRequest$json = {
  '1': 'RefreshTokenRequest',
  '2': [
    {'1': 'refresh_token', '3': 1, '4': 1, '5': 9, '10': 'refreshToken'},
  ],
};

/// Descriptor for `RefreshTokenRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List refreshTokenRequestDescriptor = $convert.base64Decode(
    'ChNSZWZyZXNoVG9rZW5SZXF1ZXN0EiMKDXJlZnJlc2hfdG9rZW4YASABKAlSDHJlZnJlc2hUb2'
    'tlbg==');

@$core.Deprecated('Use logoutRequestDescriptor instead')
const LogoutRequest$json = {
  '1': 'LogoutRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `LogoutRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List logoutRequestDescriptor = $convert
    .base64Decode('Cg1Mb2dvdXRSZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZA==');

@$core.Deprecated('Use logoutResponseDescriptor instead')
const LogoutResponse$json = {
  '1': 'LogoutResponse',
};

/// Descriptor for `LogoutResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List logoutResponseDescriptor =
    $convert.base64Decode('Cg5Mb2dvdXRSZXNwb25zZQ==');

@$core.Deprecated('Use validateTokenRequestDescriptor instead')
const ValidateTokenRequest$json = {
  '1': 'ValidateTokenRequest',
  '2': [
    {'1': 'access_token', '3': 1, '4': 1, '5': 9, '10': 'accessToken'},
  ],
};

/// Descriptor for `ValidateTokenRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validateTokenRequestDescriptor = $convert.base64Decode(
    'ChRWYWxpZGF0ZVRva2VuUmVxdWVzdBIhCgxhY2Nlc3NfdG9rZW4YASABKAlSC2FjY2Vzc1Rva2'
    'Vu');

@$core.Deprecated('Use validateTokenResponseDescriptor instead')
const ValidateTokenResponse$json = {
  '1': 'ValidateTokenResponse',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'email', '3': 2, '4': 1, '5': 9, '10': 'email'},
    {'1': 'role', '3': 3, '4': 1, '5': 9, '10': 'role'},
  ],
};

/// Descriptor for `ValidateTokenResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validateTokenResponseDescriptor = $convert.base64Decode(
    'ChVWYWxpZGF0ZVRva2VuUmVzcG9uc2USFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklkEhQKBWVtYW'
    'lsGAIgASgJUgVlbWFpbBISCgRyb2xlGAMgASgJUgRyb2xl');

@$core.Deprecated('Use socialLoginRequestDescriptor instead')
const SocialLoginRequest$json = {
  '1': 'SocialLoginRequest',
  '2': [
    {
      '1': 'provider',
      '3': 1,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.SocialProvider',
      '10': 'provider'
    },
    {'1': 'id_token', '3': 2, '4': 1, '5': 9, '10': 'idToken'},
    {'1': 'access_token', '3': 3, '4': 1, '5': 9, '10': 'accessToken'},
    {'1': 'nonce', '3': 4, '4': 1, '5': 9, '10': 'nonce'},
  ],
};

/// Descriptor for `SocialLoginRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List socialLoginRequestDescriptor = $convert.base64Decode(
    'ChJTb2NpYWxMb2dpblJlcXVlc3QSNwoIcHJvdmlkZXIYASABKA4yGy5tYW5wYXNpay52MS5Tb2'
    'NpYWxQcm92aWRlclIIcHJvdmlkZXISGQoIaWRfdG9rZW4YAiABKAlSB2lkVG9rZW4SIQoMYWNj'
    'ZXNzX3Rva2VuGAMgASgJUgthY2Nlc3NUb2tlbhIUCgVub25jZRgEIAEoCVIFbm9uY2U=');

@$core.Deprecated('Use resetPasswordRequestDescriptor instead')
const ResetPasswordRequest$json = {
  '1': 'ResetPasswordRequest',
  '2': [
    {'1': 'email', '3': 1, '4': 1, '5': 9, '10': 'email'},
  ],
};

/// Descriptor for `ResetPasswordRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resetPasswordRequestDescriptor =
    $convert.base64Decode(
        'ChRSZXNldFBhc3N3b3JkUmVxdWVzdBIUCgVlbWFpbBgBIAEoCVIFZW1haWw=');

@$core.Deprecated('Use resetPasswordResponseDescriptor instead')
const ResetPasswordResponse$json = {
  '1': 'ResetPasswordResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
    {'1': 'message', '3': 2, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResetPasswordResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resetPasswordResponseDescriptor = $convert.base64Decode(
    'ChVSZXNldFBhc3N3b3JkUmVzcG9uc2USGAoHc3VjY2VzcxgBIAEoCFIHc3VjY2VzcxIYCgdtZX'
    'NzYWdlGAIgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use startSessionRequestDescriptor instead')
const StartSessionRequest$json = {
  '1': 'StartSessionRequest',
  '2': [
    {'1': 'device_id', '3': 1, '4': 1, '5': 9, '10': 'deviceId'},
    {'1': 'cartridge_id', '3': 2, '4': 1, '5': 9, '10': 'cartridgeId'},
    {'1': 'user_id', '3': 3, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'cartridge_category',
      '3': 4,
      '4': 1,
      '5': 5,
      '10': 'cartridgeCategory'
    },
    {
      '1': 'cartridge_type_index',
      '3': 5,
      '4': 1,
      '5': 5,
      '10': 'cartridgeTypeIndex'
    },
  ],
};

/// Descriptor for `StartSessionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startSessionRequestDescriptor = $convert.base64Decode(
    'ChNTdGFydFNlc3Npb25SZXF1ZXN0EhsKCWRldmljZV9pZBgBIAEoCVIIZGV2aWNlSWQSIQoMY2'
    'FydHJpZGdlX2lkGAIgASgJUgtjYXJ0cmlkZ2VJZBIXCgd1c2VyX2lkGAMgASgJUgZ1c2VySWQS'
    'LQoSY2FydHJpZGdlX2NhdGVnb3J5GAQgASgFUhFjYXJ0cmlkZ2VDYXRlZ29yeRIwChRjYXJ0cm'
    'lkZ2VfdHlwZV9pbmRleBgFIAEoBVISY2FydHJpZGdlVHlwZUluZGV4');

@$core.Deprecated('Use startSessionResponseDescriptor instead')
const StartSessionResponse$json = {
  '1': 'StartSessionResponse',
  '2': [
    {'1': 'session_id', '3': 1, '4': 1, '5': 9, '10': 'sessionId'},
    {
      '1': 'started_at',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'startedAt'
    },
  ],
};

/// Descriptor for `StartSessionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startSessionResponseDescriptor = $convert.base64Decode(
    'ChRTdGFydFNlc3Npb25SZXNwb25zZRIdCgpzZXNzaW9uX2lkGAEgASgJUglzZXNzaW9uSWQSOQ'
    'oKc3RhcnRlZF9hdBgCIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSCXN0YXJ0ZWRB'
    'dA==');

@$core.Deprecated('Use measurementDataDescriptor instead')
const MeasurementData$json = {
  '1': 'MeasurementData',
  '2': [
    {'1': 'session_id', '3': 1, '4': 1, '5': 9, '10': 'sessionId'},
    {'1': 'raw_channels', '3': 2, '4': 3, '5': 1, '10': 'rawChannels'},
    {
      '1': 'differential',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.manpasik.v1.DifferentialCorrection',
      '10': 'differential'
    },
    {
      '1': 'env_meta',
      '3': 4,
      '4': 1,
      '5': 11,
      '6': '.manpasik.v1.EnvironmentMeta',
      '10': 'envMeta'
    },
    {
      '1': 'timestamp',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'timestamp'
    },
  ],
};

/// Descriptor for `MeasurementData`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List measurementDataDescriptor = $convert.base64Decode(
    'Cg9NZWFzdXJlbWVudERhdGESHQoKc2Vzc2lvbl9pZBgBIAEoCVIJc2Vzc2lvbklkEiEKDHJhd1'
    '9jaGFubmVscxgCIAMoAVILcmF3Q2hhbm5lbHMSRwoMZGlmZmVyZW50aWFsGAMgASgLMiMubWFu'
    'cGFzaWsudjEuRGlmZmVyZW50aWFsQ29ycmVjdGlvblIMZGlmZmVyZW50aWFsEjcKCGVudl9tZX'
    'RhGAQgASgLMhwubWFucGFzaWsudjEuRW52aXJvbm1lbnRNZXRhUgdlbnZNZXRhEjgKCXRpbWVz'
    'dGFtcBgFIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSCXRpbWVzdGFtcA==');

@$core.Deprecated('Use differentialCorrectionDescriptor instead')
const DifferentialCorrection$json = {
  '1': 'DifferentialCorrection',
  '2': [
    {'1': 's_det', '3': 1, '4': 1, '5': 1, '10': 'sDet'},
    {'1': 's_ref', '3': 2, '4': 1, '5': 1, '10': 'sRef'},
    {'1': 'alpha', '3': 3, '4': 1, '5': 1, '10': 'alpha'},
    {'1': 's_corrected', '3': 4, '4': 1, '5': 1, '10': 'sCorrected'},
  ],
};

/// Descriptor for `DifferentialCorrection`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List differentialCorrectionDescriptor = $convert.base64Decode(
    'ChZEaWZmZXJlbnRpYWxDb3JyZWN0aW9uEhMKBXNfZGV0GAEgASgBUgRzRGV0EhMKBXNfcmVmGA'
    'IgASgBUgRzUmVmEhQKBWFscGhhGAMgASgBUgVhbHBoYRIfCgtzX2NvcnJlY3RlZBgEIAEoAVIK'
    'c0NvcnJlY3RlZA==');

@$core.Deprecated('Use environmentMetaDescriptor instead')
const EnvironmentMeta$json = {
  '1': 'EnvironmentMeta',
  '2': [
    {'1': 'temp_c', '3': 1, '4': 1, '5': 2, '10': 'tempC'},
    {'1': 'humidity_pct', '3': 2, '4': 1, '5': 2, '10': 'humidityPct'},
    {'1': 'pressure_kpa', '3': 3, '4': 1, '5': 2, '10': 'pressureKpa'},
  ],
};

/// Descriptor for `EnvironmentMeta`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List environmentMetaDescriptor = $convert.base64Decode(
    'Cg9FbnZpcm9ubWVudE1ldGESFQoGdGVtcF9jGAEgASgCUgV0ZW1wQxIhCgxodW1pZGl0eV9wY3'
    'QYAiABKAJSC2h1bWlkaXR5UGN0EiEKDHByZXNzdXJlX2twYRgDIAEoAlILcHJlc3N1cmVLcGE=');

@$core.Deprecated('Use measurementResultDescriptor instead')
const MeasurementResult$json = {
  '1': 'MeasurementResult',
  '2': [
    {'1': 'session_id', '3': 1, '4': 1, '5': 9, '10': 'sessionId'},
    {'1': 'primary_value', '3': 2, '4': 1, '5': 1, '10': 'primaryValue'},
    {'1': 'unit', '3': 3, '4': 1, '5': 9, '10': 'unit'},
    {'1': 'confidence', '3': 4, '4': 1, '5': 1, '10': 'confidence'},
    {
      '1': 'fingerprint_vector',
      '3': 5,
      '4': 3,
      '5': 2,
      '10': 'fingerprintVector'
    },
    {
      '1': 'processed_at',
      '3': 6,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'processedAt'
    },
    {'1': 'evidence_status', '3': 7, '4': 1, '5': 9, '10': 'evidenceStatus'},
    {'1': 'diagnostic_ready', '3': 8, '4': 1, '5': 8, '10': 'diagnosticReady'},
    {'1': 'evidence_gaps', '3': 9, '4': 3, '5': 9, '10': 'evidenceGaps'},
  ],
};

/// Descriptor for `MeasurementResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List measurementResultDescriptor = $convert.base64Decode(
    'ChFNZWFzdXJlbWVudFJlc3VsdBIdCgpzZXNzaW9uX2lkGAEgASgJUglzZXNzaW9uSWQSIwoNcH'
    'JpbWFyeV92YWx1ZRgCIAEoAVIMcHJpbWFyeVZhbHVlEhIKBHVuaXQYAyABKAlSBHVuaXQSHgoK'
    'Y29uZmlkZW5jZRgEIAEoAVIKY29uZmlkZW5jZRItChJmaW5nZXJwcmludF92ZWN0b3IYBSADKA'
    'JSEWZpbmdlcnByaW50VmVjdG9yEj0KDHByb2Nlc3NlZF9hdBgGIAEoCzIaLmdvb2dsZS5wcm90'
    'b2J1Zi5UaW1lc3RhbXBSC3Byb2Nlc3NlZEF0EicKD2V2aWRlbmNlX3N0YXR1cxgHIAEoCVIOZX'
    'ZpZGVuY2VTdGF0dXMSKQoQZGlhZ25vc3RpY19yZWFkeRgIIAEoCFIPZGlhZ25vc3RpY1JlYWR5'
    'EiMKDWV2aWRlbmNlX2dhcHMYCSADKAlSDGV2aWRlbmNlR2Fwcw==');

@$core.Deprecated('Use endSessionRequestDescriptor instead')
const EndSessionRequest$json = {
  '1': 'EndSessionRequest',
  '2': [
    {'1': 'session_id', '3': 1, '4': 1, '5': 9, '10': 'sessionId'},
  ],
};

/// Descriptor for `EndSessionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List endSessionRequestDescriptor = $convert.base64Decode(
    'ChFFbmRTZXNzaW9uUmVxdWVzdBIdCgpzZXNzaW9uX2lkGAEgASgJUglzZXNzaW9uSWQ=');

@$core.Deprecated('Use endSessionResponseDescriptor instead')
const EndSessionResponse$json = {
  '1': 'EndSessionResponse',
  '2': [
    {'1': 'session_id', '3': 1, '4': 1, '5': 9, '10': 'sessionId'},
    {
      '1': 'total_measurements',
      '3': 2,
      '4': 1,
      '5': 5,
      '10': 'totalMeasurements'
    },
    {
      '1': 'ended_at',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'endedAt'
    },
  ],
};

/// Descriptor for `EndSessionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List endSessionResponseDescriptor = $convert.base64Decode(
    'ChJFbmRTZXNzaW9uUmVzcG9uc2USHQoKc2Vzc2lvbl9pZBgBIAEoCVIJc2Vzc2lvbklkEi0KEn'
    'RvdGFsX21lYXN1cmVtZW50cxgCIAEoBVIRdG90YWxNZWFzdXJlbWVudHMSNQoIZW5kZWRfYXQY'
    'AyABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUgdlbmRlZEF0');

@$core.Deprecated('Use getHistoryRequestDescriptor instead')
const GetHistoryRequest$json = {
  '1': 'GetHistoryRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'start_time',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'startTime'
    },
    {
      '1': 'end_time',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'endTime'
    },
    {'1': 'limit', '3': 4, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 5, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `GetHistoryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getHistoryRequestDescriptor = $convert.base64Decode(
    'ChFHZXRIaXN0b3J5UmVxdWVzdBIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQSOQoKc3RhcnRfdG'
    'ltZRgCIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSCXN0YXJ0VGltZRI1CghlbmRf'
    'dGltZRgDIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSB2VuZFRpbWUSFAoFbGltaX'
    'QYBCABKAVSBWxpbWl0EhYKBm9mZnNldBgFIAEoBVIGb2Zmc2V0');

@$core.Deprecated('Use getHistoryResponseDescriptor instead')
const GetHistoryResponse$json = {
  '1': 'GetHistoryResponse',
  '2': [
    {
      '1': 'measurements',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.MeasurementSummary',
      '10': 'measurements'
    },
    {'1': 'total_count', '3': 2, '4': 1, '5': 5, '10': 'totalCount'},
  ],
};

/// Descriptor for `GetHistoryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getHistoryResponseDescriptor = $convert.base64Decode(
    'ChJHZXRIaXN0b3J5UmVzcG9uc2USQwoMbWVhc3VyZW1lbnRzGAEgAygLMh8ubWFucGFzaWsudj'
    'EuTWVhc3VyZW1lbnRTdW1tYXJ5UgxtZWFzdXJlbWVudHMSHwoLdG90YWxfY291bnQYAiABKAVS'
    'CnRvdGFsQ291bnQ=');

@$core.Deprecated('Use measurementSummaryDescriptor instead')
const MeasurementSummary$json = {
  '1': 'MeasurementSummary',
  '2': [
    {'1': 'session_id', '3': 1, '4': 1, '5': 9, '10': 'sessionId'},
    {'1': 'cartridge_type', '3': 2, '4': 1, '5': 9, '10': 'cartridgeType'},
    {'1': 'primary_value', '3': 3, '4': 1, '5': 1, '10': 'primaryValue'},
    {'1': 'unit', '3': 4, '4': 1, '5': 9, '10': 'unit'},
    {
      '1': 'measured_at',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'measuredAt'
    },
    {'1': 'evidence_status', '3': 6, '4': 1, '5': 9, '10': 'evidenceStatus'},
    {'1': 'diagnostic_ready', '3': 7, '4': 1, '5': 8, '10': 'diagnosticReady'},
    {'1': 'evidence_gaps', '3': 8, '4': 3, '5': 9, '10': 'evidenceGaps'},
  ],
};

/// Descriptor for `MeasurementSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List measurementSummaryDescriptor = $convert.base64Decode(
    'ChJNZWFzdXJlbWVudFN1bW1hcnkSHQoKc2Vzc2lvbl9pZBgBIAEoCVIJc2Vzc2lvbklkEiUKDm'
    'NhcnRyaWRnZV90eXBlGAIgASgJUg1jYXJ0cmlkZ2VUeXBlEiMKDXByaW1hcnlfdmFsdWUYAyAB'
    'KAFSDHByaW1hcnlWYWx1ZRISCgR1bml0GAQgASgJUgR1bml0EjsKC21lYXN1cmVkX2F0GAUgAS'
    'gLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFIKbWVhc3VyZWRBdBInCg9ldmlkZW5jZV9z'
    'dGF0dXMYBiABKAlSDmV2aWRlbmNlU3RhdHVzEikKEGRpYWdub3N0aWNfcmVhZHkYByABKAhSD2'
    'RpYWdub3N0aWNSZWFkeRIjCg1ldmlkZW5jZV9nYXBzGAggAygJUgxldmlkZW5jZUdhcHM=');

@$core.Deprecated('Use registerDeviceRequestDescriptor instead')
const RegisterDeviceRequest$json = {
  '1': 'RegisterDeviceRequest',
  '2': [
    {'1': 'device_id', '3': 1, '4': 1, '5': 9, '10': 'deviceId'},
    {'1': 'serial_number', '3': 2, '4': 1, '5': 9, '10': 'serialNumber'},
    {'1': 'firmware_version', '3': 3, '4': 1, '5': 9, '10': 'firmwareVersion'},
    {'1': 'user_id', '3': 4, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `RegisterDeviceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List registerDeviceRequestDescriptor = $convert.base64Decode(
    'ChVSZWdpc3RlckRldmljZVJlcXVlc3QSGwoJZGV2aWNlX2lkGAEgASgJUghkZXZpY2VJZBIjCg'
    '1zZXJpYWxfbnVtYmVyGAIgASgJUgxzZXJpYWxOdW1iZXISKQoQZmlybXdhcmVfdmVyc2lvbhgD'
    'IAEoCVIPZmlybXdhcmVWZXJzaW9uEhcKB3VzZXJfaWQYBCABKAlSBnVzZXJJZA==');

@$core.Deprecated('Use registerDeviceResponseDescriptor instead')
const RegisterDeviceResponse$json = {
  '1': 'RegisterDeviceResponse',
  '2': [
    {'1': 'device_id', '3': 1, '4': 1, '5': 9, '10': 'deviceId'},
    {
      '1': 'registration_token',
      '3': 2,
      '4': 1,
      '5': 9,
      '10': 'registrationToken'
    },
    {
      '1': 'registered_at',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'registeredAt'
    },
  ],
};

/// Descriptor for `RegisterDeviceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List registerDeviceResponseDescriptor = $convert.base64Decode(
    'ChZSZWdpc3RlckRldmljZVJlc3BvbnNlEhsKCWRldmljZV9pZBgBIAEoCVIIZGV2aWNlSWQSLQ'
    'oScmVnaXN0cmF0aW9uX3Rva2VuGAIgASgJUhFyZWdpc3RyYXRpb25Ub2tlbhI/Cg1yZWdpc3Rl'
    'cmVkX2F0GAMgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFIMcmVnaXN0ZXJlZEF0');

@$core.Deprecated('Use listDevicesRequestDescriptor instead')
const ListDevicesRequest$json = {
  '1': 'ListDevicesRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `ListDevicesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDevicesRequestDescriptor =
    $convert.base64Decode(
        'ChJMaXN0RGV2aWNlc1JlcXVlc3QSFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklk');

@$core.Deprecated('Use listDevicesResponseDescriptor instead')
const ListDevicesResponse$json = {
  '1': 'ListDevicesResponse',
  '2': [
    {
      '1': 'devices',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.DeviceInfo',
      '10': 'devices'
    },
  ],
};

/// Descriptor for `ListDevicesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDevicesResponseDescriptor = $convert.base64Decode(
    'ChNMaXN0RGV2aWNlc1Jlc3BvbnNlEjEKB2RldmljZXMYASADKAsyFy5tYW5wYXNpay52MS5EZX'
    'ZpY2VJbmZvUgdkZXZpY2Vz');

@$core.Deprecated('Use deviceInfoDescriptor instead')
const DeviceInfo$json = {
  '1': 'DeviceInfo',
  '2': [
    {'1': 'device_id', '3': 1, '4': 1, '5': 9, '10': 'deviceId'},
    {'1': 'name', '3': 2, '4': 1, '5': 9, '10': 'name'},
    {'1': 'firmware_version', '3': 3, '4': 1, '5': 9, '10': 'firmwareVersion'},
    {
      '1': 'status',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.DeviceStatus',
      '10': 'status'
    },
    {'1': 'battery_percent', '3': 5, '4': 1, '5': 5, '10': 'batteryPercent'},
    {
      '1': 'last_seen',
      '3': 6,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'lastSeen'
    },
  ],
};

/// Descriptor for `DeviceInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deviceInfoDescriptor = $convert.base64Decode(
    'CgpEZXZpY2VJbmZvEhsKCWRldmljZV9pZBgBIAEoCVIIZGV2aWNlSWQSEgoEbmFtZRgCIAEoCV'
    'IEbmFtZRIpChBmaXJtd2FyZV92ZXJzaW9uGAMgASgJUg9maXJtd2FyZVZlcnNpb24SMQoGc3Rh'
    'dHVzGAQgASgOMhkubWFucGFzaWsudjEuRGV2aWNlU3RhdHVzUgZzdGF0dXMSJwoPYmF0dGVyeV'
    '9wZXJjZW50GAUgASgFUg5iYXR0ZXJ5UGVyY2VudBI3CglsYXN0X3NlZW4YBiABKAsyGi5nb29n'
    'bGUucHJvdG9idWYuVGltZXN0YW1wUghsYXN0U2Vlbg==');

@$core.Deprecated('Use deviceStatusUpdateDescriptor instead')
const DeviceStatusUpdate$json = {
  '1': 'DeviceStatusUpdate',
  '2': [
    {'1': 'device_id', '3': 1, '4': 1, '5': 9, '10': 'deviceId'},
    {
      '1': 'status',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.DeviceStatus',
      '10': 'status'
    },
    {'1': 'battery_percent', '3': 3, '4': 1, '5': 5, '10': 'batteryPercent'},
    {'1': 'signal_strength', '3': 4, '4': 1, '5': 5, '10': 'signalStrength'},
    {
      '1': 'timestamp',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'timestamp'
    },
  ],
};

/// Descriptor for `DeviceStatusUpdate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deviceStatusUpdateDescriptor = $convert.base64Decode(
    'ChJEZXZpY2VTdGF0dXNVcGRhdGUSGwoJZGV2aWNlX2lkGAEgASgJUghkZXZpY2VJZBIxCgZzdG'
    'F0dXMYAiABKA4yGS5tYW5wYXNpay52MS5EZXZpY2VTdGF0dXNSBnN0YXR1cxInCg9iYXR0ZXJ5'
    'X3BlcmNlbnQYAyABKAVSDmJhdHRlcnlQZXJjZW50EicKD3NpZ25hbF9zdHJlbmd0aBgEIAEoBV'
    'IOc2lnbmFsU3RyZW5ndGgSOAoJdGltZXN0YW1wGAUgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRp'
    'bWVzdGFtcFIJdGltZXN0YW1w');

@$core.Deprecated('Use deviceCommandDescriptor instead')
const DeviceCommand$json = {
  '1': 'DeviceCommand',
  '2': [
    {'1': 'command_id', '3': 1, '4': 1, '5': 9, '10': 'commandId'},
    {
      '1': 'command_type',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.CommandType',
      '10': 'commandType'
    },
    {'1': 'payload', '3': 3, '4': 1, '5': 12, '10': 'payload'},
  ],
};

/// Descriptor for `DeviceCommand`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deviceCommandDescriptor = $convert.base64Decode(
    'Cg1EZXZpY2VDb21tYW5kEh0KCmNvbW1hbmRfaWQYASABKAlSCWNvbW1hbmRJZBI7Cgxjb21tYW'
    '5kX3R5cGUYAiABKA4yGC5tYW5wYXNpay52MS5Db21tYW5kVHlwZVILY29tbWFuZFR5cGUSGAoH'
    'cGF5bG9hZBgDIAEoDFIHcGF5bG9hZA==');

@$core.Deprecated('Use otaRequestDescriptor instead')
const OtaRequest$json = {
  '1': 'OtaRequest',
  '2': [
    {'1': 'device_id', '3': 1, '4': 1, '5': 9, '10': 'deviceId'},
    {'1': 'target_version', '3': 2, '4': 1, '5': 9, '10': 'targetVersion'},
  ],
};

/// Descriptor for `OtaRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List otaRequestDescriptor = $convert.base64Decode(
    'CgpPdGFSZXF1ZXN0EhsKCWRldmljZV9pZBgBIAEoCVIIZGV2aWNlSWQSJQoOdGFyZ2V0X3Zlcn'
    'Npb24YAiABKAlSDXRhcmdldFZlcnNpb24=');

@$core.Deprecated('Use otaResponseDescriptor instead')
const OtaResponse$json = {
  '1': 'OtaResponse',
  '2': [
    {'1': 'update_id', '3': 1, '4': 1, '5': 9, '10': 'updateId'},
    {'1': 'download_url', '3': 2, '4': 1, '5': 9, '10': 'downloadUrl'},
    {'1': 'checksum', '3': 3, '4': 1, '5': 9, '10': 'checksum'},
  ],
};

/// Descriptor for `OtaResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List otaResponseDescriptor = $convert.base64Decode(
    'CgtPdGFSZXNwb25zZRIbCgl1cGRhdGVfaWQYASABKAlSCHVwZGF0ZUlkEiEKDGRvd25sb2FkX3'
    'VybBgCIAEoCVILZG93bmxvYWRVcmwSGgoIY2hlY2tzdW0YAyABKAlSCGNoZWNrc3Vt');

@$core.Deprecated('Use getProfileRequestDescriptor instead')
const GetProfileRequest$json = {
  '1': 'GetProfileRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `GetProfileRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getProfileRequestDescriptor = $convert.base64Decode(
    'ChFHZXRQcm9maWxlUmVxdWVzdBIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQ=');

@$core.Deprecated('Use updateProfileRequestDescriptor instead')
const UpdateProfileRequest$json = {
  '1': 'UpdateProfileRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'display_name', '3': 2, '4': 1, '5': 9, '10': 'displayName'},
    {'1': 'avatar_url', '3': 3, '4': 1, '5': 9, '10': 'avatarUrl'},
    {'1': 'language', '3': 4, '4': 1, '5': 9, '10': 'language'},
    {'1': 'timezone', '3': 5, '4': 1, '5': 9, '10': 'timezone'},
    {'1': 'birth_date', '3': 6, '4': 1, '5': 9, '10': 'birthDate'},
    {
      '1': 'gender',
      '3': 7,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.Gender',
      '10': 'gender'
    },
    {'1': 'blood_type', '3': 8, '4': 1, '5': 9, '10': 'bloodType'},
    {'1': 'height_cm', '3': 9, '4': 1, '5': 1, '10': 'heightCm'},
    {'1': 'weight_kg', '3': 10, '4': 1, '5': 1, '10': 'weightKg'},
    {
      '1': 'medical_conditions',
      '3': 11,
      '4': 3,
      '5': 9,
      '10': 'medicalConditions'
    },
    {'1': 'allergies', '3': 12, '4': 3, '5': 9, '10': 'allergies'},
    {
      '1': 'emergency_contact_name',
      '3': 13,
      '4': 1,
      '5': 9,
      '10': 'emergencyContactName'
    },
    {
      '1': 'emergency_contact_phone',
      '3': 14,
      '4': 1,
      '5': 9,
      '10': 'emergencyContactPhone'
    },
  ],
};

/// Descriptor for `UpdateProfileRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateProfileRequestDescriptor = $convert.base64Decode(
    'ChRVcGRhdGVQcm9maWxlUmVxdWVzdBIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQSIQoMZGlzcG'
    'xheV9uYW1lGAIgASgJUgtkaXNwbGF5TmFtZRIdCgphdmF0YXJfdXJsGAMgASgJUglhdmF0YXJV'
    'cmwSGgoIbGFuZ3VhZ2UYBCABKAlSCGxhbmd1YWdlEhoKCHRpbWV6b25lGAUgASgJUgh0aW1lem'
    '9uZRIdCgpiaXJ0aF9kYXRlGAYgASgJUgliaXJ0aERhdGUSKwoGZ2VuZGVyGAcgASgOMhMubWFu'
    'cGFzaWsudjEuR2VuZGVyUgZnZW5kZXISHQoKYmxvb2RfdHlwZRgIIAEoCVIJYmxvb2RUeXBlEh'
    'sKCWhlaWdodF9jbRgJIAEoAVIIaGVpZ2h0Q20SGwoJd2VpZ2h0X2tnGAogASgBUgh3ZWlnaHRL'
    'ZxItChJtZWRpY2FsX2NvbmRpdGlvbnMYCyADKAlSEW1lZGljYWxDb25kaXRpb25zEhwKCWFsbG'
    'VyZ2llcxgMIAMoCVIJYWxsZXJnaWVzEjQKFmVtZXJnZW5jeV9jb250YWN0X25hbWUYDSABKAlS'
    'FGVtZXJnZW5jeUNvbnRhY3ROYW1lEjYKF2VtZXJnZW5jeV9jb250YWN0X3Bob25lGA4gASgJUh'
    'VlbWVyZ2VuY3lDb250YWN0UGhvbmU=');

@$core.Deprecated('Use userProfileDescriptor instead')
const UserProfile$json = {
  '1': 'UserProfile',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'email', '3': 2, '4': 1, '5': 9, '10': 'email'},
    {'1': 'display_name', '3': 3, '4': 1, '5': 9, '10': 'displayName'},
    {'1': 'avatar_url', '3': 4, '4': 1, '5': 9, '10': 'avatarUrl'},
    {'1': 'language', '3': 5, '4': 1, '5': 9, '10': 'language'},
    {'1': 'timezone', '3': 6, '4': 1, '5': 9, '10': 'timezone'},
    {
      '1': 'subscription_tier',
      '3': 7,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.SubscriptionTier',
      '10': 'subscriptionTier'
    },
    {
      '1': 'created_at',
      '3': 8,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
    {'1': 'birth_date', '3': 9, '4': 1, '5': 9, '10': 'birthDate'},
    {
      '1': 'gender',
      '3': 10,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.Gender',
      '10': 'gender'
    },
    {'1': 'blood_type', '3': 11, '4': 1, '5': 9, '10': 'bloodType'},
    {'1': 'height_cm', '3': 12, '4': 1, '5': 1, '10': 'heightCm'},
    {'1': 'weight_kg', '3': 13, '4': 1, '5': 1, '10': 'weightKg'},
    {
      '1': 'medical_conditions',
      '3': 14,
      '4': 3,
      '5': 9,
      '10': 'medicalConditions'
    },
    {'1': 'allergies', '3': 15, '4': 3, '5': 9, '10': 'allergies'},
    {
      '1': 'emergency_contact_name',
      '3': 16,
      '4': 1,
      '5': 9,
      '10': 'emergencyContactName'
    },
    {
      '1': 'emergency_contact_phone',
      '3': 17,
      '4': 1,
      '5': 9,
      '10': 'emergencyContactPhone'
    },
    {
      '1': 'social_provider',
      '3': 18,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.SocialProvider',
      '10': 'socialProvider'
    },
    {'1': 'social_id', '3': 19, '4': 1, '5': 9, '10': 'socialId'},
  ],
};

/// Descriptor for `UserProfile`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List userProfileDescriptor = $convert.base64Decode(
    'CgtVc2VyUHJvZmlsZRIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQSFAoFZW1haWwYAiABKAlSBW'
    'VtYWlsEiEKDGRpc3BsYXlfbmFtZRgDIAEoCVILZGlzcGxheU5hbWUSHQoKYXZhdGFyX3VybBgE'
    'IAEoCVIJYXZhdGFyVXJsEhoKCGxhbmd1YWdlGAUgASgJUghsYW5ndWFnZRIaCgh0aW1lem9uZR'
    'gGIAEoCVIIdGltZXpvbmUSSgoRc3Vic2NyaXB0aW9uX3RpZXIYByABKA4yHS5tYW5wYXNpay52'
    'MS5TdWJzY3JpcHRpb25UaWVyUhBzdWJzY3JpcHRpb25UaWVyEjkKCmNyZWF0ZWRfYXQYCCABKA'
    'syGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUgljcmVhdGVkQXQSHQoKYmlydGhfZGF0ZRgJ'
    'IAEoCVIJYmlydGhEYXRlEisKBmdlbmRlchgKIAEoDjITLm1hbnBhc2lrLnYxLkdlbmRlclIGZ2'
    'VuZGVyEh0KCmJsb29kX3R5cGUYCyABKAlSCWJsb29kVHlwZRIbCgloZWlnaHRfY20YDCABKAFS'
    'CGhlaWdodENtEhsKCXdlaWdodF9rZxgNIAEoAVIId2VpZ2h0S2cSLQoSbWVkaWNhbF9jb25kaX'
    'Rpb25zGA4gAygJUhFtZWRpY2FsQ29uZGl0aW9ucxIcCglhbGxlcmdpZXMYDyADKAlSCWFsbGVy'
    'Z2llcxI0ChZlbWVyZ2VuY3lfY29udGFjdF9uYW1lGBAgASgJUhRlbWVyZ2VuY3lDb250YWN0Tm'
    'FtZRI2ChdlbWVyZ2VuY3lfY29udGFjdF9waG9uZRgRIAEoCVIVZW1lcmdlbmN5Q29udGFjdFBo'
    'b25lEkQKD3NvY2lhbF9wcm92aWRlchgSIAEoDjIbLm1hbnBhc2lrLnYxLlNvY2lhbFByb3ZpZG'
    'VyUg5zb2NpYWxQcm92aWRlchIbCglzb2NpYWxfaWQYEyABKAlSCHNvY2lhbElk');

@$core.Deprecated('Use getSubscriptionRequestDescriptor instead')
const GetSubscriptionRequest$json = {
  '1': 'GetSubscriptionRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `GetSubscriptionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSubscriptionRequestDescriptor =
    $convert.base64Decode(
        'ChZHZXRTdWJzY3JpcHRpb25SZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZA==');

@$core.Deprecated('Use subscriptionInfoDescriptor instead')
const SubscriptionInfo$json = {
  '1': 'SubscriptionInfo',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'tier',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.SubscriptionTier',
      '10': 'tier'
    },
    {
      '1': 'started_at',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'startedAt'
    },
    {
      '1': 'expires_at',
      '3': 4,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'expiresAt'
    },
    {'1': 'max_devices', '3': 5, '4': 1, '5': 5, '10': 'maxDevices'},
    {
      '1': 'max_family_members',
      '3': 6,
      '4': 1,
      '5': 5,
      '10': 'maxFamilyMembers'
    },
    {
      '1': 'ai_coaching_enabled',
      '3': 7,
      '4': 1,
      '5': 8,
      '10': 'aiCoachingEnabled'
    },
    {
      '1': 'telemedicine_enabled',
      '3': 8,
      '4': 1,
      '5': 8,
      '10': 'telemedicineEnabled'
    },
  ],
};

/// Descriptor for `SubscriptionInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List subscriptionInfoDescriptor = $convert.base64Decode(
    'ChBTdWJzY3JpcHRpb25JbmZvEhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZBIxCgR0aWVyGAIgAS'
    'gOMh0ubWFucGFzaWsudjEuU3Vic2NyaXB0aW9uVGllclIEdGllchI5CgpzdGFydGVkX2F0GAMg'
    'ASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFIJc3RhcnRlZEF0EjkKCmV4cGlyZXNfYX'
    'QYBCABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUglleHBpcmVzQXQSHwoLbWF4X2Rl'
    'dmljZXMYBSABKAVSCm1heERldmljZXMSLAoSbWF4X2ZhbWlseV9tZW1iZXJzGAYgASgFUhBtYX'
    'hGYW1pbHlNZW1iZXJzEi4KE2FpX2NvYWNoaW5nX2VuYWJsZWQYByABKAhSEWFpQ29hY2hpbmdF'
    'bmFibGVkEjEKFHRlbGVtZWRpY2luZV9lbmFibGVkGAggASgIUhN0ZWxlbWVkaWNpbmVFbmFibG'
    'Vk');

@$core.Deprecated('Use createSubscriptionRequestDescriptor instead')
const CreateSubscriptionRequest$json = {
  '1': 'CreateSubscriptionRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'tier',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.SubscriptionTier',
      '10': 'tier'
    },
  ],
};

/// Descriptor for `CreateSubscriptionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createSubscriptionRequestDescriptor =
    $convert.base64Decode(
        'ChlDcmVhdGVTdWJzY3JpcHRpb25SZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZBIxCg'
        'R0aWVyGAIgASgOMh0ubWFucGFzaWsudjEuU3Vic2NyaXB0aW9uVGllclIEdGllcg==');

@$core.Deprecated('Use getSubscriptionDetailRequestDescriptor instead')
const GetSubscriptionDetailRequest$json = {
  '1': 'GetSubscriptionDetailRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `GetSubscriptionDetailRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSubscriptionDetailRequestDescriptor =
    $convert.base64Decode(
        'ChxHZXRTdWJzY3JpcHRpb25EZXRhaWxSZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZA'
        '==');

@$core.Deprecated('Use updateSubscriptionRequestDescriptor instead')
const UpdateSubscriptionRequest$json = {
  '1': 'UpdateSubscriptionRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'new_tier',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.SubscriptionTier',
      '10': 'newTier'
    },
    {'1': 'payment_id', '3': 3, '4': 1, '5': 9, '10': 'paymentId'},
  ],
};

/// Descriptor for `UpdateSubscriptionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateSubscriptionRequestDescriptor = $convert.base64Decode(
    'ChlVcGRhdGVTdWJzY3JpcHRpb25SZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZBI4Cg'
    'huZXdfdGllchgCIAEoDjIdLm1hbnBhc2lrLnYxLlN1YnNjcmlwdGlvblRpZXJSB25ld1RpZXIS'
    'HQoKcGF5bWVudF9pZBgDIAEoCVIJcGF5bWVudElk');

@$core.Deprecated('Use cancelSubscriptionRequestDescriptor instead')
const CancelSubscriptionRequest$json = {
  '1': 'CancelSubscriptionRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'reason', '3': 2, '4': 1, '5': 9, '10': 'reason'},
  ],
};

/// Descriptor for `CancelSubscriptionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelSubscriptionRequestDescriptor =
    $convert.base64Decode(
        'ChlDYW5jZWxTdWJzY3JpcHRpb25SZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZBIWCg'
        'ZyZWFzb24YAiABKAlSBnJlYXNvbg==');

@$core.Deprecated('Use cancelSubscriptionResponseDescriptor instead')
const CancelSubscriptionResponse$json = {
  '1': 'CancelSubscriptionResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
    {
      '1': 'cancelled_at',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'cancelledAt'
    },
    {
      '1': 'effective_until',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'effectiveUntil'
    },
  ],
};

/// Descriptor for `CancelSubscriptionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelSubscriptionResponseDescriptor = $convert.base64Decode(
    'ChpDYW5jZWxTdWJzY3JpcHRpb25SZXNwb25zZRIYCgdzdWNjZXNzGAEgASgIUgdzdWNjZXNzEj'
    '0KDGNhbmNlbGxlZF9hdBgCIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSC2NhbmNl'
    'bGxlZEF0EkMKD2VmZmVjdGl2ZV91bnRpbBgDIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3'
    'RhbXBSDmVmZmVjdGl2ZVVudGls');

@$core.Deprecated('Use subscriptionDetailDescriptor instead')
const SubscriptionDetail$json = {
  '1': 'SubscriptionDetail',
  '2': [
    {'1': 'subscription_id', '3': 1, '4': 1, '5': 9, '10': 'subscriptionId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'tier',
      '3': 3,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.SubscriptionTier',
      '10': 'tier'
    },
    {
      '1': 'status',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.SubscriptionStatus',
      '10': 'status'
    },
    {
      '1': 'started_at',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'startedAt'
    },
    {
      '1': 'expires_at',
      '3': 6,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'expiresAt'
    },
    {
      '1': 'cancelled_at',
      '3': 7,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'cancelledAt'
    },
    {'1': 'max_devices', '3': 8, '4': 1, '5': 5, '10': 'maxDevices'},
    {
      '1': 'max_family_members',
      '3': 9,
      '4': 1,
      '5': 5,
      '10': 'maxFamilyMembers'
    },
    {
      '1': 'ai_coaching_enabled',
      '3': 10,
      '4': 1,
      '5': 8,
      '10': 'aiCoachingEnabled'
    },
    {
      '1': 'telemedicine_enabled',
      '3': 11,
      '4': 1,
      '5': 8,
      '10': 'telemedicineEnabled'
    },
    {
      '1': 'monthly_price_krw',
      '3': 12,
      '4': 1,
      '5': 5,
      '10': 'monthlyPriceKrw'
    },
    {'1': 'auto_renew', '3': 13, '4': 1, '5': 8, '10': 'autoRenew'},
  ],
};

/// Descriptor for `SubscriptionDetail`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List subscriptionDetailDescriptor = $convert.base64Decode(
    'ChJTdWJzY3JpcHRpb25EZXRhaWwSJwoPc3Vic2NyaXB0aW9uX2lkGAEgASgJUg5zdWJzY3JpcH'
    'Rpb25JZBIXCgd1c2VyX2lkGAIgASgJUgZ1c2VySWQSMQoEdGllchgDIAEoDjIdLm1hbnBhc2lr'
    'LnYxLlN1YnNjcmlwdGlvblRpZXJSBHRpZXISNwoGc3RhdHVzGAQgASgOMh8ubWFucGFzaWsudj'
    'EuU3Vic2NyaXB0aW9uU3RhdHVzUgZzdGF0dXMSOQoKc3RhcnRlZF9hdBgFIAEoCzIaLmdvb2ds'
    'ZS5wcm90b2J1Zi5UaW1lc3RhbXBSCXN0YXJ0ZWRBdBI5CgpleHBpcmVzX2F0GAYgASgLMhouZ2'
    '9vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFIJZXhwaXJlc0F0Ej0KDGNhbmNlbGxlZF9hdBgHIAEo'
    'CzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSC2NhbmNlbGxlZEF0Eh8KC21heF9kZXZpY2'
    'VzGAggASgFUgptYXhEZXZpY2VzEiwKEm1heF9mYW1pbHlfbWVtYmVycxgJIAEoBVIQbWF4RmFt'
    'aWx5TWVtYmVycxIuChNhaV9jb2FjaGluZ19lbmFibGVkGAogASgIUhFhaUNvYWNoaW5nRW5hYm'
    'xlZBIxChR0ZWxlbWVkaWNpbmVfZW5hYmxlZBgLIAEoCFITdGVsZW1lZGljaW5lRW5hYmxlZBIq'
    'ChFtb250aGx5X3ByaWNlX2tydxgMIAEoBVIPbW9udGhseVByaWNlS3J3Eh0KCmF1dG9fcmVuZX'
    'cYDSABKAhSCWF1dG9SZW5ldw==');

@$core.Deprecated('Use checkFeatureAccessRequestDescriptor instead')
const CheckFeatureAccessRequest$json = {
  '1': 'CheckFeatureAccessRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'feature_name', '3': 2, '4': 1, '5': 9, '10': 'featureName'},
  ],
};

/// Descriptor for `CheckFeatureAccessRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List checkFeatureAccessRequestDescriptor =
    $convert.base64Decode(
        'ChlDaGVja0ZlYXR1cmVBY2Nlc3NSZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZBIhCg'
        'xmZWF0dXJlX25hbWUYAiABKAlSC2ZlYXR1cmVOYW1l');

@$core.Deprecated('Use checkFeatureAccessResponseDescriptor instead')
const CheckFeatureAccessResponse$json = {
  '1': 'CheckFeatureAccessResponse',
  '2': [
    {'1': 'allowed', '3': 1, '4': 1, '5': 8, '10': 'allowed'},
    {
      '1': 'required_tier',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.SubscriptionTier',
      '10': 'requiredTier'
    },
    {
      '1': 'current_tier',
      '3': 3,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.SubscriptionTier',
      '10': 'currentTier'
    },
    {'1': 'message', '3': 4, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CheckFeatureAccessResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List checkFeatureAccessResponseDescriptor = $convert.base64Decode(
    'ChpDaGVja0ZlYXR1cmVBY2Nlc3NSZXNwb25zZRIYCgdhbGxvd2VkGAEgASgIUgdhbGxvd2VkEk'
    'IKDXJlcXVpcmVkX3RpZXIYAiABKA4yHS5tYW5wYXNpay52MS5TdWJzY3JpcHRpb25UaWVyUgxy'
    'ZXF1aXJlZFRpZXISQAoMY3VycmVudF90aWVyGAMgASgOMh0ubWFucGFzaWsudjEuU3Vic2NyaX'
    'B0aW9uVGllclILY3VycmVudFRpZXISGAoHbWVzc2FnZRgEIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use listSubscriptionPlansRequestDescriptor instead')
const ListSubscriptionPlansRequest$json = {
  '1': 'ListSubscriptionPlansRequest',
};

/// Descriptor for `ListSubscriptionPlansRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSubscriptionPlansRequestDescriptor =
    $convert.base64Decode('ChxMaXN0U3Vic2NyaXB0aW9uUGxhbnNSZXF1ZXN0');

@$core.Deprecated('Use listSubscriptionPlansResponseDescriptor instead')
const ListSubscriptionPlansResponse$json = {
  '1': 'ListSubscriptionPlansResponse',
  '2': [
    {
      '1': 'plans',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.SubscriptionPlan',
      '10': 'plans'
    },
  ],
};

/// Descriptor for `ListSubscriptionPlansResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSubscriptionPlansResponseDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0U3Vic2NyaXB0aW9uUGxhbnNSZXNwb25zZRIzCgVwbGFucxgBIAMoCzIdLm1hbnBhc2'
        'lrLnYxLlN1YnNjcmlwdGlvblBsYW5SBXBsYW5z');

@$core.Deprecated('Use subscriptionPlanDescriptor instead')
const SubscriptionPlan$json = {
  '1': 'SubscriptionPlan',
  '2': [
    {
      '1': 'tier',
      '3': 1,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.SubscriptionTier',
      '10': 'tier'
    },
    {'1': 'name', '3': 2, '4': 1, '5': 9, '10': 'name'},
    {'1': 'description', '3': 3, '4': 1, '5': 9, '10': 'description'},
    {'1': 'monthly_price_krw', '3': 4, '4': 1, '5': 5, '10': 'monthlyPriceKrw'},
    {'1': 'max_devices', '3': 5, '4': 1, '5': 5, '10': 'maxDevices'},
    {
      '1': 'max_family_members',
      '3': 6,
      '4': 1,
      '5': 5,
      '10': 'maxFamilyMembers'
    },
    {
      '1': 'ai_coaching_enabled',
      '3': 7,
      '4': 1,
      '5': 8,
      '10': 'aiCoachingEnabled'
    },
    {
      '1': 'telemedicine_enabled',
      '3': 8,
      '4': 1,
      '5': 8,
      '10': 'telemedicineEnabled'
    },
    {'1': 'features', '3': 9, '4': 3, '5': 9, '10': 'features'},
  ],
};

/// Descriptor for `SubscriptionPlan`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List subscriptionPlanDescriptor = $convert.base64Decode(
    'ChBTdWJzY3JpcHRpb25QbGFuEjEKBHRpZXIYASABKA4yHS5tYW5wYXNpay52MS5TdWJzY3JpcH'
    'Rpb25UaWVyUgR0aWVyEhIKBG5hbWUYAiABKAlSBG5hbWUSIAoLZGVzY3JpcHRpb24YAyABKAlS'
    'C2Rlc2NyaXB0aW9uEioKEW1vbnRobHlfcHJpY2Vfa3J3GAQgASgFUg9tb250aGx5UHJpY2VLcn'
    'cSHwoLbWF4X2RldmljZXMYBSABKAVSCm1heERldmljZXMSLAoSbWF4X2ZhbWlseV9tZW1iZXJz'
    'GAYgASgFUhBtYXhGYW1pbHlNZW1iZXJzEi4KE2FpX2NvYWNoaW5nX2VuYWJsZWQYByABKAhSEW'
    'FpQ29hY2hpbmdFbmFibGVkEjEKFHRlbGVtZWRpY2luZV9lbmFibGVkGAggASgIUhN0ZWxlbWVk'
    'aWNpbmVFbmFibGVkEhoKCGZlYXR1cmVzGAkgAygJUghmZWF0dXJlcw==');

@$core.Deprecated('Use listProductsRequestDescriptor instead')
const ListProductsRequest$json = {
  '1': 'ListProductsRequest',
  '2': [
    {
      '1': 'category',
      '3': 1,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.ProductCategory',
      '10': 'category'
    },
    {'1': 'limit', '3': 2, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 3, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `ListProductsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listProductsRequestDescriptor = $convert.base64Decode(
    'ChNMaXN0UHJvZHVjdHNSZXF1ZXN0EjgKCGNhdGVnb3J5GAEgASgOMhwubWFucGFzaWsudjEuUH'
    'JvZHVjdENhdGVnb3J5UghjYXRlZ29yeRIUCgVsaW1pdBgCIAEoBVIFbGltaXQSFgoGb2Zmc2V0'
    'GAMgASgFUgZvZmZzZXQ=');

@$core.Deprecated('Use listProductsResponseDescriptor instead')
const ListProductsResponse$json = {
  '1': 'ListProductsResponse',
  '2': [
    {
      '1': 'products',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.Product',
      '10': 'products'
    },
    {'1': 'total_count', '3': 2, '4': 1, '5': 5, '10': 'totalCount'},
  ],
};

/// Descriptor for `ListProductsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listProductsResponseDescriptor = $convert.base64Decode(
    'ChRMaXN0UHJvZHVjdHNSZXNwb25zZRIwCghwcm9kdWN0cxgBIAMoCzIULm1hbnBhc2lrLnYxLl'
    'Byb2R1Y3RSCHByb2R1Y3RzEh8KC3RvdGFsX2NvdW50GAIgASgFUgp0b3RhbENvdW50');

@$core.Deprecated('Use getProductRequestDescriptor instead')
const GetProductRequest$json = {
  '1': 'GetProductRequest',
  '2': [
    {'1': 'product_id', '3': 1, '4': 1, '5': 9, '10': 'productId'},
  ],
};

/// Descriptor for `GetProductRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getProductRequestDescriptor = $convert.base64Decode(
    'ChFHZXRQcm9kdWN0UmVxdWVzdBIdCgpwcm9kdWN0X2lkGAEgASgJUglwcm9kdWN0SWQ=');

@$core.Deprecated('Use productDescriptor instead')
const Product$json = {
  '1': 'Product',
  '2': [
    {'1': 'product_id', '3': 1, '4': 1, '5': 9, '10': 'productId'},
    {'1': 'name', '3': 2, '4': 1, '5': 9, '10': 'name'},
    {'1': 'description', '3': 3, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'category',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.ProductCategory',
      '10': 'category'
    },
    {'1': 'price_krw', '3': 5, '4': 1, '5': 5, '10': 'priceKrw'},
    {'1': 'stock', '3': 6, '4': 1, '5': 5, '10': 'stock'},
    {'1': 'image_url', '3': 7, '4': 1, '5': 9, '10': 'imageUrl'},
    {'1': 'is_active', '3': 8, '4': 1, '5': 8, '10': 'isActive'},
    {
      '1': 'created_at',
      '3': 9,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
  ],
};

/// Descriptor for `Product`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List productDescriptor = $convert.base64Decode(
    'CgdQcm9kdWN0Eh0KCnByb2R1Y3RfaWQYASABKAlSCXByb2R1Y3RJZBISCgRuYW1lGAIgASgJUg'
    'RuYW1lEiAKC2Rlc2NyaXB0aW9uGAMgASgJUgtkZXNjcmlwdGlvbhI4CghjYXRlZ29yeRgEIAEo'
    'DjIcLm1hbnBhc2lrLnYxLlByb2R1Y3RDYXRlZ29yeVIIY2F0ZWdvcnkSGwoJcHJpY2Vfa3J3GA'
    'UgASgFUghwcmljZUtydxIUCgVzdG9jaxgGIAEoBVIFc3RvY2sSGwoJaW1hZ2VfdXJsGAcgASgJ'
    'UghpbWFnZVVybBIbCglpc19hY3RpdmUYCCABKAhSCGlzQWN0aXZlEjkKCmNyZWF0ZWRfYXQYCS'
    'ABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUgljcmVhdGVkQXQ=');

@$core.Deprecated('Use addToCartRequestDescriptor instead')
const AddToCartRequest$json = {
  '1': 'AddToCartRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'product_id', '3': 2, '4': 1, '5': 9, '10': 'productId'},
    {'1': 'quantity', '3': 3, '4': 1, '5': 5, '10': 'quantity'},
  ],
};

/// Descriptor for `AddToCartRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addToCartRequestDescriptor = $convert.base64Decode(
    'ChBBZGRUb0NhcnRSZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZBIdCgpwcm9kdWN0X2'
    'lkGAIgASgJUglwcm9kdWN0SWQSGgoIcXVhbnRpdHkYAyABKAVSCHF1YW50aXR5');

@$core.Deprecated('Use getCartRequestDescriptor instead')
const GetCartRequest$json = {
  '1': 'GetCartRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `GetCartRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCartRequestDescriptor = $convert
    .base64Decode('Cg5HZXRDYXJ0UmVxdWVzdBIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQ=');

@$core.Deprecated('Use removeFromCartRequestDescriptor instead')
const RemoveFromCartRequest$json = {
  '1': 'RemoveFromCartRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'cart_item_id', '3': 2, '4': 1, '5': 9, '10': 'cartItemId'},
  ],
};

/// Descriptor for `RemoveFromCartRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removeFromCartRequestDescriptor = $convert.base64Decode(
    'ChVSZW1vdmVGcm9tQ2FydFJlcXVlc3QSFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklkEiAKDGNhcn'
    'RfaXRlbV9pZBgCIAEoCVIKY2FydEl0ZW1JZA==');

@$core.Deprecated('Use cartDescriptor instead')
const Cart$json = {
  '1': 'Cart',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'items',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.CartItem',
      '10': 'items'
    },
    {'1': 'total_price_krw', '3': 3, '4': 1, '5': 5, '10': 'totalPriceKrw'},
  ],
};

/// Descriptor for `Cart`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cartDescriptor = $convert.base64Decode(
    'CgRDYXJ0EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZBIrCgVpdGVtcxgCIAMoCzIVLm1hbnBhc2'
    'lrLnYxLkNhcnRJdGVtUgVpdGVtcxImCg90b3RhbF9wcmljZV9rcncYAyABKAVSDXRvdGFsUHJp'
    'Y2VLcnc=');

@$core.Deprecated('Use cartItemDescriptor instead')
const CartItem$json = {
  '1': 'CartItem',
  '2': [
    {'1': 'cart_item_id', '3': 1, '4': 1, '5': 9, '10': 'cartItemId'},
    {'1': 'product_id', '3': 2, '4': 1, '5': 9, '10': 'productId'},
    {'1': 'product_name', '3': 3, '4': 1, '5': 9, '10': 'productName'},
    {'1': 'quantity', '3': 4, '4': 1, '5': 5, '10': 'quantity'},
    {'1': 'unit_price_krw', '3': 5, '4': 1, '5': 5, '10': 'unitPriceKrw'},
    {'1': 'total_price_krw', '3': 6, '4': 1, '5': 5, '10': 'totalPriceKrw'},
  ],
};

/// Descriptor for `CartItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cartItemDescriptor = $convert.base64Decode(
    'CghDYXJ0SXRlbRIgCgxjYXJ0X2l0ZW1faWQYASABKAlSCmNhcnRJdGVtSWQSHQoKcHJvZHVjdF'
    '9pZBgCIAEoCVIJcHJvZHVjdElkEiEKDHByb2R1Y3RfbmFtZRgDIAEoCVILcHJvZHVjdE5hbWUS'
    'GgoIcXVhbnRpdHkYBCABKAVSCHF1YW50aXR5EiQKDnVuaXRfcHJpY2Vfa3J3GAUgASgFUgx1bm'
    'l0UHJpY2VLcncSJgoPdG90YWxfcHJpY2Vfa3J3GAYgASgFUg10b3RhbFByaWNlS3J3');

@$core.Deprecated('Use createOrderRequestDescriptor instead')
const CreateOrderRequest$json = {
  '1': 'CreateOrderRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'shipping_address', '3': 2, '4': 1, '5': 9, '10': 'shippingAddress'},
    {'1': 'payment_method', '3': 3, '4': 1, '5': 9, '10': 'paymentMethod'},
  ],
};

/// Descriptor for `CreateOrderRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createOrderRequestDescriptor = $convert.base64Decode(
    'ChJDcmVhdGVPcmRlclJlcXVlc3QSFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklkEikKEHNoaXBwaW'
    '5nX2FkZHJlc3MYAiABKAlSD3NoaXBwaW5nQWRkcmVzcxIlCg5wYXltZW50X21ldGhvZBgDIAEo'
    'CVINcGF5bWVudE1ldGhvZA==');

@$core.Deprecated('Use getOrderRequestDescriptor instead')
const GetOrderRequest$json = {
  '1': 'GetOrderRequest',
  '2': [
    {'1': 'order_id', '3': 1, '4': 1, '5': 9, '10': 'orderId'},
  ],
};

/// Descriptor for `GetOrderRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getOrderRequestDescriptor = $convert.base64Decode(
    'Cg9HZXRPcmRlclJlcXVlc3QSGQoIb3JkZXJfaWQYASABKAlSB29yZGVySWQ=');

@$core.Deprecated('Use listOrdersRequestDescriptor instead')
const ListOrdersRequest$json = {
  '1': 'ListOrdersRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'limit', '3': 2, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 3, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `ListOrdersRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listOrdersRequestDescriptor = $convert.base64Decode(
    'ChFMaXN0T3JkZXJzUmVxdWVzdBIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQSFAoFbGltaXQYAi'
    'ABKAVSBWxpbWl0EhYKBm9mZnNldBgDIAEoBVIGb2Zmc2V0');

@$core.Deprecated('Use listOrdersResponseDescriptor instead')
const ListOrdersResponse$json = {
  '1': 'ListOrdersResponse',
  '2': [
    {
      '1': 'orders',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.Order',
      '10': 'orders'
    },
    {'1': 'total_count', '3': 2, '4': 1, '5': 5, '10': 'totalCount'},
  ],
};

/// Descriptor for `ListOrdersResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listOrdersResponseDescriptor = $convert.base64Decode(
    'ChJMaXN0T3JkZXJzUmVzcG9uc2USKgoGb3JkZXJzGAEgAygLMhIubWFucGFzaWsudjEuT3JkZX'
    'JSBm9yZGVycxIfCgt0b3RhbF9jb3VudBgCIAEoBVIKdG90YWxDb3VudA==');

@$core.Deprecated('Use orderDescriptor instead')
const Order$json = {
  '1': 'Order',
  '2': [
    {'1': 'order_id', '3': 1, '4': 1, '5': 9, '10': 'orderId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'items',
      '3': 3,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.OrderItem',
      '10': 'items'
    },
    {'1': 'total_price_krw', '3': 4, '4': 1, '5': 5, '10': 'totalPriceKrw'},
    {
      '1': 'status',
      '3': 5,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.OrderStatus',
      '10': 'status'
    },
    {'1': 'shipping_address', '3': 6, '4': 1, '5': 9, '10': 'shippingAddress'},
    {'1': 'payment_id', '3': 7, '4': 1, '5': 9, '10': 'paymentId'},
    {
      '1': 'created_at',
      '3': 8,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
    {
      '1': 'updated_at',
      '3': 9,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'updatedAt'
    },
  ],
};

/// Descriptor for `Order`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List orderDescriptor = $convert.base64Decode(
    'CgVPcmRlchIZCghvcmRlcl9pZBgBIAEoCVIHb3JkZXJJZBIXCgd1c2VyX2lkGAIgASgJUgZ1c2'
    'VySWQSLAoFaXRlbXMYAyADKAsyFi5tYW5wYXNpay52MS5PcmRlckl0ZW1SBWl0ZW1zEiYKD3Rv'
    'dGFsX3ByaWNlX2tydxgEIAEoBVINdG90YWxQcmljZUtydxIwCgZzdGF0dXMYBSABKA4yGC5tYW'
    '5wYXNpay52MS5PcmRlclN0YXR1c1IGc3RhdHVzEikKEHNoaXBwaW5nX2FkZHJlc3MYBiABKAlS'
    'D3NoaXBwaW5nQWRkcmVzcxIdCgpwYXltZW50X2lkGAcgASgJUglwYXltZW50SWQSOQoKY3JlYX'
    'RlZF9hdBgIIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSCWNyZWF0ZWRBdBI5Cgp1'
    'cGRhdGVkX2F0GAkgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFIJdXBkYXRlZEF0');

@$core.Deprecated('Use orderItemDescriptor instead')
const OrderItem$json = {
  '1': 'OrderItem',
  '2': [
    {'1': 'product_id', '3': 1, '4': 1, '5': 9, '10': 'productId'},
    {'1': 'product_name', '3': 2, '4': 1, '5': 9, '10': 'productName'},
    {'1': 'quantity', '3': 3, '4': 1, '5': 5, '10': 'quantity'},
    {'1': 'unit_price_krw', '3': 4, '4': 1, '5': 5, '10': 'unitPriceKrw'},
    {'1': 'total_price_krw', '3': 5, '4': 1, '5': 5, '10': 'totalPriceKrw'},
  ],
};

/// Descriptor for `OrderItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List orderItemDescriptor = $convert.base64Decode(
    'CglPcmRlckl0ZW0SHQoKcHJvZHVjdF9pZBgBIAEoCVIJcHJvZHVjdElkEiEKDHByb2R1Y3Rfbm'
    'FtZRgCIAEoCVILcHJvZHVjdE5hbWUSGgoIcXVhbnRpdHkYAyABKAVSCHF1YW50aXR5EiQKDnVu'
    'aXRfcHJpY2Vfa3J3GAQgASgFUgx1bml0UHJpY2VLcncSJgoPdG90YWxfcHJpY2Vfa3J3GAUgAS'
    'gFUg10b3RhbFByaWNlS3J3');

@$core.Deprecated('Use createPaymentRequestDescriptor instead')
const CreatePaymentRequest$json = {
  '1': 'CreatePaymentRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'order_id', '3': 2, '4': 1, '5': 9, '10': 'orderId'},
    {'1': 'subscription_id', '3': 3, '4': 1, '5': 9, '10': 'subscriptionId'},
    {
      '1': 'payment_type',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.PaymentType',
      '10': 'paymentType'
    },
    {'1': 'amount_krw', '3': 5, '4': 1, '5': 5, '10': 'amountKrw'},
    {'1': 'payment_method', '3': 6, '4': 1, '5': 9, '10': 'paymentMethod'},
  ],
};

/// Descriptor for `CreatePaymentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createPaymentRequestDescriptor = $convert.base64Decode(
    'ChRDcmVhdGVQYXltZW50UmVxdWVzdBIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQSGQoIb3JkZX'
    'JfaWQYAiABKAlSB29yZGVySWQSJwoPc3Vic2NyaXB0aW9uX2lkGAMgASgJUg5zdWJzY3JpcHRp'
    'b25JZBI7CgxwYXltZW50X3R5cGUYBCABKA4yGC5tYW5wYXNpay52MS5QYXltZW50VHlwZVILcG'
    'F5bWVudFR5cGUSHQoKYW1vdW50X2tydxgFIAEoBVIJYW1vdW50S3J3EiUKDnBheW1lbnRfbWV0'
    'aG9kGAYgASgJUg1wYXltZW50TWV0aG9k');

@$core.Deprecated('Use confirmPaymentRequestDescriptor instead')
const ConfirmPaymentRequest$json = {
  '1': 'ConfirmPaymentRequest',
  '2': [
    {'1': 'payment_id', '3': 1, '4': 1, '5': 9, '10': 'paymentId'},
    {'1': 'pg_transaction_id', '3': 2, '4': 1, '5': 9, '10': 'pgTransactionId'},
    {'1': 'pg_provider', '3': 3, '4': 1, '5': 9, '10': 'pgProvider'},
    {'1': 'payment_key', '3': 4, '4': 1, '5': 9, '10': 'paymentKey'},
  ],
};

/// Descriptor for `ConfirmPaymentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List confirmPaymentRequestDescriptor = $convert.base64Decode(
    'ChVDb25maXJtUGF5bWVudFJlcXVlc3QSHQoKcGF5bWVudF9pZBgBIAEoCVIJcGF5bWVudElkEi'
    'oKEXBnX3RyYW5zYWN0aW9uX2lkGAIgASgJUg9wZ1RyYW5zYWN0aW9uSWQSHwoLcGdfcHJvdmlk'
    'ZXIYAyABKAlSCnBnUHJvdmlkZXISHwoLcGF5bWVudF9rZXkYBCABKAlSCnBheW1lbnRLZXk=');

@$core.Deprecated('Use getPaymentRequestDescriptor instead')
const GetPaymentRequest$json = {
  '1': 'GetPaymentRequest',
  '2': [
    {'1': 'payment_id', '3': 1, '4': 1, '5': 9, '10': 'paymentId'},
  ],
};

/// Descriptor for `GetPaymentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPaymentRequestDescriptor = $convert.base64Decode(
    'ChFHZXRQYXltZW50UmVxdWVzdBIdCgpwYXltZW50X2lkGAEgASgJUglwYXltZW50SWQ=');

@$core.Deprecated('Use listPaymentsRequestDescriptor instead')
const ListPaymentsRequest$json = {
  '1': 'ListPaymentsRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'limit', '3': 2, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 3, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `ListPaymentsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPaymentsRequestDescriptor = $convert.base64Decode(
    'ChNMaXN0UGF5bWVudHNSZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZBIUCgVsaW1pdB'
    'gCIAEoBVIFbGltaXQSFgoGb2Zmc2V0GAMgASgFUgZvZmZzZXQ=');

@$core.Deprecated('Use listPaymentsResponseDescriptor instead')
const ListPaymentsResponse$json = {
  '1': 'ListPaymentsResponse',
  '2': [
    {
      '1': 'payments',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.PaymentDetail',
      '10': 'payments'
    },
    {'1': 'total_count', '3': 2, '4': 1, '5': 5, '10': 'totalCount'},
  ],
};

/// Descriptor for `ListPaymentsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPaymentsResponseDescriptor = $convert.base64Decode(
    'ChRMaXN0UGF5bWVudHNSZXNwb25zZRI2CghwYXltZW50cxgBIAMoCzIaLm1hbnBhc2lrLnYxLl'
    'BheW1lbnREZXRhaWxSCHBheW1lbnRzEh8KC3RvdGFsX2NvdW50GAIgASgFUgp0b3RhbENvdW50');

@$core.Deprecated('Use paymentDetailDescriptor instead')
const PaymentDetail$json = {
  '1': 'PaymentDetail',
  '2': [
    {'1': 'payment_id', '3': 1, '4': 1, '5': 9, '10': 'paymentId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'order_id', '3': 3, '4': 1, '5': 9, '10': 'orderId'},
    {'1': 'subscription_id', '3': 4, '4': 1, '5': 9, '10': 'subscriptionId'},
    {
      '1': 'payment_type',
      '3': 5,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.PaymentType',
      '10': 'paymentType'
    },
    {'1': 'amount_krw', '3': 6, '4': 1, '5': 5, '10': 'amountKrw'},
    {
      '1': 'status',
      '3': 7,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.PaymentStatus',
      '10': 'status'
    },
    {'1': 'payment_method', '3': 8, '4': 1, '5': 9, '10': 'paymentMethod'},
    {'1': 'pg_transaction_id', '3': 9, '4': 1, '5': 9, '10': 'pgTransactionId'},
    {'1': 'pg_provider', '3': 10, '4': 1, '5': 9, '10': 'pgProvider'},
    {
      '1': 'created_at',
      '3': 11,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
    {
      '1': 'completed_at',
      '3': 12,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'completedAt'
    },
  ],
};

/// Descriptor for `PaymentDetail`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List paymentDetailDescriptor = $convert.base64Decode(
    'Cg1QYXltZW50RGV0YWlsEh0KCnBheW1lbnRfaWQYASABKAlSCXBheW1lbnRJZBIXCgd1c2VyX2'
    'lkGAIgASgJUgZ1c2VySWQSGQoIb3JkZXJfaWQYAyABKAlSB29yZGVySWQSJwoPc3Vic2NyaXB0'
    'aW9uX2lkGAQgASgJUg5zdWJzY3JpcHRpb25JZBI7CgxwYXltZW50X3R5cGUYBSABKA4yGC5tYW'
    '5wYXNpay52MS5QYXltZW50VHlwZVILcGF5bWVudFR5cGUSHQoKYW1vdW50X2tydxgGIAEoBVIJ'
    'YW1vdW50S3J3EjIKBnN0YXR1cxgHIAEoDjIaLm1hbnBhc2lrLnYxLlBheW1lbnRTdGF0dXNSBn'
    'N0YXR1cxIlCg5wYXltZW50X21ldGhvZBgIIAEoCVINcGF5bWVudE1ldGhvZBIqChFwZ190cmFu'
    'c2FjdGlvbl9pZBgJIAEoCVIPcGdUcmFuc2FjdGlvbklkEh8KC3BnX3Byb3ZpZGVyGAogASgJUg'
    'pwZ1Byb3ZpZGVyEjkKCmNyZWF0ZWRfYXQYCyABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0'
    'YW1wUgljcmVhdGVkQXQSPQoMY29tcGxldGVkX2F0GAwgASgLMhouZ29vZ2xlLnByb3RvYnVmLl'
    'RpbWVzdGFtcFILY29tcGxldGVkQXQ=');

@$core.Deprecated('Use refundPaymentRequestDescriptor instead')
const RefundPaymentRequest$json = {
  '1': 'RefundPaymentRequest',
  '2': [
    {'1': 'payment_id', '3': 1, '4': 1, '5': 9, '10': 'paymentId'},
    {'1': 'refund_amount_krw', '3': 2, '4': 1, '5': 5, '10': 'refundAmountKrw'},
    {'1': 'reason', '3': 3, '4': 1, '5': 9, '10': 'reason'},
  ],
};

/// Descriptor for `RefundPaymentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List refundPaymentRequestDescriptor = $convert.base64Decode(
    'ChRSZWZ1bmRQYXltZW50UmVxdWVzdBIdCgpwYXltZW50X2lkGAEgASgJUglwYXltZW50SWQSKg'
    'oRcmVmdW5kX2Ftb3VudF9rcncYAiABKAVSD3JlZnVuZEFtb3VudEtydxIWCgZyZWFzb24YAyAB'
    'KAlSBnJlYXNvbg==');

@$core.Deprecated('Use refundResponseDescriptor instead')
const RefundResponse$json = {
  '1': 'RefundResponse',
  '2': [
    {'1': 'refund_id', '3': 1, '4': 1, '5': 9, '10': 'refundId'},
    {'1': 'payment_id', '3': 2, '4': 1, '5': 9, '10': 'paymentId'},
    {'1': 'refund_amount_krw', '3': 3, '4': 1, '5': 5, '10': 'refundAmountKrw'},
    {
      '1': 'payment_status',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.PaymentStatus',
      '10': 'paymentStatus'
    },
    {
      '1': 'refunded_at',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'refundedAt'
    },
  ],
};

/// Descriptor for `RefundResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List refundResponseDescriptor = $convert.base64Decode(
    'Cg5SZWZ1bmRSZXNwb25zZRIbCglyZWZ1bmRfaWQYASABKAlSCHJlZnVuZElkEh0KCnBheW1lbn'
    'RfaWQYAiABKAlSCXBheW1lbnRJZBIqChFyZWZ1bmRfYW1vdW50X2tydxgDIAEoBVIPcmVmdW5k'
    'QW1vdW50S3J3EkEKDnBheW1lbnRfc3RhdHVzGAQgASgOMhoubWFucGFzaWsudjEuUGF5bWVudF'
    'N0YXR1c1INcGF5bWVudFN0YXR1cxI7CgtyZWZ1bmRlZF9hdBgFIAEoCzIaLmdvb2dsZS5wcm90'
    'b2J1Zi5UaW1lc3RhbXBSCnJlZnVuZGVkQXQ=');

@$core.Deprecated('Use analyzeMeasurementRequestDescriptor instead')
const AnalyzeMeasurementRequest$json = {
  '1': 'AnalyzeMeasurementRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'measurement_id', '3': 2, '4': 1, '5': 9, '10': 'measurementId'},
    {
      '1': 'models',
      '3': 3,
      '4': 3,
      '5': 14,
      '6': '.manpasik.v1.AiModelType',
      '10': 'models'
    },
  ],
};

/// Descriptor for `AnalyzeMeasurementRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List analyzeMeasurementRequestDescriptor = $convert.base64Decode(
    'ChlBbmFseXplTWVhc3VyZW1lbnRSZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZBIlCg'
    '5tZWFzdXJlbWVudF9pZBgCIAEoCVINbWVhc3VyZW1lbnRJZBIwCgZtb2RlbHMYAyADKA4yGC5t'
    'YW5wYXNpay52MS5BaU1vZGVsVHlwZVIGbW9kZWxz');

@$core.Deprecated('Use biomarkerResultDescriptor instead')
const BiomarkerResult$json = {
  '1': 'BiomarkerResult',
  '2': [
    {'1': 'biomarker_name', '3': 1, '4': 1, '5': 9, '10': 'biomarkerName'},
    {'1': 'value', '3': 2, '4': 1, '5': 1, '10': 'value'},
    {'1': 'unit', '3': 3, '4': 1, '5': 9, '10': 'unit'},
    {'1': 'classification', '3': 4, '4': 1, '5': 9, '10': 'classification'},
    {'1': 'confidence', '3': 5, '4': 1, '5': 1, '10': 'confidence'},
    {
      '1': 'risk_level',
      '3': 6,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.RiskLevel',
      '10': 'riskLevel'
    },
    {'1': 'reference_range', '3': 7, '4': 1, '5': 9, '10': 'referenceRange'},
  ],
};

/// Descriptor for `BiomarkerResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List biomarkerResultDescriptor = $convert.base64Decode(
    'Cg9CaW9tYXJrZXJSZXN1bHQSJQoOYmlvbWFya2VyX25hbWUYASABKAlSDWJpb21hcmtlck5hbW'
    'USFAoFdmFsdWUYAiABKAFSBXZhbHVlEhIKBHVuaXQYAyABKAlSBHVuaXQSJgoOY2xhc3NpZmlj'
    'YXRpb24YBCABKAlSDmNsYXNzaWZpY2F0aW9uEh4KCmNvbmZpZGVuY2UYBSABKAFSCmNvbmZpZG'
    'VuY2USNQoKcmlza19sZXZlbBgGIAEoDjIWLm1hbnBhc2lrLnYxLlJpc2tMZXZlbFIJcmlza0xl'
    'dmVsEicKD3JlZmVyZW5jZV9yYW5nZRgHIAEoCVIOcmVmZXJlbmNlUmFuZ2U=');

@$core.Deprecated('Use anomalyFlagDescriptor instead')
const AnomalyFlag$json = {
  '1': 'AnomalyFlag',
  '2': [
    {'1': 'metric_name', '3': 1, '4': 1, '5': 9, '10': 'metricName'},
    {'1': 'value', '3': 2, '4': 1, '5': 1, '10': 'value'},
    {'1': 'expected_min', '3': 3, '4': 1, '5': 1, '10': 'expectedMin'},
    {'1': 'expected_max', '3': 4, '4': 1, '5': 1, '10': 'expectedMax'},
    {'1': 'anomaly_score', '3': 5, '4': 1, '5': 1, '10': 'anomalyScore'},
    {'1': 'description', '3': 6, '4': 1, '5': 9, '10': 'description'},
  ],
};

/// Descriptor for `AnomalyFlag`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List anomalyFlagDescriptor = $convert.base64Decode(
    'CgtBbm9tYWx5RmxhZxIfCgttZXRyaWNfbmFtZRgBIAEoCVIKbWV0cmljTmFtZRIUCgV2YWx1ZR'
    'gCIAEoAVIFdmFsdWUSIQoMZXhwZWN0ZWRfbWluGAMgASgBUgtleHBlY3RlZE1pbhIhCgxleHBl'
    'Y3RlZF9tYXgYBCABKAFSC2V4cGVjdGVkTWF4EiMKDWFub21hbHlfc2NvcmUYBSABKAFSDGFub2'
    '1hbHlTY29yZRIgCgtkZXNjcmlwdGlvbhgGIAEoCVILZGVzY3JpcHRpb24=');

@$core.Deprecated('Use analysisResultDescriptor instead')
const AnalysisResult$json = {
  '1': 'AnalysisResult',
  '2': [
    {'1': 'analysis_id', '3': 1, '4': 1, '5': 9, '10': 'analysisId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'measurement_id', '3': 3, '4': 1, '5': 9, '10': 'measurementId'},
    {
      '1': 'biomarkers',
      '3': 4,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.BiomarkerResult',
      '10': 'biomarkers'
    },
    {
      '1': 'anomalies',
      '3': 5,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.AnomalyFlag',
      '10': 'anomalies'
    },
    {
      '1': 'overall_health_score',
      '3': 6,
      '4': 1,
      '5': 1,
      '10': 'overallHealthScore'
    },
    {'1': 'summary', '3': 7, '4': 1, '5': 9, '10': 'summary'},
    {
      '1': 'analyzed_at',
      '3': 8,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'analyzedAt'
    },
  ],
};

/// Descriptor for `AnalysisResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List analysisResultDescriptor = $convert.base64Decode(
    'Cg5BbmFseXNpc1Jlc3VsdBIfCgthbmFseXNpc19pZBgBIAEoCVIKYW5hbHlzaXNJZBIXCgd1c2'
    'VyX2lkGAIgASgJUgZ1c2VySWQSJQoObWVhc3VyZW1lbnRfaWQYAyABKAlSDW1lYXN1cmVtZW50'
    'SWQSPAoKYmlvbWFya2VycxgEIAMoCzIcLm1hbnBhc2lrLnYxLkJpb21hcmtlclJlc3VsdFIKYm'
    'lvbWFya2VycxI2Cglhbm9tYWxpZXMYBSADKAsyGC5tYW5wYXNpay52MS5Bbm9tYWx5RmxhZ1IJ'
    'YW5vbWFsaWVzEjAKFG92ZXJhbGxfaGVhbHRoX3Njb3JlGAYgASgBUhJvdmVyYWxsSGVhbHRoU2'
    'NvcmUSGAoHc3VtbWFyeRgHIAEoCVIHc3VtbWFyeRI7CgthbmFseXplZF9hdBgIIAEoCzIaLmdv'
    'b2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSCmFuYWx5emVkQXQ=');

@$core.Deprecated('Use getHealthScoreRequestDescriptor instead')
const GetHealthScoreRequest$json = {
  '1': 'GetHealthScoreRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'days', '3': 2, '4': 1, '5': 5, '10': 'days'},
  ],
};

/// Descriptor for `GetHealthScoreRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getHealthScoreRequestDescriptor = $convert.base64Decode(
    'ChVHZXRIZWFsdGhTY29yZVJlcXVlc3QSFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklkEhIKBGRheX'
    'MYAiABKAVSBGRheXM=');

@$core.Deprecated('Use healthScoreResponseDescriptor instead')
const HealthScoreResponse$json = {
  '1': 'HealthScoreResponse',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'overall_score', '3': 2, '4': 1, '5': 1, '10': 'overallScore'},
    {
      '1': 'category_scores',
      '3': 3,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.HealthScoreResponse.CategoryScoresEntry',
      '10': 'categoryScores'
    },
    {'1': 'trend', '3': 4, '4': 1, '5': 9, '10': 'trend'},
    {'1': 'recommendation', '3': 5, '4': 1, '5': 9, '10': 'recommendation'},
    {
      '1': 'calculated_at',
      '3': 6,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'calculatedAt'
    },
  ],
  '3': [HealthScoreResponse_CategoryScoresEntry$json],
};

@$core.Deprecated('Use healthScoreResponseDescriptor instead')
const HealthScoreResponse_CategoryScoresEntry$json = {
  '1': 'CategoryScoresEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 1, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `HealthScoreResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List healthScoreResponseDescriptor = $convert.base64Decode(
    'ChNIZWFsdGhTY29yZVJlc3BvbnNlEhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZBIjCg1vdmVyYW'
    'xsX3Njb3JlGAIgASgBUgxvdmVyYWxsU2NvcmUSXQoPY2F0ZWdvcnlfc2NvcmVzGAMgAygLMjQu'
    'bWFucGFzaWsudjEuSGVhbHRoU2NvcmVSZXNwb25zZS5DYXRlZ29yeVNjb3Jlc0VudHJ5Ug5jYX'
    'RlZ29yeVNjb3JlcxIUCgV0cmVuZBgEIAEoCVIFdHJlbmQSJgoOcmVjb21tZW5kYXRpb24YBSAB'
    'KAlSDnJlY29tbWVuZGF0aW9uEj8KDWNhbGN1bGF0ZWRfYXQYBiABKAsyGi5nb29nbGUucHJvdG'
    '9idWYuVGltZXN0YW1wUgxjYWxjdWxhdGVkQXQaQQoTQ2F0ZWdvcnlTY29yZXNFbnRyeRIQCgNr'
    'ZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoAVIFdmFsdWU6AjgB');

@$core.Deprecated('Use predictTrendRequestDescriptor instead')
const PredictTrendRequest$json = {
  '1': 'PredictTrendRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'metric_name', '3': 2, '4': 1, '5': 9, '10': 'metricName'},
    {'1': 'history_days', '3': 3, '4': 1, '5': 5, '10': 'historyDays'},
    {'1': 'prediction_days', '3': 4, '4': 1, '5': 5, '10': 'predictionDays'},
  ],
};

/// Descriptor for `PredictTrendRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List predictTrendRequestDescriptor = $convert.base64Decode(
    'ChNQcmVkaWN0VHJlbmRSZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZBIfCgttZXRyaW'
    'NfbmFtZRgCIAEoCVIKbWV0cmljTmFtZRIhCgxoaXN0b3J5X2RheXMYAyABKAVSC2hpc3RvcnlE'
    'YXlzEicKD3ByZWRpY3Rpb25fZGF5cxgEIAEoBVIOcHJlZGljdGlvbkRheXM=');

@$core.Deprecated('Use trendDataPointDescriptor instead')
const TrendDataPoint$json = {
  '1': 'TrendDataPoint',
  '2': [
    {
      '1': 'timestamp',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'timestamp'
    },
    {'1': 'value', '3': 2, '4': 1, '5': 1, '10': 'value'},
    {'1': 'lower_bound', '3': 3, '4': 1, '5': 1, '10': 'lowerBound'},
    {'1': 'upper_bound', '3': 4, '4': 1, '5': 1, '10': 'upperBound'},
  ],
};

/// Descriptor for `TrendDataPoint`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List trendDataPointDescriptor = $convert.base64Decode(
    'Cg5UcmVuZERhdGFQb2ludBI4Cgl0aW1lc3RhbXAYASABKAsyGi5nb29nbGUucHJvdG9idWYuVG'
    'ltZXN0YW1wUgl0aW1lc3RhbXASFAoFdmFsdWUYAiABKAFSBXZhbHVlEh8KC2xvd2VyX2JvdW5k'
    'GAMgASgBUgpsb3dlckJvdW5kEh8KC3VwcGVyX2JvdW5kGAQgASgBUgp1cHBlckJvdW5k');

@$core.Deprecated('Use trendPredictionDescriptor instead')
const TrendPrediction$json = {
  '1': 'TrendPrediction',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'metric_name', '3': 2, '4': 1, '5': 9, '10': 'metricName'},
    {
      '1': 'historical',
      '3': 3,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.TrendDataPoint',
      '10': 'historical'
    },
    {
      '1': 'predicted',
      '3': 4,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.TrendDataPoint',
      '10': 'predicted'
    },
    {'1': 'confidence', '3': 5, '4': 1, '5': 1, '10': 'confidence'},
    {'1': 'direction', '3': 6, '4': 1, '5': 9, '10': 'direction'},
    {'1': 'insight', '3': 7, '4': 1, '5': 9, '10': 'insight'},
  ],
};

/// Descriptor for `TrendPrediction`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List trendPredictionDescriptor = $convert.base64Decode(
    'Cg9UcmVuZFByZWRpY3Rpb24SFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklkEh8KC21ldHJpY19uYW'
    '1lGAIgASgJUgptZXRyaWNOYW1lEjsKCmhpc3RvcmljYWwYAyADKAsyGy5tYW5wYXNpay52MS5U'
    'cmVuZERhdGFQb2ludFIKaGlzdG9yaWNhbBI5CglwcmVkaWN0ZWQYBCADKAsyGy5tYW5wYXNpay'
    '52MS5UcmVuZERhdGFQb2ludFIJcHJlZGljdGVkEh4KCmNvbmZpZGVuY2UYBSABKAFSCmNvbmZp'
    'ZGVuY2USHAoJZGlyZWN0aW9uGAYgASgJUglkaXJlY3Rpb24SGAoHaW5zaWdodBgHIAEoCVIHaW'
    '5zaWdodA==');

@$core.Deprecated('Use getModelInfoRequestDescriptor instead')
const GetModelInfoRequest$json = {
  '1': 'GetModelInfoRequest',
  '2': [
    {
      '1': 'model_type',
      '3': 1,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.AiModelType',
      '10': 'modelType'
    },
  ],
};

/// Descriptor for `GetModelInfoRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getModelInfoRequestDescriptor = $convert.base64Decode(
    'ChNHZXRNb2RlbEluZm9SZXF1ZXN0EjcKCm1vZGVsX3R5cGUYASABKA4yGC5tYW5wYXNpay52MS'
    '5BaU1vZGVsVHlwZVIJbW9kZWxUeXBl');

@$core.Deprecated('Use modelInfoDescriptor instead')
const ModelInfo$json = {
  '1': 'ModelInfo',
  '2': [
    {
      '1': 'model_type',
      '3': 1,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.AiModelType',
      '10': 'modelType'
    },
    {'1': 'name', '3': 2, '4': 1, '5': 9, '10': 'name'},
    {'1': 'version', '3': 3, '4': 1, '5': 9, '10': 'version'},
    {'1': 'description', '3': 4, '4': 1, '5': 9, '10': 'description'},
    {'1': 'accuracy', '3': 5, '4': 1, '5': 1, '10': 'accuracy'},
    {
      '1': 'last_trained',
      '3': 6,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'lastTrained'
    },
    {'1': 'status', '3': 7, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `ModelInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List modelInfoDescriptor = $convert.base64Decode(
    'CglNb2RlbEluZm8SNwoKbW9kZWxfdHlwZRgBIAEoDjIYLm1hbnBhc2lrLnYxLkFpTW9kZWxUeX'
    'BlUgltb2RlbFR5cGUSEgoEbmFtZRgCIAEoCVIEbmFtZRIYCgd2ZXJzaW9uGAMgASgJUgd2ZXJz'
    'aW9uEiAKC2Rlc2NyaXB0aW9uGAQgASgJUgtkZXNjcmlwdGlvbhIaCghhY2N1cmFjeRgFIAEoAV'
    'IIYWNjdXJhY3kSPQoMbGFzdF90cmFpbmVkGAYgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVz'
    'dGFtcFILbGFzdFRyYWluZWQSFgoGc3RhdHVzGAcgASgJUgZzdGF0dXM=');

@$core.Deprecated('Use listModelsRequestDescriptor instead')
const ListModelsRequest$json = {
  '1': 'ListModelsRequest',
};

/// Descriptor for `ListModelsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listModelsRequestDescriptor =
    $convert.base64Decode('ChFMaXN0TW9kZWxzUmVxdWVzdA==');

@$core.Deprecated('Use listModelsResponseDescriptor instead')
const ListModelsResponse$json = {
  '1': 'ListModelsResponse',
  '2': [
    {
      '1': 'models',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.ModelInfo',
      '10': 'models'
    },
  ],
};

/// Descriptor for `ListModelsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listModelsResponseDescriptor = $convert.base64Decode(
    'ChJMaXN0TW9kZWxzUmVzcG9uc2USLgoGbW9kZWxzGAEgAygLMhYubWFucGFzaWsudjEuTW9kZW'
    'xJbmZvUgZtb2RlbHM=');

@$core.Deprecated('Use readCartridgeRequestDescriptor instead')
const ReadCartridgeRequest$json = {
  '1': 'ReadCartridgeRequest',
  '2': [
    {'1': 'nfc_tag_data', '3': 1, '4': 1, '5': 12, '10': 'nfcTagData'},
    {'1': 'tag_version', '3': 2, '4': 1, '5': 5, '10': 'tagVersion'},
  ],
};

/// Descriptor for `ReadCartridgeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List readCartridgeRequestDescriptor = $convert.base64Decode(
    'ChRSZWFkQ2FydHJpZGdlUmVxdWVzdBIgCgxuZmNfdGFnX2RhdGEYASABKAxSCm5mY1RhZ0RhdG'
    'ESHwoLdGFnX3ZlcnNpb24YAiABKAVSCnRhZ1ZlcnNpb24=');

@$core.Deprecated('Use cartridgeDetailDescriptor instead')
const CartridgeDetail$json = {
  '1': 'CartridgeDetail',
  '2': [
    {'1': 'cartridge_uid', '3': 1, '4': 1, '5': 9, '10': 'cartridgeUid'},
    {'1': 'category_code', '3': 2, '4': 1, '5': 5, '10': 'categoryCode'},
    {'1': 'type_index', '3': 3, '4': 1, '5': 5, '10': 'typeIndex'},
    {'1': 'legacy_code', '3': 4, '4': 1, '5': 5, '10': 'legacyCode'},
    {'1': 'name_ko', '3': 5, '4': 1, '5': 9, '10': 'nameKo'},
    {'1': 'name_en', '3': 6, '4': 1, '5': 9, '10': 'nameEn'},
    {'1': 'lot_id', '3': 7, '4': 1, '5': 9, '10': 'lotId'},
    {'1': 'expiry_date', '3': 8, '4': 1, '5': 9, '10': 'expiryDate'},
    {'1': 'remaining_uses', '3': 9, '4': 1, '5': 5, '10': 'remainingUses'},
    {'1': 'max_uses', '3': 10, '4': 1, '5': 5, '10': 'maxUses'},
    {
      '1': 'alpha_coefficient',
      '3': 11,
      '4': 1,
      '5': 1,
      '10': 'alphaCoefficient'
    },
    {'1': 'temp_coefficient', '3': 12, '4': 1, '5': 1, '10': 'tempCoefficient'},
    {
      '1': 'humidity_coefficient',
      '3': 13,
      '4': 1,
      '5': 1,
      '10': 'humidityCoefficient'
    },
    {
      '1': 'required_channels',
      '3': 14,
      '4': 1,
      '5': 5,
      '10': 'requiredChannels'
    },
    {'1': 'measurement_secs', '3': 15, '4': 1, '5': 5, '10': 'measurementSecs'},
    {'1': 'unit', '3': 16, '4': 1, '5': 9, '10': 'unit'},
    {'1': 'reference_range', '3': 17, '4': 1, '5': 9, '10': 'referenceRange'},
    {'1': 'is_valid', '3': 18, '4': 1, '5': 8, '10': 'isValid'},
    {
      '1': 'validation_message',
      '3': 19,
      '4': 1,
      '5': 9,
      '10': 'validationMessage'
    },
  ],
};

/// Descriptor for `CartridgeDetail`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cartridgeDetailDescriptor = $convert.base64Decode(
    'Cg9DYXJ0cmlkZ2VEZXRhaWwSIwoNY2FydHJpZGdlX3VpZBgBIAEoCVIMY2FydHJpZGdlVWlkEi'
    'MKDWNhdGVnb3J5X2NvZGUYAiABKAVSDGNhdGVnb3J5Q29kZRIdCgp0eXBlX2luZGV4GAMgASgF'
    'Ugl0eXBlSW5kZXgSHwoLbGVnYWN5X2NvZGUYBCABKAVSCmxlZ2FjeUNvZGUSFwoHbmFtZV9rbx'
    'gFIAEoCVIGbmFtZUtvEhcKB25hbWVfZW4YBiABKAlSBm5hbWVFbhIVCgZsb3RfaWQYByABKAlS'
    'BWxvdElkEh8KC2V4cGlyeV9kYXRlGAggASgJUgpleHBpcnlEYXRlEiUKDnJlbWFpbmluZ191c2'
    'VzGAkgASgFUg1yZW1haW5pbmdVc2VzEhkKCG1heF91c2VzGAogASgFUgdtYXhVc2VzEisKEWFs'
    'cGhhX2NvZWZmaWNpZW50GAsgASgBUhBhbHBoYUNvZWZmaWNpZW50EikKEHRlbXBfY29lZmZpY2'
    'llbnQYDCABKAFSD3RlbXBDb2VmZmljaWVudBIxChRodW1pZGl0eV9jb2VmZmljaWVudBgNIAEo'
    'AVITaHVtaWRpdHlDb2VmZmljaWVudBIrChFyZXF1aXJlZF9jaGFubmVscxgOIAEoBVIQcmVxdW'
    'lyZWRDaGFubmVscxIpChBtZWFzdXJlbWVudF9zZWNzGA8gASgFUg9tZWFzdXJlbWVudFNlY3MS'
    'EgoEdW5pdBgQIAEoCVIEdW5pdBInCg9yZWZlcmVuY2VfcmFuZ2UYESABKAlSDnJlZmVyZW5jZV'
    'JhbmdlEhkKCGlzX3ZhbGlkGBIgASgIUgdpc1ZhbGlkEi0KEnZhbGlkYXRpb25fbWVzc2FnZRgT'
    'IAEoCVIRdmFsaWRhdGlvbk1lc3NhZ2U=');

@$core.Deprecated('Use recordUsageRequestDescriptor instead')
const RecordUsageRequest$json = {
  '1': 'RecordUsageRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'session_id', '3': 2, '4': 1, '5': 9, '10': 'sessionId'},
    {'1': 'cartridge_uid', '3': 3, '4': 1, '5': 9, '10': 'cartridgeUid'},
    {'1': 'category_code', '3': 4, '4': 1, '5': 5, '10': 'categoryCode'},
    {'1': 'type_index', '3': 5, '4': 1, '5': 5, '10': 'typeIndex'},
  ],
};

/// Descriptor for `RecordUsageRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recordUsageRequestDescriptor = $convert.base64Decode(
    'ChJSZWNvcmRVc2FnZVJlcXVlc3QSFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklkEh0KCnNlc3Npb2'
    '5faWQYAiABKAlSCXNlc3Npb25JZBIjCg1jYXJ0cmlkZ2VfdWlkGAMgASgJUgxjYXJ0cmlkZ2VV'
    'aWQSIwoNY2F0ZWdvcnlfY29kZRgEIAEoBVIMY2F0ZWdvcnlDb2RlEh0KCnR5cGVfaW5kZXgYBS'
    'ABKAVSCXR5cGVJbmRleA==');

@$core.Deprecated('Use recordUsageResponseDescriptor instead')
const RecordUsageResponse$json = {
  '1': 'RecordUsageResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
    {'1': 'remaining_uses', '3': 2, '4': 1, '5': 5, '10': 'remainingUses'},
    {'1': 'remaining_daily', '3': 3, '4': 1, '5': 5, '10': 'remainingDaily'},
    {
      '1': 'remaining_monthly',
      '3': 4,
      '4': 1,
      '5': 5,
      '10': 'remainingMonthly'
    },
  ],
};

/// Descriptor for `RecordUsageResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recordUsageResponseDescriptor = $convert.base64Decode(
    'ChNSZWNvcmRVc2FnZVJlc3BvbnNlEhgKB3N1Y2Nlc3MYASABKAhSB3N1Y2Nlc3MSJQoOcmVtYW'
    'luaW5nX3VzZXMYAiABKAVSDXJlbWFpbmluZ1VzZXMSJwoPcmVtYWluaW5nX2RhaWx5GAMgASgF'
    'Ug5yZW1haW5pbmdEYWlseRIrChFyZW1haW5pbmdfbW9udGhseRgEIAEoBVIQcmVtYWluaW5nTW'
    '9udGhseQ==');

@$core.Deprecated('Use getUsageHistoryRequestDescriptor instead')
const GetUsageHistoryRequest$json = {
  '1': 'GetUsageHistoryRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'limit', '3': 2, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 3, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `GetUsageHistoryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getUsageHistoryRequestDescriptor =
    $convert.base64Decode(
        'ChZHZXRVc2FnZUhpc3RvcnlSZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZBIUCgVsaW'
        '1pdBgCIAEoBVIFbGltaXQSFgoGb2Zmc2V0GAMgASgFUgZvZmZzZXQ=');

@$core.Deprecated('Use getUsageHistoryResponseDescriptor instead')
const GetUsageHistoryResponse$json = {
  '1': 'GetUsageHistoryResponse',
  '2': [
    {
      '1': 'records',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.CartridgeUsageRecord',
      '10': 'records'
    },
    {'1': 'total_count', '3': 2, '4': 1, '5': 5, '10': 'totalCount'},
  ],
};

/// Descriptor for `GetUsageHistoryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getUsageHistoryResponseDescriptor = $convert.base64Decode(
    'ChdHZXRVc2FnZUhpc3RvcnlSZXNwb25zZRI7CgdyZWNvcmRzGAEgAygLMiEubWFucGFzaWsudj'
    'EuQ2FydHJpZGdlVXNhZ2VSZWNvcmRSB3JlY29yZHMSHwoLdG90YWxfY291bnQYAiABKAVSCnRv'
    'dGFsQ291bnQ=');

@$core.Deprecated('Use cartridgeUsageRecordDescriptor instead')
const CartridgeUsageRecord$json = {
  '1': 'CartridgeUsageRecord',
  '2': [
    {'1': 'record_id', '3': 1, '4': 1, '5': 9, '10': 'recordId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'session_id', '3': 3, '4': 1, '5': 9, '10': 'sessionId'},
    {'1': 'cartridge_uid', '3': 4, '4': 1, '5': 9, '10': 'cartridgeUid'},
    {'1': 'category_code', '3': 5, '4': 1, '5': 5, '10': 'categoryCode'},
    {'1': 'type_index', '3': 6, '4': 1, '5': 5, '10': 'typeIndex'},
    {'1': 'type_name_ko', '3': 7, '4': 1, '5': 9, '10': 'typeNameKo'},
    {
      '1': 'used_at',
      '3': 8,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'usedAt'
    },
  ],
};

/// Descriptor for `CartridgeUsageRecord`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cartridgeUsageRecordDescriptor = $convert.base64Decode(
    'ChRDYXJ0cmlkZ2VVc2FnZVJlY29yZBIbCglyZWNvcmRfaWQYASABKAlSCHJlY29yZElkEhcKB3'
    'VzZXJfaWQYAiABKAlSBnVzZXJJZBIdCgpzZXNzaW9uX2lkGAMgASgJUglzZXNzaW9uSWQSIwoN'
    'Y2FydHJpZGdlX3VpZBgEIAEoCVIMY2FydHJpZGdlVWlkEiMKDWNhdGVnb3J5X2NvZGUYBSABKA'
    'VSDGNhdGVnb3J5Q29kZRIdCgp0eXBlX2luZGV4GAYgASgFUgl0eXBlSW5kZXgSIAoMdHlwZV9u'
    'YW1lX2tvGAcgASgJUgp0eXBlTmFtZUtvEjMKB3VzZWRfYXQYCCABKAsyGi5nb29nbGUucHJvdG'
    '9idWYuVGltZXN0YW1wUgZ1c2VkQXQ=');

@$core.Deprecated('Use getCartridgeTypeRequestDescriptor instead')
const GetCartridgeTypeRequest$json = {
  '1': 'GetCartridgeTypeRequest',
  '2': [
    {'1': 'category_code', '3': 1, '4': 1, '5': 5, '10': 'categoryCode'},
    {'1': 'type_index', '3': 2, '4': 1, '5': 5, '10': 'typeIndex'},
  ],
};

/// Descriptor for `GetCartridgeTypeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCartridgeTypeRequestDescriptor =
    $convert.base64Decode(
        'ChdHZXRDYXJ0cmlkZ2VUeXBlUmVxdWVzdBIjCg1jYXRlZ29yeV9jb2RlGAEgASgFUgxjYXRlZ2'
        '9yeUNvZGUSHQoKdHlwZV9pbmRleBgCIAEoBVIJdHlwZUluZGV4');

@$core.Deprecated('Use listCategoriesRequestDescriptor instead')
const ListCategoriesRequest$json = {
  '1': 'ListCategoriesRequest',
};

/// Descriptor for `ListCategoriesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCategoriesRequestDescriptor =
    $convert.base64Decode('ChVMaXN0Q2F0ZWdvcmllc1JlcXVlc3Q=');

@$core.Deprecated('Use listCategoriesResponseDescriptor instead')
const ListCategoriesResponse$json = {
  '1': 'ListCategoriesResponse',
  '2': [
    {
      '1': 'categories',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.CartridgeCategoryInfo',
      '10': 'categories'
    },
  ],
};

/// Descriptor for `ListCategoriesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCategoriesResponseDescriptor =
    $convert.base64Decode(
        'ChZMaXN0Q2F0ZWdvcmllc1Jlc3BvbnNlEkIKCmNhdGVnb3JpZXMYASADKAsyIi5tYW5wYXNpay'
        '52MS5DYXJ0cmlkZ2VDYXRlZ29yeUluZm9SCmNhdGVnb3JpZXM=');

@$core.Deprecated('Use listTypesByCategoryRequestDescriptor instead')
const ListTypesByCategoryRequest$json = {
  '1': 'ListTypesByCategoryRequest',
  '2': [
    {'1': 'category_code', '3': 1, '4': 1, '5': 5, '10': 'categoryCode'},
  ],
};

/// Descriptor for `ListTypesByCategoryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTypesByCategoryRequestDescriptor =
    $convert.base64Decode(
        'ChpMaXN0VHlwZXNCeUNhdGVnb3J5UmVxdWVzdBIjCg1jYXRlZ29yeV9jb2RlGAEgASgFUgxjYX'
        'RlZ29yeUNvZGU=');

@$core.Deprecated('Use listTypesByCategoryResponseDescriptor instead')
const ListTypesByCategoryResponse$json = {
  '1': 'ListTypesByCategoryResponse',
  '2': [
    {
      '1': 'types',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.CartridgeTypeInfo',
      '10': 'types'
    },
  ],
};

/// Descriptor for `ListTypesByCategoryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTypesByCategoryResponseDescriptor =
    $convert.base64Decode(
        'ChtMaXN0VHlwZXNCeUNhdGVnb3J5UmVzcG9uc2USNAoFdHlwZXMYASADKAsyHi5tYW5wYXNpay'
        '52MS5DYXJ0cmlkZ2VUeXBlSW5mb1IFdHlwZXM=');

@$core.Deprecated('Use getRemainingUsesRequestDescriptor instead')
const GetRemainingUsesRequest$json = {
  '1': 'GetRemainingUsesRequest',
  '2': [
    {'1': 'cartridge_uid', '3': 1, '4': 1, '5': 9, '10': 'cartridgeUid'},
  ],
};

/// Descriptor for `GetRemainingUsesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRemainingUsesRequestDescriptor =
    $convert.base64Decode(
        'ChdHZXRSZW1haW5pbmdVc2VzUmVxdWVzdBIjCg1jYXJ0cmlkZ2VfdWlkGAEgASgJUgxjYXJ0cm'
        'lkZ2VVaWQ=');

@$core.Deprecated('Use getRemainingUsesResponseDescriptor instead')
const GetRemainingUsesResponse$json = {
  '1': 'GetRemainingUsesResponse',
  '2': [
    {'1': 'cartridge_uid', '3': 1, '4': 1, '5': 9, '10': 'cartridgeUid'},
    {'1': 'remaining_uses', '3': 2, '4': 1, '5': 5, '10': 'remainingUses'},
    {'1': 'max_uses', '3': 3, '4': 1, '5': 5, '10': 'maxUses'},
    {'1': 'expiry_date', '3': 4, '4': 1, '5': 9, '10': 'expiryDate'},
    {'1': 'is_expired', '3': 5, '4': 1, '5': 8, '10': 'isExpired'},
  ],
};

/// Descriptor for `GetRemainingUsesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRemainingUsesResponseDescriptor = $convert.base64Decode(
    'ChhHZXRSZW1haW5pbmdVc2VzUmVzcG9uc2USIwoNY2FydHJpZGdlX3VpZBgBIAEoCVIMY2FydH'
    'JpZGdlVWlkEiUKDnJlbWFpbmluZ191c2VzGAIgASgFUg1yZW1haW5pbmdVc2VzEhkKCG1heF91'
    'c2VzGAMgASgFUgdtYXhVc2VzEh8KC2V4cGlyeV9kYXRlGAQgASgJUgpleHBpcnlEYXRlEh0KCm'
    'lzX2V4cGlyZWQYBSABKAhSCWlzRXhwaXJlZA==');

@$core.Deprecated('Use validateCartridgeRequestDescriptor instead')
const ValidateCartridgeRequest$json = {
  '1': 'ValidateCartridgeRequest',
  '2': [
    {'1': 'cartridge_uid', '3': 1, '4': 1, '5': 9, '10': 'cartridgeUid'},
    {'1': 'category_code', '3': 2, '4': 1, '5': 5, '10': 'categoryCode'},
    {'1': 'type_index', '3': 3, '4': 1, '5': 5, '10': 'typeIndex'},
    {'1': 'user_id', '3': 4, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `ValidateCartridgeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validateCartridgeRequestDescriptor = $convert.base64Decode(
    'ChhWYWxpZGF0ZUNhcnRyaWRnZVJlcXVlc3QSIwoNY2FydHJpZGdlX3VpZBgBIAEoCVIMY2FydH'
    'JpZGdlVWlkEiMKDWNhdGVnb3J5X2NvZGUYAiABKAVSDGNhdGVnb3J5Q29kZRIdCgp0eXBlX2lu'
    'ZGV4GAMgASgFUgl0eXBlSW5kZXgSFwoHdXNlcl9pZBgEIAEoCVIGdXNlcklk');

@$core.Deprecated('Use validateCartridgeResponseDescriptor instead')
const ValidateCartridgeResponse$json = {
  '1': 'ValidateCartridgeResponse',
  '2': [
    {'1': 'is_valid', '3': 1, '4': 1, '5': 8, '10': 'isValid'},
    {'1': 'reason', '3': 2, '4': 1, '5': 9, '10': 'reason'},
    {'1': 'remaining_uses', '3': 3, '4': 1, '5': 5, '10': 'remainingUses'},
    {
      '1': 'access_level',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.CartridgeAccessLevel',
      '10': 'accessLevel'
    },
    {
      '1': 'detail',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.manpasik.v1.CartridgeDetail',
      '10': 'detail'
    },
  ],
};

/// Descriptor for `ValidateCartridgeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validateCartridgeResponseDescriptor = $convert.base64Decode(
    'ChlWYWxpZGF0ZUNhcnRyaWRnZVJlc3BvbnNlEhkKCGlzX3ZhbGlkGAEgASgIUgdpc1ZhbGlkEh'
    'YKBnJlYXNvbhgCIAEoCVIGcmVhc29uEiUKDnJlbWFpbmluZ191c2VzGAMgASgFUg1yZW1haW5p'
    'bmdVc2VzEkQKDGFjY2Vzc19sZXZlbBgEIAEoDjIhLm1hbnBhc2lrLnYxLkNhcnRyaWRnZUFjY2'
    'Vzc0xldmVsUgthY2Nlc3NMZXZlbBI0CgZkZXRhaWwYBSABKAsyHC5tYW5wYXNpay52MS5DYXJ0'
    'cmlkZ2VEZXRhaWxSBmRldGFpbA==');

@$core.Deprecated('Use registerFactoryCalibrationRequestDescriptor instead')
const RegisterFactoryCalibrationRequest$json = {
  '1': 'RegisterFactoryCalibrationRequest',
  '2': [
    {'1': 'device_id', '3': 1, '4': 1, '5': 9, '10': 'deviceId'},
    {
      '1': 'cartridge_category',
      '3': 2,
      '4': 1,
      '5': 5,
      '10': 'cartridgeCategory'
    },
    {
      '1': 'cartridge_type_index',
      '3': 3,
      '4': 1,
      '5': 5,
      '10': 'cartridgeTypeIndex'
    },
    {'1': 'alpha', '3': 4, '4': 1, '5': 1, '10': 'alpha'},
    {'1': 'channel_offsets', '3': 5, '4': 3, '5': 1, '10': 'channelOffsets'},
    {'1': 'channel_gains', '3': 6, '4': 3, '5': 1, '10': 'channelGains'},
    {'1': 'temp_coefficient', '3': 7, '4': 1, '5': 1, '10': 'tempCoefficient'},
    {
      '1': 'humidity_coefficient',
      '3': 8,
      '4': 1,
      '5': 1,
      '10': 'humidityCoefficient'
    },
    {
      '1': 'reference_standard',
      '3': 9,
      '4': 1,
      '5': 9,
      '10': 'referenceStandard'
    },
    {'1': 'calibrated_by', '3': 10, '4': 1, '5': 9, '10': 'calibratedBy'},
  ],
};

/// Descriptor for `RegisterFactoryCalibrationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List registerFactoryCalibrationRequestDescriptor = $convert.base64Decode(
    'CiFSZWdpc3RlckZhY3RvcnlDYWxpYnJhdGlvblJlcXVlc3QSGwoJZGV2aWNlX2lkGAEgASgJUg'
    'hkZXZpY2VJZBItChJjYXJ0cmlkZ2VfY2F0ZWdvcnkYAiABKAVSEWNhcnRyaWRnZUNhdGVnb3J5'
    'EjAKFGNhcnRyaWRnZV90eXBlX2luZGV4GAMgASgFUhJjYXJ0cmlkZ2VUeXBlSW5kZXgSFAoFYW'
    'xwaGEYBCABKAFSBWFscGhhEicKD2NoYW5uZWxfb2Zmc2V0cxgFIAMoAVIOY2hhbm5lbE9mZnNl'
    'dHMSIwoNY2hhbm5lbF9nYWlucxgGIAMoAVIMY2hhbm5lbEdhaW5zEikKEHRlbXBfY29lZmZpY2'
    'llbnQYByABKAFSD3RlbXBDb2VmZmljaWVudBIxChRodW1pZGl0eV9jb2VmZmljaWVudBgIIAEo'
    'AVITaHVtaWRpdHlDb2VmZmljaWVudBItChJyZWZlcmVuY2Vfc3RhbmRhcmQYCSABKAlSEXJlZm'
    'VyZW5jZVN0YW5kYXJkEiMKDWNhbGlicmF0ZWRfYnkYCiABKAlSDGNhbGlicmF0ZWRCeQ==');

@$core.Deprecated('Use performFieldCalibrationRequestDescriptor instead')
const PerformFieldCalibrationRequest$json = {
  '1': 'PerformFieldCalibrationRequest',
  '2': [
    {'1': 'device_id', '3': 1, '4': 1, '5': 9, '10': 'deviceId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'cartridge_category',
      '3': 3,
      '4': 1,
      '5': 5,
      '10': 'cartridgeCategory'
    },
    {
      '1': 'cartridge_type_index',
      '3': 4,
      '4': 1,
      '5': 5,
      '10': 'cartridgeTypeIndex'
    },
    {'1': 'reference_values', '3': 5, '4': 3, '5': 1, '10': 'referenceValues'},
    {'1': 'measured_values', '3': 6, '4': 3, '5': 1, '10': 'measuredValues'},
    {'1': 'temperature_c', '3': 7, '4': 1, '5': 1, '10': 'temperatureC'},
    {'1': 'humidity_pct', '3': 8, '4': 1, '5': 1, '10': 'humidityPct'},
  ],
};

/// Descriptor for `PerformFieldCalibrationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List performFieldCalibrationRequestDescriptor = $convert.base64Decode(
    'Ch5QZXJmb3JtRmllbGRDYWxpYnJhdGlvblJlcXVlc3QSGwoJZGV2aWNlX2lkGAEgASgJUghkZX'
    'ZpY2VJZBIXCgd1c2VyX2lkGAIgASgJUgZ1c2VySWQSLQoSY2FydHJpZGdlX2NhdGVnb3J5GAMg'
    'ASgFUhFjYXJ0cmlkZ2VDYXRlZ29yeRIwChRjYXJ0cmlkZ2VfdHlwZV9pbmRleBgEIAEoBVISY2'
    'FydHJpZGdlVHlwZUluZGV4EikKEHJlZmVyZW5jZV92YWx1ZXMYBSADKAFSD3JlZmVyZW5jZVZh'
    'bHVlcxInCg9tZWFzdXJlZF92YWx1ZXMYBiADKAFSDm1lYXN1cmVkVmFsdWVzEiMKDXRlbXBlcm'
    'F0dXJlX2MYByABKAFSDHRlbXBlcmF0dXJlQxIhCgxodW1pZGl0eV9wY3QYCCABKAFSC2h1bWlk'
    'aXR5UGN0');

@$core.Deprecated('Use getCalibrationRequestDescriptor instead')
const GetCalibrationRequest$json = {
  '1': 'GetCalibrationRequest',
  '2': [
    {'1': 'device_id', '3': 1, '4': 1, '5': 9, '10': 'deviceId'},
    {
      '1': 'cartridge_category',
      '3': 2,
      '4': 1,
      '5': 5,
      '10': 'cartridgeCategory'
    },
    {
      '1': 'cartridge_type_index',
      '3': 3,
      '4': 1,
      '5': 5,
      '10': 'cartridgeTypeIndex'
    },
  ],
};

/// Descriptor for `GetCalibrationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCalibrationRequestDescriptor = $convert.base64Decode(
    'ChVHZXRDYWxpYnJhdGlvblJlcXVlc3QSGwoJZGV2aWNlX2lkGAEgASgJUghkZXZpY2VJZBItCh'
    'JjYXJ0cmlkZ2VfY2F0ZWdvcnkYAiABKAVSEWNhcnRyaWRnZUNhdGVnb3J5EjAKFGNhcnRyaWRn'
    'ZV90eXBlX2luZGV4GAMgASgFUhJjYXJ0cmlkZ2VUeXBlSW5kZXg=');

@$core.Deprecated('Use calibrationRecordDescriptor instead')
const CalibrationRecord$json = {
  '1': 'CalibrationRecord',
  '2': [
    {'1': 'calibration_id', '3': 1, '4': 1, '5': 9, '10': 'calibrationId'},
    {'1': 'device_id', '3': 2, '4': 1, '5': 9, '10': 'deviceId'},
    {
      '1': 'cartridge_category',
      '3': 3,
      '4': 1,
      '5': 5,
      '10': 'cartridgeCategory'
    },
    {
      '1': 'cartridge_type_index',
      '3': 4,
      '4': 1,
      '5': 5,
      '10': 'cartridgeTypeIndex'
    },
    {
      '1': 'calibration_type',
      '3': 5,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.CalibrationType',
      '10': 'calibrationType'
    },
    {'1': 'alpha', '3': 6, '4': 1, '5': 1, '10': 'alpha'},
    {'1': 'channel_offsets', '3': 7, '4': 3, '5': 1, '10': 'channelOffsets'},
    {'1': 'channel_gains', '3': 8, '4': 3, '5': 1, '10': 'channelGains'},
    {'1': 'temp_coefficient', '3': 9, '4': 1, '5': 1, '10': 'tempCoefficient'},
    {
      '1': 'humidity_coefficient',
      '3': 10,
      '4': 1,
      '5': 1,
      '10': 'humidityCoefficient'
    },
    {'1': 'accuracy_score', '3': 11, '4': 1, '5': 1, '10': 'accuracyScore'},
    {
      '1': 'reference_standard',
      '3': 12,
      '4': 1,
      '5': 9,
      '10': 'referenceStandard'
    },
    {'1': 'calibrated_by', '3': 13, '4': 1, '5': 9, '10': 'calibratedBy'},
    {
      '1': 'calibrated_at',
      '3': 14,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'calibratedAt'
    },
    {
      '1': 'expires_at',
      '3': 15,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'expiresAt'
    },
    {
      '1': 'status',
      '3': 16,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.CalibrationStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `CalibrationRecord`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List calibrationRecordDescriptor = $convert.base64Decode(
    'ChFDYWxpYnJhdGlvblJlY29yZBIlCg5jYWxpYnJhdGlvbl9pZBgBIAEoCVINY2FsaWJyYXRpb2'
    '5JZBIbCglkZXZpY2VfaWQYAiABKAlSCGRldmljZUlkEi0KEmNhcnRyaWRnZV9jYXRlZ29yeRgD'
    'IAEoBVIRY2FydHJpZGdlQ2F0ZWdvcnkSMAoUY2FydHJpZGdlX3R5cGVfaW5kZXgYBCABKAVSEm'
    'NhcnRyaWRnZVR5cGVJbmRleBJHChBjYWxpYnJhdGlvbl90eXBlGAUgASgOMhwubWFucGFzaWsu'
    'djEuQ2FsaWJyYXRpb25UeXBlUg9jYWxpYnJhdGlvblR5cGUSFAoFYWxwaGEYBiABKAFSBWFscG'
    'hhEicKD2NoYW5uZWxfb2Zmc2V0cxgHIAMoAVIOY2hhbm5lbE9mZnNldHMSIwoNY2hhbm5lbF9n'
    'YWlucxgIIAMoAVIMY2hhbm5lbEdhaW5zEikKEHRlbXBfY29lZmZpY2llbnQYCSABKAFSD3RlbX'
    'BDb2VmZmljaWVudBIxChRodW1pZGl0eV9jb2VmZmljaWVudBgKIAEoAVITaHVtaWRpdHlDb2Vm'
    'ZmljaWVudBIlCg5hY2N1cmFjeV9zY29yZRgLIAEoAVINYWNjdXJhY3lTY29yZRItChJyZWZlcm'
    'VuY2Vfc3RhbmRhcmQYDCABKAlSEXJlZmVyZW5jZVN0YW5kYXJkEiMKDWNhbGlicmF0ZWRfYnkY'
    'DSABKAlSDGNhbGlicmF0ZWRCeRI/Cg1jYWxpYnJhdGVkX2F0GA4gASgLMhouZ29vZ2xlLnByb3'
    'RvYnVmLlRpbWVzdGFtcFIMY2FsaWJyYXRlZEF0EjkKCmV4cGlyZXNfYXQYDyABKAsyGi5nb29n'
    'bGUucHJvdG9idWYuVGltZXN0YW1wUglleHBpcmVzQXQSNgoGc3RhdHVzGBAgASgOMh4ubWFucG'
    'FzaWsudjEuQ2FsaWJyYXRpb25TdGF0dXNSBnN0YXR1cw==');

@$core.Deprecated('Use listCalibrationHistoryRequestDescriptor instead')
const ListCalibrationHistoryRequest$json = {
  '1': 'ListCalibrationHistoryRequest',
  '2': [
    {'1': 'device_id', '3': 1, '4': 1, '5': 9, '10': 'deviceId'},
    {'1': 'limit', '3': 2, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 3, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `ListCalibrationHistoryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCalibrationHistoryRequestDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0Q2FsaWJyYXRpb25IaXN0b3J5UmVxdWVzdBIbCglkZXZpY2VfaWQYASABKAlSCGRldm'
        'ljZUlkEhQKBWxpbWl0GAIgASgFUgVsaW1pdBIWCgZvZmZzZXQYAyABKAVSBm9mZnNldA==');

@$core.Deprecated('Use listCalibrationHistoryResponseDescriptor instead')
const ListCalibrationHistoryResponse$json = {
  '1': 'ListCalibrationHistoryResponse',
  '2': [
    {
      '1': 'records',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.CalibrationRecord',
      '10': 'records'
    },
    {'1': 'total_count', '3': 2, '4': 1, '5': 5, '10': 'totalCount'},
  ],
};

/// Descriptor for `ListCalibrationHistoryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCalibrationHistoryResponseDescriptor =
    $convert.base64Decode(
        'Ch5MaXN0Q2FsaWJyYXRpb25IaXN0b3J5UmVzcG9uc2USOAoHcmVjb3JkcxgBIAMoCzIeLm1hbn'
        'Bhc2lrLnYxLkNhbGlicmF0aW9uUmVjb3JkUgdyZWNvcmRzEh8KC3RvdGFsX2NvdW50GAIgASgF'
        'Ugp0b3RhbENvdW50');

@$core.Deprecated('Use checkCalibrationStatusRequestDescriptor instead')
const CheckCalibrationStatusRequest$json = {
  '1': 'CheckCalibrationStatusRequest',
  '2': [
    {'1': 'device_id', '3': 1, '4': 1, '5': 9, '10': 'deviceId'},
    {
      '1': 'cartridge_category',
      '3': 2,
      '4': 1,
      '5': 5,
      '10': 'cartridgeCategory'
    },
    {
      '1': 'cartridge_type_index',
      '3': 3,
      '4': 1,
      '5': 5,
      '10': 'cartridgeTypeIndex'
    },
  ],
};

/// Descriptor for `CheckCalibrationStatusRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List checkCalibrationStatusRequestDescriptor =
    $convert.base64Decode(
        'Ch1DaGVja0NhbGlicmF0aW9uU3RhdHVzUmVxdWVzdBIbCglkZXZpY2VfaWQYASABKAlSCGRldm'
        'ljZUlkEi0KEmNhcnRyaWRnZV9jYXRlZ29yeRgCIAEoBVIRY2FydHJpZGdlQ2F0ZWdvcnkSMAoU'
        'Y2FydHJpZGdlX3R5cGVfaW5kZXgYAyABKAVSEmNhcnRyaWRnZVR5cGVJbmRleA==');

@$core.Deprecated('Use calibrationStatusResponseDescriptor instead')
const CalibrationStatusResponse$json = {
  '1': 'CalibrationStatusResponse',
  '2': [
    {
      '1': 'status',
      '3': 1,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.CalibrationStatus',
      '10': 'status'
    },
    {'1': 'device_id', '3': 2, '4': 1, '5': 9, '10': 'deviceId'},
    {'1': 'message', '3': 3, '4': 1, '5': 9, '10': 'message'},
    {
      '1': 'last_calibrated_at',
      '3': 4,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'lastCalibratedAt'
    },
    {
      '1': 'expires_at',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'expiresAt'
    },
    {'1': 'days_until_expiry', '3': 6, '4': 1, '5': 5, '10': 'daysUntilExpiry'},
    {
      '1': 'latest_record',
      '3': 7,
      '4': 1,
      '5': 11,
      '6': '.manpasik.v1.CalibrationRecord',
      '10': 'latestRecord'
    },
  ],
};

/// Descriptor for `CalibrationStatusResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List calibrationStatusResponseDescriptor = $convert.base64Decode(
    'ChlDYWxpYnJhdGlvblN0YXR1c1Jlc3BvbnNlEjYKBnN0YXR1cxgBIAEoDjIeLm1hbnBhc2lrLn'
    'YxLkNhbGlicmF0aW9uU3RhdHVzUgZzdGF0dXMSGwoJZGV2aWNlX2lkGAIgASgJUghkZXZpY2VJ'
    'ZBIYCgdtZXNzYWdlGAMgASgJUgdtZXNzYWdlEkgKEmxhc3RfY2FsaWJyYXRlZF9hdBgEIAEoCz'
    'IaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSEGxhc3RDYWxpYnJhdGVkQXQSOQoKZXhwaXJl'
    'c19hdBgFIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSCWV4cGlyZXNBdBIqChFkYX'
    'lzX3VudGlsX2V4cGlyeRgGIAEoBVIPZGF5c1VudGlsRXhwaXJ5EkMKDWxhdGVzdF9yZWNvcmQY'
    'ByABKAsyHi5tYW5wYXNpay52MS5DYWxpYnJhdGlvblJlY29yZFIMbGF0ZXN0UmVjb3Jk');

@$core.Deprecated('Use listCalibrationModelsRequestDescriptor instead')
const ListCalibrationModelsRequest$json = {
  '1': 'ListCalibrationModelsRequest',
};

/// Descriptor for `ListCalibrationModelsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCalibrationModelsRequestDescriptor =
    $convert.base64Decode('ChxMaXN0Q2FsaWJyYXRpb25Nb2RlbHNSZXF1ZXN0');

@$core.Deprecated('Use calibrationModelDescriptor instead')
const CalibrationModel$json = {
  '1': 'CalibrationModel',
  '2': [
    {'1': 'model_id', '3': 1, '4': 1, '5': 9, '10': 'modelId'},
    {
      '1': 'cartridge_category',
      '3': 2,
      '4': 1,
      '5': 5,
      '10': 'cartridgeCategory'
    },
    {
      '1': 'cartridge_type_index',
      '3': 3,
      '4': 1,
      '5': 5,
      '10': 'cartridgeTypeIndex'
    },
    {'1': 'name', '3': 4, '4': 1, '5': 9, '10': 'name'},
    {'1': 'version', '3': 5, '4': 1, '5': 9, '10': 'version'},
    {'1': 'default_alpha', '3': 6, '4': 1, '5': 1, '10': 'defaultAlpha'},
    {'1': 'validity_days', '3': 7, '4': 1, '5': 5, '10': 'validityDays'},
    {'1': 'description', '3': 8, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'created_at',
      '3': 9,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
  ],
};

/// Descriptor for `CalibrationModel`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List calibrationModelDescriptor = $convert.base64Decode(
    'ChBDYWxpYnJhdGlvbk1vZGVsEhkKCG1vZGVsX2lkGAEgASgJUgdtb2RlbElkEi0KEmNhcnRyaW'
    'RnZV9jYXRlZ29yeRgCIAEoBVIRY2FydHJpZGdlQ2F0ZWdvcnkSMAoUY2FydHJpZGdlX3R5cGVf'
    'aW5kZXgYAyABKAVSEmNhcnRyaWRnZVR5cGVJbmRleBISCgRuYW1lGAQgASgJUgRuYW1lEhgKB3'
    'ZlcnNpb24YBSABKAlSB3ZlcnNpb24SIwoNZGVmYXVsdF9hbHBoYRgGIAEoAVIMZGVmYXVsdEFs'
    'cGhhEiMKDXZhbGlkaXR5X2RheXMYByABKAVSDHZhbGlkaXR5RGF5cxIgCgtkZXNjcmlwdGlvbh'
    'gIIAEoCVILZGVzY3JpcHRpb24SOQoKY3JlYXRlZF9hdBgJIAEoCzIaLmdvb2dsZS5wcm90b2J1'
    'Zi5UaW1lc3RhbXBSCWNyZWF0ZWRBdA==');

@$core.Deprecated('Use listCalibrationModelsResponseDescriptor instead')
const ListCalibrationModelsResponse$json = {
  '1': 'ListCalibrationModelsResponse',
  '2': [
    {
      '1': 'models',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.CalibrationModel',
      '10': 'models'
    },
  ],
};

/// Descriptor for `ListCalibrationModelsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCalibrationModelsResponseDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0Q2FsaWJyYXRpb25Nb2RlbHNSZXNwb25zZRI1CgZtb2RlbHMYASADKAsyHS5tYW5wYX'
        'Npay52MS5DYWxpYnJhdGlvbk1vZGVsUgZtb2RlbHM=');

@$core.Deprecated('Use setHealthGoalRequestDescriptor instead')
const SetHealthGoalRequest$json = {
  '1': 'SetHealthGoalRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'category',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.GoalCategory',
      '10': 'category'
    },
    {'1': 'metric_name', '3': 3, '4': 1, '5': 9, '10': 'metricName'},
    {'1': 'target_value', '3': 4, '4': 1, '5': 1, '10': 'targetValue'},
    {'1': 'unit', '3': 5, '4': 1, '5': 9, '10': 'unit'},
    {'1': 'description', '3': 6, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'target_date',
      '3': 7,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'targetDate'
    },
  ],
};

/// Descriptor for `SetHealthGoalRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setHealthGoalRequestDescriptor = $convert.base64Decode(
    'ChRTZXRIZWFsdGhHb2FsUmVxdWVzdBIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQSNQoIY2F0ZW'
    'dvcnkYAiABKA4yGS5tYW5wYXNpay52MS5Hb2FsQ2F0ZWdvcnlSCGNhdGVnb3J5Eh8KC21ldHJp'
    'Y19uYW1lGAMgASgJUgptZXRyaWNOYW1lEiEKDHRhcmdldF92YWx1ZRgEIAEoAVILdGFyZ2V0Vm'
    'FsdWUSEgoEdW5pdBgFIAEoCVIEdW5pdBIgCgtkZXNjcmlwdGlvbhgGIAEoCVILZGVzY3JpcHRp'
    'b24SOwoLdGFyZ2V0X2RhdGUYByABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUgp0YX'
    'JnZXREYXRl');

@$core.Deprecated('Use healthGoalDescriptor instead')
const HealthGoal$json = {
  '1': 'HealthGoal',
  '2': [
    {'1': 'goal_id', '3': 1, '4': 1, '5': 9, '10': 'goalId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'category',
      '3': 3,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.GoalCategory',
      '10': 'category'
    },
    {'1': 'metric_name', '3': 4, '4': 1, '5': 9, '10': 'metricName'},
    {'1': 'target_value', '3': 5, '4': 1, '5': 1, '10': 'targetValue'},
    {'1': 'current_value', '3': 6, '4': 1, '5': 1, '10': 'currentValue'},
    {'1': 'unit', '3': 7, '4': 1, '5': 9, '10': 'unit'},
    {'1': 'progress_pct', '3': 8, '4': 1, '5': 1, '10': 'progressPct'},
    {
      '1': 'status',
      '3': 9,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.GoalStatus',
      '10': 'status'
    },
    {'1': 'description', '3': 10, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'created_at',
      '3': 11,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
    {
      '1': 'target_date',
      '3': 12,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'targetDate'
    },
    {
      '1': 'achieved_at',
      '3': 13,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'achievedAt'
    },
  ],
};

/// Descriptor for `HealthGoal`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List healthGoalDescriptor = $convert.base64Decode(
    'CgpIZWFsdGhHb2FsEhcKB2dvYWxfaWQYASABKAlSBmdvYWxJZBIXCgd1c2VyX2lkGAIgASgJUg'
    'Z1c2VySWQSNQoIY2F0ZWdvcnkYAyABKA4yGS5tYW5wYXNpay52MS5Hb2FsQ2F0ZWdvcnlSCGNh'
    'dGVnb3J5Eh8KC21ldHJpY19uYW1lGAQgASgJUgptZXRyaWNOYW1lEiEKDHRhcmdldF92YWx1ZR'
    'gFIAEoAVILdGFyZ2V0VmFsdWUSIwoNY3VycmVudF92YWx1ZRgGIAEoAVIMY3VycmVudFZhbHVl'
    'EhIKBHVuaXQYByABKAlSBHVuaXQSIQoMcHJvZ3Jlc3NfcGN0GAggASgBUgtwcm9ncmVzc1BjdB'
    'IvCgZzdGF0dXMYCSABKA4yFy5tYW5wYXNpay52MS5Hb2FsU3RhdHVzUgZzdGF0dXMSIAoLZGVz'
    'Y3JpcHRpb24YCiABKAlSC2Rlc2NyaXB0aW9uEjkKCmNyZWF0ZWRfYXQYCyABKAsyGi5nb29nbG'
    'UucHJvdG9idWYuVGltZXN0YW1wUgljcmVhdGVkQXQSOwoLdGFyZ2V0X2RhdGUYDCABKAsyGi5n'
    'b29nbGUucHJvdG9idWYuVGltZXN0YW1wUgp0YXJnZXREYXRlEjsKC2FjaGlldmVkX2F0GA0gAS'
    'gLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFIKYWNoaWV2ZWRBdA==');

@$core.Deprecated('Use getHealthGoalsRequestDescriptor instead')
const GetHealthGoalsRequest$json = {
  '1': 'GetHealthGoalsRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'status_filter',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.GoalStatus',
      '10': 'statusFilter'
    },
  ],
};

/// Descriptor for `GetHealthGoalsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getHealthGoalsRequestDescriptor = $convert.base64Decode(
    'ChVHZXRIZWFsdGhHb2Fsc1JlcXVlc3QSFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklkEjwKDXN0YX'
    'R1c19maWx0ZXIYAiABKA4yFy5tYW5wYXNpay52MS5Hb2FsU3RhdHVzUgxzdGF0dXNGaWx0ZXI=');

@$core.Deprecated('Use getHealthGoalsResponseDescriptor instead')
const GetHealthGoalsResponse$json = {
  '1': 'GetHealthGoalsResponse',
  '2': [
    {
      '1': 'goals',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.HealthGoal',
      '10': 'goals'
    },
  ],
};

/// Descriptor for `GetHealthGoalsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getHealthGoalsResponseDescriptor =
    $convert.base64Decode(
        'ChZHZXRIZWFsdGhHb2Fsc1Jlc3BvbnNlEi0KBWdvYWxzGAEgAygLMhcubWFucGFzaWsudjEuSG'
        'VhbHRoR29hbFIFZ29hbHM=');

@$core.Deprecated('Use generateCoachingRequestDescriptor instead')
const GenerateCoachingRequest$json = {
  '1': 'GenerateCoachingRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'measurement_id', '3': 2, '4': 1, '5': 9, '10': 'measurementId'},
    {
      '1': 'coaching_type',
      '3': 3,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.CoachingType',
      '10': 'coachingType'
    },
  ],
};

/// Descriptor for `GenerateCoachingRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List generateCoachingRequestDescriptor = $convert.base64Decode(
    'ChdHZW5lcmF0ZUNvYWNoaW5nUmVxdWVzdBIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQSJQoObW'
    'Vhc3VyZW1lbnRfaWQYAiABKAlSDW1lYXN1cmVtZW50SWQSPgoNY29hY2hpbmdfdHlwZRgDIAEo'
    'DjIZLm1hbnBhc2lrLnYxLkNvYWNoaW5nVHlwZVIMY29hY2hpbmdUeXBl');

@$core.Deprecated('Use coachingMessageDescriptor instead')
const CoachingMessage$json = {
  '1': 'CoachingMessage',
  '2': [
    {'1': 'message_id', '3': 1, '4': 1, '5': 9, '10': 'messageId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'coaching_type',
      '3': 3,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.CoachingType',
      '10': 'coachingType'
    },
    {'1': 'title', '3': 4, '4': 1, '5': 9, '10': 'title'},
    {'1': 'body', '3': 5, '4': 1, '5': 9, '10': 'body'},
    {
      '1': 'risk_level',
      '3': 6,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.RiskLevel',
      '10': 'riskLevel'
    },
    {'1': 'action_items', '3': 7, '4': 3, '5': 9, '10': 'actionItems'},
    {'1': 'related_metric', '3': 8, '4': 1, '5': 9, '10': 'relatedMetric'},
    {'1': 'related_value', '3': 9, '4': 1, '5': 1, '10': 'relatedValue'},
    {
      '1': 'created_at',
      '3': 10,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
  ],
};

/// Descriptor for `CoachingMessage`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List coachingMessageDescriptor = $convert.base64Decode(
    'Cg9Db2FjaGluZ01lc3NhZ2USHQoKbWVzc2FnZV9pZBgBIAEoCVIJbWVzc2FnZUlkEhcKB3VzZX'
    'JfaWQYAiABKAlSBnVzZXJJZBI+Cg1jb2FjaGluZ190eXBlGAMgASgOMhkubWFucGFzaWsudjEu'
    'Q29hY2hpbmdUeXBlUgxjb2FjaGluZ1R5cGUSFAoFdGl0bGUYBCABKAlSBXRpdGxlEhIKBGJvZH'
    'kYBSABKAlSBGJvZHkSNQoKcmlza19sZXZlbBgGIAEoDjIWLm1hbnBhc2lrLnYxLlJpc2tMZXZl'
    'bFIJcmlza0xldmVsEiEKDGFjdGlvbl9pdGVtcxgHIAMoCVILYWN0aW9uSXRlbXMSJQoOcmVsYX'
    'RlZF9tZXRyaWMYCCABKAlSDXJlbGF0ZWRNZXRyaWMSIwoNcmVsYXRlZF92YWx1ZRgJIAEoAVIM'
    'cmVsYXRlZFZhbHVlEjkKCmNyZWF0ZWRfYXQYCiABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZX'
    'N0YW1wUgljcmVhdGVkQXQ=');

@$core.Deprecated('Use listCoachingMessagesRequestDescriptor instead')
const ListCoachingMessagesRequest$json = {
  '1': 'ListCoachingMessagesRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'type_filter',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.CoachingType',
      '10': 'typeFilter'
    },
    {'1': 'limit', '3': 3, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 4, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `ListCoachingMessagesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCoachingMessagesRequestDescriptor =
    $convert.base64Decode(
        'ChtMaXN0Q29hY2hpbmdNZXNzYWdlc1JlcXVlc3QSFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklkEj'
        'oKC3R5cGVfZmlsdGVyGAIgASgOMhkubWFucGFzaWsudjEuQ29hY2hpbmdUeXBlUgp0eXBlRmls'
        'dGVyEhQKBWxpbWl0GAMgASgFUgVsaW1pdBIWCgZvZmZzZXQYBCABKAVSBm9mZnNldA==');

@$core.Deprecated('Use listCoachingMessagesResponseDescriptor instead')
const ListCoachingMessagesResponse$json = {
  '1': 'ListCoachingMessagesResponse',
  '2': [
    {
      '1': 'messages',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.CoachingMessage',
      '10': 'messages'
    },
    {'1': 'total_count', '3': 2, '4': 1, '5': 5, '10': 'totalCount'},
  ],
};

/// Descriptor for `ListCoachingMessagesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCoachingMessagesResponseDescriptor =
    $convert.base64Decode(
        'ChxMaXN0Q29hY2hpbmdNZXNzYWdlc1Jlc3BvbnNlEjgKCG1lc3NhZ2VzGAEgAygLMhwubWFucG'
        'FzaWsudjEuQ29hY2hpbmdNZXNzYWdlUghtZXNzYWdlcxIfCgt0b3RhbF9jb3VudBgCIAEoBVIK'
        'dG90YWxDb3VudA==');

@$core.Deprecated('Use generateDailyReportRequestDescriptor instead')
const GenerateDailyReportRequest$json = {
  '1': 'GenerateDailyReportRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'date',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'date'
    },
  ],
};

/// Descriptor for `GenerateDailyReportRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List generateDailyReportRequestDescriptor =
    $convert.base64Decode(
        'ChpHZW5lcmF0ZURhaWx5UmVwb3J0UmVxdWVzdBIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQSLg'
        'oEZGF0ZRgCIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSBGRhdGU=');

@$core.Deprecated('Use dailyHealthReportDescriptor instead')
const DailyHealthReport$json = {
  '1': 'DailyHealthReport',
  '2': [
    {'1': 'report_id', '3': 1, '4': 1, '5': 9, '10': 'reportId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'report_date',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'reportDate'
    },
    {'1': 'overall_score', '3': 4, '4': 1, '5': 1, '10': 'overallScore'},
    {
      '1': 'measurements_count',
      '3': 5,
      '4': 1,
      '5': 5,
      '10': 'measurementsCount'
    },
    {
      '1': 'highlights',
      '3': 6,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.CoachingMessage',
      '10': 'highlights'
    },
    {'1': 'summary', '3': 7, '4': 1, '5': 9, '10': 'summary'},
    {'1': 'recommendations', '3': 8, '4': 3, '5': 9, '10': 'recommendations'},
  ],
};

/// Descriptor for `DailyHealthReport`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dailyHealthReportDescriptor = $convert.base64Decode(
    'ChFEYWlseUhlYWx0aFJlcG9ydBIbCglyZXBvcnRfaWQYASABKAlSCHJlcG9ydElkEhcKB3VzZX'
    'JfaWQYAiABKAlSBnVzZXJJZBI7CgtyZXBvcnRfZGF0ZRgDIAEoCzIaLmdvb2dsZS5wcm90b2J1'
    'Zi5UaW1lc3RhbXBSCnJlcG9ydERhdGUSIwoNb3ZlcmFsbF9zY29yZRgEIAEoAVIMb3ZlcmFsbF'
    'Njb3JlEi0KEm1lYXN1cmVtZW50c19jb3VudBgFIAEoBVIRbWVhc3VyZW1lbnRzQ291bnQSPAoK'
    'aGlnaGxpZ2h0cxgGIAMoCzIcLm1hbnBhc2lrLnYxLkNvYWNoaW5nTWVzc2FnZVIKaGlnaGxpZ2'
    'h0cxIYCgdzdW1tYXJ5GAcgASgJUgdzdW1tYXJ5EigKD3JlY29tbWVuZGF0aW9ucxgIIAMoCVIP'
    'cmVjb21tZW5kYXRpb25z');

@$core.Deprecated('Use getWeeklyReportRequestDescriptor instead')
const GetWeeklyReportRequest$json = {
  '1': 'GetWeeklyReportRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'week_start',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'weekStart'
    },
  ],
};

/// Descriptor for `GetWeeklyReportRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getWeeklyReportRequestDescriptor = $convert.base64Decode(
    'ChZHZXRXZWVrbHlSZXBvcnRSZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZBI5Cgp3ZW'
    'VrX3N0YXJ0GAIgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFIJd2Vla1N0YXJ0');

@$core.Deprecated('Use weeklyHealthReportDescriptor instead')
const WeeklyHealthReport$json = {
  '1': 'WeeklyHealthReport',
  '2': [
    {'1': 'report_id', '3': 1, '4': 1, '5': 9, '10': 'reportId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'week_start',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'weekStart'
    },
    {
      '1': 'week_end',
      '3': 4,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'weekEnd'
    },
    {'1': 'average_score', '3': 5, '4': 1, '5': 1, '10': 'averageScore'},
    {'1': 'score_trend', '3': 6, '4': 1, '5': 9, '10': 'scoreTrend'},
    {
      '1': 'total_measurements',
      '3': 7,
      '4': 1,
      '5': 5,
      '10': 'totalMeasurements'
    },
    {'1': 'goals_achieved', '3': 8, '4': 1, '5': 5, '10': 'goalsAchieved'},
    {'1': 'goals_active', '3': 9, '4': 1, '5': 5, '10': 'goalsActive'},
    {
      '1': 'daily_reports',
      '3': 10,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.DailyHealthReport',
      '10': 'dailyReports'
    },
    {'1': 'weekly_summary', '3': 11, '4': 1, '5': 9, '10': 'weeklySummary'},
    {'1': 'key_insights', '3': 12, '4': 3, '5': 9, '10': 'keyInsights'},
  ],
};

/// Descriptor for `WeeklyHealthReport`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List weeklyHealthReportDescriptor = $convert.base64Decode(
    'ChJXZWVrbHlIZWFsdGhSZXBvcnQSGwoJcmVwb3J0X2lkGAEgASgJUghyZXBvcnRJZBIXCgd1c2'
    'VyX2lkGAIgASgJUgZ1c2VySWQSOQoKd2Vla19zdGFydBgDIAEoCzIaLmdvb2dsZS5wcm90b2J1'
    'Zi5UaW1lc3RhbXBSCXdlZWtTdGFydBI1Cgh3ZWVrX2VuZBgEIAEoCzIaLmdvb2dsZS5wcm90b2'
    'J1Zi5UaW1lc3RhbXBSB3dlZWtFbmQSIwoNYXZlcmFnZV9zY29yZRgFIAEoAVIMYXZlcmFnZVNj'
    'b3JlEh8KC3Njb3JlX3RyZW5kGAYgASgJUgpzY29yZVRyZW5kEi0KEnRvdGFsX21lYXN1cmVtZW'
    '50cxgHIAEoBVIRdG90YWxNZWFzdXJlbWVudHMSJQoOZ29hbHNfYWNoaWV2ZWQYCCABKAVSDWdv'
    'YWxzQWNoaWV2ZWQSIQoMZ29hbHNfYWN0aXZlGAkgASgFUgtnb2Fsc0FjdGl2ZRJDCg1kYWlseV'
    '9yZXBvcnRzGAogAygLMh4ubWFucGFzaWsudjEuRGFpbHlIZWFsdGhSZXBvcnRSDGRhaWx5UmVw'
    'b3J0cxIlCg53ZWVrbHlfc3VtbWFyeRgLIAEoCVINd2Vla2x5U3VtbWFyeRIhCgxrZXlfaW5zaW'
    'dodHMYDCADKAlSC2tleUluc2lnaHRz');

@$core.Deprecated('Use getRecommendationsRequestDescriptor instead')
const GetRecommendationsRequest$json = {
  '1': 'GetRecommendationsRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'type_filter',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.RecommendationType',
      '10': 'typeFilter'
    },
    {'1': 'limit', '3': 3, '4': 1, '5': 5, '10': 'limit'},
  ],
};

/// Descriptor for `GetRecommendationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRecommendationsRequestDescriptor = $convert.base64Decode(
    'ChlHZXRSZWNvbW1lbmRhdGlvbnNSZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZBJACg'
    't0eXBlX2ZpbHRlchgCIAEoDjIfLm1hbnBhc2lrLnYxLlJlY29tbWVuZGF0aW9uVHlwZVIKdHlw'
    'ZUZpbHRlchIUCgVsaW1pdBgDIAEoBVIFbGltaXQ=');

@$core.Deprecated('Use recommendationDescriptor instead')
const Recommendation$json = {
  '1': 'Recommendation',
  '2': [
    {
      '1': 'recommendation_id',
      '3': 1,
      '4': 1,
      '5': 9,
      '10': 'recommendationId'
    },
    {
      '1': 'type',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.RecommendationType',
      '10': 'type'
    },
    {'1': 'title', '3': 3, '4': 1, '5': 9, '10': 'title'},
    {'1': 'description', '3': 4, '4': 1, '5': 9, '10': 'description'},
    {'1': 'reason', '3': 5, '4': 1, '5': 9, '10': 'reason'},
    {
      '1': 'priority',
      '3': 6,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.RiskLevel',
      '10': 'priority'
    },
    {'1': 'action_steps', '3': 7, '4': 3, '5': 9, '10': 'actionSteps'},
    {'1': 'related_metric', '3': 8, '4': 1, '5': 9, '10': 'relatedMetric'},
    {
      '1': 'created_at',
      '3': 9,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
  ],
};

/// Descriptor for `Recommendation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recommendationDescriptor = $convert.base64Decode(
    'Cg5SZWNvbW1lbmRhdGlvbhIrChFyZWNvbW1lbmRhdGlvbl9pZBgBIAEoCVIQcmVjb21tZW5kYX'
    'Rpb25JZBIzCgR0eXBlGAIgASgOMh8ubWFucGFzaWsudjEuUmVjb21tZW5kYXRpb25UeXBlUgR0'
    'eXBlEhQKBXRpdGxlGAMgASgJUgV0aXRsZRIgCgtkZXNjcmlwdGlvbhgEIAEoCVILZGVzY3JpcH'
    'Rpb24SFgoGcmVhc29uGAUgASgJUgZyZWFzb24SMgoIcHJpb3JpdHkYBiABKA4yFi5tYW5wYXNp'
    'ay52MS5SaXNrTGV2ZWxSCHByaW9yaXR5EiEKDGFjdGlvbl9zdGVwcxgHIAMoCVILYWN0aW9uU3'
    'RlcHMSJQoOcmVsYXRlZF9tZXRyaWMYCCABKAlSDXJlbGF0ZWRNZXRyaWMSOQoKY3JlYXRlZF9h'
    'dBgJIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSCWNyZWF0ZWRBdA==');

@$core.Deprecated('Use getRecommendationsResponseDescriptor instead')
const GetRecommendationsResponse$json = {
  '1': 'GetRecommendationsResponse',
  '2': [
    {
      '1': 'recommendations',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.Recommendation',
      '10': 'recommendations'
    },
  ],
};

/// Descriptor for `GetRecommendationsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRecommendationsResponseDescriptor =
    $convert.base64Decode(
        'ChpHZXRSZWNvbW1lbmRhdGlvbnNSZXNwb25zZRJFCg9yZWNvbW1lbmRhdGlvbnMYASADKAsyGy'
        '5tYW5wYXNpay52MS5SZWNvbW1lbmRhdGlvblIPcmVjb21tZW5kYXRpb25z');

@$core.Deprecated('Use cartridgeCategoryInfoDescriptor instead')
const CartridgeCategoryInfo$json = {
  '1': 'CartridgeCategoryInfo',
  '2': [
    {'1': 'code', '3': 1, '4': 1, '5': 5, '10': 'code'},
    {'1': 'name_en', '3': 2, '4': 1, '5': 9, '10': 'nameEn'},
    {'1': 'name_ko', '3': 3, '4': 1, '5': 9, '10': 'nameKo'},
    {'1': 'description', '3': 4, '4': 1, '5': 9, '10': 'description'},
    {'1': 'type_count', '3': 5, '4': 1, '5': 5, '10': 'typeCount'},
    {'1': 'is_active', '3': 6, '4': 1, '5': 8, '10': 'isActive'},
  ],
};

/// Descriptor for `CartridgeCategoryInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cartridgeCategoryInfoDescriptor = $convert.base64Decode(
    'ChVDYXJ0cmlkZ2VDYXRlZ29yeUluZm8SEgoEY29kZRgBIAEoBVIEY29kZRIXCgduYW1lX2VuGA'
    'IgASgJUgZuYW1lRW4SFwoHbmFtZV9rbxgDIAEoCVIGbmFtZUtvEiAKC2Rlc2NyaXB0aW9uGAQg'
    'ASgJUgtkZXNjcmlwdGlvbhIdCgp0eXBlX2NvdW50GAUgASgFUgl0eXBlQ291bnQSGwoJaXNfYW'
    'N0aXZlGAYgASgIUghpc0FjdGl2ZQ==');

@$core.Deprecated('Use cartridgeTypeInfoDescriptor instead')
const CartridgeTypeInfo$json = {
  '1': 'CartridgeTypeInfo',
  '2': [
    {'1': 'category_code', '3': 1, '4': 1, '5': 5, '10': 'categoryCode'},
    {'1': 'type_index', '3': 2, '4': 1, '5': 5, '10': 'typeIndex'},
    {'1': 'legacy_code', '3': 3, '4': 1, '5': 5, '10': 'legacyCode'},
    {'1': 'name_en', '3': 4, '4': 1, '5': 9, '10': 'nameEn'},
    {'1': 'name_ko', '3': 5, '4': 1, '5': 9, '10': 'nameKo'},
    {'1': 'description', '3': 6, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'required_channels',
      '3': 7,
      '4': 1,
      '5': 5,
      '10': 'requiredChannels'
    },
    {'1': 'measurement_secs', '3': 8, '4': 1, '5': 5, '10': 'measurementSecs'},
    {'1': 'unit', '3': 9, '4': 1, '5': 9, '10': 'unit'},
    {'1': 'reference_range', '3': 10, '4': 1, '5': 9, '10': 'referenceRange'},
    {'1': 'is_active', '3': 11, '4': 1, '5': 8, '10': 'isActive'},
    {'1': 'is_beta', '3': 12, '4': 1, '5': 8, '10': 'isBeta'},
    {'1': 'manufacturer', '3': 13, '4': 1, '5': 9, '10': 'manufacturer'},
  ],
};

/// Descriptor for `CartridgeTypeInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cartridgeTypeInfoDescriptor = $convert.base64Decode(
    'ChFDYXJ0cmlkZ2VUeXBlSW5mbxIjCg1jYXRlZ29yeV9jb2RlGAEgASgFUgxjYXRlZ29yeUNvZG'
    'USHQoKdHlwZV9pbmRleBgCIAEoBVIJdHlwZUluZGV4Eh8KC2xlZ2FjeV9jb2RlGAMgASgFUgps'
    'ZWdhY3lDb2RlEhcKB25hbWVfZW4YBCABKAlSBm5hbWVFbhIXCgduYW1lX2tvGAUgASgJUgZuYW'
    '1lS28SIAoLZGVzY3JpcHRpb24YBiABKAlSC2Rlc2NyaXB0aW9uEisKEXJlcXVpcmVkX2NoYW5u'
    'ZWxzGAcgASgFUhByZXF1aXJlZENoYW5uZWxzEikKEG1lYXN1cmVtZW50X3NlY3MYCCABKAVSD2'
    '1lYXN1cmVtZW50U2VjcxISCgR1bml0GAkgASgJUgR1bml0EicKD3JlZmVyZW5jZV9yYW5nZRgK'
    'IAEoCVIOcmVmZXJlbmNlUmFuZ2USGwoJaXNfYWN0aXZlGAsgASgIUghpc0FjdGl2ZRIXCgdpc1'
    '9iZXRhGAwgASgIUgZpc0JldGESIgoMbWFudWZhY3R1cmVyGA0gASgJUgxtYW51ZmFjdHVyZXI=');

@$core.Deprecated('Use checkCartridgeAccessRequestDescriptor instead')
const CheckCartridgeAccessRequest$json = {
  '1': 'CheckCartridgeAccessRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'category_code', '3': 2, '4': 1, '5': 5, '10': 'categoryCode'},
    {'1': 'type_index', '3': 3, '4': 1, '5': 5, '10': 'typeIndex'},
  ],
};

/// Descriptor for `CheckCartridgeAccessRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List checkCartridgeAccessRequestDescriptor =
    $convert.base64Decode(
        'ChtDaGVja0NhcnRyaWRnZUFjY2Vzc1JlcXVlc3QSFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklkEi'
        'MKDWNhdGVnb3J5X2NvZGUYAiABKAVSDGNhdGVnb3J5Q29kZRIdCgp0eXBlX2luZGV4GAMgASgF'
        'Ugl0eXBlSW5kZXg=');

@$core.Deprecated('Use checkCartridgeAccessResponseDescriptor instead')
const CheckCartridgeAccessResponse$json = {
  '1': 'CheckCartridgeAccessResponse',
  '2': [
    {'1': 'allowed', '3': 1, '4': 1, '5': 8, '10': 'allowed'},
    {
      '1': 'access_level',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.CartridgeAccessLevel',
      '10': 'accessLevel'
    },
    {'1': 'remaining_daily', '3': 3, '4': 1, '5': 5, '10': 'remainingDaily'},
    {
      '1': 'remaining_monthly',
      '3': 4,
      '4': 1,
      '5': 5,
      '10': 'remainingMonthly'
    },
    {
      '1': 'required_tier',
      '3': 5,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.SubscriptionTier',
      '10': 'requiredTier'
    },
    {
      '1': 'current_tier',
      '3': 6,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.SubscriptionTier',
      '10': 'currentTier'
    },
    {'1': 'message', '3': 7, '4': 1, '5': 9, '10': 'message'},
    {'1': 'addon_price_krw', '3': 8, '4': 1, '5': 5, '10': 'addonPriceKrw'},
  ],
};

/// Descriptor for `CheckCartridgeAccessResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List checkCartridgeAccessResponseDescriptor = $convert.base64Decode(
    'ChxDaGVja0NhcnRyaWRnZUFjY2Vzc1Jlc3BvbnNlEhgKB2FsbG93ZWQYASABKAhSB2FsbG93ZW'
    'QSRAoMYWNjZXNzX2xldmVsGAIgASgOMiEubWFucGFzaWsudjEuQ2FydHJpZGdlQWNjZXNzTGV2'
    'ZWxSC2FjY2Vzc0xldmVsEicKD3JlbWFpbmluZ19kYWlseRgDIAEoBVIOcmVtYWluaW5nRGFpbH'
    'kSKwoRcmVtYWluaW5nX21vbnRobHkYBCABKAVSEHJlbWFpbmluZ01vbnRobHkSQgoNcmVxdWly'
    'ZWRfdGllchgFIAEoDjIdLm1hbnBhc2lrLnYxLlN1YnNjcmlwdGlvblRpZXJSDHJlcXVpcmVkVG'
    'llchJACgxjdXJyZW50X3RpZXIYBiABKA4yHS5tYW5wYXNpay52MS5TdWJzY3JpcHRpb25UaWVy'
    'UgtjdXJyZW50VGllchIYCgdtZXNzYWdlGAcgASgJUgdtZXNzYWdlEiYKD2FkZG9uX3ByaWNlX2'
    'tydxgIIAEoBVINYWRkb25QcmljZUtydw==');

@$core.Deprecated('Use listAccessibleCartridgesRequestDescriptor instead')
const ListAccessibleCartridgesRequest$json = {
  '1': 'ListAccessibleCartridgesRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `ListAccessibleCartridgesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAccessibleCartridgesRequestDescriptor =
    $convert.base64Decode(
        'Ch9MaXN0QWNjZXNzaWJsZUNhcnRyaWRnZXNSZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZX'
        'JJZA==');

@$core.Deprecated('Use listAccessibleCartridgesResponseDescriptor instead')
const ListAccessibleCartridgesResponse$json = {
  '1': 'ListAccessibleCartridgesResponse',
  '2': [
    {
      '1': 'entries',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.CartridgeAccessEntry',
      '10': 'entries'
    },
  ],
};

/// Descriptor for `ListAccessibleCartridgesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAccessibleCartridgesResponseDescriptor =
    $convert.base64Decode(
        'CiBMaXN0QWNjZXNzaWJsZUNhcnRyaWRnZXNSZXNwb25zZRI7CgdlbnRyaWVzGAEgAygLMiEubW'
        'FucGFzaWsudjEuQ2FydHJpZGdlQWNjZXNzRW50cnlSB2VudHJpZXM=');

@$core.Deprecated('Use cartridgeAccessEntryDescriptor instead')
const CartridgeAccessEntry$json = {
  '1': 'CartridgeAccessEntry',
  '2': [
    {
      '1': 'type_info',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.manpasik.v1.CartridgeTypeInfo',
      '10': 'typeInfo'
    },
    {
      '1': 'access_level',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.CartridgeAccessLevel',
      '10': 'accessLevel'
    },
    {'1': 'remaining_daily', '3': 3, '4': 1, '5': 5, '10': 'remainingDaily'},
    {
      '1': 'remaining_monthly',
      '3': 4,
      '4': 1,
      '5': 5,
      '10': 'remainingMonthly'
    },
  ],
};

/// Descriptor for `CartridgeAccessEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cartridgeAccessEntryDescriptor = $convert.base64Decode(
    'ChRDYXJ0cmlkZ2VBY2Nlc3NFbnRyeRI7Cgl0eXBlX2luZm8YASABKAsyHi5tYW5wYXNpay52MS'
    '5DYXJ0cmlkZ2VUeXBlSW5mb1IIdHlwZUluZm8SRAoMYWNjZXNzX2xldmVsGAIgASgOMiEubWFu'
    'cGFzaWsudjEuQ2FydHJpZGdlQWNjZXNzTGV2ZWxSC2FjY2Vzc0xldmVsEicKD3JlbWFpbmluZ1'
    '9kYWlseRgDIAEoBVIOcmVtYWluaW5nRGFpbHkSKwoRcmVtYWluaW5nX21vbnRobHkYBCABKAVS'
    'EHJlbWFpbmluZ01vbnRobHk=');

@$core.Deprecated('Use searchFacilitiesRequestDescriptor instead')
const SearchFacilitiesRequest$json = {
  '1': 'SearchFacilitiesRequest',
  '2': [
    {'1': 'latitude', '3': 1, '4': 1, '5': 1, '10': 'latitude'},
    {'1': 'longitude', '3': 2, '4': 1, '5': 1, '10': 'longitude'},
    {'1': 'radius_km', '3': 3, '4': 1, '5': 1, '10': 'radiusKm'},
    {
      '1': 'type',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.FacilityType',
      '10': 'type'
    },
    {
      '1': 'specialty',
      '3': 5,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.DoctorSpecialty',
      '10': 'specialty'
    },
    {'1': 'query', '3': 6, '4': 1, '5': 9, '10': 'query'},
    {'1': 'limit', '3': 7, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 8, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `SearchFacilitiesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List searchFacilitiesRequestDescriptor = $convert.base64Decode(
    'ChdTZWFyY2hGYWNpbGl0aWVzUmVxdWVzdBIaCghsYXRpdHVkZRgBIAEoAVIIbGF0aXR1ZGUSHA'
    'oJbG9uZ2l0dWRlGAIgASgBUglsb25naXR1ZGUSGwoJcmFkaXVzX2ttGAMgASgBUghyYWRpdXNL'
    'bRItCgR0eXBlGAQgASgOMhkubWFucGFzaWsudjEuRmFjaWxpdHlUeXBlUgR0eXBlEjoKCXNwZW'
    'NpYWx0eRgFIAEoDjIcLm1hbnBhc2lrLnYxLkRvY3RvclNwZWNpYWx0eVIJc3BlY2lhbHR5EhQK'
    'BXF1ZXJ5GAYgASgJUgVxdWVyeRIUCgVsaW1pdBgHIAEoBVIFbGltaXQSFgoGb2Zmc2V0GAggAS'
    'gFUgZvZmZzZXQ=');

@$core.Deprecated('Use searchFacilitiesResponseDescriptor instead')
const SearchFacilitiesResponse$json = {
  '1': 'SearchFacilitiesResponse',
  '2': [
    {
      '1': 'facilities',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.Facility',
      '10': 'facilities'
    },
    {'1': 'total_count', '3': 2, '4': 1, '5': 5, '10': 'totalCount'},
  ],
};

/// Descriptor for `SearchFacilitiesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List searchFacilitiesResponseDescriptor = $convert.base64Decode(
    'ChhTZWFyY2hGYWNpbGl0aWVzUmVzcG9uc2USNQoKZmFjaWxpdGllcxgBIAMoCzIVLm1hbnBhc2'
    'lrLnYxLkZhY2lsaXR5UgpmYWNpbGl0aWVzEh8KC3RvdGFsX2NvdW50GAIgASgFUgp0b3RhbENv'
    'dW50');

@$core.Deprecated('Use getFacilityRequestDescriptor instead')
const GetFacilityRequest$json = {
  '1': 'GetFacilityRequest',
  '2': [
    {'1': 'facility_id', '3': 1, '4': 1, '5': 9, '10': 'facilityId'},
  ],
};

/// Descriptor for `GetFacilityRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFacilityRequestDescriptor = $convert.base64Decode(
    'ChJHZXRGYWNpbGl0eVJlcXVlc3QSHwoLZmFjaWxpdHlfaWQYASABKAlSCmZhY2lsaXR5SWQ=');

@$core.Deprecated('Use facilityDescriptor instead')
const Facility$json = {
  '1': 'Facility',
  '2': [
    {'1': 'facility_id', '3': 1, '4': 1, '5': 9, '10': 'facilityId'},
    {'1': 'name', '3': 2, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'type',
      '3': 3,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.FacilityType',
      '10': 'type'
    },
    {'1': 'address', '3': 4, '4': 1, '5': 9, '10': 'address'},
    {'1': 'latitude', '3': 5, '4': 1, '5': 1, '10': 'latitude'},
    {'1': 'longitude', '3': 6, '4': 1, '5': 1, '10': 'longitude'},
    {'1': 'phone', '3': 7, '4': 1, '5': 9, '10': 'phone'},
    {
      '1': 'specialties',
      '3': 8,
      '4': 3,
      '5': 14,
      '6': '.manpasik.v1.DoctorSpecialty',
      '10': 'specialties'
    },
    {'1': 'rating', '3': 9, '4': 1, '5': 2, '10': 'rating'},
    {'1': 'is_open_now', '3': 10, '4': 1, '5': 8, '10': 'isOpenNow'},
    {
      '1': 'accepts_reservation',
      '3': 11,
      '4': 1,
      '5': 8,
      '10': 'acceptsReservation'
    },
    {'1': 'operating_hours', '3': 12, '4': 1, '5': 9, '10': 'operatingHours'},
    {'1': 'distance_km', '3': 13, '4': 1, '5': 1, '10': 'distanceKm'},
  ],
};

/// Descriptor for `Facility`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List facilityDescriptor = $convert.base64Decode(
    'CghGYWNpbGl0eRIfCgtmYWNpbGl0eV9pZBgBIAEoCVIKZmFjaWxpdHlJZBISCgRuYW1lGAIgAS'
    'gJUgRuYW1lEi0KBHR5cGUYAyABKA4yGS5tYW5wYXNpay52MS5GYWNpbGl0eVR5cGVSBHR5cGUS'
    'GAoHYWRkcmVzcxgEIAEoCVIHYWRkcmVzcxIaCghsYXRpdHVkZRgFIAEoAVIIbGF0aXR1ZGUSHA'
    'oJbG9uZ2l0dWRlGAYgASgBUglsb25naXR1ZGUSFAoFcGhvbmUYByABKAlSBXBob25lEj4KC3Nw'
    'ZWNpYWx0aWVzGAggAygOMhwubWFucGFzaWsudjEuRG9jdG9yU3BlY2lhbHR5UgtzcGVjaWFsdG'
    'llcxIWCgZyYXRpbmcYCSABKAJSBnJhdGluZxIeCgtpc19vcGVuX25vdxgKIAEoCFIJaXNPcGVu'
    'Tm93Ei8KE2FjY2VwdHNfcmVzZXJ2YXRpb24YCyABKAhSEmFjY2VwdHNSZXNlcnZhdGlvbhInCg'
    '9vcGVyYXRpbmdfaG91cnMYDCABKAlSDm9wZXJhdGluZ0hvdXJzEh8KC2Rpc3RhbmNlX2ttGA0g'
    'ASgBUgpkaXN0YW5jZUtt');

@$core.Deprecated('Use getAvailableSlotsRequestDescriptor instead')
const GetAvailableSlotsRequest$json = {
  '1': 'GetAvailableSlotsRequest',
  '2': [
    {'1': 'facility_id', '3': 1, '4': 1, '5': 9, '10': 'facilityId'},
    {'1': 'date', '3': 2, '4': 1, '5': 9, '10': 'date'},
    {'1': 'doctor_id', '3': 3, '4': 1, '5': 9, '10': 'doctorId'},
    {
      '1': 'specialty',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.DoctorSpecialty',
      '10': 'specialty'
    },
  ],
};

/// Descriptor for `GetAvailableSlotsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAvailableSlotsRequestDescriptor = $convert.base64Decode(
    'ChhHZXRBdmFpbGFibGVTbG90c1JlcXVlc3QSHwoLZmFjaWxpdHlfaWQYASABKAlSCmZhY2lsaX'
    'R5SWQSEgoEZGF0ZRgCIAEoCVIEZGF0ZRIbCglkb2N0b3JfaWQYAyABKAlSCGRvY3RvcklkEjoK'
    'CXNwZWNpYWx0eRgEIAEoDjIcLm1hbnBhc2lrLnYxLkRvY3RvclNwZWNpYWx0eVIJc3BlY2lhbH'
    'R5');

@$core.Deprecated('Use getAvailableSlotsResponseDescriptor instead')
const GetAvailableSlotsResponse$json = {
  '1': 'GetAvailableSlotsResponse',
  '2': [
    {
      '1': 'slots',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.TimeSlot',
      '10': 'slots'
    },
  ],
};

/// Descriptor for `GetAvailableSlotsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAvailableSlotsResponseDescriptor =
    $convert.base64Decode(
        'ChlHZXRBdmFpbGFibGVTbG90c1Jlc3BvbnNlEisKBXNsb3RzGAEgAygLMhUubWFucGFzaWsudj'
        'EuVGltZVNsb3RSBXNsb3Rz');

@$core.Deprecated('Use timeSlotDescriptor instead')
const TimeSlot$json = {
  '1': 'TimeSlot',
  '2': [
    {'1': 'slot_id', '3': 1, '4': 1, '5': 9, '10': 'slotId'},
    {
      '1': 'start_time',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'startTime'
    },
    {
      '1': 'end_time',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'endTime'
    },
    {'1': 'is_available', '3': 4, '4': 1, '5': 8, '10': 'isAvailable'},
    {'1': 'doctor_name', '3': 5, '4': 1, '5': 9, '10': 'doctorName'},
    {
      '1': 'specialty',
      '3': 6,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.DoctorSpecialty',
      '10': 'specialty'
    },
  ],
};

/// Descriptor for `TimeSlot`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List timeSlotDescriptor = $convert.base64Decode(
    'CghUaW1lU2xvdBIXCgdzbG90X2lkGAEgASgJUgZzbG90SWQSOQoKc3RhcnRfdGltZRgCIAEoCz'
    'IaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSCXN0YXJ0VGltZRI1CghlbmRfdGltZRgDIAEo'
    'CzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSB2VuZFRpbWUSIQoMaXNfYXZhaWxhYmxlGA'
    'QgASgIUgtpc0F2YWlsYWJsZRIfCgtkb2N0b3JfbmFtZRgFIAEoCVIKZG9jdG9yTmFtZRI6Cglz'
    'cGVjaWFsdHkYBiABKA4yHC5tYW5wYXNpay52MS5Eb2N0b3JTcGVjaWFsdHlSCXNwZWNpYWx0eQ'
    '==');

@$core.Deprecated('Use createReservationRequestDescriptor instead')
const CreateReservationRequest$json = {
  '1': 'CreateReservationRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'facility_id', '3': 2, '4': 1, '5': 9, '10': 'facilityId'},
    {'1': 'slot_id', '3': 3, '4': 1, '5': 9, '10': 'slotId'},
    {'1': 'doctor_id', '3': 4, '4': 1, '5': 9, '10': 'doctorId'},
    {
      '1': 'specialty',
      '3': 5,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.DoctorSpecialty',
      '10': 'specialty'
    },
    {'1': 'reason', '3': 6, '4': 1, '5': 9, '10': 'reason'},
    {'1': 'notes', '3': 7, '4': 1, '5': 9, '10': 'notes'},
  ],
};

/// Descriptor for `CreateReservationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createReservationRequestDescriptor = $convert.base64Decode(
    'ChhDcmVhdGVSZXNlcnZhdGlvblJlcXVlc3QSFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklkEh8KC2'
    'ZhY2lsaXR5X2lkGAIgASgJUgpmYWNpbGl0eUlkEhcKB3Nsb3RfaWQYAyABKAlSBnNsb3RJZBIb'
    'Cglkb2N0b3JfaWQYBCABKAlSCGRvY3RvcklkEjoKCXNwZWNpYWx0eRgFIAEoDjIcLm1hbnBhc2'
    'lrLnYxLkRvY3RvclNwZWNpYWx0eVIJc3BlY2lhbHR5EhYKBnJlYXNvbhgGIAEoCVIGcmVhc29u'
    'EhQKBW5vdGVzGAcgASgJUgVub3Rlcw==');

@$core.Deprecated('Use reservationDescriptor instead')
const Reservation$json = {
  '1': 'Reservation',
  '2': [
    {'1': 'reservation_id', '3': 1, '4': 1, '5': 9, '10': 'reservationId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'facility_id', '3': 3, '4': 1, '5': 9, '10': 'facilityId'},
    {'1': 'facility_name', '3': 4, '4': 1, '5': 9, '10': 'facilityName'},
    {'1': 'doctor_id', '3': 5, '4': 1, '5': 9, '10': 'doctorId'},
    {'1': 'doctor_name', '3': 6, '4': 1, '5': 9, '10': 'doctorName'},
    {
      '1': 'specialty',
      '3': 7,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.DoctorSpecialty',
      '10': 'specialty'
    },
    {
      '1': 'appointment_time',
      '3': 8,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'appointmentTime'
    },
    {
      '1': 'status',
      '3': 9,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.ReservationStatus',
      '10': 'status'
    },
    {'1': 'reason', '3': 10, '4': 1, '5': 9, '10': 'reason'},
    {'1': 'notes', '3': 11, '4': 1, '5': 9, '10': 'notes'},
    {
      '1': 'created_at',
      '3': 12,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
    {
      '1': 'updated_at',
      '3': 13,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'updatedAt'
    },
  ],
};

/// Descriptor for `Reservation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List reservationDescriptor = $convert.base64Decode(
    'CgtSZXNlcnZhdGlvbhIlCg5yZXNlcnZhdGlvbl9pZBgBIAEoCVINcmVzZXJ2YXRpb25JZBIXCg'
    'd1c2VyX2lkGAIgASgJUgZ1c2VySWQSHwoLZmFjaWxpdHlfaWQYAyABKAlSCmZhY2lsaXR5SWQS'
    'IwoNZmFjaWxpdHlfbmFtZRgEIAEoCVIMZmFjaWxpdHlOYW1lEhsKCWRvY3Rvcl9pZBgFIAEoCV'
    'IIZG9jdG9ySWQSHwoLZG9jdG9yX25hbWUYBiABKAlSCmRvY3Rvck5hbWUSOgoJc3BlY2lhbHR5'
    'GAcgASgOMhwubWFucGFzaWsudjEuRG9jdG9yU3BlY2lhbHR5UglzcGVjaWFsdHkSRQoQYXBwb2'
    'ludG1lbnRfdGltZRgIIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSD2FwcG9pbnRt'
    'ZW50VGltZRI2CgZzdGF0dXMYCSABKA4yHi5tYW5wYXNpay52MS5SZXNlcnZhdGlvblN0YXR1c1'
    'IGc3RhdHVzEhYKBnJlYXNvbhgKIAEoCVIGcmVhc29uEhQKBW5vdGVzGAsgASgJUgVub3RlcxI5'
    'CgpjcmVhdGVkX2F0GAwgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFIJY3JlYXRlZE'
    'F0EjkKCnVwZGF0ZWRfYXQYDSABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUgl1cGRh'
    'dGVkQXQ=');

@$core.Deprecated('Use getReservationRequestDescriptor instead')
const GetReservationRequest$json = {
  '1': 'GetReservationRequest',
  '2': [
    {'1': 'reservation_id', '3': 1, '4': 1, '5': 9, '10': 'reservationId'},
  ],
};

/// Descriptor for `GetReservationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getReservationRequestDescriptor = $convert.base64Decode(
    'ChVHZXRSZXNlcnZhdGlvblJlcXVlc3QSJQoOcmVzZXJ2YXRpb25faWQYASABKAlSDXJlc2Vydm'
    'F0aW9uSWQ=');

@$core.Deprecated('Use listReservationsRequestDescriptor instead')
const ListReservationsRequest$json = {
  '1': 'ListReservationsRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'status',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.ReservationStatus',
      '10': 'status'
    },
    {'1': 'limit', '3': 3, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 4, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `ListReservationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listReservationsRequestDescriptor = $convert.base64Decode(
    'ChdMaXN0UmVzZXJ2YXRpb25zUmVxdWVzdBIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQSNgoGc3'
    'RhdHVzGAIgASgOMh4ubWFucGFzaWsudjEuUmVzZXJ2YXRpb25TdGF0dXNSBnN0YXR1cxIUCgVs'
    'aW1pdBgDIAEoBVIFbGltaXQSFgoGb2Zmc2V0GAQgASgFUgZvZmZzZXQ=');

@$core.Deprecated('Use listReservationsResponseDescriptor instead')
const ListReservationsResponse$json = {
  '1': 'ListReservationsResponse',
  '2': [
    {
      '1': 'reservations',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.Reservation',
      '10': 'reservations'
    },
    {'1': 'total_count', '3': 2, '4': 1, '5': 5, '10': 'totalCount'},
  ],
};

/// Descriptor for `ListReservationsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listReservationsResponseDescriptor = $convert.base64Decode(
    'ChhMaXN0UmVzZXJ2YXRpb25zUmVzcG9uc2USPAoMcmVzZXJ2YXRpb25zGAEgAygLMhgubWFucG'
    'FzaWsudjEuUmVzZXJ2YXRpb25SDHJlc2VydmF0aW9ucxIfCgt0b3RhbF9jb3VudBgCIAEoBVIK'
    'dG90YWxDb3VudA==');

@$core.Deprecated('Use cancelReservationRequestDescriptor instead')
const CancelReservationRequest$json = {
  '1': 'CancelReservationRequest',
  '2': [
    {'1': 'reservation_id', '3': 1, '4': 1, '5': 9, '10': 'reservationId'},
    {'1': 'reason', '3': 2, '4': 1, '5': 9, '10': 'reason'},
  ],
};

/// Descriptor for `CancelReservationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelReservationRequestDescriptor =
    $convert.base64Decode(
        'ChhDYW5jZWxSZXNlcnZhdGlvblJlcXVlc3QSJQoOcmVzZXJ2YXRpb25faWQYASABKAlSDXJlc2'
        'VydmF0aW9uSWQSFgoGcmVhc29uGAIgASgJUgZyZWFzb24=');

@$core.Deprecated('Use cancelReservationResponseDescriptor instead')
const CancelReservationResponse$json = {
  '1': 'CancelReservationResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
    {'1': 'message', '3': 2, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CancelReservationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelReservationResponseDescriptor =
    $convert.base64Decode(
        'ChlDYW5jZWxSZXNlcnZhdGlvblJlc3BvbnNlEhgKB3N1Y2Nlc3MYASABKAhSB3N1Y2Nlc3MSGA'
        'oHbWVzc2FnZRgCIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use createAdminRequestDescriptor instead')
const CreateAdminRequest$json = {
  '1': 'CreateAdminRequest',
  '2': [
    {'1': 'email', '3': 1, '4': 1, '5': 9, '10': 'email'},
    {'1': 'password', '3': 2, '4': 1, '5': 9, '10': 'password'},
    {'1': 'display_name', '3': 3, '4': 1, '5': 9, '10': 'displayName'},
    {
      '1': 'role',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.AdminRole',
      '10': 'role'
    },
    {'1': 'region', '3': 5, '4': 1, '5': 9, '10': 'region'},
    {'1': 'branch', '3': 6, '4': 1, '5': 9, '10': 'branch'},
  ],
};

/// Descriptor for `CreateAdminRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createAdminRequestDescriptor = $convert.base64Decode(
    'ChJDcmVhdGVBZG1pblJlcXVlc3QSFAoFZW1haWwYASABKAlSBWVtYWlsEhoKCHBhc3N3b3JkGA'
    'IgASgJUghwYXNzd29yZBIhCgxkaXNwbGF5X25hbWUYAyABKAlSC2Rpc3BsYXlOYW1lEioKBHJv'
    'bGUYBCABKA4yFi5tYW5wYXNpay52MS5BZG1pblJvbGVSBHJvbGUSFgoGcmVnaW9uGAUgASgJUg'
    'ZyZWdpb24SFgoGYnJhbmNoGAYgASgJUgZicmFuY2g=');

@$core.Deprecated('Use getAdminRequestDescriptor instead')
const GetAdminRequest$json = {
  '1': 'GetAdminRequest',
  '2': [
    {'1': 'admin_id', '3': 1, '4': 1, '5': 9, '10': 'adminId'},
  ],
};

/// Descriptor for `GetAdminRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAdminRequestDescriptor = $convert.base64Decode(
    'Cg9HZXRBZG1pblJlcXVlc3QSGQoIYWRtaW5faWQYASABKAlSB2FkbWluSWQ=');

@$core.Deprecated('Use listAdminsRequestDescriptor instead')
const ListAdminsRequest$json = {
  '1': 'ListAdminsRequest',
  '2': [
    {
      '1': 'role_filter',
      '3': 1,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.AdminRole',
      '10': 'roleFilter'
    },
    {'1': 'region', '3': 2, '4': 1, '5': 9, '10': 'region'},
    {'1': 'limit', '3': 3, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 4, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `ListAdminsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAdminsRequestDescriptor = $convert.base64Decode(
    'ChFMaXN0QWRtaW5zUmVxdWVzdBI3Cgtyb2xlX2ZpbHRlchgBIAEoDjIWLm1hbnBhc2lrLnYxLk'
    'FkbWluUm9sZVIKcm9sZUZpbHRlchIWCgZyZWdpb24YAiABKAlSBnJlZ2lvbhIUCgVsaW1pdBgD'
    'IAEoBVIFbGltaXQSFgoGb2Zmc2V0GAQgASgFUgZvZmZzZXQ=');

@$core.Deprecated('Use listAdminsResponseDescriptor instead')
const ListAdminsResponse$json = {
  '1': 'ListAdminsResponse',
  '2': [
    {
      '1': 'admins',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.AdminUser',
      '10': 'admins'
    },
    {'1': 'total_count', '3': 2, '4': 1, '5': 5, '10': 'totalCount'},
  ],
};

/// Descriptor for `ListAdminsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAdminsResponseDescriptor = $convert.base64Decode(
    'ChJMaXN0QWRtaW5zUmVzcG9uc2USLgoGYWRtaW5zGAEgAygLMhYubWFucGFzaWsudjEuQWRtaW'
    '5Vc2VyUgZhZG1pbnMSHwoLdG90YWxfY291bnQYAiABKAVSCnRvdGFsQ291bnQ=');

@$core.Deprecated('Use updateAdminRoleRequestDescriptor instead')
const UpdateAdminRoleRequest$json = {
  '1': 'UpdateAdminRoleRequest',
  '2': [
    {'1': 'admin_id', '3': 1, '4': 1, '5': 9, '10': 'adminId'},
    {
      '1': 'new_role',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.AdminRole',
      '10': 'newRole'
    },
    {'1': 'region', '3': 3, '4': 1, '5': 9, '10': 'region'},
    {'1': 'branch', '3': 4, '4': 1, '5': 9, '10': 'branch'},
  ],
};

/// Descriptor for `UpdateAdminRoleRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateAdminRoleRequestDescriptor = $convert.base64Decode(
    'ChZVcGRhdGVBZG1pblJvbGVSZXF1ZXN0EhkKCGFkbWluX2lkGAEgASgJUgdhZG1pbklkEjEKCG'
    '5ld19yb2xlGAIgASgOMhYubWFucGFzaWsudjEuQWRtaW5Sb2xlUgduZXdSb2xlEhYKBnJlZ2lv'
    'bhgDIAEoCVIGcmVnaW9uEhYKBmJyYW5jaBgEIAEoCVIGYnJhbmNo');

@$core.Deprecated('Use deactivateAdminRequestDescriptor instead')
const DeactivateAdminRequest$json = {
  '1': 'DeactivateAdminRequest',
  '2': [
    {'1': 'admin_id', '3': 1, '4': 1, '5': 9, '10': 'adminId'},
    {'1': 'reason', '3': 2, '4': 1, '5': 9, '10': 'reason'},
  ],
};

/// Descriptor for `DeactivateAdminRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deactivateAdminRequestDescriptor =
    $convert.base64Decode(
        'ChZEZWFjdGl2YXRlQWRtaW5SZXF1ZXN0EhkKCGFkbWluX2lkGAEgASgJUgdhZG1pbklkEhYKBn'
        'JlYXNvbhgCIAEoCVIGcmVhc29u');

@$core.Deprecated('Use adminUserDescriptor instead')
const AdminUser$json = {
  '1': 'AdminUser',
  '2': [
    {'1': 'admin_id', '3': 1, '4': 1, '5': 9, '10': 'adminId'},
    {'1': 'email', '3': 2, '4': 1, '5': 9, '10': 'email'},
    {'1': 'display_name', '3': 3, '4': 1, '5': 9, '10': 'displayName'},
    {
      '1': 'role',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.AdminRole',
      '10': 'role'
    },
    {'1': 'region', '3': 5, '4': 1, '5': 9, '10': 'region'},
    {'1': 'branch', '3': 6, '4': 1, '5': 9, '10': 'branch'},
    {'1': 'is_active', '3': 7, '4': 1, '5': 8, '10': 'isActive'},
    {
      '1': 'created_at',
      '3': 8,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
    {
      '1': 'last_login_at',
      '3': 9,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'lastLoginAt'
    },
  ],
};

/// Descriptor for `AdminUser`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminUserDescriptor = $convert.base64Decode(
    'CglBZG1pblVzZXISGQoIYWRtaW5faWQYASABKAlSB2FkbWluSWQSFAoFZW1haWwYAiABKAlSBW'
    'VtYWlsEiEKDGRpc3BsYXlfbmFtZRgDIAEoCVILZGlzcGxheU5hbWUSKgoEcm9sZRgEIAEoDjIW'
    'Lm1hbnBhc2lrLnYxLkFkbWluUm9sZVIEcm9sZRIWCgZyZWdpb24YBSABKAlSBnJlZ2lvbhIWCg'
    'ZicmFuY2gYBiABKAlSBmJyYW5jaBIbCglpc19hY3RpdmUYByABKAhSCGlzQWN0aXZlEjkKCmNy'
    'ZWF0ZWRfYXQYCCABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUgljcmVhdGVkQXQSPg'
    'oNbGFzdF9sb2dpbl9hdBgJIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSC2xhc3RM'
    'b2dpbkF0');

@$core.Deprecated('Use adminListUsersRequestDescriptor instead')
const AdminListUsersRequest$json = {
  '1': 'AdminListUsersRequest',
  '2': [
    {'1': 'query', '3': 1, '4': 1, '5': 9, '10': 'query'},
    {
      '1': 'tier_filter',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.SubscriptionTier',
      '10': 'tierFilter'
    },
    {'1': 'active_only', '3': 3, '4': 1, '5': 8, '10': 'activeOnly'},
    {'1': 'limit', '3': 4, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 5, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `AdminListUsersRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminListUsersRequestDescriptor = $convert.base64Decode(
    'ChVBZG1pbkxpc3RVc2Vyc1JlcXVlc3QSFAoFcXVlcnkYASABKAlSBXF1ZXJ5Ej4KC3RpZXJfZm'
    'lsdGVyGAIgASgOMh0ubWFucGFzaWsudjEuU3Vic2NyaXB0aW9uVGllclIKdGllckZpbHRlchIf'
    'CgthY3RpdmVfb25seRgDIAEoCFIKYWN0aXZlT25seRIUCgVsaW1pdBgEIAEoBVIFbGltaXQSFg'
    'oGb2Zmc2V0GAUgASgFUgZvZmZzZXQ=');

@$core.Deprecated('Use adminListUsersResponseDescriptor instead')
const AdminListUsersResponse$json = {
  '1': 'AdminListUsersResponse',
  '2': [
    {
      '1': 'users',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.AdminUserSummary',
      '10': 'users'
    },
    {'1': 'total_count', '3': 2, '4': 1, '5': 5, '10': 'totalCount'},
  ],
};

/// Descriptor for `AdminListUsersResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminListUsersResponseDescriptor = $convert.base64Decode(
    'ChZBZG1pbkxpc3RVc2Vyc1Jlc3BvbnNlEjMKBXVzZXJzGAEgAygLMh0ubWFucGFzaWsudjEuQW'
    'RtaW5Vc2VyU3VtbWFyeVIFdXNlcnMSHwoLdG90YWxfY291bnQYAiABKAVSCnRvdGFsQ291bnQ=');

@$core.Deprecated('Use adminUserSummaryDescriptor instead')
const AdminUserSummary$json = {
  '1': 'AdminUserSummary',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'email', '3': 2, '4': 1, '5': 9, '10': 'email'},
    {'1': 'display_name', '3': 3, '4': 1, '5': 9, '10': 'displayName'},
    {
      '1': 'tier',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.SubscriptionTier',
      '10': 'tier'
    },
    {'1': 'is_active', '3': 5, '4': 1, '5': 8, '10': 'isActive'},
    {'1': 'device_count', '3': 6, '4': 1, '5': 5, '10': 'deviceCount'},
    {
      '1': 'measurement_count',
      '3': 7,
      '4': 1,
      '5': 5,
      '10': 'measurementCount'
    },
    {
      '1': 'created_at',
      '3': 8,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
    {
      '1': 'last_active_at',
      '3': 9,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'lastActiveAt'
    },
  ],
};

/// Descriptor for `AdminUserSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List adminUserSummaryDescriptor = $convert.base64Decode(
    'ChBBZG1pblVzZXJTdW1tYXJ5EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZBIUCgVlbWFpbBgCIA'
    'EoCVIFZW1haWwSIQoMZGlzcGxheV9uYW1lGAMgASgJUgtkaXNwbGF5TmFtZRIxCgR0aWVyGAQg'
    'ASgOMh0ubWFucGFzaWsudjEuU3Vic2NyaXB0aW9uVGllclIEdGllchIbCglpc19hY3RpdmUYBS'
    'ABKAhSCGlzQWN0aXZlEiEKDGRldmljZV9jb3VudBgGIAEoBVILZGV2aWNlQ291bnQSKwoRbWVh'
    'c3VyZW1lbnRfY291bnQYByABKAVSEG1lYXN1cmVtZW50Q291bnQSOQoKY3JlYXRlZF9hdBgIIA'
    'EoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSCWNyZWF0ZWRBdBJACg5sYXN0X2FjdGl2'
    'ZV9hdBgJIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSDGxhc3RBY3RpdmVBdA==');

@$core.Deprecated('Use getSystemStatsRequestDescriptor instead')
const GetSystemStatsRequest$json = {
  '1': 'GetSystemStatsRequest',
};

/// Descriptor for `GetSystemStatsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSystemStatsRequestDescriptor =
    $convert.base64Decode('ChVHZXRTeXN0ZW1TdGF0c1JlcXVlc3Q=');

@$core.Deprecated('Use getSystemStatsResponseDescriptor instead')
const GetSystemStatsResponse$json = {
  '1': 'GetSystemStatsResponse',
  '2': [
    {'1': 'total_users', '3': 1, '4': 1, '5': 5, '10': 'totalUsers'},
    {'1': 'active_users', '3': 2, '4': 1, '5': 5, '10': 'activeUsers'},
    {'1': 'total_devices', '3': 3, '4': 1, '5': 5, '10': 'totalDevices'},
    {
      '1': 'total_measurements',
      '3': 4,
      '4': 1,
      '5': 5,
      '10': 'totalMeasurements'
    },
    {
      '1': 'users_by_tier',
      '3': 5,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.GetSystemStatsResponse.UsersByTierEntry',
      '10': 'usersByTier'
    },
    {
      '1': 'measurements_by_type',
      '3': 6,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.GetSystemStatsResponse.MeasurementsByTypeEntry',
      '10': 'measurementsByType'
    },
    {
      '1': 'generated_at',
      '3': 7,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'generatedAt'
    },
  ],
  '3': [
    GetSystemStatsResponse_UsersByTierEntry$json,
    GetSystemStatsResponse_MeasurementsByTypeEntry$json
  ],
};

@$core.Deprecated('Use getSystemStatsResponseDescriptor instead')
const GetSystemStatsResponse_UsersByTierEntry$json = {
  '1': 'UsersByTierEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 5, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use getSystemStatsResponseDescriptor instead')
const GetSystemStatsResponse_MeasurementsByTypeEntry$json = {
  '1': 'MeasurementsByTypeEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 5, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetSystemStatsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSystemStatsResponseDescriptor = $convert.base64Decode(
    'ChZHZXRTeXN0ZW1TdGF0c1Jlc3BvbnNlEh8KC3RvdGFsX3VzZXJzGAEgASgFUgp0b3RhbFVzZX'
    'JzEiEKDGFjdGl2ZV91c2VycxgCIAEoBVILYWN0aXZlVXNlcnMSIwoNdG90YWxfZGV2aWNlcxgD'
    'IAEoBVIMdG90YWxEZXZpY2VzEi0KEnRvdGFsX21lYXN1cmVtZW50cxgEIAEoBVIRdG90YWxNZW'
    'FzdXJlbWVudHMSWAoNdXNlcnNfYnlfdGllchgFIAMoCzI0Lm1hbnBhc2lrLnYxLkdldFN5c3Rl'
    'bVN0YXRzUmVzcG9uc2UuVXNlcnNCeVRpZXJFbnRyeVILdXNlcnNCeVRpZXISbQoUbWVhc3VyZW'
    '1lbnRzX2J5X3R5cGUYBiADKAsyOy5tYW5wYXNpay52MS5HZXRTeXN0ZW1TdGF0c1Jlc3BvbnNl'
    'Lk1lYXN1cmVtZW50c0J5VHlwZUVudHJ5UhJtZWFzdXJlbWVudHNCeVR5cGUSPQoMZ2VuZXJhdG'
    'VkX2F0GAcgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFILZ2VuZXJhdGVkQXQaPgoQ'
    'VXNlcnNCeVRpZXJFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoBVIFdmFsdW'
    'U6AjgBGkUKF01lYXN1cmVtZW50c0J5VHlwZUVudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZh'
    'bHVlGAIgASgFUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use getAuditLogRequestDescriptor instead')
const GetAuditLogRequest$json = {
  '1': 'GetAuditLogRequest',
  '2': [
    {'1': 'admin_id', '3': 1, '4': 1, '5': 9, '10': 'adminId'},
    {
      '1': 'action_filter',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.AuditAction',
      '10': 'actionFilter'
    },
    {
      '1': 'start_time',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'startTime'
    },
    {
      '1': 'end_time',
      '3': 4,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'endTime'
    },
    {'1': 'limit', '3': 5, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 6, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `GetAuditLogRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAuditLogRequestDescriptor = $convert.base64Decode(
    'ChJHZXRBdWRpdExvZ1JlcXVlc3QSGQoIYWRtaW5faWQYASABKAlSB2FkbWluSWQSPQoNYWN0aW'
    '9uX2ZpbHRlchgCIAEoDjIYLm1hbnBhc2lrLnYxLkF1ZGl0QWN0aW9uUgxhY3Rpb25GaWx0ZXIS'
    'OQoKc3RhcnRfdGltZRgDIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSCXN0YXJ0VG'
    'ltZRI1CghlbmRfdGltZRgEIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSB2VuZFRp'
    'bWUSFAoFbGltaXQYBSABKAVSBWxpbWl0EhYKBm9mZnNldBgGIAEoBVIGb2Zmc2V0');

@$core.Deprecated('Use getAuditLogResponseDescriptor instead')
const GetAuditLogResponse$json = {
  '1': 'GetAuditLogResponse',
  '2': [
    {
      '1': 'entries',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.AuditLogEntry',
      '10': 'entries'
    },
    {'1': 'total_count', '3': 2, '4': 1, '5': 5, '10': 'totalCount'},
  ],
};

/// Descriptor for `GetAuditLogResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAuditLogResponseDescriptor = $convert.base64Decode(
    'ChNHZXRBdWRpdExvZ1Jlc3BvbnNlEjQKB2VudHJpZXMYASADKAsyGi5tYW5wYXNpay52MS5BdW'
    'RpdExvZ0VudHJ5UgdlbnRyaWVzEh8KC3RvdGFsX2NvdW50GAIgASgFUgp0b3RhbENvdW50');

@$core.Deprecated('Use auditLogEntryDescriptor instead')
const AuditLogEntry$json = {
  '1': 'AuditLogEntry',
  '2': [
    {'1': 'entry_id', '3': 1, '4': 1, '5': 9, '10': 'entryId'},
    {'1': 'admin_id', '3': 2, '4': 1, '5': 9, '10': 'adminId'},
    {'1': 'admin_email', '3': 3, '4': 1, '5': 9, '10': 'adminEmail'},
    {
      '1': 'action',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.AuditAction',
      '10': 'action'
    },
    {'1': 'resource_type', '3': 5, '4': 1, '5': 9, '10': 'resourceType'},
    {'1': 'resource_id', '3': 6, '4': 1, '5': 9, '10': 'resourceId'},
    {'1': 'details', '3': 7, '4': 1, '5': 9, '10': 'details'},
    {'1': 'ip_address', '3': 8, '4': 1, '5': 9, '10': 'ipAddress'},
    {
      '1': 'timestamp',
      '3': 9,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'timestamp'
    },
  ],
};

/// Descriptor for `AuditLogEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List auditLogEntryDescriptor = $convert.base64Decode(
    'Cg1BdWRpdExvZ0VudHJ5EhkKCGVudHJ5X2lkGAEgASgJUgdlbnRyeUlkEhkKCGFkbWluX2lkGA'
    'IgASgJUgdhZG1pbklkEh8KC2FkbWluX2VtYWlsGAMgASgJUgphZG1pbkVtYWlsEjAKBmFjdGlv'
    'bhgEIAEoDjIYLm1hbnBhc2lrLnYxLkF1ZGl0QWN0aW9uUgZhY3Rpb24SIwoNcmVzb3VyY2VfdH'
    'lwZRgFIAEoCVIMcmVzb3VyY2VUeXBlEh8KC3Jlc291cmNlX2lkGAYgASgJUgpyZXNvdXJjZUlk'
    'EhgKB2RldGFpbHMYByABKAlSB2RldGFpbHMSHQoKaXBfYWRkcmVzcxgIIAEoCVIJaXBBZGRyZX'
    'NzEjgKCXRpbWVzdGFtcBgJIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSCXRpbWVz'
    'dGFtcA==');

@$core.Deprecated('Use setSystemConfigRequestDescriptor instead')
const SetSystemConfigRequest$json = {
  '1': 'SetSystemConfigRequest',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
    {'1': 'description', '3': 3, '4': 1, '5': 9, '10': 'description'},
  ],
};

/// Descriptor for `SetSystemConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setSystemConfigRequestDescriptor =
    $convert.base64Decode(
        'ChZTZXRTeXN0ZW1Db25maWdSZXF1ZXN0EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgAS'
        'gJUgV2YWx1ZRIgCgtkZXNjcmlwdGlvbhgDIAEoCVILZGVzY3JpcHRpb24=');

@$core.Deprecated('Use getSystemConfigRequestDescriptor instead')
const GetSystemConfigRequest$json = {
  '1': 'GetSystemConfigRequest',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
  ],
};

/// Descriptor for `GetSystemConfigRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSystemConfigRequestDescriptor = $convert
    .base64Decode('ChZHZXRTeXN0ZW1Db25maWdSZXF1ZXN0EhAKA2tleRgBIAEoCVIDa2V5');

@$core.Deprecated('Use systemConfigDescriptor instead')
const SystemConfig$json = {
  '1': 'SystemConfig',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
    {'1': 'description', '3': 3, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'updated_at',
      '3': 4,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'updatedAt'
    },
    {'1': 'updated_by', '3': 5, '4': 1, '5': 9, '10': 'updatedBy'},
  ],
};

/// Descriptor for `SystemConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List systemConfigDescriptor = $convert.base64Decode(
    'CgxTeXN0ZW1Db25maWcSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlEi'
    'AKC2Rlc2NyaXB0aW9uGAMgASgJUgtkZXNjcmlwdGlvbhI5Cgp1cGRhdGVkX2F0GAQgASgLMhou'
    'Z29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFIJdXBkYXRlZEF0Eh0KCnVwZGF0ZWRfYnkYBSABKA'
    'lSCXVwZGF0ZWRCeQ==');

@$core.Deprecated('Use createFamilyGroupRequestDescriptor instead')
const CreateFamilyGroupRequest$json = {
  '1': 'CreateFamilyGroupRequest',
  '2': [
    {'1': 'owner_user_id', '3': 1, '4': 1, '5': 9, '10': 'ownerUserId'},
    {'1': 'name', '3': 2, '4': 1, '5': 9, '10': 'name'},
    {'1': 'description', '3': 3, '4': 1, '5': 9, '10': 'description'},
  ],
};

/// Descriptor for `CreateFamilyGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createFamilyGroupRequestDescriptor = $convert.base64Decode(
    'ChhDcmVhdGVGYW1pbHlHcm91cFJlcXVlc3QSIgoNb3duZXJfdXNlcl9pZBgBIAEoCVILb3duZX'
    'JVc2VySWQSEgoEbmFtZRgCIAEoCVIEbmFtZRIgCgtkZXNjcmlwdGlvbhgDIAEoCVILZGVzY3Jp'
    'cHRpb24=');

@$core.Deprecated('Use getFamilyGroupRequestDescriptor instead')
const GetFamilyGroupRequest$json = {
  '1': 'GetFamilyGroupRequest',
  '2': [
    {'1': 'group_id', '3': 1, '4': 1, '5': 9, '10': 'groupId'},
  ],
};

/// Descriptor for `GetFamilyGroupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFamilyGroupRequestDescriptor =
    $convert.base64Decode(
        'ChVHZXRGYW1pbHlHcm91cFJlcXVlc3QSGQoIZ3JvdXBfaWQYASABKAlSB2dyb3VwSWQ=');

@$core.Deprecated('Use familyGroupDescriptor instead')
const FamilyGroup$json = {
  '1': 'FamilyGroup',
  '2': [
    {'1': 'group_id', '3': 1, '4': 1, '5': 9, '10': 'groupId'},
    {'1': 'name', '3': 2, '4': 1, '5': 9, '10': 'name'},
    {'1': 'description', '3': 3, '4': 1, '5': 9, '10': 'description'},
    {'1': 'owner_user_id', '3': 4, '4': 1, '5': 9, '10': 'ownerUserId'},
    {'1': 'member_count', '3': 5, '4': 1, '5': 5, '10': 'memberCount'},
    {'1': 'max_members', '3': 6, '4': 1, '5': 5, '10': 'maxMembers'},
    {
      '1': 'members',
      '3': 7,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.FamilyMember',
      '10': 'members'
    },
    {
      '1': 'created_at',
      '3': 8,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
  ],
};

/// Descriptor for `FamilyGroup`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List familyGroupDescriptor = $convert.base64Decode(
    'CgtGYW1pbHlHcm91cBIZCghncm91cF9pZBgBIAEoCVIHZ3JvdXBJZBISCgRuYW1lGAIgASgJUg'
    'RuYW1lEiAKC2Rlc2NyaXB0aW9uGAMgASgJUgtkZXNjcmlwdGlvbhIiCg1vd25lcl91c2VyX2lk'
    'GAQgASgJUgtvd25lclVzZXJJZBIhCgxtZW1iZXJfY291bnQYBSABKAVSC21lbWJlckNvdW50Eh'
    '8KC21heF9tZW1iZXJzGAYgASgFUgptYXhNZW1iZXJzEjMKB21lbWJlcnMYByADKAsyGS5tYW5w'
    'YXNpay52MS5GYW1pbHlNZW1iZXJSB21lbWJlcnMSOQoKY3JlYXRlZF9hdBgIIAEoCzIaLmdvb2'
    'dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSCWNyZWF0ZWRBdA==');

@$core.Deprecated('Use familyMemberDescriptor instead')
const FamilyMember$json = {
  '1': 'FamilyMember',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'display_name', '3': 2, '4': 1, '5': 9, '10': 'displayName'},
    {
      '1': 'role',
      '3': 3,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.FamilyRole',
      '10': 'role'
    },
    {'1': 'relationship', '3': 4, '4': 1, '5': 9, '10': 'relationship'},
    {
      '1': 'can_view_health_data',
      '3': 5,
      '4': 1,
      '5': 8,
      '10': 'canViewHealthData'
    },
    {
      '1': 'can_manage_devices',
      '3': 6,
      '4': 1,
      '5': 8,
      '10': 'canManageDevices'
    },
    {
      '1': 'joined_at',
      '3': 7,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'joinedAt'
    },
  ],
};

/// Descriptor for `FamilyMember`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List familyMemberDescriptor = $convert.base64Decode(
    'CgxGYW1pbHlNZW1iZXISFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklkEiEKDGRpc3BsYXlfbmFtZR'
    'gCIAEoCVILZGlzcGxheU5hbWUSKwoEcm9sZRgDIAEoDjIXLm1hbnBhc2lrLnYxLkZhbWlseVJv'
    'bGVSBHJvbGUSIgoMcmVsYXRpb25zaGlwGAQgASgJUgxyZWxhdGlvbnNoaXASLwoUY2FuX3ZpZX'
    'dfaGVhbHRoX2RhdGEYBSABKAhSEWNhblZpZXdIZWFsdGhEYXRhEiwKEmNhbl9tYW5hZ2VfZGV2'
    'aWNlcxgGIAEoCFIQY2FuTWFuYWdlRGV2aWNlcxI3Cglqb2luZWRfYXQYByABKAsyGi5nb29nbG'
    'UucHJvdG9idWYuVGltZXN0YW1wUghqb2luZWRBdA==');

@$core.Deprecated('Use inviteMemberRequestDescriptor instead')
const InviteMemberRequest$json = {
  '1': 'InviteMemberRequest',
  '2': [
    {'1': 'group_id', '3': 1, '4': 1, '5': 9, '10': 'groupId'},
    {'1': 'inviter_user_id', '3': 2, '4': 1, '5': 9, '10': 'inviterUserId'},
    {'1': 'invitee_email', '3': 3, '4': 1, '5': 9, '10': 'inviteeEmail'},
    {
      '1': 'role',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.FamilyRole',
      '10': 'role'
    },
    {'1': 'relationship', '3': 5, '4': 1, '5': 9, '10': 'relationship'},
  ],
};

/// Descriptor for `InviteMemberRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inviteMemberRequestDescriptor = $convert.base64Decode(
    'ChNJbnZpdGVNZW1iZXJSZXF1ZXN0EhkKCGdyb3VwX2lkGAEgASgJUgdncm91cElkEiYKD2ludm'
    'l0ZXJfdXNlcl9pZBgCIAEoCVINaW52aXRlclVzZXJJZBIjCg1pbnZpdGVlX2VtYWlsGAMgASgJ'
    'UgxpbnZpdGVlRW1haWwSKwoEcm9sZRgEIAEoDjIXLm1hbnBhc2lrLnYxLkZhbWlseVJvbGVSBH'
    'JvbGUSIgoMcmVsYXRpb25zaGlwGAUgASgJUgxyZWxhdGlvbnNoaXA=');

@$core.Deprecated('Use familyInvitationDescriptor instead')
const FamilyInvitation$json = {
  '1': 'FamilyInvitation',
  '2': [
    {'1': 'invitation_id', '3': 1, '4': 1, '5': 9, '10': 'invitationId'},
    {'1': 'group_id', '3': 2, '4': 1, '5': 9, '10': 'groupId'},
    {'1': 'group_name', '3': 3, '4': 1, '5': 9, '10': 'groupName'},
    {'1': 'inviter_user_id', '3': 4, '4': 1, '5': 9, '10': 'inviterUserId'},
    {'1': 'inviter_name', '3': 5, '4': 1, '5': 9, '10': 'inviterName'},
    {'1': 'invitee_email', '3': 6, '4': 1, '5': 9, '10': 'inviteeEmail'},
    {
      '1': 'role',
      '3': 7,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.FamilyRole',
      '10': 'role'
    },
    {
      '1': 'status',
      '3': 8,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.InvitationStatus',
      '10': 'status'
    },
    {
      '1': 'created_at',
      '3': 9,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
    {
      '1': 'expires_at',
      '3': 10,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'expiresAt'
    },
  ],
};

/// Descriptor for `FamilyInvitation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List familyInvitationDescriptor = $convert.base64Decode(
    'ChBGYW1pbHlJbnZpdGF0aW9uEiMKDWludml0YXRpb25faWQYASABKAlSDGludml0YXRpb25JZB'
    'IZCghncm91cF9pZBgCIAEoCVIHZ3JvdXBJZBIdCgpncm91cF9uYW1lGAMgASgJUglncm91cE5h'
    'bWUSJgoPaW52aXRlcl91c2VyX2lkGAQgASgJUg1pbnZpdGVyVXNlcklkEiEKDGludml0ZXJfbm'
    'FtZRgFIAEoCVILaW52aXRlck5hbWUSIwoNaW52aXRlZV9lbWFpbBgGIAEoCVIMaW52aXRlZUVt'
    'YWlsEisKBHJvbGUYByABKA4yFy5tYW5wYXNpay52MS5GYW1pbHlSb2xlUgRyb2xlEjUKBnN0YX'
    'R1cxgIIAEoDjIdLm1hbnBhc2lrLnYxLkludml0YXRpb25TdGF0dXNSBnN0YXR1cxI5CgpjcmVh'
    'dGVkX2F0GAkgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFIJY3JlYXRlZEF0EjkKCm'
    'V4cGlyZXNfYXQYCiABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUglleHBpcmVzQXQ=');

@$core.Deprecated('Use respondToInvitationRequestDescriptor instead')
const RespondToInvitationRequest$json = {
  '1': 'RespondToInvitationRequest',
  '2': [
    {'1': 'invitation_id', '3': 1, '4': 1, '5': 9, '10': 'invitationId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'accept', '3': 3, '4': 1, '5': 8, '10': 'accept'},
  ],
};

/// Descriptor for `RespondToInvitationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List respondToInvitationRequestDescriptor =
    $convert.base64Decode(
        'ChpSZXNwb25kVG9JbnZpdGF0aW9uUmVxdWVzdBIjCg1pbnZpdGF0aW9uX2lkGAEgASgJUgxpbn'
        'ZpdGF0aW9uSWQSFwoHdXNlcl9pZBgCIAEoCVIGdXNlcklkEhYKBmFjY2VwdBgDIAEoCFIGYWNj'
        'ZXB0');

@$core.Deprecated('Use respondToInvitationResponseDescriptor instead')
const RespondToInvitationResponse$json = {
  '1': 'RespondToInvitationResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
    {'1': 'message', '3': 2, '4': 1, '5': 9, '10': 'message'},
    {
      '1': 'group',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.manpasik.v1.FamilyGroup',
      '10': 'group'
    },
  ],
};

/// Descriptor for `RespondToInvitationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List respondToInvitationResponseDescriptor =
    $convert.base64Decode(
        'ChtSZXNwb25kVG9JbnZpdGF0aW9uUmVzcG9uc2USGAoHc3VjY2VzcxgBIAEoCFIHc3VjY2Vzcx'
        'IYCgdtZXNzYWdlGAIgASgJUgdtZXNzYWdlEi4KBWdyb3VwGAMgASgLMhgubWFucGFzaWsudjEu'
        'RmFtaWx5R3JvdXBSBWdyb3Vw');

@$core.Deprecated('Use removeMemberRequestDescriptor instead')
const RemoveMemberRequest$json = {
  '1': 'RemoveMemberRequest',
  '2': [
    {'1': 'group_id', '3': 1, '4': 1, '5': 9, '10': 'groupId'},
    {'1': 'requester_user_id', '3': 2, '4': 1, '5': 9, '10': 'requesterUserId'},
    {'1': 'target_user_id', '3': 3, '4': 1, '5': 9, '10': 'targetUserId'},
  ],
};

/// Descriptor for `RemoveMemberRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removeMemberRequestDescriptor = $convert.base64Decode(
    'ChNSZW1vdmVNZW1iZXJSZXF1ZXN0EhkKCGdyb3VwX2lkGAEgASgJUgdncm91cElkEioKEXJlcX'
    'Vlc3Rlcl91c2VyX2lkGAIgASgJUg9yZXF1ZXN0ZXJVc2VySWQSJAoOdGFyZ2V0X3VzZXJfaWQY'
    'AyABKAlSDHRhcmdldFVzZXJJZA==');

@$core.Deprecated('Use removeMemberResponseDescriptor instead')
const RemoveMemberResponse$json = {
  '1': 'RemoveMemberResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
    {'1': 'message', '3': 2, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `RemoveMemberResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removeMemberResponseDescriptor = $convert.base64Decode(
    'ChRSZW1vdmVNZW1iZXJSZXNwb25zZRIYCgdzdWNjZXNzGAEgASgIUgdzdWNjZXNzEhgKB21lc3'
    'NhZ2UYAiABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use updateMemberRoleRequestDescriptor instead')
const UpdateMemberRoleRequest$json = {
  '1': 'UpdateMemberRoleRequest',
  '2': [
    {'1': 'group_id', '3': 1, '4': 1, '5': 9, '10': 'groupId'},
    {'1': 'requester_user_id', '3': 2, '4': 1, '5': 9, '10': 'requesterUserId'},
    {'1': 'target_user_id', '3': 3, '4': 1, '5': 9, '10': 'targetUserId'},
    {
      '1': 'new_role',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.FamilyRole',
      '10': 'newRole'
    },
  ],
};

/// Descriptor for `UpdateMemberRoleRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateMemberRoleRequestDescriptor = $convert.base64Decode(
    'ChdVcGRhdGVNZW1iZXJSb2xlUmVxdWVzdBIZCghncm91cF9pZBgBIAEoCVIHZ3JvdXBJZBIqCh'
    'FyZXF1ZXN0ZXJfdXNlcl9pZBgCIAEoCVIPcmVxdWVzdGVyVXNlcklkEiQKDnRhcmdldF91c2Vy'
    'X2lkGAMgASgJUgx0YXJnZXRVc2VySWQSMgoIbmV3X3JvbGUYBCABKA4yFy5tYW5wYXNpay52MS'
    '5GYW1pbHlSb2xlUgduZXdSb2xl');

@$core.Deprecated('Use listFamilyMembersRequestDescriptor instead')
const ListFamilyMembersRequest$json = {
  '1': 'ListFamilyMembersRequest',
  '2': [
    {'1': 'group_id', '3': 1, '4': 1, '5': 9, '10': 'groupId'},
  ],
};

/// Descriptor for `ListFamilyMembersRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listFamilyMembersRequestDescriptor =
    $convert.base64Decode(
        'ChhMaXN0RmFtaWx5TWVtYmVyc1JlcXVlc3QSGQoIZ3JvdXBfaWQYASABKAlSB2dyb3VwSWQ=');

@$core.Deprecated('Use listFamilyMembersResponseDescriptor instead')
const ListFamilyMembersResponse$json = {
  '1': 'ListFamilyMembersResponse',
  '2': [
    {
      '1': 'members',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.FamilyMember',
      '10': 'members'
    },
  ],
};

/// Descriptor for `ListFamilyMembersResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listFamilyMembersResponseDescriptor =
    $convert.base64Decode(
        'ChlMaXN0RmFtaWx5TWVtYmVyc1Jlc3BvbnNlEjMKB21lbWJlcnMYASADKAsyGS5tYW5wYXNpay'
        '52MS5GYW1pbHlNZW1iZXJSB21lbWJlcnM=');

@$core.Deprecated('Use setSharingPreferencesRequestDescriptor instead')
const SetSharingPreferencesRequest$json = {
  '1': 'SetSharingPreferencesRequest',
  '2': [
    {'1': 'group_id', '3': 1, '4': 1, '5': 9, '10': 'groupId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'share_measurements',
      '3': 3,
      '4': 1,
      '5': 8,
      '10': 'shareMeasurements'
    },
    {
      '1': 'share_health_records',
      '3': 4,
      '4': 1,
      '5': 8,
      '10': 'shareHealthRecords'
    },
    {
      '1': 'share_prescriptions',
      '3': 5,
      '4': 1,
      '5': 8,
      '10': 'sharePrescriptions'
    },
    {'1': 'share_coaching', '3': 6, '4': 1, '5': 8, '10': 'shareCoaching'},
    {
      '1': 'shared_with_user_ids',
      '3': 7,
      '4': 3,
      '5': 9,
      '10': 'sharedWithUserIds'
    },
    {
      '1': 'share_health_score',
      '3': 8,
      '4': 1,
      '5': 8,
      '10': 'shareHealthScore'
    },
    {'1': 'share_goals', '3': 9, '4': 1, '5': 8, '10': 'shareGoals'},
    {'1': 'share_alerts', '3': 10, '4': 1, '5': 8, '10': 'shareAlerts'},
    {
      '1': 'measurement_days_limit',
      '3': 11,
      '4': 1,
      '5': 5,
      '10': 'measurementDaysLimit'
    },
    {
      '1': 'allowed_biomarkers',
      '3': 12,
      '4': 3,
      '5': 9,
      '10': 'allowedBiomarkers'
    },
    {'1': 'require_approval', '3': 13, '4': 1, '5': 8, '10': 'requireApproval'},
  ],
};

/// Descriptor for `SetSharingPreferencesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setSharingPreferencesRequestDescriptor = $convert.base64Decode(
    'ChxTZXRTaGFyaW5nUHJlZmVyZW5jZXNSZXF1ZXN0EhkKCGdyb3VwX2lkGAEgASgJUgdncm91cE'
    'lkEhcKB3VzZXJfaWQYAiABKAlSBnVzZXJJZBItChJzaGFyZV9tZWFzdXJlbWVudHMYAyABKAhS'
    'EXNoYXJlTWVhc3VyZW1lbnRzEjAKFHNoYXJlX2hlYWx0aF9yZWNvcmRzGAQgASgIUhJzaGFyZU'
    'hlYWx0aFJlY29yZHMSLwoTc2hhcmVfcHJlc2NyaXB0aW9ucxgFIAEoCFISc2hhcmVQcmVzY3Jp'
    'cHRpb25zEiUKDnNoYXJlX2NvYWNoaW5nGAYgASgIUg1zaGFyZUNvYWNoaW5nEi8KFHNoYXJlZF'
    '93aXRoX3VzZXJfaWRzGAcgAygJUhFzaGFyZWRXaXRoVXNlcklkcxIsChJzaGFyZV9oZWFsdGhf'
    'c2NvcmUYCCABKAhSEHNoYXJlSGVhbHRoU2NvcmUSHwoLc2hhcmVfZ29hbHMYCSABKAhSCnNoYX'
    'JlR29hbHMSIQoMc2hhcmVfYWxlcnRzGAogASgIUgtzaGFyZUFsZXJ0cxI0ChZtZWFzdXJlbWVu'
    'dF9kYXlzX2xpbWl0GAsgASgFUhRtZWFzdXJlbWVudERheXNMaW1pdBItChJhbGxvd2VkX2Jpb2'
    '1hcmtlcnMYDCADKAlSEWFsbG93ZWRCaW9tYXJrZXJzEikKEHJlcXVpcmVfYXBwcm92YWwYDSAB'
    'KAhSD3JlcXVpcmVBcHByb3ZhbA==');

@$core.Deprecated('Use sharingPreferencesDescriptor instead')
const SharingPreferences$json = {
  '1': 'SharingPreferences',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'group_id', '3': 2, '4': 1, '5': 9, '10': 'groupId'},
    {
      '1': 'share_measurements',
      '3': 3,
      '4': 1,
      '5': 8,
      '10': 'shareMeasurements'
    },
    {
      '1': 'share_health_records',
      '3': 4,
      '4': 1,
      '5': 8,
      '10': 'shareHealthRecords'
    },
    {
      '1': 'share_prescriptions',
      '3': 5,
      '4': 1,
      '5': 8,
      '10': 'sharePrescriptions'
    },
    {'1': 'share_coaching', '3': 6, '4': 1, '5': 8, '10': 'shareCoaching'},
    {
      '1': 'shared_with_user_ids',
      '3': 7,
      '4': 3,
      '5': 9,
      '10': 'sharedWithUserIds'
    },
    {
      '1': 'updated_at',
      '3': 8,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'updatedAt'
    },
    {
      '1': 'share_health_score',
      '3': 9,
      '4': 1,
      '5': 8,
      '10': 'shareHealthScore'
    },
    {'1': 'share_goals', '3': 10, '4': 1, '5': 8, '10': 'shareGoals'},
    {'1': 'share_alerts', '3': 11, '4': 1, '5': 8, '10': 'shareAlerts'},
    {
      '1': 'measurement_days_limit',
      '3': 12,
      '4': 1,
      '5': 5,
      '10': 'measurementDaysLimit'
    },
    {
      '1': 'allowed_biomarkers',
      '3': 13,
      '4': 3,
      '5': 9,
      '10': 'allowedBiomarkers'
    },
    {'1': 'require_approval', '3': 14, '4': 1, '5': 8, '10': 'requireApproval'},
  ],
};

/// Descriptor for `SharingPreferences`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sharingPreferencesDescriptor = $convert.base64Decode(
    'ChJTaGFyaW5nUHJlZmVyZW5jZXMSFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklkEhkKCGdyb3VwX2'
    'lkGAIgASgJUgdncm91cElkEi0KEnNoYXJlX21lYXN1cmVtZW50cxgDIAEoCFIRc2hhcmVNZWFz'
    'dXJlbWVudHMSMAoUc2hhcmVfaGVhbHRoX3JlY29yZHMYBCABKAhSEnNoYXJlSGVhbHRoUmVjb3'
    'JkcxIvChNzaGFyZV9wcmVzY3JpcHRpb25zGAUgASgIUhJzaGFyZVByZXNjcmlwdGlvbnMSJQoO'
    'c2hhcmVfY29hY2hpbmcYBiABKAhSDXNoYXJlQ29hY2hpbmcSLwoUc2hhcmVkX3dpdGhfdXNlcl'
    '9pZHMYByADKAlSEXNoYXJlZFdpdGhVc2VySWRzEjkKCnVwZGF0ZWRfYXQYCCABKAsyGi5nb29n'
    'bGUucHJvdG9idWYuVGltZXN0YW1wUgl1cGRhdGVkQXQSLAoSc2hhcmVfaGVhbHRoX3Njb3JlGA'
    'kgASgIUhBzaGFyZUhlYWx0aFNjb3JlEh8KC3NoYXJlX2dvYWxzGAogASgIUgpzaGFyZUdvYWxz'
    'EiEKDHNoYXJlX2FsZXJ0cxgLIAEoCFILc2hhcmVBbGVydHMSNAoWbWVhc3VyZW1lbnRfZGF5c1'
    '9saW1pdBgMIAEoBVIUbWVhc3VyZW1lbnREYXlzTGltaXQSLQoSYWxsb3dlZF9iaW9tYXJrZXJz'
    'GA0gAygJUhFhbGxvd2VkQmlvbWFya2VycxIpChByZXF1aXJlX2FwcHJvdmFsGA4gASgIUg9yZX'
    'F1aXJlQXBwcm92YWw=');

@$core.Deprecated('Use getSharedHealthDataRequestDescriptor instead')
const GetSharedHealthDataRequest$json = {
  '1': 'GetSharedHealthDataRequest',
  '2': [
    {'1': 'group_id', '3': 1, '4': 1, '5': 9, '10': 'groupId'},
    {'1': 'requester_user_id', '3': 2, '4': 1, '5': 9, '10': 'requesterUserId'},
    {'1': 'target_user_id', '3': 3, '4': 1, '5': 9, '10': 'targetUserId'},
  ],
};

/// Descriptor for `GetSharedHealthDataRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSharedHealthDataRequestDescriptor =
    $convert.base64Decode(
        'ChpHZXRTaGFyZWRIZWFsdGhEYXRhUmVxdWVzdBIZCghncm91cF9pZBgBIAEoCVIHZ3JvdXBJZB'
        'IqChFyZXF1ZXN0ZXJfdXNlcl9pZBgCIAEoCVIPcmVxdWVzdGVyVXNlcklkEiQKDnRhcmdldF91'
        'c2VyX2lkGAMgASgJUgx0YXJnZXRVc2VySWQ=');

@$core.Deprecated('Use getSharedHealthDataResponseDescriptor instead')
const GetSharedHealthDataResponse$json = {
  '1': 'GetSharedHealthDataResponse',
  '2': [
    {'1': 'target_user_id', '3': 1, '4': 1, '5': 9, '10': 'targetUserId'},
    {
      '1': 'target_display_name',
      '3': 2,
      '4': 1,
      '5': 9,
      '10': 'targetDisplayName'
    },
    {
      '1': 'recent_measurements',
      '3': 3,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.MeasurementSummary',
      '10': 'recentMeasurements'
    },
    {'1': 'health_score', '3': 4, '4': 1, '5': 1, '10': 'healthScore'},
    {'1': 'last_active', '3': 5, '4': 1, '5': 9, '10': 'lastActive'},
  ],
};

/// Descriptor for `GetSharedHealthDataResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSharedHealthDataResponseDescriptor = $convert.base64Decode(
    'ChtHZXRTaGFyZWRIZWFsdGhEYXRhUmVzcG9uc2USJAoOdGFyZ2V0X3VzZXJfaWQYASABKAlSDH'
    'RhcmdldFVzZXJJZBIuChN0YXJnZXRfZGlzcGxheV9uYW1lGAIgASgJUhF0YXJnZXREaXNwbGF5'
    'TmFtZRJQChNyZWNlbnRfbWVhc3VyZW1lbnRzGAMgAygLMh8ubWFucGFzaWsudjEuTWVhc3VyZW'
    '1lbnRTdW1tYXJ5UhJyZWNlbnRNZWFzdXJlbWVudHMSIQoMaGVhbHRoX3Njb3JlGAQgASgBUgto'
    'ZWFsdGhTY29yZRIfCgtsYXN0X2FjdGl2ZRgFIAEoCVIKbGFzdEFjdGl2ZQ==');

@$core.Deprecated('Use createHealthRecordRequestDescriptor instead')
const CreateHealthRecordRequest$json = {
  '1': 'CreateHealthRecordRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'record_type',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.HealthRecordType',
      '10': 'recordType'
    },
    {'1': 'title', '3': 3, '4': 1, '5': 9, '10': 'title'},
    {'1': 'description', '3': 4, '4': 1, '5': 9, '10': 'description'},
    {'1': 'provider', '3': 5, '4': 1, '5': 9, '10': 'provider'},
    {
      '1': 'metadata',
      '3': 6,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.CreateHealthRecordRequest.MetadataEntry',
      '10': 'metadata'
    },
    {'1': 'measurement_id', '3': 7, '4': 1, '5': 9, '10': 'measurementId'},
  ],
  '3': [CreateHealthRecordRequest_MetadataEntry$json],
};

@$core.Deprecated('Use createHealthRecordRequestDescriptor instead')
const CreateHealthRecordRequest_MetadataEntry$json = {
  '1': 'MetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateHealthRecordRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createHealthRecordRequestDescriptor = $convert.base64Decode(
    'ChlDcmVhdGVIZWFsdGhSZWNvcmRSZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZBI+Cg'
    'tyZWNvcmRfdHlwZRgCIAEoDjIdLm1hbnBhc2lrLnYxLkhlYWx0aFJlY29yZFR5cGVSCnJlY29y'
    'ZFR5cGUSFAoFdGl0bGUYAyABKAlSBXRpdGxlEiAKC2Rlc2NyaXB0aW9uGAQgASgJUgtkZXNjcm'
    'lwdGlvbhIaCghwcm92aWRlchgFIAEoCVIIcHJvdmlkZXISUAoIbWV0YWRhdGEYBiADKAsyNC5t'
    'YW5wYXNpay52MS5DcmVhdGVIZWFsdGhSZWNvcmRSZXF1ZXN0Lk1ldGFkYXRhRW50cnlSCG1ldG'
    'FkYXRhEiUKDm1lYXN1cmVtZW50X2lkGAcgASgJUg1tZWFzdXJlbWVudElkGjsKDU1ldGFkYXRh'
    'RW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use getHealthRecordRequestDescriptor instead')
const GetHealthRecordRequest$json = {
  '1': 'GetHealthRecordRequest',
  '2': [
    {'1': 'record_id', '3': 1, '4': 1, '5': 9, '10': 'recordId'},
  ],
};

/// Descriptor for `GetHealthRecordRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getHealthRecordRequestDescriptor =
    $convert.base64Decode(
        'ChZHZXRIZWFsdGhSZWNvcmRSZXF1ZXN0EhsKCXJlY29yZF9pZBgBIAEoCVIIcmVjb3JkSWQ=');

@$core.Deprecated('Use listHealthRecordsRequestDescriptor instead')
const ListHealthRecordsRequest$json = {
  '1': 'ListHealthRecordsRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'type_filter',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.HealthRecordType',
      '10': 'typeFilter'
    },
    {'1': 'limit', '3': 3, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 4, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `ListHealthRecordsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listHealthRecordsRequestDescriptor = $convert.base64Decode(
    'ChhMaXN0SGVhbHRoUmVjb3Jkc1JlcXVlc3QSFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklkEj4KC3'
    'R5cGVfZmlsdGVyGAIgASgOMh0ubWFucGFzaWsudjEuSGVhbHRoUmVjb3JkVHlwZVIKdHlwZUZp'
    'bHRlchIUCgVsaW1pdBgDIAEoBVIFbGltaXQSFgoGb2Zmc2V0GAQgASgFUgZvZmZzZXQ=');

@$core.Deprecated('Use listHealthRecordsResponseDescriptor instead')
const ListHealthRecordsResponse$json = {
  '1': 'ListHealthRecordsResponse',
  '2': [
    {
      '1': 'records',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.HealthRecord',
      '10': 'records'
    },
    {'1': 'total_count', '3': 2, '4': 1, '5': 5, '10': 'totalCount'},
  ],
};

/// Descriptor for `ListHealthRecordsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listHealthRecordsResponseDescriptor = $convert.base64Decode(
    'ChlMaXN0SGVhbHRoUmVjb3Jkc1Jlc3BvbnNlEjMKB3JlY29yZHMYASADKAsyGS5tYW5wYXNpay'
    '52MS5IZWFsdGhSZWNvcmRSB3JlY29yZHMSHwoLdG90YWxfY291bnQYAiABKAVSCnRvdGFsQ291'
    'bnQ=');

@$core.Deprecated('Use updateHealthRecordRequestDescriptor instead')
const UpdateHealthRecordRequest$json = {
  '1': 'UpdateHealthRecordRequest',
  '2': [
    {'1': 'record_id', '3': 1, '4': 1, '5': 9, '10': 'recordId'},
    {'1': 'title', '3': 2, '4': 1, '5': 9, '10': 'title'},
    {'1': 'description', '3': 3, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'metadata',
      '3': 4,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.UpdateHealthRecordRequest.MetadataEntry',
      '10': 'metadata'
    },
  ],
  '3': [UpdateHealthRecordRequest_MetadataEntry$json],
};

@$core.Deprecated('Use updateHealthRecordRequestDescriptor instead')
const UpdateHealthRecordRequest_MetadataEntry$json = {
  '1': 'MetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `UpdateHealthRecordRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateHealthRecordRequestDescriptor = $convert.base64Decode(
    'ChlVcGRhdGVIZWFsdGhSZWNvcmRSZXF1ZXN0EhsKCXJlY29yZF9pZBgBIAEoCVIIcmVjb3JkSW'
    'QSFAoFdGl0bGUYAiABKAlSBXRpdGxlEiAKC2Rlc2NyaXB0aW9uGAMgASgJUgtkZXNjcmlwdGlv'
    'bhJQCghtZXRhZGF0YRgEIAMoCzI0Lm1hbnBhc2lrLnYxLlVwZGF0ZUhlYWx0aFJlY29yZFJlcX'
    'Vlc3QuTWV0YWRhdGFFbnRyeVIIbWV0YWRhdGEaOwoNTWV0YWRhdGFFbnRyeRIQCgNrZXkYASAB'
    'KAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use deleteHealthRecordRequestDescriptor instead')
const DeleteHealthRecordRequest$json = {
  '1': 'DeleteHealthRecordRequest',
  '2': [
    {'1': 'record_id', '3': 1, '4': 1, '5': 9, '10': 'recordId'},
  ],
};

/// Descriptor for `DeleteHealthRecordRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteHealthRecordRequestDescriptor =
    $convert.base64Decode(
        'ChlEZWxldGVIZWFsdGhSZWNvcmRSZXF1ZXN0EhsKCXJlY29yZF9pZBgBIAEoCVIIcmVjb3JkSW'
        'Q=');

@$core.Deprecated('Use deleteHealthRecordResponseDescriptor instead')
const DeleteHealthRecordResponse$json = {
  '1': 'DeleteHealthRecordResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
    {'1': 'message', '3': 2, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DeleteHealthRecordResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteHealthRecordResponseDescriptor =
    $convert.base64Decode(
        'ChpEZWxldGVIZWFsdGhSZWNvcmRSZXNwb25zZRIYCgdzdWNjZXNzGAEgASgIUgdzdWNjZXNzEh'
        'gKB21lc3NhZ2UYAiABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use healthRecordDescriptor instead')
const HealthRecord$json = {
  '1': 'HealthRecord',
  '2': [
    {'1': 'record_id', '3': 1, '4': 1, '5': 9, '10': 'recordId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'record_type',
      '3': 3,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.HealthRecordType',
      '10': 'recordType'
    },
    {'1': 'title', '3': 4, '4': 1, '5': 9, '10': 'title'},
    {'1': 'description', '3': 5, '4': 1, '5': 9, '10': 'description'},
    {'1': 'provider', '3': 6, '4': 1, '5': 9, '10': 'provider'},
    {
      '1': 'metadata',
      '3': 7,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.HealthRecord.MetadataEntry',
      '10': 'metadata'
    },
    {'1': 'measurement_id', '3': 8, '4': 1, '5': 9, '10': 'measurementId'},
    {
      '1': 'created_at',
      '3': 9,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
    {
      '1': 'updated_at',
      '3': 10,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'updatedAt'
    },
  ],
  '3': [HealthRecord_MetadataEntry$json],
};

@$core.Deprecated('Use healthRecordDescriptor instead')
const HealthRecord_MetadataEntry$json = {
  '1': 'MetadataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `HealthRecord`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List healthRecordDescriptor = $convert.base64Decode(
    'CgxIZWFsdGhSZWNvcmQSGwoJcmVjb3JkX2lkGAEgASgJUghyZWNvcmRJZBIXCgd1c2VyX2lkGA'
    'IgASgJUgZ1c2VySWQSPgoLcmVjb3JkX3R5cGUYAyABKA4yHS5tYW5wYXNpay52MS5IZWFsdGhS'
    'ZWNvcmRUeXBlUgpyZWNvcmRUeXBlEhQKBXRpdGxlGAQgASgJUgV0aXRsZRIgCgtkZXNjcmlwdG'
    'lvbhgFIAEoCVILZGVzY3JpcHRpb24SGgoIcHJvdmlkZXIYBiABKAlSCHByb3ZpZGVyEkMKCG1l'
    'dGFkYXRhGAcgAygLMicubWFucGFzaWsudjEuSGVhbHRoUmVjb3JkLk1ldGFkYXRhRW50cnlSCG'
    '1ldGFkYXRhEiUKDm1lYXN1cmVtZW50X2lkGAggASgJUg1tZWFzdXJlbWVudElkEjkKCmNyZWF0'
    'ZWRfYXQYCSABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUgljcmVhdGVkQXQSOQoKdX'
    'BkYXRlZF9hdBgKIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSCXVwZGF0ZWRBdBo7'
    'Cg1NZXRhZGF0YUVudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZT'
    'oCOAE=');

@$core.Deprecated('Use exportToFHIRRequestDescriptor instead')
const ExportToFHIRRequest$json = {
  '1': 'ExportToFHIRRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'record_ids', '3': 2, '4': 3, '5': 9, '10': 'recordIds'},
    {
      '1': 'target_type',
      '3': 3,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.FHIRResourceType',
      '10': 'targetType'
    },
  ],
};

/// Descriptor for `ExportToFHIRRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List exportToFHIRRequestDescriptor = $convert.base64Decode(
    'ChNFeHBvcnRUb0ZISVJSZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZBIdCgpyZWNvcm'
    'RfaWRzGAIgAygJUglyZWNvcmRJZHMSPgoLdGFyZ2V0X3R5cGUYAyABKA4yHS5tYW5wYXNpay52'
    'MS5GSElSUmVzb3VyY2VUeXBlUgp0YXJnZXRUeXBl');

@$core.Deprecated('Use exportToFHIRResponseDescriptor instead')
const ExportToFHIRResponse$json = {
  '1': 'ExportToFHIRResponse',
  '2': [
    {'1': 'fhir_bundle_json', '3': 1, '4': 1, '5': 9, '10': 'fhirBundleJson'},
    {
      '1': 'resource_type',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.FHIRResourceType',
      '10': 'resourceType'
    },
    {'1': 'resource_count', '3': 3, '4': 1, '5': 5, '10': 'resourceCount'},
  ],
};

/// Descriptor for `ExportToFHIRResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List exportToFHIRResponseDescriptor = $convert.base64Decode(
    'ChRFeHBvcnRUb0ZISVJSZXNwb25zZRIoChBmaGlyX2J1bmRsZV9qc29uGAEgASgJUg5maGlyQn'
    'VuZGxlSnNvbhJCCg1yZXNvdXJjZV90eXBlGAIgASgOMh0ubWFucGFzaWsudjEuRkhJUlJlc291'
    'cmNlVHlwZVIMcmVzb3VyY2VUeXBlEiUKDnJlc291cmNlX2NvdW50GAMgASgFUg1yZXNvdXJjZU'
    'NvdW50');

@$core.Deprecated('Use importFromFHIRRequestDescriptor instead')
const ImportFromFHIRRequest$json = {
  '1': 'ImportFromFHIRRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'fhir_bundle_json', '3': 2, '4': 1, '5': 9, '10': 'fhirBundleJson'},
  ],
};

/// Descriptor for `ImportFromFHIRRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importFromFHIRRequestDescriptor = $convert.base64Decode(
    'ChVJbXBvcnRGcm9tRkhJUlJlcXVlc3QSFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklkEigKEGZoaX'
    'JfYnVuZGxlX2pzb24YAiABKAlSDmZoaXJCdW5kbGVKc29u');

@$core.Deprecated('Use importFromFHIRResponseDescriptor instead')
const ImportFromFHIRResponse$json = {
  '1': 'ImportFromFHIRResponse',
  '2': [
    {'1': 'imported_count', '3': 1, '4': 1, '5': 5, '10': 'importedCount'},
    {'1': 'skipped_count', '3': 2, '4': 1, '5': 5, '10': 'skippedCount'},
    {
      '1': 'imported_record_ids',
      '3': 3,
      '4': 3,
      '5': 9,
      '10': 'importedRecordIds'
    },
    {'1': 'errors', '3': 4, '4': 3, '5': 9, '10': 'errors'},
  ],
};

/// Descriptor for `ImportFromFHIRResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importFromFHIRResponseDescriptor = $convert.base64Decode(
    'ChZJbXBvcnRGcm9tRkhJUlJlc3BvbnNlEiUKDmltcG9ydGVkX2NvdW50GAEgASgFUg1pbXBvcn'
    'RlZENvdW50EiMKDXNraXBwZWRfY291bnQYAiABKAVSDHNraXBwZWRDb3VudBIuChNpbXBvcnRl'
    'ZF9yZWNvcmRfaWRzGAMgAygJUhFpbXBvcnRlZFJlY29yZElkcxIWCgZlcnJvcnMYBCADKAlSBm'
    'Vycm9ycw==');

@$core.Deprecated('Use getHealthSummaryRequestDescriptor instead')
const GetHealthSummaryRequest$json = {
  '1': 'GetHealthSummaryRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `GetHealthSummaryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getHealthSummaryRequestDescriptor =
    $convert.base64Decode(
        'ChdHZXRIZWFsdGhTdW1tYXJ5UmVxdWVzdBIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQ=');

@$core.Deprecated('Use getHealthSummaryResponseDescriptor instead')
const GetHealthSummaryResponse$json = {
  '1': 'GetHealthSummaryResponse',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'total_records', '3': 2, '4': 1, '5': 5, '10': 'totalRecords'},
    {
      '1': 'records_by_type',
      '3': 3,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.GetHealthSummaryResponse.RecordsByTypeEntry',
      '10': 'recordsByType'
    },
    {
      '1': 'recent_records',
      '3': 4,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.HealthRecord',
      '10': 'recentRecords'
    },
    {
      '1': 'last_updated',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'lastUpdated'
    },
  ],
  '3': [GetHealthSummaryResponse_RecordsByTypeEntry$json],
};

@$core.Deprecated('Use getHealthSummaryResponseDescriptor instead')
const GetHealthSummaryResponse_RecordsByTypeEntry$json = {
  '1': 'RecordsByTypeEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 5, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetHealthSummaryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getHealthSummaryResponseDescriptor = $convert.base64Decode(
    'ChhHZXRIZWFsdGhTdW1tYXJ5UmVzcG9uc2USFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklkEiMKDX'
    'RvdGFsX3JlY29yZHMYAiABKAVSDHRvdGFsUmVjb3JkcxJgCg9yZWNvcmRzX2J5X3R5cGUYAyAD'
    'KAsyOC5tYW5wYXNpay52MS5HZXRIZWFsdGhTdW1tYXJ5UmVzcG9uc2UuUmVjb3Jkc0J5VHlwZU'
    'VudHJ5Ug1yZWNvcmRzQnlUeXBlEkAKDnJlY2VudF9yZWNvcmRzGAQgAygLMhkubWFucGFzaWsu'
    'djEuSGVhbHRoUmVjb3JkUg1yZWNlbnRSZWNvcmRzEj0KDGxhc3RfdXBkYXRlZBgFIAEoCzIaLm'
    'dvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSC2xhc3RVcGRhdGVkGkAKElJlY29yZHNCeVR5cGVF'
    'bnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoBVIFdmFsdWU6AjgB');

@$core.Deprecated('Use createPrescriptionRequestDescriptor instead')
const CreatePrescriptionRequest$json = {
  '1': 'CreatePrescriptionRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'doctor_id', '3': 2, '4': 1, '5': 9, '10': 'doctorId'},
    {'1': 'doctor_name', '3': 3, '4': 1, '5': 9, '10': 'doctorName'},
    {'1': 'facility_id', '3': 4, '4': 1, '5': 9, '10': 'facilityId'},
    {'1': 'diagnosis', '3': 5, '4': 1, '5': 9, '10': 'diagnosis'},
    {'1': 'notes', '3': 6, '4': 1, '5': 9, '10': 'notes'},
    {
      '1': 'medications',
      '3': 7,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.Medication',
      '10': 'medications'
    },
  ],
};

/// Descriptor for `CreatePrescriptionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createPrescriptionRequestDescriptor = $convert.base64Decode(
    'ChlDcmVhdGVQcmVzY3JpcHRpb25SZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZBIbCg'
    'lkb2N0b3JfaWQYAiABKAlSCGRvY3RvcklkEh8KC2RvY3Rvcl9uYW1lGAMgASgJUgpkb2N0b3JO'
    'YW1lEh8KC2ZhY2lsaXR5X2lkGAQgASgJUgpmYWNpbGl0eUlkEhwKCWRpYWdub3NpcxgFIAEoCV'
    'IJZGlhZ25vc2lzEhQKBW5vdGVzGAYgASgJUgVub3RlcxI5CgttZWRpY2F0aW9ucxgHIAMoCzIX'
    'Lm1hbnBhc2lrLnYxLk1lZGljYXRpb25SC21lZGljYXRpb25z');

@$core.Deprecated('Use getPrescriptionRequestDescriptor instead')
const GetPrescriptionRequest$json = {
  '1': 'GetPrescriptionRequest',
  '2': [
    {'1': 'prescription_id', '3': 1, '4': 1, '5': 9, '10': 'prescriptionId'},
  ],
};

/// Descriptor for `GetPrescriptionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPrescriptionRequestDescriptor =
    $convert.base64Decode(
        'ChZHZXRQcmVzY3JpcHRpb25SZXF1ZXN0EicKD3ByZXNjcmlwdGlvbl9pZBgBIAEoCVIOcHJlc2'
        'NyaXB0aW9uSWQ=');

@$core.Deprecated('Use listPrescriptionsRequestDescriptor instead')
const ListPrescriptionsRequest$json = {
  '1': 'ListPrescriptionsRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'status_filter',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.PrescriptionStatus',
      '10': 'statusFilter'
    },
    {'1': 'limit', '3': 3, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 4, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `ListPrescriptionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPrescriptionsRequestDescriptor = $convert.base64Decode(
    'ChhMaXN0UHJlc2NyaXB0aW9uc1JlcXVlc3QSFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklkEkQKDX'
    'N0YXR1c19maWx0ZXIYAiABKA4yHy5tYW5wYXNpay52MS5QcmVzY3JpcHRpb25TdGF0dXNSDHN0'
    'YXR1c0ZpbHRlchIUCgVsaW1pdBgDIAEoBVIFbGltaXQSFgoGb2Zmc2V0GAQgASgFUgZvZmZzZX'
    'Q=');

@$core.Deprecated('Use listPrescriptionsResponseDescriptor instead')
const ListPrescriptionsResponse$json = {
  '1': 'ListPrescriptionsResponse',
  '2': [
    {
      '1': 'prescriptions',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.Prescription',
      '10': 'prescriptions'
    },
    {'1': 'total_count', '3': 2, '4': 1, '5': 5, '10': 'totalCount'},
  ],
};

/// Descriptor for `ListPrescriptionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPrescriptionsResponseDescriptor = $convert.base64Decode(
    'ChlMaXN0UHJlc2NyaXB0aW9uc1Jlc3BvbnNlEj8KDXByZXNjcmlwdGlvbnMYASADKAsyGS5tYW'
    '5wYXNpay52MS5QcmVzY3JpcHRpb25SDXByZXNjcmlwdGlvbnMSHwoLdG90YWxfY291bnQYAiAB'
    'KAVSCnRvdGFsQ291bnQ=');

@$core.Deprecated('Use updatePrescriptionStatusRequestDescriptor instead')
const UpdatePrescriptionStatusRequest$json = {
  '1': 'UpdatePrescriptionStatusRequest',
  '2': [
    {'1': 'prescription_id', '3': 1, '4': 1, '5': 9, '10': 'prescriptionId'},
    {
      '1': 'new_status',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.PrescriptionStatus',
      '10': 'newStatus'
    },
  ],
};

/// Descriptor for `UpdatePrescriptionStatusRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updatePrescriptionStatusRequestDescriptor =
    $convert.base64Decode(
        'Ch9VcGRhdGVQcmVzY3JpcHRpb25TdGF0dXNSZXF1ZXN0EicKD3ByZXNjcmlwdGlvbl9pZBgBIA'
        'EoCVIOcHJlc2NyaXB0aW9uSWQSPgoKbmV3X3N0YXR1cxgCIAEoDjIfLm1hbnBhc2lrLnYxLlBy'
        'ZXNjcmlwdGlvblN0YXR1c1IJbmV3U3RhdHVz');

@$core.Deprecated('Use addMedicationRequestDescriptor instead')
const AddMedicationRequest$json = {
  '1': 'AddMedicationRequest',
  '2': [
    {'1': 'prescription_id', '3': 1, '4': 1, '5': 9, '10': 'prescriptionId'},
    {
      '1': 'medication',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.manpasik.v1.Medication',
      '10': 'medication'
    },
  ],
};

/// Descriptor for `AddMedicationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addMedicationRequestDescriptor = $convert.base64Decode(
    'ChRBZGRNZWRpY2F0aW9uUmVxdWVzdBInCg9wcmVzY3JpcHRpb25faWQYASABKAlSDnByZXNjcm'
    'lwdGlvbklkEjcKCm1lZGljYXRpb24YAiABKAsyFy5tYW5wYXNpay52MS5NZWRpY2F0aW9uUgpt'
    'ZWRpY2F0aW9u');

@$core.Deprecated('Use removeMedicationRequestDescriptor instead')
const RemoveMedicationRequest$json = {
  '1': 'RemoveMedicationRequest',
  '2': [
    {'1': 'prescription_id', '3': 1, '4': 1, '5': 9, '10': 'prescriptionId'},
    {'1': 'medication_id', '3': 2, '4': 1, '5': 9, '10': 'medicationId'},
  ],
};

/// Descriptor for `RemoveMedicationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removeMedicationRequestDescriptor =
    $convert.base64Decode(
        'ChdSZW1vdmVNZWRpY2F0aW9uUmVxdWVzdBInCg9wcmVzY3JpcHRpb25faWQYASABKAlSDnByZX'
        'NjcmlwdGlvbklkEiMKDW1lZGljYXRpb25faWQYAiABKAlSDG1lZGljYXRpb25JZA==');

@$core.Deprecated('Use prescriptionDescriptor instead')
const Prescription$json = {
  '1': 'Prescription',
  '2': [
    {'1': 'prescription_id', '3': 1, '4': 1, '5': 9, '10': 'prescriptionId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'doctor_id', '3': 3, '4': 1, '5': 9, '10': 'doctorId'},
    {'1': 'doctor_name', '3': 4, '4': 1, '5': 9, '10': 'doctorName'},
    {'1': 'facility_id', '3': 5, '4': 1, '5': 9, '10': 'facilityId'},
    {'1': 'diagnosis', '3': 6, '4': 1, '5': 9, '10': 'diagnosis'},
    {'1': 'notes', '3': 7, '4': 1, '5': 9, '10': 'notes'},
    {
      '1': 'status',
      '3': 8,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.PrescriptionStatus',
      '10': 'status'
    },
    {
      '1': 'medications',
      '3': 9,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.Medication',
      '10': 'medications'
    },
    {
      '1': 'prescribed_at',
      '3': 10,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'prescribedAt'
    },
    {
      '1': 'expires_at',
      '3': 11,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'expiresAt'
    },
    {
      '1': 'updated_at',
      '3': 12,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'updatedAt'
    },
  ],
};

/// Descriptor for `Prescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List prescriptionDescriptor = $convert.base64Decode(
    'CgxQcmVzY3JpcHRpb24SJwoPcHJlc2NyaXB0aW9uX2lkGAEgASgJUg5wcmVzY3JpcHRpb25JZB'
    'IXCgd1c2VyX2lkGAIgASgJUgZ1c2VySWQSGwoJZG9jdG9yX2lkGAMgASgJUghkb2N0b3JJZBIf'
    'Cgtkb2N0b3JfbmFtZRgEIAEoCVIKZG9jdG9yTmFtZRIfCgtmYWNpbGl0eV9pZBgFIAEoCVIKZm'
    'FjaWxpdHlJZBIcCglkaWFnbm9zaXMYBiABKAlSCWRpYWdub3NpcxIUCgVub3RlcxgHIAEoCVIF'
    'bm90ZXMSNwoGc3RhdHVzGAggASgOMh8ubWFucGFzaWsudjEuUHJlc2NyaXB0aW9uU3RhdHVzUg'
    'ZzdGF0dXMSOQoLbWVkaWNhdGlvbnMYCSADKAsyFy5tYW5wYXNpay52MS5NZWRpY2F0aW9uUgtt'
    'ZWRpY2F0aW9ucxI/Cg1wcmVzY3JpYmVkX2F0GAogASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbW'
    'VzdGFtcFIMcHJlc2NyaWJlZEF0EjkKCmV4cGlyZXNfYXQYCyABKAsyGi5nb29nbGUucHJvdG9i'
    'dWYuVGltZXN0YW1wUglleHBpcmVzQXQSOQoKdXBkYXRlZF9hdBgMIAEoCzIaLmdvb2dsZS5wcm'
    '90b2J1Zi5UaW1lc3RhbXBSCXVwZGF0ZWRBdA==');

@$core.Deprecated('Use medicationDescriptor instead')
const Medication$json = {
  '1': 'Medication',
  '2': [
    {'1': 'medication_id', '3': 1, '4': 1, '5': 9, '10': 'medicationId'},
    {'1': 'name', '3': 2, '4': 1, '5': 9, '10': 'name'},
    {'1': 'dosage', '3': 3, '4': 1, '5': 9, '10': 'dosage'},
    {'1': 'frequency', '3': 4, '4': 1, '5': 9, '10': 'frequency'},
    {'1': 'route', '3': 5, '4': 1, '5': 9, '10': 'route'},
    {'1': 'duration_days', '3': 6, '4': 1, '5': 5, '10': 'durationDays'},
    {'1': 'instructions', '3': 7, '4': 1, '5': 9, '10': 'instructions'},
    {'1': 'is_critical', '3': 8, '4': 1, '5': 8, '10': 'isCritical'},
  ],
};

/// Descriptor for `Medication`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List medicationDescriptor = $convert.base64Decode(
    'CgpNZWRpY2F0aW9uEiMKDW1lZGljYXRpb25faWQYASABKAlSDG1lZGljYXRpb25JZBISCgRuYW'
    '1lGAIgASgJUgRuYW1lEhYKBmRvc2FnZRgDIAEoCVIGZG9zYWdlEhwKCWZyZXF1ZW5jeRgEIAEo'
    'CVIJZnJlcXVlbmN5EhQKBXJvdXRlGAUgASgJUgVyb3V0ZRIjCg1kdXJhdGlvbl9kYXlzGAYgAS'
    'gFUgxkdXJhdGlvbkRheXMSIgoMaW5zdHJ1Y3Rpb25zGAcgASgJUgxpbnN0cnVjdGlvbnMSHwoL'
    'aXNfY3JpdGljYWwYCCABKAhSCmlzQ3JpdGljYWw=');

@$core.Deprecated('Use checkDrugInteractionRequestDescriptor instead')
const CheckDrugInteractionRequest$json = {
  '1': 'CheckDrugInteractionRequest',
  '2': [
    {'1': 'medication_names', '3': 1, '4': 3, '5': 9, '10': 'medicationNames'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `CheckDrugInteractionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List checkDrugInteractionRequestDescriptor =
    $convert.base64Decode(
        'ChtDaGVja0RydWdJbnRlcmFjdGlvblJlcXVlc3QSKQoQbWVkaWNhdGlvbl9uYW1lcxgBIAMoCV'
        'IPbWVkaWNhdGlvbk5hbWVzEhcKB3VzZXJfaWQYAiABKAlSBnVzZXJJZA==');

@$core.Deprecated('Use checkDrugInteractionResponseDescriptor instead')
const CheckDrugInteractionResponse$json = {
  '1': 'CheckDrugInteractionResponse',
  '2': [
    {
      '1': 'interactions',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.DrugInteraction',
      '10': 'interactions'
    },
    {'1': 'has_critical', '3': 2, '4': 1, '5': 8, '10': 'hasCritical'},
  ],
};

/// Descriptor for `CheckDrugInteractionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List checkDrugInteractionResponseDescriptor =
    $convert.base64Decode(
        'ChxDaGVja0RydWdJbnRlcmFjdGlvblJlc3BvbnNlEkAKDGludGVyYWN0aW9ucxgBIAMoCzIcLm'
        '1hbnBhc2lrLnYxLkRydWdJbnRlcmFjdGlvblIMaW50ZXJhY3Rpb25zEiEKDGhhc19jcml0aWNh'
        'bBgCIAEoCFILaGFzQ3JpdGljYWw=');

@$core.Deprecated('Use drugInteractionDescriptor instead')
const DrugInteraction$json = {
  '1': 'DrugInteraction',
  '2': [
    {'1': 'drug_a', '3': 1, '4': 1, '5': 9, '10': 'drugA'},
    {'1': 'drug_b', '3': 2, '4': 1, '5': 9, '10': 'drugB'},
    {
      '1': 'severity',
      '3': 3,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.DrugInteractionSeverity',
      '10': 'severity'
    },
    {'1': 'description', '3': 4, '4': 1, '5': 9, '10': 'description'},
    {'1': 'recommendation', '3': 5, '4': 1, '5': 9, '10': 'recommendation'},
  ],
};

/// Descriptor for `DrugInteraction`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List drugInteractionDescriptor = $convert.base64Decode(
    'Cg9EcnVnSW50ZXJhY3Rpb24SFQoGZHJ1Z19hGAEgASgJUgVkcnVnQRIVCgZkcnVnX2IYAiABKA'
    'lSBWRydWdCEkAKCHNldmVyaXR5GAMgASgOMiQubWFucGFzaWsudjEuRHJ1Z0ludGVyYWN0aW9u'
    'U2V2ZXJpdHlSCHNldmVyaXR5EiAKC2Rlc2NyaXB0aW9uGAQgASgJUgtkZXNjcmlwdGlvbhImCg'
    '5yZWNvbW1lbmRhdGlvbhgFIAEoCVIOcmVjb21tZW5kYXRpb24=');

@$core.Deprecated('Use getMedicationRemindersRequestDescriptor instead')
const GetMedicationRemindersRequest$json = {
  '1': 'GetMedicationRemindersRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'date', '3': 2, '4': 1, '5': 9, '10': 'date'},
  ],
};

/// Descriptor for `GetMedicationRemindersRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMedicationRemindersRequestDescriptor =
    $convert.base64Decode(
        'Ch1HZXRNZWRpY2F0aW9uUmVtaW5kZXJzUmVxdWVzdBIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySW'
        'QSEgoEZGF0ZRgCIAEoCVIEZGF0ZQ==');

@$core.Deprecated('Use getMedicationRemindersResponseDescriptor instead')
const GetMedicationRemindersResponse$json = {
  '1': 'GetMedicationRemindersResponse',
  '2': [
    {
      '1': 'reminders',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.MedicationReminder',
      '10': 'reminders'
    },
  ],
};

/// Descriptor for `GetMedicationRemindersResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMedicationRemindersResponseDescriptor =
    $convert.base64Decode(
        'Ch5HZXRNZWRpY2F0aW9uUmVtaW5kZXJzUmVzcG9uc2USPQoJcmVtaW5kZXJzGAEgAygLMh8ubW'
        'FucGFzaWsudjEuTWVkaWNhdGlvblJlbWluZGVyUglyZW1pbmRlcnM=');

@$core.Deprecated('Use medicationReminderDescriptor instead')
const MedicationReminder$json = {
  '1': 'MedicationReminder',
  '2': [
    {'1': 'reminder_id', '3': 1, '4': 1, '5': 9, '10': 'reminderId'},
    {'1': 'prescription_id', '3': 2, '4': 1, '5': 9, '10': 'prescriptionId'},
    {'1': 'medication_name', '3': 3, '4': 1, '5': 9, '10': 'medicationName'},
    {'1': 'dosage', '3': 4, '4': 1, '5': 9, '10': 'dosage'},
    {'1': 'scheduled_time', '3': 5, '4': 1, '5': 9, '10': 'scheduledTime'},
    {'1': 'is_taken', '3': 6, '4': 1, '5': 8, '10': 'isTaken'},
    {'1': 'instructions', '3': 7, '4': 1, '5': 9, '10': 'instructions'},
  ],
};

/// Descriptor for `MedicationReminder`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List medicationReminderDescriptor = $convert.base64Decode(
    'ChJNZWRpY2F0aW9uUmVtaW5kZXISHwoLcmVtaW5kZXJfaWQYASABKAlSCnJlbWluZGVySWQSJw'
    'oPcHJlc2NyaXB0aW9uX2lkGAIgASgJUg5wcmVzY3JpcHRpb25JZBInCg9tZWRpY2F0aW9uX25h'
    'bWUYAyABKAlSDm1lZGljYXRpb25OYW1lEhYKBmRvc2FnZRgEIAEoCVIGZG9zYWdlEiUKDnNjaG'
    'VkdWxlZF90aW1lGAUgASgJUg1zY2hlZHVsZWRUaW1lEhkKCGlzX3Rha2VuGAYgASgIUgdpc1Rh'
    'a2VuEiIKDGluc3RydWN0aW9ucxgHIAEoCVIMaW5zdHJ1Y3Rpb25z');

@$core.Deprecated('Use createPostRequestDescriptor instead')
const CreatePostRequest$json = {
  '1': 'CreatePostRequest',
  '2': [
    {'1': 'author_id', '3': 1, '4': 1, '5': 9, '10': 'authorId'},
    {'1': 'title', '3': 2, '4': 1, '5': 9, '10': 'title'},
    {'1': 'content', '3': 3, '4': 1, '5': 9, '10': 'content'},
    {
      '1': 'category',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.PostCategory',
      '10': 'category'
    },
    {'1': 'tags', '3': 5, '4': 3, '5': 9, '10': 'tags'},
  ],
};

/// Descriptor for `CreatePostRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createPostRequestDescriptor = $convert.base64Decode(
    'ChFDcmVhdGVQb3N0UmVxdWVzdBIbCglhdXRob3JfaWQYASABKAlSCGF1dGhvcklkEhQKBXRpdG'
    'xlGAIgASgJUgV0aXRsZRIYCgdjb250ZW50GAMgASgJUgdjb250ZW50EjUKCGNhdGVnb3J5GAQg'
    'ASgOMhkubWFucGFzaWsudjEuUG9zdENhdGVnb3J5UghjYXRlZ29yeRISCgR0YWdzGAUgAygJUg'
    'R0YWdz');

@$core.Deprecated('Use getPostRequestDescriptor instead')
const GetPostRequest$json = {
  '1': 'GetPostRequest',
  '2': [
    {'1': 'post_id', '3': 1, '4': 1, '5': 9, '10': 'postId'},
  ],
};

/// Descriptor for `GetPostRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPostRequestDescriptor = $convert
    .base64Decode('Cg5HZXRQb3N0UmVxdWVzdBIXCgdwb3N0X2lkGAEgASgJUgZwb3N0SWQ=');

@$core.Deprecated('Use listPostsRequestDescriptor instead')
const ListPostsRequest$json = {
  '1': 'ListPostsRequest',
  '2': [
    {
      '1': 'category',
      '3': 1,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.PostCategory',
      '10': 'category'
    },
    {'1': 'author_id', '3': 2, '4': 1, '5': 9, '10': 'authorId'},
    {'1': 'query', '3': 3, '4': 1, '5': 9, '10': 'query'},
    {'1': 'limit', '3': 4, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 5, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `ListPostsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPostsRequestDescriptor = $convert.base64Decode(
    'ChBMaXN0UG9zdHNSZXF1ZXN0EjUKCGNhdGVnb3J5GAEgASgOMhkubWFucGFzaWsudjEuUG9zdE'
    'NhdGVnb3J5UghjYXRlZ29yeRIbCglhdXRob3JfaWQYAiABKAlSCGF1dGhvcklkEhQKBXF1ZXJ5'
    'GAMgASgJUgVxdWVyeRIUCgVsaW1pdBgEIAEoBVIFbGltaXQSFgoGb2Zmc2V0GAUgASgFUgZvZm'
    'ZzZXQ=');

@$core.Deprecated('Use listPostsResponseDescriptor instead')
const ListPostsResponse$json = {
  '1': 'ListPostsResponse',
  '2': [
    {
      '1': 'posts',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.Post',
      '10': 'posts'
    },
    {'1': 'total_count', '3': 2, '4': 1, '5': 5, '10': 'totalCount'},
  ],
};

/// Descriptor for `ListPostsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPostsResponseDescriptor = $convert.base64Decode(
    'ChFMaXN0UG9zdHNSZXNwb25zZRInCgVwb3N0cxgBIAMoCzIRLm1hbnBhc2lrLnYxLlBvc3RSBX'
    'Bvc3RzEh8KC3RvdGFsX2NvdW50GAIgASgFUgp0b3RhbENvdW50');

@$core.Deprecated('Use postDescriptor instead')
const Post$json = {
  '1': 'Post',
  '2': [
    {'1': 'post_id', '3': 1, '4': 1, '5': 9, '10': 'postId'},
    {'1': 'author_id', '3': 2, '4': 1, '5': 9, '10': 'authorId'},
    {'1': 'author_name', '3': 3, '4': 1, '5': 9, '10': 'authorName'},
    {'1': 'title', '3': 4, '4': 1, '5': 9, '10': 'title'},
    {'1': 'content', '3': 5, '4': 1, '5': 9, '10': 'content'},
    {
      '1': 'category',
      '3': 6,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.PostCategory',
      '10': 'category'
    },
    {'1': 'tags', '3': 7, '4': 3, '5': 9, '10': 'tags'},
    {'1': 'like_count', '3': 8, '4': 1, '5': 5, '10': 'likeCount'},
    {'1': 'comment_count', '3': 9, '4': 1, '5': 5, '10': 'commentCount'},
    {'1': 'view_count', '3': 10, '4': 1, '5': 5, '10': 'viewCount'},
    {
      '1': 'created_at',
      '3': 11,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
    {
      '1': 'updated_at',
      '3': 12,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'updatedAt'
    },
  ],
};

/// Descriptor for `Post`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List postDescriptor = $convert.base64Decode(
    'CgRQb3N0EhcKB3Bvc3RfaWQYASABKAlSBnBvc3RJZBIbCglhdXRob3JfaWQYAiABKAlSCGF1dG'
    'hvcklkEh8KC2F1dGhvcl9uYW1lGAMgASgJUgphdXRob3JOYW1lEhQKBXRpdGxlGAQgASgJUgV0'
    'aXRsZRIYCgdjb250ZW50GAUgASgJUgdjb250ZW50EjUKCGNhdGVnb3J5GAYgASgOMhkubWFucG'
    'FzaWsudjEuUG9zdENhdGVnb3J5UghjYXRlZ29yeRISCgR0YWdzGAcgAygJUgR0YWdzEh0KCmxp'
    'a2VfY291bnQYCCABKAVSCWxpa2VDb3VudBIjCg1jb21tZW50X2NvdW50GAkgASgFUgxjb21tZW'
    '50Q291bnQSHQoKdmlld19jb3VudBgKIAEoBVIJdmlld0NvdW50EjkKCmNyZWF0ZWRfYXQYCyAB'
    'KAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUgljcmVhdGVkQXQSOQoKdXBkYXRlZF9hdB'
    'gMIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSCXVwZGF0ZWRBdA==');

@$core.Deprecated('Use likePostRequestDescriptor instead')
const LikePostRequest$json = {
  '1': 'LikePostRequest',
  '2': [
    {'1': 'post_id', '3': 1, '4': 1, '5': 9, '10': 'postId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `LikePostRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List likePostRequestDescriptor = $convert.base64Decode(
    'Cg9MaWtlUG9zdFJlcXVlc3QSFwoHcG9zdF9pZBgBIAEoCVIGcG9zdElkEhcKB3VzZXJfaWQYAi'
    'ABKAlSBnVzZXJJZA==');

@$core.Deprecated('Use likePostResponseDescriptor instead')
const LikePostResponse$json = {
  '1': 'LikePostResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
    {'1': 'like_count', '3': 2, '4': 1, '5': 5, '10': 'likeCount'},
    {'1': 'is_liked', '3': 3, '4': 1, '5': 8, '10': 'isLiked'},
  ],
};

/// Descriptor for `LikePostResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List likePostResponseDescriptor = $convert.base64Decode(
    'ChBMaWtlUG9zdFJlc3BvbnNlEhgKB3N1Y2Nlc3MYASABKAhSB3N1Y2Nlc3MSHQoKbGlrZV9jb3'
    'VudBgCIAEoBVIJbGlrZUNvdW50EhkKCGlzX2xpa2VkGAMgASgIUgdpc0xpa2Vk');

@$core.Deprecated('Use createCommentRequestDescriptor instead')
const CreateCommentRequest$json = {
  '1': 'CreateCommentRequest',
  '2': [
    {'1': 'post_id', '3': 1, '4': 1, '5': 9, '10': 'postId'},
    {'1': 'author_id', '3': 2, '4': 1, '5': 9, '10': 'authorId'},
    {'1': 'content', '3': 3, '4': 1, '5': 9, '10': 'content'},
    {'1': 'parent_comment_id', '3': 4, '4': 1, '5': 9, '10': 'parentCommentId'},
  ],
};

/// Descriptor for `CreateCommentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createCommentRequestDescriptor = $convert.base64Decode(
    'ChRDcmVhdGVDb21tZW50UmVxdWVzdBIXCgdwb3N0X2lkGAEgASgJUgZwb3N0SWQSGwoJYXV0aG'
    '9yX2lkGAIgASgJUghhdXRob3JJZBIYCgdjb250ZW50GAMgASgJUgdjb250ZW50EioKEXBhcmVu'
    'dF9jb21tZW50X2lkGAQgASgJUg9wYXJlbnRDb21tZW50SWQ=');

@$core.Deprecated('Use listCommentsRequestDescriptor instead')
const ListCommentsRequest$json = {
  '1': 'ListCommentsRequest',
  '2': [
    {'1': 'post_id', '3': 1, '4': 1, '5': 9, '10': 'postId'},
    {'1': 'limit', '3': 2, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 3, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `ListCommentsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCommentsRequestDescriptor = $convert.base64Decode(
    'ChNMaXN0Q29tbWVudHNSZXF1ZXN0EhcKB3Bvc3RfaWQYASABKAlSBnBvc3RJZBIUCgVsaW1pdB'
    'gCIAEoBVIFbGltaXQSFgoGb2Zmc2V0GAMgASgFUgZvZmZzZXQ=');

@$core.Deprecated('Use listCommentsResponseDescriptor instead')
const ListCommentsResponse$json = {
  '1': 'ListCommentsResponse',
  '2': [
    {
      '1': 'comments',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.Comment',
      '10': 'comments'
    },
    {'1': 'total_count', '3': 2, '4': 1, '5': 5, '10': 'totalCount'},
  ],
};

/// Descriptor for `ListCommentsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listCommentsResponseDescriptor = $convert.base64Decode(
    'ChRMaXN0Q29tbWVudHNSZXNwb25zZRIwCghjb21tZW50cxgBIAMoCzIULm1hbnBhc2lrLnYxLk'
    'NvbW1lbnRSCGNvbW1lbnRzEh8KC3RvdGFsX2NvdW50GAIgASgFUgp0b3RhbENvdW50');

@$core.Deprecated('Use commentDescriptor instead')
const Comment$json = {
  '1': 'Comment',
  '2': [
    {'1': 'comment_id', '3': 1, '4': 1, '5': 9, '10': 'commentId'},
    {'1': 'post_id', '3': 2, '4': 1, '5': 9, '10': 'postId'},
    {'1': 'author_id', '3': 3, '4': 1, '5': 9, '10': 'authorId'},
    {'1': 'author_name', '3': 4, '4': 1, '5': 9, '10': 'authorName'},
    {'1': 'content', '3': 5, '4': 1, '5': 9, '10': 'content'},
    {'1': 'parent_comment_id', '3': 6, '4': 1, '5': 9, '10': 'parentCommentId'},
    {'1': 'like_count', '3': 7, '4': 1, '5': 5, '10': 'likeCount'},
    {
      '1': 'created_at',
      '3': 8,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
  ],
};

/// Descriptor for `Comment`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List commentDescriptor = $convert.base64Decode(
    'CgdDb21tZW50Eh0KCmNvbW1lbnRfaWQYASABKAlSCWNvbW1lbnRJZBIXCgdwb3N0X2lkGAIgAS'
    'gJUgZwb3N0SWQSGwoJYXV0aG9yX2lkGAMgASgJUghhdXRob3JJZBIfCgthdXRob3JfbmFtZRgE'
    'IAEoCVIKYXV0aG9yTmFtZRIYCgdjb250ZW50GAUgASgJUgdjb250ZW50EioKEXBhcmVudF9jb2'
    '1tZW50X2lkGAYgASgJUg9wYXJlbnRDb21tZW50SWQSHQoKbGlrZV9jb3VudBgHIAEoBVIJbGlr'
    'ZUNvdW50EjkKCmNyZWF0ZWRfYXQYCCABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUg'
    'ljcmVhdGVkQXQ=');

@$core.Deprecated('Use createChallengeRequestDescriptor instead')
const CreateChallengeRequest$json = {
  '1': 'CreateChallengeRequest',
  '2': [
    {'1': 'creator_id', '3': 1, '4': 1, '5': 9, '10': 'creatorId'},
    {'1': 'title', '3': 2, '4': 1, '5': 9, '10': 'title'},
    {'1': 'description', '3': 3, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'challenge_type',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.ChallengeType',
      '10': 'challengeType'
    },
    {'1': 'target_value', '3': 5, '4': 1, '5': 1, '10': 'targetValue'},
    {'1': 'unit', '3': 6, '4': 1, '5': 9, '10': 'unit'},
    {
      '1': 'start_date',
      '3': 7,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'startDate'
    },
    {
      '1': 'end_date',
      '3': 8,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'endDate'
    },
    {'1': 'max_participants', '3': 9, '4': 1, '5': 5, '10': 'maxParticipants'},
  ],
};

/// Descriptor for `CreateChallengeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createChallengeRequestDescriptor = $convert.base64Decode(
    'ChZDcmVhdGVDaGFsbGVuZ2VSZXF1ZXN0Eh0KCmNyZWF0b3JfaWQYASABKAlSCWNyZWF0b3JJZB'
    'IUCgV0aXRsZRgCIAEoCVIFdGl0bGUSIAoLZGVzY3JpcHRpb24YAyABKAlSC2Rlc2NyaXB0aW9u'
    'EkEKDmNoYWxsZW5nZV90eXBlGAQgASgOMhoubWFucGFzaWsudjEuQ2hhbGxlbmdlVHlwZVINY2'
    'hhbGxlbmdlVHlwZRIhCgx0YXJnZXRfdmFsdWUYBSABKAFSC3RhcmdldFZhbHVlEhIKBHVuaXQY'
    'BiABKAlSBHVuaXQSOQoKc3RhcnRfZGF0ZRgHIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3'
    'RhbXBSCXN0YXJ0RGF0ZRI1CghlbmRfZGF0ZRgIIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1l'
    'c3RhbXBSB2VuZERhdGUSKQoQbWF4X3BhcnRpY2lwYW50cxgJIAEoBVIPbWF4UGFydGljaXBhbn'
    'Rz');

@$core.Deprecated('Use getChallengeRequestDescriptor instead')
const GetChallengeRequest$json = {
  '1': 'GetChallengeRequest',
  '2': [
    {'1': 'challenge_id', '3': 1, '4': 1, '5': 9, '10': 'challengeId'},
  ],
};

/// Descriptor for `GetChallengeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getChallengeRequestDescriptor = $convert.base64Decode(
    'ChNHZXRDaGFsbGVuZ2VSZXF1ZXN0EiEKDGNoYWxsZW5nZV9pZBgBIAEoCVILY2hhbGxlbmdlSW'
    'Q=');

@$core.Deprecated('Use challengeDescriptor instead')
const Challenge$json = {
  '1': 'Challenge',
  '2': [
    {'1': 'challenge_id', '3': 1, '4': 1, '5': 9, '10': 'challengeId'},
    {'1': 'creator_id', '3': 2, '4': 1, '5': 9, '10': 'creatorId'},
    {'1': 'title', '3': 3, '4': 1, '5': 9, '10': 'title'},
    {'1': 'description', '3': 4, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'challenge_type',
      '3': 5,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.ChallengeType',
      '10': 'challengeType'
    },
    {
      '1': 'status',
      '3': 6,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.ChallengeStatus',
      '10': 'status'
    },
    {'1': 'target_value', '3': 7, '4': 1, '5': 1, '10': 'targetValue'},
    {'1': 'unit', '3': 8, '4': 1, '5': 9, '10': 'unit'},
    {
      '1': 'participant_count',
      '3': 9,
      '4': 1,
      '5': 5,
      '10': 'participantCount'
    },
    {'1': 'max_participants', '3': 10, '4': 1, '5': 5, '10': 'maxParticipants'},
    {
      '1': 'start_date',
      '3': 11,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'startDate'
    },
    {
      '1': 'end_date',
      '3': 12,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'endDate'
    },
    {
      '1': 'created_at',
      '3': 13,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
  ],
};

/// Descriptor for `Challenge`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List challengeDescriptor = $convert.base64Decode(
    'CglDaGFsbGVuZ2USIQoMY2hhbGxlbmdlX2lkGAEgASgJUgtjaGFsbGVuZ2VJZBIdCgpjcmVhdG'
    '9yX2lkGAIgASgJUgljcmVhdG9ySWQSFAoFdGl0bGUYAyABKAlSBXRpdGxlEiAKC2Rlc2NyaXB0'
    'aW9uGAQgASgJUgtkZXNjcmlwdGlvbhJBCg5jaGFsbGVuZ2VfdHlwZRgFIAEoDjIaLm1hbnBhc2'
    'lrLnYxLkNoYWxsZW5nZVR5cGVSDWNoYWxsZW5nZVR5cGUSNAoGc3RhdHVzGAYgASgOMhwubWFu'
    'cGFzaWsudjEuQ2hhbGxlbmdlU3RhdHVzUgZzdGF0dXMSIQoMdGFyZ2V0X3ZhbHVlGAcgASgBUg'
    't0YXJnZXRWYWx1ZRISCgR1bml0GAggASgJUgR1bml0EisKEXBhcnRpY2lwYW50X2NvdW50GAkg'
    'ASgFUhBwYXJ0aWNpcGFudENvdW50EikKEG1heF9wYXJ0aWNpcGFudHMYCiABKAVSD21heFBhcn'
    'RpY2lwYW50cxI5CgpzdGFydF9kYXRlGAsgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFt'
    'cFIJc3RhcnREYXRlEjUKCGVuZF9kYXRlGAwgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdG'
    'FtcFIHZW5kRGF0ZRI5CgpjcmVhdGVkX2F0GA0gASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVz'
    'dGFtcFIJY3JlYXRlZEF0');

@$core.Deprecated('Use joinChallengeRequestDescriptor instead')
const JoinChallengeRequest$json = {
  '1': 'JoinChallengeRequest',
  '2': [
    {'1': 'challenge_id', '3': 1, '4': 1, '5': 9, '10': 'challengeId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `JoinChallengeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List joinChallengeRequestDescriptor = $convert.base64Decode(
    'ChRKb2luQ2hhbGxlbmdlUmVxdWVzdBIhCgxjaGFsbGVuZ2VfaWQYASABKAlSC2NoYWxsZW5nZU'
    'lkEhcKB3VzZXJfaWQYAiABKAlSBnVzZXJJZA==');

@$core.Deprecated('Use joinChallengeResponseDescriptor instead')
const JoinChallengeResponse$json = {
  '1': 'JoinChallengeResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
    {'1': 'message', '3': 2, '4': 1, '5': 9, '10': 'message'},
    {
      '1': 'participant_count',
      '3': 3,
      '4': 1,
      '5': 5,
      '10': 'participantCount'
    },
  ],
};

/// Descriptor for `JoinChallengeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List joinChallengeResponseDescriptor = $convert.base64Decode(
    'ChVKb2luQ2hhbGxlbmdlUmVzcG9uc2USGAoHc3VjY2VzcxgBIAEoCFIHc3VjY2VzcxIYCgdtZX'
    'NzYWdlGAIgASgJUgdtZXNzYWdlEisKEXBhcnRpY2lwYW50X2NvdW50GAMgASgFUhBwYXJ0aWNp'
    'cGFudENvdW50');

@$core.Deprecated('Use listChallengesRequestDescriptor instead')
const ListChallengesRequest$json = {
  '1': 'ListChallengesRequest',
  '2': [
    {
      '1': 'type_filter',
      '3': 1,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.ChallengeType',
      '10': 'typeFilter'
    },
    {
      '1': 'status_filter',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.ChallengeStatus',
      '10': 'statusFilter'
    },
    {'1': 'limit', '3': 3, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 4, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `ListChallengesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listChallengesRequestDescriptor = $convert.base64Decode(
    'ChVMaXN0Q2hhbGxlbmdlc1JlcXVlc3QSOwoLdHlwZV9maWx0ZXIYASABKA4yGi5tYW5wYXNpay'
    '52MS5DaGFsbGVuZ2VUeXBlUgp0eXBlRmlsdGVyEkEKDXN0YXR1c19maWx0ZXIYAiABKA4yHC5t'
    'YW5wYXNpay52MS5DaGFsbGVuZ2VTdGF0dXNSDHN0YXR1c0ZpbHRlchIUCgVsaW1pdBgDIAEoBV'
    'IFbGltaXQSFgoGb2Zmc2V0GAQgASgFUgZvZmZzZXQ=');

@$core.Deprecated('Use listChallengesResponseDescriptor instead')
const ListChallengesResponse$json = {
  '1': 'ListChallengesResponse',
  '2': [
    {
      '1': 'challenges',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.Challenge',
      '10': 'challenges'
    },
    {'1': 'total_count', '3': 2, '4': 1, '5': 5, '10': 'totalCount'},
  ],
};

/// Descriptor for `ListChallengesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listChallengesResponseDescriptor = $convert.base64Decode(
    'ChZMaXN0Q2hhbGxlbmdlc1Jlc3BvbnNlEjYKCmNoYWxsZW5nZXMYASADKAsyFi5tYW5wYXNpay'
    '52MS5DaGFsbGVuZ2VSCmNoYWxsZW5nZXMSHwoLdG90YWxfY291bnQYAiABKAVSCnRvdGFsQ291'
    'bnQ=');

@$core.Deprecated('Use createRoomRequestDescriptor instead')
const CreateRoomRequest$json = {
  '1': 'CreateRoomRequest',
  '2': [
    {'1': 'host_user_id', '3': 1, '4': 1, '5': 9, '10': 'hostUserId'},
    {
      '1': 'room_type',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.RoomType',
      '10': 'roomType'
    },
    {'1': 'title', '3': 3, '4': 1, '5': 9, '10': 'title'},
    {'1': 'max_participants', '3': 4, '4': 1, '5': 5, '10': 'maxParticipants'},
    {'1': 'scheduled_at', '3': 5, '4': 1, '5': 9, '10': 'scheduledAt'},
    {'1': 'reservation_id', '3': 6, '4': 1, '5': 9, '10': 'reservationId'},
  ],
};

/// Descriptor for `CreateRoomRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createRoomRequestDescriptor = $convert.base64Decode(
    'ChFDcmVhdGVSb29tUmVxdWVzdBIgCgxob3N0X3VzZXJfaWQYASABKAlSCmhvc3RVc2VySWQSMg'
    'oJcm9vbV90eXBlGAIgASgOMhUubWFucGFzaWsudjEuUm9vbVR5cGVSCHJvb21UeXBlEhQKBXRp'
    'dGxlGAMgASgJUgV0aXRsZRIpChBtYXhfcGFydGljaXBhbnRzGAQgASgFUg9tYXhQYXJ0aWNpcG'
    'FudHMSIQoMc2NoZWR1bGVkX2F0GAUgASgJUgtzY2hlZHVsZWRBdBIlCg5yZXNlcnZhdGlvbl9p'
    'ZBgGIAEoCVINcmVzZXJ2YXRpb25JZA==');

@$core.Deprecated('Use getRoomRequestDescriptor instead')
const GetRoomRequest$json = {
  '1': 'GetRoomRequest',
  '2': [
    {'1': 'room_id', '3': 1, '4': 1, '5': 9, '10': 'roomId'},
  ],
};

/// Descriptor for `GetRoomRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRoomRequestDescriptor = $convert
    .base64Decode('Cg5HZXRSb29tUmVxdWVzdBIXCgdyb29tX2lkGAEgASgJUgZyb29tSWQ=');

@$core.Deprecated('Use roomDescriptor instead')
const Room$json = {
  '1': 'Room',
  '2': [
    {'1': 'room_id', '3': 1, '4': 1, '5': 9, '10': 'roomId'},
    {'1': 'host_user_id', '3': 2, '4': 1, '5': 9, '10': 'hostUserId'},
    {
      '1': 'room_type',
      '3': 3,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.RoomType',
      '10': 'roomType'
    },
    {
      '1': 'status',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.RoomStatus',
      '10': 'status'
    },
    {'1': 'title', '3': 5, '4': 1, '5': 9, '10': 'title'},
    {
      '1': 'participant_count',
      '3': 6,
      '4': 1,
      '5': 5,
      '10': 'participantCount'
    },
    {'1': 'max_participants', '3': 7, '4': 1, '5': 5, '10': 'maxParticipants'},
    {'1': 'reservation_id', '3': 8, '4': 1, '5': 9, '10': 'reservationId'},
    {
      '1': 'created_at',
      '3': 9,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
    {
      '1': 'started_at',
      '3': 10,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'startedAt'
    },
    {
      '1': 'ended_at',
      '3': 11,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'endedAt'
    },
  ],
};

/// Descriptor for `Room`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List roomDescriptor = $convert.base64Decode(
    'CgRSb29tEhcKB3Jvb21faWQYASABKAlSBnJvb21JZBIgCgxob3N0X3VzZXJfaWQYAiABKAlSCm'
    'hvc3RVc2VySWQSMgoJcm9vbV90eXBlGAMgASgOMhUubWFucGFzaWsudjEuUm9vbVR5cGVSCHJv'
    'b21UeXBlEi8KBnN0YXR1cxgEIAEoDjIXLm1hbnBhc2lrLnYxLlJvb21TdGF0dXNSBnN0YXR1cx'
    'IUCgV0aXRsZRgFIAEoCVIFdGl0bGUSKwoRcGFydGljaXBhbnRfY291bnQYBiABKAVSEHBhcnRp'
    'Y2lwYW50Q291bnQSKQoQbWF4X3BhcnRpY2lwYW50cxgHIAEoBVIPbWF4UGFydGljaXBhbnRzEi'
    'UKDnJlc2VydmF0aW9uX2lkGAggASgJUg1yZXNlcnZhdGlvbklkEjkKCmNyZWF0ZWRfYXQYCSAB'
    'KAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUgljcmVhdGVkQXQSOQoKc3RhcnRlZF9hdB'
    'gKIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSCXN0YXJ0ZWRBdBI1CghlbmRlZF9h'
    'dBgLIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSB2VuZGVkQXQ=');

@$core.Deprecated('Use joinRoomRequestDescriptor instead')
const JoinRoomRequest$json = {
  '1': 'JoinRoomRequest',
  '2': [
    {'1': 'room_id', '3': 1, '4': 1, '5': 9, '10': 'roomId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `JoinRoomRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List joinRoomRequestDescriptor = $convert.base64Decode(
    'Cg9Kb2luUm9vbVJlcXVlc3QSFwoHcm9vbV9pZBgBIAEoCVIGcm9vbUlkEhcKB3VzZXJfaWQYAi'
    'ABKAlSBnVzZXJJZA==');

@$core.Deprecated('Use joinRoomResponseDescriptor instead')
const JoinRoomResponse$json = {
  '1': 'JoinRoomResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
    {'1': 'token', '3': 2, '4': 1, '5': 9, '10': 'token'},
    {
      '1': 'room',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.manpasik.v1.Room',
      '10': 'room'
    },
    {
      '1': 'participants',
      '3': 4,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.Participant',
      '10': 'participants'
    },
  ],
};

/// Descriptor for `JoinRoomResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List joinRoomResponseDescriptor = $convert.base64Decode(
    'ChBKb2luUm9vbVJlc3BvbnNlEhgKB3N1Y2Nlc3MYASABKAhSB3N1Y2Nlc3MSFAoFdG9rZW4YAi'
    'ABKAlSBXRva2VuEiUKBHJvb20YAyABKAsyES5tYW5wYXNpay52MS5Sb29tUgRyb29tEjwKDHBh'
    'cnRpY2lwYW50cxgEIAMoCzIYLm1hbnBhc2lrLnYxLlBhcnRpY2lwYW50UgxwYXJ0aWNpcGFudH'
    'M=');

@$core.Deprecated('Use leaveRoomRequestDescriptor instead')
const LeaveRoomRequest$json = {
  '1': 'LeaveRoomRequest',
  '2': [
    {'1': 'room_id', '3': 1, '4': 1, '5': 9, '10': 'roomId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `LeaveRoomRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List leaveRoomRequestDescriptor = $convert.base64Decode(
    'ChBMZWF2ZVJvb21SZXF1ZXN0EhcKB3Jvb21faWQYASABKAlSBnJvb21JZBIXCgd1c2VyX2lkGA'
    'IgASgJUgZ1c2VySWQ=');

@$core.Deprecated('Use leaveRoomResponseDescriptor instead')
const LeaveRoomResponse$json = {
  '1': 'LeaveRoomResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
    {'1': 'message', '3': 2, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `LeaveRoomResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List leaveRoomResponseDescriptor = $convert.base64Decode(
    'ChFMZWF2ZVJvb21SZXNwb25zZRIYCgdzdWNjZXNzGAEgASgIUgdzdWNjZXNzEhgKB21lc3NhZ2'
    'UYAiABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use endRoomRequestDescriptor instead')
const EndRoomRequest$json = {
  '1': 'EndRoomRequest',
  '2': [
    {'1': 'room_id', '3': 1, '4': 1, '5': 9, '10': 'roomId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `EndRoomRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List endRoomRequestDescriptor = $convert.base64Decode(
    'Cg5FbmRSb29tUmVxdWVzdBIXCgdyb29tX2lkGAEgASgJUgZyb29tSWQSFwoHdXNlcl9pZBgCIA'
    'EoCVIGdXNlcklk');

@$core.Deprecated('Use participantDescriptor instead')
const Participant$json = {
  '1': 'Participant',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'display_name', '3': 2, '4': 1, '5': 9, '10': 'displayName'},
    {'1': 'is_host', '3': 3, '4': 1, '5': 8, '10': 'isHost'},
    {'1': 'is_muted', '3': 4, '4': 1, '5': 8, '10': 'isMuted'},
    {'1': 'is_video_on', '3': 5, '4': 1, '5': 8, '10': 'isVideoOn'},
    {
      '1': 'joined_at',
      '3': 6,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'joinedAt'
    },
  ],
};

/// Descriptor for `Participant`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List participantDescriptor = $convert.base64Decode(
    'CgtQYXJ0aWNpcGFudBIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQSIQoMZGlzcGxheV9uYW1lGA'
    'IgASgJUgtkaXNwbGF5TmFtZRIXCgdpc19ob3N0GAMgASgIUgZpc0hvc3QSGQoIaXNfbXV0ZWQY'
    'BCABKAhSB2lzTXV0ZWQSHgoLaXNfdmlkZW9fb24YBSABKAhSCWlzVmlkZW9PbhI3Cglqb2luZW'
    'RfYXQYBiABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUghqb2luZWRBdA==');

@$core.Deprecated('Use sendSignalRequestDescriptor instead')
const SendSignalRequest$json = {
  '1': 'SendSignalRequest',
  '2': [
    {'1': 'room_id', '3': 1, '4': 1, '5': 9, '10': 'roomId'},
    {'1': 'from_user_id', '3': 2, '4': 1, '5': 9, '10': 'fromUserId'},
    {'1': 'to_user_id', '3': 3, '4': 1, '5': 9, '10': 'toUserId'},
    {
      '1': 'signal_type',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.SignalType',
      '10': 'signalType'
    },
    {'1': 'payload', '3': 5, '4': 1, '5': 9, '10': 'payload'},
  ],
};

/// Descriptor for `SendSignalRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendSignalRequestDescriptor = $convert.base64Decode(
    'ChFTZW5kU2lnbmFsUmVxdWVzdBIXCgdyb29tX2lkGAEgASgJUgZyb29tSWQSIAoMZnJvbV91c2'
    'VyX2lkGAIgASgJUgpmcm9tVXNlcklkEhwKCnRvX3VzZXJfaWQYAyABKAlSCHRvVXNlcklkEjgK'
    'C3NpZ25hbF90eXBlGAQgASgOMhcubWFucGFzaWsudjEuU2lnbmFsVHlwZVIKc2lnbmFsVHlwZR'
    'IYCgdwYXlsb2FkGAUgASgJUgdwYXlsb2Fk');

@$core.Deprecated('Use sendSignalResponseDescriptor instead')
const SendSignalResponse$json = {
  '1': 'SendSignalResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
  ],
};

/// Descriptor for `SendSignalResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendSignalResponseDescriptor =
    $convert.base64Decode(
        'ChJTZW5kU2lnbmFsUmVzcG9uc2USGAoHc3VjY2VzcxgBIAEoCFIHc3VjY2Vzcw==');

@$core.Deprecated('Use listParticipantsRequestDescriptor instead')
const ListParticipantsRequest$json = {
  '1': 'ListParticipantsRequest',
  '2': [
    {'1': 'room_id', '3': 1, '4': 1, '5': 9, '10': 'roomId'},
  ],
};

/// Descriptor for `ListParticipantsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listParticipantsRequestDescriptor =
    $convert.base64Decode(
        'ChdMaXN0UGFydGljaXBhbnRzUmVxdWVzdBIXCgdyb29tX2lkGAEgASgJUgZyb29tSWQ=');

@$core.Deprecated('Use listParticipantsResponseDescriptor instead')
const ListParticipantsResponse$json = {
  '1': 'ListParticipantsResponse',
  '2': [
    {
      '1': 'participants',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.Participant',
      '10': 'participants'
    },
  ],
};

/// Descriptor for `ListParticipantsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listParticipantsResponseDescriptor =
    $convert.base64Decode(
        'ChhMaXN0UGFydGljaXBhbnRzUmVzcG9uc2USPAoMcGFydGljaXBhbnRzGAEgAygLMhgubWFucG'
        'FzaWsudjEuUGFydGljaXBhbnRSDHBhcnRpY2lwYW50cw==');

@$core.Deprecated('Use getRoomStatsRequestDescriptor instead')
const GetRoomStatsRequest$json = {
  '1': 'GetRoomStatsRequest',
  '2': [
    {'1': 'room_id', '3': 1, '4': 1, '5': 9, '10': 'roomId'},
  ],
};

/// Descriptor for `GetRoomStatsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRoomStatsRequestDescriptor =
    $convert.base64Decode(
        'ChNHZXRSb29tU3RhdHNSZXF1ZXN0EhcKB3Jvb21faWQYASABKAlSBnJvb21JZA==');

@$core.Deprecated('Use getRoomStatsResponseDescriptor instead')
const GetRoomStatsResponse$json = {
  '1': 'GetRoomStatsResponse',
  '2': [
    {'1': 'room_id', '3': 1, '4': 1, '5': 9, '10': 'roomId'},
    {
      '1': 'total_participants',
      '3': 2,
      '4': 1,
      '5': 5,
      '10': 'totalParticipants'
    },
    {
      '1': 'current_participants',
      '3': 3,
      '4': 1,
      '5': 5,
      '10': 'currentParticipants'
    },
    {'1': 'duration_seconds', '3': 4, '4': 1, '5': 3, '10': 'durationSeconds'},
    {'1': 'signal_count', '3': 5, '4': 1, '5': 5, '10': 'signalCount'},
    {
      '1': 'started_at',
      '3': 6,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'startedAt'
    },
  ],
};

/// Descriptor for `GetRoomStatsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRoomStatsResponseDescriptor = $convert.base64Decode(
    'ChRHZXRSb29tU3RhdHNSZXNwb25zZRIXCgdyb29tX2lkGAEgASgJUgZyb29tSWQSLQoSdG90YW'
    'xfcGFydGljaXBhbnRzGAIgASgFUhF0b3RhbFBhcnRpY2lwYW50cxIxChRjdXJyZW50X3BhcnRp'
    'Y2lwYW50cxgDIAEoBVITY3VycmVudFBhcnRpY2lwYW50cxIpChBkdXJhdGlvbl9zZWNvbmRzGA'
    'QgASgDUg9kdXJhdGlvblNlY29uZHMSIQoMc2lnbmFsX2NvdW50GAUgASgFUgtzaWduYWxDb3Vu'
    'dBI5CgpzdGFydGVkX2F0GAYgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFIJc3Rhcn'
    'RlZEF0');

@$core.Deprecated('Use sendNotificationRequestDescriptor instead')
const SendNotificationRequest$json = {
  '1': 'SendNotificationRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'type',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.NotificationType',
      '10': 'type'
    },
    {'1': 'title', '3': 3, '4': 1, '5': 9, '10': 'title'},
    {'1': 'body', '3': 4, '4': 1, '5': 9, '10': 'body'},
    {
      '1': 'priority',
      '3': 5,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.NotificationPriority',
      '10': 'priority'
    },
    {
      '1': 'channel',
      '3': 6,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.NotificationChannel',
      '10': 'channel'
    },
    {
      '1': 'data',
      '3': 7,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.SendNotificationRequest.DataEntry',
      '10': 'data'
    },
    {'1': 'action_url', '3': 8, '4': 1, '5': 9, '10': 'actionUrl'},
  ],
  '3': [SendNotificationRequest_DataEntry$json],
};

@$core.Deprecated('Use sendNotificationRequestDescriptor instead')
const SendNotificationRequest_DataEntry$json = {
  '1': 'DataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `SendNotificationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendNotificationRequestDescriptor = $convert.base64Decode(
    'ChdTZW5kTm90aWZpY2F0aW9uUmVxdWVzdBIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQSMQoEdH'
    'lwZRgCIAEoDjIdLm1hbnBhc2lrLnYxLk5vdGlmaWNhdGlvblR5cGVSBHR5cGUSFAoFdGl0bGUY'
    'AyABKAlSBXRpdGxlEhIKBGJvZHkYBCABKAlSBGJvZHkSPQoIcHJpb3JpdHkYBSABKA4yIS5tYW'
    '5wYXNpay52MS5Ob3RpZmljYXRpb25Qcmlvcml0eVIIcHJpb3JpdHkSOgoHY2hhbm5lbBgGIAEo'
    'DjIgLm1hbnBhc2lrLnYxLk5vdGlmaWNhdGlvbkNoYW5uZWxSB2NoYW5uZWwSQgoEZGF0YRgHIA'
    'MoCzIuLm1hbnBhc2lrLnYxLlNlbmROb3RpZmljYXRpb25SZXF1ZXN0LkRhdGFFbnRyeVIEZGF0'
    'YRIdCgphY3Rpb25fdXJsGAggASgJUglhY3Rpb25VcmwaNwoJRGF0YUVudHJ5EhAKA2tleRgBIA'
    'EoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use notificationDescriptor instead')
const Notification$json = {
  '1': 'Notification',
  '2': [
    {'1': 'notification_id', '3': 1, '4': 1, '5': 9, '10': 'notificationId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'type',
      '3': 3,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.NotificationType',
      '10': 'type'
    },
    {'1': 'title', '3': 4, '4': 1, '5': 9, '10': 'title'},
    {'1': 'body', '3': 5, '4': 1, '5': 9, '10': 'body'},
    {
      '1': 'priority',
      '3': 6,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.NotificationPriority',
      '10': 'priority'
    },
    {
      '1': 'channel',
      '3': 7,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.NotificationChannel',
      '10': 'channel'
    },
    {'1': 'is_read', '3': 8, '4': 1, '5': 8, '10': 'isRead'},
    {
      '1': 'data',
      '3': 9,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.Notification.DataEntry',
      '10': 'data'
    },
    {'1': 'action_url', '3': 10, '4': 1, '5': 9, '10': 'actionUrl'},
    {
      '1': 'created_at',
      '3': 11,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
    {
      '1': 'read_at',
      '3': 12,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'readAt'
    },
  ],
  '3': [Notification_DataEntry$json],
};

@$core.Deprecated('Use notificationDescriptor instead')
const Notification_DataEntry$json = {
  '1': 'DataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `Notification`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List notificationDescriptor = $convert.base64Decode(
    'CgxOb3RpZmljYXRpb24SJwoPbm90aWZpY2F0aW9uX2lkGAEgASgJUg5ub3RpZmljYXRpb25JZB'
    'IXCgd1c2VyX2lkGAIgASgJUgZ1c2VySWQSMQoEdHlwZRgDIAEoDjIdLm1hbnBhc2lrLnYxLk5v'
    'dGlmaWNhdGlvblR5cGVSBHR5cGUSFAoFdGl0bGUYBCABKAlSBXRpdGxlEhIKBGJvZHkYBSABKA'
    'lSBGJvZHkSPQoIcHJpb3JpdHkYBiABKA4yIS5tYW5wYXNpay52MS5Ob3RpZmljYXRpb25Qcmlv'
    'cml0eVIIcHJpb3JpdHkSOgoHY2hhbm5lbBgHIAEoDjIgLm1hbnBhc2lrLnYxLk5vdGlmaWNhdG'
    'lvbkNoYW5uZWxSB2NoYW5uZWwSFwoHaXNfcmVhZBgIIAEoCFIGaXNSZWFkEjcKBGRhdGEYCSAD'
    'KAsyIy5tYW5wYXNpay52MS5Ob3RpZmljYXRpb24uRGF0YUVudHJ5UgRkYXRhEh0KCmFjdGlvbl'
    '91cmwYCiABKAlSCWFjdGlvblVybBI5CgpjcmVhdGVkX2F0GAsgASgLMhouZ29vZ2xlLnByb3Rv'
    'YnVmLlRpbWVzdGFtcFIJY3JlYXRlZEF0EjMKB3JlYWRfYXQYDCABKAsyGi5nb29nbGUucHJvdG'
    '9idWYuVGltZXN0YW1wUgZyZWFkQXQaNwoJRGF0YUVudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQK'
    'BXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use listNotificationsRequestDescriptor instead')
const ListNotificationsRequest$json = {
  '1': 'ListNotificationsRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'type_filter',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.NotificationType',
      '10': 'typeFilter'
    },
    {'1': 'unread_only', '3': 3, '4': 1, '5': 8, '10': 'unreadOnly'},
    {'1': 'limit', '3': 4, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 5, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `ListNotificationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listNotificationsRequestDescriptor = $convert.base64Decode(
    'ChhMaXN0Tm90aWZpY2F0aW9uc1JlcXVlc3QSFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklkEj4KC3'
    'R5cGVfZmlsdGVyGAIgASgOMh0ubWFucGFzaWsudjEuTm90aWZpY2F0aW9uVHlwZVIKdHlwZUZp'
    'bHRlchIfCgt1bnJlYWRfb25seRgDIAEoCFIKdW5yZWFkT25seRIUCgVsaW1pdBgEIAEoBVIFbG'
    'ltaXQSFgoGb2Zmc2V0GAUgASgFUgZvZmZzZXQ=');

@$core.Deprecated('Use listNotificationsResponseDescriptor instead')
const ListNotificationsResponse$json = {
  '1': 'ListNotificationsResponse',
  '2': [
    {
      '1': 'notifications',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.Notification',
      '10': 'notifications'
    },
    {'1': 'total_count', '3': 2, '4': 1, '5': 5, '10': 'totalCount'},
  ],
};

/// Descriptor for `ListNotificationsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listNotificationsResponseDescriptor = $convert.base64Decode(
    'ChlMaXN0Tm90aWZpY2F0aW9uc1Jlc3BvbnNlEj8KDW5vdGlmaWNhdGlvbnMYASADKAsyGS5tYW'
    '5wYXNpay52MS5Ob3RpZmljYXRpb25SDW5vdGlmaWNhdGlvbnMSHwoLdG90YWxfY291bnQYAiAB'
    'KAVSCnRvdGFsQ291bnQ=');

@$core.Deprecated('Use markAsReadRequestDescriptor instead')
const MarkAsReadRequest$json = {
  '1': 'MarkAsReadRequest',
  '2': [
    {'1': 'notification_id', '3': 1, '4': 1, '5': 9, '10': 'notificationId'},
  ],
};

/// Descriptor for `MarkAsReadRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List markAsReadRequestDescriptor = $convert.base64Decode(
    'ChFNYXJrQXNSZWFkUmVxdWVzdBInCg9ub3RpZmljYXRpb25faWQYASABKAlSDm5vdGlmaWNhdG'
    'lvbklk');

@$core.Deprecated('Use markAsReadResponseDescriptor instead')
const MarkAsReadResponse$json = {
  '1': 'MarkAsReadResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
  ],
};

/// Descriptor for `MarkAsReadResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List markAsReadResponseDescriptor =
    $convert.base64Decode(
        'ChJNYXJrQXNSZWFkUmVzcG9uc2USGAoHc3VjY2VzcxgBIAEoCFIHc3VjY2Vzcw==');

@$core.Deprecated('Use markAllAsReadRequestDescriptor instead')
const MarkAllAsReadRequest$json = {
  '1': 'MarkAllAsReadRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `MarkAllAsReadRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List markAllAsReadRequestDescriptor =
    $convert.base64Decode(
        'ChRNYXJrQWxsQXNSZWFkUmVxdWVzdBIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQ=');

@$core.Deprecated('Use markAllAsReadResponseDescriptor instead')
const MarkAllAsReadResponse$json = {
  '1': 'MarkAllAsReadResponse',
  '2': [
    {'1': 'marked_count', '3': 1, '4': 1, '5': 5, '10': 'markedCount'},
  ],
};

/// Descriptor for `MarkAllAsReadResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List markAllAsReadResponseDescriptor = $convert.base64Decode(
    'ChVNYXJrQWxsQXNSZWFkUmVzcG9uc2USIQoMbWFya2VkX2NvdW50GAEgASgFUgttYXJrZWRDb3'
    'VudA==');

@$core.Deprecated('Use getUnreadCountRequestDescriptor instead')
const GetUnreadCountRequest$json = {
  '1': 'GetUnreadCountRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `GetUnreadCountRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getUnreadCountRequestDescriptor =
    $convert.base64Decode(
        'ChVHZXRVbnJlYWRDb3VudFJlcXVlc3QSFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklk');

@$core.Deprecated('Use getUnreadCountResponseDescriptor instead')
const GetUnreadCountResponse$json = {
  '1': 'GetUnreadCountResponse',
  '2': [
    {'1': 'count', '3': 1, '4': 1, '5': 5, '10': 'count'},
  ],
};

/// Descriptor for `GetUnreadCountResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getUnreadCountResponseDescriptor =
    $convert.base64Decode(
        'ChZHZXRVbnJlYWRDb3VudFJlc3BvbnNlEhQKBWNvdW50GAEgASgFUgVjb3VudA==');

@$core.Deprecated('Use updateNotificationPreferencesRequestDescriptor instead')
const UpdateNotificationPreferencesRequest$json = {
  '1': 'UpdateNotificationPreferencesRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'push_enabled', '3': 2, '4': 1, '5': 8, '10': 'pushEnabled'},
    {'1': 'email_enabled', '3': 3, '4': 1, '5': 8, '10': 'emailEnabled'},
    {'1': 'sms_enabled', '3': 4, '4': 1, '5': 8, '10': 'smsEnabled'},
    {
      '1': 'measurement_alerts',
      '3': 5,
      '4': 1,
      '5': 8,
      '10': 'measurementAlerts'
    },
    {'1': 'health_alerts', '3': 6, '4': 1, '5': 8, '10': 'healthAlerts'},
    {
      '1': 'appointment_reminders',
      '3': 7,
      '4': 1,
      '5': 8,
      '10': 'appointmentReminders'
    },
    {
      '1': 'community_updates',
      '3': 8,
      '4': 1,
      '5': 8,
      '10': 'communityUpdates'
    },
    {'1': 'promotions', '3': 9, '4': 1, '5': 8, '10': 'promotions'},
    {
      '1': 'quiet_hours_start',
      '3': 10,
      '4': 1,
      '5': 9,
      '10': 'quietHoursStart'
    },
    {'1': 'quiet_hours_end', '3': 11, '4': 1, '5': 9, '10': 'quietHoursEnd'},
  ],
};

/// Descriptor for `UpdateNotificationPreferencesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateNotificationPreferencesRequestDescriptor = $convert.base64Decode(
    'CiRVcGRhdGVOb3RpZmljYXRpb25QcmVmZXJlbmNlc1JlcXVlc3QSFwoHdXNlcl9pZBgBIAEoCV'
    'IGdXNlcklkEiEKDHB1c2hfZW5hYmxlZBgCIAEoCFILcHVzaEVuYWJsZWQSIwoNZW1haWxfZW5h'
    'YmxlZBgDIAEoCFIMZW1haWxFbmFibGVkEh8KC3Ntc19lbmFibGVkGAQgASgIUgpzbXNFbmFibG'
    'VkEi0KEm1lYXN1cmVtZW50X2FsZXJ0cxgFIAEoCFIRbWVhc3VyZW1lbnRBbGVydHMSIwoNaGVh'
    'bHRoX2FsZXJ0cxgGIAEoCFIMaGVhbHRoQWxlcnRzEjMKFWFwcG9pbnRtZW50X3JlbWluZGVycx'
    'gHIAEoCFIUYXBwb2ludG1lbnRSZW1pbmRlcnMSKwoRY29tbXVuaXR5X3VwZGF0ZXMYCCABKAhS'
    'EGNvbW11bml0eVVwZGF0ZXMSHgoKcHJvbW90aW9ucxgJIAEoCFIKcHJvbW90aW9ucxIqChFxdW'
    'lldF9ob3Vyc19zdGFydBgKIAEoCVIPcXVpZXRIb3Vyc1N0YXJ0EiYKD3F1aWV0X2hvdXJzX2Vu'
    'ZBgLIAEoCVINcXVpZXRIb3Vyc0VuZA==');

@$core.Deprecated('Use getNotificationPreferencesRequestDescriptor instead')
const GetNotificationPreferencesRequest$json = {
  '1': 'GetNotificationPreferencesRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `GetNotificationPreferencesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getNotificationPreferencesRequestDescriptor =
    $convert.base64Decode(
        'CiFHZXROb3RpZmljYXRpb25QcmVmZXJlbmNlc1JlcXVlc3QSFwoHdXNlcl9pZBgBIAEoCVIGdX'
        'Nlcklk');

@$core.Deprecated('Use notificationPreferencesDescriptor instead')
const NotificationPreferences$json = {
  '1': 'NotificationPreferences',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'push_enabled', '3': 2, '4': 1, '5': 8, '10': 'pushEnabled'},
    {'1': 'email_enabled', '3': 3, '4': 1, '5': 8, '10': 'emailEnabled'},
    {'1': 'sms_enabled', '3': 4, '4': 1, '5': 8, '10': 'smsEnabled'},
    {
      '1': 'measurement_alerts',
      '3': 5,
      '4': 1,
      '5': 8,
      '10': 'measurementAlerts'
    },
    {'1': 'health_alerts', '3': 6, '4': 1, '5': 8, '10': 'healthAlerts'},
    {
      '1': 'appointment_reminders',
      '3': 7,
      '4': 1,
      '5': 8,
      '10': 'appointmentReminders'
    },
    {
      '1': 'community_updates',
      '3': 8,
      '4': 1,
      '5': 8,
      '10': 'communityUpdates'
    },
    {'1': 'promotions', '3': 9, '4': 1, '5': 8, '10': 'promotions'},
    {
      '1': 'quiet_hours_start',
      '3': 10,
      '4': 1,
      '5': 9,
      '10': 'quietHoursStart'
    },
    {'1': 'quiet_hours_end', '3': 11, '4': 1, '5': 9, '10': 'quietHoursEnd'},
    {
      '1': 'updated_at',
      '3': 12,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'updatedAt'
    },
  ],
};

/// Descriptor for `NotificationPreferences`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List notificationPreferencesDescriptor = $convert.base64Decode(
    'ChdOb3RpZmljYXRpb25QcmVmZXJlbmNlcxIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQSIQoMcH'
    'VzaF9lbmFibGVkGAIgASgIUgtwdXNoRW5hYmxlZBIjCg1lbWFpbF9lbmFibGVkGAMgASgIUgxl'
    'bWFpbEVuYWJsZWQSHwoLc21zX2VuYWJsZWQYBCABKAhSCnNtc0VuYWJsZWQSLQoSbWVhc3VyZW'
    '1lbnRfYWxlcnRzGAUgASgIUhFtZWFzdXJlbWVudEFsZXJ0cxIjCg1oZWFsdGhfYWxlcnRzGAYg'
    'ASgIUgxoZWFsdGhBbGVydHMSMwoVYXBwb2ludG1lbnRfcmVtaW5kZXJzGAcgASgIUhRhcHBvaW'
    '50bWVudFJlbWluZGVycxIrChFjb21tdW5pdHlfdXBkYXRlcxgIIAEoCFIQY29tbXVuaXR5VXBk'
    'YXRlcxIeCgpwcm9tb3Rpb25zGAkgASgIUgpwcm9tb3Rpb25zEioKEXF1aWV0X2hvdXJzX3N0YX'
    'J0GAogASgJUg9xdWlldEhvdXJzU3RhcnQSJgoPcXVpZXRfaG91cnNfZW5kGAsgASgJUg1xdWll'
    'dEhvdXJzRW5kEjkKCnVwZGF0ZWRfYXQYDCABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW'
    '1wUgl1cGRhdGVkQXQ=');

@$core.Deprecated('Use translateTextRequestDescriptor instead')
const TranslateTextRequest$json = {
  '1': 'TranslateTextRequest',
  '2': [
    {'1': 'text', '3': 1, '4': 1, '5': 9, '10': 'text'},
    {'1': 'source_language', '3': 2, '4': 1, '5': 9, '10': 'sourceLanguage'},
    {'1': 'target_language', '3': 3, '4': 1, '5': 9, '10': 'targetLanguage'},
    {'1': 'is_medical', '3': 4, '4': 1, '5': 8, '10': 'isMedical'},
    {'1': 'context', '3': 5, '4': 1, '5': 9, '10': 'context'},
  ],
};

/// Descriptor for `TranslateTextRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List translateTextRequestDescriptor = $convert.base64Decode(
    'ChRUcmFuc2xhdGVUZXh0UmVxdWVzdBISCgR0ZXh0GAEgASgJUgR0ZXh0EicKD3NvdXJjZV9sYW'
    '5ndWFnZRgCIAEoCVIOc291cmNlTGFuZ3VhZ2USJwoPdGFyZ2V0X2xhbmd1YWdlGAMgASgJUg50'
    'YXJnZXRMYW5ndWFnZRIdCgppc19tZWRpY2FsGAQgASgIUglpc01lZGljYWwSGAoHY29udGV4dB'
    'gFIAEoCVIHY29udGV4dA==');

@$core.Deprecated('Use translateTextResponseDescriptor instead')
const TranslateTextResponse$json = {
  '1': 'TranslateTextResponse',
  '2': [
    {'1': 'translated_text', '3': 1, '4': 1, '5': 9, '10': 'translatedText'},
    {'1': 'source_language', '3': 2, '4': 1, '5': 9, '10': 'sourceLanguage'},
    {'1': 'target_language', '3': 3, '4': 1, '5': 9, '10': 'targetLanguage'},
    {'1': 'confidence', '3': 4, '4': 1, '5': 1, '10': 'confidence'},
    {'1': 'is_medical_term', '3': 5, '4': 1, '5': 8, '10': 'isMedicalTerm'},
    {'1': 'original_text', '3': 6, '4': 1, '5': 9, '10': 'originalText'},
  ],
};

/// Descriptor for `TranslateTextResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List translateTextResponseDescriptor = $convert.base64Decode(
    'ChVUcmFuc2xhdGVUZXh0UmVzcG9uc2USJwoPdHJhbnNsYXRlZF90ZXh0GAEgASgJUg50cmFuc2'
    'xhdGVkVGV4dBInCg9zb3VyY2VfbGFuZ3VhZ2UYAiABKAlSDnNvdXJjZUxhbmd1YWdlEicKD3Rh'
    'cmdldF9sYW5ndWFnZRgDIAEoCVIOdGFyZ2V0TGFuZ3VhZ2USHgoKY29uZmlkZW5jZRgEIAEoAV'
    'IKY29uZmlkZW5jZRImCg9pc19tZWRpY2FsX3Rlcm0YBSABKAhSDWlzTWVkaWNhbFRlcm0SIwoN'
    'b3JpZ2luYWxfdGV4dBgGIAEoCVIMb3JpZ2luYWxUZXh0');

@$core.Deprecated('Use detectLanguageRequestDescriptor instead')
const DetectLanguageRequest$json = {
  '1': 'DetectLanguageRequest',
  '2': [
    {'1': 'text', '3': 1, '4': 1, '5': 9, '10': 'text'},
  ],
};

/// Descriptor for `DetectLanguageRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List detectLanguageRequestDescriptor =
    $convert.base64Decode(
        'ChVEZXRlY3RMYW5ndWFnZVJlcXVlc3QSEgoEdGV4dBgBIAEoCVIEdGV4dA==');

@$core.Deprecated('Use detectLanguageResponseDescriptor instead')
const DetectLanguageResponse$json = {
  '1': 'DetectLanguageResponse',
  '2': [
    {
      '1': 'languages',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.DetectedLanguage',
      '10': 'languages'
    },
  ],
};

/// Descriptor for `DetectLanguageResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List detectLanguageResponseDescriptor =
    $convert.base64Decode(
        'ChZEZXRlY3RMYW5ndWFnZVJlc3BvbnNlEjsKCWxhbmd1YWdlcxgBIAMoCzIdLm1hbnBhc2lrLn'
        'YxLkRldGVjdGVkTGFuZ3VhZ2VSCWxhbmd1YWdlcw==');

@$core.Deprecated('Use detectedLanguageDescriptor instead')
const DetectedLanguage$json = {
  '1': 'DetectedLanguage',
  '2': [
    {'1': 'language_code', '3': 1, '4': 1, '5': 9, '10': 'languageCode'},
    {'1': 'language_name', '3': 2, '4': 1, '5': 9, '10': 'languageName'},
    {'1': 'confidence', '3': 3, '4': 1, '5': 1, '10': 'confidence'},
  ],
};

/// Descriptor for `DetectedLanguage`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List detectedLanguageDescriptor = $convert.base64Decode(
    'ChBEZXRlY3RlZExhbmd1YWdlEiMKDWxhbmd1YWdlX2NvZGUYASABKAlSDGxhbmd1YWdlQ29kZR'
    'IjCg1sYW5ndWFnZV9uYW1lGAIgASgJUgxsYW5ndWFnZU5hbWUSHgoKY29uZmlkZW5jZRgDIAEo'
    'AVIKY29uZmlkZW5jZQ==');

@$core.Deprecated('Use listSupportedLanguagesRequestDescriptor instead')
const ListSupportedLanguagesRequest$json = {
  '1': 'ListSupportedLanguagesRequest',
};

/// Descriptor for `ListSupportedLanguagesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSupportedLanguagesRequestDescriptor =
    $convert.base64Decode('Ch1MaXN0U3VwcG9ydGVkTGFuZ3VhZ2VzUmVxdWVzdA==');

@$core.Deprecated('Use listSupportedLanguagesResponseDescriptor instead')
const ListSupportedLanguagesResponse$json = {
  '1': 'ListSupportedLanguagesResponse',
  '2': [
    {
      '1': 'languages',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.SupportedLanguage',
      '10': 'languages'
    },
  ],
};

/// Descriptor for `ListSupportedLanguagesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSupportedLanguagesResponseDescriptor =
    $convert.base64Decode(
        'Ch5MaXN0U3VwcG9ydGVkTGFuZ3VhZ2VzUmVzcG9uc2USPAoJbGFuZ3VhZ2VzGAEgAygLMh4ubW'
        'FucGFzaWsudjEuU3VwcG9ydGVkTGFuZ3VhZ2VSCWxhbmd1YWdlcw==');

@$core.Deprecated('Use supportedLanguageDescriptor instead')
const SupportedLanguage$json = {
  '1': 'SupportedLanguage',
  '2': [
    {'1': 'language_code', '3': 1, '4': 1, '5': 9, '10': 'languageCode'},
    {'1': 'language_name', '3': 2, '4': 1, '5': 9, '10': 'languageName'},
    {'1': 'native_name', '3': 3, '4': 1, '5': 9, '10': 'nativeName'},
    {'1': 'supports_medical', '3': 4, '4': 1, '5': 8, '10': 'supportsMedical'},
  ],
};

/// Descriptor for `SupportedLanguage`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List supportedLanguageDescriptor = $convert.base64Decode(
    'ChFTdXBwb3J0ZWRMYW5ndWFnZRIjCg1sYW5ndWFnZV9jb2RlGAEgASgJUgxsYW5ndWFnZUNvZG'
    'USIwoNbGFuZ3VhZ2VfbmFtZRgCIAEoCVIMbGFuZ3VhZ2VOYW1lEh8KC25hdGl2ZV9uYW1lGAMg'
    'ASgJUgpuYXRpdmVOYW1lEikKEHN1cHBvcnRzX21lZGljYWwYBCABKAhSD3N1cHBvcnRzTWVkaW'
    'NhbA==');

@$core.Deprecated('Use translateBatchRequestDescriptor instead')
const TranslateBatchRequest$json = {
  '1': 'TranslateBatchRequest',
  '2': [
    {'1': 'texts', '3': 1, '4': 3, '5': 9, '10': 'texts'},
    {'1': 'source_language', '3': 2, '4': 1, '5': 9, '10': 'sourceLanguage'},
    {'1': 'target_language', '3': 3, '4': 1, '5': 9, '10': 'targetLanguage'},
    {'1': 'is_medical', '3': 4, '4': 1, '5': 8, '10': 'isMedical'},
  ],
};

/// Descriptor for `TranslateBatchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List translateBatchRequestDescriptor = $convert.base64Decode(
    'ChVUcmFuc2xhdGVCYXRjaFJlcXVlc3QSFAoFdGV4dHMYASADKAlSBXRleHRzEicKD3NvdXJjZV'
    '9sYW5ndWFnZRgCIAEoCVIOc291cmNlTGFuZ3VhZ2USJwoPdGFyZ2V0X2xhbmd1YWdlGAMgASgJ'
    'Ug50YXJnZXRMYW5ndWFnZRIdCgppc19tZWRpY2FsGAQgASgIUglpc01lZGljYWw=');

@$core.Deprecated('Use translateBatchResponseDescriptor instead')
const TranslateBatchResponse$json = {
  '1': 'TranslateBatchResponse',
  '2': [
    {
      '1': 'translations',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.TranslateTextResponse',
      '10': 'translations'
    },
  ],
};

/// Descriptor for `TranslateBatchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List translateBatchResponseDescriptor =
    $convert.base64Decode(
        'ChZUcmFuc2xhdGVCYXRjaFJlc3BvbnNlEkYKDHRyYW5zbGF0aW9ucxgBIAMoCzIiLm1hbnBhc2'
        'lrLnYxLlRyYW5zbGF0ZVRleHRSZXNwb25zZVIMdHJhbnNsYXRpb25z');

@$core.Deprecated('Use getTranslationHistoryRequestDescriptor instead')
const GetTranslationHistoryRequest$json = {
  '1': 'GetTranslationHistoryRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'limit', '3': 2, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 3, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `GetTranslationHistoryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTranslationHistoryRequestDescriptor =
    $convert.base64Decode(
        'ChxHZXRUcmFuc2xhdGlvbkhpc3RvcnlSZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZB'
        'IUCgVsaW1pdBgCIAEoBVIFbGltaXQSFgoGb2Zmc2V0GAMgASgFUgZvZmZzZXQ=');

@$core.Deprecated('Use getTranslationHistoryResponseDescriptor instead')
const GetTranslationHistoryResponse$json = {
  '1': 'GetTranslationHistoryResponse',
  '2': [
    {
      '1': 'records',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.TranslationRecord',
      '10': 'records'
    },
    {'1': 'total_count', '3': 2, '4': 1, '5': 5, '10': 'totalCount'},
  ],
};

/// Descriptor for `GetTranslationHistoryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTranslationHistoryResponseDescriptor =
    $convert.base64Decode(
        'Ch1HZXRUcmFuc2xhdGlvbkhpc3RvcnlSZXNwb25zZRI4CgdyZWNvcmRzGAEgAygLMh4ubWFucG'
        'FzaWsudjEuVHJhbnNsYXRpb25SZWNvcmRSB3JlY29yZHMSHwoLdG90YWxfY291bnQYAiABKAVS'
        'CnRvdGFsQ291bnQ=');

@$core.Deprecated('Use translationRecordDescriptor instead')
const TranslationRecord$json = {
  '1': 'TranslationRecord',
  '2': [
    {'1': 'record_id', '3': 1, '4': 1, '5': 9, '10': 'recordId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'source_text', '3': 3, '4': 1, '5': 9, '10': 'sourceText'},
    {'1': 'translated_text', '3': 4, '4': 1, '5': 9, '10': 'translatedText'},
    {'1': 'source_language', '3': 5, '4': 1, '5': 9, '10': 'sourceLanguage'},
    {'1': 'target_language', '3': 6, '4': 1, '5': 9, '10': 'targetLanguage'},
    {'1': 'is_medical', '3': 7, '4': 1, '5': 8, '10': 'isMedical'},
    {
      '1': 'created_at',
      '3': 8,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
  ],
};

/// Descriptor for `TranslationRecord`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List translationRecordDescriptor = $convert.base64Decode(
    'ChFUcmFuc2xhdGlvblJlY29yZBIbCglyZWNvcmRfaWQYASABKAlSCHJlY29yZElkEhcKB3VzZX'
    'JfaWQYAiABKAlSBnVzZXJJZBIfCgtzb3VyY2VfdGV4dBgDIAEoCVIKc291cmNlVGV4dBInCg90'
    'cmFuc2xhdGVkX3RleHQYBCABKAlSDnRyYW5zbGF0ZWRUZXh0EicKD3NvdXJjZV9sYW5ndWFnZR'
    'gFIAEoCVIOc291cmNlTGFuZ3VhZ2USJwoPdGFyZ2V0X2xhbmd1YWdlGAYgASgJUg50YXJnZXRM'
    'YW5ndWFnZRIdCgppc19tZWRpY2FsGAcgASgIUglpc01lZGljYWwSOQoKY3JlYXRlZF9hdBgIIA'
    'EoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSCWNyZWF0ZWRBdA==');

@$core.Deprecated('Use getTranslationUsageRequestDescriptor instead')
const GetTranslationUsageRequest$json = {
  '1': 'GetTranslationUsageRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `GetTranslationUsageRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTranslationUsageRequestDescriptor =
    $convert.base64Decode(
        'ChpHZXRUcmFuc2xhdGlvblVzYWdlUmVxdWVzdBIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQ=');

@$core.Deprecated('Use getTranslationUsageResponseDescriptor instead')
const GetTranslationUsageResponse$json = {
  '1': 'GetTranslationUsageResponse',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'total_translations',
      '3': 2,
      '4': 1,
      '5': 5,
      '10': 'totalTranslations'
    },
    {
      '1': 'monthly_translations',
      '3': 3,
      '4': 1,
      '5': 5,
      '10': 'monthlyTranslations'
    },
    {'1': 'monthly_limit', '3': 4, '4': 1, '5': 5, '10': 'monthlyLimit'},
    {
      '1': 'medical_translations',
      '3': 5,
      '4': 1,
      '5': 5,
      '10': 'medicalTranslations'
    },
    {
      '1': 'by_language_pair',
      '3': 6,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.GetTranslationUsageResponse.ByLanguagePairEntry',
      '10': 'byLanguagePair'
    },
  ],
  '3': [GetTranslationUsageResponse_ByLanguagePairEntry$json],
};

@$core.Deprecated('Use getTranslationUsageResponseDescriptor instead')
const GetTranslationUsageResponse_ByLanguagePairEntry$json = {
  '1': 'ByLanguagePairEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 5, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetTranslationUsageResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTranslationUsageResponseDescriptor = $convert.base64Decode(
    'ChtHZXRUcmFuc2xhdGlvblVzYWdlUmVzcG9uc2USFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklkEi'
    '0KEnRvdGFsX3RyYW5zbGF0aW9ucxgCIAEoBVIRdG90YWxUcmFuc2xhdGlvbnMSMQoUbW9udGhs'
    'eV90cmFuc2xhdGlvbnMYAyABKAVSE21vbnRobHlUcmFuc2xhdGlvbnMSIwoNbW9udGhseV9saW'
    '1pdBgEIAEoBVIMbW9udGhseUxpbWl0EjEKFG1lZGljYWxfdHJhbnNsYXRpb25zGAUgASgFUhNt'
    'ZWRpY2FsVHJhbnNsYXRpb25zEmYKEGJ5X2xhbmd1YWdlX3BhaXIYBiADKAsyPC5tYW5wYXNpay'
    '52MS5HZXRUcmFuc2xhdGlvblVzYWdlUmVzcG9uc2UuQnlMYW5ndWFnZVBhaXJFbnRyeVIOYnlM'
    'YW5ndWFnZVBhaXIaQQoTQnlMYW5ndWFnZVBhaXJFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCg'
    'V2YWx1ZRgCIAEoBVIFdmFsdWU6AjgB');

@$core.Deprecated('Use listDoctorsByFacilityRequestDescriptor instead')
const ListDoctorsByFacilityRequest$json = {
  '1': 'ListDoctorsByFacilityRequest',
  '2': [
    {'1': 'facility_id', '3': 1, '4': 1, '5': 9, '10': 'facilityId'},
  ],
};

/// Descriptor for `ListDoctorsByFacilityRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDoctorsByFacilityRequestDescriptor =
    $convert.base64Decode(
        'ChxMaXN0RG9jdG9yc0J5RmFjaWxpdHlSZXF1ZXN0Eh8KC2ZhY2lsaXR5X2lkGAEgASgJUgpmYW'
        'NpbGl0eUlk');

@$core.Deprecated('Use listDoctorsByFacilityResponseDescriptor instead')
const ListDoctorsByFacilityResponse$json = {
  '1': 'ListDoctorsByFacilityResponse',
  '2': [
    {
      '1': 'doctors',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.Doctor',
      '10': 'doctors'
    },
  ],
};

/// Descriptor for `ListDoctorsByFacilityResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDoctorsByFacilityResponseDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0RG9jdG9yc0J5RmFjaWxpdHlSZXNwb25zZRItCgdkb2N0b3JzGAEgAygLMhMubWFucG'
        'FzaWsudjEuRG9jdG9yUgdkb2N0b3Jz');

@$core.Deprecated('Use doctorDescriptor instead')
const Doctor$json = {
  '1': 'Doctor',
  '2': [
    {'1': 'doctor_id', '3': 1, '4': 1, '5': 9, '10': 'doctorId'},
    {'1': 'name', '3': 2, '4': 1, '5': 9, '10': 'name'},
    {'1': 'specialty', '3': 3, '4': 1, '5': 9, '10': 'specialty'},
    {'1': 'facility_id', '3': 4, '4': 1, '5': 9, '10': 'facilityId'},
    {'1': 'rating', '3': 5, '4': 1, '5': 2, '10': 'rating'},
    {'1': 'consultation_fee', '3': 6, '4': 1, '5': 5, '10': 'consultationFee'},
    {
      '1': 'accepts_telemedicine',
      '3': 7,
      '4': 1,
      '5': 8,
      '10': 'acceptsTelemedicine'
    },
    {
      '1': 'available_region_codes',
      '3': 8,
      '4': 3,
      '5': 9,
      '10': 'availableRegionCodes'
    },
    {'1': 'next_available_at', '3': 9, '4': 1, '5': 9, '10': 'nextAvailableAt'},
  ],
};

/// Descriptor for `Doctor`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List doctorDescriptor = $convert.base64Decode(
    'CgZEb2N0b3ISGwoJZG9jdG9yX2lkGAEgASgJUghkb2N0b3JJZBISCgRuYW1lGAIgASgJUgRuYW'
    '1lEhwKCXNwZWNpYWx0eRgDIAEoCVIJc3BlY2lhbHR5Eh8KC2ZhY2lsaXR5X2lkGAQgASgJUgpm'
    'YWNpbGl0eUlkEhYKBnJhdGluZxgFIAEoAlIGcmF0aW5nEikKEGNvbnN1bHRhdGlvbl9mZWUYBi'
    'ABKAVSD2NvbnN1bHRhdGlvbkZlZRIxChRhY2NlcHRzX3RlbGVtZWRpY2luZRgHIAEoCFITYWNj'
    'ZXB0c1RlbGVtZWRpY2luZRI0ChZhdmFpbGFibGVfcmVnaW9uX2NvZGVzGAggAygJUhRhdmFpbG'
    'FibGVSZWdpb25Db2RlcxIqChFuZXh0X2F2YWlsYWJsZV9hdBgJIAEoCVIPbmV4dEF2YWlsYWJs'
    'ZUF0');

@$core.Deprecated('Use getDoctorAvailabilityRequestDescriptor instead')
const GetDoctorAvailabilityRequest$json = {
  '1': 'GetDoctorAvailabilityRequest',
  '2': [
    {'1': 'doctor_id', '3': 1, '4': 1, '5': 9, '10': 'doctorId'},
    {'1': 'date', '3': 2, '4': 1, '5': 9, '10': 'date'},
  ],
};

/// Descriptor for `GetDoctorAvailabilityRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDoctorAvailabilityRequestDescriptor =
    $convert.base64Decode(
        'ChxHZXREb2N0b3JBdmFpbGFiaWxpdHlSZXF1ZXN0EhsKCWRvY3Rvcl9pZBgBIAEoCVIIZG9jdG'
        '9ySWQSEgoEZGF0ZRgCIAEoCVIEZGF0ZQ==');

@$core.Deprecated('Use getDoctorAvailabilityResponseDescriptor instead')
const GetDoctorAvailabilityResponse$json = {
  '1': 'GetDoctorAvailabilityResponse',
  '2': [
    {
      '1': 'slots',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.TimeSlotDetail',
      '10': 'slots'
    },
  ],
};

/// Descriptor for `GetDoctorAvailabilityResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDoctorAvailabilityResponseDescriptor =
    $convert.base64Decode(
        'Ch1HZXREb2N0b3JBdmFpbGFiaWxpdHlSZXNwb25zZRIxCgVzbG90cxgBIAMoCzIbLm1hbnBhc2'
        'lrLnYxLlRpbWVTbG90RGV0YWlsUgVzbG90cw==');

@$core.Deprecated('Use timeSlotDetailDescriptor instead')
const TimeSlotDetail$json = {
  '1': 'TimeSlotDetail',
  '2': [
    {'1': 'start_time', '3': 1, '4': 1, '5': 9, '10': 'startTime'},
    {'1': 'end_time', '3': 2, '4': 1, '5': 9, '10': 'endTime'},
    {'1': 'is_available', '3': 3, '4': 1, '5': 8, '10': 'isAvailable'},
  ],
};

/// Descriptor for `TimeSlotDetail`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List timeSlotDetailDescriptor = $convert.base64Decode(
    'Cg5UaW1lU2xvdERldGFpbBIdCgpzdGFydF90aW1lGAEgASgJUglzdGFydFRpbWUSGQoIZW5kX3'
    'RpbWUYAiABKAlSB2VuZFRpbWUSIQoMaXNfYXZhaWxhYmxlGAMgASgIUgtpc0F2YWlsYWJsZQ==');

@$core.Deprecated('Use selectDoctorRequestDescriptor instead')
const SelectDoctorRequest$json = {
  '1': 'SelectDoctorRequest',
  '2': [
    {'1': 'facility_id', '3': 1, '4': 1, '5': 9, '10': 'facilityId'},
    {'1': 'doctor_id', '3': 2, '4': 1, '5': 9, '10': 'doctorId'},
    {'1': 'user_id', '3': 3, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `SelectDoctorRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List selectDoctorRequestDescriptor = $convert.base64Decode(
    'ChNTZWxlY3REb2N0b3JSZXF1ZXN0Eh8KC2ZhY2lsaXR5X2lkGAEgASgJUgpmYWNpbGl0eUlkEh'
    'sKCWRvY3Rvcl9pZBgCIAEoCVIIZG9jdG9ySWQSFwoHdXNlcl9pZBgDIAEoCVIGdXNlcklk');

@$core.Deprecated('Use selectDoctorResponseDescriptor instead')
const SelectDoctorResponse$json = {
  '1': 'SelectDoctorResponse',
  '2': [
    {
      '1': 'doctor',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.manpasik.v1.Doctor',
      '10': 'doctor'
    },
    {'1': 'success', '3': 2, '4': 1, '5': 8, '10': 'success'},
  ],
};

/// Descriptor for `SelectDoctorResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List selectDoctorResponseDescriptor = $convert.base64Decode(
    'ChRTZWxlY3REb2N0b3JSZXNwb25zZRIrCgZkb2N0b3IYASABKAsyEy5tYW5wYXNpay52MS5Eb2'
    'N0b3JSBmRvY3RvchIYCgdzdWNjZXNzGAIgASgIUgdzdWNjZXNz');

@$core.Deprecated('Use selectPharmacyRequestDescriptor instead')
const SelectPharmacyRequest$json = {
  '1': 'SelectPharmacyRequest',
  '2': [
    {'1': 'prescription_id', '3': 1, '4': 1, '5': 9, '10': 'prescriptionId'},
    {'1': 'pharmacy_id', '3': 2, '4': 1, '5': 9, '10': 'pharmacyId'},
    {'1': 'pharmacy_name', '3': 3, '4': 1, '5': 9, '10': 'pharmacyName'},
    {'1': 'fulfillment_type', '3': 4, '4': 1, '5': 9, '10': 'fulfillmentType'},
    {'1': 'shipping_address', '3': 5, '4': 1, '5': 9, '10': 'shippingAddress'},
  ],
};

/// Descriptor for `SelectPharmacyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List selectPharmacyRequestDescriptor = $convert.base64Decode(
    'ChVTZWxlY3RQaGFybWFjeVJlcXVlc3QSJwoPcHJlc2NyaXB0aW9uX2lkGAEgASgJUg5wcmVzY3'
    'JpcHRpb25JZBIfCgtwaGFybWFjeV9pZBgCIAEoCVIKcGhhcm1hY3lJZBIjCg1waGFybWFjeV9u'
    'YW1lGAMgASgJUgxwaGFybWFjeU5hbWUSKQoQZnVsZmlsbG1lbnRfdHlwZRgEIAEoCVIPZnVsZm'
    'lsbG1lbnRUeXBlEikKEHNoaXBwaW5nX2FkZHJlc3MYBSABKAlSD3NoaXBwaW5nQWRkcmVzcw==');

@$core.Deprecated('Use selectPharmacyResponseDescriptor instead')
const SelectPharmacyResponse$json = {
  '1': 'SelectPharmacyResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
    {'1': 'message', '3': 2, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `SelectPharmacyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List selectPharmacyResponseDescriptor =
    $convert.base64Decode(
        'ChZTZWxlY3RQaGFybWFjeVJlc3BvbnNlEhgKB3N1Y2Nlc3MYASABKAhSB3N1Y2Nlc3MSGAoHbW'
        'Vzc2FnZRgCIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use sendToPharmacyRequestDescriptor instead')
const SendToPharmacyRequest$json = {
  '1': 'SendToPharmacyRequest',
  '2': [
    {'1': 'prescription_id', '3': 1, '4': 1, '5': 9, '10': 'prescriptionId'},
  ],
};

/// Descriptor for `SendToPharmacyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendToPharmacyRequestDescriptor = $convert.base64Decode(
    'ChVTZW5kVG9QaGFybWFjeVJlcXVlc3QSJwoPcHJlc2NyaXB0aW9uX2lkGAEgASgJUg5wcmVzY3'
    'JpcHRpb25JZA==');

@$core.Deprecated('Use sendToPharmacyResponseDescriptor instead')
const SendToPharmacyResponse$json = {
  '1': 'SendToPharmacyResponse',
  '2': [
    {
      '1': 'fulfillment_token',
      '3': 1,
      '4': 1,
      '5': 9,
      '10': 'fulfillmentToken'
    },
    {'1': 'expires_at', '3': 2, '4': 1, '5': 9, '10': 'expiresAt'},
    {'1': 'success', '3': 3, '4': 1, '5': 8, '10': 'success'},
  ],
};

/// Descriptor for `SendToPharmacyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendToPharmacyResponseDescriptor = $convert.base64Decode(
    'ChZTZW5kVG9QaGFybWFjeVJlc3BvbnNlEisKEWZ1bGZpbGxtZW50X3Rva2VuGAEgASgJUhBmdW'
    'xmaWxsbWVudFRva2VuEh0KCmV4cGlyZXNfYXQYAiABKAlSCWV4cGlyZXNBdBIYCgdzdWNjZXNz'
    'GAMgASgIUgdzdWNjZXNz');

@$core.Deprecated('Use getByTokenRequestDescriptor instead')
const GetByTokenRequest$json = {
  '1': 'GetByTokenRequest',
  '2': [
    {
      '1': 'fulfillment_token',
      '3': 1,
      '4': 1,
      '5': 9,
      '10': 'fulfillmentToken'
    },
  ],
};

/// Descriptor for `GetByTokenRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getByTokenRequestDescriptor = $convert.base64Decode(
    'ChFHZXRCeVRva2VuUmVxdWVzdBIrChFmdWxmaWxsbWVudF90b2tlbhgBIAEoCVIQZnVsZmlsbG'
    '1lbnRUb2tlbg==');

@$core.Deprecated('Use updateDispensaryStatusRequestDescriptor instead')
const UpdateDispensaryStatusRequest$json = {
  '1': 'UpdateDispensaryStatusRequest',
  '2': [
    {'1': 'prescription_id', '3': 1, '4': 1, '5': 9, '10': 'prescriptionId'},
    {'1': 'status', '3': 2, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `UpdateDispensaryStatusRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDispensaryStatusRequestDescriptor =
    $convert.base64Decode(
        'Ch1VcGRhdGVEaXNwZW5zYXJ5U3RhdHVzUmVxdWVzdBInCg9wcmVzY3JpcHRpb25faWQYASABKA'
        'lSDnByZXNjcmlwdGlvbklkEhYKBnN0YXR1cxgCIAEoCVIGc3RhdHVz');

@$core.Deprecated('Use createConsentRequestDescriptor instead')
const CreateConsentRequest$json = {
  '1': 'CreateConsentRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'provider_id', '3': 2, '4': 1, '5': 9, '10': 'providerId'},
    {'1': 'provider_name', '3': 3, '4': 1, '5': 9, '10': 'providerName'},
    {'1': 'consent_type', '3': 4, '4': 1, '5': 9, '10': 'consentType'},
    {'1': 'scope', '3': 5, '4': 3, '5': 9, '10': 'scope'},
    {'1': 'purpose', '3': 6, '4': 1, '5': 9, '10': 'purpose'},
  ],
};

/// Descriptor for `CreateConsentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createConsentRequestDescriptor = $convert.base64Decode(
    'ChRDcmVhdGVDb25zZW50UmVxdWVzdBIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQSHwoLcHJvdm'
    'lkZXJfaWQYAiABKAlSCnByb3ZpZGVySWQSIwoNcHJvdmlkZXJfbmFtZRgDIAEoCVIMcHJvdmlk'
    'ZXJOYW1lEiEKDGNvbnNlbnRfdHlwZRgEIAEoCVILY29uc2VudFR5cGUSFAoFc2NvcGUYBSADKA'
    'lSBXNjb3BlEhgKB3B1cnBvc2UYBiABKAlSB3B1cnBvc2U=');

@$core.Deprecated('Use dataSharingConsentDescriptor instead')
const DataSharingConsent$json = {
  '1': 'DataSharingConsent',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'provider_id', '3': 3, '4': 1, '5': 9, '10': 'providerId'},
    {'1': 'provider_name', '3': 4, '4': 1, '5': 9, '10': 'providerName'},
    {'1': 'consent_type', '3': 5, '4': 1, '5': 9, '10': 'consentType'},
    {'1': 'scope', '3': 6, '4': 3, '5': 9, '10': 'scope'},
    {'1': 'purpose', '3': 7, '4': 1, '5': 9, '10': 'purpose'},
    {'1': 'status', '3': 8, '4': 1, '5': 9, '10': 'status'},
    {'1': 'granted_at', '3': 9, '4': 1, '5': 9, '10': 'grantedAt'},
    {'1': 'expires_at', '3': 10, '4': 1, '5': 9, '10': 'expiresAt'},
    {'1': 'revoked_at', '3': 11, '4': 1, '5': 9, '10': 'revokedAt'},
    {'1': 'revoke_reason', '3': 12, '4': 1, '5': 9, '10': 'revokeReason'},
  ],
};

/// Descriptor for `DataSharingConsent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dataSharingConsentDescriptor = $convert.base64Decode(
    'ChJEYXRhU2hhcmluZ0NvbnNlbnQSDgoCaWQYASABKAlSAmlkEhcKB3VzZXJfaWQYAiABKAlSBn'
    'VzZXJJZBIfCgtwcm92aWRlcl9pZBgDIAEoCVIKcHJvdmlkZXJJZBIjCg1wcm92aWRlcl9uYW1l'
    'GAQgASgJUgxwcm92aWRlck5hbWUSIQoMY29uc2VudF90eXBlGAUgASgJUgtjb25zZW50VHlwZR'
    'IUCgVzY29wZRgGIAMoCVIFc2NvcGUSGAoHcHVycG9zZRgHIAEoCVIHcHVycG9zZRIWCgZzdGF0'
    'dXMYCCABKAlSBnN0YXR1cxIdCgpncmFudGVkX2F0GAkgASgJUglncmFudGVkQXQSHQoKZXhwaX'
    'Jlc19hdBgKIAEoCVIJZXhwaXJlc0F0Eh0KCnJldm9rZWRfYXQYCyABKAlSCXJldm9rZWRBdBIj'
    'Cg1yZXZva2VfcmVhc29uGAwgASgJUgxyZXZva2VSZWFzb24=');

@$core.Deprecated('Use revokeConsentRequestDescriptor instead')
const RevokeConsentRequest$json = {
  '1': 'RevokeConsentRequest',
  '2': [
    {'1': 'consent_id', '3': 1, '4': 1, '5': 9, '10': 'consentId'},
    {'1': 'reason', '3': 2, '4': 1, '5': 9, '10': 'reason'},
  ],
};

/// Descriptor for `RevokeConsentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List revokeConsentRequestDescriptor = $convert.base64Decode(
    'ChRSZXZva2VDb25zZW50UmVxdWVzdBIdCgpjb25zZW50X2lkGAEgASgJUgljb25zZW50SWQSFg'
    'oGcmVhc29uGAIgASgJUgZyZWFzb24=');

@$core.Deprecated('Use revokeConsentResponseDescriptor instead')
const RevokeConsentResponse$json = {
  '1': 'RevokeConsentResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
  ],
};

/// Descriptor for `RevokeConsentResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List revokeConsentResponseDescriptor =
    $convert.base64Decode(
        'ChVSZXZva2VDb25zZW50UmVzcG9uc2USGAoHc3VjY2VzcxgBIAEoCFIHc3VjY2Vzcw==');

@$core.Deprecated('Use listConsentsRequestDescriptor instead')
const ListConsentsRequest$json = {
  '1': 'ListConsentsRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `ListConsentsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listConsentsRequestDescriptor =
    $convert.base64Decode(
        'ChNMaXN0Q29uc2VudHNSZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZA==');

@$core.Deprecated('Use listConsentsResponseDescriptor instead')
const ListConsentsResponse$json = {
  '1': 'ListConsentsResponse',
  '2': [
    {
      '1': 'consents',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.DataSharingConsent',
      '10': 'consents'
    },
  ],
};

/// Descriptor for `ListConsentsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listConsentsResponseDescriptor = $convert.base64Decode(
    'ChRMaXN0Q29uc2VudHNSZXNwb25zZRI7Cghjb25zZW50cxgBIAMoCzIfLm1hbnBhc2lrLnYxLk'
    'RhdGFTaGFyaW5nQ29uc2VudFIIY29uc2VudHM=');

@$core.Deprecated('Use shareWithProviderRequestDescriptor instead')
const ShareWithProviderRequest$json = {
  '1': 'ShareWithProviderRequest',
  '2': [
    {'1': 'consent_id', '3': 1, '4': 1, '5': 9, '10': 'consentId'},
  ],
};

/// Descriptor for `ShareWithProviderRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List shareWithProviderRequestDescriptor =
    $convert.base64Decode(
        'ChhTaGFyZVdpdGhQcm92aWRlclJlcXVlc3QSHQoKY29uc2VudF9pZBgBIAEoCVIJY29uc2VudE'
        'lk');

@$core.Deprecated('Use shareWithProviderResponseDescriptor instead')
const ShareWithProviderResponse$json = {
  '1': 'ShareWithProviderResponse',
  '2': [
    {'1': 'fhir_bundle_json', '3': 1, '4': 1, '5': 9, '10': 'fhirBundleJson'},
    {'1': 'resource_count', '3': 2, '4': 1, '5': 5, '10': 'resourceCount'},
    {'1': 'shared_at', '3': 3, '4': 1, '5': 9, '10': 'sharedAt'},
  ],
};

/// Descriptor for `ShareWithProviderResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List shareWithProviderResponseDescriptor = $convert.base64Decode(
    'ChlTaGFyZVdpdGhQcm92aWRlclJlc3BvbnNlEigKEGZoaXJfYnVuZGxlX2pzb24YASABKAlSDm'
    'ZoaXJCdW5kbGVKc29uEiUKDnJlc291cmNlX2NvdW50GAIgASgFUg1yZXNvdXJjZUNvdW50EhsK'
    'CXNoYXJlZF9hdBgDIAEoCVIIc2hhcmVkQXQ=');

@$core.Deprecated('Use getDataAccessLogRequestDescriptor instead')
const GetDataAccessLogRequest$json = {
  '1': 'GetDataAccessLogRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'limit', '3': 2, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 3, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `GetDataAccessLogRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDataAccessLogRequestDescriptor =
    $convert.base64Decode(
        'ChdHZXREYXRhQWNjZXNzTG9nUmVxdWVzdBIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQSFAoFbG'
        'ltaXQYAiABKAVSBWxpbWl0EhYKBm9mZnNldBgDIAEoBVIGb2Zmc2V0');

@$core.Deprecated('Use getDataAccessLogResponseDescriptor instead')
const GetDataAccessLogResponse$json = {
  '1': 'GetDataAccessLogResponse',
  '2': [
    {
      '1': 'entries',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.DataAccessLogEntry',
      '10': 'entries'
    },
    {'1': 'total', '3': 2, '4': 1, '5': 5, '10': 'total'},
  ],
};

/// Descriptor for `GetDataAccessLogResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDataAccessLogResponseDescriptor =
    $convert.base64Decode(
        'ChhHZXREYXRhQWNjZXNzTG9nUmVzcG9uc2USOQoHZW50cmllcxgBIAMoCzIfLm1hbnBhc2lrLn'
        'YxLkRhdGFBY2Nlc3NMb2dFbnRyeVIHZW50cmllcxIUCgV0b3RhbBgCIAEoBVIFdG90YWw=');

@$core.Deprecated('Use dataAccessLogEntryDescriptor instead')
const DataAccessLogEntry$json = {
  '1': 'DataAccessLogEntry',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'consent_id', '3': 2, '4': 1, '5': 9, '10': 'consentId'},
    {'1': 'user_id', '3': 3, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'provider_id', '3': 4, '4': 1, '5': 9, '10': 'providerId'},
    {'1': 'action', '3': 5, '4': 1, '5': 9, '10': 'action'},
    {'1': 'resource_type', '3': 6, '4': 1, '5': 9, '10': 'resourceType'},
    {'1': 'resource_ids', '3': 7, '4': 3, '5': 9, '10': 'resourceIds'},
    {'1': 'accessed_at', '3': 8, '4': 1, '5': 9, '10': 'accessedAt'},
  ],
};

/// Descriptor for `DataAccessLogEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dataAccessLogEntryDescriptor = $convert.base64Decode(
    'ChJEYXRhQWNjZXNzTG9nRW50cnkSDgoCaWQYASABKAlSAmlkEh0KCmNvbnNlbnRfaWQYAiABKA'
    'lSCWNvbnNlbnRJZBIXCgd1c2VyX2lkGAMgASgJUgZ1c2VySWQSHwoLcHJvdmlkZXJfaWQYBCAB'
    'KAlSCnByb3ZpZGVySWQSFgoGYWN0aW9uGAUgASgJUgZhY3Rpb24SIwoNcmVzb3VyY2VfdHlwZR'
    'gGIAEoCVIMcmVzb3VyY2VUeXBlEiEKDHJlc291cmNlX2lkcxgHIAMoCVILcmVzb3VyY2VJZHMS'
    'HwoLYWNjZXNzZWRfYXQYCCABKAlSCmFjY2Vzc2VkQXQ=');

@$core.Deprecated('Use exportSingleMeasurementRequestDescriptor instead')
const ExportSingleMeasurementRequest$json = {
  '1': 'ExportSingleMeasurementRequest',
  '2': [
    {'1': 'session_id', '3': 1, '4': 1, '5': 9, '10': 'sessionId'},
  ],
};

/// Descriptor for `ExportSingleMeasurementRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List exportSingleMeasurementRequestDescriptor =
    $convert.base64Decode(
        'Ch5FeHBvcnRTaW5nbGVNZWFzdXJlbWVudFJlcXVlc3QSHQoKc2Vzc2lvbl9pZBgBIAEoCVIJc2'
        'Vzc2lvbklk');

@$core.Deprecated('Use exportToFHIRObservationsRequestDescriptor instead')
const ExportToFHIRObservationsRequest$json = {
  '1': 'ExportToFHIRObservationsRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'from_date', '3': 2, '4': 1, '5': 9, '10': 'fromDate'},
    {'1': 'to_date', '3': 3, '4': 1, '5': 9, '10': 'toDate'},
  ],
};

/// Descriptor for `ExportToFHIRObservationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List exportToFHIRObservationsRequestDescriptor =
    $convert.base64Decode(
        'Ch9FeHBvcnRUb0ZISVJPYnNlcnZhdGlvbnNSZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZX'
        'JJZBIbCglmcm9tX2RhdGUYAiABKAlSCGZyb21EYXRlEhcKB3RvX2RhdGUYAyABKAlSBnRvRGF0'
        'ZQ==');

@$core.Deprecated('Use exportFHIRResponseDescriptor instead')
const ExportFHIRResponse$json = {
  '1': 'ExportFHIRResponse',
  '2': [
    {'1': 'fhir_bundle_json', '3': 1, '4': 1, '5': 9, '10': 'fhirBundleJson'},
    {'1': 'resource_count', '3': 2, '4': 1, '5': 5, '10': 'resourceCount'},
  ],
};

/// Descriptor for `ExportFHIRResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List exportFHIRResponseDescriptor = $convert.base64Decode(
    'ChJFeHBvcnRGSElSUmVzcG9uc2USKAoQZmhpcl9idW5kbGVfanNvbhgBIAEoCVIOZmhpckJ1bm'
    'RsZUpzb24SJQoOcmVzb3VyY2VfY291bnQYAiABKAVSDXJlc291cmNlQ291bnQ=');

@$core.Deprecated('Use syncDigitalTwinRequestDescriptor instead')
const SyncDigitalTwinRequest$json = {
  '1': 'SyncDigitalTwinRequest',
  '2': [
    {'1': 'session_id', '3': 1, '4': 1, '5': 9, '10': 'sessionId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'device_id', '3': 3, '4': 1, '5': 9, '10': 'deviceId'},
    {'1': 'residuals', '3': 4, '4': 3, '5': 1, '10': 'residuals'},
    {'1': 'ewma_value', '3': 5, '4': 1, '5': 1, '10': 'ewmaValue'},
    {'1': 'cusum_pos', '3': 6, '4': 1, '5': 1, '10': 'cusumPos'},
    {'1': 'cusum_neg', '3': 7, '4': 1, '5': 1, '10': 'cusumNeg'},
    {'1': 'health_state', '3': 8, '4': 1, '5': 9, '10': 'healthState'},
    {'1': 'drift_score', '3': 9, '4': 1, '5': 1, '10': 'driftScore'},
    {
      '1': 'remaining_measurements',
      '3': 10,
      '4': 1,
      '5': 5,
      '10': 'remainingMeasurements'
    },
    {
      '1': 'fingerprint_vector',
      '3': 11,
      '4': 3,
      '5': 2,
      '10': 'fingerprintVector'
    },
    {'1': 'fingerprint_dim', '3': 12, '4': 1, '5': 5, '10': 'fingerprintDim'},
  ],
};

/// Descriptor for `SyncDigitalTwinRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List syncDigitalTwinRequestDescriptor = $convert.base64Decode(
    'ChZTeW5jRGlnaXRhbFR3aW5SZXF1ZXN0Eh0KCnNlc3Npb25faWQYASABKAlSCXNlc3Npb25JZB'
    'IXCgd1c2VyX2lkGAIgASgJUgZ1c2VySWQSGwoJZGV2aWNlX2lkGAMgASgJUghkZXZpY2VJZBIc'
    'CglyZXNpZHVhbHMYBCADKAFSCXJlc2lkdWFscxIdCgpld21hX3ZhbHVlGAUgASgBUglld21hVm'
    'FsdWUSGwoJY3VzdW1fcG9zGAYgASgBUghjdXN1bVBvcxIbCgljdXN1bV9uZWcYByABKAFSCGN1'
    'c3VtTmVnEiEKDGhlYWx0aF9zdGF0ZRgIIAEoCVILaGVhbHRoU3RhdGUSHwoLZHJpZnRfc2Nvcm'
    'UYCSABKAFSCmRyaWZ0U2NvcmUSNQoWcmVtYWluaW5nX21lYXN1cmVtZW50cxgKIAEoBVIVcmVt'
    'YWluaW5nTWVhc3VyZW1lbnRzEi0KEmZpbmdlcnByaW50X3ZlY3RvchgLIAMoAlIRZmluZ2VycH'
    'JpbnRWZWN0b3ISJwoPZmluZ2VycHJpbnRfZGltGAwgASgFUg5maW5nZXJwcmludERpbQ==');

@$core.Deprecated('Use syncDigitalTwinResponseDescriptor instead')
const SyncDigitalTwinResponse$json = {
  '1': 'SyncDigitalTwinResponse',
  '2': [
    {'1': 'accepted', '3': 1, '4': 1, '5': 8, '10': 'accepted'},
    {'1': 'twin_id', '3': 2, '4': 1, '5': 9, '10': 'twinId'},
    {
      '1': 'recommended_action',
      '3': 3,
      '4': 1,
      '5': 9,
      '10': 'recommendedAction'
    },
    {
      '1': 'next_calibration_drift',
      '3': 4,
      '4': 1,
      '5': 1,
      '10': 'nextCalibrationDrift'
    },
    {
      '1': 'synced_at',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'syncedAt'
    },
  ],
};

/// Descriptor for `SyncDigitalTwinResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List syncDigitalTwinResponseDescriptor = $convert.base64Decode(
    'ChdTeW5jRGlnaXRhbFR3aW5SZXNwb25zZRIaCghhY2NlcHRlZBgBIAEoCFIIYWNjZXB0ZWQSFw'
    'oHdHdpbl9pZBgCIAEoCVIGdHdpbklkEi0KEnJlY29tbWVuZGVkX2FjdGlvbhgDIAEoCVIRcmVj'
    'b21tZW5kZWRBY3Rpb24SNAoWbmV4dF9jYWxpYnJhdGlvbl9kcmlmdBgEIAEoAVIUbmV4dENhbG'
    'licmF0aW9uRHJpZnQSNwoJc3luY2VkX2F0GAUgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVz'
    'dGFtcFIIc3luY2VkQXQ=');

@$core.Deprecated('Use getCalibrationStatusRequestDescriptor instead')
const GetCalibrationStatusRequest$json = {
  '1': 'GetCalibrationStatusRequest',
  '2': [
    {'1': 'session_id', '3': 1, '4': 1, '5': 9, '10': 'sessionId'},
    {'1': 'device_id', '3': 2, '4': 1, '5': 9, '10': 'deviceId'},
  ],
};

/// Descriptor for `GetCalibrationStatusRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCalibrationStatusRequestDescriptor =
    $convert.base64Decode(
        'ChtHZXRDYWxpYnJhdGlvblN0YXR1c1JlcXVlc3QSHQoKc2Vzc2lvbl9pZBgBIAEoCVIJc2Vzc2'
        'lvbklkEhsKCWRldmljZV9pZBgCIAEoCVIIZGV2aWNlSWQ=');

@$core.Deprecated('Use getCalibrationStatusResponseDescriptor instead')
const GetCalibrationStatusResponse$json = {
  '1': 'GetCalibrationStatusResponse',
  '2': [
    {
      '1': 'calibration_state',
      '3': 1,
      '4': 1,
      '5': 9,
      '10': 'calibrationState'
    },
    {'1': 'drift_score', '3': 2, '4': 1, '5': 1, '10': 'driftScore'},
    {
      '1': 'measurements_since_cal',
      '3': 3,
      '4': 1,
      '5': 5,
      '10': 'measurementsSinceCal'
    },
    {
      '1': 'max_measurements_before_cal',
      '3': 4,
      '4': 1,
      '5': 5,
      '10': 'maxMeasurementsBeforeCal'
    },
    {
      '1': 'last_calibrated_at',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'lastCalibratedAt'
    },
    {
      '1': 'next_calibration_at',
      '3': 6,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'nextCalibrationAt'
    },
  ],
};

/// Descriptor for `GetCalibrationStatusResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCalibrationStatusResponseDescriptor = $convert.base64Decode(
    'ChxHZXRDYWxpYnJhdGlvblN0YXR1c1Jlc3BvbnNlEisKEWNhbGlicmF0aW9uX3N0YXRlGAEgAS'
    'gJUhBjYWxpYnJhdGlvblN0YXRlEh8KC2RyaWZ0X3Njb3JlGAIgASgBUgpkcmlmdFNjb3JlEjQK'
    'Fm1lYXN1cmVtZW50c19zaW5jZV9jYWwYAyABKAVSFG1lYXN1cmVtZW50c1NpbmNlQ2FsEj0KG2'
    '1heF9tZWFzdXJlbWVudHNfYmVmb3JlX2NhbBgEIAEoBVIYbWF4TWVhc3VyZW1lbnRzQmVmb3Jl'
    'Q2FsEkgKEmxhc3RfY2FsaWJyYXRlZF9hdBgFIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3'
    'RhbXBSEGxhc3RDYWxpYnJhdGVkQXQSSgoTbmV4dF9jYWxpYnJhdGlvbl9hdBgGIAEoCzIaLmdv'
    'b2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSEW5leHRDYWxpYnJhdGlvbkF0');

@$core.Deprecated('Use updateDeviceStatusRequestDescriptor instead')
const UpdateDeviceStatusRequest$json = {
  '1': 'UpdateDeviceStatusRequest',
  '2': [
    {'1': 'device_id', '3': 1, '4': 1, '5': 9, '10': 'deviceId'},
    {'1': 'status', '3': 2, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `UpdateDeviceStatusRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDeviceStatusRequestDescriptor =
    $convert.base64Decode(
        'ChlVcGRhdGVEZXZpY2VTdGF0dXNSZXF1ZXN0EhsKCWRldmljZV9pZBgBIAEoCVIIZGV2aWNlSW'
        'QSFgoGc3RhdHVzGAIgASgJUgZzdGF0dXM=');

@$core.Deprecated('Use updateDeviceStatusResponseDescriptor instead')
const UpdateDeviceStatusResponse$json = {
  '1': 'UpdateDeviceStatusResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
    {'1': 'message', '3': 2, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UpdateDeviceStatusResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateDeviceStatusResponseDescriptor =
    $convert.base64Decode(
        'ChpVcGRhdGVEZXZpY2VTdGF0dXNSZXNwb25zZRIYCgdzdWNjZXNzGAEgASgIUgdzdWNjZXNzEh'
        'gKB21lc3NhZ2UYAiABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use listAdminsByRegionRequestDescriptor instead')
const ListAdminsByRegionRequest$json = {
  '1': 'ListAdminsByRegionRequest',
  '2': [
    {'1': 'country_code', '3': 1, '4': 1, '5': 9, '10': 'countryCode'},
    {'1': 'region_code', '3': 2, '4': 1, '5': 9, '10': 'regionCode'},
  ],
};

/// Descriptor for `ListAdminsByRegionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAdminsByRegionRequestDescriptor =
    $convert.base64Decode(
        'ChlMaXN0QWRtaW5zQnlSZWdpb25SZXF1ZXN0EiEKDGNvdW50cnlfY29kZRgBIAEoCVILY291bn'
        'RyeUNvZGUSHwoLcmVnaW9uX2NvZGUYAiABKAlSCnJlZ2lvbkNvZGU=');

@$core.Deprecated('Use listSystemConfigsRequestDescriptor instead')
const ListSystemConfigsRequest$json = {
  '1': 'ListSystemConfigsRequest',
  '2': [
    {'1': 'language_code', '3': 1, '4': 1, '5': 9, '10': 'languageCode'},
    {'1': 'category', '3': 2, '4': 1, '5': 9, '10': 'category'},
    {'1': 'include_secrets', '3': 3, '4': 1, '5': 8, '10': 'includeSecrets'},
  ],
};

/// Descriptor for `ListSystemConfigsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSystemConfigsRequestDescriptor = $convert.base64Decode(
    'ChhMaXN0U3lzdGVtQ29uZmlnc1JlcXVlc3QSIwoNbGFuZ3VhZ2VfY29kZRgBIAEoCVIMbGFuZ3'
    'VhZ2VDb2RlEhoKCGNhdGVnb3J5GAIgASgJUghjYXRlZ29yeRInCg9pbmNsdWRlX3NlY3JldHMY'
    'AyABKAhSDmluY2x1ZGVTZWNyZXRz');

@$core.Deprecated('Use listSystemConfigsResponseDescriptor instead')
const ListSystemConfigsResponse$json = {
  '1': 'ListSystemConfigsResponse',
  '2': [
    {
      '1': 'configs',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.ConfigWithMeta',
      '10': 'configs'
    },
    {
      '1': 'category_counts',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.ListSystemConfigsResponse.CategoryCountsEntry',
      '10': 'categoryCounts'
    },
  ],
  '3': [ListSystemConfigsResponse_CategoryCountsEntry$json],
};

@$core.Deprecated('Use listSystemConfigsResponseDescriptor instead')
const ListSystemConfigsResponse_CategoryCountsEntry$json = {
  '1': 'CategoryCountsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 5, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ListSystemConfigsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSystemConfigsResponseDescriptor = $convert.base64Decode(
    'ChlMaXN0U3lzdGVtQ29uZmlnc1Jlc3BvbnNlEjUKB2NvbmZpZ3MYASADKAsyGy5tYW5wYXNpay'
    '52MS5Db25maWdXaXRoTWV0YVIHY29uZmlncxJjCg9jYXRlZ29yeV9jb3VudHMYAiADKAsyOi5t'
    'YW5wYXNpay52MS5MaXN0U3lzdGVtQ29uZmlnc1Jlc3BvbnNlLkNhdGVnb3J5Q291bnRzRW50cn'
    'lSDmNhdGVnb3J5Q291bnRzGkEKE0NhdGVnb3J5Q291bnRzRW50cnkSEAoDa2V5GAEgASgJUgNr'
    'ZXkSFAoFdmFsdWUYAiABKAVSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use configWithMetaDescriptor instead')
const ConfigWithMeta$json = {
  '1': 'ConfigWithMeta',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
    {'1': 'raw_value', '3': 3, '4': 1, '5': 9, '10': 'rawValue'},
    {'1': 'category', '3': 4, '4': 1, '5': 9, '10': 'category'},
    {'1': 'value_type', '3': 5, '4': 1, '5': 9, '10': 'valueType'},
    {'1': 'security_level', '3': 6, '4': 1, '5': 9, '10': 'securityLevel'},
    {'1': 'is_required', '3': 7, '4': 1, '5': 8, '10': 'isRequired'},
    {'1': 'default_value', '3': 8, '4': 1, '5': 9, '10': 'defaultValue'},
    {'1': 'allowed_values', '3': 9, '4': 3, '5': 9, '10': 'allowedValues'},
    {'1': 'validation_regex', '3': 10, '4': 1, '5': 9, '10': 'validationRegex'},
    {'1': 'validation_min', '3': 11, '4': 1, '5': 1, '10': 'validationMin'},
    {'1': 'validation_max', '3': 12, '4': 1, '5': 1, '10': 'validationMax'},
    {'1': 'depends_on', '3': 13, '4': 1, '5': 9, '10': 'dependsOn'},
    {'1': 'depends_value', '3': 14, '4': 1, '5': 9, '10': 'dependsValue'},
    {'1': 'env_var_name', '3': 15, '4': 1, '5': 9, '10': 'envVarName'},
    {'1': 'service_name', '3': 16, '4': 1, '5': 9, '10': 'serviceName'},
    {'1': 'restart_required', '3': 17, '4': 1, '5': 8, '10': 'restartRequired'},
    {'1': 'display_name', '3': 20, '4': 1, '5': 9, '10': 'displayName'},
    {'1': 'description', '3': 21, '4': 1, '5': 9, '10': 'description'},
    {'1': 'placeholder', '3': 22, '4': 1, '5': 9, '10': 'placeholder'},
    {'1': 'help_text', '3': 23, '4': 1, '5': 9, '10': 'helpText'},
    {
      '1': 'validation_message',
      '3': 24,
      '4': 1,
      '5': 9,
      '10': 'validationMessage'
    },
    {'1': 'updated_by', '3': 30, '4': 1, '5': 9, '10': 'updatedBy'},
    {
      '1': 'updated_at',
      '3': 31,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'updatedAt'
    },
  ],
};

/// Descriptor for `ConfigWithMeta`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List configWithMetaDescriptor = $convert.base64Decode(
    'Cg5Db25maWdXaXRoTWV0YRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdW'
    'USGwoJcmF3X3ZhbHVlGAMgASgJUghyYXdWYWx1ZRIaCghjYXRlZ29yeRgEIAEoCVIIY2F0ZWdv'
    'cnkSHQoKdmFsdWVfdHlwZRgFIAEoCVIJdmFsdWVUeXBlEiUKDnNlY3VyaXR5X2xldmVsGAYgAS'
    'gJUg1zZWN1cml0eUxldmVsEh8KC2lzX3JlcXVpcmVkGAcgASgIUgppc1JlcXVpcmVkEiMKDWRl'
    'ZmF1bHRfdmFsdWUYCCABKAlSDGRlZmF1bHRWYWx1ZRIlCg5hbGxvd2VkX3ZhbHVlcxgJIAMoCV'
    'INYWxsb3dlZFZhbHVlcxIpChB2YWxpZGF0aW9uX3JlZ2V4GAogASgJUg92YWxpZGF0aW9uUmVn'
    'ZXgSJQoOdmFsaWRhdGlvbl9taW4YCyABKAFSDXZhbGlkYXRpb25NaW4SJQoOdmFsaWRhdGlvbl'
    '9tYXgYDCABKAFSDXZhbGlkYXRpb25NYXgSHQoKZGVwZW5kc19vbhgNIAEoCVIJZGVwZW5kc09u'
    'EiMKDWRlcGVuZHNfdmFsdWUYDiABKAlSDGRlcGVuZHNWYWx1ZRIgCgxlbnZfdmFyX25hbWUYDy'
    'ABKAlSCmVudlZhck5hbWUSIQoMc2VydmljZV9uYW1lGBAgASgJUgtzZXJ2aWNlTmFtZRIpChBy'
    'ZXN0YXJ0X3JlcXVpcmVkGBEgASgIUg9yZXN0YXJ0UmVxdWlyZWQSIQoMZGlzcGxheV9uYW1lGB'
    'QgASgJUgtkaXNwbGF5TmFtZRIgCgtkZXNjcmlwdGlvbhgVIAEoCVILZGVzY3JpcHRpb24SIAoL'
    'cGxhY2Vob2xkZXIYFiABKAlSC3BsYWNlaG9sZGVyEhsKCWhlbHBfdGV4dBgXIAEoCVIIaGVscF'
    'RleHQSLQoSdmFsaWRhdGlvbl9tZXNzYWdlGBggASgJUhF2YWxpZGF0aW9uTWVzc2FnZRIdCgp1'
    'cGRhdGVkX2J5GB4gASgJUgl1cGRhdGVkQnkSOQoKdXBkYXRlZF9hdBgfIAEoCzIaLmdvb2dsZS'
    '5wcm90b2J1Zi5UaW1lc3RhbXBSCXVwZGF0ZWRBdA==');

@$core.Deprecated('Use getConfigWithMetaRequestDescriptor instead')
const GetConfigWithMetaRequest$json = {
  '1': 'GetConfigWithMetaRequest',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'language_code', '3': 2, '4': 1, '5': 9, '10': 'languageCode'},
  ],
};

/// Descriptor for `GetConfigWithMetaRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getConfigWithMetaRequestDescriptor =
    $convert.base64Decode(
        'ChhHZXRDb25maWdXaXRoTWV0YVJlcXVlc3QSEAoDa2V5GAEgASgJUgNrZXkSIwoNbGFuZ3VhZ2'
        'VfY29kZRgCIAEoCVIMbGFuZ3VhZ2VDb2Rl');

@$core.Deprecated('Use validateConfigValueRequestDescriptor instead')
const ValidateConfigValueRequest$json = {
  '1': 'ValidateConfigValueRequest',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `ValidateConfigValueRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validateConfigValueRequestDescriptor =
    $convert.base64Decode(
        'ChpWYWxpZGF0ZUNvbmZpZ1ZhbHVlUmVxdWVzdBIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZR'
        'gCIAEoCVIFdmFsdWU=');

@$core.Deprecated('Use validateConfigValueResponseDescriptor instead')
const ValidateConfigValueResponse$json = {
  '1': 'ValidateConfigValueResponse',
  '2': [
    {'1': 'valid', '3': 1, '4': 1, '5': 8, '10': 'valid'},
    {'1': 'error_message', '3': 2, '4': 1, '5': 9, '10': 'errorMessage'},
    {'1': 'suggestions', '3': 3, '4': 3, '5': 9, '10': 'suggestions'},
  ],
};

/// Descriptor for `ValidateConfigValueResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validateConfigValueResponseDescriptor =
    $convert.base64Decode(
        'ChtWYWxpZGF0ZUNvbmZpZ1ZhbHVlUmVzcG9uc2USFAoFdmFsaWQYASABKAhSBXZhbGlkEiMKDW'
        'Vycm9yX21lc3NhZ2UYAiABKAlSDGVycm9yTWVzc2FnZRIgCgtzdWdnZXN0aW9ucxgDIAMoCVIL'
        'c3VnZ2VzdGlvbnM=');

@$core.Deprecated('Use bulkSetConfigsRequestDescriptor instead')
const BulkSetConfigsRequest$json = {
  '1': 'BulkSetConfigsRequest',
  '2': [
    {
      '1': 'configs',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.SetSystemConfigRequest',
      '10': 'configs'
    },
    {'1': 'reason', '3': 2, '4': 1, '5': 9, '10': 'reason'},
  ],
};

/// Descriptor for `BulkSetConfigsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List bulkSetConfigsRequestDescriptor = $convert.base64Decode(
    'ChVCdWxrU2V0Q29uZmlnc1JlcXVlc3QSPQoHY29uZmlncxgBIAMoCzIjLm1hbnBhc2lrLnYxLl'
    'NldFN5c3RlbUNvbmZpZ1JlcXVlc3RSB2NvbmZpZ3MSFgoGcmVhc29uGAIgASgJUgZyZWFzb24=');

@$core.Deprecated('Use bulkSetConfigsResponseDescriptor instead')
const BulkSetConfigsResponse$json = {
  '1': 'BulkSetConfigsResponse',
  '2': [
    {
      '1': 'results',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.ConfigChangeResult',
      '10': 'results'
    },
    {'1': 'success_count', '3': 2, '4': 1, '5': 5, '10': 'successCount'},
    {'1': 'failure_count', '3': 3, '4': 1, '5': 5, '10': 'failureCount'},
  ],
};

/// Descriptor for `BulkSetConfigsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List bulkSetConfigsResponseDescriptor = $convert.base64Decode(
    'ChZCdWxrU2V0Q29uZmlnc1Jlc3BvbnNlEjkKB3Jlc3VsdHMYASADKAsyHy5tYW5wYXNpay52MS'
    '5Db25maWdDaGFuZ2VSZXN1bHRSB3Jlc3VsdHMSIwoNc3VjY2Vzc19jb3VudBgCIAEoBVIMc3Vj'
    'Y2Vzc0NvdW50EiMKDWZhaWx1cmVfY291bnQYAyABKAVSDGZhaWx1cmVDb3VudA==');

@$core.Deprecated('Use configChangeResultDescriptor instead')
const ConfigChangeResult$json = {
  '1': 'ConfigChangeResult',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'success', '3': 2, '4': 1, '5': 8, '10': 'success'},
    {'1': 'error_message', '3': 3, '4': 1, '5': 9, '10': 'errorMessage'},
  ],
};

/// Descriptor for `ConfigChangeResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List configChangeResultDescriptor = $convert.base64Decode(
    'ChJDb25maWdDaGFuZ2VSZXN1bHQSEAoDa2V5GAEgASgJUgNrZXkSGAoHc3VjY2VzcxgCIAEoCF'
    'IHc3VjY2VzcxIjCg1lcnJvcl9tZXNzYWdlGAMgASgJUgxlcnJvck1lc3NhZ2U=');

@$core.Deprecated('Use consultationDescriptor instead')
const Consultation$json = {
  '1': 'Consultation',
  '2': [
    {'1': 'consultation_id', '3': 1, '4': 1, '5': 9, '10': 'consultationId'},
    {'1': 'patient_user_id', '3': 2, '4': 1, '5': 9, '10': 'patientUserId'},
    {'1': 'doctor_id', '3': 3, '4': 1, '5': 9, '10': 'doctorId'},
    {
      '1': 'specialty',
      '3': 4,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.DoctorSpecialty',
      '10': 'specialty'
    },
    {'1': 'chief_complaint', '3': 5, '4': 1, '5': 9, '10': 'chiefComplaint'},
    {'1': 'description', '3': 6, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'status',
      '3': 7,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.ConsultationStatus',
      '10': 'status'
    },
    {'1': 'diagnosis', '3': 8, '4': 1, '5': 9, '10': 'diagnosis'},
    {'1': 'doctor_notes', '3': 9, '4': 1, '5': 9, '10': 'doctorNotes'},
    {'1': 'prescription_id', '3': 10, '4': 1, '5': 9, '10': 'prescriptionId'},
    {'1': 'duration_minutes', '3': 11, '4': 1, '5': 5, '10': 'durationMinutes'},
    {'1': 'rating', '3': 12, '4': 1, '5': 1, '10': 'rating'},
    {
      '1': 'created_at',
      '3': 13,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
    {
      '1': 'scheduled_at',
      '3': 14,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'scheduledAt'
    },
    {
      '1': 'started_at',
      '3': 15,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'startedAt'
    },
    {
      '1': 'ended_at',
      '3': 16,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'endedAt'
    },
  ],
};

/// Descriptor for `Consultation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List consultationDescriptor = $convert.base64Decode(
    'CgxDb25zdWx0YXRpb24SJwoPY29uc3VsdGF0aW9uX2lkGAEgASgJUg5jb25zdWx0YXRpb25JZB'
    'ImCg9wYXRpZW50X3VzZXJfaWQYAiABKAlSDXBhdGllbnRVc2VySWQSGwoJZG9jdG9yX2lkGAMg'
    'ASgJUghkb2N0b3JJZBI6CglzcGVjaWFsdHkYBCABKA4yHC5tYW5wYXNpay52MS5Eb2N0b3JTcG'
    'VjaWFsdHlSCXNwZWNpYWx0eRInCg9jaGllZl9jb21wbGFpbnQYBSABKAlSDmNoaWVmQ29tcGxh'
    'aW50EiAKC2Rlc2NyaXB0aW9uGAYgASgJUgtkZXNjcmlwdGlvbhI3CgZzdGF0dXMYByABKA4yHy'
    '5tYW5wYXNpay52MS5Db25zdWx0YXRpb25TdGF0dXNSBnN0YXR1cxIcCglkaWFnbm9zaXMYCCAB'
    'KAlSCWRpYWdub3NpcxIhCgxkb2N0b3Jfbm90ZXMYCSABKAlSC2RvY3Rvck5vdGVzEicKD3ByZX'
    'NjcmlwdGlvbl9pZBgKIAEoCVIOcHJlc2NyaXB0aW9uSWQSKQoQZHVyYXRpb25fbWludXRlcxgL'
    'IAEoBVIPZHVyYXRpb25NaW51dGVzEhYKBnJhdGluZxgMIAEoAVIGcmF0aW5nEjkKCmNyZWF0ZW'
    'RfYXQYDSABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUgljcmVhdGVkQXQSPQoMc2No'
    'ZWR1bGVkX2F0GA4gASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFILc2NoZWR1bGVkQX'
    'QSOQoKc3RhcnRlZF9hdBgPIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSCXN0YXJ0'
    'ZWRBdBI1CghlbmRlZF9hdBgQIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSB2VuZG'
    'VkQXQ=');

@$core.Deprecated('Use createConsultationRequestDescriptor instead')
const CreateConsultationRequest$json = {
  '1': 'CreateConsultationRequest',
  '2': [
    {'1': 'patient_user_id', '3': 1, '4': 1, '5': 9, '10': 'patientUserId'},
    {
      '1': 'specialty',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.DoctorSpecialty',
      '10': 'specialty'
    },
    {'1': 'chief_complaint', '3': 3, '4': 1, '5': 9, '10': 'chiefComplaint'},
    {'1': 'description', '3': 4, '4': 1, '5': 9, '10': 'description'},
  ],
};

/// Descriptor for `CreateConsultationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createConsultationRequestDescriptor = $convert.base64Decode(
    'ChlDcmVhdGVDb25zdWx0YXRpb25SZXF1ZXN0EiYKD3BhdGllbnRfdXNlcl9pZBgBIAEoCVINcG'
    'F0aWVudFVzZXJJZBI6CglzcGVjaWFsdHkYAiABKA4yHC5tYW5wYXNpay52MS5Eb2N0b3JTcGVj'
    'aWFsdHlSCXNwZWNpYWx0eRInCg9jaGllZl9jb21wbGFpbnQYAyABKAlSDmNoaWVmQ29tcGxhaW'
    '50EiAKC2Rlc2NyaXB0aW9uGAQgASgJUgtkZXNjcmlwdGlvbg==');

@$core.Deprecated('Use getConsultationRequestDescriptor instead')
const GetConsultationRequest$json = {
  '1': 'GetConsultationRequest',
  '2': [
    {'1': 'consultation_id', '3': 1, '4': 1, '5': 9, '10': 'consultationId'},
  ],
};

/// Descriptor for `GetConsultationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getConsultationRequestDescriptor =
    $convert.base64Decode(
        'ChZHZXRDb25zdWx0YXRpb25SZXF1ZXN0EicKD2NvbnN1bHRhdGlvbl9pZBgBIAEoCVIOY29uc3'
        'VsdGF0aW9uSWQ=');

@$core.Deprecated('Use listConsultationsRequestDescriptor instead')
const ListConsultationsRequest$json = {
  '1': 'ListConsultationsRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'status_filter',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.ConsultationStatus',
      '10': 'statusFilter'
    },
    {'1': 'limit', '3': 3, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 4, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `ListConsultationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listConsultationsRequestDescriptor = $convert.base64Decode(
    'ChhMaXN0Q29uc3VsdGF0aW9uc1JlcXVlc3QSFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklkEkQKDX'
    'N0YXR1c19maWx0ZXIYAiABKA4yHy5tYW5wYXNpay52MS5Db25zdWx0YXRpb25TdGF0dXNSDHN0'
    'YXR1c0ZpbHRlchIUCgVsaW1pdBgDIAEoBVIFbGltaXQSFgoGb2Zmc2V0GAQgASgFUgZvZmZzZX'
    'Q=');

@$core.Deprecated('Use listConsultationsResponseDescriptor instead')
const ListConsultationsResponse$json = {
  '1': 'ListConsultationsResponse',
  '2': [
    {
      '1': 'consultations',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.Consultation',
      '10': 'consultations'
    },
    {'1': 'total_count', '3': 2, '4': 1, '5': 5, '10': 'totalCount'},
  ],
};

/// Descriptor for `ListConsultationsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listConsultationsResponseDescriptor = $convert.base64Decode(
    'ChlMaXN0Q29uc3VsdGF0aW9uc1Jlc3BvbnNlEj8KDWNvbnN1bHRhdGlvbnMYASADKAsyGS5tYW'
    '5wYXNpay52MS5Db25zdWx0YXRpb25SDWNvbnN1bHRhdGlvbnMSHwoLdG90YWxfY291bnQYAiAB'
    'KAVSCnRvdGFsQ291bnQ=');

@$core.Deprecated('Use matchDoctorRequestDescriptor instead')
const MatchDoctorRequest$json = {
  '1': 'MatchDoctorRequest',
  '2': [
    {
      '1': 'specialty',
      '3': 1,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.DoctorSpecialty',
      '10': 'specialty'
    },
    {'1': 'language', '3': 2, '4': 1, '5': 9, '10': 'language'},
  ],
};

/// Descriptor for `MatchDoctorRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List matchDoctorRequestDescriptor = $convert.base64Decode(
    'ChJNYXRjaERvY3RvclJlcXVlc3QSOgoJc3BlY2lhbHR5GAEgASgOMhwubWFucGFzaWsudjEuRG'
    '9jdG9yU3BlY2lhbHR5UglzcGVjaWFsdHkSGgoIbGFuZ3VhZ2UYAiABKAlSCGxhbmd1YWdl');

@$core.Deprecated('Use matchDoctorResponseDescriptor instead')
const MatchDoctorResponse$json = {
  '1': 'MatchDoctorResponse',
  '2': [
    {
      '1': 'doctors',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.DoctorProfile',
      '10': 'doctors'
    },
    {'1': 'total_available', '3': 2, '4': 1, '5': 5, '10': 'totalAvailable'},
  ],
};

/// Descriptor for `MatchDoctorResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List matchDoctorResponseDescriptor = $convert.base64Decode(
    'ChNNYXRjaERvY3RvclJlc3BvbnNlEjQKB2RvY3RvcnMYASADKAsyGi5tYW5wYXNpay52MS5Eb2'
    'N0b3JQcm9maWxlUgdkb2N0b3JzEicKD3RvdGFsX2F2YWlsYWJsZRgCIAEoBVIOdG90YWxBdmFp'
    'bGFibGU=');

@$core.Deprecated('Use doctorProfileDescriptor instead')
const DoctorProfile$json = {
  '1': 'DoctorProfile',
  '2': [
    {'1': 'doctor_id', '3': 1, '4': 1, '5': 9, '10': 'doctorId'},
    {'1': 'name', '3': 2, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'specialty',
      '3': 3,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.DoctorSpecialty',
      '10': 'specialty'
    },
    {'1': 'hospital', '3': 4, '4': 1, '5': 9, '10': 'hospital'},
    {'1': 'license_number', '3': 5, '4': 1, '5': 9, '10': 'licenseNumber'},
    {'1': 'experience_years', '3': 6, '4': 1, '5': 5, '10': 'experienceYears'},
    {'1': 'rating', '3': 7, '4': 1, '5': 1, '10': 'rating'},
    {
      '1': 'total_consultations',
      '3': 8,
      '4': 1,
      '5': 5,
      '10': 'totalConsultations'
    },
    {'1': 'is_available', '3': 9, '4': 1, '5': 8, '10': 'isAvailable'},
    {'1': 'languages', '3': 10, '4': 3, '5': 9, '10': 'languages'},
    {
      '1': 'profile_image_url',
      '3': 11,
      '4': 1,
      '5': 9,
      '10': 'profileImageUrl'
    },
  ],
};

/// Descriptor for `DoctorProfile`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List doctorProfileDescriptor = $convert.base64Decode(
    'Cg1Eb2N0b3JQcm9maWxlEhsKCWRvY3Rvcl9pZBgBIAEoCVIIZG9jdG9ySWQSEgoEbmFtZRgCIA'
    'EoCVIEbmFtZRI6CglzcGVjaWFsdHkYAyABKA4yHC5tYW5wYXNpay52MS5Eb2N0b3JTcGVjaWFs'
    'dHlSCXNwZWNpYWx0eRIaCghob3NwaXRhbBgEIAEoCVIIaG9zcGl0YWwSJQoObGljZW5zZV9udW'
    '1iZXIYBSABKAlSDWxpY2Vuc2VOdW1iZXISKQoQZXhwZXJpZW5jZV95ZWFycxgGIAEoBVIPZXhw'
    'ZXJpZW5jZVllYXJzEhYKBnJhdGluZxgHIAEoAVIGcmF0aW5nEi8KE3RvdGFsX2NvbnN1bHRhdG'
    'lvbnMYCCABKAVSEnRvdGFsQ29uc3VsdGF0aW9ucxIhCgxpc19hdmFpbGFibGUYCSABKAhSC2lz'
    'QXZhaWxhYmxlEhwKCWxhbmd1YWdlcxgKIAMoCVIJbGFuZ3VhZ2VzEioKEXByb2ZpbGVfaW1hZ2'
    'VfdXJsGAsgASgJUg9wcm9maWxlSW1hZ2VVcmw=');

@$core.Deprecated('Use startVideoSessionRequestDescriptor instead')
const StartVideoSessionRequest$json = {
  '1': 'StartVideoSessionRequest',
  '2': [
    {'1': 'consultation_id', '3': 1, '4': 1, '5': 9, '10': 'consultationId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `StartVideoSessionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startVideoSessionRequestDescriptor =
    $convert.base64Decode(
        'ChhTdGFydFZpZGVvU2Vzc2lvblJlcXVlc3QSJwoPY29uc3VsdGF0aW9uX2lkGAEgASgJUg5jb2'
        '5zdWx0YXRpb25JZBIXCgd1c2VyX2lkGAIgASgJUgZ1c2VySWQ=');

@$core.Deprecated('Use videoSessionDescriptor instead')
const VideoSession$json = {
  '1': 'VideoSession',
  '2': [
    {'1': 'session_id', '3': 1, '4': 1, '5': 9, '10': 'sessionId'},
    {'1': 'consultation_id', '3': 2, '4': 1, '5': 9, '10': 'consultationId'},
    {'1': 'room_url', '3': 3, '4': 1, '5': 9, '10': 'roomUrl'},
    {'1': 'token', '3': 4, '4': 1, '5': 9, '10': 'token'},
    {
      '1': 'status',
      '3': 5,
      '4': 1,
      '5': 14,
      '6': '.manpasik.v1.VideoSessionStatus',
      '10': 'status'
    },
    {'1': 'duration_seconds', '3': 6, '4': 1, '5': 5, '10': 'durationSeconds'},
    {
      '1': 'started_at',
      '3': 7,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'startedAt'
    },
    {
      '1': 'ended_at',
      '3': 8,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'endedAt'
    },
  ],
};

/// Descriptor for `VideoSession`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List videoSessionDescriptor = $convert.base64Decode(
    'CgxWaWRlb1Nlc3Npb24SHQoKc2Vzc2lvbl9pZBgBIAEoCVIJc2Vzc2lvbklkEicKD2NvbnN1bH'
    'RhdGlvbl9pZBgCIAEoCVIOY29uc3VsdGF0aW9uSWQSGQoIcm9vbV91cmwYAyABKAlSB3Jvb21V'
    'cmwSFAoFdG9rZW4YBCABKAlSBXRva2VuEjcKBnN0YXR1cxgFIAEoDjIfLm1hbnBhc2lrLnYxLl'
    'ZpZGVvU2Vzc2lvblN0YXR1c1IGc3RhdHVzEikKEGR1cmF0aW9uX3NlY29uZHMYBiABKAVSD2R1'
    'cmF0aW9uU2Vjb25kcxI5CgpzdGFydGVkX2F0GAcgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbW'
    'VzdGFtcFIJc3RhcnRlZEF0EjUKCGVuZGVkX2F0GAggASgLMhouZ29vZ2xlLnByb3RvYnVmLlRp'
    'bWVzdGFtcFIHZW5kZWRBdA==');

@$core.Deprecated('Use endVideoSessionRequestDescriptor instead')
const EndVideoSessionRequest$json = {
  '1': 'EndVideoSessionRequest',
  '2': [
    {'1': 'session_id', '3': 1, '4': 1, '5': 9, '10': 'sessionId'},
    {'1': 'consultation_id', '3': 2, '4': 1, '5': 9, '10': 'consultationId'},
    {'1': 'doctor_notes', '3': 3, '4': 1, '5': 9, '10': 'doctorNotes'},
    {'1': 'diagnosis', '3': 4, '4': 1, '5': 9, '10': 'diagnosis'},
  ],
};

/// Descriptor for `EndVideoSessionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List endVideoSessionRequestDescriptor = $convert.base64Decode(
    'ChZFbmRWaWRlb1Nlc3Npb25SZXF1ZXN0Eh0KCnNlc3Npb25faWQYASABKAlSCXNlc3Npb25JZB'
    'InCg9jb25zdWx0YXRpb25faWQYAiABKAlSDmNvbnN1bHRhdGlvbklkEiEKDGRvY3Rvcl9ub3Rl'
    'cxgDIAEoCVILZG9jdG9yTm90ZXMSHAoJZGlhZ25vc2lzGAQgASgJUglkaWFnbm9zaXM=');

@$core.Deprecated('Use rateConsultationRequestDescriptor instead')
const RateConsultationRequest$json = {
  '1': 'RateConsultationRequest',
  '2': [
    {'1': 'consultation_id', '3': 1, '4': 1, '5': 9, '10': 'consultationId'},
    {'1': 'rating', '3': 2, '4': 1, '5': 1, '10': 'rating'},
  ],
};

/// Descriptor for `RateConsultationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List rateConsultationRequestDescriptor =
    $convert.base64Decode(
        'ChdSYXRlQ29uc3VsdGF0aW9uUmVxdWVzdBInCg9jb25zdWx0YXRpb25faWQYASABKAlSDmNvbn'
        'N1bHRhdGlvbklkEhYKBnJhdGluZxgCIAEoAVIGcmF0aW5n');

@$core.Deprecated('Use rateConsultationResponseDescriptor instead')
const RateConsultationResponse$json = {
  '1': 'RateConsultationResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
    {
      '1': 'new_average_rating',
      '3': 2,
      '4': 1,
      '5': 1,
      '10': 'newAverageRating'
    },
  ],
};

/// Descriptor for `RateConsultationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List rateConsultationResponseDescriptor =
    $convert.base64Decode(
        'ChhSYXRlQ29uc3VsdGF0aW9uUmVzcG9uc2USGAoHc3VjY2VzcxgBIAEoCFIHc3VjY2VzcxIsCh'
        'JuZXdfYXZlcmFnZV9yYXRpbmcYAiABKAFSEG5ld0F2ZXJhZ2VSYXRpbmc=');

@$core.Deprecated('Use sendFromTemplateRequestDescriptor instead')
const SendFromTemplateRequest$json = {
  '1': 'SendFromTemplateRequest',
  '2': [
    {'1': 'template_key', '3': 1, '4': 1, '5': 9, '10': 'templateKey'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'data',
      '3': 3,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.SendFromTemplateRequest.DataEntry',
      '10': 'data'
    },
  ],
  '3': [SendFromTemplateRequest_DataEntry$json],
};

@$core.Deprecated('Use sendFromTemplateRequestDescriptor instead')
const SendFromTemplateRequest_DataEntry$json = {
  '1': 'DataEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `SendFromTemplateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendFromTemplateRequestDescriptor = $convert.base64Decode(
    'ChdTZW5kRnJvbVRlbXBsYXRlUmVxdWVzdBIhCgx0ZW1wbGF0ZV9rZXkYASABKAlSC3RlbXBsYX'
    'RlS2V5EhcKB3VzZXJfaWQYAiABKAlSBnVzZXJJZBJCCgRkYXRhGAMgAygLMi4ubWFucGFzaWsu'
    'djEuU2VuZEZyb21UZW1wbGF0ZVJlcXVlc3QuRGF0YUVudHJ5UgRkYXRhGjcKCURhdGFFbnRyeR'
    'IQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use validateSharingAccessRequestDescriptor instead')
const ValidateSharingAccessRequest$json = {
  '1': 'ValidateSharingAccessRequest',
  '2': [
    {'1': 'group_id', '3': 1, '4': 1, '5': 9, '10': 'groupId'},
    {'1': 'requester_user_id', '3': 2, '4': 1, '5': 9, '10': 'requesterUserId'},
    {'1': 'target_user_id', '3': 3, '4': 1, '5': 9, '10': 'targetUserId'},
    {'1': 'biomarker', '3': 4, '4': 1, '5': 9, '10': 'biomarker'},
  ],
};

/// Descriptor for `ValidateSharingAccessRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validateSharingAccessRequestDescriptor = $convert.base64Decode(
    'ChxWYWxpZGF0ZVNoYXJpbmdBY2Nlc3NSZXF1ZXN0EhkKCGdyb3VwX2lkGAEgASgJUgdncm91cE'
    'lkEioKEXJlcXVlc3Rlcl91c2VyX2lkGAIgASgJUg9yZXF1ZXN0ZXJVc2VySWQSJAoOdGFyZ2V0'
    'X3VzZXJfaWQYAyABKAlSDHRhcmdldFVzZXJJZBIcCgliaW9tYXJrZXIYBCABKAlSCWJpb21hcm'
    'tlcg==');

@$core.Deprecated('Use validateSharingAccessResponseDescriptor instead')
const ValidateSharingAccessResponse$json = {
  '1': 'ValidateSharingAccessResponse',
  '2': [
    {'1': 'allowed', '3': 1, '4': 1, '5': 8, '10': 'allowed'},
    {'1': 'deny_reason', '3': 2, '4': 1, '5': 9, '10': 'denyReason'},
  ],
};

/// Descriptor for `ValidateSharingAccessResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validateSharingAccessResponseDescriptor =
    $convert.base64Decode(
        'Ch1WYWxpZGF0ZVNoYXJpbmdBY2Nlc3NSZXNwb25zZRIYCgdhbGxvd2VkGAEgASgIUgdhbGxvd2'
        'VkEh8KC2RlbnlfcmVhc29uGAIgASgJUgpkZW55UmVhc29u');

@$core.Deprecated('Use getAuditLogDetailsRequestDescriptor instead')
const GetAuditLogDetailsRequest$json = {
  '1': 'GetAuditLogDetailsRequest',
  '2': [
    {'1': 'limit', '3': 1, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 2, '4': 1, '5': 5, '10': 'offset'},
    {'1': 'admin_id', '3': 3, '4': 1, '5': 9, '10': 'adminId'},
    {'1': 'action', '3': 4, '4': 1, '5': 9, '10': 'action'},
  ],
};

/// Descriptor for `GetAuditLogDetailsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAuditLogDetailsRequestDescriptor = $convert.base64Decode(
    'ChlHZXRBdWRpdExvZ0RldGFpbHNSZXF1ZXN0EhQKBWxpbWl0GAEgASgFUgVsaW1pdBIWCgZvZm'
    'ZzZXQYAiABKAVSBm9mZnNldBIZCghhZG1pbl9pZBgDIAEoCVIHYWRtaW5JZBIWCgZhY3Rpb24Y'
    'BCABKAlSBmFjdGlvbg==');

@$core.Deprecated('Use getAuditLogDetailsResponseDescriptor instead')
const GetAuditLogDetailsResponse$json = {
  '1': 'GetAuditLogDetailsResponse',
  '2': [
    {
      '1': 'details',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.AuditLogDetail',
      '10': 'details'
    },
    {'1': 'total', '3': 2, '4': 1, '5': 5, '10': 'total'},
  ],
};

/// Descriptor for `GetAuditLogDetailsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAuditLogDetailsResponseDescriptor =
    $convert.base64Decode(
        'ChpHZXRBdWRpdExvZ0RldGFpbHNSZXNwb25zZRI1CgdkZXRhaWxzGAEgAygLMhsubWFucGFzaW'
        'sudjEuQXVkaXRMb2dEZXRhaWxSB2RldGFpbHMSFAoFdG90YWwYAiABKAVSBXRvdGFs');

@$core.Deprecated('Use auditLogDetailDescriptor instead')
const AuditLogDetail$json = {
  '1': 'AuditLogDetail',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'admin_id', '3': 2, '4': 1, '5': 9, '10': 'adminId'},
    {'1': 'action', '3': 3, '4': 1, '5': 9, '10': 'action'},
    {'1': 'resource_type', '3': 4, '4': 1, '5': 9, '10': 'resourceType'},
    {'1': 'resource_id', '3': 5, '4': 1, '5': 9, '10': 'resourceId'},
    {'1': 'old_value', '3': 6, '4': 1, '5': 9, '10': 'oldValue'},
    {'1': 'new_value', '3': 7, '4': 1, '5': 9, '10': 'newValue'},
    {'1': 'description', '3': 8, '4': 1, '5': 9, '10': 'description'},
    {'1': 'ip_address', '3': 9, '4': 1, '5': 9, '10': 'ipAddress'},
    {
      '1': 'created_at',
      '3': 10,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
  ],
};

/// Descriptor for `AuditLogDetail`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List auditLogDetailDescriptor = $convert.base64Decode(
    'Cg5BdWRpdExvZ0RldGFpbBIOCgJpZBgBIAEoCVICaWQSGQoIYWRtaW5faWQYAiABKAlSB2FkbW'
    'luSWQSFgoGYWN0aW9uGAMgASgJUgZhY3Rpb24SIwoNcmVzb3VyY2VfdHlwZRgEIAEoCVIMcmVz'
    'b3VyY2VUeXBlEh8KC3Jlc291cmNlX2lkGAUgASgJUgpyZXNvdXJjZUlkEhsKCW9sZF92YWx1ZR'
    'gGIAEoCVIIb2xkVmFsdWUSGwoJbmV3X3ZhbHVlGAcgASgJUghuZXdWYWx1ZRIgCgtkZXNjcmlw'
    'dGlvbhgIIAEoCVILZGVzY3JpcHRpb24SHQoKaXBfYWRkcmVzcxgJIAEoCVIJaXBBZGRyZXNzEj'
    'kKCmNyZWF0ZWRfYXQYCiABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUgljcmVhdGVk'
    'QXQ=');

@$core.Deprecated('Use streamChatRequestDescriptor instead')
const StreamChatRequest$json = {
  '1': 'StreamChatRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'message', '3': 2, '4': 1, '5': 9, '10': 'message'},
    {'1': 'session_id', '3': 3, '4': 1, '5': 9, '10': 'sessionId'},
    {
      '1': 'context_measurement_ids',
      '3': 4,
      '4': 3,
      '5': 9,
      '10': 'contextMeasurementIds'
    },
  ],
};

/// Descriptor for `StreamChatRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List streamChatRequestDescriptor = $convert.base64Decode(
    'ChFTdHJlYW1DaGF0UmVxdWVzdBIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQSGAoHbWVzc2FnZR'
    'gCIAEoCVIHbWVzc2FnZRIdCgpzZXNzaW9uX2lkGAMgASgJUglzZXNzaW9uSWQSNgoXY29udGV4'
    'dF9tZWFzdXJlbWVudF9pZHMYBCADKAlSFWNvbnRleHRNZWFzdXJlbWVudElkcw==');

@$core.Deprecated('Use streamChatResponseDescriptor instead')
const StreamChatResponse$json = {
  '1': 'StreamChatResponse',
  '2': [
    {'1': 'chunk', '3': 1, '4': 1, '5': 9, '10': 'chunk'},
    {'1': 'is_final', '3': 2, '4': 1, '5': 8, '10': 'isFinal'},
    {'1': 'session_id', '3': 3, '4': 1, '5': 9, '10': 'sessionId'},
    {'1': 'full_response', '3': 4, '4': 1, '5': 9, '10': 'fullResponse'},
  ],
};

/// Descriptor for `StreamChatResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List streamChatResponseDescriptor = $convert.base64Decode(
    'ChJTdHJlYW1DaGF0UmVzcG9uc2USFAoFY2h1bmsYASABKAlSBWNodW5rEhkKCGlzX2ZpbmFsGA'
    'IgASgIUgdpc0ZpbmFsEh0KCnNlc3Npb25faWQYAyABKAlSCXNlc3Npb25JZBIjCg1mdWxsX3Jl'
    'c3BvbnNlGAQgASgJUgxmdWxsUmVzcG9uc2U=');

@$core.Deprecated('Use getChallengeLeaderboardRequestDescriptor instead')
const GetChallengeLeaderboardRequest$json = {
  '1': 'GetChallengeLeaderboardRequest',
  '2': [
    {'1': 'challenge_id', '3': 1, '4': 1, '5': 9, '10': 'challengeId'},
    {'1': 'limit', '3': 2, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 3, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `GetChallengeLeaderboardRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getChallengeLeaderboardRequestDescriptor =
    $convert.base64Decode(
        'Ch5HZXRDaGFsbGVuZ2VMZWFkZXJib2FyZFJlcXVlc3QSIQoMY2hhbGxlbmdlX2lkGAEgASgJUg'
        'tjaGFsbGVuZ2VJZBIUCgVsaW1pdBgCIAEoBVIFbGltaXQSFgoGb2Zmc2V0GAMgASgFUgZvZmZz'
        'ZXQ=');

@$core.Deprecated('Use getChallengeLeaderboardResponseDescriptor instead')
const GetChallengeLeaderboardResponse$json = {
  '1': 'GetChallengeLeaderboardResponse',
  '2': [
    {
      '1': 'entries',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.LeaderboardEntry',
      '10': 'entries'
    },
    {
      '1': 'total_participants',
      '3': 2,
      '4': 1,
      '5': 5,
      '10': 'totalParticipants'
    },
    {
      '1': 'my_entry',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.manpasik.v1.LeaderboardEntry',
      '10': 'myEntry'
    },
  ],
};

/// Descriptor for `GetChallengeLeaderboardResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getChallengeLeaderboardResponseDescriptor =
    $convert.base64Decode(
        'Ch9HZXRDaGFsbGVuZ2VMZWFkZXJib2FyZFJlc3BvbnNlEjcKB2VudHJpZXMYASADKAsyHS5tYW'
        '5wYXNpay52MS5MZWFkZXJib2FyZEVudHJ5UgdlbnRyaWVzEi0KEnRvdGFsX3BhcnRpY2lwYW50'
        'cxgCIAEoBVIRdG90YWxQYXJ0aWNpcGFudHMSOAoIbXlfZW50cnkYAyABKAsyHS5tYW5wYXNpay'
        '52MS5MZWFkZXJib2FyZEVudHJ5UgdteUVudHJ5');

@$core.Deprecated('Use leaderboardEntryDescriptor instead')
const LeaderboardEntry$json = {
  '1': 'LeaderboardEntry',
  '2': [
    {'1': 'rank', '3': 1, '4': 1, '5': 5, '10': 'rank'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'display_name', '3': 3, '4': 1, '5': 9, '10': 'displayName'},
    {'1': 'avatar_url', '3': 4, '4': 1, '5': 9, '10': 'avatarUrl'},
    {'1': 'progress_value', '3': 5, '4': 1, '5': 1, '10': 'progressValue'},
    {'1': 'target_value', '3': 6, '4': 1, '5': 1, '10': 'targetValue'},
    {'1': 'progress_pct', '3': 7, '4': 1, '5': 1, '10': 'progressPct'},
    {'1': 'streak_days', '3': 8, '4': 1, '5': 5, '10': 'streakDays'},
    {
      '1': 'last_updated',
      '3': 9,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'lastUpdated'
    },
  ],
};

/// Descriptor for `LeaderboardEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List leaderboardEntryDescriptor = $convert.base64Decode(
    'ChBMZWFkZXJib2FyZEVudHJ5EhIKBHJhbmsYASABKAVSBHJhbmsSFwoHdXNlcl9pZBgCIAEoCV'
    'IGdXNlcklkEiEKDGRpc3BsYXlfbmFtZRgDIAEoCVILZGlzcGxheU5hbWUSHQoKYXZhdGFyX3Vy'
    'bBgEIAEoCVIJYXZhdGFyVXJsEiUKDnByb2dyZXNzX3ZhbHVlGAUgASgBUg1wcm9ncmVzc1ZhbH'
    'VlEiEKDHRhcmdldF92YWx1ZRgGIAEoAVILdGFyZ2V0VmFsdWUSIQoMcHJvZ3Jlc3NfcGN0GAcg'
    'ASgBUgtwcm9ncmVzc1BjdBIfCgtzdHJlYWtfZGF5cxgIIAEoBVIKc3RyZWFrRGF5cxI9CgxsYX'
    'N0X3VwZGF0ZWQYCSABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUgtsYXN0VXBkYXRl'
    'ZA==');

@$core.Deprecated('Use updateChallengeProgressRequestDescriptor instead')
const UpdateChallengeProgressRequest$json = {
  '1': 'UpdateChallengeProgressRequest',
  '2': [
    {'1': 'challenge_id', '3': 1, '4': 1, '5': 9, '10': 'challengeId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'value', '3': 3, '4': 1, '5': 1, '10': 'value'},
    {'1': 'measurement_id', '3': 4, '4': 1, '5': 9, '10': 'measurementId'},
  ],
};

/// Descriptor for `UpdateChallengeProgressRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateChallengeProgressRequestDescriptor =
    $convert.base64Decode(
        'Ch5VcGRhdGVDaGFsbGVuZ2VQcm9ncmVzc1JlcXVlc3QSIQoMY2hhbGxlbmdlX2lkGAEgASgJUg'
        'tjaGFsbGVuZ2VJZBIXCgd1c2VyX2lkGAIgASgJUgZ1c2VySWQSFAoFdmFsdWUYAyABKAFSBXZh'
        'bHVlEiUKDm1lYXN1cmVtZW50X2lkGAQgASgJUg1tZWFzdXJlbWVudElk');

@$core.Deprecated('Use updateChallengeProgressResponseDescriptor instead')
const UpdateChallengeProgressResponse$json = {
  '1': 'UpdateChallengeProgressResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
    {'1': 'new_progress', '3': 2, '4': 1, '5': 1, '10': 'newProgress'},
    {'1': 'target_value', '3': 3, '4': 1, '5': 1, '10': 'targetValue'},
    {'1': 'new_rank', '3': 4, '4': 1, '5': 5, '10': 'newRank'},
  ],
};

/// Descriptor for `UpdateChallengeProgressResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateChallengeProgressResponseDescriptor =
    $convert.base64Decode(
        'Ch9VcGRhdGVDaGFsbGVuZ2VQcm9ncmVzc1Jlc3BvbnNlEhgKB3N1Y2Nlc3MYASABKAhSB3N1Y2'
        'Nlc3MSIQoMbmV3X3Byb2dyZXNzGAIgASgBUgtuZXdQcm9ncmVzcxIhCgx0YXJnZXRfdmFsdWUY'
        'AyABKAFSC3RhcmdldFZhbHVlEhkKCG5ld19yYW5rGAQgASgFUgduZXdSYW5r');

@$core.Deprecated('Use getRevenueStatsRequestDescriptor instead')
const GetRevenueStatsRequest$json = {
  '1': 'GetRevenueStatsRequest',
  '2': [
    {'1': 'period', '3': 1, '4': 1, '5': 9, '10': 'period'},
    {'1': 'start_date', '3': 2, '4': 1, '5': 9, '10': 'startDate'},
    {'1': 'end_date', '3': 3, '4': 1, '5': 9, '10': 'endDate'},
  ],
};

/// Descriptor for `GetRevenueStatsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRevenueStatsRequestDescriptor = $convert.base64Decode(
    'ChZHZXRSZXZlbnVlU3RhdHNSZXF1ZXN0EhYKBnBlcmlvZBgBIAEoCVIGcGVyaW9kEh0KCnN0YX'
    'J0X2RhdGUYAiABKAlSCXN0YXJ0RGF0ZRIZCghlbmRfZGF0ZRgDIAEoCVIHZW5kRGF0ZQ==');

@$core.Deprecated('Use getRevenueStatsResponseDescriptor instead')
const GetRevenueStatsResponse$json = {
  '1': 'GetRevenueStatsResponse',
  '2': [
    {'1': 'total_revenue_krw', '3': 1, '4': 1, '5': 3, '10': 'totalRevenueKrw'},
    {
      '1': 'subscription_revenue_krw',
      '3': 2,
      '4': 1,
      '5': 3,
      '10': 'subscriptionRevenueKrw'
    },
    {
      '1': 'product_revenue_krw',
      '3': 3,
      '4': 1,
      '5': 3,
      '10': 'productRevenueKrw'
    },
    {
      '1': 'total_transactions',
      '3': 4,
      '4': 1,
      '5': 5,
      '10': 'totalTransactions'
    },
    {
      '1': 'periods',
      '3': 5,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.RevenuePeriod',
      '10': 'periods'
    },
    {
      '1': 'revenue_by_tier',
      '3': 6,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.GetRevenueStatsResponse.RevenueByTierEntry',
      '10': 'revenueByTier'
    },
  ],
  '3': [GetRevenueStatsResponse_RevenueByTierEntry$json],
};

@$core.Deprecated('Use getRevenueStatsResponseDescriptor instead')
const GetRevenueStatsResponse_RevenueByTierEntry$json = {
  '1': 'RevenueByTierEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 3, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetRevenueStatsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRevenueStatsResponseDescriptor = $convert.base64Decode(
    'ChdHZXRSZXZlbnVlU3RhdHNSZXNwb25zZRIqChF0b3RhbF9yZXZlbnVlX2tydxgBIAEoA1IPdG'
    '90YWxSZXZlbnVlS3J3EjgKGHN1YnNjcmlwdGlvbl9yZXZlbnVlX2tydxgCIAEoA1IWc3Vic2Ny'
    'aXB0aW9uUmV2ZW51ZUtydxIuChNwcm9kdWN0X3JldmVudWVfa3J3GAMgASgDUhFwcm9kdWN0Um'
    'V2ZW51ZUtydxItChJ0b3RhbF90cmFuc2FjdGlvbnMYBCABKAVSEXRvdGFsVHJhbnNhY3Rpb25z'
    'EjQKB3BlcmlvZHMYBSADKAsyGi5tYW5wYXNpay52MS5SZXZlbnVlUGVyaW9kUgdwZXJpb2RzEl'
    '8KD3JldmVudWVfYnlfdGllchgGIAMoCzI3Lm1hbnBhc2lrLnYxLkdldFJldmVudWVTdGF0c1Jl'
    'c3BvbnNlLlJldmVudWVCeVRpZXJFbnRyeVINcmV2ZW51ZUJ5VGllchpAChJSZXZlbnVlQnlUaW'
    'VyRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKANSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use revenuePeriodDescriptor instead')
const RevenuePeriod$json = {
  '1': 'RevenuePeriod',
  '2': [
    {'1': 'label', '3': 1, '4': 1, '5': 9, '10': 'label'},
    {'1': 'revenue_krw', '3': 2, '4': 1, '5': 3, '10': 'revenueKrw'},
    {
      '1': 'transaction_count',
      '3': 3,
      '4': 1,
      '5': 5,
      '10': 'transactionCount'
    },
  ],
};

/// Descriptor for `RevenuePeriod`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List revenuePeriodDescriptor = $convert.base64Decode(
    'Cg1SZXZlbnVlUGVyaW9kEhQKBWxhYmVsGAEgASgJUgVsYWJlbBIfCgtyZXZlbnVlX2tydxgCIA'
    'EoA1IKcmV2ZW51ZUtydxIrChF0cmFuc2FjdGlvbl9jb3VudBgDIAEoBVIQdHJhbnNhY3Rpb25D'
    'b3VudA==');

@$core.Deprecated('Use getInventoryStatsRequestDescriptor instead')
const GetInventoryStatsRequest$json = {
  '1': 'GetInventoryStatsRequest',
  '2': [
    {'1': 'category_filter', '3': 1, '4': 1, '5': 5, '10': 'categoryFilter'},
  ],
};

/// Descriptor for `GetInventoryStatsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getInventoryStatsRequestDescriptor =
    $convert.base64Decode(
        'ChhHZXRJbnZlbnRvcnlTdGF0c1JlcXVlc3QSJwoPY2F0ZWdvcnlfZmlsdGVyGAEgASgFUg5jYX'
        'RlZ29yeUZpbHRlcg==');

@$core.Deprecated('Use getInventoryStatsResponseDescriptor instead')
const GetInventoryStatsResponse$json = {
  '1': 'GetInventoryStatsResponse',
  '2': [
    {
      '1': 'items',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.InventoryItem',
      '10': 'items'
    },
    {'1': 'total_products', '3': 2, '4': 1, '5': 5, '10': 'totalProducts'},
    {'1': 'low_stock_count', '3': 3, '4': 1, '5': 5, '10': 'lowStockCount'},
    {
      '1': 'out_of_stock_count',
      '3': 4,
      '4': 1,
      '5': 5,
      '10': 'outOfStockCount'
    },
  ],
};

/// Descriptor for `GetInventoryStatsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getInventoryStatsResponseDescriptor = $convert.base64Decode(
    'ChlHZXRJbnZlbnRvcnlTdGF0c1Jlc3BvbnNlEjAKBWl0ZW1zGAEgAygLMhoubWFucGFzaWsudj'
    'EuSW52ZW50b3J5SXRlbVIFaXRlbXMSJQoOdG90YWxfcHJvZHVjdHMYAiABKAVSDXRvdGFsUHJv'
    'ZHVjdHMSJgoPbG93X3N0b2NrX2NvdW50GAMgASgFUg1sb3dTdG9ja0NvdW50EisKEm91dF9vZl'
    '9zdG9ja19jb3VudBgEIAEoBVIPb3V0T2ZTdG9ja0NvdW50');

@$core.Deprecated('Use inventoryItemDescriptor instead')
const InventoryItem$json = {
  '1': 'InventoryItem',
  '2': [
    {'1': 'product_id', '3': 1, '4': 1, '5': 9, '10': 'productId'},
    {'1': 'product_name', '3': 2, '4': 1, '5': 9, '10': 'productName'},
    {'1': 'category', '3': 3, '4': 1, '5': 5, '10': 'category'},
    {'1': 'current_stock', '3': 4, '4': 1, '5': 5, '10': 'currentStock'},
    {
      '1': 'min_stock_threshold',
      '3': 5,
      '4': 1,
      '5': 5,
      '10': 'minStockThreshold'
    },
    {'1': 'monthly_sales', '3': 6, '4': 1, '5': 5, '10': 'monthlySales'},
    {'1': 'price_krw', '3': 7, '4': 1, '5': 5, '10': 'priceKrw'},
    {'1': 'status', '3': 8, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `InventoryItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inventoryItemDescriptor = $convert.base64Decode(
    'Cg1JbnZlbnRvcnlJdGVtEh0KCnByb2R1Y3RfaWQYASABKAlSCXByb2R1Y3RJZBIhCgxwcm9kdW'
    'N0X25hbWUYAiABKAlSC3Byb2R1Y3ROYW1lEhoKCGNhdGVnb3J5GAMgASgFUghjYXRlZ29yeRIj'
    'Cg1jdXJyZW50X3N0b2NrGAQgASgFUgxjdXJyZW50U3RvY2sSLgoTbWluX3N0b2NrX3RocmVzaG'
    '9sZBgFIAEoBVIRbWluU3RvY2tUaHJlc2hvbGQSIwoNbW9udGhseV9zYWxlcxgGIAEoBVIMbW9u'
    'dGhseVNhbGVzEhsKCXByaWNlX2tydxgHIAEoBVIIcHJpY2VLcncSFgoGc3RhdHVzGAggASgJUg'
    'ZzdGF0dXM=');

@$core.Deprecated('Use translateRealtimeRequestDescriptor instead')
const TranslateRealtimeRequest$json = {
  '1': 'TranslateRealtimeRequest',
  '2': [
    {'1': 'text', '3': 1, '4': 1, '5': 9, '10': 'text'},
    {'1': 'source_language', '3': 2, '4': 1, '5': 9, '10': 'sourceLanguage'},
    {'1': 'target_language', '3': 3, '4': 1, '5': 9, '10': 'targetLanguage'},
    {'1': 'is_medical', '3': 4, '4': 1, '5': 8, '10': 'isMedical'},
    {'1': 'context', '3': 5, '4': 1, '5': 9, '10': 'context'},
    {'1': 'session_id', '3': 6, '4': 1, '5': 9, '10': 'sessionId'},
  ],
};

/// Descriptor for `TranslateRealtimeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List translateRealtimeRequestDescriptor = $convert.base64Decode(
    'ChhUcmFuc2xhdGVSZWFsdGltZVJlcXVlc3QSEgoEdGV4dBgBIAEoCVIEdGV4dBInCg9zb3VyY2'
    'VfbGFuZ3VhZ2UYAiABKAlSDnNvdXJjZUxhbmd1YWdlEicKD3RhcmdldF9sYW5ndWFnZRgDIAEo'
    'CVIOdGFyZ2V0TGFuZ3VhZ2USHQoKaXNfbWVkaWNhbBgEIAEoCFIJaXNNZWRpY2FsEhgKB2Nvbn'
    'RleHQYBSABKAlSB2NvbnRleHQSHQoKc2Vzc2lvbl9pZBgGIAEoCVIJc2Vzc2lvbklk');

@$core.Deprecated('Use translateRealtimeResponseDescriptor instead')
const TranslateRealtimeResponse$json = {
  '1': 'TranslateRealtimeResponse',
  '2': [
    {'1': 'translated_text', '3': 1, '4': 1, '5': 9, '10': 'translatedText'},
    {'1': 'source_language', '3': 2, '4': 1, '5': 9, '10': 'sourceLanguage'},
    {'1': 'target_language', '3': 3, '4': 1, '5': 9, '10': 'targetLanguage'},
    {'1': 'confidence', '3': 4, '4': 1, '5': 1, '10': 'confidence'},
    {'1': 'is_medical_term', '3': 5, '4': 1, '5': 8, '10': 'isMedicalTerm'},
    {'1': 'latency_ms', '3': 6, '4': 1, '5': 3, '10': 'latencyMs'},
    {
      '1': 'medical_terms',
      '3': 7,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.MedicalTermMapping',
      '10': 'medicalTerms'
    },
  ],
};

/// Descriptor for `TranslateRealtimeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List translateRealtimeResponseDescriptor = $convert.base64Decode(
    'ChlUcmFuc2xhdGVSZWFsdGltZVJlc3BvbnNlEicKD3RyYW5zbGF0ZWRfdGV4dBgBIAEoCVIOdH'
    'JhbnNsYXRlZFRleHQSJwoPc291cmNlX2xhbmd1YWdlGAIgASgJUg5zb3VyY2VMYW5ndWFnZRIn'
    'Cg90YXJnZXRfbGFuZ3VhZ2UYAyABKAlSDnRhcmdldExhbmd1YWdlEh4KCmNvbmZpZGVuY2UYBC'
    'ABKAFSCmNvbmZpZGVuY2USJgoPaXNfbWVkaWNhbF90ZXJtGAUgASgIUg1pc01lZGljYWxUZXJt'
    'Eh0KCmxhdGVuY3lfbXMYBiABKANSCWxhdGVuY3lNcxJECg1tZWRpY2FsX3Rlcm1zGAcgAygLMh'
    '8ubWFucGFzaWsudjEuTWVkaWNhbFRlcm1NYXBwaW5nUgxtZWRpY2FsVGVybXM=');

@$core.Deprecated('Use medicalTermMappingDescriptor instead')
const MedicalTermMapping$json = {
  '1': 'MedicalTermMapping',
  '2': [
    {'1': 'original', '3': 1, '4': 1, '5': 9, '10': 'original'},
    {'1': 'translated', '3': 2, '4': 1, '5': 9, '10': 'translated'},
    {'1': 'category', '3': 3, '4': 1, '5': 9, '10': 'category'},
  ],
};

/// Descriptor for `MedicalTermMapping`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List medicalTermMappingDescriptor = $convert.base64Decode(
    'ChJNZWRpY2FsVGVybU1hcHBpbmcSGgoIb3JpZ2luYWwYASABKAlSCG9yaWdpbmFsEh4KCnRyYW'
    '5zbGF0ZWQYAiABKAlSCnRyYW5zbGF0ZWQSGgoIY2F0ZWdvcnkYAyABKAlSCGNhdGVnb3J5');

@$core.Deprecated('Use assistantCommandRequestDescriptor instead')
const AssistantCommandRequest$json = {
  '1': 'AssistantCommandRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'text', '3': 2, '4': 1, '5': 9, '10': 'text'},
    {'1': 'session_id', '3': 3, '4': 1, '5': 9, '10': 'sessionId'},
  ],
};

/// Descriptor for `AssistantCommandRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List assistantCommandRequestDescriptor =
    $convert.base64Decode(
        'ChdBc3Npc3RhbnRDb21tYW5kUmVxdWVzdBIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQSEgoEdG'
        'V4dBgCIAEoCVIEdGV4dBIdCgpzZXNzaW9uX2lkGAMgASgJUglzZXNzaW9uSWQ=');

@$core.Deprecated('Use assistantCommandResponseDescriptor instead')
const AssistantCommandResponse$json = {
  '1': 'AssistantCommandResponse',
  '2': [
    {'1': 'session_id', '3': 1, '4': 1, '5': 9, '10': 'sessionId'},
    {'1': 'turn_id', '3': 2, '4': 1, '5': 9, '10': 'turnId'},
    {'1': 'response_text', '3': 3, '4': 1, '5': 9, '10': 'responseText'},
    {'1': 'intent', '3': 4, '4': 1, '5': 9, '10': 'intent'},
    {'1': 'action_type', '3': 5, '4': 1, '5': 9, '10': 'actionType'},
    {'1': 'action_result', '3': 6, '4': 1, '5': 9, '10': 'actionResult'},
    {
      '1': 'requires_confirmation',
      '3': 7,
      '4': 1,
      '5': 8,
      '10': 'requiresConfirmation'
    },
    {
      '1': 'confirmation_prompt',
      '3': 8,
      '4': 1,
      '5': 9,
      '10': 'confirmationPrompt'
    },
  ],
};

/// Descriptor for `AssistantCommandResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List assistantCommandResponseDescriptor = $convert.base64Decode(
    'ChhBc3Npc3RhbnRDb21tYW5kUmVzcG9uc2USHQoKc2Vzc2lvbl9pZBgBIAEoCVIJc2Vzc2lvbk'
    'lkEhcKB3R1cm5faWQYAiABKAlSBnR1cm5JZBIjCg1yZXNwb25zZV90ZXh0GAMgASgJUgxyZXNw'
    'b25zZVRleHQSFgoGaW50ZW50GAQgASgJUgZpbnRlbnQSHwoLYWN0aW9uX3R5cGUYBSABKAlSCm'
    'FjdGlvblR5cGUSIwoNYWN0aW9uX3Jlc3VsdBgGIAEoCVIMYWN0aW9uUmVzdWx0EjMKFXJlcXVp'
    'cmVzX2NvbmZpcm1hdGlvbhgHIAEoCFIUcmVxdWlyZXNDb25maXJtYXRpb24SLwoTY29uZmlybW'
    'F0aW9uX3Byb21wdBgIIAEoCVISY29uZmlybWF0aW9uUHJvbXB0');

@$core.Deprecated('Use assistantSessionDescriptor instead')
const AssistantSession$json = {
  '1': 'AssistantSession',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'title', '3': 3, '4': 1, '5': 9, '10': 'title'},
    {'1': 'status', '3': 4, '4': 1, '5': 9, '10': 'status'},
    {'1': 'turn_count', '3': 5, '4': 1, '5': 5, '10': 'turnCount'},
    {
      '1': 'created_at',
      '3': 6,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
    {
      '1': 'updated_at',
      '3': 7,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'updatedAt'
    },
  ],
};

/// Descriptor for `AssistantSession`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List assistantSessionDescriptor = $convert.base64Decode(
    'ChBBc3Npc3RhbnRTZXNzaW9uEg4KAmlkGAEgASgJUgJpZBIXCgd1c2VyX2lkGAIgASgJUgZ1c2'
    'VySWQSFAoFdGl0bGUYAyABKAlSBXRpdGxlEhYKBnN0YXR1cxgEIAEoCVIGc3RhdHVzEh0KCnR1'
    'cm5fY291bnQYBSABKAVSCXR1cm5Db3VudBI5CgpjcmVhdGVkX2F0GAYgASgLMhouZ29vZ2xlLn'
    'Byb3RvYnVmLlRpbWVzdGFtcFIJY3JlYXRlZEF0EjkKCnVwZGF0ZWRfYXQYByABKAsyGi5nb29n'
    'bGUucHJvdG9idWYuVGltZXN0YW1wUgl1cGRhdGVkQXQ=');

@$core.Deprecated('Use assistantTurnDescriptor instead')
const AssistantTurn$json = {
  '1': 'AssistantTurn',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'session_id', '3': 2, '4': 1, '5': 9, '10': 'sessionId'},
    {'1': 'role', '3': 3, '4': 1, '5': 9, '10': 'role'},
    {'1': 'content', '3': 4, '4': 1, '5': 9, '10': 'content'},
    {'1': 'intent', '3': 5, '4': 1, '5': 9, '10': 'intent'},
    {'1': 'action_type', '3': 6, '4': 1, '5': 9, '10': 'actionType'},
    {'1': 'action_result', '3': 7, '4': 1, '5': 9, '10': 'actionResult'},
    {
      '1': 'created_at',
      '3': 8,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
  ],
};

/// Descriptor for `AssistantTurn`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List assistantTurnDescriptor = $convert.base64Decode(
    'Cg1Bc3Npc3RhbnRUdXJuEg4KAmlkGAEgASgJUgJpZBIdCgpzZXNzaW9uX2lkGAIgASgJUglzZX'
    'NzaW9uSWQSEgoEcm9sZRgDIAEoCVIEcm9sZRIYCgdjb250ZW50GAQgASgJUgdjb250ZW50EhYK'
    'BmludGVudBgFIAEoCVIGaW50ZW50Eh8KC2FjdGlvbl90eXBlGAYgASgJUgphY3Rpb25UeXBlEi'
    'MKDWFjdGlvbl9yZXN1bHQYByABKAlSDGFjdGlvblJlc3VsdBI5CgpjcmVhdGVkX2F0GAggASgL'
    'MhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFIJY3JlYXRlZEF0');

@$core.Deprecated('Use getAssistantSessionRequestDescriptor instead')
const GetAssistantSessionRequest$json = {
  '1': 'GetAssistantSessionRequest',
  '2': [
    {'1': 'session_id', '3': 1, '4': 1, '5': 9, '10': 'sessionId'},
  ],
};

/// Descriptor for `GetAssistantSessionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAssistantSessionRequestDescriptor =
    $convert.base64Decode(
        'ChpHZXRBc3Npc3RhbnRTZXNzaW9uUmVxdWVzdBIdCgpzZXNzaW9uX2lkGAEgASgJUglzZXNzaW'
        '9uSWQ=');

@$core.Deprecated('Use listAssistantSessionsRequestDescriptor instead')
const ListAssistantSessionsRequest$json = {
  '1': 'ListAssistantSessionsRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'limit', '3': 2, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 3, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `ListAssistantSessionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAssistantSessionsRequestDescriptor =
    $convert.base64Decode(
        'ChxMaXN0QXNzaXN0YW50U2Vzc2lvbnNSZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZB'
        'IUCgVsaW1pdBgCIAEoBVIFbGltaXQSFgoGb2Zmc2V0GAMgASgFUgZvZmZzZXQ=');

@$core.Deprecated('Use listAssistantSessionsResponseDescriptor instead')
const ListAssistantSessionsResponse$json = {
  '1': 'ListAssistantSessionsResponse',
  '2': [
    {
      '1': 'sessions',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.AssistantSession',
      '10': 'sessions'
    },
    {'1': 'total', '3': 2, '4': 1, '5': 5, '10': 'total'},
  ],
};

/// Descriptor for `ListAssistantSessionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAssistantSessionsResponseDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0QXNzaXN0YW50U2Vzc2lvbnNSZXNwb25zZRI5CghzZXNzaW9ucxgBIAMoCzIdLm1hbn'
        'Bhc2lrLnYxLkFzc2lzdGFudFNlc3Npb25SCHNlc3Npb25zEhQKBXRvdGFsGAIgASgFUgV0b3Rh'
        'bA==');

@$core.Deprecated('Use listAssistantTurnsRequestDescriptor instead')
const ListAssistantTurnsRequest$json = {
  '1': 'ListAssistantTurnsRequest',
  '2': [
    {'1': 'session_id', '3': 1, '4': 1, '5': 9, '10': 'sessionId'},
    {'1': 'limit', '3': 2, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 3, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `ListAssistantTurnsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAssistantTurnsRequestDescriptor =
    $convert.base64Decode(
        'ChlMaXN0QXNzaXN0YW50VHVybnNSZXF1ZXN0Eh0KCnNlc3Npb25faWQYASABKAlSCXNlc3Npb2'
        '5JZBIUCgVsaW1pdBgCIAEoBVIFbGltaXQSFgoGb2Zmc2V0GAMgASgFUgZvZmZzZXQ=');

@$core.Deprecated('Use listAssistantTurnsResponseDescriptor instead')
const ListAssistantTurnsResponse$json = {
  '1': 'ListAssistantTurnsResponse',
  '2': [
    {
      '1': 'turns',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.AssistantTurn',
      '10': 'turns'
    },
    {'1': 'total', '3': 2, '4': 1, '5': 5, '10': 'total'},
  ],
};

/// Descriptor for `ListAssistantTurnsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAssistantTurnsResponseDescriptor =
    $convert.base64Decode(
        'ChpMaXN0QXNzaXN0YW50VHVybnNSZXNwb25zZRIwCgV0dXJucxgBIAMoCzIaLm1hbnBhc2lrLn'
        'YxLkFzc2lzdGFudFR1cm5SBXR1cm5zEhQKBXRvdGFsGAIgASgFUgV0b3RhbA==');

@$core.Deprecated('Use deleteAssistantSessionRequestDescriptor instead')
const DeleteAssistantSessionRequest$json = {
  '1': 'DeleteAssistantSessionRequest',
  '2': [
    {'1': 'session_id', '3': 1, '4': 1, '5': 9, '10': 'sessionId'},
  ],
};

/// Descriptor for `DeleteAssistantSessionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteAssistantSessionRequestDescriptor =
    $convert.base64Decode(
        'Ch1EZWxldGVBc3Npc3RhbnRTZXNzaW9uUmVxdWVzdBIdCgpzZXNzaW9uX2lkGAEgASgJUglzZX'
        'NzaW9uSWQ=');

@$core.Deprecated('Use deleteAssistantSessionResponseDescriptor instead')
const DeleteAssistantSessionResponse$json = {
  '1': 'DeleteAssistantSessionResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
  ],
};

/// Descriptor for `DeleteAssistantSessionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteAssistantSessionResponseDescriptor =
    $convert.base64Decode(
        'Ch5EZWxldGVBc3Npc3RhbnRTZXNzaW9uUmVzcG9uc2USGAoHc3VjY2VzcxgBIAEoCFIHc3VjY2'
        'Vzcw==');

@$core.Deprecated('Use analyzeFoodRequestDescriptor instead')
const AnalyzeFoodRequest$json = {
  '1': 'AnalyzeFoodRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'image_url', '3': 2, '4': 1, '5': 9, '10': 'imageUrl'},
    {'1': 'meal_type', '3': 3, '4': 1, '5': 9, '10': 'mealType'},
  ],
};

/// Descriptor for `AnalyzeFoodRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List analyzeFoodRequestDescriptor = $convert.base64Decode(
    'ChJBbmFseXplRm9vZFJlcXVlc3QSFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklkEhsKCWltYWdlX3'
    'VybBgCIAEoCVIIaW1hZ2VVcmwSGwoJbWVhbF90eXBlGAMgASgJUghtZWFsVHlwZQ==');

@$core.Deprecated('Use foodAnalysisResultDescriptor instead')
const FoodAnalysisResult$json = {
  '1': 'FoodAnalysisResult',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'food_items',
      '3': 3,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.FoodItem',
      '10': 'foodItems'
    },
    {'1': 'total_calories', '3': 4, '4': 1, '5': 1, '10': 'totalCalories'},
    {'1': 'total_carbs_g', '3': 5, '4': 1, '5': 1, '10': 'totalCarbsG'},
    {'1': 'total_protein_g', '3': 6, '4': 1, '5': 1, '10': 'totalProteinG'},
    {'1': 'total_fat_g', '3': 7, '4': 1, '5': 1, '10': 'totalFatG'},
    {'1': 'meal_type', '3': 8, '4': 1, '5': 9, '10': 'mealType'},
    {'1': 'image_url', '3': 9, '4': 1, '5': 9, '10': 'imageUrl'},
    {
      '1': 'analyzed_at',
      '3': 10,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'analyzedAt'
    },
  ],
};

/// Descriptor for `FoodAnalysisResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List foodAnalysisResultDescriptor = $convert.base64Decode(
    'ChJGb29kQW5hbHlzaXNSZXN1bHQSDgoCaWQYASABKAlSAmlkEhcKB3VzZXJfaWQYAiABKAlSBn'
    'VzZXJJZBI0Cgpmb29kX2l0ZW1zGAMgAygLMhUubWFucGFzaWsudjEuRm9vZEl0ZW1SCWZvb2RJ'
    'dGVtcxIlCg50b3RhbF9jYWxvcmllcxgEIAEoAVINdG90YWxDYWxvcmllcxIiCg10b3RhbF9jYX'
    'Jic19nGAUgASgBUgt0b3RhbENhcmJzRxImCg90b3RhbF9wcm90ZWluX2cYBiABKAFSDXRvdGFs'
    'UHJvdGVpbkcSHgoLdG90YWxfZmF0X2cYByABKAFSCXRvdGFsRmF0RxIbCgltZWFsX3R5cGUYCC'
    'ABKAlSCG1lYWxUeXBlEhsKCWltYWdlX3VybBgJIAEoCVIIaW1hZ2VVcmwSOwoLYW5hbHl6ZWRf'
    'YXQYCiABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUgphbmFseXplZEF0');

@$core.Deprecated('Use foodItemDescriptor instead')
const FoodItem$json = {
  '1': 'FoodItem',
  '2': [
    {'1': 'name', '3': 1, '4': 1, '5': 9, '10': 'name'},
    {'1': 'calories', '3': 2, '4': 1, '5': 1, '10': 'calories'},
    {'1': 'portion_g', '3': 3, '4': 1, '5': 1, '10': 'portionG'},
    {'1': 'carbs_g', '3': 4, '4': 1, '5': 1, '10': 'carbsG'},
    {'1': 'protein_g', '3': 5, '4': 1, '5': 1, '10': 'proteinG'},
    {'1': 'fat_g', '3': 6, '4': 1, '5': 1, '10': 'fatG'},
    {'1': 'confidence', '3': 7, '4': 1, '5': 1, '10': 'confidence'},
  ],
};

/// Descriptor for `FoodItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List foodItemDescriptor = $convert.base64Decode(
    'CghGb29kSXRlbRISCgRuYW1lGAEgASgJUgRuYW1lEhoKCGNhbG9yaWVzGAIgASgBUghjYWxvcm'
    'llcxIbCglwb3J0aW9uX2cYAyABKAFSCHBvcnRpb25HEhcKB2NhcmJzX2cYBCABKAFSBmNhcmJz'
    'RxIbCglwcm90ZWluX2cYBSABKAFSCHByb3RlaW5HEhMKBWZhdF9nGAYgASgBUgRmYXRHEh4KCm'
    'NvbmZpZGVuY2UYByABKAFSCmNvbmZpZGVuY2U=');

@$core.Deprecated('Use getFoodAnalysisRequestDescriptor instead')
const GetFoodAnalysisRequest$json = {
  '1': 'GetFoodAnalysisRequest',
  '2': [
    {'1': 'analysis_id', '3': 1, '4': 1, '5': 9, '10': 'analysisId'},
  ],
};

/// Descriptor for `GetFoodAnalysisRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFoodAnalysisRequestDescriptor =
    $convert.base64Decode(
        'ChZHZXRGb29kQW5hbHlzaXNSZXF1ZXN0Eh8KC2FuYWx5c2lzX2lkGAEgASgJUgphbmFseXNpc0'
        'lk');

@$core.Deprecated('Use listFoodAnalysesRequestDescriptor instead')
const ListFoodAnalysesRequest$json = {
  '1': 'ListFoodAnalysesRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'start_date', '3': 2, '4': 1, '5': 9, '10': 'startDate'},
    {'1': 'end_date', '3': 3, '4': 1, '5': 9, '10': 'endDate'},
    {'1': 'limit', '3': 4, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 5, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `ListFoodAnalysesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listFoodAnalysesRequestDescriptor = $convert.base64Decode(
    'ChdMaXN0Rm9vZEFuYWx5c2VzUmVxdWVzdBIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQSHQoKc3'
    'RhcnRfZGF0ZRgCIAEoCVIJc3RhcnREYXRlEhkKCGVuZF9kYXRlGAMgASgJUgdlbmREYXRlEhQK'
    'BWxpbWl0GAQgASgFUgVsaW1pdBIWCgZvZmZzZXQYBSABKAVSBm9mZnNldA==');

@$core.Deprecated('Use listFoodAnalysesResponseDescriptor instead')
const ListFoodAnalysesResponse$json = {
  '1': 'ListFoodAnalysesResponse',
  '2': [
    {
      '1': 'analyses',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.FoodAnalysisResult',
      '10': 'analyses'
    },
    {'1': 'total', '3': 2, '4': 1, '5': 5, '10': 'total'},
  ],
};

/// Descriptor for `ListFoodAnalysesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listFoodAnalysesResponseDescriptor = $convert.base64Decode(
    'ChhMaXN0Rm9vZEFuYWx5c2VzUmVzcG9uc2USOwoIYW5hbHlzZXMYASADKAsyHy5tYW5wYXNpay'
    '52MS5Gb29kQW5hbHlzaXNSZXN1bHRSCGFuYWx5c2VzEhQKBXRvdGFsGAIgASgFUgV0b3RhbA==');

@$core.Deprecated('Use dailyNutritionSummaryDescriptor instead')
const DailyNutritionSummary$json = {
  '1': 'DailyNutritionSummary',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'date', '3': 2, '4': 1, '5': 9, '10': 'date'},
    {'1': 'total_calories', '3': 3, '4': 1, '5': 1, '10': 'totalCalories'},
    {'1': 'total_carbs_g', '3': 4, '4': 1, '5': 1, '10': 'totalCarbsG'},
    {'1': 'total_protein_g', '3': 5, '4': 1, '5': 1, '10': 'totalProteinG'},
    {'1': 'total_fat_g', '3': 6, '4': 1, '5': 1, '10': 'totalFatG'},
    {'1': 'meal_count', '3': 7, '4': 1, '5': 5, '10': 'mealCount'},
    {
      '1': 'meals',
      '3': 8,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.MealLog',
      '10': 'meals'
    },
  ],
};

/// Descriptor for `DailyNutritionSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dailyNutritionSummaryDescriptor = $convert.base64Decode(
    'ChVEYWlseU51dHJpdGlvblN1bW1hcnkSFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklkEhIKBGRhdG'
    'UYAiABKAlSBGRhdGUSJQoOdG90YWxfY2Fsb3JpZXMYAyABKAFSDXRvdGFsQ2Fsb3JpZXMSIgoN'
    'dG90YWxfY2FyYnNfZxgEIAEoAVILdG90YWxDYXJic0cSJgoPdG90YWxfcHJvdGVpbl9nGAUgAS'
    'gBUg10b3RhbFByb3RlaW5HEh4KC3RvdGFsX2ZhdF9nGAYgASgBUgl0b3RhbEZhdEcSHQoKbWVh'
    'bF9jb3VudBgHIAEoBVIJbWVhbENvdW50EioKBW1lYWxzGAggAygLMhQubWFucGFzaWsudjEuTW'
    'VhbExvZ1IFbWVhbHM=');

@$core.Deprecated('Use getDailyNutritionSummaryRequestDescriptor instead')
const GetDailyNutritionSummaryRequest$json = {
  '1': 'GetDailyNutritionSummaryRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'date', '3': 2, '4': 1, '5': 9, '10': 'date'},
  ],
};

/// Descriptor for `GetDailyNutritionSummaryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDailyNutritionSummaryRequestDescriptor =
    $convert.base64Decode(
        'Ch9HZXREYWlseU51dHJpdGlvblN1bW1hcnlSZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZX'
        'JJZBISCgRkYXRlGAIgASgJUgRkYXRl');

@$core.Deprecated('Use logMealRequestDescriptor instead')
const LogMealRequest$json = {
  '1': 'LogMealRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'meal_type', '3': 2, '4': 1, '5': 9, '10': 'mealType'},
    {'1': 'description', '3': 3, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'items',
      '3': 4,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.FoodItem',
      '10': 'items'
    },
    {'1': 'analysis_id', '3': 5, '4': 1, '5': 9, '10': 'analysisId'},
  ],
};

/// Descriptor for `LogMealRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List logMealRequestDescriptor = $convert.base64Decode(
    'Cg5Mb2dNZWFsUmVxdWVzdBIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQSGwoJbWVhbF90eXBlGA'
    'IgASgJUghtZWFsVHlwZRIgCgtkZXNjcmlwdGlvbhgDIAEoCVILZGVzY3JpcHRpb24SKwoFaXRl'
    'bXMYBCADKAsyFS5tYW5wYXNpay52MS5Gb29kSXRlbVIFaXRlbXMSHwoLYW5hbHlzaXNfaWQYBS'
    'ABKAlSCmFuYWx5c2lzSWQ=');

@$core.Deprecated('Use mealLogDescriptor instead')
const MealLog$json = {
  '1': 'MealLog',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'meal_type', '3': 3, '4': 1, '5': 9, '10': 'mealType'},
    {'1': 'description', '3': 4, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'items',
      '3': 5,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.FoodItem',
      '10': 'items'
    },
    {'1': 'total_calories', '3': 6, '4': 1, '5': 1, '10': 'totalCalories'},
    {'1': 'analysis_id', '3': 7, '4': 1, '5': 9, '10': 'analysisId'},
    {
      '1': 'logged_at',
      '3': 8,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'loggedAt'
    },
  ],
};

/// Descriptor for `MealLog`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mealLogDescriptor = $convert.base64Decode(
    'CgdNZWFsTG9nEg4KAmlkGAEgASgJUgJpZBIXCgd1c2VyX2lkGAIgASgJUgZ1c2VySWQSGwoJbW'
    'VhbF90eXBlGAMgASgJUghtZWFsVHlwZRIgCgtkZXNjcmlwdGlvbhgEIAEoCVILZGVzY3JpcHRp'
    'b24SKwoFaXRlbXMYBSADKAsyFS5tYW5wYXNpay52MS5Gb29kSXRlbVIFaXRlbXMSJQoOdG90YW'
    'xfY2Fsb3JpZXMYBiABKAFSDXRvdGFsQ2Fsb3JpZXMSHwoLYW5hbHlzaXNfaWQYByABKAlSCmFu'
    'YWx5c2lzSWQSNwoJbG9nZ2VkX2F0GAggASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcF'
    'IIbG9nZ2VkQXQ=');

@$core.Deprecated('Use getMealHistoryRequestDescriptor instead')
const GetMealHistoryRequest$json = {
  '1': 'GetMealHistoryRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'start_date', '3': 2, '4': 1, '5': 9, '10': 'startDate'},
    {'1': 'end_date', '3': 3, '4': 1, '5': 9, '10': 'endDate'},
    {'1': 'limit', '3': 4, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 5, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `GetMealHistoryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMealHistoryRequestDescriptor = $convert.base64Decode(
    'ChVHZXRNZWFsSGlzdG9yeVJlcXVlc3QSFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklkEh0KCnN0YX'
    'J0X2RhdGUYAiABKAlSCXN0YXJ0RGF0ZRIZCghlbmRfZGF0ZRgDIAEoCVIHZW5kRGF0ZRIUCgVs'
    'aW1pdBgEIAEoBVIFbGltaXQSFgoGb2Zmc2V0GAUgASgFUgZvZmZzZXQ=');

@$core.Deprecated('Use getMealHistoryResponseDescriptor instead')
const GetMealHistoryResponse$json = {
  '1': 'GetMealHistoryResponse',
  '2': [
    {
      '1': 'meals',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.MealLog',
      '10': 'meals'
    },
    {'1': 'total', '3': 2, '4': 1, '5': 5, '10': 'total'},
  ],
};

/// Descriptor for `GetMealHistoryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMealHistoryResponseDescriptor =
    $convert.base64Decode(
        'ChZHZXRNZWFsSGlzdG9yeVJlc3BvbnNlEioKBW1lYWxzGAEgAygLMhQubWFucGFzaWsudjEuTW'
        'VhbExvZ1IFbWVhbHMSFAoFdG90YWwYAiABKAVSBXRvdGFs');

@$core.Deprecated('Use conceptDescriptor instead')
const Concept$json = {
  '1': 'Concept',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'name', '3': 2, '4': 1, '5': 9, '10': 'name'},
    {'1': 'description', '3': 3, '4': 1, '5': 9, '10': 'description'},
    {'1': 'category', '3': 4, '4': 1, '5': 9, '10': 'category'},
    {'1': 'icon_url', '3': 5, '4': 1, '5': 9, '10': 'iconUrl'},
    {'1': 'owner_id', '3': 6, '4': 1, '5': 9, '10': 'ownerId'},
    {'1': 'device_ids', '3': 7, '4': 3, '5': 9, '10': 'deviceIds'},
    {
      '1': 'created_at',
      '3': 8,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
  ],
};

/// Descriptor for `Concept`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List conceptDescriptor = $convert.base64Decode(
    'CgdDb25jZXB0Eg4KAmlkGAEgASgJUgJpZBISCgRuYW1lGAIgASgJUgRuYW1lEiAKC2Rlc2NyaX'
    'B0aW9uGAMgASgJUgtkZXNjcmlwdGlvbhIaCghjYXRlZ29yeRgEIAEoCVIIY2F0ZWdvcnkSGQoI'
    'aWNvbl91cmwYBSABKAlSB2ljb25VcmwSGQoIb3duZXJfaWQYBiABKAlSB293bmVySWQSHQoKZG'
    'V2aWNlX2lkcxgHIAMoCVIJZGV2aWNlSWRzEjkKCmNyZWF0ZWRfYXQYCCABKAsyGi5nb29nbGUu'
    'cHJvdG9idWYuVGltZXN0YW1wUgljcmVhdGVkQXQ=');

@$core.Deprecated('Use listConceptsRequestDescriptor instead')
const ListConceptsRequest$json = {
  '1': 'ListConceptsRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'category', '3': 2, '4': 1, '5': 9, '10': 'category'},
  ],
};

/// Descriptor for `ListConceptsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listConceptsRequestDescriptor = $convert.base64Decode(
    'ChNMaXN0Q29uY2VwdHNSZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZBIaCghjYXRlZ2'
    '9yeRgCIAEoCVIIY2F0ZWdvcnk=');

@$core.Deprecated('Use listConceptsResponseDescriptor instead')
const ListConceptsResponse$json = {
  '1': 'ListConceptsResponse',
  '2': [
    {
      '1': 'concepts',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.Concept',
      '10': 'concepts'
    },
  ],
};

/// Descriptor for `ListConceptsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listConceptsResponseDescriptor = $convert.base64Decode(
    'ChRMaXN0Q29uY2VwdHNSZXNwb25zZRIwCghjb25jZXB0cxgBIAMoCzIULm1hbnBhc2lrLnYxLk'
    'NvbmNlcHRSCGNvbmNlcHRz');

@$core.Deprecated('Use getConceptRequestDescriptor instead')
const GetConceptRequest$json = {
  '1': 'GetConceptRequest',
  '2': [
    {'1': 'concept_id', '3': 1, '4': 1, '5': 9, '10': 'conceptId'},
  ],
};

/// Descriptor for `GetConceptRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getConceptRequestDescriptor = $convert.base64Decode(
    'ChFHZXRDb25jZXB0UmVxdWVzdBIdCgpjb25jZXB0X2lkGAEgASgJUgljb25jZXB0SWQ=');

@$core.Deprecated('Use createConceptRequestDescriptor instead')
const CreateConceptRequest$json = {
  '1': 'CreateConceptRequest',
  '2': [
    {'1': 'name', '3': 1, '4': 1, '5': 9, '10': 'name'},
    {'1': 'description', '3': 2, '4': 1, '5': 9, '10': 'description'},
    {'1': 'category', '3': 3, '4': 1, '5': 9, '10': 'category'},
    {'1': 'owner_id', '3': 4, '4': 1, '5': 9, '10': 'ownerId'},
  ],
};

/// Descriptor for `CreateConceptRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createConceptRequestDescriptor = $convert.base64Decode(
    'ChRDcmVhdGVDb25jZXB0UmVxdWVzdBISCgRuYW1lGAEgASgJUgRuYW1lEiAKC2Rlc2NyaXB0aW'
    '9uGAIgASgJUgtkZXNjcmlwdGlvbhIaCghjYXRlZ29yeRgDIAEoCVIIY2F0ZWdvcnkSGQoIb3du'
    'ZXJfaWQYBCABKAlSB293bmVySWQ=');

@$core.Deprecated('Use assignDeviceToConceptRequestDescriptor instead')
const AssignDeviceToConceptRequest$json = {
  '1': 'AssignDeviceToConceptRequest',
  '2': [
    {'1': 'concept_id', '3': 1, '4': 1, '5': 9, '10': 'conceptId'},
    {'1': 'device_id', '3': 2, '4': 1, '5': 9, '10': 'deviceId'},
  ],
};

/// Descriptor for `AssignDeviceToConceptRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List assignDeviceToConceptRequestDescriptor =
    $convert.base64Decode(
        'ChxBc3NpZ25EZXZpY2VUb0NvbmNlcHRSZXF1ZXN0Eh0KCmNvbmNlcHRfaWQYASABKAlSCWNvbm'
        'NlcHRJZBIbCglkZXZpY2VfaWQYAiABKAlSCGRldmljZUlk');

@$core.Deprecated('Use assignDeviceToConceptResponseDescriptor instead')
const AssignDeviceToConceptResponse$json = {
  '1': 'AssignDeviceToConceptResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
  ],
};

/// Descriptor for `AssignDeviceToConceptResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List assignDeviceToConceptResponseDescriptor =
    $convert.base64Decode(
        'Ch1Bc3NpZ25EZXZpY2VUb0NvbmNlcHRSZXNwb25zZRIYCgdzdWNjZXNzGAEgASgIUgdzdWNjZX'
        'Nz');

@$core.Deprecated('Use getConceptStatsRequestDescriptor instead')
const GetConceptStatsRequest$json = {
  '1': 'GetConceptStatsRequest',
  '2': [
    {'1': 'concept_id', '3': 1, '4': 1, '5': 9, '10': 'conceptId'},
    {'1': 'period', '3': 2, '4': 1, '5': 9, '10': 'period'},
  ],
};

/// Descriptor for `GetConceptStatsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getConceptStatsRequestDescriptor =
    $convert.base64Decode(
        'ChZHZXRDb25jZXB0U3RhdHNSZXF1ZXN0Eh0KCmNvbmNlcHRfaWQYASABKAlSCWNvbmNlcHRJZB'
        'IWCgZwZXJpb2QYAiABKAlSBnBlcmlvZA==');

@$core.Deprecated('Use conceptStatsDescriptor instead')
const ConceptStats$json = {
  '1': 'ConceptStats',
  '2': [
    {'1': 'concept_id', '3': 1, '4': 1, '5': 9, '10': 'conceptId'},
    {
      '1': 'measurement_count',
      '3': 2,
      '4': 1,
      '5': 5,
      '10': 'measurementCount'
    },
    {'1': 'device_count', '3': 3, '4': 1, '5': 5, '10': 'deviceCount'},
    {'1': 'avg_health_score', '3': 4, '4': 1, '5': 1, '10': 'avgHealthScore'},
    {'1': 'period', '3': 5, '4': 1, '5': 9, '10': 'period'},
  ],
};

/// Descriptor for `ConceptStats`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List conceptStatsDescriptor = $convert.base64Decode(
    'CgxDb25jZXB0U3RhdHMSHQoKY29uY2VwdF9pZBgBIAEoCVIJY29uY2VwdElkEisKEW1lYXN1cm'
    'VtZW50X2NvdW50GAIgASgFUhBtZWFzdXJlbWVudENvdW50EiEKDGRldmljZV9jb3VudBgDIAEo'
    'BVILZGV2aWNlQ291bnQSKAoQYXZnX2hlYWx0aF9zY29yZRgEIAEoAVIOYXZnSGVhbHRoU2Nvcm'
    'USFgoGcGVyaW9kGAUgASgJUgZwZXJpb2Q=');

@$core.Deprecated('Use getConceptDashboardRequestDescriptor instead')
const GetConceptDashboardRequest$json = {
  '1': 'GetConceptDashboardRequest',
  '2': [
    {'1': 'concept_id', '3': 1, '4': 1, '5': 9, '10': 'conceptId'},
  ],
};

/// Descriptor for `GetConceptDashboardRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getConceptDashboardRequestDescriptor =
    $convert.base64Decode(
        'ChpHZXRDb25jZXB0RGFzaGJvYXJkUmVxdWVzdBIdCgpjb25jZXB0X2lkGAEgASgJUgljb25jZX'
        'B0SWQ=');

@$core.Deprecated('Use conceptDashboardDescriptor instead')
const ConceptDashboard$json = {
  '1': 'ConceptDashboard',
  '2': [
    {'1': 'concept_id', '3': 1, '4': 1, '5': 9, '10': 'conceptId'},
    {'1': 'concept_name', '3': 2, '4': 1, '5': 9, '10': 'conceptName'},
    {
      '1': 'stats',
      '3': 3,
      '4': 1,
      '5': 11,
      '6': '.manpasik.v1.ConceptStats',
      '10': 'stats'
    },
    {'1': 'recent_alerts', '3': 4, '4': 3, '5': 9, '10': 'recentAlerts'},
  ],
};

/// Descriptor for `ConceptDashboard`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List conceptDashboardDescriptor = $convert.base64Decode(
    'ChBDb25jZXB0RGFzaGJvYXJkEh0KCmNvbmNlcHRfaWQYASABKAlSCWNvbmNlcHRJZBIhCgxjb2'
    '5jZXB0X25hbWUYAiABKAlSC2NvbmNlcHROYW1lEi8KBXN0YXRzGAMgASgLMhkubWFucGFzaWsu'
    'djEuQ29uY2VwdFN0YXRzUgVzdGF0cxIjCg1yZWNlbnRfYWxlcnRzGAQgAygJUgxyZWNlbnRBbG'
    'VydHM=');

@$core.Deprecated('Use organizationDescriptor instead')
const Organization$json = {
  '1': 'Organization',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'name', '3': 2, '4': 1, '5': 9, '10': 'name'},
    {'1': 'type', '3': 3, '4': 1, '5': 9, '10': 'type'},
    {'1': 'description', '3': 4, '4': 1, '5': 9, '10': 'description'},
    {'1': 'owner_id', '3': 5, '4': 1, '5': 9, '10': 'ownerId'},
    {'1': 'member_count', '3': 6, '4': 1, '5': 5, '10': 'memberCount'},
    {
      '1': 'created_at',
      '3': 7,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
  ],
};

/// Descriptor for `Organization`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List organizationDescriptor = $convert.base64Decode(
    'CgxPcmdhbml6YXRpb24SDgoCaWQYASABKAlSAmlkEhIKBG5hbWUYAiABKAlSBG5hbWUSEgoEdH'
    'lwZRgDIAEoCVIEdHlwZRIgCgtkZXNjcmlwdGlvbhgEIAEoCVILZGVzY3JpcHRpb24SGQoIb3du'
    'ZXJfaWQYBSABKAlSB293bmVySWQSIQoMbWVtYmVyX2NvdW50GAYgASgFUgttZW1iZXJDb3VudB'
    'I5CgpjcmVhdGVkX2F0GAcgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFIJY3JlYXRl'
    'ZEF0');

@$core.Deprecated('Use createOrganizationRequestDescriptor instead')
const CreateOrganizationRequest$json = {
  '1': 'CreateOrganizationRequest',
  '2': [
    {'1': 'name', '3': 1, '4': 1, '5': 9, '10': 'name'},
    {'1': 'type', '3': 2, '4': 1, '5': 9, '10': 'type'},
    {'1': 'description', '3': 3, '4': 1, '5': 9, '10': 'description'},
    {'1': 'owner_id', '3': 4, '4': 1, '5': 9, '10': 'ownerId'},
  ],
};

/// Descriptor for `CreateOrganizationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createOrganizationRequestDescriptor = $convert.base64Decode(
    'ChlDcmVhdGVPcmdhbml6YXRpb25SZXF1ZXN0EhIKBG5hbWUYASABKAlSBG5hbWUSEgoEdHlwZR'
    'gCIAEoCVIEdHlwZRIgCgtkZXNjcmlwdGlvbhgDIAEoCVILZGVzY3JpcHRpb24SGQoIb3duZXJf'
    'aWQYBCABKAlSB293bmVySWQ=');

@$core.Deprecated('Use getOrganizationRequestDescriptor instead')
const GetOrganizationRequest$json = {
  '1': 'GetOrganizationRequest',
  '2': [
    {'1': 'organization_id', '3': 1, '4': 1, '5': 9, '10': 'organizationId'},
  ],
};

/// Descriptor for `GetOrganizationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getOrganizationRequestDescriptor =
    $convert.base64Decode(
        'ChZHZXRPcmdhbml6YXRpb25SZXF1ZXN0EicKD29yZ2FuaXphdGlvbl9pZBgBIAEoCVIOb3JnYW'
        '5pemF0aW9uSWQ=');

@$core.Deprecated('Use listOrganizationsRequestDescriptor instead')
const ListOrganizationsRequest$json = {
  '1': 'ListOrganizationsRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `ListOrganizationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listOrganizationsRequestDescriptor =
    $convert.base64Decode(
        'ChhMaXN0T3JnYW5pemF0aW9uc1JlcXVlc3QSFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklk');

@$core.Deprecated('Use listOrganizationsResponseDescriptor instead')
const ListOrganizationsResponse$json = {
  '1': 'ListOrganizationsResponse',
  '2': [
    {
      '1': 'organizations',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.Organization',
      '10': 'organizations'
    },
  ],
};

/// Descriptor for `ListOrganizationsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listOrganizationsResponseDescriptor =
    $convert.base64Decode(
        'ChlMaXN0T3JnYW5pemF0aW9uc1Jlc3BvbnNlEj8KDW9yZ2FuaXphdGlvbnMYASADKAsyGS5tYW'
        '5wYXNpay52MS5Pcmdhbml6YXRpb25SDW9yZ2FuaXphdGlvbnM=');

@$core.Deprecated('Use addOrgMemberRequestDescriptor instead')
const AddOrgMemberRequest$json = {
  '1': 'AddOrgMemberRequest',
  '2': [
    {'1': 'organization_id', '3': 1, '4': 1, '5': 9, '10': 'organizationId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'role', '3': 3, '4': 1, '5': 9, '10': 'role'},
  ],
};

/// Descriptor for `AddOrgMemberRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addOrgMemberRequestDescriptor = $convert.base64Decode(
    'ChNBZGRPcmdNZW1iZXJSZXF1ZXN0EicKD29yZ2FuaXphdGlvbl9pZBgBIAEoCVIOb3JnYW5pem'
    'F0aW9uSWQSFwoHdXNlcl9pZBgCIAEoCVIGdXNlcklkEhIKBHJvbGUYAyABKAlSBHJvbGU=');

@$core.Deprecated('Use addOrgMemberResponseDescriptor instead')
const AddOrgMemberResponse$json = {
  '1': 'AddOrgMemberResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
  ],
};

/// Descriptor for `AddOrgMemberResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addOrgMemberResponseDescriptor =
    $convert.base64Decode(
        'ChRBZGRPcmdNZW1iZXJSZXNwb25zZRIYCgdzdWNjZXNzGAEgASgIUgdzdWNjZXNz');

@$core.Deprecated('Use removeOrgMemberRequestDescriptor instead')
const RemoveOrgMemberRequest$json = {
  '1': 'RemoveOrgMemberRequest',
  '2': [
    {'1': 'organization_id', '3': 1, '4': 1, '5': 9, '10': 'organizationId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `RemoveOrgMemberRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removeOrgMemberRequestDescriptor =
    $convert.base64Decode(
        'ChZSZW1vdmVPcmdNZW1iZXJSZXF1ZXN0EicKD29yZ2FuaXphdGlvbl9pZBgBIAEoCVIOb3JnYW'
        '5pemF0aW9uSWQSFwoHdXNlcl9pZBgCIAEoCVIGdXNlcklk');

@$core.Deprecated('Use removeOrgMemberResponseDescriptor instead')
const RemoveOrgMemberResponse$json = {
  '1': 'RemoveOrgMemberResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
  ],
};

/// Descriptor for `RemoveOrgMemberResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removeOrgMemberResponseDescriptor =
    $convert.base64Decode(
        'ChdSZW1vdmVPcmdNZW1iZXJSZXNwb25zZRIYCgdzdWNjZXNzGAEgASgIUgdzdWNjZXNz');

@$core.Deprecated('Use developerProfileDescriptor instead')
const DeveloperProfile$json = {
  '1': 'DeveloperProfile',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'company_name', '3': 3, '4': 1, '5': 9, '10': 'companyName'},
    {'1': 'tier', '3': 4, '4': 1, '5': 9, '10': 'tier'},
    {'1': 'status', '3': 5, '4': 1, '5': 9, '10': 'status'},
    {'1': 'cartridge_count', '3': 6, '4': 1, '5': 5, '10': 'cartridgeCount'},
    {
      '1': 'created_at',
      '3': 7,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
  ],
};

/// Descriptor for `DeveloperProfile`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List developerProfileDescriptor = $convert.base64Decode(
    'ChBEZXZlbG9wZXJQcm9maWxlEg4KAmlkGAEgASgJUgJpZBIXCgd1c2VyX2lkGAIgASgJUgZ1c2'
    'VySWQSIQoMY29tcGFueV9uYW1lGAMgASgJUgtjb21wYW55TmFtZRISCgR0aWVyGAQgASgJUgR0'
    'aWVyEhYKBnN0YXR1cxgFIAEoCVIGc3RhdHVzEicKD2NhcnRyaWRnZV9jb3VudBgGIAEoBVIOY2'
    'FydHJpZGdlQ291bnQSOQoKY3JlYXRlZF9hdBgHIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1l'
    'c3RhbXBSCWNyZWF0ZWRBdA==');

@$core.Deprecated('Use registerDeveloperRequestDescriptor instead')
const RegisterDeveloperRequest$json = {
  '1': 'RegisterDeveloperRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'company_name', '3': 2, '4': 1, '5': 9, '10': 'companyName'},
    {'1': 'tier', '3': 3, '4': 1, '5': 9, '10': 'tier'},
  ],
};

/// Descriptor for `RegisterDeveloperRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List registerDeveloperRequestDescriptor =
    $convert.base64Decode(
        'ChhSZWdpc3RlckRldmVsb3BlclJlcXVlc3QSFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklkEiEKDG'
        'NvbXBhbnlfbmFtZRgCIAEoCVILY29tcGFueU5hbWUSEgoEdGllchgDIAEoCVIEdGllcg==');

@$core.Deprecated('Use getDeveloperProfileRequestDescriptor instead')
const GetDeveloperProfileRequest$json = {
  '1': 'GetDeveloperProfileRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `GetDeveloperProfileRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDeveloperProfileRequestDescriptor =
    $convert.base64Decode(
        'ChpHZXREZXZlbG9wZXJQcm9maWxlUmVxdWVzdBIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQ=');

@$core.Deprecated('Use createApiKeyRequestDescriptor instead')
const CreateApiKeyRequest$json = {
  '1': 'CreateApiKeyRequest',
  '2': [
    {'1': 'developer_id', '3': 1, '4': 1, '5': 9, '10': 'developerId'},
    {'1': 'name', '3': 2, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `CreateApiKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createApiKeyRequestDescriptor = $convert.base64Decode(
    'ChNDcmVhdGVBcGlLZXlSZXF1ZXN0EiEKDGRldmVsb3Blcl9pZBgBIAEoCVILZGV2ZWxvcGVySW'
    'QSEgoEbmFtZRgCIAEoCVIEbmFtZQ==');

@$core.Deprecated('Use apiKeyResponseDescriptor instead')
const ApiKeyResponse$json = {
  '1': 'ApiKeyResponse',
  '2': [
    {'1': 'key_id', '3': 1, '4': 1, '5': 9, '10': 'keyId'},
    {'1': 'api_key', '3': 2, '4': 1, '5': 9, '10': 'apiKey'},
    {'1': 'name', '3': 3, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'created_at',
      '3': 4,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
  ],
};

/// Descriptor for `ApiKeyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List apiKeyResponseDescriptor = $convert.base64Decode(
    'Cg5BcGlLZXlSZXNwb25zZRIVCgZrZXlfaWQYASABKAlSBWtleUlkEhcKB2FwaV9rZXkYAiABKA'
    'lSBmFwaUtleRISCgRuYW1lGAMgASgJUgRuYW1lEjkKCmNyZWF0ZWRfYXQYBCABKAsyGi5nb29n'
    'bGUucHJvdG9idWYuVGltZXN0YW1wUgljcmVhdGVkQXQ=');

@$core.Deprecated('Use submitCartridgeRequestDescriptor instead')
const SubmitCartridgeRequest$json = {
  '1': 'SubmitCartridgeRequest',
  '2': [
    {'1': 'developer_id', '3': 1, '4': 1, '5': 9, '10': 'developerId'},
    {'1': 'cartridge_name', '3': 2, '4': 1, '5': 9, '10': 'cartridgeName'},
    {'1': 'category', '3': 3, '4': 1, '5': 9, '10': 'category'},
    {'1': 'description', '3': 4, '4': 1, '5': 9, '10': 'description'},
    {'1': 'price_krw', '3': 5, '4': 1, '5': 5, '10': 'priceKrw'},
    {'1': 'package_url', '3': 6, '4': 1, '5': 9, '10': 'packageUrl'},
  ],
};

/// Descriptor for `SubmitCartridgeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List submitCartridgeRequestDescriptor = $convert.base64Decode(
    'ChZTdWJtaXRDYXJ0cmlkZ2VSZXF1ZXN0EiEKDGRldmVsb3Blcl9pZBgBIAEoCVILZGV2ZWxvcG'
    'VySWQSJQoOY2FydHJpZGdlX25hbWUYAiABKAlSDWNhcnRyaWRnZU5hbWUSGgoIY2F0ZWdvcnkY'
    'AyABKAlSCGNhdGVnb3J5EiAKC2Rlc2NyaXB0aW9uGAQgASgJUgtkZXNjcmlwdGlvbhIbCglwcm'
    'ljZV9rcncYBSABKAVSCHByaWNlS3J3Eh8KC3BhY2thZ2VfdXJsGAYgASgJUgpwYWNrYWdlVXJs');

@$core.Deprecated('Use submitCartridgeResponseDescriptor instead')
const SubmitCartridgeResponse$json = {
  '1': 'SubmitCartridgeResponse',
  '2': [
    {'1': 'submission_id', '3': 1, '4': 1, '5': 9, '10': 'submissionId'},
    {'1': 'status', '3': 2, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `SubmitCartridgeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List submitCartridgeResponseDescriptor =
    $convert.base64Decode(
        'ChdTdWJtaXRDYXJ0cmlkZ2VSZXNwb25zZRIjCg1zdWJtaXNzaW9uX2lkGAEgASgJUgxzdWJtaX'
        'NzaW9uSWQSFgoGc3RhdHVzGAIgASgJUgZzdGF0dXM=');

@$core.Deprecated('Use listSubmissionsRequestDescriptor instead')
const ListSubmissionsRequest$json = {
  '1': 'ListSubmissionsRequest',
  '2': [
    {'1': 'developer_id', '3': 1, '4': 1, '5': 9, '10': 'developerId'},
    {'1': 'limit', '3': 2, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 3, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `ListSubmissionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSubmissionsRequestDescriptor =
    $convert.base64Decode(
        'ChZMaXN0U3VibWlzc2lvbnNSZXF1ZXN0EiEKDGRldmVsb3Blcl9pZBgBIAEoCVILZGV2ZWxvcG'
        'VySWQSFAoFbGltaXQYAiABKAVSBWxpbWl0EhYKBm9mZnNldBgDIAEoBVIGb2Zmc2V0');

@$core.Deprecated('Use listSubmissionsResponseDescriptor instead')
const ListSubmissionsResponse$json = {
  '1': 'ListSubmissionsResponse',
  '2': [
    {
      '1': 'submissions',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.CartridgeSubmission',
      '10': 'submissions'
    },
    {'1': 'total', '3': 2, '4': 1, '5': 5, '10': 'total'},
  ],
};

/// Descriptor for `ListSubmissionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSubmissionsResponseDescriptor = $convert.base64Decode(
    'ChdMaXN0U3VibWlzc2lvbnNSZXNwb25zZRJCCgtzdWJtaXNzaW9ucxgBIAMoCzIgLm1hbnBhc2'
    'lrLnYxLkNhcnRyaWRnZVN1Ym1pc3Npb25SC3N1Ym1pc3Npb25zEhQKBXRvdGFsGAIgASgFUgV0'
    'b3RhbA==');

@$core.Deprecated('Use cartridgeSubmissionDescriptor instead')
const CartridgeSubmission$json = {
  '1': 'CartridgeSubmission',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'cartridge_name', '3': 2, '4': 1, '5': 9, '10': 'cartridgeName'},
    {'1': 'status', '3': 3, '4': 1, '5': 9, '10': 'status'},
    {'1': 'review_comment', '3': 4, '4': 1, '5': 9, '10': 'reviewComment'},
    {
      '1': 'submitted_at',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'submittedAt'
    },
  ],
};

/// Descriptor for `CartridgeSubmission`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cartridgeSubmissionDescriptor = $convert.base64Decode(
    'ChNDYXJ0cmlkZ2VTdWJtaXNzaW9uEg4KAmlkGAEgASgJUgJpZBIlCg5jYXJ0cmlkZ2VfbmFtZR'
    'gCIAEoCVINY2FydHJpZGdlTmFtZRIWCgZzdGF0dXMYAyABKAlSBnN0YXR1cxIlCg5yZXZpZXdf'
    'Y29tbWVudBgEIAEoCVINcmV2aWV3Q29tbWVudBI9CgxzdWJtaXR0ZWRfYXQYBSABKAsyGi5nb2'
    '9nbGUucHJvdG9idWYuVGltZXN0YW1wUgtzdWJtaXR0ZWRBdA==');

@$core.Deprecated('Use registerCartridgeTypeRequestDescriptor instead')
const RegisterCartridgeTypeRequest$json = {
  '1': 'RegisterCartridgeTypeRequest',
  '2': [
    {'1': 'developer_id', '3': 1, '4': 1, '5': 9, '10': 'developerId'},
    {'1': 'category_code', '3': 2, '4': 1, '5': 5, '10': 'categoryCode'},
    {'1': 'type_index', '3': 3, '4': 1, '5': 5, '10': 'typeIndex'},
    {'1': 'name', '3': 4, '4': 1, '5': 9, '10': 'name'},
    {'1': 'version', '3': 5, '4': 1, '5': 9, '10': 'version'},
    {'1': 'csi_version', '3': 6, '4': 1, '5': 9, '10': 'csiVersion'},
    {'1': 'pin_config_json', '3': 7, '4': 1, '5': 9, '10': 'pinConfigJson'},
    {'1': 'afe_blocks', '3': 8, '4': 3, '5': 9, '10': 'afeBlocks'},
    {
      '1': 'measurement_modes',
      '3': 9,
      '4': 3,
      '5': 9,
      '10': 'measurementModes'
    },
    {
      '1': 'calibration_schema_json',
      '3': 10,
      '4': 1,
      '5': 9,
      '10': 'calibrationSchemaJson'
    },
    {'1': 'alpha_default', '3': 11, '4': 1, '5': 1, '10': 'alphaDefault'},
    {
      '1': 'temp_compensation',
      '3': 12,
      '4': 1,
      '5': 8,
      '10': 'tempCompensation'
    },
    {
      '1': 'humidity_compensation',
      '3': 13,
      '4': 1,
      '5': 8,
      '10': 'humidityCompensation'
    },
    {'1': 'regulatory_class', '3': 14, '4': 1, '5': 9, '10': 'regulatoryClass'},
    {'1': 'regulatory_path', '3': 15, '4': 1, '5': 9, '10': 'regulatoryPath'},
  ],
};

/// Descriptor for `RegisterCartridgeTypeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List registerCartridgeTypeRequestDescriptor = $convert.base64Decode(
    'ChxSZWdpc3RlckNhcnRyaWRnZVR5cGVSZXF1ZXN0EiEKDGRldmVsb3Blcl9pZBgBIAEoCVILZG'
    'V2ZWxvcGVySWQSIwoNY2F0ZWdvcnlfY29kZRgCIAEoBVIMY2F0ZWdvcnlDb2RlEh0KCnR5cGVf'
    'aW5kZXgYAyABKAVSCXR5cGVJbmRleBISCgRuYW1lGAQgASgJUgRuYW1lEhgKB3ZlcnNpb24YBS'
    'ABKAlSB3ZlcnNpb24SHwoLY3NpX3ZlcnNpb24YBiABKAlSCmNzaVZlcnNpb24SJgoPcGluX2Nv'
    'bmZpZ19qc29uGAcgASgJUg1waW5Db25maWdKc29uEh0KCmFmZV9ibG9ja3MYCCADKAlSCWFmZU'
    'Jsb2NrcxIrChFtZWFzdXJlbWVudF9tb2RlcxgJIAMoCVIQbWVhc3VyZW1lbnRNb2RlcxI2Chdj'
    'YWxpYnJhdGlvbl9zY2hlbWFfanNvbhgKIAEoCVIVY2FsaWJyYXRpb25TY2hlbWFKc29uEiMKDW'
    'FscGhhX2RlZmF1bHQYCyABKAFSDGFscGhhRGVmYXVsdBIrChF0ZW1wX2NvbXBlbnNhdGlvbhgM'
    'IAEoCFIQdGVtcENvbXBlbnNhdGlvbhIzChVodW1pZGl0eV9jb21wZW5zYXRpb24YDSABKAhSFG'
    'h1bWlkaXR5Q29tcGVuc2F0aW9uEikKEHJlZ3VsYXRvcnlfY2xhc3MYDiABKAlSD3JlZ3VsYXRv'
    'cnlDbGFzcxInCg9yZWd1bGF0b3J5X3BhdGgYDyABKAlSDnJlZ3VsYXRvcnlQYXRo');

@$core.Deprecated('Use registerCartridgeTypeResponseDescriptor instead')
const RegisterCartridgeTypeResponse$json = {
  '1': 'RegisterCartridgeTypeResponse',
  '2': [
    {'1': 'definition_id', '3': 1, '4': 1, '5': 9, '10': 'definitionId'},
    {'1': 'status', '3': 2, '4': 1, '5': 9, '10': 'status'},
    {
      '1': 'csi_compatibility',
      '3': 3,
      '4': 1,
      '5': 9,
      '10': 'csiCompatibility'
    },
    {
      '1': 'validation_warnings',
      '3': 4,
      '4': 3,
      '5': 9,
      '10': 'validationWarnings'
    },
    {
      '1': 'created_at',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
  ],
};

/// Descriptor for `RegisterCartridgeTypeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List registerCartridgeTypeResponseDescriptor = $convert.base64Decode(
    'Ch1SZWdpc3RlckNhcnRyaWRnZVR5cGVSZXNwb25zZRIjCg1kZWZpbml0aW9uX2lkGAEgASgJUg'
    'xkZWZpbml0aW9uSWQSFgoGc3RhdHVzGAIgASgJUgZzdGF0dXMSKwoRY3NpX2NvbXBhdGliaWxp'
    'dHkYAyABKAlSEGNzaUNvbXBhdGliaWxpdHkSLwoTdmFsaWRhdGlvbl93YXJuaW5ncxgEIAMoCV'
    'ISdmFsaWRhdGlvbldhcm5pbmdzEjkKCmNyZWF0ZWRfYXQYBSABKAsyGi5nb29nbGUucHJvdG9i'
    'dWYuVGltZXN0YW1wUgljcmVhdGVkQXQ=');

@$core.Deprecated('Use getCartridgeDefinitionRequestDescriptor instead')
const GetCartridgeDefinitionRequest$json = {
  '1': 'GetCartridgeDefinitionRequest',
  '2': [
    {'1': 'definition_id', '3': 1, '4': 1, '5': 9, '10': 'definitionId'},
  ],
};

/// Descriptor for `GetCartridgeDefinitionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCartridgeDefinitionRequestDescriptor =
    $convert.base64Decode(
        'Ch1HZXRDYXJ0cmlkZ2VEZWZpbml0aW9uUmVxdWVzdBIjCg1kZWZpbml0aW9uX2lkGAEgASgJUg'
        'xkZWZpbml0aW9uSWQ=');

@$core.Deprecated('Use cartridgeDefinitionDetailDescriptor instead')
const CartridgeDefinitionDetail$json = {
  '1': 'CartridgeDefinitionDetail',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'developer_id', '3': 2, '4': 1, '5': 9, '10': 'developerId'},
    {'1': 'category_code', '3': 3, '4': 1, '5': 5, '10': 'categoryCode'},
    {'1': 'type_index', '3': 4, '4': 1, '5': 5, '10': 'typeIndex'},
    {'1': 'name', '3': 5, '4': 1, '5': 9, '10': 'name'},
    {'1': 'version', '3': 6, '4': 1, '5': 9, '10': 'version'},
    {'1': 'csi_version', '3': 7, '4': 1, '5': 9, '10': 'csiVersion'},
    {'1': 'pin_config_json', '3': 8, '4': 1, '5': 9, '10': 'pinConfigJson'},
    {'1': 'afe_blocks', '3': 9, '4': 3, '5': 9, '10': 'afeBlocks'},
    {
      '1': 'measurement_modes',
      '3': 10,
      '4': 3,
      '5': 9,
      '10': 'measurementModes'
    },
    {
      '1': 'calibration_schema_json',
      '3': 11,
      '4': 1,
      '5': 9,
      '10': 'calibrationSchemaJson'
    },
    {'1': 'alpha_default', '3': 12, '4': 1, '5': 1, '10': 'alphaDefault'},
    {
      '1': 'temp_compensation',
      '3': 13,
      '4': 1,
      '5': 8,
      '10': 'tempCompensation'
    },
    {
      '1': 'humidity_compensation',
      '3': 14,
      '4': 1,
      '5': 8,
      '10': 'humidityCompensation'
    },
    {'1': 'regulatory_class', '3': 15, '4': 1, '5': 9, '10': 'regulatoryClass'},
    {'1': 'regulatory_path', '3': 16, '4': 1, '5': 9, '10': 'regulatoryPath'},
    {'1': 'approval_status', '3': 17, '4': 1, '5': 9, '10': 'approvalStatus'},
    {
      '1': 'created_at',
      '3': 18,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
    {
      '1': 'approved_at',
      '3': 19,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'approvedAt'
    },
  ],
};

/// Descriptor for `CartridgeDefinitionDetail`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cartridgeDefinitionDetailDescriptor = $convert.base64Decode(
    'ChlDYXJ0cmlkZ2VEZWZpbml0aW9uRGV0YWlsEg4KAmlkGAEgASgJUgJpZBIhCgxkZXZlbG9wZX'
    'JfaWQYAiABKAlSC2RldmVsb3BlcklkEiMKDWNhdGVnb3J5X2NvZGUYAyABKAVSDGNhdGVnb3J5'
    'Q29kZRIdCgp0eXBlX2luZGV4GAQgASgFUgl0eXBlSW5kZXgSEgoEbmFtZRgFIAEoCVIEbmFtZR'
    'IYCgd2ZXJzaW9uGAYgASgJUgd2ZXJzaW9uEh8KC2NzaV92ZXJzaW9uGAcgASgJUgpjc2lWZXJz'
    'aW9uEiYKD3Bpbl9jb25maWdfanNvbhgIIAEoCVINcGluQ29uZmlnSnNvbhIdCgphZmVfYmxvY2'
    'tzGAkgAygJUglhZmVCbG9ja3MSKwoRbWVhc3VyZW1lbnRfbW9kZXMYCiADKAlSEG1lYXN1cmVt'
    'ZW50TW9kZXMSNgoXY2FsaWJyYXRpb25fc2NoZW1hX2pzb24YCyABKAlSFWNhbGlicmF0aW9uU2'
    'NoZW1hSnNvbhIjCg1hbHBoYV9kZWZhdWx0GAwgASgBUgxhbHBoYURlZmF1bHQSKwoRdGVtcF9j'
    'b21wZW5zYXRpb24YDSABKAhSEHRlbXBDb21wZW5zYXRpb24SMwoVaHVtaWRpdHlfY29tcGVuc2'
    'F0aW9uGA4gASgIUhRodW1pZGl0eUNvbXBlbnNhdGlvbhIpChByZWd1bGF0b3J5X2NsYXNzGA8g'
    'ASgJUg9yZWd1bGF0b3J5Q2xhc3MSJwoPcmVndWxhdG9yeV9wYXRoGBAgASgJUg5yZWd1bGF0b3'
    'J5UGF0aBInCg9hcHByb3ZhbF9zdGF0dXMYESABKAlSDmFwcHJvdmFsU3RhdHVzEjkKCmNyZWF0'
    'ZWRfYXQYEiABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUgljcmVhdGVkQXQSOwoLYX'
    'Bwcm92ZWRfYXQYEyABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUgphcHByb3ZlZEF0');

@$core.Deprecated('Use storeItemDescriptor instead')
const StoreItem$json = {
  '1': 'StoreItem',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'cartridge_name', '3': 2, '4': 1, '5': 9, '10': 'cartridgeName'},
    {'1': 'developer_name', '3': 3, '4': 1, '5': 9, '10': 'developerName'},
    {'1': 'category', '3': 4, '4': 1, '5': 9, '10': 'category'},
    {'1': 'description', '3': 5, '4': 1, '5': 9, '10': 'description'},
    {'1': 'price_krw', '3': 6, '4': 1, '5': 5, '10': 'priceKrw'},
    {'1': 'rating', '3': 7, '4': 1, '5': 1, '10': 'rating'},
    {'1': 'review_count', '3': 8, '4': 1, '5': 5, '10': 'reviewCount'},
    {'1': 'download_count', '3': 9, '4': 1, '5': 5, '10': 'downloadCount'},
    {'1': 'icon_url', '3': 10, '4': 1, '5': 9, '10': 'iconUrl'},
    {
      '1': 'published_at',
      '3': 11,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'publishedAt'
    },
  ],
};

/// Descriptor for `StoreItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List storeItemDescriptor = $convert.base64Decode(
    'CglTdG9yZUl0ZW0SDgoCaWQYASABKAlSAmlkEiUKDmNhcnRyaWRnZV9uYW1lGAIgASgJUg1jYX'
    'J0cmlkZ2VOYW1lEiUKDmRldmVsb3Blcl9uYW1lGAMgASgJUg1kZXZlbG9wZXJOYW1lEhoKCGNh'
    'dGVnb3J5GAQgASgJUghjYXRlZ29yeRIgCgtkZXNjcmlwdGlvbhgFIAEoCVILZGVzY3JpcHRpb2'
    '4SGwoJcHJpY2Vfa3J3GAYgASgFUghwcmljZUtydxIWCgZyYXRpbmcYByABKAFSBnJhdGluZxIh'
    'CgxyZXZpZXdfY291bnQYCCABKAVSC3Jldmlld0NvdW50EiUKDmRvd25sb2FkX2NvdW50GAkgAS'
    'gFUg1kb3dubG9hZENvdW50EhkKCGljb25fdXJsGAogASgJUgdpY29uVXJsEj0KDHB1Ymxpc2hl'
    'ZF9hdBgLIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSC3B1Ymxpc2hlZEF0');

@$core.Deprecated('Use listStoreItemsRequestDescriptor instead')
const ListStoreItemsRequest$json = {
  '1': 'ListStoreItemsRequest',
  '2': [
    {'1': 'category', '3': 1, '4': 1, '5': 9, '10': 'category'},
    {'1': 'sort_by', '3': 2, '4': 1, '5': 9, '10': 'sortBy'},
    {'1': 'limit', '3': 3, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 4, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `ListStoreItemsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listStoreItemsRequestDescriptor = $convert.base64Decode(
    'ChVMaXN0U3RvcmVJdGVtc1JlcXVlc3QSGgoIY2F0ZWdvcnkYASABKAlSCGNhdGVnb3J5EhcKB3'
    'NvcnRfYnkYAiABKAlSBnNvcnRCeRIUCgVsaW1pdBgDIAEoBVIFbGltaXQSFgoGb2Zmc2V0GAQg'
    'ASgFUgZvZmZzZXQ=');

@$core.Deprecated('Use listStoreItemsResponseDescriptor instead')
const ListStoreItemsResponse$json = {
  '1': 'ListStoreItemsResponse',
  '2': [
    {
      '1': 'items',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.StoreItem',
      '10': 'items'
    },
    {'1': 'total', '3': 2, '4': 1, '5': 5, '10': 'total'},
  ],
};

/// Descriptor for `ListStoreItemsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listStoreItemsResponseDescriptor =
    $convert.base64Decode(
        'ChZMaXN0U3RvcmVJdGVtc1Jlc3BvbnNlEiwKBWl0ZW1zGAEgAygLMhYubWFucGFzaWsudjEuU3'
        'RvcmVJdGVtUgVpdGVtcxIUCgV0b3RhbBgCIAEoBVIFdG90YWw=');

@$core.Deprecated('Use searchCartridgesRequestDescriptor instead')
const SearchCartridgesRequest$json = {
  '1': 'SearchCartridgesRequest',
  '2': [
    {'1': 'query', '3': 1, '4': 1, '5': 9, '10': 'query'},
    {'1': 'category', '3': 2, '4': 1, '5': 9, '10': 'category'},
    {'1': 'limit', '3': 3, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 4, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `SearchCartridgesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List searchCartridgesRequestDescriptor = $convert.base64Decode(
    'ChdTZWFyY2hDYXJ0cmlkZ2VzUmVxdWVzdBIUCgVxdWVyeRgBIAEoCVIFcXVlcnkSGgoIY2F0ZW'
    'dvcnkYAiABKAlSCGNhdGVnb3J5EhQKBWxpbWl0GAMgASgFUgVsaW1pdBIWCgZvZmZzZXQYBCAB'
    'KAVSBm9mZnNldA==');

@$core.Deprecated('Use searchCartridgesResponseDescriptor instead')
const SearchCartridgesResponse$json = {
  '1': 'SearchCartridgesResponse',
  '2': [
    {
      '1': 'results',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.StoreItem',
      '10': 'results'
    },
    {'1': 'total', '3': 2, '4': 1, '5': 5, '10': 'total'},
  ],
};

/// Descriptor for `SearchCartridgesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List searchCartridgesResponseDescriptor =
    $convert.base64Decode(
        'ChhTZWFyY2hDYXJ0cmlkZ2VzUmVzcG9uc2USMAoHcmVzdWx0cxgBIAMoCzIWLm1hbnBhc2lrLn'
        'YxLlN0b3JlSXRlbVIHcmVzdWx0cxIUCgV0b3RhbBgCIAEoBVIFdG90YWw=');

@$core.Deprecated('Use getStoreItemRequestDescriptor instead')
const GetStoreItemRequest$json = {
  '1': 'GetStoreItemRequest',
  '2': [
    {'1': 'item_id', '3': 1, '4': 1, '5': 9, '10': 'itemId'},
  ],
};

/// Descriptor for `GetStoreItemRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getStoreItemRequestDescriptor =
    $convert.base64Decode(
        'ChNHZXRTdG9yZUl0ZW1SZXF1ZXN0EhcKB2l0ZW1faWQYASABKAlSBml0ZW1JZA==');

@$core.Deprecated('Use purchaseCartridgeRequestDescriptor instead')
const PurchaseCartridgeRequest$json = {
  '1': 'PurchaseCartridgeRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'item_id', '3': 2, '4': 1, '5': 9, '10': 'itemId'},
    {'1': 'payment_method', '3': 3, '4': 1, '5': 9, '10': 'paymentMethod'},
  ],
};

/// Descriptor for `PurchaseCartridgeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List purchaseCartridgeRequestDescriptor = $convert.base64Decode(
    'ChhQdXJjaGFzZUNhcnRyaWRnZVJlcXVlc3QSFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklkEhcKB2'
    'l0ZW1faWQYAiABKAlSBml0ZW1JZBIlCg5wYXltZW50X21ldGhvZBgDIAEoCVINcGF5bWVudE1l'
    'dGhvZA==');

@$core.Deprecated('Use purchaseCartridgeResponseDescriptor instead')
const PurchaseCartridgeResponse$json = {
  '1': 'PurchaseCartridgeResponse',
  '2': [
    {'1': 'purchase_id', '3': 1, '4': 1, '5': 9, '10': 'purchaseId'},
    {'1': 'success', '3': 2, '4': 1, '5': 8, '10': 'success'},
    {'1': 'message', '3': 3, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `PurchaseCartridgeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List purchaseCartridgeResponseDescriptor = $convert.base64Decode(
    'ChlQdXJjaGFzZUNhcnRyaWRnZVJlc3BvbnNlEh8KC3B1cmNoYXNlX2lkGAEgASgJUgpwdXJjaG'
    'FzZUlkEhgKB3N1Y2Nlc3MYAiABKAhSB3N1Y2Nlc3MSGAoHbWVzc2FnZRgDIAEoCVIHbWVzc2Fn'
    'ZQ==');

@$core.Deprecated('Use getPurchaseHistoryRequestDescriptor instead')
const GetPurchaseHistoryRequest$json = {
  '1': 'GetPurchaseHistoryRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'limit', '3': 2, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 3, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `GetPurchaseHistoryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPurchaseHistoryRequestDescriptor =
    $convert.base64Decode(
        'ChlHZXRQdXJjaGFzZUhpc3RvcnlSZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZBIUCg'
        'VsaW1pdBgCIAEoBVIFbGltaXQSFgoGb2Zmc2V0GAMgASgFUgZvZmZzZXQ=');

@$core.Deprecated('Use getPurchaseHistoryResponseDescriptor instead')
const GetPurchaseHistoryResponse$json = {
  '1': 'GetPurchaseHistoryResponse',
  '2': [
    {
      '1': 'purchases',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.CartridgePurchase',
      '10': 'purchases'
    },
    {'1': 'total', '3': 2, '4': 1, '5': 5, '10': 'total'},
  ],
};

/// Descriptor for `GetPurchaseHistoryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPurchaseHistoryResponseDescriptor =
    $convert.base64Decode(
        'ChpHZXRQdXJjaGFzZUhpc3RvcnlSZXNwb25zZRI8CglwdXJjaGFzZXMYASADKAsyHi5tYW5wYX'
        'Npay52MS5DYXJ0cmlkZ2VQdXJjaGFzZVIJcHVyY2hhc2VzEhQKBXRvdGFsGAIgASgFUgV0b3Rh'
        'bA==');

@$core.Deprecated('Use cartridgePurchaseDescriptor instead')
const CartridgePurchase$json = {
  '1': 'CartridgePurchase',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'item_id', '3': 2, '4': 1, '5': 9, '10': 'itemId'},
    {'1': 'cartridge_name', '3': 3, '4': 1, '5': 9, '10': 'cartridgeName'},
    {'1': 'price_krw', '3': 4, '4': 1, '5': 5, '10': 'priceKrw'},
    {
      '1': 'purchased_at',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'purchasedAt'
    },
  ],
};

/// Descriptor for `CartridgePurchase`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cartridgePurchaseDescriptor = $convert.base64Decode(
    'ChFDYXJ0cmlkZ2VQdXJjaGFzZRIOCgJpZBgBIAEoCVICaWQSFwoHaXRlbV9pZBgCIAEoCVIGaX'
    'RlbUlkEiUKDmNhcnRyaWRnZV9uYW1lGAMgASgJUg1jYXJ0cmlkZ2VOYW1lEhsKCXByaWNlX2ty'
    'dxgEIAEoBVIIcHJpY2VLcncSPQoMcHVyY2hhc2VkX2F0GAUgASgLMhouZ29vZ2xlLnByb3RvYn'
    'VmLlRpbWVzdGFtcFILcHVyY2hhc2VkQXQ=');

@$core.Deprecated('Use reviewStatusDescriptor instead')
const ReviewStatus$json = {
  '1': 'ReviewStatus',
  '2': [
    {'1': 'submission_id', '3': 1, '4': 1, '5': 9, '10': 'submissionId'},
    {'1': 'status', '3': 2, '4': 1, '5': 9, '10': 'status'},
    {'1': 'reviewer_comment', '3': 3, '4': 1, '5': 9, '10': 'reviewerComment'},
    {'1': 'checklist_passed', '3': 4, '4': 3, '5': 9, '10': 'checklistPassed'},
    {'1': 'checklist_failed', '3': 5, '4': 3, '5': 9, '10': 'checklistFailed'},
    {
      '1': 'updated_at',
      '3': 6,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'updatedAt'
    },
  ],
};

/// Descriptor for `ReviewStatus`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List reviewStatusDescriptor = $convert.base64Decode(
    'CgxSZXZpZXdTdGF0dXMSIwoNc3VibWlzc2lvbl9pZBgBIAEoCVIMc3VibWlzc2lvbklkEhYKBn'
    'N0YXR1cxgCIAEoCVIGc3RhdHVzEikKEHJldmlld2VyX2NvbW1lbnQYAyABKAlSD3Jldmlld2Vy'
    'Q29tbWVudBIpChBjaGVja2xpc3RfcGFzc2VkGAQgAygJUg9jaGVja2xpc3RQYXNzZWQSKQoQY2'
    'hlY2tsaXN0X2ZhaWxlZBgFIAMoCVIPY2hlY2tsaXN0RmFpbGVkEjkKCnVwZGF0ZWRfYXQYBiAB'
    'KAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUgl1cGRhdGVkQXQ=');

@$core.Deprecated('Use submitForReviewRequestDescriptor instead')
const SubmitForReviewRequest$json = {
  '1': 'SubmitForReviewRequest',
  '2': [
    {'1': 'submission_id', '3': 1, '4': 1, '5': 9, '10': 'submissionId'},
  ],
};

/// Descriptor for `SubmitForReviewRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List submitForReviewRequestDescriptor =
    $convert.base64Decode(
        'ChZTdWJtaXRGb3JSZXZpZXdSZXF1ZXN0EiMKDXN1Ym1pc3Npb25faWQYASABKAlSDHN1Ym1pc3'
        'Npb25JZA==');

@$core.Deprecated('Use getReviewStatusRequestDescriptor instead')
const GetReviewStatusRequest$json = {
  '1': 'GetReviewStatusRequest',
  '2': [
    {'1': 'submission_id', '3': 1, '4': 1, '5': 9, '10': 'submissionId'},
  ],
};

/// Descriptor for `GetReviewStatusRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getReviewStatusRequestDescriptor =
    $convert.base64Decode(
        'ChZHZXRSZXZpZXdTdGF0dXNSZXF1ZXN0EiMKDXN1Ym1pc3Npb25faWQYASABKAlSDHN1Ym1pc3'
        'Npb25JZA==');

@$core.Deprecated('Use approveCartridgeRequestDescriptor instead')
const ApproveCartridgeRequest$json = {
  '1': 'ApproveCartridgeRequest',
  '2': [
    {'1': 'submission_id', '3': 1, '4': 1, '5': 9, '10': 'submissionId'},
    {'1': 'reviewer_id', '3': 2, '4': 1, '5': 9, '10': 'reviewerId'},
    {'1': 'comment', '3': 3, '4': 1, '5': 9, '10': 'comment'},
  ],
};

/// Descriptor for `ApproveCartridgeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List approveCartridgeRequestDescriptor = $convert.base64Decode(
    'ChdBcHByb3ZlQ2FydHJpZGdlUmVxdWVzdBIjCg1zdWJtaXNzaW9uX2lkGAEgASgJUgxzdWJtaX'
    'NzaW9uSWQSHwoLcmV2aWV3ZXJfaWQYAiABKAlSCnJldmlld2VySWQSGAoHY29tbWVudBgDIAEo'
    'CVIHY29tbWVudA==');

@$core.Deprecated('Use rejectCartridgeRequestDescriptor instead')
const RejectCartridgeRequest$json = {
  '1': 'RejectCartridgeRequest',
  '2': [
    {'1': 'submission_id', '3': 1, '4': 1, '5': 9, '10': 'submissionId'},
    {'1': 'reviewer_id', '3': 2, '4': 1, '5': 9, '10': 'reviewerId'},
    {'1': 'comment', '3': 3, '4': 1, '5': 9, '10': 'comment'},
    {'1': 'reasons', '3': 4, '4': 3, '5': 9, '10': 'reasons'},
  ],
};

/// Descriptor for `RejectCartridgeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List rejectCartridgeRequestDescriptor = $convert.base64Decode(
    'ChZSZWplY3RDYXJ0cmlkZ2VSZXF1ZXN0EiMKDXN1Ym1pc3Npb25faWQYASABKAlSDHN1Ym1pc3'
    'Npb25JZBIfCgtyZXZpZXdlcl9pZBgCIAEoCVIKcmV2aWV3ZXJJZBIYCgdjb21tZW50GAMgASgJ'
    'Ugdjb21tZW50EhgKB3JlYXNvbnMYBCADKAlSB3JlYXNvbnM=');

@$core.Deprecated('Use getRegionStatsRequestDescriptor instead')
const GetRegionStatsRequest$json = {
  '1': 'GetRegionStatsRequest',
  '2': [
    {'1': 'region_code', '3': 1, '4': 1, '5': 9, '10': 'regionCode'},
    {'1': 'period', '3': 2, '4': 1, '5': 9, '10': 'period'},
    {'1': 'biomarker', '3': 3, '4': 1, '5': 9, '10': 'biomarker'},
  ],
};

/// Descriptor for `GetRegionStatsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRegionStatsRequestDescriptor = $convert.base64Decode(
    'ChVHZXRSZWdpb25TdGF0c1JlcXVlc3QSHwoLcmVnaW9uX2NvZGUYASABKAlSCnJlZ2lvbkNvZG'
    'USFgoGcGVyaW9kGAIgASgJUgZwZXJpb2QSHAoJYmlvbWFya2VyGAMgASgJUgliaW9tYXJrZXI=');

@$core.Deprecated('Use regionStatsDescriptor instead')
const RegionStats$json = {
  '1': 'RegionStats',
  '2': [
    {'1': 'region_code', '3': 1, '4': 1, '5': 9, '10': 'regionCode'},
    {'1': 'region_name', '3': 2, '4': 1, '5': 9, '10': 'regionName'},
    {'1': 'period', '3': 3, '4': 1, '5': 9, '10': 'period'},
    {'1': 'avg_value', '3': 4, '4': 1, '5': 1, '10': 'avgValue'},
    {'1': 'min_value', '3': 5, '4': 1, '5': 1, '10': 'minValue'},
    {'1': 'max_value', '3': 6, '4': 1, '5': 1, '10': 'maxValue'},
    {'1': 'sample_count', '3': 7, '4': 1, '5': 5, '10': 'sampleCount'},
    {'1': 'risk_level', '3': 8, '4': 1, '5': 9, '10': 'riskLevel'},
  ],
};

/// Descriptor for `RegionStats`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List regionStatsDescriptor = $convert.base64Decode(
    'CgtSZWdpb25TdGF0cxIfCgtyZWdpb25fY29kZRgBIAEoCVIKcmVnaW9uQ29kZRIfCgtyZWdpb2'
    '5fbmFtZRgCIAEoCVIKcmVnaW9uTmFtZRIWCgZwZXJpb2QYAyABKAlSBnBlcmlvZBIbCglhdmdf'
    'dmFsdWUYBCABKAFSCGF2Z1ZhbHVlEhsKCW1pbl92YWx1ZRgFIAEoAVIIbWluVmFsdWUSGwoJbW'
    'F4X3ZhbHVlGAYgASgBUghtYXhWYWx1ZRIhCgxzYW1wbGVfY291bnQYByABKAVSC3NhbXBsZUNv'
    'dW50Eh0KCnJpc2tfbGV2ZWwYCCABKAlSCXJpc2tMZXZlbA==');

@$core.Deprecated('Use listHotspotsRequestDescriptor instead')
const ListHotspotsRequest$json = {
  '1': 'ListHotspotsRequest',
  '2': [
    {'1': 'biomarker', '3': 1, '4': 1, '5': 9, '10': 'biomarker'},
    {'1': 'risk_level', '3': 2, '4': 1, '5': 9, '10': 'riskLevel'},
    {'1': 'limit', '3': 3, '4': 1, '5': 5, '10': 'limit'},
  ],
};

/// Descriptor for `ListHotspotsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listHotspotsRequestDescriptor = $convert.base64Decode(
    'ChNMaXN0SG90c3BvdHNSZXF1ZXN0EhwKCWJpb21hcmtlchgBIAEoCVIJYmlvbWFya2VyEh0KCn'
    'Jpc2tfbGV2ZWwYAiABKAlSCXJpc2tMZXZlbBIUCgVsaW1pdBgDIAEoBVIFbGltaXQ=');

@$core.Deprecated('Use listHotspotsResponseDescriptor instead')
const ListHotspotsResponse$json = {
  '1': 'ListHotspotsResponse',
  '2': [
    {
      '1': 'hotspots',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.RegionStats',
      '10': 'hotspots'
    },
  ],
};

/// Descriptor for `ListHotspotsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listHotspotsResponseDescriptor = $convert.base64Decode(
    'ChRMaXN0SG90c3BvdHNSZXNwb25zZRI0Cghob3RzcG90cxgBIAMoCzIYLm1hbnBhc2lrLnYxLl'
    'JlZ2lvblN0YXRzUghob3RzcG90cw==');

@$core.Deprecated('Use getTrendByRegionRequestDescriptor instead')
const GetTrendByRegionRequest$json = {
  '1': 'GetTrendByRegionRequest',
  '2': [
    {'1': 'region_code', '3': 1, '4': 1, '5': 9, '10': 'regionCode'},
    {'1': 'biomarker', '3': 2, '4': 1, '5': 9, '10': 'biomarker'},
    {'1': 'start_date', '3': 3, '4': 1, '5': 9, '10': 'startDate'},
    {'1': 'end_date', '3': 4, '4': 1, '5': 9, '10': 'endDate'},
  ],
};

/// Descriptor for `GetTrendByRegionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTrendByRegionRequestDescriptor = $convert.base64Decode(
    'ChdHZXRUcmVuZEJ5UmVnaW9uUmVxdWVzdBIfCgtyZWdpb25fY29kZRgBIAEoCVIKcmVnaW9uQ2'
    '9kZRIcCgliaW9tYXJrZXIYAiABKAlSCWJpb21hcmtlchIdCgpzdGFydF9kYXRlGAMgASgJUglz'
    'dGFydERhdGUSGQoIZW5kX2RhdGUYBCABKAlSB2VuZERhdGU=');

@$core.Deprecated('Use regionTrendDescriptor instead')
const RegionTrend$json = {
  '1': 'RegionTrend',
  '2': [
    {'1': 'region_code', '3': 1, '4': 1, '5': 9, '10': 'regionCode'},
    {'1': 'biomarker', '3': 2, '4': 1, '5': 9, '10': 'biomarker'},
    {
      '1': 'points',
      '3': 3,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.RegionTrendPoint',
      '10': 'points'
    },
    {'1': 'trend_direction', '3': 4, '4': 1, '5': 9, '10': 'trendDirection'},
  ],
};

/// Descriptor for `RegionTrend`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List regionTrendDescriptor = $convert.base64Decode(
    'CgtSZWdpb25UcmVuZBIfCgtyZWdpb25fY29kZRgBIAEoCVIKcmVnaW9uQ29kZRIcCgliaW9tYX'
    'JrZXIYAiABKAlSCWJpb21hcmtlchI1CgZwb2ludHMYAyADKAsyHS5tYW5wYXNpay52MS5SZWdp'
    'b25UcmVuZFBvaW50UgZwb2ludHMSJwoPdHJlbmRfZGlyZWN0aW9uGAQgASgJUg50cmVuZERpcm'
    'VjdGlvbg==');

@$core.Deprecated('Use regionTrendPointDescriptor instead')
const RegionTrendPoint$json = {
  '1': 'RegionTrendPoint',
  '2': [
    {'1': 'date', '3': 1, '4': 1, '5': 9, '10': 'date'},
    {'1': 'value', '3': 2, '4': 1, '5': 1, '10': 'value'},
    {'1': 'sample_count', '3': 3, '4': 1, '5': 5, '10': 'sampleCount'},
  ],
};

/// Descriptor for `RegionTrendPoint`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List regionTrendPointDescriptor = $convert.base64Decode(
    'ChBSZWdpb25UcmVuZFBvaW50EhIKBGRhdGUYASABKAlSBGRhdGUSFAoFdmFsdWUYAiABKAFSBX'
    'ZhbHVlEiEKDHNhbXBsZV9jb3VudBgDIAEoBVILc2FtcGxlQ291bnQ=');

@$core.Deprecated('Use listDatasetsRequestDescriptor instead')
const ListDatasetsRequest$json = {
  '1': 'ListDatasetsRequest',
  '2': [
    {'1': 'dataset_type', '3': 1, '4': 1, '5': 9, '10': 'datasetType'},
    {'1': 'limit', '3': 2, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 3, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `ListDatasetsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDatasetsRequestDescriptor = $convert.base64Decode(
    'ChNMaXN0RGF0YXNldHNSZXF1ZXN0EiEKDGRhdGFzZXRfdHlwZRgBIAEoCVILZGF0YXNldFR5cG'
    'USFAoFbGltaXQYAiABKAVSBWxpbWl0EhYKBm9mZnNldBgDIAEoBVIGb2Zmc2V0');

@$core.Deprecated('Use listDatasetsResponseDescriptor instead')
const ListDatasetsResponse$json = {
  '1': 'ListDatasetsResponse',
  '2': [
    {
      '1': 'datasets',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.AnonymizedDataset',
      '10': 'datasets'
    },
    {'1': 'total', '3': 2, '4': 1, '5': 5, '10': 'total'},
  ],
};

/// Descriptor for `ListDatasetsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDatasetsResponseDescriptor = $convert.base64Decode(
    'ChRMaXN0RGF0YXNldHNSZXNwb25zZRI6CghkYXRhc2V0cxgBIAMoCzIeLm1hbnBhc2lrLnYxLk'
    'Fub255bWl6ZWREYXRhc2V0UghkYXRhc2V0cxIUCgV0b3RhbBgCIAEoBVIFdG90YWw=');

@$core.Deprecated('Use anonymizedDatasetDescriptor instead')
const AnonymizedDataset$json = {
  '1': 'AnonymizedDataset',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'dataset_type', '3': 2, '4': 1, '5': 9, '10': 'datasetType'},
    {'1': 'description', '3': 3, '4': 1, '5': 9, '10': 'description'},
    {'1': 'record_count', '3': 4, '4': 1, '5': 5, '10': 'recordCount'},
    {'1': 'date_range', '3': 5, '4': 1, '5': 9, '10': 'dateRange'},
    {
      '1': 'created_at',
      '3': 6,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
  ],
};

/// Descriptor for `AnonymizedDataset`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List anonymizedDatasetDescriptor = $convert.base64Decode(
    'ChFBbm9ueW1pemVkRGF0YXNldBIOCgJpZBgBIAEoCVICaWQSIQoMZGF0YXNldF90eXBlGAIgAS'
    'gJUgtkYXRhc2V0VHlwZRIgCgtkZXNjcmlwdGlvbhgDIAEoCVILZGVzY3JpcHRpb24SIQoMcmVj'
    'b3JkX2NvdW50GAQgASgFUgtyZWNvcmRDb3VudBIdCgpkYXRlX3JhbmdlGAUgASgJUglkYXRlUm'
    'FuZ2USOQoKY3JlYXRlZF9hdBgGIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBSCWNy'
    'ZWF0ZWRBdA==');

@$core.Deprecated('Use getAnonymizedDatasetRequestDescriptor instead')
const GetAnonymizedDatasetRequest$json = {
  '1': 'GetAnonymizedDatasetRequest',
  '2': [
    {'1': 'dataset_id', '3': 1, '4': 1, '5': 9, '10': 'datasetId'},
  ],
};

/// Descriptor for `GetAnonymizedDatasetRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAnonymizedDatasetRequestDescriptor =
    $convert.base64Decode(
        'ChtHZXRBbm9ueW1pemVkRGF0YXNldFJlcXVlc3QSHQoKZGF0YXNldF9pZBgBIAEoCVIJZGF0YX'
        'NldElk');

@$core.Deprecated('Use requestDataAccessRequestDescriptor instead')
const RequestDataAccessRequest$json = {
  '1': 'RequestDataAccessRequest',
  '2': [
    {'1': 'requester_id', '3': 1, '4': 1, '5': 9, '10': 'requesterId'},
    {'1': 'dataset_id', '3': 2, '4': 1, '5': 9, '10': 'datasetId'},
    {'1': 'purpose', '3': 3, '4': 1, '5': 9, '10': 'purpose'},
  ],
};

/// Descriptor for `RequestDataAccessRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List requestDataAccessRequestDescriptor = $convert.base64Decode(
    'ChhSZXF1ZXN0RGF0YUFjY2Vzc1JlcXVlc3QSIQoMcmVxdWVzdGVyX2lkGAEgASgJUgtyZXF1ZX'
    'N0ZXJJZBIdCgpkYXRhc2V0X2lkGAIgASgJUglkYXRhc2V0SWQSGAoHcHVycG9zZRgDIAEoCVIH'
    'cHVycG9zZQ==');

@$core.Deprecated('Use dataAccessResponseDescriptor instead')
const DataAccessResponse$json = {
  '1': 'DataAccessResponse',
  '2': [
    {'1': 'request_id', '3': 1, '4': 1, '5': 9, '10': 'requestId'},
    {'1': 'status', '3': 2, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `DataAccessResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dataAccessResponseDescriptor = $convert.base64Decode(
    'ChJEYXRhQWNjZXNzUmVzcG9uc2USHQoKcmVxdWVzdF9pZBgBIAEoCVIJcmVxdWVzdElkEhYKBn'
    'N0YXR1cxgCIAEoCVIGc3RhdHVz');

@$core.Deprecated('Use getConsentRequestDescriptor instead')
const GetConsentRequest$json = {
  '1': 'GetConsentRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `GetConsentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getConsentRequestDescriptor = $convert.base64Decode(
    'ChFHZXRDb25zZW50UmVxdWVzdBIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQ=');

@$core.Deprecated('Use dataSharingConsentInfoDescriptor instead')
const DataSharingConsentInfo$json = {
  '1': 'DataSharingConsentInfo',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'anonymous_research',
      '3': 2,
      '4': 1,
      '5': 8,
      '10': 'anonymousResearch'
    },
    {'1': 'region_stats', '3': 3, '4': 1, '5': 8, '10': 'regionStats'},
    {'1': 'enterprise_data', '3': 4, '4': 1, '5': 8, '10': 'enterpriseData'},
    {
      '1': 'updated_at',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'updatedAt'
    },
  ],
};

/// Descriptor for `DataSharingConsentInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dataSharingConsentInfoDescriptor = $convert.base64Decode(
    'ChZEYXRhU2hhcmluZ0NvbnNlbnRJbmZvEhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZBItChJhbm'
    '9ueW1vdXNfcmVzZWFyY2gYAiABKAhSEWFub255bW91c1Jlc2VhcmNoEiEKDHJlZ2lvbl9zdGF0'
    'cxgDIAEoCFILcmVnaW9uU3RhdHMSJwoPZW50ZXJwcmlzZV9kYXRhGAQgASgIUg5lbnRlcnByaX'
    'NlRGF0YRI5Cgp1cGRhdGVkX2F0GAUgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcFIJ'
    'dXBkYXRlZEF0');

@$core.Deprecated('Use updateConsentRequestDescriptor instead')
const UpdateConsentRequest$json = {
  '1': 'UpdateConsentRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {
      '1': 'anonymous_research',
      '3': 2,
      '4': 1,
      '5': 8,
      '10': 'anonymousResearch'
    },
    {'1': 'region_stats', '3': 3, '4': 1, '5': 8, '10': 'regionStats'},
    {'1': 'enterprise_data', '3': 4, '4': 1, '5': 8, '10': 'enterpriseData'},
  ],
};

/// Descriptor for `UpdateConsentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateConsentRequestDescriptor = $convert.base64Decode(
    'ChRVcGRhdGVDb25zZW50UmVxdWVzdBIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQSLQoSYW5vbn'
    'ltb3VzX3Jlc2VhcmNoGAIgASgIUhFhbm9ueW1vdXNSZXNlYXJjaBIhCgxyZWdpb25fc3RhdHMY'
    'AyABKAhSC3JlZ2lvblN0YXRzEicKD2VudGVycHJpc2VfZGF0YRgEIAEoCFIOZW50ZXJwcmlzZU'
    'RhdGE=');

@$core.Deprecated('Use triggerAggregationRequestDescriptor instead')
const TriggerAggregationRequest$json = {
  '1': 'TriggerAggregationRequest',
  '2': [
    {'1': 'admin_id', '3': 1, '4': 1, '5': 9, '10': 'adminId'},
    {'1': 'biomarker', '3': 2, '4': 1, '5': 9, '10': 'biomarker'},
    {'1': 'date', '3': 3, '4': 1, '5': 9, '10': 'date'},
  ],
};

/// Descriptor for `TriggerAggregationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List triggerAggregationRequestDescriptor =
    $convert.base64Decode(
        'ChlUcmlnZ2VyQWdncmVnYXRpb25SZXF1ZXN0EhkKCGFkbWluX2lkGAEgASgJUgdhZG1pbklkEh'
        'wKCWJpb21hcmtlchgCIAEoCVIJYmlvbWFya2VyEhIKBGRhdGUYAyABKAlSBGRhdGU=');

@$core.Deprecated('Use triggerAggregationResponseDescriptor instead')
const TriggerAggregationResponse$json = {
  '1': 'TriggerAggregationResponse',
  '2': [
    {'1': 'aggregation_id', '3': 1, '4': 1, '5': 9, '10': 'aggregationId'},
    {'1': 'status', '3': 2, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `TriggerAggregationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List triggerAggregationResponseDescriptor =
    $convert.base64Decode(
        'ChpUcmlnZ2VyQWdncmVnYXRpb25SZXNwb25zZRIlCg5hZ2dyZWdhdGlvbl9pZBgBIAEoCVINYW'
        'dncmVnYXRpb25JZBIWCgZzdGF0dXMYAiABKAlSBnN0YXR1cw==');

@$core.Deprecated('Use getInsightsRequestDescriptor instead')
const GetInsightsRequest$json = {
  '1': 'GetInsightsRequest',
  '2': [
    {'1': 'biomarker', '3': 1, '4': 1, '5': 9, '10': 'biomarker'},
    {'1': 'limit', '3': 2, '4': 1, '5': 5, '10': 'limit'},
  ],
};

/// Descriptor for `GetInsightsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getInsightsRequestDescriptor = $convert.base64Decode(
    'ChJHZXRJbnNpZ2h0c1JlcXVlc3QSHAoJYmlvbWFya2VyGAEgASgJUgliaW9tYXJrZXISFAoFbG'
    'ltaXQYAiABKAVSBWxpbWl0');

@$core.Deprecated('Use getInsightsResponseDescriptor instead')
const GetInsightsResponse$json = {
  '1': 'GetInsightsResponse',
  '2': [
    {
      '1': 'insights',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.AggregatedInsight',
      '10': 'insights'
    },
  ],
};

/// Descriptor for `GetInsightsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getInsightsResponseDescriptor = $convert.base64Decode(
    'ChNHZXRJbnNpZ2h0c1Jlc3BvbnNlEjoKCGluc2lnaHRzGAEgAygLMh4ubWFucGFzaWsudjEuQW'
    'dncmVnYXRlZEluc2lnaHRSCGluc2lnaHRz');

@$core.Deprecated('Use aggregatedInsightDescriptor instead')
const AggregatedInsight$json = {
  '1': 'AggregatedInsight',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'biomarker', '3': 2, '4': 1, '5': 9, '10': 'biomarker'},
    {'1': 'insight_type', '3': 3, '4': 1, '5': 9, '10': 'insightType'},
    {'1': 'summary', '3': 4, '4': 1, '5': 9, '10': 'summary'},
    {
      '1': 'created_at',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
  ],
};

/// Descriptor for `AggregatedInsight`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List aggregatedInsightDescriptor = $convert.base64Decode(
    'ChFBZ2dyZWdhdGVkSW5zaWdodBIOCgJpZBgBIAEoCVICaWQSHAoJYmlvbWFya2VyGAIgASgJUg'
    'liaW9tYXJrZXISIQoMaW5zaWdodF90eXBlGAMgASgJUgtpbnNpZ2h0VHlwZRIYCgdzdW1tYXJ5'
    'GAQgASgJUgdzdW1tYXJ5EjkKCmNyZWF0ZWRfYXQYBSABKAsyGi5nb29nbGUucHJvdG9idWYuVG'
    'ltZXN0YW1wUgljcmVhdGVkQXQ=');

@$core.Deprecated('Use getSalesReportRequestDescriptor instead')
const GetSalesReportRequest$json = {
  '1': 'GetSalesReportRequest',
  '2': [
    {'1': 'developer_id', '3': 1, '4': 1, '5': 9, '10': 'developerId'},
    {'1': 'period', '3': 2, '4': 1, '5': 9, '10': 'period'},
  ],
};

/// Descriptor for `GetSalesReportRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSalesReportRequestDescriptor = $convert.base64Decode(
    'ChVHZXRTYWxlc1JlcG9ydFJlcXVlc3QSIQoMZGV2ZWxvcGVyX2lkGAEgASgJUgtkZXZlbG9wZX'
    'JJZBIWCgZwZXJpb2QYAiABKAlSBnBlcmlvZA==');

@$core.Deprecated('Use salesReportDescriptor instead')
const SalesReport$json = {
  '1': 'SalesReport',
  '2': [
    {'1': 'developer_id', '3': 1, '4': 1, '5': 9, '10': 'developerId'},
    {'1': 'period', '3': 2, '4': 1, '5': 9, '10': 'period'},
    {'1': 'total_revenue_krw', '3': 3, '4': 1, '5': 3, '10': 'totalRevenueKrw'},
    {'1': 'platform_fee_krw', '3': 4, '4': 1, '5': 3, '10': 'platformFeeKrw'},
    {
      '1': 'developer_earning_krw',
      '3': 5,
      '4': 1,
      '5': 3,
      '10': 'developerEarningKrw'
    },
    {'1': 'total_sales', '3': 6, '4': 1, '5': 5, '10': 'totalSales'},
    {
      '1': 'top_items',
      '3': 7,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.SalesItem',
      '10': 'topItems'
    },
  ],
};

/// Descriptor for `SalesReport`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List salesReportDescriptor = $convert.base64Decode(
    'CgtTYWxlc1JlcG9ydBIhCgxkZXZlbG9wZXJfaWQYASABKAlSC2RldmVsb3BlcklkEhYKBnBlcm'
    'lvZBgCIAEoCVIGcGVyaW9kEioKEXRvdGFsX3JldmVudWVfa3J3GAMgASgDUg90b3RhbFJldmVu'
    'dWVLcncSKAoQcGxhdGZvcm1fZmVlX2tydxgEIAEoA1IOcGxhdGZvcm1GZWVLcncSMgoVZGV2ZW'
    'xvcGVyX2Vhcm5pbmdfa3J3GAUgASgDUhNkZXZlbG9wZXJFYXJuaW5nS3J3Eh8KC3RvdGFsX3Nh'
    'bGVzGAYgASgFUgp0b3RhbFNhbGVzEjMKCXRvcF9pdGVtcxgHIAMoCzIWLm1hbnBhc2lrLnYxLl'
    'NhbGVzSXRlbVIIdG9wSXRlbXM=');

@$core.Deprecated('Use salesItemDescriptor instead')
const SalesItem$json = {
  '1': 'SalesItem',
  '2': [
    {'1': 'item_id', '3': 1, '4': 1, '5': 9, '10': 'itemId'},
    {'1': 'cartridge_name', '3': 2, '4': 1, '5': 9, '10': 'cartridgeName'},
    {'1': 'sales_count', '3': 3, '4': 1, '5': 5, '10': 'salesCount'},
    {'1': 'revenue_krw', '3': 4, '4': 1, '5': 3, '10': 'revenueKrw'},
  ],
};

/// Descriptor for `SalesItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List salesItemDescriptor = $convert.base64Decode(
    'CglTYWxlc0l0ZW0SFwoHaXRlbV9pZBgBIAEoCVIGaXRlbUlkEiUKDmNhcnRyaWRnZV9uYW1lGA'
    'IgASgJUg1jYXJ0cmlkZ2VOYW1lEh8KC3NhbGVzX2NvdW50GAMgASgFUgpzYWxlc0NvdW50Eh8K'
    'C3JldmVudWVfa3J3GAQgASgDUgpyZXZlbnVlS3J3');

@$core.Deprecated('Use getPayoutHistoryRequestDescriptor instead')
const GetPayoutHistoryRequest$json = {
  '1': 'GetPayoutHistoryRequest',
  '2': [
    {'1': 'developer_id', '3': 1, '4': 1, '5': 9, '10': 'developerId'},
    {'1': 'limit', '3': 2, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 3, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `GetPayoutHistoryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPayoutHistoryRequestDescriptor = $convert.base64Decode(
    'ChdHZXRQYXlvdXRIaXN0b3J5UmVxdWVzdBIhCgxkZXZlbG9wZXJfaWQYASABKAlSC2RldmVsb3'
    'BlcklkEhQKBWxpbWl0GAIgASgFUgVsaW1pdBIWCgZvZmZzZXQYAyABKAVSBm9mZnNldA==');

@$core.Deprecated('Use getPayoutHistoryResponseDescriptor instead')
const GetPayoutHistoryResponse$json = {
  '1': 'GetPayoutHistoryResponse',
  '2': [
    {
      '1': 'payouts',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.PayoutRecord',
      '10': 'payouts'
    },
    {'1': 'total', '3': 2, '4': 1, '5': 5, '10': 'total'},
  ],
};

/// Descriptor for `GetPayoutHistoryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPayoutHistoryResponseDescriptor =
    $convert.base64Decode(
        'ChhHZXRQYXlvdXRIaXN0b3J5UmVzcG9uc2USMwoHcGF5b3V0cxgBIAMoCzIZLm1hbnBhc2lrLn'
        'YxLlBheW91dFJlY29yZFIHcGF5b3V0cxIUCgV0b3RhbBgCIAEoBVIFdG90YWw=');

@$core.Deprecated('Use payoutRecordDescriptor instead')
const PayoutRecord$json = {
  '1': 'PayoutRecord',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'amount_krw', '3': 2, '4': 1, '5': 3, '10': 'amountKrw'},
    {'1': 'status', '3': 3, '4': 1, '5': 9, '10': 'status'},
    {
      '1': 'requested_at',
      '3': 4,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'requestedAt'
    },
    {
      '1': 'completed_at',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'completedAt'
    },
  ],
};

/// Descriptor for `PayoutRecord`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List payoutRecordDescriptor = $convert.base64Decode(
    'CgxQYXlvdXRSZWNvcmQSDgoCaWQYASABKAlSAmlkEh0KCmFtb3VudF9rcncYAiABKANSCWFtb3'
    'VudEtydxIWCgZzdGF0dXMYAyABKAlSBnN0YXR1cxI9CgxyZXF1ZXN0ZWRfYXQYBCABKAsyGi5n'
    'b29nbGUucHJvdG9idWYuVGltZXN0YW1wUgtyZXF1ZXN0ZWRBdBI9Cgxjb21wbGV0ZWRfYXQYBS'
    'ABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wUgtjb21wbGV0ZWRBdA==');

@$core.Deprecated('Use configureRevenueSplitRequestDescriptor instead')
const ConfigureRevenueSplitRequest$json = {
  '1': 'ConfigureRevenueSplitRequest',
  '2': [
    {'1': 'developer_id', '3': 1, '4': 1, '5': 9, '10': 'developerId'},
    {'1': 'developer_pct', '3': 2, '4': 1, '5': 5, '10': 'developerPct'},
    {'1': 'platform_pct', '3': 3, '4': 1, '5': 5, '10': 'platformPct'},
  ],
};

/// Descriptor for `ConfigureRevenueSplitRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List configureRevenueSplitRequestDescriptor =
    $convert.base64Decode(
        'ChxDb25maWd1cmVSZXZlbnVlU3BsaXRSZXF1ZXN0EiEKDGRldmVsb3Blcl9pZBgBIAEoCVILZG'
        'V2ZWxvcGVySWQSIwoNZGV2ZWxvcGVyX3BjdBgCIAEoBVIMZGV2ZWxvcGVyUGN0EiEKDHBsYXRm'
        'b3JtX3BjdBgDIAEoBVILcGxhdGZvcm1QY3Q=');

@$core.Deprecated('Use revenueSplitConfigDescriptor instead')
const RevenueSplitConfig$json = {
  '1': 'RevenueSplitConfig',
  '2': [
    {'1': 'developer_id', '3': 1, '4': 1, '5': 9, '10': 'developerId'},
    {'1': 'developer_pct', '3': 2, '4': 1, '5': 5, '10': 'developerPct'},
    {'1': 'platform_pct', '3': 3, '4': 1, '5': 5, '10': 'platformPct'},
  ],
};

/// Descriptor for `RevenueSplitConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List revenueSplitConfigDescriptor = $convert.base64Decode(
    'ChJSZXZlbnVlU3BsaXRDb25maWcSIQoMZGV2ZWxvcGVyX2lkGAEgASgJUgtkZXZlbG9wZXJJZB'
    'IjCg1kZXZlbG9wZXJfcGN0GAIgASgFUgxkZXZlbG9wZXJQY3QSIQoMcGxhdGZvcm1fcGN0GAMg'
    'ASgFUgtwbGF0Zm9ybVBjdA==');

@$core.Deprecated('Use requestPayoutRequestDescriptor instead')
const RequestPayoutRequest$json = {
  '1': 'RequestPayoutRequest',
  '2': [
    {'1': 'developer_id', '3': 1, '4': 1, '5': 9, '10': 'developerId'},
    {'1': 'amount_krw', '3': 2, '4': 1, '5': 3, '10': 'amountKrw'},
  ],
};

/// Descriptor for `RequestPayoutRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List requestPayoutRequestDescriptor = $convert.base64Decode(
    'ChRSZXF1ZXN0UGF5b3V0UmVxdWVzdBIhCgxkZXZlbG9wZXJfaWQYASABKAlSC2RldmVsb3Blck'
    'lkEh0KCmFtb3VudF9rcncYAiABKANSCWFtb3VudEtydw==');

@$core.Deprecated('Use payoutResponseDescriptor instead')
const PayoutResponse$json = {
  '1': 'PayoutResponse',
  '2': [
    {'1': 'payout_id', '3': 1, '4': 1, '5': 9, '10': 'payoutId'},
    {'1': 'status', '3': 2, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `PayoutResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List payoutResponseDescriptor = $convert.base64Decode(
    'Cg5QYXlvdXRSZXNwb25zZRIbCglwYXlvdXRfaWQYASABKAlSCHBheW91dElkEhYKBnN0YXR1cx'
    'gCIAEoCVIGc3RhdHVz');

@$core.Deprecated('Use getCartridgeUsageStatsRequestDescriptor instead')
const GetCartridgeUsageStatsRequest$json = {
  '1': 'GetCartridgeUsageStatsRequest',
  '2': [
    {'1': 'listing_id', '3': 1, '4': 1, '5': 9, '10': 'listingId'},
    {'1': 'period', '3': 2, '4': 1, '5': 9, '10': 'period'},
  ],
};

/// Descriptor for `GetCartridgeUsageStatsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCartridgeUsageStatsRequestDescriptor =
    $convert.base64Decode(
        'Ch1HZXRDYXJ0cmlkZ2VVc2FnZVN0YXRzUmVxdWVzdBIdCgpsaXN0aW5nX2lkGAEgASgJUglsaX'
        'N0aW5nSWQSFgoGcGVyaW9kGAIgASgJUgZwZXJpb2Q=');

@$core.Deprecated('Use cartridgeUsageStatsResponseDescriptor instead')
const CartridgeUsageStatsResponse$json = {
  '1': 'CartridgeUsageStatsResponse',
  '2': [
    {'1': 'listing_id', '3': 1, '4': 1, '5': 9, '10': 'listingId'},
    {'1': 'total_uses', '3': 2, '4': 1, '5': 5, '10': 'totalUses'},
    {'1': 'unique_users', '3': 3, '4': 1, '5': 5, '10': 'uniqueUsers'},
    {'1': 'active_users_30d', '3': 4, '4': 1, '5': 5, '10': 'activeUsers30d'},
    {'1': 'avg_uses_per_user', '3': 5, '4': 1, '5': 1, '10': 'avgUsesPerUser'},
  ],
};

/// Descriptor for `CartridgeUsageStatsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cartridgeUsageStatsResponseDescriptor = $convert.base64Decode(
    'ChtDYXJ0cmlkZ2VVc2FnZVN0YXRzUmVzcG9uc2USHQoKbGlzdGluZ19pZBgBIAEoCVIJbGlzdG'
    'luZ0lkEh0KCnRvdGFsX3VzZXMYAiABKAVSCXRvdGFsVXNlcxIhCgx1bmlxdWVfdXNlcnMYAyAB'
    'KAVSC3VuaXF1ZVVzZXJzEigKEGFjdGl2ZV91c2Vyc18zMGQYBCABKAVSDmFjdGl2ZVVzZXJzMz'
    'BkEikKEWF2Z191c2VzX3Blcl91c2VyGAUgASgBUg5hdmdVc2VzUGVyVXNlcg==');

@$core.Deprecated('Use getCartridgeRatingsRequestDescriptor instead')
const GetCartridgeRatingsRequest$json = {
  '1': 'GetCartridgeRatingsRequest',
  '2': [
    {'1': 'listing_id', '3': 1, '4': 1, '5': 9, '10': 'listingId'},
    {'1': 'limit', '3': 2, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'offset', '3': 3, '4': 1, '5': 5, '10': 'offset'},
  ],
};

/// Descriptor for `GetCartridgeRatingsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCartridgeRatingsRequestDescriptor =
    $convert.base64Decode(
        'ChpHZXRDYXJ0cmlkZ2VSYXRpbmdzUmVxdWVzdBIdCgpsaXN0aW5nX2lkGAEgASgJUglsaXN0aW'
        '5nSWQSFAoFbGltaXQYAiABKAVSBWxpbWl0EhYKBm9mZnNldBgDIAEoBVIGb2Zmc2V0');

@$core.Deprecated('Use cartridgeRatingsResponseDescriptor instead')
const CartridgeRatingsResponse$json = {
  '1': 'CartridgeRatingsResponse',
  '2': [
    {'1': 'avg_rating', '3': 1, '4': 1, '5': 1, '10': 'avgRating'},
    {'1': 'total_reviews', '3': 2, '4': 1, '5': 5, '10': 'totalReviews'},
    {
      '1': 'reviews',
      '3': 3,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.UserReview',
      '10': 'reviews'
    },
  ],
};

/// Descriptor for `CartridgeRatingsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cartridgeRatingsResponseDescriptor = $convert.base64Decode(
    'ChhDYXJ0cmlkZ2VSYXRpbmdzUmVzcG9uc2USHQoKYXZnX3JhdGluZxgBIAEoAVIJYXZnUmF0aW'
    '5nEiMKDXRvdGFsX3Jldmlld3MYAiABKAVSDHRvdGFsUmV2aWV3cxIxCgdyZXZpZXdzGAMgAygL'
    'MhcubWFucGFzaWsudjEuVXNlclJldmlld1IHcmV2aWV3cw==');

@$core.Deprecated('Use userReviewDescriptor instead')
const UserReview$json = {
  '1': 'UserReview',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'display_name', '3': 3, '4': 1, '5': 9, '10': 'displayName'},
    {'1': 'rating', '3': 4, '4': 1, '5': 5, '10': 'rating'},
    {'1': 'comment', '3': 5, '4': 1, '5': 9, '10': 'comment'},
    {
      '1': 'created_at',
      '3': 6,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
  ],
};

/// Descriptor for `UserReview`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List userReviewDescriptor = $convert.base64Decode(
    'CgpVc2VyUmV2aWV3Eg4KAmlkGAEgASgJUgJpZBIXCgd1c2VyX2lkGAIgASgJUgZ1c2VySWQSIQ'
    'oMZGlzcGxheV9uYW1lGAMgASgJUgtkaXNwbGF5TmFtZRIWCgZyYXRpbmcYBCABKAVSBnJhdGlu'
    'ZxIYCgdjb21tZW50GAUgASgJUgdjb21tZW50EjkKCmNyZWF0ZWRfYXQYBiABKAsyGi5nb29nbG'
    'UucHJvdG9idWYuVGltZXN0YW1wUgljcmVhdGVkQXQ=');

@$core.Deprecated('Use submitUserReviewRequestDescriptor instead')
const SubmitUserReviewRequest$json = {
  '1': 'SubmitUserReviewRequest',
  '2': [
    {'1': 'listing_id', '3': 1, '4': 1, '5': 9, '10': 'listingId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'rating', '3': 3, '4': 1, '5': 5, '10': 'rating'},
    {'1': 'comment', '3': 4, '4': 1, '5': 9, '10': 'comment'},
  ],
};

/// Descriptor for `SubmitUserReviewRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List submitUserReviewRequestDescriptor = $convert.base64Decode(
    'ChdTdWJtaXRVc2VyUmV2aWV3UmVxdWVzdBIdCgpsaXN0aW5nX2lkGAEgASgJUglsaXN0aW5nSW'
    'QSFwoHdXNlcl9pZBgCIAEoCVIGdXNlcklkEhYKBnJhdGluZxgDIAEoBVIGcmF0aW5nEhgKB2Nv'
    'bW1lbnQYBCABKAlSB2NvbW1lbnQ=');

@$core.Deprecated('Use submitUserReviewResponseDescriptor instead')
const SubmitUserReviewResponse$json = {
  '1': 'SubmitUserReviewResponse',
  '2': [
    {'1': 'review_id', '3': 1, '4': 1, '5': 9, '10': 'reviewId'},
    {'1': 'success', '3': 2, '4': 1, '5': 8, '10': 'success'},
  ],
};

/// Descriptor for `SubmitUserReviewResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List submitUserReviewResponseDescriptor =
    $convert.base64Decode(
        'ChhTdWJtaXRVc2VyUmV2aWV3UmVzcG9uc2USGwoJcmV2aWV3X2lkGAEgASgJUghyZXZpZXdJZB'
        'IYCgdzdWNjZXNzGAIgASgIUgdzdWNjZXNz');

@$core.Deprecated('Use getDeveloperAnalyticsRequestDescriptor instead')
const GetDeveloperAnalyticsRequest$json = {
  '1': 'GetDeveloperAnalyticsRequest',
  '2': [
    {'1': 'developer_id', '3': 1, '4': 1, '5': 9, '10': 'developerId'},
    {'1': 'period', '3': 2, '4': 1, '5': 9, '10': 'period'},
  ],
};

/// Descriptor for `GetDeveloperAnalyticsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDeveloperAnalyticsRequestDescriptor =
    $convert.base64Decode(
        'ChxHZXREZXZlbG9wZXJBbmFseXRpY3NSZXF1ZXN0EiEKDGRldmVsb3Blcl9pZBgBIAEoCVILZG'
        'V2ZWxvcGVySWQSFgoGcGVyaW9kGAIgASgJUgZwZXJpb2Q=');

@$core.Deprecated('Use developerAnalyticsDescriptor instead')
const DeveloperAnalytics$json = {
  '1': 'DeveloperAnalytics',
  '2': [
    {'1': 'developer_id', '3': 1, '4': 1, '5': 9, '10': 'developerId'},
    {'1': 'total_cartridges', '3': 2, '4': 1, '5': 5, '10': 'totalCartridges'},
    {'1': 'total_downloads', '3': 3, '4': 1, '5': 5, '10': 'totalDownloads'},
    {'1': 'avg_rating', '3': 4, '4': 1, '5': 1, '10': 'avgRating'},
    {'1': 'total_revenue_krw', '3': 5, '4': 1, '5': 3, '10': 'totalRevenueKrw'},
    {
      '1': 'top_cartridges',
      '3': 6,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.CartridgePerformance',
      '10': 'topCartridges'
    },
  ],
};

/// Descriptor for `DeveloperAnalytics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List developerAnalyticsDescriptor = $convert.base64Decode(
    'ChJEZXZlbG9wZXJBbmFseXRpY3MSIQoMZGV2ZWxvcGVyX2lkGAEgASgJUgtkZXZlbG9wZXJJZB'
    'IpChB0b3RhbF9jYXJ0cmlkZ2VzGAIgASgFUg90b3RhbENhcnRyaWRnZXMSJwoPdG90YWxfZG93'
    'bmxvYWRzGAMgASgFUg50b3RhbERvd25sb2FkcxIdCgphdmdfcmF0aW5nGAQgASgBUglhdmdSYX'
    'RpbmcSKgoRdG90YWxfcmV2ZW51ZV9rcncYBSABKANSD3RvdGFsUmV2ZW51ZUtydxJICg50b3Bf'
    'Y2FydHJpZGdlcxgGIAMoCzIhLm1hbnBhc2lrLnYxLkNhcnRyaWRnZVBlcmZvcm1hbmNlUg10b3'
    'BDYXJ0cmlkZ2Vz');

@$core.Deprecated('Use cartridgePerformanceDescriptor instead')
const CartridgePerformance$json = {
  '1': 'CartridgePerformance',
  '2': [
    {'1': 'listing_id', '3': 1, '4': 1, '5': 9, '10': 'listingId'},
    {'1': 'name', '3': 2, '4': 1, '5': 9, '10': 'name'},
    {'1': 'downloads', '3': 3, '4': 1, '5': 5, '10': 'downloads'},
    {'1': 'rating', '3': 4, '4': 1, '5': 1, '10': 'rating'},
    {'1': 'revenue_krw', '3': 5, '4': 1, '5': 3, '10': 'revenueKrw'},
  ],
};

/// Descriptor for `CartridgePerformance`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cartridgePerformanceDescriptor = $convert.base64Decode(
    'ChRDYXJ0cmlkZ2VQZXJmb3JtYW5jZRIdCgpsaXN0aW5nX2lkGAEgASgJUglsaXN0aW5nSWQSEg'
    'oEbmFtZRgCIAEoCVIEbmFtZRIcCglkb3dubG9hZHMYAyABKAVSCWRvd25sb2FkcxIWCgZyYXRp'
    'bmcYBCABKAFSBnJhdGluZxIfCgtyZXZlbnVlX2tydxgFIAEoA1IKcmV2ZW51ZUtydw==');

@$core.Deprecated('Use createVoiceProfileRequestDescriptor instead')
const CreateVoiceProfileRequest$json = {
  '1': 'CreateVoiceProfileRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'profile_name', '3': 2, '4': 1, '5': 9, '10': 'profileName'},
    {'1': 'voice_sample', '3': 3, '4': 1, '5': 12, '10': 'voiceSample'},
    {'1': 'language', '3': 4, '4': 1, '5': 9, '10': 'language'},
  ],
};

/// Descriptor for `CreateVoiceProfileRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createVoiceProfileRequestDescriptor = $convert.base64Decode(
    'ChlDcmVhdGVWb2ljZVByb2ZpbGVSZXF1ZXN0EhcKB3VzZXJfaWQYASABKAlSBnVzZXJJZBIhCg'
    'xwcm9maWxlX25hbWUYAiABKAlSC3Byb2ZpbGVOYW1lEiEKDHZvaWNlX3NhbXBsZRgDIAEoDFIL'
    'dm9pY2VTYW1wbGUSGgoIbGFuZ3VhZ2UYBCABKAlSCGxhbmd1YWdl');

@$core.Deprecated('Use voiceProfileDescriptor instead')
const VoiceProfile$json = {
  '1': 'VoiceProfile',
  '2': [
    {'1': 'profile_id', '3': 1, '4': 1, '5': 9, '10': 'profileId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'profile_name', '3': 3, '4': 1, '5': 9, '10': 'profileName'},
    {'1': 'language', '3': 4, '4': 1, '5': 9, '10': 'language'},
    {'1': 'status', '3': 5, '4': 1, '5': 9, '10': 'status'},
    {'1': 'model_url', '3': 6, '4': 1, '5': 9, '10': 'modelUrl'},
    {
      '1': 'created_at',
      '3': 7,
      '4': 1,
      '5': 11,
      '6': '.google.protobuf.Timestamp',
      '10': 'createdAt'
    },
  ],
};

/// Descriptor for `VoiceProfile`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List voiceProfileDescriptor = $convert.base64Decode(
    'CgxWb2ljZVByb2ZpbGUSHQoKcHJvZmlsZV9pZBgBIAEoCVIJcHJvZmlsZUlkEhcKB3VzZXJfaW'
    'QYAiABKAlSBnVzZXJJZBIhCgxwcm9maWxlX25hbWUYAyABKAlSC3Byb2ZpbGVOYW1lEhoKCGxh'
    'bmd1YWdlGAQgASgJUghsYW5ndWFnZRIWCgZzdGF0dXMYBSABKAlSBnN0YXR1cxIbCgltb2RlbF'
    '91cmwYBiABKAlSCG1vZGVsVXJsEjkKCmNyZWF0ZWRfYXQYByABKAsyGi5nb29nbGUucHJvdG9i'
    'dWYuVGltZXN0YW1wUgljcmVhdGVkQXQ=');

@$core.Deprecated('Use getVoiceProfileRequestDescriptor instead')
const GetVoiceProfileRequest$json = {
  '1': 'GetVoiceProfileRequest',
  '2': [
    {'1': 'profile_id', '3': 1, '4': 1, '5': 9, '10': 'profileId'},
  ],
};

/// Descriptor for `GetVoiceProfileRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getVoiceProfileRequestDescriptor =
    $convert.base64Decode(
        'ChZHZXRWb2ljZVByb2ZpbGVSZXF1ZXN0Eh0KCnByb2ZpbGVfaWQYASABKAlSCXByb2ZpbGVJZA'
        '==');

@$core.Deprecated('Use listVoiceProfilesRequestDescriptor instead')
const ListVoiceProfilesRequest$json = {
  '1': 'ListVoiceProfilesRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
  ],
};

/// Descriptor for `ListVoiceProfilesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listVoiceProfilesRequestDescriptor =
    $convert.base64Decode(
        'ChhMaXN0Vm9pY2VQcm9maWxlc1JlcXVlc3QSFwoHdXNlcl9pZBgBIAEoCVIGdXNlcklk');

@$core.Deprecated('Use listVoiceProfilesResponseDescriptor instead')
const ListVoiceProfilesResponse$json = {
  '1': 'ListVoiceProfilesResponse',
  '2': [
    {
      '1': 'profiles',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.v1.VoiceProfile',
      '10': 'profiles'
    },
  ],
};

/// Descriptor for `ListVoiceProfilesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listVoiceProfilesResponseDescriptor =
    $convert.base64Decode(
        'ChlMaXN0Vm9pY2VQcm9maWxlc1Jlc3BvbnNlEjUKCHByb2ZpbGVzGAEgAygLMhkubWFucGFzaW'
        'sudjEuVm9pY2VQcm9maWxlUghwcm9maWxlcw==');

@$core.Deprecated('Use synthesizeTranslationRequestDescriptor instead')
const SynthesizeTranslationRequest$json = {
  '1': 'SynthesizeTranslationRequest',
  '2': [
    {'1': 'profile_id', '3': 1, '4': 1, '5': 9, '10': 'profileId'},
    {'1': 'text', '3': 2, '4': 1, '5': 9, '10': 'text'},
    {'1': 'target_language', '3': 3, '4': 1, '5': 9, '10': 'targetLanguage'},
  ],
};

/// Descriptor for `SynthesizeTranslationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List synthesizeTranslationRequestDescriptor =
    $convert.base64Decode(
        'ChxTeW50aGVzaXplVHJhbnNsYXRpb25SZXF1ZXN0Eh0KCnByb2ZpbGVfaWQYASABKAlSCXByb2'
        'ZpbGVJZBISCgR0ZXh0GAIgASgJUgR0ZXh0EicKD3RhcmdldF9sYW5ndWFnZRgDIAEoCVIOdGFy'
        'Z2V0TGFuZ3VhZ2U=');

@$core.Deprecated('Use synthesizeTranslationResponseDescriptor instead')
const SynthesizeTranslationResponse$json = {
  '1': 'SynthesizeTranslationResponse',
  '2': [
    {'1': 'audio_data', '3': 1, '4': 1, '5': 12, '10': 'audioData'},
    {'1': 'translated_text', '3': 2, '4': 1, '5': 9, '10': 'translatedText'},
    {'1': 'language', '3': 3, '4': 1, '5': 9, '10': 'language'},
    {'1': 'similarity_score', '3': 4, '4': 1, '5': 1, '10': 'similarityScore'},
  ],
};

/// Descriptor for `SynthesizeTranslationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List synthesizeTranslationResponseDescriptor = $convert.base64Decode(
    'Ch1TeW50aGVzaXplVHJhbnNsYXRpb25SZXNwb25zZRIdCgphdWRpb19kYXRhGAEgASgMUglhdW'
    'Rpb0RhdGESJwoPdHJhbnNsYXRlZF90ZXh0GAIgASgJUg50cmFuc2xhdGVkVGV4dBIaCghsYW5n'
    'dWFnZRgDIAEoCVIIbGFuZ3VhZ2USKQoQc2ltaWxhcml0eV9zY29yZRgEIAEoAVIPc2ltaWxhcm'
    'l0eVNjb3Jl');

@$core.Deprecated('Use deleteVoiceProfileRequestDescriptor instead')
const DeleteVoiceProfileRequest$json = {
  '1': 'DeleteVoiceProfileRequest',
  '2': [
    {'1': 'profile_id', '3': 1, '4': 1, '5': 9, '10': 'profileId'},
  ],
};

/// Descriptor for `DeleteVoiceProfileRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteVoiceProfileRequestDescriptor =
    $convert.base64Decode(
        'ChlEZWxldGVWb2ljZVByb2ZpbGVSZXF1ZXN0Eh0KCnByb2ZpbGVfaWQYASABKAlSCXByb2ZpbG'
        'VJZA==');

@$core.Deprecated('Use deleteVoiceProfileResponseDescriptor instead')
const DeleteVoiceProfileResponse$json = {
  '1': 'DeleteVoiceProfileResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
  ],
};

/// Descriptor for `DeleteVoiceProfileResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteVoiceProfileResponseDescriptor =
    $convert.base64Decode(
        'ChpEZWxldGVWb2ljZVByb2ZpbGVSZXNwb25zZRIYCgdzdWNjZXNzGAEgASgIUgdzdWNjZXNz');
