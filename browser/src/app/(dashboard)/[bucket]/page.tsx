import { ObjectBrowser } from "@/components/ObjectBrowser";
import { formatDate, humanSize, rpc } from "@/lib/rpc";

interface Props {
  params: Promise<{ bucket: string }>;
  searchParams: Promise<{ prefix?: string }>;
}

export default async function BucketPage({ params, searchParams }: Props) {
  const { bucket } = await params;
  const { prefix = "" } = await searchParams;
  const bucketName = decodeURIComponent(bucket);

  let objects: { name: string; size: number; lastModified: string; contentType: string }[] = [];
  let error = "";

  try {
    const result = await rpc<{
      objects: { name: string; size: number; lastModified: string; contentType: string }[] | null;
    }>("ListObjects", { bucketName, prefix, marker: "" });
    objects = result.objects || [];
  } catch (err) {
    error = err instanceof Error ? err.message : "Failed to load objects";
  }

  // Separate folders (prefixes) and files
  const folders = objects.filter((o) => o.name.endsWith("/"));
  const files = objects.filter((o) => !o.name.endsWith("/"));

  return (
    <div>
      {/* Breadcrumb */}
      <div className="mb-6 flex items-center gap-2">
        <a
          href={`/${encodeURIComponent(bucketName)}`}
          className="font-mono text-accent text-sm transition-colors hover:text-accent-bright"
        >
          {bucketName}
        </a>
        {prefix &&
          prefix
            .split("/")
            .filter(Boolean)
            .map((part, i, arr) => {
              const path = arr.slice(0, i + 1).join("/") + "/";
              return (
                <span key={path} className="flex items-center gap-2">
                  <span className="text-text-muted">/</span>
                  <a
                    href={`/${encodeURIComponent(bucketName)}?prefix=${encodeURIComponent(path)}`}
                    className="font-mono text-sm text-text-secondary transition-colors hover:text-text-primary"
                  >
                    {part}
                  </a>
                </span>
              );
            })}
      </div>

      {error ? (
        <div className="flex items-center gap-2 rounded-lg border border-danger/20 bg-danger/5 px-4 py-3">
          <span className="icon-[lucide--alert-circle] text-danger text-sm" />
          <span className="font-body text-danger text-sm">{error}</span>
        </div>
      ) : (
        <ObjectBrowser
          bucketName={bucketName}
          prefix={prefix}
          folders={folders.map((f) => f.name)}
          files={files.map((f) => ({
            name: f.name,
            size: humanSize(f.size),
            sizeBytes: f.size,
            lastModified: formatDate(f.lastModified),
            contentType: f.contentType || "",
          }))}
        />
      )}
    </div>
  );
}
