/** Pre-flight mic — pattern từ Memora transcribe-shared.ts */
export function getMediaSupportIssue(): string | null {
  if (typeof navigator === "undefined") {
    return "Micro chỉ dùng được trên trình duyệt.";
  }
  if (typeof window !== "undefined" && !window.isSecureContext) {
    return "Cần HTTPS hoặc localhost để dùng micro.";
  }
  if (!("mediaDevices" in navigator) || navigator.mediaDevices == null) {
    return "Trình duyệt không hỗ trợ micro — thử Chrome, Safari hoặc Firefox.";
  }
  if (typeof navigator.mediaDevices.getUserMedia !== "function") {
    return "Trình duyệt không hỗ trợ getUserMedia.";
  }
  return null;
}
