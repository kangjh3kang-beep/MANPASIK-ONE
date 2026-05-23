import { handlers } from "@mmup/auth";
import type { NextRequest } from "next/server";
import { NextResponse } from "next/server";

export const GET = async (req: NextRequest) => {
  try {
    return await handlers.GET(req as any);
  } catch {
    // Cloudflare Workers에서 NextAuth 초기화 실패 시 빈 세션 반환
    const url = new URL(req.url);
    if (url.pathname.endsWith('/session')) {
      return NextResponse.json({});
    }
    return NextResponse.json({ error: 'Auth configuration error' }, { status: 500 });
  }
};

export const POST = async (req: NextRequest) => {
  try {
    return await handlers.POST(req as any);
  } catch {
    return NextResponse.json({ error: 'Auth configuration error' }, { status: 500 });
  }
};
