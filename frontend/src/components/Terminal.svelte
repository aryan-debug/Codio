<script lang="ts">
	import { state } from '../store.svelte';
</script>

<pre class="terminal">{state.result.output}{state.result.error}</pre>

<style>
	.terminal {
		flex: 1;
		background-color: black;
		border: 1px solid var(--light-light-grey);
		border-radius: 10px;
		margin: 0 2em 0 0;
		color: white;
		padding: 0.5em;
		font-size: 1.2em;
	}
</style>
