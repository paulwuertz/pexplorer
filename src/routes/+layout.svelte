<script>
  //import '../app.css';
  
  let { children } = $props();
  
  import { base } from '$app/paths';
  import { browser } from '$app/environment';
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  
  import {
      Collapse,
      Nav,
      Navbar,
      NavbarBrand,
      NavItem,
      NavLink,
      NavbarToggler
  } from '@sveltestrap/sveltestrap';
  
  import { symbols } from './symbols.svelte.js';
  
  let isOpen = $state(true);
  
  function handleUpdate(event) {
      isOpen = event.detail.isOpen;
  }
</script>

<nav>
  <Navbar color="light" light expand="md" container="md">
    <NavbarBrand href="{base}/">Puncover</NavbarBrand>
    <NavbarToggler on:click={() => (isOpen = !isOpen)} />
    <Collapse {isOpen} navbar expand="md" on:update={handleUpdate}>
        <Nav class="ms-auto" navbar>
            <NavItem><NavLink href="{base}/">Browse (TODO)</NavLink></NavItem>
            <NavItem><NavLink href="{base}/diff">Compare Versions</NavLink></NavItem>
            <NavItem><NavLink href="{base}/">Interrupts (TODO)</NavLink></NavItem>
            <NavItem><NavLink href="{base}/symbols">Symbols</NavLink></NavItem>
            <NavItem><NavLink href="{base}/">Thread (TODO)</NavLink></NavItem>
            <NavItem><NavLink href="{base}/">Settings (TODO)</NavLink></NavItem>
            <NavItem><NavLink href="{base}/">About (TODO)</NavLink></NavItem>
        </Nav>
    </Collapse>
  </Navbar>
  <!--- TODO add breadcrumbs-->
</nav>

{@render children()}

<hr>

<div id="page-footer">
    brought to you with &#128150;
    by <a href="https://twitter.com/hbehrens">Heiko Behrens</a>, Paul Würtz
    <span class="secondary">– MIT license, copyright &copy; 2014-2025</span>
</div>
