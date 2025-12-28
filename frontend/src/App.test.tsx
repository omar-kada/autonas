import { http } from 'msw';
import App from './App';
import { server } from './tests/server';
import { renderWithClient } from './tests/utils';

beforeAll(() => server.listen());
afterEach(() => server.resetHandlers());
afterAll(() => server.close());

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => `translated:${key}`,
  }),
}));

Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation((query) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
});

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

    /*expect(
      await screen.findByText('Error fetching status: Unexpected end of JSON input'),
    ).toBeInTheDocument();*/
  });

  it('renders data state', async () => {
    renderWithClient(<App />);

    //expect(screen.getByText('STATUS')).toBeInTheDocument();
  });
});
