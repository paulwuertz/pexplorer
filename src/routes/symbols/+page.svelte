<script>
	import { onMount } from 'svelte';
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import { browser } from '$app/environment';
	import { writable } from 'svelte/store';
	// ui stuff
	import { DataTable } from '@careswitch/svelte-data-table';
	import {
		Badge,
		Button,
		Col,
		Container,
		FormGroup,
		Input,
		Label,
		Row,
		Table
	} from '@sveltestrap/sveltestrap';

	import { symbols } from '../symbols.svelte.js';
    import * as helpers from '../helpers.js';
	import FunctionSymbolTable from '../../components/FunctionSymbolTable.svelte';

	let { data } = $props();
	let files = $state();
	let versions = $derived(Object.keys(symbols.symbols));
	let all_symbols = $derived(symbols.symbols);
	let selected_version = $derived(symbols.selected_version);
	let selected_symbol_array = $derived(
		Object.hasOwn(all_symbols, selected_version) ? all_symbols[selected_version]['symbols'] : []
	);
	let selected_symbols = $derived(helpers.symbolsToMap($state.snapshot(selected_symbol_array)));
	let variable_table_data = $derived(helpers.symbolsToVariableMap(selected_symbols));
	let variable_table = $derived(
		new DataTable({
			pageSize: 999999, // TODO
			data: variable_table_data,
			columns: [
				{ id: 'name', key: 'display_name', name: 'Name' },
				{ id: 'remark', key: 'remark', name: 'Remarks' },
				{ id: 'size', key: 'size', name: 'Static size' }
			]
		})
	);

	const updateSelectedSymbols = () => {
		variable_table = new DataTable({
			pageSize: variable_table_data.length,
			data: variable_table_data,
			columns: [
				{ id: 'name', key: 'name', name: 'Name' },
				{ id: 'remark', key: 'remark', name: 'Remarks' },
				{ id: 'size', key: 'size', name: 'Static size' }
			]
		});
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
				const selected_version = localStorage.getItem('selected_version');
				if (symbols.selected_version) {
					updateSelectedSymbols();
				} else if (selected_version) {
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

	<Container fluid>
		{#if !symbols.elfDataProvided && files && !files[0]}
			<label for="elfinput">Upload a puncover .json file:</label>
			<input accept="*/json" bind:files id="elfinput" name="elfinput" type="file" />
		{:else if !symbols.selected_version}
			<h3>Select a version to browse elf symbols :)</h3>
		{:else}
            <h3>Function symbols for {symbols.selected_version}</h3>

			<FunctionSymbolTable
                fnSymbols={selected_symbols}
                selected_version={selected_version}
            />

			<h3>Variable symbols for {symbols.selected_version}</h3>

			<Table>
				<thead>
					<tr>
						{#each variable_table.columns as column (column.name)}
							<th>
								{column.name}
								<button
									class="flex items-center"
									onclick={() => variable_table.toggleSort(column.id)}
									disabled={!variable_table.isSortable(column.id)}
								>
									{#if variable_table.isSortable(column.id)}
										<span class="ml-2">
											{#if variable_table.getSortState(column.id) === 'asc'}
												↑
											{:else if variable_table.getSortState(column.id) === 'desc'}
												↓
											{:else}
												↕
											{/if}
										</span>
									{/if}
								</button>
							</th>
						{/each}
					</tr>
				</thead>
				<tbody>
					{#each variable_table.rows as row (row.file + row.name)}
						<tr>
							{#each variable_table.columns as column (column.name)}
								<td>{row[column.key]}</td>
							{/each}
						</tr>
					{/each}
				</tbody>
			</Table>
		{/if}
	</Container>
</div>

<style>
	td {
		min-width: 125px;
	}
</style>
