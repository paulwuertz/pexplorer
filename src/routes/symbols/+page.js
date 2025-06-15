import { goto } from '$app/navigation';
import { symbols } from '../symbols.svelte.js';

export const trailingSlash = 'always';
export const prerender = false;
export const ssr = false;
export const csr = true;

export async function load({ url, parent }) {
    // load elf data
    let componentData = await parent();
    console.log("sym",symbols)

    if (!symbols.elfDataProvided){
        goto("/");
        return;
    }
	return componentData;
};
