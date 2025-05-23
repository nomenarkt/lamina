import { render, screen, fireEvent } from '@testing-library/react';
import { AuthCardLayout } from '../AuthCardLayout';

describe('AuthCardLayout', () => {
  it('renders title, children and footer', () => {
    render(
      <AuthCardLayout
        title="Welcome"
        onSubmit={() => {}}
        footer={<div>Footer Content</div>}
      >
        <input placeholder="test input" />
      </AuthCardLayout>
    );

    expect(screen.getByRole('heading', { name: /welcome/i })).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/test input/i)).toBeInTheDocument();
    expect(screen.getByText(/footer content/i)).toBeInTheDocument();
  });

  it('calls onSubmit handler', () => {
    const handleSubmit = jest.fn((e) => e.preventDefault());
    render(
      <AuthCardLayout
        title="Test"
        onSubmit={handleSubmit}
        footer={<div>Footer</div>}
      >
        <button type="submit">Submit</button>
      </AuthCardLayout>
    );

    fireEvent.click(screen.getByRole('button', { name: /submit/i }));
    expect(handleSubmit).toHaveBeenCalled();
  });
});
