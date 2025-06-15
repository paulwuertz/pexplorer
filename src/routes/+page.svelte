
<script lang="ts">
    import { onMount } from 'svelte';
    import { browser } from '$app/environment';
    // ui stuff
    import {
      Button, Container,
      Card, CardBody, CardFooter, CardHeader, CardSubtitle, CardText, CardTitle,
      Input, InputGroup
    } from '@sveltestrap/sveltestrap';

  	import { symbols } from './symbols.svelte.js';

    let { data } = $props();

    let files = $state();
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

    <Container fluid>
      <Card>
        <CardHeader>
          <CardTitle>Add your firmwares symbol files</CardTitle>
        </CardHeader>
        <CardBody>
          <CardSubtitle><b>By link:</b></CardSubtitle>
          <CardText>Adding the symbol via links saves them in your browsers local storage so you can continue browsing the same file when you come back.</CardText>

          <InputGroup>
            <Input type="url" placeholder="enter a link to your firmwares symbol json..." />
            <Button size="md" color="success">Download symbols</Button>
          </InputGroup>

          <br>

          <CardSubtitle><b>By file:</b></CardSubtitle>
          <CardText>Uploading the symbol file is session based and is reset when refreshing or returning later.</CardText>
          <InputGroup>
            <Input type="file" accept="*/json" bind:files id="elfinput" name="elfinput" />
            <Button size="md" color="success">Upload symbols</Button>
          </InputGroup>
        </CardBody>
        <CardFooter>
          Currently provided symbol via links:

          <ul>
            <li>link</li>
          </ul>

          and temporarily provided symbols via file upload:

          <ul>
            <li>file</li>
          </ul>

          <Button size="md" color="danger">Clear all links and files (TODO)</Button>
        </CardFooter>
      </Card>

    </Container>
</div>



