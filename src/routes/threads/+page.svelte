<script>
    import { onMount } from 'svelte';
    import { base } from '$app/paths';
    import { page } from '$app/stores';
    import { browser } from '$app/environment';
    import { writable } from "svelte/store";
    // ui stuff
    import { Badge, Button, Col, Container, Input, Row, Table } from '@sveltestrap/sveltestrap';

	import { symbols } from '../symbols.svelte.js';

	let { data } = $props();
    let files = $state();
    let versions = $derived(Object.keys(symbols.symbols));
    let selected_thread_stat = $state({});

    let symbolsToMap = (syms) => {
        let symMap = {};
        for (const sym of syms) {
            sym.remark = sym.called_from_other_file ? "x-module" : "";
            sym.newSymbols = false;
            sym.deletedSymbols = false;
            symMap[sym.file+sym.display_name] = sym;
        }
        return symMap;
    }

    let symbolsToFunctionMap = (symMap) => {
        return Object.values(symMap).filter((e) => {return e["type"] === "function";})
    }

    let symbolsToVariableMap = (symMap) => {
        return Object.values(symMap).filter((e) => {return e["type"] === "variable";})
    }

    let symMapToSymNameSet = (symMap) => {
        return new Set(Object.keys(symMap));
    }

    const updateSelectedSymbols = () => {
        let versionObj = symbols.symbols[symbols.selected_version];
        if(!versionObj || !versionObj.hasOwnProperty("stack_reports")) return;
        console.log(versionObj["stack_reports"]);
        for (const [k,v] of Object.entries(versionObj["stack_reports"])) {
            selected_thread_stat[k] = v;
        }
    };

    const updateSelectedVersion = () => {
        localStorage.selected_version = symbols.selected_version;
        updateSelectedSymbols();
    };

    onMount(async () => {
        if (browser) {
            // load elf data
            if (Object.keys(symbols.symbols).length == 0) {
                console.log("No ELF data URL passed or stored, please upload it as a file then :)");
            } else {
                if(symbols.selected_version && symbols.selected_versions_to_compare)
                {
                    updateSelectedSymbols();
                } else {
                    console.log("ELF loaded, please select which version to show :)");
                }
            }
        }
    });
</script>

<style>
    /*
    @import 'static/css/style.css';
    */
    @import 'https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css';
</style>

<div class="container" id="content">

    <Row>
        <Col>
            Select a version of the .elf you want to see:
            <Input type="select"
                bind:value={symbols.selected_version}
                on:change={updateSelectedVersion}
            >
                {#each versions as version}
                    <option>{version}</option>
                {/each}
            </Input>
        </Col>
      </Row>

      <hr>

    <Container fluid>
        {#if !symbols.selected_version}
            <h3>Select a version to browse thread info :)</h3>
        {:else}

            <h3>Stack data for {symbols.selected_version}</h3>

            {#key selected_thread_stat}

            <ul>
                {#each Object.keys(selected_thread_stat) as thread_name, index (thread_name+index)}
                <h4>Thread stats for '{thread_name}'</h4>

                <Table>
                    <thead>
                        <tr>
                            <th width="100%">Name</th>
                            <th>Stack Size</th>
                            <!-- <th colspan="2" style="text-align: center">Code</th> -->
                        </tr>
                        <!-- <tr>
                            <th class="col_size">
                                &sum; = {selected_thread_stat[thread_name].max_static_stack_size} / {selected_thread_stat[thread_name].max_stack_size}
                            </th>
                            <th class="col_size"></th>
                            <th class="col_size">&sum;= TODO not exported atm...</th>
                            <th class="col_size"></th>
                        </tr> -->
                    </thead>
                    <tbody>
                        {#each selected_thread_stat[thread_name]["call_stack"] as fn, fn_index  (thread_name+"_"+fn.name+"_"+index+"_"+fn_index)}
                            <tr>
                                <td>{fn.name}</td>
                                <td>{fn.stack_size}</td>
                            </tr>
                        {/each}
                    </tbody>
                    <tfoot>
                      <tr>
                        <td><b>Total stack usage</b></td>
                        <td><b>
                            &sum; = {selected_thread_stat[thread_name].max_static_stack_size} / {selected_thread_stat[thread_name].max_stack_size}
                        </b></td>
                      </tr>
                    </tfoot>
                </Table>

                <hr>
                {/each}
            </ul>
            {/key}
        {/if}
    </Container>
</div>



