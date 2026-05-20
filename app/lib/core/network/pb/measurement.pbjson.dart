// This is a generated file - do not edit.
//
// Generated from measurement.proto.

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

@$core.Deprecated('Use syncMeasurementsRequestDescriptor instead')
const SyncMeasurementsRequest$json = {
  '1': 'SyncMeasurementsRequest',
  '2': [
    {
      '1': 'records',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.manpasik.api.v1.MeasurementRecord',
      '10': 'records'
    },
  ],
};

/// Descriptor for `SyncMeasurementsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List syncMeasurementsRequestDescriptor =
    $convert.base64Decode(
        'ChdTeW5jTWVhc3VyZW1lbnRzUmVxdWVzdBI8CgdyZWNvcmRzGAEgAygLMiIubWFucGFzaWsuYX'
        'BpLnYxLk1lYXN1cmVtZW50UmVjb3JkUgdyZWNvcmRz');

@$core.Deprecated('Use syncMeasurementsResponseDescriptor instead')
const SyncMeasurementsResponse$json = {
  '1': 'SyncMeasurementsResponse',
  '2': [
    {'1': 'synced_count', '3': 1, '4': 1, '5': 5, '10': 'syncedCount'},
    {'1': 'failed_ids', '3': 2, '4': 3, '5': 9, '10': 'failedIds'},
  ],
};

/// Descriptor for `SyncMeasurementsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List syncMeasurementsResponseDescriptor =
    $convert.base64Decode(
        'ChhTeW5jTWVhc3VyZW1lbnRzUmVzcG9uc2USIQoMc3luY2VkX2NvdW50GAEgASgFUgtzeW5jZW'
        'RDb3VudBIdCgpmYWlsZWRfaWRzGAIgAygJUglmYWlsZWRJZHM=');

@$core.Deprecated('Use measurementRecordDescriptor instead')
const MeasurementRecord$json = {
  '1': 'MeasurementRecord',
  '2': [
    {'1': 'id', '3': 1, '4': 1, '5': 9, '10': 'id'},
    {'1': 'device_mac', '3': 2, '4': 1, '5': 9, '10': 'deviceMac'},
    {'1': 'timestamp', '3': 3, '4': 1, '5': 3, '10': 'timestamp'},
    {'1': 'diff_signal', '3': 4, '4': 3, '5': 1, '10': 'diffSignal'},
    {'1': 'fingerprint', '3': 5, '4': 3, '5': 1, '10': 'fingerprint'},
    {'1': 'health_score', '3': 6, '4': 1, '5': 5, '10': 'healthScore'},
    {'1': 'risk_label', '3': 7, '4': 1, '5': 9, '10': 'riskLabel'},
  ],
};

/// Descriptor for `MeasurementRecord`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List measurementRecordDescriptor = $convert.base64Decode(
    'ChFNZWFzdXJlbWVudFJlY29yZBIOCgJpZBgBIAEoCVICaWQSHQoKZGV2aWNlX21hYxgCIAEoCV'
    'IJZGV2aWNlTWFjEhwKCXRpbWVzdGFtcBgDIAEoA1IJdGltZXN0YW1wEh8KC2RpZmZfc2lnbmFs'
    'GAQgAygBUgpkaWZmU2lnbmFsEiAKC2ZpbmdlcnByaW50GAUgAygBUgtmaW5nZXJwcmludBIhCg'
    'xoZWFsdGhfc2NvcmUYBiABKAVSC2hlYWx0aFNjb3JlEh0KCnJpc2tfbGFiZWwYByABKAlSCXJp'
    'c2tMYWJlbA==');

@$core.Deprecated('Use submitCheckupRequestDescriptor instead')
const SubmitCheckupRequest$json = {
  '1': 'SubmitCheckupRequest',
  '2': [
    {'1': 'session_id', '3': 1, '4': 1, '5': 9, '10': 'sessionId'},
    {'1': 'user_id', '3': 2, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'package_type', '3': 3, '4': 1, '5': 9, '10': 'packageType'},
    {
      '1': 'results',
      '3': 4,
      '4': 3,
      '5': 11,
      '6': '.manpasik.api.v1.MeasurementRecord',
      '10': 'results'
    },
  ],
};

/// Descriptor for `SubmitCheckupRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List submitCheckupRequestDescriptor = $convert.base64Decode(
    'ChRTdWJtaXRDaGVja3VwUmVxdWVzdBIdCgpzZXNzaW9uX2lkGAEgASgJUglzZXNzaW9uSWQSFw'
    'oHdXNlcl9pZBgCIAEoCVIGdXNlcklkEiEKDHBhY2thZ2VfdHlwZRgDIAEoCVILcGFja2FnZVR5'
    'cGUSPAoHcmVzdWx0cxgEIAMoCzIiLm1hbnBhc2lrLmFwaS52MS5NZWFzdXJlbWVudFJlY29yZF'
    'IHcmVzdWx0cw==');

@$core.Deprecated('Use submitCheckupResponseDescriptor instead')
const SubmitCheckupResponse$json = {
  '1': 'SubmitCheckupResponse',
  '2': [
    {'1': 'success', '3': 1, '4': 1, '5': 8, '10': 'success'},
    {'1': 'composite_score', '3': 2, '4': 1, '5': 1, '10': 'compositeScore'},
  ],
};

/// Descriptor for `SubmitCheckupResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List submitCheckupResponseDescriptor = $convert.base64Decode(
    'ChVTdWJtaXRDaGVja3VwUmVzcG9uc2USGAoHc3VjY2VzcxgBIAEoCFIHc3VjY2VzcxInCg9jb2'
    '1wb3NpdGVfc2NvcmUYAiABKAFSDmNvbXBvc2l0ZVNjb3Jl');

@$core.Deprecated('Use contextRequestDescriptor instead')
const ContextRequest$json = {
  '1': 'ContextRequest',
  '2': [
    {'1': 'user_id', '3': 1, '4': 1, '5': 9, '10': 'userId'},
    {'1': 'requested_cards', '3': 2, '4': 1, '5': 5, '10': 'requestedCards'},
  ],
};

/// Descriptor for `ContextRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List contextRequestDescriptor = $convert.base64Decode(
    'Cg5Db250ZXh0UmVxdWVzdBIXCgd1c2VyX2lkGAEgASgJUgZ1c2VySWQSJwoPcmVxdWVzdGVkX2'
    'NhcmRzGAIgASgFUg5yZXF1ZXN0ZWRDYXJkcw==');

@$core.Deprecated('Use contextCardResponseDescriptor instead')
const ContextCardResponse$json = {
  '1': 'ContextCardResponse',
  '2': [
    {'1': 'card_id', '3': 1, '4': 1, '5': 9, '10': 'cardId'},
    {'1': 'type', '3': 2, '4': 1, '5': 9, '10': 'type'},
    {'1': 'title', '3': 3, '4': 1, '5': 9, '10': 'title'},
    {'1': 'body', '3': 4, '4': 1, '5': 9, '10': 'body'},
    {'1': 'priority', '3': 5, '4': 1, '5': 5, '10': 'priority'},
  ],
};

/// Descriptor for `ContextCardResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List contextCardResponseDescriptor = $convert.base64Decode(
    'ChNDb250ZXh0Q2FyZFJlc3BvbnNlEhcKB2NhcmRfaWQYASABKAlSBmNhcmRJZBISCgR0eXBlGA'
    'IgASgJUgR0eXBlEhQKBXRpdGxlGAMgASgJUgV0aXRsZRISCgRib2R5GAQgASgJUgRib2R5EhoK'
    'CHByaW9yaXR5GAUgASgFUghwcmlvcml0eQ==');

@$core.Deprecated('Use healthReportRequestDescriptor instead')
const HealthReportRequest$json = {
  '1': 'HealthReportRequest',
  '2': [
    {'1': 'device_id', '3': 1, '4': 1, '5': 9, '10': 'deviceId'},
    {'1': 'overall_score', '3': 2, '4': 1, '5': 1, '10': 'overallScore'},
    {'1': 'hw_score', '3': 3, '4': 1, '5': 1, '10': 'hwScore'},
    {'1': 'rust_score', '3': 4, '4': 1, '5': 1, '10': 'rustScore'},
    {'1': 'timestamp', '3': 5, '4': 1, '5': 3, '10': 'timestamp'},
  ],
};

/// Descriptor for `HealthReportRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List healthReportRequestDescriptor = $convert.base64Decode(
    'ChNIZWFsdGhSZXBvcnRSZXF1ZXN0EhsKCWRldmljZV9pZBgBIAEoCVIIZGV2aWNlSWQSIwoNb3'
    'ZlcmFsbF9zY29yZRgCIAEoAVIMb3ZlcmFsbFNjb3JlEhkKCGh3X3Njb3JlGAMgASgBUgdod1Nj'
    'b3JlEh0KCnJ1c3Rfc2NvcmUYBCABKAFSCXJ1c3RTY29yZRIcCgl0aW1lc3RhbXAYBSABKANSCX'
    'RpbWVzdGFtcA==');

@$core.Deprecated('Use healthReportResponseDescriptor instead')
const HealthReportResponse$json = {
  '1': 'HealthReportResponse',
  '2': [
    {
      '1': 'requires_service_ticket',
      '3': 1,
      '4': 1,
      '5': 8,
      '10': 'requiresServiceTicket'
    },
    {'1': 'ticket_id', '3': 2, '4': 1, '5': 9, '10': 'ticketId'},
  ],
};

/// Descriptor for `HealthReportResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List healthReportResponseDescriptor = $convert.base64Decode(
    'ChRIZWFsdGhSZXBvcnRSZXNwb25zZRI2ChdyZXF1aXJlc19zZXJ2aWNlX3RpY2tldBgBIAEoCF'
    'IVcmVxdWlyZXNTZXJ2aWNlVGlja2V0EhsKCXRpY2tldF9pZBgCIAEoCVIIdGlja2V0SWQ=');

@$core.Deprecated('Use healingEventRequestDescriptor instead')
const HealingEventRequest$json = {
  '1': 'HealingEventRequest',
  '2': [
    {'1': 'device_id', '3': 1, '4': 1, '5': 9, '10': 'deviceId'},
    {'1': 'layer', '3': 2, '4': 1, '5': 9, '10': 'layer'},
    {'1': 'strategy', '3': 3, '4': 1, '5': 9, '10': 'strategy'},
    {'1': 'success', '3': 4, '4': 1, '5': 8, '10': 'success'},
    {'1': 'timestamp', '3': 5, '4': 1, '5': 3, '10': 'timestamp'},
  ],
};

/// Descriptor for `HealingEventRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List healingEventRequestDescriptor = $convert.base64Decode(
    'ChNIZWFsaW5nRXZlbnRSZXF1ZXN0EhsKCWRldmljZV9pZBgBIAEoCVIIZGV2aWNlSWQSFAoFbG'
    'F5ZXIYAiABKAlSBWxheWVyEhoKCHN0cmF0ZWd5GAMgASgJUghzdHJhdGVneRIYCgdzdWNjZXNz'
    'GAQgASgIUgdzdWNjZXNzEhwKCXRpbWVzdGFtcBgFIAEoA1IJdGltZXN0YW1w');

@$core.Deprecated('Use healingEventResponseDescriptor instead')
const HealingEventResponse$json = {
  '1': 'HealingEventResponse',
  '2': [
    {'1': 'ack', '3': 1, '4': 1, '5': 8, '10': 'ack'},
  ],
};

/// Descriptor for `HealingEventResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List healingEventResponseDescriptor = $convert
    .base64Decode('ChRIZWFsaW5nRXZlbnRSZXNwb25zZRIQCgNhY2sYASABKAhSA2Fjaw==');
