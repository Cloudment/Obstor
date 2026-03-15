import { rpc, humanSize } from "@/lib/rpc";

interface StorageResult {
  used: number;
}

interface BucketResult {
  buckets: { name: string; creationDate: string }[] | null;
}

interface ServerResult {
  MinioVersion: string;
  MinioPlatform: string;
  MinioRuntime: string;
}

export default async function DashboardHome() {
  let used = 0;
  let bucketCount = 0;
  let version = "";
  let platform = "";
  let runtime = "";

  try {
    const s = await rpc<StorageResult>("StorageInfo");
    used = s.used;
  } catch {}

  try {
    const b = await rpc<BucketResult>("ListBuckets");
    bucketCount = b.buckets?.length || 0;
  } catch {}

  try {
    const sv = await rpc<ServerResult>("ServerInfo");
    version = sv.MinioVersion;
    platform = sv.MinioPlatform;
    runtime = sv.MinioRuntime;
  } catch {}

  return (
    <div>
      <h1 className="mb-6 font-display text-2xl font-semibold text-text-primary">Cluster Overview</h1>

      {/* Stats grid */}
      <div className="mb-6 grid grid-cols-2 gap-4 lg:grid-cols-4">
        <div className="rounded-xl border border-border bg-abyss p-5">
          <div className="mb-2 flex items-center gap-2">
            <span className="icon-[lucide--database] text-sm text-accent" />
            <span className="font-mono text-[10px] uppercase tracking-wider text-text-muted">
              Storage Used
            </span>
          </div>
          <p className="font-display text-2xl font-bold text-text-primary">{humanSize(used)}</p>
          <p className="mt-1 font-mono text-[10px] text-text-muted">
            Across all nodes &middot; Replicated
          </p>
        </div>

        <div className="rounded-xl border border-border bg-abyss p-5">
          <div className="mb-2 flex items-center gap-2">
            <span className="icon-[lucide--hard-drive] text-sm text-accent" />
            <span className="font-mono text-[10px] uppercase tracking-wider text-text-muted">
              Buckets
            </span>
          </div>
          <p className="font-display text-2xl font-bold text-text-primary">{bucketCount}</p>
          <p className="mt-1 font-mono text-[10px] text-text-muted">
            Sharded across cluster
          </p>
        </div>

        <div className="rounded-xl border border-border bg-abyss p-5">
          <div className="mb-2 flex items-center gap-2">
            <span className="icon-[lucide--server] text-sm text-accent" />
            <span className="font-mono text-[10px] uppercase tracking-wider text-text-muted">
              Replication
            </span>
          </div>
          <p className="font-display text-2xl font-bold text-text-primary">2x min</p>
          <p className="mt-1 font-mono text-[10px] text-text-muted">
            Objects stored in 2+ locations
          </p>
        </div>

        <div className="rounded-xl border border-border bg-abyss p-5">
          <div className="mb-2 flex items-center gap-2">
            <span className="icon-[lucide--shield-check] text-sm text-accent" />
            <span className="font-mono text-[10px] uppercase tracking-wider text-text-muted">
              Erasure Coding
            </span>
          </div>
          <p className="font-display text-2xl font-bold text-text-primary">Active</p>
          <p className="mt-1 font-mono text-[10px] text-text-muted">
            Data protected against node failure
          </p>
        </div>
      </div>

      {/* Server info */}
      <div className="mb-6 rounded-xl border border-border bg-abyss">
        <div className="flex items-center gap-2 border-b border-border px-5 py-3">
          <span className="icon-[lucide--monitor] text-sm text-text-muted" />
          <span className="font-display text-sm font-semibold text-text-primary">
            Server Info
          </span>
        </div>
        <div className="grid gap-px bg-border sm:grid-cols-3">
          <div className="bg-abyss px-5 py-4">
            <span className="font-mono text-[10px] uppercase tracking-wider text-text-muted">
              Version
            </span>
            <p className="mt-1 font-mono text-xs text-text-secondary">{version || "—"}</p>
          </div>
          <div className="bg-abyss px-5 py-4">
            <span className="font-mono text-[10px] uppercase tracking-wider text-text-muted">
              Platform
            </span>
            <p className="mt-1 font-mono text-xs text-text-secondary">{platform || "—"}</p>
          </div>
          <div className="bg-abyss px-5 py-4">
            <span className="font-mono text-[10px] uppercase tracking-wider text-text-muted">
              Runtime
            </span>
            <p className="mt-1 font-mono text-xs text-text-secondary">{runtime || "—"}</p>
          </div>
        </div>
      </div>

      {/* Distribution note */}
      <div className="rounded-xl border border-border bg-abyss p-5">
        <div className="flex items-start gap-3">
          <span className="icon-[lucide--globe] mt-0.5 shrink-0 text-base text-accent" />
          <div>
            <h3 className="mb-1 font-display text-sm font-semibold text-text-primary">
              Distributed Storage
            </h3>
            <p className="font-body text-xs leading-relaxed text-text-muted">
              Objects are automatically sharded and replicated across connected
              nodes. Each object is stored in at least 2 global locations.
              Node capacity varies — the cluster aggregates all available
              storage and balances data transparently.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
