export const API_URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export const WS_URL =
  process.env.NEXT_PUBLIC_WS_URL ?? "ws://localhost:8080";

export const DEFAULT_USER_ID = 1;

/** True in `next dev`; false in production build. */
export const IS_DEV = process.env.NODE_ENV === "development";
