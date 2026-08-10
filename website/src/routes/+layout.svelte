<script lang="ts">
	import "./layout.css";
	import favicon from "$lib/assets/favicon.svg";
	import { onMount } from "svelte";
	import { auth } from "$lib/core/auth";
	import { authStatus } from "$lib/stores/auth-status";
	import { Toaster } from "svelte-sonner";

	let { children } = $props();

	onMount(() => {
		const { data: listener } = auth.onAuthStateChange((_, session) => {
			authStatus.set({
				session,
				isReady: true,
			});
		});

		auth.getSession().then(({ data }) => {
			authStatus.set({
				session: data.session,
				isReady: true,
			});
		});

		return () => {
			listener.subscription.unsubscribe();
		};
	});
</script>

<svelte:head>
	<title>Project EG</title>
	<link rel="icon" href={favicon} />
</svelte:head>

<Toaster position="top-center"/>
{@render children()}
