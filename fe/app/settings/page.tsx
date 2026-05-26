import { RequireAuth } from "@/components/RequireAuth";
import { SettingsView } from "@/components/settings/SettingsView";

export default function SettingsPage() {
  return (
    <RequireAuth>
      <SettingsView />
    </RequireAuth>
  );
}
