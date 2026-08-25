<script lang="ts">
  import type { Snippet } from "svelte";
  import { ChevronRight } from "@lucide/svelte";

  let {
    title,
    subtitle,
    iconBg,
    icon,
    trailing,
    onclick,
    disabled = false,
  }: {
    title: string;
    subtitle: string;
    iconBg: string;
    icon: Snippet;
    trailing?: Snippet;
    onclick?: () => void;
    disabled?: boolean;
  } = $props();
</script>

<button
  type="button"
  {onclick}
  {disabled}
  class="flex w-full items-center justify-between rounded-2xl border border-black/[0.06] bg-white px-4 py-4 text-left shadow-sm transition duration-150 enabled:hover:-translate-y-0.5 enabled:hover:border-black/10 enabled:hover:shadow-md disabled:cursor-default disabled:opacity-60"
>
  <div class="flex items-center gap-3">
    <div class="flex h-11 w-11 items-center justify-center rounded-xl {iconBg}">
      {@render icon()}
    </div>
    <div class="flex flex-col gap-1">
      <div class="font-bold leading-tight">{title}</div>
      <div class="text-xs font-medium text-neutral-500">{subtitle}</div>
    </div>
  </div>
  {#if trailing}
    {@render trailing()}
  {:else}
    <ChevronRight class="h-5 w-5 text-neutral-300" />
  {/if}
</button>
