<script lang="ts">
	import { goto } from '$app/navigation';

	let username = '';
	let password = '';
	let error = '';
	let loading = false;

	async function handleLogin() {
		// Client-side validation (JS)
		if (!username.trim() || !password.trim()) {
			error = 'Por favor ingresa usuario y contraseña.';
			return;
		}

		error = '';
		loading = true;

		try {
			const res = await fetch('http://localhost:8080/api/login', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ username, password }),
				credentials: 'include' // Allow cookies
			});

			if (!res.ok) {
				const text = await res.text();
				error = text || 'Credenciales inválidas.';
			} else {
				// Navigate to cart
				goto('/cart');
			}
		} catch (err) {
			error = 'Error de conexión con el servidor.';
		} finally {
			loading = false;
		}
	}
</script>

<main>
	<div class="glass-card">
		<h1>Iniciar Sesión</h1>
		
		<form on:submit|preventDefault={handleLogin}>
			<div class="form-group">
				<label for="username">Usuario</label>
				<input 
					type="text" 
					id="username" 
					bind:value={username} 
					placeholder="admin"
				/>
			</div>
			
			<div class="form-group">
				<label for="password">Contraseña</label>
				<input 
					type="password" 
					id="password" 
					bind:value={password} 
					placeholder="1234"
				/>
			</div>

			{#if error}
				<div class="error-msg">{error}</div>
			{/if}

			<button type="submit" disabled={loading}>
				{loading ? 'Iniciando...' : 'Entrar'}
			</button>
		</form>
	</div>
</main>
