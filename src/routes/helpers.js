export let symbolsToMap = (syms) => {
	let symMap = {};
	for (const sym of syms) {
		sym.remark = sym.called_from_other_file ? 'x-module' : '';
		sym.newSymbols = false;
		sym.deletedSymbols = false;
		symMap[sym.file + sym.display_name] = sym;
	}
	return symMap;
};

export let symbolsToFunctionMap = (symMap) => {
	return Object.values(symMap).filter((e) => {
		return e['type'] === 'function';
	});
};

export let symbolsToVariableMap = (symMap) => {
	return Object.values(symMap).filter((e) => {
		return e['type'] === 'variable';
	});
};

export let symMapToSymNameSet = (symMap) => {
	return new Set(Object.keys(symMap));
};

export let get_all_threads_names = (allSymVersions) => {
	let threads = new Set();
	for (let symVersion of Object.values(allSymVersions)) {
		for (let thread_name of Object.keys(symVersion['stack_reports'])) {
			threads.add(thread_name);
		}
	}
	return threads;
};

export let get_all_threads_function_names_on_stacks = (allSymVersions) => {
	let functions = new Set();
	for (let symVersion of Object.values(allSymVersions)) {
		for (let thread_obj of Object.values(symVersion['stack_reports'])) {
			if (typeof thread_obj != 'object') continue;
			for (let function_obj of thread_obj['call_stack']) {
				functions.add(function_obj['function']);
			}
		}
	}
	return functions;
};

/**
 * Gets all versions ordered by timestamp
 * @param allSymVersions
 */
export let get_versions_ordered_by_timestamps = (allSymVersions) => {
	let versions = [];
	// get all versions
	for (const [versionStr, symVersion] of Object.entries(allSymVersions)) {
		const timestamp = new Date(symVersion['timestamp']);
		versions.push({ version: versionStr, timestamp: timestamp });
	}
	versions.sort(function (a, b) {
		return a['timestamp'] > b['timestamp'];
	});
	return versions;
};

export let get_max_stack_sizes_of_thread = (allSymVersions, threadname) => {
	let threads = new Set();
	for (let symVersion of Object.values(allSymVersions)) {
		for (let thread_name of Object.keys(symVersion['stack_reports'])) {
			threads.add(thread_name);
		}
	}
	// alert(threads)
	return threads;
};

export const callxrs_text_to_links = (base, symbol_version, callxrs_text) => {
    return base+"/browse/" + symbol_version + "/" +  callxrs_text
};

export const callxrs_text_to_symname = (callxrs_text) => {
    let callxrs_slugs = callxrs_text.split("/")
    let sym_name = callxrs_slugs[ callxrs_slugs.length - 1 ]
    return sym_name
};
