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

	const { var_childs, symbol_version, SymPathByAddr, VariableTypes, sections } = $props();

	let sym_data = $derived(var_childs[0]);
	// to display
	let symbol_path_and_name = $derived(sym_data.file + sym_data.name);
	let address = $derived(sym_data.address);
	let code_size = $derived(sym_data.size);
	let staticInitDataBase64 = $derived(sym_data.staticInitData);
	let staticInitData = $derived(
		(staticInitDataBase64 && Uint8Array.fromBase64(staticInitDataBase64)) || undefined
	);
	let section = $derived(sections[sym_data.secidx]);

	let type = $derived(sym_data.type);
	let typeStruct = $derived(
		VariableTypes.find((val, i, arr) => {
			return val.name == type;
		})
	);

	let getMemberBytes = (data, byte_offset, size) => {
		let b = [];
		let val = 0;
		for (let i = byte_offset; i < data.length && i < byte_offset + size; i++) {
			b.push(data[i]);
			val += data[i] * Math.pow(256, i - byte_offset);
		}
		return JSON.stringify(b) + ' - 0x' + val.toString(16);
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
			<td><b>Variable size</b>:</td>
			<td>
				{code_size} bytes
			</td>
		</tr>

		<tr>
			<td><b>Type</b>:</td>
			<td>
				{type}
			</td>
		</tr>
		{#if section}
			<tr>
				<td><b>Section</b>:</td>
				<td>
					{sections[sym_data.secidx].name}
				</td>
			</tr>
		{/if}
	</tbody>
</Table>

{#if staticInitData}
	<h4>Init Data</h4>
	<!-- TODO figure out better export to highlight and diff... <Highlight language={armasm} {asm_code} /> -->
	<pre>
	Base64: {staticInitDataBase64}
	arrayindex: {staticInitData.map((e, i) => i)}
	Uint8Array: {staticInitData}
	<!-- // .map(function(e, i, a){return "index "+i+" => "+ e}) -->
</pre>
{/if}

{#if typeStruct && staticInitData}
	<h4>Interpretation</h4>

	<Table style="word-break: break-all;" hover bordered>
		<thead>
			<tr>
				<th># byte</th>
				<th>Name</th>
				<th>Size</th>
				<th>Value</th>
			</tr>

			{#each typeStruct.members as m (m.name)}
				<tr>
					<td>{m.byte_offset}</td>
					<td
						>{#if m.is_pointer}*{/if}{m.name}</td
					>
					<td>{m.size}</td>
					<td>{getMemberBytes(staticInitData, m.byte_offset, m.size)}</td>
				</tr>
				<!-- {#if row[column.key] === 'var'}
                            <img src="{base}/icons/Method_16x.svg" width="16px" title="Type variable" />
                            {#if Object.hasOwn(row, 'type')}
                                {row['type']}
                            {/if}
                        {:else if row[column.key] === 'fn'}
                            <img src="{base}/icons/Field_16x.svg" width="16px" title="Type function" />
                        {/if} -->
			{/each}
		</thead>
	</Table>
{/if}

<style>
	pre {
		background-color: #f5f5f5;
		border: 1px solid #ccc;
		border-radius: 4px;
	}
</style>
