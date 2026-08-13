<script>
	import { onMount } from 'svelte';
	import { base } from '$app/paths';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';

	import * as echarts from 'echarts';
	import * as helpers from '../routes/helpers.js';

	var option;

	const { sym_data, fn_calltree } = $props();
	let data = $derived(helpers.symbols_to_call_tree_data(fn_calltree.tree));
	onMount(async () => {
		var chartDom = document.getElementById('calltree');
		var myChart = echarts.init(chartDom);
		let option = {
			tooltip: {
				trigger: 'item',
				triggerOn: 'mousemove'
			},
			series: [
				{
					type: 'tree',
					id: 0,
					name: 'tree1',
					data: [data],
					top: '10%',
					left: '8%',
					bottom: '22%',
					right: '20%',
					symbolSize: 7,
					edgeShape: 'polyline',
					edgeForkPosition: '63%',
					initialTreeDepth: 1,
					lineStyle: {
						width: 2
					},
					label: {
						backgroundColor: '#fff',
						position: 'left',
						verticalAlign: 'middle',
						align: 'right'
					},
					leaves: {
						label: {
							position: 'right',
							verticalAlign: 'middle',
							align: 'left'
						}
					},
					emphasis: {
						focus: 'descendant'
					},
					expandAndCollapse: true,
					animationDuration: 550,
					animationDurationUpdate: 750
				}
			]
		};
		myChart.setOption(option);
	});
</script>

<div id="calltree" style="width: 100%;height:600px;"></div>
