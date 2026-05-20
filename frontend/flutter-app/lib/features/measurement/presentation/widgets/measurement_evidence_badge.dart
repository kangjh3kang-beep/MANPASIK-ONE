import 'package:flutter/material.dart';
import 'package:manpasik/features/measurement/domain/measurement_evidence_presentation.dart';

class MeasurementEvidenceBadge extends StatelessWidget {
  const MeasurementEvidenceBadge({
    super.key,
    required this.evidenceStatus,
    required this.diagnosticReady,
    required this.evidenceGaps,
  });

  final String evidenceStatus;
  final bool diagnosticReady;
  final List<String> evidenceGaps;

  @override
  Widget build(BuildContext context) {
    final copy = MeasurementEvidencePresentation.from(
      evidenceStatus: evidenceStatus,
      diagnosticReady: diagnosticReady,
      evidenceGaps: evidenceGaps,
    );

    return DecoratedBox(
      decoration: BoxDecoration(
        color: Colors.white.withValues(alpha: 0.08),
        border: Border.all(color: Colors.white.withValues(alpha: 0.18)),
        borderRadius: BorderRadius.circular(999),
      ),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
        child: Text(
          copy.badgeLabel,
          style: const TextStyle(
            color: Colors.white,
            fontSize: 12,
            fontWeight: FontWeight.w700,
          ),
        ),
      ),
    );
  }
}
