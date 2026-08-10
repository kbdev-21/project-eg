import type { Session } from "@supabase/supabase-js";
import { writable } from "svelte/store";

type AuthStatus = {
  session: Session | null;
  isReady: boolean;
}

export const authStatus = writable<AuthStatus>({
  session: null,
  isReady: false
});