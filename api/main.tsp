import "@typespec/http";
import "@typespec/versioning";

using Http;
using Versioning;

@service(#{ title: "AutoNAS API" })
@versioned(Versions)
@route("/api")
namespace DeploymentService;
enum Versions {
  v1: "1.0",
}

model Event {
  time: utcDateTime;
  msg: string;
  level: "ERROR" | "INFO" | "DEBUG" | "WARN";
}

enum DeploymentStatus {
  planned: "planned",
  running: "running",
  success: "success",
  error: "error",
}

enum ContainerHealth {
  healthy: "healthy",
  unhealthy: "unhealthy",
  starting: "starting",
  none: "none",
}

model FileDiff {
  oldFile: string;
  newFile: string;
  diff: string;
}

model Deployment {
  id: string;
  title: string; // Commit message
  author: string; // Commit message
  time: utcDateTime;
  endTime: utcDateTime;
  diff: string;
  files: FileDiff[];
  status: DeploymentStatus;
  events: Event[];
}

model StackStatus {
  stackId: string;
  name: string;
  services: Array<ContainerStatus>;
}

model ContainerStatus {
  containerId: string;
  state:
    | "created"
    | "running"
    | "paused"
    | "restarting"
    | "removing"
    | "exited"
    | "dead";
  name: string;
  health: "healthy" | "unhealthy" | "starting" | "none";
  startedAt: utcDateTime;
}

model Stats {
  error: int32;
  success: int32;
  nextDeploy: utcDateTime;
  lastDeploy: utcDateTime;
  status: DeploymentStatus;
  health: ContainerHealth;
  author: string;
}

model Page<T> {
  data: T[];
  pageInfo: PageInfo;
}

model PageInfo {
  hasNextPage: boolean;
  endCursor: string;
}

@error
model Error {
  code: int32;
  message: string;
}

@route("/deployment")
@tag("Deployment")
interface DeployementAPI {
  /** List deployments */
  @get list(
    @query first: int32,
    @query after?: string,
  ): Page<Deployment> | Error;
  /** Read deployments */
  @get read(@path id: string): Deployment | Error;
  /** Sync and deploy if changes */
  @post sync(): Deployment | Error;
}

@route("/diff")
@tag("Diff")
interface DiffAPI {
  /** Get diff with repo */
  @get get(): FileDiff[] | Error;
}

@route("/stats")
@tag("Stats")
interface StatsAPI {
  /** Get deployment statistics */
  @get get(@path days: int32): Stats | Error;
}

@route("/status")
@tag("Status")
interface StatusAPI {
  /** Get current status */
  @get get(): StackStatus[] | Error;
}
