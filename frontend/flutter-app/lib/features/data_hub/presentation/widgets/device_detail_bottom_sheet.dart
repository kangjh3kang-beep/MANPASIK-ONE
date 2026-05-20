import 'dart:ui';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:manpasik/core/theme/sanggam_theme.dart';
import 'package:manpasik/features/devices/domain/device_repository.dart';

// ───────────────────────────────────────────────────
// DeviceDetailBottomSheet — Sanggam Orbit 기기 상세
//
// [Rule 4] app_theme → sanggam_theme
// [Rule 4] Theme.of(context) isDark 분기 제거 (항상 다크)
// [Rule 4] isDark 파라미터 전체 제거 (15+ 위젯)
// [Rule 4] AppTheme.sanggamGold → SanggamTheme.primary
// [Rule 4] AppTheme.waveCyan → SanggamTheme.jagaeCyan
// [Rule 4] Color(0xFF00E676) → SanggamTheme.jagaeCyan
// [Rule 4] Color(0xFFFF4D4D) → SanggamTheme.error
// [Rule 4] Color(0xFF1B2640/050B14) → SanggamTheme 상수
// [Rule 4] Colors.redAccent → SanggamTheme.error
// [Rule 4] Colors.grey → SanggamTheme.onSurfaceDim
// ───────────────────────────────────────────────────

/// 모니터링 화면 내 기기 상세 Dialog
void showDeviceDetailSheet(
    BuildContext context,
    ConnectedDevice device) {
  showDialog(
    context: context,
    barrierColor:
        Colors.black.withValues(alpha: 0.7),
    builder: (context) {
      return Center(
        child: Material(
          color: Colors.transparent,
          child: Container(
            width: MediaQuery.of(context)
                    .size
                    .width *
                0.9,
            constraints:
                const BoxConstraints(
                    maxWidth: 420,
                    maxHeight: 680),
            margin:
                const EdgeInsets.all(24),
            child: _DeviceDetailContent(
              device: device,
              scrollController:
                  ScrollController(),
            ),
          ),
        ),
      );
    },
  );
}

class _DeviceDetailContent
    extends StatelessWidget {
  final ConnectedDevice device;
  final ScrollController scrollController;

  const _DeviceDetailContent({
    required this.device,
    required this.scrollController,
  });

  @override
  Widget build(BuildContext context) {
    final isConnected = device.status ==
        DeviceConnectionStatus.connected;
    final statusColor = isConnected
        ? SanggamTheme.jagaeCyan
        : SanggamTheme.onSurfaceDim;

    return ClipRRect(
      borderRadius:
          BorderRadius.circular(24),
      child: BackdropFilter(
        filter: ImageFilter.blur(
            sigmaX: 20, sigmaY: 20),
        child: Container(
          decoration: BoxDecoration(
            borderRadius:
                BorderRadius.circular(24),
            gradient: LinearGradient(
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
              colors: [
                SanggamTheme.surfaceVariant
                    .withValues(alpha: 0.85),
                SanggamTheme.background
                    .withValues(alpha: 0.95),
              ],
            ),
            border: Border.all(
              color: SanggamTheme.primary
                  .withValues(alpha: 0.3),
              width: 1.5,
            ),
            boxShadow: [
              BoxShadow(
                color: Colors.black
                    .withValues(alpha: 0.5),
                blurRadius: 30,
                spreadRadius: 5,
              ),
            ],
          ),
          child: Column(
            children: [
              // 닫기 버튼 헤더
              Padding(
                padding:
                    const EdgeInsets.fromLTRB(
                        20, 16, 16, 0),
                child: Row(
                  mainAxisAlignment:
                      MainAxisAlignment.end,
                  children: [
                    IconButton(
                      icon: const Icon(
                          Icons.close,
                          color: SanggamTheme
                              .onSurfaceDim),
                      tooltip: '닫기',
                      onPressed: () =>
                          Navigator.of(
                                  context)
                              .pop(),
                      padding:
                          EdgeInsets.zero,
                      constraints:
                          const BoxConstraints(),
                      splashRadius: 20,
                    ),
                  ],
                ),
              ),

              // 스크롤 가능한 컨텐츠
              Expanded(
                child: ListView(
                  controller:
                      scrollController,
                  padding:
                      const EdgeInsets
                          .fromLTRB(
                          24, 0, 24, 24),
                  children: [
                    _buildHeader(
                        isConnected,
                        statusColor),
                    const SizedBox(
                        height: 20),
                    _buildTypeSpecificSection(),
                    const SizedBox(
                        height: 20),
                    _buildSectionTitle(
                        '상태 정보'),
                    const SizedBox(
                        height: 8),
                    _StatusBar(
                      icon: device.batteryLevel <
                              20
                          ? Icons
                              .battery_alert
                          : Icons
                              .battery_full,
                      label: '배터리',
                      value: device
                          .batteryLevel,
                      color: device.batteryLevel <
                              20
                          ? SanggamTheme
                              .error
                          : SanggamTheme
                              .jagaeCyan,
                    ),
                    const SizedBox(
                        height: 8),
                    _StatusBar(
                      icon: Icons
                          .signal_cellular_alt,
                      label: '신호 강도',
                      value: device
                          .signalStrength,
                      color: device.signalStrength >
                              60
                          ? SanggamTheme
                              .jagaeCyan
                          : SanggamTheme
                              .primary,
                    ),
                    const SizedBox(
                        height: 20),
                    if (device.latestReadings
                            .length >=
                        2) ...[
                      _buildSectionTitle(
                          '측정 추이'),
                      const SizedBox(
                          height: 8),
                      Container(
                        height: 100,
                        width:
                            double.infinity,
                        padding:
                            const EdgeInsets
                                .all(12),
                        decoration:
                            BoxDecoration(
                          color: Colors.white
                              .withValues(
                                  alpha:
                                      0.04),
                          borderRadius:
                              BorderRadius
                                  .circular(
                                      14),
                          border:
                              Border.all(
                            color: Colors
                                .white
                                .withValues(
                                    alpha:
                                        0.06),
                          ),
                        ),
                        child: CustomPaint(
                          painter:
                              _DetailSparklinePainter(
                            data: device
                                .latestReadings,
                            color: SanggamTheme
                                .jagaeCyan,
                          ),
                        ),
                      ),
                      const SizedBox(
                          height: 20),
                    ],
                    _buildSectionTitle(
                        '기기 정보'),
                    const SizedBox(
                        height: 8),
                    Container(
                      padding:
                          const EdgeInsets
                              .all(12),
                      decoration:
                          BoxDecoration(
                        color: Colors.white
                            .withValues(
                                alpha:
                                    0.04),
                        borderRadius:
                            BorderRadius
                                .circular(
                                    12),
                      ),
                      child: Column(
                        children: [
                          _InfoRow(
                              label: 'ID',
                              value:
                                  device
                                      .id),
                          _InfoRow(
                              label: '타입',
                              value: _deviceTypeName(
                                  device
                                      .type)),
                          _InfoRow(
                              label: '상태',
                              value: isConnected
                                  ? '연결됨'
                                  : '연결 끊김'),
                        ],
                      ),
                    ),
                    const SizedBox(
                        height: 16),
                    SizedBox(
                      width:
                          double.infinity,
                      height: 44,
                      child: FilledButton
                          .icon(
                        onPressed: () {
                          Navigator.pop(
                              context);
                          context.push(
                              '/devices/${device.id}');
                        },
                        icon: const Icon(
                            Icons.settings,
                            size: 18),
                        label: const Text(
                            '기기 관리'),
                        style: FilledButton
                            .styleFrom(
                          backgroundColor:
                              SanggamTheme
                                  .primary,
                          foregroundColor:
                              SanggamTheme
                                  .background,
                          shape:
                              RoundedRectangleBorder(
                            borderRadius:
                                BorderRadius
                                    .circular(
                                        12),
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildHeader(
      bool isConnected, Color statusColor) {
    return Row(
      children: [
        Container(
          width: 48,
          height: 48,
          decoration: BoxDecoration(
            shape: BoxShape.circle,
            color: statusColor.withValues(
                alpha: 0.15),
            border: Border.all(
                color: statusColor.withValues(
                    alpha: 0.4)),
          ),
          child: Icon(
            _deviceIcon(device.type),
            color: SanggamTheme.primary,
            size: 24,
          ),
        ),
        const SizedBox(width: 12),
        Expanded(
          child: Column(
            crossAxisAlignment:
                CrossAxisAlignment.start,
            children: [
              Text(
                device.name,
                style: const TextStyle(
                  color: Colors.white,
                  fontWeight: FontWeight.bold,
                  fontSize: 18,
                ),
              ),
              const SizedBox(height: 2),
              Row(
                children: [
                  Container(
                    padding:
                        const EdgeInsets
                            .symmetric(
                            horizontal: 6,
                            vertical: 2),
                    decoration:
                        BoxDecoration(
                      color: statusColor
                          .withValues(
                              alpha: 0.12),
                      borderRadius:
                          BorderRadius
                              .circular(6),
                    ),
                    child: Row(
                      mainAxisSize:
                          MainAxisSize.min,
                      children: [
                        Container(
                          width: 6,
                          height: 6,
                          decoration:
                              BoxDecoration(
                            color:
                                statusColor,
                            shape: BoxShape
                                .circle,
                          ),
                        ),
                        const SizedBox(
                            width: 4),
                        Text(
                          isConnected
                              ? '활성'
                              : '비활성',
                          style: TextStyle(
                            color:
                                statusColor,
                            fontSize: 10,
                            fontWeight:
                                FontWeight
                                    .bold,
                          ),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(width: 8),
                  Text(
                    _deviceTypeName(
                        device.type),
                    style: const TextStyle(
                      color: SanggamTheme
                          .onSurfaceDim,
                      fontSize: 12,
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildSectionTitle(String title) {
    return Text(
      title,
      style: TextStyle(
        color: SanggamTheme.onSurfaceDim
            .withValues(alpha: 0.8),
        fontSize: 10,
        fontWeight: FontWeight.w600,
        letterSpacing: 1.5,
      ),
    );
  }

  Widget _buildTypeSpecificSection() {
    switch (device.type) {
      case DeviceType.gasCartridge:
        return _GasDetailSection(
            values: device.currentValues);
      case DeviceType.envCartridge:
        return _EnvDetailSection(
            values: device.currentValues);
      case DeviceType.bioCartridge:
        return _BioDetailSection(
            values: device.currentValues);
      case DeviceType.unknown:
        return _GenericSensorSection(
            values: device.currentValues);
    }
  }

  IconData _deviceIcon(DeviceType type) {
    switch (type) {
      case DeviceType.gasCartridge:
        return Icons.cloud_outlined;
      case DeviceType.envCartridge:
        return Icons.thermostat_outlined;
      case DeviceType.bioCartridge:
        return Icons.science_outlined;
      case DeviceType.unknown:
        return Icons.device_unknown;
    }
  }

  String _deviceTypeName(DeviceType type) {
    switch (type) {
      case DeviceType.gasCartridge:
        return '가스 카트리지';
      case DeviceType.envCartridge:
        return '환경 카트리지';
      case DeviceType.bioCartridge:
        return '바이오 카트리지';
      case DeviceType.unknown:
        return '알 수 없음';
    }
  }
}

// ═══════════════════════════════════════
// 가스 카트리지 전용 영역
// ═══════════════════════════════════════

enum _DangerLevel { safe, caution, danger }

class _GasDetailSection
    extends StatelessWidget {
  final Map<String, dynamic> values;

  const _GasDetailSection(
      {required this.values});

  _DangerLevel _assessGas(
      String key, dynamic val) {
    final str =
        '$val'.replaceAll(RegExp(r'[^0-9.]'), '');
    final num = double.tryParse(str);
    if (num == null) {
      if (key == 'Smoke') {
        return val == 'None'
            ? _DangerLevel.safe
            : _DangerLevel.danger;
      }
      return _DangerLevel.safe;
    }
    switch (key) {
      case 'CO':
        if (num > 70) {
          return _DangerLevel.danger;
        }
        if (num > 30) {
          return _DangerLevel.caution;
        }
        return _DangerLevel.safe;
      case 'CO2':
        if (num > 2000) {
          return _DangerLevel.danger;
        }
        if (num > 1000) {
          return _DangerLevel.caution;
        }
        return _DangerLevel.safe;
      case 'LNG':
        if (num > 2) {
          return _DangerLevel.danger;
        }
        if (num > 0.5) {
          return _DangerLevel.caution;
        }
        return _DangerLevel.safe;
      case 'VOC':
        if (num > 0.5) {
          return _DangerLevel.danger;
        }
        if (num > 0.1) {
          return _DangerLevel.caution;
        }
        return _DangerLevel.safe;
      default:
        return _DangerLevel.safe;
    }
  }

  Color _dangerColor(_DangerLevel level) {
    switch (level) {
      case _DangerLevel.safe:
        return SanggamTheme.jagaeCyan;
      case _DangerLevel.caution:
        return SanggamTheme.primary;
      case _DangerLevel.danger:
        return SanggamTheme.error;
    }
  }

  String _dangerLabel(_DangerLevel level) {
    switch (level) {
      case _DangerLevel.safe:
        return '안전';
      case _DangerLevel.caution:
        return '주의';
      case _DangerLevel.danger:
        return '위험';
    }
  }

  @override
  Widget build(BuildContext context) {
    final entries = values.entries.toList();
    final levels = {
      for (final e in entries)
        e.key: _assessGas(e.key, e.value)
    };
    final worstLevel =
        levels.values.fold<_DangerLevel>(
      _DangerLevel.safe,
      (worst, level) =>
          level.index > worst.index
              ? level
              : worst,
    );
    final needVent =
        worstLevel != _DangerLevel.safe;

    return Column(
      crossAxisAlignment:
          CrossAxisAlignment.start,
      children: [
        Text(
          '센서 데이터',
          style: TextStyle(
            color: SanggamTheme.onSurfaceDim
                .withValues(alpha: 0.8),
            fontSize: 10,
            fontWeight: FontWeight.w600,
            letterSpacing: 1.5,
          ),
        ),
        const SizedBox(height: 8),
        Wrap(
          spacing: 8,
          runSpacing: 8,
          children: entries.map((e) {
            final level = levels[e.key]!;
            return _GasSensorTile(
              label: e.key,
              value: '${e.value}',
              level: level,
              dangerColor:
                  _dangerColor(level),
            );
          }).toList(),
        ),
        const SizedBox(height: 12),
        Container(
          padding:
              const EdgeInsets.all(10),
          decoration: BoxDecoration(
            color:
                _dangerColor(worstLevel)
                    .withValues(
                        alpha: 0.08),
            borderRadius:
                BorderRadius.circular(10),
            border: Border.all(
                color:
                    _dangerColor(worstLevel)
                        .withValues(
                            alpha: 0.25)),
          ),
          child: Row(
            children: [
              Icon(
                worstLevel ==
                        _DangerLevel.safe
                    ? Icons
                        .check_circle_outline
                    : (worstLevel ==
                            _DangerLevel
                                .caution
                        ? Icons
                            .warning_amber_rounded
                        : Icons
                            .error_outline),
                color: _dangerColor(
                    worstLevel),
                size: 18,
              ),
              const SizedBox(width: 8),
              Expanded(
                child: Column(
                  crossAxisAlignment:
                      CrossAxisAlignment
                          .start,
                  children: [
                    Text(
                      '종합 위험도: ${_dangerLabel(worstLevel)}',
                      style: TextStyle(
                        color: _dangerColor(
                            worstLevel),
                        fontSize: 12,
                        fontWeight:
                            FontWeight.bold,
                      ),
                    ),
                    const SizedBox(
                        height: 2),
                    Text(
                      needVent
                          ? '환기 권장 — 가스 수치가 정상 범위를 초과했습니다'
                          : '현재 정상 — 환기 불필요',
                      style:
                          const TextStyle(
                        color: SanggamTheme
                            .onSurfaceDim,
                        fontSize: 10,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }
}

class _GasSensorTile
    extends StatelessWidget {
  final String label;
  final String value;
  final _DangerLevel level;
  final Color dangerColor;

  const _GasSensorTile({
    required this.label,
    required this.value,
    required this.level,
    required this.dangerColor,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 90,
      padding: const EdgeInsets.symmetric(
          horizontal: 10, vertical: 8),
      decoration: BoxDecoration(
        color: dangerColor.withValues(
            alpha: 0.06),
        borderRadius:
            BorderRadius.circular(10),
        border: Border.all(
            color: dangerColor.withValues(
                alpha: 0.18)),
      ),
      child: Column(
        crossAxisAlignment:
            CrossAxisAlignment.start,
        children: [
          Text(label,
              style: const TextStyle(
                  color: SanggamTheme
                      .onSurfaceDim,
                  fontSize: 10)),
          const SizedBox(height: 2),
          Text(
            value,
            style: TextStyle(
              color: dangerColor,
              fontSize: 16,
              fontWeight: FontWeight.w700,
              fontFeatures: const [
                FontFeature
                    .tabularFigures()
              ],
            ),
          ),
          const SizedBox(height: 2),
          Container(
            width: 8,
            height: 8,
            decoration: BoxDecoration(
              color: dangerColor,
              shape: BoxShape.circle,
            ),
          ),
        ],
      ),
    );
  }
}

// ═══════════════════════════════════════
// 환경 카트리지 전용 영역
// ═══════════════════════════════════════

class _EnvDetailSection
    extends StatelessWidget {
  final Map<String, dynamic> values;

  const _EnvDetailSection(
      {required this.values});

  ({String label, Color color})
      _assessEnvParam(
          String key, dynamic val) {
    final str =
        '$val'.replaceAll(RegExp(r'[^0-9.]'), '');
    final num = double.tryParse(str);
    if (num == null) {
      return (
        label: '—',
        color: SanggamTheme.onSurfaceDim
      );
    }

    switch (key) {
      case 'Temp':
        if (num < 18) {
          return (
            label: '추움',
            color: SanggamTheme.jagaeCyan
          );
        }
        if (num > 26) {
          return (
            label: '더움',
            color: SanggamTheme.error
          );
        }
        return (
          label: '적정',
          color: SanggamTheme.jagaeCyan
        );
      case 'Humidity':
        if (num < 40) {
          return (
            label: '건조',
            color: SanggamTheme.primary
          );
        }
        if (num > 60) {
          return (
            label: '습함',
            color: SanggamTheme.jagaeCyan
          );
        }
        return (
          label: '적정',
          color: SanggamTheme.jagaeCyan
        );
      case 'Light':
        if (num < 100) {
          return (
            label: '어두움',
            color: SanggamTheme.onSurfaceDim
          );
        }
        if (num > 500) {
          return (
            label: '밝음',
            color: SanggamTheme.primary
          );
        }
        return (
          label: '적정',
          color: SanggamTheme.jagaeCyan
        );
      default:
        return (
          label: '—',
          color: SanggamTheme.onSurfaceDim
        );
    }
  }

  double _comfortScore() {
    double score = 0;
    int count = 0;

    for (final e in values.entries) {
      final str = '${e.value}'
          .replaceAll(RegExp(r'[^0-9.]'), '');
      final num = double.tryParse(str);
      if (num == null) continue;

      switch (e.key) {
        case 'Temp':
          score += 1.0 -
              ((num - 22).abs() / 8)
                  .clamp(0.0, 1.0);
          count++;
          break;
        case 'Humidity':
          score += 1.0 -
              ((num - 50).abs() / 20)
                  .clamp(0.0, 1.0);
          count++;
          break;
        case 'Light':
          score += 1.0 -
              ((num - 300).abs() / 300)
                  .clamp(0.0, 1.0);
          count++;
          break;
      }
    }
    return count > 0
        ? score / count
        : 0.5;
  }

  @override
  Widget build(BuildContext context) {
    final entries = values.entries.toList();
    final comfort = _comfortScore();
    final comfortLabel = comfort > 0.7
        ? '쾌적'
        : (comfort > 0.4 ? '보통' : '불쾌');
    final comfortColor = comfort > 0.7
        ? SanggamTheme.jagaeCyan
        : (comfort > 0.4
            ? SanggamTheme.primary
            : SanggamTheme.error);

    return Column(
      crossAxisAlignment:
          CrossAxisAlignment.start,
      children: [
        Text(
          '환경 상태',
          style: TextStyle(
            color: SanggamTheme.onSurfaceDim
                .withValues(alpha: 0.8),
            fontSize: 10,
            fontWeight: FontWeight.w600,
            letterSpacing: 1.5,
          ),
        ),
        const SizedBox(height: 8),
        Wrap(
          spacing: 8,
          runSpacing: 8,
          children: entries.map((e) {
            final assess = _assessEnvParam(
                e.key, e.value);
            return _EnvSensorTile(
              label: e.key,
              value: '${e.value}',
              statusLabel: assess.label,
              statusColor: assess.color,
            );
          }).toList(),
        ),
        const SizedBox(height: 14),
        Text(
          '쾌적도',
          style: TextStyle(
            color: SanggamTheme.onSurfaceDim
                .withValues(alpha: 0.8),
            fontSize: 10,
            fontWeight: FontWeight.w600,
            letterSpacing: 1.5,
          ),
        ),
        const SizedBox(height: 6),
        _ComfortIndicator(
          comfort: comfort,
          label: comfortLabel,
          color: comfortColor,
        ),
        const SizedBox(height: 10),
        ...entries
            .where((e) => [
                  'Temp',
                  'Humidity',
                  'Light'
                ].contains(e.key))
            .map((e) {
          final assess = _assessEnvParam(
              e.key, e.value);
          final range = e.key == 'Temp'
              ? '18~26°C'
              : (e.key == 'Humidity'
                  ? '40~60%'
                  : '100~500 lux');
          return Padding(
            padding:
                const EdgeInsets.only(
                    bottom: 4),
            child: Text(
              '${_paramKoName(e.key)} 범위: $range (${assess.label})',
              style: const TextStyle(
                color: SanggamTheme
                    .onSurfaceDim,
                fontSize: 10,
              ),
            ),
          );
        }),
      ],
    );
  }

  String _paramKoName(String key) {
    switch (key) {
      case 'Temp':
        return '온도';
      case 'Humidity':
        return '습도';
      case 'Light':
        return '조도';
      default:
        return key;
    }
  }
}

class _EnvSensorTile
    extends StatelessWidget {
  final String label;
  final String value;
  final String statusLabel;
  final Color statusColor;

  const _EnvSensorTile({
    required this.label,
    required this.value,
    required this.statusLabel,
    required this.statusColor,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 90,
      padding: const EdgeInsets.symmetric(
          horizontal: 10, vertical: 8),
      decoration: BoxDecoration(
        color: statusColor.withValues(
            alpha: 0.06),
        borderRadius:
            BorderRadius.circular(10),
        border: Border.all(
            color: statusColor.withValues(
                alpha: 0.18)),
      ),
      child: Column(
        crossAxisAlignment:
            CrossAxisAlignment.start,
        children: [
          Text(label,
              style: const TextStyle(
                  color: SanggamTheme
                      .onSurfaceDim,
                  fontSize: 10)),
          const SizedBox(height: 2),
          Text(
            value,
            style: const TextStyle(
              color: Colors.white,
              fontSize: 16,
              fontWeight: FontWeight.w700,
              fontFeatures: [
                FontFeature
                    .tabularFigures()
              ],
            ),
          ),
          const SizedBox(height: 2),
          Text(statusLabel,
              style: TextStyle(
                  color: statusColor,
                  fontSize: 9,
                  fontWeight:
                      FontWeight.w600)),
        ],
      ),
    );
  }
}

class _ComfortIndicator
    extends StatelessWidget {
  final double comfort;
  final String label;
  final Color color;

  const _ComfortIndicator({
    required this.comfort,
    required this.label,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        SizedBox(
          height: 24,
          child: Stack(
            children: [
              Positioned.fill(
                child: Container(
                  decoration: BoxDecoration(
                    borderRadius:
                        BorderRadius.circular(
                            12),
                    gradient:
                        const LinearGradient(
                            colors: [
                          Color(0xFF42A5F5),
                          Color(0xFF00E676),
                          Color(0xFFFFC107),
                          Color(0xFFFF4D4D),
                        ]),
                  ),
                ),
              ),
              Positioned(
                left: (comfort *
                        (MediaQuery.of(
                                    context)
                                .size
                                .width -
                            80))
                    .clamp(0.0,
                        double.infinity),
                top: 0,
                bottom: 0,
                child: Container(
                  width: 20,
                  height: 24,
                  decoration: BoxDecoration(
                    color: Colors.white,
                    borderRadius:
                        BorderRadius.circular(
                            10),
                    border: Border.all(
                        color: color,
                        width: 2),
                    boxShadow: [
                      BoxShadow(
                          color: color
                              .withValues(
                                  alpha:
                                      0.4),
                          blurRadius: 6)
                    ],
                  ),
                ),
              ),
            ],
          ),
        ),
        const SizedBox(height: 4),
        Row(
          mainAxisAlignment:
              MainAxisAlignment
                  .spaceBetween,
          children: [
            Text('불쾌',
                style: TextStyle(
                    fontSize: 8,
                    color: SanggamTheme
                        .onSurfaceDim
                        .withValues(
                            alpha: 0.5))),
            Text(label,
                style: TextStyle(
                    fontSize: 11,
                    color: color,
                    fontWeight:
                        FontWeight.bold)),
            Text('쾌적',
                style: TextStyle(
                    fontSize: 8,
                    color: SanggamTheme
                        .onSurfaceDim
                        .withValues(
                            alpha: 0.5))),
          ],
        ),
      ],
    );
  }
}

// ═══════════════════════════════════════
// 바이오 카트리지 전용 영역
// ═══════════════════════════════════════

class _BioDetailSection
    extends StatelessWidget {
  final Map<String, dynamic> values;

  const _BioDetailSection(
      {required this.values});

  ({String label, Color color}) _assessBio(
      String key, dynamic val) {
    final str =
        '$val'.replaceAll(RegExp(r'[^0-9.]'), '');
    final num = double.tryParse(str);

    switch (key) {
      case 'Pulse':
        if (num == null) {
          return (
            label: '—',
            color: SanggamTheme.onSurfaceDim
          );
        }
        if (num < 60) {
          return (
            label: '서맥',
            color: SanggamTheme.primary
          );
        }
        if (num > 100) {
          return (
            label: '빈맥',
            color: SanggamTheme.error
          );
        }
        return (
          label: '정상',
          color: SanggamTheme.jagaeCyan
        );
      case 'O2':
        if (num == null) {
          return (
            label: '—',
            color: SanggamTheme.onSurfaceDim
          );
        }
        if (num < 90) {
          return (
            label: '위험',
            color: SanggamTheme.error
          );
        }
        if (num < 95) {
          return (
            label: '주의',
            color: SanggamTheme.primary
          );
        }
        return (
          label: '정상',
          color: SanggamTheme.jagaeCyan
        );
      case 'Stress':
        final s = '$val'.toLowerCase();
        if (s.contains('high')) {
          return (
            label: '주의',
            color: SanggamTheme.error
          );
        }
        if (s.contains('medium')) {
          return (
            label: '보통',
            color: SanggamTheme.primary
          );
        }
        return (
          label: '양호',
          color: SanggamTheme.jagaeCyan
        );
      case 'Glucose':
        if (num == null) {
          return (
            label: '—',
            color: SanggamTheme.onSurfaceDim
          );
        }
        if (num < 70) {
          return (
            label: '저혈당',
            color: SanggamTheme.primary
          );
        }
        if (num > 140) {
          return (
            label: '고혈당',
            color: SanggamTheme.error
          );
        }
        return (
          label: '정상',
          color: SanggamTheme.jagaeCyan
        );
      default:
        return (
          label: '—',
          color: SanggamTheme.onSurfaceDim
        );
    }
  }

  String _bioIcon(String key) {
    switch (key) {
      case 'Pulse':
        return '\u2665';
      case 'O2':
        return '\u25C9';
      case 'Stress':
        return '\u25C8';
      case 'Glucose':
        return '\u2B21';
      default:
        return '\u2022';
    }
  }

  @override
  Widget build(BuildContext context) {
    final entries = values.entries.toList();

    return Column(
      crossAxisAlignment:
          CrossAxisAlignment.start,
      children: [
        Text(
          '생체 신호',
          style: TextStyle(
            color: SanggamTheme.onSurfaceDim
                .withValues(alpha: 0.8),
            fontSize: 10,
            fontWeight: FontWeight.w600,
            letterSpacing: 1.5,
          ),
        ),
        const SizedBox(height: 8),
        Wrap(
          spacing: 8,
          runSpacing: 8,
          children: entries.map((e) {
            final assess = _assessBio(
                e.key, e.value);
            return _BioSensorTile(
              icon: _bioIcon(e.key),
              label: e.key,
              value: '${e.value}',
              statusLabel: assess.label,
              statusColor: assess.color,
            );
          }).toList(),
        ),
        const SizedBox(height: 14),
        Text(
          'ECG 파형',
          style: TextStyle(
            color: SanggamTheme.onSurfaceDim
                .withValues(alpha: 0.8),
            fontSize: 10,
            fontWeight: FontWeight.w600,
            letterSpacing: 1.5,
          ),
        ),
        const SizedBox(height: 6),
        Container(
          height: 50,
          width: double.infinity,
          padding: const EdgeInsets.symmetric(
              horizontal: 8, vertical: 4),
          decoration: BoxDecoration(
            color: Colors.white
                .withValues(alpha: 0.03),
            borderRadius:
                BorderRadius.circular(10),
            border: Border.all(
              color: SanggamTheme.error
                  .withValues(alpha: 0.15),
            ),
          ),
          child: CustomPaint(
            painter: _MiniEcgPainter(
              color: SanggamTheme.error,
            ),
          ),
        ),
      ],
    );
  }
}

class _BioSensorTile
    extends StatelessWidget {
  final String icon;
  final String label;
  final String value;
  final String statusLabel;
  final Color statusColor;

  const _BioSensorTile({
    required this.icon,
    required this.label,
    required this.value,
    required this.statusLabel,
    required this.statusColor,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      width: 90,
      padding: const EdgeInsets.symmetric(
          horizontal: 10, vertical: 8),
      decoration: BoxDecoration(
        color: statusColor.withValues(
            alpha: 0.06),
        borderRadius:
            BorderRadius.circular(10),
        border: Border.all(
            color: statusColor.withValues(
                alpha: 0.18)),
      ),
      child: Column(
        crossAxisAlignment:
            CrossAxisAlignment.start,
        children: [
          Text(icon,
              style: TextStyle(
                  color: statusColor,
                  fontSize: 16)),
          const SizedBox(height: 4),
          Text(label,
              style: const TextStyle(
                  color: SanggamTheme
                      .onSurfaceDim,
                  fontSize: 10)),
          const SizedBox(height: 2),
          Text(
            value,
            style: const TextStyle(
              color: Colors.white,
              fontSize: 16,
              fontWeight: FontWeight.w700,
              fontFeatures: [
                FontFeature.tabularFigures()
              ],
            ),
          ),
          const SizedBox(height: 2),
          Text(statusLabel,
              style: TextStyle(
                  color: statusColor,
                  fontSize: 9,
                  fontWeight:
                      FontWeight.w600)),
        ],
      ),
    );
  }
}

class _MiniEcgPainter
    extends CustomPainter {
  final Color color;
  _MiniEcgPainter({required this.color});

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = color
      ..strokeWidth = 1.5
      ..style = PaintingStyle.stroke;

    final path = Path();
    final w = size.width;
    final h = size.height;
    final mid = h / 2;

    path.moveTo(0, mid);
    path.lineTo(w * 0.1, mid);
    path.quadraticBezierTo(
        w * 0.15, mid - 5, w * 0.2, mid);
    path.lineTo(w * 0.3, mid);
    path.lineTo(w * 0.35, mid + 5);
    path.lineTo(w * 0.4, mid - 20);
    path.lineTo(w * 0.45, mid + 10);
    path.lineTo(w * 0.5, mid);
    path.quadraticBezierTo(
        w * 0.6, mid - 8, w * 0.7, mid);
    path.lineTo(w, mid);

    canvas.drawPath(path, paint);
  }

  @override
  bool shouldRepaint(
          covariant CustomPainter
              oldDelegate) =>
      false;
}

// ═══════════════════════════════════════
// 공통 위젯들
// ═══════════════════════════════════════

class _GenericSensorSection
    extends StatelessWidget {
  final Map<String, dynamic> values;

  const _GenericSensorSection(
      {required this.values});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment:
          CrossAxisAlignment.start,
      children: [
        Text(
          '센서 데이터',
          style: TextStyle(
            color: SanggamTheme.onSurfaceDim
                .withValues(alpha: 0.8),
            fontSize: 10,
            fontWeight: FontWeight.w600,
            letterSpacing: 1.5,
          ),
        ),
        const SizedBox(height: 8),
        Wrap(
          spacing: 8,
          runSpacing: 8,
          children:
              values.entries.map((e) {
            return _InfoRow(
                label: e.key,
                value: '${e.value}');
          }).toList(),
        ),
      ],
    );
  }
}

class _StatusBar extends StatelessWidget {
  final IconData icon;
  final String label;
  final dynamic value;
  final Color color;

  const _StatusBar({
    required this.icon,
    required this.label,
    required this.value,
    required this.color,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(
          horizontal: 12, vertical: 8),
      decoration: BoxDecoration(
        color:
            color.withValues(alpha: 0.1),
        borderRadius:
            BorderRadius.circular(8),
      ),
      child: Row(
        children: [
          Icon(icon,
              size: 16, color: color),
          const SizedBox(width: 10),
          Text(label,
              style: const TextStyle(
                  color: SanggamTheme
                      .onSurfaceDim,
                  fontSize: 12)),
          const Spacer(),
          Text(
            '$value${label == '배터리' ? '%' : ''}',
            style: TextStyle(
                color: color,
                fontWeight: FontWeight.bold,
                fontSize: 12),
          ),
        ],
      ),
    );
  }
}

class _InfoRow extends StatelessWidget {
  final String label;
  final String value;

  const _InfoRow({
    required this.label,
    required this.value,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(
          vertical: 4),
      child: Row(
        mainAxisAlignment:
            MainAxisAlignment.spaceBetween,
        children: [
          Text(label,
              style: const TextStyle(
                  color: SanggamTheme
                      .onSurfaceDim,
                  fontSize: 12)),
          Text(value,
              style: const TextStyle(
                  color: Colors.white,
                  fontSize: 13,
                  fontWeight:
                      FontWeight.w500)),
        ],
      ),
    );
  }
}

class _DetailSparklinePainter
    extends CustomPainter {
  final List<double> data;
  final Color color;

  _DetailSparklinePainter(
      {required this.data,
      required this.color});

  @override
  void paint(Canvas canvas, Size size) {
    if (data.isEmpty) return;

    final paint = Paint()
      ..color = color
      ..strokeWidth = 2.0
      ..style = PaintingStyle.stroke;

    final path = Path();
    final w = size.width;
    final h = size.height;
    final dx = w / (data.length - 1);

    for (int i = 0; i < data.length; i++) {
      final x = i * dx;
      final y = h - (data[i] * h);
      if (i == 0) {
        path.moveTo(x, y);
      } else {
        path.lineTo(x, y);
      }
    }

    canvas.drawPath(path, paint);

    final fillPaint = Paint()
      ..shader = LinearGradient(
        begin: Alignment.topCenter,
        end: Alignment.bottomCenter,
        colors: [
          color.withValues(alpha: 0.3),
          color.withValues(alpha: 0.0),
        ],
      ).createShader(
          Rect.fromLTWH(0, 0, w, h))
      ..style = PaintingStyle.fill;

    final fillPath = Path.from(path)
      ..lineTo(w, h)
      ..lineTo(0, h)
      ..close();

    canvas.drawPath(fillPath, fillPaint);

    for (int i = 0; i < data.length; i++) {
      final x = i * dx;
      final y = h - (data[i] * h);
      canvas.drawCircle(Offset(x, y), 3,
          Paint()..color = color);
      canvas.drawCircle(Offset(x, y), 2,
          Paint()..color = Colors.white);
    }
  }

  @override
  bool shouldRepaint(
          covariant CustomPainter
              oldDelegate) =>
      false;
}
