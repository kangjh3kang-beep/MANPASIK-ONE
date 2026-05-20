import './globals.css';

export const metadata = {
  title: 'MMUP 생체 지표 예측',
  description: 'GNN 기반 다중 질환 시계열 예측 시스템',
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
