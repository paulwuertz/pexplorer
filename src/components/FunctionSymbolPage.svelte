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

	const { symbol_data, symbol_version } = $props();

	// settings
	let show_full_asm = $state(false);

	// to display
	let symbol_path_and_name = $derived(symbol_data.file + symbol_data.name);
	// let asm = $derived(Object.hasOwn(symbol_data, 'asm') ? symbol_data.asm : []);
    let asm = [];
	let asm_code_preview = $derived(asm.slice(0, 5).join('\n'));
	let address = $derived(symbol_data.address);
	let stack_size = $derived(symbol_data.stack_size);
	let stack_qualifier = $derived(symbol_data.stack_qualifiers);
	let asm_code = $derived(asm.join('\n'));
	let callers = $derived(symbol_data.callers || []);
	let callees = $derived(symbol_data.callees || []);
	let code_size = $derived(symbol_data.size);

	const worst_call_stack = () => {
		// let my_symbol = { full_symbol_path: symbol_path_and_name, stack_size: symbol_data.stack_size };
		// let stack_down = deepest_callees_tree.concat([my_symbol]);
		// return stack_down.concat(deepest_callers_tree);
        return []
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
			<td><b>Function stack size</b>:</td>
			<td>
				{stack_size} bytes - stack usage is '{stack_qualifier}' - TODO add info about qualifiers :)
			</td>
		</tr>

		<tr>
			<td><b>Callers </b> ({callers.length}):</td>
			<td>
				{#each callers as caller}
					<a href={helpers.callxrs_text_to_links(base, symbol_version, caller)}>
						<small>
							{helpers.callxrs_text_to_symname(caller)}
						</small>
					</a>{', '}
				{/each}
			</td>
		</tr>
		<tr>
			<td><b>Callees</b> ({callees.length}):</td>
			<td>
				{#each callees as callee}
					<a href={helpers.callxrs_text_to_links(base, symbol_version, callee)}>
						<small>
							{helpers.callxrs_text_to_symname(callee)}
						</small>
					</a>{', '}
				{/each}
			</td>
		</tr>
	</tbody>
</Table>

<h4>Disassembly</h4>
<!-- TODO figure out better export to highlight and diff... <Highlight language={armasm} {asm_code} /> -->
<pre>
{#if show_full_asm}
		{asm_code}
	{:else}
		{asm_code_preview}
    ...
	{/if}
<span class="center" onclick={() => (show_full_asm = !show_full_asm)}>
    {#if show_full_asm}↑ show less ↑{:else}↓ show more ↓{/if}
</span>
</pre>

<h4>Stack Worst-Case Scenarios</h4>

<Table style="word-break: break-all;" hover bordered>
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
					<a href={helpers.callxrs_text_to_links(base, symbol_version, caller.full_symbol_path)}>
						{#if symbol_path_and_name.includes(caller.full_symbol_path)}
							<small>
								<b>{helpers.callxrs_text_to_symname(caller.full_symbol_path)}</b> - (this function)
							</small>
						{:else}
							<small>
								{helpers.callxrs_text_to_symname(caller.full_symbol_path)}
							</small>
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
			<td
				><b>&sum; {symbol_data.deepest_callee_tree_size + symbol_data.deepest_caller_tree_size}</b
				></td
			>
		</tr>
	</tfoot>
</Table>

<style>
	pre {
		background-color: #f5f5f5;
		border: 1px solid #ccc;
		border-radius: 4px;
	}
</style>
