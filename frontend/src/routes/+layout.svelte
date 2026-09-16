<script lang="ts">
	import { page } from '$app/stores';
	import '../app.css';
	let { children } = $props();
</script>

<!-- Global Desktop Top Navigation -->
{#if $page.url.pathname !== '/'}
<nav class="hidden md:flex fixed top-0 inset-x-0 h-16 bg-surface/90 backdrop-blur-xl shadow-sm z-50 items-center justify-between px-margin">
	<div class="flex items-center gap-space-sm cursor-pointer" onclick={() => window.location.href='/home'}>
		<img alt="La Tiendita Logo" class="h-8 w-auto object-contain" src="https://lh3.googleusercontent.com/aida/AEtjO1VdBD_zfQ5Q4HTvgfcK5U0GnN2R44-whuA_1qYMzndpSPV9l9d0oVrZeWqSXd_HKMQOO-KwnLuhIXa4haTayD8W9M48rvvWefKRWmJ0z4IPVuEmFewqqXZYa3wdptuse-m5xeYEMABjtwKHx0ECJ4SVkSLVL54kBeUUcttNOuFZbbxjWFOWkSV5yu1tyPSlNVApOVRZla72Utj1BJ5PBlN0CLsgqR477UZ-23_3RV_bIrDLNNJkZdDMtd8"/>
		<span class="font-headline-sm text-headline-sm font-medium text-on-surface">La Tiendita</span>
	</div>
	
	<div class="flex items-center gap-space-lg h-full">
		<a href="/home" class="flex items-center gap-2 h-full border-b-2 transition-all { $page.url.pathname === '/home' ? 'border-primary text-primary font-medium' : 'border-transparent text-on-surface-variant hover:text-on-surface' }">
			<span class="font-label-md text-label-md">Inicio</span>
		</a>
		<a href="/home" class="flex items-center gap-2 h-full border-b-2 transition-all { $page.url.pathname === '/explore' ? 'border-primary text-primary font-medium' : 'border-transparent text-on-surface-variant hover:text-on-surface' }">
			<span class="font-label-md text-label-md">Explorar</span>
		</a>
	</div>

	<div class="flex items-center gap-space-sm">
		<a href="/cart" class="w-10 h-10 flex items-center justify-center rounded-full hover:bg-surface-container transition-colors relative { $page.url.pathname === '/cart' ? 'text-primary' : 'text-on-surface-variant' }">
			<span class="material-symbols-outlined text-[24px]">shopping_bag</span>
		</a>
		<a href="javascript:void(0)" class="w-10 h-10 flex items-center justify-center rounded-full">
			<img alt="Profile" class="w-8 h-8 rounded-full object-cover" src="https://lh3.googleusercontent.com/aida-public/AB6AXuAcPHqh5WXHDRm1kdiadEXUJj_ummit40yNuJIScfiZfEy_4Vnw0kYgWLrsdpeRGfp8jbi8LydFPn-mMIKVPS-urlaggjvxtlcfx0gT9K-m1zjP_7uwj_YALXykDky4FvDiPByH3pqhIHg8q1bGQCBmoYuIxc7eQ92NEH1Yabt64S9jApLnEiSAZKKfFEQ4oD_aIu4B0LVMHyDqY1KvgsKqYxQvGpxEPa3pMXc8dQExapbK81fEEz_E"/>
		</a>
	</div>
</nav>
{/if}

<div class="md:pt-16 min-h-screen">
	{@render children()}
</div>

{#if $page.url.pathname !== '/'}
<!-- Bottom Navigation Bar (Mobile only) -->
<nav class="md:hidden fixed bottom-0 inset-x-0 z-50 pb-safe bg-surface/90 backdrop-blur-xl shadow-[0_-1px_8px_rgba(0,0,0,0.03)]">
	<div class="flex justify-around items-center h-16 px-space-xs max-w-md mx-auto">
		<a 
			class="flex flex-col items-center justify-center flex-1 h-full transition-all { $page.url.pathname === '/home' ? 'text-primary font-semibold' : 'text-on-surface-variant hover:text-on-surface' }" 
			href="/home"
		>
			<span class="material-symbols-outlined text-[22px]">home</span>
			<span class="font-label-sm text-label-sm mt-1">Inicio</span>
		</a>
		<a 
			class="flex flex-col items-center justify-center flex-1 h-full transition-all { $page.url.pathname === '/explore' ? 'text-primary font-semibold' : 'text-on-surface-variant hover:text-on-surface' }" 
			href="/home"
		>
			<span class="material-symbols-outlined text-[22px]">grid_view</span>
			<span class="font-label-sm text-label-sm mt-1">Explorar</span>
		</a>
		<a 
			class="flex flex-col items-center justify-center flex-1 h-full transition-all relative { $page.url.pathname === '/cart' ? 'text-primary font-semibold' : 'text-on-surface-variant hover:text-on-surface' }" 
			href="/cart"
		>
			<div class="relative flex items-center justify-center">
				<span class="material-symbols-outlined text-[22px]">shopping_bag</span>
			</div>
			<span class="font-label-sm text-label-sm mt-1">Bolsa</span>
		</a>
		<a 
			class="flex flex-col items-center justify-center flex-1 h-full transition-all { $page.url.pathname === '/profile' ? 'text-primary font-semibold' : 'text-on-surface-variant hover:text-on-surface' }" 
			href="javascript:void(0)"
		>
			<span class="material-symbols-outlined text-[22px]">person</span>
			<span class="font-label-sm text-label-sm mt-1">Perfil</span>
		</a>
	</div>
</nav>
{/if}
