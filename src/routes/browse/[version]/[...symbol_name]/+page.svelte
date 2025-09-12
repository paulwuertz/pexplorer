<script>
	import { page } from '$app/stores';
	import { writable } from 'svelte/store';
	import { base } from '$app/paths';
    import { symbols } from '../../../symbols.svelte.js';

    import { Table } from '@sveltestrap/sveltestrap';
    import Highlight from "svelte-highlight";
    import atomOneDark from "svelte-highlight/styles/atom-one-dark";
    import armasm from "svelte-highlight/languages/armasm";

	import * as helpers from '../../../helpers.js';

	let { data, route, params } = $props();

    let show_full_asm = $state(false);

	let symbol_version = $derived(params.version);
	let symbol_path_and_name = $derived(params.symbol_name);
	let symbol_data = $derived(symbols.symbols[symbol_version].symbols.find( e => { return symbol_path_and_name.includes(e.file)} ));
	// to display
	let asm_code_preview = $derived(symbol_data.asm.slice(0, 5).join("\n"));
    let address	 = $derived(symbol_data.address);
	let asm_code = $derived(symbol_data.asm.join("\n"));
	let callers = $derived(JSON.parse(symbol_data.callers));
	let callees = $derived(JSON.parse(symbol_data.callees));
	let deepest_callers_tree = $derived(JSON.parse(symbol_data.deepest_caller_tree || false));
	let deepest_callees_tree = $derived(JSON.parse(symbol_data.deepest_callee_tree || false));
	let code_size = $derived(symbol_data.size);

	const worst_call_stack = () => {
        let my_symbol = {full_symbol_path: symbol_path_and_name, stack_size: symbol_data.stack_size};
        let stack_down = deepest_callees_tree.concat([my_symbol]);
        return stack_down.concat(deepest_callers_tree)
	};

</script>

<svelte:head>
  {@html atomOneDark}
</svelte:head>

<div class="container">

    <hr>

    <!-- TODO make it real breadcrumbs with links working :) -->
	{" / "} {symbol_path_and_name.split("/").join(" / ")}

    <hr>

    <Table hover bordered>
        <tbody>

            <tr>
                <td><b>Address</b>:</td>
                <td>
                    0x{address}
                </td>
            </tr>

            <tr>
                <td><b>Function code size</b>:</td>
                <td>
                    {code_size} bytes
                </td>
            </tr>

            <tr>
                <td><b>Callers </b> ({callers.length}):</td>
                <td>
                    {#each callers as caller}
                    <a href="{helpers.callxrs_text_to_links(base, symbol_version, caller)}">
                        <small>
                            {helpers.callxrs_text_to_symname(caller)}
                        </small>
                    </a>{", "}
                    {/each}

                </td>
            </tr>
            <tr>
                <td><b>Callees</b> ({callees.length}):</td>
                <td>
                    {#each callees as callee}
                    <a href="{helpers.callxrs_text_to_links(base, symbol_version, callee)}">
                        <small>
                            {helpers.callxrs_text_to_symname(callee)}
                        </small>
                    </a>{", "}
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
        <span class="center" onclick={() => show_full_asm=!show_full_asm}>
            {#if show_full_asm}↑ show less ↑{:else}↓ show more ↓{/if}
        </span>
    </pre>

    <h4>Stack Worst-Case Scenarios</h4>

    <Table hover bordered>
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
                    <td>{index + " "}</td>
                    <td>
                        <a href="{helpers.callxrs_text_to_links(base, symbol_version, caller.full_symbol_path)}">
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
                <td><b>&sum; {symbol_data.deepest_callee_tree_size + symbol_data.deepest_caller_tree_size}</b></td>
            </tr>
        </tfoot>
    </Table>
</div>

<style>

pre {
  background-color: #f5f5f5;
  border: 1px solid #ccc;
  border-radius: 4px;
}

</style>
