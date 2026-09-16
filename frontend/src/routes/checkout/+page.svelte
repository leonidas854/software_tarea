<script lang="ts">
	import { goto } from '$app/navigation';
</script>

<header class="md:hidden fixed top-0 inset-x-0 z-50 bg-surface/90 backdrop-blur-xl shadow-[0_1px_8px_rgba(0,0,0,0.03)] pt-safe">
	<div class="h-14 px-margin flex items-center justify-between max-w-md mx-auto">
		<div class="flex items-center gap-space-sm">
			<img alt="La Tiendita Logo" class="h-8 w-auto object-contain" src="https://lh3.googleusercontent.com/aida/AEtjO1VdBD_zfQ5Q4HTvgfcK5U0GnN2R44-whuA_1qYMzndpSPV9l9d0oVrZeWqSXd_HKMQOO-KwnLuhIXa4haTayD8W9M48rvvWefKRWmJ0z4IPVuEmFewqqXZYa3wdptuse-m5xeYEMABjtwKHx0ECJ4SVkSLVL54kBeUUcttNOuFZbbxjWFOWkSV5yu1tyPSlNVApOVRZla72Utj1BJ5PBlN0CLsgqR477UZ-23_3RV_bIrDLNNJkZdDMtd8"/>
			<span class="font-headline-sm text-headline-sm text-on-surface font-medium">Pago Seguro</span>
		</div>
		<button class="w-11 h-11 flex items-center justify-center rounded-full text-on-surface-variant hover:text-on-surface transition-colors" onclick={() => goto('/cart')}>
			<span class="material-symbols-outlined text-[22px]">close</span>
		</button>
	</div>
</header>

<main class="flex-1 flex flex-col relative w-full pt-14 pb-8 bg-surface max-w-md mx-auto min-h-screen">
	<!-- Mockup for checkout since the user provided a static one -->
	<div class="px-margin flex flex-col space-y-space-md pt-space-md">
		<!-- Express Checkout Options -->
		<div class="bg-surface-container-low rounded-xl p-space-sm shadow-sm flex flex-col gap-2">
			<span class="font-label-sm text-label-sm text-on-surface-variant uppercase tracking-wider text-center mb-1">Pago Rápido</span>
			<button class="w-full bg-[#000000] text-[#ffffff] h-12 rounded-lg flex items-center justify-center gap-2 hover:opacity-90 transition-opacity" onclick={() => goto('/home')}>
				<span class="material-symbols-outlined text-[24px]" style="font-variation-settings: 'FILL' 1;">apple</span>
				<span class="font-headline-sm text-headline-sm tracking-tight font-medium">Pay</span>
			</button>
			<button class="w-full bg-[#FFC439] text-[#000000] h-12 rounded-lg flex items-center justify-center hover:opacity-90 transition-opacity" onclick={() => goto('/home')}>
				<span class="font-headline-sm text-headline-sm font-semibold tracking-tighter italic">PayPal</span>
			</button>
		</div>

		<div class="flex items-center gap-2 my-1 opacity-70">
			<hr class="flex-1 border-outline-variant" />
			<span class="font-label-sm text-label-sm text-outline">o paga con tarjeta</span>
			<hr class="flex-1 border-outline-variant" />
		</div>

		<!-- Card Payment Form Mockup -->
		<div class="bg-surface-container-low rounded-xl p-space-md shadow-sm">
			<form class="flex flex-col gap-space-sm" onsubmit={(e) => { e.preventDefault(); goto('/home'); }}>
				
				<div class="flex flex-col gap-1">
					<label class="font-label-sm text-label-sm text-on-surface-variant uppercase tracking-wider">Correo electrónico</label>
					<input class="w-full h-11 bg-surface-container-lowest rounded-lg px-3 font-body-md text-body-md text-on-surface focus:outline-none focus:ring-1 focus:ring-primary border border-outline-variant/30" placeholder="tu@email.com" type="email" required/>
				</div>
				
				<div class="flex flex-col gap-1 pt-2">
					<label class="font-label-sm text-label-sm text-on-surface-variant uppercase tracking-wider">Información de la tarjeta</label>
					<div class="flex flex-col rounded-lg border border-outline-variant/30 overflow-hidden bg-surface-container-lowest focus-within:ring-1 focus-within:ring-primary focus-within:border-primary">
						<div class="relative flex items-center border-b border-outline-variant/30 px-3">
							<span class="material-symbols-outlined text-outline text-[20px] absolute left-3">credit_card</span>
							<input class="w-full h-11 bg-transparent pl-8 font-body-md text-body-md text-on-surface focus:outline-none" placeholder="Número de tarjeta" type="text" required/>
							<div class="flex items-center gap-1 opacity-70 absolute right-3 pointer-events-none">
								<div class="w-8 h-5 bg-surface-container-high rounded flex items-center justify-center">
									<div class="w-4 h-4 rounded-full bg-[#EB001B]/80 -mr-1.5 mix-blend-multiply"></div>
									<div class="w-4 h-4 rounded-full bg-[#F79E1B]/80 mix-blend-multiply"></div>
								</div>
							</div>
						</div>
						<div class="flex">
							<input class="w-1/2 h-11 bg-transparent px-3 border-r border-outline-variant/30 font-body-md text-body-md text-on-surface focus:outline-none" placeholder="MM / AA" type="text" required/>
							<input class="w-1/2 h-11 bg-transparent px-3 font-body-md text-body-md text-on-surface focus:outline-none" placeholder="CVC" type="text" required/>
						</div>
					</div>
				</div>

				<div class="flex flex-col gap-1 pt-2">
					<label class="font-label-sm text-label-sm text-on-surface-variant uppercase tracking-wider">Nombre en la tarjeta</label>
					<input class="w-full h-11 bg-surface-container-lowest rounded-lg px-3 font-body-md text-body-md text-on-surface focus:outline-none focus:ring-1 focus:ring-primary border border-outline-variant/30" placeholder="Nombre completo" type="text" required/>
				</div>

				<div class="flex items-center gap-2 pt-2">
					<input type="checkbox" id="saveInfo" class="w-4 h-4 rounded text-primary focus:ring-primary bg-surface-container-lowest border-outline-variant/50" />
					<label for="saveInfo" class="font-body-sm text-body-sm text-on-surface-variant">Guardar mi información para la próxima compra más rápido</label>
				</div>

				<button type="submit" class="w-full bg-primary text-on-primary font-headline-sm text-headline-sm h-12 rounded-xl shadow-md flex items-center justify-center gap-2 mt-4 hover:bg-primary-container transition-colors">
					<span class="material-symbols-outlined text-[18px]">lock</span>
					<span>Pagar Pedido</span>
				</button>
			</form>
		</div>

		<div class="flex items-center justify-center gap-2 pt-4 opacity-60">
			<span class="material-symbols-outlined text-[16px]">verified</span>
			<span class="font-label-sm text-label-sm uppercase tracking-wider">Conexión Encriptada de Extremo a Extremo</span>
		</div>
	</div>
</main>
