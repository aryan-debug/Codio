<script lang="ts">
	import CodeMirror from '../../node_modules/svelte-codemirror-editor/dist/CodeMirror.svelte';
	import { go } from '@codemirror/lang-go';
	import { barf } from '../themes/modifiedBarf';
	let value = '';
</script>

<CodeMirror bind:value lang={go()} theme={barf} />
