<script lang="ts">
  import { Circle, Send, X } from "@lucide/svelte";
  import { createQuery, useQueryClient } from "@tanstack/svelte-query";
  import { goto } from "$app/navigation";
  import { getMe, getUserById } from "$lib/api/api";
  import { AVATARS } from "$lib/core/avatars";
  import { CARO_BOARD_SIZE, CARO_MAX_MOVE_TIME_MS } from "$lib/core/ws-types";
  import { sendCaroPlayMoveMessage } from "$lib/core/websocket";
  import { authStore } from "$lib/stores/auth-store.svelte";
  import { stateStore } from "$lib/stores/state-store.svelte";

  type Mark = "X" | "O";

  const MARK_COLOR: Record<Mark, string> = {
    X: "#3B5BDB",
    O: "#d52e33",
  };

  const match = $derived(stateStore.match);
  const isEnded = $derived(!!match && match.status !== "PLAYING");

  // Vào thẳng URL khi không có trận -> về sảnh.
  // Chờ hydrated để F5 giữa trận không bị đá về sảnh trước khi PING trả state.
  $effect(() => {
    if (stateStore.hydrated && !match) goto("/caro");
  });

  const meQuery = createQuery(() => ({
    queryKey: ["me"],
    queryFn: () => getMe(authStore.session?.access_token ?? ""),
    enabled: authStore.isReady,
  }));
  const meId = $derived(meQuery.data?.id);

  const xUserQuery = createQuery(() => ({
    queryKey: ["user", match?.xPlayerId],
    queryFn: () => getUserById(match!.xPlayerId),
    enabled: !!match,
  }));
  const oUserQuery = createQuery(() => ({
    queryKey: ["user", match?.oPlayerId],
    queryFn: () => getUserById(match!.oPlayerId),
    enabled: !!match,
  }));

  const mySide: Mark = $derived(match && match.xPlayerId === meId ? "X" : "O");
  const oppSide: Mark = $derived(mySide === "X" ? "O" : "X");
  const isMyTurn = $derived(!!match && !isEnded && match.turnOf === mySide);

  const lastMove = $derived(match?.moves.at(-1));

  // --- Đồng hồ 20s/nước ---
  // Không có field deadline: suy ra từ nước cuối, fallback startedAt.
  // playedAt là giờ server nên lệch đồng hồ máy sẽ làm countdown lệch; server mới là trọng tài.
  let now = $state(Date.now());

  $effect(() => {
    if (!match || isEnded) return;
    const id = setInterval(() => (now = Date.now()), 250);
    return () => clearInterval(id);
  });

  const remainMs = $derived.by(() => {
    if (!match || isEnded) return 0;
    const deadline =
      new Date(lastMove?.playedAt ?? match.startedAt).getTime() +
      CARO_MAX_MOVE_TIME_MS;
    return Math.max(0, deadline - now);
  });

  function isActive(side: Mark) {
    return !!match && !isEnded && match.turnOf === side;
  }

  function clockOf(side: Mark) {
    if (!isActive(side)) return "0:20";
    return `0:${String(Math.ceil(remainMs / 1000)).padStart(2, "0")}`;
  }

  function progressOf(side: Mark) {
    if (!isActive(side)) return 1;
    return remainMs / CARO_MAX_MOVE_TIME_MS;
  }

  function handleCellClick(x: number, y: number) {
    if (!isMyTurn || !match) return;
    if (match.board[y][x] !== "") return;
    sendCaroPlayMoveMessage(x, y);
  }

  // --- Kết quả trận ---
  const myRatingBefore = $derived(
    mySide === "X" ? match?.xPlayerRatingBefore : match?.oPlayerRatingBefore,
  );
  const myRatingAfter = $derived(
    mySide === "X" ? match?.xPlayerRatingAfter : match?.oPlayerRatingAfter,
  );
  const ratingDelta = $derived((myRatingAfter ?? 0) - (myRatingBefore ?? 0));

  const resultLabel = $derived(
    !match || !isEnded
      ? ""
      : match.status === "DRAW"
        ? "Hoà"
        : match.status === `${mySide}_WON`
          ? "Bạn thắng 🎉"
          : "Bạn thua",
  );

  const queryClient = useQueryClient();

  function backToLobby() {
    stateStore.match = null;
    // rating vừa đổi sau trận -> lấy lại cho PlayerCard ở sảnh
    queryClient.invalidateQueries({ queryKey: ["me"] });
    goto("/caro");
  }
</script>

<!-- hàng info của 1 player: avatar, badge, tên, rating, thanh thời gian, đồng hồ -->
{#snippet playerBar(
  mark: Mark,
  name: string,
  rating: number,
  avatarSrc: string | undefined,
  active: boolean,
  progress: number,
  clock: string,
)}
  <div class="flex items-center gap-2.5">
    <div
      class="flex h-11 w-11 shrink-0 items-center justify-center overflow-hidden rounded-full bg-gray-100"
      style={active ? `box-shadow: 0 0 0 2px ${MARK_COLOR[mark]}` : undefined}
    >
      {#if avatarSrc}
        <img src={avatarSrc} alt={name} class="h-full w-full object-cover" />
      {/if}
    </div>

    <div
      class="flex h-6 w-6 shrink-0 items-center justify-center rounded-lg text-white"
      style="background-color: {MARK_COLOR[mark]}"
    >
      {#if mark === "X"}
        <X size={14} strokeWidth={3} />
      {:else}
        <Circle size={12} strokeWidth={3} />
      {/if}
    </div>

    <span class="font-bold">{name}</span>
    <span class="text-sm font-medium text-neutral-400"
      >{rating.toLocaleString()}</span
    >

    <div class="ml-2 h-1.5 flex-1 overflow-hidden rounded-full bg-gray-200">
      <div
        class="h-full rounded-full transition-[width] duration-200 ease-linear"
        style="width: {progress * 100}%; background-color: {active
          ? MARK_COLOR[mark]
          : '#D1D5DB'}"
      ></div>
    </div>

    <span class="w-9 text-right text-sm font-medium text-neutral-500">{clock}</span>
  </div>
{/snippet}

{#snippet playerBarFor(side: Mark)}
  {@const user = side === "X" ? xUserQuery.data : oUserQuery.data}
  {@const rating =
    (side === "X" ? match?.xPlayerRatingBefore : match?.oPlayerRatingBefore) ?? 0}
  {@render playerBar(
    side,
    user?.name ?? "…",
    rating,
    user ? AVATARS[user.avtCode] : undefined,
    isActive(side),
    progressOf(side),
    clockOf(side),
  )}
{/snippet}

<main class="flex min-h-screen justify-center gap-8 px-6 py-4">
  <div class="flex w-full max-w-[560px] flex-col gap-3">
    <!-- status -->
    <div
      class="flex items-center justify-center gap-2 text-xs font-bold tracking-wider text-neutral-500 uppercase"
    >
      <span class="h-1.5 w-1.5 rounded-full bg-amber-400"></span>
      {#if isEnded}
        Trận đã kết thúc
      {:else if isMyTurn}
        Lượt của bạn
      {:else}
        Lượt đối thủ
      {/if}
      <span class="h-1.5 w-1.5 rounded-full bg-amber-400"></span>
    </div>

    <!-- opponent -->
    {@render playerBarFor(oppSide)}

    <!-- board -->
    <div class="relative">
      <div
        class="grid aspect-square w-full gap-px overflow-hidden rounded-lg bg-gray-200"
        style="grid-template-columns: repeat({CARO_BOARD_SIZE}, minmax(0, 1fr)); grid-template-rows: repeat({CARO_BOARD_SIZE}, minmax(0, 1fr));"
      >
        {#each match?.board ?? [] as rowCells, y}
          {#each rowCells as piece, x}
            <button
              type="button"
              disabled={!isMyTurn || piece !== ""}
              onclick={() => handleCellClick(x, y)}
              class="flex items-center justify-center enabled:hover:bg-gray-100 disabled:cursor-default {lastMove &&
              lastMove.x === x &&
              lastMove.y === y
                ? 'bg-gray-100'
                : 'bg-white'}"
              aria-label="ô {x}-{y}"
            >
              {#if piece === "X"}
                <X size={18} strokeWidth={4} color={MARK_COLOR.X} />
              {:else if piece === "O"}
                <Circle size={16} strokeWidth={4} color={MARK_COLOR.O} />
              {/if}
            </button>
          {/each}
        {/each}
      </div>

      {#if match && isEnded}
        <div
          class="absolute inset-0 flex items-center justify-center rounded-lg bg-white/70"
        >
          <div
            class="flex w-[260px] flex-col gap-4 rounded-2xl border border-black/[0.06] bg-white p-5 shadow-lg"
          >
            <div class="text-center text-lg font-bold">{resultLabel}</div>

            {#if match.endReason === "OUT_OF_TIME"}
              <div class="-mt-3 text-center text-xs text-neutral-400">Hết giờ</div>
            {/if}

            <div class="flex items-baseline justify-center gap-2">
              <span class="text-2xl font-extrabold tracking-tight">
                {(myRatingAfter ?? 0).toLocaleString()}
              </span>
              <span
                class="text-sm font-bold {ratingDelta > 0
                  ? 'text-green-600'
                  : ratingDelta < 0
                    ? 'text-red-600'
                    : 'text-neutral-400'}"
              >
                {ratingDelta > 0 ? `+${ratingDelta}` : ratingDelta}
              </span>
            </div>

            <button
              type="button"
              onclick={backToLobby}
              class="w-full rounded-lg bg-blue-700 py-2.5 text-sm font-bold text-white transition-colors hover:bg-blue-800"
            >
              Về sảnh
            </button>
          </div>
        </div>
      {/if}
    </div>

    <!-- me -->
    {@render playerBarFor(mySide)}
  </div>

  <!-- TODO: backend chưa có chat qua WS -> panel tĩnh, chưa wire -->
  <div
    class="flex h-[700px] w-[320px] shrink-0 flex-col rounded-2xl border border-black/[0.06] bg-white shadow-sm"
  >
    <div class="flex items-center justify-between px-5 pt-5 pb-3">
      <span class="font-bold">Match chat</span>
      <span class="flex items-center gap-1.5 text-xs font-medium text-neutral-400">
        <span class="h-2 w-2 rounded-full bg-neutral-300"></span>
        Sắp có
      </span>
    </div>

    <div class="flex flex-1 items-center justify-center px-4">
      <span class="text-center text-sm text-neutral-400">
        Tính năng chat đang được phát triển
      </span>
    </div>

    <div class="flex items-center gap-2 px-4 pb-4">
      <input
        type="text"
        disabled
        placeholder="Say something..."
        class="h-11 min-w-0 flex-1 rounded-xl border border-gray-200 bg-gray-50 px-4 text-sm outline-none placeholder:text-neutral-400 disabled:cursor-default disabled:opacity-60"
      />
      <button
        type="button"
        disabled
        class="flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-amber-400 text-white disabled:cursor-default disabled:opacity-60"
        aria-label="Send"
      >
        <Send class="h-5 w-5" />
      </button>
    </div>
  </div>
</main>
