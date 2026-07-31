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
	import { json } from '@sveltejs/kit';

	const isInMemories = (sec) => {
		return sec.ram_size != 0 || sec.rom_size != 0;
	};

	let params = $props();
	// TODO - why are propsed nested +1 here on production build?!
	let parameters = $derived(params && (params.data.data || params.data));
	let version_name = $derived((parameters && parameters.version) || '');
	let version = $derived(symbols.symbols[version_name] || {});
	let functions = $derived(version.functions || []);
	let variables = $derived(version.variables || []);
	let secs = $derived(version.sections || []);
	let sections = $derived(secs.filter(isInMemories));

	let sum_section_id_size = (symbols) => {
		let sec_id_size = {};
		let nr_no_size_var = 0;
		for (let s of symbols) {
			if (!Object.hasOwn(sec_id_size, s.secidx)) {
				sec_id_size[s.secidx] = 0;
			}
			sec_id_size[s.secidx] += s.size;
		}
		console.log('nr_no_size_var: ', nr_no_size_var);
		return sec_id_size;
	};

	let sections_fill = $derived(sum_section_id_size([...functions, ...variables]));

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
					<th>usage</th>
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
						<td>{sections_fill[section.index]}</td>
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
		{JSON.stringify(sections_fill)}
	</Container>
</div>
