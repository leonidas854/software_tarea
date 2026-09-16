/// <reference types="vitest/globals" />
import '@testing-library/jest-dom';
import { vi } from 'vitest';

vi.mock('$app/navigation', () => ({
	goto: vi.fn()
}));
