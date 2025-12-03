import type { StackStatus } from '@/api/api';
import { http } from 'msw';

const mockState: StackStatus = {
  name: 'homepage',
  stackId: 'stack-1',
  services: [
    {
      containerId: '1',
      state: 'running',
      name: 'web-container',
      health: 'healthy',
      createdAt: new Date() + '',
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
