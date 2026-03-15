"use client";

import { useActionState } from "react";
import { loginAction } from "@/lib/actions";

export function LoginForm() {
  const [state, formAction, pending] = useActionState(
    async (_prev: { error?: string } | null, formData: FormData) => {
      return await loginAction(formData);
    },
    null,
  );

  return (
    <form action={formAction} className="space-y-4">
      {state?.error && (
        <div className="flex items-center gap-2 rounded-lg border border-danger/20 bg-danger/5 px-4 py-3">
          <span className="icon-[lucide--alert-circle] shrink-0 text-danger text-sm" />
          <span className="font-body text-danger text-sm">{state.error}</span>
        </div>
      )}

      <div>
        <label
          htmlFor="accessKey"
          className="mb-1.5 block font-body font-medium text-text-secondary text-xs"
        >
          Access Key
        </label>
        <input
          id="accessKey"
          name="accessKey"
          type="text"
          required
          autoFocus
          className="w-full rounded-lg border border-border bg-surface px-4 py-2.5 font-mono text-sm text-text-primary outline-none transition-colors placeholder:text-text-muted focus:border-accent"
        />
      </div>

      <div>
        <label
          htmlFor="secretKey"
          className="mb-1.5 block font-body font-medium text-text-secondary text-xs"
        >
          Secret Key
        </label>
        <input
          id="secretKey"
          name="secretKey"
          type="password"
          required
          className="w-full rounded-lg border border-border bg-surface px-4 py-2.5 font-mono text-sm text-text-primary outline-none transition-colors placeholder:text-text-muted focus:border-accent"
        />
      </div>

      <button
        type="submit"
        disabled={pending}
        className="w-full rounded-lg bg-accent py-2.5 font-body font-medium text-black text-sm transition-all hover:bg-accent-bright disabled:opacity-50"
      >
        {pending ? "Signing in..." : "Sign In"}
      </button>
    </form>
  );
}
