<script>
	import { onMount } from 'svelte';
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import { browser } from '$app/environment';
	import { writable } from 'svelte/store';
	// ui stuff
	import { DataTable } from '@careswitch/svelte-data-table';
	import { Badge, Button, Col, Container, Input, Row, Table } from '@sveltestrap/sveltestrap';

	import { symbols } from '../symbols.svelte.js';

	import * as echarts from 'echarts';
	import * as helpers from '../helpers.js';

	let { data } = $props();
	let files = $state();
	let versions = $derived(Object.keys(symbols.symbols));
	let selected_symbols = $state({});

	onMount(async () => {
		if (browser) {
			let trend_data = versions;
			let threads = helpers.get_all_threads(symbols.symbols);
			let ordered_versions_and_timestamps = helpers.get_versions_ordered_by_timestamps(
				symbols.symbols
			);
			let ordered_versions = ordered_versions_and_timestamps.map((vt) => vt['version']);
			let ordered_timestamps = ordered_versions_and_timestamps.map((vt) =>
				vt['timestamp'].getTime()
			);

			let series_maxstack_data = [];
			for (let thread of threads) {
				let maxstack_data = [];
				for (let [i, v] of ordered_versions.entries()) {
					let thread_stack = symbols.symbols[v]['stack_reports'][thread];
					if (!(thread_stack && thread_stack.hasOwnProperty('max_static_stack_size'))) continue;
					let max_stack_size = thread_stack['max_static_stack_size'];
					maxstack_data.push([ordered_timestamps[i], max_stack_size, ordered_versions[i]]);
				}
				series_maxstack_data.push({
					name: thread,
					type: 'line',
					data: maxstack_data
				});
			}
			// Create the echarts instance
			var maxStackSizeChart = echarts.init(document.getElementById('maxStackSizeChartID'), {
				width: '95%',
				height: 600
			});

			// Draw the chart
			maxStackSizeChart.setOption({
				title: {
					text: 'Static max stack sizes'
				},
				tooltip: {
					trigger: 'axis'
				},
				legend: {
					data: threads
				},
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
				xAxis: {
					type: 'time'
				},
				yAxis: {
					type: 'value'
				},
				series: series_maxstack_data
			});

			let alloc_call_data = [];
			const allocating_functions = [
				'operator new(unsigned int)',
				'operator new[](unsigned int)',
				'malloc',
				'calloc'
			];
			let alloc_map_f = {};
			for (let sym of allocating_functions) {
				alloc_map_f[sym] = [];
			}
			for (let v of ordered_versions) {
				let vdata = symbols.symbols[v];
				let vsymbols = vdata['symbols'];
				let vtimestamp = vdata['timestamp'];
				for (let sym of vsymbols) {
					if (!sym.hasOwnProperty('display_name')) continue;
					let sym_name = sym['display_name'];
					if (allocating_functions.includes(sym['display_name'])) {
						if (!sym.hasOwnProperty('callers')) continue;
						let callers_str = sym['callers'];
						let callers = JSON.parse(callers_str);
						alloc_map_f[sym_name].push([vtimestamp, callers.length]);
					}
				}
			}
			for (let sym of allocating_functions) {
				let nr_calls = alloc_map_f.hasOwnProperty(sym) ? alloc_map_f[sym] : 0;
				alloc_call_data.push({
					name: sym,
					type: 'line',
					stack: 'Total',
					data: nr_calls
				});
			}
			console.log('alloc_call_data ', alloc_call_data, alloc_map_f);
			// Create the echarts instance
			var maxStackSizeChart = echarts.init(document.getElementById('nrDynamicAllocs'), {
				width: '95%',
				height: 600
			});

			// Draw the chart
			maxStackSizeChart.setOption({
				title: {
					text: 'Number of dynamic allocations'
				},
				tooltip: {
					trigger: 'axis'
				},
				legend: {
					data: threads
				},
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
				xAxis: {
					type: 'time'
				},
				yAxis: {
					type: 'value'
				},
				series: alloc_call_data
			});
		}
	});
</script>

<div class="container" id="content">
	<div id="maxStackSizeChartID" style="width: 100%;height:600px;"></div>

	<hr />

	<div id="nrDynamicAllocs" style="width: 100%;height:600px;"></div>
</div>

<style>
	/*
    @import 'static/css/style.css';
    */
	@import 'https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css';
</style>
