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
	import * as helpers from '../helpers.js';

	let { data } = $props();
	let files = $state();
	let versions = $derived(Object.keys(symbols.symbols));
	let selected_symbols = $state({});
	let function_table_data = $state([]);
	let variable_table_data = $state([]);
	let function_table = $derived(
		new DataTable({
			pageSize: 9999, // TODO
			data: function_table_data,
			columns: [
				{ id: 'name', key: 'display_name', name: 'Name' },
				{ id: 'remark', key: 'remark', name: 'Remarks' },
				{ id: 'size', key: 'size', name: 'Code size' },
				{ id: 'stack_size', key: 'stack_size', name: 'Stack size' }
			]
		})
	);
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
		selected_symbols = helpers.symbolsToMap(symbols.symbols[symbols.selected_version]['symbols']);
		function_table_data = helpers.symbolsToFunctionMap(selected_symbols);
		variable_table_data = helpers.symbolsToVariableMap(selected_symbols);
		//alert(function_table_data.length+" !! "+variable_table_data.length)
		function_table = new DataTable({
			pageSize: function_table_data.length,
			data: function_table_data,
			columns: [
				{ id: 'name', key: 'display_name', name: 'Name' },
				{ id: 'remark', key: 'remark', name: 'Remarks' },
				{ id: 'size', key: 'size', name: 'Code size' },
				{ id: 'stack_size', key: 'stack_size', name: 'Stack size' }
			]
		});
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
		{#if !symbols.elfDataProvided && files && !files[0]}
			<label for="elfinput">Upload a puncover .json file:</label>
			<input accept="*/json" bind:files id="elfinput" name="elfinput" type="file" />
		{:else if !symbols.selected_version}
			<h3>Select a version to browse elf symbols :)</h3>
		{:else}
			<h3>Function symbols for {symbols.selected_version}</h3>

			<Table hover bordered>
				<thead>
					<tr>
						{#each function_table.columns as column (column.name)}
							<th>
								{column.name}
								<button
									class="flex items-center"
									onclick={() => function_table.toggleSort(column.id)}
									disabled={!function_table.isSortable(column.id)}
								>
									{#if function_table.isSortable(column.id)}
										<span class="ml-2">
											{#if function_table.getSortState(column.id) === 'asc'}
												↑
											{:else if function_table.getSortState(column.id) === 'desc'}
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
					{#each function_table.rows as row (row.file + row.name)}
						<tr>
							{#each function_table.columns as column (column.name)}
								{#if column.key == 'display_name'}
									<td>
                                        <a data-sveltekit-preload-data="tap" href={helpers.row2AHref(base, symbols.selected_version, row)}>
                                            {row[column.key]}
                                        </a>
                                    </td
									>
								{:else}
									<td>{row[column.key]}</td>
								{/if}
							{/each}
						</tr>
					{/each}
				</tbody>
			</Table>

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
