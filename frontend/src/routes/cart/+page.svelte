<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	interface Product {
		id: string;
		name: string;
		price: number;
		quantity: number;
	}

	let cartItems = $state<Product[]>([]);
	let subtotal = $state(0);
	let error = $state('');
	
	// Form inputs
	let prodId = $state('');
	let prodName = $state('');
	let prodPrice = $state('');
	let prodQuantity = $state(1);
	let loading = $state(false);

	async function fetchCart() {
		try {
			const res = await fetch('http://localhost:8080/api/cart', {
				credentials: 'include'
			});
			if (res.status === 401) {
				goto('/');
				return;
			}
			if (res.ok) {
				const data = await res.json();
				cartItems = data.items || [];
				subtotal = data.subtotal || 0;
			}
		} catch (err) {
			error = 'Error cargando el carrito';
		}
	}

	onMount(() => {
		fetchCart();
	});

	async function addToCart() {
		error = '';
		
		// Client-side validation
		if (!prodId.trim() || !prodName.trim() || !prodPrice) {
			error = 'Todos los campos son obligatorios.';
			return;
		}
		if (prodQuantity <= 0 || !Number.isInteger(prodQuantity)) {
			error = 'La cantidad debe ser un número entero mayor a 0.';
			return;
		}

		loading = true;
		try {
			const res = await fetch('http://localhost:8080/api/cart', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					id: prodId,
					name: prodName,
					price: parseFloat(prodPrice),
					quantity: parseInt(prodQuantity.toString())
				}),
				credentials: 'include'
			});

			if (res.status === 401) {
				goto('/');
				return;
			}

			if (!res.ok) {
				const errText = await res.text();
				error = errText || 'Error al agregar el producto.';
			} else {
				// Reset form on success
				prodId = '';
				prodName = '';
				prodPrice = '';
				prodQuantity = 1;
				// Refresh cart
				await fetchCart();
			}
		} catch (err) {
			error = 'Error de conexión.';
		} finally {
			loading = false;
		}
	}

	async function logout() {
		try {
			await fetch('http://localhost:8080/api/logout', { method: 'POST', credentials: 'include' });
		} catch(e) {}
		goto('/');
	}
</script>

<nav class="navbar">
	<a href="#" class="nav-brand">ShopApp</a>
	<div class="nav-actions">
		<button on:click={logout}>Cerrar Sesión</button>
	</div>
</nav>

<div style="width: 100%; max-width: 800px; margin: 2rem auto; display: grid; grid-template-columns: 1fr 1.5fr; gap: 2rem;">
	<!-- Add Product Form -->
	<div class="glass-card" style="margin-top: 0;">
		<h2>Agregar Producto</h2>
		<form on:submit|preventDefault={addToCart} style="margin-top: 1.5rem;">
			<div class="form-group">
				<label for="id">ID del Producto</label>
				<input type="text" id="id" bind:value={prodId} placeholder="EJ: P001">
			</div>
			<div class="form-group">
				<label for="name">Nombre</label>
				<input type="text" id="name" bind:value={prodName} placeholder="Laptop Gamer">
			</div>
			<div class="form-group">
				<label for="price">Precio ($)</label>
				<input type="number" step="0.01" id="price" bind:value={prodPrice} placeholder="1500.00">
			</div>
			<div class="form-group">
				<label for="quantity">Cantidad</label>
				<!-- HTML validation pattern, but JS handles the actual error message -->
				<input type="number" id="quantity" bind:value={prodQuantity} min="1">
			</div>

			{#if error}
				<div class="error-msg">{error}</div>
			{/if}

			<button type="submit" disabled={loading}>
				{loading ? 'Agregando...' : 'Agregar al Carrito'}
			</button>
		</form>
	</div>

	<!-- Cart View -->
	<div class="glass-card" style="margin-top: 0;">
		<h2>Tu Carrito</h2>
		<div style="margin-top: 1.5rem;">
			{#if cartItems.length === 0}
				<p style="color: #cbd5e1; text-align: center;">El carrito está vacío.</p>
			{:else}
				{#each cartItems as item}
					<div class="cart-item">
						<div>
							<strong>{item.name}</strong> (x{item.quantity})
							<div style="font-size: 0.85rem; color: #cbd5e1;">ID: {item.id}</div>
						</div>
						<div>
							${(item.price * item.quantity).toFixed(2)}
						</div>
					</div>
				{/each}
				
				<div class="cart-subtotal">
					Total: ${subtotal.toFixed(2)}
				</div>
			{/if}
		</div>
	</div>
</div>
