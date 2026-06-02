<script>
	import { onMount } from 'svelte';
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import { browser } from '$app/environment';
	import { writable } from 'svelte/store';
	// ui stuff
	import { Badge, Button, Col, Container, Input, Row, Table,  TabContent, TabPane } from '@sveltestrap/sveltestrap';

	import { symbols } from '../symbols.svelte.js';
	import * as helpers from '../helpers.js';

	import * as echarts from 'echarts';
	import SettingsEditor from '../../components/SettingsEditor.svelte';

	const PLOT_ID_PREFIX = 'stackSizeLayered_';

	let { data } = $props();
	let versions = $derived(Object.keys(symbols.symbols));
	let selected_thread_stat = $state({});
    let settings_dummy = {"threads":[{"thread_entry_name":"log_process_thread_func","stack_variable_name":"logging_stack"},{"thread_entry_name":"bg_thread_main","stack_variable_name":"z_main_stack"},{"thread_entry_name":"gs_usb_rx_thread","size":1024},{"thread_entry_name":"gs_usb_tx_thread","size":1024},{"thread_entry_name":"work_queue_main","stack_variable_name":"sys_work_q_stack"}],"dynamic_calls":[{"call_from":"work_queue_main","call_to":"led_event_triggered_work_handler"},{"call_from":"work_queue_main","call_to":"cannectivity_usb_reboot"},{"call_from":"work_queue_main","call_to":"dfu_button_poll"}]}

</script>

<div class="container" id="content">
	<hr />

	<Container fluid>
    <Button color="success">Add Config</Button>
    <TabContent vertical pills>
        <TabPane tabId="alpha" tab="Alpha asdasdasdc" active>
            <SettingsEditor settings={settings_dummy}></SettingsEditor>
        </TabPane>
        <TabPane tabId="bravo" tab="Bravo">
            <h2>Bravo</h2>
            <img
            alt="Johnny Bravo"
            src="https://upload.wikimedia.org/wikipedia/commons/thumb/b/be/Johnny_Bravo_series_logo.png/320px-Johnny_Bravo_series_logo.png"
            />
        </TabPane>
        <TabPane tabId="charlie" tab="Charlie">
            <h2>Charlie</h2>
            <img alt="Charlie Brown" src="https://upload.wikimedia.org/wikipedia/en/2/22/Charlie_Brown.png" />
        </TabPane>
        </TabContent>
	</Container>
</div>
