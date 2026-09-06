<script>
	import { onMount } from 'svelte';
	import { base } from '$app/paths';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';

	import {
		Card,
		Button,
		ButtonGroup,
		Col,
		CardHeader,
		Container,
		Icon,
		Input,
		Row,
		CardBody,
		CardText,
		CardTitle,
		Table
	} from '@sveltestrap/sveltestrap';
	import Select from 'svelte-select';
	import { symbols } from '../routes/symbols.svelte.js';
	import * as helpers from '../routes/helpers.js';
	import Dropzone from 'svelte-file-dropzone';

	const { settings } = $props();

	const itemId = 'address';
	const label = 'name';
	let versions = $derived(Object.keys(symbols.symbols));
	let version_name = $state(versions[0]);
	let restored_settings = $state();
	let version = $derived(symbols.symbols[version_name] || {});
	let setting = $derived((restored_settings && restored_settings[version.firmware_hash]) || {});
	let threads = $derived((setting && setting.threads) || []);
	let dynamic_calls = $derived((setting && setting.dynamic_calls) || []);
	// $inspect(setting, threads, dynamic_calls);
	let functions = $derived(version.functions || []);
	let variables = $derived(version.variables || []);
	//
	let functions_missing_calls = $derived(
		functions.filter((f, i) => {
			if (!Object.hasOwn(f, 'callees')) return false;
			for (const callee of f['callees'] || []) {
				if (!Object.hasOwn(callee, 'to')) {
					return true;
				}
			}
			return false;
		})
	);
	let count_fns = () => {
		let nr_fns_with_indirect_calls = 0;
		for (const f of functions_missing_calls) {
			for (const callee of f['callees'] || []) {
				if (!Object.hasOwn(callee, 'to')) {
					nr_fns_with_indirect_calls++;
					break;
				}
			}
		}
		return nr_fns_with_indirect_calls;
	};
	let functions_missing_calls_nr = $derived(count_fns());
	let missing_calls_nr = $derived(
		Object.fromEntries(
			functions_missing_calls.map((f, i) => {
				let nr_indirect_calls = 0;
				for (const callee of f['callees'] || []) {
					if (!Object.hasOwn(callee, 'to')) {
						nr_indirect_calls++;
					}
				}
				return [f['name'], nr_indirect_calls];
			})
		)
	);
	let total_functions_missing_calls_nr = $derived(
		Object.values(missing_calls_nr).reduce((acc, val) => acc + val, 0)
	);

	let selected_call_from = $state();
	let selected_call_to = $state();
	let link_caller_and_callee = (caller, callee) => {
		//console.log(caller, 'bef');

		for (let i = 0; i < len(caller['callees']); i++) {
			if (!Object.hasOwn(caller['callees'][i], 'to')) {
				caller['callees'][i]['to'] = callee['address'];
				caller['callees'][i]['to_function_name'] = callee['name'];
				break;
			}
		}

		//console.log(caller, 'af');
	};
	let add_dynamic_call = () => {
		if (!selected_call_from) {
			// TODO err
			alert('Add missing selected_call_from');
			return;
		}
		if (!selected_call_to) {
			// TODO err
			alert('Add missing selected_call_to');
			return;
		}
		dynamic_calls.push({
			call_from: selected_call_from.name,
			call_to: selected_call_to.name
		});
		generate_and_store_new_setting();
		link_caller_and_callee(selected_call_from, selected_call_to);
		selected_call_to = null; //  easier allows adding new call
	};

	let restore_active_settings = () => {
		let all_settings = helpers.restore_default_settings(version.firmware_hash);
		restored_settings = all_settings;
		if (Object.hasOwn(all_settings, version.firmware_hash)) {
			return all_settings[version.firmware_hash];
		} else {
			return {};
		}
	};

	let download_template_settings = () => {
		let all_settings = helpers.restore_default_settings(version.firmware_hash);
		restored_settings = all_settings;
		if (Object.hasOwn(all_settings, version.firmware_hash)) {
			let settings = all_settings[version.firmware_hash];
			settings['dynamic_calls'] = puncover_158_indirect_calls();
			return settings;
		} else {
			return {};
		}
	};

	onMount(() => {
		setting = restore_active_settings();
	});

	function handleFilesSelect(e) {
		const { acceptedFiles, fileRejections } = e.detail;
		console.log(acceptedFiles);

		for (let i = 0; i < acceptedFiles.length; i++) {
			let file = acceptedFiles[i];
			console.log(file);
			const reader = new FileReader();
			// Validate file existence and type
			if (!file) {
				console.log('No file selected. Please choose a file.', 'error');
				return;
			}

			let file_type = file.type.toLocaleLowerCase();
			let is_json = file_type.endsWith('json');

			if (!is_json) {
				console.log(
					file.type +
						'Unsupported file type ' +
						"'" +
						file.type +
						"'" +
						' - please select a text file.',
					'error'
				);
				return;
			}
			reader.onload = () => {
				console.log('loded', reader);
				if (is_json) {
					//TODO verify to schema
					//TODO support multiple schemas
					let settings = JSON.parse(reader.result);
					dynamic_calls = settings.dynamic_calls;
					threads = settings.threads;
					backup_settings(version.firmware_hash, {
						threads: threads,
						dynamic_calls: dynamic_calls
					});
				} else {
					//TODO
					console.log('TODO');
				}
			};
			reader.onerror = () => {
				alert('Error reading the file. Please try again.');
			};

			reader.readAsText(file);
		}
	}

	let selected_thread_entry = $state();
	let selected_stack_variable = $state();
	let selected_stack_size = $state();
	let add_thread = () => {
		if (!selected_thread_entry) {
			// TODO err
			alert('Add missing selected_thread_entry');
			return;
		}
		if (selected_stack_variable) {
			// "stack_variable_name" "size"
			threads.push({
				thread_entry_name: selected_thread_entry.name,
				stack_variable_name: selected_stack_variable.name
			});
			generate_and_store_new_setting();
		} else if (selected_stack_size) {
			threads.push({
				thread_entry_name: selected_thread_entry.name,
				size: selected_stack_size.name
			});
			generate_and_store_new_setting();
		} else {
			alert('Add missing selected_stack_variable');
			return;
		}
	};

	let backup_settings = (settings_name, new_setting_set) => {
		let stored_settings = helpers.restore_default_settings(version.firmware_hash);
		// TODO protect overwrites?
		// if (Object.hasOwn(stored_settings, settings_name)) {
		stored_settings[settings_name] = new_setting_set;
		// }
		localStorage.setItem(helpers.local_storage_key, JSON.stringify(stored_settings));
		restored_settings = stored_settings;
	};

	let generate_and_store_new_setting = () => {
		let new_setting = {
			threads: threads,
			dynamic_calls: dynamic_calls
		};
		backup_settings(version.firmware_hash, new_setting);
	};

	let puncover_158_indirect_calls = () => {
		let active_settings = restore_active_settings();
		active_settings['dynamic_calls'] = helpers.flat_calls_to_arrayed_callees(active_settings);
		return active_settings;
	};

	let download_puncover_158_indirect_calls = () => {
		let calls = {};
		let active_settings = restore_active_settings();
		// build a map
		for (const dynamic_call of active_settings['dynamic_calls']) {
			let call_from = dynamic_call['call_from'];
			let call_to = dynamic_call['call_to'];
			if (Object.hasOwn(calls, call_from)) {
				calls[call_from].push(call_to);
			} else {
				calls[call_from] = [call_to];
			}
		}
		// map to array
		let calls_arr = [];
		for (const dynamic_caller in calls) {
			calls_arr.push({
				caller: dynamic_caller,
				callees: calls[dynamic_caller]
			});
		}
		let indirect_calls = {
			version: 1,
			indirect_callees: puncover_158_indirect_calls()
		};
		helpers.download(
			'puncover-dynamic-calls-' + version_name + '.json',
			JSON.stringify(indirect_calls, 0, 4)
		);
	};
</script>

<div class="container" id="content">
	<h3>pexplorer setting:</h3>

	<div class="pt-3">Upload a config file or download the current:</div>

	<Row>
		<Col>
			<div class="pb-3 pt-3">
				<Dropzone on:drop={handleFilesSelect} accept=".json">
					<p>Drag 'n' drop JSON config here, or click to select files.</p>
				</Dropzone>
			</div>
		</Col>
		<Col>
			<Button
				color="primary"
				on:click={() =>
					helpers.download(
						'pexplorer-' + version_name + '.json',
						JSON.stringify(download_template_settings(), 0, 4)
					)}
			>
				Download settings
			</Button>

			<br />

			<ButtonGroup>
				<Button
					class="mt-3"
					color="primary"
					on:click={() => download_puncover_158_indirect_calls()}
				>
					Download settings for puncover #158
				</Button>
				<Button
					class="mt-3"
					color="warning"
					on:click={() => download_template_puncover_158_indirect_calls()}
				>
					Download dynamic call template for puncover #158
				</Button>
			</ButtonGroup>
		</Col>
	</Row>

	Select which firmware to configure:

	<Input class="mb-3 mt-3" type="select" bind:value={version_name}>
		{#each versions as option}
			<option>{option}</option>
		{/each}
	</Input>

	<h4>Threads:</h4>

	<Table bordered>
		<thead>
			<tr>
				<th>Thread entry function name</th>
				<th>Thread stack variable name</th>
				<th>Thread stack size</th>
				<th></th>
			</tr>
		</thead>
		<tbody>
			{#each threads as thread, index (thread.thread_entry_name + index)}
				<tr>
					<td>{thread.thread_entry_name}</td>
					<td>{thread.stack_variable_name}</td>
					<td>{helpers.stored_thread_settings_stack_size(thread, variables)}</td>
					<td>
						<Button
							onclick={() => {
								threads = threads.filter((v, i) => i != index);
								generate_and_store_new_setting();
							}}
							size="md"
							outline
							color="danger"
						>
							<svg height="32" width="29" viewBox="0 0 1000 875"
								><path
									d="M0 281.296l0 -68.355q1.953 -37.107 29.295 -62.496t64.449 -25.389l93.744 0l0 -31.248q0 -39.06 27.342 -66.402t66.402 -27.342l312.48 0q39.06 0 66.402 27.342t27.342 66.402l0 31.248l93.744 0q37.107 0 64.449 25.389t29.295 62.496l0 68.355q0 25.389 -18.553 43.943t-43.943 18.553l0 531.216q0 52.731 -36.13 88.862t-88.862 36.13l-499.968 0q-52.731 0 -88.862 -36.13t-36.13 -88.862l0 -531.216q-25.389 0 -43.943 -18.553t-18.553 -43.943zm62.496 0l749.952 0l0 -62.496q0 -13.671 -8.789 -22.46t-22.46 -8.789l-687.456 0q-13.671 0 -22.46 8.789t-8.789 22.46l0 62.496zm62.496 593.712q0 25.389 18.553 43.943t43.943 18.553l499.968 0q25.389 0 43.943 -18.553t18.553 -43.943l0 -531.216l-624.96 0l0 531.216zm62.496 -31.248l0 -406.224q0 -13.671 8.789 -22.46t22.46 -8.789l62.496 0q13.671 0 22.46 8.789t8.789 22.46l0 406.224q0 13.671 -8.789 22.46t-22.46 8.789l-62.496 0q-13.671 0 -22.46 -8.789t-8.789 -22.46zm31.248 0l62.496 0l0 -406.224l-62.496 0l0 406.224zm31.248 -718.704l374.976 0l0 -31.248q0 -13.671 -8.789 -22.46t-22.46 -8.789l-312.48 0q-13.671 0 -22.46 8.789t-8.789 22.46l0 31.248zm124.992 718.704l0 -406.224q0 -13.671 8.789 -22.46t22.46 -8.789l62.496 0q13.671 0 22.46 8.789t8.789 22.46l0 406.224q0 13.671 -8.789 22.46t-22.46 8.789l-62.496 0q-13.671 0 -22.46 -8.789t-8.789 -22.46zm31.248 0l62.496 0l0 -406.224l-62.496 0l0 406.224zm156.24 0l0 -406.224q0 -13.671 8.789 -22.46t22.46 -8.789l62.496 0q13.671 0 22.46 8.789t8.789 22.46l0 406.224q0 13.671 -8.789 22.46t-22.46 8.789l-62.496 0q-13.671 0 -22.46 -8.789t-8.789 -22.46zm31.248 0l62.496 0l0 -406.224l-62.496 0l0 406.224z"
								/></svg
							>
						</Button>
					</td>
				</tr>
			{/each}

			<tr>
				<td>
					<Select {itemId} {label} items={functions} bind:value={selected_thread_entry}></Select>
				</td>
				<td>
					<Select {itemId} {label} items={variables} bind:value={selected_stack_variable}></Select>
				</td>
				<td>
					<Input type="number" bind:value={selected_stack_size} />
				</td>
				<td>
					<Button color="success" on:click={add_thread}>Add thread</Button>
				</td>
			</tr>
		</tbody>
	</Table>

	<h4>
		Resolved dynamic calls (from {functions_missing_calls_nr} functions still {total_functions_missing_calls_nr}
		calls unresolved):
	</h4>

	<Table bordered>
		<thead>
			<tr>
				<th>From</th>
				<th>To</th>
				<th></th>
			</tr>
		</thead>
		<tbody>
			{#each dynamic_calls as dynamic_call, index (dynamic_call.call_from + dynamic_call.call_to + index)}
				<tr>
					<td>{dynamic_call.call_from}</td>
					<td>{dynamic_call.call_to}</td>
					<td>
						<Button
							onclick={() => {
								dynamic_calls = dynamic_calls.filter((v, i) => i != index);
								generate_and_store_new_setting();
							}}
							size="md"
							outline
							color="danger"
						>
							<svg height="32" width="29" viewBox="0 0 1000 875"
								><path
									d="M0 281.296l0 -68.355q1.953 -37.107 29.295 -62.496t64.449 -25.389l93.744 0l0 -31.248q0 -39.06 27.342 -66.402t66.402 -27.342l312.48 0q39.06 0 66.402 27.342t27.342 66.402l0 31.248l93.744 0q37.107 0 64.449 25.389t29.295 62.496l0 68.355q0 25.389 -18.553 43.943t-43.943 18.553l0 531.216q0 52.731 -36.13 88.862t-88.862 36.13l-499.968 0q-52.731 0 -88.862 -36.13t-36.13 -88.862l0 -531.216q-25.389 0 -43.943 -18.553t-18.553 -43.943zm62.496 0l749.952 0l0 -62.496q0 -13.671 -8.789 -22.46t-22.46 -8.789l-687.456 0q-13.671 0 -22.46 8.789t-8.789 22.46l0 62.496zm62.496 593.712q0 25.389 18.553 43.943t43.943 18.553l499.968 0q25.389 0 43.943 -18.553t18.553 -43.943l0 -531.216l-624.96 0l0 531.216zm62.496 -31.248l0 -406.224q0 -13.671 8.789 -22.46t22.46 -8.789l62.496 0q13.671 0 22.46 8.789t8.789 22.46l0 406.224q0 13.671 -8.789 22.46t-22.46 8.789l-62.496 0q-13.671 0 -22.46 -8.789t-8.789 -22.46zm31.248 0l62.496 0l0 -406.224l-62.496 0l0 406.224zm31.248 -718.704l374.976 0l0 -31.248q0 -13.671 -8.789 -22.46t-22.46 -8.789l-312.48 0q-13.671 0 -22.46 8.789t-8.789 22.46l0 31.248zm124.992 718.704l0 -406.224q0 -13.671 8.789 -22.46t22.46 -8.789l62.496 0q13.671 0 22.46 8.789t8.789 22.46l0 406.224q0 13.671 -8.789 22.46t-22.46 8.789l-62.496 0q-13.671 0 -22.46 -8.789t-8.789 -22.46zm31.248 0l62.496 0l0 -406.224l-62.496 0l0 406.224zm156.24 0l0 -406.224q0 -13.671 8.789 -22.46t22.46 -8.789l62.496 0q13.671 0 22.46 8.789t8.789 22.46l0 406.224q0 13.671 -8.789 22.46t-22.46 8.789l-62.496 0q-13.671 0 -22.46 -8.789t-8.789 -22.46zm31.248 0l62.496 0l0 -406.224l-62.496 0l0 406.224z"
								/></svg
							>
						</Button>
					</td>
				</tr>
			{/each}

			<tr>
				<td>
					<Select {itemId} {label} items={functions_missing_calls} bind:value={selected_call_from}
					></Select>
					<p>
						{selected_call_from &&
							missing_calls_nr[selected_call_from.name] +
								' unresolved calls left for ' +
								selected_call_from.name}
					</p>
				</td>
				<td>
					<Select {itemId} {label} items={functions} bind:value={selected_call_to}></Select>
				</td>
				<td>
					<Button color="success" on:click={add_dynamic_call}>Add dynamic call</Button>
				</td>
			</tr>
		</tbody>
	</Table>
</div>
