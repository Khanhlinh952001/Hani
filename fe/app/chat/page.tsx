import { RequireAuth } from "@/components/RequireAuth";
import { RequireCharacter } from "@/components/RequireCharacter";
import { ChatView } from "@/components/chat/ChatView";

export default function ChatPracticePage() {
  return (
    <RequireAuth>
      <RequireCharacter>
        <ChatView practiceMode="chat" />
      </RequireCharacter>
    </RequireAuth>
  );
}
