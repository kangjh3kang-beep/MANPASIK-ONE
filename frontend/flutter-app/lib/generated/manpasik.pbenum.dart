// This is a generated file - do not edit.
//
// Generated from manpasik.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class SocialProvider extends $pb.ProtobufEnum {
  static const SocialProvider SOCIAL_PROVIDER_UNSPECIFIED =
      SocialProvider._(0, _omitEnumNames ? '' : 'SOCIAL_PROVIDER_UNSPECIFIED');
  static const SocialProvider SOCIAL_PROVIDER_GOOGLE =
      SocialProvider._(1, _omitEnumNames ? '' : 'SOCIAL_PROVIDER_GOOGLE');
  static const SocialProvider SOCIAL_PROVIDER_APPLE =
      SocialProvider._(2, _omitEnumNames ? '' : 'SOCIAL_PROVIDER_APPLE');
  static const SocialProvider SOCIAL_PROVIDER_KAKAO =
      SocialProvider._(3, _omitEnumNames ? '' : 'SOCIAL_PROVIDER_KAKAO');
  static const SocialProvider SOCIAL_PROVIDER_NAVER =
      SocialProvider._(4, _omitEnumNames ? '' : 'SOCIAL_PROVIDER_NAVER');

  static const $core.List<SocialProvider> values = <SocialProvider>[
    SOCIAL_PROVIDER_UNSPECIFIED,
    SOCIAL_PROVIDER_GOOGLE,
    SOCIAL_PROVIDER_APPLE,
    SOCIAL_PROVIDER_KAKAO,
    SOCIAL_PROVIDER_NAVER,
  ];

  static final $core.List<SocialProvider?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static SocialProvider? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SocialProvider._(super.value, super.name);
}

class Gender extends $pb.ProtobufEnum {
  static const Gender GENDER_UNSPECIFIED =
      Gender._(0, _omitEnumNames ? '' : 'GENDER_UNSPECIFIED');
  static const Gender GENDER_MALE =
      Gender._(1, _omitEnumNames ? '' : 'GENDER_MALE');
  static const Gender GENDER_FEMALE =
      Gender._(2, _omitEnumNames ? '' : 'GENDER_FEMALE');
  static const Gender GENDER_OTHER =
      Gender._(3, _omitEnumNames ? '' : 'GENDER_OTHER');
  static const Gender GENDER_PREFER_NOT_TO_SAY =
      Gender._(4, _omitEnumNames ? '' : 'GENDER_PREFER_NOT_TO_SAY');

  static const $core.List<Gender> values = <Gender>[
    GENDER_UNSPECIFIED,
    GENDER_MALE,
    GENDER_FEMALE,
    GENDER_OTHER,
    GENDER_PREFER_NOT_TO_SAY,
  ];

  static final $core.List<Gender?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static Gender? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Gender._(super.value, super.name);
}

class DeviceStatus extends $pb.ProtobufEnum {
  static const DeviceStatus DEVICE_STATUS_UNKNOWN =
      DeviceStatus._(0, _omitEnumNames ? '' : 'DEVICE_STATUS_UNKNOWN');
  static const DeviceStatus DEVICE_STATUS_ONLINE =
      DeviceStatus._(1, _omitEnumNames ? '' : 'DEVICE_STATUS_ONLINE');
  static const DeviceStatus DEVICE_STATUS_OFFLINE =
      DeviceStatus._(2, _omitEnumNames ? '' : 'DEVICE_STATUS_OFFLINE');
  static const DeviceStatus DEVICE_STATUS_MEASURING =
      DeviceStatus._(3, _omitEnumNames ? '' : 'DEVICE_STATUS_MEASURING');
  static const DeviceStatus DEVICE_STATUS_UPDATING =
      DeviceStatus._(4, _omitEnumNames ? '' : 'DEVICE_STATUS_UPDATING');
  static const DeviceStatus DEVICE_STATUS_ERROR =
      DeviceStatus._(5, _omitEnumNames ? '' : 'DEVICE_STATUS_ERROR');

  static const $core.List<DeviceStatus> values = <DeviceStatus>[
    DEVICE_STATUS_UNKNOWN,
    DEVICE_STATUS_ONLINE,
    DEVICE_STATUS_OFFLINE,
    DEVICE_STATUS_MEASURING,
    DEVICE_STATUS_UPDATING,
    DEVICE_STATUS_ERROR,
  ];

  static final $core.List<DeviceStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static DeviceStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DeviceStatus._(super.value, super.name);
}

class CommandType extends $pb.ProtobufEnum {
  static const CommandType COMMAND_TYPE_UNKNOWN =
      CommandType._(0, _omitEnumNames ? '' : 'COMMAND_TYPE_UNKNOWN');
  static const CommandType COMMAND_TYPE_START_MEASUREMENT =
      CommandType._(1, _omitEnumNames ? '' : 'COMMAND_TYPE_START_MEASUREMENT');
  static const CommandType COMMAND_TYPE_STOP_MEASUREMENT =
      CommandType._(2, _omitEnumNames ? '' : 'COMMAND_TYPE_STOP_MEASUREMENT');
  static const CommandType COMMAND_TYPE_CALIBRATE =
      CommandType._(3, _omitEnumNames ? '' : 'COMMAND_TYPE_CALIBRATE');
  static const CommandType COMMAND_TYPE_REBOOT =
      CommandType._(4, _omitEnumNames ? '' : 'COMMAND_TYPE_REBOOT');
  static const CommandType COMMAND_TYPE_OTA_UPDATE =
      CommandType._(5, _omitEnumNames ? '' : 'COMMAND_TYPE_OTA_UPDATE');

  static const $core.List<CommandType> values = <CommandType>[
    COMMAND_TYPE_UNKNOWN,
    COMMAND_TYPE_START_MEASUREMENT,
    COMMAND_TYPE_STOP_MEASUREMENT,
    COMMAND_TYPE_CALIBRATE,
    COMMAND_TYPE_REBOOT,
    COMMAND_TYPE_OTA_UPDATE,
  ];

  static final $core.List<CommandType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static CommandType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CommandType._(super.value, super.name);
}

class SubscriptionTier extends $pb.ProtobufEnum {
  static const SubscriptionTier SUBSCRIPTION_TIER_FREE =
      SubscriptionTier._(0, _omitEnumNames ? '' : 'SUBSCRIPTION_TIER_FREE');
  static const SubscriptionTier SUBSCRIPTION_TIER_BASIC =
      SubscriptionTier._(1, _omitEnumNames ? '' : 'SUBSCRIPTION_TIER_BASIC');
  static const SubscriptionTier SUBSCRIPTION_TIER_PRO =
      SubscriptionTier._(2, _omitEnumNames ? '' : 'SUBSCRIPTION_TIER_PRO');
  static const SubscriptionTier SUBSCRIPTION_TIER_CLINICAL =
      SubscriptionTier._(3, _omitEnumNames ? '' : 'SUBSCRIPTION_TIER_CLINICAL');

  static const $core.List<SubscriptionTier> values = <SubscriptionTier>[
    SUBSCRIPTION_TIER_FREE,
    SUBSCRIPTION_TIER_BASIC,
    SUBSCRIPTION_TIER_PRO,
    SUBSCRIPTION_TIER_CLINICAL,
  ];

  static final $core.List<SubscriptionTier?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static SubscriptionTier? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SubscriptionTier._(super.value, super.name);
}

class SubscriptionStatus extends $pb.ProtobufEnum {
  static const SubscriptionStatus SUBSCRIPTION_STATUS_UNKNOWN =
      SubscriptionStatus._(
          0, _omitEnumNames ? '' : 'SUBSCRIPTION_STATUS_UNKNOWN');
  static const SubscriptionStatus SUBSCRIPTION_STATUS_ACTIVE =
      SubscriptionStatus._(
          1, _omitEnumNames ? '' : 'SUBSCRIPTION_STATUS_ACTIVE');
  static const SubscriptionStatus SUBSCRIPTION_STATUS_CANCELLED =
      SubscriptionStatus._(
          2, _omitEnumNames ? '' : 'SUBSCRIPTION_STATUS_CANCELLED');
  static const SubscriptionStatus SUBSCRIPTION_STATUS_EXPIRED =
      SubscriptionStatus._(
          3, _omitEnumNames ? '' : 'SUBSCRIPTION_STATUS_EXPIRED');
  static const SubscriptionStatus SUBSCRIPTION_STATUS_SUSPENDED =
      SubscriptionStatus._(
          4, _omitEnumNames ? '' : 'SUBSCRIPTION_STATUS_SUSPENDED');
  static const SubscriptionStatus SUBSCRIPTION_STATUS_TRIAL =
      SubscriptionStatus._(
          5, _omitEnumNames ? '' : 'SUBSCRIPTION_STATUS_TRIAL');

  static const $core.List<SubscriptionStatus> values = <SubscriptionStatus>[
    SUBSCRIPTION_STATUS_UNKNOWN,
    SUBSCRIPTION_STATUS_ACTIVE,
    SUBSCRIPTION_STATUS_CANCELLED,
    SUBSCRIPTION_STATUS_EXPIRED,
    SUBSCRIPTION_STATUS_SUSPENDED,
    SUBSCRIPTION_STATUS_TRIAL,
  ];

  static final $core.List<SubscriptionStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static SubscriptionStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SubscriptionStatus._(super.value, super.name);
}

class ProductCategory extends $pb.ProtobufEnum {
  static const ProductCategory PRODUCT_CATEGORY_UNKNOWN =
      ProductCategory._(0, _omitEnumNames ? '' : 'PRODUCT_CATEGORY_UNKNOWN');
  static const ProductCategory PRODUCT_CATEGORY_CARTRIDGE =
      ProductCategory._(1, _omitEnumNames ? '' : 'PRODUCT_CATEGORY_CARTRIDGE');
  static const ProductCategory PRODUCT_CATEGORY_READER =
      ProductCategory._(2, _omitEnumNames ? '' : 'PRODUCT_CATEGORY_READER');
  static const ProductCategory PRODUCT_CATEGORY_ACCESSORY =
      ProductCategory._(3, _omitEnumNames ? '' : 'PRODUCT_CATEGORY_ACCESSORY');
  static const ProductCategory PRODUCT_CATEGORY_BUNDLE =
      ProductCategory._(4, _omitEnumNames ? '' : 'PRODUCT_CATEGORY_BUNDLE');

  static const $core.List<ProductCategory> values = <ProductCategory>[
    PRODUCT_CATEGORY_UNKNOWN,
    PRODUCT_CATEGORY_CARTRIDGE,
    PRODUCT_CATEGORY_READER,
    PRODUCT_CATEGORY_ACCESSORY,
    PRODUCT_CATEGORY_BUNDLE,
  ];

  static final $core.List<ProductCategory?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static ProductCategory? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ProductCategory._(super.value, super.name);
}

class OrderStatus extends $pb.ProtobufEnum {
  static const OrderStatus ORDER_STATUS_UNKNOWN =
      OrderStatus._(0, _omitEnumNames ? '' : 'ORDER_STATUS_UNKNOWN');
  static const OrderStatus ORDER_STATUS_PENDING =
      OrderStatus._(1, _omitEnumNames ? '' : 'ORDER_STATUS_PENDING');
  static const OrderStatus ORDER_STATUS_PAID =
      OrderStatus._(2, _omitEnumNames ? '' : 'ORDER_STATUS_PAID');
  static const OrderStatus ORDER_STATUS_SHIPPED =
      OrderStatus._(3, _omitEnumNames ? '' : 'ORDER_STATUS_SHIPPED');
  static const OrderStatus ORDER_STATUS_DELIVERED =
      OrderStatus._(4, _omitEnumNames ? '' : 'ORDER_STATUS_DELIVERED');
  static const OrderStatus ORDER_STATUS_CANCELLED =
      OrderStatus._(5, _omitEnumNames ? '' : 'ORDER_STATUS_CANCELLED');
  static const OrderStatus ORDER_STATUS_REFUNDED =
      OrderStatus._(6, _omitEnumNames ? '' : 'ORDER_STATUS_REFUNDED');

  static const $core.List<OrderStatus> values = <OrderStatus>[
    ORDER_STATUS_UNKNOWN,
    ORDER_STATUS_PENDING,
    ORDER_STATUS_PAID,
    ORDER_STATUS_SHIPPED,
    ORDER_STATUS_DELIVERED,
    ORDER_STATUS_CANCELLED,
    ORDER_STATUS_REFUNDED,
  ];

  static final $core.List<OrderStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 6);
  static OrderStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OrderStatus._(super.value, super.name);
}

class PaymentType extends $pb.ProtobufEnum {
  static const PaymentType PAYMENT_TYPE_UNKNOWN =
      PaymentType._(0, _omitEnumNames ? '' : 'PAYMENT_TYPE_UNKNOWN');
  static const PaymentType PAYMENT_TYPE_ONE_TIME =
      PaymentType._(1, _omitEnumNames ? '' : 'PAYMENT_TYPE_ONE_TIME');
  static const PaymentType PAYMENT_TYPE_SUBSCRIPTION =
      PaymentType._(2, _omitEnumNames ? '' : 'PAYMENT_TYPE_SUBSCRIPTION');

  static const $core.List<PaymentType> values = <PaymentType>[
    PAYMENT_TYPE_UNKNOWN,
    PAYMENT_TYPE_ONE_TIME,
    PAYMENT_TYPE_SUBSCRIPTION,
  ];

  static final $core.List<PaymentType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static PaymentType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PaymentType._(super.value, super.name);
}

class PaymentStatus extends $pb.ProtobufEnum {
  static const PaymentStatus PAYMENT_STATUS_UNKNOWN =
      PaymentStatus._(0, _omitEnumNames ? '' : 'PAYMENT_STATUS_UNKNOWN');
  static const PaymentStatus PAYMENT_STATUS_PENDING =
      PaymentStatus._(1, _omitEnumNames ? '' : 'PAYMENT_STATUS_PENDING');
  static const PaymentStatus PAYMENT_STATUS_COMPLETED =
      PaymentStatus._(2, _omitEnumNames ? '' : 'PAYMENT_STATUS_COMPLETED');
  static const PaymentStatus PAYMENT_STATUS_FAILED =
      PaymentStatus._(3, _omitEnumNames ? '' : 'PAYMENT_STATUS_FAILED');
  static const PaymentStatus PAYMENT_STATUS_CANCELLED =
      PaymentStatus._(4, _omitEnumNames ? '' : 'PAYMENT_STATUS_CANCELLED');
  static const PaymentStatus PAYMENT_STATUS_REFUNDED =
      PaymentStatus._(5, _omitEnumNames ? '' : 'PAYMENT_STATUS_REFUNDED');
  static const PaymentStatus PAYMENT_STATUS_PARTIAL_REFUND =
      PaymentStatus._(6, _omitEnumNames ? '' : 'PAYMENT_STATUS_PARTIAL_REFUND');

  static const $core.List<PaymentStatus> values = <PaymentStatus>[
    PAYMENT_STATUS_UNKNOWN,
    PAYMENT_STATUS_PENDING,
    PAYMENT_STATUS_COMPLETED,
    PAYMENT_STATUS_FAILED,
    PAYMENT_STATUS_CANCELLED,
    PAYMENT_STATUS_REFUNDED,
    PAYMENT_STATUS_PARTIAL_REFUND,
  ];

  static final $core.List<PaymentStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 6);
  static PaymentStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PaymentStatus._(super.value, super.name);
}

class AiModelType extends $pb.ProtobufEnum {
  static const AiModelType AI_MODEL_TYPE_UNSPECIFIED =
      AiModelType._(0, _omitEnumNames ? '' : 'AI_MODEL_TYPE_UNSPECIFIED');
  static const AiModelType AI_MODEL_TYPE_BIOMARKER_CLASSIFIER = AiModelType._(
      1, _omitEnumNames ? '' : 'AI_MODEL_TYPE_BIOMARKER_CLASSIFIER');
  static const AiModelType AI_MODEL_TYPE_ANOMALY_DETECTOR =
      AiModelType._(2, _omitEnumNames ? '' : 'AI_MODEL_TYPE_ANOMALY_DETECTOR');
  static const AiModelType AI_MODEL_TYPE_TREND_PREDICTOR =
      AiModelType._(3, _omitEnumNames ? '' : 'AI_MODEL_TYPE_TREND_PREDICTOR');
  static const AiModelType AI_MODEL_TYPE_HEALTH_SCORER =
      AiModelType._(4, _omitEnumNames ? '' : 'AI_MODEL_TYPE_HEALTH_SCORER');
  static const AiModelType AI_MODEL_TYPE_FOOD_CALORIE_ESTIMATOR = AiModelType._(
      5, _omitEnumNames ? '' : 'AI_MODEL_TYPE_FOOD_CALORIE_ESTIMATOR');

  static const $core.List<AiModelType> values = <AiModelType>[
    AI_MODEL_TYPE_UNSPECIFIED,
    AI_MODEL_TYPE_BIOMARKER_CLASSIFIER,
    AI_MODEL_TYPE_ANOMALY_DETECTOR,
    AI_MODEL_TYPE_TREND_PREDICTOR,
    AI_MODEL_TYPE_HEALTH_SCORER,
    AI_MODEL_TYPE_FOOD_CALORIE_ESTIMATOR,
  ];

  static final $core.List<AiModelType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static AiModelType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AiModelType._(super.value, super.name);
}

class RiskLevel extends $pb.ProtobufEnum {
  static const RiskLevel RISK_LEVEL_UNSPECIFIED =
      RiskLevel._(0, _omitEnumNames ? '' : 'RISK_LEVEL_UNSPECIFIED');
  static const RiskLevel RISK_LEVEL_LOW =
      RiskLevel._(1, _omitEnumNames ? '' : 'RISK_LEVEL_LOW');
  static const RiskLevel RISK_LEVEL_MODERATE =
      RiskLevel._(2, _omitEnumNames ? '' : 'RISK_LEVEL_MODERATE');
  static const RiskLevel RISK_LEVEL_HIGH =
      RiskLevel._(3, _omitEnumNames ? '' : 'RISK_LEVEL_HIGH');
  static const RiskLevel RISK_LEVEL_CRITICAL =
      RiskLevel._(4, _omitEnumNames ? '' : 'RISK_LEVEL_CRITICAL');

  static const $core.List<RiskLevel> values = <RiskLevel>[
    RISK_LEVEL_UNSPECIFIED,
    RISK_LEVEL_LOW,
    RISK_LEVEL_MODERATE,
    RISK_LEVEL_HIGH,
    RISK_LEVEL_CRITICAL,
  ];

  static final $core.List<RiskLevel?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static RiskLevel? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RiskLevel._(super.value, super.name);
}

class CalibrationType extends $pb.ProtobufEnum {
  static const CalibrationType CALIBRATION_TYPE_UNKNOWN =
      CalibrationType._(0, _omitEnumNames ? '' : 'CALIBRATION_TYPE_UNKNOWN');
  static const CalibrationType CALIBRATION_TYPE_FACTORY =
      CalibrationType._(1, _omitEnumNames ? '' : 'CALIBRATION_TYPE_FACTORY');
  static const CalibrationType CALIBRATION_TYPE_FIELD =
      CalibrationType._(2, _omitEnumNames ? '' : 'CALIBRATION_TYPE_FIELD');
  static const CalibrationType CALIBRATION_TYPE_AUTO =
      CalibrationType._(3, _omitEnumNames ? '' : 'CALIBRATION_TYPE_AUTO');

  static const $core.List<CalibrationType> values = <CalibrationType>[
    CALIBRATION_TYPE_UNKNOWN,
    CALIBRATION_TYPE_FACTORY,
    CALIBRATION_TYPE_FIELD,
    CALIBRATION_TYPE_AUTO,
  ];

  static final $core.List<CalibrationType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static CalibrationType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CalibrationType._(super.value, super.name);
}

class CalibrationStatus extends $pb.ProtobufEnum {
  static const CalibrationStatus CALIBRATION_STATUS_UNKNOWN =
      CalibrationStatus._(
          0, _omitEnumNames ? '' : 'CALIBRATION_STATUS_UNKNOWN');
  static const CalibrationStatus CALIBRATION_STATUS_VALID =
      CalibrationStatus._(1, _omitEnumNames ? '' : 'CALIBRATION_STATUS_VALID');
  static const CalibrationStatus CALIBRATION_STATUS_EXPIRING =
      CalibrationStatus._(
          2, _omitEnumNames ? '' : 'CALIBRATION_STATUS_EXPIRING');
  static const CalibrationStatus CALIBRATION_STATUS_EXPIRED =
      CalibrationStatus._(
          3, _omitEnumNames ? '' : 'CALIBRATION_STATUS_EXPIRED');
  static const CalibrationStatus CALIBRATION_STATUS_NEEDED =
      CalibrationStatus._(4, _omitEnumNames ? '' : 'CALIBRATION_STATUS_NEEDED');

  static const $core.List<CalibrationStatus> values = <CalibrationStatus>[
    CALIBRATION_STATUS_UNKNOWN,
    CALIBRATION_STATUS_VALID,
    CALIBRATION_STATUS_EXPIRING,
    CALIBRATION_STATUS_EXPIRED,
    CALIBRATION_STATUS_NEEDED,
  ];

  static final $core.List<CalibrationStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static CalibrationStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CalibrationStatus._(super.value, super.name);
}

class GoalCategory extends $pb.ProtobufEnum {
  static const GoalCategory GOAL_CATEGORY_UNKNOWN =
      GoalCategory._(0, _omitEnumNames ? '' : 'GOAL_CATEGORY_UNKNOWN');
  static const GoalCategory GOAL_CATEGORY_BLOOD_GLUCOSE =
      GoalCategory._(1, _omitEnumNames ? '' : 'GOAL_CATEGORY_BLOOD_GLUCOSE');
  static const GoalCategory GOAL_CATEGORY_BLOOD_PRESSURE =
      GoalCategory._(2, _omitEnumNames ? '' : 'GOAL_CATEGORY_BLOOD_PRESSURE');
  static const GoalCategory GOAL_CATEGORY_CHOLESTEROL =
      GoalCategory._(3, _omitEnumNames ? '' : 'GOAL_CATEGORY_CHOLESTEROL');
  static const GoalCategory GOAL_CATEGORY_WEIGHT =
      GoalCategory._(4, _omitEnumNames ? '' : 'GOAL_CATEGORY_WEIGHT');
  static const GoalCategory GOAL_CATEGORY_EXERCISE =
      GoalCategory._(5, _omitEnumNames ? '' : 'GOAL_CATEGORY_EXERCISE');
  static const GoalCategory GOAL_CATEGORY_NUTRITION =
      GoalCategory._(6, _omitEnumNames ? '' : 'GOAL_CATEGORY_NUTRITION');
  static const GoalCategory GOAL_CATEGORY_SLEEP =
      GoalCategory._(7, _omitEnumNames ? '' : 'GOAL_CATEGORY_SLEEP');
  static const GoalCategory GOAL_CATEGORY_STRESS =
      GoalCategory._(8, _omitEnumNames ? '' : 'GOAL_CATEGORY_STRESS');
  static const GoalCategory GOAL_CATEGORY_CUSTOM =
      GoalCategory._(9, _omitEnumNames ? '' : 'GOAL_CATEGORY_CUSTOM');

  static const $core.List<GoalCategory> values = <GoalCategory>[
    GOAL_CATEGORY_UNKNOWN,
    GOAL_CATEGORY_BLOOD_GLUCOSE,
    GOAL_CATEGORY_BLOOD_PRESSURE,
    GOAL_CATEGORY_CHOLESTEROL,
    GOAL_CATEGORY_WEIGHT,
    GOAL_CATEGORY_EXERCISE,
    GOAL_CATEGORY_NUTRITION,
    GOAL_CATEGORY_SLEEP,
    GOAL_CATEGORY_STRESS,
    GOAL_CATEGORY_CUSTOM,
  ];

  static final $core.List<GoalCategory?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 9);
  static GoalCategory? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const GoalCategory._(super.value, super.name);
}

class GoalStatus extends $pb.ProtobufEnum {
  static const GoalStatus GOAL_STATUS_UNKNOWN =
      GoalStatus._(0, _omitEnumNames ? '' : 'GOAL_STATUS_UNKNOWN');
  static const GoalStatus GOAL_STATUS_ACTIVE =
      GoalStatus._(1, _omitEnumNames ? '' : 'GOAL_STATUS_ACTIVE');
  static const GoalStatus GOAL_STATUS_ACHIEVED =
      GoalStatus._(2, _omitEnumNames ? '' : 'GOAL_STATUS_ACHIEVED');
  static const GoalStatus GOAL_STATUS_PAUSED =
      GoalStatus._(3, _omitEnumNames ? '' : 'GOAL_STATUS_PAUSED');
  static const GoalStatus GOAL_STATUS_CANCELLED =
      GoalStatus._(4, _omitEnumNames ? '' : 'GOAL_STATUS_CANCELLED');

  static const $core.List<GoalStatus> values = <GoalStatus>[
    GOAL_STATUS_UNKNOWN,
    GOAL_STATUS_ACTIVE,
    GOAL_STATUS_ACHIEVED,
    GOAL_STATUS_PAUSED,
    GOAL_STATUS_CANCELLED,
  ];

  static final $core.List<GoalStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static GoalStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const GoalStatus._(super.value, super.name);
}

class CoachingType extends $pb.ProtobufEnum {
  static const CoachingType COACHING_TYPE_UNKNOWN =
      CoachingType._(0, _omitEnumNames ? '' : 'COACHING_TYPE_UNKNOWN');
  static const CoachingType COACHING_TYPE_MEASUREMENT_FEEDBACK = CoachingType._(
      1, _omitEnumNames ? '' : 'COACHING_TYPE_MEASUREMENT_FEEDBACK');
  static const CoachingType COACHING_TYPE_DAILY_TIP =
      CoachingType._(2, _omitEnumNames ? '' : 'COACHING_TYPE_DAILY_TIP');
  static const CoachingType COACHING_TYPE_GOAL_PROGRESS =
      CoachingType._(3, _omitEnumNames ? '' : 'COACHING_TYPE_GOAL_PROGRESS');
  static const CoachingType COACHING_TYPE_ALERT =
      CoachingType._(4, _omitEnumNames ? '' : 'COACHING_TYPE_ALERT');
  static const CoachingType COACHING_TYPE_MOTIVATION =
      CoachingType._(5, _omitEnumNames ? '' : 'COACHING_TYPE_MOTIVATION');
  static const CoachingType COACHING_TYPE_RECOMMENDATION =
      CoachingType._(6, _omitEnumNames ? '' : 'COACHING_TYPE_RECOMMENDATION');

  static const $core.List<CoachingType> values = <CoachingType>[
    COACHING_TYPE_UNKNOWN,
    COACHING_TYPE_MEASUREMENT_FEEDBACK,
    COACHING_TYPE_DAILY_TIP,
    COACHING_TYPE_GOAL_PROGRESS,
    COACHING_TYPE_ALERT,
    COACHING_TYPE_MOTIVATION,
    COACHING_TYPE_RECOMMENDATION,
  ];

  static final $core.List<CoachingType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 6);
  static CoachingType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CoachingType._(super.value, super.name);
}

class RecommendationType extends $pb.ProtobufEnum {
  static const RecommendationType RECOMMENDATION_TYPE_UNKNOWN =
      RecommendationType._(
          0, _omitEnumNames ? '' : 'RECOMMENDATION_TYPE_UNKNOWN');
  static const RecommendationType RECOMMENDATION_TYPE_FOOD =
      RecommendationType._(1, _omitEnumNames ? '' : 'RECOMMENDATION_TYPE_FOOD');
  static const RecommendationType RECOMMENDATION_TYPE_EXERCISE =
      RecommendationType._(
          2, _omitEnumNames ? '' : 'RECOMMENDATION_TYPE_EXERCISE');
  static const RecommendationType RECOMMENDATION_TYPE_SUPPLEMENT =
      RecommendationType._(
          3, _omitEnumNames ? '' : 'RECOMMENDATION_TYPE_SUPPLEMENT');
  static const RecommendationType RECOMMENDATION_TYPE_LIFESTYLE =
      RecommendationType._(
          4, _omitEnumNames ? '' : 'RECOMMENDATION_TYPE_LIFESTYLE');
  static const RecommendationType RECOMMENDATION_TYPE_CHECKUP =
      RecommendationType._(
          5, _omitEnumNames ? '' : 'RECOMMENDATION_TYPE_CHECKUP');

  static const $core.List<RecommendationType> values = <RecommendationType>[
    RECOMMENDATION_TYPE_UNKNOWN,
    RECOMMENDATION_TYPE_FOOD,
    RECOMMENDATION_TYPE_EXERCISE,
    RECOMMENDATION_TYPE_SUPPLEMENT,
    RECOMMENDATION_TYPE_LIFESTYLE,
    RECOMMENDATION_TYPE_CHECKUP,
  ];

  static final $core.List<RecommendationType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static RecommendationType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RecommendationType._(super.value, super.name);
}

/// 카트리지 접근 레벨
class CartridgeAccessLevel extends $pb.ProtobufEnum {
  static const CartridgeAccessLevel CARTRIDGE_ACCESS_UNKNOWN =
      CartridgeAccessLevel._(
          0, _omitEnumNames ? '' : 'CARTRIDGE_ACCESS_UNKNOWN');
  static const CartridgeAccessLevel CARTRIDGE_ACCESS_INCLUDED =
      CartridgeAccessLevel._(
          1, _omitEnumNames ? '' : 'CARTRIDGE_ACCESS_INCLUDED');
  static const CartridgeAccessLevel CARTRIDGE_ACCESS_LIMITED =
      CartridgeAccessLevel._(
          2, _omitEnumNames ? '' : 'CARTRIDGE_ACCESS_LIMITED');
  static const CartridgeAccessLevel CARTRIDGE_ACCESS_ADD_ON =
      CartridgeAccessLevel._(
          3, _omitEnumNames ? '' : 'CARTRIDGE_ACCESS_ADD_ON');
  static const CartridgeAccessLevel CARTRIDGE_ACCESS_RESTRICTED =
      CartridgeAccessLevel._(
          4, _omitEnumNames ? '' : 'CARTRIDGE_ACCESS_RESTRICTED');
  static const CartridgeAccessLevel CARTRIDGE_ACCESS_BETA =
      CartridgeAccessLevel._(5, _omitEnumNames ? '' : 'CARTRIDGE_ACCESS_BETA');

  static const $core.List<CartridgeAccessLevel> values = <CartridgeAccessLevel>[
    CARTRIDGE_ACCESS_UNKNOWN,
    CARTRIDGE_ACCESS_INCLUDED,
    CARTRIDGE_ACCESS_LIMITED,
    CARTRIDGE_ACCESS_ADD_ON,
    CARTRIDGE_ACCESS_RESTRICTED,
    CARTRIDGE_ACCESS_BETA,
  ];

  static final $core.List<CartridgeAccessLevel?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static CartridgeAccessLevel? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CartridgeAccessLevel._(super.value, super.name);
}

class FacilityType extends $pb.ProtobufEnum {
  static const FacilityType FACILITY_TYPE_UNKNOWN =
      FacilityType._(0, _omitEnumNames ? '' : 'FACILITY_TYPE_UNKNOWN');
  static const FacilityType FACILITY_TYPE_HOSPITAL =
      FacilityType._(1, _omitEnumNames ? '' : 'FACILITY_TYPE_HOSPITAL');
  static const FacilityType FACILITY_TYPE_CLINIC =
      FacilityType._(2, _omitEnumNames ? '' : 'FACILITY_TYPE_CLINIC');
  static const FacilityType FACILITY_TYPE_PHARMACY =
      FacilityType._(3, _omitEnumNames ? '' : 'FACILITY_TYPE_PHARMACY');
  static const FacilityType FACILITY_TYPE_DENTAL =
      FacilityType._(4, _omitEnumNames ? '' : 'FACILITY_TYPE_DENTAL');
  static const FacilityType FACILITY_TYPE_ORIENTAL =
      FacilityType._(5, _omitEnumNames ? '' : 'FACILITY_TYPE_ORIENTAL');

  static const $core.List<FacilityType> values = <FacilityType>[
    FACILITY_TYPE_UNKNOWN,
    FACILITY_TYPE_HOSPITAL,
    FACILITY_TYPE_CLINIC,
    FACILITY_TYPE_PHARMACY,
    FACILITY_TYPE_DENTAL,
    FACILITY_TYPE_ORIENTAL,
  ];

  static final $core.List<FacilityType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static FacilityType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FacilityType._(super.value, super.name);
}

class DoctorSpecialty extends $pb.ProtobufEnum {
  static const DoctorSpecialty DOCTOR_SPECIALTY_UNKNOWN =
      DoctorSpecialty._(0, _omitEnumNames ? '' : 'DOCTOR_SPECIALTY_UNKNOWN');
  static const DoctorSpecialty DOCTOR_SPECIALTY_GENERAL =
      DoctorSpecialty._(1, _omitEnumNames ? '' : 'DOCTOR_SPECIALTY_GENERAL');
  static const DoctorSpecialty DOCTOR_SPECIALTY_INTERNAL =
      DoctorSpecialty._(2, _omitEnumNames ? '' : 'DOCTOR_SPECIALTY_INTERNAL');
  static const DoctorSpecialty DOCTOR_SPECIALTY_CARDIOLOGY =
      DoctorSpecialty._(3, _omitEnumNames ? '' : 'DOCTOR_SPECIALTY_CARDIOLOGY');
  static const DoctorSpecialty DOCTOR_SPECIALTY_ENDOCRINOLOGY =
      DoctorSpecialty._(
          4, _omitEnumNames ? '' : 'DOCTOR_SPECIALTY_ENDOCRINOLOGY');
  static const DoctorSpecialty DOCTOR_SPECIALTY_DERMATOLOGY = DoctorSpecialty._(
      5, _omitEnumNames ? '' : 'DOCTOR_SPECIALTY_DERMATOLOGY');
  static const DoctorSpecialty DOCTOR_SPECIALTY_PEDIATRICS =
      DoctorSpecialty._(6, _omitEnumNames ? '' : 'DOCTOR_SPECIALTY_PEDIATRICS');
  static const DoctorSpecialty DOCTOR_SPECIALTY_PSYCHIATRY =
      DoctorSpecialty._(7, _omitEnumNames ? '' : 'DOCTOR_SPECIALTY_PSYCHIATRY');
  static const DoctorSpecialty DOCTOR_SPECIALTY_ORTHOPEDICS = DoctorSpecialty._(
      8, _omitEnumNames ? '' : 'DOCTOR_SPECIALTY_ORTHOPEDICS');
  static const DoctorSpecialty DOCTOR_SPECIALTY_OPHTHALMOLOGY =
      DoctorSpecialty._(
          9, _omitEnumNames ? '' : 'DOCTOR_SPECIALTY_OPHTHALMOLOGY');
  static const DoctorSpecialty DOCTOR_SPECIALTY_ENT =
      DoctorSpecialty._(10, _omitEnumNames ? '' : 'DOCTOR_SPECIALTY_ENT');

  static const $core.List<DoctorSpecialty> values = <DoctorSpecialty>[
    DOCTOR_SPECIALTY_UNKNOWN,
    DOCTOR_SPECIALTY_GENERAL,
    DOCTOR_SPECIALTY_INTERNAL,
    DOCTOR_SPECIALTY_CARDIOLOGY,
    DOCTOR_SPECIALTY_ENDOCRINOLOGY,
    DOCTOR_SPECIALTY_DERMATOLOGY,
    DOCTOR_SPECIALTY_PEDIATRICS,
    DOCTOR_SPECIALTY_PSYCHIATRY,
    DOCTOR_SPECIALTY_ORTHOPEDICS,
    DOCTOR_SPECIALTY_OPHTHALMOLOGY,
    DOCTOR_SPECIALTY_ENT,
  ];

  static final $core.List<DoctorSpecialty?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 10);
  static DoctorSpecialty? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DoctorSpecialty._(super.value, super.name);
}

class ReservationStatus extends $pb.ProtobufEnum {
  static const ReservationStatus RESERVATION_STATUS_UNKNOWN =
      ReservationStatus._(
          0, _omitEnumNames ? '' : 'RESERVATION_STATUS_UNKNOWN');
  static const ReservationStatus RESERVATION_STATUS_PENDING =
      ReservationStatus._(
          1, _omitEnumNames ? '' : 'RESERVATION_STATUS_PENDING');
  static const ReservationStatus RESERVATION_STATUS_CONFIRMED =
      ReservationStatus._(
          2, _omitEnumNames ? '' : 'RESERVATION_STATUS_CONFIRMED');
  static const ReservationStatus RESERVATION_STATUS_COMPLETED =
      ReservationStatus._(
          3, _omitEnumNames ? '' : 'RESERVATION_STATUS_COMPLETED');
  static const ReservationStatus RESERVATION_STATUS_CANCELLED =
      ReservationStatus._(
          4, _omitEnumNames ? '' : 'RESERVATION_STATUS_CANCELLED');
  static const ReservationStatus RESERVATION_STATUS_NO_SHOW =
      ReservationStatus._(
          5, _omitEnumNames ? '' : 'RESERVATION_STATUS_NO_SHOW');

  static const $core.List<ReservationStatus> values = <ReservationStatus>[
    RESERVATION_STATUS_UNKNOWN,
    RESERVATION_STATUS_PENDING,
    RESERVATION_STATUS_CONFIRMED,
    RESERVATION_STATUS_COMPLETED,
    RESERVATION_STATUS_CANCELLED,
    RESERVATION_STATUS_NO_SHOW,
  ];

  static final $core.List<ReservationStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static ReservationStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ReservationStatus._(super.value, super.name);
}

class AdminRole extends $pb.ProtobufEnum {
  static const AdminRole ADMIN_ROLE_UNKNOWN =
      AdminRole._(0, _omitEnumNames ? '' : 'ADMIN_ROLE_UNKNOWN');
  static const AdminRole ADMIN_ROLE_SUPER =
      AdminRole._(1, _omitEnumNames ? '' : 'ADMIN_ROLE_SUPER');
  static const AdminRole ADMIN_ROLE_NATIONAL =
      AdminRole._(2, _omitEnumNames ? '' : 'ADMIN_ROLE_NATIONAL');
  static const AdminRole ADMIN_ROLE_REGIONAL =
      AdminRole._(3, _omitEnumNames ? '' : 'ADMIN_ROLE_REGIONAL');
  static const AdminRole ADMIN_ROLE_BRANCH =
      AdminRole._(4, _omitEnumNames ? '' : 'ADMIN_ROLE_BRANCH');
  static const AdminRole ADMIN_ROLE_STORE =
      AdminRole._(5, _omitEnumNames ? '' : 'ADMIN_ROLE_STORE');

  static const $core.List<AdminRole> values = <AdminRole>[
    ADMIN_ROLE_UNKNOWN,
    ADMIN_ROLE_SUPER,
    ADMIN_ROLE_NATIONAL,
    ADMIN_ROLE_REGIONAL,
    ADMIN_ROLE_BRANCH,
    ADMIN_ROLE_STORE,
  ];

  static final $core.List<AdminRole?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static AdminRole? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AdminRole._(super.value, super.name);
}

class AuditAction extends $pb.ProtobufEnum {
  static const AuditAction AUDIT_ACTION_UNKNOWN =
      AuditAction._(0, _omitEnumNames ? '' : 'AUDIT_ACTION_UNKNOWN');
  static const AuditAction AUDIT_ACTION_LOGIN =
      AuditAction._(1, _omitEnumNames ? '' : 'AUDIT_ACTION_LOGIN');
  static const AuditAction AUDIT_ACTION_LOGOUT =
      AuditAction._(2, _omitEnumNames ? '' : 'AUDIT_ACTION_LOGOUT');
  static const AuditAction AUDIT_ACTION_CREATE =
      AuditAction._(3, _omitEnumNames ? '' : 'AUDIT_ACTION_CREATE');
  static const AuditAction AUDIT_ACTION_UPDATE =
      AuditAction._(4, _omitEnumNames ? '' : 'AUDIT_ACTION_UPDATE');
  static const AuditAction AUDIT_ACTION_DELETE =
      AuditAction._(5, _omitEnumNames ? '' : 'AUDIT_ACTION_DELETE');
  static const AuditAction AUDIT_ACTION_CONFIG_CHANGE =
      AuditAction._(6, _omitEnumNames ? '' : 'AUDIT_ACTION_CONFIG_CHANGE');
  static const AuditAction AUDIT_ACTION_ROLE_CHANGE =
      AuditAction._(7, _omitEnumNames ? '' : 'AUDIT_ACTION_ROLE_CHANGE');

  static const $core.List<AuditAction> values = <AuditAction>[
    AUDIT_ACTION_UNKNOWN,
    AUDIT_ACTION_LOGIN,
    AUDIT_ACTION_LOGOUT,
    AUDIT_ACTION_CREATE,
    AUDIT_ACTION_UPDATE,
    AUDIT_ACTION_DELETE,
    AUDIT_ACTION_CONFIG_CHANGE,
    AUDIT_ACTION_ROLE_CHANGE,
  ];

  static final $core.List<AuditAction?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static AuditAction? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AuditAction._(super.value, super.name);
}

class FamilyRole extends $pb.ProtobufEnum {
  static const FamilyRole FAMILY_ROLE_UNKNOWN =
      FamilyRole._(0, _omitEnumNames ? '' : 'FAMILY_ROLE_UNKNOWN');
  static const FamilyRole FAMILY_ROLE_OWNER =
      FamilyRole._(1, _omitEnumNames ? '' : 'FAMILY_ROLE_OWNER');
  static const FamilyRole FAMILY_ROLE_GUARDIAN =
      FamilyRole._(2, _omitEnumNames ? '' : 'FAMILY_ROLE_GUARDIAN');
  static const FamilyRole FAMILY_ROLE_MEMBER =
      FamilyRole._(3, _omitEnumNames ? '' : 'FAMILY_ROLE_MEMBER');
  static const FamilyRole FAMILY_ROLE_CHILD =
      FamilyRole._(4, _omitEnumNames ? '' : 'FAMILY_ROLE_CHILD');
  static const FamilyRole FAMILY_ROLE_ELDERLY =
      FamilyRole._(5, _omitEnumNames ? '' : 'FAMILY_ROLE_ELDERLY');

  static const $core.List<FamilyRole> values = <FamilyRole>[
    FAMILY_ROLE_UNKNOWN,
    FAMILY_ROLE_OWNER,
    FAMILY_ROLE_GUARDIAN,
    FAMILY_ROLE_MEMBER,
    FAMILY_ROLE_CHILD,
    FAMILY_ROLE_ELDERLY,
  ];

  static final $core.List<FamilyRole?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static FamilyRole? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FamilyRole._(super.value, super.name);
}

class InvitationStatus extends $pb.ProtobufEnum {
  static const InvitationStatus INVITATION_STATUS_UNKNOWN =
      InvitationStatus._(0, _omitEnumNames ? '' : 'INVITATION_STATUS_UNKNOWN');
  static const InvitationStatus INVITATION_STATUS_PENDING =
      InvitationStatus._(1, _omitEnumNames ? '' : 'INVITATION_STATUS_PENDING');
  static const InvitationStatus INVITATION_STATUS_ACCEPTED =
      InvitationStatus._(2, _omitEnumNames ? '' : 'INVITATION_STATUS_ACCEPTED');
  static const InvitationStatus INVITATION_STATUS_REJECTED =
      InvitationStatus._(3, _omitEnumNames ? '' : 'INVITATION_STATUS_REJECTED');
  static const InvitationStatus INVITATION_STATUS_EXPIRED =
      InvitationStatus._(4, _omitEnumNames ? '' : 'INVITATION_STATUS_EXPIRED');

  static const $core.List<InvitationStatus> values = <InvitationStatus>[
    INVITATION_STATUS_UNKNOWN,
    INVITATION_STATUS_PENDING,
    INVITATION_STATUS_ACCEPTED,
    INVITATION_STATUS_REJECTED,
    INVITATION_STATUS_EXPIRED,
  ];

  static final $core.List<InvitationStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static InvitationStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const InvitationStatus._(super.value, super.name);
}

class HealthRecordType extends $pb.ProtobufEnum {
  static const HealthRecordType HEALTH_RECORD_TYPE_UNKNOWN =
      HealthRecordType._(0, _omitEnumNames ? '' : 'HEALTH_RECORD_TYPE_UNKNOWN');
  static const HealthRecordType HEALTH_RECORD_TYPE_LAB_RESULT =
      HealthRecordType._(
          1, _omitEnumNames ? '' : 'HEALTH_RECORD_TYPE_LAB_RESULT');
  static const HealthRecordType HEALTH_RECORD_TYPE_IMAGING =
      HealthRecordType._(2, _omitEnumNames ? '' : 'HEALTH_RECORD_TYPE_IMAGING');
  static const HealthRecordType HEALTH_RECORD_TYPE_VITAL_SIGN =
      HealthRecordType._(
          3, _omitEnumNames ? '' : 'HEALTH_RECORD_TYPE_VITAL_SIGN');
  static const HealthRecordType HEALTH_RECORD_TYPE_ALLERGY =
      HealthRecordType._(4, _omitEnumNames ? '' : 'HEALTH_RECORD_TYPE_ALLERGY');
  static const HealthRecordType HEALTH_RECORD_TYPE_CONDITION =
      HealthRecordType._(
          5, _omitEnumNames ? '' : 'HEALTH_RECORD_TYPE_CONDITION');
  static const HealthRecordType HEALTH_RECORD_TYPE_IMMUNIZATION =
      HealthRecordType._(
          6, _omitEnumNames ? '' : 'HEALTH_RECORD_TYPE_IMMUNIZATION');
  static const HealthRecordType HEALTH_RECORD_TYPE_PROCEDURE =
      HealthRecordType._(
          7, _omitEnumNames ? '' : 'HEALTH_RECORD_TYPE_PROCEDURE');

  static const $core.List<HealthRecordType> values = <HealthRecordType>[
    HEALTH_RECORD_TYPE_UNKNOWN,
    HEALTH_RECORD_TYPE_LAB_RESULT,
    HEALTH_RECORD_TYPE_IMAGING,
    HEALTH_RECORD_TYPE_VITAL_SIGN,
    HEALTH_RECORD_TYPE_ALLERGY,
    HEALTH_RECORD_TYPE_CONDITION,
    HEALTH_RECORD_TYPE_IMMUNIZATION,
    HEALTH_RECORD_TYPE_PROCEDURE,
  ];

  static final $core.List<HealthRecordType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static HealthRecordType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const HealthRecordType._(super.value, super.name);
}

class FHIRResourceType extends $pb.ProtobufEnum {
  static const FHIRResourceType FHIR_RESOURCE_TYPE_UNKNOWN =
      FHIRResourceType._(0, _omitEnumNames ? '' : 'FHIR_RESOURCE_TYPE_UNKNOWN');
  static const FHIRResourceType FHIR_RESOURCE_TYPE_OBSERVATION =
      FHIRResourceType._(
          1, _omitEnumNames ? '' : 'FHIR_RESOURCE_TYPE_OBSERVATION');
  static const FHIRResourceType FHIR_RESOURCE_TYPE_CONDITION =
      FHIRResourceType._(
          2, _omitEnumNames ? '' : 'FHIR_RESOURCE_TYPE_CONDITION');
  static const FHIRResourceType FHIR_RESOURCE_TYPE_ALLERGY_INTOLERANCE =
      FHIRResourceType._(
          3, _omitEnumNames ? '' : 'FHIR_RESOURCE_TYPE_ALLERGY_INTOLERANCE');
  static const FHIRResourceType FHIR_RESOURCE_TYPE_IMMUNIZATION =
      FHIRResourceType._(
          4, _omitEnumNames ? '' : 'FHIR_RESOURCE_TYPE_IMMUNIZATION');
  static const FHIRResourceType FHIR_RESOURCE_TYPE_PROCEDURE =
      FHIRResourceType._(
          5, _omitEnumNames ? '' : 'FHIR_RESOURCE_TYPE_PROCEDURE');
  static const FHIRResourceType FHIR_RESOURCE_TYPE_DIAGNOSTIC_REPORT =
      FHIRResourceType._(
          6, _omitEnumNames ? '' : 'FHIR_RESOURCE_TYPE_DIAGNOSTIC_REPORT');
  static const FHIRResourceType FHIR_RESOURCE_TYPE_BUNDLE =
      FHIRResourceType._(7, _omitEnumNames ? '' : 'FHIR_RESOURCE_TYPE_BUNDLE');

  static const $core.List<FHIRResourceType> values = <FHIRResourceType>[
    FHIR_RESOURCE_TYPE_UNKNOWN,
    FHIR_RESOURCE_TYPE_OBSERVATION,
    FHIR_RESOURCE_TYPE_CONDITION,
    FHIR_RESOURCE_TYPE_ALLERGY_INTOLERANCE,
    FHIR_RESOURCE_TYPE_IMMUNIZATION,
    FHIR_RESOURCE_TYPE_PROCEDURE,
    FHIR_RESOURCE_TYPE_DIAGNOSTIC_REPORT,
    FHIR_RESOURCE_TYPE_BUNDLE,
  ];

  static final $core.List<FHIRResourceType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static FHIRResourceType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FHIRResourceType._(super.value, super.name);
}

class PrescriptionStatus extends $pb.ProtobufEnum {
  static const PrescriptionStatus PRESCRIPTION_STATUS_UNKNOWN =
      PrescriptionStatus._(
          0, _omitEnumNames ? '' : 'PRESCRIPTION_STATUS_UNKNOWN');
  static const PrescriptionStatus PRESCRIPTION_STATUS_DRAFT =
      PrescriptionStatus._(
          1, _omitEnumNames ? '' : 'PRESCRIPTION_STATUS_DRAFT');
  static const PrescriptionStatus PRESCRIPTION_STATUS_ACTIVE =
      PrescriptionStatus._(
          2, _omitEnumNames ? '' : 'PRESCRIPTION_STATUS_ACTIVE');
  static const PrescriptionStatus PRESCRIPTION_STATUS_DISPENSED =
      PrescriptionStatus._(
          3, _omitEnumNames ? '' : 'PRESCRIPTION_STATUS_DISPENSED');
  static const PrescriptionStatus PRESCRIPTION_STATUS_COMPLETED =
      PrescriptionStatus._(
          4, _omitEnumNames ? '' : 'PRESCRIPTION_STATUS_COMPLETED');
  static const PrescriptionStatus PRESCRIPTION_STATUS_CANCELLED =
      PrescriptionStatus._(
          5, _omitEnumNames ? '' : 'PRESCRIPTION_STATUS_CANCELLED');
  static const PrescriptionStatus PRESCRIPTION_STATUS_EXPIRED =
      PrescriptionStatus._(
          6, _omitEnumNames ? '' : 'PRESCRIPTION_STATUS_EXPIRED');

  static const $core.List<PrescriptionStatus> values = <PrescriptionStatus>[
    PRESCRIPTION_STATUS_UNKNOWN,
    PRESCRIPTION_STATUS_DRAFT,
    PRESCRIPTION_STATUS_ACTIVE,
    PRESCRIPTION_STATUS_DISPENSED,
    PRESCRIPTION_STATUS_COMPLETED,
    PRESCRIPTION_STATUS_CANCELLED,
    PRESCRIPTION_STATUS_EXPIRED,
  ];

  static final $core.List<PrescriptionStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 6);
  static PrescriptionStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PrescriptionStatus._(super.value, super.name);
}

class DrugInteractionSeverity extends $pb.ProtobufEnum {
  static const DrugInteractionSeverity DRUG_INTERACTION_SEVERITY_UNKNOWN =
      DrugInteractionSeverity._(
          0, _omitEnumNames ? '' : 'DRUG_INTERACTION_SEVERITY_UNKNOWN');
  static const DrugInteractionSeverity DRUG_INTERACTION_SEVERITY_NONE =
      DrugInteractionSeverity._(
          1, _omitEnumNames ? '' : 'DRUG_INTERACTION_SEVERITY_NONE');
  static const DrugInteractionSeverity DRUG_INTERACTION_SEVERITY_MINOR =
      DrugInteractionSeverity._(
          2, _omitEnumNames ? '' : 'DRUG_INTERACTION_SEVERITY_MINOR');
  static const DrugInteractionSeverity DRUG_INTERACTION_SEVERITY_MODERATE =
      DrugInteractionSeverity._(
          3, _omitEnumNames ? '' : 'DRUG_INTERACTION_SEVERITY_MODERATE');
  static const DrugInteractionSeverity DRUG_INTERACTION_SEVERITY_MAJOR =
      DrugInteractionSeverity._(
          4, _omitEnumNames ? '' : 'DRUG_INTERACTION_SEVERITY_MAJOR');
  static const DrugInteractionSeverity
      DRUG_INTERACTION_SEVERITY_CONTRAINDICATED = DrugInteractionSeverity._(
          5, _omitEnumNames ? '' : 'DRUG_INTERACTION_SEVERITY_CONTRAINDICATED');

  static const $core.List<DrugInteractionSeverity> values =
      <DrugInteractionSeverity>[
    DRUG_INTERACTION_SEVERITY_UNKNOWN,
    DRUG_INTERACTION_SEVERITY_NONE,
    DRUG_INTERACTION_SEVERITY_MINOR,
    DRUG_INTERACTION_SEVERITY_MODERATE,
    DRUG_INTERACTION_SEVERITY_MAJOR,
    DRUG_INTERACTION_SEVERITY_CONTRAINDICATED,
  ];

  static final $core.List<DrugInteractionSeverity?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static DrugInteractionSeverity? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DrugInteractionSeverity._(super.value, super.name);
}

class PostCategory extends $pb.ProtobufEnum {
  static const PostCategory POST_CATEGORY_UNKNOWN =
      PostCategory._(0, _omitEnumNames ? '' : 'POST_CATEGORY_UNKNOWN');
  static const PostCategory POST_CATEGORY_GENERAL =
      PostCategory._(1, _omitEnumNames ? '' : 'POST_CATEGORY_GENERAL');
  static const PostCategory POST_CATEGORY_HEALTH_TIP =
      PostCategory._(2, _omitEnumNames ? '' : 'POST_CATEGORY_HEALTH_TIP');
  static const PostCategory POST_CATEGORY_QNA =
      PostCategory._(3, _omitEnumNames ? '' : 'POST_CATEGORY_QNA');
  static const PostCategory POST_CATEGORY_EXPERIENCE =
      PostCategory._(4, _omitEnumNames ? '' : 'POST_CATEGORY_EXPERIENCE');
  static const PostCategory POST_CATEGORY_RECIPE =
      PostCategory._(5, _omitEnumNames ? '' : 'POST_CATEGORY_RECIPE');
  static const PostCategory POST_CATEGORY_EXERCISE =
      PostCategory._(6, _omitEnumNames ? '' : 'POST_CATEGORY_EXERCISE');

  static const $core.List<PostCategory> values = <PostCategory>[
    POST_CATEGORY_UNKNOWN,
    POST_CATEGORY_GENERAL,
    POST_CATEGORY_HEALTH_TIP,
    POST_CATEGORY_QNA,
    POST_CATEGORY_EXPERIENCE,
    POST_CATEGORY_RECIPE,
    POST_CATEGORY_EXERCISE,
  ];

  static final $core.List<PostCategory?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 6);
  static PostCategory? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const PostCategory._(super.value, super.name);
}

class ChallengeStatus extends $pb.ProtobufEnum {
  static const ChallengeStatus CHALLENGE_STATUS_UNKNOWN =
      ChallengeStatus._(0, _omitEnumNames ? '' : 'CHALLENGE_STATUS_UNKNOWN');
  static const ChallengeStatus CHALLENGE_STATUS_UPCOMING =
      ChallengeStatus._(1, _omitEnumNames ? '' : 'CHALLENGE_STATUS_UPCOMING');
  static const ChallengeStatus CHALLENGE_STATUS_ACTIVE =
      ChallengeStatus._(2, _omitEnumNames ? '' : 'CHALLENGE_STATUS_ACTIVE');
  static const ChallengeStatus CHALLENGE_STATUS_COMPLETED =
      ChallengeStatus._(3, _omitEnumNames ? '' : 'CHALLENGE_STATUS_COMPLETED');
  static const ChallengeStatus CHALLENGE_STATUS_CANCELLED =
      ChallengeStatus._(4, _omitEnumNames ? '' : 'CHALLENGE_STATUS_CANCELLED');

  static const $core.List<ChallengeStatus> values = <ChallengeStatus>[
    CHALLENGE_STATUS_UNKNOWN,
    CHALLENGE_STATUS_UPCOMING,
    CHALLENGE_STATUS_ACTIVE,
    CHALLENGE_STATUS_COMPLETED,
    CHALLENGE_STATUS_CANCELLED,
  ];

  static final $core.List<ChallengeStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static ChallengeStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ChallengeStatus._(super.value, super.name);
}

class ChallengeType extends $pb.ProtobufEnum {
  static const ChallengeType CHALLENGE_TYPE_UNKNOWN =
      ChallengeType._(0, _omitEnumNames ? '' : 'CHALLENGE_TYPE_UNKNOWN');
  static const ChallengeType CHALLENGE_TYPE_STEPS =
      ChallengeType._(1, _omitEnumNames ? '' : 'CHALLENGE_TYPE_STEPS');
  static const ChallengeType CHALLENGE_TYPE_MEASUREMENT =
      ChallengeType._(2, _omitEnumNames ? '' : 'CHALLENGE_TYPE_MEASUREMENT');
  static const ChallengeType CHALLENGE_TYPE_DIET =
      ChallengeType._(3, _omitEnumNames ? '' : 'CHALLENGE_TYPE_DIET');
  static const ChallengeType CHALLENGE_TYPE_EXERCISE =
      ChallengeType._(4, _omitEnumNames ? '' : 'CHALLENGE_TYPE_EXERCISE');
  static const ChallengeType CHALLENGE_TYPE_SLEEP =
      ChallengeType._(5, _omitEnumNames ? '' : 'CHALLENGE_TYPE_SLEEP');

  static const $core.List<ChallengeType> values = <ChallengeType>[
    CHALLENGE_TYPE_UNKNOWN,
    CHALLENGE_TYPE_STEPS,
    CHALLENGE_TYPE_MEASUREMENT,
    CHALLENGE_TYPE_DIET,
    CHALLENGE_TYPE_EXERCISE,
    CHALLENGE_TYPE_SLEEP,
  ];

  static final $core.List<ChallengeType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static ChallengeType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ChallengeType._(super.value, super.name);
}

class RoomType extends $pb.ProtobufEnum {
  static const RoomType ROOM_TYPE_UNKNOWN =
      RoomType._(0, _omitEnumNames ? '' : 'ROOM_TYPE_UNKNOWN');
  static const RoomType ROOM_TYPE_ONE_TO_ONE =
      RoomType._(1, _omitEnumNames ? '' : 'ROOM_TYPE_ONE_TO_ONE');
  static const RoomType ROOM_TYPE_GROUP =
      RoomType._(2, _omitEnumNames ? '' : 'ROOM_TYPE_GROUP');
  static const RoomType ROOM_TYPE_WEBINAR =
      RoomType._(3, _omitEnumNames ? '' : 'ROOM_TYPE_WEBINAR');
  static const RoomType ROOM_TYPE_CONSULTATION =
      RoomType._(4, _omitEnumNames ? '' : 'ROOM_TYPE_CONSULTATION');

  static const $core.List<RoomType> values = <RoomType>[
    ROOM_TYPE_UNKNOWN,
    ROOM_TYPE_ONE_TO_ONE,
    ROOM_TYPE_GROUP,
    ROOM_TYPE_WEBINAR,
    ROOM_TYPE_CONSULTATION,
  ];

  static final $core.List<RoomType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static RoomType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RoomType._(super.value, super.name);
}

class RoomStatus extends $pb.ProtobufEnum {
  static const RoomStatus ROOM_STATUS_UNKNOWN =
      RoomStatus._(0, _omitEnumNames ? '' : 'ROOM_STATUS_UNKNOWN');
  static const RoomStatus ROOM_STATUS_WAITING =
      RoomStatus._(1, _omitEnumNames ? '' : 'ROOM_STATUS_WAITING');
  static const RoomStatus ROOM_STATUS_ACTIVE =
      RoomStatus._(2, _omitEnumNames ? '' : 'ROOM_STATUS_ACTIVE');
  static const RoomStatus ROOM_STATUS_ENDED =
      RoomStatus._(3, _omitEnumNames ? '' : 'ROOM_STATUS_ENDED');
  static const RoomStatus ROOM_STATUS_FAILED =
      RoomStatus._(4, _omitEnumNames ? '' : 'ROOM_STATUS_FAILED');

  static const $core.List<RoomStatus> values = <RoomStatus>[
    ROOM_STATUS_UNKNOWN,
    ROOM_STATUS_WAITING,
    ROOM_STATUS_ACTIVE,
    ROOM_STATUS_ENDED,
    ROOM_STATUS_FAILED,
  ];

  static final $core.List<RoomStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static RoomStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RoomStatus._(super.value, super.name);
}

class SignalType extends $pb.ProtobufEnum {
  static const SignalType SIGNAL_TYPE_UNKNOWN =
      SignalType._(0, _omitEnumNames ? '' : 'SIGNAL_TYPE_UNKNOWN');
  static const SignalType SIGNAL_TYPE_OFFER =
      SignalType._(1, _omitEnumNames ? '' : 'SIGNAL_TYPE_OFFER');
  static const SignalType SIGNAL_TYPE_ANSWER =
      SignalType._(2, _omitEnumNames ? '' : 'SIGNAL_TYPE_ANSWER');
  static const SignalType SIGNAL_TYPE_ICE_CANDIDATE =
      SignalType._(3, _omitEnumNames ? '' : 'SIGNAL_TYPE_ICE_CANDIDATE');
  static const SignalType SIGNAL_TYPE_RENEGOTIATE =
      SignalType._(4, _omitEnumNames ? '' : 'SIGNAL_TYPE_RENEGOTIATE');
  static const SignalType SIGNAL_TYPE_MUTE =
      SignalType._(5, _omitEnumNames ? '' : 'SIGNAL_TYPE_MUTE');
  static const SignalType SIGNAL_TYPE_UNMUTE =
      SignalType._(6, _omitEnumNames ? '' : 'SIGNAL_TYPE_UNMUTE');

  static const $core.List<SignalType> values = <SignalType>[
    SIGNAL_TYPE_UNKNOWN,
    SIGNAL_TYPE_OFFER,
    SIGNAL_TYPE_ANSWER,
    SIGNAL_TYPE_ICE_CANDIDATE,
    SIGNAL_TYPE_RENEGOTIATE,
    SIGNAL_TYPE_MUTE,
    SIGNAL_TYPE_UNMUTE,
  ];

  static final $core.List<SignalType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 6);
  static SignalType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SignalType._(super.value, super.name);
}

class NotificationType extends $pb.ProtobufEnum {
  static const NotificationType NOTIFICATION_TYPE_UNKNOWN =
      NotificationType._(0, _omitEnumNames ? '' : 'NOTIFICATION_TYPE_UNKNOWN');
  static const NotificationType NOTIFICATION_TYPE_MEASUREMENT =
      NotificationType._(
          1, _omitEnumNames ? '' : 'NOTIFICATION_TYPE_MEASUREMENT');
  static const NotificationType NOTIFICATION_TYPE_HEALTH_ALERT =
      NotificationType._(
          2, _omitEnumNames ? '' : 'NOTIFICATION_TYPE_HEALTH_ALERT');
  static const NotificationType NOTIFICATION_TYPE_APPOINTMENT =
      NotificationType._(
          3, _omitEnumNames ? '' : 'NOTIFICATION_TYPE_APPOINTMENT');
  static const NotificationType NOTIFICATION_TYPE_PRESCRIPTION =
      NotificationType._(
          4, _omitEnumNames ? '' : 'NOTIFICATION_TYPE_PRESCRIPTION');
  static const NotificationType NOTIFICATION_TYPE_COMMUNITY =
      NotificationType._(
          5, _omitEnumNames ? '' : 'NOTIFICATION_TYPE_COMMUNITY');
  static const NotificationType NOTIFICATION_TYPE_SYSTEM =
      NotificationType._(6, _omitEnumNames ? '' : 'NOTIFICATION_TYPE_SYSTEM');
  static const NotificationType NOTIFICATION_TYPE_PROMOTION =
      NotificationType._(
          7, _omitEnumNames ? '' : 'NOTIFICATION_TYPE_PROMOTION');

  static const $core.List<NotificationType> values = <NotificationType>[
    NOTIFICATION_TYPE_UNKNOWN,
    NOTIFICATION_TYPE_MEASUREMENT,
    NOTIFICATION_TYPE_HEALTH_ALERT,
    NOTIFICATION_TYPE_APPOINTMENT,
    NOTIFICATION_TYPE_PRESCRIPTION,
    NOTIFICATION_TYPE_COMMUNITY,
    NOTIFICATION_TYPE_SYSTEM,
    NOTIFICATION_TYPE_PROMOTION,
  ];

  static final $core.List<NotificationType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static NotificationType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const NotificationType._(super.value, super.name);
}

class NotificationChannel extends $pb.ProtobufEnum {
  static const NotificationChannel NOTIFICATION_CHANNEL_UNKNOWN =
      NotificationChannel._(
          0, _omitEnumNames ? '' : 'NOTIFICATION_CHANNEL_UNKNOWN');
  static const NotificationChannel NOTIFICATION_CHANNEL_PUSH =
      NotificationChannel._(
          1, _omitEnumNames ? '' : 'NOTIFICATION_CHANNEL_PUSH');
  static const NotificationChannel NOTIFICATION_CHANNEL_EMAIL =
      NotificationChannel._(
          2, _omitEnumNames ? '' : 'NOTIFICATION_CHANNEL_EMAIL');
  static const NotificationChannel NOTIFICATION_CHANNEL_SMS =
      NotificationChannel._(
          3, _omitEnumNames ? '' : 'NOTIFICATION_CHANNEL_SMS');
  static const NotificationChannel NOTIFICATION_CHANNEL_IN_APP =
      NotificationChannel._(
          4, _omitEnumNames ? '' : 'NOTIFICATION_CHANNEL_IN_APP');

  static const $core.List<NotificationChannel> values = <NotificationChannel>[
    NOTIFICATION_CHANNEL_UNKNOWN,
    NOTIFICATION_CHANNEL_PUSH,
    NOTIFICATION_CHANNEL_EMAIL,
    NOTIFICATION_CHANNEL_SMS,
    NOTIFICATION_CHANNEL_IN_APP,
  ];

  static final $core.List<NotificationChannel?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static NotificationChannel? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const NotificationChannel._(super.value, super.name);
}

class NotificationPriority extends $pb.ProtobufEnum {
  static const NotificationPriority NOTIFICATION_PRIORITY_UNKNOWN =
      NotificationPriority._(
          0, _omitEnumNames ? '' : 'NOTIFICATION_PRIORITY_UNKNOWN');
  static const NotificationPriority NOTIFICATION_PRIORITY_LOW =
      NotificationPriority._(
          1, _omitEnumNames ? '' : 'NOTIFICATION_PRIORITY_LOW');
  static const NotificationPriority NOTIFICATION_PRIORITY_NORMAL =
      NotificationPriority._(
          2, _omitEnumNames ? '' : 'NOTIFICATION_PRIORITY_NORMAL');
  static const NotificationPriority NOTIFICATION_PRIORITY_HIGH =
      NotificationPriority._(
          3, _omitEnumNames ? '' : 'NOTIFICATION_PRIORITY_HIGH');
  static const NotificationPriority NOTIFICATION_PRIORITY_URGENT =
      NotificationPriority._(
          4, _omitEnumNames ? '' : 'NOTIFICATION_PRIORITY_URGENT');

  static const $core.List<NotificationPriority> values = <NotificationPriority>[
    NOTIFICATION_PRIORITY_UNKNOWN,
    NOTIFICATION_PRIORITY_LOW,
    NOTIFICATION_PRIORITY_NORMAL,
    NOTIFICATION_PRIORITY_HIGH,
    NOTIFICATION_PRIORITY_URGENT,
  ];

  static final $core.List<NotificationPriority?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static NotificationPriority? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const NotificationPriority._(super.value, super.name);
}

class ConsultationStatus extends $pb.ProtobufEnum {
  static const ConsultationStatus CONSULTATION_STATUS_UNKNOWN =
      ConsultationStatus._(
          0, _omitEnumNames ? '' : 'CONSULTATION_STATUS_UNKNOWN');
  static const ConsultationStatus CONSULTATION_STATUS_REQUESTED =
      ConsultationStatus._(
          1, _omitEnumNames ? '' : 'CONSULTATION_STATUS_REQUESTED');
  static const ConsultationStatus CONSULTATION_STATUS_MATCHED =
      ConsultationStatus._(
          2, _omitEnumNames ? '' : 'CONSULTATION_STATUS_MATCHED');
  static const ConsultationStatus CONSULTATION_STATUS_SCHEDULED =
      ConsultationStatus._(
          3, _omitEnumNames ? '' : 'CONSULTATION_STATUS_SCHEDULED');
  static const ConsultationStatus CONSULTATION_STATUS_IN_PROGRESS =
      ConsultationStatus._(
          4, _omitEnumNames ? '' : 'CONSULTATION_STATUS_IN_PROGRESS');
  static const ConsultationStatus CONSULTATION_STATUS_COMPLETED =
      ConsultationStatus._(
          5, _omitEnumNames ? '' : 'CONSULTATION_STATUS_COMPLETED');
  static const ConsultationStatus CONSULTATION_STATUS_CANCELLED =
      ConsultationStatus._(
          6, _omitEnumNames ? '' : 'CONSULTATION_STATUS_CANCELLED');
  static const ConsultationStatus CONSULTATION_STATUS_NO_SHOW =
      ConsultationStatus._(
          7, _omitEnumNames ? '' : 'CONSULTATION_STATUS_NO_SHOW');

  static const $core.List<ConsultationStatus> values = <ConsultationStatus>[
    CONSULTATION_STATUS_UNKNOWN,
    CONSULTATION_STATUS_REQUESTED,
    CONSULTATION_STATUS_MATCHED,
    CONSULTATION_STATUS_SCHEDULED,
    CONSULTATION_STATUS_IN_PROGRESS,
    CONSULTATION_STATUS_COMPLETED,
    CONSULTATION_STATUS_CANCELLED,
    CONSULTATION_STATUS_NO_SHOW,
  ];

  static final $core.List<ConsultationStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static ConsultationStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ConsultationStatus._(super.value, super.name);
}

class VideoSessionStatus extends $pb.ProtobufEnum {
  static const VideoSessionStatus VIDEO_SESSION_STATUS_UNKNOWN =
      VideoSessionStatus._(
          0, _omitEnumNames ? '' : 'VIDEO_SESSION_STATUS_UNKNOWN');
  static const VideoSessionStatus VIDEO_SESSION_STATUS_WAITING =
      VideoSessionStatus._(
          1, _omitEnumNames ? '' : 'VIDEO_SESSION_STATUS_WAITING');
  static const VideoSessionStatus VIDEO_SESSION_STATUS_CONNECTED =
      VideoSessionStatus._(
          2, _omitEnumNames ? '' : 'VIDEO_SESSION_STATUS_CONNECTED');
  static const VideoSessionStatus VIDEO_SESSION_STATUS_ENDED =
      VideoSessionStatus._(
          3, _omitEnumNames ? '' : 'VIDEO_SESSION_STATUS_ENDED');
  static const VideoSessionStatus VIDEO_SESSION_STATUS_FAILED =
      VideoSessionStatus._(
          4, _omitEnumNames ? '' : 'VIDEO_SESSION_STATUS_FAILED');

  static const $core.List<VideoSessionStatus> values = <VideoSessionStatus>[
    VIDEO_SESSION_STATUS_UNKNOWN,
    VIDEO_SESSION_STATUS_WAITING,
    VIDEO_SESSION_STATUS_CONNECTED,
    VIDEO_SESSION_STATUS_ENDED,
    VIDEO_SESSION_STATUS_FAILED,
  ];

  static final $core.List<VideoSessionStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static VideoSessionStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const VideoSessionStatus._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
