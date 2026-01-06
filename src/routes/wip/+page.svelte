<script>
	import { onMount } from 'svelte';
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import { browser } from '$app/environment';
	import { writable } from 'svelte/store';
	// ui stuff
	import { DataTable } from '@careswitch/svelte-data-table';

	import { symbols } from '../symbols.svelte.js';

	import * as helpers from '../helpers.js';

	let { data } = $props();

    const go = new Go(); // Defined in wasm_exec.js
    const WASM_URL = '/sELFperf.wasm';

    var wasm;

	onMount(async () => {
        fetch(WASM_URL).then(resp =>
            resp.arrayBuffer()
        ).then(bytes =>
            WebAssembly.instantiate(bytes, go.importObject).then(function (obj) {
                wasm = obj.instance;
                go.run(wasm);
            })
        )
    })
</script>
