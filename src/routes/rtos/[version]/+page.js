import { goto } from '$app/navigation';
import { symbols } from '../../symbols.svelte.js';
import { base } from '$app/paths';

export async function load({ params }) {
	if (!symbols.elfDataProvided) {
		goto(base + '/#/');
		console.log('rtos nodata version', params);
		return params;
	}
	console.log('rtos sidata version', params);
	return params;
}
