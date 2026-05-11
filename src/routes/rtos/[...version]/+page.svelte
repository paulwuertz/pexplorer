<script>
	import { onMount } from 'svelte';
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import { browser } from '$app/environment';
	import { writable } from 'svelte/store';
	// ui stuff
	import { Badge, Button, Col, Container, Input, Row, Table } from '@sveltestrap/sveltestrap';

	import { symbols } from '../../symbols.svelte.js';
	import * as helpers from '../../helpers.js';
	import VersionSelect from '../../../components/VersionSelect.svelte';

	import * as echarts from 'echarts';

	const PLOT_ID_PREFIX = 'stackSizeLayered_';

	let { data } = $props();
	let files = $state();
	let versions = $derived(Object.keys(symbols.symbols));
	let selected_thread_stat = $state({});

</script>

<div class="container" id="content">

<VersionSelect page_name={"rtos"}></VersionSelect>
	<Container fluid>
	</Container>
</div>
