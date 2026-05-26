/** Plays TTS as a queue of MP3 segments (one per sentence from the server). */
export class TtsPlayer {
  private segmentChunks: Uint8Array[] = [];
  private queue: Blob[] = [];
  private playing = false;
  private streamDone = false;
  private onFinish: (() => void) | null = null;
  private currentAudio: HTMLAudioElement | null = null;
  private unlocked = false;

  reset() {
    this.stopCurrent();
    this.segmentChunks = [];
    this.queue = [];
    this.playing = false;
    this.streamDone = false;
    this.onFinish = null;
  }

  async unlock(): Promise<void> {
    if (this.unlocked) return;

    const Ctx =
      typeof window !== "undefined"
        ? window.AudioContext ||
          (window as Window & { webkitAudioContext?: typeof AudioContext })
            .webkitAudioContext
        : undefined;

    if (Ctx) {
      const ctx = new Ctx();
      await ctx.resume();
      await ctx.close();
    }

    try {
      const silent = new Audio();
      silent.muted = true;
      silent.src =
        "data:audio/wav;base64,UklGRiQAAABXQVZFZm10IBAAAAABAAEAQB8AAEAfAAABAAgAZGF0YQAAAAA=";
      await silent.play();
      silent.pause();
    } catch {
      /* best effort */
    }

    this.unlocked = true;
  }

  startStream(onEnded?: () => void) {
    this.reset();
    this.onFinish = onEnded ?? null;
  }

  appendBase64(b64: string) {
    const binary = atob(b64);
    const bytes = new Uint8Array(binary.length);
    for (let i = 0; i < binary.length; i++) {
      bytes[i] = binary.charCodeAt(i);
    }
    this.segmentChunks.push(bytes);
  }

  /** One sentence / TTS call finished — enqueue and play in order. */
  endSegment() {
    if (this.segmentChunks.length > 0) {
      this.queue.push(
        new Blob(this.segmentChunks as BlobPart[], { type: "audio/mpeg" })
      );
      this.segmentChunks = [];
    }
    void this.drainQueue();
  }

  /** Whole assistant turn finished — no more segments after queue drains. */
  endStream() {
    this.endSegment();
    this.streamDone = true;
    void this.drainQueue();
  }

  private async drainQueue() {
    if (this.playing) return;

    this.playing = true;
    while (this.queue.length > 0) {
      const blob = this.queue.shift()!;
      const url = URL.createObjectURL(blob);
      try {
        await this.playUrl(url);
      } finally {
        URL.revokeObjectURL(url);
      }
    }
    this.playing = false;

    if (this.streamDone) {
      this.finish();
    }
  }

  private playUrl(url: string): Promise<void> {
    return new Promise((resolve) => {
      const audio = new Audio(url);
      this.currentAudio = audio;
      const done = () => {
        if (this.currentAudio === audio) {
          this.currentAudio = null;
        }
        resolve();
      };
      audio.onended = done;
      audio.onerror = done;
      void audio.play().catch(done);
    });
  }

  private stopCurrent() {
    if (this.currentAudio) {
      this.currentAudio.pause();
      this.currentAudio.src = "";
      this.currentAudio = null;
    }
  }

  private finish() {
    this.stopCurrent();
    this.onFinish?.();
    this.onFinish = null;
  }

  /** Legacy: play all buffered chunks at once. */
  async play(): Promise<void> {
    return new Promise((resolve) => {
      this.startStream(resolve);
      this.endStream();
    });
  }
}
