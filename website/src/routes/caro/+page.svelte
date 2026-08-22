<script lang="ts">
  import { getMe } from "$lib/api/api";
  import { authStore } from "$lib/stores/auth-store";
  import { createQuery } from "@tanstack/svelte-query";

  const meQuery = createQuery(() => ({
    queryKey: ["me"],
    queryFn: () => getMe($authStore.session?.access_token ?? ""),
    enabled: $authStore.isReady,
  }));
</script>

{#if meQuery.isPending}
  <div>Loading...</div>
{:else}
  <h1>Welcome to Caro Lobby {meQuery.data?.name ?? ""}!</h1>
{/if}
