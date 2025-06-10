<script>
	import Editor from '../components/Editor.svelte';
	import Header from '../components/Header.svelte';
	import Terminal from '../components/Terminal.svelte';
</script>

<Header></Header>
<div class="main">
	<Editor />
	<Terminal />
</div>

<style>
	.main {
		display: flex;
	}
</style>
