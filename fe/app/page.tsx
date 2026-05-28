import { RequireAuth } from "@/components/RequireAuth";
import { RequireCharacter } from "@/components/RequireCharacter";
import { HomeView } from "@/components/home/HomeView";

export default function Home() {
  return (
    <RequireAuth>
      <RequireCharacter>
        <HomeView />
      </RequireCharacter>
    </RequireAuth>
  );
}
