<script>
	import { onMount } from 'svelte';
	import { base } from '$app/paths';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';

	import { DataTable } from '@careswitch/svelte-data-table';
	import {
		Alert,
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
	import FlameGraph from './FlameGraph.svelte';
	import * as helpers from '../routes/helpers.js';

	const { fn_childs, symbol_version, SymPathByAddr } = $props();

	let sym_data = $derived(fn_childs[0]);
	let sym_path_by_addr = $derived(SymPathByAddr);
	// settings
	let show_full_asm = $state(false);
	// to display
	let symbol_path_and_name = $derived(sym_data.file + sym_data.name);
	let asm_code = $derived(
		sym_data.asm
			? helpers.csBase64ToASMText(sym_data.asm, sym_data.address, show_full_asm)
			: undefined
	);
	let address = $derived(sym_data.address);
	let stack_size = $derived(sym_data.stack_size);
	let max_stack_size_callees = $derived(sym_data.max_stack_size_callees);
	let stack_qualifier = $derived(sym_data.stack_qualifiers);
	let callers = $derived(sym_data.callers || []);
	let callees = $derived(sym_data.callees || []);
	let code_size = $derived(sym_data.size);

	// TODO check if from ELF or JSON...
	// or export... could take much time and disk...
	let fn_calltree = $derived(JSON.parse(get_fn_calltree(sym_data.address)));
	let unresolved = $derived(
		Object.entries(
			(fn_calltree.unresolved || []).reduce((acc, element) => {
				acc[element.from] = (acc[element.from] || 0) + 1;
				return acc;
			}, {})
		)
	);
	let branches = $derived(fn_calltree.branches);

	const worst_call_stack = () => {
		// let my_symbol = { full_symbol_path: symbol_path_and_name, stack_size: sym_data.stack_size };
		// let stack_down = deepest_callees_tree.concat([my_symbol]);
		// return stack_down.concat(deepest_callers_tree);
		return [];
	};
</script>

<Table style="word-break: break-all;" hover bordered>
	<tbody>
		<tr>
			<td><b>Address</b>:</td>
			<td>
				0x{address.toString(16)}
			</td>
		</tr>

		<tr>
			<td><b>Function code size</b>:</td>
			<td>
				{code_size} bytes
			</td>
		</tr>

		<tr>
			<td><b>Function stack usage</b>:</td>
			<td>
				{stack_size} bytes - qualified as '{stack_qualifier}' - TODO add info about qualifiers :)
			</td>
		</tr>

		<tr>
			<td><b>Callers </b> ({callers.length}):</td>
			<td>
				{#each callers as caller}
					{#if caller.from}
						<a
							href={helpers.callxrs_text_to_links(
								base,
								symbol_version,
								caller,
								sym_path_by_addr,
								true
							)}
						>
							<small>
								{helpers.callxrs_text_to_symname(caller, sym_path_by_addr, true)}
							</small>
						</a>
					{:else}
						<small>
							{helpers.callxrs_text_to_symname(caller, sym_path_by_addr, true)}
						</small>
					{/if}{', '}
				{/each}
			</td>
		</tr>
		<tr>
			<td><b>Callees</b> ({callees.length}):</td>
			<td>
				{#each callees as callee}
					{#if callee.to}
						<a
							href={helpers.callxrs_text_to_links(
								base,
								symbol_version,
								callee,
								sym_path_by_addr,
								false
							)}
						>
							<small>
								{helpers.callxrs_text_to_symname(callee, sym_path_by_addr, false)}
							</small>
						</a>{', '}
					{:else}
						<small>
							{helpers.callxrs_text_to_symname(callee, sym_path_by_addr, false)}
						</small>
					{/if}
				{/each}
			</td>
		</tr>
	</tbody>
</Table>

<h4>Disassembly</h4>
<!-- TODO figure out better export to highlight and diff... <Highlight language={armasm} {asm_code} /> -->
<pre>
{asm_code}
<span class="center" onclick={() => (show_full_asm = !show_full_asm)}>
    {#if show_full_asm}↑ show less ↑{:else}↓ show more ↓{/if}
</span>
</pre>

<h4>
	Worst-Case Stack Scenario currently found{#if max_stack_size_callees}: {max_stack_size_callees} bytes{/if}
</h4>

<Alert color="warning">
	<h4 class="alert-heading text-capitalize">Warning: experimental</h4>

	Estimating the stack use from the ELF+dwarf debug .framesection is still an experiment and needs
	testing!
</Alert>

{#key unresolved}
	{#if unresolved && unresolved.length > 0}
		<Alert color="warning">
			<h4 class="alert-heading text-capitalize">Warning: unresolved function calls</h4>

			There are {unresolved.length} unresolved function calls in this calltree.
			<br />
			These can be dynamic functions like callback.
			<br />
			Unresolved calls are from:

			<ul>
				{#each unresolved as addr_occurences}
					<li>
						<a
							href={helpers.callxrs_text_to_links(
								base,
								symbol_version,
								parseInt(addr_occurences[0]),
								sym_path_by_addr,
								true
							)}
						>
							<small>
								{helpers.callxrs_text_to_symname(
									{ from: parseInt(addr_occurences[0]) },
									sym_path_by_addr,
									true
								)}
							</small>
						</a>
						- called ({addr_occurences[1]}) time{#if addr_occurences[1] > 1}s{/if}
					</li>
				{/each}
			</ul>

			Try to manually resolve them - TODO implement then link resolve page :)
		</Alert>
	{/if}
{/key}

<h5>Top 10 callpaths with highest stack usage</h5>

<Table style="word-break: break-all;" hover bordered>
	<thead>
		<tr>
			<th>#</th>
			<th>Stack usage</th>
			<th>Callpath causing this stack</th>
		</tr>
	</thead>
	<tbody>
		{#each branches.slice(0, 10) as b, i}
			<tr>
				<td>i</td>
				<td>{JSON.stringify(b['stack_size'])}</td>
				<td>
					{#each b['call_list'] as c, j}
						{c['name'] + ' (' + c['stack_size'] + ') ->'}
					{/each}
				</td>
			</tr>
		{/each}
	</tbody>
</Table>

<FlameGraph {sym_data} {fn_calltree}></FlameGraph>

<!-- <Table style="word-break: break-all;" hover bordered>
	<thead>
		<tr>
			<th>#</th>
			<th>Name</th>
			<th>Stack size</th>
		</tr>
	</thead>
	<tbody>
		{#each worst_call_stack() as caller, index}
			<tr>
				<td>{index + ' '}</td>
				<td>
					 <a href={helpers.callxrs_text_to_links(base, symbol_version, caller, sym_path_by_addr)}>
						{#if symbol_path_and_name.includes(caller)}
							<small>
								<b>{helpers.callxrs_text_to_symname(caller, sym_path_by_addr)}</b> - (this function)
							</small>
						{:else}
						{/if}
					</a>
				</td>
				<td>
					{caller.stack_size}
				</td>
			</tr>
		{/each}
	</tbody>
	<tfoot>
		<tr>
			<td></td>
			<td></td>
			<td><b>&sum; {sym_data.deepest_callee_tree_size + sym_data.deepest_caller_tree_size}</b></td>
		</tr>
	</tfoot>
</Table> -->

<style>
	pre {
		background-color: #f5f5f5;
		border: 1px solid #ccc;
		border-radius: 4px;
	}
</style>
