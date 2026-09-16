<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	interface Product {
		id: string;
		name: string;
		description: string;
		price: number;
		image: string;
		category: string;
		subtitle: string;
	}

	let products = $state<Product[]>([]);
	let loading = $state(true);
	let activeFilter = $state('todo');
	let error = $state('');

	const API = 'http://localhost:8080/api';

	async function fetchProducts() {
		try {
			const res = await fetch(`${API}/products`, { credentials: 'include' });
			if (res.status === 401) { goto('/'); return; }
			if (res.ok) {
				products = await res.json();
			} else {
				error = 'Error cargando productos';
			}
		} catch {
			error = 'Error de conexión';
		} finally {
			loading = false;
		}
	}

	async function addToCart(productId: string, event: Event) {
		event.stopPropagation();
		event.preventDefault();
		const btn = event.currentTarget as HTMLButtonElement;
		
		btn.classList.add('bg-primary', 'text-on-primary');
		setTimeout(() => btn.classList.remove('bg-primary', 'text-on-primary'), 400);

		try {
			await fetch(`${API}/cart`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ product_id: productId, quantity: 1 }),
				credentials: 'include'
			});
			// In a real app we might update a global cart store or show a toast
			showToast();
		} catch (e) {
			console.error(e);
		}
	}

	let toastVisible = $state(false);
	let toastTimeout: ReturnType<typeof setTimeout>;
	function showToast() {
		clearTimeout(toastTimeout);
		toastVisible = true;
		toastTimeout = setTimeout(() => toastVisible = false, 2200);
	}

	onMount(() => {
		fetchProducts();
	});

	let filteredProducts = $derived(
		activeFilter === 'todo' 
			? products 
			: products.filter(p => p.category.toLowerCase().includes(activeFilter))
	);

	function setFilter(filter: string) {
		activeFilter = filter;
	}
</script>

<!-- Interactive Feedback Toast -->
<div class="fixed top-16 left-1/2 -translate-x-1/2 z-50 bg-inverse-surface text-inverse-on-surface px-space-md py-space-xs rounded-full shadow-xl flex items-center gap-space-xs pointer-events-none transition-all duration-300 transform {toastVisible ? 'opacity-100 translate-y-0' : 'opacity-0 -translate-y-2'}">
	<span class="material-symbols-outlined text-[18px] text-secondary-fixed">check_circle</span>
	<span class="font-label-sm text-label-sm">Agregado a tu bolsa con calma</span>
</div>

<!-- Header (Mobile only, Desktop uses layout nav) -->
<header class="md:hidden fixed top-0 inset-x-0 z-40 bg-surface/90 backdrop-blur-xl shadow-[0_1px_8px_rgba(0,0,0,0.03)] pt-safe">
	<div class="h-14 px-margin flex items-center justify-between max-w-md mx-auto">
		<div class="flex items-center gap-space-sm">
			<img alt="La Tiendita Logo" class="h-8 w-auto object-contain" src="https://lh3.googleusercontent.com/aida/AEtjO1VdBD_zfQ5Q4HTvgfcK5U0GnN2R44-whuA_1qYMzndpSPV9l9d0oVrZeWqSXd_HKMQOO-KwnLuhIXa4haTayD8W9M48rvvWefKRWmJ0z4IPVuEmFewqqXZYa3wdptuse-m5xeYEMABjtwKHx0ECJ4SVkSLVL54kBeUUcttNOuFZbbxjWFOWkSV5yu1tyPSlNVApOVRZla72Utj1BJ5PBlN0CLsgqR477UZ-23_3RV_bIrDLNNJkZdDMtd8"/>
			<span class="font-headline-sm text-headline-sm text-on-surface font-medium">La Tiendita</span>
		</div>
		<div class="flex items-center gap-space-xs">
			<button class="w-11 h-11 flex items-center justify-center rounded-full text-on-surface-variant hover:text-on-surface transition-colors">
				<span class="material-symbols-outlined text-[22px]">search</span>
			</button>
			<button class="w-11 h-11 flex items-center justify-center rounded-full" onclick={() => {
				fetch(`${API}/logout`, { method: 'POST', credentials: 'include' }).then(() => goto('/'));
			}}>
				<img alt="Profile" class="w-8 h-8 rounded-full object-cover" src="https://lh3.googleusercontent.com/aida-public/AB6AXuAcPHqh5WXHDRm1kdiadEXUJj_ummit40yNuJIScfiZfEy_4Vnw0kYgWLrsdpeRGfp8jbi8LydFPn-mMIKVPS-urlaggjvxtlcfx0gT9K-m1zjP_7uwj_YALXykDky4FvDiPByH3pqhIHg8q1bGQCBmoYuIxc7eQ92NEH1Yabt64S9jApLnEiSAZKKfFEQ4oD_aIu4B0LVMHyDqY1KvgsKqYxQvGpxEPa3pMXc8dQExapbK81fEEz_E"/>
			</button>
		</div>
	</div>
</header>

<main class="flex-1 flex flex-col relative w-full pt-14 md:pt-0 pb-20 md:pb-8 bg-surface max-w-7xl mx-auto">
	<div class="flex flex-col w-full pb-8">
		
		<!-- Search & Filter Bar -->
		<div class="px-margin pt-space-sm pb-space-xs max-w-xl mx-auto md:mx-0 w-full">
			<div class="bg-surface-container-low rounded-xl px-space-md py-2.5 flex items-center justify-between shadow-sm">
				<div class="flex items-center gap-space-sm flex-1 min-w-0">
					<span class="material-symbols-outlined text-[20px] text-on-surface-variant shrink-0">search</span>
					<input class="bg-transparent font-body-sm md:font-body-md text-body-sm md:text-body-md text-on-surface placeholder:text-on-surface-variant/70 focus:outline-none w-full" placeholder="Buscar objetos, ropa, café..." type="text"/>
				</div>
				<button aria-label="Filtrar" class="w-8 h-8 rounded-lg bg-surface-container flex items-center justify-center text-on-surface-variant hover:text-on-surface transition-colors shrink-0 ml-2">
					<span class="material-symbols-outlined text-[18px]">tune</span>
				</button>
			</div>
		</div>

		<!-- Editorial Warm Hero Banner -->
		<div class="px-margin pt-space-sm">
			<div class="relative overflow-hidden rounded-2xl md:rounded-3xl bg-surface-container p-space-lg md:p-12 shadow-sm min-h-[250px] md:min-h-[350px] flex flex-col justify-center">
				<div class="absolute inset-0 bg-cover bg-center md:bg-[center_top_-4rem] opacity-30 mix-blend-multiply pointer-events-none" style="background-image: url('https://lh3.googleusercontent.com/aida-public/AB6AXuAZNVmp6Zyu0QXymxzy1MvtwGqEgFHDl7MP9kLIHsjPbDHiawrTy1jEdsvcCPnJvFZkS4NOlaMQmLJvm_wnLcZOjVv7z3f7Aw1rqDDGG2SuTdWAUaSNuyw7H16qhW5zCYk-0A9-kTV3vLk2ew2CypP_eOmfbfOSsEcWN_fZ2Qgd8HzGy6BmRHxJlzTZ6F7cMTX45nZm5eryLm8lxQoV1LOwKtihTnsqIlDQHl7DrpuK2HPx1cLgiLcD')"></div>
				<div class="relative z-10 flex flex-col items-start max-w-[80%] md:max-w-md">
					<span class="bg-secondary-container text-on-secondary-container font-label-sm text-label-sm px-3 py-1 rounded-full mb-space-xs uppercase tracking-wider">Edición Limitada</span>
					<h2 class="font-headline-lg-mobile md:text-display text-headline-lg-mobile text-on-surface mb-space-xs leading-tight drop-shadow-sm">Colección Calma & Hogar</h2>
					<p class="font-body-sm md:text-body-lg text-body-sm text-on-surface-variant mb-space-md drop-shadow-sm">Piezas atemporales hechas a mano con amor y materiales nobles.</p>
					<button class="bg-primary hover:bg-primary-container active:scale-95 transition-all text-on-primary font-label-md md:text-label-lg text-label-md px-space-lg py-3 md:py-3 rounded-xl flex items-center gap-1.5 shadow-md">
						<span>Ver novedades</span>
						<span class="material-symbols-outlined text-[18px]">arrow_forward</span>
					</button>
				</div>
			</div>
		</div>

		<!-- Category Filter Horizontal Scroll -->
		<div class="pt-space-md pb-space-xs">
			<div class="flex items-center gap-2 md:gap-4 overflow-x-auto px-margin scrollbar-none" style="scrollbar-width: none;">
				<button class="shrink-0 px-4 md:px-6 py-2 md:py-2.5 rounded-full font-label-sm md:text-label-md text-label-sm transition-all shadow-sm flex items-center gap-1.5 {activeFilter === 'todo' ? 'bg-primary text-on-primary' : 'bg-surface-container text-on-surface-variant hover:bg-surface-container-high'}" onclick={() => setFilter('todo')}>
					{#if activeFilter === 'todo'}<span class="w-1.5 h-1.5 rounded-full bg-on-primary"></span>{/if}
					<span>Todo</span>
				</button>
				<button class="shrink-0 px-4 md:px-6 py-2 md:py-2.5 rounded-full font-label-sm md:text-label-md text-label-sm transition-all {activeFilter === 'gres' ? 'bg-primary text-on-primary shadow-sm' : 'bg-surface-container text-on-surface-variant hover:bg-surface-container-high'}" onclick={() => setFilter('gres')}>
					<span>Cerámica & Hogar</span>
				</button>
				<button class="shrink-0 px-4 md:px-6 py-2 md:py-2.5 rounded-full font-label-sm md:text-label-md text-label-sm transition-all {activeFilter === 'lino' ? 'bg-primary text-on-primary shadow-sm' : 'bg-surface-container text-on-surface-variant hover:bg-surface-container-high'}" onclick={() => setFilter('lino')}>
					<span>Ropa de Lino</span>
				</button>
				<button class="shrink-0 px-4 md:px-6 py-2 md:py-2.5 rounded-full font-label-sm md:text-label-md text-label-sm transition-all {activeFilter === 'infusión' ? 'bg-primary text-on-primary shadow-sm' : 'bg-surface-container text-on-surface-variant hover:bg-surface-container-high'}" onclick={() => setFilter('infusión')}>
					<span>Café & Infusiones</span>
				</button>
				<button class="shrink-0 px-4 md:px-6 py-2 md:py-2.5 rounded-full font-label-sm md:text-label-md text-label-sm transition-all {activeFilter === 'soja' ? 'bg-primary text-on-primary shadow-sm' : 'bg-surface-container text-on-surface-variant hover:bg-surface-container-high'}" onclick={() => setFilter('soja')}>
					<span>Aromas</span>
				</button>
			</div>
		</div>

		<!-- Section Header -->
		<div class="px-margin pt-space-lg pb-space-sm flex items-baseline justify-between">
			<div>
				<h3 class="font-headline-sm md:text-headline-lg text-headline-sm text-on-surface">Curaduría de la Semana</h3>
				<p class="font-body-sm md:text-body-md text-body-sm text-on-surface-variant md:mt-1">Selección serena para el día a día</p>
			</div>
			<span class="font-label-sm md:text-label-md text-label-sm text-primary font-medium tracking-wide bg-primary/10 px-3 py-1 rounded-full">{filteredProducts.length} piezas</span>
		</div>

		<!-- Product Grid (Responsive: 2 cols on mobile, 3 on md, 4 on lg) -->
		<div class="px-margin grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-space-sm md:gap-space-lg pt-space-xs">
			{#if loading}
				<div class="col-span-full text-center py-10 text-on-surface-variant">Cargando catálogo...</div>
			{:else}
				{#each filteredProducts as product}
				<a href="/product/{product.id}" class="product-card group flex flex-col bg-surface-container-low rounded-2xl md:rounded-3xl overflow-hidden p-space-xs transition-all duration-300 hover:shadow-lg hover:-translate-y-1 cursor-pointer block">
					<div class="relative w-full aspect-[4/5] rounded-xl md:rounded-2xl overflow-hidden bg-surface-container">
						<img alt={product.name} src={product.image} class="w-full h-full object-cover transition-transform duration-700 group-hover:scale-105" />
						<button aria-label="Favorito" class="absolute top-3 right-3 w-8 h-8 md:w-10 md:h-10 rounded-full bg-surface/80 backdrop-blur-md flex items-center justify-center text-on-surface-variant hover:text-primary hover:bg-surface transition-all shadow-sm" onclick={(e) => { e.preventDefault(); e.stopPropagation(); e.currentTarget.classList.toggle('text-primary'); }}>
							<span class="material-symbols-outlined text-[18px] md:text-[22px]">favorite</span>
						</button>
					</div>
					<div class="p-3 md:p-4 flex flex-col flex-1 justify-between gap-2">
						<div>
							<span class="font-label-sm text-label-sm text-secondary block mb-1">{product.category}</span>
							<h4 class="font-body-md md:text-headline-sm text-body-md text-on-surface font-medium line-clamp-1">{product.name}</h4>
						</div>
						<div class="flex items-center justify-between pt-2">
							<span class="font-label-md md:text-headline-sm text-label-md text-on-surface font-semibold">${product.price.toFixed(2)}</span>
							<button aria-label="Añadir a la bolsa" class="add-btn w-8 h-8 md:w-10 md:h-10 rounded-xl md:rounded-2xl bg-surface-container-high hover:bg-primary hover:text-on-primary text-on-surface flex items-center justify-center transition-colors active:scale-90 shadow-sm" onclick={(e) => addToCart(product.id, e)}>
								<span class="material-symbols-outlined text-[18px] md:text-[22px]">add</span>
							</button>
						</div>
					</div>
				</a>
				{/each}
			{/if}
		</div>

		<!-- Trust & Sustainability Micro-Banner -->
		<div class="px-margin pt-space-xl">
			<div class="bg-surface-container-low rounded-2xl md:rounded-3xl p-space-md md:p-space-lg flex flex-col md:flex-row md:items-center md:justify-center gap-space-md shadow-sm border border-outline-variant/20">
				<div class="w-12 h-12 md:w-16 md:h-16 rounded-full bg-secondary-container flex items-center justify-center text-secondary shrink-0 mx-auto md:mx-0">
					<span class="material-symbols-outlined text-[24px] md:text-[32px]">eco</span>
				</div>
				<div class="flex flex-col text-center md:text-left max-w-md">
					<span class="font-label-md md:text-headline-sm text-label-md text-on-surface font-medium mb-1">Filosofía Consciente</span>
					<span class="font-body-sm md:text-body-md text-body-sm text-on-surface-variant">Envíos sostenibles y empaques 100% compostables en cada orden para cuidar del hogar más grande que tenemos.</span>
				</div>
			</div>
		</div>

		<!-- Visual Delight / Wabi-sabi quote card -->
		<div class="px-margin pt-space-lg">
			<div class="text-center py-space-xl">
				<span class="material-symbols-outlined text-primary/40 text-[28px] mb-2">nest_eco_leaf</span>
				<p class="font-headline-sm md:text-headline-md text-headline-sm text-on-surface-variant/80 italic max-w-xs md:max-w-lg mx-auto">
					"La belleza reside en la silenciosa armonía de lo natural y lo simple."
				</p>
			</div>
		</div>
	</div>
</main>
