// This is a generated file - do not edit.
//
// Generated from measurement.proto.

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

import 'measurement.pb.dart' as $0;

export 'measurement.pb.dart';

/// 1. 측정 데이터 수집
@$pb.GrpcServiceName('manpasik.api.v1.MeasurementService')
class MeasurementServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  MeasurementServiceClient(super.channel, {super.options, super.interceptors});

  /// 측정 완료된 데이터를 중앙(클라우드)으로 오프라인 벌크 동기화 (CRDT)
  $grpc.ResponseFuture<$0.SyncMeasurementsResponse> syncMeasurements(
    $0.SyncMeasurementsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$syncMeasurements, request, options: options);
  }

  /// v2.1: 패키지 종합 검진 결과 수집
  $grpc.ResponseFuture<$0.SubmitCheckupResponse> submitCheckupSession(
    $0.SubmitCheckupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$submitCheckupSession, request, options: options);
  }

  // method descriptors

  static final _$syncMeasurements = $grpc.ClientMethod<
          $0.SyncMeasurementsRequest, $0.SyncMeasurementsResponse>(
      '/manpasik.api.v1.MeasurementService/SyncMeasurements',
      ($0.SyncMeasurementsRequest value) => value.writeToBuffer(),
      $0.SyncMeasurementsResponse.fromBuffer);
  static final _$submitCheckupSession =
      $grpc.ClientMethod<$0.SubmitCheckupRequest, $0.SubmitCheckupResponse>(
          '/manpasik.api.v1.MeasurementService/SubmitCheckupSession',
          ($0.SubmitCheckupRequest value) => value.writeToBuffer(),
          $0.SubmitCheckupResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.api.v1.MeasurementService')
abstract class MeasurementServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.api.v1.MeasurementService';

  MeasurementServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.SyncMeasurementsRequest,
            $0.SyncMeasurementsResponse>(
        'SyncMeasurements',
        syncMeasurements_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SyncMeasurementsRequest.fromBuffer(value),
        ($0.SyncMeasurementsResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.SubmitCheckupRequest, $0.SubmitCheckupResponse>(
            'SubmitCheckupSession',
            submitCheckupSession_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.SubmitCheckupRequest.fromBuffer(value),
            ($0.SubmitCheckupResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.SyncMeasurementsResponse> syncMeasurements_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SyncMeasurementsRequest> $request) async {
    return syncMeasurements($call, await $request);
  }

  $async.Future<$0.SyncMeasurementsResponse> syncMeasurements(
      $grpc.ServiceCall call, $0.SyncMeasurementsRequest request);

  $async.Future<$0.SubmitCheckupResponse> submitCheckupSession_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SubmitCheckupRequest> $request) async {
    return submitCheckupSession($call, await $request);
  }

  $async.Future<$0.SubmitCheckupResponse> submitCheckupSession(
      $grpc.ServiceCall call, $0.SubmitCheckupRequest request);
}

/// 2. 유기적 연동 + 모바일 상태 보강 (Context Engine)
@$pb.GrpcServiceName('manpasik.api.v1.ContextService')
class ContextServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  ContextServiceClient(super.channel, {super.options, super.interceptors});

  /// v2.2: 사용자 상태를 분석하여 컨텍스트 카드 스트림 내려줌
  $grpc.ResponseStream<$0.ContextCardResponse> getContextCards(
    $0.ContextRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createStreamingCall(
        _$getContextCards, $async.Stream.fromIterable([request]),
        options: options);
  }

  // method descriptors

  static final _$getContextCards =
      $grpc.ClientMethod<$0.ContextRequest, $0.ContextCardResponse>(
          '/manpasik.api.v1.ContextService/GetContextCards',
          ($0.ContextRequest value) => value.writeToBuffer(),
          $0.ContextCardResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.api.v1.ContextService')
abstract class ContextServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.api.v1.ContextService';

  ContextServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.ContextRequest, $0.ContextCardResponse>(
        'GetContextCards',
        getContextCards_Pre,
        false,
        true,
        ($core.List<$core.int> value) => $0.ContextRequest.fromBuffer(value),
        ($0.ContextCardResponse value) => value.writeToBuffer()));
  }

  $async.Stream<$0.ContextCardResponse> getContextCards_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ContextRequest> $request) async* {
    yield* getContextCards($call, await $request);
  }

  $async.Stream<$0.ContextCardResponse> getContextCards(
      $grpc.ServiceCall call, $0.ContextRequest request);
}

/// 3. 자가검증 파이프라인 수집기
@$pb.GrpcServiceName('manpasik.api.v1.DiagnosticService')
class DiagnosticServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  DiagnosticServiceClient(super.channel, {super.options, super.interceptors});

  /// v2.3: 시스템 통합 건강 점수(L1~L6) 및 HwAlert 레포팅
  $grpc.ResponseFuture<$0.HealthReportResponse> reportSystemHealth(
    $0.HealthReportRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$reportSystemHealth, request, options: options);
  }

  /// v2.4: 자가치유 결과 기록 (원격 에스컬레이션 티켓 연계용)
  $grpc.ResponseFuture<$0.HealingEventResponse> reportHealingEvent(
    $0.HealingEventRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$reportHealingEvent, request, options: options);
  }

  // method descriptors

  static final _$reportSystemHealth =
      $grpc.ClientMethod<$0.HealthReportRequest, $0.HealthReportResponse>(
          '/manpasik.api.v1.DiagnosticService/ReportSystemHealth',
          ($0.HealthReportRequest value) => value.writeToBuffer(),
          $0.HealthReportResponse.fromBuffer);
  static final _$reportHealingEvent =
      $grpc.ClientMethod<$0.HealingEventRequest, $0.HealingEventResponse>(
          '/manpasik.api.v1.DiagnosticService/ReportHealingEvent',
          ($0.HealingEventRequest value) => value.writeToBuffer(),
          $0.HealingEventResponse.fromBuffer);
}

@$pb.GrpcServiceName('manpasik.api.v1.DiagnosticService')
abstract class DiagnosticServiceBase extends $grpc.Service {
  $core.String get $name => 'manpasik.api.v1.DiagnosticService';

  DiagnosticServiceBase() {
    $addMethod(
        $grpc.ServiceMethod<$0.HealthReportRequest, $0.HealthReportResponse>(
            'ReportSystemHealth',
            reportSystemHealth_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.HealthReportRequest.fromBuffer(value),
            ($0.HealthReportResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.HealingEventRequest, $0.HealingEventResponse>(
            'ReportHealingEvent',
            reportHealingEvent_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.HealingEventRequest.fromBuffer(value),
            ($0.HealingEventResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.HealthReportResponse> reportSystemHealth_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.HealthReportRequest> $request) async {
    return reportSystemHealth($call, await $request);
  }

  $async.Future<$0.HealthReportResponse> reportSystemHealth(
      $grpc.ServiceCall call, $0.HealthReportRequest request);

  $async.Future<$0.HealingEventResponse> reportHealingEvent_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.HealingEventRequest> $request) async {
    return reportHealingEvent($call, await $request);
  }

  $async.Future<$0.HealingEventResponse> reportHealingEvent(
      $grpc.ServiceCall call, $0.HealingEventRequest request);
}
