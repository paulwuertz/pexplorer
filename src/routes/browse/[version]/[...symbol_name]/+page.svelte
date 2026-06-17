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
	import armasm from 'svelte-highlight/languages/armasm';
	import SymbolTable from '../../../../components/SymbolTable.svelte';
	import * as helpers from '../../../helpers.js';
	import FunctionSymbolPage from '../../../../components/FunctionSymbolPage.svelte';
	import VariableSymbolPage from '../../../../components/VariableSymbolPage.svelte';
	import SunburstTreeMapGraph from '../../../../components/SunburstTreeMapGraph.svelte';

	let { route, data } = $props();
	let params = $derived(data);

	// let sunburstGraphRenderAngle = $state(0);
	let treeMapRenderDepth = $state(6);
	let symbol_version = $derived(params.version);
	let fw_symbols = $derived(symbols.symbols[symbol_version]);
	let SymPathByAddr = $derived(fw_symbols['SymPathByAddr']);
	let SymPathByName = $derived(fw_symbols['symPathByName']);
	let VariableTypes = $derived(fw_symbols['types']);
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
	let sections = $derived(fw_symbols.sections);
	let fn_childs = $derived(
		((fw_symbols && fw_symbols.functions) || []).filter(isChildToPath).sort()
	);
	let var_childs = $derived(
		((fw_symbols && fw_symbols.variables) || []).filter(isChildToPath).sort()
	);
	let symbols_in_symbols = $derived([...fn_childs, ...var_childs]);
</script>

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
			<VariableSymbolPage
				{var_childs}
				{symbol_version}
				{SymPathByAddr}
				{VariableTypes}
				{sections}
			/>
		{:else if fn_childs}
			<SunburstTreeMapGraph {symbols_in_symbols} {treeMapRenderDepth} />
			{fn_childs.length} function and {var_childs.length} variable symbols in this path: '/{symbol_path}'

			<SymbolTable
				fnSymbols={fn_childs}
				varSymbols={var_childs}
				selected_version={symbol_version}
				{sections}
				bind:render_depth={treeMapRenderDepth}
			/>
		{:else}
			404 - nonononon
		{/if}
	{/key}
</div>
