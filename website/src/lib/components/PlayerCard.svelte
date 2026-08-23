<script lang="ts">
  import type { User } from "$lib/api/type";
  import { Pencil } from "@lucide/svelte";
  import Bunny from "$lib/assets/Bunny.png";
  import Kitty from "$lib/assets/Kitty.png";
  import Teddy from "$lib/assets/Teddy.png";
  import Hamster from "$lib/assets/Hamster.png";
  import Monkey from "$lib/assets/Monkey.png";
  import Piggy from "$lib/assets/Piggy.png";

  const avatars: Record<User["avtCode"], string> = {
    BUNNY: Bunny,
    KITTY: Kitty,
    TEDDY: Teddy,
    HAMSTER: Hamster,
    MONKEY: Monkey,
    PIGGY: Piggy,
  };

  let {
    player,
    isLoading,
  }: {
    player?: User;
    isLoading: boolean;
  } = $props();

  const avatarSrc = $derived(player ? avatars[player.avtCode] : undefined);
</script>

<div
  class="flex items-center justify-between rounded-2xl border border-black/[0.06] bg-white px-5 py-4 shadow-sm"
>
  {#if isLoading}
    <div class="flex items-center gap-3">
      <div class="h-11 w-11 animate-pulse rounded-full bg-stone-200"></div>
      <div class="flex flex-col gap-2">
        <div class="h-4 w-24 animate-pulse rounded bg-stone-200"></div>
        <div class="h-2.5 w-8 animate-pulse rounded bg-stone-200"></div>
      </div>
    </div>
    <div class="flex items-center gap-4">
      <div class="h-8 w-px bg-black/10"></div>
      <div class="flex flex-col items-end gap-2">
        <div class="h-6 w-14 animate-pulse rounded bg-stone-200"></div>
        <div class="h-2.5 w-10 animate-pulse rounded bg-stone-200"></div>
      </div>
    </div>
  {:else}
    <div class="flex items-center gap-3">
      <div
        class="flex h-12 w-12 items-center justify-center overflow-hidden rounded-full bg-stone-200"
      >
        {#if avatarSrc}
          <img src={avatarSrc} alt={player?.name} class="h-full w-full object-cover" />
        {/if}
      </div>
      <div class="flex items-center gap-2">
        <div class="flex flex-col gap-1">
          <div class="font-bold leading-tight">{player?.name}</div>
          <div class="text-[11px] font-semibold tracking-wide text-neutral-400">
            BẠN
          </div>
        </div>
      </div>
    </div>
    <div class="flex items-center gap-4">
      <div class="h-8 w-px bg-black/10"></div>
      <div class="text-right flex flex-col gap-1">
        <div class="text-xl font-extrabold leading-none">
          {player?.caroRating.toLocaleString()}
        </div>
        <div class="text-[11px] font-semibold tracking-wide text-neutral-400">
          RATING
        </div>
      </div>
    </div>
  {/if}
</div>
