import { RequireAuth } from "@/components/RequireAuth";
import { CharacterChooser } from "@/components/characters/CharacterChooser";

export default function ChooseHaniPage() {
  return (
    <RequireAuth>
      <CharacterChooser />
    </RequireAuth>
  );
}
