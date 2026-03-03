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
		InputGroup,
		Label,
		Row,
		Table
	} from '@sveltestrap/sveltestrap';
	import * as helpers from '../routes/helpers.js';

	const { fnSymbols, varSymbols, selected_version, sections } = $props();
	let show_sections = $state(true);
	let show_filepath = $state('Full filepath');
	let function_table_data = $derived([...fnSymbols, ...varSymbols]);
	let function_table = $derived(
		new DataTable({
			pageSize: 99999, // TODO
			data: function_table_data,
			columns: [
				{ id: 'symtype', key: 'symtype', name: 'Type' },
				{ id: 'file', key: 'file', name: 'Filepath' },
				{ id: 'name', key: 'name', name: 'Name' },
				// { id: 'remark', key: 'remark', name: 'Remarks' }, TODO discuss removal?
				{ id: 'secidx', key: 'secidx', name: 'Section' },
				{ id: 'address', key: 'address', name: 'Address' },
				{ id: 'size', key: 'size', name: 'Symbol size' },
				{ id: 'callers', key: 'callers', name: '#Callers' },
				{ id: 'callees', key: 'callees', name: '#Callees' },
				{ id: 'stack_size', key: 'stack_size', name: 'Stack size' },
				{ id: 'stack_qualifiers', key: 'stack_qualifiers', name: 'Stack size type' }
			]
		})
	);
</script>

<hr />

<Row>
	<!-- TODO add filter for fn, var or both :) -->
	<Col sm="12" md={4}>Filter symbols by text:</Col>
	<Col sm="12" md={4}>Show file column:</Col>
</Row>
<Row>
	<Col sm="12" md={4}>
		<Input
			type="text"
			placeholder="Enter a filter string"
			class="md:m3-auto md:max-w-[500px]"
			bind:value={function_table.globalFilter}
			style={'margin-top: 10px;'}
		/>
	</Col>

	<Col sm="12" md={4}>
		<InputGroup>
			{#each ['None', 'Filename only', 'Full filepath'] as value}
				<Input type="radio" bind:group={show_filepath} {value} label={value} />
			{/each}
		</InputGroup>
		<Input bind:checked={show_sections} type="switch" label="Show section column:" />
	</Col>
</Row>

<hr />

<Row>
	<Col>
		<p>Showing {function_table.allRows.length} / {function_table.baseRows.length} symbols</p>
	</Col>
</Row>

<Table hover bordered style="word-break: break-all;">
	<thead>
		<tr>
			{#each function_table.columns as column (column.name)}
				{#if (column.id != 'file' && column.id != 'secidx') || (column.id == 'file' && show_filepath != 'None') || (column.id == 'secidx' && show_sections)}
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
				{/if}
			{/each}
		</tr>
		<tr>
			{#each function_table.columns as column (column.name)}
				{#if (column.id != 'file' && column.id != 'secidx') || (column.id == 'file' && show_filepath != 'None') || (column.id == 'secidx' && show_sections)}
					<th>
						{#if column.id == 'name'}
							Sum of all selected symbols
						{/if}
						{#if column.id == 'size'}
							<!-- TODO add split in ROM/FLash + RAM for vars in .bss and .data -->
							&sum; = {function_table.allRows.reduce(
								(accumulator, row) => accumulator + row.size,
								0
							)}
						{/if}
					</th>
				{/if}
			{/each}
		</tr>
	</thead>
	<tbody>
		{#each function_table.rows as row, i (row.file + row.name + i)}
			<tr>
				{#each function_table.columns as column (column.name)}
					{#if column.key == 'name'}
						<td>
							<a
								data-sveltekit-preload-data="tap"
								href={helpers.row2AHref(base, selected_version, row)}
							>
								{row[column.key]}
							</a>
						</td>
					{:else if column.key == 'secidx'}
						<td>
							{#if show_sections && sections && sections['secidx']}
								{sections['secidx'].name}
							{/if}
						</td>
					{:else if column.key == 'address'}
						<td>0x{row[column.key].toString(16)}</td>
					{:else if column.key == 'file'}
						{#if show_filepath != 'None'}
							<td>
								{#if show_filepath == 'Filename only'}
									{(row[column.key] || '').substring((row[column.key] || '').lastIndexOf('/') + 1)}
								{:else}
									{row[column.key]}
								{/if}
							</td>
						{/if}
					{:else if column.id == 'symtype'}
						<td style="min-width: 0px;">
							{#if row[column.key] === 'var'}
								<img src="{base}/icons/Method_16x.svg" width="16px" title="Type variable" />
							{:else if row[column.key] === 'fn'}
								<img src="{base}/icons/Field_16x.svg" width="16px" title="Type function" />
							{/if}
						</td>
					{:else if column.id == 'callers' || column.id == 'callees'}
						<td>{(row[column.key] || []).length}</td>
					{:else}
						<td>{row[column.key]}</td>
					{/if}
				{/each}
			</tr>
		{/each}
	</tbody>
</Table>

<style>
	td {
		min-width: 130px;
		font-size: 14px;
		padding-top: 0.25rem;
		padding-bottom: 0.25rem;
	}

	/* purgecss start ignore */
	:global(.form-check) {
		margin-right: 15px;
	}
	/* purgecss end ignore */
</style>
