import { RequireAuth } from "@/components/RequireAuth";
import { HomeView } from "@/components/home/HomeView";

export default function Home() {
  return (
    <RequireAuth>
      <HomeView />
    </RequireAuth>
  );
}
