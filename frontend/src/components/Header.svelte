<script lang="ts">
	import Fa from '../../node_modules/svelte-fa/dist/fa.svelte';
	import LanguageDropdown from './LanguageDropdown.svelte';
	import { faPlay } from '@fortawesome/free-solid-svg-icons';
	import { state } from '../store.svelte';

	async function getOutput() {
		const { code, language } = state;
		const request = await fetch('http://localhost:8080/api/run', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ code, language })
		});
		state.result = await request.json();
	}
</script>

<div class="header">
	<div class="left">
		<h2 class="app-name">Codio</h2>
	</div>
	<div class="right">
		<LanguageDropdown />
		<button class="run-button" type="submit" on:click={getOutput}><Fa icon={faPlay} /></button>
	</div>
</div>

<style>
	.header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin: 0 2em;
	}
	.app-name {
		color: white;
		font-family: var(--font-inter);
		font-optical-sizing: auto;
		font-size: 2em;
		font-weight: bold;
	}

	.run-button {
		color: var(--accent-color);
		background-color: var(--light-grey-color);
		border: none;
		font-size: 1.5em;
	}

	.right {
		display: flex;
		gap: 40px;
	}
</style>
