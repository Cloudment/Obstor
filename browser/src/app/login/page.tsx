import { LoginForm } from "./LoginForm";

export default function LoginPage() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-void px-6">
      <div className="w-full max-w-sm">
        <div className="mb-8 text-center">
          <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-accent">
            <span className="icon-[lucide--database] text-2xl text-black" />
          </div>
          <h1 className="font-display text-2xl font-semibold text-text-primary">Obstor</h1>
          <p className="mt-1 font-body text-sm text-text-muted">Sign in to your storage console</p>
        </div>

        <LoginForm />

        <p className="mt-8 text-center font-mono text-[10px] text-text-muted">
          Obstor Console &middot; Apache 2.0 Licensed
        </p>
      </div>
    </div>
  );
}
