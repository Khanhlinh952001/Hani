"use client";

import { AuthProvider } from "@/hooks/useAuth";
import { SettingsProvider } from "@/hooks/useSettings";
import { ToastProvider } from "@/hooks/useToast";
import { FirebaseInit } from "@/components/firebase/FirebaseInit";
import { PwaRegistrar } from "@/components/pwa/PwaRegistrar";
import { PushRegistrar } from "@/components/pwa/PushRegistrar";

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <AuthProvider>
      <SettingsProvider>
        <ToastProvider>
          <FirebaseInit />
          <PwaRegistrar />
          <PushRegistrar />
          {children}
        </ToastProvider>
      </SettingsProvider>
    </AuthProvider>
  );
}
