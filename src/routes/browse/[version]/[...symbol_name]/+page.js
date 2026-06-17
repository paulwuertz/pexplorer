import { goto } from '$app/navigation';
import { symbols } from '../../../symbols.svelte.js';
import { base } from '$app/paths';

export async function load({ route, params }) {
	// load elf data
	console.log('sym', symbols);

	if (!symbols.elfDataProvided) {
		goto(base + '/#/');
		return;
	}
	console.log('load return', route, params);
	return route, params;
}
