// ManPaSik gRPC - 수동 생성 (protoc 미사용). proto: backend/shared/proto/manpasik.proto
// ignore_for_file: type=lint

import 'dart:core' as $core;
import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

// ---------------------------------------------------------------------------
// Auth
// ---------------------------------------------------------------------------

class RegisterRequest extends $pb.GeneratedMessage {
  factory RegisterRequest({
    $core.String? email,
    $core.String? password,
    $core.String? displayName,
  }) =>
      CreateMessage();
  RegisterRequest._() : super();
  factory RegisterRequest.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);
  factory RegisterRequest.fromJson($core.String i) =>
      CreateMessage()..mergeFromJson(i);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'RegisterRequest',
    createEmptyInstance: CreateMessage,
  )
    ..aOS(1, 'email')
    ..aOS(2, 'password')
    ..aOS(3, 'displayName');

  @$pb.TagNumber(1)
  $core.String get email => $_getSZ(0);
  @$pb.TagNumber(1)
  set email($core.String v) {
    $_setString(0, v);
  }

  @$pb.TagNumber(2)
  $core.String get password => $_getSZ(1);
  @$pb.TagNumber(2)
  set password($core.String v) {
    $_setString(1, v);
  }

  @$pb.TagNumber(3)
  $core.String get displayName => $_getSZ(2);
  @$pb.TagNumber(3)
  set displayName($core.String v) {
    $_setString(2, v);
  }

  static RegisterRequest CreateMessage() => RegisterRequest._();

  @$core.override
  RegisterRequest clone() => RegisterRequest()..mergeFromMessage(this);

  @$core.override
  RegisterRequest createEmptyInstance() => CreateMessage();

  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class RegisterResponse extends $pb.GeneratedMessage {
  factory RegisterResponse({
    $core.String? userId,
    $core.String? email,
    $core.String? displayName,
    $core.String? role,
  }) =>
      CreateMessage();
  RegisterResponse._() : super();
  factory RegisterResponse.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);
  factory RegisterResponse.fromJson($core.String i) =>
      CreateMessage()..mergeFromJson(i);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'RegisterResponse',
    createEmptyInstance: CreateMessage,
  )
    ..aOS(1, 'userId')
    ..aOS(2, 'email')
    ..aOS(3, 'displayName')
    ..aOS(4, 'role');

  @$pb.TagNumber(1)
  $core.String get userId => $_getSZ(0);
  @$pb.TagNumber(1)
  set userId($core.String v) {
    $_setString(0, v);
  }

  @$pb.TagNumber(2)
  $core.String get email => $_getSZ(1);
  @$pb.TagNumber(2)
  set email($core.String v) {
    $_setString(1, v);
  }

  @$pb.TagNumber(3)
  $core.String get displayName => $_getSZ(2);
  @$pb.TagNumber(3)
  set displayName($core.String v) {
    $_setString(2, v);
  }

  @$pb.TagNumber(4)
  $core.String get role => $_getSZ(3);
  @$pb.TagNumber(4)
  set role($core.String v) {
    $_setString(3, v);
  }

  static RegisterResponse CreateMessage() => RegisterResponse._();

  @$core.override
  RegisterResponse clone() => RegisterResponse()..mergeFromMessage(this);

  @$core.override
  RegisterResponse createEmptyInstance() => CreateMessage();

  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class LoginRequest extends $pb.GeneratedMessage {
  factory LoginRequest({
    $core.String? email,
    $core.String? password,
  }) =>
      CreateMessage();
  LoginRequest._() : super();
  factory LoginRequest.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);
  factory LoginRequest.fromJson($core.String i) => CreateMessage()..mergeFromJson(i);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'LoginRequest',
    createEmptyInstance: CreateMessage,
  )
    ..aOS(1, 'email')
    ..aOS(2, 'password');

  @$pb.TagNumber(1)
  $core.String get email => $_getSZ(0);
  @$pb.TagNumber(1)
  set email($core.String v) {
    $_setString(0, v);
  }

  @$pb.TagNumber(2)
  $core.String get password => $_getSZ(1);
  @$pb.TagNumber(2)
  set password($core.String v) {
    $_setString(1, v);
  }

  static LoginRequest CreateMessage() => LoginRequest._();

  @$core.override
  LoginRequest clone() => LoginRequest()..mergeFromMessage(this);

  @$core.override
  LoginRequest createEmptyInstance() => CreateMessage();

  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class LoginResponse extends $pb.GeneratedMessage {
  factory LoginResponse({
    $core.String? accessToken,
    $core.String? refreshToken,
    $fixnum.Int64? expiresIn,
    $core.String? tokenType,
    $core.String? userId,
  }) =>
      CreateMessage();
  LoginResponse._() : super();
  factory LoginResponse.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);
  factory LoginResponse.fromJson($core.String i) =>
      CreateMessage()..mergeFromJson(i);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'LoginResponse',
    createEmptyInstance: CreateMessage,
  )
    ..aOS(1, 'accessToken')
    ..aOS(2, 'refreshToken')
    ..aInt64(3, 'expiresIn')
    ..aOS(4, 'tokenType')
    ..aOS(5, 'userId');

  @$pb.TagNumber(1)
  $core.String get accessToken => $_getSZ(0);
  @$pb.TagNumber(1)
  set accessToken($core.String v) {
    $_setString(0, v);
  }

  @$pb.TagNumber(2)
  $core.String get refreshToken => $_getSZ(1);
  @$pb.TagNumber(2)
  set refreshToken($core.String v) {
    $_setString(1, v);
  }

  @$pb.TagNumber(3)
  $fixnum.Int64 get expiresIn => $_getI64(2);
  @$pb.TagNumber(3)
  set expiresIn($fixnum.Int64 v) {
    $_setInt64(2, v);
  }

  @$pb.TagNumber(4)
  $core.String get tokenType => $_getSZ(3);
  @$pb.TagNumber(4)
  set tokenType($core.String v) {
    $_setString(3, v);
  }

  @$pb.TagNumber(5)
  $core.String get userId => $_getSZ(4);
  @$pb.TagNumber(5)
  set userId($core.String v) {
    $_setString(4, v);
  }

  static LoginResponse CreateMessage() => LoginResponse._();

  @$core.override
  LoginResponse clone() => LoginResponse()..mergeFromMessage(this);

  @$core.override
  LoginResponse createEmptyInstance() => CreateMessage();

  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class RefreshTokenRequest extends $pb.GeneratedMessage {
  factory RefreshTokenRequest({$core.String? refreshToken}) => CreateMessage();
  RefreshTokenRequest._() : super();
  factory RefreshTokenRequest.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);
  factory RefreshTokenRequest.fromJson($core.String i) =>
      CreateMessage()..mergeFromJson(i);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'RefreshTokenRequest',
    createEmptyInstance: CreateMessage,
  )..aOS(1, 'refreshToken');

  @$pb.TagNumber(1)
  $core.String get refreshToken => $_getSZ(0);
  @$pb.TagNumber(1)
  set refreshToken($core.String v) {
    $_setString(0, v);
  }

  static RefreshTokenRequest CreateMessage() => RefreshTokenRequest._();

  @$core.override
  RefreshTokenRequest clone() => RefreshTokenRequest()..mergeFromMessage(this);

  @$core.override
  RefreshTokenRequest createEmptyInstance() => CreateMessage();

  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class LogoutRequest extends $pb.GeneratedMessage {
  factory LogoutRequest({$core.String? userId}) => CreateMessage();
  LogoutRequest._() : super();
  factory LogoutRequest.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);
  factory LogoutRequest.fromJson($core.String i) =>
      CreateMessage()..mergeFromJson(i);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'LogoutRequest',
    createEmptyInstance: CreateMessage,
  )..aOS(1, 'userId');

  @$pb.TagNumber(1)
  $core.String get userId => $_getSZ(0);
  @$pb.TagNumber(1)
  set userId($core.String v) {
    $_setString(0, v);
  }

  static LogoutRequest CreateMessage() => LogoutRequest._();

  @$core.override
  LogoutRequest clone() => LogoutRequest()..mergeFromMessage(this);

  @$core.override
  LogoutRequest createEmptyInstance() => CreateMessage();

  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class LogoutResponse extends $pb.GeneratedMessage {
  factory LogoutResponse() => CreateMessage();
  LogoutResponse._() : super();
  factory LogoutResponse.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);
  factory LogoutResponse.fromJson($core.String i) =>
      CreateMessage()..mergeFromJson(i);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'LogoutResponse',
    createEmptyInstance: CreateMessage,
  );

  static LogoutResponse CreateMessage() => LogoutResponse._();

  @$core.override
  LogoutResponse clone() => LogoutResponse()..mergeFromMessage(this);

  @$core.override
  LogoutResponse createEmptyInstance() => CreateMessage();

  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

// ---------------------------------------------------------------------------
// User (minimal for GetProfile, GetSubscription)
// ---------------------------------------------------------------------------

class GetProfileRequest extends $pb.GeneratedMessage {
  factory GetProfileRequest({$core.String? userId}) => CreateMessage();
  GetProfileRequest._() : super();
  factory GetProfileRequest.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);
  factory GetProfileRequest.fromJson($core.String i) =>
      CreateMessage()..mergeFromJson(i);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'GetProfileRequest',
    createEmptyInstance: CreateMessage,
  )..aOS(1, 'userId');

  @$pb.TagNumber(1)
  $core.String get userId => $_getSZ(0);
  @$pb.TagNumber(1)
  set userId($core.String v) {
    $_setString(0, v);
  }

  static GetProfileRequest CreateMessage() => GetProfileRequest._();

  @$core.override
  GetProfileRequest clone() => GetProfileRequest()..mergeFromMessage(this);

  @$core.override
  GetProfileRequest createEmptyInstance() => CreateMessage();

  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class UserProfile extends $pb.GeneratedMessage {
  factory UserProfile({
    $core.String? userId,
    $core.String? email,
    $core.String? displayName,
    $core.String? avatarUrl,
    $core.String? language,
    $core.String? timezone,
    $core.int? subscriptionTier,
  }) =>
      CreateMessage();
  UserProfile._() : super();
  factory UserProfile.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);
  factory UserProfile.fromJson($core.String i) =>
      CreateMessage()..mergeFromJson(i);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'UserProfile',
    createEmptyInstance: CreateMessage,
  )
    ..aOS(1, 'userId')
    ..aOS(2, 'email')
    ..aOS(3, 'displayName')
    ..aOS(4, 'avatarUrl')
    ..aOS(5, 'language')
    ..aOS(6, 'timezone')
    ..a<$core.int>(7, 'subscriptionTier', $pb.PbFieldType.O3);

  @$pb.TagNumber(1)
  $core.String get userId => $_getSZ(0);
  @$pb.TagNumber(1)
  set userId($core.String v) {
    $_setString(0, v);
  }

  @$pb.TagNumber(2)
  $core.String get email => $_getSZ(1);
  @$pb.TagNumber(2)
  set email($core.String v) {
    $_setString(1, v);
  }

  @$pb.TagNumber(3)
  $core.String get displayName => $_getSZ(2);
  @$pb.TagNumber(3)
  set displayName($core.String v) {
    $_setString(2, v);
  }

  @$pb.TagNumber(4)
  $core.String get avatarUrl => $_getSZ(3);
  @$pb.TagNumber(4)
  set avatarUrl($core.String v) {
    $_setString(3, v);
  }

  @$pb.TagNumber(5)
  $core.String get language => $_getSZ(4);
  @$pb.TagNumber(5)
  set language($core.String v) {
    $_setString(4, v);
  }

  @$pb.TagNumber(6)
  $core.String get timezone => $_getSZ(5);
  @$pb.TagNumber(6)
  set timezone($core.String v) {
    $_setString(5, v);
  }

  @$pb.TagNumber(7)
  $core.int get subscriptionTier => $_getIZ(6);
  @$pb.TagNumber(7)
  set subscriptionTier($core.int v) {
    $_setSignedInt32(6, v);
  }

  static UserProfile CreateMessage() => UserProfile._();

  @$core.override
  UserProfile clone() => UserProfile()..mergeFromMessage(this);

  @$core.override
  UserProfile createEmptyInstance() => CreateMessage();

  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class GetSubscriptionRequest extends $pb.GeneratedMessage {
  factory GetSubscriptionRequest({$core.String? userId}) => CreateMessage();
  GetSubscriptionRequest._() : super();
  factory GetSubscriptionRequest.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);
  factory GetSubscriptionRequest.fromJson($core.String i) =>
      CreateMessage()..mergeFromJson(i);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'GetSubscriptionRequest',
    createEmptyInstance: CreateMessage,
  )..aOS(1, 'userId');

  @$pb.TagNumber(1)
  $core.String get userId => $_getSZ(0);
  @$pb.TagNumber(1)
  set userId($core.String v) {
    $_setString(0, v);
  }

  static GetSubscriptionRequest CreateMessage() => GetSubscriptionRequest._();

  @$core.override
  GetSubscriptionRequest clone() => GetSubscriptionRequest()..mergeFromMessage(this);

  @$core.override
  GetSubscriptionRequest createEmptyInstance() => CreateMessage();

  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class SubscriptionInfo extends $pb.GeneratedMessage {
  factory SubscriptionInfo({
    $core.String? userId,
    $core.int? tier,
    $core.int? maxDevices,
    $core.int? maxFamilyMembers,
    $core.bool? aiCoachingEnabled,
    $core.bool? telemedicineEnabled,
  }) =>
      CreateMessage();
  SubscriptionInfo._() : super();
  factory SubscriptionInfo.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);
  factory SubscriptionInfo.fromJson($core.String i) =>
      CreateMessage()..mergeFromJson(i);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'SubscriptionInfo',
    createEmptyInstance: CreateMessage,
  )
    ..aOS(1, 'userId')
    ..a<$core.int>(2, 'tier', $pb.PbFieldType.O3)
    ..a<$core.int>(5, 'maxDevices', $pb.PbFieldType.O3)
    ..a<$core.int>(6, 'maxFamilyMembers', $pb.PbFieldType.O3)
    ..aOB(7, 'aiCoachingEnabled')
    ..aOB(8, 'telemedicineEnabled');

  @$pb.TagNumber(1)
  $core.String get userId => $_getSZ(0);
  @$pb.TagNumber(1)
  set userId($core.String v) {
    $_setString(0, v);
  }

  @$pb.TagNumber(2)
  $core.int get tier => $_getIZ(1);
  @$pb.TagNumber(2)
  set tier($core.int v) {
    $_setSignedInt32(1, v);
  }

  @$pb.TagNumber(5)
  $core.int get maxDevices => $_getIZ(2);
  @$pb.TagNumber(5)
  set maxDevices($core.int v) {
    $_setSignedInt32(2, v);
  }

  @$pb.TagNumber(6)
  $core.int get maxFamilyMembers => $_getIZ(3);
  @$pb.TagNumber(6)
  set maxFamilyMembers($core.int v) {
    $_setSignedInt32(3, v);
  }

  @$pb.TagNumber(7)
  $core.bool get aiCoachingEnabled => $_getBF(4);
  @$pb.TagNumber(7)
  set aiCoachingEnabled($core.bool v) {
    $_setBool(4, v);
  }

  @$pb.TagNumber(8)
  $core.bool get telemedicineEnabled => $_getBF(5);
  @$pb.TagNumber(8)
  set telemedicineEnabled($core.bool v) {
    $_setBool(5, v);
  }

  static SubscriptionInfo CreateMessage() => SubscriptionInfo._();

  @$core.override
  SubscriptionInfo clone() => SubscriptionInfo()..mergeFromMessage(this);

  @$core.override
  SubscriptionInfo createEmptyInstance() => CreateMessage();

  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

// ---------------------------------------------------------------------------
// Device (ListDevices)
// ---------------------------------------------------------------------------

class ListDevicesRequest extends $pb.GeneratedMessage {
  factory ListDevicesRequest({$core.String? userId}) => CreateMessage();
  ListDevicesRequest._() : super();
  factory ListDevicesRequest.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);
  factory ListDevicesRequest.fromJson($core.String i) =>
      CreateMessage()..mergeFromJson(i);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'ListDevicesRequest',
    createEmptyInstance: CreateMessage,
  )..aOS(1, 'userId');

  @$pb.TagNumber(1)
  $core.String get userId => $_getSZ(0);
  @$pb.TagNumber(1)
  set userId($core.String v) {
    $_setString(0, v);
  }

  static ListDevicesRequest CreateMessage() => ListDevicesRequest._();

  @$core.override
  ListDevicesRequest clone() => ListDevicesRequest()..mergeFromMessage(this);

  @$core.override
  ListDevicesRequest createEmptyInstance() => CreateMessage();

  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class ListDevicesResponse extends $pb.GeneratedMessage {
  factory ListDevicesResponse({$core.List<DeviceInfo>? devices}) => CreateMessage();
  ListDevicesResponse._() : super();
  factory ListDevicesResponse.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);
  factory ListDevicesResponse.fromJson($core.String i) =>
      CreateMessage()..mergeFromJson(i);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'ListDevicesResponse',
    createEmptyInstance: CreateMessage,
  )..pc<DeviceInfo>(1, 'devices', $pb.PbFieldType.PM, subBuilder: DeviceInfo.CreateMessage);

  @$pb.TagNumber(1)
  $core.List<DeviceInfo> get devices => $_getList(0);

  static ListDevicesResponse CreateMessage() => ListDevicesResponse._();

  @$core.override
  ListDevicesResponse clone() => ListDevicesResponse()..mergeFromMessage(this);

  @$core.override
  ListDevicesResponse createEmptyInstance() => CreateMessage();

  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class DeviceInfo extends $pb.GeneratedMessage {
  factory DeviceInfo({
    $core.String? deviceId,
    $core.String? name,
    $core.String? firmwareVersion,
    $core.int? status,
    $core.int? batteryPercent,
  }) =>
      CreateMessage();
  DeviceInfo._() : super();
  factory DeviceInfo.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);
  factory DeviceInfo.fromJson($core.String i) =>
      CreateMessage()..mergeFromJson(i);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'DeviceInfo',
    createEmptyInstance: CreateMessage,
  )
    ..aOS(1, 'deviceId')
    ..aOS(2, 'name')
    ..aOS(3, 'firmwareVersion')
    ..a<$core.int>(4, 'status', $pb.PbFieldType.O3)
    ..a<$core.int>(5, 'batteryPercent', $pb.PbFieldType.O3);

  @$pb.TagNumber(1)
  $core.String get deviceId => $_getSZ(0);
  @$pb.TagNumber(1)
  set deviceId($core.String v) {
    $_setString(0, v);
  }

  @$pb.TagNumber(2)
  $core.String get name => $_getSZ(1);
  @$pb.TagNumber(2)
  set name($core.String v) {
    $_setString(1, v);
  }

  @$pb.TagNumber(3)
  $core.String get firmwareVersion => $_getSZ(2);
  @$pb.TagNumber(3)
  set firmwareVersion($core.String v) {
    $_setString(2, v);
  }

  @$pb.TagNumber(4)
  $core.int get status => $_getIZ(3);
  @$pb.TagNumber(4)
  set status($core.int v) {
    $_setSignedInt32(3, v);
  }

  @$pb.TagNumber(5)
  $core.int get batteryPercent => $_getIZ(4);
  @$pb.TagNumber(5)
  set batteryPercent($core.int v) {
    $_setSignedInt32(4, v);
  }

  static DeviceInfo CreateMessage() => DeviceInfo._();

  @$core.override
  DeviceInfo clone() => DeviceInfo()..mergeFromMessage(this);

  @$core.override
  DeviceInfo createEmptyInstance() => CreateMessage();

  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

// ---------------------------------------------------------------------------
// Measurement (StartSession, EndSession, GetHistory)
// ---------------------------------------------------------------------------

class StartSessionRequest extends $pb.GeneratedMessage {
  factory StartSessionRequest({
    $core.String? deviceId,
    $core.String? cartridgeId,
    $core.String? userId,
  }) =>
      CreateMessage();
  StartSessionRequest._() : super();
  factory StartSessionRequest.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);
  factory StartSessionRequest.fromJson($core.String i) =>
      CreateMessage()..mergeFromJson(i);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'StartSessionRequest',
    createEmptyInstance: CreateMessage,
  )
    ..aOS(1, 'deviceId')
    ..aOS(2, 'cartridgeId')
    ..aOS(3, 'userId');

  @$pb.TagNumber(1)
  $core.String get deviceId => $_getSZ(0);
  @$pb.TagNumber(1)
  set deviceId($core.String v) {
    $_setString(0, v);
  }

  @$pb.TagNumber(2)
  $core.String get cartridgeId => $_getSZ(1);
  @$pb.TagNumber(2)
  set cartridgeId($core.String v) {
    $_setString(1, v);
  }

  @$pb.TagNumber(3)
  $core.String get userId => $_getSZ(2);
  @$pb.TagNumber(3)
  set userId($core.String v) {
    $_setString(2, v);
  }

  static StartSessionRequest CreateMessage() => StartSessionRequest._();

  @$core.override
  StartSessionRequest clone() => StartSessionRequest()..mergeFromMessage(this);

  @$core.override
  StartSessionRequest createEmptyInstance() => CreateMessage();

  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class StartSessionResponse extends $pb.GeneratedMessage {
  factory StartSessionResponse({$core.String? sessionId}) => CreateMessage();
  StartSessionResponse._() : super();
  factory StartSessionResponse.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);
  factory StartSessionResponse.fromJson($core.String i) =>
      CreateMessage()..mergeFromJson(i);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'StartSessionResponse',
    createEmptyInstance: CreateMessage,
  )..aOS(1, 'sessionId');

  @$pb.TagNumber(1)
  $core.String get sessionId => $_getSZ(0);
  @$pb.TagNumber(1)
  set sessionId($core.String v) {
    $_setString(0, v);
  }

  static StartSessionResponse CreateMessage() => StartSessionResponse._();

  @$core.override
  StartSessionResponse clone() => StartSessionResponse()..mergeFromMessage(this);

  @$core.override
  StartSessionResponse createEmptyInstance() => CreateMessage();

  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class EndSessionRequest extends $pb.GeneratedMessage {
  factory EndSessionRequest({$core.String? sessionId}) => CreateMessage();
  EndSessionRequest._() : super();
  factory EndSessionRequest.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);
  factory EndSessionRequest.fromJson($core.String i) =>
      CreateMessage()..mergeFromJson(i);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'EndSessionRequest',
    createEmptyInstance: CreateMessage,
  )..aOS(1, 'sessionId');

  @$pb.TagNumber(1)
  $core.String get sessionId => $_getSZ(0);
  @$pb.TagNumber(1)
  set sessionId($core.String v) {
    $_setString(0, v);
  }

  static EndSessionRequest CreateMessage() => EndSessionRequest._();

  @$core.override
  EndSessionRequest clone() => EndSessionRequest()..mergeFromMessage(this);

  @$core.override
  EndSessionRequest createEmptyInstance() => CreateMessage();

  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class EndSessionResponse extends $pb.GeneratedMessage {
  factory EndSessionResponse({
    $core.String? sessionId,
    $core.int? totalMeasurements,
  }) =>
      CreateMessage();
  EndSessionResponse._() : super();
  factory EndSessionResponse.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);
  factory EndSessionResponse.fromJson($core.String i) =>
      CreateMessage()..mergeFromJson(i);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'EndSessionResponse',
    createEmptyInstance: CreateMessage,
  )
    ..aOS(1, 'sessionId')
    ..a<$core.int>(2, 'totalMeasurements', $pb.PbFieldType.O3);

  @$pb.TagNumber(1)
  $core.String get sessionId => $_getSZ(0);
  @$pb.TagNumber(1)
  set sessionId($core.String v) {
    $_setString(0, v);
  }

  @$pb.TagNumber(2)
  $core.int get totalMeasurements => $_getIZ(1);
  @$pb.TagNumber(2)
  set totalMeasurements($core.int v) {
    $_setSignedInt32(1, v);
  }

  static EndSessionResponse CreateMessage() => EndSessionResponse._();

  @$core.override
  EndSessionResponse clone() => EndSessionResponse()..mergeFromMessage(this);

  @$core.override
  EndSessionResponse createEmptyInstance() => CreateMessage();

  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class GetHistoryRequest extends $pb.GeneratedMessage {
  factory GetHistoryRequest({
    $core.String? userId,
    $core.int? limit,
    $core.int? offset,
  }) =>
      CreateMessage();
  GetHistoryRequest._() : super();
  factory GetHistoryRequest.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);
  factory GetHistoryRequest.fromJson($core.String i) =>
      CreateMessage()..mergeFromJson(i);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'GetHistoryRequest',
    createEmptyInstance: CreateMessage,
  )
    ..aOS(1, 'userId')
    ..a<$core.int>(4, 'limit', $pb.PbFieldType.O3)
    ..a<$core.int>(5, 'offset', $pb.PbFieldType.O3);

  @$pb.TagNumber(1)
  $core.String get userId => $_getSZ(0);
  @$pb.TagNumber(1)
  set userId($core.String v) {
    $_setString(0, v);
  }

  @$pb.TagNumber(4)
  $core.int get limit => $_getIZ(1);
  @$pb.TagNumber(4)
  set limit($core.int v) {
    $_setSignedInt32(1, v);
  }

  @$pb.TagNumber(5)
  $core.int get offset => $_getIZ(2);
  @$pb.TagNumber(5)
  set offset($core.int v) {
    $_setSignedInt32(2, v);
  }

  static GetHistoryRequest CreateMessage() => GetHistoryRequest._();

  @$core.override
  GetHistoryRequest clone() => GetHistoryRequest()..mergeFromMessage(this);

  @$core.override
  GetHistoryRequest createEmptyInstance() => CreateMessage();

  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class GetHistoryResponse extends $pb.GeneratedMessage {
  factory GetHistoryResponse({
    $core.List<MeasurementSummary>? measurements,
    $core.int? totalCount,
  }) =>
      CreateMessage();
  GetHistoryResponse._() : super();
  factory GetHistoryResponse.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);
  factory GetHistoryResponse.fromJson($core.String i) =>
      CreateMessage()..mergeFromJson(i);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'GetHistoryResponse',
    createEmptyInstance: CreateMessage,
  )
    ..pc<MeasurementSummary>(
        1, 'measurements', $pb.PbFieldType.PM, subBuilder: MeasurementSummary.CreateMessage)
    ..a<$core.int>(2, 'totalCount', $pb.PbFieldType.O3);

  @$pb.TagNumber(1)
  $core.List<MeasurementSummary> get measurements => $_getList(0);

  @$pb.TagNumber(2)
  $core.int get totalCount => $_getIZ(1);
  @$pb.TagNumber(2)
  set totalCount($core.int v) {
    $_setSignedInt32(1, v);
  }

  static GetHistoryResponse CreateMessage() => GetHistoryResponse._();

  @$core.override
  GetHistoryResponse clone() => GetHistoryResponse()..mergeFromMessage(this);

  @$core.override
  GetHistoryResponse createEmptyInstance() => CreateMessage();

  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class MeasurementSummary extends $pb.GeneratedMessage {
  factory MeasurementSummary({
    $core.String? sessionId,
    $core.String? cartridgeType,
    $core.double? primaryValue,
    $core.String? unit,
  }) =>
      CreateMessage();
  MeasurementSummary._() : super();
  factory MeasurementSummary.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);
  factory MeasurementSummary.fromJson($core.String i) =>
      CreateMessage()..mergeFromJson(i);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'MeasurementSummary',
    createEmptyInstance: CreateMessage,
  )
    ..aOS(1, 'sessionId')
    ..aOS(2, 'cartridgeType')
    ..a<$core.double>(3, 'primaryValue', $pb.PbFieldType.OD)
    ..aOS(4, 'unit');

  @$pb.TagNumber(1)
  $core.String get sessionId => $_getSZ(0);
  @$pb.TagNumber(1)
  set sessionId($core.String v) {
    $_setString(0, v);
  }

  @$pb.TagNumber(2)
  $core.String get cartridgeType => $_getSZ(1);
  @$pb.TagNumber(2)
  set cartridgeType($core.String v) {
    $_setString(1, v);
  }

  @$pb.TagNumber(3)
  $core.double get primaryValue => $_getN(2);
  @$pb.TagNumber(3)
  set primaryValue($core.double v) {
    $_setDouble(2, v);
  }

  @$pb.TagNumber(4)
  $core.String get unit => $_getSZ(3);
  @$pb.TagNumber(4)
  set unit($core.String v) {
    $_setString(3, v);
  }

  static MeasurementSummary CreateMessage() => MeasurementSummary._();

  @$core.override
  MeasurementSummary clone() => MeasurementSummary()..mergeFromMessage(this);

  @$core.override
  MeasurementSummary createEmptyInstance() => CreateMessage();

  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

// ---------------------------------------------------------------------------
// Admin — System Config
// ---------------------------------------------------------------------------

class ListSystemConfigsRequest extends $pb.GeneratedMessage {
  factory ListSystemConfigsRequest({
    $core.String? languageCode,
    $core.String? category,
    $core.bool? includeSecrets,
  }) {
    final msg = CreateMessage();
    if (languageCode != null) msg.languageCode = languageCode;
    if (category != null) msg.category = category;
    if (includeSecrets != null) msg.includeSecrets = includeSecrets;
    return msg;
  }
  ListSystemConfigsRequest._() : super();
  factory ListSystemConfigsRequest.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'ListSystemConfigsRequest',
    createEmptyInstance: CreateMessage,
  )
    ..aOS(1, 'languageCode')
    ..aOS(2, 'category')
    ..aOB(3, 'includeSecrets');

  @$pb.TagNumber(1)
  $core.String get languageCode => $_getSZ(0);
  @$pb.TagNumber(1)
  set languageCode($core.String v) { $_setString(0, v); }

  @$pb.TagNumber(2)
  $core.String get category => $_getSZ(1);
  @$pb.TagNumber(2)
  set category($core.String v) { $_setString(1, v); }

  @$pb.TagNumber(3)
  $core.bool get includeSecrets => $_getBF(2);
  @$pb.TagNumber(3)
  set includeSecrets($core.bool v) { $_setBool(2, v); }

  static ListSystemConfigsRequest CreateMessage() => ListSystemConfigsRequest._();

  @$core.override
  ListSystemConfigsRequest clone() => ListSystemConfigsRequest()..mergeFromMessage(this);
  @$core.override
  ListSystemConfigsRequest createEmptyInstance() => CreateMessage();
  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class ListSystemConfigsResponse extends $pb.GeneratedMessage {
  factory ListSystemConfigsResponse() => CreateMessage();
  ListSystemConfigsResponse._() : super();
  factory ListSystemConfigsResponse.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'ListSystemConfigsResponse',
    createEmptyInstance: CreateMessage,
  )
    ..pc<ConfigWithMeta>(1, 'configs', $pb.PbFieldType.PM,
        subBuilder: ConfigWithMeta.CreateMessage)
    ..m<$core.String, $core.int>(2, 'categoryCounts',
        entryClassName: 'ListSystemConfigsResponse.CategoryCountsEntry',
        keyFieldType: $pb.PbFieldType.OS,
        valueFieldType: $pb.PbFieldType.O3);

  @$pb.TagNumber(1)
  $core.List<ConfigWithMeta> get configs => $_getList(0);

  @$pb.TagNumber(2)
  $core.Map<$core.String, $core.int> get categoryCounts => $_getMap(1);

  static ListSystemConfigsResponse CreateMessage() => ListSystemConfigsResponse._();

  @$core.override
  ListSystemConfigsResponse clone() => ListSystemConfigsResponse()..mergeFromMessage(this);
  @$core.override
  ListSystemConfigsResponse createEmptyInstance() => CreateMessage();
  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class ConfigWithMeta extends $pb.GeneratedMessage {
  factory ConfigWithMeta() => CreateMessage();
  ConfigWithMeta._() : super();
  factory ConfigWithMeta.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'ConfigWithMeta',
    createEmptyInstance: CreateMessage,
  )
    ..aOS(1, 'key')
    ..aOS(2, 'value')
    ..aOS(3, 'rawValue')
    ..aOS(4, 'category')
    ..aOS(5, 'valueType')
    ..aOS(6, 'securityLevel')
    ..aOB(7, 'isRequired')
    ..aOS(8, 'defaultValue')
    ..pPS(9, 'allowedValues')
    ..aOS(10, 'validationRegex')
    ..a<$core.double>(11, 'validationMin', $pb.PbFieldType.OD)
    ..a<$core.double>(12, 'validationMax', $pb.PbFieldType.OD)
    ..aOS(13, 'dependsOn')
    ..aOS(14, 'dependsValue')
    ..aOS(15, 'envVarName')
    ..aOS(16, 'serviceName')
    ..aOB(17, 'restartRequired')
    ..aOS(20, 'displayName')
    ..aOS(21, 'description')
    ..aOS(22, 'placeholder')
    ..aOS(23, 'helpText')
    ..aOS(24, 'validationMessage')
    ..aOS(30, 'updatedBy');

  @$pb.TagNumber(1)
  $core.String get key => $_getSZ(0);
  @$pb.TagNumber(1)
  set key($core.String v) { $_setString(0, v); }

  @$pb.TagNumber(2)
  $core.String get value => $_getSZ(1);
  @$pb.TagNumber(2)
  set value($core.String v) { $_setString(1, v); }

  @$pb.TagNumber(3)
  $core.String get rawValue => $_getSZ(2);
  @$pb.TagNumber(3)
  set rawValue($core.String v) { $_setString(2, v); }

  @$pb.TagNumber(4)
  $core.String get category => $_getSZ(3);
  @$pb.TagNumber(4)
  set category($core.String v) { $_setString(3, v); }

  @$pb.TagNumber(5)
  $core.String get valueType => $_getSZ(4);
  @$pb.TagNumber(5)
  set valueType($core.String v) { $_setString(4, v); }

  @$pb.TagNumber(6)
  $core.String get securityLevel => $_getSZ(5);
  @$pb.TagNumber(6)
  set securityLevel($core.String v) { $_setString(5, v); }

  @$pb.TagNumber(7)
  $core.bool get isRequired => $_getBF(6);
  @$pb.TagNumber(7)
  set isRequired($core.bool v) { $_setBool(6, v); }

  @$pb.TagNumber(8)
  $core.String get defaultValue => $_getSZ(7);
  @$pb.TagNumber(8)
  set defaultValue($core.String v) { $_setString(7, v); }

  @$pb.TagNumber(9)
  $core.List<$core.String> get allowedValues => $_getList(8);

  @$pb.TagNumber(10)
  $core.String get validationRegex => $_getSZ(9);
  @$pb.TagNumber(10)
  set validationRegex($core.String v) { $_setString(9, v); }

  @$pb.TagNumber(11)
  $core.double get validationMin => $_getN(10);
  @$pb.TagNumber(11)
  set validationMin($core.double v) { $_setDouble(10, v); }

  @$pb.TagNumber(12)
  $core.double get validationMax => $_getN(11);
  @$pb.TagNumber(12)
  set validationMax($core.double v) { $_setDouble(11, v); }

  @$pb.TagNumber(13)
  $core.String get dependsOn => $_getSZ(12);
  @$pb.TagNumber(13)
  set dependsOn($core.String v) { $_setString(12, v); }

  @$pb.TagNumber(14)
  $core.String get dependsValue => $_getSZ(13);
  @$pb.TagNumber(14)
  set dependsValue($core.String v) { $_setString(13, v); }

  @$pb.TagNumber(15)
  $core.String get envVarName => $_getSZ(14);
  @$pb.TagNumber(15)
  set envVarName($core.String v) { $_setString(14, v); }

  @$pb.TagNumber(16)
  $core.String get serviceName => $_getSZ(15);
  @$pb.TagNumber(16)
  set serviceName($core.String v) { $_setString(15, v); }

  @$pb.TagNumber(17)
  $core.bool get restartRequired => $_getBF(16);
  @$pb.TagNumber(17)
  set restartRequired($core.bool v) { $_setBool(16, v); }

  @$pb.TagNumber(20)
  $core.String get displayName => $_getSZ(17);
  @$pb.TagNumber(20)
  set displayName($core.String v) { $_setString(17, v); }

  @$pb.TagNumber(21)
  $core.String get description => $_getSZ(18);
  @$pb.TagNumber(21)
  set description($core.String v) { $_setString(18, v); }

  @$pb.TagNumber(22)
  $core.String get placeholder => $_getSZ(19);
  @$pb.TagNumber(22)
  set placeholder($core.String v) { $_setString(19, v); }

  @$pb.TagNumber(23)
  $core.String get helpText => $_getSZ(20);
  @$pb.TagNumber(23)
  set helpText($core.String v) { $_setString(20, v); }

  @$pb.TagNumber(24)
  $core.String get validationMessage => $_getSZ(21);
  @$pb.TagNumber(24)
  set validationMessage($core.String v) { $_setString(21, v); }

  @$pb.TagNumber(30)
  $core.String get updatedBy => $_getSZ(22);
  @$pb.TagNumber(30)
  set updatedBy($core.String v) { $_setString(22, v); }

  static ConfigWithMeta CreateMessage() => ConfigWithMeta._();

  @$core.override
  ConfigWithMeta clone() => ConfigWithMeta()..mergeFromMessage(this);
  @$core.override
  ConfigWithMeta createEmptyInstance() => CreateMessage();
  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class SetSystemConfigRequest extends $pb.GeneratedMessage {
  factory SetSystemConfigRequest({
    $core.String? key,
    $core.String? value,
    $core.String? description,
  }) {
    final msg = CreateMessage();
    if (key != null) msg.key = key;
    if (value != null) msg.value = value;
    if (description != null) msg.description = description;
    return msg;
  }
  SetSystemConfigRequest._() : super();
  factory SetSystemConfigRequest.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'SetSystemConfigRequest',
    createEmptyInstance: CreateMessage,
  )
    ..aOS(1, 'key')
    ..aOS(2, 'value')
    ..aOS(3, 'description');

  @$pb.TagNumber(1)
  $core.String get key => $_getSZ(0);
  @$pb.TagNumber(1)
  set key($core.String v) { $_setString(0, v); }

  @$pb.TagNumber(2)
  $core.String get value => $_getSZ(1);
  @$pb.TagNumber(2)
  set value($core.String v) { $_setString(1, v); }

  @$pb.TagNumber(3)
  $core.String get description => $_getSZ(2);
  @$pb.TagNumber(3)
  set description($core.String v) { $_setString(2, v); }

  static SetSystemConfigRequest CreateMessage() => SetSystemConfigRequest._();

  @$core.override
  SetSystemConfigRequest clone() => SetSystemConfigRequest()..mergeFromMessage(this);
  @$core.override
  SetSystemConfigRequest createEmptyInstance() => CreateMessage();
  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class SystemConfig extends $pb.GeneratedMessage {
  factory SystemConfig() => CreateMessage();
  SystemConfig._() : super();
  factory SystemConfig.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'SystemConfig',
    createEmptyInstance: CreateMessage,
  )
    ..aOS(1, 'key')
    ..aOS(2, 'value')
    ..aOS(3, 'description')
    ..aOS(5, 'updatedBy');

  @$pb.TagNumber(1)
  $core.String get key => $_getSZ(0);
  @$pb.TagNumber(1)
  set key($core.String v) { $_setString(0, v); }

  @$pb.TagNumber(2)
  $core.String get value => $_getSZ(1);
  @$pb.TagNumber(2)
  set value($core.String v) { $_setString(1, v); }

  @$pb.TagNumber(3)
  $core.String get description => $_getSZ(2);
  @$pb.TagNumber(3)
  set description($core.String v) { $_setString(2, v); }

  @$pb.TagNumber(5)
  $core.String get updatedBy => $_getSZ(3);
  @$pb.TagNumber(5)
  set updatedBy($core.String v) { $_setString(3, v); }

  static SystemConfig CreateMessage() => SystemConfig._();

  @$core.override
  SystemConfig clone() => SystemConfig()..mergeFromMessage(this);
  @$core.override
  SystemConfig createEmptyInstance() => CreateMessage();
  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class GetConfigWithMetaRequest extends $pb.GeneratedMessage {
  factory GetConfigWithMetaRequest({
    $core.String? key,
    $core.String? languageCode,
  }) {
    final msg = CreateMessage();
    if (key != null) msg.key = key;
    if (languageCode != null) msg.languageCode = languageCode;
    return msg;
  }
  GetConfigWithMetaRequest._() : super();
  factory GetConfigWithMetaRequest.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'GetConfigWithMetaRequest',
    createEmptyInstance: CreateMessage,
  )
    ..aOS(1, 'key')
    ..aOS(2, 'languageCode');

  @$pb.TagNumber(1)
  $core.String get key => $_getSZ(0);
  @$pb.TagNumber(1)
  set key($core.String v) { $_setString(0, v); }

  @$pb.TagNumber(2)
  $core.String get languageCode => $_getSZ(1);
  @$pb.TagNumber(2)
  set languageCode($core.String v) { $_setString(1, v); }

  static GetConfigWithMetaRequest CreateMessage() => GetConfigWithMetaRequest._();

  @$core.override
  GetConfigWithMetaRequest clone() => GetConfigWithMetaRequest()..mergeFromMessage(this);
  @$core.override
  GetConfigWithMetaRequest createEmptyInstance() => CreateMessage();
  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class ValidateConfigValueRequest extends $pb.GeneratedMessage {
  factory ValidateConfigValueRequest({
    $core.String? key,
    $core.String? value,
  }) {
    final msg = CreateMessage();
    if (key != null) msg.key = key;
    if (value != null) msg.value = value;
    return msg;
  }
  ValidateConfigValueRequest._() : super();
  factory ValidateConfigValueRequest.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'ValidateConfigValueRequest',
    createEmptyInstance: CreateMessage,
  )
    ..aOS(1, 'key')
    ..aOS(2, 'value');

  @$pb.TagNumber(1)
  $core.String get key => $_getSZ(0);
  @$pb.TagNumber(1)
  set key($core.String v) { $_setString(0, v); }

  @$pb.TagNumber(2)
  $core.String get value => $_getSZ(1);
  @$pb.TagNumber(2)
  set value($core.String v) { $_setString(1, v); }

  static ValidateConfigValueRequest CreateMessage() => ValidateConfigValueRequest._();

  @$core.override
  ValidateConfigValueRequest clone() => ValidateConfigValueRequest()..mergeFromMessage(this);
  @$core.override
  ValidateConfigValueRequest createEmptyInstance() => CreateMessage();
  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class ValidateConfigValueResponse extends $pb.GeneratedMessage {
  factory ValidateConfigValueResponse() => CreateMessage();
  ValidateConfigValueResponse._() : super();
  factory ValidateConfigValueResponse.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'ValidateConfigValueResponse',
    createEmptyInstance: CreateMessage,
  )
    ..aOB(1, 'valid')
    ..aOS(2, 'errorMessage')
    ..pPS(3, 'suggestions');

  @$pb.TagNumber(1)
  $core.bool get valid => $_getBF(0);
  @$pb.TagNumber(1)
  set valid($core.bool v) { $_setBool(0, v); }

  @$pb.TagNumber(2)
  $core.String get errorMessage => $_getSZ(1);
  @$pb.TagNumber(2)
  set errorMessage($core.String v) { $_setString(1, v); }

  @$pb.TagNumber(3)
  $core.List<$core.String> get suggestions => $_getList(2);

  static ValidateConfigValueResponse CreateMessage() => ValidateConfigValueResponse._();

  @$core.override
  ValidateConfigValueResponse clone() => ValidateConfigValueResponse()..mergeFromMessage(this);
  @$core.override
  ValidateConfigValueResponse createEmptyInstance() => CreateMessage();
  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

// ---------------------------------------------------------------------------
// AI Inference (Phase 2)
// ---------------------------------------------------------------------------

class AnalyzeMeasurementRequest extends $pb.GeneratedMessage {
  factory AnalyzeMeasurementRequest({
    $core.String? userId,
    $core.String? measurementId,
  }) =>
      CreateMessage();
  AnalyzeMeasurementRequest._() : super();
  factory AnalyzeMeasurementRequest.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'AnalyzeMeasurementRequest',
    createEmptyInstance: CreateMessage,
  )
    ..aOS(1, 'userId')
    ..aOS(2, 'measurementId');

  @$pb.TagNumber(1)
  $core.String get userId => $_getSZ(0);
  @$pb.TagNumber(1)
  set userId($core.String v) { $_setString(0, v); }

  @$pb.TagNumber(2)
  $core.String get measurementId => $_getSZ(1);
  @$pb.TagNumber(2)
  set measurementId($core.String v) { $_setString(1, v); }

  static AnalyzeMeasurementRequest CreateMessage() => AnalyzeMeasurementRequest._();

  @$core.override
  AnalyzeMeasurementRequest clone() => AnalyzeMeasurementRequest()..mergeFromMessage(this);
  @$core.override
  AnalyzeMeasurementRequest createEmptyInstance() => CreateMessage();
  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class AnalysisResult extends $pb.GeneratedMessage {
  factory AnalysisResult() => CreateMessage();
  AnalysisResult._() : super();
  factory AnalysisResult.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'AnalysisResult',
    createEmptyInstance: CreateMessage,
  )
    ..aOS(1, 'analysisId')
    ..aOS(2, 'userId')
    ..aOS(3, 'measurementId')
    ..a<$core.double>(6, 'overallHealthScore', $pb.PbFieldType.OD)
    ..aOS(7, 'summary');

  @$pb.TagNumber(1)
  $core.String get analysisId => $_getSZ(0);
  @$pb.TagNumber(1)
  set analysisId($core.String v) { $_setString(0, v); }

  @$pb.TagNumber(2)
  $core.String get userId => $_getSZ(1);
  @$pb.TagNumber(2)
  set userId($core.String v) { $_setString(1, v); }

  @$pb.TagNumber(3)
  $core.String get measurementId => $_getSZ(2);
  @$pb.TagNumber(3)
  set measurementId($core.String v) { $_setString(2, v); }

  @$pb.TagNumber(6)
  $core.double get overallHealthScore => $_getN(3);
  @$pb.TagNumber(6)
  set overallHealthScore($core.double v) { $_setDouble(3, v); }

  @$pb.TagNumber(7)
  $core.String get summary => $_getSZ(4);
  @$pb.TagNumber(7)
  set summary($core.String v) { $_setString(4, v); }

  static AnalysisResult CreateMessage() => AnalysisResult._();

  @$core.override
  AnalysisResult clone() => AnalysisResult()..mergeFromMessage(this);
  @$core.override
  AnalysisResult createEmptyInstance() => CreateMessage();
  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class GetHealthScoreRequest extends $pb.GeneratedMessage {
  factory GetHealthScoreRequest({
    $core.String? userId,
    $core.int? days,
  }) =>
      CreateMessage();
  GetHealthScoreRequest._() : super();
  factory GetHealthScoreRequest.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'GetHealthScoreRequest',
    createEmptyInstance: CreateMessage,
  )
    ..aOS(1, 'userId')
    ..a<$core.int>(2, 'days', $pb.PbFieldType.O3);

  @$pb.TagNumber(1)
  $core.String get userId => $_getSZ(0);
  @$pb.TagNumber(1)
  set userId($core.String v) { $_setString(0, v); }

  @$pb.TagNumber(2)
  $core.int get days => $_getIZ(1);
  @$pb.TagNumber(2)
  set days($core.int v) { $_setSignedInt32(1, v); }

  static GetHealthScoreRequest CreateMessage() => GetHealthScoreRequest._();

  @$core.override
  GetHealthScoreRequest clone() => GetHealthScoreRequest()..mergeFromMessage(this);
  @$core.override
  GetHealthScoreRequest createEmptyInstance() => CreateMessage();
  @$core.override
  $pb.BuilderInfo get info_ => _i;
}

class HealthScoreResponse extends $pb.GeneratedMessage {
  factory HealthScoreResponse() => CreateMessage();
  HealthScoreResponse._() : super();
  factory HealthScoreResponse.fromBuffer($core.List<$core.int> i,
          [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) =>
      CreateMessage()..mergeFromBuffer(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
    'HealthScoreResponse',
    createEmptyInstance: CreateMessage,
  )
    ..aOS(1, 'userId')
    ..a<$core.double>(2, 'overallScore', $pb.PbFieldType.OD)
    ..aOS(4, 'trend')
    ..aOS(5, 'recommendation');

  @$pb.TagNumber(1)
  $core.String get userId => $_getSZ(0);
  @$pb.TagNumber(1)
  set userId($core.String v) { $_setString(0, v); }

  @$pb.TagNumber(2)
  $core.double get overallScore => $_getN(1);
  @$pb.TagNumber(2)
  set overallScore($core.double v) { $_setDouble(1, v); }

  @$pb.TagNumber(4)
  $core.String get trend => $_getSZ(2);
  @$pb.TagNumber(4)
  set trend($core.String v) { $_setString(2, v); }

  @$pb.TagNumber(5)
  $core.String get recommendation => $_getSZ(3);
  @$pb.TagNumber(5)
  set recommendation($core.String v) { $_setString(3, v); }

  static HealthScoreResponse CreateMessage() => HealthScoreResponse._();

  @$core.override
  HealthScoreResponse clone() => HealthScoreResponse()..mergeFromMessage(this);
  @$core.override
  HealthScoreResponse createEmptyInstance() => CreateMessage();
  @$core.override
  $pb.BuilderInfo get info_ => _i;
}
