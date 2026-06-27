import type {
  ButtonHTMLAttributes,
  InputHTMLAttributes,
  ReactNode,
} from "react";

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: "primary" | "secondary";
};

export function Button({
  variant = "primary",
  className = "",
  ...props
}: ButtonProps) {
  const base =
    "inline-flex items-center justify-center px-4 py-2 text-sm font-medium rounded-[2px] transition-colors disabled:opacity-40 disabled:cursor-not-allowed";
  const styles =
    variant === "primary"
      ? "bg-accent text-white hover:bg-accent-ink"
      : "bg-transparent text-ink border border-ink hover:bg-section";
  return <button className={`${base} ${styles} ${className}`} {...props} />;
}

type FieldProps = InputHTMLAttributes<HTMLInputElement> & {
  label: string;
};

export function Field({ label, className = "", id, ...props }: FieldProps) {
  const inputId = id ?? props.name ?? label;
  return (
    <label className="block" htmlFor={inputId}>
      <span className="block text-xs uppercase tracking-wide text-muted mb-1">
        {label}
      </span>
      <input
        id={inputId}
        className={`w-full rounded-[2px] border border-line bg-surface px-3 py-2 text-sm text-ink placeholder:text-muted ${className}`}
        {...props}
      />
    </label>
  );
}

export function Card({
  children,
  className = "",
}: {
  children: ReactNode;
  className?: string;
}) {
  return (
    <div
      className={`rounded-[2px] border border-line bg-surface ${className}`}
    >
      {children}
    </div>
  );
}

export function EmptyState({
  title,
  children,
}: {
  title: string;
  children?: ReactNode;
}) {
  return (
    <div className="border border-dashed border-line rounded-[2px] bg-surface px-6 py-16 text-center">
      <p className="text-lg font-display">{title}</p>
      {children ? <div className="mt-2 text-sm text-muted">{children}</div> : null}
    </div>
  );
}

export function ErrorState({ message }: { message: string }) {
  return (
    <div className="border border-line rounded-[2px] bg-surface px-6 py-12 text-center">
      <p className="text-danger font-display text-lg">Something went wrong</p>
      <p className="mt-2 text-sm text-muted">{message}</p>
    </div>
  );
}
