import { cleanup, render, screen, fireEvent, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { api } from '../src/lib/api';
import App from '../src/App';

vi.mock('../src/lib/api', () => ({
  api: {
    register: vi.fn(),
    farms: vi.fn(),
    crops: vi.fn(),
  },
}));

const registeredFarmer = {
  id: 'farmer-1',
  full_name: 'Test Farmer',
  role: 'farmer',
  preferred_language: 'en',
};

function openRegistrationForm() {
  render(<App />);
  fireEvent.click(
    screen.getByRole('button', { name: /create farmer profile/i }),
  );
}

function fillRegistrationForm(values: {
  name?: string;
  email?: string;
  phone?: string;
  password?: string;
}) {
  if (values.name !== undefined) {
    fireEvent.change(screen.getByLabelText(/full name/i), {
      target: { value: values.name },
    });
  }
  if (values.email !== undefined) {
    fireEvent.change(screen.getByLabelText(/email address/i), {
      target: { value: values.email },
    });
  }
  if (values.phone !== undefined) {
    fireEvent.change(screen.getByLabelText(/phone number/i), {
      target: { value: values.phone },
    });
  }
  if (values.password !== undefined) {
    fireEvent.change(screen.getByLabelText(/password/i), {
      target: { value: values.password },
    });
  }
}

function submitRegistration() {
  fireEvent.click(screen.getByRole('button', { name: /create profile/i }));
}

describe('Farmer registration', () => {
  afterEach(() => {
    cleanup();
  });

  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();

    vi.mocked(api.register).mockResolvedValue({
      success: true,
      data: registeredFarmer,
    });
    vi.mocked(api.farms).mockResolvedValue({
      success: true,
      data: { data: [], total: 0 },
    });
    vi.mocked(api.crops).mockResolvedValue({
      success: true,
      data: { data: [], total: 0 },
    });
  });

  it('does not submit registration when both email and phone are missing', () => {
    openRegistrationForm();
    fillRegistrationForm({
      name: 'Test Farmer',
      password: 'password123',
    });

    submitRegistration();

    expect(api.register).not.toHaveBeenCalled();
    expect(screen.getByRole('alert')).toHaveTextContent(
      /email address or a phone number/i,
    );
  });

  it('submits registration when an email address is provided', async () => {
    openRegistrationForm();
    fillRegistrationForm({
      name: 'Test Farmer',
      email: 'farmer@example.com',
      password: 'password123',
    });

    submitRegistration();

    await waitFor(() => expect(api.register).toHaveBeenCalledTimes(1));

    const payload = vi.mocked(api.register).mock.calls[0][0];
    expect(payload).toMatchObject({
      full_name: 'Test Farmer',
      email: 'farmer@example.com',
      password: 'password123',
      preferred_language: 'en',
    });
    expect(payload.phone).toBeUndefined();
  });

  it('submits registration when a phone number is provided', async () => {
    openRegistrationForm();
    fillRegistrationForm({
      name: 'Test Farmer',
      phone: '+254712345678',
      password: 'password123',
    });

    submitRegistration();

    await waitFor(() => expect(api.register).toHaveBeenCalledTimes(1));

    const payload = vi.mocked(api.register).mock.calls[0][0];
    expect(payload).toMatchObject({
      full_name: 'Test Farmer',
      phone: '+254712345678',
      password: 'password123',
      preferred_language: 'en',
    });
    expect(payload.email).toBeUndefined();
  });

  it('shows a validation message for a password shorter than eight characters', () => {
    openRegistrationForm();
    fillRegistrationForm({
      name: 'Test Farmer',
      email: 'farmer@example.com',
      password: 'short',
    });

    const form = screen
      .getByRole('button', { name: /create profile/i })
      .closest('form');
    expect(form).not.toBeNull();

    // Submit directly so the assertion covers client-side validation rather
    // than the browser's native minLength constraint.
    fireEvent.submit(form!);

    expect(api.register).not.toHaveBeenCalled();
    expect(screen.getByRole('alert')).toHaveTextContent(
      /at least 8 characters/i,
    );
  });

  it('shows a validation message when the full name is missing', () => {
    openRegistrationForm();
    fillRegistrationForm({
      email: 'farmer@example.com',
      password: 'password123',
    });

    const form = screen
      .getByRole('button', { name: /create profile/i })
      .closest('form');
    expect(form).not.toBeNull();

    fireEvent.submit(form!);

    expect(api.register).not.toHaveBeenCalled();
    expect(screen.getByRole('alert')).toHaveTextContent(/full name/i);
  });
});
