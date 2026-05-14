<script>
	import { onMount } from 'svelte';
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import { browser } from '$app/environment';
	import { writable } from 'svelte/store';
	// ui stuff
	import { Badge, Button, Col, Container, Input, Row, Table } from '@sveltestrap/sveltestrap';

	import { symbols } from '../../symbols.svelte.js';
	import * as helpers from '../../helpers.js';
	import VersionSelect from '../../../components/VersionSelect.svelte';

	import * as echarts from 'echarts';

	const PLOT_ID_PREFIX = 'stackSizeLayered_';

	let { params, data } = $props();
	let version_name = $derived((params.version || ""));
	let version = $derived(symbols.symbols[version_name] || {});
	let variables = $derived(version.variables || []);

	let typeStruct = $derived(
		variableTypes.find((val, i, arr) => {
			return val.name == "_static_thread_data";
		})
	);

	const isStaticThread = (symbol) => {
		let secidx = symbol.secidx || null;
		let sections = version.sections || [];
		let section = sections[secidx] || {};
		let section_name = section.name || "";
        return section_name == '_static_thread_data_area';
	};

    // TODO mv to helpers
	let getMemberBytes = (data, byte_offset, size) => {
		let b = [];
		let val = 0;
		for (let i = byte_offset; i < data.length && i < byte_offset + size; i++) {
			b.push(data[i]);
			val += data[i] * Math.pow(256, i - byte_offset);
		}
		return val;
	};

	let staticThreads = $derived(variables.filter(isStaticThread));
	let variableTypes = $derived(version['types']);
    $inspect(staticThreads)
    let static_thread_args = $derived(typeStruct && typeStruct.members || [])
    $inspect(typeStruct, static_thread_args)
    let static_thread_data = $derived(staticThreads.map((t, i, a) => {
        let staticInitDataBase64 = t.staticInitData;
        let staticInitData = (
            (staticInitDataBase64 && Uint8Array.fromBase64(staticInitDataBase64)) || undefined
        );
        let thread_data = {}
        static_thread_args.forEach((member) => {
            thread_data[member.name] = getMemberBytes(staticInitData, member.byte_offset, member.size)
        });
        thread_data["name"] = t["name"]
        return thread_data;
    }));
    $inspect(typeStruct, static_thread_args, static_thread_data)
</script>

<div class="container" id="content">

{#if !version || version_name == ""}
    <VersionSelect page_name={"rtos"}></VersionSelect>
{:else}
    <Container fluid>
        <h3>version_name: {JSON.stringify(version_name)}</h3>
        <h4>Static threads by K_THREAD_DEFINE</h4>
        <ul>
        {#each static_thread_data as sTread (sTread.name)}
            <li>{sTread["name"]} -
                <ul>
                {#each Object.entries(sTread) as s}
                    <li>{s[0]} - {s[1]}</li>
                {/each}
                </ul>
            </li>
        {/each}
        </ul>
    </Container>
{/if }
</div>
