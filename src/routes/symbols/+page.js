import { goto } from '$app/navigation';

export const trailingSlash = 'always';
export const prerender = false;
export const ssr = false;
export const csr = true;

export async function load({ url }) {
    // load elf data
    let componentData = {};

    const hasElfURLData = url.searchParams.has('elfURLData');
    const storedElfURLData = localStorage.getItem("lastOpenElfURLs");
    const elfUrl = (hasElfURLData) ? [decodeURIComponent(url.searchParams.get('elfURLData'))]
                                    : storedElfURLData;

    try {
        test_if_a_valid_url = new URL(elfUrl);
        console.log("trurl")
    } catch (_) {
        console.log("furl")
        goto("/")
        return;
    }

    if(elfUrl)
    {
        // download data
        const response = await fetch(elfUrl);
        const data = await response.json();
        // persist
        localStorage.lastOpenElfURL = elfUrl;
        localStorage.elfStorageDate = new Date().toISOString();
        componentData.symbols = data;
        console.log("Loaded elf data");
        componentData.elfDataProvided = true;
    }
    else
    {
        console.log("Loaded elf data", elfUrl, Object.keys(symbols.symbols).length);
    }
    // version of the elf
    const hasSelectedVersion = url.searchParams.has('selected_version');
    if(hasSelectedVersion)
    {
        const selected_version_param = url.searchParams.get('selected_version')
        localStorage.selected_version = selected_version_param;
        componentData.selected_version = selected_version_param;
    }
    else if (localStorage.getItem("selected_version"))
    {
        componentData.selected_version = localStorage.getItem("selected_version");
    }
    // version of the elf to compare to
    const hasSelectedVersionToCompare = url.searchParams.has('selected_versions_to_compare');
    if (hasSelectedVersionToCompare)
    {
        const selected_version_to_compare_param = url.searchParams.get('selected_versions_to_compare')
        localStorage.selected_versions_to_compare = selected_version_to_compare_param;
        componentData.selected_versions_to_compare = selected_version_to_compare_param;
    }
    else if (localStorage.getItem("selected_versions_to_compare"))
    {
        componentData.selected_versions_to_compare = localStorage.getItem("selected_versions_to_compare");
    }

	return componentData;
};
