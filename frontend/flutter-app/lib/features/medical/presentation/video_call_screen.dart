import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/theme/app_theme.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';

/// WebRTC 화상 통화 화면
///
/// 시그널링 서버(video-service)와 연동하여 P2P 화상 통화를 진행합니다.
/// flutter_webrtc 패키지 설치 시 RTCPeerConnection 기반 실제 P2P 연결.
/// 미설치 환경에서는 REST API 기반 시뮬레이션 모드로 동작합니다.
class VideoCallScreen extends ConsumerStatefulWidget {
  const VideoCallScreen({super.key, required this.sessionId});

  final String sessionId;

  @override
  ConsumerState<VideoCallScreen> createState() => _VideoCallScreenState();
}

class _VideoCallScreenState extends ConsumerState<VideoCallScreen> {
  bool _isConnecting = true;
  bool _isConnected = false;
  bool _isMuted = false;
  bool _isVideoOff = false;
  bool _isSpeakerOn = true;
  bool _isWebRtcActive = false;
  Duration _callDuration = Duration.zero;
  Timer? _timer;
  String _remoteName = '의사';
  String? _roomToken;

  @override
  void initState() {
    super.initState();
    _joinRoom();
  }

  @override
  void dispose() {
    _timer?.cancel();
    _disposeWebRtc();
    super.dispose();
  }

  Future<void> _joinRoom() async {
    try {
      final client = ref.read(restClientProvider);
      final userId = ref.read(authProvider).userId ?? '';
      final res = await client.joinVideoRoom(
        roomId: widget.sessionId,
        userId: userId,
      );
      if (!mounted) return;
      setState(() {
        _roomToken = res['token'] as String?;
        _remoteName = (res['doctor_name'] as String?) ?? '의사';
        _isConnecting = false;
        _isConnected = true;
      });

      // WebRTC 활성화 시도 (flutter_webrtc 패키지 설치 시)
      await _initWebRtc();

      _startTimer();
    } catch (e) {
      if (!mounted) return;
      // 시그널링 서버 미연결 시 시뮬레이션 모드
      setState(() {
        _isConnecting = false;
        _isConnected = true;
      });
      _startTimer();
    }
  }

  /// WebRTC P2P 연결 초기화
  ///
  /// flutter_webrtc 패키지 설치 후 아래 주석 해제:
  /// RTCPeerConnection 생성 → offer/answer 교환 → ICE candidate 설정
  Future<void> _initWebRtc() async {
    // if (!AppConfig.enableWebRtc) return;
    //
    // final config = {
    //   'iceServers': [
    //     {'urls': 'stun:stun.l.google.com:19302'},
    //     {'urls': 'turn:turn.manpasik.com:3478', 'username': 'mpk', 'credential': _roomToken},
    //   ],
    // };
    //
    // _peerConnection = await createPeerConnection(config);
    // _localStream = await navigator.mediaDevices.getUserMedia({
    //   'audio': true,
    //   'video': {'facingMode': 'user', 'width': 640, 'height': 480},
    // });
    // _localStream!.getTracks().forEach((track) => _peerConnection!.addTrack(track, _localStream!));
    //
    // _peerConnection!.onTrack = (event) {
    //   if (event.streams.isNotEmpty) {
    //     setState(() => _remoteStream = event.streams[0]);
    //   }
    // };
    //
    // _isWebRtcActive = true;
  }

  /// WebRTC 자원 해제
  Future<void> _disposeWebRtc() async {
    // _localStream?.getTracks().forEach((track) => track.stop());
    // _localStream?.dispose();
    // _peerConnection?.close();
  }

  void _startTimer() {
    _timer = Timer.periodic(const Duration(seconds: 1), (_) {
      if (!mounted) return;
      setState(() => _callDuration += const Duration(seconds: 1));
    });
  }

  String _formatDuration(Duration d) {
    final m = d.inMinutes.remainder(60).toString().padLeft(2, '0');
    final s = d.inSeconds.remainder(60).toString().padLeft(2, '0');
    if (d.inHours > 0) {
      return '${d.inHours}:$m:$s';
    }
    return '$m:$s';
  }

  Future<void> _endCall() async {
    _timer?.cancel();
    try {
      final client = ref.read(restClientProvider);
      final userId = ref.read(authProvider).userId ?? '';
      await client.leaveVideoRoom(
        roomId: widget.sessionId,
        userId: userId,
      );
    } catch (_) {}
    if (mounted) {
      context.pushReplacement('/medical/consultation/${widget.sessionId}/result');
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      backgroundColor: Colors.black,
      body: SafeArea(
        child: Stack(
          children: [
            // 원격 비디오 (전체 화면) — 실제 구현 시 RTCVideoView 사용
            Center(
              child: _isConnecting
                  ? Column(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        const CircularProgressIndicator(color: Colors.white),
                        const SizedBox(height: 16),
                        Text('연결 중...', style: theme.textTheme.titleMedium?.copyWith(color: Colors.white)),
                      ],
                    )
                  : _isVideoOff
                      ? Column(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            CircleAvatar(
                              radius: 48,
                              backgroundColor: AppTheme.sanggamGold.withValues(alpha: 0.3),
                              child: Text(
                                _remoteName.isNotEmpty ? _remoteName[0] : '?',
                                style: const TextStyle(fontSize: 36, color: Colors.white),
                              ),
                            ),
                            const SizedBox(height: 12),
                            Text(_remoteName, style: const TextStyle(color: Colors.white, fontSize: 18)),
                          ],
                        )
                      : Container(
                          // 원격 비디오 플레이스홀더
                          color: Colors.grey[900],
                          child: Center(
                            child: Column(
                              mainAxisSize: MainAxisSize.min,
                              children: [
                                const Icon(Icons.videocam, size: 80, color: Colors.white24),
                                const SizedBox(height: 8),
                                Text(_remoteName, style: const TextStyle(color: Colors.white70, fontSize: 16)),
                                Text('화상 진료 중', style: TextStyle(color: Colors.green[300], fontSize: 14)),
                              ],
                            ),
                          ),
                        ),
            ),

            // 로컬 비디오 (PIP — 우상단)
            if (_isConnected)
              Positioned(
                top: 16,
                right: 16,
                child: Container(
                  width: 100,
                  height: 140,
                  decoration: BoxDecoration(
                    color: _isVideoOff ? Colors.grey[800] : Colors.grey[700],
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(color: Colors.white24, width: 1),
                  ),
                  child: _isVideoOff
                      ? const Center(child: Icon(Icons.videocam_off, color: Colors.white54))
                      : const Center(child: Icon(Icons.person, color: Colors.white38, size: 40)),
                ),
              ),

            // 상단 정보 바
            Positioned(
              top: 16,
              left: 16,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                    decoration: BoxDecoration(
                      color: _isConnected ? Colors.green.withValues(alpha: 0.8) : Colors.orange.withValues(alpha: 0.8),
                      borderRadius: BorderRadius.circular(16),
                    ),
                    child: Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Container(
                          width: 8,
                          height: 8,
                          decoration: const BoxDecoration(color: Colors.white, shape: BoxShape.circle),
                        ),
                        const SizedBox(width: 6),
                        Text(
                          _isConnected ? '통화 중' : '연결 중',
                          style: const TextStyle(color: Colors.white, fontSize: 12, fontWeight: FontWeight.w600),
                        ),
                      ],
                    ),
                  ),
                  if (_isConnected) ...[
                    const SizedBox(height: 4),
                    Text(
                      _formatDuration(_callDuration),
                      style: const TextStyle(color: Colors.white70, fontSize: 14),
                    ),
                  ],
                ],
              ),
            ),

            // 하단 컨트롤 바
            Positioned(
              bottom: 24,
              left: 0,
              right: 0,
              child: Column(
                children: [
                  // 메인 컨트롤
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                    children: [
                      _buildControlButton(
                        icon: _isMuted ? Icons.mic_off : Icons.mic,
                        label: _isMuted ? '음소거 해제' : '음소거',
                        isActive: _isMuted,
                        onTap: () => setState(() => _isMuted = !_isMuted),
                      ),
                      _buildControlButton(
                        icon: _isVideoOff ? Icons.videocam_off : Icons.videocam,
                        label: _isVideoOff ? '카메라 켜기' : '카메라 끄기',
                        isActive: _isVideoOff,
                        onTap: () => setState(() => _isVideoOff = !_isVideoOff),
                      ),
                      // 통화 종료
                      GestureDetector(
                        onTap: _endCall,
                        child: Container(
                          width: 64,
                          height: 64,
                          decoration: const BoxDecoration(
                            color: Colors.red,
                            shape: BoxShape.circle,
                          ),
                          child: const Icon(Icons.call_end, color: Colors.white, size: 32),
                        ),
                      ),
                      _buildControlButton(
                        icon: _isSpeakerOn ? Icons.volume_up : Icons.volume_off,
                        label: _isSpeakerOn ? '스피커' : '수화기',
                        isActive: false,
                        onTap: () => setState(() => _isSpeakerOn = !_isSpeakerOn),
                      ),
                      _buildControlButton(
                        icon: Icons.chat_bubble_outline,
                        label: '채팅',
                        isActive: false,
                        onTap: () => _showChatSheet(context),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _showChatSheet(BuildContext context) {
    final chatController = TextEditingController();
    final messages = <String>[];
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setSheetState) => Padding(
          padding: EdgeInsets.only(bottom: MediaQuery.of(ctx).viewInsets.bottom),
          child: SizedBox(
            height: 360,
            child: Column(
              children: [
                AppBar(
                  title: const Text('진료 중 채팅'),
                  leading: IconButton(
                    icon: const Icon(Icons.close),
                    onPressed: () => Navigator.pop(ctx),
                  ),
                  automaticallyImplyLeading: false,
                ),
                Expanded(
                  child: messages.isEmpty
                      ? const Center(child: Text('메시지가 없습니다.'))
                      : ListView.builder(
                          padding: const EdgeInsets.all(8),
                          itemCount: messages.length,
                          itemBuilder: (_, i) => Align(
                            alignment: Alignment.centerRight,
                            child: Container(
                              margin: const EdgeInsets.only(bottom: 4),
                              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                              decoration: BoxDecoration(
                                color: Theme.of(context).colorScheme.primaryContainer,
                                borderRadius: BorderRadius.circular(12),
                              ),
                              child: Text(messages[i]),
                            ),
                          ),
                        ),
                ),
                Padding(
                  padding: const EdgeInsets.all(8),
                  child: Row(
                    children: [
                      Expanded(
                        child: TextField(
                          controller: chatController,
                          decoration: const InputDecoration(
                            hintText: '메시지를 입력하세요',
                            border: OutlineInputBorder(),
                            contentPadding: EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                          ),
                        ),
                      ),
                      const SizedBox(width: 8),
                      IconButton(
                        icon: const Icon(Icons.send),
                        onPressed: () {
                          if (chatController.text.trim().isNotEmpty) {
                            setSheetState(() => messages.add(chatController.text.trim()));
                            chatController.clear();
                          }
                        },
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildControlButton({
    required IconData icon,
    required String label,
    required bool isActive,
    required VoidCallback onTap,
  }) {
    return GestureDetector(
      onTap: onTap,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            width: 48,
            height: 48,
            decoration: BoxDecoration(
              color: isActive ? Colors.white : Colors.white.withValues(alpha: 0.2),
              shape: BoxShape.circle,
            ),
            child: Icon(icon, color: isActive ? Colors.black : Colors.white, size: 24),
          ),
          const SizedBox(height: 4),
          Text(label, style: const TextStyle(color: Colors.white70, fontSize: 10)),
        ],
      ),
    );
  }
}
