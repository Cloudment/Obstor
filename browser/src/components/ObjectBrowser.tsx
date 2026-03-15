"use client";

import { useRouter } from "next/navigation";
import { useCallback, useRef, useState } from "react";
import { deleteObjectAction, getShareLink } from "@/lib/actions";

interface FileEntry {
  name: string;
  size: string;
  sizeBytes: number;
  lastModified: string;
  contentType: string;
}

interface Props {
  bucketName: string;
  prefix: string;
  folders: string[];
  files: FileEntry[];
}

type SortField = "name" | "size" | "lastModified";

export function ObjectBrowser({ bucketName, prefix, folders, files }: Props) {
  const router = useRouter();
  const fileInputRef = useRef<HTMLInputElement>(null);

  const [sortField, setSortField] = useState<SortField>("name");
  const [sortAsc, setSortAsc] = useState(true);
  const [selected, setSelected] = useState<Set<string>>(new Set());
  const [dragging, setDragging] = useState(false);
  const [uploading, setUploading] = useState<{ name: string; progress: number }[]>([]);
  const [shareModal, setShareModal] = useState<{ name: string; url: string } | null>(null);
  const [deleteConfirm, setDeleteConfirm] = useState<string | null>(null);
  const [filter, setFilter] = useState("");

  // Sorting
  const sortedFiles = [...files]
    .filter((f) => f.name.toLowerCase().includes(filter.toLowerCase()))
    .sort((a, b) => {
      let cmp = 0;
      if (sortField === "name") cmp = a.name.localeCompare(b.name);
      else if (sortField === "size") cmp = a.sizeBytes - b.sizeBytes;
      else cmp = new Date(a.lastModified).getTime() - new Date(b.lastModified).getTime();
      return sortAsc ? cmp : -cmp;
    });

  const filteredFolders = folders.filter((f) => f.toLowerCase().includes(filter.toLowerCase()));

  const toggleSort = (field: SortField) => {
    if (sortField === field) setSortAsc(!sortAsc);
    else {
      setSortField(field);
      setSortAsc(true);
    }
  };

  const toggleSelect = (name: string) => {
    const next = new Set(selected);
    if (next.has(name)) next.delete(name);
    else next.add(name);
    setSelected(next);
  };

  const selectAll = () => {
    if (selected.size === files.length) setSelected(new Set());
    else setSelected(new Set(files.map((f) => f.name)));
  };

  // Upload
  const uploadFiles = useCallback(
    async (fileList: FileList | File[]) => {
      const arr = Array.from(fileList);
      const entries = arr.map((f) => ({ name: f.name, progress: 0 }));
      setUploading(entries);

      for (let i = 0; i < arr.length; i++) {
        const file = arr[i];
        const objectName = prefix + file.name;
        const uploadPath = `/api/upload/${encodeURIComponent(bucketName)}/${objectName.split("/").map(encodeURIComponent).join("/")}`;

        try {
          await new Promise<void>((resolve, reject) => {
            const xhr = new XMLHttpRequest();
            xhr.open("PUT", uploadPath);
            xhr.upload.onprogress = (e) => {
              if (e.lengthComputable) {
                setUploading((prev) =>
                  prev.map((u, idx) =>
                    idx === i ? { ...u, progress: Math.round((e.loaded / e.total) * 100) } : u,
                  ),
                );
              }
            };
            xhr.onload = () =>
              xhr.status < 400 ? resolve() : reject(new Error(`HTTP ${xhr.status}`));
            xhr.onerror = () => reject(new Error("Upload failed"));
            xhr.send(file);
          });
        } catch {
          // Continue uploading remaining files
        }
      }

      setUploading([]);
      router.refresh();
    },
    [bucketName, prefix, router],
  );

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault();
      setDragging(false);
      if (e.dataTransfer.files.length > 0) uploadFiles(e.dataTransfer.files);
    },
    [uploadFiles],
  );

  // Actions
  const handleDownload = (objectName: string) => {
    const path = `/api/download/${encodeURIComponent(bucketName)}/${objectName.split("/").map(encodeURIComponent).join("/")}`;
    window.open(path, "_blank");
  };

  const handleShare = async (objectName: string) => {
    const result = await getShareLink(bucketName, objectName);
    if (result.url) setShareModal({ name: objectName, url: result.url });
  };

  const handleDelete = async (objectName: string) => {
    await deleteObjectAction(bucketName, objectName);
    setDeleteConfirm(null);
    setSelected((prev) => {
      const next = new Set(prev);
      next.delete(objectName);
      return next;
    });
    router.refresh();
  };

  const handleBulkDelete = async () => {
    for (const name of selected) {
      await deleteObjectAction(bucketName, name);
    }
    setSelected(new Set());
    router.refresh();
  };

  const displayName = (name: string) => {
    const withoutPrefix = name.startsWith(prefix) ? name.slice(prefix.length) : name;
    return withoutPrefix.replace(/\/$/, "");
  };

  return (
    <div
      onDragOver={(e) => {
        e.preventDefault();
        setDragging(true);
      }}
      onDragLeave={() => setDragging(false)}
      onDrop={handleDrop}
      className={`rounded-xl border transition-colors ${
        dragging ? "border-accent bg-accent-subtle" : "border-border bg-abyss"
      }`}
    >
      {/* Toolbar */}
      <div className="flex items-center justify-between border-border border-b px-4 py-3">
        <div className="flex items-center gap-3">
          {/* Search */}
          <div className="relative">
            <span className="icon-[lucide--search] absolute top-1/2 left-2.5 -translate-y-1/2 text-text-muted text-xs" />
            <input
              type="text"
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
              placeholder="Filter objects"
              className="rounded-md border border-border bg-surface py-1.5 pr-3 pl-8 font-mono text-text-primary text-xs outline-none placeholder:text-text-muted focus:border-accent"
            />
          </div>

          {selected.size > 0 && (
            <button
              onClick={handleBulkDelete}
              className="flex items-center gap-1.5 rounded-md border border-danger/20 bg-danger/5 px-3 py-1.5 font-body text-danger text-xs transition-colors hover:bg-danger/10"
            >
              <span className="icon-[lucide--trash-2] text-xs" />
              Delete {selected.size}
            </button>
          )}
        </div>

        <div className="flex items-center gap-2">
          <button
            onClick={() => fileInputRef.current?.click()}
            className="flex items-center gap-1.5 rounded-md bg-accent px-3 py-1.5 font-body font-medium text-black text-xs transition-colors hover:bg-accent-bright"
          >
            <span className="icon-[lucide--upload] text-xs" />
            Upload
          </button>
          <input
            ref={fileInputRef}
            type="file"
            multiple
            className="hidden"
            onChange={(e) => {
              if (e.target.files?.length) uploadFiles(e.target.files);
              e.target.value = "";
            }}
          />
        </div>
      </div>

      {/* Upload progress */}
      {uploading.length > 0 && (
        <div className="border-border border-b bg-surface px-4 py-3">
          <div className="mb-2 flex items-center gap-2">
            <span className="icon-[lucide--upload] text-accent text-xs" />
            <span className="font-body text-text-secondary text-xs">
              Uploading {uploading.length} file{uploading.length > 1 ? "s" : ""}
            </span>
          </div>
          {uploading.map((u) => (
            <div key={u.name} className="mb-1.5 last:mb-0">
              <div className="mb-1 flex items-center justify-between">
                <span className="truncate font-mono text-[11px] text-text-secondary">{u.name}</span>
                <span className="font-mono text-[10px] text-text-muted">{u.progress}%</span>
              </div>
              <div className="h-1 overflow-hidden rounded-full bg-surface-overlay">
                <div
                  className="h-full rounded-full bg-accent transition-all"
                  style={{ width: `${u.progress}%` }}
                />
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Drag overlay */}
      {dragging && (
        <div className="flex items-center justify-center border-accent/20 border-b bg-accent-subtle py-8">
          <div className="text-center">
            <span className="icon-[lucide--upload-cloud] mx-auto mb-2 block text-3xl text-accent" />
            <p className="font-body text-accent text-sm">Drop files to upload</p>
          </div>
        </div>
      )}

      {/* Column headers */}
      <div className="grid grid-cols-[auto_1fr_100px_160px_80px] items-center gap-4 border-border border-b px-4 py-2">
        <input
          type="checkbox"
          checked={selected.size === files.length && files.length > 0}
          onChange={selectAll}
          className="h-3.5 w-3.5 accent-accent"
        />
        <button
          onClick={() => toggleSort("name")}
          className="flex items-center gap-1 font-mono text-[10px] text-text-muted uppercase tracking-wider hover:text-text-secondary"
        >
          Name
          {sortField === "name" && (
            <span
              className={`text-[10px] text-accent ${sortAsc ? "icon-[lucide--chevron-up]" : "icon-[lucide--chevron-down]"}`}
            />
          )}
        </button>
        <button
          onClick={() => toggleSort("size")}
          className="flex items-center gap-1 font-mono text-[10px] text-text-muted uppercase tracking-wider hover:text-text-secondary"
        >
          Size
          {sortField === "size" && (
            <span
              className={`text-[10px] text-accent ${sortAsc ? "icon-[lucide--chevron-up]" : "icon-[lucide--chevron-down]"}`}
            />
          )}
        </button>
        <button
          onClick={() => toggleSort("lastModified")}
          className="flex items-center gap-1 font-mono text-[10px] text-text-muted uppercase tracking-wider hover:text-text-secondary"
        >
          Modified
          {sortField === "lastModified" && (
            <span
              className={`text-[10px] text-accent ${sortAsc ? "icon-[lucide--chevron-up]" : "icon-[lucide--chevron-down]"}`}
            />
          )}
        </button>
        <span />
      </div>

      {/* Rows */}
      <div className="divide-y divide-border">
        {/* Go up */}
        {prefix && (
          <a
            href={`/${encodeURIComponent(bucketName)}${
              prefix.split("/").filter(Boolean).length > 1
                ? `?prefix=${encodeURIComponent(prefix.split("/").slice(0, -2).join("/") + "/")}`
                : ""
            }`}
            className="grid grid-cols-[auto_1fr_100px_160px_80px] items-center gap-4 px-4 py-2.5 transition-colors hover:bg-surface"
          >
            <span className="h-3.5 w-3.5" />
            <span className="flex items-center gap-2 font-mono text-sm text-text-secondary">
              <span className="icon-[lucide--corner-left-up] text-sm text-text-muted" />
              ..
            </span>
            <span />
            <span />
            <span />
          </a>
        )}

        {/* Folders */}
        {filteredFolders.map((folder) => (
          <a
            key={folder}
            href={`/${encodeURIComponent(bucketName)}?prefix=${encodeURIComponent(folder)}`}
            className="grid grid-cols-[auto_1fr_100px_160px_80px] items-center gap-4 px-4 py-2.5 transition-colors hover:bg-surface"
          >
            <span className="h-3.5 w-3.5" />
            <span className="flex items-center gap-2 truncate font-mono text-sm text-text-primary">
              <span className="icon-[lucide--folder] text-accent text-sm" />
              {displayName(folder)}
            </span>
            <span className="font-mono text-text-muted text-xs">—</span>
            <span className="font-mono text-text-muted text-xs">—</span>
            <span />
          </a>
        ))}

        {/* Files */}
        {sortedFiles.map((file) => (
          <div
            key={file.name}
            className="grid grid-cols-[auto_1fr_100px_160px_80px] items-center gap-4 px-4 py-2.5 transition-colors hover:bg-surface"
          >
            <input
              type="checkbox"
              checked={selected.has(file.name)}
              onChange={() => toggleSelect(file.name)}
              className="h-3.5 w-3.5 accent-accent"
            />
            <span className="flex items-center gap-2 truncate font-mono text-sm text-text-primary">
              <span className="icon-[lucide--file] text-sm text-text-muted" />
              {displayName(file.name)}
            </span>
            <span className="font-mono text-text-secondary text-xs">{file.size}</span>
            <span className="font-mono text-text-muted text-xs">{file.lastModified}</span>
            <div className="flex items-center gap-1">
              <button
                onClick={() => handleDownload(file.name)}
                className="flex h-7 w-7 items-center justify-center rounded-md text-text-muted transition-colors hover:bg-surface-overlay hover:text-text-secondary"
                title="Download"
              >
                <span className="icon-[lucide--download] text-xs" />
              </button>
              <button
                onClick={() => handleShare(file.name)}
                className="flex h-7 w-7 items-center justify-center rounded-md text-text-muted transition-colors hover:bg-surface-overlay hover:text-text-secondary"
                title="Share"
              >
                <span className="icon-[lucide--link] text-xs" />
              </button>
              <button
                onClick={() => setDeleteConfirm(file.name)}
                className="flex h-7 w-7 items-center justify-center rounded-md text-text-muted transition-colors hover:bg-danger/10 hover:text-danger"
                title="Delete"
              >
                <span className="icon-[lucide--trash-2] text-xs" />
              </button>
            </div>
          </div>
        ))}

        {/* Empty state */}
        {filteredFolders.length === 0 && sortedFiles.length === 0 && !prefix && (
          <div className="py-12 text-center">
            <span className="icon-[lucide--upload-cloud] mx-auto mb-3 block text-3xl text-text-muted" />
            <p className="font-body text-sm text-text-muted">
              This bucket is empty. Upload files or drag and drop.
            </p>
          </div>
        )}
      </div>

      {/* Share modal */}
      {shareModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-void/80 backdrop-blur-sm">
          <div className="w-full max-w-md rounded-xl border border-border bg-surface p-6">
            <div className="mb-4 flex items-center justify-between">
              <h3 className="font-display font-semibold text-base text-text-primary">
                Share Object
              </h3>
              <button
                onClick={() => setShareModal(null)}
                className="flex h-7 w-7 items-center justify-center rounded-md text-text-muted hover:bg-surface-overlay hover:text-text-primary"
              >
                <span className="icon-[lucide--x] text-sm" />
              </button>
            </div>
            <p className="mb-3 truncate font-mono text-text-muted text-xs">{shareModal.name}</p>
            <div className="mb-4 flex items-center gap-2 rounded-lg border border-border bg-abyss p-3">
              <input
                type="text"
                value={shareModal.url}
                readOnly
                className="flex-1 bg-transparent font-mono text-text-secondary text-xs outline-none"
              />
              <button
                onClick={() => navigator.clipboard.writeText(shareModal.url)}
                className="shrink-0 rounded-md bg-accent px-3 py-1.5 font-body font-medium text-black text-xs"
              >
                Copy
              </button>
            </div>
            <p className="font-body text-[11px] text-text-muted">This link expires in 24 hours.</p>
          </div>
        </div>
      )}

      {/* Delete confirmation */}
      {deleteConfirm && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-void/80 backdrop-blur-sm">
          <div className="w-full max-w-sm rounded-xl border border-border bg-surface p-6">
            <div className="mb-1 flex items-center gap-2">
              <span className="icon-[lucide--alert-triangle] text-base text-danger" />
              <h3 className="font-display font-semibold text-base text-text-primary">
                Delete Object
              </h3>
            </div>
            <p className="mb-4 font-body text-sm text-text-secondary">
              Are you sure you want to delete{" "}
              <span className="font-mono text-text-primary">{displayName(deleteConfirm)}</span>?
              This cannot be undone.
            </p>
            <div className="flex gap-2">
              <button
                onClick={() => handleDelete(deleteConfirm)}
                className="flex-1 rounded-md bg-danger px-4 py-2 font-body font-medium text-sm text-white"
              >
                Delete
              </button>
              <button
                onClick={() => setDeleteConfirm(null)}
                className="rounded-md border border-border px-4 py-2 font-body text-sm text-text-muted"
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
