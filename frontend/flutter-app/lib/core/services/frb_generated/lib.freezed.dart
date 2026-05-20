// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'lib.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

T _$identity<T>(T value) => value;

final _privateConstructorUsedError = UnsupportedError(
    'It seems like you constructed your class using `MyClass._()`. This constructor is only meant to be used by freezed and you are not supposed to need it nor use it.\nPlease check the documentation here for more information: https://github.com/rrousselGit/freezed#adding-getters-and-methods-to-our-models');

/// @nodoc
mixin _$AiAnalysisResultDto {
  String get riskLevel => throw _privateConstructorUsedError;
  double get healthScore => throw _privateConstructorUsedError;
  String get summary => throw _privateConstructorUsedError;
  List<String> get recommendations => throw _privateConstructorUsedError;
  String get trend => throw _privateConstructorUsedError;

  /// Create a copy of AiAnalysisResultDto
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $AiAnalysisResultDtoCopyWith<AiAnalysisResultDto> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $AiAnalysisResultDtoCopyWith<$Res> {
  factory $AiAnalysisResultDtoCopyWith(
          AiAnalysisResultDto value, $Res Function(AiAnalysisResultDto) then) =
      _$AiAnalysisResultDtoCopyWithImpl<$Res, AiAnalysisResultDto>;
  @useResult
  $Res call(
      {String riskLevel,
      double healthScore,
      String summary,
      List<String> recommendations,
      String trend});
}

/// @nodoc
class _$AiAnalysisResultDtoCopyWithImpl<$Res, $Val extends AiAnalysisResultDto>
    implements $AiAnalysisResultDtoCopyWith<$Res> {
  _$AiAnalysisResultDtoCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of AiAnalysisResultDto
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? riskLevel = null,
    Object? healthScore = null,
    Object? summary = null,
    Object? recommendations = null,
    Object? trend = null,
  }) {
    return _then(_value.copyWith(
      riskLevel: null == riskLevel
          ? _value.riskLevel
          : riskLevel // ignore: cast_nullable_to_non_nullable
              as String,
      healthScore: null == healthScore
          ? _value.healthScore
          : healthScore // ignore: cast_nullable_to_non_nullable
              as double,
      summary: null == summary
          ? _value.summary
          : summary // ignore: cast_nullable_to_non_nullable
              as String,
      recommendations: null == recommendations
          ? _value.recommendations
          : recommendations // ignore: cast_nullable_to_non_nullable
              as List<String>,
      trend: null == trend
          ? _value.trend
          : trend // ignore: cast_nullable_to_non_nullable
              as String,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$AiAnalysisResultDtoImplCopyWith<$Res>
    implements $AiAnalysisResultDtoCopyWith<$Res> {
  factory _$$AiAnalysisResultDtoImplCopyWith(_$AiAnalysisResultDtoImpl value,
          $Res Function(_$AiAnalysisResultDtoImpl) then) =
      __$$AiAnalysisResultDtoImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {String riskLevel,
      double healthScore,
      String summary,
      List<String> recommendations,
      String trend});
}

/// @nodoc
class __$$AiAnalysisResultDtoImplCopyWithImpl<$Res>
    extends _$AiAnalysisResultDtoCopyWithImpl<$Res, _$AiAnalysisResultDtoImpl>
    implements _$$AiAnalysisResultDtoImplCopyWith<$Res> {
  __$$AiAnalysisResultDtoImplCopyWithImpl(_$AiAnalysisResultDtoImpl _value,
      $Res Function(_$AiAnalysisResultDtoImpl) _then)
      : super(_value, _then);

  /// Create a copy of AiAnalysisResultDto
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? riskLevel = null,
    Object? healthScore = null,
    Object? summary = null,
    Object? recommendations = null,
    Object? trend = null,
  }) {
    return _then(_$AiAnalysisResultDtoImpl(
      riskLevel: null == riskLevel
          ? _value.riskLevel
          : riskLevel // ignore: cast_nullable_to_non_nullable
              as String,
      healthScore: null == healthScore
          ? _value.healthScore
          : healthScore // ignore: cast_nullable_to_non_nullable
              as double,
      summary: null == summary
          ? _value.summary
          : summary // ignore: cast_nullable_to_non_nullable
              as String,
      recommendations: null == recommendations
          ? _value._recommendations
          : recommendations // ignore: cast_nullable_to_non_nullable
              as List<String>,
      trend: null == trend
          ? _value.trend
          : trend // ignore: cast_nullable_to_non_nullable
              as String,
    ));
  }
}

/// @nodoc

class _$AiAnalysisResultDtoImpl implements _AiAnalysisResultDto {
  const _$AiAnalysisResultDtoImpl(
      {required this.riskLevel,
      required this.healthScore,
      required this.summary,
      required final List<String> recommendations,
      required this.trend})
      : _recommendations = recommendations;

  @override
  final String riskLevel;
  @override
  final double healthScore;
  @override
  final String summary;
  final List<String> _recommendations;
  @override
  List<String> get recommendations {
    if (_recommendations is EqualUnmodifiableListView) return _recommendations;
    // ignore: implicit_dynamic_type
    return EqualUnmodifiableListView(_recommendations);
  }

  @override
  final String trend;

  @override
  String toString() {
    return 'AiAnalysisResultDto(riskLevel: $riskLevel, healthScore: $healthScore, summary: $summary, recommendations: $recommendations, trend: $trend)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$AiAnalysisResultDtoImpl &&
            (identical(other.riskLevel, riskLevel) ||
                other.riskLevel == riskLevel) &&
            (identical(other.healthScore, healthScore) ||
                other.healthScore == healthScore) &&
            (identical(other.summary, summary) || other.summary == summary) &&
            const DeepCollectionEquality()
                .equals(other._recommendations, _recommendations) &&
            (identical(other.trend, trend) || other.trend == trend));
  }

  @override
  int get hashCode => Object.hash(runtimeType, riskLevel, healthScore, summary,
      const DeepCollectionEquality().hash(_recommendations), trend);

  /// Create a copy of AiAnalysisResultDto
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$AiAnalysisResultDtoImplCopyWith<_$AiAnalysisResultDtoImpl> get copyWith =>
      __$$AiAnalysisResultDtoImplCopyWithImpl<_$AiAnalysisResultDtoImpl>(
          this, _$identity);
}

abstract class _AiAnalysisResultDto implements AiAnalysisResultDto {
  const factory _AiAnalysisResultDto(
      {required final String riskLevel,
      required final double healthScore,
      required final String summary,
      required final List<String> recommendations,
      required final String trend}) = _$AiAnalysisResultDtoImpl;

  @override
  String get riskLevel;
  @override
  double get healthScore;
  @override
  String get summary;
  @override
  List<String> get recommendations;
  @override
  String get trend;

  /// Create a copy of AiAnalysisResultDto
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$AiAnalysisResultDtoImplCopyWith<_$AiAnalysisResultDtoImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
mixin _$CartridgeInfoDto {
  String get cartridgeId => throw _privateConstructorUsedError;
  String get cartridgeType => throw _privateConstructorUsedError;
  String get lotId => throw _privateConstructorUsedError;
  String get expiryDate => throw _privateConstructorUsedError;
  int get remainingUses => throw _privateConstructorUsedError;

  /// Create a copy of CartridgeInfoDto
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $CartridgeInfoDtoCopyWith<CartridgeInfoDto> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $CartridgeInfoDtoCopyWith<$Res> {
  factory $CartridgeInfoDtoCopyWith(
          CartridgeInfoDto value, $Res Function(CartridgeInfoDto) then) =
      _$CartridgeInfoDtoCopyWithImpl<$Res, CartridgeInfoDto>;
  @useResult
  $Res call(
      {String cartridgeId,
      String cartridgeType,
      String lotId,
      String expiryDate,
      int remainingUses});
}

/// @nodoc
class _$CartridgeInfoDtoCopyWithImpl<$Res, $Val extends CartridgeInfoDto>
    implements $CartridgeInfoDtoCopyWith<$Res> {
  _$CartridgeInfoDtoCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of CartridgeInfoDto
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? cartridgeId = null,
    Object? cartridgeType = null,
    Object? lotId = null,
    Object? expiryDate = null,
    Object? remainingUses = null,
  }) {
    return _then(_value.copyWith(
      cartridgeId: null == cartridgeId
          ? _value.cartridgeId
          : cartridgeId // ignore: cast_nullable_to_non_nullable
              as String,
      cartridgeType: null == cartridgeType
          ? _value.cartridgeType
          : cartridgeType // ignore: cast_nullable_to_non_nullable
              as String,
      lotId: null == lotId
          ? _value.lotId
          : lotId // ignore: cast_nullable_to_non_nullable
              as String,
      expiryDate: null == expiryDate
          ? _value.expiryDate
          : expiryDate // ignore: cast_nullable_to_non_nullable
              as String,
      remainingUses: null == remainingUses
          ? _value.remainingUses
          : remainingUses // ignore: cast_nullable_to_non_nullable
              as int,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$CartridgeInfoDtoImplCopyWith<$Res>
    implements $CartridgeInfoDtoCopyWith<$Res> {
  factory _$$CartridgeInfoDtoImplCopyWith(_$CartridgeInfoDtoImpl value,
          $Res Function(_$CartridgeInfoDtoImpl) then) =
      __$$CartridgeInfoDtoImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {String cartridgeId,
      String cartridgeType,
      String lotId,
      String expiryDate,
      int remainingUses});
}

/// @nodoc
class __$$CartridgeInfoDtoImplCopyWithImpl<$Res>
    extends _$CartridgeInfoDtoCopyWithImpl<$Res, _$CartridgeInfoDtoImpl>
    implements _$$CartridgeInfoDtoImplCopyWith<$Res> {
  __$$CartridgeInfoDtoImplCopyWithImpl(_$CartridgeInfoDtoImpl _value,
      $Res Function(_$CartridgeInfoDtoImpl) _then)
      : super(_value, _then);

  /// Create a copy of CartridgeInfoDto
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? cartridgeId = null,
    Object? cartridgeType = null,
    Object? lotId = null,
    Object? expiryDate = null,
    Object? remainingUses = null,
  }) {
    return _then(_$CartridgeInfoDtoImpl(
      cartridgeId: null == cartridgeId
          ? _value.cartridgeId
          : cartridgeId // ignore: cast_nullable_to_non_nullable
              as String,
      cartridgeType: null == cartridgeType
          ? _value.cartridgeType
          : cartridgeType // ignore: cast_nullable_to_non_nullable
              as String,
      lotId: null == lotId
          ? _value.lotId
          : lotId // ignore: cast_nullable_to_non_nullable
              as String,
      expiryDate: null == expiryDate
          ? _value.expiryDate
          : expiryDate // ignore: cast_nullable_to_non_nullable
              as String,
      remainingUses: null == remainingUses
          ? _value.remainingUses
          : remainingUses // ignore: cast_nullable_to_non_nullable
              as int,
    ));
  }
}

/// @nodoc

class _$CartridgeInfoDtoImpl implements _CartridgeInfoDto {
  const _$CartridgeInfoDtoImpl(
      {required this.cartridgeId,
      required this.cartridgeType,
      required this.lotId,
      required this.expiryDate,
      required this.remainingUses});

  @override
  final String cartridgeId;
  @override
  final String cartridgeType;
  @override
  final String lotId;
  @override
  final String expiryDate;
  @override
  final int remainingUses;

  @override
  String toString() {
    return 'CartridgeInfoDto(cartridgeId: $cartridgeId, cartridgeType: $cartridgeType, lotId: $lotId, expiryDate: $expiryDate, remainingUses: $remainingUses)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$CartridgeInfoDtoImpl &&
            (identical(other.cartridgeId, cartridgeId) ||
                other.cartridgeId == cartridgeId) &&
            (identical(other.cartridgeType, cartridgeType) ||
                other.cartridgeType == cartridgeType) &&
            (identical(other.lotId, lotId) || other.lotId == lotId) &&
            (identical(other.expiryDate, expiryDate) ||
                other.expiryDate == expiryDate) &&
            (identical(other.remainingUses, remainingUses) ||
                other.remainingUses == remainingUses));
  }

  @override
  int get hashCode => Object.hash(runtimeType, cartridgeId, cartridgeType,
      lotId, expiryDate, remainingUses);

  /// Create a copy of CartridgeInfoDto
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$CartridgeInfoDtoImplCopyWith<_$CartridgeInfoDtoImpl> get copyWith =>
      __$$CartridgeInfoDtoImplCopyWithImpl<_$CartridgeInfoDtoImpl>(
          this, _$identity);
}

abstract class _CartridgeInfoDto implements CartridgeInfoDto {
  const factory _CartridgeInfoDto(
      {required final String cartridgeId,
      required final String cartridgeType,
      required final String lotId,
      required final String expiryDate,
      required final int remainingUses}) = _$CartridgeInfoDtoImpl;

  @override
  String get cartridgeId;
  @override
  String get cartridgeType;
  @override
  String get lotId;
  @override
  String get expiryDate;
  @override
  int get remainingUses;

  /// Create a copy of CartridgeInfoDto
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$CartridgeInfoDtoImplCopyWith<_$CartridgeInfoDtoImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
mixin _$DeviceInfoDto {
  String get deviceId => throw _privateConstructorUsedError;
  String get name => throw _privateConstructorUsedError;
  int get rssi => throw _privateConstructorUsedError;
  String get state => throw _privateConstructorUsedError;

  /// Create a copy of DeviceInfoDto
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $DeviceInfoDtoCopyWith<DeviceInfoDto> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $DeviceInfoDtoCopyWith<$Res> {
  factory $DeviceInfoDtoCopyWith(
          DeviceInfoDto value, $Res Function(DeviceInfoDto) then) =
      _$DeviceInfoDtoCopyWithImpl<$Res, DeviceInfoDto>;
  @useResult
  $Res call({String deviceId, String name, int rssi, String state});
}

/// @nodoc
class _$DeviceInfoDtoCopyWithImpl<$Res, $Val extends DeviceInfoDto>
    implements $DeviceInfoDtoCopyWith<$Res> {
  _$DeviceInfoDtoCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of DeviceInfoDto
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? deviceId = null,
    Object? name = null,
    Object? rssi = null,
    Object? state = null,
  }) {
    return _then(_value.copyWith(
      deviceId: null == deviceId
          ? _value.deviceId
          : deviceId // ignore: cast_nullable_to_non_nullable
              as String,
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      rssi: null == rssi
          ? _value.rssi
          : rssi // ignore: cast_nullable_to_non_nullable
              as int,
      state: null == state
          ? _value.state
          : state // ignore: cast_nullable_to_non_nullable
              as String,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$DeviceInfoDtoImplCopyWith<$Res>
    implements $DeviceInfoDtoCopyWith<$Res> {
  factory _$$DeviceInfoDtoImplCopyWith(
          _$DeviceInfoDtoImpl value, $Res Function(_$DeviceInfoDtoImpl) then) =
      __$$DeviceInfoDtoImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({String deviceId, String name, int rssi, String state});
}

/// @nodoc
class __$$DeviceInfoDtoImplCopyWithImpl<$Res>
    extends _$DeviceInfoDtoCopyWithImpl<$Res, _$DeviceInfoDtoImpl>
    implements _$$DeviceInfoDtoImplCopyWith<$Res> {
  __$$DeviceInfoDtoImplCopyWithImpl(
      _$DeviceInfoDtoImpl _value, $Res Function(_$DeviceInfoDtoImpl) _then)
      : super(_value, _then);

  /// Create a copy of DeviceInfoDto
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? deviceId = null,
    Object? name = null,
    Object? rssi = null,
    Object? state = null,
  }) {
    return _then(_$DeviceInfoDtoImpl(
      deviceId: null == deviceId
          ? _value.deviceId
          : deviceId // ignore: cast_nullable_to_non_nullable
              as String,
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      rssi: null == rssi
          ? _value.rssi
          : rssi // ignore: cast_nullable_to_non_nullable
              as int,
      state: null == state
          ? _value.state
          : state // ignore: cast_nullable_to_non_nullable
              as String,
    ));
  }
}

/// @nodoc

class _$DeviceInfoDtoImpl implements _DeviceInfoDto {
  const _$DeviceInfoDtoImpl(
      {required this.deviceId,
      required this.name,
      required this.rssi,
      required this.state});

  @override
  final String deviceId;
  @override
  final String name;
  @override
  final int rssi;
  @override
  final String state;

  @override
  String toString() {
    return 'DeviceInfoDto(deviceId: $deviceId, name: $name, rssi: $rssi, state: $state)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$DeviceInfoDtoImpl &&
            (identical(other.deviceId, deviceId) ||
                other.deviceId == deviceId) &&
            (identical(other.name, name) || other.name == name) &&
            (identical(other.rssi, rssi) || other.rssi == rssi) &&
            (identical(other.state, state) || other.state == state));
  }

  @override
  int get hashCode => Object.hash(runtimeType, deviceId, name, rssi, state);

  /// Create a copy of DeviceInfoDto
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$DeviceInfoDtoImplCopyWith<_$DeviceInfoDtoImpl> get copyWith =>
      __$$DeviceInfoDtoImplCopyWithImpl<_$DeviceInfoDtoImpl>(this, _$identity);
}

abstract class _DeviceInfoDto implements DeviceInfoDto {
  const factory _DeviceInfoDto(
      {required final String deviceId,
      required final String name,
      required final int rssi,
      required final String state}) = _$DeviceInfoDtoImpl;

  @override
  String get deviceId;
  @override
  String get name;
  @override
  int get rssi;
  @override
  String get state;

  /// Create a copy of DeviceInfoDto
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$DeviceInfoDtoImplCopyWith<_$DeviceInfoDtoImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
mixin _$DifferentialCorrectionDto {
  double get sDet => throw _privateConstructorUsedError;
  double get sRef => throw _privateConstructorUsedError;
  double get alpha => throw _privateConstructorUsedError;
  double get sCorrected => throw _privateConstructorUsedError;

  /// Create a copy of DifferentialCorrectionDto
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $DifferentialCorrectionDtoCopyWith<DifferentialCorrectionDto> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $DifferentialCorrectionDtoCopyWith<$Res> {
  factory $DifferentialCorrectionDtoCopyWith(DifferentialCorrectionDto value,
          $Res Function(DifferentialCorrectionDto) then) =
      _$DifferentialCorrectionDtoCopyWithImpl<$Res, DifferentialCorrectionDto>;
  @useResult
  $Res call({double sDet, double sRef, double alpha, double sCorrected});
}

/// @nodoc
class _$DifferentialCorrectionDtoCopyWithImpl<$Res,
        $Val extends DifferentialCorrectionDto>
    implements $DifferentialCorrectionDtoCopyWith<$Res> {
  _$DifferentialCorrectionDtoCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of DifferentialCorrectionDto
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? sDet = null,
    Object? sRef = null,
    Object? alpha = null,
    Object? sCorrected = null,
  }) {
    return _then(_value.copyWith(
      sDet: null == sDet
          ? _value.sDet
          : sDet // ignore: cast_nullable_to_non_nullable
              as double,
      sRef: null == sRef
          ? _value.sRef
          : sRef // ignore: cast_nullable_to_non_nullable
              as double,
      alpha: null == alpha
          ? _value.alpha
          : alpha // ignore: cast_nullable_to_non_nullable
              as double,
      sCorrected: null == sCorrected
          ? _value.sCorrected
          : sCorrected // ignore: cast_nullable_to_non_nullable
              as double,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$DifferentialCorrectionDtoImplCopyWith<$Res>
    implements $DifferentialCorrectionDtoCopyWith<$Res> {
  factory _$$DifferentialCorrectionDtoImplCopyWith(
          _$DifferentialCorrectionDtoImpl value,
          $Res Function(_$DifferentialCorrectionDtoImpl) then) =
      __$$DifferentialCorrectionDtoImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({double sDet, double sRef, double alpha, double sCorrected});
}

/// @nodoc
class __$$DifferentialCorrectionDtoImplCopyWithImpl<$Res>
    extends _$DifferentialCorrectionDtoCopyWithImpl<$Res,
        _$DifferentialCorrectionDtoImpl>
    implements _$$DifferentialCorrectionDtoImplCopyWith<$Res> {
  __$$DifferentialCorrectionDtoImplCopyWithImpl(
      _$DifferentialCorrectionDtoImpl _value,
      $Res Function(_$DifferentialCorrectionDtoImpl) _then)
      : super(_value, _then);

  /// Create a copy of DifferentialCorrectionDto
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? sDet = null,
    Object? sRef = null,
    Object? alpha = null,
    Object? sCorrected = null,
  }) {
    return _then(_$DifferentialCorrectionDtoImpl(
      sDet: null == sDet
          ? _value.sDet
          : sDet // ignore: cast_nullable_to_non_nullable
              as double,
      sRef: null == sRef
          ? _value.sRef
          : sRef // ignore: cast_nullable_to_non_nullable
              as double,
      alpha: null == alpha
          ? _value.alpha
          : alpha // ignore: cast_nullable_to_non_nullable
              as double,
      sCorrected: null == sCorrected
          ? _value.sCorrected
          : sCorrected // ignore: cast_nullable_to_non_nullable
              as double,
    ));
  }
}

/// @nodoc

class _$DifferentialCorrectionDtoImpl implements _DifferentialCorrectionDto {
  const _$DifferentialCorrectionDtoImpl(
      {required this.sDet,
      required this.sRef,
      required this.alpha,
      required this.sCorrected});

  @override
  final double sDet;
  @override
  final double sRef;
  @override
  final double alpha;
  @override
  final double sCorrected;

  @override
  String toString() {
    return 'DifferentialCorrectionDto(sDet: $sDet, sRef: $sRef, alpha: $alpha, sCorrected: $sCorrected)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$DifferentialCorrectionDtoImpl &&
            (identical(other.sDet, sDet) || other.sDet == sDet) &&
            (identical(other.sRef, sRef) || other.sRef == sRef) &&
            (identical(other.alpha, alpha) || other.alpha == alpha) &&
            (identical(other.sCorrected, sCorrected) ||
                other.sCorrected == sCorrected));
  }

  @override
  int get hashCode => Object.hash(runtimeType, sDet, sRef, alpha, sCorrected);

  /// Create a copy of DifferentialCorrectionDto
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$DifferentialCorrectionDtoImplCopyWith<_$DifferentialCorrectionDtoImpl>
      get copyWith => __$$DifferentialCorrectionDtoImplCopyWithImpl<
          _$DifferentialCorrectionDtoImpl>(this, _$identity);
}

abstract class _DifferentialCorrectionDto implements DifferentialCorrectionDto {
  const factory _DifferentialCorrectionDto(
      {required final double sDet,
      required final double sRef,
      required final double alpha,
      required final double sCorrected}) = _$DifferentialCorrectionDtoImpl;

  @override
  double get sDet;
  @override
  double get sRef;
  @override
  double get alpha;
  @override
  double get sCorrected;

  /// Create a copy of DifferentialCorrectionDto
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$DifferentialCorrectionDtoImplCopyWith<_$DifferentialCorrectionDtoImpl>
      get copyWith => throw _privateConstructorUsedError;
}

/// @nodoc
mixin _$FingerprintDto {
  Float32List get data => throw _privateConstructorUsedError;
  BigInt get dimension => throw _privateConstructorUsedError;
  String get measurementType => throw _privateConstructorUsedError;
  bool get normalized => throw _privateConstructorUsedError;

  /// Create a copy of FingerprintDto
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $FingerprintDtoCopyWith<FingerprintDto> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $FingerprintDtoCopyWith<$Res> {
  factory $FingerprintDtoCopyWith(
          FingerprintDto value, $Res Function(FingerprintDto) then) =
      _$FingerprintDtoCopyWithImpl<$Res, FingerprintDto>;
  @useResult
  $Res call(
      {Float32List data,
      BigInt dimension,
      String measurementType,
      bool normalized});
}

/// @nodoc
class _$FingerprintDtoCopyWithImpl<$Res, $Val extends FingerprintDto>
    implements $FingerprintDtoCopyWith<$Res> {
  _$FingerprintDtoCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of FingerprintDto
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? data = null,
    Object? dimension = null,
    Object? measurementType = null,
    Object? normalized = null,
  }) {
    return _then(_value.copyWith(
      data: null == data
          ? _value.data
          : data // ignore: cast_nullable_to_non_nullable
              as Float32List,
      dimension: null == dimension
          ? _value.dimension
          : dimension // ignore: cast_nullable_to_non_nullable
              as BigInt,
      measurementType: null == measurementType
          ? _value.measurementType
          : measurementType // ignore: cast_nullable_to_non_nullable
              as String,
      normalized: null == normalized
          ? _value.normalized
          : normalized // ignore: cast_nullable_to_non_nullable
              as bool,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$FingerprintDtoImplCopyWith<$Res>
    implements $FingerprintDtoCopyWith<$Res> {
  factory _$$FingerprintDtoImplCopyWith(_$FingerprintDtoImpl value,
          $Res Function(_$FingerprintDtoImpl) then) =
      __$$FingerprintDtoImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {Float32List data,
      BigInt dimension,
      String measurementType,
      bool normalized});
}

/// @nodoc
class __$$FingerprintDtoImplCopyWithImpl<$Res>
    extends _$FingerprintDtoCopyWithImpl<$Res, _$FingerprintDtoImpl>
    implements _$$FingerprintDtoImplCopyWith<$Res> {
  __$$FingerprintDtoImplCopyWithImpl(
      _$FingerprintDtoImpl _value, $Res Function(_$FingerprintDtoImpl) _then)
      : super(_value, _then);

  /// Create a copy of FingerprintDto
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? data = null,
    Object? dimension = null,
    Object? measurementType = null,
    Object? normalized = null,
  }) {
    return _then(_$FingerprintDtoImpl(
      data: null == data
          ? _value.data
          : data // ignore: cast_nullable_to_non_nullable
              as Float32List,
      dimension: null == dimension
          ? _value.dimension
          : dimension // ignore: cast_nullable_to_non_nullable
              as BigInt,
      measurementType: null == measurementType
          ? _value.measurementType
          : measurementType // ignore: cast_nullable_to_non_nullable
              as String,
      normalized: null == normalized
          ? _value.normalized
          : normalized // ignore: cast_nullable_to_non_nullable
              as bool,
    ));
  }
}

/// @nodoc

class _$FingerprintDtoImpl implements _FingerprintDto {
  const _$FingerprintDtoImpl(
      {required this.data,
      required this.dimension,
      required this.measurementType,
      required this.normalized});

  @override
  final Float32List data;
  @override
  final BigInt dimension;
  @override
  final String measurementType;
  @override
  final bool normalized;

  @override
  String toString() {
    return 'FingerprintDto(data: $data, dimension: $dimension, measurementType: $measurementType, normalized: $normalized)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$FingerprintDtoImpl &&
            const DeepCollectionEquality().equals(other.data, data) &&
            (identical(other.dimension, dimension) ||
                other.dimension == dimension) &&
            (identical(other.measurementType, measurementType) ||
                other.measurementType == measurementType) &&
            (identical(other.normalized, normalized) ||
                other.normalized == normalized));
  }

  @override
  int get hashCode => Object.hash(
      runtimeType,
      const DeepCollectionEquality().hash(data),
      dimension,
      measurementType,
      normalized);

  /// Create a copy of FingerprintDto
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$FingerprintDtoImplCopyWith<_$FingerprintDtoImpl> get copyWith =>
      __$$FingerprintDtoImplCopyWithImpl<_$FingerprintDtoImpl>(
          this, _$identity);
}

abstract class _FingerprintDto implements FingerprintDto {
  const factory _FingerprintDto(
      {required final Float32List data,
      required final BigInt dimension,
      required final String measurementType,
      required final bool normalized}) = _$FingerprintDtoImpl;

  @override
  Float32List get data;
  @override
  BigInt get dimension;
  @override
  String get measurementType;
  @override
  bool get normalized;

  /// Create a copy of FingerprintDto
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$FingerprintDtoImplCopyWith<_$FingerprintDtoImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
mixin _$MeasurementPipelineDto {
  double get primaryValue => throw _privateConstructorUsedError;
  double get referenceValue => throw _privateConstructorUsedError;
  double get differentialValue => throw _privateConstructorUsedError;
  double get snr => throw _privateConstructorUsedError;
  double get confidence => throw _privateConstructorUsedError;
  String get biomarker => throw _privateConstructorUsedError;
  String get unit => throw _privateConstructorUsedError;
  String get riskLevel => throw _privateConstructorUsedError;
  double get healthScore => throw _privateConstructorUsedError;
  List<String> get recommendations => throw _privateConstructorUsedError;
  double get pipelineDurationMs => throw _privateConstructorUsedError;

  /// Create a copy of MeasurementPipelineDto
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $MeasurementPipelineDtoCopyWith<MeasurementPipelineDto> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $MeasurementPipelineDtoCopyWith<$Res> {
  factory $MeasurementPipelineDtoCopyWith(MeasurementPipelineDto value,
          $Res Function(MeasurementPipelineDto) then) =
      _$MeasurementPipelineDtoCopyWithImpl<$Res, MeasurementPipelineDto>;
  @useResult
  $Res call(
      {double primaryValue,
      double referenceValue,
      double differentialValue,
      double snr,
      double confidence,
      String biomarker,
      String unit,
      String riskLevel,
      double healthScore,
      List<String> recommendations,
      double pipelineDurationMs});
}

/// @nodoc
class _$MeasurementPipelineDtoCopyWithImpl<$Res,
        $Val extends MeasurementPipelineDto>
    implements $MeasurementPipelineDtoCopyWith<$Res> {
  _$MeasurementPipelineDtoCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of MeasurementPipelineDto
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? primaryValue = null,
    Object? referenceValue = null,
    Object? differentialValue = null,
    Object? snr = null,
    Object? confidence = null,
    Object? biomarker = null,
    Object? unit = null,
    Object? riskLevel = null,
    Object? healthScore = null,
    Object? recommendations = null,
    Object? pipelineDurationMs = null,
  }) {
    return _then(_value.copyWith(
      primaryValue: null == primaryValue
          ? _value.primaryValue
          : primaryValue // ignore: cast_nullable_to_non_nullable
              as double,
      referenceValue: null == referenceValue
          ? _value.referenceValue
          : referenceValue // ignore: cast_nullable_to_non_nullable
              as double,
      differentialValue: null == differentialValue
          ? _value.differentialValue
          : differentialValue // ignore: cast_nullable_to_non_nullable
              as double,
      snr: null == snr
          ? _value.snr
          : snr // ignore: cast_nullable_to_non_nullable
              as double,
      confidence: null == confidence
          ? _value.confidence
          : confidence // ignore: cast_nullable_to_non_nullable
              as double,
      biomarker: null == biomarker
          ? _value.biomarker
          : biomarker // ignore: cast_nullable_to_non_nullable
              as String,
      unit: null == unit
          ? _value.unit
          : unit // ignore: cast_nullable_to_non_nullable
              as String,
      riskLevel: null == riskLevel
          ? _value.riskLevel
          : riskLevel // ignore: cast_nullable_to_non_nullable
              as String,
      healthScore: null == healthScore
          ? _value.healthScore
          : healthScore // ignore: cast_nullable_to_non_nullable
              as double,
      recommendations: null == recommendations
          ? _value.recommendations
          : recommendations // ignore: cast_nullable_to_non_nullable
              as List<String>,
      pipelineDurationMs: null == pipelineDurationMs
          ? _value.pipelineDurationMs
          : pipelineDurationMs // ignore: cast_nullable_to_non_nullable
              as double,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$MeasurementPipelineDtoImplCopyWith<$Res>
    implements $MeasurementPipelineDtoCopyWith<$Res> {
  factory _$$MeasurementPipelineDtoImplCopyWith(
          _$MeasurementPipelineDtoImpl value,
          $Res Function(_$MeasurementPipelineDtoImpl) then) =
      __$$MeasurementPipelineDtoImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {double primaryValue,
      double referenceValue,
      double differentialValue,
      double snr,
      double confidence,
      String biomarker,
      String unit,
      String riskLevel,
      double healthScore,
      List<String> recommendations,
      double pipelineDurationMs});
}

/// @nodoc
class __$$MeasurementPipelineDtoImplCopyWithImpl<$Res>
    extends _$MeasurementPipelineDtoCopyWithImpl<$Res,
        _$MeasurementPipelineDtoImpl>
    implements _$$MeasurementPipelineDtoImplCopyWith<$Res> {
  __$$MeasurementPipelineDtoImplCopyWithImpl(
      _$MeasurementPipelineDtoImpl _value,
      $Res Function(_$MeasurementPipelineDtoImpl) _then)
      : super(_value, _then);

  /// Create a copy of MeasurementPipelineDto
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? primaryValue = null,
    Object? referenceValue = null,
    Object? differentialValue = null,
    Object? snr = null,
    Object? confidence = null,
    Object? biomarker = null,
    Object? unit = null,
    Object? riskLevel = null,
    Object? healthScore = null,
    Object? recommendations = null,
    Object? pipelineDurationMs = null,
  }) {
    return _then(_$MeasurementPipelineDtoImpl(
      primaryValue: null == primaryValue
          ? _value.primaryValue
          : primaryValue // ignore: cast_nullable_to_non_nullable
              as double,
      referenceValue: null == referenceValue
          ? _value.referenceValue
          : referenceValue // ignore: cast_nullable_to_non_nullable
              as double,
      differentialValue: null == differentialValue
          ? _value.differentialValue
          : differentialValue // ignore: cast_nullable_to_non_nullable
              as double,
      snr: null == snr
          ? _value.snr
          : snr // ignore: cast_nullable_to_non_nullable
              as double,
      confidence: null == confidence
          ? _value.confidence
          : confidence // ignore: cast_nullable_to_non_nullable
              as double,
      biomarker: null == biomarker
          ? _value.biomarker
          : biomarker // ignore: cast_nullable_to_non_nullable
              as String,
      unit: null == unit
          ? _value.unit
          : unit // ignore: cast_nullable_to_non_nullable
              as String,
      riskLevel: null == riskLevel
          ? _value.riskLevel
          : riskLevel // ignore: cast_nullable_to_non_nullable
              as String,
      healthScore: null == healthScore
          ? _value.healthScore
          : healthScore // ignore: cast_nullable_to_non_nullable
              as double,
      recommendations: null == recommendations
          ? _value._recommendations
          : recommendations // ignore: cast_nullable_to_non_nullable
              as List<String>,
      pipelineDurationMs: null == pipelineDurationMs
          ? _value.pipelineDurationMs
          : pipelineDurationMs // ignore: cast_nullable_to_non_nullable
              as double,
    ));
  }
}

/// @nodoc

class _$MeasurementPipelineDtoImpl implements _MeasurementPipelineDto {
  const _$MeasurementPipelineDtoImpl(
      {required this.primaryValue,
      required this.referenceValue,
      required this.differentialValue,
      required this.snr,
      required this.confidence,
      required this.biomarker,
      required this.unit,
      required this.riskLevel,
      required this.healthScore,
      required final List<String> recommendations,
      required this.pipelineDurationMs})
      : _recommendations = recommendations;

  @override
  final double primaryValue;
  @override
  final double referenceValue;
  @override
  final double differentialValue;
  @override
  final double snr;
  @override
  final double confidence;
  @override
  final String biomarker;
  @override
  final String unit;
  @override
  final String riskLevel;
  @override
  final double healthScore;
  final List<String> _recommendations;
  @override
  List<String> get recommendations {
    if (_recommendations is EqualUnmodifiableListView) return _recommendations;
    // ignore: implicit_dynamic_type
    return EqualUnmodifiableListView(_recommendations);
  }

  @override
  final double pipelineDurationMs;

  @override
  String toString() {
    return 'MeasurementPipelineDto(primaryValue: $primaryValue, referenceValue: $referenceValue, differentialValue: $differentialValue, snr: $snr, confidence: $confidence, biomarker: $biomarker, unit: $unit, riskLevel: $riskLevel, healthScore: $healthScore, recommendations: $recommendations, pipelineDurationMs: $pipelineDurationMs)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$MeasurementPipelineDtoImpl &&
            (identical(other.primaryValue, primaryValue) ||
                other.primaryValue == primaryValue) &&
            (identical(other.referenceValue, referenceValue) ||
                other.referenceValue == referenceValue) &&
            (identical(other.differentialValue, differentialValue) ||
                other.differentialValue == differentialValue) &&
            (identical(other.snr, snr) || other.snr == snr) &&
            (identical(other.confidence, confidence) ||
                other.confidence == confidence) &&
            (identical(other.biomarker, biomarker) ||
                other.biomarker == biomarker) &&
            (identical(other.unit, unit) || other.unit == unit) &&
            (identical(other.riskLevel, riskLevel) ||
                other.riskLevel == riskLevel) &&
            (identical(other.healthScore, healthScore) ||
                other.healthScore == healthScore) &&
            const DeepCollectionEquality()
                .equals(other._recommendations, _recommendations) &&
            (identical(other.pipelineDurationMs, pipelineDurationMs) ||
                other.pipelineDurationMs == pipelineDurationMs));
  }

  @override
  int get hashCode => Object.hash(
      runtimeType,
      primaryValue,
      referenceValue,
      differentialValue,
      snr,
      confidence,
      biomarker,
      unit,
      riskLevel,
      healthScore,
      const DeepCollectionEquality().hash(_recommendations),
      pipelineDurationMs);

  /// Create a copy of MeasurementPipelineDto
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$MeasurementPipelineDtoImplCopyWith<_$MeasurementPipelineDtoImpl>
      get copyWith => __$$MeasurementPipelineDtoImplCopyWithImpl<
          _$MeasurementPipelineDtoImpl>(this, _$identity);
}

abstract class _MeasurementPipelineDto implements MeasurementPipelineDto {
  const factory _MeasurementPipelineDto(
      {required final double primaryValue,
      required final double referenceValue,
      required final double differentialValue,
      required final double snr,
      required final double confidence,
      required final String biomarker,
      required final String unit,
      required final String riskLevel,
      required final double healthScore,
      required final List<String> recommendations,
      required final double pipelineDurationMs}) = _$MeasurementPipelineDtoImpl;

  @override
  double get primaryValue;
  @override
  double get referenceValue;
  @override
  double get differentialValue;
  @override
  double get snr;
  @override
  double get confidence;
  @override
  String get biomarker;
  @override
  String get unit;
  @override
  String get riskLevel;
  @override
  double get healthScore;
  @override
  List<String> get recommendations;
  @override
  double get pipelineDurationMs;

  /// Create a copy of MeasurementPipelineDto
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$MeasurementPipelineDtoImplCopyWith<_$MeasurementPipelineDtoImpl>
      get copyWith => throw _privateConstructorUsedError;
}
