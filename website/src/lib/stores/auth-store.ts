import type { Session } from "@supabase/supabase-js";
import { writable } from "svelte/store";

type AuthState = {
  session: Session | null;
  isReady: boolean;
}

export const authStore = writable<AuthState>({
  session: null,
  isReady: false
});