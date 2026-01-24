<script>
	import { onMount } from 'svelte';
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import { browser } from '$app/environment';
	import { writable } from 'svelte/store';
	import {
		Card,
		Button,
		Col,
		CardHeader,
		Container,
		Input,
		Row,
		CardBody,
		CardText,
		CardTitle
	} from '@sveltestrap/sveltestrap';
	// ui stuff
	import { DataTable } from '@careswitch/svelte-data-table';

	import * as echarts from 'echarts';
	import { symbols } from '../symbols.svelte.js';
	import * as helpers from '../helpers.js';

	let { data } = $props();
	let files = $state();
	let versions = $derived(Object.keys(symbols.symbols));
	let selected_symbols = $state({});
</script>

<hr />

<div class="container" id="content">
	<h3>Select a version to browse its symbols:</h3>

	<Row cols={{ lg: 3, md: 2, sm: 1 }}>
		{#each versions as version}
			<div class="pb-3 px-3">
				<Card>
					<CardHeader>
						<CardTitle>{version}</CardTitle>
					</CardHeader>
					<CardBody>
						<CardText>
							<!-- TODO add source link symbol json -->
							Buildtime: {symbols.symbols[version].timestamp}
						</CardText>
						<Button href={base+"/browse/"+version}>Browse</Button>
					</CardBody>
				</Card>
			</div>
		{/each}
	</Row>
</div>
