import { RequireAuth } from "@/components/RequireAuth";
import { RequireCharacter } from "@/components/RequireCharacter";
import { MemoryView } from "@/components/memory/MemoryView";

export default function MemoryPage() {
  return (
    <RequireAuth>
      <RequireCharacter>
        <MemoryView />
      </RequireCharacter>
    </RequireAuth>
  );
}
