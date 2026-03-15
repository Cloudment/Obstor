import { cookies } from "next/headers";
import { NextRequest, NextResponse } from "next/server";
import { rpc } from "@/lib/rpc";

const ENDPOINT = process.env.OBSTOR_ENDPOINT || "http://localhost:9000";

export async function GET(
  _request: NextRequest,
  { params }: { params: Promise<{ bucket: string; object: string[] }> },
) {
  const cookieStore = await cookies();
  const token = cookieStore.get("obstor_token")?.value;
  if (!token) return NextResponse.json({ error: "Unauthorized" }, { status: 401 });

  const { bucket, object } = await params;
  const objectPath = object.join("/");

  // Get a short-lived URL token for the download
  let urlToken: string;
  try {
    const result = await rpc<{ token: string }>("CreateURLToken");
    urlToken = result.token;
  } catch {
    return NextResponse.json({ error: "Failed to create download token" }, { status: 500 });
  }

  const downloadUrl = `${ENDPOINT}/obstor/download/${bucket}/${objectPath}?token=${urlToken}`;

  const res = await fetch(downloadUrl, {
    headers: { "User-Agent": "Mozilla/5.0 Obstor Console" },
  });

  if (!res.ok) {
    return NextResponse.json({ error: "Download failed" }, { status: res.status });
  }

  const headers = new Headers();
  headers.set("Content-Type", res.headers.get("content-type") || "application/octet-stream");
  headers.set("Content-Disposition", `attachment; filename="${encodeURIComponent(objectPath.split("/").pop() || "file")}"`);
  if (res.headers.get("content-length")) {
    headers.set("Content-Length", res.headers.get("content-length")!);
  }

  return new NextResponse(res.body, { headers });
}
