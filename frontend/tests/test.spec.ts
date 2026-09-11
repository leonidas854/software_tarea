import { test, expect } from '@playwright/test';

test.describe('Shopping Cart Flow', () => {
	test('user can login and add items to cart', async ({ page }) => {
		// Go to login page
		await page.goto('/');

		// Check title
		await expect(page.getByRole('heading', { name: 'Iniciar Sesión' })).toBeVisible();

		// Fill credentials
		await page.getByLabel('Usuario').fill('admin');
		await page.getByLabel('Contraseña').fill('1234');
		await page.getByRole('button', { name: 'Entrar' }).click();

		// Check redirection to cart
		await expect(page).toHaveURL(/.*\/cart/);
		await expect(page.getByRole('heading', { name: 'Tu Carrito' })).toBeVisible();

		// Add product
		await page.getByLabel('ID del Producto').fill('TEST1');
		await page.getByLabel('Nombre').fill('Playwright Test Item');
		await page.getByLabel('Precio ($)').fill('50.50');
		await page.getByLabel('Cantidad').fill('2');
		await page.getByRole('button', { name: 'Agregar al Carrito' }).click();

		// Check if product is in cart
		await expect(page.getByText('Playwright Test Item')).toBeVisible();
		await expect(page.getByText('Total: $101.00')).toBeVisible();

		// Logout
		await page.getByRole('button', { name: 'Cerrar Sesión' }).click();
		await expect(page).toHaveURL('http://localhost:5173/');
	});
});
