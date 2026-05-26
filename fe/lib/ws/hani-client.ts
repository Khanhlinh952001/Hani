import { WS_URL } from "../config";
import { ClientEvents, ServerMessage } from "./events";

export type HaniClientCallbacks = {
  onMessage: (msg: ServerMessage) => void;
  onOpen?: () => void;
  onClose?: (ev: CloseEvent) => void;
  onError?: (err: Event) => void;
};

export class HaniWsClient {
  private ws: WebSocket | null = null;
  private pingTimer: ReturnType<typeof setInterval> | null = null;

  connect(
    token: string,
    sessionId: string | undefined,
    callbacks: HaniClientCallbacks,
    prefs?: {
      ttsProvider?: "openai" | "soniox";
      ttsVoice?: string;
      ttsLanguage?: string;
      showVietnamese?: boolean;
      practiceMode?: "speak" | "chat";
    }
  ) {
    this.disconnect();

    const params = new URLSearchParams({ token });
    if (sessionId) params.set("session_id", sessionId);
    if (prefs?.ttsProvider) params.set("tts_provider", prefs.ttsProvider);
    if (prefs?.ttsVoice) params.set("tts_voice", prefs.ttsVoice);
    if (prefs?.ttsLanguage) params.set("tts_language", prefs.ttsLanguage);
    if (prefs?.showVietnamese === false) {
      params.set("show_vietnamese", "0");
    }
    if (prefs?.practiceMode === "chat") {
      params.set("practice_mode", "chat");
    }

    const url = `${WS_URL}/api/ws/chat?${params}`;
    const ws = new WebSocket(url);
    ws.binaryType = "arraybuffer";
    this.ws = ws;

    ws.onopen = () => {
      this.pingTimer = setInterval(() => {
        this.sendJson({ type: ClientEvents.Ping });
      }, 30000);
      callbacks.onOpen?.();
    };

    ws.onmessage = (ev) => {
      if (typeof ev.data !== "string") return;
      try {
        const msg = JSON.parse(ev.data) as ServerMessage;
        callbacks.onMessage(msg);
      } catch {
        /* ignore */
      }
    };

    ws.onerror = (e) => callbacks.onError?.(e);
    ws.onclose = (ev) => {
      this.clearPing();
      callbacks.onClose?.(ev);
    };
  }

  startListening() {
    this.sendJson({ type: ClientEvents.StartListening });
  }

  stopSpeaking(transcript: string, translation?: string) {
    this.sendJson({
      type: ClientEvents.StopSpeaking,
      text: transcript,
      ...(translation ? { translation } : {}),
    });
  }

  endSession() {
    this.sendJson({ type: ClientEvents.SessionEnd });
  }

  private sendJson(payload: object) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(payload));
    }
  }

  disconnect() {
    this.clearPing();
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  private clearPing() {
    if (this.pingTimer) {
      clearInterval(this.pingTimer);
      this.pingTimer = null;
    }
  }

  get isOpen() {
    return this.ws?.readyState === WebSocket.OPEN;
  }
}
