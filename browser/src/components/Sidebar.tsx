"use client";

import { usePathname, useRouter } from "next/navigation";
import { useState } from "react";
import { createBucketAction, logoutAction } from "@/lib/actions";

interface Props {
  buckets: { name: string; creationDate: string }[];
  storageUsed: string;
  storageBytes: number;
  bucketCount: number;
  serverVersion: string;
  serverPlatform: string;
}

export function Sidebar({
  buckets,
  storageUsed,
  storageBytes,
  bucketCount,
  serverVersion,
  serverPlatform,
}: Props) {
  const pathname = usePathname();
  const router = useRouter();
  const activeBucket = decodeURIComponent(pathname.split("/")[1] || "");
  const [filter, setFilter] = useState("");
  const [showCreate, setShowCreate] = useState(false);
  const [createError, setCreateError] = useState("");

  const filtered = buckets.filter((b) => b.name.toLowerCase().includes(filter.toLowerCase()));

  return (
    <aside className="flex h-full w-64 shrink-0 flex-col border-r border-border bg-abyss">
      {/* Logo */}
      <div className="flex items-center gap-2.5 border-b border-border px-4 py-4">
        <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-accent">
          <span className="icon-[lucide--database] text-lg text-black" />
        </div>
        <div>
          <span className="font-display text-base font-semibold text-text-primary">Obstor</span>
          {serverVersion && (
            <p className="font-mono text-[9px] text-text-muted">{serverVersion}</p>
          )}
        </div>
      </div>

      {/* Search + Create */}
      <div className="border-b border-border p-3">
        <div className="mb-2 flex gap-2">
          <div className="relative flex-1">
            <span className="icon-[lucide--search] absolute left-2.5 top-1/2 -translate-y-1/2 text-xs text-text-muted" />
            <input
              type="text"
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
              placeholder="Filter buckets"
              className="w-full rounded-md border border-border bg-surface py-1.5 pl-8 pr-3 font-mono text-xs text-text-primary outline-none placeholder:text-text-muted focus:border-accent"
            />
          </div>
          <button
            onClick={() => setShowCreate(true)}
            className="flex h-[30px] w-[30px] shrink-0 items-center justify-center rounded-md bg-accent text-black transition-colors hover:bg-accent-bright"
            title="Create bucket"
          >
            <span className="icon-[lucide--plus] text-sm" />
          </button>
        </div>
      </div>

      {/* Create bucket */}
      {showCreate && (
        <div className="border-b border-border bg-surface p-3">
          <form
            action={async (formData) => {
              const result = await createBucketAction(formData);
              if (result?.error) {
                setCreateError(result.error);
              } else {
                setShowCreate(false);
                setCreateError("");
                router.refresh();
              }
            }}
          >
            <input
              name="bucketName"
              type="text"
              placeholder="Bucket name"
              required
              autoFocus
              pattern="[a-z0-9][a-z0-9.-]{1,61}[a-z0-9]"
              className="mb-2 w-full rounded-md border border-border bg-abyss px-3 py-1.5 font-mono text-xs text-text-primary outline-none placeholder:text-text-muted focus:border-accent"
            />
            {createError && <p className="mb-2 font-body text-[11px] text-danger">{createError}</p>}
            <div className="flex gap-2">
              <button
                type="submit"
                className="flex-1 rounded-md bg-accent px-3 py-1.5 font-body text-xs font-medium text-black"
              >
                Create
              </button>
              <button
                type="button"
                onClick={() => {
                  setShowCreate(false);
                  setCreateError("");
                }}
                className="rounded-md border border-border px-3 py-1.5 font-body text-xs text-text-muted"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Bucket list */}
      <nav className="flex-1 overflow-y-auto p-2">
        {filtered.length === 0 ? (
          <p className="px-2 py-4 text-center font-body text-xs text-text-muted">
            {buckets.length === 0 ? "No buckets yet" : "No matches"}
          </p>
        ) : (
          filtered.map((bucket) => {
            const isActive = bucket.name === activeBucket;
            return (
              <a
                key={bucket.name}
                href={`/${encodeURIComponent(bucket.name)}`}
                className={`group mb-0.5 flex items-center gap-2 truncate rounded-lg px-3 py-2 font-mono text-xs transition-colors ${
                  isActive
                    ? "bg-accent-subtle text-accent"
                    : "text-text-secondary hover:bg-surface hover:text-text-primary"
                }`}
              >
                <span className="icon-[lucide--hard-drive] shrink-0 text-sm" />
                {bucket.name}
              </a>
            );
          })
        )}
      </nav>

      {/* Stats + Logout */}
      <div className="border-t border-border p-3">
        {/* Cluster stats */}
        <div className="mb-3 space-y-1.5">
          <div className="flex items-center justify-between rounded-md bg-surface px-3 py-2">
            <span className="flex items-center gap-1.5 font-mono text-[10px] text-text-muted">
              <span className="icon-[lucide--hard-drive] text-[10px]" />
              Buckets
            </span>
            <span className="font-mono text-[10px] text-text-secondary">{bucketCount}</span>
          </div>
          <div className="flex items-center justify-between rounded-md bg-surface px-3 py-2">
            <span className="flex items-center gap-1.5 font-mono text-[10px] text-text-muted">
              <span className="icon-[lucide--database] text-[10px]" />
              Storage Used
            </span>
            <span className="font-mono text-[10px] text-text-secondary">{storageUsed}</span>
          </div>
          {serverPlatform && (
            <div className="rounded-md bg-surface px-3 py-2">
              <span className="font-mono text-[9px] leading-relaxed text-text-muted">
                {serverPlatform}
              </span>
            </div>
          )}
        </div>

        <form action={logoutAction}>
          <button
            type="submit"
            className="flex w-full items-center justify-center gap-2 rounded-md border border-border py-2 font-body text-xs text-text-muted transition-colors hover:border-border-bright hover:text-text-secondary"
          >
            <span className="icon-[lucide--log-out] text-xs" />
            Sign Out
          </button>
        </form>
      </div>
    </aside>
  );
}
