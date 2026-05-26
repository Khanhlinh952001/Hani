# Hani Realtime WebSocket Architecture

## Folder structure

```
cmd/api/main.go
internal/
  config/          env helpers
  db/              global GORM connection
  websocket/       hub, events, session pipeline, handler
  stt/             (legacy) Soniox helpers — STT runs in browser
  tts/             OpenAI or Soniox TTS (TTS_PROVIDER)
  ai/              Bilingual replies: Korean + «VI» + Vietnamese (no separate translate API)
  ai/              OpenAI chat stream + embeddings + Hani persona
  memory/          pgvector retrieval for prompts
  conversation/    sessions + messages persistence
  modules/         REST CRUD (users, sessions, messages, memories)
```

## Connect

```
ws://localhost:8080/api/ws/chat?user_id=1
ws://localhost:8080/api/ws/chat?user_id=1&session_id=<uuid>   # reconnect
```

## Client → Server

| Payload | Description |
|---------|-------------|
| `{"type":"start_listening"}` | Request STT context for browser Soniox |
| `{"type":"stop_speaking","text":"..."}` | Final transcript from client STT → AI pipeline |
| `{"type":"session_end"}` | Close conversation session |
| `{"type":"ping"}` | Keepalive |

## Server → Client

| Event | Description |
|-------|-------------|
| `ready` | Connected; `session_id`, `stt_context` |
| `listening` | Ack + refreshed `stt_context` |
| `final_transcript` | User utterance saved (echo) |
| `typing_start` / `typing_end` | Hani thinking indicator |
| `ai_response` | Streaming text (`delta`) |
| `subtitle` | Full Korean reply line |
| `ai_audio_chunk` | Base64 MP3 segment (per sentence stream) |
| `ai_audio_segment_end` | One sentence TTS finished — FE enqueues next segment |
| `ai_audio_end` | Full reply TTS complete |
| `error` | Failure |

## Pipeline (low latency)

```
browser mic → Soniox STT → stop_speaking(text)
  → save user message
  → embed query → pgvector memory search (top 5)
  → OpenAI stream (persona + memories + last 8 turns only)
  → save assistant message
  → OpenAI TTS stream → ai_audio_chunk
```

## Env

```
OPENAI_API_KEY=
OPENAI_MODEL=gpt-4o-mini
OPENAI_EMBEDDING_MODEL=text-embedding-3-small

# TTS — openai (default) | soniox
TTS_PROVIDER=openai
SONIOX_API_KEY=
SONIOX_TTS_MODEL=tts-rt-v1
SONIOX_TTS_LANGUAGE=ko
SONIOX_TTS_VOICE=Kenji
SONIOX_TTS_AUDIO_FORMAT=mp3

# OpenAI TTS (only if TTS_PROVIDER=openai)
OPENAI_TTS_MODEL=tts-1
OPENAI_TTS_VOICE=nova
```

## Design notes

- **No full chat history** in prompt — only recent turns + retrieved memories.
- **Memories** are not created on every utterance (extract separately when worth remembering).
- **Hub** tracks concurrent connections; one process scales horizontally with sticky sessions later.
- REST modules remain for admin/debug; realtime path is optimized for conversation.
