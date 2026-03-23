import { cookies } from "next/headers";
import { type NextRequest, NextResponse } from "next/server";

const ENDPOINT = process.env.OBSTOR_ENDPOINT || "http://localhost:9000";

export async function PUT(
  request: NextRequest,
  { params }: { params: Promise<{ bucket: string; object: string[] }> },
) {
  const cookieStore = await cookies();
  const token = cookieStore.get("obstor_token")?.value;
  if (!token) return NextResponse.json({ error: "Unauthorized" }, { status: 401 });

  const { bucket, object } = await params;
  const objectPath = object.join("/");

  const body = await request.arrayBuffer();

  const res = await fetch(`${ENDPOINT}/obstor/upload/${bucket}/${objectPath}`, {
    method: "PUT",
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": request.headers.get("content-type") || "application/octet-stream",
      "User-Agent": "Mozilla/5.0 Obstor Console",
    },
    body,
  });

  if (!res.ok) {
    const text = await res.text();
    return NextResponse.json({ error: text }, { status: res.status });
  }

  return NextResponse.json({ success: true });
}
