import type { Metadata } from "next";
import { Noto_Sans_KR } from "next/font/google";
import "./globals.css";

const notoSansKR = Noto_Sans_KR({
  variable: "--font-noto-sans-kr",
  subsets: ["latin"],
  weight: ["300", "400", "500", "700"],
});

export const metadata: Metadata = {
  title: "만파식 건강관리연구소 | MANPASIK Health Management Lab",
  description: "초정밀 차동 계측 데이터 기반 통합 건강 관제 대시보드",
  keywords: ["MANPASIK", "만파식", "건강관리연구소", "Health Management Lab", "차동 계측"],
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ko">
      <body className={`${notoSansKR.variable} antialiased`}>
        {children}
      </body>
    </html>
  );
}
