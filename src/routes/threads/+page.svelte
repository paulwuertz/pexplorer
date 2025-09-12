<script>
	import { onMount } from 'svelte';
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import { browser } from '$app/environment';
	import { writable } from 'svelte/store';
	// ui stuff
	import { Badge, Button, Col, Container, Input, Row, Table } from '@sveltestrap/sveltestrap';

	import { symbols } from '../symbols.svelte.js';
	import * as helpers from '../helpers.js';

	import * as echarts from 'echarts';

	const PLOT_ID_PREFIX = 'stackSizeLayered_';

	let { data } = $props();
	let files = $state();
	let versions = $derived(Object.keys(symbols.symbols));
	let selected_thread_stat = $state({});

	const updateSelectedSymbols = () => {
		let versionObj = symbols.symbols[symbols.selected_version];
		if (!versionObj || !versionObj.hasOwnProperty('stack_reports')) return;
		let plotsContainer = document.getElementById('plotsContainer');
		plotsContainer.innerHTML = '';

		let firmware_versions = helpers.get_versions_ordered_by_timestamps(symbols.symbols);
		let all_threads_names = helpers.get_all_threads_names(symbols.symbols);
		let all_fuctions_names = helpers.get_all_threads_function_names_on_stacks(symbols.symbols);
		let function_size_series_map = {};
		// all sizes default to 0
		for (const fuctions_names of all_fuctions_names) {
			function_size_series_map[fuctions_names] = Array.from(
				{ length: firmware_versions.length },
				(v, i) => 0
			);
		}
		// set sizes for function for all versions where it is defined
		for (const functions_name of all_fuctions_names) {
			let found = false;
			for (const [i, version] of firmware_versions.entries()) {
				let versionStr = version.version;
				// console.log(versionStr, symbols.symbols[versionStr]);

				let versionReport = symbols.symbols[versionStr]['stack_reports'];
				for (const fn_call_stack_info of Object.values(versionReport)) {
					let fn_call_stack = fn_call_stack_info.call_stack;
					let symVersion = fn_call_stack.find((e) => e['function'] == functions_name);
					if (!symVersion) {
						continue;
					}
					function_size_series_map[functions_name][i] = symVersion['stack_size'];
				}
			}
		}
		for (const [thread_name, threadObj] of Object.entries(versionObj['stack_reports'])) {
			const plotContainerID = PLOT_ID_PREFIX + thread_name;
			selected_thread_stat[thread_name] = threadObj;
			plotsContainer.innerHTML +=
				`
            <div class="col">
                <h6>` +
				thread_name +
				`</h6>
                <div style="width: 100%;height:400px;" id="` +
				plotContainerID +
				`"></div>
            </div>
            `;
		}

		for (const [thread_name, thread_info] of Object.entries(selected_thread_stat)) {
			const plotContainerID = PLOT_ID_PREFIX + thread_name;
			console.log('idddd', plotContainerID, document.getElementById(plotContainerID));

			// Create the echarts instance
			var maxStackSizeChart = echarts.init(document.getElementById(plotContainerID), {
				width: '95%',
				height: 400
			});

			let thread_function_stack_series = [];
			for (const function_info of thread_info['call_stack']) {
				let function_name = function_info['function'];
				let function_stack_size_by_versions = function_size_series_map[function_name];
				thread_function_stack_series.push({
					name: function_name,
					type: 'line',
					stack: 'Total',
					areaStyle: {},
					emphasis: {
						focus: 'series'
					},
					step: true,
					data: function_stack_size_by_versions
				});
			}

			// Draw the chart
			maxStackSizeChart.setOption({
				title: { text: '' },
				colorBy: 'series',
				color: [
					'#e5f5e0',
					'#c7e9c0',
					'#a1d99b',
					'#74c476',
					'#41ab5d',
					'#238b45',
					'#006d2c',
					'#00441b'
				],
				tooltip: {
					trigger: 'axis'
				},
				// legend: {
				// 	data: ['Email', 'Union Ads']
				// },
				grid: {
					left: '3%',
					right: '4%',
					bottom: '3%',
					containLabel: true
				},
				toolbox: {
					feature: {
						saveAsImage: {}
					}
				},
				xAxis: [
					{
						type: 'category',
						boundaryGap: false,
						data: firmware_versions
					}
				],
				yAxis: {
					type: 'value'
				},
				series: thread_function_stack_series
			});
		}
	};

	const updateSelectedVersion = () => {
		localStorage.selected_version = symbols.selected_version;
		updateSelectedSymbols();
	};

	onMount(async () => {
		if (browser) {
			// load elf data
			if (Object.keys(symbols.symbols).length == 0) {
				console.log('No ELF data URL passed or stored, please upload it as a file then :)');
			} else {
				if (symbols.selected_version && symbols.selected_versions_to_compare) {
					updateSelectedSymbols();
				} else {
					console.log('ELF loaded, please select which version to show :)');
				}
			}
		}
	});
</script>

<div class="container" id="content">
	<Row>
		<Col>
			Select a version of the .elf you want to see:
			<Input type="select" bind:value={symbols.selected_version} on:change={updateSelectedVersion}>
				{#each versions as version}
					<option>{version}</option>
				{/each}
			</Input>
		</Col>
	</Row>

	<hr />

	<Container fluid>
		{#if !symbols.selected_version}
			<h3>Select a version to browse thread info :)</h3>
		{:else}
			<h3>Stack data for {symbols.selected_version}</h3>

			<Row cols={{ lg: 3, md: 2, sm: 1 }} id="plotsContainer"></Row>

			{#key selected_thread_stat}
				<ul>
					{#each Object.keys(selected_thread_stat) as thread_name, index (thread_name + index)}
						<h4>{thread_name}</h4>

						<Table hover bordered>
							<thead>
								<tr>
									<th width="100%">Name</th>
									<th>Stack Size</th>
									<!-- <th colspan="2" style="text-align: center">Code</th> -->
								</tr>
								<!-- <tr>
                            <th class="col_size">
                                &sum; = {selected_thread_stat[thread_name].max_static_stack_size} / {selected_thread_stat[thread_name].max_stack_size}
                            </th>
                            <th class="col_size"></th>
                            <th class="col_size">&sum;= TODO not exported atm...</th>
                            <th class="col_size"></th>
                        </tr> -->
							</thead>
							<tbody>
								{#each selected_thread_stat[thread_name]['call_stack'] as fn, fn_index (thread_name + '_' + fn.name + '_' + index + '_' + fn_index)}
									<tr>
										<td>{fn.name}</td>
										<td>{fn.stack_size}</td>
									</tr>
								{/each}
							</tbody>
							<tfoot>
								<tr>
									<td><b>Total stack usage</b></td>
									<td
										><b>
											&sum; = {selected_thread_stat[thread_name].max_static_stack_size} / {selected_thread_stat[
												thread_name
											].max_stack_size}
										</b></td
									>
								</tr>
							</tfoot>
						</Table>

						<hr />
					{/each}
				</ul>
			{/key}
		{/if}
	</Container>
</div>

<style>
	/*
    @import 'static/css/style.css';
    */
	@import 'https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css';
</style>
