import { handlers } from "@mmup/auth";
import type { NextRequest } from "next/server";

export const GET = async (req: NextRequest) => {
  return handlers.GET(req as any);
};

export const POST = async (req: NextRequest) => {
  return handlers.POST(req as any);
};
