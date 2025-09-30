<script lang="ts">
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';
	// ui stuff
	import {
		Button,
		ButtonGroup,
		Container,
		Card,
		CardBody,
		CardFooter,
		CardHeader,
		CardSubtitle,
		CardText,
		CardTitle,
		Input,
		InputGroup
	} from '@sveltestrap/sveltestrap';

	import { symbols } from './symbols.svelte.js';
	import * as helpers from './helpers.js';

	let CANNECTIVITY_SAMPLE_URL = 'https://p4w5.eu/report.json';
	let ZEPHYR_HELLO_SAMPLE_URL = 'https://p4w5.eu/reportHelloWorld.json';
	let ZEPHYR_MQTT_SAMPLE_URL = 'https://p4w5.eu/reportMQTTPublisher.json';
	let { data } = $props();

	let files = $state();
	let link_input_field = $state();
	let symbol_links = $state(symbols.symbolLinks);
	let versions = $derived(Object.keys(symbols.symbols));
	let selected_symbols = $state({});
	let function_table_data = $state([]);
	let variable_table_data = $state([]);

	const updateSelectedSymbols = () => {
		selected_symbols = helpers.symbolsToMap(symbols.symbols[symbols.selected_version]['symbols']);
		function_table_data = helpers.symbolsToFunctionMap(selected_symbols);
		variable_table_data = helpers.symbolsToVariableMap(selected_symbols);
	};

	const updateSelectedVersion = () => {
		localStorage.selected_version = symbols.selected_version;
		updateSelectedSymbols();
	};

	$effect(() => {
		if (files) {
			// Note that `files` is of type `FileList`, not an Array:
			// https://developer.mozilla.org/en-US/docs/Web/API/FileList
			console.log('files ' + files);
			const file = files[0];

			// Validate file existence and type
			if (!file) {
				console.log('No file selected. Please choose a file.', 'error');
				return;
			}

			if (!(file.type.endsWith('JSON') || file.type.endsWith('json'))) {
				console.log(file.type + 'Unsupported file type. Please select a text file.', 'error');
				return;
			}

			// Read the file
			const reader = new FileReader();
			reader.onload = () => {
				symbols.symbols = JSON.parse(reader.result);
			};
			reader.onerror = () => {
				alert('Error reading the file. Please try again.', 'error');
			};
			reader.readAsText(file);
		}
	});

	onMount(async () => {
		if (browser) {
			// load elf data
			if (symbols.symbols && Object.keys(symbols.symbols).length == 0) {
				console.log('No ELF data URL passed or stored, please upload it as a file then :)');
			} else {
				if (symbols.selected_version && symbols.selected_versions_to_compare) {
					updateSelectedSymbols();
				} else {
					console.log('ELF loaded, please select which version to show :)');
				}
			}
		}
	});

	async function addFirmwareByLink(link) {
		if (!symbol_links.includes(link)) {
			const response = await fetch(link);
			const data = await response.json();
			symbol_links.push(link);
			localStorage.lastOpenElfURLs = JSON.stringify(symbol_links);
			symbols.elfDataProvided = true;
			for (const versionKey of Object.keys(data)) {
				console.log('Loaded version data', versionKey);
				symbols.symbols[versionKey] = data[versionKey];
			}
		} else {
			alert('Link already added :)');
		}
	}

	function addLocalSample() {
		addFirmwareByLink('/report.json');
	}

	function addCanncectifitySample() {
		addFirmwareByLink(CANNECTIVITY_SAMPLE_URL);
	}

	function addZephyrSampleHELLO() {
		addFirmwareByLink(ZEPHYR_HELLO_SAMPLE_URL);
	}

	function addZephyrSampleMQTT() {
		addFirmwareByLink(ZEPHYR_MQTT_SAMPLE_URL);
	}

	function addLink() {
		addFirmwareByLink(link_input_field);
		link_input_field = null;
	}

	function resetLinks() {
		localStorage.removeItem('lastOpenElfURLs');
		symbol_links = [];
		versions = [];
		selected_symbols = {};
		symbols.symbols = {};
		symbols.symbolLinks = [];
		symbols.selected_version = null;
		symbols.selected_versions_to_compare = null;
		symbols.elfDataProvided = false;
	}
</script>

<div class="container" id="content">
	<Container fluid>
		<Card>
			<CardHeader>
				<CardTitle>Add your firmwares symbol files</CardTitle>
			</CardHeader>
			<CardBody>
				<CardSubtitle><b>By link:</b></CardSubtitle>
				<CardText
					>Adding the symbol via links saves them in your browsers local storage so you can continue
					browsing the same file when you come back.</CardText
				>

				<InputGroup>
					<Input
						type="url"
						bind:value={link_input_field}
						placeholder="enter a link to your firmwares symbol json..."
					/>
					<Button size="md" color="success" onclick={addLink}>Download symbols</Button>
				</InputGroup>

				<br />

				<CardSubtitle><b>By file:</b></CardSubtitle>
				<CardText
					>Uploading the symbol file is session based and is reset when refreshing or returning
					later.</CardText
				>
				<InputGroup>
					<Input type="file" accept="*/json" bind:files id="elfinput" name="elfinput" />
					<Button size="md" color="success">Upload symbols</Button>
				</InputGroup>

				<br />

				<CardSubtitle><b>Load a sample:</b></CardSubtitle>
				<CardText>Do not have any and just want to see a demo? Then load a sample :)</CardText>
				<ButtonGroup>
					<Button color="light" onclick={addCanncectifitySample}>cannectivity Releases</Button>
					<Button color="light" onclick={addZephyrSampleHELLO}>zephyr "hello world"</Button>
					<Button color="light" onclick={addZephyrSampleMQTT}>zephyr MQTT pub</Button>
					<Button color="light" onclick={addLocalSample}>Local report.json sample</Button>
				</ButtonGroup>
			</CardBody>
			<CardFooter>
				Currently provided symbol via links:

				<ul>
					{#each symbol_links as symbol_link, i ('link-' + symbol_link)}
						<li>{symbol_link}</li>
					{:else}
						<p>No links given yet.</p>
					{/each}
				</ul>

				and temporarily provided symbols via file upload:

				<ul>
					{#each files as file, i (files)}
						<li>{file}</li>
					{:else}
						<p>No files given yet.</p>
					{/each}
				</ul>

				yielding these firmware versions to view:

				<ul>
					{#each Object.keys(symbols.symbols) as sym_version, i ('ver-' + sym_version)}
						<li>{sym_version}</li>
					{:else}
						<p>No versions yet.</p>
					{/each}
				</ul>

				<Button size="md" color="danger" onclick={resetLinks}>Clear all links and files</Button>
			</CardFooter>
		</Card>
	</Container>
</div>

<style>
    #content {
        margin-top: 20px;
    }
</style>
