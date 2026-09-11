import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import LoginPage from './+page.svelte';

// Mock the fetch API globally
global.fetch = vi.fn();

describe('Login Page', () => {
	it('renders login form correctly', () => {
		render(LoginPage);
		
		expect(screen.getByText('Iniciar Sesión')).toBeInTheDocument();
		expect(screen.getByLabelText('Usuario')).toBeInTheDocument();
		expect(screen.getByLabelText('Contraseña')).toBeInTheDocument();
		expect(screen.getByRole('button', { name: 'Entrar' })).toBeInTheDocument();
	});

	it('shows error message if fields are empty on submit', async () => {
		render(LoginPage);
		
		const button = screen.getByRole('button', { name: 'Entrar' });
		await fireEvent.click(button);
		
		await waitFor(() => {
			expect(screen.getByText('Por favor ingresa usuario y contraseña.')).toBeInTheDocument();
		});
	});

	it('calls fetch on valid submission', async () => {
		render(LoginPage);
		
		const usernameInput = screen.getByLabelText('Usuario');
		const passwordInput = screen.getByLabelText('Contraseña');
		const button = screen.getByRole('button', { name: 'Entrar' });

		await fireEvent.input(usernameInput, { target: { value: 'admin' } });
		await fireEvent.input(passwordInput, { target: { value: '1234' } });
		
		// Mock a successful login response
		(global.fetch as any).mockResolvedValueOnce({
			ok: true,
			json: async () => ({ message: 'Logged in successfully' })
		});

		await fireEvent.click(button);
		
		expect(global.fetch).toHaveBeenCalledWith('http://localhost:8080/api/login', expect.objectContaining({
			method: 'POST',
			body: JSON.stringify({ username: 'admin', password: '1234' })
		}));
	});
});
