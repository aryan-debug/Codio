<script>
	import Editor from '../components/Editor.svelte';
	import Header from '../components/Header.svelte';
</script>

<Header></Header>
<Editor />

<style>
	:global(:root) {
		--light-grey: #1a1a1a;
		--accent-color: #33bf69;
		--font-inter: 'Inter', sans-serif;
	}

	:global(body) {
		background-color: var(--light-grey);
	}
</style>
