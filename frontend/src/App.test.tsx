import { screen } from '@testing-library/react';
import { http } from 'msw';
import App from './App';
import { server } from './tests/server';
import { renderWithClient } from './tests/utils';

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

describe('App', () => {
  it('renders error state', async () => {
    server.use(
      http.get('/api/status', () => {
        return new Response('', {
          status: 500,
        });
      }),
    );

    renderWithClient(<App />);

    expect(
      await screen.findByText('Error fetching status: Unexpected end of JSON input'),
    ).toBeInTheDocument();
  });

  it('renders data state', async () => {
    renderWithClient(<App />);

    //expect(screen.getByText('STATUS')).toBeInTheDocument();
  });
});
