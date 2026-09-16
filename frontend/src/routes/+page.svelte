<script lang="ts">
	import { goto } from '$app/navigation';

	let username = $state('');
	let password = $state('');
	let error = $state('');
	let loading = $state(false);

	async function handleLogin(e: SubmitEvent) {
		e.preventDefault();

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
				credentials: 'include'
			});

			if (!res.ok) {
				const text = await res.text();
				error = text || 'Credenciales inválidas.';
			} else {
				goto('/home');
			}
		} catch (err) {
			error = 'Error de conexión con el servidor.';
		} finally {
			loading = false;
		}
	}
</script>

<div class="min-h-screen flex flex-col md:flex-row bg-surface">
	<!-- Left Side: Image (Desktop only) or Background (Mobile) -->
	<div class="relative w-full md:w-1/2 flex-1 md:flex-none flex flex-col justify-center bg-[url('https://lh3.googleusercontent.com/aida-public/AB6AXuAZNVmp6Zyu0QXymxzy1MvtwGqEgFHDl7MP9kLIHsjPbDHiawrTy1jEdsvcCPnJvFZkS4NOlaMQmLJvm_wnLcZOjVv7z3f7Aw1rqDDGG2SuTdWAUaSNuyw7H16qhW5zCYk-0A9-kTV3vLk2ew2CypP_eOmfbfOSsEcWN_fZ2Qgd8HzGy6BmRHxJlzTZ6F7cMTX45nZm5eryLm8lxQoV1LOwKtihTnsqIlDQHl7DrpuK2HPx1cLgiLcD')] bg-cover bg-center bg-no-repeat min-h-[40vh] md:min-h-screen">
		<!-- Mobile backdrop -->
		<div class="absolute inset-0 bg-surface/85 backdrop-blur-md md:hidden"></div>
		
		<!-- Desktop subtle overlay -->
		<div class="hidden md:block absolute inset-0 bg-black/20"></div>

		<!-- Branding (Centered on both) -->
		<div class="relative z-10 w-full flex flex-col items-center gap-space-xs p-margin">
			<img alt="La Tiendita Logo" class="h-20 w-auto object-contain drop-shadow-md" src="https://lh3.googleusercontent.com/aida/AEtjO1VdBD_zfQ5Q4HTvgfcK5U0GnN2R44-whuA_1qYMzndpSPV9l9d0oVrZeWqSXd_HKMQOO-KwnLuhIXa4haTayD8W9M48rvvWefKRWmJ0z4IPVuEmFewqqXZYa3wdptuse-m5xeYEMABjtwKHx0ECJ4SVkSLVL54kBeUUcttNOuFZbbxjWFOWkSV5yu1tyPSlNVApOVRZla72Utj1BJ5PBlN0CLsgqR477UZ-23_3RV_bIrDLNNJkZdDMtd8"/>
			<h1 class="font-display text-display-mobile md:text-display text-on-surface md:text-white text-center drop-shadow-lg">Bienvenido</h1>
			<p class="font-body-sm text-body-sm text-on-surface-variant md:text-white/90 text-center max-w-xs md:max-w-sm drop-shadow-md">
				Ingresa con <span class="font-bold text-primary md:text-primary-fixed">admin / 1234</span> para acceder a tu sesión.
			</p>
		</div>
	</div>

	<!-- Right Side: Form -->
	<div class="w-full md:w-1/2 flex items-center justify-center p-margin -mt-10 md:mt-0 relative z-20">
		<div class="w-full max-w-sm bg-surface-container-low md:bg-transparent rounded-2xl md:rounded-none p-space-lg md:p-0 shadow-xl md:shadow-none border md:border-none border-outline-variant/30 flex flex-col gap-space-md">
			<div class="hidden md:block mb-space-sm">
				<h2 class="font-headline-md text-headline-md text-on-surface">Iniciar Sesión</h2>
				<p class="font-body-md text-body-md text-on-surface-variant mt-1">Ingresa tus credenciales para continuar.</p>
			</div>

			<form onsubmit={handleLogin} class="flex flex-col gap-space-md">
				<div class="flex flex-col gap-1.5">
					<label for="username" class="font-label-sm text-label-sm text-on-surface-variant uppercase tracking-wider">Usuario</label>
					<div class="relative flex items-center">
						<span class="absolute left-3 top-1/2 -translate-y-1/2 material-symbols-outlined text-outline text-[20px] pointer-events-none">person</span>
						<input 
							type="text" 
							id="username" 
							bind:value={username} 
							placeholder="admin"
							class="w-full h-12 bg-surface-container-lowest rounded-xl pl-11 pr-4 font-body-md text-body-md text-on-surface placeholder:text-outline focus:outline-none focus:ring-2 focus:ring-primary transition-all border border-outline-variant/30 focus:border-transparent"
						/>
					</div>
				</div>
				
				<div class="flex flex-col gap-1.5">
					<label for="password" class="font-label-sm text-label-sm text-on-surface-variant uppercase tracking-wider">Contraseña</label>
					<div class="relative flex items-center">
						<span class="absolute left-3 top-1/2 -translate-y-1/2 material-symbols-outlined text-outline text-[20px] pointer-events-none">lock</span>
						<input 
							type="password" 
							id="password" 
							bind:value={password} 
							placeholder="1234"
							class="w-full h-12 bg-surface-container-lowest rounded-xl pl-11 pr-4 font-body-md text-body-md text-on-surface placeholder:text-outline focus:outline-none focus:ring-2 focus:ring-primary transition-all border border-outline-variant/30 focus:border-transparent"
						/>
					</div>
				</div>

				{#if error}
					<div class="bg-error-container text-on-error-container font-label-sm text-label-sm p-3 rounded-xl flex items-center gap-2 mt-2 animate-pulse">
						<span class="material-symbols-outlined text-[18px]">error</span>
						{error}
					</div>
				{/if}

				<button 
					type="submit" 
					disabled={loading}
					class="w-full h-12 bg-primary text-on-primary font-headline-sm text-headline-sm rounded-xl flex items-center justify-center gap-2 shadow-md active:scale-[0.98] transition-all hover:bg-primary-container hover:shadow-lg disabled:opacity-60 disabled:scale-100 mt-4"
				>
					{#if loading}
						<span class="material-symbols-outlined animate-spin text-[20px]">progress_activity</span>
						Entrando...
					{:else}
						<span>Entrar a la tienda</span>
						<span class="material-symbols-outlined text-[18px]">arrow_forward</span>
					{/if}
				</button>
			</form>
		</div>
	</div>
</div>
