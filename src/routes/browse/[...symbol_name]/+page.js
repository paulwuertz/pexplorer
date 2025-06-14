export const trailingSlash = 'always';
export const prerender = false;
export const ssr = false;
export const csr = true;

export async function load({ url, params }) {
    // load elf data
    let componentData = {};

    const hasElfURLData = url.searchParams.has('elfURLData');
    const storedElfURLData = localStorage.getItem("lastOpenElfURL");
    const elfUrl = (hasElfURLData) ? decodeURIComponent(url.searchParams.get('elfURLData'))
                                    : storedElfURLData;
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
    const hasSelectedVersionToCompare = url.searchParams.has('symbols.selected_versions_to_compare');
    if (hasSelectedVersionToCompare)
    {
        const selected_version_to_compare_param = url.searchParams.get('symbols.selected_versions_to_compare')
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

    const currentPath = params.symbol_name;
    const currentPathDepth = params.symbol_name.split("/").length;
    console.log(currentPath);

    componentData.currentSubfolders = new Set(componentData.symbols[componentData.selected_version].symbols.map((sym) => {
        let trailingPath = sym.file.substr(currentPath.length);
        let pathArr = trailingPath.split("/")
        console.log(currentPath, sym.file, currentPathDepth, (pathArr.length - currentPathDepth), trailingPath, pathArr, pathArr.length)
        return (sym.file.startsWith(currentPath) && (pathArr.length - currentPathDepth) >= 2) ? (pathArr[0]) : "";
    }).filter(folder => folder!==""));
    componentData.currentSubfiles = componentData.symbols[componentData.selected_version].symbols.filter((sym) => {
        return sym.file.startsWith(currentPath)
    });

	return componentData;
};
