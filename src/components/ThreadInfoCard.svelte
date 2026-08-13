<script>
	import { onMount } from 'svelte';
	import { base } from '$app/paths';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';

	import {
		Card,
		Button,
		ButtonGroup,
		Col,
		CardHeader,
		Container,
		Icon,
		Input,
		Progress,
		Row,
		CardBody,
		CardSubtitle,
		CardText,
		CardTitle,
		Table
	} from '@sveltestrap/sveltestrap';
	import { symbols } from '../routes/symbols.svelte.js';
	import * as helpers from '../routes/helpers.js';
	import Dropzone from 'svelte-file-dropzone';

	const { sTread, version_name } = $props();
	let thread = $derived(sTread);

	let stackLevelToColor = (stackUse, stackSize) => {
		let percentage = (100.0 * stackUse) / stackSize;
		if (percentage < 50.0) {
			return '';
		}
		if (percentage < 80.0) {
			return 'success';
		}
		if (percentage < 100.0) {
			return 'warning';
		} else {
			return 'danger';
		}
	};

	let versions = $derived(Object.keys(symbols.symbols));
	let symInOtherVersions = $derived(
		Object.entries(symbols.symbols).map((v, i) => {
			let version = v[0];
			let fw_report = v[1];
			if (version == version_name) return null;
			for (let symIndex = 0; symIndex < fw_report.functions.length; symIndex++) {
				const f = fw_report.functions[symIndex];
				if (f.name == thread.name) {
					console.log('found ', f, 'in version');
				}
			}
		})
	);
</script>

<Card>
	<CardHeader>
		<CardTitle>{thread['name']}</CardTitle>
	</CardHeader>
	<CardBody>
		<CardText>
			<!-- TODO add source link symbol json -->
			<!-- Buildtime: {symbols.symbols[version].timestamp} -->
			<div class="pb-3">
				<b>Thread entry function:</b>
				<br />
				<span>
					<a
						data-sveltekit-preload-data="tap"
						href={'#/browse/' + version_name + thread['file'] + '/' + thread['name']}
					>
						{thread['name']}
					</a>
				</span>
			</div>
			<div class="pb-3">
				<b>Configured stack size:</b>
				<br />
				<span>init_stack_size bytes</span>
			</div>
			<div class="pb-3">
				<b>Unresolved dynamic calls:</b>
				<br />
				<span
					>{thread['from_nr_functions']} function contain {thread['unresolved_calls']} unresolved calls
					over the currently known calltree of the task.</span
				>
			</div>
			<!-- <div class="pb-3">
                <b>Functions missing stack-use info:</b>
                <br />
                <span>TODO</span>
            </div> -->

			<CardSubtitle>Minimum stack use scenario found:</CardSubtitle>
			<div class="pt-3">
				<Progress
					color={stackLevelToColor(thread['max_stack_size_callees'], thread['init_stack_size'])}
					value={thread['max_stack_size_callees']}
					max={thread['init_stack_size']}
					class="mb-2"
				>
					{((100 * thread['max_stack_size_callees']) / thread['init_stack_size']).toFixed(2)}% - {thread[
						'max_stack_size_callees'
					]} / {thread['init_stack_size']} bytes
				</Progress>
			</div>
			<div class="pt-3">
				<CardSubtitle class="pb-3">in "CANnectivity 1.4 - GCC stm32g0b1xx":</CardSubtitle>
				<Progress
					color={stackLevelToColor(thread['max_stack_size_callees'], thread['init_stack_size'])}
					value={816}
					max={thread['init_stack_size']}
					class="mb-2"
				>
					{((100 * thread['max_stack_size_callees']) / thread['init_stack_size']).toFixed(2)}% - {thread[
						'max_stack_size_callees'
					]} / {thread['init_stack_size']} bytes
				</Progress>
			</div>
			<div class="pt-3">
				<CardSubtitle class="pb-3">in "CANnectivity 1.4 - GCC LPC55S16v16"</CardSubtitle>
				<Progress
					color={stackLevelToColor(892, thread['init_stack_size'])}
					value={892}
					max={thread['init_stack_size']}
					class="mb-2"
				>
					{((100 * thread['max_stack_size_callees']) / thread['init_stack_size']).toFixed(2)}% - {892}
					/ {thread['init_stack_size']} bytes
				</Progress>
			</div>
		</CardText>
	</CardBody>
</Card>
