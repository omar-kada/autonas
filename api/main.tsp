import "@typespec/http";
import "@typespec/versioning";

using Http;
using Versioning;

@service(#{ title: "AutoNAS API" })
@versioned(Versions)
@route("/api")
@useAuth(BearerAuth /*| OpenIdConnectAuth<"">*/)
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
  _unknown: "unknown",
}

model FileDiff {
  oldFile: string;
  newFile: string;
  diff: string;
}

model Deployment {
  id: string;
  title: string;
  author: string;
  time: utcDateTime;
  endTime: utcDateTime;
  diff: string;
  status: DeploymentStatus;
}

model DeploymentWithDetails is Deployment {
  files: FileDiff[];
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
  health: ContainerHealth;
  startedAt: utcDateTime;
}

model Stats {
  error: int32;
  success: int32;
  nextDeploy: utcDateTime;
  lastDeploy: utcDateTime;
  author: string;
  status: DeploymentStatus;
  health: ContainerHealth;
}

model Settings {
  repo: string;
  branch?: string;
  cron?: string;
  username?: string;
  token?: string;
}

model Config {
  globalVariables: Record<string>;
  services: Record<Record<string>>;
}

model Features {
  displayConfig: boolean;
  editConfig: boolean;
  displayLogs: boolean;
  editSettings: boolean;
}

model Page<T> {
  items: T[];
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

@route("/register")
@tag("Register")
@useAuth(NoAuth)
interface RegisterAPI {
  @get registered(): {
    registered: boolean;
  } | Error;
  @post register(@body user: Credentials): BooleanResponse | Error;
}

model User {
  username: string;
}

model Credentials {
  username: string;
  password: string;
}

model BooleanResponse {
  success: boolean;
}

@route("/login")
@tag("Login")
@useAuth(NoAuth)
interface AuthAPI {
  @post login(@body credentials: Credentials): BooleanResponse | Error;
}

@route("/logout")
@tag("Logout")
interface LogoutAPI {
  @post logout(): BooleanResponse | Error;
}

@route("/user")
@tag("user")
interface UserAPI {
  @get get(): User | Error;
  //@post update(user: User): BooleanResponse | Error;
  @post @route("change-password") changePassword(
    oldPass: string,
    newPass: string,
  ): BooleanResponse | Error;
  @delete delete(): BooleanResponse | Error;
}

@route("/deployment")
@tag("Deployment")
interface DeployementAPI {
  /** List deployments */
  @get list(
    @query limit: int32,
    @query offset?: string,
  ): Page<Deployment> | Error;
  /** Read deployment */
  @get read(@path id: string): DeploymentWithDetails | Error;
  /** Sync and deploy if changes */
  @post sync(): DeploymentWithDetails | void | Error;
}

@route("/settings")
@tag("Settings")
interface SettingsAPI {
  @get get(): Settings | Error;
  @post set(@body settings: Settings): Settings | Error;
}

@route("/config")
@tag("Config")
interface ConfigAPI {
  @get get(): Config | Error;
  @post set(@body config: Config): Config | Error;
}

@route("/features")
@tag("Features")
interface FeaturesAPI {
  @get get(): Features | Error;
}

@route("/diff")
@tag("Diff")
interface DiffAPI {
  @get get(): FileDiff[] | Error;
}

@route("/stats")
@tag("Stats")
interface StatsAPI {
  @get get(@path days: int32): Stats | Error;
}

@route("/status")
@tag("Status")
interface StatusAPI {
  @get get(): StackStatus[] | Error;
}
