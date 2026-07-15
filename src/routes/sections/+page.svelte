<script>
	import { onMount } from 'svelte';
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import { browser } from '$app/environment';
	import { writable } from 'svelte/store';
	// ui stuff
	import {
		Badge,
		Button,
		ButtonGroup,
		Col,
		Container,
		Input,
		Row,
		Table,
		Card,
		CardHeader,
		CardBody,
		CardText,
		CardTitle,
		CardSubtitle,
		Progress,
		CardFooter
	} from '@sveltestrap/sveltestrap';

	import { symbols } from '../symbols.svelte.js';
	import * as helpers from '../helpers.js';
	import VersionSelect from '../../components/VersionSelect.svelte';
</script>

<div class="container" id="content">
	<VersionSelect page_name={'sections'}></VersionSelect>
</div>
