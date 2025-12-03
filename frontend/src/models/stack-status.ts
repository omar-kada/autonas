export type ContainerState =
  | 'created'
  | 'running'
  | 'paused'
  | 'restarting'
  | 'removing'
  | 'exited'
  | 'dead';

export type HealthStatus = 'healthy' | 'unhealthy' | 'starting' | 'none';

export type ContainerStatus = {
  containerId: string;
  state: ContainerState;
  name: string;
  health: HealthStatus;
  createdAt: Date;
};

export interface StackState {
  [key: string]: Array<ContainerStatus>;
}
