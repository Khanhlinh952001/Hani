import { RequireAuth } from "@/components/RequireAuth";
import { ChatView } from "@/components/chat/ChatView";

export default function ChatPracticePage() {
  return (
    <RequireAuth>
      <ChatView practiceMode="chat" />
    </RequireAuth>
  );
}
