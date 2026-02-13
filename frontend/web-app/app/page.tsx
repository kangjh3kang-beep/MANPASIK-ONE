export default function Home() {
  return (
    <div className="relative flex min-h-screen flex-col items-center justify-center p-8">
      {/* ═══ Wave Ripple Background ═══ */}
      <div className="wave-ripple-bg" />

      {/* ═══ Header ═══ */}
      <header className="mb-8 text-center relative z-10">
        <h1 className="text-5xl font-bold tracking-tight mb-4">
          만파식 건강관리연구소
        </h1>
        <p className="data-label text-lg">
          MANPASIK Health Management Lab
        </p>
      </header>

      {/* ═══ Sanggam Glow Line ═══ */}
      <div className="sanggam-glow-line w-full max-w-7xl mb-8" />

      {/* ═══ Dashboard Grid ═══ */}
      <main className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 w-full max-w-7xl relative z-10">

        {/* Panel 1: 실시간 데이터 분석 */}
        <div className="sanggam-panel p-8 lg:col-span-2">
          <p className="data-label mb-3">실시간 데이터 분석</p>
          <h2 className="text-2xl mb-4">통합 관제 대시보드</h2>
          <p className="text-[var(--hanji-white)] opacity-70 leading-relaxed">
            896차원 핑거프린트 데이터의 심층 분석 및 실시간 모니터링.
            전 세계 리더기 연결 현황을 3D Globe 위젯으로 시각화합니다.
          </p>
          <div className="mt-6 flex gap-4">
            <div className="flex-1 rounded-lg bg-[var(--glass-navy)] p-4 text-center">
              <p className="data-label">연결된 디바이스</p>
              <p className="text-3xl font-bold text-[var(--wave-cyan)] mt-2">247</p>
            </div>
            <div className="flex-1 rounded-lg bg-[var(--glass-navy)] p-4 text-center">
              <p className="data-label">오늘의 측정</p>
              <p className="text-3xl font-bold text-[var(--sanggam-gold)] mt-2">1,842</p>
            </div>
            <div className="flex-1 rounded-lg bg-[var(--glass-navy)] p-4 text-center">
              <p className="data-label">시스템 상태</p>
              <p className="text-3xl font-bold text-[var(--wave-cyan)] mt-2">정상</p>
            </div>
          </div>
        </div>

        {/* Panel 2: AI 인사이트 (Breathing Animation) */}
        <div className="sanggam-panel p-8 breathing-animate">
          <p className="data-label mb-3">AI 인사이트</p>
          <h2 className="text-xl mb-4">파동 분석 요약</h2>
          <p className="text-[var(--hanji-white)] opacity-70 text-sm leading-relaxed">
            최근 24시간 동안 수집된 데이터에서 비정상 패턴이 감지되지 않았습니다.
            모든 센서가 안정적으로 작동 중입니다.
          </p>
          <div className="mt-6">
            <button className="sanggam-btn w-full">
              상세 보고서
            </button>
          </div>
        </div>

        {/* Panel 3: 건강 트렌드 */}
        <div className="sanggam-panel p-8">
          <p className="data-label mb-3">건강 트렌드</p>
          <h2 className="text-xl mb-4">히트맵 분석</h2>
          <div className="grid grid-cols-7 gap-1">
            {Array.from({ length: 35 }).map((_, i) => (
              <div
                key={i}
                className="aspect-square rounded-sm"
                style={{
                  background: `rgba(100, 255, 218, ${Math.random() * 0.8 + 0.1})`,
                }}
              />
            ))}
          </div>
          <p className="data-label mt-4">최근 5주간 바이오마커 변화</p>
        </div>

        {/* Panel 4: 카트리지 현황 */}
        <div className="sanggam-panel p-8">
          <p className="data-label mb-3">카트리지 관리</p>
          <h2 className="text-xl mb-4">재고 현황</h2>
          <ul className="space-y-3">
            {[
              { name: "건강 바이오마커", count: 128, color: "var(--sanggam-gold)" },
              { name: "환경 분석", count: 64, color: "var(--celadon-teal)" },
              { name: "식품 안전", count: 32, color: "var(--wave-cyan)" },
            ].map((item) => (
              <li key={item.name} className="flex items-center justify-between">
                <span className="flex items-center gap-2">
                  <span
                    className="w-2 h-2 rounded-full"
                    style={{ background: item.color }}
                  />
                  <span className="text-sm text-[var(--hanji-white)] opacity-80">{item.name}</span>
                </span>
                <span className="data-label">{item.count}개</span>
              </li>
            ))}
          </ul>
        </div>

        {/* Panel 5: 빠른 실행 (Wave Ripple) */}
        <div className="sanggam-panel p-8">
          <p className="data-label mb-3">빠른 실행</p>
          <h2 className="text-xl mb-4">연구 도구</h2>
          <div className="space-y-3">
            <button className="sanggam-btn w-full text-sm">측정 시작</button>
            <button className="sanggam-btn w-full text-sm">데이터 내보내기</button>
            <button className="sanggam-btn w-full text-sm">연구원 관리</button>
          </div>
        </div>
      </main>

      {/* ═══ Sanggam Glow Line ═══ */}
      <div className="sanggam-glow-line w-full max-w-7xl mt-8" />

      {/* ═══ Footer ═══ */}
      <footer className="mt-8 text-center relative z-10">
        <p className="data-label opacity-50">
          © 2026 MANPASIK World · 만파식 건강관리연구소
        </p>
      </footer>
    </div>
  );
}

