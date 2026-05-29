"use client";

import { AuthProvider } from "@/hooks/useAuth";
import { SettingsProvider } from "@/hooks/useSettings";
import { FirebaseInit } from "@/components/firebase/FirebaseInit";
import { PwaRegistrar } from "@/components/pwa/PwaRegistrar";
import { PushRegistrar } from "@/components/pwa/PushRegistrar";

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <AuthProvider>
      <SettingsProvider>
        <FirebaseInit />
        <PwaRegistrar />
        <PushRegistrar />
        {children}
      </SettingsProvider>
    </AuthProvider>
  );
}
