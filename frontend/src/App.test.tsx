import { screen } from '@testing-library/react';
import { renderWithClient } from './tests/utils';
import App from './App';
import { server } from './tests/server';
import { http } from 'msw';

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

    expect(await screen.findByText('count is 0')).toBeInTheDocument();
    expect(screen.getByText('TITLE')).toBeInTheDocument();
  });
});
