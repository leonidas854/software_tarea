<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	interface CartItem {
		product_id: string;
		name: string;
		price: number;
		quantity: number;
		image: string;
		subtitle: string;
	}

	let cartItems = $state<CartItem[]>([]);
	let subtotal = $state(0);
	let loading = $state(true);
	let error = $state('');

	// Local state for UI
	let discountApplied = $state(false);
	let discountRate = 0.10;
	let couponInput = $state('');
	let couponStatus = $state('Aplicar');

	const API = 'http://localhost:8080/api';

	async function fetchCart() {
		try {
			const res = await fetch(`${API}/cart`, { credentials: 'include' });
			if (res.status === 401) { goto('/'); return; }
			if (res.ok) {
				const data = await res.json();
				// Ensure cart items are sorted or consistent
				cartItems = data.items || [];
				cartItems.sort((a, b) => a.name.localeCompare(b.name));
				subtotal = data.subtotal || 0;
			}
		} catch {
			error = 'Error cargando el carrito';
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		fetchCart();
	});

	async function updateQuantity(productId: string, qty: number) {
		if (qty < 1) {
			removeItem(productId);
			return;
		}

		// Optimistic update
		const item = cartItems.find(i => i.product_id === productId);
		if (item) {
			item.quantity = qty;
			recalcLocalSubtotal();
		}

		try {
			const res = await fetch(`${API}/cart`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ product_id: productId, quantity: qty }),
				credentials: 'include'
			});
			if (res.ok) {
				const data = await res.json();
				cartItems = data.items || [];
				cartItems.sort((a: any, b: any) => a.name.localeCompare(b.name));
				subtotal = data.subtotal || 0;
			}
		} catch (e) {
			console.error('Failed to update quantity');
		}
	}

	async function removeItem(productId: string) {
		// Optimistic update
		cartItems = cartItems.filter(i => i.product_id !== productId);
		recalcLocalSubtotal();

		try {
			const res = await fetch(`${API}/cart`, {
				method: 'DELETE',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ product_id: productId }),
				credentials: 'include'
			});
			if (res.ok) {
				const data = await res.json();
				cartItems = data.items || [];
				cartItems.sort((a: any, b: any) => a.name.localeCompare(b.name));
				subtotal = data.subtotal || 0;
			}
		} catch (e) {
			console.error('Failed to remove item');
		}
	}

	function recalcLocalSubtotal() {
		subtotal = cartItems.reduce((sum, item) => sum + (item.price * item.quantity), 0);
	}

	function applyCoupon() {
		if (couponInput.trim().toUpperCase() === 'TIENDITA10') {
			discountApplied = true;
			couponStatus = 'Aplicado';
		} else {
			// Trigger shake animation via DOM in real life, here we just change text briefly
			const old = couponStatus;
			couponStatus = 'Inválido';
			setTimeout(() => couponStatus = old, 1500);
		}
	}

	let discount = $derived(discountApplied ? subtotal * discountRate : 0);
	let shipping = $derived(subtotal >= 50 ? 0 : (subtotal > 0 ? 5 : 0));
	let total = $derived(Math.max(0, subtotal - discount + shipping));
	let totalItems = $derived(cartItems.reduce((acc, i) => acc + i.quantity, 0));
	let progress = $derived(subtotal >= 50 ? 100 : (subtotal / 50) * 100);

</script>

<header class="md:hidden fixed top-0 inset-x-0 z-50 bg-surface/90 backdrop-blur-xl shadow-[0_1px_8px_rgba(0,0,0,0.03)] pt-safe">
	<div class="h-14 px-margin flex items-center justify-between max-w-md mx-auto">
		<div class="flex items-center gap-space-sm">
			<img alt="La Tiendita Logo" class="h-8 w-auto object-contain" src="https://lh3.googleusercontent.com/aida/AEtjO1VdBD_zfQ5Q4HTvgfcK5U0GnN2R44-whuA_1qYMzndpSPV9l9d0oVrZeWqSXd_HKMQOO-KwnLuhIXa4haTayD8W9M48rvvWefKRWmJ0z4IPVuEmFewqqXZYa3wdptuse-m5xeYEMABjtwKHx0ECJ4SVkSLVL54kBeUUcttNOuFZbbxjWFOWkSV5yu1tyPSlNVApOVRZla72Utj1BJ5PBlN0CLsgqR477UZ-23_3RV_bIrDLNNJkZdDMtd8"/>
			<span class="font-headline-sm text-headline-sm text-on-surface font-medium">La Tiendita</span>
		</div>
		<div class="flex items-center gap-space-xs">
			<button class="w-11 h-11 flex items-center justify-center rounded-full text-on-surface-variant hover:text-on-surface transition-colors" onclick={() => goto('/home')}>
				<span class="material-symbols-outlined text-[22px]">search</span>
			</button>
			<a class="w-11 h-11 flex items-center justify-center rounded-full" href="javascript:void(0)">
				<img alt="Profile" class="w-8 h-8 rounded-full object-cover" src="https://lh3.googleusercontent.com/aida-public/AB6AXuAcPHqh5WXHDRm1kdiadEXUJj_ummit40yNuJIScfiZfEy_4Vnw0kYgWLrsdpeRGfp8jbi8LydFPn-mMIKVPS-urlaggjvxtlcfx0gT9K-m1zjP_7uwj_YALXykDky4FvDiPByH3pqhIHg8q1bGQCBmoYuIxc7eQ92NEH1Yabt64S9jApLnEiSAZKKfFEQ4oD_aIu4B0LVMHyDqY1KvgsKqYxQvGpxEPa3pMXc8dQExapbK81fEEz_E"/>
			</a>
		</div>
	</div>
</header>

<main class="flex-1 flex flex-col relative w-full pt-14 md:pt-4 pb-20 md:pb-12 bg-surface max-w-6xl mx-auto">
	<div class="flex flex-col w-full px-margin pb-space-lg">
		
		<!-- Header / Breadcrumb Minimal Info -->
		<div class="flex items-baseline justify-between pt-space-sm md:pt-space-lg mb-space-md">
			<h1 class="font-headline-lg-mobile md:text-display text-headline-lg-mobile text-on-surface tracking-tight">Mi Bolsa</h1>
			<span class="font-label-md md:text-label-lg text-label-md text-on-surface-variant bg-surface-container px-3 py-1.5 rounded-full">{totalItems} piezas</span>
		</div>

		<!-- 2-Column Layout on Desktop -->
		<div class="flex flex-col lg:flex-row gap-space-lg lg:gap-12 items-start">
			
			<!-- Left Column: Items & Shipping -->
			<div class="w-full lg:w-[60%] flex flex-col space-y-space-md">
				
				<!-- Free Shipping Milestone Bar -->
				<div class="bg-surface-container-low p-space-md md:p-6 rounded-2xl shadow-sm border border-outline-variant/20 flex flex-col space-y-3">
					<div class="flex items-center justify-between text-on-surface">
						<div class="flex items-center space-x-3">
							<span class="material-symbols-outlined text-primary text-[24px]" style="font-variation-settings: 'FILL' 1;">local_shipping</span>
							<p class="font-body-sm md:text-body-md text-body-sm">
								{#if subtotal >= 50}
									¡Felicidades! Tienes <strong class="font-medium text-on-surface">Envío Gratis</strong>
								{:else}
									Te faltan <span class="font-label-md md:text-label-lg text-label-md text-primary font-semibold">${(50 - subtotal).toFixed(2)}</span> para <strong class="font-medium text-on-surface">Envío Gratis</strong>
								{/if}
							</p>
						</div>
						<span class="font-label-sm md:text-label-md text-label-sm text-on-surface-variant bg-surface-container px-2 py-0.5 rounded-md">{Math.min(Math.round(progress), 100)}%</span>
					</div>
					<div class="w-full bg-surface-container-highest h-2 rounded-full overflow-hidden">
						<div class="bg-primary h-full rounded-full transition-all duration-500 ease-out" style="width: {progress}%;"></div>
					</div>
				</div>

				<!-- Cart Items List -->
				<div class="flex flex-col space-y-space-sm">
					{#if loading}
						<div class="text-center py-12 text-on-surface-variant">Cargando bolsa...</div>
					{:else if cartItems.length === 0}
						<div class="text-center py-16 bg-surface-container-lowest rounded-2xl border border-outline-variant/30 text-on-surface-variant flex flex-col items-center">
							<span class="material-symbols-outlined text-[48px] text-surface-container-highest mb-4">shopping_bag</span>
							<p class="font-headline-sm">Tu bolsa está vacía.</p>
							<p class="font-body-sm mt-1 mb-6">Explora nuestra colección para llenarla de calma.</p>
							<button class="bg-primary text-on-primary font-label-md text-label-md py-3.5 rounded-xl w-48 shadow-md hover:bg-primary-container transition-colors" onclick={() => goto('/home')}>Explorar catálogo</button>
						</div>
					{:else}
						{#each cartItems as item (item.product_id)}
							<div class="bg-surface-container-lowest p-space-sm md:p-4 rounded-2xl shadow-sm border border-outline-variant/20 flex space-x-space-md md:space-x-6 relative transition-transform duration-200 hover:shadow-md">
								<div class="w-24 h-28 md:w-32 md:h-36 rounded-xl overflow-hidden bg-surface-container-low flex-shrink-0 cursor-pointer" onclick={() => goto(`/product/${item.product_id}`)}>
									<img alt={item.name} class="w-full h-full object-cover hover:scale-105 transition-transform" src={item.image}/>
								</div>
								<div class="flex flex-col justify-between flex-1 min-w-0 py-0.5 md:py-2">
									<div class="flex justify-between items-start">
										<div class="min-w-0 pr-2">
											<h2 class="font-headline-sm md:text-headline-md text-headline-sm text-on-surface truncate cursor-pointer hover:text-primary transition-colors" onclick={() => goto(`/product/${item.product_id}`)}>{item.name}</h2>
											<p class="font-body-sm md:text-body-md text-body-sm text-on-surface-variant mt-1">{item.subtitle}</p>
										</div>
										<button aria-label="Eliminar" class="text-outline hover:text-error hover:bg-error-container transition-colors p-1.5 flex items-center justify-center rounded-full" onclick={() => removeItem(item.product_id)}>
											<span class="material-symbols-outlined text-[20px]">close</span>
										</button>
									</div>
									<div class="flex items-center justify-between mt-2 md:mt-auto">
										<div class="flex items-center bg-surface-container-low border border-outline-variant/30 rounded-lg h-9 md:h-11 px-1 space-x-1 md:space-x-2">
											<button aria-label="Disminuir" class="w-7 h-7 md:w-9 md:h-9 flex items-center justify-center text-on-surface hover:bg-surface-container-highest rounded active:scale-95 transition-all" onclick={() => updateQuantity(item.product_id, item.quantity - 1)}>
												<span class="material-symbols-outlined text-[16px] md:text-[20px]">remove</span>
											</button>
											<span class="font-label-md md:text-label-lg text-label-md text-on-surface w-4 md:w-6 text-center">{item.quantity}</span>
											<button aria-label="Aumentar" class="w-7 h-7 md:w-9 md:h-9 flex items-center justify-center text-on-surface hover:bg-surface-container-highest rounded active:scale-95 transition-all" onclick={() => updateQuantity(item.product_id, item.quantity + 1)}>
												<span class="material-symbols-outlined text-[16px] md:text-[20px]">add</span>
											</button>
										</div>
										<span class="font-headline-sm md:text-headline-md text-headline-sm text-on-surface font-semibold tracking-tight">${(item.price * item.quantity).toFixed(2)}</span>
									</div>
								</div>
							</div>
						{/each}
					{/if}
				</div>
			</div>

			<!-- Right Column: Summary (Sticky on Desktop) -->
			{#if cartItems.length > 0}
				<div class="w-full lg:w-[40%] flex flex-col space-y-space-md lg:sticky lg:top-24">
					<!-- Coupon Code Component -->
					<div class="bg-surface-container-lowest p-space-sm md:p-3 rounded-2xl shadow-sm border border-outline-variant/20 flex items-center space-x-2">
						<div class="relative flex-1 flex items-center">
							<span class="material-symbols-outlined absolute left-3 md:left-4 text-outline text-[20px]">sell</span>
							<input bind:value={couponInput} class="w-full bg-surface-container-low pl-10 md:pl-12 pr-3 py-3 md:py-3.5 rounded-xl text-on-surface font-body-sm md:text-body-md text-body-sm focus:outline-none focus:ring-1 focus:ring-primary uppercase tracking-wider placeholder:text-outline/70 border border-transparent transition-all" placeholder="Código promocional" type="text"/>
						</div>
						<button class="{discountApplied ? 'bg-secondary-container text-on-secondary-container' : 'bg-secondary hover:bg-secondary/90 text-on-secondary'} font-label-md md:text-label-lg text-label-md px-4 md:px-6 py-3 md:py-3.5 rounded-xl active:scale-95 transition-transform" onclick={applyCoupon}>
							{couponStatus}
						</button>
					</div>

					<!-- Order Breakdown -->
					<div class="bg-surface-container-low p-space-md md:p-6 rounded-2xl space-y-4 shadow-sm border border-outline-variant/20">
						<h3 class="font-label-md md:text-label-lg text-label-md text-on-surface-variant uppercase tracking-wider mb-2">Resumen del Pedido</h3>
						<div class="flex justify-between items-center">
							<span class="font-body-md md:text-body-lg text-body-md text-on-surface-variant">Subtotal</span>
							<span class="font-label-md md:text-label-lg text-label-md text-on-surface">${subtotal.toFixed(2)}</span>
						</div>
						{#if discountApplied && discount > 0}
							<div class="flex justify-between items-center">
								<span class="font-body-md md:text-body-lg text-body-md text-primary flex items-center gap-1.5">
									<span class="material-symbols-outlined text-[18px]">local_offer</span> Descuento (10%)
								</span>
								<span class="font-label-md md:text-label-lg text-label-md text-primary font-semibold">-${discount.toFixed(2)}</span>
							</div>
						{/if}
						<div class="flex justify-between items-center">
							<div class="flex items-center space-x-1.5">
								<span class="font-body-md md:text-body-lg text-body-md text-on-surface-variant">Envío estimado</span>
								<span class="material-symbols-outlined text-outline text-[18px]" title="Calculado para estándar">info</span>
							</div>
							<span class="font-label-md md:text-label-lg text-label-md text-on-surface">{shipping === 0 ? 'Gratis' : `$${shipping.toFixed(2)}`}</span>
						</div>
						<hr class="border-outline-variant/30 my-2" />
						<div class="flex justify-between items-center pt-1">
							<div>
								<p class="font-headline-sm md:text-headline-md text-headline-sm text-on-surface font-medium">Total</p>
								<p class="font-body-sm text-body-sm text-on-surface-variant">Impuestos incluidos</p>
							</div>
							<span class="font-display-mobile md:text-display text-display-mobile text-on-surface font-medium">${total.toFixed(2)}</span>
						</div>
					</div>

					<!-- Trust & Security Badges -->
					<div class="flex flex-col items-center justify-center py-2 space-y-2.5">
						<div class="flex items-center space-x-4 text-on-surface-variant bg-surface-container-lowest px-4 py-2 rounded-full border border-outline-variant/20 shadow-sm">
							<div class="flex items-center space-x-1.5">
								<span class="material-symbols-outlined text-[20px] text-secondary">verified_user</span>
								<span class="font-label-sm md:text-label-md text-label-sm">Pago 100% Seguro</span>
							</div>
							<span class="text-outline text-sm opacity-50">•</span>
							<div class="flex items-center space-x-2 opacity-75">
								<span class="material-symbols-outlined text-[20px]">credit_card</span>
								<span class="material-symbols-outlined text-[20px]">contactless</span>
								<span class="material-symbols-outlined text-[20px]">account_balance_wallet</span>
							</div>
						</div>
					</div>

					<!-- Primary Sticky/Prominent Conversion Button -->
					<div class="pt-2">
						<button class="w-full bg-primary text-on-primary font-headline-sm md:text-headline-md text-headline-sm py-4 rounded-2xl shadow-md hover:shadow-xl hover:bg-primary-container hover:text-on-primary-container flex items-center justify-center space-x-2.5 active:scale-[0.98] transition-all duration-200" onclick={() => goto('/checkout')}>
							<span class="material-symbols-outlined text-[22px]">lock</span>
							<span>Proceder al Pago • ${total.toFixed(2)}</span>
						</button>
					</div>
				</div>
			{/if}
		</div>
	</div>
</main>
