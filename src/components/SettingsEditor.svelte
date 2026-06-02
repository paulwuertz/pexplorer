<script>
	import { onMount } from 'svelte';
	import { base } from '$app/paths';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';

	import {
		Card,
		Button,
		Col,
		CardHeader,
		Container,
        Icon,
		Input,
		Row,
		CardBody,
		CardText,
		CardTitle,
        Table,
	} from '@sveltestrap/sveltestrap';
	import { symbols } from '../routes/symbols.svelte.js';
	import * as helpers from '../routes/helpers.js';

	const { settings } = $props();

    let thread_stack_size = (thread) => {
        if (Object.hasOwn(thread, "stack_variable_name")) {
            return 1024 // TODO lookup symbolname + size
        } else if (Object.hasOwn(thread, "size")) {
            return thread.size
        } else {
            return 0 // TODO error
        }
    }

    let threads = settings.threads || []
    let dynamic_calls = settings.dynamic_calls || []
	let versions = $derived(Object.keys(symbols.symbols));
</script>

<div class="container" id="content">
	<h3>pexplorer setting:</h3>
	<h4>threads:</h4>

    <Table bordered>
    <thead>
        <tr>
            <th>Thread entry function name</th>
            <th>Thread stack variable name</th>
            <th>Thread stack size</th>
            <th></th>
        </tr>
    </thead>
    <tbody>
        {#each threads as thread (thread.thread_entry_name)}
        <tr>
            <td>{thread.thread_entry_name}</td>
            <td>{thread.stack_variable_name}</td>
            <td>{thread_stack_size(thread)}</td>
            <td>
                <Button size="md" outline color="danger">
                    <svg height="32" width="29" viewBox="0 0 1000 875"><path d="M0 281.296l0 -68.355q1.953 -37.107 29.295 -62.496t64.449 -25.389l93.744 0l0 -31.248q0 -39.06 27.342 -66.402t66.402 -27.342l312.48 0q39.06 0 66.402 27.342t27.342 66.402l0 31.248l93.744 0q37.107 0 64.449 25.389t29.295 62.496l0 68.355q0 25.389 -18.553 43.943t-43.943 18.553l0 531.216q0 52.731 -36.13 88.862t-88.862 36.13l-499.968 0q-52.731 0 -88.862 -36.13t-36.13 -88.862l0 -531.216q-25.389 0 -43.943 -18.553t-18.553 -43.943zm62.496 0l749.952 0l0 -62.496q0 -13.671 -8.789 -22.46t-22.46 -8.789l-687.456 0q-13.671 0 -22.46 8.789t-8.789 22.46l0 62.496zm62.496 593.712q0 25.389 18.553 43.943t43.943 18.553l499.968 0q25.389 0 43.943 -18.553t18.553 -43.943l0 -531.216l-624.96 0l0 531.216zm62.496 -31.248l0 -406.224q0 -13.671 8.789 -22.46t22.46 -8.789l62.496 0q13.671 0 22.46 8.789t8.789 22.46l0 406.224q0 13.671 -8.789 22.46t-22.46 8.789l-62.496 0q-13.671 0 -22.46 -8.789t-8.789 -22.46zm31.248 0l62.496 0l0 -406.224l-62.496 0l0 406.224zm31.248 -718.704l374.976 0l0 -31.248q0 -13.671 -8.789 -22.46t-22.46 -8.789l-312.48 0q-13.671 0 -22.46 8.789t-8.789 22.46l0 31.248zm124.992 718.704l0 -406.224q0 -13.671 8.789 -22.46t22.46 -8.789l62.496 0q13.671 0 22.46 8.789t8.789 22.46l0 406.224q0 13.671 -8.789 22.46t-22.46 8.789l-62.496 0q-13.671 0 -22.46 -8.789t-8.789 -22.46zm31.248 0l62.496 0l0 -406.224l-62.496 0l0 406.224zm156.24 0l0 -406.224q0 -13.671 8.789 -22.46t22.46 -8.789l62.496 0q13.671 0 22.46 8.789t8.789 22.46l0 406.224q0 13.671 -8.789 22.46t-22.46 8.789l-62.496 0q-13.671 0 -22.46 -8.789t-8.789 -22.46zm31.248 0l62.496 0l0 -406.224l-62.496 0l0 406.224z"/></svg>
                </Button>
            </td>
        </tr>
        {/each}

        <tr>
            <td>
                <Input type="select">
                    <option></option>
                    {#each [1, 2, 3, 4, 5] as option}
                    <option>{option}</option>
                    {/each}
                </Input>
            </td>
            <td>
                <Input type="select">
                    <option></option>
                    {#each [1, 2, 3, 4, 5] as option}
                    <option>{option}</option>
                    {/each}
                </Input>
            </td>
            <td>
                <Input type="number"/>
            </td>
            <td>
                <Button color="success">Add thread</Button>
            </td>
        </tr>
    </tbody>
    </Table>

	<h4>Dynamic calls:</h4>


    <Table bordered>
    <thead>
        <tr>
            <th>From</th>
            <th>To</th>
            <th></th>
        </tr>
    </thead>
    <tbody>
        {#each dynamic_calls as dynamic_call, index (dynamic_call.call_from+dynamic_call.call_to+index)}
        <tr>
            <td>{dynamic_call.call_from}</td>
            <td>{dynamic_call.call_to}</td>
            <td>
                <Button size="md" outline color="danger">
                    <svg height="32" width="29" viewBox="0 0 1000 875"><path d="M0 281.296l0 -68.355q1.953 -37.107 29.295 -62.496t64.449 -25.389l93.744 0l0 -31.248q0 -39.06 27.342 -66.402t66.402 -27.342l312.48 0q39.06 0 66.402 27.342t27.342 66.402l0 31.248l93.744 0q37.107 0 64.449 25.389t29.295 62.496l0 68.355q0 25.389 -18.553 43.943t-43.943 18.553l0 531.216q0 52.731 -36.13 88.862t-88.862 36.13l-499.968 0q-52.731 0 -88.862 -36.13t-36.13 -88.862l0 -531.216q-25.389 0 -43.943 -18.553t-18.553 -43.943zm62.496 0l749.952 0l0 -62.496q0 -13.671 -8.789 -22.46t-22.46 -8.789l-687.456 0q-13.671 0 -22.46 8.789t-8.789 22.46l0 62.496zm62.496 593.712q0 25.389 18.553 43.943t43.943 18.553l499.968 0q25.389 0 43.943 -18.553t18.553 -43.943l0 -531.216l-624.96 0l0 531.216zm62.496 -31.248l0 -406.224q0 -13.671 8.789 -22.46t22.46 -8.789l62.496 0q13.671 0 22.46 8.789t8.789 22.46l0 406.224q0 13.671 -8.789 22.46t-22.46 8.789l-62.496 0q-13.671 0 -22.46 -8.789t-8.789 -22.46zm31.248 0l62.496 0l0 -406.224l-62.496 0l0 406.224zm31.248 -718.704l374.976 0l0 -31.248q0 -13.671 -8.789 -22.46t-22.46 -8.789l-312.48 0q-13.671 0 -22.46 8.789t-8.789 22.46l0 31.248zm124.992 718.704l0 -406.224q0 -13.671 8.789 -22.46t22.46 -8.789l62.496 0q13.671 0 22.46 8.789t8.789 22.46l0 406.224q0 13.671 -8.789 22.46t-22.46 8.789l-62.496 0q-13.671 0 -22.46 -8.789t-8.789 -22.46zm31.248 0l62.496 0l0 -406.224l-62.496 0l0 406.224zm156.24 0l0 -406.224q0 -13.671 8.789 -22.46t22.46 -8.789l62.496 0q13.671 0 22.46 8.789t8.789 22.46l0 406.224q0 13.671 -8.789 22.46t-22.46 8.789l-62.496 0q-13.671 0 -22.46 -8.789t-8.789 -22.46zm31.248 0l62.496 0l0 -406.224l-62.496 0l0 406.224z"/></svg>
                </Button>
            </td>
        </tr>
        {/each}

        <tr>
            <td>
                <Input type="select">
                    <option></option>
                    {#each [1, 2, 3, 4, 5] as option}
                    <option>{option}</option>
                    {/each}
                </Input>
            </td>
            <td>
                <Input type="select">
                    <option></option>
                    {#each [1, 2, 3, 4, 5] as option}
                    <option>{option}</option>
                    {/each}
                </Input>
            </td>
            <td>
                <Button color="success">Add dynamic call</Button>
            </td>
        </tr>
    </tbody>
    </Table>
</div>
