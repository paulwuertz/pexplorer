<script>
	import { onMount } from 'svelte';
	import { base } from '$app/paths';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';

	import { Col, Row, Button, ButtonGroup } from '@sveltestrap/sveltestrap';
	import * as echarts from 'echarts';
	import * as helpers from '../routes/helpers.js';

	const { symbols_in_symbols, treeMapRenderDepth } = $props();

	let symbols = $derived(symbols_in_symbols ? symbols_in_symbols : []);
	let render_depth = $derived(treeMapRenderDepth);
	let chartStyle = $state('treemap');
	let sunburst_chart_rom = $state(null);
	let sunburst_chart_ram = $state(null);
	let romChart = $state(null);
	let ramChart = $state(null);

	// TODO rm hack - think about how to distinguish better memory regions for all controller types...
	let rom_syms = $derived(
		symbols.filter(
			// TODO non bss section
			(e) =>
				!(
					e.address.toString(16).startsWith(2) && // not 0x2... in RAM
					e.address >= 0x20000000
				)
		)
	);
	let ram_syms = $derived(
		symbols.filter((e) => {
			return (
				e.address.toString(16).startsWith(2) && // 0x2... in RAM
				e.address >= 0x20000000
			);
		})
	);
	let rom_sunburst_data = $derived(helpers.symbols_to_sunburst_tree_data(rom_syms, 'size'));
	let ram_sunburst_data = $derived(helpers.symbols_to_sunburst_tree_data(ram_syms, 'size'));
	let updateCharts = () => {
		var option = {
			// visualMap: {
			//     type: 'continuous',
			//     inRange: {
			//         color: ['#2F93C8', '#AEC48F', '#FFDB5C', '#F98862']
			//     }
			// },
			series: {
				// color: [
				//     "#5070dd", "#b6d634", "#505372", "#ff994d", "#0ca8df", "#ffd10a",
				//     "#fb628b", "#785db0", "#3fbe95"
				// ],
				// colorMappingBy: "id",
				data: rom_sunburst_data,
				label: {
					show: false,
					// rotate: 'radial',
					// minAngle: 15,
					formatter: '{b} - {c} bytes'
				}
			},
			tooltip: {
				show: true,
				trigger: 'item'
			}
		};
		console.log(rom_sunburst_data);

		if (chartStyle == 'treemap') {
			option.series.type = 'treemap';
			option.series.roam = false;
			option.series.nodeClick = undefined;
			option.series.label.show = true;
			// TODO expose leafDepth as setting
			option.series.leafDepth = render_depth;
			// TODO add? visibleMin = 300;
		} else {
			option.series.type = 'sunburst';
			option.series.radius = [0, '90%'];
		}
		//
		option && romChart && romChart.setOption(option);

		let ram_options = JSON.parse(JSON.stringify(option));
		ram_options.series.data = ram_sunburst_data;
		ram_options && ramChart && ramChart.setOption(ram_options);
		return [ramChart, romChart];
	};
	let charts = $derived(() => {
		return updateCharts();
	});
	let btnStyle = (btn) => {
		if (btn == chartStyle) {
			return 'success';
		} else {
			return 'secondary';
		}
	};

	$effect(() => {
		romChart = echarts.init(sunburst_chart_rom);
		ramChart = echarts.init(sunburst_chart_ram);
		console.log('onmounty2');
		updateCharts();
	});
</script>

<Row cols={{ md: 2, sm: 1 }}>
	<div>
		<h3>Flash usage</h3>
		<div bind:this={sunburst_chart_rom} style="width: 100%;height:600px;"></div>
	</div>
	<div>
		<h3>Static RAM usage</h3>
		<div bind:this={sunburst_chart_ram} style="width: 100%;height:600px;"></div>
	</div>
</Row>

<ButtonGroup class="d-flex justify-content-center">
	<Button
		on:click={() => {
			chartStyle = 'treemap';
			updateCharts();
		}}
		color={btnStyle('treemap')}>Treemap</Button
	>
	<Button
		on:click={() => {
			chartStyle = 'sunburst';
			updateCharts();
		}}
		color={btnStyle('sunburst')}>Sunburst</Button
	>
	<!-- <Button on:click={()=> { chartStyle="flamegraph"; updateCharts()} } color={btnStyle("flamegraph")}>Flamegraph</Button> -->
</ButtonGroup>

<style>
</style>
