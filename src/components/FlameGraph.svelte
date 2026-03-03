<script>
	import { onMount } from 'svelte';
	import { base } from '$app/paths';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';

	import * as echarts from 'echarts';
	import * as helpers from '../routes/helpers.js';

	var ROOT_PATH = base + '/t.json';

	var option;

	const { sym_data, fn_calltree } = $props();

	onMount(async () => {
		var chartDom = document.getElementById('main');
		var myChart = echarts.init(chartDom);
		// const ColorTypes = {
		//     root: '#8fd3e8',
		//     genunix: '#d95850',
		//     unix: '#eb8146',
		//     ufs: '#ffb248',
		//     FSS: '#f2d643',
		//     namefs: '#ebdba4',
		//     doorfs: '#fcce10',
		//     lofs: '#b5c334',
		//     zfs: '#1bca93'
		// };
		const filterJson = (json, id) => {
			if (id == null) {
				return json;
			}
			const recur = (item, id) => {
				if (item.name === id) {
					return item;
				}
				for (const child of item.calls || []) {
					const temp = recur(child, id);
					if (temp) {
						item.calls = [temp];
						item.max_stack_size_callees = temp.max_stack_size_callees; // change the parents' values
						return item;
					}
				}
			};
			return recur(json, id) || json;
		};
		const recursionJson = (jsonObj, id) => {
			const data = [];
			const filteredJson = filterJson(structuredClone(jsonObj), id);
			const rootCalls = filteredJson.calls;
			const rootVal = filteredJson.max_stack_size_callees;
			const recur = (item, start = 0, level = 0, factor) => {
				let max_calle_stack_size = item.max_stack_size_callees * factor;
				const temp = {
					name: item.name,
					// [level, start_val, end_val, name, stacksize, max_callee_stack]
					value: [
						level,
						start,
						start + max_calle_stack_size,
						item.name,
						item.stack_size,
						item.max_stack_size_callees
					],
					itemStyle: {
						// TODO make recursive calls red :)
						// color: ColorTypes[item.name.split(' ')[0]]
					}
				};

				data.push(temp);
				let prevStart = start;
				let childSum = item.calls
					? item.calls.reduce((c, e) => c + (e.max_stack_size_callees || 0), 0)
					: 0;
				let width_left_to_children = item.max_stack_size_callees - item.stack_size;
				console.log(
					level,
					start,
					' lev}\n\t',
					temp.value,
					item.max_stack_size_callees,
					factor,
					childSum
				);
				for (const child of item.calls || []) {
					if (child.max_stack_size_callees) {
						let childFactor = (child.max_stack_size_callees || 0) / childSum;
						recur(child, prevStart, level + 1, factor * childFactor);
						prevStart = prevStart + width_left_to_children * childFactor;
					}
				}
			};
			recur(filteredJson, 0, 0, 1.0);
			return data;
		};
		const heightOfJson = (json) => {
			const recur = (item, level = 0) => {
				if ((item.calls || []).length === 0) {
					return level;
				}
				let maxLevel = level;
				for (const child of item.calls) {
					const tempLevel = recur(child, level + 1);
					maxLevel = Math.max(maxLevel, tempLevel);
				}
				return maxLevel;
			};
			return recur(json);
		};
		const renderItem = (params, api) => {
			const level = api.value(0);
			const start = api.coord([api.value(1), level]);
			const end = api.coord([api.value(2), level]);
			const height = ((api.size && api.size([0, 1])) || [0, 20])[1];
			// [level, start_val, end_val, name, stacksize, max_callee_stack]
			const width = end[0] - start[0];
			return {
				type: 'rect',
				transition: ['shape'],
				shape: {
					x: start[0],
					y: start[1] - height / 2,
					width,
					height: height - 2 /* itemGap */,
					r: 2
				},
				style: { fill: api.visual('color') },
				emphasis: { style: { stroke: '#000' } },
				textConfig: {
					position: 'insideLeft'
				},
				textContent: {
					style: {
						text: api.value(3),
						fontFamily: 'Verdana',
						fill: '#000',
						width: width - 4,
						overflow: 'truncate',
						ellipsis: '..',
						truncateMinChar: 1
					},
					emphasis: {
						style: {
							stroke: '#000',
							lineWidth: 0.5
						}
					}
				}
			};
		};
		myChart.showLoading();
		let stackTrace = fn_calltree; //await fetch(ROOT_PATH).then((resp) => resp.json())
		myChart.hideLoading();
		console.log(stackTrace);

		const levelOfOriginalJson = heightOfJson(stackTrace.tree);
		option = {
			backgroundColor: {
				type: 'linear',
				x: 0,
				y: 0,
				x2: 0,
				y2: 1,
				colorStops: [
					{ offset: 0.05, color: '#eee' },
					{ offset: 0.95, color: '#eeeeb0' }
				]
			},
			tooltip: {
				formatter: (params) => {
					const wcf = params.value[5];
					return `${params.marker} ${params.value[3]}: (${params.value[4]} bytes of stack use itself and ${wcf} bytes worst case call stack usage)`;
				}
			},
			title: [
				{
					text: sym_data.name,
					left: 'center',
					top: 10,
					textStyle: {
						fontFamily: 'Verdana',
						fontWeight: 'normal',
						fontSize: 20
					}
				}
			],
			toolbox: {
				feature: { restore: {} },
				right: 20,
				top: 10
			},
			xAxis: { show: false },
			yAxis: {
				show: false,
				max: levelOfOriginalJson
			},
			series: [
				{
					type: 'custom',
					renderItem,
					encode: {
						x: [0, 1, 2],
						y: 0
					},
					data: recursionJson(stackTrace.tree)
				}
			]
		};
		myChart.setOption(option);
		myChart.on('click', (params) => {
			const data = recursionJson(stackTrace.tree, params.data.name);
			const rootValue = data[0].value[2];
			myChart.setOption({
				xAxis: { max: rootValue },
				series: [{ data }]
			});
		});
	});
</script>

<div id="main" style="width: 100%;height:600px;"></div>
