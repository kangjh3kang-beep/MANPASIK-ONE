import 'dart:async';
import 'dart:convert';
import 'package:flutter/foundation.dart';

import 'package:manpasik/core/config/app_config.dart';
import 'package:manpasik/core/services/rest_client.dart';

/// WebRTC 연결 상태
enum WebRtcConnectionState { disconnected, connecting, connected, failed }

/// WebRTC 시그널링 메시지 타입
enum SignalType { offer, answer, iceCandidate }

/// WebRTC P2P 연결 서비스 인터페이스
///
/// 화상 통화를 위한 WebRTC P2P 연결을 관리합니다.
/// flutter_webrtc 패키지 설치 시 실제 WebRTC, 미설치 시 REST 시뮬레이션.
abstract class WebRtcService {
  /// 현재 연결 상태
  WebRtcConnectionState get connectionState;

  /// 연결 상태 변경 스트림
  Stream<WebRtcConnectionState> get onStateChanged;

  /// P2P 연결 초기화
  Future<void> initialize({
    required String roomId,
    required String userId,
    required String token,
  });

  /// 로컬 미디어(카메라/마이크) 시작
  Future<void> startLocalMedia({bool audio = true, bool video = true});

  /// 음소거 토글
  Future<void> toggleMute(bool muted);

  /// 비디오 토글
  Future<void> toggleVideo(bool videoOff);

  /// 연결 종료 및 리소스 해제
  Future<void> dispose();

  /// 팩토리: 환경에 맞는 구현체 반환
  static WebRtcService create({required ManPaSikRestClient restClient}) {
    // flutter_webrtc 패키지 설치 시 RealWebRtcService 반환
    // 현재는 REST 기반 시그널링 서비스 사용
    return RestSignalingWebRtcService(restClient: restClient);
  }
}

/// REST 시그널링 기반 WebRTC 서비스
///
/// REST API를 통해 SDP offer/answer 및 ICE candidate를 교환합니다.
/// flutter_webrtc 패키지 설치 시 RTCPeerConnection을 사용하고,
/// 미설치 시 시그널링만 동작하는 시뮬레이션 모드.
class RestSignalingWebRtcService implements WebRtcService {
  RestSignalingWebRtcService({required this.restClient});

  final ManPaSikRestClient restClient;

  String _roomId = '';
  String _userId = '';
  String _token = '';
  Timer? _signalingPollTimer;

  WebRtcConnectionState _state = WebRtcConnectionState.disconnected;
  final _stateController =
      StreamController<WebRtcConnectionState>.broadcast();

  bool _isMuted = false;
  bool _isVideoOff = false;

  @override
  WebRtcConnectionState get connectionState => _state;

  @override
  Stream<WebRtcConnectionState> get onStateChanged => _stateController.stream;

  void _updateState(WebRtcConnectionState newState) {
    _state = newState;
    _stateController.add(newState);
  }

  @override
  Future<void> initialize({
    required String roomId,
    required String userId,
    required String token,
  }) async {
    _roomId = roomId;
    _userId = userId;
    _token = token;
    _updateState(WebRtcConnectionState.connecting);

    try {
      if (AppConfig.enableWebRtc && !kIsWeb) {
        await _initializeRealWebRtc();
      } else {
        await _initializeSimulated();
      }
    } catch (e) {
      debugPrint('[WebRTC] 초기화 실패, 시뮬레이션 모드: $e');
      await _initializeSimulated();
    }
  }

  /// 실제 WebRTC 초기화 (flutter_webrtc 패키지 설치 시)
  Future<void> _initializeRealWebRtc() async {
    // flutter_webrtc 패키지 설치 후 아래 코드 활성화:
    //
    // final config = {
    //   'iceServers': [
    //     {'urls': 'stun:stun.l.google.com:19302'},
    //     {
    //       'urls': 'turn:turn.manpasik.com:3478',
    //       'username': 'mpk',
    //       'credential': _token,
    //     },
    //   ],
    // };
    //
    // _peerConnection = await createPeerConnection(config);
    //
    // // 로컬 스트림
    // _localStream = await navigator.mediaDevices.getUserMedia({
    //   'audio': true,
    //   'video': {'facingMode': 'user', 'width': 640, 'height': 480},
    // });
    // _localStream!.getTracks().forEach((track) {
    //   _peerConnection!.addTrack(track, _localStream!);
    // });
    //
    // // 리모트 스트림 수신
    // _peerConnection!.onTrack = (event) {
    //   if (event.streams.isNotEmpty) {
    //     _remoteStream = event.streams[0];
    //   }
    // };
    //
    // // ICE candidate
    // _peerConnection!.onIceCandidate = (candidate) {
    //   _sendSignal(SignalType.iceCandidate, jsonEncode(candidate.toMap()));
    // };
    //
    // // 연결 상태 모니터링
    // _peerConnection!.onConnectionState = (state) {
    //   switch (state) {
    //     case RTCPeerConnectionState.RTCPeerConnectionStateConnected:
    //       _updateState(WebRtcConnectionState.connected);
    //     case RTCPeerConnectionState.RTCPeerConnectionStateFailed:
    //       _updateState(WebRtcConnectionState.failed);
    //     case RTCPeerConnectionState.RTCPeerConnectionStateDisconnected:
    //       _updateState(WebRtcConnectionState.disconnected);
    //     default:
    //       break;
    //   }
    // };

    // 시그널링 폴링 시작
    _startSignalingPoll();
    _updateState(WebRtcConnectionState.connected);
  }

  /// 시뮬레이션 모드 초기화
  Future<void> _initializeSimulated() async {
    debugPrint('[WebRTC:Sim] 시뮬레이션 모드로 연결: roomId=$_roomId');
    await Future.delayed(const Duration(milliseconds: 500));
    _startSignalingPoll();
    _updateState(WebRtcConnectionState.connected);
  }

  /// REST 폴링으로 시그널링 메시지 교환
  void _startSignalingPoll() {
    _signalingPollTimer?.cancel();
    _signalingPollTimer = Timer.periodic(
      const Duration(seconds: 2),
      (_) => _pollSignaling(),
    );
  }

  Future<void> _pollSignaling() async {
    try {
      final res = await restClient.getVideoSignals(
        roomId: _roomId,
        userId: _userId,
      );
      final signals = res['signals'] as List<dynamic>? ?? [];
      for (final signal in signals) {
        final s = signal as Map<String, dynamic>;
        final type = s['type'] as String? ?? '';
        final data = s['data'] as String? ?? '';
        debugPrint('[WebRTC] 시그널 수신: $type (${data.length} bytes)');

        // flutter_webrtc 설치 시 SDP/ICE 처리:
        // if (type == 'offer') {
        //   final desc = RTCSessionDescription(data, 'offer');
        //   await _peerConnection!.setRemoteDescription(desc);
        //   final answer = await _peerConnection!.createAnswer();
        //   await _peerConnection!.setLocalDescription(answer);
        //   _sendSignal(SignalType.answer, answer.sdp!);
        // } else if (type == 'answer') {
        //   final desc = RTCSessionDescription(data, 'answer');
        //   await _peerConnection!.setRemoteDescription(desc);
        // } else if (type == 'ice-candidate') {
        //   final candidate = RTCIceCandidate.fromMap(jsonDecode(data));
        //   await _peerConnection!.addCandidate(candidate);
        // }
      }
    } catch (e) {
      // 폴링 실패는 무시 (다음 주기에 재시도)
    }
  }

  Future<void> _sendSignal(SignalType type, String data) async {
    final typeName = switch (type) {
      SignalType.offer => 'offer',
      SignalType.answer => 'answer',
      SignalType.iceCandidate => 'ice-candidate',
    };
    try {
      await restClient.sendVideoSignal(
        roomId: _roomId,
        signalType: typeName,
        payload: data,
      );
    } catch (e) {
      debugPrint('[WebRTC] 시그널 전송 실패: $e');
    }
  }

  @override
  Future<void> startLocalMedia({bool audio = true, bool video = true}) async {
    // flutter_webrtc 설치 시:
    // _localStream = await navigator.mediaDevices.getUserMedia({
    //   'audio': audio,
    //   'video': video ? {'facingMode': 'user'} : false,
    // });
    debugPrint('[WebRTC] 로컬 미디어 시작: audio=$audio, video=$video');
  }

  @override
  Future<void> toggleMute(bool muted) async {
    _isMuted = muted;
    // flutter_webrtc 설치 시:
    // _localStream?.getAudioTracks().forEach((t) => t.enabled = !muted);
    debugPrint('[WebRTC] 음소거: $muted');
  }

  @override
  Future<void> toggleVideo(bool videoOff) async {
    _isVideoOff = videoOff;
    // flutter_webrtc 설치 시:
    // _localStream?.getVideoTracks().forEach((t) => t.enabled = !videoOff);
    debugPrint('[WebRTC] 비디오 끄기: $videoOff');
  }

  @override
  Future<void> dispose() async {
    _signalingPollTimer?.cancel();
    // flutter_webrtc 설치 시:
    // _localStream?.getTracks().forEach((t) => t.stop());
    // _localStream?.dispose();
    // _peerConnection?.close();
    _updateState(WebRtcConnectionState.disconnected);
    await _stateController.close();
  }
}
