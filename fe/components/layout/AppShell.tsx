import { cn } from "@/lib/utils";

type Props = {
  children: React.ReactNode;
  className?: string;
  footer?: React.ReactNode;
};

export function AppShell({ children, className, footer }: Props) {
  return (
    <div
      className={cn(
        "hani-app-frame relative z-10 mx-auto flex h-dvh max-h-dvh w-full max-w-md flex-col overflow-hidden rounded-none border-x border-primary/12 bg-card/92 shadow-2xl shadow-primary/20 backdrop-blur-md sm:my-2 sm:h-[calc(100dvh-1rem)] sm:max-h-[calc(100dvh-1rem)] sm:rounded-[1.75rem]",
        className
      )}
    >
      <div className="hani-shell-inner">{children}</div>
      {footer ? <footer className="hani-footer">{footer}</footer> : null}
    </div>
  );
}
