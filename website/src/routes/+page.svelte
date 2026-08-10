<script lang="ts">
  import { auth } from "$lib/core/auth";
  import { authStatus } from "$lib/stores/auth-status";
  import { Button } from "bits-ui";
  import { toast } from "svelte-sonner";

  async function onLogInClick() {
    await auth.signInWithOAuth({
      provider: "google",
      options: {
        redirectTo: window.location.origin,
      },
    });
  }
</script>

<h1>Welcome to Project EG</h1>
{#if !$authStatus.isReady}
  <div>Loading...</div>
{:else if $authStatus.session}
  <h1>Hello {$authStatus.session.user.email}</h1>
  <Button.Root class="cursor-pointer bg-red-500" onclick={() => {
    // toast("U want to sign out?", {
    //   duration: Infinity,
    //   cancel: {
    //     label: "Cancel",
    //     onClick: () => {}
    //   }
    // });
    auth.signOut();
  }}>Log out</Button.Root>
{:else}
  <Button.Root onclick={onLogInClick}>Log in with Google</Button.Root>
{/if}
