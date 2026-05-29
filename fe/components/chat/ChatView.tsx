"use client";

import { useRef } from "react";
import type { PracticeMode } from "@/lib/practice/mode";
import { useHaniChat } from "@/hooks/useHaniChat";
import { useSettings } from "@/hooks/useSettings";
import { CompanionLayout } from "@/components/layout/CompanionLayout";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { ConnectionBar } from "./ConnectionBar";
import { MessageList } from "./MessageList";
import { PushToTalkButton } from "./PushToTalkButton";
import { TextComposer } from "./TextComposer";

type Props = {
  practiceMode: PracticeMode;
};

export function ChatView({ practiceMode }: Props) {
  const { showVietnamese } = useSettings();
  const chat = useHaniChat(practiceMode);
  const scrollRef = useRef<HTMLDivElement>(null);
  const isSpeakMode = practiceMode === "speak";

  const footer = isSpeakMode ? (
    <PushToTalkButton
      status={chat.status}
      isConnected={chat.isConnected}
      canPress={chat.canPress}
      holding={chat.holding}
      busy={chat.busy}
      showVietnamese={showVietnamese}
      onPressStart={chat.pressStart}
      onPressEnd={chat.pressEnd}
      isSourceMuted={chat.isSourceMuted}
    />
  ) : (
    <TextComposer
      status={chat.status}
      isConnected={chat.isConnected}
      showVietnamese={showVietnamese}
      onSend={chat.sendText}
      chatOnly
    />
  );

  return (
    <CompanionLayout hideNav footer={footer} className="flex flex-col">
      <ConnectionBar
        practiceMode={practiceMode}
        status={chat.status}
        user={chat.user}
        isConnected={chat.isConnected}
      />

      {chat.error && (
        <div className="px-3 pt-2">
          <Alert variant="destructive" className="animate-in fade-in slide-in-from-top-1">
            <AlertDescription>{chat.error}</AlertDescription>
          </Alert>
        </div>
      )}

      <div className="message-scroll-wrap min-h-0 flex-1">
        <div
          ref={scrollRef}
          className="message-scroll h-full overflow-y-auto overscroll-contain px-3"
        >
          <MessageList
            user={chat.user}
            messages={chat.messages}
            partial={chat.partial}
            partialVi={showVietnamese ? chat.partialVi : undefined}
            status={chat.status}
            holding={chat.holding}
            showVietnamese={showVietnamese}
            practiceMode={practiceMode}
            hasMoreHistory={chat.hasMoreHistory}
            loadingMoreHistory={chat.loadingMoreHistory}
            onLoadOlder={() => void chat.loadOlderMessages()}
            scrollRef={scrollRef}
          />
        </div>
      </div>
    </CompanionLayout>
  );
}
