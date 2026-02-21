import 'dart:async';
import 'dart:convert';
import 'dart:typed_data';

import 'package:camera/camera.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:permission_handler/permission_handler.dart';

import 'package:manpasik/core/theme/app_theme.dart';

// ─────────────────────────────────────────────
// 1. 데이터 모델
// ─────────────────────────────────────────────

/// Rust api::AnalysisResult 에 대응하는 Dart 모델
class VisionAnalysisResult {
  final bool detected;
  final String biomarker;
  final double value;
  final double confidence;
  final List<int> dominantRgb;
  final String warning;

  const VisionAnalysisResult({
    required this.detected,
    required this.biomarker,
    required this.value,
    required this.confidence,
    required this.dominantRgb,
    required this.warning,
  });

  factory VisionAnalysisResult.fromJson(Map<String, dynamic> json) {
    return VisionAnalysisResult(
      detected: json['detected'] as bool? ?? false,
      biomarker: json['biomarker'] as String? ?? '',
      value: (json['value'] as num?)?.toDouble() ?? 0.0,
      confidence: (json['confidence'] as num?)?.toDouble() ?? 0.0,
      dominantRgb: (json['dominant_rgb'] as List<dynamic>?)
              ?.map((e) => (e as num).toInt())
              .toList() ??
          [0, 0, 0],
      warning: json['warning'] as String? ?? '',
    );
  }

  /// 신뢰도를 백분율 문자열로 반환
  String get confidencePercent => '${(confidence * 100).toStringAsFixed(1)}%';

  /// 감지된 색상을 Color로 변환
  Color get dominantColor => dominantRgb.length >= 3
      ? Color.fromRGBO(dominantRgb[0], dominantRgb[1], dominantRgb[2], 1.0)
      : Colors.grey;
}

/// 비전 분석기 상태
enum VisionAnalyzerStatus {
  /// 초기 상태 — 카메라 미활성
  idle,

  /// 카메라 권한 요청 중
  requestingPermission,

  /// 카메라 권한 거부됨
  permissionDenied,

  /// 카메라 권한 영구 거부됨 (설정 이동 필요)
  permissionPermanentlyDenied,

  /// 카메라 초기화 중
  initializing,

  /// 카메라 스트림 수신 중 (분석 실행 중)
  streaming,

  /// 분석 완료 (결과 확정)
  completed,

  /// 오류 발생
  error,
}

/// 비전 분석기 전체 상태
class VisionAnalyzerState {
  final VisionAnalyzerStatus status;
  final VisionAnalysisResult? latestResult;
  final int framesProcessed;
  final String? errorMessage;
  final double fps;

  const VisionAnalyzerState({
    this.status = VisionAnalyzerStatus.idle,
    this.latestResult,
    this.framesProcessed = 0,
    this.errorMessage,
    this.fps = 0.0,
  });

  VisionAnalyzerState copyWith({
    VisionAnalyzerStatus? status,
    VisionAnalysisResult? latestResult,
    int? framesProcessed,
    String? errorMessage,
    double? fps,
  }) {
    return VisionAnalyzerState(
      status: status ?? this.status,
      latestResult: latestResult ?? this.latestResult,
      framesProcessed: framesProcessed ?? this.framesProcessed,
      errorMessage: errorMessage ?? this.errorMessage,
      fps: fps ?? this.fps,
    );
  }
}

// ─────────────────────────────────────────────
// 2. Rust FFI 브릿지 (시뮬레이션 폴백 포함)
// ─────────────────────────────────────────────

/// Rust process_camera_frame FFI 호출 래퍼
///
/// flutter_rust_bridge 빌드가 완료되면 네이티브 호출을 활성화합니다.
/// 빌드 전까지는 시뮬레이션 모드로 동작합니다.
class VisionRustBridge {
  VisionRustBridge._();
  static final VisionRustBridge instance = VisionRustBridge._();

  /// 네이티브 Rust 엔진 사용 가능 여부
  bool _useNative = false;
  bool get isNativeEnabled => _useNative;

  /// 초기화
  Future<void> init() async {
    // flutter_rust_bridge 빌드 완료 후 아래 주석 해제:
    // try {
    //   await RustLib.init();
    //   _useNative = true;
    // } catch (_) {
    //   _useNative = false;
    // }
    _useNative = false;
  }

  /// 카메라 프레임 분석 (Rust FFI)
  ///
  /// [bytes]: BGRA8888 원시 픽셀 데이터
  /// [width]: 프레임 너비
  /// [height]: 프레임 높이
  ///
  /// Returns: JSON 문자열 형태의 AnalysisResult
  Future<String> processCameraFrame(
    Uint8List bytes,
    int width,
    int height,
  ) async {
    if (_useNative) {
      // flutter_rust_bridge 빌드 완료 후:
      // return await api.processCameraFrame(
      //   bytes: bytes,
      //   width: width,
      //   height: height,
      // );
    }

    // 시뮬레이션 모드: 픽셀 데이터에서 간단한 색상 분석 수행
    return _simulateAnalysis(bytes, width, height);
  }

  /// 시뮬레이션 분석 — 중앙 ROI의 평균 RGB로 결과 생성
  String _simulateAnalysis(Uint8List bytes, int width, int height) {
    if (bytes.isEmpty || width == 0 || height == 0) {
      return jsonEncode({
        'detected': false,
        'biomarker': '',
        'value': 0.0,
        'confidence': 0.0,
        'dominant_rgb': [0, 0, 0],
        'warning': '프레임 데이터가 비어있습니다.',
      });
    }

    // 중앙 ROI 영역의 평균 RGB 계산
    final roiX = (width * 0.3).toInt();
    final roiY = (height * 0.3).toInt();
    final roiW = (width * 0.4).toInt();
    final roiH = (height * 0.4).toInt();

    int rSum = 0, gSum = 0, bSum = 0;
    int count = 0;
    final bytesPerRow = width * 4; // BGRA

    for (int dy = 0; dy < roiH && (roiY + dy) < height; dy++) {
      for (int dx = 0; dx < roiW && (roiX + dx) < width; dx++) {
        final offset = (roiY + dy) * bytesPerRow + (roiX + dx) * 4;
        if (offset + 3 < bytes.length) {
          // BGRA 순서
          bSum += bytes[offset];
          gSum += bytes[offset + 1];
          rSum += bytes[offset + 2];
          count++;
        }
      }
    }

    if (count == 0) {
      return jsonEncode({
        'detected': false,
        'biomarker': '',
        'value': 0.0,
        'confidence': 0.0,
        'dominant_rgb': [0, 0, 0],
        'warning': 'ROI 영역에 유효한 픽셀이 없습니다.',
      });
    }

    final avgR = rSum ~/ count;
    final avgG = gSum ~/ count;
    final avgB = bSum ~/ count;

    // 간단한 색상 기반 카트리지 반응 판별 (Rust 로직과 동일)
    bool detected = avgR > 30 || avgG > 30 || avgB > 30;
    String biomarker = '';
    double value = 0.0;
    double confidence = 0.0;

    if (detected) {
      if (avgR > 150 && avgG < 100 && avgB < 100) {
        biomarker = 'pH';
        value = 4.0 + (avgR - 150) / 105.0 * 3.0;
        confidence = 0.85;
      } else if (avgR < 100 && avgG > 150 && avgB < 100) {
        biomarker = 'Glucose';
        value = 70.0 + (avgG - 150) / 105.0 * 130.0;
        confidence = 0.82;
      } else if (avgR < 100 && avgG < 100 && avgB > 150) {
        biomarker = 'Protein';
        value = (avgB - 150) / 105.0 * 30.0;
        confidence = 0.78;
      } else {
        biomarker = 'General';
        value = (avgR + avgG + avgB) / (3.0 * 255.0) * 100.0;
        confidence = 0.6;
      }
    }

    final warning = !detected
        ? '카트리지가 감지되지 않았습니다. 카메라를 카트리지 위에 정확히 위치시켜 주세요.'
        : confidence < 0.5
            ? '측정 신뢰도가 낮습니다. 조명을 확인해 주세요.'
            : '';

    return jsonEncode({
      'detected': detected,
      'biomarker': biomarker,
      'value': value,
      'confidence': confidence,
      'dominant_rgb': [avgR, avgG, avgB],
      'warning': warning,
    });
  }
}

// ─────────────────────────────────────────────
// 3. Riverpod 프로바이더
// ─────────────────────────────────────────────

/// 비전 분석기 상태 프로바이더
final visionAnalyzerProvider =
    StateNotifierProvider.autoDispose<VisionAnalyzerNotifier, VisionAnalyzerState>(
  (ref) => VisionAnalyzerNotifier(),
);

/// 비전 분석기 상태 노티파이어
class VisionAnalyzerNotifier extends StateNotifier<VisionAnalyzerState> {
  VisionAnalyzerNotifier() : super(const VisionAnalyzerState());

  CameraController? _cameraController;
  Timer? _throttleTimer;
  DateTime _lastFrameTime = DateTime.now();
  int _frameCount = 0;
  final _bridge = VisionRustBridge.instance;
  bool _isProcessing = false;

  /// 5fps 쓰로틀 간격 (200ms)
  static const _throttleInterval = Duration(milliseconds: 200);

  /// 카메라 초기화 및 스트림 시작
  Future<void> startCamera() async {
    // 1) 카메라 권한 요청
    state = state.copyWith(status: VisionAnalyzerStatus.requestingPermission);

    final permissionStatus = await Permission.camera.request();

    if (permissionStatus.isDenied) {
      state = state.copyWith(
        status: VisionAnalyzerStatus.permissionDenied,
        errorMessage: '카메라 권한이 거부되었습니다. 카트리지 분석을 위해 카메라 권한이 필요합니다.',
      );
      return;
    }

    if (permissionStatus.isPermanentlyDenied) {
      state = state.copyWith(
        status: VisionAnalyzerStatus.permissionPermanentlyDenied,
        errorMessage: '카메라 권한이 영구적으로 거부되었습니다. 설정에서 권한을 허용해주세요.',
      );
      return;
    }

    // 2) 카메라 초기화
    state = state.copyWith(status: VisionAnalyzerStatus.initializing);

    try {
      final cameras = await availableCameras();
      if (cameras.isEmpty) {
        state = state.copyWith(
          status: VisionAnalyzerStatus.error,
          errorMessage: '사용 가능한 카메라가 없습니다.',
        );
        return;
      }

      // 후면 카메라 우선 선택
      final camera = cameras.firstWhere(
        (c) => c.lensDirection == CameraLensDirection.back,
        orElse: () => cameras.first,
      );

      _cameraController = CameraController(
        camera,
        ResolutionPreset.medium,
        enableAudio: false,
        imageFormatGroup: ImageFormatGroup.bgra8888,
      );

      await _cameraController!.initialize();

      // Rust 브릿지 초기화
      await _bridge.init();

      // 3) 이미지 스트림 시작 (5fps 쓰로틀)
      _frameCount = 0;
      _lastFrameTime = DateTime.now();

      await _cameraController!.startImageStream(_onCameraFrame);

      state = state.copyWith(status: VisionAnalyzerStatus.streaming);
    } catch (e) {
      state = state.copyWith(
        status: VisionAnalyzerStatus.error,
        errorMessage: '카메라 초기화 실패: $e',
      );
    }
  }

  /// 카메라 프레임 콜백 — 5fps 쓰로틀 적용
  void _onCameraFrame(CameraImage image) {
    final now = DateTime.now();
    final elapsed = now.difference(_lastFrameTime);

    // 5fps 쓰로틀: 200ms 미만이면 프레임 무시
    if (elapsed < _throttleInterval) return;
    if (_isProcessing) return;

    _lastFrameTime = now;
    _isProcessing = true;

    _processFrame(image).then((_) {
      _isProcessing = false;
    }).catchError((_) {
      _isProcessing = false;
    });
  }

  /// 프레임 처리 — Rust FFI 호출
  Future<void> _processFrame(CameraImage image) async {
    try {
      // CameraImage → Uint8List (BGRA8888)
      final bytes = _cameraImageToBytes(image);
      final width = image.width;
      final height = image.height;

      // Rust FFI 호출 (또는 시뮬레이션)
      final jsonResult = await _bridge.processCameraFrame(bytes, width, height);

      // JSON 파싱
      final Map<String, dynamic> decoded = jsonDecode(jsonResult);
      final result = VisionAnalysisResult.fromJson(decoded);

      _frameCount++;

      // FPS 계산
      final fps = _frameCount > 1
          ? 1000.0 / _throttleInterval.inMilliseconds
          : 0.0;

      state = state.copyWith(
        latestResult: result,
        framesProcessed: _frameCount,
        fps: fps,
      );
    } catch (e) {
      // 개별 프레임 오류는 스트림을 중단하지 않음
      state = state.copyWith(
        errorMessage: '프레임 처리 오류: $e',
      );
    }
  }

  /// CameraImage → BGRA Uint8List 변환
  Uint8List _cameraImageToBytes(CameraImage image) {
    if (image.planes.isEmpty) return Uint8List(0);

    // 단일 plane (BGRA8888)
    if (image.planes.length == 1) {
      return image.planes[0].bytes;
    }

    // 다중 plane (YUV 등) → 결합
    int totalBytes = 0;
    for (final plane in image.planes) {
      totalBytes += plane.bytes.length;
    }

    final result = Uint8List(totalBytes);
    int offset = 0;
    for (final plane in image.planes) {
      result.setRange(offset, offset + plane.bytes.length, plane.bytes);
      offset += plane.bytes.length;
    }
    return result;
  }

  /// 분석 확정 — 현재 결과를 최종 결과로 저장
  void confirmResult() {
    if (state.latestResult != null && state.latestResult!.detected) {
      stopCamera();
      state = state.copyWith(status: VisionAnalyzerStatus.completed);
    }
  }

  /// 카메라 중지 및 리소스 해제
  void stopCamera() {
    _throttleTimer?.cancel();
    _throttleTimer = null;

    if (_cameraController != null) {
      if (_cameraController!.value.isStreamingImages) {
        _cameraController!.stopImageStream();
      }
      _cameraController!.dispose();
      _cameraController = null;
    }
  }

  /// 시스템 설정 열기 (영구 거부 시)
  Future<void> openPermissionSettings() async {
    await openAppSettings();
  }

  /// 카메라 컨트롤러 접근 (프리뷰 표시용)
  CameraController? get cameraController => _cameraController;

  @override
  void dispose() {
    stopCamera();
    super.dispose();
  }
}

// ─────────────────────────────────────────────
// 4. 비전 분석기 화면 위젯
// ─────────────────────────────────────────────

/// 카트리지 비전 분석 화면
///
/// 카메라 프리뷰 + 실시간 분석 결과 오버레이.
/// Riverpod 상태 관리, 카메라 권한 거부 처리 포함.
class VisionAnalyzerScreen extends ConsumerStatefulWidget {
  const VisionAnalyzerScreen({super.key});

  @override
  ConsumerState<VisionAnalyzerScreen> createState() =>
      _VisionAnalyzerScreenState();
}

class _VisionAnalyzerScreenState extends ConsumerState<VisionAnalyzerScreen>
    with WidgetsBindingObserver {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addObserver(this);
    // 화면 진입 시 자동 카메라 시작
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(visionAnalyzerProvider.notifier).startCamera();
    });
  }

  @override
  void dispose() {
    WidgetsBinding.instance.removeObserver(this);
    super.dispose();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    final notifier = ref.read(visionAnalyzerProvider.notifier);
    if (state == AppLifecycleState.inactive ||
        state == AppLifecycleState.paused) {
      notifier.stopCamera();
    } else if (state == AppLifecycleState.resumed) {
      final currentState = ref.read(visionAnalyzerProvider);
      if (currentState.status == VisionAnalyzerStatus.streaming) {
        notifier.startCamera();
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(visionAnalyzerProvider);
    final notifier = ref.read(visionAnalyzerProvider.notifier);
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('카트리지 비전 분석'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () {
            notifier.stopCamera();
            context.pop();
          },
        ),
        actions: [
          if (state.status == VisionAnalyzerStatus.streaming)
            Padding(
              padding: const EdgeInsets.only(right: 12),
              child: Center(
                child: Text(
                  '${state.fps.toStringAsFixed(1)} fps',
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: AppTheme.waveCyan,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
            ),
        ],
      ),
      body: _buildBody(context, state, notifier, theme),
    );
  }

  Widget _buildBody(
    BuildContext context,
    VisionAnalyzerState state,
    VisionAnalyzerNotifier notifier,
    ThemeData theme,
  ) {
    switch (state.status) {
      case VisionAnalyzerStatus.idle:
      case VisionAnalyzerStatus.requestingPermission:
      case VisionAnalyzerStatus.initializing:
        return _buildLoadingView(theme, state.status);

      case VisionAnalyzerStatus.permissionDenied:
        return _buildPermissionDeniedView(theme, notifier, isPermanent: false);

      case VisionAnalyzerStatus.permissionPermanentlyDenied:
        return _buildPermissionDeniedView(theme, notifier, isPermanent: true);

      case VisionAnalyzerStatus.streaming:
        return _buildStreamingView(state, notifier, theme);

      case VisionAnalyzerStatus.completed:
        return _buildCompletedView(state, theme);

      case VisionAnalyzerStatus.error:
        return _buildErrorView(state, notifier, theme);
    }
  }

  /// 로딩 뷰 (초기화 / 권한 요청 중)
  Widget _buildLoadingView(ThemeData theme, VisionAnalyzerStatus status) {
    final message = switch (status) {
      VisionAnalyzerStatus.requestingPermission => '카메라 권한을 요청하는 중...',
      VisionAnalyzerStatus.initializing => '카메라를 초기화하는 중...',
      _ => '준비 중...',
    };

    return Center(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          const CircularProgressIndicator(color: AppTheme.sanggamGold),
          const SizedBox(height: 24),
          Text(message, style: theme.textTheme.bodyLarge),
        ],
      ),
    );
  }

  /// 카메라 권한 거부 뷰
  Widget _buildPermissionDeniedView(
    ThemeData theme,
    VisionAnalyzerNotifier notifier, {
    required bool isPermanent,
  }) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              Icons.camera_alt_outlined,
              size: 72,
              color: Colors.grey.withValues(alpha: 0.5),
            ),
            const SizedBox(height: 24),
            Text(
              '카메라 권한 필요',
              style: theme.textTheme.titleLarge?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 12),
            Text(
              isPermanent
                  ? '카메라 권한이 영구적으로 거부되었습니다.\n설정에서 직접 권한을 허용해주세요.'
                  : '카트리지 분석을 위해 카메라 접근 권한이 필요합니다.\n권한을 허용해주세요.',
              textAlign: TextAlign.center,
              style: theme.textTheme.bodyMedium?.copyWith(
                color: theme.textTheme.bodySmall?.color,
              ),
            ),
            const SizedBox(height: 32),
            if (isPermanent)
              FilledButton.icon(
                onPressed: () => notifier.openPermissionSettings(),
                icon: const Icon(Icons.settings),
                label: const Text('설정 열기'),
                style: FilledButton.styleFrom(
                  backgroundColor: AppTheme.sanggamGold,
                  minimumSize: const Size(200, 48),
                ),
              )
            else
              FilledButton.icon(
                onPressed: () => notifier.startCamera(),
                icon: const Icon(Icons.camera_alt),
                label: const Text('권한 다시 요청'),
                style: FilledButton.styleFrom(
                  backgroundColor: AppTheme.sanggamGold,
                  minimumSize: const Size(200, 48),
                ),
              ),
            const SizedBox(height: 12),
            TextButton(
              onPressed: () => context.pop(),
              child: const Text('돌아가기'),
            ),
          ],
        ),
      ),
    );
  }

  /// 스트리밍 뷰 (카메라 프리뷰 + 분석 결과 오버레이)
  Widget _buildStreamingView(
    VisionAnalyzerState state,
    VisionAnalyzerNotifier notifier,
    ThemeData theme,
  ) {
    final controller = notifier.cameraController;
    if (controller == null || !controller.value.isInitialized) {
      return _buildLoadingView(theme, VisionAnalyzerStatus.initializing);
    }

    return Stack(
      fit: StackFit.expand,
      children: [
        // 카메라 프리뷰
        ClipRRect(
          borderRadius: BorderRadius.zero,
          child: CameraPreview(controller),
        ),

        // ROI 가이드 오버레이
        CustomPaint(
          painter: _RoiGuidePainter(
            detected: state.latestResult?.detected ?? false,
          ),
        ),

        // 상단 상태 표시
        Positioned(
          top: 16,
          left: 16,
          right: 16,
          child: _buildStatusBar(state, theme),
        ),

        // 하단 분석 결과 카드
        if (state.latestResult != null)
          Positioned(
            bottom: 0,
            left: 0,
            right: 0,
            child: _buildResultCard(state.latestResult!, notifier, theme),
          ),
      ],
    );
  }

  /// 상단 상태 바
  Widget _buildStatusBar(VisionAnalyzerState state, ThemeData theme) {
    final result = state.latestResult;
    final detected = result?.detected ?? false;

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
      decoration: BoxDecoration(
        color: Colors.black.withValues(alpha: 0.6),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        children: [
          Icon(
            detected ? Icons.check_circle : Icons.search,
            color: detected ? Colors.greenAccent : AppTheme.sanggamGold,
            size: 20,
          ),
          const SizedBox(width: 8),
          Expanded(
            child: Text(
              detected
                  ? '카트리지 감지됨 — ${result!.biomarker}'
                  : '카트리지를 카메라에 비춰주세요...',
              style: theme.textTheme.bodySmall?.copyWith(
                color: Colors.white,
                fontWeight: FontWeight.w500,
              ),
            ),
          ),
          Text(
            '프레임: ${state.framesProcessed}',
            style: theme.textTheme.bodySmall?.copyWith(
              color: Colors.white70,
              fontSize: 11,
            ),
          ),
        ],
      ),
    );
  }

  /// 하단 분석 결과 카드
  Widget _buildResultCard(
    VisionAnalysisResult result,
    VisionAnalyzerNotifier notifier,
    ThemeData theme,
  ) {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: theme.scaffoldBackgroundColor.withValues(alpha: 0.95),
        borderRadius: const BorderRadius.vertical(top: Radius.circular(20)),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.3),
            blurRadius: 20,
            offset: const Offset(0, -4),
          ),
        ],
      ),
      child: SafeArea(
        top: false,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // 감지 색상 미리보기 + 바이오마커
            Row(
              children: [
                Container(
                  width: 40,
                  height: 40,
                  decoration: BoxDecoration(
                    color: result.dominantColor,
                    borderRadius: BorderRadius.circular(8),
                    border: Border.all(color: Colors.white24),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        result.detected
                            ? result.biomarker.isEmpty
                                ? '분석 중...'
                                : result.biomarker
                            : '미감지',
                        style: theme.textTheme.titleMedium?.copyWith(
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      if (result.detected)
                        Text(
                          '값: ${result.value.toStringAsFixed(2)}  |  신뢰도: ${result.confidencePercent}',
                          style: theme.textTheme.bodySmall?.copyWith(
                            color: AppTheme.waveCyan,
                          ),
                        ),
                    ],
                  ),
                ),
                // 신뢰도 원형 표시
                if (result.detected)
                  SizedBox(
                    width: 48,
                    height: 48,
                    child: Stack(
                      fit: StackFit.expand,
                      children: [
                        CircularProgressIndicator(
                          value: result.confidence,
                          strokeWidth: 4,
                          backgroundColor: Colors.grey.withValues(alpha: 0.2),
                          valueColor: AlwaysStoppedAnimation<Color>(
                            result.confidence > 0.7
                                ? Colors.greenAccent
                                : result.confidence > 0.4
                                    ? AppTheme.sanggamGold
                                    : Colors.redAccent,
                          ),
                        ),
                        Center(
                          child: Text(
                            '${(result.confidence * 100).toInt()}',
                            style: theme.textTheme.bodySmall?.copyWith(
                              fontWeight: FontWeight.bold,
                              fontSize: 12,
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
              ],
            ),

            // 경고 메시지
            if (result.warning.isNotEmpty) ...[
              const SizedBox(height: 12),
              Container(
                padding: const EdgeInsets.all(10),
                decoration: BoxDecoration(
                  color: Colors.orange.withValues(alpha: 0.15),
                  borderRadius: BorderRadius.circular(8),
                  border: Border.all(
                    color: Colors.orange.withValues(alpha: 0.3),
                  ),
                ),
                child: Row(
                  children: [
                    const Icon(Icons.warning_amber, color: Colors.orange, size: 18),
                    const SizedBox(width: 8),
                    Expanded(
                      child: Text(
                        result.warning,
                        style: theme.textTheme.bodySmall?.copyWith(
                          color: Colors.orange,
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ],

            // 확정 버튼
            if (result.detected && result.confidence > 0.5) ...[
              const SizedBox(height: 16),
              SizedBox(
                width: double.infinity,
                child: FilledButton.icon(
                  onPressed: () => notifier.confirmResult(),
                  icon: const Icon(Icons.check),
                  label: const Text('분석 결과 확정'),
                  style: FilledButton.styleFrom(
                    backgroundColor: AppTheme.sanggamGold,
                    minimumSize: const Size.fromHeight(48),
                  ),
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }

  /// 완료 뷰
  Widget _buildCompletedView(VisionAnalyzerState state, ThemeData theme) {
    final result = state.latestResult!;

    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              Icons.check_circle_outline,
              size: 80,
              color: Colors.greenAccent.withValues(alpha: 0.8),
            ),
            const SizedBox(height: 24),
            Text(
              '분석 완료',
              style: theme.textTheme.headlineSmall?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 16),
            Card(
              child: Padding(
                padding: const EdgeInsets.all(20),
                child: Column(
                  children: [
                    _buildResultRow(theme, '바이오마커', result.biomarker),
                    const Divider(height: 20),
                    _buildResultRow(theme, '측정값', result.value.toStringAsFixed(2)),
                    const Divider(height: 20),
                    _buildResultRow(theme, '신뢰도', result.confidencePercent),
                    const Divider(height: 20),
                    _buildResultRow(theme, '처리 프레임', '${state.framesProcessed}'),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 24),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                OutlinedButton.icon(
                  onPressed: () {
                    ref.read(visionAnalyzerProvider.notifier).startCamera();
                  },
                  icon: const Icon(Icons.refresh),
                  label: const Text('다시 측정'),
                ),
                const SizedBox(width: 12),
                FilledButton.icon(
                  onPressed: () => context.pop(),
                  icon: const Icon(Icons.done),
                  label: const Text('완료'),
                  style: FilledButton.styleFrom(
                    backgroundColor: AppTheme.sanggamGold,
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildResultRow(ThemeData theme, String label, String value) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(label, style: theme.textTheme.bodyMedium?.copyWith(
          color: theme.textTheme.bodySmall?.color,
        )),
        Text(value, style: theme.textTheme.bodyMedium?.copyWith(
          fontWeight: FontWeight.bold,
        )),
      ],
    );
  }

  /// 오류 뷰
  Widget _buildErrorView(
    VisionAnalyzerState state,
    VisionAnalyzerNotifier notifier,
    ThemeData theme,
  ) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(Icons.error_outline, size: 64, color: Colors.redAccent),
            const SizedBox(height: 24),
            Text(
              '오류 발생',
              style: theme.textTheme.titleLarge?.copyWith(
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 12),
            Text(
              state.errorMessage ?? '알 수 없는 오류가 발생했습니다.',
              textAlign: TextAlign.center,
              style: theme.textTheme.bodyMedium,
            ),
            const SizedBox(height: 24),
            FilledButton.icon(
              onPressed: () => notifier.startCamera(),
              icon: const Icon(Icons.refresh),
              label: const Text('다시 시도'),
              style: FilledButton.styleFrom(
                backgroundColor: AppTheme.sanggamGold,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

// ─────────────────────────────────────────────
// 5. ROI 가이드 오버레이 페인터
// ─────────────────────────────────────────────

/// 카메라 프리뷰 위에 ROI(관심 영역) 가이드를 그리는 페인터
class _RoiGuidePainter extends CustomPainter {
  final bool detected;

  _RoiGuidePainter({required this.detected});

  @override
  void paint(Canvas canvas, Size size) {
    // 중앙 40% ROI 영역
    final roiRect = Rect.fromLTWH(
      size.width * 0.3,
      size.height * 0.3,
      size.width * 0.4,
      size.height * 0.4,
    );

    // 반투명 어둡게 (ROI 외부)
    final dimPaint = Paint()
      ..color = Colors.black.withValues(alpha: 0.5)
      ..style = PaintingStyle.fill;

    // ROI 외부 영역 어둡게
    final fullRect = Rect.fromLTWH(0, 0, size.width, size.height);
    canvas.drawPath(
      Path.combine(
        PathOperation.difference,
        Path()..addRect(fullRect),
        Path()..addRRect(RRect.fromRectAndRadius(roiRect, const Radius.circular(12))),
      ),
      dimPaint,
    );

    // ROI 테두리
    final borderPaint = Paint()
      ..color = detected ? Colors.greenAccent : AppTheme.sanggamGold
      ..style = PaintingStyle.stroke
      ..strokeWidth = 2.5;

    canvas.drawRRect(
      RRect.fromRectAndRadius(roiRect, const Radius.circular(12)),
      borderPaint,
    );

    // 코너 액센트 (L자 형태)
    final cornerPaint = Paint()
      ..color = detected ? Colors.greenAccent : AppTheme.sanggamGold
      ..style = PaintingStyle.stroke
      ..strokeWidth = 4
      ..strokeCap = StrokeCap.round;

    const cornerLen = 24.0;
    final corners = [
      // 좌상단
      [Offset(roiRect.left, roiRect.top + cornerLen), Offset(roiRect.left, roiRect.top), Offset(roiRect.left + cornerLen, roiRect.top)],
      // 우상단
      [Offset(roiRect.right - cornerLen, roiRect.top), Offset(roiRect.right, roiRect.top), Offset(roiRect.right, roiRect.top + cornerLen)],
      // 좌하단
      [Offset(roiRect.left, roiRect.bottom - cornerLen), Offset(roiRect.left, roiRect.bottom), Offset(roiRect.left + cornerLen, roiRect.bottom)],
      // 우하단
      [Offset(roiRect.right - cornerLen, roiRect.bottom), Offset(roiRect.right, roiRect.bottom), Offset(roiRect.right, roiRect.bottom - cornerLen)],
    ];

    for (final corner in corners) {
      canvas.drawLine(corner[0], corner[1], cornerPaint);
      canvas.drawLine(corner[1], corner[2], cornerPaint);
    }
  }

  @override
  bool shouldRepaint(_RoiGuidePainter oldDelegate) =>
      oldDelegate.detected != detected;
}
