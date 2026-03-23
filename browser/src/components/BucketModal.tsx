"use client";

import { AnimatePresence, domAnimation, LazyMotion, m } from "framer-motion";
import { useCallback, useEffect, useMemo, useReducer, useState } from "react";
import {
  createBucketWithSettingsAction,
  getBucketSettingsAction,
  updateBucketSettingsAction,
} from "@/lib/actions";

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

export interface BucketSettings {
  name: string;
  versioning: boolean;
  objectLocking: boolean;
  accessPolicy: "private" | "public-read" | "public-read-write" | "custom";
  customPolicy: string;
  sftpEnabled: boolean;
  s3Enabled: boolean;
  anonymousAccess: boolean;
  quotaEnabled: boolean;
  quotaType: "hard" | "fifo";
  quotaSize: string;
  quotaUnit: "GB" | "TB" | "PB";
  encryptionEnabled: boolean;
  encryptionType: "SSE-S3" | "SSE-KMS";
  kmsKeyId: string;
  tags: { key: string; value: string }[];
  placementStrategy: "smart" | "custom";
  regions: string[];
}

interface Props {
  open: boolean;
  onClose: () => void;
  onSuccess: () => void;
  editBucket?: string | null;
}

const TABS = [
  { id: "general", label: "General", icon: "icon-[lucide--settings-2]" },
  { id: "access", label: "Access", icon: "icon-[lucide--shield]" },
  { id: "features", label: "Features", icon: "icon-[lucide--toggle-right]" },
  { id: "quota", label: "Quota", icon: "icon-[lucide--gauge]" },
  { id: "encryption", label: "Encryption", icon: "icon-[lucide--lock]" },
  { id: "tags", label: "Tags", icon: "icon-[lucide--tag]" },
  { id: "region", label: "Region", icon: "icon-[lucide--map-pin]" },
] as const;

type TabId = (typeof TABS)[number]["id"];

const DEFAULT_POLICY_TEMPLATE = (bucketName: string) =>
  JSON.stringify(
    {
      Version: "2012-10-17",
      Statement: [
        {
          Effect: "Allow",
          Action: [
            "admin:DataUsageInfo",
            "admin:GetBucketQuota",
            "admin:GetBucketTarget",
            "admin:Heal",
            "admin:SetBucketTarget",
            "admin:TopLocksInfo",
          ],
          Resource: [`arn:aws:s3:::${bucketName || "BUCKET_NAME"}*`],
        },
        {
          Effect: "Allow",
          Action: ["s3:*"],
          Resource: [`arn:aws:s3:::${bucketName || "BUCKET_NAME"}*`],
        },
      ],
    },
    null,
    2,
  );

const EMPTY_SETTINGS: BucketSettings = {
  name: "",
  versioning: false,
  objectLocking: false,
  accessPolicy: "private",
  customPolicy: "",
  sftpEnabled: false,
  s3Enabled: true,
  anonymousAccess: false,
  quotaEnabled: false,
  quotaType: "hard",
  quotaSize: "",
  quotaUnit: "GB",
  encryptionEnabled: false,
  encryptionType: "SSE-S3",
  kmsKeyId: "",
  tags: [],
  placementStrategy: "smart",
  regions: [],
};

// ---------------------------------------------------------------------------
// Sub-components
// ---------------------------------------------------------------------------

function Toggle({
  checked,
  onChange,
  disabled,
}: {
  checked: boolean;
  onChange: (v: boolean) => void;
  disabled?: boolean;
}) {
  return (
    <button
      type="button"
      role="switch"
      aria-checked={checked}
      disabled={disabled}
      onClick={() => onChange(!checked)}
      className={`relative inline-flex h-6 w-11 shrink-0 cursor-pointer items-center rounded-full transition-colors duration-200 ${
        disabled ? "cursor-not-allowed opacity-40" : ""
      } ${checked ? "bg-accent" : "bg-surface-overlay"}`}
    >
      <span
        className={`pointer-events-none inline-block h-4 w-4 rounded-full bg-text-primary shadow-sm transition-transform duration-200 ${
          checked ? "translate-x-6" : "translate-x-1"
        }`}
      />
    </button>
  );
}

function SettingRow({
  label,
  description,
  children,
}: {
  label: string;
  description?: string;
  children: React.ReactNode;
}) {
  return (
    <div className="flex items-center justify-between gap-4 rounded-lg border border-border bg-surface/50 px-4 py-3">
      <div className="min-w-0">
        <p className="font-body text-sm text-text-primary">{label}</p>
        {description && (
          <p className="mt-0.5 font-body text-[11px] text-text-muted leading-relaxed">
            {description}
          </p>
        )}
      </div>
      <div className="shrink-0">{children}</div>
    </div>
  );
}

function SectionLabel({ children }: { children: React.ReactNode }) {
  return (
    <p className="mb-2 font-mono text-[10px] text-text-muted uppercase tracking-wider">
      {children}
    </p>
  );
}

// ---------------------------------------------------------------------------
// Tab content components
// ---------------------------------------------------------------------------

function GeneralTab({
  settings,
  onChange,
  isEdit,
}: {
  settings: BucketSettings;
  onChange: (s: Partial<BucketSettings>) => void;
  isEdit: boolean;
}) {
  return (
    <div className="space-y-4">
      <div>
        <SectionLabel>Bucket Name</SectionLabel>
        <input
          type="text"
          value={settings.name}
          disabled={isEdit}
          onChange={(e) =>
            onChange({ name: e.target.value.toLowerCase().replace(/[^a-z0-9.-]/g, "") })
          }
          placeholder="my-bucket"
          className="w-full rounded-lg border border-border bg-surface px-4 py-2.5 font-mono text-sm text-text-primary outline-none transition-colors placeholder:text-text-muted focus:border-accent disabled:cursor-not-allowed disabled:opacity-50"
        />
        {!isEdit && settings.name && (
          <p className="mt-1.5 font-mono text-[10px] text-text-muted">
            arn:aws:s3:::{settings.name}
          </p>
        )}
      </div>

      <div className="space-y-2">
        <SectionLabel>Storage Configuration</SectionLabel>
        <SettingRow
          label="Versioning"
          description="Keep multiple versions of an object in the same bucket"
        >
          <Toggle checked={settings.versioning} onChange={(v) => onChange({ versioning: v })} />
        </SettingRow>
        <SettingRow
          label="Object Locking"
          description="Prevent objects from being deleted. Can only be enabled at bucket creation."
        >
          <Toggle
            checked={settings.objectLocking}
            onChange={(v) => onChange({ objectLocking: v })}
            disabled={isEdit}
          />
        </SettingRow>
      </div>
    </div>
  );
}

function AccessTab({
  settings,
  onChange,
}: {
  settings: BucketSettings;
  onChange: (s: Partial<BucketSettings>) => void;
}) {
  const policies = [
    { id: "private", label: "Private", desc: "No public access" },
    { id: "public-read", label: "Public Read", desc: "Anyone can read objects" },
    { id: "public-read-write", label: "Public Read/Write", desc: "Anyone can read and write" },
    { id: "custom", label: "Custom Policy", desc: "Define a custom IAM policy" },
  ] as const;

  return (
    <div className="space-y-4">
      <div>
        <SectionLabel>Access Policy</SectionLabel>
        <div className="grid grid-cols-2 gap-2">
          {policies.map((p) => (
            <button
              key={p.id}
              type="button"
              onClick={() => {
                const update: Partial<BucketSettings> = { accessPolicy: p.id };
                if (p.id === "custom" && !settings.customPolicy) {
                  update.customPolicy = DEFAULT_POLICY_TEMPLATE(settings.name);
                }
                onChange(update);
              }}
              className={`rounded-lg border px-3 py-2.5 text-left transition-all ${
                settings.accessPolicy === p.id
                  ? "border-accent/40 bg-accent-subtle"
                  : "border-border bg-surface/50 hover:border-border-bright"
              }`}
            >
              <p
                className={`font-body font-medium text-xs ${settings.accessPolicy === p.id ? "text-accent-bright" : "text-text-primary"}`}
              >
                {p.label}
              </p>
              <p className="mt-0.5 font-body text-[10px] text-text-muted">{p.desc}</p>
            </button>
          ))}
        </div>
      </div>

      <AnimatePresence>
        {settings.accessPolicy === "custom" && (
          <m.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: "auto", opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.2 }}
            className="overflow-hidden"
          >
            <SectionLabel>Policy Document (JSON)</SectionLabel>
            <div className="relative">
              <textarea
                value={settings.customPolicy}
                onChange={(e) => onChange({ customPolicy: e.target.value })}
                rows={16}
                spellCheck={false}
                className="w-full rounded-lg border border-border bg-void p-4 font-mono text-[11px] text-text-secondary leading-relaxed outline-none transition-colors focus:border-accent"
              />
              <button
                type="button"
                onClick={() => {
                  try {
                    const formatted = JSON.stringify(JSON.parse(settings.customPolicy), null, 2);
                    onChange({ customPolicy: formatted });
                  } catch {
                    /* invalid JSON, ignore */
                  }
                }}
                className="absolute top-2 right-2 rounded-md bg-surface-overlay px-2 py-1 font-mono text-[10px] text-text-muted transition-colors hover:text-text-secondary"
              >
                Format
              </button>
            </div>
          </m.div>
        )}
      </AnimatePresence>
    </div>
  );
}

function FeaturesTab({
  settings,
  onChange,
}: {
  settings: BucketSettings;
  onChange: (s: Partial<BucketSettings>) => void;
}) {
  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <SectionLabel>Protocol Access</SectionLabel>
        <SettingRow label="S3 API" description="Access bucket via the S3-compatible API">
          <Toggle checked={settings.s3Enabled} onChange={(v) => onChange({ s3Enabled: v })} />
        </SettingRow>
        <SettingRow label="SFTP Access" description="Allow file access over SFTP protocol">
          <Toggle checked={settings.sftpEnabled} onChange={(v) => onChange({ sftpEnabled: v })} />
        </SettingRow>
      </div>

      <div className="space-y-2">
        <SectionLabel>Visibility</SectionLabel>
        <SettingRow
          label="Public URL Access"
          description="Make bucket contents accessible via a direct public URL without authentication"
        >
          <Toggle
            checked={settings.accessPolicy !== "private"}
            onChange={(v) =>
              onChange({
                accessPolicy: v ? "public-read" : "private",
                anonymousAccess: v,
              })
            }
          />
        </SettingRow>
      </div>
    </div>
  );
}

function QuotaTab({
  settings,
  onChange,
}: {
  settings: BucketSettings;
  onChange: (s: Partial<BucketSettings>) => void;
}) {
  return (
    <div className="space-y-4">
      <SettingRow
        label="Enable Quota"
        description="Limit the amount of data that can be stored in this bucket"
      >
        <Toggle checked={settings.quotaEnabled} onChange={(v) => onChange({ quotaEnabled: v })} />
      </SettingRow>

      <AnimatePresence>
        {settings.quotaEnabled && (
          <m.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: "auto", opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.2 }}
            className="space-y-3 overflow-hidden"
          >
            <div>
              <SectionLabel>Quota Type</SectionLabel>
              <div className="grid grid-cols-2 gap-2">
                {(["hard", "fifo"] as const).map((t) => (
                  <button
                    key={t}
                    type="button"
                    onClick={() => onChange({ quotaType: t })}
                    className={`rounded-lg border px-3 py-2.5 text-left transition-all ${
                      settings.quotaType === t
                        ? "border-accent/40 bg-accent-subtle"
                        : "border-border bg-surface/50 hover:border-border-bright"
                    }`}
                  >
                    <p
                      className={`font-body font-medium text-xs ${settings.quotaType === t ? "text-accent-bright" : "text-text-primary"}`}
                    >
                      {t === "hard" ? "Hard Limit" : "FIFO"}
                    </p>
                    <p className="mt-0.5 font-body text-[10px] text-text-muted">
                      {t === "hard"
                        ? "Reject writes when limit is reached"
                        : "Delete oldest objects when limit is reached"}
                    </p>
                  </button>
                ))}
              </div>
            </div>

            <div>
              <SectionLabel>Quota Size</SectionLabel>
              <div className="flex gap-2">
                <input
                  type="number"
                  min="1"
                  value={settings.quotaSize}
                  onChange={(e) => onChange({ quotaSize: e.target.value })}
                  placeholder="100"
                  className="flex-1 rounded-lg border border-border bg-surface px-4 py-2.5 font-mono text-sm text-text-primary outline-none transition-colors placeholder:text-text-muted focus:border-accent"
                />
                <select
                  value={settings.quotaUnit}
                  onChange={(e) => onChange({ quotaUnit: e.target.value as "GB" | "TB" | "PB" })}
                  className="rounded-lg border border-border bg-surface px-3 py-2.5 font-mono text-sm text-text-primary outline-none transition-colors focus:border-accent"
                >
                  <option value="GB">GB</option>
                  <option value="TB">TB</option>
                  <option value="PB">PB</option>
                </select>
              </div>
            </div>
          </m.div>
        )}
      </AnimatePresence>
    </div>
  );
}

function EncryptionTab({
  settings,
  onChange,
}: {
  settings: BucketSettings;
  onChange: (s: Partial<BucketSettings>) => void;
}) {
  return (
    <div className="space-y-4">
      <SettingRow
        label="Server-Side Encryption"
        description="Automatically encrypt objects when stored"
      >
        <Toggle
          checked={settings.encryptionEnabled}
          onChange={(v) => onChange({ encryptionEnabled: v })}
        />
      </SettingRow>

      <AnimatePresence>
        {settings.encryptionEnabled && (
          <m.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: "auto", opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.2 }}
            className="space-y-3 overflow-hidden"
          >
            <div>
              <SectionLabel>Encryption Type</SectionLabel>
              <div className="grid grid-cols-2 gap-2">
                {(["SSE-S3", "SSE-KMS"] as const).map((t) => (
                  <button
                    key={t}
                    type="button"
                    onClick={() => onChange({ encryptionType: t })}
                    className={`rounded-lg border px-3 py-2.5 text-left transition-all ${
                      settings.encryptionType === t
                        ? "border-accent/40 bg-accent-subtle"
                        : "border-border bg-surface/50 hover:border-border-bright"
                    }`}
                  >
                    <p
                      className={`font-body font-medium text-xs ${settings.encryptionType === t ? "text-accent-bright" : "text-text-primary"}`}
                    >
                      {t}
                    </p>
                    <p className="mt-0.5 font-body text-[10px] text-text-muted">
                      {t === "SSE-S3"
                        ? "Server-managed encryption keys"
                        : "External KMS-managed keys"}
                    </p>
                  </button>
                ))}
              </div>
            </div>

            <AnimatePresence>
              {settings.encryptionType === "SSE-KMS" && (
                <m.div
                  initial={{ height: 0, opacity: 0 }}
                  animate={{ height: "auto", opacity: 1 }}
                  exit={{ height: 0, opacity: 0 }}
                  transition={{ duration: 0.15 }}
                  className="overflow-hidden"
                >
                  <SectionLabel>KMS Key ID</SectionLabel>
                  <input
                    type="text"
                    value={settings.kmsKeyId}
                    onChange={(e) => onChange({ kmsKeyId: e.target.value })}
                    placeholder="arn:aws:kms:region:account:key/key-id"
                    className="w-full rounded-lg border border-border bg-surface px-4 py-2.5 font-mono text-sm text-text-primary outline-none transition-colors placeholder:text-text-muted focus:border-accent"
                  />
                </m.div>
              )}
            </AnimatePresence>
          </m.div>
        )}
      </AnimatePresence>
    </div>
  );
}

function TagsTab({
  settings,
  onChange,
}: {
  settings: BucketSettings;
  onChange: (s: Partial<BucketSettings>) => void;
}) {
  const addTag = () => onChange({ tags: [...settings.tags, { key: "", value: "" }] });
  const removeTag = (i: number) => onChange({ tags: settings.tags.filter((_, idx) => idx !== i) });
  const updateTag = (i: number, field: "key" | "value", val: string) =>
    onChange({
      tags: settings.tags.map((t, idx) => (idx === i ? { ...t, [field]: val } : t)),
    });

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <SectionLabel>Bucket Tags</SectionLabel>
        <button
          type="button"
          onClick={addTag}
          className="flex items-center gap-1.5 rounded-md bg-surface-overlay px-2.5 py-1 font-body text-[11px] text-text-secondary transition-colors hover:text-text-primary"
        >
          <span className="icon-[lucide--plus] text-[10px]" />
          Add Tag
        </button>
      </div>

      {settings.tags.length === 0 ? (
        <div className="flex flex-col items-center justify-center rounded-lg border border-border border-dashed py-8">
          <span className="icon-[lucide--tag] mb-2 text-lg text-text-muted" />
          <p className="font-body text-text-muted text-xs">No tags configured</p>
          <button
            type="button"
            onClick={addTag}
            className="mt-2 font-body text-accent text-xs transition-colors hover:text-accent-bright"
          >
            Add your first tag
          </button>
        </div>
      ) : (
        <div className="space-y-2">
          {settings.tags.map((tag, i) => (
            <div key={`tag-${tag.key || i}`} className="flex items-center gap-2">
              <input
                type="text"
                value={tag.key}
                onChange={(e) => updateTag(i, "key", e.target.value)}
                placeholder="Key"
                className="flex-1 rounded-lg border border-border bg-surface px-3 py-2 font-mono text-text-primary text-xs outline-none transition-colors placeholder:text-text-muted focus:border-accent"
              />
              <input
                type="text"
                value={tag.value}
                onChange={(e) => updateTag(i, "value", e.target.value)}
                placeholder="Value"
                className="flex-1 rounded-lg border border-border bg-surface px-3 py-2 font-mono text-text-primary text-xs outline-none transition-colors placeholder:text-text-muted focus:border-accent"
              />
              <button
                type="button"
                onClick={() => removeTag(i)}
                className="flex h-8 w-8 shrink-0 items-center justify-center rounded-md text-text-muted transition-colors hover:bg-danger/10 hover:text-danger"
              >
                <span className="icon-[lucide--x] text-xs" />
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

function RegionTab({
  settings,
  onChange,
}: {
  settings: BucketSettings;
  onChange: (s: Partial<BucketSettings>) => void;
}) {
  const [regionInput, setRegionInput] = useState("");

  const strategies = [
    {
      id: "smart",
      label: "Smart Placement",
      desc: "Automatically distributes data across available regions for optimal durability and performance",
    },
    {
      id: "custom",
      label: "Custom Regions",
      desc: "Manually select which regions your data should be stored in",
    },
  ] as const;

  const addRegion = (value: string) => {
    const trimmed = value.trim().toLowerCase();
    if (trimmed && !settings.regions.includes(trimmed)) {
      onChange({ regions: [...settings.regions, trimmed] });
    }
    setRegionInput("");
  };

  const removeRegion = (index: number) => {
    onChange({ regions: settings.regions.filter((_, i) => i !== index) });
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      e.preventDefault();
      addRegion(regionInput);
    }
  };

  return (
    <div className="space-y-4">
      <div>
        <SectionLabel>Placement Strategy</SectionLabel>
        <div className="grid grid-cols-2 gap-2">
          {strategies.map((s) => (
            <button
              key={s.id}
              type="button"
              onClick={() => onChange({ placementStrategy: s.id })}
              className={`rounded-lg border px-3 py-2.5 text-left transition-all ${
                settings.placementStrategy === s.id
                  ? "border-accent/40 bg-accent-subtle"
                  : "border-border bg-surface/50 hover:border-border-bright"
              }`}
            >
              <p
                className={`font-body font-medium text-xs ${settings.placementStrategy === s.id ? "text-accent-bright" : "text-text-primary"}`}
              >
                {s.label}
              </p>
              <p className="mt-0.5 font-body text-[10px] text-text-muted">{s.desc}</p>
            </button>
          ))}
        </div>
      </div>

      <AnimatePresence>
        {settings.placementStrategy === "custom" && (
          <m.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: "auto", opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.2 }}
            className="space-y-3 overflow-hidden"
          >
            <div>
              <SectionLabel>Regions</SectionLabel>
              <div className="flex flex-wrap items-center gap-1.5 rounded-lg border border-border bg-surface px-3 py-2">
                {settings.regions.map((region, i) => (
                  <span
                    key={region}
                    className="flex items-center gap-1 rounded bg-surface-overlay px-2 py-0.5 font-mono text-[11px] text-text-secondary"
                  >
                    {region}
                    <button
                      type="button"
                      onClick={() => removeRegion(i)}
                      className="ml-0.5 text-text-muted transition-colors hover:text-text-primary"
                    >
                      <span className="icon-[lucide--x] text-[9px]" />
                    </button>
                  </span>
                ))}
                <input
                  type="text"
                  value={regionInput}
                  onChange={(e) => setRegionInput(e.target.value)}
                  onKeyDown={handleKeyDown}
                  placeholder="Type a region and press Enter..."
                  className="min-w-[180px] flex-1 bg-transparent py-0.5 font-mono text-sm text-text-primary outline-none placeholder:text-text-muted"
                />
              </div>
            </div>

            <p className="font-body text-[11px] text-text-muted leading-relaxed">
              Data will only be replicated to the selected regions. At least 2 regions are
              recommended for durability.
            </p>
          </m.div>
        )}
      </AnimatePresence>
    </div>
  );
}

// ---------------------------------------------------------------------------
// Reducer for main modal state
// ---------------------------------------------------------------------------

interface ModalState {
  tab: TabId;
  settings: BucketSettings;
  loading: boolean;
  saving: boolean;
  error: string;
}

type ModalAction =
  | { type: "SET_TAB"; tab: TabId }
  | { type: "UPDATE_SETTINGS"; partial: Partial<BucketSettings> }
  | { type: "SET_SETTINGS"; settings: BucketSettings }
  | { type: "SET_LOADING"; loading: boolean }
  | { type: "SET_SAVING"; saving: boolean }
  | { type: "SET_ERROR"; error: string }
  | { type: "RESET"; editBucket: string | null };

function modalReducer(state: ModalState, action: ModalAction): ModalState {
  switch (action.type) {
    case "SET_TAB":
      return { ...state, tab: action.tab };
    case "UPDATE_SETTINGS":
      return { ...state, settings: { ...state.settings, ...action.partial }, error: "" };
    case "SET_SETTINGS":
      return { ...state, settings: action.settings };
    case "SET_LOADING":
      return { ...state, loading: action.loading };
    case "SET_SAVING":
      return { ...state, saving: action.saving };
    case "SET_ERROR":
      return { ...state, error: action.error };
    case "RESET":
      return {
        tab: "general",
        settings: EMPTY_SETTINGS,
        loading: !!action.editBucket,
        saving: false,
        error: "",
      };
    default:
      return state;
  }
}

const INITIAL_STATE: ModalState = {
  tab: "general",
  settings: EMPTY_SETTINGS,
  loading: false,
  saving: false,
  error: "",
};

// ---------------------------------------------------------------------------
// Main Modal
// ---------------------------------------------------------------------------

export function BucketModal({ open, onClose, onSuccess, editBucket }: Props) {
  const isEdit = !!editBucket;
  const [state, dispatch] = useReducer(modalReducer, INITIAL_STATE);
  const { tab, settings, loading, saving, error } = state;

  const update = useCallback((partial: Partial<BucketSettings>) => {
    dispatch({ type: "UPDATE_SETTINGS", partial });
  }, []);

  // Load settings when editing
  useEffect(() => {
    if (!open) return;
    dispatch({ type: "RESET", editBucket: editBucket ?? null });

    if (editBucket) {
      getBucketSettingsAction(editBucket)
        .then((res) => {
          if ("error" in res) {
            dispatch({ type: "SET_ERROR", error: res.error });
          } else {
            dispatch({ type: "SET_SETTINGS", settings: { ...EMPTY_SETTINGS, ...res } });
          }
        })
        .finally(() => dispatch({ type: "SET_LOADING", loading: false }));
    }
  }, [open, editBucket]);

  // Update custom policy ARN when name changes
  const policyWithBucketName = useMemo(() => {
    if (settings.accessPolicy !== "custom" || isEdit) return settings.customPolicy;
    try {
      const parsed = JSON.parse(settings.customPolicy);
      let changed = false;
      for (const stmt of parsed.Statement || []) {
        const newResource = [`arn:aws:s3:::${settings.name || "BUCKET_NAME"}*`];
        if (JSON.stringify(stmt.Resource) !== JSON.stringify(newResource)) {
          stmt.Resource = newResource;
          changed = true;
        }
      }
      return changed ? JSON.stringify(parsed, null, 2) : settings.customPolicy;
    } catch {
      return settings.customPolicy;
    }
  }, [settings.customPolicy, settings.name, settings.accessPolicy, isEdit]);

  useEffect(() => {
    if (policyWithBucketName !== settings.customPolicy && !isEdit) {
      dispatch({
        type: "SET_SETTINGS",
        settings: { ...settings, customPolicy: policyWithBucketName },
      });
    }
  }, [policyWithBucketName, settings, isEdit]);

  const handleSubmit = async () => {
    if (!isEdit && !settings.name) {
      dispatch({ type: "SET_ERROR", error: "Bucket name is required" });
      dispatch({ type: "SET_TAB", tab: "general" });
      return;
    }

    dispatch({ type: "SET_SAVING", saving: true });
    dispatch({ type: "SET_ERROR", error: "" });

    try {
      const result = isEdit
        ? await updateBucketSettingsAction(settings)
        : await createBucketWithSettingsAction(settings);

      if (result && "error" in result) {
        dispatch({ type: "SET_ERROR", error: result.error });
      } else {
        onSuccess();
        onClose();
      }
    } catch (err) {
      dispatch({
        type: "SET_ERROR",
        error: err instanceof Error ? err.message : "An error occurred",
      });
    } finally {
      dispatch({ type: "SET_SAVING", saving: false });
    }
  };

  // Close on Escape
  useEffect(() => {
    if (!open) return;
    const handler = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, [open, onClose]);

  return (
    <AnimatePresence>
      {open && (
        <LazyMotion features={domAnimation}>
          <m.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.15 }}
            className="fixed inset-0 z-50 flex items-center justify-center p-4"
          >
            {/* Backdrop */}
            <m.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              className="absolute inset-0 bg-void/80 backdrop-blur-sm"
              onClick={onClose}
            />

            {/* Modal */}
            <m.div
              initial={{ opacity: 0, scale: 0.96, y: 8 }}
              animate={{ opacity: 1, scale: 1, y: 0 }}
              exit={{ opacity: 0, scale: 0.96, y: 8 }}
              transition={{ duration: 0.2, ease: [0.16, 1, 0.3, 1] }}
              className="relative flex h-[70vh] w-full max-w-2xl flex-col overflow-hidden rounded-xl border border-border bg-abyss shadow-2xl shadow-black/50"
            >
              {/* Header */}
              <div className="flex items-center justify-between border-border border-b px-6 py-4">
                <div className="flex items-center gap-3">
                  <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-accent/10">
                    <span
                      className={`text-accent text-sm ${isEdit ? "icon-[lucide--settings]" : "icon-[lucide--plus-circle]"}`}
                    />
                  </div>
                  <div>
                    <h2 className="font-display font-semibold text-base text-text-primary">
                      {isEdit ? "Bucket Settings" : "Create Bucket"}
                    </h2>
                    {isEdit && (
                      <p className="font-mono text-[10px] text-text-muted">{editBucket}</p>
                    )}
                  </div>
                </div>
                <button
                  type="button"
                  onClick={onClose}
                  className="flex h-8 w-8 items-center justify-center rounded-md text-text-muted transition-colors hover:bg-surface-overlay hover:text-text-secondary"
                >
                  <span className="icon-[lucide--x] text-sm" />
                </button>
              </div>

              {/* Tabs */}
              <div className="flex gap-1 overflow-x-auto border-border border-b bg-surface/30 px-6 py-1.5">
                {TABS.map((t) => (
                  <button
                    key={t.id}
                    type="button"
                    onClick={() => dispatch({ type: "SET_TAB", tab: t.id })}
                    className={`flex items-center gap-1.5 whitespace-nowrap rounded-md px-3 py-1.5 font-body text-xs transition-all ${
                      tab === t.id
                        ? "bg-accent/10 text-accent"
                        : "text-text-muted hover:bg-surface-overlay hover:text-text-secondary"
                    }`}
                  >
                    <span className={`${t.icon} text-[11px]`} />
                    {t.label}
                  </button>
                ))}
              </div>

              {/* Content */}
              <div className="flex-1 overflow-y-auto px-6 py-5">
                {loading ? (
                  <div className="flex items-center justify-center py-12">
                    <div className="h-5 w-5 animate-spin rounded-full border-2 border-accent border-t-transparent" />
                  </div>
                ) : (
                  <>
                    {tab === "general" && (
                      <GeneralTab settings={settings} onChange={update} isEdit={isEdit} />
                    )}
                    {tab === "access" && <AccessTab settings={settings} onChange={update} />}
                    {tab === "features" && <FeaturesTab settings={settings} onChange={update} />}
                    {tab === "quota" && <QuotaTab settings={settings} onChange={update} />}
                    {tab === "encryption" && (
                      <EncryptionTab settings={settings} onChange={update} />
                    )}
                    {tab === "tags" && <TagsTab settings={settings} onChange={update} />}
                    {tab === "region" && <RegionTab settings={settings} onChange={update} />}
                  </>
                )}
              </div>

              {/* Footer */}
              <div className="flex items-center justify-between border-border border-t bg-surface/20 px-6 py-4">
                <div className="min-w-0 flex-1">
                  <AnimatePresence mode="wait">
                    {error && (
                      <m.div
                        initial={{ opacity: 0, x: -4 }}
                        animate={{ opacity: 1, x: 0 }}
                        exit={{ opacity: 0 }}
                        className="flex items-center gap-2"
                      >
                        <span className="icon-[lucide--alert-circle] shrink-0 text-danger text-xs" />
                        <p className="truncate font-body text-danger text-xs">{error}</p>
                      </m.div>
                    )}
                  </AnimatePresence>
                </div>
                <div className="flex items-center gap-2">
                  <button
                    type="button"
                    onClick={onClose}
                    className="rounded-lg border border-border px-4 py-2 font-body text-sm text-text-muted transition-colors hover:bg-surface-overlay hover:text-text-secondary"
                  >
                    Cancel
                  </button>
                  <button
                    type="button"
                    onClick={handleSubmit}
                    disabled={saving || loading}
                    className="rounded-lg bg-accent px-5 py-2 font-body font-medium text-black text-sm transition-all hover:bg-accent-bright disabled:opacity-50"
                  >
                    {saving ? (
                      <span className="flex items-center gap-2">
                        <span className="h-3.5 w-3.5 animate-spin rounded-full border-2 border-black/30 border-t-black" />
                        {isEdit ? "Saving..." : "Creating..."}
                      </span>
                    ) : isEdit ? (
                      "Save Changes"
                    ) : (
                      "Create Bucket"
                    )}
                  </button>
                </div>
              </div>
            </m.div>
          </m.div>
        </LazyMotion>
      )}
    </AnimatePresence>
  );
}
