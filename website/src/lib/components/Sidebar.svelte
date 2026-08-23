<script lang="ts">
  import { Dialog } from "bits-ui";
  import { Hash, LogOut, ALargeSmall, Menu, X } from "@lucide/svelte";

  type NavItem = {
    label: string;
    icon: typeof Hash;
    iconBg: string;
    active?: boolean;
  };

  const nav: NavItem[] = [
    { label: "Cờ Caro", icon: Hash, iconBg: "bg-blue-700", active: true },
    { label: "Nối từ", icon: ALargeSmall, iconBg: "bg-emerald-500" },
  ];

  let open = $state(false);
</script>

<Dialog.Root bind:open>
  <!-- Menu trigger -->
  <Dialog.Trigger
    class="fixed left-4 top-4 z-30 flex h-10 w-10 items-center justify-center rounded-xl border border-black/[0.06] bg-white text-neutral-900 shadow-sm transition duration-75 hover:bg-black/5"
    aria-label="Mở menu"
  >
    <Menu class="h-5 w-5" />
  </Dialog.Trigger>

  <Dialog.Portal>
    <Dialog.Overlay class="sheet-overlay fixed inset-0 z-40 bg-black/40" />
    <Dialog.Content
      class="sheet-content fixed left-0 top-0 z-50 flex h-screen w-64 flex-col border-r border-black/[0.06] bg-white px-4 py-5 shadow-xl focus:outline-none"
    >
      <!-- Header -->
      <div class="mb-4 flex items-center justify-between px-2">
        <Dialog.Title class="text-lg font-extrabold tracking-tight">Dicey.gg</Dialog.Title>
        <Dialog.Close
          class="flex h-8 w-8 items-center justify-center rounded-lg text-neutral-500 transition duration-75 hover:bg-black/5 hover:text-neutral-900"
          aria-label="Đóng menu"
        >
          <X class="h-5 w-5" />
        </Dialog.Close>
      </div>
      <Dialog.Description class="sr-only">Điều hướng chính</Dialog.Description>

      <!-- Nav -->
      <nav class="flex flex-col gap-1">
        {#each nav as item (item.label)}
          <Dialog.Close
            class="flex items-center gap-3 rounded-xl px-2 py-2 text-left text-sm font-semibold transition duration-75
              text-neutral-900 hover:bg-black/5 {item.active ? 'bg-black/5' : ''}"
          >
            <span class="flex h-8 w-8 items-center justify-center rounded-lg text-white shadow-sm {item.iconBg}">
              <item.icon class="h-4 w-4" />
            </span>
            {item.label}
          </Dialog.Close>
        {/each}
      </nav>

      <!-- Bottom -->
      <div class="mt-auto">
        <button
          type="button"
          class="flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-left text-sm font-semibold text-neutral-900 transition duration-75 hover:bg-black/5"
        >
          <LogOut class="h-5 w-5" />
          Đăng xuất
        </button>
      </div>
    </Dialog.Content>
  </Dialog.Portal>
</Dialog.Root>
