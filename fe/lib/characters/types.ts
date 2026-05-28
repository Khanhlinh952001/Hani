export type PublicCharacter = {
  id: string;
  name: string;
  display_name: string;
  gender: "female" | "male" | string;
  avatar_url: string;
  intro_message_ko: string;
  intro_message_vi: string;
  speaking_style: string;
  emotion_style: string;
  emoji_style?: string;
  typing_style?: string;
  voice_provider: string;
  voice_id: string;
};
