<script lang="ts">
	import "./layout.css";
	import favicon from "$lib/assets/favicon.svg";
	import { onMount } from "svelte";
	import { auth } from "$lib/core/auth";
	import { authStore } from "$lib/stores/auth-store.svelte";
	import { Toaster } from "svelte-sonner";
	import { QueryClient, QueryClientProvider } from "@tanstack/svelte-query";
	import { browser } from "$app/environment";
	import { ensureAuth } from "$lib/core/ensure-auth";
	import Sidebar from "$lib/components/Sidebar.svelte";

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
			authStore.session = session;
			authStore.isReady = true;
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
	<div class="min-h-screen bg-stone-100 text-neutral-900">
		<Sidebar />
		{@render children()}
	</div>
</QueryClientProvider>
