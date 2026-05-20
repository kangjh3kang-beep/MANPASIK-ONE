/**
 * @mmup-axis 1 유니버설 측정
 * @mmup-stage 1 측정
 * @sb SB-1
 * @standard IEC 62304 Class B
 *
 * Clinical 루트 레이아웃 — Tailwind CSS 적용
 */
import './globals.css';

export const metadata = {
  title: 'MMUP 임상 데이터 콘솔',
  description: 'MMUP 실시간 패킷 검증 및 해시 체인 무결성 대시보드',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ko" suppressHydrationWarning>
      <body suppressHydrationWarning className="antialiased">
        {children}
      </body>
    </html>
  );
}
