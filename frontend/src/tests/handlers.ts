import { ContainerHealth, DeploymentStatus, type StackStatus, type Stats } from '@/api/api';
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
      startedAt: `${new Date()}`,
    },
  ],
};

const mockStats: Stats = {
  error: 5,
  success: 15,
  nextDeploy: new Date().toString(),
  lastDeploy: new Date().toString(),
  status: DeploymentStatus.success,
  health: ContainerHealth.healthy,
  author: 'Test',
};

export const handlers = [
  http.get('/api/status', () => {
    return new Response(JSON.stringify(mockState), {
      status: 200,
    });
  }),
  http.get('/api/stats/:days', () => {
    return new Response(JSON.stringify(mockStats), {
      status: 200,
    });
  }),
];
