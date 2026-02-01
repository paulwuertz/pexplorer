import * as cs from "@alexaltea/capstone-js/dist/capstone.min.js"

export let symbolsToMap = (syms) => {
	let symMap = {};
	for (const sym of syms) {
		sym.remark = sym.called_from_other_file ? 'linked-from-library' : '';
		sym.newSymbols = false;
		sym.deletedSymbols = false;
		symMap[sym.file + sym.name] = sym;
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
	return base + '/browse/' + symbol_version + '/' + callxrs_text;
};

export const callxrs_text_to_symname = (callxrs_text) => {
	let callxrs_slugs = callxrs_text.split('/');
	let sym_name = callxrs_slugs[callxrs_slugs.length - 1];
	return sym_name;
};

export const symbols_total_size = (symbols) => {
	let summed_size = 0;
	symbols.forEach((element) => {
		summed_size += element.size;
	});
	return summed_size;
};

let add_sym_to_object_tree = (symbol, path_array, data_tree, data_field) => {
	let entry = {
		name: symbol.name,
		value: symbol[data_field]
	};
	let folder_name = path_array.shift();
	let current_branch = data_tree;
	while (folder_name) {
		let sub_branch = current_branch.find((e) => e.name == folder_name);
		if (!sub_branch) {
			let new_branch = { name: folder_name, children: [] };
			current_branch.push(new_branch);
			current_branch = new_branch.children;
			// console.log("added", folder_name);
		} else {
			current_branch = sub_branch.children;
			// console.log("enter+found", folder_name);
		}
		// next sub-branch
		folder_name = path_array.shift();
	}
	current_branch.push(entry);
};

export const symbols_to_sunburst_tree_data = (symbols, data_field) => {
	let data = [];
	for (const symbol of symbols) {
		if (typeof symbol.file !== 'string') {
			continue;
		}
		let path_elements = symbol.file ? symbol.file.split('/').slice(1) : ["unlocatable"];
		add_sym_to_object_tree(symbol, path_elements, data, data_field);
	}
	return data;
};

export const row2AHref = (base, selected_version, row_data) => {
	return base + '/browse/' + selected_version + '/' + row_data.file + '/' + row_data.name;
};

var d = new cs.Capstone(cs.ARCH_ARM, cs.MODE_THUMB+cs.MODE_MCLASS);

export const csBase64ToASMText = (base64text, baseAddr) =>  {
    let ASM = Uint8Array.fromBase64(base64text)
    console.log("ASM: "+ASM, base64text);
    let disasmData = d.disasm(ASM, baseAddr)
    console.log(JSON.stringify(disasmData, null, 4));
    // Display results;
    let result = "\tAddr\tINSTR\tOP\n";
    disasmData.forEach(function (instr) {
        result += "\t0x" + instr.address.toString(16) + ":\t" + instr.mnemonic + "\t" + instr.op_str + "\n"
    });
    console.log(result, null, 4);

    // Delete decoder
    // d.close();
    return result
}

export const getDisasmFnMap = (asmReport) => {
	let fn2Disasm = {}
	for (let i= 0; i < asmReport["functions"].length; i++) {
		const f = asmReport["functions"][i];
		if (!Object.hasOwn(f, "asm") || !Object.hasOwn(f, "name")) {
			console.log(f, "has no asm for getting its calltree...");
			continue;
		}
		let fName = f["name"]
		let baseAddr = f["address"]
		let fFile = f["file"] || ""
		let ASM = Uint8Array.fromBase64(f["asm"])
		let disasm = [];
		console.log(fName+fFile+" ASM: "+ASM);
		let disasmData = d.disasm(ASM, baseAddr)
		// Display results;
		disasmData.forEach(function (instr) {
			fn2Disasm[baseAddr] = {
				"addr": instr.address, 
				"instruction": instr.mnemonic, 
				"opstr": instr.op_str,
				"insBytes": instr.bytes
			}
		});
	}
	console.log(fn2Disasm, JSON.stringify(fn2Disasm, null, 4));
	return fn2Disasm
}