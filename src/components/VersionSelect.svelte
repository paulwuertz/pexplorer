<script>
	import { onMount } from 'svelte';
	import { base } from '$app/paths';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';

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
	import { symbols } from '../routes/symbols.svelte.js';
	import * as helpers from '../routes/helpers.js';

	const { page_name } = $props();

	let versions = $derived(Object.keys(symbols.symbols));
</script>

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
						<Button href={base + '/#/' + page_name + '/' + version}>Browse</Button>
					</CardBody>
				</Card>
			</div>
		{/each}
	</Row>
</div>
