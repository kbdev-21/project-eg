import { authStore } from "$lib/stores/auth-store";
import { auth } from "./auth";

export async function ensureAuth() {
  const { data } = await auth.getSession();
  if (!data.session) {
    const { data } = await auth.signInAnonymously();
    authStore.set({
      session: data.session,
      isReady: true
    });
  }
}