<script lang="ts">
  import { getMe } from "$lib/api/api";
  import ActionCard from "$lib/components/ActionCard.svelte";
  import PlayerCard from "$lib/components/PlayerCard.svelte";
  import { Hash, Search, Bot, UserPlus, Swords } from "@lucide/svelte";
  import { authStore } from "$lib/stores/auth-store.svelte";
  import { createQuery } from "@tanstack/svelte-query";

  const meQuery = createQuery(() => ({
    queryKey: ["me"],
    queryFn: () => getMe(authStore.session?.access_token ?? ""),
    enabled: authStore.isReady,
  }));
</script>

<div class="min-h-screen">
  <main class="mx-auto w-full max-w-[460px] px-4 pt-10">
    <!-- Title -->
    <div class="mb-8 flex items-center justify-center gap-3">
      <div
        class="flex h-12 w-12 items-center justify-center rounded-2xl bg-blue-700 text-white shadow-sm"
      >
        <Hash class="h-6 w-6" />
      </div>
      <h1 class="text-3xl font-extrabold tracking-tight">Cờ Caro</h1>
    </div>

    <!-- Profile card -->
    <PlayerCard player={meQuery.data} isLoading={meQuery.isPending} />

    <!-- PLAY header -->
    <div class="mt-7 mb-3 flex items-center justify-between px-1">
      <span class="text-xs font-bold tracking-widest text-neutral-500"
        >CHƠI</span
      >
      <span class="flex items-center gap-1.5 text-xs text-neutral-400">
        <span class="h-2 w-2 rounded-full bg-green-500"></span>
        <span class="font-bold text-neutral-600">3,147</span>
        người đang online
      </span>
    </div>

    <!-- Matchmaking -->
    <div class="mb-3">
      <ActionCard
        title="Chơi Online"
        subtitle="Tìm trận và đấu với đối thủ cùng trình độ"
        iconBg="bg-amber-400 text-white"
      >
        {#snippet icon()}
          <Swords class="h-5 w-5" />
        {/snippet}
        {#snippet trailing()}
          <span
            class="rounded-full bg-amber-400 px-2.5 py-1 text-[11px] font-bold tracking-wide text-neutral-900"
            >RANKED</span
          >
        {/snippet}
      </ActionCard>
    </div>

    <!-- Practice with bot -->
    <div class="mb-3">
      <ActionCard
        title="Chơi với máy"
        subtitle="Luyện tập với máy. Không giới hạn"
        iconBg="bg-green-600 text-white"
      >
        {#snippet icon()}
          <Bot class="h-5 w-5" />
        {/snippet}
      </ActionCard>
    </div>

    <!-- Play a friend -->
    <ActionCard
      title="Chơi với bạn bè"
      subtitle="Tạo phòng riêng và thách đấu với bạn bè"
      iconBg="bg-blue-700 text-white"
    >
      {#snippet icon()}
        <UserPlus class="h-5 w-5" />
      {/snippet}
    </ActionCard>
  </main>
</div>
