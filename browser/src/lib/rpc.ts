import { cookies } from "next/headers";

const RPC_ENDPOINT = process.env.OBSTOR_ENDPOINT || "http://localhost:9000";

export async function rpc<T = unknown>(
  method: string,
  params: Record<string, unknown> = {},
): Promise<T> {
  const cookieStore = await cookies();
  const token = cookieStore.get("obstor_token")?.value;

  const res = await fetch(`${RPC_ENDPOINT}/obstor/webrpc`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "User-Agent": "Mozilla/5.0 Obstor Dashboard",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    body: JSON.stringify({
      id: 1,
      jsonrpc: "2.0",
      method: `web.${method}`,
      params,
    }),
    cache: "no-store",
  });

  const text = await res.text();

  if (res.status === 401 || res.status === 403) {
    throw new Error("Unauthorized");
  }

  let data: { result?: T; error?: { message?: string } };
  try {
    data = JSON.parse(text);
  } catch {
    throw new Error(`Server returned non-JSON: ${res.status}`);
  }

  if (data.error) throw new Error(data.error.message || "RPC error");
  return data.result as T;
}

export function humanSize(bytes: number): string {
  if (bytes === 0) return "0 B";
  const units = ["B", "KB", "MB", "GB", "TB", "PB"];
  const i = Math.floor(Math.log(bytes) / Math.log(1024));
  return `${(bytes / 1024 ** i).toFixed(i > 0 ? 1 : 0)} ${units[i]}`;
}

export function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}
