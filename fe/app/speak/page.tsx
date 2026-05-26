import dynamic from "next/dynamic";
import { RequireAuth } from "@/components/RequireAuth";
import { ChatView } from "@/components/chat/ChatView";

export default function SpeakPage() {
  return (
    <RequireAuth>
      <ChatView practiceMode="speak" />
    </RequireAuth>
  );
}
