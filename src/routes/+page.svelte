
<script lang="ts">
    import { onMount } from 'svelte';
    import { base } from '$app/paths';
    import { page } from '$app/stores';
    import { browser } from '$app/environment';
    import { writable } from "svelte/store";
    // ui stuff
    import { Badge, Col, Container, Input, Row } from '@sveltestrap/sveltestrap';

	import { symbols } from './symbols.svelte.js';

    import { cubicOut } from 'svelte/easing';
    //import { hierarchy, type HierarchyNode, type HierarchyRectangularNode } from 'd3-hierarchy';
    import { scaleSequential, scaleOrdinal } from 'd3-scale';
    import * as chromatic from 'd3-scale-chromatic';
    import { hsl } from 'd3-color';

    import { Arc, Bounds, Chart, Partition, Svg, Tooltip, findAncestor } from 'layerchart';

    import { Breadcrumb, Button, Field, ToggleGroup, ToggleOption } from 'svelte-ux';
    import { format, sortFunc, compoundSortFunc } from '@layerstack/utils';

    //import Preview from '$lib/docs/Preview.svelte';

	let { data } = $props();

    // const complexHierarchy = hierarchy(data.flare)
    //     .sum((d) => d.value)
    //     .sort(compoundSortFunc(sortFunc('height', 'desc'), sortFunc('value', 'desc')));

    // let colorBy = 'parent';

    //let selected: HierarchyRectangularNode<any> = complexHierarchy as HierarchyRectangularNode<any>; // select root initially

    const sequentialColor = scaleSequential([4, -1], chromatic.interpolateGnBu);
    // filter out hard to see yellow and green
    const ordinalColor = scaleOrdinal(
        chromatic.schemeSpectral[9].filter((c) => hsl(c).h < 60 || hsl(c).h > 90)
    );
    // const ordinalColor = scaleOrdinal(chromatic.schemeCategory10)

    // function getNodeColor(node: HierarchyNode<any>, colorBy: string) {
    //     switch (colorBy) {
    //     case 'children':
    //         return node.children ? '#ccc' : '#ddd';
    //     case 'depth':
    //         return sequentialColor(node.depth).toString();
    //     case 'parent':
    //         const colorParent = findAncestor(node, (n) => n.depth === 1);
    //         return colorParent
    //         ? hsl(ordinalColor(colorParent.data.name))
    //             .brighter(node.depth * 0.3)
    //             .toString()
    //         : '#ddd';
    //     }
    //     return '';
    // }

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

        <!-- <Breadcrumb items={selected?.ancestors().reverse() ?? []}>
            <Button
              slot="item"
              let:item
              on:click={() => (selected = item)}
              base
              class="px-2 py-1 rounded"
            >
              <div class="text-left">
                <div class="text-sm">{item.data.name}</div>
                <div class="text-xs text-surface-content/50">
                  {format(item.value ?? 0, "integer")}
                </div>
              </div>
            </Button>
          </Breadcrumb>
          <div class="h-[600px] p-4 border rounded">
            <Chart data={complexHierarchy} let:tooltip>
              <Svg center>
                <Bounds
                  domain={{
                    x0: selected?.x0 ?? 0,
                    x1: selected?.x1 ?? 1,
                    y0: selected?.y0 ?? 0,
                    y1: 1,
                  }}
                  range={({ height }) => ({
                    x0: 0,
                    x1: 2 * Math.PI,
                    y0: selected?.y0 ? 20 : 0,
                    y1: height / 2,
                  })}
                  tweened={{ duration: 800, easing: cubicOut }}
                  let:xScale
                  let:yScale
                >
                  <Partition size={[1, 1]} let:nodes>
                    {#each nodes as node}
                      {@const nodeColor = getNodeColor(node, colorBy)}
                      <Arc
                        value={node.value}
                        startAngle={Math.max(0, Math.min(2 * Math.PI, xScale(node.x0)))}
                        endAngle={Math.max(0, Math.min(2 * Math.PI, xScale(node.x1)))}
                        innerRadius={Math.max(0, yScale(node.y0))}
                        outerRadius={Math.max(0, yScale(node.y1))}
                        fill={nodeColor}
                        _stroke={hsl(nodeColor).darker(colorBy === "children" ? 0.5 : 2)}
                        stroke="hsl(0 0% 20%)"
                        class="cursor-pointer"
                        let:centroid
                        onclick={() => {
                          selected = node;
                        }}
                        onpointermove={(e) => tooltip.show(e, node)}
                        onpointerleave={tooltip.hide}
                      >
                        <text x={centroid[0]} y={centroid[1]}>{node.data.name}</text>
                      </Arc>
                    {/each}
                  </Partition>
                </Bounds>
              </Svg>

              <Tooltip.Root let:data>
                <Tooltip.Header>{data.data.name}</Tooltip.Header>
                <Tooltip.List>
                  <Tooltip.Item label="value" value={data.value} format="integer" />
                </Tooltip.List>
              </Tooltip.Root>
            </Chart>
          </div> -->

    </Container>
</div>



