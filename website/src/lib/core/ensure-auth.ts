import { authStore } from "$lib/stores/auth-store.svelte";
import { auth } from "./auth";

export async function ensureAuth() {
  const { data } = await auth.getSession();
  if (!data.session) {
    const { data } = await auth.signInAnonymously();
    authStore.session = data.session;
    authStore.isReady = true;
  }
}