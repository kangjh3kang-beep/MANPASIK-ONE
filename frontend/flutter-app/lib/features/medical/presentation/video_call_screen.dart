import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'package:manpasik/core/providers/grpc_provider.dart';
import 'package:manpasik/core/services/webrtc_service.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/shared/providers/auth_provider.dart';
import 'package:manpasik/features/medical/presentation/widgets/vital_signs_hud.dart';

// ───────────────────────────────────────────────────
// VideoCallScreen — Sanggam Orbit 화상 통화
//
// [Rule 4] app_theme.dart → sanggam_theme.dart
// [Rule 4] Theme.of(context) 2x 제거
// [Rule 4] theme.textTheme.* → 직접 TextStyle
// [Rule 4] theme.colorScheme.* → SanggamTheme 상수
// [Rule 4] AppTheme.sanggamGold → SanggamTheme.primary
// [Rule 4] Colors.green → SanggamTheme.jagaeCyan
// [Rule 4] Colors.orange → SanggamTheme.primary
// [Rule 4] Colors.red → SanggamTheme.error
// [Rule 4] Colors.grey[*] → SanggamTheme.surface/surfaceVariant
// [Rule 4] Colors.black → SanggamTheme.background
// [Rule 4] AppBar (chatSheet) → 커스텀 헤더
// [Rule 2] borderRadius:12→16, padding:10→16
// ───────────────────────────────────────────────────

/// WebRTC 화상 통화 화면
///
/// 시그널링 서버(video-service)와 연동하여 P2P 화상 통화를 진행합니다.
/// flutter_webrtc 패키지 설치 시 RTCPeerConnection 기반 실제 P2P 연결.
/// 미설치 환경에서는 REST API 기반 시뮬레이션 모드로 동작합니다.
class VideoCallScreen extends ConsumerStatefulWidget {
  const VideoCallScreen(
      {super.key, required this.sessionId});

  final String sessionId;

  @override
  ConsumerState<VideoCallScreen> createState() =>
      _VideoCallScreenState();
}

class _VideoCallScreenState
    extends ConsumerState<VideoCallScreen> {
  bool _isConnecting = true;
  bool _isConnected = false;
  bool _isMuted = false;
  bool _isVideoOff = false;
  bool _isSpeakerOn = true;
  bool _isHudExpanded = false;
  bool _showHud = true;
  Duration _callDuration = Duration.zero;
  Timer? _timer;
  String _remoteName = '의사';
  String? _roomToken;
  WebRtcService? _webRtcService;

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
      final userId =
          ref.read(authProvider).userId ?? '';
      final res = await client.joinVideoRoom(
        roomId: widget.sessionId,
        userId: userId,
      );
      if (!mounted) return;
      setState(() {
        _roomToken = res['token'] as String?;
        _remoteName =
            (res['doctor_name'] as String?) ?? '의사';
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
  Future<void> _initWebRtc() async {
    final client = ref.read(restClientProvider);
    final userId =
        ref.read(authProvider).userId ?? '';
    _webRtcService =
        WebRtcService.create(restClient: client);

    await _webRtcService!.initialize(
      roomId: widget.sessionId,
      userId: userId,
      token: _roomToken ?? '',
    );

    await _webRtcService!.startLocalMedia();

    _webRtcService!.onStateChanged.listen((state) {
      if (!mounted) return;
      if (state == WebRtcConnectionState.failed) {
        setState(() => _isConnected = false);
      }
    });
  }

  /// WebRTC 자원 해제
  Future<void> _disposeWebRtc() async {
    await _webRtcService?.dispose();
  }

  void _startTimer() {
    _timer = Timer.periodic(
        const Duration(seconds: 1), (_) {
      if (!mounted) return;
      setState(() => _callDuration +=
          const Duration(seconds: 1));
    });
  }

  String _formatDuration(Duration d) {
    final m = d.inMinutes
        .remainder(60)
        .toString()
        .padLeft(2, '0');
    final s = d.inSeconds
        .remainder(60)
        .toString()
        .padLeft(2, '0');
    if (d.inHours > 0) {
      return '${d.inHours}:$m:$s';
    }
    return '$m:$s';
  }

  Future<void> _endCall() async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        backgroundColor: SanggamTheme.surface,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(16),
        ),
        title: const Text('통화 종료',
            style: TextStyle(color: Colors.white)),
        content: const Text('화상 상담을 종료하시겠습니까?',
            style: TextStyle(
                color: SanggamTheme.onSurfaceDim)),
        actions: [
          TextButton(
            onPressed: () =>
                Navigator.of(ctx).pop(false),
            child: const Text('계속 통화',
                style: TextStyle(
                    color:
                        SanggamTheme.onSurfaceDim)),
          ),
          FilledButton(
            onPressed: () =>
                Navigator.of(ctx).pop(true),
            style: FilledButton.styleFrom(
              backgroundColor: SanggamTheme.error,
              foregroundColor: Colors.white,
              shape: RoundedRectangleBorder(
                borderRadius:
                    BorderRadius.circular(16),
              ),
            ),
            child: const Text('종료'),
          ),
        ],
      ),
    );
    if (confirmed != true) return;

    _timer?.cancel();
    try {
      final client = ref.read(restClientProvider);
      final userId =
          ref.read(authProvider).userId ?? '';
      await client.leaveVideoRoom(
        roomId: widget.sessionId,
        userId: userId,
      );
    } catch (_) {}
    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
            content: Text('상담이 종료되었습니다'),
            duration: Duration(seconds: 2)),
      );
      context.pushReplacement(
          '/medical/consultation/${widget.sessionId}/result');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: SanggamTheme.background,
      body: SafeArea(
        child: Stack(
          children: [
            // 원격 비디오 (전체 화면)
            Center(
              child: _isConnecting
                  ? const Column(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        CircularProgressIndicator(
                            color: Colors.white),
                        SizedBox(height: 16),
                        Text('연결 중...',
                            style: TextStyle(
                              color: Colors.white,
                              fontSize: 16,
                              fontWeight:
                                  FontWeight.w500,
                            )),
                      ],
                    )
                  : _isVideoOff
                      ? Column(
                          mainAxisSize:
                              MainAxisSize.min,
                          children: [
                            CircleAvatar(
                              radius: 48,
                              backgroundColor:
                                  SanggamTheme
                                      .primary
                                      .withValues(
                                          alpha:
                                              0.3),
                              child: Text(
                                _remoteName
                                        .isNotEmpty
                                    ? _remoteName[0]
                                    : '?',
                                style:
                                    const TextStyle(
                                        fontSize:
                                            36,
                                        color: Colors
                                            .white),
                              ),
                            ),
                            const SizedBox(
                                height: 16),
                            Text(_remoteName,
                                style:
                                    const TextStyle(
                                        color: Colors
                                            .white,
                                        fontSize:
                                            18)),
                          ],
                        )
                      : Container(
                          color:
                              SanggamTheme.surface,
                          child: Center(
                            child: Column(
                              mainAxisSize:
                                  MainAxisSize.min,
                              children: [
                                Icon(
                                    Icons.videocam,
                                    size: 80,
                                    color: Colors
                                        .white
                                        .withValues(
                                            alpha:
                                                0.15)),
                                const SizedBox(
                                    height: 8),
                                Text(_remoteName,
                                    style: TextStyle(
                                        color: Colors
                                            .white
                                            .withValues(
                                                alpha:
                                                    0.7),
                                        fontSize:
                                            16)),
                                const Text(
                                    '화상 진료 중',
                                    style:
                                        TextStyle(
                                      color: SanggamTheme
                                          .jagaeCyan,
                                      fontSize: 14,
                                    )),
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
                    color: _isVideoOff
                        ? SanggamTheme.surface
                        : SanggamTheme
                            .surfaceVariant,
                    borderRadius:
                        BorderRadius.circular(16),
                    border: Border.all(
                        color: Colors.white
                            .withValues(
                                alpha: 0.15),
                        width: 1),
                  ),
                  child: _isVideoOff
                      ? Icon(Icons.videocam_off,
                          color: Colors.white
                              .withValues(
                                  alpha: 0.5))
                      : Icon(Icons.person,
                          color: Colors.white
                              .withValues(
                                  alpha: 0.3),
                          size: 40),
                ),
              ),

            // 상단 정보 바
            Positioned(
              top: 16,
              left: 16,
              child: Column(
                crossAxisAlignment:
                    CrossAxisAlignment.start,
                children: [
                  Container(
                    padding:
                        const EdgeInsets.symmetric(
                            horizontal: 16,
                            vertical: 8),
                    decoration: BoxDecoration(
                      color: _isConnected
                          ? SanggamTheme.jagaeCyan
                              .withValues(
                                  alpha: 0.8)
                          : SanggamTheme.primary
                              .withValues(
                                  alpha: 0.8),
                      borderRadius:
                          BorderRadius.circular(16),
                    ),
                    child: Row(
                      mainAxisSize:
                          MainAxisSize.min,
                      children: [
                        Container(
                          width: 8,
                          height: 8,
                          decoration:
                              const BoxDecoration(
                                  color:
                                      Colors.white,
                                  shape: BoxShape
                                      .circle),
                        ),
                        const SizedBox(width: 8),
                        Text(
                          _isConnected
                              ? '통화 중'
                              : '연결 중',
                          style: const TextStyle(
                              color: Colors.white,
                              fontSize: 12,
                              fontWeight:
                                  FontWeight.w600),
                        ),
                      ],
                    ),
                  ),
                  if (_isConnected) ...[
                    const SizedBox(height: 8),
                    Text(
                      _formatDuration(
                          _callDuration),
                      style: TextStyle(
                          color: Colors.white
                              .withValues(
                                  alpha: 0.7),
                          fontSize: 14),
                    ),
                  ],
                ],
              ),
            ),

            // 생체 신호 HUD 오버레이
            if (_isConnected && _showHud)
              Positioned(
                bottom: 110,
                left: 0,
                right: 0,
                child: VitalSignsHud(
                  isExpanded: _isHudExpanded,
                  onToggle: () => setState(() =>
                      _isHudExpanded =
                          !_isHudExpanded),
                ),
              ),

            // 하단 컨트롤 바
            Positioned(
              bottom: 24,
              left: 0,
              right: 0,
              child: Column(
                children: [
                  Row(
                    mainAxisAlignment:
                        MainAxisAlignment
                            .spaceEvenly,
                    children: [
                      _buildControlButton(
                        icon: _isMuted
                            ? Icons.mic_off
                            : Icons.mic,
                        label: _isMuted
                            ? '음소거 해제'
                            : '음소거',
                        isActive: _isMuted,
                        onTap: () {
                          setState(() => _isMuted =
                              !_isMuted);
                          _webRtcService
                              ?.toggleMute(
                                  _isMuted);
                        },
                      ),
                      _buildControlButton(
                        icon: _isVideoOff
                            ? Icons.videocam_off
                            : Icons.videocam,
                        label: _isVideoOff
                            ? '카메라 켜기'
                            : '카메라 끄기',
                        isActive: _isVideoOff,
                        onTap: () {
                          setState(() =>
                              _isVideoOff =
                                  !_isVideoOff);
                          _webRtcService
                              ?.toggleVideo(
                                  _isVideoOff);
                        },
                      ),
                      // 통화 종료
                      Tooltip(
                        message: '통화 종료',
                        child: GestureDetector(
                          onTap: _endCall,
                          child: Container(
                            width: 64,
                            height: 64,
                            decoration:
                                const BoxDecoration(
                              color: SanggamTheme
                                  .error,
                              shape:
                                  BoxShape.circle,
                            ),
                            child: const Icon(
                                Icons.call_end,
                                color: Colors.white,
                                size: 32),
                          ),
                        ),
                      ),
                      _buildControlButton(
                        icon: _isSpeakerOn
                            ? Icons.volume_up
                            : Icons.volume_off,
                        label: _isSpeakerOn
                            ? '스피커'
                            : '수화기',
                        isActive: false,
                        onTap: () => setState(() =>
                            _isSpeakerOn =
                                !_isSpeakerOn),
                      ),
                      _buildControlButton(
                        icon: Icons
                            .chat_bubble_outline,
                        label: '채팅',
                        isActive: false,
                        onTap: () =>
                            _showChatSheet(context),
                      ),
                      _buildControlButton(
                        icon: _showHud
                            ? Icons.monitor_heart
                            : Icons
                                .monitor_heart_outlined,
                        label: _showHud
                            ? '바이탈'
                            : '바이탈 끔',
                        isActive: _showHud,
                        onTap: () =>
                            setState(() {
                          _showHud = !_showHud;
                          if (!_showHud) {
                            _isHudExpanded = false;
                          }
                        }),
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
      backgroundColor: SanggamTheme.surface,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(
            top: Radius.circular(16)),
      ),
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setSheetState) => Padding(
          padding: EdgeInsets.only(
              bottom: MediaQuery.of(ctx)
                  .viewInsets
                  .bottom),
          child: SizedBox(
            height: 360,
            child: Column(
              children: [
                // 커스텀 헤더
                Padding(
                  padding:
                      const EdgeInsets.symmetric(
                          horizontal: 8,
                          vertical: 8),
                  child: Row(
                    children: [
                      IconButton(
                        icon: const Icon(
                            Icons.close,
                            color: Colors.white),
                        tooltip: '닫기',
                        onPressed: () =>
                            Navigator.pop(ctx),
                      ),
                      const SizedBox(width: 8),
                      const Expanded(
                        child: Text(
                          '진료 중 채팅',
                          style: TextStyle(
                            color: Colors.white,
                            fontSize: 16,
                            fontWeight:
                                FontWeight.bold,
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
                const Divider(
                    height: 1,
                    color: SanggamTheme
                        .surfaceVariant),
                Expanded(
                  child: messages.isEmpty
                      ? const Center(
                          child: Text(
                              '메시지가 없습니다.',
                              style: TextStyle(
                                  color: SanggamTheme
                                      .onSurfaceDim)))
                      : ListView.builder(
                          padding:
                              const EdgeInsets.all(
                                  8),
                          itemCount:
                              messages.length,
                          itemBuilder: (_, i) =>
                              Align(
                            alignment: Alignment
                                .centerRight,
                            child: Container(
                              margin:
                                  const EdgeInsets
                                      .only(
                                      bottom: 8),
                              padding:
                                  const EdgeInsets
                                      .symmetric(
                                      horizontal:
                                          16,
                                      vertical: 8),
                              decoration:
                                  BoxDecoration(
                                color: SanggamTheme
                                    .primary
                                    .withValues(
                                        alpha:
                                            0.15),
                                borderRadius:
                                    BorderRadius
                                        .circular(
                                            16),
                              ),
                              child: Text(
                                  messages[i],
                                  style: const TextStyle(
                                      color: Colors
                                          .white)),
                            ),
                          ),
                        ),
                ),
                Padding(
                  padding:
                      const EdgeInsets.all(8),
                  child: Row(
                    children: [
                      Expanded(
                        child: TextField(
                          controller:
                              chatController,
                          style: const TextStyle(
                              color: Colors.white),
                          decoration:
                              InputDecoration(
                            hintText: '메시지를 입력하세요',
                            hintStyle: const TextStyle(
                                color: SanggamTheme
                                    .onSurfaceDim),
                            filled: true,
                            fillColor: SanggamTheme
                                .surfaceVariant,
                            border:
                                OutlineInputBorder(
                              borderRadius:
                                  BorderRadius
                                      .circular(
                                          16),
                              borderSide:
                                  BorderSide.none,
                            ),
                            contentPadding:
                                const EdgeInsets
                                    .symmetric(
                                    horizontal: 16,
                                    vertical: 8),
                          ),
                        ),
                      ),
                      const SizedBox(width: 8),
                      IconButton(
                        icon: const Icon(
                            Icons.send,
                            color: SanggamTheme
                                .primary),
                        tooltip: '메시지 전송',
                        onPressed: () {
                          if (chatController.text
                              .trim()
                              .isNotEmpty) {
                            setSheetState(() =>
                                messages.add(
                                    chatController
                                        .text
                                        .trim()));
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
    return Tooltip(
      message: label,
      child: GestureDetector(
        onTap: onTap,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              width: 48,
              height: 48,
              decoration: BoxDecoration(
                color: isActive
                    ? Colors.white
                    : Colors.white
                        .withValues(alpha: 0.2),
                shape: BoxShape.circle,
              ),
              child: Icon(icon,
                  color: isActive
                      ? SanggamTheme.background
                      : Colors.white,
                  size: 24),
            ),
            const SizedBox(height: 4),
            Text(label,
                style: TextStyle(
                    color: Colors.white
                        .withValues(alpha: 0.7),
                    fontSize: 10)),
          ],
        ),
      ),
    );
  }
}
