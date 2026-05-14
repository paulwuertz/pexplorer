import { goto } from '$app/navigation';
import { symbols } from '../../symbols.svelte.js';
import { base } from '$app/paths';

export async function load({ url, params, parent }) {
	let componentData = await parent();

	let versions = Object.keys(symbols.symbols);
	console.log('browse', symbols, params.version, versions.length);
	if (versions.length == 1 && !params.version) {
		let selVers = versions[0];
		symbols.elfDataProvided = true;
		symbols.selected_version = selVers;
		goto(base + '/#/rtos/' + selVers);
	} else if (!symbols.elfDataProvided) {
		goto(base + '/#/');
		return;
	}
	return (params, componentData);
}
