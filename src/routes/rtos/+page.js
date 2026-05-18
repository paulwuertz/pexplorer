import { goto } from '$app/navigation';
import { symbols } from '../symbols.svelte.js';
import { base } from '$app/paths';

export async function load({ route, params }) {
	let versions = Object.keys(symbols.symbols);
    let version = params && Object.hasOwn(params, "version") && params.version
	console.log('rtos versel', symbols, params, version, versions.length);
	if (versions.length == 1 && !version) {
		let selVers = versions[0];
		symbols.elfDataProvided = true;
		symbols.selected_version = selVers;
		goto(base + '/#/rtos/' + selVers);
		return params;
	} else if (!symbols.elfDataProvided) {
		goto(base + '/#/');
		return params;
	}
	console.log('rtos versel', route, params);
	return params;
}
