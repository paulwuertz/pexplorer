<script lang="ts">
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
	import * as helpers from './helpers.js';
	import { base } from '$app/paths';
	const testFW = [
		{
			name: 'CANnectivity v1.2 LLVM @ lpc55s16',
			url: 'https://media.githubusercontent.com/media/paulwuertz/pexplorer/refs/heads/main/testdata/zephyr_cannectivity_12_llvm_lpc55s16.elf'
		},
		{
			name: 'CANnectivity v1.2 GCC @ lpc55s16',
			url: 'https://media.githubusercontent.com/media/paulwuertz/pexplorer/refs/heads/main/testdata/zephyr_cannectivity_12_gcc_lpc55s16.elf'
		},
		{
			name: 'CANnectivity v1.3 LLVM @ lpc55s16',
			url: 'https://media.githubusercontent.com/media/paulwuertz/pexplorer/refs/heads/main/testdata/zephyr_cannectivity_13_llvm_lpc55s16.elf'
		},
		{
			name: 'CANnectivity v1.3 GCC @ lpc55s16',
			url: 'https://media.githubusercontent.com/media/paulwuertz/pexplorer/refs/heads/main/testdata/zephyr_cannectivity_13_gcc_lpc55s16.elf'
		},
		{
			name: 'Prusa Buddy - Core One v6.4.0',
			url: 'https://media.githubusercontent.com/media/paulwuertz/pexplorer/refs/heads/main/testdata/prusa_buddy_boot_64.elf'
		},
		{
			name: 'IronOS Pinecilv1 EN v2.23 (RISC-V WiP)',
			url: 'https://media.githubusercontent.com/media/paulwuertz/pexplorer/refs/heads/main/testdata/Pinecilv1_EN_v2_23.elf'
		},
		{
			name: 'IronOS Pinecilv2 EN v2.23 (RISC-V WiP)',
			url: 'https://media.githubusercontent.com/media/paulwuertz/pexplorer/refs/heads/main/testdata/Pinecilv2_EN_v2_23.elf'
		},
		{
			name: 'ZSWatch v0.7.0',
			url: 'https://media.githubusercontent.com/media/paulwuertz/pexplorer/refs/heads/main/testdata/zswatch_nrf5340_07.elf'
		}
	];

	let files = $state({
		accepted: [],
		rejected: []
	});
	let link_input_field = $state();
	let symbol_map = $state(symbols.symbols);
	let symbol_links = $state(symbols.symbolLinks);
	let versions = $derived(Object.keys(symbol_map));
	let nr_versions_provided = $derived(versions.length);
	let selected_primary_versions = $state(symbols.selected_version);
	let selected_secondary_versions = $state(symbols.selected_versions_to_compare);

	function loadReportsfromJson(jsonReports) {
		symbols.symbols = JSON.parse(jsonReports);
	}

	function loadReportFromELF(symID, elfBinary) {
		console.log(elfBinary);
		const go = new Go(); // Defined in wasm_exec.js
		const WASM_URL = base + '/sELFperf.wasm';
		var wasm;
		fetch(WASM_URL)
			.then((resp) => resp.arrayBuffer())
			.then((bytes) =>
				WebAssembly.instantiate(bytes, go.importObject).then(function (obj) {
					wasm = obj.instance;
					go.run(wasm);
					let reportJSONstr = get_elf_report(elfBinary);
					let reportJSON = JSON.parse(reportJSONstr);
					console.log(reportJSON);
					// end TODO :)
					if (reportJSON.hasOwnProperty('singlefirmware')) {
						let disasmFnMap = helpers.getDisasmFnMap(reportJSON);
						let disasmFnMapArg = Uint8Array.fromBase64(btoa(JSON.stringify(disasmFnMap)));
						reportJSONstr = add_fn_calls_from_disasm(disasmFnMapArg);
						reportJSON = JSON.parse(reportJSONstr);
						// Todo mv somewhere better
						let symPathByAddr = {};
						let symPathByName = {};
						let reportFns = reportJSON['functions'];
						let reportVars = reportJSON['variables'];
						for (let i = 0; i < reportFns.length; i++) {
							let addr = reportFns[i]['address'];
							// TODO what about syms with unknown path - can they be eliminated ^^?
							let urlPath = reportFns[i]['file'] + '/' + reportFns[i]['name'];
							reportFns[i]['symtype'] = 'fn';
							symPathByAddr[addr] = urlPath;
							symPathByName[urlPath] = reportFns[i];
						}
						for (let i = 0; i < reportVars.length; i++) {
							reportVars[i]['symtype'] = 'var';
						}
						reportJSON['SymPathByAddr'] = symPathByAddr;
						reportJSON['symPathByName'] = symPathByAddr;
						reportJSON['goWasmLoaded'] = true;
						// TODO how much space saving it pre-calculated?
						console.log(reportJSON);
						symbols.symbols[symID] = reportJSON;
						symbols.elfDataProvided = true;
					} else {
					}
				})
			);
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
					let symID = files.accepted[files.accepted.length - 1].name + '_' + nr_versions_provided;
					loadReportFromELF(symID, new Uint8Array(reader.result));
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

	async function addFirmwareByLink(name, link) {
		if (!symbol_links.includes(link)) {
			const response = await fetch(link);
			if (link.endsWith('.elf')) {
				loadReportFromELF(name, await response.bytes());
				return;
			}
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

	function addFWSample(name, url) {
		addFirmwareByLink(name, url);
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
					<Col sm="12" md={8}>
						<CardSubtitle><b>...by file (easiest):</b></CardSubtitle>
						<CardText>
							<div class="uploadfield">
								<div>
									Uploaded files never leave your browser, processing happens locally. Refreshing
									the site resets everything.
								</div>

								<br />

								<Dropzone on:drop={handleFilesSelect} accept=".json,.elf" style="min-height:200px">
									<p>Drag 'n' drop ELF files or JSON-report here, or click to select files.</p>
									<p>ELF files must be build with debug info (-g) to get best results.</p>
								</Dropzone>
							</div>
						</CardText>
					</Col>
					<Col sm="12" md={4}>
						<Row>
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
						</Row>
						<Row>
							<CardSubtitle style="margin-top:10px;"><b>...OR by loading a sample:</b></CardSubtitle
							>
							<CardText
								>Do not have any and just want to see a demo? <br /> Then load a sample to see some features
								:)</CardText
							>
							<div>
								{#each testFW as fw, i ('link-' + fw.name)}
									<Button color="light" onclick={() => addFWSample(fw.name, fw.url)}
										>{fw.name}</Button
									>
								{/each}
							</div>
						</Row>
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
					<Col sm="12" md={4}></Col>
					<Col sm="12" md={4}>
						<!-- <a href="/#/wip">WIP - new feature playground :)</a> -->
					</Col>
				</Row>
			</CardFooter>
		</Card>
	</Container>
</div>

<style>
	#content {
		margin-top: 20px;
	}
</style>
