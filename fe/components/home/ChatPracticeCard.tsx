"use client";

import Link from "next/link";
import { MessageCircle } from "lucide-react";

type Props = {
  preview: string;
  unread: boolean;
  onNavigate?: () => void;
};

export function ChatPracticeCard({ preview, unread, onNavigate }: Props) {
  return (
    <Link href="/chat" className="home-action-card home-action-chat" onClick={onNavigate}>
      <div className="home-action-top">
        <div className="home-action-icon home-action-icon-chat">
          <MessageCircle className="size-5" strokeWidth={2.25} />
        </div>
        {unread ? (
          <span className="home-unread-badge" aria-label="Tin nhắn mới">
            1
          </span>
        ) : null}
      </div>

      <div className="home-action-body">
        <p className="home-action-ko">채팅 연습</p>
        <p className="home-action-title">Chat Practice</p>
        <p className="home-action-desc">Nhắn tin như KakaoTalk</p>
      </div>

      <div className="home-preview-bubble">
        <span className="home-preview-from">Hani</span>
        <p className="home-preview-text">{preview}</p>
      </div>
    </Link>
  );
}
