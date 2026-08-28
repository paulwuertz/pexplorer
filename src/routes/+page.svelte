<script>
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
		Toast,
		ToastBody,
		ToastHeader,
		Col
	} from '@sveltestrap/sveltestrap';
	import Dropzone from 'svelte-file-dropzone';

	import { symbols } from './symbols.svelte.js';
	import * as helpers from './helpers.js';
	import { onMount } from 'svelte';
	import { base } from '$app/paths';
	import { text } from '@sveltejs/kit';
	import { color } from 'echarts';

	const testFWsets = [
		{
			common_name: 'CANnectivity - NXP LPC55S16',
			versions: [
				{
					name: 'v1.2 - LLVM',
					url: 'https://media.githubusercontent.com/media/paulwuertz/pexplorer/refs/heads/main/testdata/elf_testdata/zephyr_cannectivity_12_llvm_lpc55s16.elf'
				},
				{
					name: 'v1.2 - GCC',
					url: 'https://media.githubusercontent.com/media/paulwuertz/pexplorer/refs/heads/main/testdata/elf_testdata/zephyr_cannectivity_12_gcc_lpc55s16.elf'
				},
				{
					name: 'v1.3 - LLVM',
					url: 'https://media.githubusercontent.com/media/paulwuertz/pexplorer/refs/heads/main/testdata/elf_testdata/zephyr_cannectivity_13_llvm_lpc55s16.elf'
				},
				{
					name: 'v1.3 - GCC',
					url: 'https://media.githubusercontent.com/media/paulwuertz/pexplorer/refs/heads/main/testdata/elf_testdata/zephyr_cannectivity_13_gcc_lpc55s16.elf'
				},
				{
					name: 'v1.4 - GCC',
					url: 'https://media.githubusercontent.com/media/paulwuertz/pexplorer/refs/heads/main/testdata/elf_testdata/cannectivity_lpcxpresso55s16_1.4.elf'
				}
			]
		},
		{
			common_name: 'CANnectivity 1.4 - GCC',
			versions: [
				{
					name: 'nucleo_h723zg',
					url: 'https://media.githubusercontent.com/media/paulwuertz/pexplorer/refs/heads/main/testdata/elf_testdata/cannectivity_nucleo_h723zg_1.4.elf'
				},
				{
					name: 'frdm_mcxn947',
					url: 'https://media.githubusercontent.com/media/paulwuertz/pexplorer/refs/heads/main/testdata/elf_testdata/cannectivity_frdm_mcxn947_1.4.elf'
				},
				{
					name: 'stm32g0b1xx',
					url: 'https://media.githubusercontent.com/media/paulwuertz/pexplorer/refs/heads/main/testdata/elf_testdata/cannectivity_candlelightfd_stm32g0b1xx_dual_1.4.elf'
				},
				{
					name: 'same70n20b',
					url: 'https://media.githubusercontent.com/media/paulwuertz/pexplorer/refs/heads/main/testdata/elf_testdata/cannectivity_canbardo_same70n20b_1.4.elf'
				}
			]
		},
		{
			common_name: 'ZSWatch - Legacy',
			versions: [
				{
					name: 'v0.7.0',
					url: 'https://media.githubusercontent.com/media/paulwuertz/pexplorer/refs/heads/main/testdata/elf_testdata/zswatch_nrf5340_070.elf'
				},
				{
					name: 'v0.8.0',
					url: 'https://media.githubusercontent.com/media/paulwuertz/pexplorer/refs/heads/main/testdata/elf_testdata/zswatch_nrf5340_080.elf'
				},
				{
					name: 'v0.8.1',
					url: 'https://media.githubusercontent.com/media/paulwuertz/pexplorer/refs/heads/main/testdata/elf_testdata/zswatch_nrf5340_081.elf'
				}
			]
		},
		{
			common_name: 'Pinecil IronOS (RISC-V WiP)',
			versions: [
				{
					name: 'V1 EN v2.23',
					url: 'https://media.githubusercontent.com/media/paulwuertz/pexplorer/refs/heads/main/testdata/elf_testdata/Pinecilv1_EN_v2_23.elf'
				},
				{
					name: 'V2 EN v2.23',
					url: 'https://media.githubusercontent.com/media/paulwuertz/pexplorer/refs/heads/main/testdata/elf_testdata/Pinecilv2_EN_v2_23.elf'
				}
			]
		}
	];
	const testFW = [
		{
			name: 'Prusa Buddy - Core One v6.4.0',
			url: 'https://media.githubusercontent.com/media/paulwuertz/pexplorer/refs/heads/main/testdata/elf_testdata/prusa_buddy_boot_64.elf'
		},
		{
			name: 'Libresolar BMS vX.Y.Z',
			url: 'https://media.githubusercontent.com/media/paulwuertz/pexplorer/refs/heads/main/testdata/elf_testdata/libre_solar_zephyr_43_b953.elf'
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
	let toast_messages = $state([]);

	function loadReportsfromJson(jsonReports) {
		try {
			let reports = JSON.parse(jsonReports);
			if (Object.hasOwn(reports, 'singlefirmware')) {
				for (const r of Object.keys(reports)) {
					[reports[r]['SymPathByAddr'], reports[r]['symPathByName']] = helpers.fn2symPathLookups(
						reports[r]['functions']
					);
					console.log('SymPathByAddr:', reports[r]['SymPathByAddr']);

					symbols.symbols[r] = reports[r];
				}
			} else if (Object.hasOwn(reports, 'multifirmware')) {
				const startTime = performance.now();
				for (const r of Object.keys(reports.symbols)) {
					console.log('add:', r, 'from multi-json report');
					symbols.symbols[r] = reports.symbols[r];
				}
				const endTime = performance.now();
				toast_messages.push({
					title: 'Finished. You can start to explore what we found on your program :)',
					text: `Loading the report took ${(endTime - startTime) / 1000} seconds.`
				});
				toast_messages.push({
					title: 'Note',
					color: 'danger',
					text:
						'JSON export for call graphs is not fully implemented yet.<br>' +
						'Upload your ELF if you want to see fully analysis. Fix comes soon :)'
				});
			}
		} catch (error) {
			// TODO UI error
			console.log('error from JSON report:', error);
			return;
		}
		symbols.elfDataProvided = true;
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
					const startTime1 = performance.now();
					let reportJSONstr = get_elf_report(elfBinary);
					const endTime1 = performance.now();
					console.log(`Report generation took ${(endTime1 - startTime1) / 1000} seconds`);
					let reportJSON = JSON.parse(reportJSONstr);
					console.log(reportJSON);
					// end TODO :)
					if (reportJSON.hasOwnProperty('singlefirmware')) {
						const startTime = performance.now();
						let disasmFnMap = helpers.getDisasmFnMap(reportJSON);
						// hacky bigint to int addr...
						for (let f of Object.values(disasmFnMap)) {
							for (let i of f) {
								i.addr = parseInt(i.addr);
							}
						}
						let disasmFnMapArg = Uint8Array.fromBase64(btoa(JSON.stringify(disasmFnMap)));
						reportJSONstr = add_fn_calls_from_disasm(disasmFnMapArg);
						const endTime = performance.now();
						toast_messages = [];
						toast_messages.push({
							title: 'Finished. You can start to explore what we found on your program :)',
							text: `Report generation took ${(endTime1 - startTime1) / 1000} seconds. <br>
                                    Call annotation took ${(endTime - startTime) / 1000} seconds`
						});
						console.log(`After call annotation it took ${(endTime - startTime) / 1000} seconds`);
						reportJSON = JSON.parse(reportJSONstr);
						// Todo mv somewhere better
						let reportFns = reportJSON['functions'];
						let reportVars = reportJSON['variables'];
						for (let i = 0; i < reportFns.length; i++) {
							// TODO what about syms with unknown path - can they be eliminated ^^?
							reportFns[i]['symtype'] = 'fn';
						}
						for (let i = 0; i < reportVars.length; i++) {
							reportVars[i]['symtype'] = 'var';
						}
						[reportJSON['SymPathByAddr'], reportJSON['symPathByName']] =
							helpers.fn2symPathLookups(reportFns);
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
			toast_messages.push({
				title: 'Firmware analysis running...',
				text: 'Fetching firmware... <br>Start extracting data from the ELF file... <br>Analysing the call data... <br> Be a little patient please :)'
			});
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
			toast_messages.push({
				title: 'Nothing to do',
				text: 'Link already added :)'
			});
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

	function downloadLinks() {
		helpers.download(
			'report.json',
			JSON.stringify({
				multifirmware: true,
				...symbols
			})
		);
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

	onMount(() => {
		MCapstone().then((cs) => {
			helpers.loadCapstone(cs);
		});
	});
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
					<Col sm="12" md={7}>
						<CardSubtitle><b>...by uploading your ELF file:</b></CardSubtitle>
						<CardText>
							<div class="uploadfield">
								<br />

								<Dropzone on:drop={handleFilesSelect} accept=".json,.elf" style="min-height:200px">
									<svg viewBox="0 0 100 100" width="25%" style="margin-bottom: 20px;">
										<g fill="#bdbdbd">
											<path
												d="m69.271 42.085c-5.941-7.138-11.883-14.276-17.824-21.414-0.437-0.524-0.951-0.715-1.447-0.681-0.496-0.034-1.011 0.157-1.447 0.681-5.941 7.137-11.883 14.276-17.824 21.414-1.1 1.319-0.461 3.494 1.446 3.494h7.479v32.386c0 1.116 0.931 2.047 2.047 2.047h16.598c1.116 0 2.047-0.931 2.047-2.047v-32.386h7.479c1.907 0 2.546-2.175 1.446-3.494z"
											/>
											<path
												d="m50 0c-27.614 0-50 22.386-50 50s22.386 50 50 50 50-22.386 50-50-22.386-50-50-50zm0 92c-23.196 0-42-18.805-42-42 0-23.196 18.804-42 42-42 23.195 0 42 18.804 42 42 0 23.195-18.805 42-42 42z"
											/>
										</g>
									</svg>
									<p>Drag 'n' drop ELF files or JSON-report here, or click to select files.</p>
									<p>Uploaded files never leave your browser, processing happens locally.</p>
									<p>Refreshing the site resets everything.</p>
									<p>ELF files must be build with debug info (-g) to get best results.</p>
								</Dropzone>
							</div>
						</CardText>
					</Col>
					<Col sm="12" md={5}>
						<!-- <Row>
							<CardSubtitle><b>...by link:</b></CardSubtitle>

							<CardText>
								<br />
								Adding the symbol via links saves them in your browsers local storage so you can continue
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
						</Row> -->
						<Row>
							<CardSubtitle><b>...OR load a sample:</b></CardSubtitle
							>
							<CardText
								>Do not have any and just want to see a demo? <br /> Then load a sample to see some features
								:)</CardText
							>
							<div>
								{#each testFWsets as fws, i ('fwset-' + fws.common_name)}
									<div><b><small>{fws.common_name}:</small></b></div>
									<div class="horizontal capitalize">
										{#each fws.versions as fw, i ('fw-version-' + fw.name)}
											<Button
												class="m-1"
												size="sm"
												color="light"
												onclick={() => addFWSample(fws.common_name + fw.name, fw.url)}
												>{fw.name}</Button
											>
										{/each}
									</div>
								{/each}
								<div><b><small>Individual firmware samples:</small></b></div>
								<div class="horizontal capitalize">
									{#each testFW as fw, i ('link-' + fw.name)}
										<Button
											size="sm"
											class="m-1"
											color="light"
											onclick={() => addFWSample(fw.name, fw.url)}>{fw.name}</Button
										>
									{/each}
								</div>
							</div>
						</Row>
					</Col>
				</Row>
			</CardBody>
			<CardFooter>
				<Row>
					<Col sm="12" md={5}>
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

						<Button size="md" color="success" onclick={downloadLinks}>Download json report</Button>
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

<div class="toast-container position-fixed bottom-0 end-0 p-3">
	<div class="p-3 mb-3">
		{#each toast_messages as toast_message, i (toast_message.title)}
			<Toast class="me-1 mb-3">
				<ToastHeader icon={toast_message.color}>{toast_message.title}</ToastHeader>
				{#if Object.hasOwn(toast_message, 'text')}
					<ToastBody>
						{@html toast_message.text}
					</ToastBody>
				{/if}
			</Toast>
		{/each}
	</div>
</div>

<style>
	#content {
		margin-top: 20px;
	}
	:global(.example-btn) {
		width: 100%;
		margin-top: 5px;
	}
</style>
