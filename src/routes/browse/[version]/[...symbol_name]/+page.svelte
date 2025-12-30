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
	import * as helpers from '../../../helpers.js';

	let { route, data } = $props();
	let params = $derived(data);
	let show_full_asm = $state(false);

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

	// to display
	let asm = $derived(Object.hasOwn(symbol_data, 'asm') ? symbol_data.asm : []);
	let asm_code_preview = $derived(asm.slice(0, 5).join('\n'));
	let address = $derived(symbol_data.address);
	let asm_code = $derived(asm.join('\n'));
	let callers = $derived(JSON.parse(symbol_data.callers));
	let callees = $derived(JSON.parse(symbol_data.callees));
	let deepest_callers_tree = $derived(JSON.parse(symbol_data.deepest_caller_tree || false));
	let deepest_callees_tree = $derived(JSON.parse(symbol_data.deepest_callee_tree || false));
	let code_size = $derived(symbol_data.size);

	const worst_call_stack = () => {
		let my_symbol = { full_symbol_path: symbol_path_and_name, stack_size: symbol_data.stack_size };
		let stack_down = deepest_callees_tree.concat([my_symbol]);
		return stack_down.concat(deepest_callers_tree);
	};

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
        <Table style="word-break: break-all;" hover bordered>
			<tbody>
				<tr>
					<td><b>Address</b>:</td>
					<td>
						0x{address}
					</td>
				</tr>

				<tr>
					<td><b>Function code size</b>:</td>
					<td>
						{code_size} bytes
					</td>
				</tr>

				<tr>
					<td><b>Callers </b> ({callers.length}):</td>
					<td>
						{#each callers as caller}
							<a href={helpers.callxrs_text_to_links(base, symbol_version, caller)}>
								<small>
									{helpers.callxrs_text_to_symname(caller)}
								</small>
							</a>{', '}
						{/each}
					</td>
				</tr>
				<tr>
					<td><b>Callees</b> ({callees.length}):</td>
					<td>
						{#each callees as callee}
							<a href={helpers.callxrs_text_to_links(base, symbol_version, callee)}>
								<small>
									{helpers.callxrs_text_to_symname(callee)}
								</small>
							</a>{', '}
						{/each}
					</td>
				</tr>
			</tbody>
		</Table>

		<h4>Disassembly</h4>
		<!-- TODO figure out better export to highlight and diff... <Highlight language={armasm} {asm_code} /> -->
		<pre>
        {#if show_full_asm}
				{asm_code}
			{:else}
				{asm_code_preview}
            ...
			{/if}
        <span class="center" onclick={() => (show_full_asm = !show_full_asm)}>
            {#if show_full_asm}↑ show less ↑{:else}↓ show more ↓{/if}
        </span>
    </pre>

		<h4>Stack Worst-Case Scenarios</h4>

		<Table style="word-break: break-all;" hover bordered>
			<thead>
				<tr>
					<th>#</th>
					<th>Name</th>
					<th>Stack size</th>
				</tr>
			</thead>
			<tbody>
				{#each worst_call_stack() as caller, index}
					<tr>
						<td>{index + ' '}</td>
						<td>
							<a
								href={helpers.callxrs_text_to_links(base, symbol_version, caller.full_symbol_path)}
							>
								{#if symbol_path_and_name.includes(caller.full_symbol_path)}
									<small>
										<b>{helpers.callxrs_text_to_symname(caller.full_symbol_path)}</b> - (this function)
									</small>
								{:else}
									<small>
										{helpers.callxrs_text_to_symname(caller.full_symbol_path)}
									</small>
								{/if}
							</a>
						</td>
						<td>
							{caller.stack_size}
						</td>
					</tr>
				{/each}
			</tbody>
			<tfoot>
				<tr>
					<td></td>
					<td></td>
					<td
						><b
							>&sum; {symbol_data.deepest_callee_tree_size +
								symbol_data.deepest_caller_tree_size}</b
						></td
					>
				</tr>
			</tfoot>
		</Table>
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

		<!-- TODO reuse table from symbol page -->
		<Table style="word-break: break-all;" hover bordered>
			<tbody>
				<tr>
					<td><b>Symbol name</b>:</td>
				</tr>
				{#each symbol_childs as child}
					<tr>
						<td
							><a
								data-sveltekit-preload-data="tap"
								href={helpers.row2AHref(base, symbol_version, child)}
							>
								/{child.file + ':' + child.display_name + ':' + child.line}
							</a></td
						>
					</tr>
				{/each}
			</tbody>
		</Table>
	{:else}
		404 - nonononon
	{/if}
    {/key}
</div>

<style>
	pre {
		background-color: #f5f5f5;
		border: 1px solid #ccc;
		border-radius: 4px;
	}
</style>
