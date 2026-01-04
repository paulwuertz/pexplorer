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
		InputGroup,
		Row,
		Col
	} from '@sveltestrap/sveltestrap';
	import Dropzone from 'svelte-file-dropzone';

	import { symbols } from './symbols.svelte.js';

	let CANNECTIVITY_SAMPLE_URL = 'https://p4w5.eu/report.json';
	let ZEPHYR_HELLO_SAMPLE_URL = 'https://p4w5.eu/reportHelloWorld.json';
	let ZEPHYR_MQTT_SAMPLE_URL = 'https://p4w5.eu/reportMQTTPublisher.json';

	let files = $state({
		accepted: [],
		rejected: []
	});
	let link_input_field = $state();
	let symbol_map = $state(symbols.symbols);
	let symbol_links = $state(symbols.symbolLinks);
	let versions = $derived(Object.keys(symbol_map));
	let selected_primary_versions = $state(symbols.selected_version);
	let selected_secondary_versions = $state(symbols.selected_versions_to_compare);

	function loadReportsfromJson(jsonReports) {
		symbols.symbols = JSON.parse(jsonReports);
	}

	function loadReportFromELF(elfBinary) {
		console.log(elfBinary);
	}

	function handleFilesSelect(e) {
		const { acceptedFiles, fileRejections } = e.detail;
		files.accepted = [...files.accepted, ...acceptedFiles];
		files.rejected = [...files.rejected, ...fileRejections];

		for (let i = 0; i < files.accepted.length; i++) {
			let file = files.accepted[i];
			const reader = new FileReader();
			// Validate file existence and type
			if (!file) {
				console.log('No file selected. Please choose a file.', 'error');
				return;
			}

			let file_type = file.type.toLocaleLowerCase();
			let is_json = file_type.endsWith('json'),
				is_elf = file.name.endsWith('elf');

			if (!(is_json || is_elf)) {
				console.log(
					file.type +
						'Unsupported file type ' +
						"'" +
						file.type +
						"'" +
						' - please select a text file.',
					'error'
				);
				return;
			}
			reader.onload = () => {
				if (is_json) {
					loadReportsfromJson(reader.result);
				} else if (is_elf) {
					loadReportFromELF(reader.result);
				} else {
				}
			};
			reader.onerror = () => {
				alert('Error reading the file. Please try again.');
			};

			if (is_json) {
				reader.readAsText(file);
			} else if (is_elf) {
				reader.readAsArrayBuffer(file);
			} else {
			}
		}
		// todo - handle files.rejected
	}

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
		symbols.symbols = {};
		symbols.symbolLinks = [];
		symbols.elfDataProvided = false;
		symbols.selected_version = undefined;
		symbols.selected_versions_to_compare = undefined;
		localStorage.removeItem('selected_version');
		localStorage.removeItem('selected_versions_to_compare');
		selected_primary_versions = undefined;
		selected_secondary_versions = undefined;
	}

	const reset_selected_versions = () => {
		selected_primary_versions = undefined;
		selected_secondary_versions = undefined;
		localStorage.removeItem('selected_version');
		localStorage.removeItem('selected_versions_to_compare');
	};

	const reset_secondary_versions = () => {
		selected_secondary_versions = undefined;
		localStorage.removeItem('selected_versions_to_compare');
	};

	const select_version = (version) => {
		if (!selected_primary_versions) {
			selected_primary_versions = version;
			localStorage.selected_version = version;
		} else if (!selected_secondary_versions) {
			selected_secondary_versions = version;
			localStorage.selected_versions_to_compare = version;
		} else {
			alert(
				'Only one version to view and a 2nd to compare can be selected. \n' +
					'To view another one reset selection and select a new one.'
			);
		}
		return;
	};
</script>

<div class="container" id="content">
	<Container fluid>
		<Card>
			<CardHeader>
				<CardTitle>
					Step 1: Add your firmwares symbol or Elf files, either...
					{#if versions.length}
						✅
					{/if}
				</CardTitle>
			</CardHeader>
			<CardBody>
				<Row>
					<Col sm="12" md={4}>
						<CardSubtitle><b>...by link:</b></CardSubtitle>
						<CardText
							>Adding the symbol via links saves them in your browsers local storage so you can
							continue browsing the same file when you come back.</CardText
						>
						<InputGroup>
							<Input
								type="url"
								bind:value={link_input_field}
								placeholder="enter a link to your firmwares symbol json..."
							/>
							<Button size="md" color="success" onclick={addLink}>Download symbols</Button>
						</InputGroup>
					</Col>
					<Col sm="12" md={4}>
						<CardSubtitle><b>...by file:</b></CardSubtitle>
						<CardText>
							<div class="uploadfield">
								<div>
									Uploading the symbol file is session based and is reset when refreshing or
									returning later.
								</div>
								<div>Any uploads are only processed locally in your browser!</div>

								<br />

								<Dropzone on:drop={handleFilesSelect} accept=".json,.elf">
									<p>Drag 'n' drop JSON-report or ELF files here, or click to select files</p>
								</Dropzone>
							</div>
						</CardText>
					</Col>
					<Col sm="12" md={4}>
						<CardSubtitle><b>...OR by loading a sample:</b></CardSubtitle>
						<CardText
							>Do not have any and just want to see a demo? <br /> Then load a sample to see some features
							:)</CardText
						>
						<div>
							<Button color="light" onclick={addCanncectifitySample}>cannectivity Releases</Button>
							<Button color="light" onclick={addZephyrSampleHELLO}>zephyr "hello world"</Button>
							<Button color="light" onclick={addZephyrSampleMQTT}>zephyr MQTT pub</Button>
							<Button color="light" onclick={addLocalSample}>Local report.json sample</Button>
						</div>
					</Col>
				</Row>
			</CardBody>
			<CardFooter>
				<Row>
					<Col sm="12" md={4}>
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
							{#each files.accepted as item}
								<li>{item.name}</li>
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
					</Col>
					<Col sm="12" md={4}>
                    </Col>
					<Col sm="12" md={4}>
                        <a href="/wip">WIP - new feature playground :)</a>
                    </Col>
				</Row>
			</CardFooter>
		</Card>

		<!-- <Card-- id="step2" class="mt-3">
			<CardHeader>
				<CardTitle>
					Step 2: Select a version to explore (mandatory)
					{#if selected_primary_versions}
						✅
					{/if}
					+ Select a second version to compare (optional)
					{#if selected_secondary_versions}
						✅
					{/if}
				</CardTitle>
			</CardHeader>
			<CardBody>
				<Row cols={{ lg: 3, md: 2, sm: 1 }}>
					{#each versions as version ('verselbtns-' + version)}
						<div class="pb-3 px-3">
							<Card>
								<CardHeader>
									<CardTitle>{version}</CardTitle>
								</CardHeader>
								<CardBody>
									<CardText>
										<TODO add source link symbol json >
										Buildtime: {symbols.symbols[version].timestamp}
									</CardText>
									{#if selected_primary_versions == version}
										<Button color="success" onclick={reset_selected_versions}>
											Selected to view
										</Button>
									{:else if selected_secondary_versions == version}
										<Button color="primary" onclick={reset_secondary_versions}>
											Selected to compare
										</Button>
									{:else if selected_primary_versions && selected_secondary_versions}
										<Button onclick={() => select_version(version)}>---</Button>
									{:else}
										<Button onclick={() => select_version(version)}>Browse</Button>
									{/if}
								</CardBody>
							</Card>
						</div>
					{/each}
				</Row>
			</CardBody>

			<CardFooter>
				<Button size="md" color="danger" onclick={reset_selected_versions}
					>Clear selected versions</Button
				>
			</CardFooter>
		</Card-->
	</Container>
</div>

<style>
	#content {
		margin-top: 20px;
	}
</style>
