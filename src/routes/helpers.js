let Capstone = null;

export async function loadCapstone(cs, reject) {
	Capstone = cs;
}

export let symbolsToMap = (syms) => {
	let symMap = {};
	for (const sym of syms) {
		sym.remark = sym.called_from_other_file ? 'linked-from-library' : '';
		sym.newSymbols = false;
		sym.deletedSymbols = false;
		symMap[sym.name] = sym;
	}
	return symMap;
};

export let symbolsToFunctionMap = (symMap) => {
	return Object.values(symMap).filter((e) => {
		return e['symtype'] === 'fn';
	});
};

export let symbolsToVariableMap = (symMap) => {
	return Object.values(symMap).filter((e) => {
		return e['symtype'] === 'var';
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

export const sympath_to_link = (base, symbol_version, callxrs_text) => {
	return base + '/#/browse/' + symbol_version + '/' + callxrs_text;
};

export const try_get_callee_link_by_addr = (
	base,
	symbol_version,
	call_to_addr,
	sym_path_by_addr
) => {
	let callee_name = sym_path_by_addr[call_to_addr];
	let link_href = base + '/#/browse/' + symbol_version + '/' + callee_name;
	if (callee_name == undefined) {
		return ['/', '???'];
	}
	let callxrs_slugs = callee_name.split('/');
	let sym_name = callxrs_slugs[callxrs_slugs.length - 1];
	return [link_href, sym_name];
};

export const callxrs_text_to_links = (
	base,
	symbol_version,
	callxrs,
	sym_path_by_addr,
	isCaller
) => {
	let direction = isCaller ? 'from' : 'to';
	let callxrs_addr = callxrs[direction];
	let callxrs_text = sym_path_by_addr[callxrs_addr];
	return base + '/#/browse/' + symbol_version + '/' + callxrs_text;
};

export const callxrs_text_to_symname = (callxrs, sym_path_by_addr, isCaller) => {
	let direction = isCaller ? 'from' : 'to';
	if (direction == 'to' && !Object.hasOwn(callxrs, 'to') && callxrs['dynamic']) {
		return 'Unresolved dynamic call, ';
	}
	let callxrs_addr = callxrs[direction];
	let callxrs_text = sym_path_by_addr[callxrs_addr];
	if (callxrs_text == undefined) {
		return '???';
	}
	//console.log('callxrs_text', callxrs_addr, callxrs_text, callxrs_text);
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
		const noPath = ['<unknown>'];
		let path_elements;
		if (!symbol.file || typeof symbol.file !== 'string') {
			path_elements = noPath;
		} else {
			path_elements = symbol.file.split('/').slice(1);
		}
		add_sym_to_object_tree(symbol, path_elements, data, data_field);
	}
	// unpack to first level with > 1 element
	while (data.length == 1 && data[0].children) {
		data = data[0].children;
	}
	return data;
};

export const symbols_to_call_tree_data = (entry_fn) => {
	let data = {
		name: entry_fn.name + '(' + entry_fn.stack_size + ')',
		value: entry_fn.stack_size,
		children: []
	};
	let callees = Object.hasOwn(entry_fn, 'calls') ? entry_fn['calls'] : [];
	for (const callee of callees) {
		data.children.push(symbols_to_call_tree_data(callee));
	}
	return data;
};

export const row2AHref = (base, selected_version, row_data) => {
	if (row_data.file) {
		return base + '/#/browse/' + selected_version + row_data.file + '/' + row_data.name;
	} else {
		return base + '/#/browse/' + selected_version + '/0x' + row_data.address.toString(16);
	}
};

export const isCallInstr = (instr) => {
	if (instr === 'bl') {
		// && ARM
		return true;
	}
	return false;
};

export const getCallAddr = (instr) => {
	if (true) {
		// ARM
		// const like '#0xae39'
		return parseInt(instr.op_str.substr(1), 16);
	}
	return 0;
};

export const checkCallInstrLink = (instr, base, symbol_version, sym_path_by_addr) => {
	let link = '';
	let isCall = isCallInstr(instr.mnemonic);
	if (isCall) {
		let addr = getCallAddr(instr);
		let [link_href, callee_name] = try_get_callee_link_by_addr(
			base,
			symbol_version,
			addr,
			sym_path_by_addr
		);
		link += ' => call to <a href="' + link_href + '">' + callee_name + '</a>';
	}
	return link;
};

export const download = (file, text) => {
	//creating an invisible element
	let element = document.createElement('a');
	element.setAttribute('href', 'data:text/plain;charset=utf-8, ' + encodeURIComponent(text));
	element.setAttribute('download', file);
	document.body.appendChild(element);
	element.click();

	document.body.removeChild(element);
};

export const csBase64ToASMText = (
	base64text,
	baseAddr,
	show_full_asm,
	base,
	symbol_version,
	sym_path_by_addr
) => {
	var d = new Capstone.Capstone(Capstone.ARCH_ARM, Capstone.MODE_THUMB + Capstone.MODE_MCLASS);
	let ASM = Uint8Array.fromBase64(base64text);
	// console.log('ASM: ' + ASM, base64text);
	let disasmData = d.disasm(ASM, baseAddr);
	// console.log(JSON.stringify(disasmData, null, 4));
	// Display results;
	let result = '\tAddr\tINSTR bytes     mnemonic\tOP\n';
	if (!show_full_asm) disasmData = disasmData.slice(0, 10);
	disasmData.forEach(function (instr) {
		let linkToCall = checkCallInstrLink(instr, base, symbol_version, sym_path_by_addr);
		result +=
			'\t0x' +
			instr.address.toString(16) +
			':\t' +
			instr.bytes
				.map((e) => e.toString(16))
				.join('')
				.padEnd(15, ' ') +
			'\t' +
			instr.mnemonic +
			'\t\t' +
			instr.op_str +
			linkToCall +
			'\n';
	});
	if (!show_full_asm) result += '...';

	// Delete decoder
	// d.close();
	return result;
};

export const getDisasmFnMap = (asmReport) => {
	let fn2Disasm = {};
	for (let i = 0; i < asmReport['functions'].length; i++) {
		const f = asmReport['functions'][i];
		if (!Object.hasOwn(f, 'asm') || !Object.hasOwn(f, 'name')) {
			console.log(f, 'has no asm for getting its calltree...');
			continue;
		}
		let fName = f['name'];
		let baseAddr = f['address'];
		let fFile = f['file'] || '';
		let ASM = Uint8Array.fromBase64(f['asm']);
		let disasm = [];
		let disasmData;
		try {
			var d = new Capstone.Capstone(Capstone.ARCH_ARM, Capstone.MODE_THUMB + Capstone.MODE_MCLASS);
			disasmData = d.disasm(ASM, baseAddr);
		} catch (error) {
			console.log(fName + fFile + ' ASM: ' + ASM);
			console.log(f, disasmData);
			continue;
		}
		// Display results;
		let fnInstr = [];
		disasmData.forEach(function (instr) {
			fnInstr.push({
				addr: instr.address,
				instruction: instr.mnemonic,
				opstr: instr.op_str,
				insBytes: instr.bytes
			});
		});
		fn2Disasm[baseAddr] = fnInstr;
	}
	// console.log(fn2Disasm, JSON.stringify(fn2Disasm, null, 4));
	return fn2Disasm;
};

export const fn2symPathLookups = (reportFns) => {
	let symPathByAddr = {};
	let symPathByName = {};
	for (let i = 0; i < reportFns.length; i++) {
		let addr = reportFns[i]['address'];
		// TODO what about syms with unknown path - can they be eliminated ^^?
		let urlPath = reportFns[i]['file'] + '/' + reportFns[i]['name'];
		reportFns[i]['symtype'] = 'fn';
		symPathByAddr[addr] = urlPath;
		symPathByName[urlPath] = reportFns[i];
	}
	return [symPathByAddr, symPathByName];
};

export const stored_thread_settings_stack_size = (thread, variables) => {
	if (Object.hasOwn(thread, 'stack_variable_name')) {
		let stack_name = thread.stack_variable_name;
		let stack_var = variables.find((v) => v.name == stack_name);
		if (stack_var) {
			return stack_var.size;
		} else {
			alert(stack_name + ' var could not be found for thread' + thread.name);
			return 0;
		}
	} else if (Object.hasOwn(thread, 'size')) {
		return thread.size;
	} else {
		alert(thread.name + ' has no associated stack');
		return 0; // TODO error
	}
};

let store_default_settings = (firmware_hash) => {
	localStorage.setItem(version.firmware_hash, '{}');
	return '{}';
};

export const local_storage_key = 'pexplorer_settings';
export const restore_default_settings = (firmware_hash) => {
	let stored_settings_str = localStorage.getItem(local_storage_key);
	let no_settings_backed_up = !stored_settings_str;
	if (no_settings_backed_up) {
		stored_settings_str = store_default_settings(firmware_hash);
	}
	console.log('stored_settings_str', stored_settings_str);
	let stored_settings = JSON.parse(stored_settings_str);
	return stored_settings;
};

export const flat_calls_to_arrayed_callees = (stored_settings) => {
	let calls = {};
	// build a map
	for (const dynamic_call of stored_settings['dynamic_calls']) {
		let call_from = dynamic_call['call_from'];
		let call_to = dynamic_call['call_to'];
		if (Object.hasOwn(calls, call_from)) {
			calls[call_from].push(call_to);
		} else {
			calls[call_from] = [call_to];
		}
	}
	// map to array
	let calls_arr = [];
	for (const dynamic_caller in calls) {
		calls_arr.push({
			caller: dynamic_caller,
			callees: calls[dynamic_caller]
		});
	}
	return calls_arr;
};
