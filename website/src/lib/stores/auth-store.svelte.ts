import type { Session } from "@supabase/supabase-js";

export const authStore = $state<{
  session: Session | null;
  isReady: boolean;
}>({
  session: null,
  isReady: false,
});
