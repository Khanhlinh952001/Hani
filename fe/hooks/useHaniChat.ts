"use client";

import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useAuth } from "@/hooks/useAuth";
import { useSettings } from "@/hooks/useSettings";
import { useSonioxPushToTalk } from "@/hooks/useSonioxPushToTalk";
import { TtsPlayer } from "@/lib/audio/tts-player";
import { createSonioxKeyFetcher } from "@/lib/stt/get-soniox-key";
import { sanitizeTranscript } from "@/lib/stt/transcript-text";
import {
  ChatMessage,
  ConnectionStatus,
  ServerEvents,
  ServerMessage,
} from "@/lib/ws/events";
import type { PracticeMode } from "@/lib/practice/mode";
import { clearConversationHistory } from "@/lib/sessions/api";
import { HaniWsClient } from "@/lib/ws/hani-client";

const SESSION_KEY = "hani_session_id";

function uid() {
  return `${Date.now()}-${Math.random().toString(36).slice(2, 9)}`;
}

export function useHaniChat(practiceMode: PracticeMode) {
  const { user, token } = useAuth();
  const { showVietnamese, ttsProvider, ttsVoice, ttsLanguage } = useSettings();
  const voiceEnabled = practiceMode === "speak";
  const getSonioxApiKey = useMemo(
    () => createSonioxKeyFetcher(token ?? ""),
    [token]
  );
  const voiceEnabledRef = useRef(voiceEnabled);
  voiceEnabledRef.current = voiceEnabled;
  const clientRef = useRef<HaniWsClient | null>(null);
  const sttContextRef = useRef("");
  const [sttContext, setSttContext] = useState("");
  const statusRef = useRef<ConnectionStatus>("disconnected");
  const ttsRef = useRef(new TtsPlayer());
  const ttsStreamingRef = useRef(false);
  const assistantIdRef = useRef<string | null>(null);
  const pttCancelRef = useRef<() => void>(() => {});

  const [sessionId, setSessionId] = useState<string | null>(() =>
    typeof window !== "undefined" ? localStorage.getItem(SESSION_KEY) : null
  );
  const [status, setStatusState] = useState<ConnectionStatus>("disconnected");
  const setStatus = useCallback((s: ConnectionStatus) => {
    statusRef.current = s;
    setStatusState(s);
  }, []);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [error, setError] = useState<string | null>(null);
  const connectingRef = useRef(false);

  const finalizeAssistant = useCallback(
    (text: string, translationVi?: string) => {
      const aid = assistantIdRef.current;
      if (aid) {
        setMessages((m) =>
          m.map((msg) =>
            msg.id === aid
              ? {
                  ...msg,
                  content: text,
                  translationVi: translationVi || msg.translationVi,
                  streaming: false,
                  justArrived: !voiceEnabledRef.current,
                }
              : msg
          )
        );
      } else {
        setMessages((m) => [
          ...m,
          {
            id: uid(),
            role: "assistant",
            content: text,
            translationVi,
            justArrived: !voiceEnabledRef.current,
          },
        ]);
      }
      assistantIdRef.current = null;
    },
    []
  );

  const ptt = useSonioxPushToTalk({
    enabled: voiceEnabled,
    getApiKey: getSonioxApiKey,
    getContextText: () => sttContext,
    translateToVi: showVietnamese,
    canStart: status === "ready",
    onBegin: () => {
      setError(null);
      ttsRef.current.reset();
      assistantIdRef.current = null;
      clientRef.current?.startListening();
      setStatus("listening");
    },
    onComplete: (transcript, translation) => {
      setError(null);
      clientRef.current?.stopSpeaking(transcript, translation);
      setStatus("thinking");
      assistantIdRef.current = null;
    },
    onAbort: () => {
      if (statusRef.current === "listening") {
        setStatus("ready");
      }
    },
    onError: (message) => {
      setError(message);
      if (statusRef.current === "listening") {
        setStatus("ready");
      }
    },
  });

  pttCancelRef.current = ptt.cancel;

  const handleServerMessage = useCallback(
    (msg: ServerMessage) => {
      switch (msg.type) {
        case ServerEvents.Ready:
          setSessionId(msg.session_id ?? null);
          if (msg.session_id) {
            localStorage.setItem(SESSION_KEY, msg.session_id);
          }
          if (msg.stt_context) {
            sttContextRef.current = msg.stt_context;
            setSttContext(msg.stt_context);
          }
          if (msg.messages?.length) {
            const history = msg.messages.slice(-3);
            setMessages(
              history.map((m, i) => ({
                id: m.id ?? `hist-${i}-${msg.session_id ?? "s"}`,
                role: m.role === "assistant" ? "assistant" : "user",
                content: m.content,
                translationVi: m.translation,
              }))
            );
            setStatus("ready");
          } else {
            setStatus("thinking");
          }
          setError(null);
          break;

        case ServerEvents.Listening:
          if (msg.stt_context !== undefined) {
            sttContextRef.current = msg.stt_context;
            setSttContext(msg.stt_context);
          }
          break;

        case ServerEvents.FinalTranscript: {
          const text = sanitizeTranscript(msg.text ?? "");
          if (!text) break;
          const translationVi = msg.translation;
          ptt.resetPartial();
          if (!voiceEnabledRef.current) {
            setMessages((m) => {
              const last = m[m.length - 1];
              if (last?.role === "user" && last.content === text) {
                return m.map((item, i) =>
                  i === m.length - 1 ? { ...item, translationVi } : item
                );
              }
              return [
                ...m,
                {
                  id: uid(),
                  role: "user",
                  content: text,
                  translationVi,
                  justArrived: true,
                },
              ];
            });
          } else {
            const id = uid();
            setMessages((m) => [
              ...m,
              {
                id,
                role: "user",
                content: text,
                translationVi: msg.translation,
              },
            ]);
          }
          setStatus("thinking");
          assistantIdRef.current = null;
          break;
        }

        case ServerEvents.TypingStart:
          setStatus("thinking");
          assistantIdRef.current = null;
          ttsStreamingRef.current = false;
          ttsRef.current.reset();
          break;

        case ServerEvents.AIResponse:
          // Voice mode: text arrives in one shot via Subtitle (KO+VI together).
          if (msg.delta && !voiceEnabledRef.current) {
            if (!assistantIdRef.current) {
              assistantIdRef.current = uid();
            }
            setStatus("thinking");
          }
          break;

        case ServerEvents.Subtitle: {
          const text = msg.text ?? msg.full_text ?? "";
          finalizeAssistant(text, msg.translation);
          if (voiceEnabledRef.current) {
            setStatus("speaking");
          } else {
            setStatus("ready");
          }
          break;
        }

        case ServerEvents.TypingEnd:
          break;

        case ServerEvents.AIAudioChunk:
          if (!voiceEnabledRef.current) break;
          if (msg.audio) {
            setStatus("speaking");
            if (!ttsStreamingRef.current) {
              ttsStreamingRef.current = true;
              ttsRef.current.startStream(() => {
                ttsStreamingRef.current = false;
                ttsRef.current.reset();
                setStatus("ready");
              });
            }
            ttsRef.current.appendBase64(msg.audio);
          }
          break;

        case ServerEvents.AIAudioSegmentEnd:
          if (!voiceEnabledRef.current) break;
          if (ttsStreamingRef.current) {
            ttsRef.current.endSegment();
          }
          break;

        case ServerEvents.AIAudioEnd:
          if (!voiceEnabledRef.current) {
            setStatus("ready");
            break;
          }
          if (ttsStreamingRef.current) {
            ttsRef.current.endStream();
          } else {
            setStatus("ready");
          }
          break;

        case ServerEvents.Error:
          setError(msg.message ?? "unknown error");
          setStatus("error");
          break;

        case ServerEvents.SessionEnded:
          localStorage.removeItem(SESSION_KEY);
          setSessionId(null);
          setStatus("disconnected");
          break;
      }
    },
    [finalizeAssistant, ptt, setStatus]
  );

  const connect = useCallback(() => {
    if (!token) {
      setError("Chưa đăng nhập");
      return;
    }
    if (connectingRef.current || clientRef.current?.isOpen) {
      return;
    }

    connectingRef.current = true;
    setStatus("connecting");
    setError(null);
    if (voiceEnabled) {
      void ttsRef.current.unlock();
    } else {
      ttsRef.current.reset();
      ttsStreamingRef.current = false;
    }

    const client = new HaniWsClient();
    clientRef.current = client;
    let gotReady = false;

    client.connect(
      token,
      undefined,
      {
        onMessage: (msg) => {
          if (msg.type === "ready") {
            gotReady = true;
            connectingRef.current = false;
          }
          handleServerMessage(msg);
        },
        onClose: (ev) => {
          connectingRef.current = false;
          setStatus("disconnected");
          pttCancelRef.current();
          if (!gotReady) {
            localStorage.removeItem(SESSION_KEY);
            setSessionId(null);
            setError(
              ev.code === 1006
                ? "Không kết nối được — backend đang chạy chưa?"
                : "Phiên hết hạn — đăng nhập lại hoặc «새 대화»"
            );
            setStatus("error");
          }
        },
        onError: () => {
          connectingRef.current = false;
          if (!gotReady) {
            setError("WebSocket lỗi — kiểm tra token / backend");
            setStatus("error");
          }
        },
      },
      {
        ttsProvider,
        ttsVoice,
        ttsLanguage,
        showVietnamese,
        practiceMode,
      }
    );
  }, [
    token,
    ttsProvider,
    ttsVoice,
    ttsLanguage,
    showVietnamese,
    practiceMode,
    voiceEnabled,
    handleServerMessage,
    setStatus,
  ]);

  const disconnect = useCallback(() => {
    connectingRef.current = false;
    clientRef.current?.disconnect();
    clientRef.current = null;
    pttCancelRef.current();
    setStatus("disconnected");
  }, [setStatus]);

  const connectRef = useRef(connect);
  const disconnectRef = useRef(disconnect);
  connectRef.current = connect;
  disconnectRef.current = disconnect;

  useEffect(() => {
    ttsRef.current.reset();
    ttsStreamingRef.current = false;
  }, [practiceMode]);

  useEffect(() => {
    if (!token) return;
    connectRef.current();
    return () => disconnectRef.current();
  }, [token, practiceMode, ttsProvider, ttsVoice, ttsLanguage]);

  const sendText = useCallback(
    (text: string) => {
      const trimmed = text.trim();
      if (!trimmed || !clientRef.current?.isOpen) return;

      ptt.resetPartial();
      setError(null);
      if (!voiceEnabledRef.current) {
        setMessages((m) => [
          ...m,
          {
            id: uid(),
            role: "user",
            content: trimmed,
            justArrived: true,
          },
        ]);
      }
      clientRef.current.stopSpeaking(trimmed);
      setStatus("thinking");
      assistantIdRef.current = null;
    },
    [ptt, setStatus]
  );

  const clearHistory = useCallback(async () => {
    if (!token) return;
    if (
      !confirm(
        "Xóa toàn bộ lịch sử và ký ức? Tin nhắn và vector nhớ sẽ bị xóa hết — Hani chào lại từ đầu."
      )
    ) {
      return;
    }
    try {
      await clearConversationHistory();
      setMessages([]);
      setSessionId(null);
      localStorage.removeItem(SESSION_KEY);
      disconnect();
      window.setTimeout(() => connect(), 50);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Không xóa được lịch sử");
      setStatus("error");
    }
  }, [token, disconnect, connect, setStatus]);

  return {
    user,
    sessionId,
    status,
    messages,
    partial: ptt.partial,
    partialVi: ptt.partialVi,
    error,
    connect,
    disconnect,
    pressStart: ptt.pressStart,
    pressEnd: ptt.pressEnd,
    sendText,
    clearHistory,
    holding: ptt.holding,
    busy: ptt.busy,
    canPress: status === "ready" && !ptt.busy && !ptt.holding,
    isSourceMuted: ptt.isSourceMuted,
    isConnected:
      status === "ready" ||
      status === "listening" ||
      status === "thinking" ||
      status === "speaking",
  };
}
