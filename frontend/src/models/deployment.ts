type DeploymentStatus = 'scheduled' | 'running' | 'stopped' | 'error' | 'success';

export type Deployment = {
  id: string;
  name: string;
  time: string;
  status: DeploymentStatus;
  diff: string;
};
