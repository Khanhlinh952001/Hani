"use client";

import {
  useCallback,
  useEffect,
  useRef,
  useState,
  type ReactNode,
} from "react";
import { X } from "lucide-react";
import { cn } from "@/lib/utils";

type Props = {
  src?: string;
  alt: string;
  children: ReactNode;
  disabled?: boolean;
  className?: string;
};

export function AvatarPreview({
  src,
  alt,
  children,
  disabled,
  className,
}: Props) {
  const dialogRef = useRef<HTMLDialogElement>(null);
  const [open, setOpen] = useState(false);

  const canPreview = Boolean(src?.trim()) && !disabled;

  const close = useCallback(() => setOpen(false), []);

  useEffect(() => {
    const el = dialogRef.current;
    if (!el) return;
    if (open && !el.open) el.showModal();
    else if (!open && el.open) el.close();
  }, [open]);

  if (!canPreview) {
    return <>{children}</>;
  }

  return (
    <>
      <button
        type="button"
        onClick={() => setOpen(true)}
        className={cn(
          "shrink-0 cursor-zoom-in rounded-full border-0 bg-transparent p-0 outline-none ring-offset-background transition-transform hover:scale-[1.03] active:scale-[0.98] focus-visible:ring-2 focus-visible:ring-ring",
          className
        )}
        aria-label={`Xem ảnh ${alt}`}
      >
        {children}
      </button>

      <dialog
        ref={dialogRef}
        className="avatar-preview-dialog"
        onClose={close}
        onCancel={close}
        onClick={(e) => {
          if (e.target === e.currentTarget) close();
        }}
      >
        <div
          className="avatar-preview-panel"
          onClick={(e) => e.stopPropagation()}
        >
          <button
            type="button"
            className="avatar-preview-close"
            onClick={close}
            aria-label="Đóng"
          >
            <X className="size-5" />
          </button>
          {/* Full-size view — not using next/image (dynamic upload URLs) */}
          <img src={src} alt={alt} className="avatar-preview-img" />
          <p className="avatar-preview-caption">{alt}</p>
        </div>
      </dialog>
    </>
  );
}
