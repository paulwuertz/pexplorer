# pexplorer

A firmware explorer inspired by [puncover](https://github.com/HBehrens/puncover/), trying to bring the most information of your build, by statically analysing your ELF file.
Try it yourself [locally in your browser](https://paulwuertz.github.io/pexplorer), upload you ELF (build with debug info "-g") and explore your firmware memory footprint or compare firmware versions.

## Functions

### "MAP file viewer" like browsing

Upload ELF files and get a MAP file viewer like experience browsing all function and variable symbols of your firmware. Sort and filter to focus on certain aspects of your firmware. File-system like exploration of consumed memory by sunburst diagrams.

![upload and explore](/static/img/upload+explore.png)

### Diff-ing your firmware builds and static analysis

Uploading a second firmware version you can see the impact changes have on your memory, by high-lighting all changes in flash, RAM and stack size of any symbols.

The ASM binary of each function is dissambled with [capstone-js](https://alexaltea.github.io/capstone.js/). Static function calls are detected from their call instructions. Stack size is estimated, by the debugging information, needed to unwind the stack on breakpoints. Dynamic function calls are detected (and planned to be resolvable by hand). Currently only an experiment only working on Cortex-M controller firmware.

![diff and function detail page](/static/img/diff+fndetails.png)

## Developing

To try it on your ELF, use the [github page export](https://paulwuertz.github.io/pexplorer).

For local deployment install dependencies with `npm install` (or `pnpm install` or `yarn`) and start a development server: `npm run dev`

## Building

To create a production version we need to build the WASM go module first. TODO better describtion...

Then the svelte web app can be build: `npm run build`

You can preview the production build with `npm run preview`.

## What features are planned?

Have a look at the [project board](https://github.com/users/paulwuertz/projects/1).

Some of the big features I want to realize down the line:

* static analysis for...
    * ...stack sizes that is tested and exact
    * ...callgraph extraction and manual dynamic call resolution
    * ...worst case stack size for each function
    * ...worst case exection time - at least some instruction counting estimate
    * ...Interrupts - runtime and stack sizes for these are especially important
* RTOS feature detection
    * automatically detect the most popular RTOSes zephyr, FreeRTOS, ...
    * extract a list of all threads, stacks,
    * get stack analysis of threads with actual configured stack sizes, where possible
* better diffs and trend detection
    * diff of dis-ASM for functions across versions
    * compare more then two versions and get more stat plot out of them
* configuration
    * support static - sharable links for every page
    * add settings and resolution pages
* build a CLI command utilities, to be used in CI automation for stack health checks
