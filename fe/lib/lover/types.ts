export type PersonalityTemplate = {
  id: string;
  name_ko: string;
  name_vi: string;
  description_ko: string;
  description_vi: string;
  icon: string;
  emoji_density: number;
  typing_speed: string;
  flirting_level: number;
};

export type VoiceProfile = {
  id: string;
  name_ko: string;
  name_vi: string;
  gender: string;
  provider: string;
  voice_id: string;
  emotion: string;
  speed: string;
  preview_text_ko: string;
};

export type SpeakingStyleOption = {
  id: string;
  label_ko: string;
  label_vi: string;
  description: string;
};

export type LoverProfile = {
  id: string;
  display_name: string;
  companion_gender: string;
  personality_template_id: string;
  speaking_style_tags: string[];
  voice_profile_id: string;
  tts_voice?: string;
  voice_provider?: string;
  avatar_url: string;
  intro_message_ko: string;
  intro_message_vi: string;
  preset_slug?: string;
  personality_name_ko?: string;
  personality_name_vi?: string;
  voice_name_ko?: string;
};

export type CreateProfilePayload = {
  companion_gender: string;
  personality_template_id: string;
  speaking_style_tags: string[];
  voice_profile_id: string;
  display_name: string;
  avatar_url?: string;
};
