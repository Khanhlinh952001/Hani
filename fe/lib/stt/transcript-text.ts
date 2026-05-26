import type { TokenGroup } from "@soniox/react";

type StreamToken = {
  text: string;
  translation_status?: "none" | "original" | "translation";
};

export function sanitizeTranscript(text: string): string {
  return text.replace(/<fin>/gi, "").trim();
}

type RecordingTextSnapshot = {
  groups: Readonly<Record<string, TokenGroup>>;
  text: string;
  finalText: string;
  partialText: string;
  finalTokens: readonly StreamToken[];
  partialTokens: readonly StreamToken[];
};

function textFromTokens(
  tokens: readonly StreamToken[],
  kind: "ko" | "vi"
): string {
  let out = "";
  for (const t of tokens) {
    if (!t.text) continue;
    const status = t.translation_status ?? "none";
    if (kind === "vi") {
      if (status !== "translation") continue;
    } else if (status === "translation") {
      continue;
    }
    out += t.text;
  }
  return out.trim();
}

function textFromSnapshot(
  snapshot: RecordingTextSnapshot,
  kind: "ko" | "vi"
): string {
  const merged = [...snapshot.finalTokens, ...snapshot.partialTokens];
  const fromTokens = textFromTokens(merged, kind);
  if (fromTokens) return fromTokens;

  if (kind === "vi") {
    return (snapshot.groups.translation?.text ?? "").trim();
  }
  return (
    snapshot.groups.original?.text ??
    snapshot.text ??
    `${snapshot.finalText}${snapshot.partialText}`
  ).trim();
}

export function transcriptFromRecording(
  snapshot: RecordingTextSnapshot
): { ko: string; vi: string } {
  return {
    ko: textFromSnapshot(snapshot, "ko"),
    vi: textFromSnapshot(snapshot, "vi"),
  };
}

export function partialFromRecording(
  snapshot: RecordingTextSnapshot
): { ko: string; vi: string } {
  return transcriptFromRecording(snapshot);
}
