// ManPaSik gRPC 클라이언트 스텁 - 수동 생성
// ignore_for_file: type=lint

import 'package:grpc/grpc.dart';
import 'manpasik.pb.dart';

const _authService = 'manpasik.v1.AuthService';
const _userService = 'manpasik.v1.UserService';
const _deviceService = 'manpasik.v1.DeviceService';
const _measurementService = 'manpasik.v1.MeasurementService';
const _adminService = 'manpasik.v1.AdminService';
const _aiInferenceService = 'manpasik.v1.AiInferenceService';

/// Auth 서비스 gRPC 클라이언트
class AuthServiceClient extends Client {
  AuthServiceClient(
    super.channel, {
    super.interceptors,
    super.options,
  });

  static final _register =
      ClientMethod<RegisterRequest, RegisterResponse>(
    '/$_authService/Register',
    (RegisterRequest v) => v.writeToBuffer(),
    (List<int> v) => RegisterResponse.fromBuffer(v),
  );
  static final _login = ClientMethod<LoginRequest, LoginResponse>(
    '/$_authService/Login',
    (LoginRequest v) => v.writeToBuffer(),
    (List<int> v) => LoginResponse.fromBuffer(v),
  );
  static final _refreshToken =
      ClientMethod<RefreshTokenRequest, LoginResponse>(
    '/$_authService/RefreshToken',
    (RefreshTokenRequest v) => v.writeToBuffer(),
    (List<int> v) => LoginResponse.fromBuffer(v),
  );
  static final _logout = ClientMethod<LogoutRequest, LogoutResponse>(
    '/$_authService/Logout',
    (LogoutRequest v) => v.writeToBuffer(),
    (List<int> v) => LogoutResponse.fromBuffer(v),
  );

  ResponseFuture<RegisterResponse> register(
    RegisterRequest request, {
    CallOptions? options,
  }) =>
      $createUnaryCall(_register, request, options: options);

  ResponseFuture<LoginResponse> login(
    LoginRequest request, {
    CallOptions? options,
  }) =>
      $createUnaryCall(_login, request, options: options);

  ResponseFuture<LoginResponse> refreshToken(
    RefreshTokenRequest request, {
    CallOptions? options,
  }) =>
      $createUnaryCall(_refreshToken, request, options: options);

  ResponseFuture<LogoutResponse> logout(
    LogoutRequest request, {
    CallOptions? options,
  }) =>
      $createUnaryCall(_logout, request, options: options);
}

/// User 서비스 gRPC 클라이언트
class UserServiceClient extends Client {
  UserServiceClient(
    super.channel, {
    super.interceptors,
    super.options,
  });

  static final _getProfile =
      ClientMethod<GetProfileRequest, UserProfile>(
    '/$_userService/GetProfile',
    (GetProfileRequest v) => v.writeToBuffer(),
    (List<int> v) => UserProfile.fromBuffer(v),
  );
  static final _getSubscription =
      ClientMethod<GetSubscriptionRequest, SubscriptionInfo>(
    '/$_userService/GetSubscription',
    (GetSubscriptionRequest v) => v.writeToBuffer(),
    (List<int> v) => SubscriptionInfo.fromBuffer(v),
  );

  ResponseFuture<UserProfile> getProfile(
    GetProfileRequest request, {
    CallOptions? options,
  }) =>
      $createUnaryCall(_getProfile, request, options: options);

  ResponseFuture<SubscriptionInfo> getSubscription(
    GetSubscriptionRequest request, {
    CallOptions? options,
  }) =>
      $createUnaryCall(_getSubscription, request, options: options);
}

/// Device 서비스 gRPC 클라이언트
class DeviceServiceClient extends Client {
  DeviceServiceClient(
    super.channel, {
    super.interceptors,
    super.options,
  });

  static final _listDevices =
      ClientMethod<ListDevicesRequest, ListDevicesResponse>(
    '/$_deviceService/ListDevices',
    (ListDevicesRequest v) => v.writeToBuffer(),
    (List<int> v) => ListDevicesResponse.fromBuffer(v),
  );

  ResponseFuture<ListDevicesResponse> listDevices(
    ListDevicesRequest request, {
    CallOptions? options,
  }) =>
      $createUnaryCall(_listDevices, request, options: options);
}

/// Measurement 서비스 gRPC 클라이언트
class MeasurementServiceClient extends Client {
  MeasurementServiceClient(
    super.channel, {
    super.interceptors,
    super.options,
  });

  static final _startSession =
      ClientMethod<StartSessionRequest, StartSessionResponse>(
    '/$_measurementService/StartSession',
    (StartSessionRequest v) => v.writeToBuffer(),
    (List<int> v) => StartSessionResponse.fromBuffer(v),
  );
  static final _endSession =
      ClientMethod<EndSessionRequest, EndSessionResponse>(
    '/$_measurementService/EndSession',
    (EndSessionRequest v) => v.writeToBuffer(),
    (List<int> v) => EndSessionResponse.fromBuffer(v),
  );
  static final _getMeasurementHistory =
      ClientMethod<GetHistoryRequest, GetHistoryResponse>(
    '/$_measurementService/GetMeasurementHistory',
    (GetHistoryRequest v) => v.writeToBuffer(),
    (List<int> v) => GetHistoryResponse.fromBuffer(v),
  );

  ResponseFuture<StartSessionResponse> startSession(
    StartSessionRequest request, {
    CallOptions? options,
  }) =>
      $createUnaryCall(_startSession, request, options: options);

  ResponseFuture<EndSessionResponse> endSession(
    EndSessionRequest request, {
    CallOptions? options,
  }) =>
      $createUnaryCall(_endSession, request, options: options);

  ResponseFuture<GetHistoryResponse> getMeasurementHistory(
    GetHistoryRequest request, {
    CallOptions? options,
  }) =>
      $createUnaryCall(_getMeasurementHistory, request, options: options);
}

/// Admin 서비스 gRPC 클라이언트
class AdminServiceClient extends Client {
  AdminServiceClient(
    super.channel, {
    super.interceptors,
    super.options,
  });

  static final _listSystemConfigs =
      ClientMethod<ListSystemConfigsRequest, ListSystemConfigsResponse>(
    '/$_adminService/ListSystemConfigs',
    (ListSystemConfigsRequest v) => v.writeToBuffer(),
    (List<int> v) => ListSystemConfigsResponse.fromBuffer(v),
  );
  static final _getConfigWithMeta =
      ClientMethod<GetConfigWithMetaRequest, ConfigWithMeta>(
    '/$_adminService/GetConfigWithMeta',
    (GetConfigWithMetaRequest v) => v.writeToBuffer(),
    (List<int> v) => ConfigWithMeta.fromBuffer(v),
  );
  static final _setSystemConfig =
      ClientMethod<SetSystemConfigRequest, SystemConfig>(
    '/$_adminService/SetSystemConfig',
    (SetSystemConfigRequest v) => v.writeToBuffer(),
    (List<int> v) => SystemConfig.fromBuffer(v),
  );
  static final _validateConfigValue =
      ClientMethod<ValidateConfigValueRequest, ValidateConfigValueResponse>(
    '/$_adminService/ValidateConfigValue',
    (ValidateConfigValueRequest v) => v.writeToBuffer(),
    (List<int> v) => ValidateConfigValueResponse.fromBuffer(v),
  );

  ResponseFuture<ListSystemConfigsResponse> listSystemConfigs(
    ListSystemConfigsRequest request, {
    CallOptions? options,
  }) =>
      $createUnaryCall(_listSystemConfigs, request, options: options);

  ResponseFuture<ConfigWithMeta> getConfigWithMeta(
    GetConfigWithMetaRequest request, {
    CallOptions? options,
  }) =>
      $createUnaryCall(_getConfigWithMeta, request, options: options);

  ResponseFuture<SystemConfig> setSystemConfig(
    SetSystemConfigRequest request, {
    CallOptions? options,
  }) =>
      $createUnaryCall(_setSystemConfig, request, options: options);

  ResponseFuture<ValidateConfigValueResponse> validateConfigValue(
    ValidateConfigValueRequest request, {
    CallOptions? options,
  }) =>
      $createUnaryCall(_validateConfigValue, request, options: options);
}

/// AI Inference 서비스 gRPC 클라이언트
class AIInferenceServiceClient extends Client {
  AIInferenceServiceClient(
    super.channel, {
    super.interceptors,
    super.options,
  });

  static final _analyzeMeasurement =
      ClientMethod<AnalyzeMeasurementRequest, AnalysisResult>(
    '/$_aiInferenceService/AnalyzeMeasurement',
    (AnalyzeMeasurementRequest v) => v.writeToBuffer(),
    (List<int> v) => AnalysisResult.fromBuffer(v),
  );
  static final _getHealthScore =
      ClientMethod<GetHealthScoreRequest, HealthScoreResponse>(
    '/$_aiInferenceService/GetHealthScore',
    (GetHealthScoreRequest v) => v.writeToBuffer(),
    (List<int> v) => HealthScoreResponse.fromBuffer(v),
  );

  ResponseFuture<AnalysisResult> analyzeMeasurement(
    AnalyzeMeasurementRequest request, {
    CallOptions? options,
  }) =>
      $createUnaryCall(_analyzeMeasurement, request, options: options);

  ResponseFuture<HealthScoreResponse> getHealthScore(
    GetHealthScoreRequest request, {
    CallOptions? options,
  }) =>
      $createUnaryCall(_getHealthScore, request, options: options);
}
