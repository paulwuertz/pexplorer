<script>
    import { page } from '$app/stores';
    import { writable } from "svelte/store";
    import { base } from '$app/paths';

	let { data } = $props();
</script>

<div class="flex">
    {JSON.stringify(data.currentSubfolders)}
    <a href="/symbols/home/p4w5/git/cannecti/cannectivity/app/">test</a>
    <ul>
        {#each data.currentSubfolders as sym}
        <li><a href="./{sym}">{sym}</a></li>
		<!--li><a href="/symbols/{sym.file}">{sym.display_name}</a></li-->
        {/each}
    </ul>

</div>

<style>
	.flex {
		display: flex;
		height: 100%;
		flex-direction: column;
		align-items: center;
		justify-content: center;
	}

	p {
		margin: 0.5rem 0;
		line-height: 1;
	}
</style>
