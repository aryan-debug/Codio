<script lang="ts">
	import CodeMirror from '../../node_modules/svelte-codemirror-editor/dist/CodeMirror.svelte';
	import { go } from '@codemirror/lang-go';
	import { barf } from '../themes/modifiedBarf';
	import { state } from '../store.svelte';
</script>

<CodeMirror class="editor" bind:value={state.code} lang={go()} theme={barf} />

<style>
	:global(.editor) {
		width: 55%;
		margin: 0 2em;
	}
</style>
