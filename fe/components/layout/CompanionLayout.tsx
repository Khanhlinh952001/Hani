import { AppShell } from "./AppShell";
import { BottomNav } from "./BottomNav";

type Props = {
  children: React.ReactNode;
  hideNav?: boolean;
  footer?: React.ReactNode;
  className?: string;
};

export function CompanionLayout({
  children,
  hideNav = false,
  footer,
  className,
}: Props) {
  return (
    <AppShell
      className={className}
      footer={footer}
      bottomNav={hideNav ? undefined : <BottomNav />}
    >
      {children}
    </AppShell>
  );
}
