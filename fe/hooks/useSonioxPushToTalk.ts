"use client";

import { useRecording } from "@soniox/react";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { getMediaSupportIssue } from "@/lib/stt/media-support";
import {
  partialFromRecording,
  transcriptFromRecording,
} from "@/lib/stt/transcript-text";

type Options = {
  enabled: boolean;
  getApiKey: () => Promise<string>;
  getContextText: () => string;
  translateToVi: boolean;
  canStart: boolean;
  onBegin: () => void;
  onComplete: (transcript: string, translation?: string) => void;
  onAbort: () => void;
  onError: (message: string) => void;
};

const UNSUPPORTED_MESSAGES: Record<string, string> = {
  ssr: "Micro chỉ dùng được trên trình duyệt.",
  "no-mediadevices": "Trình duyệt không hỗ trợ micro.",
  "no-getusermedia": "Trình duyệt không hỗ trợ ghi âm.",
  "insecure-context": "Cần HTTPS hoặc localhost để dùng micro.",
};

export function useSonioxPushToTalk({
  enabled,
  getApiKey,
  getContextText,
  translateToVi,
  canStart,
  onBegin,
  onComplete,
  onAbort,
  onError,
}: Options) {
  const [holding, setHolding] = useState(false);
  const [finishing, setFinishing] = useState(false);
  const finishRequestedRef = useRef(false);
  const releasePendingRef = useRef(false);
  const callbacksRef = useRef({ onBegin, onComplete, onAbort, onError });
  callbacksRef.current = { onBegin, onComplete, onAbort, onError };

  const contextText = enabled ? getContextText() : "";

  const recordingOptions = useMemo(
    () => ({
      apiKey: getApiKey,
      model: "stt-rt-v4" as const,
      language_hints: ["ko"],
      enable_endpoint_detection: false,
      ...(contextText ? { context: { text: contextText } } : {}),
      ...(translateToVi
        ? {
            translation: { type: "one_way" as const, target_language: "vi" },
            groupBy: "translation" as const,
          }
        : {}),
    }),
    [getApiKey, contextText, translateToVi]
  );

  const recording = useRecording({
    ...recordingOptions,
    resetOnStart: true,
    onFinished: () => {
      if (!finishRequestedRef.current) return;
      finishRequestedRef.current = false;
      setFinishing(false);

      const { ko, vi } = transcriptFromRecording(recordingRef.current);
      if (!ko) {
        callbacksRef.current.onError(
          "Không nghe thấy giọng — giữ nút, nói rồi thả để gửi"
        );
        callbacksRef.current.onAbort();
        return;
      }
      callbacksRef.current.onComplete(ko, vi || undefined);
    },
    onError: (err) => {
      if (!finishRequestedRef.current && !holdingRef.current) return;
      finishRequestedRef.current = false;
      releasePendingRef.current = false;
      setFinishing(false);
      setHolding(false);
      callbacksRef.current.onError(err.message);
      callbacksRef.current.onAbort();
    },
  });

  const recordingRef = useRef(recording);
  recordingRef.current = recording;
  const holdingRef = useRef(holding);
  holdingRef.current = holding;

  const { ko: partial, vi: partialVi } = partialFromRecording({
    groups: recording.groups,
    text: recording.text,
    finalText: recording.finalText,
    partialText: recording.partialText,
    finalTokens: recording.finalTokens,
    partialTokens: recording.partialTokens,
  });

  useEffect(() => {
    if (releasePendingRef.current && recording.isRecording) {
      releasePendingRef.current = false;
      finishRequestedRef.current = true;
      setFinishing(true);
      void recording.stop();
    }
  }, [recording.isRecording, recording.stop]);

  const pressStart = useCallback(() => {
    if (
      !enabled ||
      !canStart ||
      holdingRef.current ||
      finishing ||
      recording.isActive
    ) {
      return;
    }

    const mediaIssue = getMediaSupportIssue();
    if (mediaIssue) {
      onError(mediaIssue);
      return;
    }
    if (!recording.isSupported) {
      const reason = recording.unsupportedReason ?? "ssr";
      onError(UNSUPPORTED_MESSAGES[reason] ?? "Micro không khả dụng");
      return;
    }

    releasePendingRef.current = false;
    finishRequestedRef.current = false;
    setHolding(true);
    callbacksRef.current.onBegin();
    recording.clearTranscript();
    recording.start();
  }, [enabled, canStart, finishing, recording, onError]);

  const pressEnd = useCallback(async () => {
    if (!holdingRef.current && !recording.isActive) return;

    setHolding(false);

    if (recording.isRecording) {
      finishRequestedRef.current = true;
      setFinishing(true);
      await recording.stop();
      return;
    }

    if (recording.state === "connecting" || recording.state === "starting") {
      releasePendingRef.current = true;
      return;
    }

    if (recording.isActive) {
      finishRequestedRef.current = false;
      recording.cancel();
      setFinishing(false);
      callbacksRef.current.onAbort();
    }
  }, [recording]);

  const cancel = useCallback(() => {
    finishRequestedRef.current = false;
    releasePendingRef.current = false;
    setHolding(false);
    setFinishing(false);
    recording.cancel();
    recording.clearTranscript();
  }, [recording]);

  const resetPartial = useCallback(() => {
    recording.clearTranscript();
  }, [recording]);

  const busy =
    finishing ||
    recording.state === "connecting" ||
    recording.state === "starting" ||
    recording.state === "stopping";

  return {
    holding,
    busy,
    partial,
    partialVi,
    isSourceMuted: recording.isSourceMuted,
    pressStart,
    pressEnd,
    cancel,
    resetPartial,
  };
}
