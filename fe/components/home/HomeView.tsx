"use client";

import { CompanionLayout } from "@/components/layout/CompanionLayout";
import { useSettings } from "@/hooks/useSettings";
import { useCompanionHome } from "@/hooks/useCompanionHome";
import { HomeHeader } from "./HomeHeader";
import { GreetingBubble } from "./GreetingBubble";
import { VoicePracticeCard } from "./VoicePracticeCard";
import { ChatPracticeCard } from "./ChatPracticeCard";
import { EmotionalStrip } from "./EmotionalStrip";

export function HomeView() {
  const { showVietnamese } = useSettings();
  const home = useCompanionHome();

  return (
    <CompanionLayout>
      <HomeHeader
        subtitleKo={home.subtitle.ko}
        subtitleVi={showVietnamese ? home.subtitle.vi : undefined}
        moodEmoji={home.mood.emoji}
        moodKo={home.mood.ko}
      />

      <main className="home-main">
        <GreetingBubble showVietnamese={showVietnamese} />

        <div className="home-actions">
          <VoicePracticeCard voiceMinutes={home.voiceMinutes} />
          <ChatPracticeCard
            preview={home.chatPreview}
            unread={home.unread}
            onNavigate={home.clearUnread}
          />
        </div>

        <EmotionalStrip
          streak={home.streak}
          dailyLabel={home.daily.ko}
          dailyPhrase={home.daily.phrase}
          dailyVi={home.daily.vi}
          proactiveKo={home.proactive.ko}
          proactiveVi={showVietnamese ? home.proactive.vi : undefined}
        />
      </main>
    </CompanionLayout>
  );
}
