"use client";

import { AuthProvider } from "@/hooks/useAuth";
import { SettingsProvider } from "@/hooks/useSettings";
import { PwaRegistrar } from "@/components/pwa/PwaRegistrar";

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <AuthProvider>
      <SettingsProvider>
        <PwaRegistrar />
        {children}
      </SettingsProvider>
    </AuthProvider>
  );
}
