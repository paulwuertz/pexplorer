<script>
    import { onMount } from 'svelte';
    import { base } from '$app/paths';
    import { page } from '$app/stores';
    import { browser } from '$app/environment';
    import { writable } from "svelte/store";
    // ui stuff
    import { DataTable } from '@careswitch/svelte-data-table';
    import { Badge, Button, Col, Container, Input, Row, Table } from '@sveltestrap/sveltestrap';

	import { symbols } from '../symbols.svelte.js';

	import * as echarts from 'echarts';

	let { data } = $props();
    let files = $state();
    let versions = $derived(Object.keys(symbols.symbols));
    let selected_symbols = $state({});

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

    let get_all_threads = (allSymVersions) => {
        let threads = new Set();
            console.log(allSymVersions)
        for(let symVersion of Object.values(allSymVersions)){
            for(let thread_name of Object.keys(symVersion["stack_reports"])){
                threads.add(thread_name);
            }
        }
        // alert(threads)
        return threads;
    }

    /**
     * Gets all versions ordered by timestamp
     * @param allSymVersions
     */
    let get_versions_ordered_by_timestamps = (allSymVersions) => {
        let versions = [];
        // get all versions
        for(const [versionStr, symVersion] of Object.entries(allSymVersions)){
            const timestamp = new Date(symVersion["timestamp"]);
            versions.push({"version": versionStr, "timestamp": timestamp})
        }
        versions.sort(function(a, b) {
            return a["timestamp"] > b["timestamp"];
        });
        return versions;
    }

    let get_max_stack_sizes_of_thread = (allSymVersions, threadname) => {
        let threads = new Set();
        for(let symVersion of Object.values(allSymVersions)){
            for(let thread_name of Object.keys(symVersion["stack_reports"])){
                threads.add(thread_name);
            }
        }
        // alert(threads)
        return threads;
    }

    onMount(async () => {
        if (browser) {
            let trend_data = versions;
            let threads = get_all_threads(symbols.symbols)
            let ordered_versions_and_timestamps = get_versions_ordered_by_timestamps(symbols.symbols);
            let ordered_versions = ordered_versions_and_timestamps.map( vt => vt["version"])
            let ordered_timestamps = ordered_versions_and_timestamps.map( vt => vt["timestamp"])

            let series_maxstack_data = [];
            for(let thread of threads){
                let maxstack_data = [];
                for(let v of ordered_versions){
                    let max_stack_size = symbols.symbols[v]["stack_reports"][thread]["max_static_stack_size"]
                    maxstack_data.push(max_stack_size)
                }
                console.log(thread, maxstack_data)
                series_maxstack_data.push({
                    name: thread,
                    type: 'line',
                    stack: 'Total',
                    data: maxstack_data
                })
            }
            console.log("dana ",ordered_versions,ordered_timestamps)
            // Create the echarts instance
            var myChart = echarts.init(
                    document.getElementById('main'),
                {width: "100%", height: 600}
            );

            // Draw the chart
            myChart.setOption({
                title: {
                    text: 'Stacked Line'
                },
                tooltip: {
                    trigger: 'axis'
                },
                legend: {
                    data: threads
                },
                grid: {
                    left: '3%',
                    right: '4%',
                    bottom: '3%',
                    containLabel: true
                },
                toolbox: {
                    feature: {
                    saveAsImage: {}
                    }
                },
                xAxis: {
                    type: 'category',
                    boundaryGap: false,
                    data: ordered_timestamps
                },
                yAxis: {
                    type: 'value'
                },
                series: series_maxstack_data
            });
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
    <div id="main" style="width: 100%;height:600px;"></div>
</div>



