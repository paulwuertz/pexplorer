import { symbols } from './symbols.svelte.js';

export async function load({ fetch, url }) {
    // load elf data
    let componentData = {};
    const hasElfURLData = url.searchParams.has('elfURLData[]'); // TODO make array param
    const storedElfURLData = JSON.parse(localStorage.getItem("lastOpenElfURLs"));
    const elfUrls = (hasElfURLData) ? url.searchParams.getAll('elfURLData[]').map(e => decodeURIComponent(e))
                                    : storedElfURLData;

    if(elfUrls)
    {
        // download data
        localStorage.lastOpenElfURLs = JSON.stringify(elfUrls);
        localStorage.elfStorageDate = new Date().toISOString();
        symbols.symbolLinks = elfUrls;
        componentData.symbols = {};
        for (const elfUrl of elfUrls) {
            const response = await fetch(elfUrl);
            const data = await response.json();
            // persist
            console.log("Loaded elf data from", elfUrl);
            componentData.elfDataProvided = true;
            for (const versionKey of Object.keys(data)) {
                console.log("Loaded version data", versionKey);
                componentData.symbols[versionKey] = data[versionKey];
                symbols.symbols[versionKey] = data[versionKey];
            }
        }
        symbols.selected_version = null;
        symbols.selected_versions_to_compare = null;
        symbols.elfDataProvided = true;
    }
    else
    {
        symbols.elfDataProvided = false;
        console.log("No loaded elf data", elfUrls, Object.keys(symbols.symbols).length);
        return {};
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

    // localstorage has 5-10 MB max so split TODO test compression...
    // let versions = Object.keys(data);
    // localStorage.elfVersions = versions;
    // for (const version of versions) {
    //     localStorage[version] = JSON.stringify(data[version]);
    // }
    // else if (localStorage.getItem("elfURLData")) {
    //     alert("elfURLData found in storage"+localStorage.getItem("elfURLData"))
    //     elfDataProvided = false;
    // }

	return componentData;
};