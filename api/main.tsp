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
  status: "planned" | "running" | "success" | "error";
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

@error
model Error {
  code: int32;
  message: string;
}

model AnalyzeResult {
  id: string;
  analysis: string;
}

@route("/deployment")
@tag("Deployment")
interface DeployementAPI {
  /** List widgets */
  @get list(): Deployment[] | Error;
  /** Read widgets */
  @get read(@path id: string): Deployment | Error;
}

@route("/status")
@tag("Status")
interface StatusAPI {
  /** List widgets */
  @get get(): StackStatus[] | Error;
}
