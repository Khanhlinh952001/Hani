import { RequireAuth } from "@/components/RequireAuth";
import { CreateLoverWizard } from "@/components/lover/CreateLoverWizard";

export default function CreateLoverPage() {
  return (
    <RequireAuth>
      <CreateLoverWizard />
    </RequireAuth>
  );
}
