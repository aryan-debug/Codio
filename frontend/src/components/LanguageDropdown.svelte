<script lang="ts">
	import { state } from '../store.svelte';

	const languages = [
		{ id: 1, name: 'Python', value: 'python' },
		{ id: 2, name: 'C++', value: 'cpp' },
		{ id: 3, name: 'Go', value: 'go' }
	];
</script>

<select class="dropdown" bind:value={state.language}>
	{#each languages as language (language.id)}
		<option value={language.value}>{language.name}</option>
	{/each}
</select>

<style>
	.dropdown {
		font-weight: 500;
		font-family: var(--font-inter);
		font-size: 1.5em;
		color: white;
		background-color: var(--light-grey);
		border: none;
	}
</style>
