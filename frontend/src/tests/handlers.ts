import type { StackState } from '@/models/stack-status';
import { http } from 'msw';

const mockState: StackState = {
  homepage: [
    {
      ID: '1',
      State: 'running',
      Name: 'web-container',
      Health: 'healthy',
    },
  ],
};
export const handlers = [
  http.get('/api/status', () => {
    return new Response(JSON.stringify(mockState), {
      status: 200,
    });
  }),
];
