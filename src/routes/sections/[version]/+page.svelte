<script>
	import { onMount } from 'svelte';
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import { browser } from '$app/environment';
	import { writable } from 'svelte/store';
	// ui stuff
	import {
		Badge,
		Button,
		ButtonGroup,
		Col,
		Container,
		Input,
		Row,
		Table,
		Card,
		CardHeader,
		CardBody,
		CardText,
		CardTitle,
		CardSubtitle,
		Progress,
		CardFooter
	} from '@sveltestrap/sveltestrap';

	import { symbols } from '../../symbols.svelte.js';
	import * as helpers from '../../helpers.js';
	import * as echarts from 'echarts';
	import ThreadInfoCard from '../../../components/ThreadInfoCard.svelte';

	const isInMemories = (sec) => {
		return sec.ram_size != 0 || sec.rom_size != 0;
	};

	let params = $props();
	// TODO - why are propsed nested +1 here on production build?!
	let parameters = $derived(params && (params.data.data || params.data));
	let version_name = $derived((parameters && parameters.version) || '');
	let version = $derived(symbols.symbols[version_name] || {});
	let secs = $derived(version.sections || []);
	let sections = $derived(secs.filter(isInMemories));

	$inspect(sections);
</script>

<div class="container" id="content">
	<Container fluid>
		<h3>version_name: {JSON.stringify(version_name)}</h3>
		<h4>Static threads by K_THREAD_DEFINE</h4>

		<table>
			<thead>
				<tr>
					<th>name</th>
					<th>address from</th>
					<th></th>
					<th>address to</th>
					<th>size</th>
					<th>remark</th>
				</tr>
			</thead>
			<tbody>
				{#each sections.sort((a, b) => a.address - b.address) as section (section.name)}
					<tr>
						<td>{section.name}</td>
						<td>0x{section.address.toString(16)}</td>
						<td> - </td>
						<td>0x{(section.address + section.size).toString(16)}</td>
						<td>{section.size}</td>
						<td>
							{#if section.ram_size != 0}
								RAM
							{/if}
							{#if section.rom_size != 0}
								ROM
							{/if}
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</Container>
</div>
