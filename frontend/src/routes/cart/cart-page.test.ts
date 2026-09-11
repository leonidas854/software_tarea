import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import CartPage from './+page.svelte';

global.fetch = vi.fn();

describe('Cart Page', () => {
	it('renders empty cart correctly', async () => {
		(global.fetch as any).mockResolvedValueOnce({
			ok: true,
			json: async () => ({ items: [], subtotal: 0 })
		});

		render(CartPage);
		
		expect(screen.getByText('Agregar Producto')).toBeInTheDocument();
		expect(screen.getByText('Tu Carrito')).toBeInTheDocument();
		
		// Wait for fetch to complete and empty cart text to show
		await waitFor(() => {
			expect(screen.getByText('El carrito está vacío.')).toBeInTheDocument();
		});
	});

	it('shows client-side validation error for invalid quantity', async () => {
		(global.fetch as any).mockResolvedValueOnce({
			ok: true,
			json: async () => ({ items: [], subtotal: 0 })
		});
		
		render(CartPage);

		const idInput = screen.getByLabelText('ID del Producto');
		const nameInput = screen.getByLabelText('Nombre');
		const priceInput = screen.getByLabelText('Precio ($)');
		const quantityInput = screen.getByLabelText('Cantidad');
		
		await fireEvent.input(idInput, { target: { value: 'P01' } });
		await fireEvent.input(nameInput, { target: { value: 'Test' } });
		await fireEvent.input(priceInput, { target: { value: '10' } });
		await fireEvent.input(quantityInput, { target: { value: '' } });
		
		const addButton = screen.getByRole('button', { name: 'Agregar al Carrito' });
		await fireEvent.click(addButton);
		
		await waitFor(() => {
			expect(screen.getByText('La cantidad debe ser un número entero mayor a 0.')).toBeInTheDocument();
		});
	});

	it('adds product correctly and refreshes cart', async () => {
		// Mock initial cart fetch
		(global.fetch as any).mockResolvedValueOnce({
			ok: true,
			json: async () => ({ items: [], subtotal: 0 })
		});
		
		render(CartPage);
		
		const idInput = screen.getByLabelText('ID del Producto');
		const nameInput = screen.getByLabelText('Nombre');
		const priceInput = screen.getByLabelText('Precio ($)');
		const quantityInput = screen.getByLabelText('Cantidad');
		
		await fireEvent.input(idInput, { target: { value: 'P01' } });
		await fireEvent.input(nameInput, { target: { value: 'Test Product' } });
		await fireEvent.input(priceInput, { target: { value: '10' } });
		await fireEvent.input(quantityInput, { target: { value: '2' } });

		// Mock the add to cart POST response
		(global.fetch as any).mockResolvedValueOnce({
			ok: true,
			text: async () => 'Success'
		});
		
		// Mock the subsequent cart refresh GET response
		(global.fetch as any).mockResolvedValueOnce({
			ok: true,
			json: async () => ({
				items: [{ id: 'P01', name: 'Test Product', price: 10.0, quantity: 2 }],
				subtotal: 20.0
			})
		});
		
		const addButton = screen.getByRole('button', { name: 'Agregar al Carrito' });
		
		await fireEvent.click(addButton);
		
		// The fetchCart function is called after the POST succeeds, we must wait for it to render
		await waitFor(() => {
			expect(screen.getByText('Test Product')).toBeInTheDocument();
			expect(screen.getByText('Total: $20.00')).toBeInTheDocument();
		}, { timeout: 2000 });
	});
});
