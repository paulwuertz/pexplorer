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
		Progress
	} from '@sveltestrap/sveltestrap';

	import { symbols } from '../../symbols.svelte.js';
	import * as helpers from '../../helpers.js';
	import VersionSelect from '../../../components/VersionSelect.svelte';

	import * as echarts from 'echarts';

	const PLOT_ID_PREFIX = 'stackSizeLayered_';

	let { params, data } = $props();
	let version_name = $derived(params.version || '');
	let version = $derived(symbols.symbols[version_name] || {});
	let functions = $derived(version.functions || []);
	let variables = $derived(version.variables || []);

	let typeStruct = $derived(
		variableTypes.find((val, i, arr) => {
			return val.name == '_static_thread_data';
		})
	);

	const isStaticThread = (symbol) => {
		let secidx = symbol.secidx || null;
		let sections = version.sections || [];
		let section = sections[secidx] || {};
		let section_name = section.name || '';
		return section_name == '_static_thread_data_area';
	};

	// TODO mv to helpers
	let getMemberBytes = (data, byte_offset, size) => {
		let b = [];
		let val = 0;
		for (let i = byte_offset; i < data.length && i < byte_offset + size; i++) {
			b.push(data[i]);
			val += data[i] * Math.pow(256, i - byte_offset);
		}
		return val;
	};

	let stackLevelToColor = (stackUse, stackSize) => {
		let percentage = (100.0 * stackUse) / stackSize;
		if (percentage < 50.0) {
			return '';
		}
		if (percentage < 80.0) {
			return 'success';
		}
		if (percentage < 100.0) {
			return 'warning';
		} else {
			return 'danger';
		}
	};

	let staticThreads = $derived(variables.filter(isStaticThread));
	let variableTypes = $derived(version['types']);
	let static_thread_args = $derived((typeStruct && typeStruct.members) || []);
	let static_thread_data = $derived(
		staticThreads.map((t, i, a) => {
			let staticInitDataBase64 = t.staticInitData;
			let staticInitData =
				(staticInitDataBase64 && Uint8Array.fromBase64(staticInitDataBase64)) || undefined;
			let thread_data = {};
			static_thread_args.forEach((member) => {
				thread_data[member.name] = getMemberBytes(staticInitData, member.byte_offset, member.size);
			});
			thread_data['name'] = t['name'];

			let entry_function = functions.find((val, i, arr) => {
				return val['address'] == thread_data['init_entry'];
			});
			if (entry_function) {
				thread_data['entry_function'] = entry_function;
				thread_data['max_stack_size_callees'] = entry_function['max_stack_size_callees'];
			}
			let thread_name = variables.find((val, i, arr) => {
				return val['address'] == thread_data['init_name'];
			});
			if (thread_name) {
				thread_data['thread_name'] = entry_function;
			}
			return thread_data;
		})
	);
	$inspect(static_thread_data);
</script>

<div class="container" id="content">
	{#if !version || version_name == ''}
		<VersionSelect page_name={'rtos'}></VersionSelect>
	{:else}
		<Container fluid>
			<h3>version_name: {JSON.stringify(version_name)}</h3>
			<h4>Static threads by K_THREAD_DEFINE</h4>
			<Row cols={{ lg: 3, md: 2, sm: 1 }}>
				{#each static_thread_data as sTread (sTread.name)}
					<div class="pb-3 px-3">
						<Card>
							<CardHeader>
								<CardTitle>{sTread['name']}</CardTitle>
							</CardHeader>
							<CardBody>
								<CardText>
									<!-- {#each Object.entries(sTread) as s}
                                        <div>{s[0]} - {s[1]}</div>
                                        {"init_thread":536872864,"init_stack":536877904,"init_stack_size":256,"init_entry":134227277,"init_p1":0,"init_p2":0,"init_p3":0,"init_prio":4,"init_options":0,"init_name":134288890,"init_delay":0,"name":"_k_thread_data_leds"}
                                    {/each} -->
									<!-- TODO add source link symbol json -->
									<!-- Buildtime: {symbols.symbols[version].timestamp} -->
									<div class="pb-3">
										<b>Thread entry function:</b>
										<br />
										<span>
											<a data-sveltekit-preload-data="tap" href={'#'}> _k_thread_data_leds </a>
										</span>
									</div>
									<div class="pb-3">
										<b>Configured stack size:</b>
										<br />
										<span>init_stack_size bytes</span>
									</div>
									<div class="pb-3">
										<b>Unresolved dynamic calls:</b>
										<br />
										<span>TODO</span>
									</div>
									<div class="pb-3">
										<b>Functions missing stack-use info:</b>
										<br />
										<span>TODO</span>
									</div>

									<CardSubtitle>Minimum stack use scenario found:</CardSubtitle>
									<div class="pt-3">
										<Progress
											color={stackLevelToColor(
												sTread['max_stack_size_callees'],
												sTread['init_stack_size']
											)}
											value={sTread['max_stack_size_callees']}
											max={sTread['init_stack_size']}
											class="mb-2"
										>
											{(100 * sTread['max_stack_size_callees']) / sTread['init_stack_size']}% - {sTread[
												'max_stack_size_callees'
											]} / {sTread['init_stack_size']} bytes)
										</Progress>
									</div>
								</CardText>
							</CardBody>
						</Card>
					</div>
				{/each}
			</Row>
		</Container>
	{/if}
</div>
