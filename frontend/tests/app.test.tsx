import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import App from '../src/App';

describe('SoilCycle application shell', () => {
  it('renders the public onboarding and product navigation', () => {
    render(<App />);
    expect(screen.getByText(/Make every season/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /Create farmer profile/i })).toBeInTheDocument();
    expect(screen.getByText(/SoilCycle by Earth2Farm/i)).toBeInTheDocument();
  });
});
