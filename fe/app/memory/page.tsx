import { RequireAuth } from "@/components/RequireAuth";
import { MemoryView } from "@/components/memory/MemoryView";

export default function MemoryPage() {
  return (
    <RequireAuth>
      <MemoryView />
    </RequireAuth>
  );
}
