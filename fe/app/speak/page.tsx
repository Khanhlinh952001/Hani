import { RequireAuth } from "@/components/RequireAuth";
import { RequireCharacter } from "@/components/RequireCharacter";
import { ChatView } from "@/components/chat/ChatView";

export default function SpeakPage() {
  return (
    <RequireAuth>
      <RequireCharacter>
        <ChatView practiceMode="speak" />
      </RequireCharacter>
    </RequireAuth>
  );
}
