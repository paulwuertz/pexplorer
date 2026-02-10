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

	const { var_childs, symbol_version, SymPathByAddr } = $props();

	let sym_data = $derived(var_childs[0]);
	// to display
	let symbol_path_and_name = $derived(sym_data.file + sym_data.name);
	let address = $derived(sym_data.address);
	let code_size = $derived(sym_data.size);
	let staticInitDataBase64 = $derived(sym_data.staticInitData);
	let staticInitData = $derived(Uint8Array.fromBase64(staticInitDataBase64) || []);
	let type = $derived(sym_data.type);
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

		<tr>
			<td><b>Segment</b>:</td>
			<td>
				<!-- {code_size} bytes -->
			</td>
		</tr>
	</tbody>
</Table>

<h4>Init Data</h4>
<!-- TODO figure out better export to highlight and diff... <Highlight language={armasm} {asm_code} /> -->
<pre>
	Base64: {staticInitDataBase64}
	arrayindex: {staticInitData.map((e, i) => i)}
	Uint8Array: {staticInitData}
	<!-- // .map(function(e, i, a){return "index "+i+" => "+ e}) -->
</pre>

<h4>Interpretation</h4>

<Table style="word-break: break-all;" hover bordered>
	<thead>
		<tr>
			<th>#</th>
			<th>Name</th>
			<th>Stack size</th>
		</tr>
	</thead>
</Table>

<style>
	pre {
		background-color: #f5f5f5;
		border: 1px solid #ccc;
		border-radius: 4px;
	}
</style>
