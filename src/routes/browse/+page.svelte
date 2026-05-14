<script>
	import { onMount } from 'svelte';
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import { browser } from '$app/environment';
	import { writable } from 'svelte/store';

	import * as echarts from 'echarts';
	import * as helpers from '../helpers.js';
	import VersionSelect from '../../components/VersionSelect.svelte';
</script>

<hr />

<VersionSelect page_name={'browse'}></VersionSelect>
