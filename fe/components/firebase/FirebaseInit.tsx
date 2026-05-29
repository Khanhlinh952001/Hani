"use client";

import { useEffect } from "react";
import { getFirebaseAnalytics, getFirebaseApp } from "@/lib/firebase/app";

/** Mount once — initializes Firebase app + Analytics (client-side). */
export function FirebaseInit() {
  useEffect(() => {
    getFirebaseApp();
    void getFirebaseAnalytics();
  }, []);

  return null;
}
