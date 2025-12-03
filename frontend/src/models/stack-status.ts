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
  ID: string;
  State: ContainerState;
  Name: string;
  Health: HealthStatus;
};

export interface StackState {
  [key: string]: Array<ContainerStatus>;
}
