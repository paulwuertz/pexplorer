import { goto } from '$app/navigation';
import { symbols } from '../symbols.svelte.js';
import { base } from '$app/paths';

export const trailingSlash = 'always';
export const prerender = false;
export const ssr = false;
export const csr = true;

export async function load({ url, parent }) {
	// load elf data
	let componentData = await parent();
	
	let versions = Object.keys(symbols.symbols)
	console.log('browse', symbols, versions, versions.length);
	if (versions.length == 1) {
		let selVers = versions[0]
		symbols.elfDataProvided = true
		symbols.selected_version = selVers
		goto(base + '/browse/' + selVers + "/");
	}else if (!symbols.elfDataProvided) {
		goto(base + '/');
		return;
	}
	return componentData;
}
