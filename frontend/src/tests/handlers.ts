import { http } from 'msw';

export const handlers = [
  http.get('/api/status', () => {
    return new Response(JSON.stringify({ status: 'active' }), {
      status: 200,
    });
  }),
];
