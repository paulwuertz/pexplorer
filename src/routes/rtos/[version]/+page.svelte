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
		ButtonGroup,
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
		Progress,
		CardFooter
	} from '@sveltestrap/sveltestrap';

	import { symbols } from '../../symbols.svelte.js';
	import * as helpers from '../../helpers.js';
	import * as echarts from 'echarts';
	import ThreadInfoCard from '../../../components/ThreadInfoCard.svelte';

	let local_storage_key = 'pexplorer_settings';
	let DUMMY_STANDARD_SETTINGS_NAME = 'test';
	let active_settings = DUMMY_STANDARD_SETTINGS_NAME;

	let params = $props();
	// TODO - why are propsed nested +1 here on production build?!
	let parameters = $derived(params && (params.data.data || params.data));
	let version_name = $derived((parameters && parameters.version) || '');
	let version = $derived(symbols.symbols[version_name] || {});
	let functions = $derived(version.functions || []);
	let variables = $derived(version.variables || []);

	let typeStruct = $derived(
		variableTypes.find((val, i, arr) => {
			return val.name == '_static_thread_data';
		})
	);

	let restore_default_settings = () => {
		let stored_settings_str = localStorage.getItem(local_storage_key);
		let no_settings_backed_up = !stored_settings_str;
		if (no_settings_backed_up) {
			stored_settings_str = store_default_settings();
		}
		console.log('stored_settings_str', stored_settings_str);
		let stored_settings = JSON.parse(stored_settings_str);
		return stored_settings;
	};

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

	let staticThreads = $derived(variables.filter(isStaticThread));
	let variableTypes = $derived(version['types']);
	let manuallyEntered = $state(restore_default_settings());
	let manuallyEnteredThreads = $derived(manuallyEntered['test']['threads']);
	let manuallyEnteredThreadInfo = $derived(
		manuallyEnteredThreads
			.map((t, i, a) => {
				let entry_function = functions.find((val, i, arr) => {
					return val['name'] == t['thread_entry_name'];
				});
				if (entry_function) {
					let thread_data = {};
					thread_data['name'] = entry_function['name'];
					thread_data['file'] = entry_function['file'];
					thread_data['init_stack_size'] = helpers.stored_thread_settings_stack_size(t, variables);
					thread_data['max_stack_size_callees'] = entry_function['max_stack_size_callees'];
					thread_data['entry_function'] = entry_function;
					let fn_calltree = JSON.parse(get_fn_calltree(entry_function.address));
					if (Object.hasOwn(fn_calltree, 'unresolved')) {
						thread_data['unresolved_calls'] = fn_calltree.unresolved.length;
						// TODO from-to+dynamic is anoying...
						// if elf is the central format maby it does not matter to much
						// but for diff/comparing call names would be nice, but also
						// could take more memory - anyway think about extending the type...
						let function_unresolved_calls_from = {};
						(fn_calltree.unresolved || []).forEach((call) => {
							function_unresolved_calls_from[call.from] = 'just counting :)';
						});
						thread_data['from_nr_functions'] = Object.keys(function_unresolved_calls_from).length;
					} else {
						thread_data['unresolved_calls'] = 0;
					}
					thread_data['fn_calltree'] = fn_calltree;
					return thread_data;
				}
				return null;
			})
			.filter((v) => v != null)
	);
	$inspect(manuallyEnteredThreadInfo);
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

			let entry_function = functions.find((val, i, arr) => {
				return val['address'] == thread_data['init_entry'];
			});
			if (entry_function) {
				thread_data['name'] = entry_function['name'];
				thread_data['file'] = entry_function['file'];
				thread_data['entry_function'] = entry_function;
				if (typeof get_fn_calltree === 'function') {
					let fn_calltree = JSON.parse(get_fn_calltree(entry_function.address));
					if (Object.hasOwn(fn_calltree, 'unresolved')) {
						thread_data['unresolved_calls'] = fn_calltree.unresolved.length;
						// TODO from-to+dynamic is anoying...
						// if elf is the central format maby it does not matter to much
						// but for diff/comparing call names would be nice, but also
						// could take more memory - anyway think about extending the type...
						let function_unresolved_calls_from = {};
						(fn_calltree.unresolved || []).forEach((call) => {
							function_unresolved_calls_from[call.from] = 'just counting :)';
						});
						thread_data['from_nr_functions'] = Object.keys(function_unresolved_calls_from).length;
					} else {
						thread_data['unresolved_calls'] = 0;
					}
					thread_data['fn_calltree'] = fn_calltree;
				}
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

	const isUnassociatedStaticStack = (symbol) => {
		let sym_type = symbol.type || '';
		let sym_addr = symbol.address || '';
		let alreadyAssociated = static_thread_data.find((element) => element.init_stack == sym_addr);
		return sym_type == 'z_thread_stack_element' && !alreadyAssociated;
	};
	let staticStacks = $derived(variables.filter(isUnassociatedStaticStack));
	$inspect(version_name, static_thread_data);
</script>

<div class="container" id="content">
	<Container fluid>
		<h3>version_name: {JSON.stringify(version_name)}</h3>
		<h4>Static threads by K_THREAD_DEFINE</h4>
		<Row cols={{ lg: 3, md: 2, sm: 1 }}>
			{#each static_thread_data as sTread (sTread.name)}
				<div class="pb-3 px-3">
					<ThreadInfoCard {sTread} {version_name}></ThreadInfoCard>
				</div>
			{/each}
		</Row>
		<h4>Manually configured threads:</h4>

		<Row cols={{ lg: 3, md: 2, sm: 1 }}>
			{#each manuallyEnteredThreadInfo as sTread (sTread.name)}
				<div class="pb-3 px-3">
					<ThreadInfoCard {sTread} {version_name}></ThreadInfoCard>
				</div>
			{/each}
		</Row>

		<h4>Detected static stacks without an associated thread:</h4>

		<ul>
			{#each staticStacks as sStacks (sStacks.name)}
				<li>
					{sStacks.name} - size: {sStacks.size} bytes - addr: 0x{sStacks.address.toString(16)}
				</li>
			{/each}
		</ul>
	</Container>
</div>
