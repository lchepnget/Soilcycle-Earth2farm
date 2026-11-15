import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { api } from '../src/lib/api';
import App from '../src/App';

vi.mock('../src/lib/api', () => ({
  api: {
    register: vi.fn(),
  },
}));

describe('Farmer registration', () => {
  it('should not submit registration when both email and phone are missing', () => {
    render(<App />);

    // Open the registration form.
    fireEvent.click(
      screen.getByRole('button', { name: /create farmer profile/i }),
    );

    // Provide the other required fields.
    fireEvent.change(screen.getByLabelText(/full name/i), {
      target: { value: 'Test Farmer' },
    });

    fireEvent.change(screen.getByLabelText(/password/i), {
      target: { value: 'password123' },
    });

    // Attempt to submit without email or phone.
    fireEvent.click(
      screen.getByRole('button', { name: /create profile/i }),
    );

    // Registration should not reach the backend.
    expect(api.register).not.toHaveBeenCalled();
  });
});
