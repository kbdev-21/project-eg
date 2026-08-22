<script lang="ts">
	import "./layout.css";
	import favicon from "$lib/assets/favicon.svg";
	import { onMount } from "svelte";
	import { auth } from "$lib/core/auth";
	import { authStore } from "$lib/stores/auth-store";
	import { Toaster } from "svelte-sonner";
	import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
	import { browser } from "$app/environment";
	import { ensureAuth } from "$lib/core/ensure-auth";

	let { children } = $props();

	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				enabled: browser,
			},
		},
	});

	onMount(() => {
		const { data: listener } = auth.onAuthStateChange((_, session) => {
			authStore.set({
				session,
				isReady: true,
			});
		});

		ensureAuth();

		return () => {
			listener.subscription.unsubscribe();
		};
	});
</script>

<svelte:head>
	<title>Project EG</title>
	<link rel="icon" href={favicon} />
</svelte:head>

<QueryClientProvider client={queryClient}>
	<Toaster position="top-center" />
	{@render children()}
</QueryClientProvider>
