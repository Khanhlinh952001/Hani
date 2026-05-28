export type UsageSnapshot = {
  plan: string;
  daily_messages: number;
  daily_messages_limit?: number | null;
  daily_voice_seconds: number;
  daily_voice_limit?: number | null;
  warning?: boolean;
};

export type PlanLimit = {
  plan: string;
  daily_messages?: number | null;
  daily_voice_seconds?: number | null;
  allow_voice: boolean;
  allow_memory: boolean;
  allow_premium_voices: boolean;
};
