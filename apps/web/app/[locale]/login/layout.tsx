/**
 * @mmup-axis 9 전체
 * @mmup-stage 1 측정
 * @sb SB-AUTH
 *
 * 로그인 페이지 전용 layout — root layout 의 Provider 트리를 우회하여
 * 가벼운 인증 진입점만 제공.
 *
 * 인증 미들웨어가 /domains/* 접근 시 본 페이지로 redirect 시키므로, 사용자는
 * URL 이 /ko/domains/clinical 였더라도 화면에는 /login 콘텐츠가 표시됨.
 * 이때 default Next.js title 이 노출되면 사용자가 "다른 페이지가 떴다" 고
 * 인식할 수 있어, 명시적인 한국어 title 로 교체.
 */
export const metadata = {
  title: 'MMUP — 로그인',
  description: 'MMUP 통합 의료 생태계 플랫폼 — 로그인 후 페르소나별 도메인에 접근합니다.',
}

export default function LoginLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="ko">
      <body>{children}</body>
    </html>
  )
}
