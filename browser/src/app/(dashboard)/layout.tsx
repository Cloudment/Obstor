import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { Header } from "@/components/Header";
import { Sidebar } from "@/components/Sidebar";
import { humanSize, rpc } from "@/lib/rpc";

interface StorageResult {
  used: number;
  uiVersion: string;
}

interface ServerResult {
  MinioVersion: string;
  MinioPlatform: string;
  MinioRuntime: string;
}

interface BucketResult {
  buckets: { name: string; creationDate: string }[] | null;
}

export default async function DashboardLayout({ children }: { children: React.ReactNode }) {
  const cookieStore = await cookies();
  const token = cookieStore.get("obstor_token");
  if (!token) redirect("/login");

  let buckets: { name: string; creationDate: string }[] = [];
  let storageUsed = "0 B";
  let storageBytes = 0;
  let serverVersion = "";
  let serverPlatform = "";
  let authFailed = false;

  try {
    const bucketsRes = await rpc<BucketResult>("ListBuckets");
    buckets = bucketsRes.buckets || [];
  } catch (err) {
    const msg = err instanceof Error ? err.message : "";
    if (msg.includes("Unauthorized") || msg.includes("token")) authFailed = true;
  }

  if (authFailed) {
    const cs = await cookies();
    cs.delete("obstor_token");
    redirect("/login");
  }

  try {
    const storageRes = await rpc<StorageResult>("StorageInfo");
    storageBytes = storageRes.used;
    storageUsed = humanSize(storageRes.used);
  } catch {
    // non-critical
  }

  try {
    const serverRes = await rpc<ServerResult>("ServerInfo");
    serverVersion = serverRes.MinioVersion;
    serverPlatform = serverRes.MinioPlatform;
  } catch {
    // non-critical
  }

  return (
    <div className="flex h-screen overflow-hidden">
      <Sidebar
        buckets={buckets}
        storageUsed={storageUsed}
        storageBytes={storageBytes}
        bucketCount={buckets.length}
        serverVersion={serverVersion}
        serverPlatform={serverPlatform}
      />
      <div className="flex flex-1 flex-col overflow-hidden">
        <Header />
        <main className="flex-1 overflow-y-auto bg-void p-6">{children}</main>
      </div>
    </div>
  );
}
