<script>
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';
	import { page } from '$app/stores';
	import { writable } from 'svelte/store';
	import { base } from '$app/paths';
	import { symbols } from '../../../symbols.svelte.js';

	import { Breadcrumb, BreadcrumbItem, Row } from '@sveltestrap/sveltestrap';
	import { Table } from '@sveltestrap/sveltestrap';
	import Highlight from 'svelte-highlight';
	import atomOneDark from 'svelte-highlight/styles/atom-one-dark';
	import armasm from 'svelte-highlight/languages/armasm';
	import * as echarts from 'echarts';
	import FunctionSymbolTable from '../../../../components/FunctionSymbolTable.svelte';
	import * as helpers from '../../../helpers.js';
	import FunctionSymbolPage from '../../../../components/FunctionSymbolPage.svelte';

	let { route, data } = $props();
	let params = $derived(data);

	let symbol_version = $derived(params.version);
	let symbol_path_and_name = $derived(params.symbol_name);
	let symbol_path_elements_and_parent_links = $derived(
		symbol_path_and_name
			.split('/')
			.slice(0, -1)
			.map((ele, index, arr) => {
				let parent_path = arr.slice(0, index + 1).join('/');
				let parent_link = helpers.callxrs_text_to_links(base, symbol_version, parent_path);
				return [ele, parent_link];
			})
	);
	let symbol_path_active = $derived(symbol_path_and_name.split('/').slice(-1));
	let symbol_data = $derived(
		symbols.symbols[symbol_version].symbols.find((e) => {
			return symbol_path_and_name.includes(e.file);
		})
	);
	const isChildToPath = (symbol) => symbol.file.includes(symbol_path_and_name);
	let symbol_childs = $derived(
		symbols.symbols[symbol_version].symbols.filter(isChildToPath).sort()
	);

	onMount(async () => {
		if (browser && symbol_childs) {
			var romChartDom = document.getElementById('sunburst_chart_rom');
			var romChart = echarts.init(romChartDom);
			var ramChartDom = document.getElementById('sunburst_chart_ram');
			var ramChart = echarts.init(ramChartDom);
			var option;

			// TODO rm hack - think about how to distinguish better memory regions for all controller types...
			let rom_syms = symbol_childs.filter(
				(e) => e.address.startsWith(0) || e.address.startsWith(8)
			);
			var sunburst_data = helpers.symbols_to_sunburst_tree_data(rom_syms, 'size');
			console.log('sunburst_data', sunburst_data);
			option = {
				// visualMap: {
				//     type: 'continuous',
				//     inRange: {
				//         color: ['#2F93C8', '#AEC48F', '#FFDB5C', '#F98862']
				//     }
				// },
				series: {
					type: 'sunburst',
					data: sunburst_data,
					radius: [0, '90%'],
					label: {
						rotate: 'radial',
						minAngle: 15,
						formatter: '{b} - {c}'
					}
				}
			};
			option && romChart.setOption(option);

			let ram_syms = symbol_childs.filter((e) => e.address.startsWith(2));
			var sunburst_data = helpers.symbols_to_sunburst_tree_data(ram_syms, 'size');
			console.log('sunburst_data', sunburst_data);
			option = {
				// visualMap: {
				//     type: 'continuous',
				//     inRange: {
				//         color: ['#2F93C8', '#AEC48F', '#FFDB5C', '#F98862']
				//     }
				// },
				series: {
					type: 'sunburst',
					data: sunburst_data,
					radius: [0, '90%'],
					label: {
						rotate: 'radial',
						minAngle: 15,
						formatter: '{b} - {c}'
					}
				}
			};
			option && ramChart.setOption(option);
		}
	});
</script>

<svelte:head>
	{@html atomOneDark}
</svelte:head>

<div class="container">
	<hr />

	<Breadcrumb divider="/">
		{#each symbol_path_elements_and_parent_links as path_element_and_parent_links}
			<BreadcrumbItem>
				<a href={path_element_and_parent_links[1]}>{path_element_and_parent_links[0]}</a>
			</BreadcrumbItem>
		{/each}
		<BreadcrumbItem active>{symbol_path_active}</BreadcrumbItem>
	</Breadcrumb>

	<hr />

    {#key symbol_path_and_name}
	{#if symbol_data}
        <FunctionSymbolPage
            symbol_data={symbol_data}
            symbol_version={symbol_version}
        />
    {:else if symbol_childs}
		<Row cols={{ md: 2, sm: 1 }}>
			<div>
				<h3>Flash usage</h3>
				<div id="sunburst_chart_rom" style="width: 100%;height:600px;"></div>
			</div>
			<div>
				<h3>Static RAM usage</h3>
				<div id="sunburst_chart_ram" style="width: 100%;height:600px;"></div>
			</div>
		</Row>
		{symbol_childs.length} child symbols in this path.

        <FunctionSymbolTable
            fnSymbols={symbol_childs}
            selected_version={symbol_version}
        />
	{:else}
		404 - nonononon
	{/if}
    {/key}
</div>
