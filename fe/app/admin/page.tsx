import { RequireAuth } from "@/components/RequireAuth";
import { AdminView } from "@/components/admin/AdminView";
import { RequireAdmin } from "@/components/admin/RequireAdmin";

export default function AdminPage() {
  return (
    <RequireAuth>
      <RequireAdmin>
        <AdminView />
      </RequireAdmin>
    </RequireAuth>
  );
}
