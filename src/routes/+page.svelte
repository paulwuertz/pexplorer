
<script lang="ts">
    import { onMount } from 'svelte';
    import { browser } from '$app/environment';
    // ui stuff
    import { Col, Container, Input, Row } from '@sveltestrap/sveltestrap';

  	import { symbols } from './symbols.svelte.js';

    import { scaleSequential, scaleOrdinal } from 'd3-scale';
    import * as chromatic from 'd3-scale-chromatic';
    import { hsl } from 'd3-color';

	  let { data } = $props();

    const sequentialColor = scaleSequential([4, -1], chromatic.interpolateGnBu);
    // filter out hard to see yellow and green
    const ordinalColor = scaleOrdinal(
        chromatic.schemeSpectral[9].filter((c) => hsl(c).h < 60 || hsl(c).h > 90)
    );
	  // let { data } = $props();
    // symbols.symbols = data.symbols;
    // symbols.selected_version = data.selected_version;
    // symbols.selected_versions_to_compare = data.selected_versions_to_compare;
    // symbols.elfDataProvided = data.elfDataProvided;
    let files = $state();
    let versions = $derived(Object.keys(symbols.symbols || []));
    let selected_symbols = $state({});
    let function_table_data = $state([]);
    let variable_table_data = $state([]);

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
        selected_symbols = symbolsToMap(symbols.symbols[symbols.selected_version]["symbols"]);
        function_table_data = symbolsToFunctionMap(selected_symbols);
        variable_table_data = symbolsToVariableMap(selected_symbols);
    };

    const updateSelectedVersion = () => {
        localStorage.selected_version = symbols.selected_version;
        updateSelectedSymbols();
    };

	$effect(() => {
		if (files) {
			// Note that `files` is of type `FileList`, not an Array:
			// https://developer.mozilla.org/en-US/docs/Web/API/FileList
			console.log("files "+files);
            const file = files[0];

            // Validate file existence and type
            if (!file) {
                console.log("No file selected. Please choose a file.", "error");
                return;
            }

            if (!(file.type.endsWith("JSON") || file.type.endsWith("json"))) {
                console.log(file.type+"Unsupported file type. Please select a text file.", "error");
                return;
            }

            // Read the file
            const reader = new FileReader();
            reader.onload = () => {
                symbols.symbols = JSON.parse(reader.result);
            };
            reader.onerror = () => {
                showMessage("Error reading the file. Please try again.", "error");
            };
            reader.readAsText(file);
		}
	});

    onMount(async () => {
        if (browser) {
            // load elf data
            if (symbols.symbols && Object.keys(symbols.symbols).length == 0) {
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
        {#if !symbols.elfDataProvided && files && !files[0]}
            <label for="elfinput">Upload a puncover .json file:</label>
            <input accept="*/json" bind:files id="elfinput" name="elfinput" type="file" />
        {:else if !symbols.selected_version}
            <h3>Select a version to browse elf symbols :)</h3>
        {:else}
            // TODO
        {/if}

    </Container>
</div>



