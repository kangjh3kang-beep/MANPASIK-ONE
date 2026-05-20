// This is a generated file - do not edit.
//
// Generated from measurement.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

class SyncMeasurementsRequest extends $pb.GeneratedMessage {
  factory SyncMeasurementsRequest({
    $core.Iterable<MeasurementRecord>? records,
  }) {
    final result = create();
    if (records != null) result.records.addAll(records);
    return result;
  }

  SyncMeasurementsRequest._();

  factory SyncMeasurementsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SyncMeasurementsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SyncMeasurementsRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'manpasik.api.v1'),
      createEmptyInstance: create)
    ..pPM<MeasurementRecord>(1, _omitFieldNames ? '' : 'records',
        subBuilder: MeasurementRecord.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SyncMeasurementsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SyncMeasurementsRequest copyWith(
          void Function(SyncMeasurementsRequest) updates) =>
      super.copyWith((message) => updates(message as SyncMeasurementsRequest))
          as SyncMeasurementsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SyncMeasurementsRequest create() => SyncMeasurementsRequest._();
  @$core.override
  SyncMeasurementsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SyncMeasurementsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SyncMeasurementsRequest>(create);
  static SyncMeasurementsRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $pb.PbList<MeasurementRecord> get records => $_getList(0);
}

class SyncMeasurementsResponse extends $pb.GeneratedMessage {
  factory SyncMeasurementsResponse({
    $core.int? syncedCount,
    $core.Iterable<$core.String>? failedIds,
  }) {
    final result = create();
    if (syncedCount != null) result.syncedCount = syncedCount;
    if (failedIds != null) result.failedIds.addAll(failedIds);
    return result;
  }

  SyncMeasurementsResponse._();

  factory SyncMeasurementsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SyncMeasurementsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SyncMeasurementsResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'manpasik.api.v1'),
      createEmptyInstance: create)
    ..aI(1, _omitFieldNames ? '' : 'syncedCount')
    ..pPS(2, _omitFieldNames ? '' : 'failedIds')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SyncMeasurementsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SyncMeasurementsResponse copyWith(
          void Function(SyncMeasurementsResponse) updates) =>
      super.copyWith((message) => updates(message as SyncMeasurementsResponse))
          as SyncMeasurementsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SyncMeasurementsResponse create() => SyncMeasurementsResponse._();
  @$core.override
  SyncMeasurementsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SyncMeasurementsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SyncMeasurementsResponse>(create);
  static SyncMeasurementsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get syncedCount => $_getIZ(0);
  @$pb.TagNumber(1)
  set syncedCount($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSyncedCount() => $_has(0);
  @$pb.TagNumber(1)
  void clearSyncedCount() => $_clearField(1);

  @$pb.TagNumber(2)
  $pb.PbList<$core.String> get failedIds => $_getList(1);
}

class MeasurementRecord extends $pb.GeneratedMessage {
  factory MeasurementRecord({
    $core.String? id,
    $core.String? deviceMac,
    $fixnum.Int64? timestamp,
    $core.Iterable<$core.double>? diffSignal,
    $core.Iterable<$core.double>? fingerprint,
    $core.int? healthScore,
    $core.String? riskLabel,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (deviceMac != null) result.deviceMac = deviceMac;
    if (timestamp != null) result.timestamp = timestamp;
    if (diffSignal != null) result.diffSignal.addAll(diffSignal);
    if (fingerprint != null) result.fingerprint.addAll(fingerprint);
    if (healthScore != null) result.healthScore = healthScore;
    if (riskLabel != null) result.riskLabel = riskLabel;
    return result;
  }

  MeasurementRecord._();

  factory MeasurementRecord.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MeasurementRecord.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MeasurementRecord',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'manpasik.api.v1'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'deviceMac')
    ..aInt64(3, _omitFieldNames ? '' : 'timestamp')
    ..p<$core.double>(
        4, _omitFieldNames ? '' : 'diffSignal', $pb.PbFieldType.KD)
    ..p<$core.double>(
        5, _omitFieldNames ? '' : 'fingerprint', $pb.PbFieldType.KD)
    ..aI(6, _omitFieldNames ? '' : 'healthScore')
    ..aOS(7, _omitFieldNames ? '' : 'riskLabel')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MeasurementRecord clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MeasurementRecord copyWith(void Function(MeasurementRecord) updates) =>
      super.copyWith((message) => updates(message as MeasurementRecord))
          as MeasurementRecord;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MeasurementRecord create() => MeasurementRecord._();
  @$core.override
  MeasurementRecord createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MeasurementRecord getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MeasurementRecord>(create);
  static MeasurementRecord? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get deviceMac => $_getSZ(1);
  @$pb.TagNumber(2)
  set deviceMac($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasDeviceMac() => $_has(1);
  @$pb.TagNumber(2)
  void clearDeviceMac() => $_clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get timestamp => $_getI64(2);
  @$pb.TagNumber(3)
  set timestamp($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasTimestamp() => $_has(2);
  @$pb.TagNumber(3)
  void clearTimestamp() => $_clearField(3);

  @$pb.TagNumber(4)
  $pb.PbList<$core.double> get diffSignal => $_getList(3);

  @$pb.TagNumber(5)
  $pb.PbList<$core.double> get fingerprint => $_getList(4);

  @$pb.TagNumber(6)
  $core.int get healthScore => $_getIZ(5);
  @$pb.TagNumber(6)
  set healthScore($core.int value) => $_setSignedInt32(5, value);
  @$pb.TagNumber(6)
  $core.bool hasHealthScore() => $_has(5);
  @$pb.TagNumber(6)
  void clearHealthScore() => $_clearField(6);

  @$pb.TagNumber(7)
  $core.String get riskLabel => $_getSZ(6);
  @$pb.TagNumber(7)
  set riskLabel($core.String value) => $_setString(6, value);
  @$pb.TagNumber(7)
  $core.bool hasRiskLabel() => $_has(6);
  @$pb.TagNumber(7)
  void clearRiskLabel() => $_clearField(7);
}

class SubmitCheckupRequest extends $pb.GeneratedMessage {
  factory SubmitCheckupRequest({
    $core.String? sessionId,
    $core.String? userId,
    $core.String? packageType,
    $core.Iterable<MeasurementRecord>? results,
  }) {
    final result = create();
    if (sessionId != null) result.sessionId = sessionId;
    if (userId != null) result.userId = userId;
    if (packageType != null) result.packageType = packageType;
    if (results != null) result.results.addAll(results);
    return result;
  }

  SubmitCheckupRequest._();

  factory SubmitCheckupRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SubmitCheckupRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SubmitCheckupRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'manpasik.api.v1'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'sessionId')
    ..aOS(2, _omitFieldNames ? '' : 'userId')
    ..aOS(3, _omitFieldNames ? '' : 'packageType')
    ..pPM<MeasurementRecord>(4, _omitFieldNames ? '' : 'results',
        subBuilder: MeasurementRecord.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SubmitCheckupRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SubmitCheckupRequest copyWith(void Function(SubmitCheckupRequest) updates) =>
      super.copyWith((message) => updates(message as SubmitCheckupRequest))
          as SubmitCheckupRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SubmitCheckupRequest create() => SubmitCheckupRequest._();
  @$core.override
  SubmitCheckupRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SubmitCheckupRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SubmitCheckupRequest>(create);
  static SubmitCheckupRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get sessionId => $_getSZ(0);
  @$pb.TagNumber(1)
  set sessionId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSessionId() => $_has(0);
  @$pb.TagNumber(1)
  void clearSessionId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get userId => $_getSZ(1);
  @$pb.TagNumber(2)
  set userId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasUserId() => $_has(1);
  @$pb.TagNumber(2)
  void clearUserId() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get packageType => $_getSZ(2);
  @$pb.TagNumber(3)
  set packageType($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasPackageType() => $_has(2);
  @$pb.TagNumber(3)
  void clearPackageType() => $_clearField(3);

  @$pb.TagNumber(4)
  $pb.PbList<MeasurementRecord> get results => $_getList(3);
}

class SubmitCheckupResponse extends $pb.GeneratedMessage {
  factory SubmitCheckupResponse({
    $core.bool? success,
    $core.double? compositeScore,
  }) {
    final result = create();
    if (success != null) result.success = success;
    if (compositeScore != null) result.compositeScore = compositeScore;
    return result;
  }

  SubmitCheckupResponse._();

  factory SubmitCheckupResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SubmitCheckupResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SubmitCheckupResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'manpasik.api.v1'),
      createEmptyInstance: create)
    ..aOB(1, _omitFieldNames ? '' : 'success')
    ..aD(2, _omitFieldNames ? '' : 'compositeScore')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SubmitCheckupResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SubmitCheckupResponse copyWith(
          void Function(SubmitCheckupResponse) updates) =>
      super.copyWith((message) => updates(message as SubmitCheckupResponse))
          as SubmitCheckupResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SubmitCheckupResponse create() => SubmitCheckupResponse._();
  @$core.override
  SubmitCheckupResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SubmitCheckupResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SubmitCheckupResponse>(create);
  static SubmitCheckupResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.bool get success => $_getBF(0);
  @$pb.TagNumber(1)
  set success($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSuccess() => $_has(0);
  @$pb.TagNumber(1)
  void clearSuccess() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.double get compositeScore => $_getN(1);
  @$pb.TagNumber(2)
  set compositeScore($core.double value) => $_setDouble(1, value);
  @$pb.TagNumber(2)
  $core.bool hasCompositeScore() => $_has(1);
  @$pb.TagNumber(2)
  void clearCompositeScore() => $_clearField(2);
}

class ContextRequest extends $pb.GeneratedMessage {
  factory ContextRequest({
    $core.String? userId,
    $core.int? requestedCards,
  }) {
    final result = create();
    if (userId != null) result.userId = userId;
    if (requestedCards != null) result.requestedCards = requestedCards;
    return result;
  }

  ContextRequest._();

  factory ContextRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ContextRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ContextRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'manpasik.api.v1'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'userId')
    ..aI(2, _omitFieldNames ? '' : 'requestedCards')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ContextRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ContextRequest copyWith(void Function(ContextRequest) updates) =>
      super.copyWith((message) => updates(message as ContextRequest))
          as ContextRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ContextRequest create() => ContextRequest._();
  @$core.override
  ContextRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ContextRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ContextRequest>(create);
  static ContextRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get userId => $_getSZ(0);
  @$pb.TagNumber(1)
  set userId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasUserId() => $_has(0);
  @$pb.TagNumber(1)
  void clearUserId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.int get requestedCards => $_getIZ(1);
  @$pb.TagNumber(2)
  set requestedCards($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasRequestedCards() => $_has(1);
  @$pb.TagNumber(2)
  void clearRequestedCards() => $_clearField(2);
}

class ContextCardResponse extends $pb.GeneratedMessage {
  factory ContextCardResponse({
    $core.String? cardId,
    $core.String? type,
    $core.String? title,
    $core.String? body,
    $core.int? priority,
  }) {
    final result = create();
    if (cardId != null) result.cardId = cardId;
    if (type != null) result.type = type;
    if (title != null) result.title = title;
    if (body != null) result.body = body;
    if (priority != null) result.priority = priority;
    return result;
  }

  ContextCardResponse._();

  factory ContextCardResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ContextCardResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ContextCardResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'manpasik.api.v1'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'cardId')
    ..aOS(2, _omitFieldNames ? '' : 'type')
    ..aOS(3, _omitFieldNames ? '' : 'title')
    ..aOS(4, _omitFieldNames ? '' : 'body')
    ..aI(5, _omitFieldNames ? '' : 'priority')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ContextCardResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ContextCardResponse copyWith(void Function(ContextCardResponse) updates) =>
      super.copyWith((message) => updates(message as ContextCardResponse))
          as ContextCardResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ContextCardResponse create() => ContextCardResponse._();
  @$core.override
  ContextCardResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ContextCardResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ContextCardResponse>(create);
  static ContextCardResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get cardId => $_getSZ(0);
  @$pb.TagNumber(1)
  set cardId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasCardId() => $_has(0);
  @$pb.TagNumber(1)
  void clearCardId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get type => $_getSZ(1);
  @$pb.TagNumber(2)
  set type($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasType() => $_has(1);
  @$pb.TagNumber(2)
  void clearType() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get title => $_getSZ(2);
  @$pb.TagNumber(3)
  set title($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasTitle() => $_has(2);
  @$pb.TagNumber(3)
  void clearTitle() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.String get body => $_getSZ(3);
  @$pb.TagNumber(4)
  set body($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasBody() => $_has(3);
  @$pb.TagNumber(4)
  void clearBody() => $_clearField(4);

  @$pb.TagNumber(5)
  $core.int get priority => $_getIZ(4);
  @$pb.TagNumber(5)
  set priority($core.int value) => $_setSignedInt32(4, value);
  @$pb.TagNumber(5)
  $core.bool hasPriority() => $_has(4);
  @$pb.TagNumber(5)
  void clearPriority() => $_clearField(5);
}

class HealthReportRequest extends $pb.GeneratedMessage {
  factory HealthReportRequest({
    $core.String? deviceId,
    $core.double? overallScore,
    $core.double? hwScore,
    $core.double? rustScore,
    $fixnum.Int64? timestamp,
  }) {
    final result = create();
    if (deviceId != null) result.deviceId = deviceId;
    if (overallScore != null) result.overallScore = overallScore;
    if (hwScore != null) result.hwScore = hwScore;
    if (rustScore != null) result.rustScore = rustScore;
    if (timestamp != null) result.timestamp = timestamp;
    return result;
  }

  HealthReportRequest._();

  factory HealthReportRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HealthReportRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HealthReportRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'manpasik.api.v1'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'deviceId')
    ..aD(2, _omitFieldNames ? '' : 'overallScore')
    ..aD(3, _omitFieldNames ? '' : 'hwScore')
    ..aD(4, _omitFieldNames ? '' : 'rustScore')
    ..aInt64(5, _omitFieldNames ? '' : 'timestamp')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HealthReportRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HealthReportRequest copyWith(void Function(HealthReportRequest) updates) =>
      super.copyWith((message) => updates(message as HealthReportRequest))
          as HealthReportRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HealthReportRequest create() => HealthReportRequest._();
  @$core.override
  HealthReportRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HealthReportRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HealthReportRequest>(create);
  static HealthReportRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get deviceId => $_getSZ(0);
  @$pb.TagNumber(1)
  set deviceId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasDeviceId() => $_has(0);
  @$pb.TagNumber(1)
  void clearDeviceId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.double get overallScore => $_getN(1);
  @$pb.TagNumber(2)
  set overallScore($core.double value) => $_setDouble(1, value);
  @$pb.TagNumber(2)
  $core.bool hasOverallScore() => $_has(1);
  @$pb.TagNumber(2)
  void clearOverallScore() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.double get hwScore => $_getN(2);
  @$pb.TagNumber(3)
  set hwScore($core.double value) => $_setDouble(2, value);
  @$pb.TagNumber(3)
  $core.bool hasHwScore() => $_has(2);
  @$pb.TagNumber(3)
  void clearHwScore() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.double get rustScore => $_getN(3);
  @$pb.TagNumber(4)
  set rustScore($core.double value) => $_setDouble(3, value);
  @$pb.TagNumber(4)
  $core.bool hasRustScore() => $_has(3);
  @$pb.TagNumber(4)
  void clearRustScore() => $_clearField(4);

  @$pb.TagNumber(5)
  $fixnum.Int64 get timestamp => $_getI64(4);
  @$pb.TagNumber(5)
  set timestamp($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(5)
  $core.bool hasTimestamp() => $_has(4);
  @$pb.TagNumber(5)
  void clearTimestamp() => $_clearField(5);
}

class HealthReportResponse extends $pb.GeneratedMessage {
  factory HealthReportResponse({
    $core.bool? requiresServiceTicket,
    $core.String? ticketId,
  }) {
    final result = create();
    if (requiresServiceTicket != null)
      result.requiresServiceTicket = requiresServiceTicket;
    if (ticketId != null) result.ticketId = ticketId;
    return result;
  }

  HealthReportResponse._();

  factory HealthReportResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HealthReportResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HealthReportResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'manpasik.api.v1'),
      createEmptyInstance: create)
    ..aOB(1, _omitFieldNames ? '' : 'requiresServiceTicket')
    ..aOS(2, _omitFieldNames ? '' : 'ticketId')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HealthReportResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HealthReportResponse copyWith(void Function(HealthReportResponse) updates) =>
      super.copyWith((message) => updates(message as HealthReportResponse))
          as HealthReportResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HealthReportResponse create() => HealthReportResponse._();
  @$core.override
  HealthReportResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HealthReportResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HealthReportResponse>(create);
  static HealthReportResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.bool get requiresServiceTicket => $_getBF(0);
  @$pb.TagNumber(1)
  set requiresServiceTicket($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(1)
  $core.bool hasRequiresServiceTicket() => $_has(0);
  @$pb.TagNumber(1)
  void clearRequiresServiceTicket() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get ticketId => $_getSZ(1);
  @$pb.TagNumber(2)
  set ticketId($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasTicketId() => $_has(1);
  @$pb.TagNumber(2)
  void clearTicketId() => $_clearField(2);
}

class HealingEventRequest extends $pb.GeneratedMessage {
  factory HealingEventRequest({
    $core.String? deviceId,
    $core.String? layer,
    $core.String? strategy,
    $core.bool? success,
    $fixnum.Int64? timestamp,
  }) {
    final result = create();
    if (deviceId != null) result.deviceId = deviceId;
    if (layer != null) result.layer = layer;
    if (strategy != null) result.strategy = strategy;
    if (success != null) result.success = success;
    if (timestamp != null) result.timestamp = timestamp;
    return result;
  }

  HealingEventRequest._();

  factory HealingEventRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HealingEventRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HealingEventRequest',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'manpasik.api.v1'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'deviceId')
    ..aOS(2, _omitFieldNames ? '' : 'layer')
    ..aOS(3, _omitFieldNames ? '' : 'strategy')
    ..aOB(4, _omitFieldNames ? '' : 'success')
    ..aInt64(5, _omitFieldNames ? '' : 'timestamp')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HealingEventRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HealingEventRequest copyWith(void Function(HealingEventRequest) updates) =>
      super.copyWith((message) => updates(message as HealingEventRequest))
          as HealingEventRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HealingEventRequest create() => HealingEventRequest._();
  @$core.override
  HealingEventRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HealingEventRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HealingEventRequest>(create);
  static HealingEventRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get deviceId => $_getSZ(0);
  @$pb.TagNumber(1)
  set deviceId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasDeviceId() => $_has(0);
  @$pb.TagNumber(1)
  void clearDeviceId() => $_clearField(1);

  @$pb.TagNumber(2)
  $core.String get layer => $_getSZ(1);
  @$pb.TagNumber(2)
  set layer($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasLayer() => $_has(1);
  @$pb.TagNumber(2)
  void clearLayer() => $_clearField(2);

  @$pb.TagNumber(3)
  $core.String get strategy => $_getSZ(2);
  @$pb.TagNumber(3)
  set strategy($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasStrategy() => $_has(2);
  @$pb.TagNumber(3)
  void clearStrategy() => $_clearField(3);

  @$pb.TagNumber(4)
  $core.bool get success => $_getBF(3);
  @$pb.TagNumber(4)
  set success($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(4)
  $core.bool hasSuccess() => $_has(3);
  @$pb.TagNumber(4)
  void clearSuccess() => $_clearField(4);

  @$pb.TagNumber(5)
  $fixnum.Int64 get timestamp => $_getI64(4);
  @$pb.TagNumber(5)
  set timestamp($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(5)
  $core.bool hasTimestamp() => $_has(4);
  @$pb.TagNumber(5)
  void clearTimestamp() => $_clearField(5);
}

class HealingEventResponse extends $pb.GeneratedMessage {
  factory HealingEventResponse({
    $core.bool? ack,
  }) {
    final result = create();
    if (ack != null) result.ack = ack;
    return result;
  }

  HealingEventResponse._();

  factory HealingEventResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HealingEventResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HealingEventResponse',
      package:
          const $pb.PackageName(_omitMessageNames ? '' : 'manpasik.api.v1'),
      createEmptyInstance: create)
    ..aOB(1, _omitFieldNames ? '' : 'ack')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HealingEventResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HealingEventResponse copyWith(void Function(HealingEventResponse) updates) =>
      super.copyWith((message) => updates(message as HealingEventResponse))
          as HealingEventResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HealingEventResponse create() => HealingEventResponse._();
  @$core.override
  HealingEventResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HealingEventResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HealingEventResponse>(create);
  static HealingEventResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.bool get ack => $_getBF(0);
  @$pb.TagNumber(1)
  set ack($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(1)
  $core.bool hasAck() => $_has(0);
  @$pb.TagNumber(1)
  void clearAck() => $_clearField(1);
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
