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

</script>

