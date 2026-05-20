// This is a generated file - do not edit.
//
// Generated from manpasik.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:grpc/service_api.dart' as $grpc;
import 'package:protobuf/protobuf.dart' as $pb;

import 'manpasik.pb.dart' as $0;

export 'manpasik.pb.dart';

@$pb.GrpcServiceName('manpasik.v1.AuthService')
class AuthServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  AuthServiceClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.RegisterResponse> register(
    $0.RegisterRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$register, request, options: options);
  }

  $grpc.ResponseFuture<$0.LoginResponse> login(
    $0.LoginRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$login, request, options: options);
  }

  $grpc.ResponseFuture<$0.LoginResponse> socialLogin(
    $0.SocialLoginRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$socialLogin, request, options: options);
  }

  $grpc.ResponseFuture<$0.LoginResponse> refreshToken(
    $0.RefreshTokenRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$refreshToken, request, options: options);
  }

  $grpc.ResponseFuture<$0.LogoutResponse> logout(
    $0.LogoutRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$logout, request, options: options);
  }

  $grpc.ResponseFuture<$0.ValidateTokenResponse> validateToken(
    $0.ValidateTokenRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$validateToken, request, options: options);
  }

  $grpc.ResponseFuture<$0.ResetPasswordResponse> resetPassword(
    $0.ResetPasswordRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$resetPassword, request, options: options);
  }

  // method descriptors

  static final _$register =
      $grpc.ClientMethod<$0.RegisterRequest, $0.RegisterResponse>(
          '/manpasik.v1.AuthService/Register',
          ($0.RegisterRequest value) => value.writeToBuffer(),
          $0.RegisterResponse.fromBuffer);
  static final _$login = $grpc.ClientMethod<$0.LoginRequest, $0.LoginResponse>(
      '/manpasik.v1.AuthService/Login',
      ($0.LoginRequest value) => value.writeToBuffer(),
      $0.LoginResponse.fromBuffer);
  static final _$socialLogin =
      $grpc.ClientMethod<$0.SocialLoginRequest, $0.LoginResponse>(
          '/manpasik.v1.AuthService/SocialLogin',
          ($0.SocialLoginRequest value) => value.writeToBuffer(),
          $0.LoginResponse.fromBuffer);
  static final _$refreshToken =
      $grpc.ClientMethod<$0.RefreshTokenRequest, $0.LoginResponse>(
          '/manpasik.v1.AuthService/RefreshToken',
          ($0.RefreshTokenRequest value) => value.writeToBuffer(),
          $0.LoginResponse.fromBuffer);
  static final _$logout =
      $grpc.ClientMethod<$0.LogoutRequest, $0.LogoutResponse>(
          '/manpasik.v1.AuthService/Logout',
          ($0.LogoutRequest value) => value.writeToBuffer(),
          $0.LogoutResponse.fromBuffer);
  static final _$validateToken =
      $grpc.ClientMethod<$0.ValidateTokenRequest, $0.ValidateTokenResponse>(
          '/manpasik.v1.AuthService/ValidateToken',
          ($0.ValidateTokenRequest value) => value.writeToBuffer(),
          $0.ValidateTokenResponse.fromBuffer);
  static final _$resetPassword =
      $grpc.ClientMethod<$0.ResetPasswordRequest, $0.ResetPasswordResponse>(
          '/manpasik.v1.AuthService/ResetPassword',
          ($0.ResetPasswordRequest value) => value.writeToBuffer(),
          $0.ResetPasswordResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.AuthService')
abstract class AuthServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.AuthService';

  AuthServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.RegisterRequest, $0.RegisterResponse>(
        'Register',
        register_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.RegisterRequest.fromBuffer(value),
        ($0.RegisterResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.LoginRequest, $0.LoginResponse>(
        'Login',
        login_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.LoginRequest.fromBuffer(value),
        ($0.LoginResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SocialLoginRequest, $0.LoginResponse>(
        'SocialLogin',
        socialLogin_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SocialLoginRequest.fromBuffer(value),
        ($0.LoginResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RefreshTokenRequest, $0.LoginResponse>(
        'RefreshToken',
        refreshToken_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RefreshTokenRequest.fromBuffer(value),
        ($0.LoginResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.LogoutRequest, $0.LogoutResponse>(
        'Logout',
        logout_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.LogoutRequest.fromBuffer(value),
        ($0.LogoutResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ValidateTokenRequest, $0.ValidateTokenResponse>(
            'ValidateToken',
            validateToken_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ValidateTokenRequest.fromBuffer(value),
            ($0.ValidateTokenResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ResetPasswordRequest, $0.ResetPasswordResponse>(
            'ResetPassword',
            resetPassword_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ResetPasswordRequest.fromBuffer(value),
            ($0.ResetPasswordResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.RegisterResponse> register_Pre($grpc.ServiceCall $call,
      $async.Future<$0.RegisterRequest> $request) async {
    return register($call, await $request);
  }

  $async.Future<$0.RegisterResponse> register(
      $grpc.ServiceCall call, $0.RegisterRequest request);

  $async.Future<$0.LoginResponse> login_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.LoginRequest> $request) async {
    return login($call, await $request);
  }

  $async.Future<$0.LoginResponse> login(
      $grpc.ServiceCall call, $0.LoginRequest request);

  $async.Future<$0.LoginResponse> socialLogin_Pre($grpc.ServiceCall $call,
      $async.Future<$0.SocialLoginRequest> $request) async {
    return socialLogin($call, await $request);
  }

  $async.Future<$0.LoginResponse> socialLogin(
      $grpc.ServiceCall call, $0.SocialLoginRequest request);

  $async.Future<$0.LoginResponse> refreshToken_Pre($grpc.ServiceCall $call,
      $async.Future<$0.RefreshTokenRequest> $request) async {
    return refreshToken($call, await $request);
  }

  $async.Future<$0.LoginResponse> refreshToken(
      $grpc.ServiceCall call, $0.RefreshTokenRequest request);

  $async.Future<$0.LogoutResponse> logout_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.LogoutRequest> $request) async {
    return logout($call, await $request);
  }

  $async.Future<$0.LogoutResponse> logout(
      $grpc.ServiceCall call, $0.LogoutRequest request);

  $async.Future<$0.ValidateTokenResponse> validateToken_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ValidateTokenRequest> $request) async {
    return validateToken($call, await $request);
  }

  $async.Future<$0.ValidateTokenResponse> validateToken(
      $grpc.ServiceCall call, $0.ValidateTokenRequest request);

  $async.Future<$0.ResetPasswordResponse> resetPassword_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ResetPasswordRequest> $request) async {
    return resetPassword($call, await $request);
  }

  $async.Future<$0.ResetPasswordResponse> resetPassword(
      $grpc.ServiceCall call, $0.ResetPasswordRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.MeasurementService')
class MeasurementServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  MeasurementServiceClient(super.channel, {super.options, super.interceptors});

  /// 측정 세션 시작
  $grpc.ResponseFuture<$0.StartSessionResponse> startSession(
    $0.StartSessionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startSession, request, options: options);
  }

  /// 측정 데이터 스트림
  $grpc.ResponseStream<$0.MeasurementResult> streamMeasurement(
    $async.Stream<$0.MeasurementData> request, {
    $grpc.CallOptions? options,
  }) {
    return $createStreamingCall(_$streamMeasurement, request, options: options);
  }

  /// 측정 세션 종료
  $grpc.ResponseFuture<$0.EndSessionResponse> endSession(
    $0.EndSessionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$endSession, request, options: options);
  }

  /// 측정 기록 조회
  $grpc.ResponseFuture<$0.GetHistoryResponse> getMeasurementHistory(
    $0.GetHistoryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getMeasurementHistory, request, options: options);
  }

  /// Phase 4: FHIR export
  $grpc.ResponseFuture<$0.ExportFHIRResponse> exportSingleMeasurement(
    $0.ExportSingleMeasurementRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$exportSingleMeasurement, request,
        options: options);
  }

  $grpc.ResponseFuture<$0.ExportFHIRResponse> exportToFHIRObservations(
    $0.ExportToFHIRObservationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$exportToFHIRObservations, request,
        options: options);
  }

  /// Phase 5: 디지털 트윈 동기화
  $grpc.ResponseFuture<$0.SyncDigitalTwinResponse> syncDigitalTwin(
    $0.SyncDigitalTwinRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$syncDigitalTwin, request, options: options);
  }

  /// Phase 5: 보정 상태 조회
  $grpc.ResponseFuture<$0.GetCalibrationStatusResponse> getCalibrationStatus(
    $0.GetCalibrationStatusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getCalibrationStatus, request, options: options);
  }

  // method descriptors

  static final _$startSession =
      $grpc.ClientMethod<$0.StartSessionRequest, $0.StartSessionResponse>(
          '/manpasik.v1.MeasurementService/StartSession',
          ($0.StartSessionRequest value) => value.writeToBuffer(),
          $0.StartSessionResponse.fromBuffer);
  static final _$streamMeasurement =
      $grpc.ClientMethod<$0.MeasurementData, $0.MeasurementResult>(
          '/manpasik.v1.MeasurementService/StreamMeasurement',
          ($0.MeasurementData value) => value.writeToBuffer(),
          $0.MeasurementResult.fromBuffer);
  static final _$endSession =
      $grpc.ClientMethod<$0.EndSessionRequest, $0.EndSessionResponse>(
          '/manpasik.v1.MeasurementService/EndSession',
          ($0.EndSessionRequest value) => value.writeToBuffer(),
          $0.EndSessionResponse.fromBuffer);
  static final _$getMeasurementHistory =
      $grpc.ClientMethod<$0.GetHistoryRequest, $0.GetHistoryResponse>(
          '/manpasik.v1.MeasurementService/GetMeasurementHistory',
          ($0.GetHistoryRequest value) => value.writeToBuffer(),
          $0.GetHistoryResponse.fromBuffer);
  static final _$exportSingleMeasurement = $grpc.ClientMethod<
          $0.ExportSingleMeasurementRequest, $0.ExportFHIRResponse>(
      '/manpasik.v1.MeasurementService/ExportSingleMeasurement',
      ($0.ExportSingleMeasurementRequest value) => value.writeToBuffer(),
      $0.ExportFHIRResponse.fromBuffer);
  static final _$exportToFHIRObservations = $grpc.ClientMethod<
          $0.ExportToFHIRObservationsRequest, $0.ExportFHIRResponse>(
      '/manpasik.v1.MeasurementService/ExportToFHIRObservations',
      ($0.ExportToFHIRObservationsRequest value) => value.writeToBuffer(),
      $0.ExportFHIRResponse.fromBuffer);
  static final _$syncDigitalTwin =
      $grpc.ClientMethod<$0.SyncDigitalTwinRequest, $0.SyncDigitalTwinResponse>(
          '/manpasik.v1.MeasurementService/SyncDigitalTwin',
          ($0.SyncDigitalTwinRequest value) => value.writeToBuffer(),
          $0.SyncDigitalTwinResponse.fromBuffer);
  static final _$getCalibrationStatus = $grpc.ClientMethod<
          $0.GetCalibrationStatusRequest, $0.GetCalibrationStatusResponse>(
      '/manpasik.v1.MeasurementService/GetCalibrationStatus',
      ($0.GetCalibrationStatusRequest value) => value.writeToBuffer(),
      $0.GetCalibrationStatusResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.MeasurementService')
abstract class MeasurementServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.MeasurementService';

  MeasurementServiceBase() {
    $addMethod(
        $grpc.ServiceMethod<$0.StartSessionRequest, $0.StartSessionResponse>(
            'StartSession',
            startSession_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.StartSessionRequest.fromBuffer(value),
            ($0.StartSessionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.MeasurementData, $0.MeasurementResult>(
        'StreamMeasurement',
        streamMeasurement,
        true,
        true,
        ($core.List<$core.int> value) => $0.MeasurementData.fromBuffer(value),
        ($0.MeasurementResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.EndSessionRequest, $0.EndSessionResponse>(
        'EndSession',
        endSession_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.EndSessionRequest.fromBuffer(value),
        ($0.EndSessionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetHistoryRequest, $0.GetHistoryResponse>(
        'GetMeasurementHistory',
        getMeasurementHistory_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetHistoryRequest.fromBuffer(value),
        ($0.GetHistoryResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ExportSingleMeasurementRequest,
            $0.ExportFHIRResponse>(
        'ExportSingleMeasurement',
        exportSingleMeasurement_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ExportSingleMeasurementRequest.fromBuffer(value),
        ($0.ExportFHIRResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ExportToFHIRObservationsRequest,
            $0.ExportFHIRResponse>(
        'ExportToFHIRObservations',
        exportToFHIRObservations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ExportToFHIRObservationsRequest.fromBuffer(value),
        ($0.ExportFHIRResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SyncDigitalTwinRequest,
            $0.SyncDigitalTwinResponse>(
        'SyncDigitalTwin',
        syncDigitalTwin_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SyncDigitalTwinRequest.fromBuffer(value),
        ($0.SyncDigitalTwinResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetCalibrationStatusRequest,
            $0.GetCalibrationStatusResponse>(
        'GetCalibrationStatus',
        getCalibrationStatus_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetCalibrationStatusRequest.fromBuffer(value),
        ($0.GetCalibrationStatusResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.StartSessionResponse> startSession_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StartSessionRequest> $request) async {
    return startSession($call, await $request);
  }

  $async.Future<$0.StartSessionResponse> startSession(
      $grpc.ServiceCall call, $0.StartSessionRequest request);

  $async.Stream<$0.MeasurementResult> streamMeasurement(
      $grpc.ServiceCall call, $async.Stream<$0.MeasurementData> request);

  $async.Future<$0.EndSessionResponse> endSession_Pre($grpc.ServiceCall $call,
      $async.Future<$0.EndSessionRequest> $request) async {
    return endSession($call, await $request);
  }

  $async.Future<$0.EndSessionResponse> endSession(
      $grpc.ServiceCall call, $0.EndSessionRequest request);

  $async.Future<$0.GetHistoryResponse> getMeasurementHistory_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetHistoryRequest> $request) async {
    return getMeasurementHistory($call, await $request);
  }

  $async.Future<$0.GetHistoryResponse> getMeasurementHistory(
      $grpc.ServiceCall call, $0.GetHistoryRequest request);

  $async.Future<$0.ExportFHIRResponse> exportSingleMeasurement_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ExportSingleMeasurementRequest> $request) async {
    return exportSingleMeasurement($call, await $request);
  }

  $async.Future<$0.ExportFHIRResponse> exportSingleMeasurement(
      $grpc.ServiceCall call, $0.ExportSingleMeasurementRequest request);

  $async.Future<$0.ExportFHIRResponse> exportToFHIRObservations_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ExportToFHIRObservationsRequest> $request) async {
    return exportToFHIRObservations($call, await $request);
  }

  $async.Future<$0.ExportFHIRResponse> exportToFHIRObservations(
      $grpc.ServiceCall call, $0.ExportToFHIRObservationsRequest request);

  $async.Future<$0.SyncDigitalTwinResponse> syncDigitalTwin_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SyncDigitalTwinRequest> $request) async {
    return syncDigitalTwin($call, await $request);
  }

  $async.Future<$0.SyncDigitalTwinResponse> syncDigitalTwin(
      $grpc.ServiceCall call, $0.SyncDigitalTwinRequest request);

  $async.Future<$0.GetCalibrationStatusResponse> getCalibrationStatus_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetCalibrationStatusRequest> $request) async {
    return getCalibrationStatus($call, await $request);
  }

  $async.Future<$0.GetCalibrationStatusResponse> getCalibrationStatus(
      $grpc.ServiceCall call, $0.GetCalibrationStatusRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.DeviceService')
class DeviceServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  DeviceServiceClient(super.channel, {super.options, super.interceptors});

  /// 디바이스 등록
  $grpc.ResponseFuture<$0.RegisterDeviceResponse> registerDevice(
    $0.RegisterDeviceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$registerDevice, request, options: options);
  }

  /// 디바이스 목록 조회
  $grpc.ResponseFuture<$0.ListDevicesResponse> listDevices(
    $0.ListDevicesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDevices, request, options: options);
  }

  /// 디바이스 상태 업데이트 스트림
  $grpc.ResponseStream<$0.DeviceCommand> streamDeviceStatus(
    $async.Stream<$0.DeviceStatusUpdate> request, {
    $grpc.CallOptions? options,
  }) {
    return $createStreamingCall(_$streamDeviceStatus, request,
        options: options);
  }

  /// 펌웨어 OTA 업데이트
  $grpc.ResponseFuture<$0.OtaResponse> requestOtaUpdate(
    $0.OtaRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$requestOtaUpdate, request, options: options);
  }

  /// Phase 4: Device status management
  $grpc.ResponseFuture<$0.UpdateDeviceStatusResponse> updateDeviceStatus(
    $0.UpdateDeviceStatusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateDeviceStatus, request, options: options);
  }

  // method descriptors

  static final _$registerDevice =
      $grpc.ClientMethod<$0.RegisterDeviceRequest, $0.RegisterDeviceResponse>(
          '/manpasik.v1.DeviceService/RegisterDevice',
          ($0.RegisterDeviceRequest value) => value.writeToBuffer(),
          $0.RegisterDeviceResponse.fromBuffer);
  static final _$listDevices =
      $grpc.ClientMethod<$0.ListDevicesRequest, $0.ListDevicesResponse>(
          '/manpasik.v1.DeviceService/ListDevices',
          ($0.ListDevicesRequest value) => value.writeToBuffer(),
          $0.ListDevicesResponse.fromBuffer);
  static final _$streamDeviceStatus =
      $grpc.ClientMethod<$0.DeviceStatusUpdate, $0.DeviceCommand>(
          '/manpasik.v1.DeviceService/StreamDeviceStatus',
          ($0.DeviceStatusUpdate value) => value.writeToBuffer(),
          $0.DeviceCommand.fromBuffer);
  static final _$requestOtaUpdate =
      $grpc.ClientMethod<$0.OtaRequest, $0.OtaResponse>(
          '/manpasik.v1.DeviceService/RequestOtaUpdate',
          ($0.OtaRequest value) => value.writeToBuffer(),
          $0.OtaResponse.fromBuffer);
  static final _$updateDeviceStatus = $grpc.ClientMethod<
          $0.UpdateDeviceStatusRequest, $0.UpdateDeviceStatusResponse>(
      '/manpasik.v1.DeviceService/UpdateDeviceStatus',
      ($0.UpdateDeviceStatusRequest value) => value.writeToBuffer(),
      $0.UpdateDeviceStatusResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.DeviceService')
abstract class DeviceServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.DeviceService';

  DeviceServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.RegisterDeviceRequest,
            $0.RegisterDeviceResponse>(
        'RegisterDevice',
        registerDevice_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RegisterDeviceRequest.fromBuffer(value),
        ($0.RegisterDeviceResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListDevicesRequest, $0.ListDevicesResponse>(
            'ListDevices',
            listDevices_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListDevicesRequest.fromBuffer(value),
            ($0.ListDevicesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeviceStatusUpdate, $0.DeviceCommand>(
        'StreamDeviceStatus',
        streamDeviceStatus,
        true,
        true,
        ($core.List<$core.int> value) =>
            $0.DeviceStatusUpdate.fromBuffer(value),
        ($0.DeviceCommand value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.OtaRequest, $0.OtaResponse>(
        'RequestOtaUpdate',
        requestOtaUpdate_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.OtaRequest.fromBuffer(value),
        ($0.OtaResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateDeviceStatusRequest,
            $0.UpdateDeviceStatusResponse>(
        'UpdateDeviceStatus',
        updateDeviceStatus_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateDeviceStatusRequest.fromBuffer(value),
        ($0.UpdateDeviceStatusResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.RegisterDeviceResponse> registerDevice_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RegisterDeviceRequest> $request) async {
    return registerDevice($call, await $request);
  }

  $async.Future<$0.RegisterDeviceResponse> registerDevice(
      $grpc.ServiceCall call, $0.RegisterDeviceRequest request);

  $async.Future<$0.ListDevicesResponse> listDevices_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListDevicesRequest> $request) async {
    return listDevices($call, await $request);
  }

  $async.Future<$0.ListDevicesResponse> listDevices(
      $grpc.ServiceCall call, $0.ListDevicesRequest request);

  $async.Stream<$0.DeviceCommand> streamDeviceStatus(
      $grpc.ServiceCall call, $async.Stream<$0.DeviceStatusUpdate> request);

  $async.Future<$0.OtaResponse> requestOtaUpdate_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.OtaRequest> $request) async {
    return requestOtaUpdate($call, await $request);
  }

  $async.Future<$0.OtaResponse> requestOtaUpdate(
      $grpc.ServiceCall call, $0.OtaRequest request);

  $async.Future<$0.UpdateDeviceStatusResponse> updateDeviceStatus_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateDeviceStatusRequest> $request) async {
    return updateDeviceStatus($call, await $request);
  }

  $async.Future<$0.UpdateDeviceStatusResponse> updateDeviceStatus(
      $grpc.ServiceCall call, $0.UpdateDeviceStatusRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.UserService')
class UserServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  UserServiceClient(super.channel, {super.options, super.interceptors});

  /// 사용자 프로필 조회
  $grpc.ResponseFuture<$0.UserProfile> getProfile(
    $0.GetProfileRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getProfile, request, options: options);
  }

  /// 사용자 프로필 업데이트
  $grpc.ResponseFuture<$0.UserProfile> updateProfile(
    $0.UpdateProfileRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateProfile, request, options: options);
  }

  /// 구독 정보 조회
  $grpc.ResponseFuture<$0.SubscriptionInfo> getSubscription(
    $0.GetSubscriptionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSubscription, request, options: options);
  }

  // method descriptors

  static final _$getProfile =
      $grpc.ClientMethod<$0.GetProfileRequest, $0.UserProfile>(
          '/manpasik.v1.UserService/GetProfile',
          ($0.GetProfileRequest value) => value.writeToBuffer(),
          $0.UserProfile.fromBuffer);
  static final _$updateProfile =
      $grpc.ClientMethod<$0.UpdateProfileRequest, $0.UserProfile>(
          '/manpasik.v1.UserService/UpdateProfile',
          ($0.UpdateProfileRequest value) => value.writeToBuffer(),
          $0.UserProfile.fromBuffer);
  static final _$getSubscription =
      $grpc.ClientMethod<$0.GetSubscriptionRequest, $0.SubscriptionInfo>(
          '/manpasik.v1.UserService/GetSubscription',
          ($0.GetSubscriptionRequest value) => value.writeToBuffer(),
          $0.SubscriptionInfo.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.UserService')
abstract class UserServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.UserService';

  UserServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.GetProfileRequest, $0.UserProfile>(
        'GetProfile',
        getProfile_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetProfileRequest.fromBuffer(value),
        ($0.UserProfile value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateProfileRequest, $0.UserProfile>(
        'UpdateProfile',
        updateProfile_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateProfileRequest.fromBuffer(value),
        ($0.UserProfile value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetSubscriptionRequest, $0.SubscriptionInfo>(
            'GetSubscription',
            getSubscription_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetSubscriptionRequest.fromBuffer(value),
            ($0.SubscriptionInfo value) => value.writeToBuffer()));
  }

  $async.Future<$0.UserProfile> getProfile_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetProfileRequest> $request) async {
    return getProfile($call, await $request);
  }

  $async.Future<$0.UserProfile> getProfile(
      $grpc.ServiceCall call, $0.GetProfileRequest request);

  $async.Future<$0.UserProfile> updateProfile_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateProfileRequest> $request) async {
    return updateProfile($call, await $request);
  }

  $async.Future<$0.UserProfile> updateProfile(
      $grpc.ServiceCall call, $0.UpdateProfileRequest request);

  $async.Future<$0.SubscriptionInfo> getSubscription_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetSubscriptionRequest> $request) async {
    return getSubscription($call, await $request);
  }

  $async.Future<$0.SubscriptionInfo> getSubscription(
      $grpc.ServiceCall call, $0.GetSubscriptionRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.SubscriptionService')
class SubscriptionServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  SubscriptionServiceClient(super.channel, {super.options, super.interceptors});

  /// 구독 생성 (회원가입 시 Free 자동 생성)
  $grpc.ResponseFuture<$0.SubscriptionDetail> createSubscription(
    $0.CreateSubscriptionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createSubscription, request, options: options);
  }

  /// 구독 정보 조회
  $grpc.ResponseFuture<$0.SubscriptionDetail> getSubscription(
    $0.GetSubscriptionDetailRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSubscription, request, options: options);
  }

  /// 구독 업데이트 (티어 변경)
  $grpc.ResponseFuture<$0.SubscriptionDetail> updateSubscription(
    $0.UpdateSubscriptionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateSubscription, request, options: options);
  }

  /// 구독 해지
  $grpc.ResponseFuture<$0.CancelSubscriptionResponse> cancelSubscription(
    $0.CancelSubscriptionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$cancelSubscription, request, options: options);
  }

  /// 기능 접근 권한 확인
  $grpc.ResponseFuture<$0.CheckFeatureAccessResponse> checkFeatureAccess(
    $0.CheckFeatureAccessRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$checkFeatureAccess, request, options: options);
  }

  /// 구독 플랜 목록 조회
  $grpc.ResponseFuture<$0.ListSubscriptionPlansResponse> listSubscriptionPlans(
    $0.ListSubscriptionPlansRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listSubscriptionPlans, request, options: options);
  }

  /// 카트리지 접근 권한 확인 (등급별 접근 제어)
  $grpc.ResponseFuture<$0.CheckCartridgeAccessResponse> checkCartridgeAccess(
    $0.CheckCartridgeAccessRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$checkCartridgeAccess, request, options: options);
  }

  /// 사용자별 접근 가능 카트리지 목록
  $grpc.ResponseFuture<$0.ListAccessibleCartridgesResponse>
      listAccessibleCartridges(
    $0.ListAccessibleCartridgesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listAccessibleCartridges, request,
        options: options);
  }

  // method descriptors

  static final _$createSubscription =
      $grpc.ClientMethod<$0.CreateSubscriptionRequest, $0.SubscriptionDetail>(
          '/manpasik.v1.SubscriptionService/CreateSubscription',
          ($0.CreateSubscriptionRequest value) => value.writeToBuffer(),
          $0.SubscriptionDetail.fromBuffer);
  static final _$getSubscription = $grpc.ClientMethod<
          $0.GetSubscriptionDetailRequest, $0.SubscriptionDetail>(
      '/manpasik.v1.SubscriptionService/GetSubscription',
      ($0.GetSubscriptionDetailRequest value) => value.writeToBuffer(),
      $0.SubscriptionDetail.fromBuffer);
  static final _$updateSubscription =
      $grpc.ClientMethod<$0.UpdateSubscriptionRequest, $0.SubscriptionDetail>(
          '/manpasik.v1.SubscriptionService/UpdateSubscription',
          ($0.UpdateSubscriptionRequest value) => value.writeToBuffer(),
          $0.SubscriptionDetail.fromBuffer);
  static final _$cancelSubscription = $grpc.ClientMethod<
          $0.CancelSubscriptionRequest, $0.CancelSubscriptionResponse>(
      '/manpasik.v1.SubscriptionService/CancelSubscription',
      ($0.CancelSubscriptionRequest value) => value.writeToBuffer(),
      $0.CancelSubscriptionResponse.fromBuffer);
  static final _$checkFeatureAccess = $grpc.ClientMethod<
          $0.CheckFeatureAccessRequest, $0.CheckFeatureAccessResponse>(
      '/manpasik.v1.SubscriptionService/CheckFeatureAccess',
      ($0.CheckFeatureAccessRequest value) => value.writeToBuffer(),
      $0.CheckFeatureAccessResponse.fromBuffer);
  static final _$listSubscriptionPlans = $grpc.ClientMethod<
          $0.ListSubscriptionPlansRequest, $0.ListSubscriptionPlansResponse>(
      '/manpasik.v1.SubscriptionService/ListSubscriptionPlans',
      ($0.ListSubscriptionPlansRequest value) => value.writeToBuffer(),
      $0.ListSubscriptionPlansResponse.fromBuffer);
  static final _$checkCartridgeAccess = $grpc.ClientMethod<
          $0.CheckCartridgeAccessRequest, $0.CheckCartridgeAccessResponse>(
      '/manpasik.v1.SubscriptionService/CheckCartridgeAccess',
      ($0.CheckCartridgeAccessRequest value) => value.writeToBuffer(),
      $0.CheckCartridgeAccessResponse.fromBuffer);
  static final _$listAccessibleCartridges = $grpc.ClientMethod<
          $0.ListAccessibleCartridgesRequest,
          $0.ListAccessibleCartridgesResponse>(
      '/manpasik.v1.SubscriptionService/ListAccessibleCartridges',
      ($0.ListAccessibleCartridgesRequest value) => value.writeToBuffer(),
      $0.ListAccessibleCartridgesResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.SubscriptionService')
abstract class SubscriptionServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.SubscriptionService';

  SubscriptionServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.CreateSubscriptionRequest,
            $0.SubscriptionDetail>(
        'CreateSubscription',
        createSubscription_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateSubscriptionRequest.fromBuffer(value),
        ($0.SubscriptionDetail value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetSubscriptionDetailRequest,
            $0.SubscriptionDetail>(
        'GetSubscription',
        getSubscription_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetSubscriptionDetailRequest.fromBuffer(value),
        ($0.SubscriptionDetail value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateSubscriptionRequest,
            $0.SubscriptionDetail>(
        'UpdateSubscription',
        updateSubscription_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateSubscriptionRequest.fromBuffer(value),
        ($0.SubscriptionDetail value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CancelSubscriptionRequest,
            $0.CancelSubscriptionResponse>(
        'CancelSubscription',
        cancelSubscription_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CancelSubscriptionRequest.fromBuffer(value),
        ($0.CancelSubscriptionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CheckFeatureAccessRequest,
            $0.CheckFeatureAccessResponse>(
        'CheckFeatureAccess',
        checkFeatureAccess_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CheckFeatureAccessRequest.fromBuffer(value),
        ($0.CheckFeatureAccessResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListSubscriptionPlansRequest,
            $0.ListSubscriptionPlansResponse>(
        'ListSubscriptionPlans',
        listSubscriptionPlans_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListSubscriptionPlansRequest.fromBuffer(value),
        ($0.ListSubscriptionPlansResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CheckCartridgeAccessRequest,
            $0.CheckCartridgeAccessResponse>(
        'CheckCartridgeAccess',
        checkCartridgeAccess_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CheckCartridgeAccessRequest.fromBuffer(value),
        ($0.CheckCartridgeAccessResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListAccessibleCartridgesRequest,
            $0.ListAccessibleCartridgesResponse>(
        'ListAccessibleCartridges',
        listAccessibleCartridges_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListAccessibleCartridgesRequest.fromBuffer(value),
        ($0.ListAccessibleCartridgesResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.SubscriptionDetail> createSubscription_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateSubscriptionRequest> $request) async {
    return createSubscription($call, await $request);
  }

  $async.Future<$0.SubscriptionDetail> createSubscription(
      $grpc.ServiceCall call, $0.CreateSubscriptionRequest request);

  $async.Future<$0.SubscriptionDetail> getSubscription_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetSubscriptionDetailRequest> $request) async {
    return getSubscription($call, await $request);
  }

  $async.Future<$0.SubscriptionDetail> getSubscription(
      $grpc.ServiceCall call, $0.GetSubscriptionDetailRequest request);

  $async.Future<$0.SubscriptionDetail> updateSubscription_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateSubscriptionRequest> $request) async {
    return updateSubscription($call, await $request);
  }

  $async.Future<$0.SubscriptionDetail> updateSubscription(
      $grpc.ServiceCall call, $0.UpdateSubscriptionRequest request);

  $async.Future<$0.CancelSubscriptionResponse> cancelSubscription_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CancelSubscriptionRequest> $request) async {
    return cancelSubscription($call, await $request);
  }

  $async.Future<$0.CancelSubscriptionResponse> cancelSubscription(
      $grpc.ServiceCall call, $0.CancelSubscriptionRequest request);

  $async.Future<$0.CheckFeatureAccessResponse> checkFeatureAccess_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CheckFeatureAccessRequest> $request) async {
    return checkFeatureAccess($call, await $request);
  }

  $async.Future<$0.CheckFeatureAccessResponse> checkFeatureAccess(
      $grpc.ServiceCall call, $0.CheckFeatureAccessRequest request);

  $async.Future<$0.ListSubscriptionPlansResponse> listSubscriptionPlans_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListSubscriptionPlansRequest> $request) async {
    return listSubscriptionPlans($call, await $request);
  }

  $async.Future<$0.ListSubscriptionPlansResponse> listSubscriptionPlans(
      $grpc.ServiceCall call, $0.ListSubscriptionPlansRequest request);

  $async.Future<$0.CheckCartridgeAccessResponse> checkCartridgeAccess_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CheckCartridgeAccessRequest> $request) async {
    return checkCartridgeAccess($call, await $request);
  }

  $async.Future<$0.CheckCartridgeAccessResponse> checkCartridgeAccess(
      $grpc.ServiceCall call, $0.CheckCartridgeAccessRequest request);

  $async.Future<$0.ListAccessibleCartridgesResponse>
      listAccessibleCartridges_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListAccessibleCartridgesRequest> $request) async {
    return listAccessibleCartridges($call, await $request);
  }

  $async.Future<$0.ListAccessibleCartridgesResponse> listAccessibleCartridges(
      $grpc.ServiceCall call, $0.ListAccessibleCartridgesRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.ShopService')
class ShopServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  ShopServiceClient(super.channel, {super.options, super.interceptors});

  /// 상품 목록 조회
  $grpc.ResponseFuture<$0.ListProductsResponse> listProducts(
    $0.ListProductsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listProducts, request, options: options);
  }

  /// 상품 상세 조회
  $grpc.ResponseFuture<$0.Product> getProduct(
    $0.GetProductRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getProduct, request, options: options);
  }

  /// 장바구니 추가
  $grpc.ResponseFuture<$0.Cart> addToCart(
    $0.AddToCartRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$addToCart, request, options: options);
  }

  /// 장바구니 조회
  $grpc.ResponseFuture<$0.Cart> getCart(
    $0.GetCartRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getCart, request, options: options);
  }

  /// 장바구니 항목 제거
  $grpc.ResponseFuture<$0.Cart> removeFromCart(
    $0.RemoveFromCartRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$removeFromCart, request, options: options);
  }

  /// 주문 생성
  $grpc.ResponseFuture<$0.Order> createOrder(
    $0.CreateOrderRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createOrder, request, options: options);
  }

  /// 주문 상세 조회
  $grpc.ResponseFuture<$0.Order> getOrder(
    $0.GetOrderRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getOrder, request, options: options);
  }

  /// 주문 이력 조회
  $grpc.ResponseFuture<$0.ListOrdersResponse> listOrders(
    $0.ListOrdersRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listOrders, request, options: options);
  }

  // method descriptors

  static final _$listProducts =
      $grpc.ClientMethod<$0.ListProductsRequest, $0.ListProductsResponse>(
          '/manpasik.v1.ShopService/ListProducts',
          ($0.ListProductsRequest value) => value.writeToBuffer(),
          $0.ListProductsResponse.fromBuffer);
  static final _$getProduct =
      $grpc.ClientMethod<$0.GetProductRequest, $0.Product>(
          '/manpasik.v1.ShopService/GetProduct',
          ($0.GetProductRequest value) => value.writeToBuffer(),
          $0.Product.fromBuffer);
  static final _$addToCart = $grpc.ClientMethod<$0.AddToCartRequest, $0.Cart>(
      '/manpasik.v1.ShopService/AddToCart',
      ($0.AddToCartRequest value) => value.writeToBuffer(),
      $0.Cart.fromBuffer);
  static final _$getCart = $grpc.ClientMethod<$0.GetCartRequest, $0.Cart>(
      '/manpasik.v1.ShopService/GetCart',
      ($0.GetCartRequest value) => value.writeToBuffer(),
      $0.Cart.fromBuffer);
  static final _$removeFromCart =
      $grpc.ClientMethod<$0.RemoveFromCartRequest, $0.Cart>(
          '/manpasik.v1.ShopService/RemoveFromCart',
          ($0.RemoveFromCartRequest value) => value.writeToBuffer(),
          $0.Cart.fromBuffer);
  static final _$createOrder =
      $grpc.ClientMethod<$0.CreateOrderRequest, $0.Order>(
          '/manpasik.v1.ShopService/CreateOrder',
          ($0.CreateOrderRequest value) => value.writeToBuffer(),
          $0.Order.fromBuffer);
  static final _$getOrder = $grpc.ClientMethod<$0.GetOrderRequest, $0.Order>(
      '/manpasik.v1.ShopService/GetOrder',
      ($0.GetOrderRequest value) => value.writeToBuffer(),
      $0.Order.fromBuffer);
  static final _$listOrders =
      $grpc.ClientMethod<$0.ListOrdersRequest, $0.ListOrdersResponse>(
          '/manpasik.v1.ShopService/ListOrders',
          ($0.ListOrdersRequest value) => value.writeToBuffer(),
          $0.ListOrdersResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.ShopService')
abstract class ShopServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.ShopService';

  ShopServiceBase() {
    $addMethod(
        $grpc.ServiceMethod<$0.ListProductsRequest, $0.ListProductsResponse>(
            'ListProducts',
            listProducts_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListProductsRequest.fromBuffer(value),
            ($0.ListProductsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetProductRequest, $0.Product>(
        'GetProduct',
        getProduct_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetProductRequest.fromBuffer(value),
        ($0.Product value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AddToCartRequest, $0.Cart>(
        'AddToCart',
        addToCart_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.AddToCartRequest.fromBuffer(value),
        ($0.Cart value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetCartRequest, $0.Cart>(
        'GetCart',
        getCart_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetCartRequest.fromBuffer(value),
        ($0.Cart value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RemoveFromCartRequest, $0.Cart>(
        'RemoveFromCart',
        removeFromCart_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RemoveFromCartRequest.fromBuffer(value),
        ($0.Cart value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateOrderRequest, $0.Order>(
        'CreateOrder',
        createOrder_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateOrderRequest.fromBuffer(value),
        ($0.Order value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetOrderRequest, $0.Order>(
        'GetOrder',
        getOrder_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetOrderRequest.fromBuffer(value),
        ($0.Order value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListOrdersRequest, $0.ListOrdersResponse>(
        'ListOrders',
        listOrders_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListOrdersRequest.fromBuffer(value),
        ($0.ListOrdersResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.ListProductsResponse> listProducts_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListProductsRequest> $request) async {
    return listProducts($call, await $request);
  }

  $async.Future<$0.ListProductsResponse> listProducts(
      $grpc.ServiceCall call, $0.ListProductsRequest request);

  $async.Future<$0.Product> getProduct_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetProductRequest> $request) async {
    return getProduct($call, await $request);
  }

  $async.Future<$0.Product> getProduct(
      $grpc.ServiceCall call, $0.GetProductRequest request);

  $async.Future<$0.Cart> addToCart_Pre($grpc.ServiceCall $call,
      $async.Future<$0.AddToCartRequest> $request) async {
    return addToCart($call, await $request);
  }

  $async.Future<$0.Cart> addToCart(
      $grpc.ServiceCall call, $0.AddToCartRequest request);

  $async.Future<$0.Cart> getCart_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetCartRequest> $request) async {
    return getCart($call, await $request);
  }

  $async.Future<$0.Cart> getCart(
      $grpc.ServiceCall call, $0.GetCartRequest request);

  $async.Future<$0.Cart> removeFromCart_Pre($grpc.ServiceCall $call,
      $async.Future<$0.RemoveFromCartRequest> $request) async {
    return removeFromCart($call, await $request);
  }

  $async.Future<$0.Cart> removeFromCart(
      $grpc.ServiceCall call, $0.RemoveFromCartRequest request);

  $async.Future<$0.Order> createOrder_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateOrderRequest> $request) async {
    return createOrder($call, await $request);
  }

  $async.Future<$0.Order> createOrder(
      $grpc.ServiceCall call, $0.CreateOrderRequest request);

  $async.Future<$0.Order> getOrder_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetOrderRequest> $request) async {
    return getOrder($call, await $request);
  }

  $async.Future<$0.Order> getOrder(
      $grpc.ServiceCall call, $0.GetOrderRequest request);

  $async.Future<$0.ListOrdersResponse> listOrders_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListOrdersRequest> $request) async {
    return listOrders($call, await $request);
  }

  $async.Future<$0.ListOrdersResponse> listOrders(
      $grpc.ServiceCall call, $0.ListOrdersRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.PaymentService')
class PaymentServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  PaymentServiceClient(super.channel, {super.options, super.interceptors});

  /// 결제 요청 생성
  $grpc.ResponseFuture<$0.PaymentDetail> createPayment(
    $0.CreatePaymentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createPayment, request, options: options);
  }

  /// 결제 확인 (PG 콜백)
  $grpc.ResponseFuture<$0.PaymentDetail> confirmPayment(
    $0.ConfirmPaymentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$confirmPayment, request, options: options);
  }

  /// 결제 상세 조회
  $grpc.ResponseFuture<$0.PaymentDetail> getPayment(
    $0.GetPaymentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getPayment, request, options: options);
  }

  /// 결제 이력 조회
  $grpc.ResponseFuture<$0.ListPaymentsResponse> listPayments(
    $0.ListPaymentsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listPayments, request, options: options);
  }

  /// 환불 처리
  $grpc.ResponseFuture<$0.RefundResponse> refundPayment(
    $0.RefundPaymentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$refundPayment, request, options: options);
  }

  // method descriptors

  static final _$createPayment =
      $grpc.ClientMethod<$0.CreatePaymentRequest, $0.PaymentDetail>(
          '/manpasik.v1.PaymentService/CreatePayment',
          ($0.CreatePaymentRequest value) => value.writeToBuffer(),
          $0.PaymentDetail.fromBuffer);
  static final _$confirmPayment =
      $grpc.ClientMethod<$0.ConfirmPaymentRequest, $0.PaymentDetail>(
          '/manpasik.v1.PaymentService/ConfirmPayment',
          ($0.ConfirmPaymentRequest value) => value.writeToBuffer(),
          $0.PaymentDetail.fromBuffer);
  static final _$getPayment =
      $grpc.ClientMethod<$0.GetPaymentRequest, $0.PaymentDetail>(
          '/manpasik.v1.PaymentService/GetPayment',
          ($0.GetPaymentRequest value) => value.writeToBuffer(),
          $0.PaymentDetail.fromBuffer);
  static final _$listPayments =
      $grpc.ClientMethod<$0.ListPaymentsRequest, $0.ListPaymentsResponse>(
          '/manpasik.v1.PaymentService/ListPayments',
          ($0.ListPaymentsRequest value) => value.writeToBuffer(),
          $0.ListPaymentsResponse.fromBuffer);
  static final _$refundPayment =
      $grpc.ClientMethod<$0.RefundPaymentRequest, $0.RefundResponse>(
          '/manpasik.v1.PaymentService/RefundPayment',
          ($0.RefundPaymentRequest value) => value.writeToBuffer(),
          $0.RefundResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.PaymentService')
abstract class PaymentServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.PaymentService';

  PaymentServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.CreatePaymentRequest, $0.PaymentDetail>(
        'CreatePayment',
        createPayment_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreatePaymentRequest.fromBuffer(value),
        ($0.PaymentDetail value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ConfirmPaymentRequest, $0.PaymentDetail>(
        'ConfirmPayment',
        confirmPayment_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ConfirmPaymentRequest.fromBuffer(value),
        ($0.PaymentDetail value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetPaymentRequest, $0.PaymentDetail>(
        'GetPayment',
        getPayment_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetPaymentRequest.fromBuffer(value),
        ($0.PaymentDetail value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListPaymentsRequest, $0.ListPaymentsResponse>(
            'ListPayments',
            listPayments_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListPaymentsRequest.fromBuffer(value),
            ($0.ListPaymentsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RefundPaymentRequest, $0.RefundResponse>(
        'RefundPayment',
        refundPayment_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RefundPaymentRequest.fromBuffer(value),
        ($0.RefundResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.PaymentDetail> createPayment_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreatePaymentRequest> $request) async {
    return createPayment($call, await $request);
  }

  $async.Future<$0.PaymentDetail> createPayment(
      $grpc.ServiceCall call, $0.CreatePaymentRequest request);

  $async.Future<$0.PaymentDetail> confirmPayment_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ConfirmPaymentRequest> $request) async {
    return confirmPayment($call, await $request);
  }

  $async.Future<$0.PaymentDetail> confirmPayment(
      $grpc.ServiceCall call, $0.ConfirmPaymentRequest request);

  $async.Future<$0.PaymentDetail> getPayment_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetPaymentRequest> $request) async {
    return getPayment($call, await $request);
  }

  $async.Future<$0.PaymentDetail> getPayment(
      $grpc.ServiceCall call, $0.GetPaymentRequest request);

  $async.Future<$0.ListPaymentsResponse> listPayments_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListPaymentsRequest> $request) async {
    return listPayments($call, await $request);
  }

  $async.Future<$0.ListPaymentsResponse> listPayments(
      $grpc.ServiceCall call, $0.ListPaymentsRequest request);

  $async.Future<$0.RefundResponse> refundPayment_Pre($grpc.ServiceCall $call,
      $async.Future<$0.RefundPaymentRequest> $request) async {
    return refundPayment($call, await $request);
  }

  $async.Future<$0.RefundResponse> refundPayment(
      $grpc.ServiceCall call, $0.RefundPaymentRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.AiInferenceService')
class AiInferenceServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  AiInferenceServiceClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.AnalysisResult> analyzeMeasurement(
    $0.AnalyzeMeasurementRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$analyzeMeasurement, request, options: options);
  }

  $grpc.ResponseFuture<$0.HealthScoreResponse> getHealthScore(
    $0.GetHealthScoreRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getHealthScore, request, options: options);
  }

  $grpc.ResponseFuture<$0.TrendPrediction> predictTrend(
    $0.PredictTrendRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$predictTrend, request, options: options);
  }

  $grpc.ResponseFuture<$0.ModelInfo> getModelInfo(
    $0.GetModelInfoRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getModelInfo, request, options: options);
  }

  $grpc.ResponseFuture<$0.ListModelsResponse> listModels(
    $0.ListModelsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listModels, request, options: options);
  }

  /// Phase 9: AI 스트리밍 채팅 (C1)
  $grpc.ResponseStream<$0.StreamChatResponse> streamChat(
    $0.StreamChatRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createStreamingCall(
        _$streamChat, $async.Stream.fromIterable([request]),
        options: options);
  }

  // method descriptors

  static final _$analyzeMeasurement =
      $grpc.ClientMethod<$0.AnalyzeMeasurementRequest, $0.AnalysisResult>(
          '/manpasik.v1.AiInferenceService/AnalyzeMeasurement',
          ($0.AnalyzeMeasurementRequest value) => value.writeToBuffer(),
          $0.AnalysisResult.fromBuffer);
  static final _$getHealthScore =
      $grpc.ClientMethod<$0.GetHealthScoreRequest, $0.HealthScoreResponse>(
          '/manpasik.v1.AiInferenceService/GetHealthScore',
          ($0.GetHealthScoreRequest value) => value.writeToBuffer(),
          $0.HealthScoreResponse.fromBuffer);
  static final _$predictTrend =
      $grpc.ClientMethod<$0.PredictTrendRequest, $0.TrendPrediction>(
          '/manpasik.v1.AiInferenceService/PredictTrend',
          ($0.PredictTrendRequest value) => value.writeToBuffer(),
          $0.TrendPrediction.fromBuffer);
  static final _$getModelInfo =
      $grpc.ClientMethod<$0.GetModelInfoRequest, $0.ModelInfo>(
          '/manpasik.v1.AiInferenceService/GetModelInfo',
          ($0.GetModelInfoRequest value) => value.writeToBuffer(),
          $0.ModelInfo.fromBuffer);
  static final _$listModels =
      $grpc.ClientMethod<$0.ListModelsRequest, $0.ListModelsResponse>(
          '/manpasik.v1.AiInferenceService/ListModels',
          ($0.ListModelsRequest value) => value.writeToBuffer(),
          $0.ListModelsResponse.fromBuffer);
  static final _$streamChat =
      $grpc.ClientMethod<$0.StreamChatRequest, $0.StreamChatResponse>(
          '/manpasik.v1.AiInferenceService/StreamChat',
          ($0.StreamChatRequest value) => value.writeToBuffer(),
          $0.StreamChatResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.AiInferenceService')
abstract class AiInferenceServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.AiInferenceService';

  AiInferenceServiceBase() {
    $addMethod(
        $grpc.ServiceMethod<$0.AnalyzeMeasurementRequest, $0.AnalysisResult>(
            'AnalyzeMeasurement',
            analyzeMeasurement_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.AnalyzeMeasurementRequest.fromBuffer(value),
            ($0.AnalysisResult value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetHealthScoreRequest, $0.HealthScoreResponse>(
            'GetHealthScore',
            getHealthScore_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetHealthScoreRequest.fromBuffer(value),
            ($0.HealthScoreResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PredictTrendRequest, $0.TrendPrediction>(
        'PredictTrend',
        predictTrend_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PredictTrendRequest.fromBuffer(value),
        ($0.TrendPrediction value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetModelInfoRequest, $0.ModelInfo>(
        'GetModelInfo',
        getModelInfo_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetModelInfoRequest.fromBuffer(value),
        ($0.ModelInfo value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListModelsRequest, $0.ListModelsResponse>(
        'ListModels',
        listModels_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListModelsRequest.fromBuffer(value),
        ($0.ListModelsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StreamChatRequest, $0.StreamChatResponse>(
        'StreamChat',
        streamChat_Pre,
        false,
        true,
        ($core.List<$core.int> value) => $0.StreamChatRequest.fromBuffer(value),
        ($0.StreamChatResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.AnalysisResult> analyzeMeasurement_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AnalyzeMeasurementRequest> $request) async {
    return analyzeMeasurement($call, await $request);
  }

  $async.Future<$0.AnalysisResult> analyzeMeasurement(
      $grpc.ServiceCall call, $0.AnalyzeMeasurementRequest request);

  $async.Future<$0.HealthScoreResponse> getHealthScore_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetHealthScoreRequest> $request) async {
    return getHealthScore($call, await $request);
  }

  $async.Future<$0.HealthScoreResponse> getHealthScore(
      $grpc.ServiceCall call, $0.GetHealthScoreRequest request);

  $async.Future<$0.TrendPrediction> predictTrend_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PredictTrendRequest> $request) async {
    return predictTrend($call, await $request);
  }

  $async.Future<$0.TrendPrediction> predictTrend(
      $grpc.ServiceCall call, $0.PredictTrendRequest request);

  $async.Future<$0.ModelInfo> getModelInfo_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetModelInfoRequest> $request) async {
    return getModelInfo($call, await $request);
  }

  $async.Future<$0.ModelInfo> getModelInfo(
      $grpc.ServiceCall call, $0.GetModelInfoRequest request);

  $async.Future<$0.ListModelsResponse> listModels_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListModelsRequest> $request) async {
    return listModels($call, await $request);
  }

  $async.Future<$0.ListModelsResponse> listModels(
      $grpc.ServiceCall call, $0.ListModelsRequest request);

  $async.Stream<$0.StreamChatResponse> streamChat_Pre($grpc.ServiceCall $call,
      $async.Future<$0.StreamChatRequest> $request) async* {
    yield* streamChat($call, await $request);
  }

  $async.Stream<$0.StreamChatResponse> streamChat(
      $grpc.ServiceCall call, $0.StreamChatRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.CartridgeService')
class CartridgeServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  CartridgeServiceClient(super.channel, {super.options, super.interceptors});

  /// NFC 태그 읽기 → 카트리지 정보 반환
  $grpc.ResponseFuture<$0.CartridgeDetail> readCartridge(
    $0.ReadCartridgeRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$readCartridge, request, options: options);
  }

  /// 카트리지 사용 기록
  $grpc.ResponseFuture<$0.RecordUsageResponse> recordUsage(
    $0.RecordUsageRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$recordUsage, request, options: options);
  }

  /// 카트리지 사용 이력 조회
  $grpc.ResponseFuture<$0.GetUsageHistoryResponse> getUsageHistory(
    $0.GetUsageHistoryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getUsageHistory, request, options: options);
  }

  /// 카트리지 타입 정보 조회 (레지스트리)
  $grpc.ResponseFuture<$0.CartridgeTypeInfo> getCartridgeType(
    $0.GetCartridgeTypeRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getCartridgeType, request, options: options);
  }

  /// 카테고리 목록 조회
  $grpc.ResponseFuture<$0.ListCategoriesResponse> listCategories(
    $0.ListCategoriesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listCategories, request, options: options);
  }

  /// 카테고리별 타입 목록 조회
  $grpc.ResponseFuture<$0.ListTypesByCategoryResponse> listTypesByCategory(
    $0.ListTypesByCategoryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTypesByCategory, request, options: options);
  }

  /// 카트리지 잔여 사용 횟수 조회
  $grpc.ResponseFuture<$0.GetRemainingUsesResponse> getRemainingUses(
    $0.GetRemainingUsesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getRemainingUses, request, options: options);
  }

  /// 카트리지 유효성 검증 (NFC UID + 유효기간 + 잔여횟수)
  $grpc.ResponseFuture<$0.ValidateCartridgeResponse> validateCartridge(
    $0.ValidateCartridgeRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$validateCartridge, request, options: options);
  }

  // method descriptors

  static final _$readCartridge =
      $grpc.ClientMethod<$0.ReadCartridgeRequest, $0.CartridgeDetail>(
          '/manpasik.v1.CartridgeService/ReadCartridge',
          ($0.ReadCartridgeRequest value) => value.writeToBuffer(),
          $0.CartridgeDetail.fromBuffer);
  static final _$recordUsage =
      $grpc.ClientMethod<$0.RecordUsageRequest, $0.RecordUsageResponse>(
          '/manpasik.v1.CartridgeService/RecordUsage',
          ($0.RecordUsageRequest value) => value.writeToBuffer(),
          $0.RecordUsageResponse.fromBuffer);
  static final _$getUsageHistory =
      $grpc.ClientMethod<$0.GetUsageHistoryRequest, $0.GetUsageHistoryResponse>(
          '/manpasik.v1.CartridgeService/GetUsageHistory',
          ($0.GetUsageHistoryRequest value) => value.writeToBuffer(),
          $0.GetUsageHistoryResponse.fromBuffer);
  static final _$getCartridgeType =
      $grpc.ClientMethod<$0.GetCartridgeTypeRequest, $0.CartridgeTypeInfo>(
          '/manpasik.v1.CartridgeService/GetCartridgeType',
          ($0.GetCartridgeTypeRequest value) => value.writeToBuffer(),
          $0.CartridgeTypeInfo.fromBuffer);
  static final _$listCategories =
      $grpc.ClientMethod<$0.ListCategoriesRequest, $0.ListCategoriesResponse>(
          '/manpasik.v1.CartridgeService/ListCategories',
          ($0.ListCategoriesRequest value) => value.writeToBuffer(),
          $0.ListCategoriesResponse.fromBuffer);
  static final _$listTypesByCategory = $grpc.ClientMethod<
          $0.ListTypesByCategoryRequest, $0.ListTypesByCategoryResponse>(
      '/manpasik.v1.CartridgeService/ListTypesByCategory',
      ($0.ListTypesByCategoryRequest value) => value.writeToBuffer(),
      $0.ListTypesByCategoryResponse.fromBuffer);
  static final _$getRemainingUses = $grpc.ClientMethod<
          $0.GetRemainingUsesRequest, $0.GetRemainingUsesResponse>(
      '/manpasik.v1.CartridgeService/GetRemainingUses',
      ($0.GetRemainingUsesRequest value) => value.writeToBuffer(),
      $0.GetRemainingUsesResponse.fromBuffer);
  static final _$validateCartridge = $grpc.ClientMethod<
          $0.ValidateCartridgeRequest, $0.ValidateCartridgeResponse>(
      '/manpasik.v1.CartridgeService/ValidateCartridge',
      ($0.ValidateCartridgeRequest value) => value.writeToBuffer(),
      $0.ValidateCartridgeResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.CartridgeService')
abstract class CartridgeServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.CartridgeService';

  CartridgeServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.ReadCartridgeRequest, $0.CartridgeDetail>(
        'ReadCartridge',
        readCartridge_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ReadCartridgeRequest.fromBuffer(value),
        ($0.CartridgeDetail value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.RecordUsageRequest, $0.RecordUsageResponse>(
            'RecordUsage',
            recordUsage_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.RecordUsageRequest.fromBuffer(value),
            ($0.RecordUsageResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetUsageHistoryRequest,
            $0.GetUsageHistoryResponse>(
        'GetUsageHistory',
        getUsageHistory_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetUsageHistoryRequest.fromBuffer(value),
        ($0.GetUsageHistoryResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetCartridgeTypeRequest, $0.CartridgeTypeInfo>(
            'GetCartridgeType',
            getCartridgeType_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetCartridgeTypeRequest.fromBuffer(value),
            ($0.CartridgeTypeInfo value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListCategoriesRequest,
            $0.ListCategoriesResponse>(
        'ListCategories',
        listCategories_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListCategoriesRequest.fromBuffer(value),
        ($0.ListCategoriesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTypesByCategoryRequest,
            $0.ListTypesByCategoryResponse>(
        'ListTypesByCategory',
        listTypesByCategory_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTypesByCategoryRequest.fromBuffer(value),
        ($0.ListTypesByCategoryResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetRemainingUsesRequest,
            $0.GetRemainingUsesResponse>(
        'GetRemainingUses',
        getRemainingUses_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetRemainingUsesRequest.fromBuffer(value),
        ($0.GetRemainingUsesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ValidateCartridgeRequest,
            $0.ValidateCartridgeResponse>(
        'ValidateCartridge',
        validateCartridge_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ValidateCartridgeRequest.fromBuffer(value),
        ($0.ValidateCartridgeResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.CartridgeDetail> readCartridge_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ReadCartridgeRequest> $request) async {
    return readCartridge($call, await $request);
  }

  $async.Future<$0.CartridgeDetail> readCartridge(
      $grpc.ServiceCall call, $0.ReadCartridgeRequest request);

  $async.Future<$0.RecordUsageResponse> recordUsage_Pre($grpc.ServiceCall $call,
      $async.Future<$0.RecordUsageRequest> $request) async {
    return recordUsage($call, await $request);
  }

  $async.Future<$0.RecordUsageResponse> recordUsage(
      $grpc.ServiceCall call, $0.RecordUsageRequest request);

  $async.Future<$0.GetUsageHistoryResponse> getUsageHistory_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetUsageHistoryRequest> $request) async {
    return getUsageHistory($call, await $request);
  }

  $async.Future<$0.GetUsageHistoryResponse> getUsageHistory(
      $grpc.ServiceCall call, $0.GetUsageHistoryRequest request);

  $async.Future<$0.CartridgeTypeInfo> getCartridgeType_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetCartridgeTypeRequest> $request) async {
    return getCartridgeType($call, await $request);
  }

  $async.Future<$0.CartridgeTypeInfo> getCartridgeType(
      $grpc.ServiceCall call, $0.GetCartridgeTypeRequest request);

  $async.Future<$0.ListCategoriesResponse> listCategories_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListCategoriesRequest> $request) async {
    return listCategories($call, await $request);
  }

  $async.Future<$0.ListCategoriesResponse> listCategories(
      $grpc.ServiceCall call, $0.ListCategoriesRequest request);

  $async.Future<$0.ListTypesByCategoryResponse> listTypesByCategory_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTypesByCategoryRequest> $request) async {
    return listTypesByCategory($call, await $request);
  }

  $async.Future<$0.ListTypesByCategoryResponse> listTypesByCategory(
      $grpc.ServiceCall call, $0.ListTypesByCategoryRequest request);

  $async.Future<$0.GetRemainingUsesResponse> getRemainingUses_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetRemainingUsesRequest> $request) async {
    return getRemainingUses($call, await $request);
  }

  $async.Future<$0.GetRemainingUsesResponse> getRemainingUses(
      $grpc.ServiceCall call, $0.GetRemainingUsesRequest request);

  $async.Future<$0.ValidateCartridgeResponse> validateCartridge_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ValidateCartridgeRequest> $request) async {
    return validateCartridge($call, await $request);
  }

  $async.Future<$0.ValidateCartridgeResponse> validateCartridge(
      $grpc.ServiceCall call, $0.ValidateCartridgeRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.CalibrationService')
class CalibrationServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  CalibrationServiceClient(super.channel, {super.options, super.interceptors});

  /// 팩토리 보정 데이터 등록
  $grpc.ResponseFuture<$0.CalibrationRecord> registerFactoryCalibration(
    $0.RegisterFactoryCalibrationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$registerFactoryCalibration, request,
        options: options);
  }

  /// 현장 보정 수행 (사용자 보정)
  $grpc.ResponseFuture<$0.CalibrationRecord> performFieldCalibration(
    $0.PerformFieldCalibrationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$performFieldCalibration, request,
        options: options);
  }

  /// 보정 데이터 조회 (디바이스 + 카트리지 타입별 최신)
  $grpc.ResponseFuture<$0.CalibrationRecord> getCalibration(
    $0.GetCalibrationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getCalibration, request, options: options);
  }

  /// 보정 이력 조회
  $grpc.ResponseFuture<$0.ListCalibrationHistoryResponse>
      listCalibrationHistory(
    $0.ListCalibrationHistoryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listCalibrationHistory, request,
        options: options);
  }

  /// 보정 상태 확인 (보정 필요 여부 판단)
  $grpc.ResponseFuture<$0.CalibrationStatusResponse> checkCalibrationStatus(
    $0.CheckCalibrationStatusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$checkCalibrationStatus, request,
        options: options);
  }

  /// 보정 모델 목록 조회
  $grpc.ResponseFuture<$0.ListCalibrationModelsResponse> listCalibrationModels(
    $0.ListCalibrationModelsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listCalibrationModels, request, options: options);
  }

  // method descriptors

  static final _$registerFactoryCalibration = $grpc.ClientMethod<
          $0.RegisterFactoryCalibrationRequest, $0.CalibrationRecord>(
      '/manpasik.v1.CalibrationService/RegisterFactoryCalibration',
      ($0.RegisterFactoryCalibrationRequest value) => value.writeToBuffer(),
      $0.CalibrationRecord.fromBuffer);
  static final _$performFieldCalibration = $grpc.ClientMethod<
          $0.PerformFieldCalibrationRequest, $0.CalibrationRecord>(
      '/manpasik.v1.CalibrationService/PerformFieldCalibration',
      ($0.PerformFieldCalibrationRequest value) => value.writeToBuffer(),
      $0.CalibrationRecord.fromBuffer);
  static final _$getCalibration =
      $grpc.ClientMethod<$0.GetCalibrationRequest, $0.CalibrationRecord>(
          '/manpasik.v1.CalibrationService/GetCalibration',
          ($0.GetCalibrationRequest value) => value.writeToBuffer(),
          $0.CalibrationRecord.fromBuffer);
  static final _$listCalibrationHistory = $grpc.ClientMethod<
          $0.ListCalibrationHistoryRequest, $0.ListCalibrationHistoryResponse>(
      '/manpasik.v1.CalibrationService/ListCalibrationHistory',
      ($0.ListCalibrationHistoryRequest value) => value.writeToBuffer(),
      $0.ListCalibrationHistoryResponse.fromBuffer);
  static final _$checkCalibrationStatus = $grpc.ClientMethod<
          $0.CheckCalibrationStatusRequest, $0.CalibrationStatusResponse>(
      '/manpasik.v1.CalibrationService/CheckCalibrationStatus',
      ($0.CheckCalibrationStatusRequest value) => value.writeToBuffer(),
      $0.CalibrationStatusResponse.fromBuffer);
  static final _$listCalibrationModels = $grpc.ClientMethod<
          $0.ListCalibrationModelsRequest, $0.ListCalibrationModelsResponse>(
      '/manpasik.v1.CalibrationService/ListCalibrationModels',
      ($0.ListCalibrationModelsRequest value) => value.writeToBuffer(),
      $0.ListCalibrationModelsResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.CalibrationService')
abstract class CalibrationServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.CalibrationService';

  CalibrationServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.RegisterFactoryCalibrationRequest,
            $0.CalibrationRecord>(
        'RegisterFactoryCalibration',
        registerFactoryCalibration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RegisterFactoryCalibrationRequest.fromBuffer(value),
        ($0.CalibrationRecord value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PerformFieldCalibrationRequest,
            $0.CalibrationRecord>(
        'PerformFieldCalibration',
        performFieldCalibration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PerformFieldCalibrationRequest.fromBuffer(value),
        ($0.CalibrationRecord value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetCalibrationRequest, $0.CalibrationRecord>(
            'GetCalibration',
            getCalibration_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetCalibrationRequest.fromBuffer(value),
            ($0.CalibrationRecord value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListCalibrationHistoryRequest,
            $0.ListCalibrationHistoryResponse>(
        'ListCalibrationHistory',
        listCalibrationHistory_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListCalibrationHistoryRequest.fromBuffer(value),
        ($0.ListCalibrationHistoryResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CheckCalibrationStatusRequest,
            $0.CalibrationStatusResponse>(
        'CheckCalibrationStatus',
        checkCalibrationStatus_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CheckCalibrationStatusRequest.fromBuffer(value),
        ($0.CalibrationStatusResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListCalibrationModelsRequest,
            $0.ListCalibrationModelsResponse>(
        'ListCalibrationModels',
        listCalibrationModels_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListCalibrationModelsRequest.fromBuffer(value),
        ($0.ListCalibrationModelsResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.CalibrationRecord> registerFactoryCalibration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RegisterFactoryCalibrationRequest> $request) async {
    return registerFactoryCalibration($call, await $request);
  }

  $async.Future<$0.CalibrationRecord> registerFactoryCalibration(
      $grpc.ServiceCall call, $0.RegisterFactoryCalibrationRequest request);

  $async.Future<$0.CalibrationRecord> performFieldCalibration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PerformFieldCalibrationRequest> $request) async {
    return performFieldCalibration($call, await $request);
  }

  $async.Future<$0.CalibrationRecord> performFieldCalibration(
      $grpc.ServiceCall call, $0.PerformFieldCalibrationRequest request);

  $async.Future<$0.CalibrationRecord> getCalibration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetCalibrationRequest> $request) async {
    return getCalibration($call, await $request);
  }

  $async.Future<$0.CalibrationRecord> getCalibration(
      $grpc.ServiceCall call, $0.GetCalibrationRequest request);

  $async.Future<$0.ListCalibrationHistoryResponse> listCalibrationHistory_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListCalibrationHistoryRequest> $request) async {
    return listCalibrationHistory($call, await $request);
  }

  $async.Future<$0.ListCalibrationHistoryResponse> listCalibrationHistory(
      $grpc.ServiceCall call, $0.ListCalibrationHistoryRequest request);

  $async.Future<$0.CalibrationStatusResponse> checkCalibrationStatus_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CheckCalibrationStatusRequest> $request) async {
    return checkCalibrationStatus($call, await $request);
  }

  $async.Future<$0.CalibrationStatusResponse> checkCalibrationStatus(
      $grpc.ServiceCall call, $0.CheckCalibrationStatusRequest request);

  $async.Future<$0.ListCalibrationModelsResponse> listCalibrationModels_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListCalibrationModelsRequest> $request) async {
    return listCalibrationModels($call, await $request);
  }

  $async.Future<$0.ListCalibrationModelsResponse> listCalibrationModels(
      $grpc.ServiceCall call, $0.ListCalibrationModelsRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.CoachingService')
class CoachingServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  CoachingServiceClient(super.channel, {super.options, super.interceptors});

  /// 건강 목표 설정
  $grpc.ResponseFuture<$0.HealthGoal> setHealthGoal(
    $0.SetHealthGoalRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$setHealthGoal, request, options: options);
  }

  /// 건강 목표 조회
  $grpc.ResponseFuture<$0.GetHealthGoalsResponse> getHealthGoals(
    $0.GetHealthGoalsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getHealthGoals, request, options: options);
  }

  /// AI 코칭 메시지 생성 (측정 결과 기반)
  $grpc.ResponseFuture<$0.CoachingMessage> generateCoaching(
    $0.GenerateCoachingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$generateCoaching, request, options: options);
  }

  /// 코칭 메시지 이력 조회
  $grpc.ResponseFuture<$0.ListCoachingMessagesResponse> listCoachingMessages(
    $0.ListCoachingMessagesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listCoachingMessages, request, options: options);
  }

  /// 일일 건강 리포트 생성
  $grpc.ResponseFuture<$0.DailyHealthReport> generateDailyReport(
    $0.GenerateDailyReportRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$generateDailyReport, request, options: options);
  }

  /// 주간 건강 리포트 조회
  $grpc.ResponseFuture<$0.WeeklyHealthReport> getWeeklyReport(
    $0.GetWeeklyReportRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getWeeklyReport, request, options: options);
  }

  /// 개인화 추천 조회
  $grpc.ResponseFuture<$0.GetRecommendationsResponse> getRecommendations(
    $0.GetRecommendationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getRecommendations, request, options: options);
  }

  // method descriptors

  static final _$setHealthGoal =
      $grpc.ClientMethod<$0.SetHealthGoalRequest, $0.HealthGoal>(
          '/manpasik.v1.CoachingService/SetHealthGoal',
          ($0.SetHealthGoalRequest value) => value.writeToBuffer(),
          $0.HealthGoal.fromBuffer);
  static final _$getHealthGoals =
      $grpc.ClientMethod<$0.GetHealthGoalsRequest, $0.GetHealthGoalsResponse>(
          '/manpasik.v1.CoachingService/GetHealthGoals',
          ($0.GetHealthGoalsRequest value) => value.writeToBuffer(),
          $0.GetHealthGoalsResponse.fromBuffer);
  static final _$generateCoaching =
      $grpc.ClientMethod<$0.GenerateCoachingRequest, $0.CoachingMessage>(
          '/manpasik.v1.CoachingService/GenerateCoaching',
          ($0.GenerateCoachingRequest value) => value.writeToBuffer(),
          $0.CoachingMessage.fromBuffer);
  static final _$listCoachingMessages = $grpc.ClientMethod<
          $0.ListCoachingMessagesRequest, $0.ListCoachingMessagesResponse>(
      '/manpasik.v1.CoachingService/ListCoachingMessages',
      ($0.ListCoachingMessagesRequest value) => value.writeToBuffer(),
      $0.ListCoachingMessagesResponse.fromBuffer);
  static final _$generateDailyReport =
      $grpc.ClientMethod<$0.GenerateDailyReportRequest, $0.DailyHealthReport>(
          '/manpasik.v1.CoachingService/GenerateDailyReport',
          ($0.GenerateDailyReportRequest value) => value.writeToBuffer(),
          $0.DailyHealthReport.fromBuffer);
  static final _$getWeeklyReport =
      $grpc.ClientMethod<$0.GetWeeklyReportRequest, $0.WeeklyHealthReport>(
          '/manpasik.v1.CoachingService/GetWeeklyReport',
          ($0.GetWeeklyReportRequest value) => value.writeToBuffer(),
          $0.WeeklyHealthReport.fromBuffer);
  static final _$getRecommendations = $grpc.ClientMethod<
          $0.GetRecommendationsRequest, $0.GetRecommendationsResponse>(
      '/manpasik.v1.CoachingService/GetRecommendations',
      ($0.GetRecommendationsRequest value) => value.writeToBuffer(),
      $0.GetRecommendationsResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.CoachingService')
abstract class CoachingServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.CoachingService';

  CoachingServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.SetHealthGoalRequest, $0.HealthGoal>(
        'SetHealthGoal',
        setHealthGoal_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SetHealthGoalRequest.fromBuffer(value),
        ($0.HealthGoal value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetHealthGoalsRequest,
            $0.GetHealthGoalsResponse>(
        'GetHealthGoals',
        getHealthGoals_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetHealthGoalsRequest.fromBuffer(value),
        ($0.GetHealthGoalsResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GenerateCoachingRequest, $0.CoachingMessage>(
            'GenerateCoaching',
            generateCoaching_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GenerateCoachingRequest.fromBuffer(value),
            ($0.CoachingMessage value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListCoachingMessagesRequest,
            $0.ListCoachingMessagesResponse>(
        'ListCoachingMessages',
        listCoachingMessages_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListCoachingMessagesRequest.fromBuffer(value),
        ($0.ListCoachingMessagesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GenerateDailyReportRequest,
            $0.DailyHealthReport>(
        'GenerateDailyReport',
        generateDailyReport_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GenerateDailyReportRequest.fromBuffer(value),
        ($0.DailyHealthReport value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetWeeklyReportRequest, $0.WeeklyHealthReport>(
            'GetWeeklyReport',
            getWeeklyReport_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetWeeklyReportRequest.fromBuffer(value),
            ($0.WeeklyHealthReport value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetRecommendationsRequest,
            $0.GetRecommendationsResponse>(
        'GetRecommendations',
        getRecommendations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetRecommendationsRequest.fromBuffer(value),
        ($0.GetRecommendationsResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.HealthGoal> setHealthGoal_Pre($grpc.ServiceCall $call,
      $async.Future<$0.SetHealthGoalRequest> $request) async {
    return setHealthGoal($call, await $request);
  }

  $async.Future<$0.HealthGoal> setHealthGoal(
      $grpc.ServiceCall call, $0.SetHealthGoalRequest request);

  $async.Future<$0.GetHealthGoalsResponse> getHealthGoals_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetHealthGoalsRequest> $request) async {
    return getHealthGoals($call, await $request);
  }

  $async.Future<$0.GetHealthGoalsResponse> getHealthGoals(
      $grpc.ServiceCall call, $0.GetHealthGoalsRequest request);

  $async.Future<$0.CoachingMessage> generateCoaching_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GenerateCoachingRequest> $request) async {
    return generateCoaching($call, await $request);
  }

  $async.Future<$0.CoachingMessage> generateCoaching(
      $grpc.ServiceCall call, $0.GenerateCoachingRequest request);

  $async.Future<$0.ListCoachingMessagesResponse> listCoachingMessages_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListCoachingMessagesRequest> $request) async {
    return listCoachingMessages($call, await $request);
  }

  $async.Future<$0.ListCoachingMessagesResponse> listCoachingMessages(
      $grpc.ServiceCall call, $0.ListCoachingMessagesRequest request);

  $async.Future<$0.DailyHealthReport> generateDailyReport_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GenerateDailyReportRequest> $request) async {
    return generateDailyReport($call, await $request);
  }

  $async.Future<$0.DailyHealthReport> generateDailyReport(
      $grpc.ServiceCall call, $0.GenerateDailyReportRequest request);

  $async.Future<$0.WeeklyHealthReport> getWeeklyReport_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetWeeklyReportRequest> $request) async {
    return getWeeklyReport($call, await $request);
  }

  $async.Future<$0.WeeklyHealthReport> getWeeklyReport(
      $grpc.ServiceCall call, $0.GetWeeklyReportRequest request);

  $async.Future<$0.GetRecommendationsResponse> getRecommendations_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetRecommendationsRequest> $request) async {
    return getRecommendations($call, await $request);
  }

  $async.Future<$0.GetRecommendationsResponse> getRecommendations(
      $grpc.ServiceCall call, $0.GetRecommendationsRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.ReservationService')
class ReservationServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  ReservationServiceClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.SearchFacilitiesResponse> searchFacilities(
    $0.SearchFacilitiesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$searchFacilities, request, options: options);
  }

  $grpc.ResponseFuture<$0.Facility> getFacility(
    $0.GetFacilityRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getFacility, request, options: options);
  }

  $grpc.ResponseFuture<$0.GetAvailableSlotsResponse> getAvailableSlots(
    $0.GetAvailableSlotsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getAvailableSlots, request, options: options);
  }

  $grpc.ResponseFuture<$0.Reservation> createReservation(
    $0.CreateReservationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createReservation, request, options: options);
  }

  $grpc.ResponseFuture<$0.Reservation> getReservation(
    $0.GetReservationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getReservation, request, options: options);
  }

  $grpc.ResponseFuture<$0.ListReservationsResponse> listReservations(
    $0.ListReservationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listReservations, request, options: options);
  }

  $grpc.ResponseFuture<$0.CancelReservationResponse> cancelReservation(
    $0.CancelReservationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$cancelReservation, request, options: options);
  }

  /// Phase 4: Doctor & region extensions
  $grpc.ResponseFuture<$0.ListDoctorsByFacilityResponse> listDoctorsByFacility(
    $0.ListDoctorsByFacilityRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDoctorsByFacility, request, options: options);
  }

  $grpc.ResponseFuture<$0.GetDoctorAvailabilityResponse> getDoctorAvailability(
    $0.GetDoctorAvailabilityRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDoctorAvailability, request, options: options);
  }

  $grpc.ResponseFuture<$0.SelectDoctorResponse> selectDoctor(
    $0.SelectDoctorRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$selectDoctor, request, options: options);
  }

  // method descriptors

  static final _$searchFacilities = $grpc.ClientMethod<
          $0.SearchFacilitiesRequest, $0.SearchFacilitiesResponse>(
      '/manpasik.v1.ReservationService/SearchFacilities',
      ($0.SearchFacilitiesRequest value) => value.writeToBuffer(),
      $0.SearchFacilitiesResponse.fromBuffer);
  static final _$getFacility =
      $grpc.ClientMethod<$0.GetFacilityRequest, $0.Facility>(
          '/manpasik.v1.ReservationService/GetFacility',
          ($0.GetFacilityRequest value) => value.writeToBuffer(),
          $0.Facility.fromBuffer);
  static final _$getAvailableSlots = $grpc.ClientMethod<
          $0.GetAvailableSlotsRequest, $0.GetAvailableSlotsResponse>(
      '/manpasik.v1.ReservationService/GetAvailableSlots',
      ($0.GetAvailableSlotsRequest value) => value.writeToBuffer(),
      $0.GetAvailableSlotsResponse.fromBuffer);
  static final _$createReservation =
      $grpc.ClientMethod<$0.CreateReservationRequest, $0.Reservation>(
          '/manpasik.v1.ReservationService/CreateReservation',
          ($0.CreateReservationRequest value) => value.writeToBuffer(),
          $0.Reservation.fromBuffer);
  static final _$getReservation =
      $grpc.ClientMethod<$0.GetReservationRequest, $0.Reservation>(
          '/manpasik.v1.ReservationService/GetReservation',
          ($0.GetReservationRequest value) => value.writeToBuffer(),
          $0.Reservation.fromBuffer);
  static final _$listReservations = $grpc.ClientMethod<
          $0.ListReservationsRequest, $0.ListReservationsResponse>(
      '/manpasik.v1.ReservationService/ListReservations',
      ($0.ListReservationsRequest value) => value.writeToBuffer(),
      $0.ListReservationsResponse.fromBuffer);
  static final _$cancelReservation = $grpc.ClientMethod<
          $0.CancelReservationRequest, $0.CancelReservationResponse>(
      '/manpasik.v1.ReservationService/CancelReservation',
      ($0.CancelReservationRequest value) => value.writeToBuffer(),
      $0.CancelReservationResponse.fromBuffer);
  static final _$listDoctorsByFacility = $grpc.ClientMethod<
          $0.ListDoctorsByFacilityRequest, $0.ListDoctorsByFacilityResponse>(
      '/manpasik.v1.ReservationService/ListDoctorsByFacility',
      ($0.ListDoctorsByFacilityRequest value) => value.writeToBuffer(),
      $0.ListDoctorsByFacilityResponse.fromBuffer);
  static final _$getDoctorAvailability = $grpc.ClientMethod<
          $0.GetDoctorAvailabilityRequest, $0.GetDoctorAvailabilityResponse>(
      '/manpasik.v1.ReservationService/GetDoctorAvailability',
      ($0.GetDoctorAvailabilityRequest value) => value.writeToBuffer(),
      $0.GetDoctorAvailabilityResponse.fromBuffer);
  static final _$selectDoctor =
      $grpc.ClientMethod<$0.SelectDoctorRequest, $0.SelectDoctorResponse>(
          '/manpasik.v1.ReservationService/SelectDoctor',
          ($0.SelectDoctorRequest value) => value.writeToBuffer(),
          $0.SelectDoctorResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.ReservationService')
abstract class ReservationServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.ReservationService';

  ReservationServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.SearchFacilitiesRequest,
            $0.SearchFacilitiesResponse>(
        'SearchFacilities',
        searchFacilities_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SearchFacilitiesRequest.fromBuffer(value),
        ($0.SearchFacilitiesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetFacilityRequest, $0.Facility>(
        'GetFacility',
        getFacility_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetFacilityRequest.fromBuffer(value),
        ($0.Facility value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetAvailableSlotsRequest,
            $0.GetAvailableSlotsResponse>(
        'GetAvailableSlots',
        getAvailableSlots_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetAvailableSlotsRequest.fromBuffer(value),
        ($0.GetAvailableSlotsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateReservationRequest, $0.Reservation>(
        'CreateReservation',
        createReservation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateReservationRequest.fromBuffer(value),
        ($0.Reservation value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetReservationRequest, $0.Reservation>(
        'GetReservation',
        getReservation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetReservationRequest.fromBuffer(value),
        ($0.Reservation value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListReservationsRequest,
            $0.ListReservationsResponse>(
        'ListReservations',
        listReservations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListReservationsRequest.fromBuffer(value),
        ($0.ListReservationsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CancelReservationRequest,
            $0.CancelReservationResponse>(
        'CancelReservation',
        cancelReservation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CancelReservationRequest.fromBuffer(value),
        ($0.CancelReservationResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListDoctorsByFacilityRequest,
            $0.ListDoctorsByFacilityResponse>(
        'ListDoctorsByFacility',
        listDoctorsByFacility_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDoctorsByFacilityRequest.fromBuffer(value),
        ($0.ListDoctorsByFacilityResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDoctorAvailabilityRequest,
            $0.GetDoctorAvailabilityResponse>(
        'GetDoctorAvailability',
        getDoctorAvailability_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDoctorAvailabilityRequest.fromBuffer(value),
        ($0.GetDoctorAvailabilityResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.SelectDoctorRequest, $0.SelectDoctorResponse>(
            'SelectDoctor',
            selectDoctor_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.SelectDoctorRequest.fromBuffer(value),
            ($0.SelectDoctorResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.SearchFacilitiesResponse> searchFacilities_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SearchFacilitiesRequest> $request) async {
    return searchFacilities($call, await $request);
  }

  $async.Future<$0.SearchFacilitiesResponse> searchFacilities(
      $grpc.ServiceCall call, $0.SearchFacilitiesRequest request);

  $async.Future<$0.Facility> getFacility_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetFacilityRequest> $request) async {
    return getFacility($call, await $request);
  }

  $async.Future<$0.Facility> getFacility(
      $grpc.ServiceCall call, $0.GetFacilityRequest request);

  $async.Future<$0.GetAvailableSlotsResponse> getAvailableSlots_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetAvailableSlotsRequest> $request) async {
    return getAvailableSlots($call, await $request);
  }

  $async.Future<$0.GetAvailableSlotsResponse> getAvailableSlots(
      $grpc.ServiceCall call, $0.GetAvailableSlotsRequest request);

  $async.Future<$0.Reservation> createReservation_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateReservationRequest> $request) async {
    return createReservation($call, await $request);
  }

  $async.Future<$0.Reservation> createReservation(
      $grpc.ServiceCall call, $0.CreateReservationRequest request);

  $async.Future<$0.Reservation> getReservation_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetReservationRequest> $request) async {
    return getReservation($call, await $request);
  }

  $async.Future<$0.Reservation> getReservation(
      $grpc.ServiceCall call, $0.GetReservationRequest request);

  $async.Future<$0.ListReservationsResponse> listReservations_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListReservationsRequest> $request) async {
    return listReservations($call, await $request);
  }

  $async.Future<$0.ListReservationsResponse> listReservations(
      $grpc.ServiceCall call, $0.ListReservationsRequest request);

  $async.Future<$0.CancelReservationResponse> cancelReservation_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CancelReservationRequest> $request) async {
    return cancelReservation($call, await $request);
  }

  $async.Future<$0.CancelReservationResponse> cancelReservation(
      $grpc.ServiceCall call, $0.CancelReservationRequest request);

  $async.Future<$0.ListDoctorsByFacilityResponse> listDoctorsByFacility_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListDoctorsByFacilityRequest> $request) async {
    return listDoctorsByFacility($call, await $request);
  }

  $async.Future<$0.ListDoctorsByFacilityResponse> listDoctorsByFacility(
      $grpc.ServiceCall call, $0.ListDoctorsByFacilityRequest request);

  $async.Future<$0.GetDoctorAvailabilityResponse> getDoctorAvailability_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDoctorAvailabilityRequest> $request) async {
    return getDoctorAvailability($call, await $request);
  }

  $async.Future<$0.GetDoctorAvailabilityResponse> getDoctorAvailability(
      $grpc.ServiceCall call, $0.GetDoctorAvailabilityRequest request);

  $async.Future<$0.SelectDoctorResponse> selectDoctor_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SelectDoctorRequest> $request) async {
    return selectDoctor($call, await $request);
  }

  $async.Future<$0.SelectDoctorResponse> selectDoctor(
      $grpc.ServiceCall call, $0.SelectDoctorRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.AdminService')
class AdminServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  AdminServiceClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.AdminUser> createAdmin(
    $0.CreateAdminRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createAdmin, request, options: options);
  }

  $grpc.ResponseFuture<$0.AdminUser> getAdmin(
    $0.GetAdminRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getAdmin, request, options: options);
  }

  $grpc.ResponseFuture<$0.ListAdminsResponse> listAdmins(
    $0.ListAdminsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listAdmins, request, options: options);
  }

  $grpc.ResponseFuture<$0.AdminUser> updateAdminRole(
    $0.UpdateAdminRoleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateAdminRole, request, options: options);
  }

  $grpc.ResponseFuture<$0.AdminUser> deactivateAdmin(
    $0.DeactivateAdminRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deactivateAdmin, request, options: options);
  }

  $grpc.ResponseFuture<$0.AdminListUsersResponse> listUsers(
    $0.AdminListUsersRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listUsers, request, options: options);
  }

  $grpc.ResponseFuture<$0.GetSystemStatsResponse> getSystemStats(
    $0.GetSystemStatsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSystemStats, request, options: options);
  }

  $grpc.ResponseFuture<$0.GetAuditLogResponse> getAuditLog(
    $0.GetAuditLogRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getAuditLog, request, options: options);
  }

  $grpc.ResponseFuture<$0.SystemConfig> setSystemConfig(
    $0.SetSystemConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$setSystemConfig, request, options: options);
  }

  $grpc.ResponseFuture<$0.SystemConfig> getSystemConfig(
    $0.GetSystemConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSystemConfig, request, options: options);
  }

  /// Phase 4: Region-based admin
  $grpc.ResponseFuture<$0.ListAdminsResponse> listAdminsByRegion(
    $0.ListAdminsByRegionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listAdminsByRegion, request, options: options);
  }

  /// Phase 5: 설정 관리 확장
  $grpc.ResponseFuture<$0.ListSystemConfigsResponse> listSystemConfigs(
    $0.ListSystemConfigsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listSystemConfigs, request, options: options);
  }

  $grpc.ResponseFuture<$0.ConfigWithMeta> getConfigWithMeta(
    $0.GetConfigWithMetaRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getConfigWithMeta, request, options: options);
  }

  $grpc.ResponseFuture<$0.ValidateConfigValueResponse> validateConfigValue(
    $0.ValidateConfigValueRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$validateConfigValue, request, options: options);
  }

  $grpc.ResponseFuture<$0.BulkSetConfigsResponse> bulkSetConfigs(
    $0.BulkSetConfigsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$bulkSetConfigs, request, options: options);
  }

  /// Phase 6: 감사 로그 상세 조회
  $grpc.ResponseFuture<$0.GetAuditLogDetailsResponse> getAuditLogDetails(
    $0.GetAuditLogDetailsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getAuditLogDetails, request, options: options);
  }

  /// Phase 9: 매출/재고 통계 (C12)
  $grpc.ResponseFuture<$0.GetRevenueStatsResponse> getRevenueStats(
    $0.GetRevenueStatsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getRevenueStats, request, options: options);
  }

  $grpc.ResponseFuture<$0.GetInventoryStatsResponse> getInventoryStats(
    $0.GetInventoryStatsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getInventoryStats, request, options: options);
  }

  // method descriptors

  static final _$createAdmin =
      $grpc.ClientMethod<$0.CreateAdminRequest, $0.AdminUser>(
          '/manpasik.v1.AdminService/CreateAdmin',
          ($0.CreateAdminRequest value) => value.writeToBuffer(),
          $0.AdminUser.fromBuffer);
  static final _$getAdmin =
      $grpc.ClientMethod<$0.GetAdminRequest, $0.AdminUser>(
          '/manpasik.v1.AdminService/GetAdmin',
          ($0.GetAdminRequest value) => value.writeToBuffer(),
          $0.AdminUser.fromBuffer);
  static final _$listAdmins =
      $grpc.ClientMethod<$0.ListAdminsRequest, $0.ListAdminsResponse>(
          '/manpasik.v1.AdminService/ListAdmins',
          ($0.ListAdminsRequest value) => value.writeToBuffer(),
          $0.ListAdminsResponse.fromBuffer);
  static final _$updateAdminRole =
      $grpc.ClientMethod<$0.UpdateAdminRoleRequest, $0.AdminUser>(
          '/manpasik.v1.AdminService/UpdateAdminRole',
          ($0.UpdateAdminRoleRequest value) => value.writeToBuffer(),
          $0.AdminUser.fromBuffer);
  static final _$deactivateAdmin =
      $grpc.ClientMethod<$0.DeactivateAdminRequest, $0.AdminUser>(
          '/manpasik.v1.AdminService/DeactivateAdmin',
          ($0.DeactivateAdminRequest value) => value.writeToBuffer(),
          $0.AdminUser.fromBuffer);
  static final _$listUsers =
      $grpc.ClientMethod<$0.AdminListUsersRequest, $0.AdminListUsersResponse>(
          '/manpasik.v1.AdminService/ListUsers',
          ($0.AdminListUsersRequest value) => value.writeToBuffer(),
          $0.AdminListUsersResponse.fromBuffer);
  static final _$getSystemStats =
      $grpc.ClientMethod<$0.GetSystemStatsRequest, $0.GetSystemStatsResponse>(
          '/manpasik.v1.AdminService/GetSystemStats',
          ($0.GetSystemStatsRequest value) => value.writeToBuffer(),
          $0.GetSystemStatsResponse.fromBuffer);
  static final _$getAuditLog =
      $grpc.ClientMethod<$0.GetAuditLogRequest, $0.GetAuditLogResponse>(
          '/manpasik.v1.AdminService/GetAuditLog',
          ($0.GetAuditLogRequest value) => value.writeToBuffer(),
          $0.GetAuditLogResponse.fromBuffer);
  static final _$setSystemConfig =
      $grpc.ClientMethod<$0.SetSystemConfigRequest, $0.SystemConfig>(
          '/manpasik.v1.AdminService/SetSystemConfig',
          ($0.SetSystemConfigRequest value) => value.writeToBuffer(),
          $0.SystemConfig.fromBuffer);
  static final _$getSystemConfig =
      $grpc.ClientMethod<$0.GetSystemConfigRequest, $0.SystemConfig>(
          '/manpasik.v1.AdminService/GetSystemConfig',
          ($0.GetSystemConfigRequest value) => value.writeToBuffer(),
          $0.SystemConfig.fromBuffer);
  static final _$listAdminsByRegion =
      $grpc.ClientMethod<$0.ListAdminsByRegionRequest, $0.ListAdminsResponse>(
          '/manpasik.v1.AdminService/ListAdminsByRegion',
          ($0.ListAdminsByRegionRequest value) => value.writeToBuffer(),
          $0.ListAdminsResponse.fromBuffer);
  static final _$listSystemConfigs = $grpc.ClientMethod<
          $0.ListSystemConfigsRequest, $0.ListSystemConfigsResponse>(
      '/manpasik.v1.AdminService/ListSystemConfigs',
      ($0.ListSystemConfigsRequest value) => value.writeToBuffer(),
      $0.ListSystemConfigsResponse.fromBuffer);
  static final _$getConfigWithMeta =
      $grpc.ClientMethod<$0.GetConfigWithMetaRequest, $0.ConfigWithMeta>(
          '/manpasik.v1.AdminService/GetConfigWithMeta',
          ($0.GetConfigWithMetaRequest value) => value.writeToBuffer(),
          $0.ConfigWithMeta.fromBuffer);
  static final _$validateConfigValue = $grpc.ClientMethod<
          $0.ValidateConfigValueRequest, $0.ValidateConfigValueResponse>(
      '/manpasik.v1.AdminService/ValidateConfigValue',
      ($0.ValidateConfigValueRequest value) => value.writeToBuffer(),
      $0.ValidateConfigValueResponse.fromBuffer);
  static final _$bulkSetConfigs =
      $grpc.ClientMethod<$0.BulkSetConfigsRequest, $0.BulkSetConfigsResponse>(
          '/manpasik.v1.AdminService/BulkSetConfigs',
          ($0.BulkSetConfigsRequest value) => value.writeToBuffer(),
          $0.BulkSetConfigsResponse.fromBuffer);
  static final _$getAuditLogDetails = $grpc.ClientMethod<
          $0.GetAuditLogDetailsRequest, $0.GetAuditLogDetailsResponse>(
      '/manpasik.v1.AdminService/GetAuditLogDetails',
      ($0.GetAuditLogDetailsRequest value) => value.writeToBuffer(),
      $0.GetAuditLogDetailsResponse.fromBuffer);
  static final _$getRevenueStats =
      $grpc.ClientMethod<$0.GetRevenueStatsRequest, $0.GetRevenueStatsResponse>(
          '/manpasik.v1.AdminService/GetRevenueStats',
          ($0.GetRevenueStatsRequest value) => value.writeToBuffer(),
          $0.GetRevenueStatsResponse.fromBuffer);
  static final _$getInventoryStats = $grpc.ClientMethod<
          $0.GetInventoryStatsRequest, $0.GetInventoryStatsResponse>(
      '/manpasik.v1.AdminService/GetInventoryStats',
      ($0.GetInventoryStatsRequest value) => value.writeToBuffer(),
      $0.GetInventoryStatsResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.AdminService')
abstract class AdminServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.AdminService';

  AdminServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.CreateAdminRequest, $0.AdminUser>(
        'CreateAdmin',
        createAdmin_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateAdminRequest.fromBuffer(value),
        ($0.AdminUser value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetAdminRequest, $0.AdminUser>(
        'GetAdmin',
        getAdmin_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetAdminRequest.fromBuffer(value),
        ($0.AdminUser value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListAdminsRequest, $0.ListAdminsResponse>(
        'ListAdmins',
        listAdmins_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListAdminsRequest.fromBuffer(value),
        ($0.ListAdminsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateAdminRoleRequest, $0.AdminUser>(
        'UpdateAdminRole',
        updateAdminRole_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateAdminRoleRequest.fromBuffer(value),
        ($0.AdminUser value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeactivateAdminRequest, $0.AdminUser>(
        'DeactivateAdmin',
        deactivateAdmin_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeactivateAdminRequest.fromBuffer(value),
        ($0.AdminUser value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AdminListUsersRequest,
            $0.AdminListUsersResponse>(
        'ListUsers',
        listUsers_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AdminListUsersRequest.fromBuffer(value),
        ($0.AdminListUsersResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetSystemStatsRequest,
            $0.GetSystemStatsResponse>(
        'GetSystemStats',
        getSystemStats_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetSystemStatsRequest.fromBuffer(value),
        ($0.GetSystemStatsResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetAuditLogRequest, $0.GetAuditLogResponse>(
            'GetAuditLog',
            getAuditLog_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetAuditLogRequest.fromBuffer(value),
            ($0.GetAuditLogResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SetSystemConfigRequest, $0.SystemConfig>(
        'SetSystemConfig',
        setSystemConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SetSystemConfigRequest.fromBuffer(value),
        ($0.SystemConfig value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetSystemConfigRequest, $0.SystemConfig>(
        'GetSystemConfig',
        getSystemConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetSystemConfigRequest.fromBuffer(value),
        ($0.SystemConfig value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListAdminsByRegionRequest,
            $0.ListAdminsResponse>(
        'ListAdminsByRegion',
        listAdminsByRegion_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListAdminsByRegionRequest.fromBuffer(value),
        ($0.ListAdminsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListSystemConfigsRequest,
            $0.ListSystemConfigsResponse>(
        'ListSystemConfigs',
        listSystemConfigs_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListSystemConfigsRequest.fromBuffer(value),
        ($0.ListSystemConfigsResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetConfigWithMetaRequest, $0.ConfigWithMeta>(
            'GetConfigWithMeta',
            getConfigWithMeta_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetConfigWithMetaRequest.fromBuffer(value),
            ($0.ConfigWithMeta value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ValidateConfigValueRequest,
            $0.ValidateConfigValueResponse>(
        'ValidateConfigValue',
        validateConfigValue_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ValidateConfigValueRequest.fromBuffer(value),
        ($0.ValidateConfigValueResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.BulkSetConfigsRequest,
            $0.BulkSetConfigsResponse>(
        'BulkSetConfigs',
        bulkSetConfigs_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.BulkSetConfigsRequest.fromBuffer(value),
        ($0.BulkSetConfigsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetAuditLogDetailsRequest,
            $0.GetAuditLogDetailsResponse>(
        'GetAuditLogDetails',
        getAuditLogDetails_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetAuditLogDetailsRequest.fromBuffer(value),
        ($0.GetAuditLogDetailsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetRevenueStatsRequest,
            $0.GetRevenueStatsResponse>(
        'GetRevenueStats',
        getRevenueStats_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetRevenueStatsRequest.fromBuffer(value),
        ($0.GetRevenueStatsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetInventoryStatsRequest,
            $0.GetInventoryStatsResponse>(
        'GetInventoryStats',
        getInventoryStats_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetInventoryStatsRequest.fromBuffer(value),
        ($0.GetInventoryStatsResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.AdminUser> createAdmin_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateAdminRequest> $request) async {
    return createAdmin($call, await $request);
  }

  $async.Future<$0.AdminUser> createAdmin(
      $grpc.ServiceCall call, $0.CreateAdminRequest request);

  $async.Future<$0.AdminUser> getAdmin_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetAdminRequest> $request) async {
    return getAdmin($call, await $request);
  }

  $async.Future<$0.AdminUser> getAdmin(
      $grpc.ServiceCall call, $0.GetAdminRequest request);

  $async.Future<$0.ListAdminsResponse> listAdmins_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListAdminsRequest> $request) async {
    return listAdmins($call, await $request);
  }

  $async.Future<$0.ListAdminsResponse> listAdmins(
      $grpc.ServiceCall call, $0.ListAdminsRequest request);

  $async.Future<$0.AdminUser> updateAdminRole_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateAdminRoleRequest> $request) async {
    return updateAdminRole($call, await $request);
  }

  $async.Future<$0.AdminUser> updateAdminRole(
      $grpc.ServiceCall call, $0.UpdateAdminRoleRequest request);

  $async.Future<$0.AdminUser> deactivateAdmin_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeactivateAdminRequest> $request) async {
    return deactivateAdmin($call, await $request);
  }

  $async.Future<$0.AdminUser> deactivateAdmin(
      $grpc.ServiceCall call, $0.DeactivateAdminRequest request);

  $async.Future<$0.AdminListUsersResponse> listUsers_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AdminListUsersRequest> $request) async {
    return listUsers($call, await $request);
  }

  $async.Future<$0.AdminListUsersResponse> listUsers(
      $grpc.ServiceCall call, $0.AdminListUsersRequest request);

  $async.Future<$0.GetSystemStatsResponse> getSystemStats_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetSystemStatsRequest> $request) async {
    return getSystemStats($call, await $request);
  }

  $async.Future<$0.GetSystemStatsResponse> getSystemStats(
      $grpc.ServiceCall call, $0.GetSystemStatsRequest request);

  $async.Future<$0.GetAuditLogResponse> getAuditLog_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetAuditLogRequest> $request) async {
    return getAuditLog($call, await $request);
  }

  $async.Future<$0.GetAuditLogResponse> getAuditLog(
      $grpc.ServiceCall call, $0.GetAuditLogRequest request);

  $async.Future<$0.SystemConfig> setSystemConfig_Pre($grpc.ServiceCall $call,
      $async.Future<$0.SetSystemConfigRequest> $request) async {
    return setSystemConfig($call, await $request);
  }

  $async.Future<$0.SystemConfig> setSystemConfig(
      $grpc.ServiceCall call, $0.SetSystemConfigRequest request);

  $async.Future<$0.SystemConfig> getSystemConfig_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetSystemConfigRequest> $request) async {
    return getSystemConfig($call, await $request);
  }

  $async.Future<$0.SystemConfig> getSystemConfig(
      $grpc.ServiceCall call, $0.GetSystemConfigRequest request);

  $async.Future<$0.ListAdminsResponse> listAdminsByRegion_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListAdminsByRegionRequest> $request) async {
    return listAdminsByRegion($call, await $request);
  }

  $async.Future<$0.ListAdminsResponse> listAdminsByRegion(
      $grpc.ServiceCall call, $0.ListAdminsByRegionRequest request);

  $async.Future<$0.ListSystemConfigsResponse> listSystemConfigs_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListSystemConfigsRequest> $request) async {
    return listSystemConfigs($call, await $request);
  }

  $async.Future<$0.ListSystemConfigsResponse> listSystemConfigs(
      $grpc.ServiceCall call, $0.ListSystemConfigsRequest request);

  $async.Future<$0.ConfigWithMeta> getConfigWithMeta_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetConfigWithMetaRequest> $request) async {
    return getConfigWithMeta($call, await $request);
  }

  $async.Future<$0.ConfigWithMeta> getConfigWithMeta(
      $grpc.ServiceCall call, $0.GetConfigWithMetaRequest request);

  $async.Future<$0.ValidateConfigValueResponse> validateConfigValue_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ValidateConfigValueRequest> $request) async {
    return validateConfigValue($call, await $request);
  }

  $async.Future<$0.ValidateConfigValueResponse> validateConfigValue(
      $grpc.ServiceCall call, $0.ValidateConfigValueRequest request);

  $async.Future<$0.BulkSetConfigsResponse> bulkSetConfigs_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.BulkSetConfigsRequest> $request) async {
    return bulkSetConfigs($call, await $request);
  }

  $async.Future<$0.BulkSetConfigsResponse> bulkSetConfigs(
      $grpc.ServiceCall call, $0.BulkSetConfigsRequest request);

  $async.Future<$0.GetAuditLogDetailsResponse> getAuditLogDetails_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetAuditLogDetailsRequest> $request) async {
    return getAuditLogDetails($call, await $request);
  }

  $async.Future<$0.GetAuditLogDetailsResponse> getAuditLogDetails(
      $grpc.ServiceCall call, $0.GetAuditLogDetailsRequest request);

  $async.Future<$0.GetRevenueStatsResponse> getRevenueStats_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetRevenueStatsRequest> $request) async {
    return getRevenueStats($call, await $request);
  }

  $async.Future<$0.GetRevenueStatsResponse> getRevenueStats(
      $grpc.ServiceCall call, $0.GetRevenueStatsRequest request);

  $async.Future<$0.GetInventoryStatsResponse> getInventoryStats_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetInventoryStatsRequest> $request) async {
    return getInventoryStats($call, await $request);
  }

  $async.Future<$0.GetInventoryStatsResponse> getInventoryStats(
      $grpc.ServiceCall call, $0.GetInventoryStatsRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.FamilyService')
class FamilyServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  FamilyServiceClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.FamilyGroup> createFamilyGroup(
    $0.CreateFamilyGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createFamilyGroup, request, options: options);
  }

  $grpc.ResponseFuture<$0.FamilyGroup> getFamilyGroup(
    $0.GetFamilyGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getFamilyGroup, request, options: options);
  }

  $grpc.ResponseFuture<$0.FamilyInvitation> inviteMember(
    $0.InviteMemberRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$inviteMember, request, options: options);
  }

  $grpc.ResponseFuture<$0.RespondToInvitationResponse> respondToInvitation(
    $0.RespondToInvitationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$respondToInvitation, request, options: options);
  }

  $grpc.ResponseFuture<$0.RemoveMemberResponse> removeMember(
    $0.RemoveMemberRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$removeMember, request, options: options);
  }

  $grpc.ResponseFuture<$0.FamilyMember> updateMemberRole(
    $0.UpdateMemberRoleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateMemberRole, request, options: options);
  }

  $grpc.ResponseFuture<$0.ListFamilyMembersResponse> listFamilyMembers(
    $0.ListFamilyMembersRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listFamilyMembers, request, options: options);
  }

  $grpc.ResponseFuture<$0.SharingPreferences> setSharingPreferences(
    $0.SetSharingPreferencesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$setSharingPreferences, request, options: options);
  }

  $grpc.ResponseFuture<$0.GetSharedHealthDataResponse> getSharedHealthData(
    $0.GetSharedHealthDataRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSharedHealthData, request, options: options);
  }

  /// Phase 6: 공유 접근 검증
  $grpc.ResponseFuture<$0.ValidateSharingAccessResponse> validateSharingAccess(
    $0.ValidateSharingAccessRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$validateSharingAccess, request, options: options);
  }

  // method descriptors

  static final _$createFamilyGroup =
      $grpc.ClientMethod<$0.CreateFamilyGroupRequest, $0.FamilyGroup>(
          '/manpasik.v1.FamilyService/CreateFamilyGroup',
          ($0.CreateFamilyGroupRequest value) => value.writeToBuffer(),
          $0.FamilyGroup.fromBuffer);
  static final _$getFamilyGroup =
      $grpc.ClientMethod<$0.GetFamilyGroupRequest, $0.FamilyGroup>(
          '/manpasik.v1.FamilyService/GetFamilyGroup',
          ($0.GetFamilyGroupRequest value) => value.writeToBuffer(),
          $0.FamilyGroup.fromBuffer);
  static final _$inviteMember =
      $grpc.ClientMethod<$0.InviteMemberRequest, $0.FamilyInvitation>(
          '/manpasik.v1.FamilyService/InviteMember',
          ($0.InviteMemberRequest value) => value.writeToBuffer(),
          $0.FamilyInvitation.fromBuffer);
  static final _$respondToInvitation = $grpc.ClientMethod<
          $0.RespondToInvitationRequest, $0.RespondToInvitationResponse>(
      '/manpasik.v1.FamilyService/RespondToInvitation',
      ($0.RespondToInvitationRequest value) => value.writeToBuffer(),
      $0.RespondToInvitationResponse.fromBuffer);
  static final _$removeMember =
      $grpc.ClientMethod<$0.RemoveMemberRequest, $0.RemoveMemberResponse>(
          '/manpasik.v1.FamilyService/RemoveMember',
          ($0.RemoveMemberRequest value) => value.writeToBuffer(),
          $0.RemoveMemberResponse.fromBuffer);
  static final _$updateMemberRole =
      $grpc.ClientMethod<$0.UpdateMemberRoleRequest, $0.FamilyMember>(
          '/manpasik.v1.FamilyService/UpdateMemberRole',
          ($0.UpdateMemberRoleRequest value) => value.writeToBuffer(),
          $0.FamilyMember.fromBuffer);
  static final _$listFamilyMembers = $grpc.ClientMethod<
          $0.ListFamilyMembersRequest, $0.ListFamilyMembersResponse>(
      '/manpasik.v1.FamilyService/ListFamilyMembers',
      ($0.ListFamilyMembersRequest value) => value.writeToBuffer(),
      $0.ListFamilyMembersResponse.fromBuffer);
  static final _$setSharingPreferences = $grpc.ClientMethod<
          $0.SetSharingPreferencesRequest, $0.SharingPreferences>(
      '/manpasik.v1.FamilyService/SetSharingPreferences',
      ($0.SetSharingPreferencesRequest value) => value.writeToBuffer(),
      $0.SharingPreferences.fromBuffer);
  static final _$getSharedHealthData = $grpc.ClientMethod<
          $0.GetSharedHealthDataRequest, $0.GetSharedHealthDataResponse>(
      '/manpasik.v1.FamilyService/GetSharedHealthData',
      ($0.GetSharedHealthDataRequest value) => value.writeToBuffer(),
      $0.GetSharedHealthDataResponse.fromBuffer);
  static final _$validateSharingAccess = $grpc.ClientMethod<
          $0.ValidateSharingAccessRequest, $0.ValidateSharingAccessResponse>(
      '/manpasik.v1.FamilyService/ValidateSharingAccess',
      ($0.ValidateSharingAccessRequest value) => value.writeToBuffer(),
      $0.ValidateSharingAccessResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.FamilyService')
abstract class FamilyServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.FamilyService';

  FamilyServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.CreateFamilyGroupRequest, $0.FamilyGroup>(
        'CreateFamilyGroup',
        createFamilyGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateFamilyGroupRequest.fromBuffer(value),
        ($0.FamilyGroup value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetFamilyGroupRequest, $0.FamilyGroup>(
        'GetFamilyGroup',
        getFamilyGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetFamilyGroupRequest.fromBuffer(value),
        ($0.FamilyGroup value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.InviteMemberRequest, $0.FamilyInvitation>(
        'InviteMember',
        inviteMember_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.InviteMemberRequest.fromBuffer(value),
        ($0.FamilyInvitation value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RespondToInvitationRequest,
            $0.RespondToInvitationResponse>(
        'RespondToInvitation',
        respondToInvitation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RespondToInvitationRequest.fromBuffer(value),
        ($0.RespondToInvitationResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.RemoveMemberRequest, $0.RemoveMemberResponse>(
            'RemoveMember',
            removeMember_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.RemoveMemberRequest.fromBuffer(value),
            ($0.RemoveMemberResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateMemberRoleRequest, $0.FamilyMember>(
        'UpdateMemberRole',
        updateMemberRole_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateMemberRoleRequest.fromBuffer(value),
        ($0.FamilyMember value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListFamilyMembersRequest,
            $0.ListFamilyMembersResponse>(
        'ListFamilyMembers',
        listFamilyMembers_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListFamilyMembersRequest.fromBuffer(value),
        ($0.ListFamilyMembersResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SetSharingPreferencesRequest,
            $0.SharingPreferences>(
        'SetSharingPreferences',
        setSharingPreferences_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SetSharingPreferencesRequest.fromBuffer(value),
        ($0.SharingPreferences value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetSharedHealthDataRequest,
            $0.GetSharedHealthDataResponse>(
        'GetSharedHealthData',
        getSharedHealthData_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetSharedHealthDataRequest.fromBuffer(value),
        ($0.GetSharedHealthDataResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ValidateSharingAccessRequest,
            $0.ValidateSharingAccessResponse>(
        'ValidateSharingAccess',
        validateSharingAccess_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ValidateSharingAccessRequest.fromBuffer(value),
        ($0.ValidateSharingAccessResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.FamilyGroup> createFamilyGroup_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateFamilyGroupRequest> $request) async {
    return createFamilyGroup($call, await $request);
  }

  $async.Future<$0.FamilyGroup> createFamilyGroup(
      $grpc.ServiceCall call, $0.CreateFamilyGroupRequest request);

  $async.Future<$0.FamilyGroup> getFamilyGroup_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetFamilyGroupRequest> $request) async {
    return getFamilyGroup($call, await $request);
  }

  $async.Future<$0.FamilyGroup> getFamilyGroup(
      $grpc.ServiceCall call, $0.GetFamilyGroupRequest request);

  $async.Future<$0.FamilyInvitation> inviteMember_Pre($grpc.ServiceCall $call,
      $async.Future<$0.InviteMemberRequest> $request) async {
    return inviteMember($call, await $request);
  }

  $async.Future<$0.FamilyInvitation> inviteMember(
      $grpc.ServiceCall call, $0.InviteMemberRequest request);

  $async.Future<$0.RespondToInvitationResponse> respondToInvitation_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RespondToInvitationRequest> $request) async {
    return respondToInvitation($call, await $request);
  }

  $async.Future<$0.RespondToInvitationResponse> respondToInvitation(
      $grpc.ServiceCall call, $0.RespondToInvitationRequest request);

  $async.Future<$0.RemoveMemberResponse> removeMember_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RemoveMemberRequest> $request) async {
    return removeMember($call, await $request);
  }

  $async.Future<$0.RemoveMemberResponse> removeMember(
      $grpc.ServiceCall call, $0.RemoveMemberRequest request);

  $async.Future<$0.FamilyMember> updateMemberRole_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateMemberRoleRequest> $request) async {
    return updateMemberRole($call, await $request);
  }

  $async.Future<$0.FamilyMember> updateMemberRole(
      $grpc.ServiceCall call, $0.UpdateMemberRoleRequest request);

  $async.Future<$0.ListFamilyMembersResponse> listFamilyMembers_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListFamilyMembersRequest> $request) async {
    return listFamilyMembers($call, await $request);
  }

  $async.Future<$0.ListFamilyMembersResponse> listFamilyMembers(
      $grpc.ServiceCall call, $0.ListFamilyMembersRequest request);

  $async.Future<$0.SharingPreferences> setSharingPreferences_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SetSharingPreferencesRequest> $request) async {
    return setSharingPreferences($call, await $request);
  }

  $async.Future<$0.SharingPreferences> setSharingPreferences(
      $grpc.ServiceCall call, $0.SetSharingPreferencesRequest request);

  $async.Future<$0.GetSharedHealthDataResponse> getSharedHealthData_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetSharedHealthDataRequest> $request) async {
    return getSharedHealthData($call, await $request);
  }

  $async.Future<$0.GetSharedHealthDataResponse> getSharedHealthData(
      $grpc.ServiceCall call, $0.GetSharedHealthDataRequest request);

  $async.Future<$0.ValidateSharingAccessResponse> validateSharingAccess_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ValidateSharingAccessRequest> $request) async {
    return validateSharingAccess($call, await $request);
  }

  $async.Future<$0.ValidateSharingAccessResponse> validateSharingAccess(
      $grpc.ServiceCall call, $0.ValidateSharingAccessRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.HealthRecordService')
class HealthRecordServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  HealthRecordServiceClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.HealthRecord> createRecord(
    $0.CreateHealthRecordRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createRecord, request, options: options);
  }

  $grpc.ResponseFuture<$0.HealthRecord> getRecord(
    $0.GetHealthRecordRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getRecord, request, options: options);
  }

  $grpc.ResponseFuture<$0.ListHealthRecordsResponse> listRecords(
    $0.ListHealthRecordsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listRecords, request, options: options);
  }

  $grpc.ResponseFuture<$0.HealthRecord> updateRecord(
    $0.UpdateHealthRecordRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateRecord, request, options: options);
  }

  $grpc.ResponseFuture<$0.DeleteHealthRecordResponse> deleteRecord(
    $0.DeleteHealthRecordRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteRecord, request, options: options);
  }

  $grpc.ResponseFuture<$0.ExportToFHIRResponse> exportToFHIR(
    $0.ExportToFHIRRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$exportToFHIR, request, options: options);
  }

  $grpc.ResponseFuture<$0.ImportFromFHIRResponse> importFromFHIR(
    $0.ImportFromFHIRRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$importFromFHIR, request, options: options);
  }

  $grpc.ResponseFuture<$0.GetHealthSummaryResponse> getHealthSummary(
    $0.GetHealthSummaryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getHealthSummary, request, options: options);
  }

  /// Phase 4: Data sharing & consent
  $grpc.ResponseFuture<$0.DataSharingConsent> createDataSharingConsent(
    $0.CreateConsentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createDataSharingConsent, request,
        options: options);
  }

  $grpc.ResponseFuture<$0.RevokeConsentResponse> revokeDataSharingConsent(
    $0.RevokeConsentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$revokeDataSharingConsent, request,
        options: options);
  }

  $grpc.ResponseFuture<$0.ListConsentsResponse> listDataSharingConsents(
    $0.ListConsentsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDataSharingConsents, request,
        options: options);
  }

  $grpc.ResponseFuture<$0.ShareWithProviderResponse> shareWithProvider(
    $0.ShareWithProviderRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$shareWithProvider, request, options: options);
  }

  $grpc.ResponseFuture<$0.GetDataAccessLogResponse> getDataAccessLog(
    $0.GetDataAccessLogRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDataAccessLog, request, options: options);
  }

  // method descriptors

  static final _$createRecord =
      $grpc.ClientMethod<$0.CreateHealthRecordRequest, $0.HealthRecord>(
          '/manpasik.v1.HealthRecordService/CreateRecord',
          ($0.CreateHealthRecordRequest value) => value.writeToBuffer(),
          $0.HealthRecord.fromBuffer);
  static final _$getRecord =
      $grpc.ClientMethod<$0.GetHealthRecordRequest, $0.HealthRecord>(
          '/manpasik.v1.HealthRecordService/GetRecord',
          ($0.GetHealthRecordRequest value) => value.writeToBuffer(),
          $0.HealthRecord.fromBuffer);
  static final _$listRecords = $grpc.ClientMethod<$0.ListHealthRecordsRequest,
          $0.ListHealthRecordsResponse>(
      '/manpasik.v1.HealthRecordService/ListRecords',
      ($0.ListHealthRecordsRequest value) => value.writeToBuffer(),
      $0.ListHealthRecordsResponse.fromBuffer);
  static final _$updateRecord =
      $grpc.ClientMethod<$0.UpdateHealthRecordRequest, $0.HealthRecord>(
          '/manpasik.v1.HealthRecordService/UpdateRecord',
          ($0.UpdateHealthRecordRequest value) => value.writeToBuffer(),
          $0.HealthRecord.fromBuffer);
  static final _$deleteRecord = $grpc.ClientMethod<$0.DeleteHealthRecordRequest,
          $0.DeleteHealthRecordResponse>(
      '/manpasik.v1.HealthRecordService/DeleteRecord',
      ($0.DeleteHealthRecordRequest value) => value.writeToBuffer(),
      $0.DeleteHealthRecordResponse.fromBuffer);
  static final _$exportToFHIR =
      $grpc.ClientMethod<$0.ExportToFHIRRequest, $0.ExportToFHIRResponse>(
          '/manpasik.v1.HealthRecordService/ExportToFHIR',
          ($0.ExportToFHIRRequest value) => value.writeToBuffer(),
          $0.ExportToFHIRResponse.fromBuffer);
  static final _$importFromFHIR =
      $grpc.ClientMethod<$0.ImportFromFHIRRequest, $0.ImportFromFHIRResponse>(
          '/manpasik.v1.HealthRecordService/ImportFromFHIR',
          ($0.ImportFromFHIRRequest value) => value.writeToBuffer(),
          $0.ImportFromFHIRResponse.fromBuffer);
  static final _$getHealthSummary = $grpc.ClientMethod<
          $0.GetHealthSummaryRequest, $0.GetHealthSummaryResponse>(
      '/manpasik.v1.HealthRecordService/GetHealthSummary',
      ($0.GetHealthSummaryRequest value) => value.writeToBuffer(),
      $0.GetHealthSummaryResponse.fromBuffer);
  static final _$createDataSharingConsent =
      $grpc.ClientMethod<$0.CreateConsentRequest, $0.DataSharingConsent>(
          '/manpasik.v1.HealthRecordService/CreateDataSharingConsent',
          ($0.CreateConsentRequest value) => value.writeToBuffer(),
          $0.DataSharingConsent.fromBuffer);
  static final _$revokeDataSharingConsent =
      $grpc.ClientMethod<$0.RevokeConsentRequest, $0.RevokeConsentResponse>(
          '/manpasik.v1.HealthRecordService/RevokeDataSharingConsent',
          ($0.RevokeConsentRequest value) => value.writeToBuffer(),
          $0.RevokeConsentResponse.fromBuffer);
  static final _$listDataSharingConsents =
      $grpc.ClientMethod<$0.ListConsentsRequest, $0.ListConsentsResponse>(
          '/manpasik.v1.HealthRecordService/ListDataSharingConsents',
          ($0.ListConsentsRequest value) => value.writeToBuffer(),
          $0.ListConsentsResponse.fromBuffer);
  static final _$shareWithProvider = $grpc.ClientMethod<
          $0.ShareWithProviderRequest, $0.ShareWithProviderResponse>(
      '/manpasik.v1.HealthRecordService/ShareWithProvider',
      ($0.ShareWithProviderRequest value) => value.writeToBuffer(),
      $0.ShareWithProviderResponse.fromBuffer);
  static final _$getDataAccessLog = $grpc.ClientMethod<
          $0.GetDataAccessLogRequest, $0.GetDataAccessLogResponse>(
      '/manpasik.v1.HealthRecordService/GetDataAccessLog',
      ($0.GetDataAccessLogRequest value) => value.writeToBuffer(),
      $0.GetDataAccessLogResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.HealthRecordService')
abstract class HealthRecordServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.HealthRecordService';

  HealthRecordServiceBase() {
    $addMethod(
        $grpc.ServiceMethod<$0.CreateHealthRecordRequest, $0.HealthRecord>(
            'CreateRecord',
            createRecord_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateHealthRecordRequest.fromBuffer(value),
            ($0.HealthRecord value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetHealthRecordRequest, $0.HealthRecord>(
        'GetRecord',
        getRecord_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetHealthRecordRequest.fromBuffer(value),
        ($0.HealthRecord value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListHealthRecordsRequest,
            $0.ListHealthRecordsResponse>(
        'ListRecords',
        listRecords_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListHealthRecordsRequest.fromBuffer(value),
        ($0.ListHealthRecordsResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateHealthRecordRequest, $0.HealthRecord>(
            'UpdateRecord',
            updateRecord_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateHealthRecordRequest.fromBuffer(value),
            ($0.HealthRecord value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteHealthRecordRequest,
            $0.DeleteHealthRecordResponse>(
        'DeleteRecord',
        deleteRecord_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteHealthRecordRequest.fromBuffer(value),
        ($0.DeleteHealthRecordResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ExportToFHIRRequest, $0.ExportToFHIRResponse>(
            'ExportToFHIR',
            exportToFHIR_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ExportToFHIRRequest.fromBuffer(value),
            ($0.ExportToFHIRResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ImportFromFHIRRequest,
            $0.ImportFromFHIRResponse>(
        'ImportFromFHIR',
        importFromFHIR_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ImportFromFHIRRequest.fromBuffer(value),
        ($0.ImportFromFHIRResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetHealthSummaryRequest,
            $0.GetHealthSummaryResponse>(
        'GetHealthSummary',
        getHealthSummary_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetHealthSummaryRequest.fromBuffer(value),
        ($0.GetHealthSummaryResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateConsentRequest, $0.DataSharingConsent>(
            'CreateDataSharingConsent',
            createDataSharingConsent_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateConsentRequest.fromBuffer(value),
            ($0.DataSharingConsent value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.RevokeConsentRequest, $0.RevokeConsentResponse>(
            'RevokeDataSharingConsent',
            revokeDataSharingConsent_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.RevokeConsentRequest.fromBuffer(value),
            ($0.RevokeConsentResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListConsentsRequest, $0.ListConsentsResponse>(
            'ListDataSharingConsents',
            listDataSharingConsents_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListConsentsRequest.fromBuffer(value),
            ($0.ListConsentsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ShareWithProviderRequest,
            $0.ShareWithProviderResponse>(
        'ShareWithProvider',
        shareWithProvider_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ShareWithProviderRequest.fromBuffer(value),
        ($0.ShareWithProviderResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDataAccessLogRequest,
            $0.GetDataAccessLogResponse>(
        'GetDataAccessLog',
        getDataAccessLog_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDataAccessLogRequest.fromBuffer(value),
        ($0.GetDataAccessLogResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.HealthRecord> createRecord_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateHealthRecordRequest> $request) async {
    return createRecord($call, await $request);
  }

  $async.Future<$0.HealthRecord> createRecord(
      $grpc.ServiceCall call, $0.CreateHealthRecordRequest request);

  $async.Future<$0.HealthRecord> getRecord_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetHealthRecordRequest> $request) async {
    return getRecord($call, await $request);
  }

  $async.Future<$0.HealthRecord> getRecord(
      $grpc.ServiceCall call, $0.GetHealthRecordRequest request);

  $async.Future<$0.ListHealthRecordsResponse> listRecords_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListHealthRecordsRequest> $request) async {
    return listRecords($call, await $request);
  }

  $async.Future<$0.ListHealthRecordsResponse> listRecords(
      $grpc.ServiceCall call, $0.ListHealthRecordsRequest request);

  $async.Future<$0.HealthRecord> updateRecord_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateHealthRecordRequest> $request) async {
    return updateRecord($call, await $request);
  }

  $async.Future<$0.HealthRecord> updateRecord(
      $grpc.ServiceCall call, $0.UpdateHealthRecordRequest request);

  $async.Future<$0.DeleteHealthRecordResponse> deleteRecord_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteHealthRecordRequest> $request) async {
    return deleteRecord($call, await $request);
  }

  $async.Future<$0.DeleteHealthRecordResponse> deleteRecord(
      $grpc.ServiceCall call, $0.DeleteHealthRecordRequest request);

  $async.Future<$0.ExportToFHIRResponse> exportToFHIR_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ExportToFHIRRequest> $request) async {
    return exportToFHIR($call, await $request);
  }

  $async.Future<$0.ExportToFHIRResponse> exportToFHIR(
      $grpc.ServiceCall call, $0.ExportToFHIRRequest request);

  $async.Future<$0.ImportFromFHIRResponse> importFromFHIR_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ImportFromFHIRRequest> $request) async {
    return importFromFHIR($call, await $request);
  }

  $async.Future<$0.ImportFromFHIRResponse> importFromFHIR(
      $grpc.ServiceCall call, $0.ImportFromFHIRRequest request);

  $async.Future<$0.GetHealthSummaryResponse> getHealthSummary_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetHealthSummaryRequest> $request) async {
    return getHealthSummary($call, await $request);
  }

  $async.Future<$0.GetHealthSummaryResponse> getHealthSummary(
      $grpc.ServiceCall call, $0.GetHealthSummaryRequest request);

  $async.Future<$0.DataSharingConsent> createDataSharingConsent_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateConsentRequest> $request) async {
    return createDataSharingConsent($call, await $request);
  }

  $async.Future<$0.DataSharingConsent> createDataSharingConsent(
      $grpc.ServiceCall call, $0.CreateConsentRequest request);

  $async.Future<$0.RevokeConsentResponse> revokeDataSharingConsent_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RevokeConsentRequest> $request) async {
    return revokeDataSharingConsent($call, await $request);
  }

  $async.Future<$0.RevokeConsentResponse> revokeDataSharingConsent(
      $grpc.ServiceCall call, $0.RevokeConsentRequest request);

  $async.Future<$0.ListConsentsResponse> listDataSharingConsents_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListConsentsRequest> $request) async {
    return listDataSharingConsents($call, await $request);
  }

  $async.Future<$0.ListConsentsResponse> listDataSharingConsents(
      $grpc.ServiceCall call, $0.ListConsentsRequest request);

  $async.Future<$0.ShareWithProviderResponse> shareWithProvider_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ShareWithProviderRequest> $request) async {
    return shareWithProvider($call, await $request);
  }

  $async.Future<$0.ShareWithProviderResponse> shareWithProvider(
      $grpc.ServiceCall call, $0.ShareWithProviderRequest request);

  $async.Future<$0.GetDataAccessLogResponse> getDataAccessLog_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDataAccessLogRequest> $request) async {
    return getDataAccessLog($call, await $request);
  }

  $async.Future<$0.GetDataAccessLogResponse> getDataAccessLog(
      $grpc.ServiceCall call, $0.GetDataAccessLogRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.PrescriptionService')
class PrescriptionServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  PrescriptionServiceClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.Prescription> createPrescription(
    $0.CreatePrescriptionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createPrescription, request, options: options);
  }

  $grpc.ResponseFuture<$0.Prescription> getPrescription(
    $0.GetPrescriptionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getPrescription, request, options: options);
  }

  $grpc.ResponseFuture<$0.ListPrescriptionsResponse> listPrescriptions(
    $0.ListPrescriptionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listPrescriptions, request, options: options);
  }

  $grpc.ResponseFuture<$0.Prescription> updatePrescriptionStatus(
    $0.UpdatePrescriptionStatusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updatePrescriptionStatus, request,
        options: options);
  }

  $grpc.ResponseFuture<$0.Prescription> addMedication(
    $0.AddMedicationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$addMedication, request, options: options);
  }

  $grpc.ResponseFuture<$0.Prescription> removeMedication(
    $0.RemoveMedicationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$removeMedication, request, options: options);
  }

  $grpc.ResponseFuture<$0.CheckDrugInteractionResponse> checkDrugInteraction(
    $0.CheckDrugInteractionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$checkDrugInteraction, request, options: options);
  }

  $grpc.ResponseFuture<$0.GetMedicationRemindersResponse>
      getMedicationReminders(
    $0.GetMedicationRemindersRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getMedicationReminders, request,
        options: options);
  }

  /// Phase 4: Pharmacy fulfillment flow
  $grpc.ResponseFuture<$0.SelectPharmacyResponse> selectPharmacyAndFulfillment(
    $0.SelectPharmacyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$selectPharmacyAndFulfillment, request,
        options: options);
  }

  $grpc.ResponseFuture<$0.SendToPharmacyResponse> sendPrescriptionToPharmacy(
    $0.SendToPharmacyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$sendPrescriptionToPharmacy, request,
        options: options);
  }

  $grpc.ResponseFuture<$0.Prescription> getPrescriptionByToken(
    $0.GetByTokenRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getPrescriptionByToken, request,
        options: options);
  }

  $grpc.ResponseFuture<$0.Prescription> updateDispensaryStatus(
    $0.UpdateDispensaryStatusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateDispensaryStatus, request,
        options: options);
  }

  // method descriptors

  static final _$createPrescription =
      $grpc.ClientMethod<$0.CreatePrescriptionRequest, $0.Prescription>(
          '/manpasik.v1.PrescriptionService/CreatePrescription',
          ($0.CreatePrescriptionRequest value) => value.writeToBuffer(),
          $0.Prescription.fromBuffer);
  static final _$getPrescription =
      $grpc.ClientMethod<$0.GetPrescriptionRequest, $0.Prescription>(
          '/manpasik.v1.PrescriptionService/GetPrescription',
          ($0.GetPrescriptionRequest value) => value.writeToBuffer(),
          $0.Prescription.fromBuffer);
  static final _$listPrescriptions = $grpc.ClientMethod<
          $0.ListPrescriptionsRequest, $0.ListPrescriptionsResponse>(
      '/manpasik.v1.PrescriptionService/ListPrescriptions',
      ($0.ListPrescriptionsRequest value) => value.writeToBuffer(),
      $0.ListPrescriptionsResponse.fromBuffer);
  static final _$updatePrescriptionStatus =
      $grpc.ClientMethod<$0.UpdatePrescriptionStatusRequest, $0.Prescription>(
          '/manpasik.v1.PrescriptionService/UpdatePrescriptionStatus',
          ($0.UpdatePrescriptionStatusRequest value) => value.writeToBuffer(),
          $0.Prescription.fromBuffer);
  static final _$addMedication =
      $grpc.ClientMethod<$0.AddMedicationRequest, $0.Prescription>(
          '/manpasik.v1.PrescriptionService/AddMedication',
          ($0.AddMedicationRequest value) => value.writeToBuffer(),
          $0.Prescription.fromBuffer);
  static final _$removeMedication =
      $grpc.ClientMethod<$0.RemoveMedicationRequest, $0.Prescription>(
          '/manpasik.v1.PrescriptionService/RemoveMedication',
          ($0.RemoveMedicationRequest value) => value.writeToBuffer(),
          $0.Prescription.fromBuffer);
  static final _$checkDrugInteraction = $grpc.ClientMethod<
          $0.CheckDrugInteractionRequest, $0.CheckDrugInteractionResponse>(
      '/manpasik.v1.PrescriptionService/CheckDrugInteraction',
      ($0.CheckDrugInteractionRequest value) => value.writeToBuffer(),
      $0.CheckDrugInteractionResponse.fromBuffer);
  static final _$getMedicationReminders = $grpc.ClientMethod<
          $0.GetMedicationRemindersRequest, $0.GetMedicationRemindersResponse>(
      '/manpasik.v1.PrescriptionService/GetMedicationReminders',
      ($0.GetMedicationRemindersRequest value) => value.writeToBuffer(),
      $0.GetMedicationRemindersResponse.fromBuffer);
  static final _$selectPharmacyAndFulfillment =
      $grpc.ClientMethod<$0.SelectPharmacyRequest, $0.SelectPharmacyResponse>(
          '/manpasik.v1.PrescriptionService/SelectPharmacyAndFulfillment',
          ($0.SelectPharmacyRequest value) => value.writeToBuffer(),
          $0.SelectPharmacyResponse.fromBuffer);
  static final _$sendPrescriptionToPharmacy =
      $grpc.ClientMethod<$0.SendToPharmacyRequest, $0.SendToPharmacyResponse>(
          '/manpasik.v1.PrescriptionService/SendPrescriptionToPharmacy',
          ($0.SendToPharmacyRequest value) => value.writeToBuffer(),
          $0.SendToPharmacyResponse.fromBuffer);
  static final _$getPrescriptionByToken =
      $grpc.ClientMethod<$0.GetByTokenRequest, $0.Prescription>(
          '/manpasik.v1.PrescriptionService/GetPrescriptionByToken',
          ($0.GetByTokenRequest value) => value.writeToBuffer(),
          $0.Prescription.fromBuffer);
  static final _$updateDispensaryStatus =
      $grpc.ClientMethod<$0.UpdateDispensaryStatusRequest, $0.Prescription>(
          '/manpasik.v1.PrescriptionService/UpdateDispensaryStatus',
          ($0.UpdateDispensaryStatusRequest value) => value.writeToBuffer(),
          $0.Prescription.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.PrescriptionService')
abstract class PrescriptionServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.PrescriptionService';

  PrescriptionServiceBase() {
    $addMethod(
        $grpc.ServiceMethod<$0.CreatePrescriptionRequest, $0.Prescription>(
            'CreatePrescription',
            createPrescription_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreatePrescriptionRequest.fromBuffer(value),
            ($0.Prescription value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetPrescriptionRequest, $0.Prescription>(
        'GetPrescription',
        getPrescription_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetPrescriptionRequest.fromBuffer(value),
        ($0.Prescription value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListPrescriptionsRequest,
            $0.ListPrescriptionsResponse>(
        'ListPrescriptions',
        listPrescriptions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListPrescriptionsRequest.fromBuffer(value),
        ($0.ListPrescriptionsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdatePrescriptionStatusRequest,
            $0.Prescription>(
        'UpdatePrescriptionStatus',
        updatePrescriptionStatus_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdatePrescriptionStatusRequest.fromBuffer(value),
        ($0.Prescription value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AddMedicationRequest, $0.Prescription>(
        'AddMedication',
        addMedication_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AddMedicationRequest.fromBuffer(value),
        ($0.Prescription value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RemoveMedicationRequest, $0.Prescription>(
        'RemoveMedication',
        removeMedication_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RemoveMedicationRequest.fromBuffer(value),
        ($0.Prescription value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CheckDrugInteractionRequest,
            $0.CheckDrugInteractionResponse>(
        'CheckDrugInteraction',
        checkDrugInteraction_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CheckDrugInteractionRequest.fromBuffer(value),
        ($0.CheckDrugInteractionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetMedicationRemindersRequest,
            $0.GetMedicationRemindersResponse>(
        'GetMedicationReminders',
        getMedicationReminders_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetMedicationRemindersRequest.fromBuffer(value),
        ($0.GetMedicationRemindersResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SelectPharmacyRequest,
            $0.SelectPharmacyResponse>(
        'SelectPharmacyAndFulfillment',
        selectPharmacyAndFulfillment_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SelectPharmacyRequest.fromBuffer(value),
        ($0.SelectPharmacyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SendToPharmacyRequest,
            $0.SendToPharmacyResponse>(
        'SendPrescriptionToPharmacy',
        sendPrescriptionToPharmacy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SendToPharmacyRequest.fromBuffer(value),
        ($0.SendToPharmacyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetByTokenRequest, $0.Prescription>(
        'GetPrescriptionByToken',
        getPrescriptionByToken_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetByTokenRequest.fromBuffer(value),
        ($0.Prescription value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateDispensaryStatusRequest, $0.Prescription>(
            'UpdateDispensaryStatus',
            updateDispensaryStatus_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateDispensaryStatusRequest.fromBuffer(value),
            ($0.Prescription value) => value.writeToBuffer()));
  }

  $async.Future<$0.Prescription> createPrescription_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreatePrescriptionRequest> $request) async {
    return createPrescription($call, await $request);
  }

  $async.Future<$0.Prescription> createPrescription(
      $grpc.ServiceCall call, $0.CreatePrescriptionRequest request);

  $async.Future<$0.Prescription> getPrescription_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetPrescriptionRequest> $request) async {
    return getPrescription($call, await $request);
  }

  $async.Future<$0.Prescription> getPrescription(
      $grpc.ServiceCall call, $0.GetPrescriptionRequest request);

  $async.Future<$0.ListPrescriptionsResponse> listPrescriptions_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListPrescriptionsRequest> $request) async {
    return listPrescriptions($call, await $request);
  }

  $async.Future<$0.ListPrescriptionsResponse> listPrescriptions(
      $grpc.ServiceCall call, $0.ListPrescriptionsRequest request);

  $async.Future<$0.Prescription> updatePrescriptionStatus_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdatePrescriptionStatusRequest> $request) async {
    return updatePrescriptionStatus($call, await $request);
  }

  $async.Future<$0.Prescription> updatePrescriptionStatus(
      $grpc.ServiceCall call, $0.UpdatePrescriptionStatusRequest request);

  $async.Future<$0.Prescription> addMedication_Pre($grpc.ServiceCall $call,
      $async.Future<$0.AddMedicationRequest> $request) async {
    return addMedication($call, await $request);
  }

  $async.Future<$0.Prescription> addMedication(
      $grpc.ServiceCall call, $0.AddMedicationRequest request);

  $async.Future<$0.Prescription> removeMedication_Pre($grpc.ServiceCall $call,
      $async.Future<$0.RemoveMedicationRequest> $request) async {
    return removeMedication($call, await $request);
  }

  $async.Future<$0.Prescription> removeMedication(
      $grpc.ServiceCall call, $0.RemoveMedicationRequest request);

  $async.Future<$0.CheckDrugInteractionResponse> checkDrugInteraction_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CheckDrugInteractionRequest> $request) async {
    return checkDrugInteraction($call, await $request);
  }

  $async.Future<$0.CheckDrugInteractionResponse> checkDrugInteraction(
      $grpc.ServiceCall call, $0.CheckDrugInteractionRequest request);

  $async.Future<$0.GetMedicationRemindersResponse> getMedicationReminders_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetMedicationRemindersRequest> $request) async {
    return getMedicationReminders($call, await $request);
  }

  $async.Future<$0.GetMedicationRemindersResponse> getMedicationReminders(
      $grpc.ServiceCall call, $0.GetMedicationRemindersRequest request);

  $async.Future<$0.SelectPharmacyResponse> selectPharmacyAndFulfillment_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SelectPharmacyRequest> $request) async {
    return selectPharmacyAndFulfillment($call, await $request);
  }

  $async.Future<$0.SelectPharmacyResponse> selectPharmacyAndFulfillment(
      $grpc.ServiceCall call, $0.SelectPharmacyRequest request);

  $async.Future<$0.SendToPharmacyResponse> sendPrescriptionToPharmacy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SendToPharmacyRequest> $request) async {
    return sendPrescriptionToPharmacy($call, await $request);
  }

  $async.Future<$0.SendToPharmacyResponse> sendPrescriptionToPharmacy(
      $grpc.ServiceCall call, $0.SendToPharmacyRequest request);

  $async.Future<$0.Prescription> getPrescriptionByToken_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetByTokenRequest> $request) async {
    return getPrescriptionByToken($call, await $request);
  }

  $async.Future<$0.Prescription> getPrescriptionByToken(
      $grpc.ServiceCall call, $0.GetByTokenRequest request);

  $async.Future<$0.Prescription> updateDispensaryStatus_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateDispensaryStatusRequest> $request) async {
    return updateDispensaryStatus($call, await $request);
  }

  $async.Future<$0.Prescription> updateDispensaryStatus(
      $grpc.ServiceCall call, $0.UpdateDispensaryStatusRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.CommunityService')
class CommunityServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  CommunityServiceClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.Post> createPost(
    $0.CreatePostRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createPost, request, options: options);
  }

  $grpc.ResponseFuture<$0.Post> getPost(
    $0.GetPostRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getPost, request, options: options);
  }

  $grpc.ResponseFuture<$0.ListPostsResponse> listPosts(
    $0.ListPostsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listPosts, request, options: options);
  }

  $grpc.ResponseFuture<$0.LikePostResponse> likePost(
    $0.LikePostRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$likePost, request, options: options);
  }

  $grpc.ResponseFuture<$0.Comment> createComment(
    $0.CreateCommentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createComment, request, options: options);
  }

  $grpc.ResponseFuture<$0.ListCommentsResponse> listComments(
    $0.ListCommentsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listComments, request, options: options);
  }

  $grpc.ResponseFuture<$0.Challenge> createChallenge(
    $0.CreateChallengeRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createChallenge, request, options: options);
  }

  $grpc.ResponseFuture<$0.Challenge> getChallenge(
    $0.GetChallengeRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getChallenge, request, options: options);
  }

  $grpc.ResponseFuture<$0.JoinChallengeResponse> joinChallenge(
    $0.JoinChallengeRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$joinChallenge, request, options: options);
  }

  $grpc.ResponseFuture<$0.ListChallengesResponse> listChallenges(
    $0.ListChallengesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listChallenges, request, options: options);
  }

  /// Phase 9: 챌린지 리더보드 (C8)
  $grpc.ResponseFuture<$0.GetChallengeLeaderboardResponse>
      getChallengeLeaderboard(
    $0.GetChallengeLeaderboardRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getChallengeLeaderboard, request,
        options: options);
  }

  $grpc.ResponseFuture<$0.UpdateChallengeProgressResponse>
      updateChallengeProgress(
    $0.UpdateChallengeProgressRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateChallengeProgress, request,
        options: options);
  }

  // method descriptors

  static final _$createPost = $grpc.ClientMethod<$0.CreatePostRequest, $0.Post>(
      '/manpasik.v1.CommunityService/CreatePost',
      ($0.CreatePostRequest value) => value.writeToBuffer(),
      $0.Post.fromBuffer);
  static final _$getPost = $grpc.ClientMethod<$0.GetPostRequest, $0.Post>(
      '/manpasik.v1.CommunityService/GetPost',
      ($0.GetPostRequest value) => value.writeToBuffer(),
      $0.Post.fromBuffer);
  static final _$listPosts =
      $grpc.ClientMethod<$0.ListPostsRequest, $0.ListPostsResponse>(
          '/manpasik.v1.CommunityService/ListPosts',
          ($0.ListPostsRequest value) => value.writeToBuffer(),
          $0.ListPostsResponse.fromBuffer);
  static final _$likePost =
      $grpc.ClientMethod<$0.LikePostRequest, $0.LikePostResponse>(
          '/manpasik.v1.CommunityService/LikePost',
          ($0.LikePostRequest value) => value.writeToBuffer(),
          $0.LikePostResponse.fromBuffer);
  static final _$createComment =
      $grpc.ClientMethod<$0.CreateCommentRequest, $0.Comment>(
          '/manpasik.v1.CommunityService/CreateComment',
          ($0.CreateCommentRequest value) => value.writeToBuffer(),
          $0.Comment.fromBuffer);
  static final _$listComments =
      $grpc.ClientMethod<$0.ListCommentsRequest, $0.ListCommentsResponse>(
          '/manpasik.v1.CommunityService/ListComments',
          ($0.ListCommentsRequest value) => value.writeToBuffer(),
          $0.ListCommentsResponse.fromBuffer);
  static final _$createChallenge =
      $grpc.ClientMethod<$0.CreateChallengeRequest, $0.Challenge>(
          '/manpasik.v1.CommunityService/CreateChallenge',
          ($0.CreateChallengeRequest value) => value.writeToBuffer(),
          $0.Challenge.fromBuffer);
  static final _$getChallenge =
      $grpc.ClientMethod<$0.GetChallengeRequest, $0.Challenge>(
          '/manpasik.v1.CommunityService/GetChallenge',
          ($0.GetChallengeRequest value) => value.writeToBuffer(),
          $0.Challenge.fromBuffer);
  static final _$joinChallenge =
      $grpc.ClientMethod<$0.JoinChallengeRequest, $0.JoinChallengeResponse>(
          '/manpasik.v1.CommunityService/JoinChallenge',
          ($0.JoinChallengeRequest value) => value.writeToBuffer(),
          $0.JoinChallengeResponse.fromBuffer);
  static final _$listChallenges =
      $grpc.ClientMethod<$0.ListChallengesRequest, $0.ListChallengesResponse>(
          '/manpasik.v1.CommunityService/ListChallenges',
          ($0.ListChallengesRequest value) => value.writeToBuffer(),
          $0.ListChallengesResponse.fromBuffer);
  static final _$getChallengeLeaderboard = $grpc.ClientMethod<
          $0.GetChallengeLeaderboardRequest,
          $0.GetChallengeLeaderboardResponse>(
      '/manpasik.v1.CommunityService/GetChallengeLeaderboard',
      ($0.GetChallengeLeaderboardRequest value) => value.writeToBuffer(),
      $0.GetChallengeLeaderboardResponse.fromBuffer);
  static final _$updateChallengeProgress = $grpc.ClientMethod<
          $0.UpdateChallengeProgressRequest,
          $0.UpdateChallengeProgressResponse>(
      '/manpasik.v1.CommunityService/UpdateChallengeProgress',
      ($0.UpdateChallengeProgressRequest value) => value.writeToBuffer(),
      $0.UpdateChallengeProgressResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.CommunityService')
abstract class CommunityServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.CommunityService';

  CommunityServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.CreatePostRequest, $0.Post>(
        'CreatePost',
        createPost_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.CreatePostRequest.fromBuffer(value),
        ($0.Post value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetPostRequest, $0.Post>(
        'GetPost',
        getPost_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetPostRequest.fromBuffer(value),
        ($0.Post value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListPostsRequest, $0.ListPostsResponse>(
        'ListPosts',
        listPosts_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListPostsRequest.fromBuffer(value),
        ($0.ListPostsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.LikePostRequest, $0.LikePostResponse>(
        'LikePost',
        likePost_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.LikePostRequest.fromBuffer(value),
        ($0.LikePostResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateCommentRequest, $0.Comment>(
        'CreateComment',
        createComment_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateCommentRequest.fromBuffer(value),
        ($0.Comment value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListCommentsRequest, $0.ListCommentsResponse>(
            'ListComments',
            listComments_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListCommentsRequest.fromBuffer(value),
            ($0.ListCommentsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateChallengeRequest, $0.Challenge>(
        'CreateChallenge',
        createChallenge_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateChallengeRequest.fromBuffer(value),
        ($0.Challenge value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetChallengeRequest, $0.Challenge>(
        'GetChallenge',
        getChallenge_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetChallengeRequest.fromBuffer(value),
        ($0.Challenge value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.JoinChallengeRequest, $0.JoinChallengeResponse>(
            'JoinChallenge',
            joinChallenge_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.JoinChallengeRequest.fromBuffer(value),
            ($0.JoinChallengeResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListChallengesRequest,
            $0.ListChallengesResponse>(
        'ListChallenges',
        listChallenges_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListChallengesRequest.fromBuffer(value),
        ($0.ListChallengesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetChallengeLeaderboardRequest,
            $0.GetChallengeLeaderboardResponse>(
        'GetChallengeLeaderboard',
        getChallengeLeaderboard_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetChallengeLeaderboardRequest.fromBuffer(value),
        ($0.GetChallengeLeaderboardResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateChallengeProgressRequest,
            $0.UpdateChallengeProgressResponse>(
        'UpdateChallengeProgress',
        updateChallengeProgress_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateChallengeProgressRequest.fromBuffer(value),
        ($0.UpdateChallengeProgressResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.Post> createPost_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreatePostRequest> $request) async {
    return createPost($call, await $request);
  }

  $async.Future<$0.Post> createPost(
      $grpc.ServiceCall call, $0.CreatePostRequest request);

  $async.Future<$0.Post> getPost_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetPostRequest> $request) async {
    return getPost($call, await $request);
  }

  $async.Future<$0.Post> getPost(
      $grpc.ServiceCall call, $0.GetPostRequest request);

  $async.Future<$0.ListPostsResponse> listPosts_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListPostsRequest> $request) async {
    return listPosts($call, await $request);
  }

  $async.Future<$0.ListPostsResponse> listPosts(
      $grpc.ServiceCall call, $0.ListPostsRequest request);

  $async.Future<$0.LikePostResponse> likePost_Pre($grpc.ServiceCall $call,
      $async.Future<$0.LikePostRequest> $request) async {
    return likePost($call, await $request);
  }

  $async.Future<$0.LikePostResponse> likePost(
      $grpc.ServiceCall call, $0.LikePostRequest request);

  $async.Future<$0.Comment> createComment_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateCommentRequest> $request) async {
    return createComment($call, await $request);
  }

  $async.Future<$0.Comment> createComment(
      $grpc.ServiceCall call, $0.CreateCommentRequest request);

  $async.Future<$0.ListCommentsResponse> listComments_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListCommentsRequest> $request) async {
    return listComments($call, await $request);
  }

  $async.Future<$0.ListCommentsResponse> listComments(
      $grpc.ServiceCall call, $0.ListCommentsRequest request);

  $async.Future<$0.Challenge> createChallenge_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateChallengeRequest> $request) async {
    return createChallenge($call, await $request);
  }

  $async.Future<$0.Challenge> createChallenge(
      $grpc.ServiceCall call, $0.CreateChallengeRequest request);

  $async.Future<$0.Challenge> getChallenge_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetChallengeRequest> $request) async {
    return getChallenge($call, await $request);
  }

  $async.Future<$0.Challenge> getChallenge(
      $grpc.ServiceCall call, $0.GetChallengeRequest request);

  $async.Future<$0.JoinChallengeResponse> joinChallenge_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.JoinChallengeRequest> $request) async {
    return joinChallenge($call, await $request);
  }

  $async.Future<$0.JoinChallengeResponse> joinChallenge(
      $grpc.ServiceCall call, $0.JoinChallengeRequest request);

  $async.Future<$0.ListChallengesResponse> listChallenges_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListChallengesRequest> $request) async {
    return listChallenges($call, await $request);
  }

  $async.Future<$0.ListChallengesResponse> listChallenges(
      $grpc.ServiceCall call, $0.ListChallengesRequest request);

  $async.Future<$0.GetChallengeLeaderboardResponse> getChallengeLeaderboard_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetChallengeLeaderboardRequest> $request) async {
    return getChallengeLeaderboard($call, await $request);
  }

  $async.Future<$0.GetChallengeLeaderboardResponse> getChallengeLeaderboard(
      $grpc.ServiceCall call, $0.GetChallengeLeaderboardRequest request);

  $async.Future<$0.UpdateChallengeProgressResponse> updateChallengeProgress_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateChallengeProgressRequest> $request) async {
    return updateChallengeProgress($call, await $request);
  }

  $async.Future<$0.UpdateChallengeProgressResponse> updateChallengeProgress(
      $grpc.ServiceCall call, $0.UpdateChallengeProgressRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.VideoService')
class VideoServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  VideoServiceClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.Room> createRoom(
    $0.CreateRoomRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createRoom, request, options: options);
  }

  $grpc.ResponseFuture<$0.Room> getRoom(
    $0.GetRoomRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getRoom, request, options: options);
  }

  $grpc.ResponseFuture<$0.JoinRoomResponse> joinRoom(
    $0.JoinRoomRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$joinRoom, request, options: options);
  }

  $grpc.ResponseFuture<$0.LeaveRoomResponse> leaveRoom(
    $0.LeaveRoomRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$leaveRoom, request, options: options);
  }

  $grpc.ResponseFuture<$0.Room> endRoom(
    $0.EndRoomRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$endRoom, request, options: options);
  }

  $grpc.ResponseFuture<$0.SendSignalResponse> sendSignal(
    $0.SendSignalRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$sendSignal, request, options: options);
  }

  $grpc.ResponseFuture<$0.ListParticipantsResponse> listParticipants(
    $0.ListParticipantsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listParticipants, request, options: options);
  }

  $grpc.ResponseFuture<$0.GetRoomStatsResponse> getRoomStats(
    $0.GetRoomStatsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getRoomStats, request, options: options);
  }

  // method descriptors

  static final _$createRoom = $grpc.ClientMethod<$0.CreateRoomRequest, $0.Room>(
      '/manpasik.v1.VideoService/CreateRoom',
      ($0.CreateRoomRequest value) => value.writeToBuffer(),
      $0.Room.fromBuffer);
  static final _$getRoom = $grpc.ClientMethod<$0.GetRoomRequest, $0.Room>(
      '/manpasik.v1.VideoService/GetRoom',
      ($0.GetRoomRequest value) => value.writeToBuffer(),
      $0.Room.fromBuffer);
  static final _$joinRoom =
      $grpc.ClientMethod<$0.JoinRoomRequest, $0.JoinRoomResponse>(
          '/manpasik.v1.VideoService/JoinRoom',
          ($0.JoinRoomRequest value) => value.writeToBuffer(),
          $0.JoinRoomResponse.fromBuffer);
  static final _$leaveRoom =
      $grpc.ClientMethod<$0.LeaveRoomRequest, $0.LeaveRoomResponse>(
          '/manpasik.v1.VideoService/LeaveRoom',
          ($0.LeaveRoomRequest value) => value.writeToBuffer(),
          $0.LeaveRoomResponse.fromBuffer);
  static final _$endRoom = $grpc.ClientMethod<$0.EndRoomRequest, $0.Room>(
      '/manpasik.v1.VideoService/EndRoom',
      ($0.EndRoomRequest value) => value.writeToBuffer(),
      $0.Room.fromBuffer);
  static final _$sendSignal =
      $grpc.ClientMethod<$0.SendSignalRequest, $0.SendSignalResponse>(
          '/manpasik.v1.VideoService/SendSignal',
          ($0.SendSignalRequest value) => value.writeToBuffer(),
          $0.SendSignalResponse.fromBuffer);
  static final _$listParticipants = $grpc.ClientMethod<
          $0.ListParticipantsRequest, $0.ListParticipantsResponse>(
      '/manpasik.v1.VideoService/ListParticipants',
      ($0.ListParticipantsRequest value) => value.writeToBuffer(),
      $0.ListParticipantsResponse.fromBuffer);
  static final _$getRoomStats =
      $grpc.ClientMethod<$0.GetRoomStatsRequest, $0.GetRoomStatsResponse>(
          '/manpasik.v1.VideoService/GetRoomStats',
          ($0.GetRoomStatsRequest value) => value.writeToBuffer(),
          $0.GetRoomStatsResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.VideoService')
abstract class VideoServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.VideoService';

  VideoServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.CreateRoomRequest, $0.Room>(
        'CreateRoom',
        createRoom_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.CreateRoomRequest.fromBuffer(value),
        ($0.Room value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetRoomRequest, $0.Room>(
        'GetRoom',
        getRoom_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetRoomRequest.fromBuffer(value),
        ($0.Room value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.JoinRoomRequest, $0.JoinRoomResponse>(
        'JoinRoom',
        joinRoom_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.JoinRoomRequest.fromBuffer(value),
        ($0.JoinRoomResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.LeaveRoomRequest, $0.LeaveRoomResponse>(
        'LeaveRoom',
        leaveRoom_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.LeaveRoomRequest.fromBuffer(value),
        ($0.LeaveRoomResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.EndRoomRequest, $0.Room>(
        'EndRoom',
        endRoom_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.EndRoomRequest.fromBuffer(value),
        ($0.Room value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SendSignalRequest, $0.SendSignalResponse>(
        'SendSignal',
        sendSignal_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.SendSignalRequest.fromBuffer(value),
        ($0.SendSignalResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListParticipantsRequest,
            $0.ListParticipantsResponse>(
        'ListParticipants',
        listParticipants_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListParticipantsRequest.fromBuffer(value),
        ($0.ListParticipantsResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetRoomStatsRequest, $0.GetRoomStatsResponse>(
            'GetRoomStats',
            getRoomStats_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetRoomStatsRequest.fromBuffer(value),
            ($0.GetRoomStatsResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.Room> createRoom_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateRoomRequest> $request) async {
    return createRoom($call, await $request);
  }

  $async.Future<$0.Room> createRoom(
      $grpc.ServiceCall call, $0.CreateRoomRequest request);

  $async.Future<$0.Room> getRoom_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetRoomRequest> $request) async {
    return getRoom($call, await $request);
  }

  $async.Future<$0.Room> getRoom(
      $grpc.ServiceCall call, $0.GetRoomRequest request);

  $async.Future<$0.JoinRoomResponse> joinRoom_Pre($grpc.ServiceCall $call,
      $async.Future<$0.JoinRoomRequest> $request) async {
    return joinRoom($call, await $request);
  }

  $async.Future<$0.JoinRoomResponse> joinRoom(
      $grpc.ServiceCall call, $0.JoinRoomRequest request);

  $async.Future<$0.LeaveRoomResponse> leaveRoom_Pre($grpc.ServiceCall $call,
      $async.Future<$0.LeaveRoomRequest> $request) async {
    return leaveRoom($call, await $request);
  }

  $async.Future<$0.LeaveRoomResponse> leaveRoom(
      $grpc.ServiceCall call, $0.LeaveRoomRequest request);

  $async.Future<$0.Room> endRoom_Pre($grpc.ServiceCall $call,
      $async.Future<$0.EndRoomRequest> $request) async {
    return endRoom($call, await $request);
  }

  $async.Future<$0.Room> endRoom(
      $grpc.ServiceCall call, $0.EndRoomRequest request);

  $async.Future<$0.SendSignalResponse> sendSignal_Pre($grpc.ServiceCall $call,
      $async.Future<$0.SendSignalRequest> $request) async {
    return sendSignal($call, await $request);
  }

  $async.Future<$0.SendSignalResponse> sendSignal(
      $grpc.ServiceCall call, $0.SendSignalRequest request);

  $async.Future<$0.ListParticipantsResponse> listParticipants_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListParticipantsRequest> $request) async {
    return listParticipants($call, await $request);
  }

  $async.Future<$0.ListParticipantsResponse> listParticipants(
      $grpc.ServiceCall call, $0.ListParticipantsRequest request);

  $async.Future<$0.GetRoomStatsResponse> getRoomStats_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetRoomStatsRequest> $request) async {
    return getRoomStats($call, await $request);
  }

  $async.Future<$0.GetRoomStatsResponse> getRoomStats(
      $grpc.ServiceCall call, $0.GetRoomStatsRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.NotificationService')
class NotificationServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  NotificationServiceClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.Notification> sendNotification(
    $0.SendNotificationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$sendNotification, request, options: options);
  }

  $grpc.ResponseFuture<$0.ListNotificationsResponse> listNotifications(
    $0.ListNotificationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listNotifications, request, options: options);
  }

  $grpc.ResponseFuture<$0.MarkAsReadResponse> markAsRead(
    $0.MarkAsReadRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$markAsRead, request, options: options);
  }

  $grpc.ResponseFuture<$0.MarkAllAsReadResponse> markAllAsRead(
    $0.MarkAllAsReadRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$markAllAsRead, request, options: options);
  }

  $grpc.ResponseFuture<$0.GetUnreadCountResponse> getUnreadCount(
    $0.GetUnreadCountRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getUnreadCount, request, options: options);
  }

  $grpc.ResponseFuture<$0.NotificationPreferences>
      updateNotificationPreferences(
    $0.UpdateNotificationPreferencesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateNotificationPreferences, request,
        options: options);
  }

  $grpc.ResponseFuture<$0.NotificationPreferences> getNotificationPreferences(
    $0.GetNotificationPreferencesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getNotificationPreferences, request,
        options: options);
  }

  /// Phase 6: 템플릿 기반 알림 발송
  $grpc.ResponseFuture<$0.Notification> sendFromTemplate(
    $0.SendFromTemplateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$sendFromTemplate, request, options: options);
  }

  // method descriptors

  static final _$sendNotification =
      $grpc.ClientMethod<$0.SendNotificationRequest, $0.Notification>(
          '/manpasik.v1.NotificationService/SendNotification',
          ($0.SendNotificationRequest value) => value.writeToBuffer(),
          $0.Notification.fromBuffer);
  static final _$listNotifications = $grpc.ClientMethod<
          $0.ListNotificationsRequest, $0.ListNotificationsResponse>(
      '/manpasik.v1.NotificationService/ListNotifications',
      ($0.ListNotificationsRequest value) => value.writeToBuffer(),
      $0.ListNotificationsResponse.fromBuffer);
  static final _$markAsRead =
      $grpc.ClientMethod<$0.MarkAsReadRequest, $0.MarkAsReadResponse>(
          '/manpasik.v1.NotificationService/MarkAsRead',
          ($0.MarkAsReadRequest value) => value.writeToBuffer(),
          $0.MarkAsReadResponse.fromBuffer);
  static final _$markAllAsRead =
      $grpc.ClientMethod<$0.MarkAllAsReadRequest, $0.MarkAllAsReadResponse>(
          '/manpasik.v1.NotificationService/MarkAllAsRead',
          ($0.MarkAllAsReadRequest value) => value.writeToBuffer(),
          $0.MarkAllAsReadResponse.fromBuffer);
  static final _$getUnreadCount =
      $grpc.ClientMethod<$0.GetUnreadCountRequest, $0.GetUnreadCountResponse>(
          '/manpasik.v1.NotificationService/GetUnreadCount',
          ($0.GetUnreadCountRequest value) => value.writeToBuffer(),
          $0.GetUnreadCountResponse.fromBuffer);
  static final _$updateNotificationPreferences = $grpc.ClientMethod<
          $0.UpdateNotificationPreferencesRequest, $0.NotificationPreferences>(
      '/manpasik.v1.NotificationService/UpdateNotificationPreferences',
      ($0.UpdateNotificationPreferencesRequest value) => value.writeToBuffer(),
      $0.NotificationPreferences.fromBuffer);
  static final _$getNotificationPreferences = $grpc.ClientMethod<
          $0.GetNotificationPreferencesRequest, $0.NotificationPreferences>(
      '/manpasik.v1.NotificationService/GetNotificationPreferences',
      ($0.GetNotificationPreferencesRequest value) => value.writeToBuffer(),
      $0.NotificationPreferences.fromBuffer);
  static final _$sendFromTemplate =
      $grpc.ClientMethod<$0.SendFromTemplateRequest, $0.Notification>(
          '/manpasik.v1.NotificationService/SendFromTemplate',
          ($0.SendFromTemplateRequest value) => value.writeToBuffer(),
          $0.Notification.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.NotificationService')
abstract class NotificationServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.NotificationService';

  NotificationServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.SendNotificationRequest, $0.Notification>(
        'SendNotification',
        sendNotification_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SendNotificationRequest.fromBuffer(value),
        ($0.Notification value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListNotificationsRequest,
            $0.ListNotificationsResponse>(
        'ListNotifications',
        listNotifications_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListNotificationsRequest.fromBuffer(value),
        ($0.ListNotificationsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.MarkAsReadRequest, $0.MarkAsReadResponse>(
        'MarkAsRead',
        markAsRead_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.MarkAsReadRequest.fromBuffer(value),
        ($0.MarkAsReadResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.MarkAllAsReadRequest, $0.MarkAllAsReadResponse>(
            'MarkAllAsRead',
            markAllAsRead_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.MarkAllAsReadRequest.fromBuffer(value),
            ($0.MarkAllAsReadResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetUnreadCountRequest,
            $0.GetUnreadCountResponse>(
        'GetUnreadCount',
        getUnreadCount_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetUnreadCountRequest.fromBuffer(value),
        ($0.GetUnreadCountResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateNotificationPreferencesRequest,
            $0.NotificationPreferences>(
        'UpdateNotificationPreferences',
        updateNotificationPreferences_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateNotificationPreferencesRequest.fromBuffer(value),
        ($0.NotificationPreferences value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetNotificationPreferencesRequest,
            $0.NotificationPreferences>(
        'GetNotificationPreferences',
        getNotificationPreferences_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetNotificationPreferencesRequest.fromBuffer(value),
        ($0.NotificationPreferences value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SendFromTemplateRequest, $0.Notification>(
        'SendFromTemplate',
        sendFromTemplate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SendFromTemplateRequest.fromBuffer(value),
        ($0.Notification value) => value.writeToBuffer()));
  }

  $async.Future<$0.Notification> sendNotification_Pre($grpc.ServiceCall $call,
      $async.Future<$0.SendNotificationRequest> $request) async {
    return sendNotification($call, await $request);
  }

  $async.Future<$0.Notification> sendNotification(
      $grpc.ServiceCall call, $0.SendNotificationRequest request);

  $async.Future<$0.ListNotificationsResponse> listNotifications_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListNotificationsRequest> $request) async {
    return listNotifications($call, await $request);
  }

  $async.Future<$0.ListNotificationsResponse> listNotifications(
      $grpc.ServiceCall call, $0.ListNotificationsRequest request);

  $async.Future<$0.MarkAsReadResponse> markAsRead_Pre($grpc.ServiceCall $call,
      $async.Future<$0.MarkAsReadRequest> $request) async {
    return markAsRead($call, await $request);
  }

  $async.Future<$0.MarkAsReadResponse> markAsRead(
      $grpc.ServiceCall call, $0.MarkAsReadRequest request);

  $async.Future<$0.MarkAllAsReadResponse> markAllAsRead_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.MarkAllAsReadRequest> $request) async {
    return markAllAsRead($call, await $request);
  }

  $async.Future<$0.MarkAllAsReadResponse> markAllAsRead(
      $grpc.ServiceCall call, $0.MarkAllAsReadRequest request);

  $async.Future<$0.GetUnreadCountResponse> getUnreadCount_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetUnreadCountRequest> $request) async {
    return getUnreadCount($call, await $request);
  }

  $async.Future<$0.GetUnreadCountResponse> getUnreadCount(
      $grpc.ServiceCall call, $0.GetUnreadCountRequest request);

  $async.Future<$0.NotificationPreferences> updateNotificationPreferences_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateNotificationPreferencesRequest> $request) async {
    return updateNotificationPreferences($call, await $request);
  }

  $async.Future<$0.NotificationPreferences> updateNotificationPreferences(
      $grpc.ServiceCall call, $0.UpdateNotificationPreferencesRequest request);

  $async.Future<$0.NotificationPreferences> getNotificationPreferences_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetNotificationPreferencesRequest> $request) async {
    return getNotificationPreferences($call, await $request);
  }

  $async.Future<$0.NotificationPreferences> getNotificationPreferences(
      $grpc.ServiceCall call, $0.GetNotificationPreferencesRequest request);

  $async.Future<$0.Notification> sendFromTemplate_Pre($grpc.ServiceCall $call,
      $async.Future<$0.SendFromTemplateRequest> $request) async {
    return sendFromTemplate($call, await $request);
  }

  $async.Future<$0.Notification> sendFromTemplate(
      $grpc.ServiceCall call, $0.SendFromTemplateRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.TranslationService')
class TranslationServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  TranslationServiceClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.TranslateTextResponse> translateText(
    $0.TranslateTextRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$translateText, request, options: options);
  }

  $grpc.ResponseFuture<$0.DetectLanguageResponse> detectLanguage(
    $0.DetectLanguageRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$detectLanguage, request, options: options);
  }

  $grpc.ResponseFuture<$0.ListSupportedLanguagesResponse>
      listSupportedLanguages(
    $0.ListSupportedLanguagesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listSupportedLanguages, request,
        options: options);
  }

  $grpc.ResponseFuture<$0.TranslateBatchResponse> translateBatch(
    $0.TranslateBatchRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$translateBatch, request, options: options);
  }

  $grpc.ResponseFuture<$0.GetTranslationHistoryResponse> getTranslationHistory(
    $0.GetTranslationHistoryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getTranslationHistory, request, options: options);
  }

  $grpc.ResponseFuture<$0.GetTranslationUsageResponse> getTranslationUsage(
    $0.GetTranslationUsageRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getTranslationUsage, request, options: options);
  }

  /// Phase 9: 실시간 번역 (C6)
  $grpc.ResponseFuture<$0.TranslateRealtimeResponse> translateRealtime(
    $0.TranslateRealtimeRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$translateRealtime, request, options: options);
  }

  // method descriptors

  static final _$translateText =
      $grpc.ClientMethod<$0.TranslateTextRequest, $0.TranslateTextResponse>(
          '/manpasik.v1.TranslationService/TranslateText',
          ($0.TranslateTextRequest value) => value.writeToBuffer(),
          $0.TranslateTextResponse.fromBuffer);
  static final _$detectLanguage =
      $grpc.ClientMethod<$0.DetectLanguageRequest, $0.DetectLanguageResponse>(
          '/manpasik.v1.TranslationService/DetectLanguage',
          ($0.DetectLanguageRequest value) => value.writeToBuffer(),
          $0.DetectLanguageResponse.fromBuffer);
  static final _$listSupportedLanguages = $grpc.ClientMethod<
          $0.ListSupportedLanguagesRequest, $0.ListSupportedLanguagesResponse>(
      '/manpasik.v1.TranslationService/ListSupportedLanguages',
      ($0.ListSupportedLanguagesRequest value) => value.writeToBuffer(),
      $0.ListSupportedLanguagesResponse.fromBuffer);
  static final _$translateBatch =
      $grpc.ClientMethod<$0.TranslateBatchRequest, $0.TranslateBatchResponse>(
          '/manpasik.v1.TranslationService/TranslateBatch',
          ($0.TranslateBatchRequest value) => value.writeToBuffer(),
          $0.TranslateBatchResponse.fromBuffer);
  static final _$getTranslationHistory = $grpc.ClientMethod<
          $0.GetTranslationHistoryRequest, $0.GetTranslationHistoryResponse>(
      '/manpasik.v1.TranslationService/GetTranslationHistory',
      ($0.GetTranslationHistoryRequest value) => value.writeToBuffer(),
      $0.GetTranslationHistoryResponse.fromBuffer);
  static final _$getTranslationUsage = $grpc.ClientMethod<
          $0.GetTranslationUsageRequest, $0.GetTranslationUsageResponse>(
      '/manpasik.v1.TranslationService/GetTranslationUsage',
      ($0.GetTranslationUsageRequest value) => value.writeToBuffer(),
      $0.GetTranslationUsageResponse.fromBuffer);
  static final _$translateRealtime = $grpc.ClientMethod<
          $0.TranslateRealtimeRequest, $0.TranslateRealtimeResponse>(
      '/manpasik.v1.TranslationService/TranslateRealtime',
      ($0.TranslateRealtimeRequest value) => value.writeToBuffer(),
      $0.TranslateRealtimeResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.TranslationService')
abstract class TranslationServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.TranslationService';

  TranslationServiceBase() {
    $addMethod(
        $grpc.ServiceMethod<$0.TranslateTextRequest, $0.TranslateTextResponse>(
            'TranslateText',
            translateText_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.TranslateTextRequest.fromBuffer(value),
            ($0.TranslateTextResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DetectLanguageRequest,
            $0.DetectLanguageResponse>(
        'DetectLanguage',
        detectLanguage_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DetectLanguageRequest.fromBuffer(value),
        ($0.DetectLanguageResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListSupportedLanguagesRequest,
            $0.ListSupportedLanguagesResponse>(
        'ListSupportedLanguages',
        listSupportedLanguages_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListSupportedLanguagesRequest.fromBuffer(value),
        ($0.ListSupportedLanguagesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TranslateBatchRequest,
            $0.TranslateBatchResponse>(
        'TranslateBatch',
        translateBatch_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.TranslateBatchRequest.fromBuffer(value),
        ($0.TranslateBatchResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetTranslationHistoryRequest,
            $0.GetTranslationHistoryResponse>(
        'GetTranslationHistory',
        getTranslationHistory_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetTranslationHistoryRequest.fromBuffer(value),
        ($0.GetTranslationHistoryResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetTranslationUsageRequest,
            $0.GetTranslationUsageResponse>(
        'GetTranslationUsage',
        getTranslationUsage_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetTranslationUsageRequest.fromBuffer(value),
        ($0.GetTranslationUsageResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TranslateRealtimeRequest,
            $0.TranslateRealtimeResponse>(
        'TranslateRealtime',
        translateRealtime_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.TranslateRealtimeRequest.fromBuffer(value),
        ($0.TranslateRealtimeResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.TranslateTextResponse> translateText_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.TranslateTextRequest> $request) async {
    return translateText($call, await $request);
  }

  $async.Future<$0.TranslateTextResponse> translateText(
      $grpc.ServiceCall call, $0.TranslateTextRequest request);

  $async.Future<$0.DetectLanguageResponse> detectLanguage_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DetectLanguageRequest> $request) async {
    return detectLanguage($call, await $request);
  }

  $async.Future<$0.DetectLanguageResponse> detectLanguage(
      $grpc.ServiceCall call, $0.DetectLanguageRequest request);

  $async.Future<$0.ListSupportedLanguagesResponse> listSupportedLanguages_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListSupportedLanguagesRequest> $request) async {
    return listSupportedLanguages($call, await $request);
  }

  $async.Future<$0.ListSupportedLanguagesResponse> listSupportedLanguages(
      $grpc.ServiceCall call, $0.ListSupportedLanguagesRequest request);

  $async.Future<$0.TranslateBatchResponse> translateBatch_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.TranslateBatchRequest> $request) async {
    return translateBatch($call, await $request);
  }

  $async.Future<$0.TranslateBatchResponse> translateBatch(
      $grpc.ServiceCall call, $0.TranslateBatchRequest request);

  $async.Future<$0.GetTranslationHistoryResponse> getTranslationHistory_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetTranslationHistoryRequest> $request) async {
    return getTranslationHistory($call, await $request);
  }

  $async.Future<$0.GetTranslationHistoryResponse> getTranslationHistory(
      $grpc.ServiceCall call, $0.GetTranslationHistoryRequest request);

  $async.Future<$0.GetTranslationUsageResponse> getTranslationUsage_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetTranslationUsageRequest> $request) async {
    return getTranslationUsage($call, await $request);
  }

  $async.Future<$0.GetTranslationUsageResponse> getTranslationUsage(
      $grpc.ServiceCall call, $0.GetTranslationUsageRequest request);

  $async.Future<$0.TranslateRealtimeResponse> translateRealtime_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.TranslateRealtimeRequest> $request) async {
    return translateRealtime($call, await $request);
  }

  $async.Future<$0.TranslateRealtimeResponse> translateRealtime(
      $grpc.ServiceCall call, $0.TranslateRealtimeRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.TelemedicineService')
class TelemedicineServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  TelemedicineServiceClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.Consultation> createConsultation(
    $0.CreateConsultationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createConsultation, request, options: options);
  }

  $grpc.ResponseFuture<$0.Consultation> getConsultation(
    $0.GetConsultationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getConsultation, request, options: options);
  }

  $grpc.ResponseFuture<$0.ListConsultationsResponse> listConsultations(
    $0.ListConsultationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listConsultations, request, options: options);
  }

  $grpc.ResponseFuture<$0.MatchDoctorResponse> matchDoctor(
    $0.MatchDoctorRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$matchDoctor, request, options: options);
  }

  $grpc.ResponseFuture<$0.VideoSession> startVideoSession(
    $0.StartVideoSessionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startVideoSession, request, options: options);
  }

  $grpc.ResponseFuture<$0.VideoSession> endVideoSession(
    $0.EndVideoSessionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$endVideoSession, request, options: options);
  }

  $grpc.ResponseFuture<$0.RateConsultationResponse> rateConsultation(
    $0.RateConsultationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$rateConsultation, request, options: options);
  }

  // method descriptors

  static final _$createConsultation =
      $grpc.ClientMethod<$0.CreateConsultationRequest, $0.Consultation>(
          '/manpasik.v1.TelemedicineService/CreateConsultation',
          ($0.CreateConsultationRequest value) => value.writeToBuffer(),
          $0.Consultation.fromBuffer);
  static final _$getConsultation =
      $grpc.ClientMethod<$0.GetConsultationRequest, $0.Consultation>(
          '/manpasik.v1.TelemedicineService/GetConsultation',
          ($0.GetConsultationRequest value) => value.writeToBuffer(),
          $0.Consultation.fromBuffer);
  static final _$listConsultations = $grpc.ClientMethod<
          $0.ListConsultationsRequest, $0.ListConsultationsResponse>(
      '/manpasik.v1.TelemedicineService/ListConsultations',
      ($0.ListConsultationsRequest value) => value.writeToBuffer(),
      $0.ListConsultationsResponse.fromBuffer);
  static final _$matchDoctor =
      $grpc.ClientMethod<$0.MatchDoctorRequest, $0.MatchDoctorResponse>(
          '/manpasik.v1.TelemedicineService/MatchDoctor',
          ($0.MatchDoctorRequest value) => value.writeToBuffer(),
          $0.MatchDoctorResponse.fromBuffer);
  static final _$startVideoSession =
      $grpc.ClientMethod<$0.StartVideoSessionRequest, $0.VideoSession>(
          '/manpasik.v1.TelemedicineService/StartVideoSession',
          ($0.StartVideoSessionRequest value) => value.writeToBuffer(),
          $0.VideoSession.fromBuffer);
  static final _$endVideoSession =
      $grpc.ClientMethod<$0.EndVideoSessionRequest, $0.VideoSession>(
          '/manpasik.v1.TelemedicineService/EndVideoSession',
          ($0.EndVideoSessionRequest value) => value.writeToBuffer(),
          $0.VideoSession.fromBuffer);
  static final _$rateConsultation = $grpc.ClientMethod<
          $0.RateConsultationRequest, $0.RateConsultationResponse>(
      '/manpasik.v1.TelemedicineService/RateConsultation',
      ($0.RateConsultationRequest value) => value.writeToBuffer(),
      $0.RateConsultationResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.TelemedicineService')
abstract class TelemedicineServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.TelemedicineService';

  TelemedicineServiceBase() {
    $addMethod(
        $grpc.ServiceMethod<$0.CreateConsultationRequest, $0.Consultation>(
            'CreateConsultation',
            createConsultation_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateConsultationRequest.fromBuffer(value),
            ($0.Consultation value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetConsultationRequest, $0.Consultation>(
        'GetConsultation',
        getConsultation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetConsultationRequest.fromBuffer(value),
        ($0.Consultation value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListConsultationsRequest,
            $0.ListConsultationsResponse>(
        'ListConsultations',
        listConsultations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListConsultationsRequest.fromBuffer(value),
        ($0.ListConsultationsResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.MatchDoctorRequest, $0.MatchDoctorResponse>(
            'MatchDoctor',
            matchDoctor_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.MatchDoctorRequest.fromBuffer(value),
            ($0.MatchDoctorResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.StartVideoSessionRequest, $0.VideoSession>(
            'StartVideoSession',
            startVideoSession_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.StartVideoSessionRequest.fromBuffer(value),
            ($0.VideoSession value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.EndVideoSessionRequest, $0.VideoSession>(
        'EndVideoSession',
        endVideoSession_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.EndVideoSessionRequest.fromBuffer(value),
        ($0.VideoSession value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RateConsultationRequest,
            $0.RateConsultationResponse>(
        'RateConsultation',
        rateConsultation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RateConsultationRequest.fromBuffer(value),
        ($0.RateConsultationResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.Consultation> createConsultation_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateConsultationRequest> $request) async {
    return createConsultation($call, await $request);
  }

  $async.Future<$0.Consultation> createConsultation(
      $grpc.ServiceCall call, $0.CreateConsultationRequest request);

  $async.Future<$0.Consultation> getConsultation_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetConsultationRequest> $request) async {
    return getConsultation($call, await $request);
  }

  $async.Future<$0.Consultation> getConsultation(
      $grpc.ServiceCall call, $0.GetConsultationRequest request);

  $async.Future<$0.ListConsultationsResponse> listConsultations_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListConsultationsRequest> $request) async {
    return listConsultations($call, await $request);
  }

  $async.Future<$0.ListConsultationsResponse> listConsultations(
      $grpc.ServiceCall call, $0.ListConsultationsRequest request);

  $async.Future<$0.MatchDoctorResponse> matchDoctor_Pre($grpc.ServiceCall $call,
      $async.Future<$0.MatchDoctorRequest> $request) async {
    return matchDoctor($call, await $request);
  }

  $async.Future<$0.MatchDoctorResponse> matchDoctor(
      $grpc.ServiceCall call, $0.MatchDoctorRequest request);

  $async.Future<$0.VideoSession> startVideoSession_Pre($grpc.ServiceCall $call,
      $async.Future<$0.StartVideoSessionRequest> $request) async {
    return startVideoSession($call, await $request);
  }

  $async.Future<$0.VideoSession> startVideoSession(
      $grpc.ServiceCall call, $0.StartVideoSessionRequest request);

  $async.Future<$0.VideoSession> endVideoSession_Pre($grpc.ServiceCall $call,
      $async.Future<$0.EndVideoSessionRequest> $request) async {
    return endVideoSession($call, await $request);
  }

  $async.Future<$0.VideoSession> endVideoSession(
      $grpc.ServiceCall call, $0.EndVideoSessionRequest request);

  $async.Future<$0.RateConsultationResponse> rateConsultation_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RateConsultationRequest> $request) async {
    return rateConsultation($call, await $request);
  }

  $async.Future<$0.RateConsultationResponse> rateConsultation(
      $grpc.ServiceCall call, $0.RateConsultationRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.AssistantService')
class AssistantServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  AssistantServiceClient(super.channel, {super.options, super.interceptors});

  /// 텍스트 명령 전송 → 의도 파싱 → 액션 실행 → 결과 반환
  $grpc.ResponseFuture<$0.AssistantCommandResponse> sendCommand(
    $0.AssistantCommandRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$sendCommand, request, options: options);
  }

  /// 세션 조회
  $grpc.ResponseFuture<$0.AssistantSession> getSession(
    $0.GetAssistantSessionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSession, request, options: options);
  }

  /// 세션 목록
  $grpc.ResponseFuture<$0.ListAssistantSessionsResponse> listSessions(
    $0.ListAssistantSessionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listSessions, request, options: options);
  }

  /// 세션 내 턴(대화) 목록
  $grpc.ResponseFuture<$0.ListAssistantTurnsResponse> listTurns(
    $0.ListAssistantTurnsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTurns, request, options: options);
  }

  /// 세션 삭제
  $grpc.ResponseFuture<$0.DeleteAssistantSessionResponse> deleteSession(
    $0.DeleteAssistantSessionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteSession, request, options: options);
  }

  // method descriptors

  static final _$sendCommand = $grpc.ClientMethod<$0.AssistantCommandRequest,
          $0.AssistantCommandResponse>(
      '/manpasik.v1.AssistantService/SendCommand',
      ($0.AssistantCommandRequest value) => value.writeToBuffer(),
      $0.AssistantCommandResponse.fromBuffer);
  static final _$getSession =
      $grpc.ClientMethod<$0.GetAssistantSessionRequest, $0.AssistantSession>(
          '/manpasik.v1.AssistantService/GetSession',
          ($0.GetAssistantSessionRequest value) => value.writeToBuffer(),
          $0.AssistantSession.fromBuffer);
  static final _$listSessions = $grpc.ClientMethod<
          $0.ListAssistantSessionsRequest, $0.ListAssistantSessionsResponse>(
      '/manpasik.v1.AssistantService/ListSessions',
      ($0.ListAssistantSessionsRequest value) => value.writeToBuffer(),
      $0.ListAssistantSessionsResponse.fromBuffer);
  static final _$listTurns = $grpc.ClientMethod<$0.ListAssistantTurnsRequest,
          $0.ListAssistantTurnsResponse>(
      '/manpasik.v1.AssistantService/ListTurns',
      ($0.ListAssistantTurnsRequest value) => value.writeToBuffer(),
      $0.ListAssistantTurnsResponse.fromBuffer);
  static final _$deleteSession = $grpc.ClientMethod<
          $0.DeleteAssistantSessionRequest, $0.DeleteAssistantSessionResponse>(
      '/manpasik.v1.AssistantService/DeleteSession',
      ($0.DeleteAssistantSessionRequest value) => value.writeToBuffer(),
      $0.DeleteAssistantSessionResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.AssistantService')
abstract class AssistantServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.AssistantService';

  AssistantServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.AssistantCommandRequest,
            $0.AssistantCommandResponse>(
        'SendCommand',
        sendCommand_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AssistantCommandRequest.fromBuffer(value),
        ($0.AssistantCommandResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetAssistantSessionRequest, $0.AssistantSession>(
            'GetSession',
            getSession_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetAssistantSessionRequest.fromBuffer(value),
            ($0.AssistantSession value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListAssistantSessionsRequest,
            $0.ListAssistantSessionsResponse>(
        'ListSessions',
        listSessions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListAssistantSessionsRequest.fromBuffer(value),
        ($0.ListAssistantSessionsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListAssistantTurnsRequest,
            $0.ListAssistantTurnsResponse>(
        'ListTurns',
        listTurns_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListAssistantTurnsRequest.fromBuffer(value),
        ($0.ListAssistantTurnsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteAssistantSessionRequest,
            $0.DeleteAssistantSessionResponse>(
        'DeleteSession',
        deleteSession_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteAssistantSessionRequest.fromBuffer(value),
        ($0.DeleteAssistantSessionResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.AssistantCommandResponse> sendCommand_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AssistantCommandRequest> $request) async {
    return sendCommand($call, await $request);
  }

  $async.Future<$0.AssistantCommandResponse> sendCommand(
      $grpc.ServiceCall call, $0.AssistantCommandRequest request);

  $async.Future<$0.AssistantSession> getSession_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetAssistantSessionRequest> $request) async {
    return getSession($call, await $request);
  }

  $async.Future<$0.AssistantSession> getSession(
      $grpc.ServiceCall call, $0.GetAssistantSessionRequest request);

  $async.Future<$0.ListAssistantSessionsResponse> listSessions_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListAssistantSessionsRequest> $request) async {
    return listSessions($call, await $request);
  }

  $async.Future<$0.ListAssistantSessionsResponse> listSessions(
      $grpc.ServiceCall call, $0.ListAssistantSessionsRequest request);

  $async.Future<$0.ListAssistantTurnsResponse> listTurns_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListAssistantTurnsRequest> $request) async {
    return listTurns($call, await $request);
  }

  $async.Future<$0.ListAssistantTurnsResponse> listTurns(
      $grpc.ServiceCall call, $0.ListAssistantTurnsRequest request);

  $async.Future<$0.DeleteAssistantSessionResponse> deleteSession_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteAssistantSessionRequest> $request) async {
    return deleteSession($call, await $request);
  }

  $async.Future<$0.DeleteAssistantSessionResponse> deleteSession(
      $grpc.ServiceCall call, $0.DeleteAssistantSessionRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.VisionService')
class VisionServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  VisionServiceClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.FoodAnalysisResult> analyzeFood(
    $0.AnalyzeFoodRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$analyzeFood, request, options: options);
  }

  $grpc.ResponseFuture<$0.FoodAnalysisResult> getFoodAnalysis(
    $0.GetFoodAnalysisRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getFoodAnalysis, request, options: options);
  }

  $grpc.ResponseFuture<$0.ListFoodAnalysesResponse> listFoodAnalyses(
    $0.ListFoodAnalysesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listFoodAnalyses, request, options: options);
  }

  $grpc.ResponseFuture<$0.DailyNutritionSummary> getDailyNutritionSummary(
    $0.GetDailyNutritionSummaryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDailyNutritionSummary, request,
        options: options);
  }

  $grpc.ResponseFuture<$0.MealLog> logMeal(
    $0.LogMealRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$logMeal, request, options: options);
  }

  $grpc.ResponseFuture<$0.GetMealHistoryResponse> getMealHistory(
    $0.GetMealHistoryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getMealHistory, request, options: options);
  }

  // method descriptors

  static final _$analyzeFood =
      $grpc.ClientMethod<$0.AnalyzeFoodRequest, $0.FoodAnalysisResult>(
          '/manpasik.v1.VisionService/AnalyzeFood',
          ($0.AnalyzeFoodRequest value) => value.writeToBuffer(),
          $0.FoodAnalysisResult.fromBuffer);
  static final _$getFoodAnalysis =
      $grpc.ClientMethod<$0.GetFoodAnalysisRequest, $0.FoodAnalysisResult>(
          '/manpasik.v1.VisionService/GetFoodAnalysis',
          ($0.GetFoodAnalysisRequest value) => value.writeToBuffer(),
          $0.FoodAnalysisResult.fromBuffer);
  static final _$listFoodAnalyses = $grpc.ClientMethod<
          $0.ListFoodAnalysesRequest, $0.ListFoodAnalysesResponse>(
      '/manpasik.v1.VisionService/ListFoodAnalyses',
      ($0.ListFoodAnalysesRequest value) => value.writeToBuffer(),
      $0.ListFoodAnalysesResponse.fromBuffer);
  static final _$getDailyNutritionSummary = $grpc.ClientMethod<
          $0.GetDailyNutritionSummaryRequest, $0.DailyNutritionSummary>(
      '/manpasik.v1.VisionService/GetDailyNutritionSummary',
      ($0.GetDailyNutritionSummaryRequest value) => value.writeToBuffer(),
      $0.DailyNutritionSummary.fromBuffer);
  static final _$logMeal = $grpc.ClientMethod<$0.LogMealRequest, $0.MealLog>(
      '/manpasik.v1.VisionService/LogMeal',
      ($0.LogMealRequest value) => value.writeToBuffer(),
      $0.MealLog.fromBuffer);
  static final _$getMealHistory =
      $grpc.ClientMethod<$0.GetMealHistoryRequest, $0.GetMealHistoryResponse>(
          '/manpasik.v1.VisionService/GetMealHistory',
          ($0.GetMealHistoryRequest value) => value.writeToBuffer(),
          $0.GetMealHistoryResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.VisionService')
abstract class VisionServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.VisionService';

  VisionServiceBase() {
    $addMethod(
        $grpc.ServiceMethod<$0.AnalyzeFoodRequest, $0.FoodAnalysisResult>(
            'AnalyzeFood',
            analyzeFood_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.AnalyzeFoodRequest.fromBuffer(value),
            ($0.FoodAnalysisResult value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetFoodAnalysisRequest, $0.FoodAnalysisResult>(
            'GetFoodAnalysis',
            getFoodAnalysis_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetFoodAnalysisRequest.fromBuffer(value),
            ($0.FoodAnalysisResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListFoodAnalysesRequest,
            $0.ListFoodAnalysesResponse>(
        'ListFoodAnalyses',
        listFoodAnalyses_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListFoodAnalysesRequest.fromBuffer(value),
        ($0.ListFoodAnalysesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDailyNutritionSummaryRequest,
            $0.DailyNutritionSummary>(
        'GetDailyNutritionSummary',
        getDailyNutritionSummary_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDailyNutritionSummaryRequest.fromBuffer(value),
        ($0.DailyNutritionSummary value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.LogMealRequest, $0.MealLog>(
        'LogMeal',
        logMeal_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.LogMealRequest.fromBuffer(value),
        ($0.MealLog value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetMealHistoryRequest,
            $0.GetMealHistoryResponse>(
        'GetMealHistory',
        getMealHistory_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetMealHistoryRequest.fromBuffer(value),
        ($0.GetMealHistoryResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.FoodAnalysisResult> analyzeFood_Pre($grpc.ServiceCall $call,
      $async.Future<$0.AnalyzeFoodRequest> $request) async {
    return analyzeFood($call, await $request);
  }

  $async.Future<$0.FoodAnalysisResult> analyzeFood(
      $grpc.ServiceCall call, $0.AnalyzeFoodRequest request);

  $async.Future<$0.FoodAnalysisResult> getFoodAnalysis_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetFoodAnalysisRequest> $request) async {
    return getFoodAnalysis($call, await $request);
  }

  $async.Future<$0.FoodAnalysisResult> getFoodAnalysis(
      $grpc.ServiceCall call, $0.GetFoodAnalysisRequest request);

  $async.Future<$0.ListFoodAnalysesResponse> listFoodAnalyses_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListFoodAnalysesRequest> $request) async {
    return listFoodAnalyses($call, await $request);
  }

  $async.Future<$0.ListFoodAnalysesResponse> listFoodAnalyses(
      $grpc.ServiceCall call, $0.ListFoodAnalysesRequest request);

  $async.Future<$0.DailyNutritionSummary> getDailyNutritionSummary_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDailyNutritionSummaryRequest> $request) async {
    return getDailyNutritionSummary($call, await $request);
  }

  $async.Future<$0.DailyNutritionSummary> getDailyNutritionSummary(
      $grpc.ServiceCall call, $0.GetDailyNutritionSummaryRequest request);

  $async.Future<$0.MealLog> logMeal_Pre($grpc.ServiceCall $call,
      $async.Future<$0.LogMealRequest> $request) async {
    return logMeal($call, await $request);
  }

  $async.Future<$0.MealLog> logMeal(
      $grpc.ServiceCall call, $0.LogMealRequest request);

  $async.Future<$0.GetMealHistoryResponse> getMealHistory_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetMealHistoryRequest> $request) async {
    return getMealHistory($call, await $request);
  }

  $async.Future<$0.GetMealHistoryResponse> getMealHistory(
      $grpc.ServiceCall call, $0.GetMealHistoryRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.ConceptService')
class ConceptServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  ConceptServiceClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.ListConceptsResponse> listConcepts(
    $0.ListConceptsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listConcepts, request, options: options);
  }

  $grpc.ResponseFuture<$0.Concept> getConcept(
    $0.GetConceptRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getConcept, request, options: options);
  }

  $grpc.ResponseFuture<$0.Concept> createConcept(
    $0.CreateConceptRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createConcept, request, options: options);
  }

  $grpc.ResponseFuture<$0.AssignDeviceToConceptResponse> assignDeviceToConcept(
    $0.AssignDeviceToConceptRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$assignDeviceToConcept, request, options: options);
  }

  $grpc.ResponseFuture<$0.ConceptStats> getConceptStats(
    $0.GetConceptStatsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getConceptStats, request, options: options);
  }

  $grpc.ResponseFuture<$0.ConceptDashboard> getConceptDashboard(
    $0.GetConceptDashboardRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getConceptDashboard, request, options: options);
  }

  // method descriptors

  static final _$listConcepts =
      $grpc.ClientMethod<$0.ListConceptsRequest, $0.ListConceptsResponse>(
          '/manpasik.v1.ConceptService/ListConcepts',
          ($0.ListConceptsRequest value) => value.writeToBuffer(),
          $0.ListConceptsResponse.fromBuffer);
  static final _$getConcept =
      $grpc.ClientMethod<$0.GetConceptRequest, $0.Concept>(
          '/manpasik.v1.ConceptService/GetConcept',
          ($0.GetConceptRequest value) => value.writeToBuffer(),
          $0.Concept.fromBuffer);
  static final _$createConcept =
      $grpc.ClientMethod<$0.CreateConceptRequest, $0.Concept>(
          '/manpasik.v1.ConceptService/CreateConcept',
          ($0.CreateConceptRequest value) => value.writeToBuffer(),
          $0.Concept.fromBuffer);
  static final _$assignDeviceToConcept = $grpc.ClientMethod<
          $0.AssignDeviceToConceptRequest, $0.AssignDeviceToConceptResponse>(
      '/manpasik.v1.ConceptService/AssignDeviceToConcept',
      ($0.AssignDeviceToConceptRequest value) => value.writeToBuffer(),
      $0.AssignDeviceToConceptResponse.fromBuffer);
  static final _$getConceptStats =
      $grpc.ClientMethod<$0.GetConceptStatsRequest, $0.ConceptStats>(
          '/manpasik.v1.ConceptService/GetConceptStats',
          ($0.GetConceptStatsRequest value) => value.writeToBuffer(),
          $0.ConceptStats.fromBuffer);
  static final _$getConceptDashboard =
      $grpc.ClientMethod<$0.GetConceptDashboardRequest, $0.ConceptDashboard>(
          '/manpasik.v1.ConceptService/GetConceptDashboard',
          ($0.GetConceptDashboardRequest value) => value.writeToBuffer(),
          $0.ConceptDashboard.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.ConceptService')
abstract class ConceptServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.ConceptService';

  ConceptServiceBase() {
    $addMethod(
        $grpc.ServiceMethod<$0.ListConceptsRequest, $0.ListConceptsResponse>(
            'ListConcepts',
            listConcepts_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListConceptsRequest.fromBuffer(value),
            ($0.ListConceptsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetConceptRequest, $0.Concept>(
        'GetConcept',
        getConcept_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetConceptRequest.fromBuffer(value),
        ($0.Concept value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateConceptRequest, $0.Concept>(
        'CreateConcept',
        createConcept_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateConceptRequest.fromBuffer(value),
        ($0.Concept value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AssignDeviceToConceptRequest,
            $0.AssignDeviceToConceptResponse>(
        'AssignDeviceToConcept',
        assignDeviceToConcept_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AssignDeviceToConceptRequest.fromBuffer(value),
        ($0.AssignDeviceToConceptResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetConceptStatsRequest, $0.ConceptStats>(
        'GetConceptStats',
        getConceptStats_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetConceptStatsRequest.fromBuffer(value),
        ($0.ConceptStats value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetConceptDashboardRequest, $0.ConceptDashboard>(
            'GetConceptDashboard',
            getConceptDashboard_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetConceptDashboardRequest.fromBuffer(value),
            ($0.ConceptDashboard value) => value.writeToBuffer()));
  }

  $async.Future<$0.ListConceptsResponse> listConcepts_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListConceptsRequest> $request) async {
    return listConcepts($call, await $request);
  }

  $async.Future<$0.ListConceptsResponse> listConcepts(
      $grpc.ServiceCall call, $0.ListConceptsRequest request);

  $async.Future<$0.Concept> getConcept_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetConceptRequest> $request) async {
    return getConcept($call, await $request);
  }

  $async.Future<$0.Concept> getConcept(
      $grpc.ServiceCall call, $0.GetConceptRequest request);

  $async.Future<$0.Concept> createConcept_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateConceptRequest> $request) async {
    return createConcept($call, await $request);
  }

  $async.Future<$0.Concept> createConcept(
      $grpc.ServiceCall call, $0.CreateConceptRequest request);

  $async.Future<$0.AssignDeviceToConceptResponse> assignDeviceToConcept_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AssignDeviceToConceptRequest> $request) async {
    return assignDeviceToConcept($call, await $request);
  }

  $async.Future<$0.AssignDeviceToConceptResponse> assignDeviceToConcept(
      $grpc.ServiceCall call, $0.AssignDeviceToConceptRequest request);

  $async.Future<$0.ConceptStats> getConceptStats_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetConceptStatsRequest> $request) async {
    return getConceptStats($call, await $request);
  }

  $async.Future<$0.ConceptStats> getConceptStats(
      $grpc.ServiceCall call, $0.GetConceptStatsRequest request);

  $async.Future<$0.ConceptDashboard> getConceptDashboard_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetConceptDashboardRequest> $request) async {
    return getConceptDashboard($call, await $request);
  }

  $async.Future<$0.ConceptDashboard> getConceptDashboard(
      $grpc.ServiceCall call, $0.GetConceptDashboardRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.OrganizationService')
class OrganizationServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  OrganizationServiceClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.Organization> createOrganization(
    $0.CreateOrganizationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createOrganization, request, options: options);
  }

  $grpc.ResponseFuture<$0.Organization> getOrganization(
    $0.GetOrganizationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getOrganization, request, options: options);
  }

  $grpc.ResponseFuture<$0.ListOrganizationsResponse> listOrganizations(
    $0.ListOrganizationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listOrganizations, request, options: options);
  }

  $grpc.ResponseFuture<$0.AddOrgMemberResponse> addMember(
    $0.AddOrgMemberRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$addMember, request, options: options);
  }

  $grpc.ResponseFuture<$0.RemoveOrgMemberResponse> removeMember(
    $0.RemoveOrgMemberRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$removeMember, request, options: options);
  }

  // method descriptors

  static final _$createOrganization =
      $grpc.ClientMethod<$0.CreateOrganizationRequest, $0.Organization>(
          '/manpasik.v1.OrganizationService/CreateOrganization',
          ($0.CreateOrganizationRequest value) => value.writeToBuffer(),
          $0.Organization.fromBuffer);
  static final _$getOrganization =
      $grpc.ClientMethod<$0.GetOrganizationRequest, $0.Organization>(
          '/manpasik.v1.OrganizationService/GetOrganization',
          ($0.GetOrganizationRequest value) => value.writeToBuffer(),
          $0.Organization.fromBuffer);
  static final _$listOrganizations = $grpc.ClientMethod<
          $0.ListOrganizationsRequest, $0.ListOrganizationsResponse>(
      '/manpasik.v1.OrganizationService/ListOrganizations',
      ($0.ListOrganizationsRequest value) => value.writeToBuffer(),
      $0.ListOrganizationsResponse.fromBuffer);
  static final _$addMember =
      $grpc.ClientMethod<$0.AddOrgMemberRequest, $0.AddOrgMemberResponse>(
          '/manpasik.v1.OrganizationService/AddMember',
          ($0.AddOrgMemberRequest value) => value.writeToBuffer(),
          $0.AddOrgMemberResponse.fromBuffer);
  static final _$removeMember =
      $grpc.ClientMethod<$0.RemoveOrgMemberRequest, $0.RemoveOrgMemberResponse>(
          '/manpasik.v1.OrganizationService/RemoveMember',
          ($0.RemoveOrgMemberRequest value) => value.writeToBuffer(),
          $0.RemoveOrgMemberResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.OrganizationService')
abstract class OrganizationServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.OrganizationService';

  OrganizationServiceBase() {
    $addMethod(
        $grpc.ServiceMethod<$0.CreateOrganizationRequest, $0.Organization>(
            'CreateOrganization',
            createOrganization_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateOrganizationRequest.fromBuffer(value),
            ($0.Organization value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetOrganizationRequest, $0.Organization>(
        'GetOrganization',
        getOrganization_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetOrganizationRequest.fromBuffer(value),
        ($0.Organization value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListOrganizationsRequest,
            $0.ListOrganizationsResponse>(
        'ListOrganizations',
        listOrganizations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListOrganizationsRequest.fromBuffer(value),
        ($0.ListOrganizationsResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.AddOrgMemberRequest, $0.AddOrgMemberResponse>(
            'AddMember',
            addMember_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.AddOrgMemberRequest.fromBuffer(value),
            ($0.AddOrgMemberResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RemoveOrgMemberRequest,
            $0.RemoveOrgMemberResponse>(
        'RemoveMember',
        removeMember_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RemoveOrgMemberRequest.fromBuffer(value),
        ($0.RemoveOrgMemberResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.Organization> createOrganization_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateOrganizationRequest> $request) async {
    return createOrganization($call, await $request);
  }

  $async.Future<$0.Organization> createOrganization(
      $grpc.ServiceCall call, $0.CreateOrganizationRequest request);

  $async.Future<$0.Organization> getOrganization_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetOrganizationRequest> $request) async {
    return getOrganization($call, await $request);
  }

  $async.Future<$0.Organization> getOrganization(
      $grpc.ServiceCall call, $0.GetOrganizationRequest request);

  $async.Future<$0.ListOrganizationsResponse> listOrganizations_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListOrganizationsRequest> $request) async {
    return listOrganizations($call, await $request);
  }

  $async.Future<$0.ListOrganizationsResponse> listOrganizations(
      $grpc.ServiceCall call, $0.ListOrganizationsRequest request);

  $async.Future<$0.AddOrgMemberResponse> addMember_Pre($grpc.ServiceCall $call,
      $async.Future<$0.AddOrgMemberRequest> $request) async {
    return addMember($call, await $request);
  }

  $async.Future<$0.AddOrgMemberResponse> addMember(
      $grpc.ServiceCall call, $0.AddOrgMemberRequest request);

  $async.Future<$0.RemoveOrgMemberResponse> removeMember_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RemoveOrgMemberRequest> $request) async {
    return removeMember($call, await $request);
  }

  $async.Future<$0.RemoveOrgMemberResponse> removeMember(
      $grpc.ServiceCall call, $0.RemoveOrgMemberRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.DeveloperService')
class DeveloperServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  DeveloperServiceClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.DeveloperProfile> registerDeveloper(
    $0.RegisterDeveloperRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$registerDeveloper, request, options: options);
  }

  $grpc.ResponseFuture<$0.DeveloperProfile> getDeveloperProfile(
    $0.GetDeveloperProfileRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDeveloperProfile, request, options: options);
  }

  $grpc.ResponseFuture<$0.ApiKeyResponse> createApiKey(
    $0.CreateApiKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createApiKey, request, options: options);
  }

  $grpc.ResponseFuture<$0.SubmitCartridgeResponse> submitCartridge(
    $0.SubmitCartridgeRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$submitCartridge, request, options: options);
  }

  $grpc.ResponseFuture<$0.ListSubmissionsResponse> listSubmissions(
    $0.ListSubmissionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listSubmissions, request, options: options);
  }

  /// Phase 5: 카트리지 타입 등록 (SDK 개발자용 하네스 인터페이스)
  $grpc.ResponseFuture<$0.RegisterCartridgeTypeResponse> registerCartridgeType(
    $0.RegisterCartridgeTypeRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$registerCartridgeType, request, options: options);
  }

  /// Phase 5: 카트리지 정의 조회
  $grpc.ResponseFuture<$0.CartridgeDefinitionDetail> getCartridgeDefinition(
    $0.GetCartridgeDefinitionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getCartridgeDefinition, request,
        options: options);
  }

  // method descriptors

  static final _$registerDeveloper =
      $grpc.ClientMethod<$0.RegisterDeveloperRequest, $0.DeveloperProfile>(
          '/manpasik.v1.DeveloperService/RegisterDeveloper',
          ($0.RegisterDeveloperRequest value) => value.writeToBuffer(),
          $0.DeveloperProfile.fromBuffer);
  static final _$getDeveloperProfile =
      $grpc.ClientMethod<$0.GetDeveloperProfileRequest, $0.DeveloperProfile>(
          '/manpasik.v1.DeveloperService/GetDeveloperProfile',
          ($0.GetDeveloperProfileRequest value) => value.writeToBuffer(),
          $0.DeveloperProfile.fromBuffer);
  static final _$createApiKey =
      $grpc.ClientMethod<$0.CreateApiKeyRequest, $0.ApiKeyResponse>(
          '/manpasik.v1.DeveloperService/CreateApiKey',
          ($0.CreateApiKeyRequest value) => value.writeToBuffer(),
          $0.ApiKeyResponse.fromBuffer);
  static final _$submitCartridge =
      $grpc.ClientMethod<$0.SubmitCartridgeRequest, $0.SubmitCartridgeResponse>(
          '/manpasik.v1.DeveloperService/SubmitCartridge',
          ($0.SubmitCartridgeRequest value) => value.writeToBuffer(),
          $0.SubmitCartridgeResponse.fromBuffer);
  static final _$listSubmissions =
      $grpc.ClientMethod<$0.ListSubmissionsRequest, $0.ListSubmissionsResponse>(
          '/manpasik.v1.DeveloperService/ListSubmissions',
          ($0.ListSubmissionsRequest value) => value.writeToBuffer(),
          $0.ListSubmissionsResponse.fromBuffer);
  static final _$registerCartridgeType = $grpc.ClientMethod<
          $0.RegisterCartridgeTypeRequest, $0.RegisterCartridgeTypeResponse>(
      '/manpasik.v1.DeveloperService/RegisterCartridgeType',
      ($0.RegisterCartridgeTypeRequest value) => value.writeToBuffer(),
      $0.RegisterCartridgeTypeResponse.fromBuffer);
  static final _$getCartridgeDefinition = $grpc.ClientMethod<
          $0.GetCartridgeDefinitionRequest, $0.CartridgeDefinitionDetail>(
      '/manpasik.v1.DeveloperService/GetCartridgeDefinition',
      ($0.GetCartridgeDefinitionRequest value) => value.writeToBuffer(),
      $0.CartridgeDefinitionDetail.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.DeveloperService')
abstract class DeveloperServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.DeveloperService';

  DeveloperServiceBase() {
    $addMethod(
        $grpc.ServiceMethod<$0.RegisterDeveloperRequest, $0.DeveloperProfile>(
            'RegisterDeveloper',
            registerDeveloper_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.RegisterDeveloperRequest.fromBuffer(value),
            ($0.DeveloperProfile value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetDeveloperProfileRequest, $0.DeveloperProfile>(
            'GetDeveloperProfile',
            getDeveloperProfile_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetDeveloperProfileRequest.fromBuffer(value),
            ($0.DeveloperProfile value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateApiKeyRequest, $0.ApiKeyResponse>(
        'CreateApiKey',
        createApiKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateApiKeyRequest.fromBuffer(value),
        ($0.ApiKeyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SubmitCartridgeRequest,
            $0.SubmitCartridgeResponse>(
        'SubmitCartridge',
        submitCartridge_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SubmitCartridgeRequest.fromBuffer(value),
        ($0.SubmitCartridgeResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListSubmissionsRequest,
            $0.ListSubmissionsResponse>(
        'ListSubmissions',
        listSubmissions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListSubmissionsRequest.fromBuffer(value),
        ($0.ListSubmissionsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RegisterCartridgeTypeRequest,
            $0.RegisterCartridgeTypeResponse>(
        'RegisterCartridgeType',
        registerCartridgeType_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RegisterCartridgeTypeRequest.fromBuffer(value),
        ($0.RegisterCartridgeTypeResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetCartridgeDefinitionRequest,
            $0.CartridgeDefinitionDetail>(
        'GetCartridgeDefinition',
        getCartridgeDefinition_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetCartridgeDefinitionRequest.fromBuffer(value),
        ($0.CartridgeDefinitionDetail value) => value.writeToBuffer()));
  }

  $async.Future<$0.DeveloperProfile> registerDeveloper_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RegisterDeveloperRequest> $request) async {
    return registerDeveloper($call, await $request);
  }

  $async.Future<$0.DeveloperProfile> registerDeveloper(
      $grpc.ServiceCall call, $0.RegisterDeveloperRequest request);

  $async.Future<$0.DeveloperProfile> getDeveloperProfile_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDeveloperProfileRequest> $request) async {
    return getDeveloperProfile($call, await $request);
  }

  $async.Future<$0.DeveloperProfile> getDeveloperProfile(
      $grpc.ServiceCall call, $0.GetDeveloperProfileRequest request);

  $async.Future<$0.ApiKeyResponse> createApiKey_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateApiKeyRequest> $request) async {
    return createApiKey($call, await $request);
  }

  $async.Future<$0.ApiKeyResponse> createApiKey(
      $grpc.ServiceCall call, $0.CreateApiKeyRequest request);

  $async.Future<$0.SubmitCartridgeResponse> submitCartridge_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SubmitCartridgeRequest> $request) async {
    return submitCartridge($call, await $request);
  }

  $async.Future<$0.SubmitCartridgeResponse> submitCartridge(
      $grpc.ServiceCall call, $0.SubmitCartridgeRequest request);

  $async.Future<$0.ListSubmissionsResponse> listSubmissions_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListSubmissionsRequest> $request) async {
    return listSubmissions($call, await $request);
  }

  $async.Future<$0.ListSubmissionsResponse> listSubmissions(
      $grpc.ServiceCall call, $0.ListSubmissionsRequest request);

  $async.Future<$0.RegisterCartridgeTypeResponse> registerCartridgeType_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RegisterCartridgeTypeRequest> $request) async {
    return registerCartridgeType($call, await $request);
  }

  $async.Future<$0.RegisterCartridgeTypeResponse> registerCartridgeType(
      $grpc.ServiceCall call, $0.RegisterCartridgeTypeRequest request);

  $async.Future<$0.CartridgeDefinitionDetail> getCartridgeDefinition_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetCartridgeDefinitionRequest> $request) async {
    return getCartridgeDefinition($call, await $request);
  }

  $async.Future<$0.CartridgeDefinitionDetail> getCartridgeDefinition(
      $grpc.ServiceCall call, $0.GetCartridgeDefinitionRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.StoreService')
class StoreServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  StoreServiceClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.ListStoreItemsResponse> listStoreItems(
    $0.ListStoreItemsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listStoreItems, request, options: options);
  }

  $grpc.ResponseFuture<$0.SearchCartridgesResponse> searchCartridges(
    $0.SearchCartridgesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$searchCartridges, request, options: options);
  }

  $grpc.ResponseFuture<$0.StoreItem> getStoreItem(
    $0.GetStoreItemRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getStoreItem, request, options: options);
  }

  $grpc.ResponseFuture<$0.PurchaseCartridgeResponse> purchaseCartridge(
    $0.PurchaseCartridgeRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$purchaseCartridge, request, options: options);
  }

  $grpc.ResponseFuture<$0.GetPurchaseHistoryResponse> getPurchaseHistory(
    $0.GetPurchaseHistoryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getPurchaseHistory, request, options: options);
  }

  // method descriptors

  static final _$listStoreItems =
      $grpc.ClientMethod<$0.ListStoreItemsRequest, $0.ListStoreItemsResponse>(
          '/manpasik.v1.StoreService/ListStoreItems',
          ($0.ListStoreItemsRequest value) => value.writeToBuffer(),
          $0.ListStoreItemsResponse.fromBuffer);
  static final _$searchCartridges = $grpc.ClientMethod<
          $0.SearchCartridgesRequest, $0.SearchCartridgesResponse>(
      '/manpasik.v1.StoreService/SearchCartridges',
      ($0.SearchCartridgesRequest value) => value.writeToBuffer(),
      $0.SearchCartridgesResponse.fromBuffer);
  static final _$getStoreItem =
      $grpc.ClientMethod<$0.GetStoreItemRequest, $0.StoreItem>(
          '/manpasik.v1.StoreService/GetStoreItem',
          ($0.GetStoreItemRequest value) => value.writeToBuffer(),
          $0.StoreItem.fromBuffer);
  static final _$purchaseCartridge = $grpc.ClientMethod<
          $0.PurchaseCartridgeRequest, $0.PurchaseCartridgeResponse>(
      '/manpasik.v1.StoreService/PurchaseCartridge',
      ($0.PurchaseCartridgeRequest value) => value.writeToBuffer(),
      $0.PurchaseCartridgeResponse.fromBuffer);
  static final _$getPurchaseHistory = $grpc.ClientMethod<
          $0.GetPurchaseHistoryRequest, $0.GetPurchaseHistoryResponse>(
      '/manpasik.v1.StoreService/GetPurchaseHistory',
      ($0.GetPurchaseHistoryRequest value) => value.writeToBuffer(),
      $0.GetPurchaseHistoryResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.StoreService')
abstract class StoreServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.StoreService';

  StoreServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.ListStoreItemsRequest,
            $0.ListStoreItemsResponse>(
        'ListStoreItems',
        listStoreItems_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListStoreItemsRequest.fromBuffer(value),
        ($0.ListStoreItemsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SearchCartridgesRequest,
            $0.SearchCartridgesResponse>(
        'SearchCartridges',
        searchCartridges_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SearchCartridgesRequest.fromBuffer(value),
        ($0.SearchCartridgesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetStoreItemRequest, $0.StoreItem>(
        'GetStoreItem',
        getStoreItem_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetStoreItemRequest.fromBuffer(value),
        ($0.StoreItem value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PurchaseCartridgeRequest,
            $0.PurchaseCartridgeResponse>(
        'PurchaseCartridge',
        purchaseCartridge_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PurchaseCartridgeRequest.fromBuffer(value),
        ($0.PurchaseCartridgeResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetPurchaseHistoryRequest,
            $0.GetPurchaseHistoryResponse>(
        'GetPurchaseHistory',
        getPurchaseHistory_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetPurchaseHistoryRequest.fromBuffer(value),
        ($0.GetPurchaseHistoryResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.ListStoreItemsResponse> listStoreItems_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListStoreItemsRequest> $request) async {
    return listStoreItems($call, await $request);
  }

  $async.Future<$0.ListStoreItemsResponse> listStoreItems(
      $grpc.ServiceCall call, $0.ListStoreItemsRequest request);

  $async.Future<$0.SearchCartridgesResponse> searchCartridges_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SearchCartridgesRequest> $request) async {
    return searchCartridges($call, await $request);
  }

  $async.Future<$0.SearchCartridgesResponse> searchCartridges(
      $grpc.ServiceCall call, $0.SearchCartridgesRequest request);

  $async.Future<$0.StoreItem> getStoreItem_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetStoreItemRequest> $request) async {
    return getStoreItem($call, await $request);
  }

  $async.Future<$0.StoreItem> getStoreItem(
      $grpc.ServiceCall call, $0.GetStoreItemRequest request);

  $async.Future<$0.PurchaseCartridgeResponse> purchaseCartridge_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PurchaseCartridgeRequest> $request) async {
    return purchaseCartridge($call, await $request);
  }

  $async.Future<$0.PurchaseCartridgeResponse> purchaseCartridge(
      $grpc.ServiceCall call, $0.PurchaseCartridgeRequest request);

  $async.Future<$0.GetPurchaseHistoryResponse> getPurchaseHistory_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetPurchaseHistoryRequest> $request) async {
    return getPurchaseHistory($call, await $request);
  }

  $async.Future<$0.GetPurchaseHistoryResponse> getPurchaseHistory(
      $grpc.ServiceCall call, $0.GetPurchaseHistoryRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.CartridgeReviewService')
class CartridgeReviewServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  CartridgeReviewServiceClient(super.channel,
      {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.ReviewStatus> submitForReview(
    $0.SubmitForReviewRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$submitForReview, request, options: options);
  }

  $grpc.ResponseFuture<$0.ReviewStatus> getReviewStatus(
    $0.GetReviewStatusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getReviewStatus, request, options: options);
  }

  $grpc.ResponseFuture<$0.ReviewStatus> approveCartridge(
    $0.ApproveCartridgeRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$approveCartridge, request, options: options);
  }

  $grpc.ResponseFuture<$0.ReviewStatus> rejectCartridge(
    $0.RejectCartridgeRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$rejectCartridge, request, options: options);
  }

  // method descriptors

  static final _$submitForReview =
      $grpc.ClientMethod<$0.SubmitForReviewRequest, $0.ReviewStatus>(
          '/manpasik.v1.CartridgeReviewService/SubmitForReview',
          ($0.SubmitForReviewRequest value) => value.writeToBuffer(),
          $0.ReviewStatus.fromBuffer);
  static final _$getReviewStatus =
      $grpc.ClientMethod<$0.GetReviewStatusRequest, $0.ReviewStatus>(
          '/manpasik.v1.CartridgeReviewService/GetReviewStatus',
          ($0.GetReviewStatusRequest value) => value.writeToBuffer(),
          $0.ReviewStatus.fromBuffer);
  static final _$approveCartridge =
      $grpc.ClientMethod<$0.ApproveCartridgeRequest, $0.ReviewStatus>(
          '/manpasik.v1.CartridgeReviewService/ApproveCartridge',
          ($0.ApproveCartridgeRequest value) => value.writeToBuffer(),
          $0.ReviewStatus.fromBuffer);
  static final _$rejectCartridge =
      $grpc.ClientMethod<$0.RejectCartridgeRequest, $0.ReviewStatus>(
          '/manpasik.v1.CartridgeReviewService/RejectCartridge',
          ($0.RejectCartridgeRequest value) => value.writeToBuffer(),
          $0.ReviewStatus.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.CartridgeReviewService')
abstract class CartridgeReviewServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.CartridgeReviewService';

  CartridgeReviewServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.SubmitForReviewRequest, $0.ReviewStatus>(
        'SubmitForReview',
        submitForReview_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SubmitForReviewRequest.fromBuffer(value),
        ($0.ReviewStatus value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetReviewStatusRequest, $0.ReviewStatus>(
        'GetReviewStatus',
        getReviewStatus_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetReviewStatusRequest.fromBuffer(value),
        ($0.ReviewStatus value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ApproveCartridgeRequest, $0.ReviewStatus>(
        'ApproveCartridge',
        approveCartridge_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ApproveCartridgeRequest.fromBuffer(value),
        ($0.ReviewStatus value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RejectCartridgeRequest, $0.ReviewStatus>(
        'RejectCartridge',
        rejectCartridge_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RejectCartridgeRequest.fromBuffer(value),
        ($0.ReviewStatus value) => value.writeToBuffer()));
  }

  $async.Future<$0.ReviewStatus> submitForReview_Pre($grpc.ServiceCall $call,
      $async.Future<$0.SubmitForReviewRequest> $request) async {
    return submitForReview($call, await $request);
  }

  $async.Future<$0.ReviewStatus> submitForReview(
      $grpc.ServiceCall call, $0.SubmitForReviewRequest request);

  $async.Future<$0.ReviewStatus> getReviewStatus_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetReviewStatusRequest> $request) async {
    return getReviewStatus($call, await $request);
  }

  $async.Future<$0.ReviewStatus> getReviewStatus(
      $grpc.ServiceCall call, $0.GetReviewStatusRequest request);

  $async.Future<$0.ReviewStatus> approveCartridge_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ApproveCartridgeRequest> $request) async {
    return approveCartridge($call, await $request);
  }

  $async.Future<$0.ReviewStatus> approveCartridge(
      $grpc.ServiceCall call, $0.ApproveCartridgeRequest request);

  $async.Future<$0.ReviewStatus> rejectCartridge_Pre($grpc.ServiceCall $call,
      $async.Future<$0.RejectCartridgeRequest> $request) async {
    return rejectCartridge($call, await $request);
  }

  $async.Future<$0.ReviewStatus> rejectCartridge(
      $grpc.ServiceCall call, $0.RejectCartridgeRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.LocationStatsService')
class LocationStatsServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  LocationStatsServiceClient(super.channel,
      {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.RegionStats> getRegionStats(
    $0.GetRegionStatsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getRegionStats, request, options: options);
  }

  $grpc.ResponseFuture<$0.ListHotspotsResponse> listHotspots(
    $0.ListHotspotsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listHotspots, request, options: options);
  }

  $grpc.ResponseFuture<$0.RegionTrend> getTrendByRegion(
    $0.GetTrendByRegionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getTrendByRegion, request, options: options);
  }

  // method descriptors

  static final _$getRegionStats =
      $grpc.ClientMethod<$0.GetRegionStatsRequest, $0.RegionStats>(
          '/manpasik.v1.LocationStatsService/GetRegionStats',
          ($0.GetRegionStatsRequest value) => value.writeToBuffer(),
          $0.RegionStats.fromBuffer);
  static final _$listHotspots =
      $grpc.ClientMethod<$0.ListHotspotsRequest, $0.ListHotspotsResponse>(
          '/manpasik.v1.LocationStatsService/ListHotspots',
          ($0.ListHotspotsRequest value) => value.writeToBuffer(),
          $0.ListHotspotsResponse.fromBuffer);
  static final _$getTrendByRegion =
      $grpc.ClientMethod<$0.GetTrendByRegionRequest, $0.RegionTrend>(
          '/manpasik.v1.LocationStatsService/GetTrendByRegion',
          ($0.GetTrendByRegionRequest value) => value.writeToBuffer(),
          $0.RegionTrend.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.LocationStatsService')
abstract class LocationStatsServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.LocationStatsService';

  LocationStatsServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.GetRegionStatsRequest, $0.RegionStats>(
        'GetRegionStats',
        getRegionStats_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetRegionStatsRequest.fromBuffer(value),
        ($0.RegionStats value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListHotspotsRequest, $0.ListHotspotsResponse>(
            'ListHotspots',
            listHotspots_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListHotspotsRequest.fromBuffer(value),
            ($0.ListHotspotsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetTrendByRegionRequest, $0.RegionTrend>(
        'GetTrendByRegion',
        getTrendByRegion_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetTrendByRegionRequest.fromBuffer(value),
        ($0.RegionTrend value) => value.writeToBuffer()));
  }

  $async.Future<$0.RegionStats> getRegionStats_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetRegionStatsRequest> $request) async {
    return getRegionStats($call, await $request);
  }

  $async.Future<$0.RegionStats> getRegionStats(
      $grpc.ServiceCall call, $0.GetRegionStatsRequest request);

  $async.Future<$0.ListHotspotsResponse> listHotspots_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListHotspotsRequest> $request) async {
    return listHotspots($call, await $request);
  }

  $async.Future<$0.ListHotspotsResponse> listHotspots(
      $grpc.ServiceCall call, $0.ListHotspotsRequest request);

  $async.Future<$0.RegionTrend> getTrendByRegion_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetTrendByRegionRequest> $request) async {
    return getTrendByRegion($call, await $request);
  }

  $async.Future<$0.RegionTrend> getTrendByRegion(
      $grpc.ServiceCall call, $0.GetTrendByRegionRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.DataProvisionService')
class DataProvisionServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  DataProvisionServiceClient(super.channel,
      {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.ListDatasetsResponse> listDatasets(
    $0.ListDatasetsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDatasets, request, options: options);
  }

  $grpc.ResponseFuture<$0.AnonymizedDataset> getAnonymizedDataset(
    $0.GetAnonymizedDatasetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getAnonymizedDataset, request, options: options);
  }

  $grpc.ResponseFuture<$0.DataAccessResponse> requestDataAccess(
    $0.RequestDataAccessRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$requestDataAccess, request, options: options);
  }

  $grpc.ResponseFuture<$0.DataSharingConsentInfo> getDataSharingConsent(
    $0.GetConsentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDataSharingConsent, request, options: options);
  }

  $grpc.ResponseFuture<$0.DataSharingConsentInfo> updateDataSharingConsent(
    $0.UpdateConsentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateDataSharingConsent, request,
        options: options);
  }

  $grpc.ResponseFuture<$0.TriggerAggregationResponse>
      triggerAnonymousAggregation(
    $0.TriggerAggregationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$triggerAnonymousAggregation, request,
        options: options);
  }

  $grpc.ResponseFuture<$0.GetInsightsResponse> getAggregatedInsights(
    $0.GetInsightsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getAggregatedInsights, request, options: options);
  }

  // method descriptors

  static final _$listDatasets =
      $grpc.ClientMethod<$0.ListDatasetsRequest, $0.ListDatasetsResponse>(
          '/manpasik.v1.DataProvisionService/ListDatasets',
          ($0.ListDatasetsRequest value) => value.writeToBuffer(),
          $0.ListDatasetsResponse.fromBuffer);
  static final _$getAnonymizedDataset =
      $grpc.ClientMethod<$0.GetAnonymizedDatasetRequest, $0.AnonymizedDataset>(
          '/manpasik.v1.DataProvisionService/GetAnonymizedDataset',
          ($0.GetAnonymizedDatasetRequest value) => value.writeToBuffer(),
          $0.AnonymizedDataset.fromBuffer);
  static final _$requestDataAccess =
      $grpc.ClientMethod<$0.RequestDataAccessRequest, $0.DataAccessResponse>(
          '/manpasik.v1.DataProvisionService/RequestDataAccess',
          ($0.RequestDataAccessRequest value) => value.writeToBuffer(),
          $0.DataAccessResponse.fromBuffer);
  static final _$getDataSharingConsent =
      $grpc.ClientMethod<$0.GetConsentRequest, $0.DataSharingConsentInfo>(
          '/manpasik.v1.DataProvisionService/GetDataSharingConsent',
          ($0.GetConsentRequest value) => value.writeToBuffer(),
          $0.DataSharingConsentInfo.fromBuffer);
  static final _$updateDataSharingConsent =
      $grpc.ClientMethod<$0.UpdateConsentRequest, $0.DataSharingConsentInfo>(
          '/manpasik.v1.DataProvisionService/UpdateDataSharingConsent',
          ($0.UpdateConsentRequest value) => value.writeToBuffer(),
          $0.DataSharingConsentInfo.fromBuffer);
  static final _$triggerAnonymousAggregation = $grpc.ClientMethod<
          $0.TriggerAggregationRequest, $0.TriggerAggregationResponse>(
      '/manpasik.v1.DataProvisionService/TriggerAnonymousAggregation',
      ($0.TriggerAggregationRequest value) => value.writeToBuffer(),
      $0.TriggerAggregationResponse.fromBuffer);
  static final _$getAggregatedInsights =
      $grpc.ClientMethod<$0.GetInsightsRequest, $0.GetInsightsResponse>(
          '/manpasik.v1.DataProvisionService/GetAggregatedInsights',
          ($0.GetInsightsRequest value) => value.writeToBuffer(),
          $0.GetInsightsResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.DataProvisionService')
abstract class DataProvisionServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.DataProvisionService';

  DataProvisionServiceBase() {
    $addMethod(
        $grpc.ServiceMethod<$0.ListDatasetsRequest, $0.ListDatasetsResponse>(
            'ListDatasets',
            listDatasets_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListDatasetsRequest.fromBuffer(value),
            ($0.ListDatasetsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetAnonymizedDatasetRequest,
            $0.AnonymizedDataset>(
        'GetAnonymizedDataset',
        getAnonymizedDataset_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetAnonymizedDatasetRequest.fromBuffer(value),
        ($0.AnonymizedDataset value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.RequestDataAccessRequest, $0.DataAccessResponse>(
            'RequestDataAccess',
            requestDataAccess_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.RequestDataAccessRequest.fromBuffer(value),
            ($0.DataAccessResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetConsentRequest, $0.DataSharingConsentInfo>(
            'GetDataSharingConsent',
            getDataSharingConsent_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetConsentRequest.fromBuffer(value),
            ($0.DataSharingConsentInfo value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateConsentRequest, $0.DataSharingConsentInfo>(
            'UpdateDataSharingConsent',
            updateDataSharingConsent_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateConsentRequest.fromBuffer(value),
            ($0.DataSharingConsentInfo value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TriggerAggregationRequest,
            $0.TriggerAggregationResponse>(
        'TriggerAnonymousAggregation',
        triggerAnonymousAggregation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.TriggerAggregationRequest.fromBuffer(value),
        ($0.TriggerAggregationResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetInsightsRequest, $0.GetInsightsResponse>(
            'GetAggregatedInsights',
            getAggregatedInsights_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetInsightsRequest.fromBuffer(value),
            ($0.GetInsightsResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.ListDatasetsResponse> listDatasets_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListDatasetsRequest> $request) async {
    return listDatasets($call, await $request);
  }

  $async.Future<$0.ListDatasetsResponse> listDatasets(
      $grpc.ServiceCall call, $0.ListDatasetsRequest request);

  $async.Future<$0.AnonymizedDataset> getAnonymizedDataset_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetAnonymizedDatasetRequest> $request) async {
    return getAnonymizedDataset($call, await $request);
  }

  $async.Future<$0.AnonymizedDataset> getAnonymizedDataset(
      $grpc.ServiceCall call, $0.GetAnonymizedDatasetRequest request);

  $async.Future<$0.DataAccessResponse> requestDataAccess_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RequestDataAccessRequest> $request) async {
    return requestDataAccess($call, await $request);
  }

  $async.Future<$0.DataAccessResponse> requestDataAccess(
      $grpc.ServiceCall call, $0.RequestDataAccessRequest request);

  $async.Future<$0.DataSharingConsentInfo> getDataSharingConsent_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetConsentRequest> $request) async {
    return getDataSharingConsent($call, await $request);
  }

  $async.Future<$0.DataSharingConsentInfo> getDataSharingConsent(
      $grpc.ServiceCall call, $0.GetConsentRequest request);

  $async.Future<$0.DataSharingConsentInfo> updateDataSharingConsent_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateConsentRequest> $request) async {
    return updateDataSharingConsent($call, await $request);
  }

  $async.Future<$0.DataSharingConsentInfo> updateDataSharingConsent(
      $grpc.ServiceCall call, $0.UpdateConsentRequest request);

  $async.Future<$0.TriggerAggregationResponse> triggerAnonymousAggregation_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.TriggerAggregationRequest> $request) async {
    return triggerAnonymousAggregation($call, await $request);
  }

  $async.Future<$0.TriggerAggregationResponse> triggerAnonymousAggregation(
      $grpc.ServiceCall call, $0.TriggerAggregationRequest request);

  $async.Future<$0.GetInsightsResponse> getAggregatedInsights_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetInsightsRequest> $request) async {
    return getAggregatedInsights($call, await $request);
  }

  $async.Future<$0.GetInsightsResponse> getAggregatedInsights(
      $grpc.ServiceCall call, $0.GetInsightsRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.RevenueService')
class RevenueServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  RevenueServiceClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.SalesReport> getSalesReport(
    $0.GetSalesReportRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSalesReport, request, options: options);
  }

  $grpc.ResponseFuture<$0.GetPayoutHistoryResponse> getPayoutHistory(
    $0.GetPayoutHistoryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getPayoutHistory, request, options: options);
  }

  $grpc.ResponseFuture<$0.RevenueSplitConfig> configureRevenueSplit(
    $0.ConfigureRevenueSplitRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$configureRevenueSplit, request, options: options);
  }

  $grpc.ResponseFuture<$0.PayoutResponse> requestPayout(
    $0.RequestPayoutRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$requestPayout, request, options: options);
  }

  // method descriptors

  static final _$getSalesReport =
      $grpc.ClientMethod<$0.GetSalesReportRequest, $0.SalesReport>(
          '/manpasik.v1.RevenueService/GetSalesReport',
          ($0.GetSalesReportRequest value) => value.writeToBuffer(),
          $0.SalesReport.fromBuffer);
  static final _$getPayoutHistory = $grpc.ClientMethod<
          $0.GetPayoutHistoryRequest, $0.GetPayoutHistoryResponse>(
      '/manpasik.v1.RevenueService/GetPayoutHistory',
      ($0.GetPayoutHistoryRequest value) => value.writeToBuffer(),
      $0.GetPayoutHistoryResponse.fromBuffer);
  static final _$configureRevenueSplit = $grpc.ClientMethod<
          $0.ConfigureRevenueSplitRequest, $0.RevenueSplitConfig>(
      '/manpasik.v1.RevenueService/ConfigureRevenueSplit',
      ($0.ConfigureRevenueSplitRequest value) => value.writeToBuffer(),
      $0.RevenueSplitConfig.fromBuffer);
  static final _$requestPayout =
      $grpc.ClientMethod<$0.RequestPayoutRequest, $0.PayoutResponse>(
          '/manpasik.v1.RevenueService/RequestPayout',
          ($0.RequestPayoutRequest value) => value.writeToBuffer(),
          $0.PayoutResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.RevenueService')
abstract class RevenueServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.RevenueService';

  RevenueServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.GetSalesReportRequest, $0.SalesReport>(
        'GetSalesReport',
        getSalesReport_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetSalesReportRequest.fromBuffer(value),
        ($0.SalesReport value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetPayoutHistoryRequest,
            $0.GetPayoutHistoryResponse>(
        'GetPayoutHistory',
        getPayoutHistory_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetPayoutHistoryRequest.fromBuffer(value),
        ($0.GetPayoutHistoryResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ConfigureRevenueSplitRequest,
            $0.RevenueSplitConfig>(
        'ConfigureRevenueSplit',
        configureRevenueSplit_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ConfigureRevenueSplitRequest.fromBuffer(value),
        ($0.RevenueSplitConfig value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RequestPayoutRequest, $0.PayoutResponse>(
        'RequestPayout',
        requestPayout_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RequestPayoutRequest.fromBuffer(value),
        ($0.PayoutResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.SalesReport> getSalesReport_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetSalesReportRequest> $request) async {
    return getSalesReport($call, await $request);
  }

  $async.Future<$0.SalesReport> getSalesReport(
      $grpc.ServiceCall call, $0.GetSalesReportRequest request);

  $async.Future<$0.GetPayoutHistoryResponse> getPayoutHistory_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetPayoutHistoryRequest> $request) async {
    return getPayoutHistory($call, await $request);
  }

  $async.Future<$0.GetPayoutHistoryResponse> getPayoutHistory(
      $grpc.ServiceCall call, $0.GetPayoutHistoryRequest request);

  $async.Future<$0.RevenueSplitConfig> configureRevenueSplit_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ConfigureRevenueSplitRequest> $request) async {
    return configureRevenueSplit($call, await $request);
  }

  $async.Future<$0.RevenueSplitConfig> configureRevenueSplit(
      $grpc.ServiceCall call, $0.ConfigureRevenueSplitRequest request);

  $async.Future<$0.PayoutResponse> requestPayout_Pre($grpc.ServiceCall $call,
      $async.Future<$0.RequestPayoutRequest> $request) async {
    return requestPayout($call, await $request);
  }

  $async.Future<$0.PayoutResponse> requestPayout(
      $grpc.ServiceCall call, $0.RequestPayoutRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.CartridgeAnalyticsService')
class CartridgeAnalyticsServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  CartridgeAnalyticsServiceClient(super.channel,
      {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.CartridgeUsageStatsResponse> getUsageStats(
    $0.GetCartridgeUsageStatsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getUsageStats, request, options: options);
  }

  $grpc.ResponseFuture<$0.CartridgeRatingsResponse> getRatings(
    $0.GetCartridgeRatingsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getRatings, request, options: options);
  }

  $grpc.ResponseFuture<$0.SubmitUserReviewResponse> submitUserReview(
    $0.SubmitUserReviewRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$submitUserReview, request, options: options);
  }

  $grpc.ResponseFuture<$0.DeveloperAnalytics> getDeveloperAnalytics(
    $0.GetDeveloperAnalyticsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDeveloperAnalytics, request, options: options);
  }

  // method descriptors

  static final _$getUsageStats = $grpc.ClientMethod<
          $0.GetCartridgeUsageStatsRequest, $0.CartridgeUsageStatsResponse>(
      '/manpasik.v1.CartridgeAnalyticsService/GetUsageStats',
      ($0.GetCartridgeUsageStatsRequest value) => value.writeToBuffer(),
      $0.CartridgeUsageStatsResponse.fromBuffer);
  static final _$getRatings = $grpc.ClientMethod<$0.GetCartridgeRatingsRequest,
          $0.CartridgeRatingsResponse>(
      '/manpasik.v1.CartridgeAnalyticsService/GetRatings',
      ($0.GetCartridgeRatingsRequest value) => value.writeToBuffer(),
      $0.CartridgeRatingsResponse.fromBuffer);
  static final _$submitUserReview = $grpc.ClientMethod<
          $0.SubmitUserReviewRequest, $0.SubmitUserReviewResponse>(
      '/manpasik.v1.CartridgeAnalyticsService/SubmitUserReview',
      ($0.SubmitUserReviewRequest value) => value.writeToBuffer(),
      $0.SubmitUserReviewResponse.fromBuffer);
  static final _$getDeveloperAnalytics = $grpc.ClientMethod<
          $0.GetDeveloperAnalyticsRequest, $0.DeveloperAnalytics>(
      '/manpasik.v1.CartridgeAnalyticsService/GetDeveloperAnalytics',
      ($0.GetDeveloperAnalyticsRequest value) => value.writeToBuffer(),
      $0.DeveloperAnalytics.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.CartridgeAnalyticsService')
abstract class CartridgeAnalyticsServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.CartridgeAnalyticsService';

  CartridgeAnalyticsServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.GetCartridgeUsageStatsRequest,
            $0.CartridgeUsageStatsResponse>(
        'GetUsageStats',
        getUsageStats_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetCartridgeUsageStatsRequest.fromBuffer(value),
        ($0.CartridgeUsageStatsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetCartridgeRatingsRequest,
            $0.CartridgeRatingsResponse>(
        'GetRatings',
        getRatings_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetCartridgeRatingsRequest.fromBuffer(value),
        ($0.CartridgeRatingsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SubmitUserReviewRequest,
            $0.SubmitUserReviewResponse>(
        'SubmitUserReview',
        submitUserReview_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SubmitUserReviewRequest.fromBuffer(value),
        ($0.SubmitUserReviewResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDeveloperAnalyticsRequest,
            $0.DeveloperAnalytics>(
        'GetDeveloperAnalytics',
        getDeveloperAnalytics_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDeveloperAnalyticsRequest.fromBuffer(value),
        ($0.DeveloperAnalytics value) => value.writeToBuffer()));
  }

  $async.Future<$0.CartridgeUsageStatsResponse> getUsageStats_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetCartridgeUsageStatsRequest> $request) async {
    return getUsageStats($call, await $request);
  }

  $async.Future<$0.CartridgeUsageStatsResponse> getUsageStats(
      $grpc.ServiceCall call, $0.GetCartridgeUsageStatsRequest request);

  $async.Future<$0.CartridgeRatingsResponse> getRatings_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetCartridgeRatingsRequest> $request) async {
    return getRatings($call, await $request);
  }

  $async.Future<$0.CartridgeRatingsResponse> getRatings(
      $grpc.ServiceCall call, $0.GetCartridgeRatingsRequest request);

  $async.Future<$0.SubmitUserReviewResponse> submitUserReview_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SubmitUserReviewRequest> $request) async {
    return submitUserReview($call, await $request);
  }

  $async.Future<$0.SubmitUserReviewResponse> submitUserReview(
      $grpc.ServiceCall call, $0.SubmitUserReviewRequest request);

  $async.Future<$0.DeveloperAnalytics> getDeveloperAnalytics_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDeveloperAnalyticsRequest> $request) async {
    return getDeveloperAnalytics($call, await $request);
  }

  $async.Future<$0.DeveloperAnalytics> getDeveloperAnalytics(
      $grpc.ServiceCall call, $0.GetDeveloperAnalyticsRequest request);
}

@$pb.GrpcServiceName('manpasik.v1.VoiceProfileService')
class VoiceProfileServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  VoiceProfileServiceClient(super.channel, {super.options, super.interceptors});

  $grpc.ResponseFuture<$0.VoiceProfile> createVoiceProfile(
    $0.CreateVoiceProfileRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createVoiceProfile, request, options: options);
  }

  $grpc.ResponseFuture<$0.VoiceProfile> getVoiceProfile(
    $0.GetVoiceProfileRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getVoiceProfile, request, options: options);
  }

  $grpc.ResponseFuture<$0.ListVoiceProfilesResponse> listVoiceProfiles(
    $0.ListVoiceProfilesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listVoiceProfiles, request, options: options);
  }

  $grpc.ResponseFuture<$0.SynthesizeTranslationResponse> synthesizeTranslation(
    $0.SynthesizeTranslationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$synthesizeTranslation, request, options: options);
  }

  $grpc.ResponseFuture<$0.DeleteVoiceProfileResponse> deleteVoiceProfile(
    $0.DeleteVoiceProfileRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteVoiceProfile, request, options: options);
  }

  // method descriptors

  static final _$createVoiceProfile =
      $grpc.ClientMethod<$0.CreateVoiceProfileRequest, $0.VoiceProfile>(
          '/manpasik.v1.VoiceProfileService/CreateVoiceProfile',
          ($0.CreateVoiceProfileRequest value) => value.writeToBuffer(),
          $0.VoiceProfile.fromBuffer);
  static final _$getVoiceProfile =
      $grpc.ClientMethod<$0.GetVoiceProfileRequest, $0.VoiceProfile>(
          '/manpasik.v1.VoiceProfileService/GetVoiceProfile',
          ($0.GetVoiceProfileRequest value) => value.writeToBuffer(),
          $0.VoiceProfile.fromBuffer);
  static final _$listVoiceProfiles = $grpc.ClientMethod<
          $0.ListVoiceProfilesRequest, $0.ListVoiceProfilesResponse>(
      '/manpasik.v1.VoiceProfileService/ListVoiceProfiles',
      ($0.ListVoiceProfilesRequest value) => value.writeToBuffer(),
      $0.ListVoiceProfilesResponse.fromBuffer);
  static final _$synthesizeTranslation = $grpc.ClientMethod<
          $0.SynthesizeTranslationRequest, $0.SynthesizeTranslationResponse>(
      '/manpasik.v1.VoiceProfileService/SynthesizeTranslation',
      ($0.SynthesizeTranslationRequest value) => value.writeToBuffer(),
      $0.SynthesizeTranslationResponse.fromBuffer);
  static final _$deleteVoiceProfile = $grpc.ClientMethod<
          $0.DeleteVoiceProfileRequest, $0.DeleteVoiceProfileResponse>(
      '/manpasik.v1.VoiceProfileService/DeleteVoiceProfile',
      ($0.DeleteVoiceProfileRequest value) => value.writeToBuffer(),
      $0.DeleteVoiceProfileResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.v1.VoiceProfileService')
abstract class VoiceProfileServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.v1.VoiceProfileService';

  VoiceProfileServiceBase() {
    $addMethod(
        $grpc.ServiceMethod<$0.CreateVoiceProfileRequest, $0.VoiceProfile>(
            'CreateVoiceProfile',
            createVoiceProfile_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateVoiceProfileRequest.fromBuffer(value),
            ($0.VoiceProfile value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetVoiceProfileRequest, $0.VoiceProfile>(
        'GetVoiceProfile',
        getVoiceProfile_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetVoiceProfileRequest.fromBuffer(value),
        ($0.VoiceProfile value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListVoiceProfilesRequest,
            $0.ListVoiceProfilesResponse>(
        'ListVoiceProfiles',
        listVoiceProfiles_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListVoiceProfilesRequest.fromBuffer(value),
        ($0.ListVoiceProfilesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SynthesizeTranslationRequest,
            $0.SynthesizeTranslationResponse>(
        'SynthesizeTranslation',
        synthesizeTranslation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SynthesizeTranslationRequest.fromBuffer(value),
        ($0.SynthesizeTranslationResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteVoiceProfileRequest,
            $0.DeleteVoiceProfileResponse>(
        'DeleteVoiceProfile',
        deleteVoiceProfile_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteVoiceProfileRequest.fromBuffer(value),
        ($0.DeleteVoiceProfileResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.VoiceProfile> createVoiceProfile_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateVoiceProfileRequest> $request) async {
    return createVoiceProfile($call, await $request);
  }

  $async.Future<$0.VoiceProfile> createVoiceProfile(
      $grpc.ServiceCall call, $0.CreateVoiceProfileRequest request);

  $async.Future<$0.VoiceProfile> getVoiceProfile_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetVoiceProfileRequest> $request) async {
    return getVoiceProfile($call, await $request);
  }

  $async.Future<$0.VoiceProfile> getVoiceProfile(
      $grpc.ServiceCall call, $0.GetVoiceProfileRequest request);

  $async.Future<$0.ListVoiceProfilesResponse> listVoiceProfiles_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListVoiceProfilesRequest> $request) async {
    return listVoiceProfiles($call, await $request);
  }

  $async.Future<$0.ListVoiceProfilesResponse> listVoiceProfiles(
      $grpc.ServiceCall call, $0.ListVoiceProfilesRequest request);

  $async.Future<$0.SynthesizeTranslationResponse> synthesizeTranslation_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SynthesizeTranslationRequest> $request) async {
    return synthesizeTranslation($call, await $request);
  }

  $async.Future<$0.SynthesizeTranslationResponse> synthesizeTranslation(
      $grpc.ServiceCall call, $0.SynthesizeTranslationRequest request);

  $async.Future<$0.DeleteVoiceProfileResponse> deleteVoiceProfile_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteVoiceProfileRequest> $request) async {
    return deleteVoiceProfile($call, await $request);
  }

  $async.Future<$0.DeleteVoiceProfileResponse> deleteVoiceProfile(
      $grpc.ServiceCall call, $0.DeleteVoiceProfileRequest request);
}
