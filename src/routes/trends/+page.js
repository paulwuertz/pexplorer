import { goto } from '$app/navigation';
import { symbols } from '../symbols.svelte.js';
import { base } from '$app/paths';

export async function load({ url, parent }) {
	// load elf data
	let componentData = await parent();
	console.log('trends', symbols);

	if (!symbols.elfDataProvided) {
		goto(base + '/#/');
		return;
	}
	return componentData;
}
