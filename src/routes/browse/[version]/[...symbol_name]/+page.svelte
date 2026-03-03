<script>
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';
	import { page } from '$app/stores';
	import { writable } from 'svelte/store';
	import { base } from '$app/paths';
	import { symbols } from '../../../symbols.svelte.js';

	import { Breadcrumb, BreadcrumbItem, Row, Tooltip } from '@sveltestrap/sveltestrap';
	import { Table } from '@sveltestrap/sveltestrap';
	import Highlight from 'svelte-highlight';
	import atomOneDark from 'svelte-highlight/styles/atom-one-dark';
	import armasm from 'svelte-highlight/languages/armasm';
	import * as echarts from 'echarts';
	import SymbolTable from '../../../../components/SymbolTable.svelte';
	import * as helpers from '../../../helpers.js';
	import FunctionSymbolPage from '../../../../components/FunctionSymbolPage.svelte';
	import VariableSymbolPage from '../../../../components/VariableSymbolPage.svelte';

	let { route, data } = $props();
	let params = $derived(data);

	let symbol_version = $derived(params.version);
	let SymPathByAddr = $derived(symbols.symbols[symbol_version]['SymPathByAddr']);
	let SymPathByName = $derived(symbols.symbols[symbol_version]['symPathByName']);
	let symbol_path_and_name = $derived(params.symbol_name);
	let symbol_path_elements_and_parent_links = $derived(
		symbol_path_and_name
			.split('/')
			.slice(0, -1)
			.map((ele, index, arr) => {
				let parent_path = arr.slice(0, index + 1).join('/');
				let parent_link = helpers.sympath_to_link(base, symbol_version, parent_path);
				return [ele, parent_link];
			})
	);
	let symbol_path = $derived(symbol_path_and_name.split('/').slice(0, -1).join('/'));
	let symbol_path_active = $derived(symbol_path_and_name.split('/').slice(-1));
	const isChildToPath = (symbol) => {
		let pathAsAddr = parseInt(symbol_path_and_name);
		let pathIsAddr = pathAsAddr != NaN;

		if (symbol_path_and_name === '/' || symbol_path_and_name === '') return true;
		if (pathIsAddr && pathAsAddr === symbol.address) return true;
		if (symbol_path_and_name === '/' || symbol_path_and_name === '') return true;
		else {
			let routePath = symbol_path_and_name;
			if (!routePath.startsWith('/')) routePath = '/' + routePath;
			let symPath = symbol.file + '/' + symbol.name + '/';
			return symbol.file && symPath.includes(routePath);
		}
	};
	let sections = $derived(symbols.symbols[symbol_version].sections);
	let fn_childs = $derived(symbols.symbols[symbol_version].functions.filter(isChildToPath).sort());
	let var_childs = $derived(symbols.symbols[symbol_version].variables.filter(isChildToPath).sort());
	$inspect(symbol_path_and_name, fn_childs, var_childs);
	onMount(async () => {
		if (browser && fn_childs) {
			var romChartDom = document.getElementById('sunburst_chart_rom');
			var romChart = echarts.init(romChartDom);
			var ramChartDom = document.getElementById('sunburst_chart_ram');
			var ramChart = echarts.init(ramChartDom);
			var option;

			// TODO rm hack - think about how to distinguish better memory regions for all controller types...
			let rom_syms = fn_childs.filter(
				(e) =>
					!(
						e.address.toString(16).startsWith(2) && // not 0x2... in RAM
						e.address >= 0x20000000
					)
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
						show: false
						// rotate: 'radial',
						// minAngle: 15,
						// formatter: '{b} - {c}'
					}
				},
				tooltip: {
					show: true,
					trigger: 'item'
				}
			};
			option && romChart.setOption(option);

			let ram_syms = var_childs.filter((e) => {
				return (
					e.address.toString(16).startsWith(2) && // 0x2... in RAM
					e.address >= 0x20000000
				);
			});
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
						show: false
						// rotate: 'radial',
						// minAngle: 15,
						// formatter: '{b} - {c}'
					}
				},
				tooltip: {
					show: true,
					trigger: 'item'
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
		{#if var_childs.length == 0 && fn_childs.length == 1}
			{console.log('symbol_data: ', $state.snapshot(fn_childs[0]))}
			<FunctionSymbolPage {fn_childs} {symbol_version} {SymPathByAddr} />
		{:else if var_childs.length == 1 && fn_childs.length == 0}
			{console.log('symbol_data: ', $state.snapshot(var_childs[0]))}
			<VariableSymbolPage {var_childs} {symbol_version} {SymPathByAddr} />
		{:else if fn_childs}
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
			{fn_childs.length} function and {var_childs.length} variable symbols in this path: '/{symbol_path}'

			<SymbolTable
				fnSymbols={fn_childs}
				varSymbols={var_childs}
				selected_version={symbol_version}
				{sections}
			/>
		{:else}
			404 - nonononon
		{/if}
	{/key}
</div>
