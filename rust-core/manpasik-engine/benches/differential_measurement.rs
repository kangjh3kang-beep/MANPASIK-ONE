//! 차동측정 벤치마크
//!
//! `cargo bench` 로 실행

use criterion::{black_box, criterion_group, criterion_main, Criterion};
use manpasik_engine::differential::{CorrectionParams, DifferentialEngine};
use manpasik_engine::DEFAULT_ALPHA;

fn bench_differential_88ch(c: &mut Criterion) {
    let engine = DifferentialEngine::with_defaults(88);
    let s_det: Vec<f64> = (0..88).map(|i| i as f64 * 0.1 + 1.0).collect();
    let s_ref: Vec<f64> = (0..88).map(|i| i as f64 * 0.01).collect();

    c.bench_function("differential_88ch", |b| {
        b.iter(|| engine.measure(black_box(&s_det), black_box(&s_ref)))
    });
}

fn bench_differential_896ch(c: &mut Criterion) {
    let engine = DifferentialEngine::with_defaults(896);
    let s_det: Vec<f64> = (0..896).map(|i| i as f64 * 0.1 + 1.0).collect();
    let s_ref: Vec<f64> = (0..896).map(|i| i as f64 * 0.01).collect();

    c.bench_function("differential_896ch", |b| {
        b.iter(|| engine.measure(black_box(&s_det), black_box(&s_ref)))
    });
}

fn bench_differential_1792ch(c: &mut Criterion) {
    let engine = DifferentialEngine::with_defaults(1792);
    let s_det: Vec<f64> = (0..1792).map(|i| i as f64 * 0.1 + 1.0).collect();
    let s_ref: Vec<f64> = (0..1792).map(|i| i as f64 * 0.01).collect();

    c.bench_function("differential_1792ch", |b| {
        b.iter(|| engine.measure(black_box(&s_det), black_box(&s_ref)))
    });
}

fn bench_single_measurement(c: &mut Criterion) {
    let engine = DifferentialEngine::with_defaults(1);

    c.bench_function("single_measurement", |b| {
        b.iter(|| engine.measure_single(black_box(1.234), black_box(0.012)))
    });
}

criterion_group!(
    benches,
    bench_differential_88ch,
    bench_differential_896ch,
    bench_differential_1792ch,
    bench_single_measurement
);
criterion_main!(benches);
