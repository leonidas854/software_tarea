<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	const productId = $page.params.id;
	let product: any = $state(null);
	let loading = $state(true);
	let error = $state('');

	const API = 'http://localhost:8080/api';

	let quantity = $state(1);
	let adding = $state(false);

	let toastVisible = $state(false);
	let toastTimeout: ReturnType<typeof setTimeout>;

	async function fetchProduct() {
		try {
			const res = await fetch(`${API}/products`, { credentials: 'include' });
			if (res.status === 401) { goto('/'); return; }
			if (res.ok) {
				const products = await res.json();
				product = products.find((p: any) => p.id === productId);
				if (!product) error = 'Producto no encontrado';
			} else {
				error = 'Error cargando el producto';
			}
		} catch {
			error = 'Error de conexión';
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		fetchProduct();
	});

	async function addToCart() {
		if (adding) return;
		adding = true;

		try {
			const res = await fetch(`${API}/cart`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ product_id: product.id, quantity }),
				credentials: 'include'
			});
			if (res.ok) {
				showToast(`${quantity} pieza(s) agregada(s) a tu bolsa`);
			}
		} catch (e) {
			console.error(e);
		} finally {
			adding = false;
		}
	}

	function showToast(msg: string) {
		clearTimeout(toastTimeout);
		toastVisible = true;
		toastTimeout = setTimeout(() => toastVisible = false, 2400);
	}

	let isFav = $state(false);

	let descExpanded = $state(false);
	let careExpanded = $state(false);

	let selectedColor = $state('Arena Natural');
</script>

{#if loading}
	<div class="min-h-screen flex items-center justify-center bg-surface">Cargando...</div>
{:else if error || !product}
	<div class="min-h-screen flex items-center justify-center bg-surface">
		<div class="text-center">
			<p class="text-error">{error}</p>
			<button class="mt-4 text-primary" onclick={() => history.back()}>Volver</button>
		</div>
	</div>
{:else}
	<!-- Desktop Header uses global layout nav, hide this mobile header on desktop -->
	<header class="md:hidden fixed top-0 inset-x-0 z-50 bg-surface/90 backdrop-blur-xl shadow-[0_1px_8px_rgba(0,0,0,0.03)] pt-safe">
		<div class="h-14 px-margin flex items-center justify-between max-w-md mx-auto">
			<div class="flex items-center gap-space-xs">
				<button class="w-11 h-11 flex items-center justify-center text-on-surface hover:text-primary transition-colors" onclick={() => history.back()}>
					<span class="material-symbols-outlined text-[22px]">arrow_back</span>
				</button>
				<div class="flex items-center gap-space-xs">
					<img alt="La Tiendita Logo" class="h-7 w-auto object-contain" src="https://lh3.googleusercontent.com/aida/AEtjO1VdBD_zfQ5Q4HTvgfcK5U0GnN2R44-whuA_1qYMzndpSPV9l9d0oVrZeWqSXd_HKMQOO-KwnLuhIXa4haTayD8W9M48rvvWefKRWmJ0z4IPVuEmFewqqXZYa3wdptuse-m5xeYEMABjtwKHx0ECJ4SVkSLVL54kBeUUcttNOuFZbbxjWFOWkSV5yu1tyPSlNVApOVRZla72Utj1BJ5PBlN0CLsgqR477UZ-23_3RV_bIrDLNNJkZdDMtd8"/>
					<h1 class="font-headline-sm text-headline-sm font-medium text-on-surface truncate">Detalle</h1>
				</div>
			</div>
			<div class="flex items-center gap-space-xs">
				<a class="w-11 h-11 flex items-center justify-center rounded-full" href="javascript:void(0)">
					<img alt="Profile" class="w-8 h-8 rounded-full object-cover" src="https://lh3.googleusercontent.com/aida-public/AB6AXuAcPHqh5WXHDRm1kdiadEXUJj_ummit40yNuJIScfiZfEy_4Vnw0kYgWLrsdpeRGfp8jbi8LydFPn-mMIKVPS-urlaggjvxtlcfx0gT9K-m1zjP_7uwj_YALXykDky4FvDiPByH3pqhIHg8q1bGQCBmoYuIxc7eQ92NEH1Yabt64S9jApLnEiSAZKKfFEQ4oD_aIu4B0LVMHyDqY1KvgsKqYxQvGpxEPa3pMXc8dQExapbK81fEEz_E"/>
				</a>
			</div>
		</div>
	</header>

	<main class="flex-1 flex flex-col relative w-full pt-14 md:pt-8 pb-32 bg-surface max-w-7xl mx-auto px-margin">
		<!-- Top In-Page Action Nav -->
		<div class="hidden md:flex items-center justify-between py-space-sm mb-4">
			<button class="inline-flex items-center gap-space-xs text-on-surface-variant hover:text-primary transition-colors py-space-xs" onclick={() => history.back()}>
				<span class="material-symbols-outlined text-[18px]">arrow_back</span>
				<span class="font-label-md text-label-md">Volver a tienda</span>
			</button>
		</div>

		<!-- 2-Column Grid on Desktop -->
		<div class="flex flex-col md:flex-row gap-space-lg md:gap-12 w-full">
			
			<!-- Left: Product Image Hero Gallery -->
			<div class="w-full md:w-1/2 md:sticky md:top-24 h-max">
				<div class="relative w-full aspect-[4/5] md:aspect-square rounded-2xl md:rounded-3xl overflow-hidden bg-surface-container-low shadow-sm">
					<div class="flex w-full h-full overflow-x-auto snap-x snap-mandatory scroll-smooth no-scrollbar" style="scrollbar-width: none; -ms-overflow-style: none;">
						<div class="w-full h-full flex-shrink-0 snap-center relative">
							<img class="w-full h-full object-cover hover:scale-105 transition-transform duration-700" alt={product.name} src={product.image}/>
						</div>
					</div>
					<div class="absolute top-4 left-4 bg-secondary-container/90 backdrop-blur-md px-space-sm py-1.5 rounded-full flex items-center gap-1 shadow-sm">
						<span class="material-symbols-outlined text-[14px] text-secondary">eco</span>
						<span class="font-label-sm text-label-sm text-on-secondary-container tracking-wider uppercase">Pieza Única</span>
					</div>
					<!-- Mobile back/favorite buttons overlaid (optional, since header exists) -->
					<div class="absolute top-4 right-4 flex md:hidden">
						<button aria-label="Guardar en favoritos" class="w-10 h-10 rounded-full bg-surface/80 backdrop-blur-md flex items-center justify-center text-on-surface transition-transform active:scale-90 shadow-sm" onclick={() => isFav = !isFav}>
							<span class="material-symbols-outlined text-[20px] transition-colors {isFav ? 'text-primary' : ''}" style={isFav ? "font-variation-settings: 'FILL' 1;" : ""}>favorite</span>
						</button>
					</div>
				</div>
			</div>

			<!-- Right: Product Details -->
			<div class="w-full md:w-1/2 flex flex-col pt-2 md:pt-4">
				
				<!-- Header Info & Rating -->
				<div class="flex flex-col gap-1.5 mb-space-lg">
					<div class="flex items-center justify-between">
						<span class="font-label-sm md:text-label-md text-label-sm uppercase tracking-widest text-primary">{product.category} • Hecho en Oaxaca</span>
						<div class="flex items-center gap-1 bg-surface-container px-2.5 py-1 rounded-full">
							<span class="material-symbols-outlined text-[14px] text-tertiary" style="font-variation-settings: 'FILL' 1;">star</span>
							<span class="font-label-sm text-label-sm text-on-surface">4.9</span>
							<span class="font-body-sm text-body-sm text-on-surface-variant">(128)</span>
						</div>
					</div>
					<div class="flex items-start justify-between mt-1">
						<h2 class="font-headline-lg-mobile md:text-display text-headline-lg-mobile text-on-surface leading-tight md:leading-tight">{product.name}</h2>
						<!-- Desktop Favorite Button -->
						<button aria-label="Guardar en favoritos" class="hidden md:flex w-12 h-12 rounded-full bg-surface-container-low border border-outline-variant/30 items-center justify-center text-on-surface-variant hover:text-primary transition-all shadow-sm ml-4 shrink-0" onclick={() => isFav = !isFav}>
							<span class="material-symbols-outlined text-[24px] transition-colors {isFav ? 'text-primary' : ''}" style={isFav ? "font-variation-settings: 'FILL' 1;" : ""}>favorite</span>
						</button>
					</div>
					<div class="flex items-baseline gap-space-xs mt-2 md:mt-4">
						<span class="font-display-mobile md:text-display-mobile text-display-mobile text-on-surface font-medium">${product.price.toFixed(2)}</span>
						<span class="font-label-sm text-label-sm text-on-surface-variant uppercase tracking-wider">USD</span>
					</div>
				</div>

				<hr class="border-outline-variant/30 w-full mb-space-lg" />

				<!-- Finish / Color Swatches -->
				<div class="mb-space-lg flex flex-col gap-space-xs">
					<div class="flex items-center justify-between mb-2">
						<label class="font-label-md text-label-md text-on-surface uppercase tracking-wider">Acabado</label>
						<span class="font-body-sm text-body-sm text-on-surface-variant font-medium">{selectedColor}</span>
					</div>
					<div class="flex items-center gap-space-md pt-1">
						<button class="relative w-12 h-12 md:w-14 md:h-14 rounded-full flex items-center justify-center p-0.5 shadow-sm transition-transform hover:scale-105 active:scale-95" onclick={() => selectedColor = 'Arena Natural'}>
							<span class="w-full h-full rounded-full bg-[#E5DFD3] shadow-inner flex items-center justify-center">
								{#if selectedColor === 'Arena Natural'}<span class="material-symbols-outlined text-[20px] text-inverse-surface">check</span>{/if}
							</span>
							{#if selectedColor === 'Arena Natural'}<span class="absolute -inset-1.5 rounded-full shadow-[0_0_0_2px_#994522]"></span>{/if}
						</button>
						<button class="relative w-12 h-12 md:w-14 md:h-14 rounded-full flex items-center justify-center p-0.5 shadow-sm transition-transform hover:scale-105 active:scale-95" onclick={() => selectedColor = 'Terracota Suave'}>
							<span class="w-full h-full rounded-full bg-[#B85D38] shadow-inner flex items-center justify-center">
								{#if selectedColor === 'Terracota Suave'}<span class="material-symbols-outlined text-[20px] text-on-primary">check</span>{/if}
							</span>
							{#if selectedColor === 'Terracota Suave'}<span class="absolute -inset-1.5 rounded-full shadow-[0_0_0_2px_#994522]"></span>{/if}
						</button>
					</div>
				</div>

				<!-- Collapsible Details -->
				<div class="mb-space-lg flex flex-col gap-space-xs">
					<!-- Accordion 1: Descripción -->
					<div class="bg-surface-container-low rounded-2xl overflow-hidden transition-colors border border-outline-variant/20">
						<button class="w-full flex items-center justify-between p-space-md md:p-5 text-left hover:bg-surface-container/50 transition-colors" onclick={() => descExpanded = !descExpanded}>
							<span class="font-headline-sm text-headline-sm text-on-surface">Historia & Proceso</span>
							<span class="material-symbols-outlined text-on-surface transition-transform duration-200 {descExpanded ? 'rotate-180' : ''}">expand_more</span>
						</button>
						{#if descExpanded}
						<div class="px-space-md md:px-5 pb-space-md md:pb-5 text-on-surface-variant font-body-md text-body-md flex flex-col gap-space-sm leading-relaxed">
							<p>{product.description}</p>
							<div class="grid grid-cols-2 gap-space-md pt-space-xs mt-2">
								<div class="flex items-center gap-2">
									<span class="material-symbols-outlined text-[20px] text-primary">verified</span>
									<span class="font-body-sm text-body-sm">Libre de plomo</span>
								</div>
								<div class="flex items-center gap-2">
									<span class="material-symbols-outlined text-[20px] text-primary">microwave</span>
									<span class="font-body-sm text-body-sm">Microondas apto</span>
								</div>
							</div>
						</div>
						{/if}
					</div>
				</div>

				<!-- Shipping and Returns Card -->
				<div class="mb-space-lg">
					<div class="bg-surface-container-high rounded-2xl p-space-md md:p-5 flex flex-col gap-space-md border border-outline-variant/20">
						<div class="flex items-start gap-space-sm">
							<span class="material-symbols-outlined text-[24px] text-primary mt-0.5">local_shipping</span>
							<div class="flex flex-col">
								<span class="font-label-md text-label-md text-on-surface">Envío Gratis sobre $50 USD</span>
								<span class="font-body-sm text-body-sm text-on-surface-variant mt-1">Entrega estimada en 3 a 5 días hábiles con empaque 100% reciclable y protegido.</span>
							</div>
						</div>
						<hr class="border-outline-variant/30" />
						<div class="flex items-start gap-space-sm">
							<span class="material-symbols-outlined text-[24px] text-secondary mt-0.5">assignment_return</span>
							<div class="flex flex-col">
								<span class="font-label-md text-label-md text-on-surface">Devoluciones sin costo por 30 días</span>
								<span class="font-body-sm text-body-sm text-on-surface-variant mt-1">Si la textura no te enamora, el retorno es completamente libre de trámites complicados.</span>
							</div>
						</div>
					</div>
				</div>

				<!-- Desktop Add to Cart (Inline) -->
				<div class="hidden md:flex items-center gap-space-md mt-auto pt-4 border-t border-outline-variant/30">
					<!-- Stepper -->
					<div class="flex items-center bg-surface-container-low border border-outline-variant/40 rounded-xl h-14 px-2">
						<button class="w-10 h-10 flex items-center justify-center text-on-surface hover:bg-surface-container-highest rounded-lg transition-colors" onclick={() => { if(quantity > 1) quantity-- }}>
							<span class="material-symbols-outlined text-[20px]">remove</span>
						</button>
						<span class="w-12 text-center font-label-md text-label-md text-on-surface select-none">{quantity}</span>
						<button class="w-10 h-10 flex items-center justify-center text-on-surface hover:bg-surface-container-highest rounded-lg transition-colors" onclick={() => quantity++}>
							<span class="material-symbols-outlined text-[20px]">add</span>
						</button>
					</div>
					<!-- Add to Bag CTA -->
					<button class="flex-1 h-14 bg-primary-container hover:bg-primary text-on-primary-container hover:text-on-primary font-headline-sm text-headline-sm rounded-xl flex items-center justify-center gap-3 shadow-md hover:shadow-lg active:scale-[0.98] transition-all {adding ? 'opacity-80' : ''}" onclick={addToCart} disabled={adding}>
						{#if adding}
							<span class="material-symbols-outlined animate-spin text-[20px]">progress_activity</span>
						{:else}
							<span class="material-symbols-outlined text-[20px]">shopping_bag</span>
						{/if}
						<span>Agregar a la bolsa • ${(product.price * quantity).toFixed(2)}</span>
					</button>
				</div>
			</div>
		</div>
		
		<!-- Fixed Bottom Conversion Bar (Mobile Only) -->
		<div class="md:hidden fixed bottom-0 inset-x-0 z-40 bg-surface/95 backdrop-blur-xl pb-safe shadow-[0_-8px_24px_rgba(28,27,27,0.06)]">
			<div class="max-w-md mx-auto px-margin py-space-sm flex items-center gap-space-sm">
				<!-- Stepper -->
				<div class="flex items-center bg-surface-container border border-outline-variant/20 rounded-xl h-12 px-1">
					<button class="w-9 h-9 flex items-center justify-center text-on-surface active:bg-surface-container-highest rounded-lg transition-colors" onclick={() => { if(quantity > 1) quantity-- }}>
						<span class="material-symbols-outlined text-[18px]">remove</span>
					</button>
					<span class="w-8 text-center font-label-md text-label-md text-on-surface select-none">{quantity}</span>
					<button class="w-9 h-9 flex items-center justify-center text-on-surface active:bg-surface-container-highest rounded-lg transition-colors" onclick={() => quantity++}>
						<span class="material-symbols-outlined text-[18px]">add</span>
					</button>
				</div>
				<!-- Add to Bag CTA -->
				<button class="flex-1 h-12 bg-primary-container text-on-primary-container font-label-md text-label-md rounded-xl flex items-center justify-center gap-2 shadow-md active:scale-[0.98] transition-all {adding ? 'opacity-80' : ''}" onclick={addToCart} disabled={adding}>
					{#if adding}
						<span class="material-symbols-outlined animate-spin text-[18px]">progress_activity</span>
					{:else}
						<span class="material-symbols-outlined text-[18px]">shopping_bag</span>
					{/if}
					<span>Agregar • ${(product.price * quantity).toFixed(2)}</span>
				</button>
			</div>
		</div>
		
		<!-- Toast Notification -->
		<div class="fixed top-20 md:bottom-10 md:top-auto left-1/2 -translate-x-1/2 z-[100] bg-inverse-surface text-inverse-on-surface px-space-lg py-space-sm rounded-full shadow-2xl flex items-center gap-3 pointer-events-none transition-all duration-300 transform {toastVisible ? 'opacity-100 translate-y-0' : 'opacity-0 -translate-y-4 md:translate-y-4'}">
			<span class="material-symbols-outlined text-[20px] text-primary-fixed">check_circle</span>
			<span class="font-label-md text-label-md">{quantity} pieza(s) agregada(s)</span>
		</div>
	</main>
{/if}
