<script>
	import { onMount } from 'svelte';
	import { base } from '$app/paths';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';

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
	import * as helpers from '../routes/helpers.js';

	const { fnSymbols, selected_version } = $props();
	let function_table_data = $derived(helpers.symbolsToFunctionMap(fnSymbols));
	let function_table = $derived(
		new DataTable({
			pageSize: 99999, // TODO
			data: function_table_data,
			columns: [
				{ id: 'name', key: 'display_name', name: 'Name' },
				{ id: 'remark', key: 'remark', name: 'Remarks' },
				{ id: 'size', key: 'size', name: 'Code size' },
				{ id: 'stack_size', key: 'stack_size', name: 'Stack size' },
				{ id: 'stack_qualifiers', key: 'stack_qualifiers', name: 'Stack size type' }
			]
		})
	);
</script>

<Row>
	<Col>
		Filter options:
		<FormGroup>
			<Input
				type="text"
				placeholder="Filter symbols by name"
				class="md:m3-auto md:max-w-[500px]"
				bind:value={function_table.globalFilter}
			/>
			<Label>
				<p>
					Showing {function_table.allRows.length} / {function_table.baseRows.length} functions
				</p>
			</Label>
		</FormGroup>
	</Col>
</Row>

<Table hover bordered style="word-break: break-all;">
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
		<tr>
			{#each function_table.columns as column (column.name)}
				<th>
					{#if column.id == 'name'}
						Sum of all selected symbols
					{/if}
					{#if column.id == 'size'}
						&sum; = {function_table.allRows.reduce((accumulator, row) => accumulator + row.size, 0)}
					{/if}
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
							<a
								data-sveltekit-preload-data="tap"
								href={helpers.row2AHref(base, selected_version, row)}
							>
								{row[column.key]}
							</a>
						</td>
					{:else}
						<td>{row[column.key]}</td>
					{/if}
				{/each}
			</tr>
		{/each}
	</tbody>
</Table>
