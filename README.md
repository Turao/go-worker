# go-worker

### Features

**Dispatch Job**
```yaml
Endpoint: POST /job
Description: |-
  Allows users to create and execute a unix process.
  This is a not a blocking operation (i.e. does not wait for the process to start).
 
Request:
{
  "name": "<your-command-here>",
  "args": [
    "<your-commands-first-argument>",
    "<your-commands-second-argument>",
    "<...>"
  ]
}

Response:
{
  "id": "<uuid>"
}
```

**Stop Job**
```yaml
Endpoint: POST /job/{id}/stop
Description: |-
  Allows users to stop a running job.
  This is a not a blocking operation (i.e. does not wait for the process to stop).

Response:
  - HTTP Status 200 (Sucessful)
  - HTTP Status 404 (Not Found)
  - HTTP Status 500 (Internal Error)
```

**Query Job Info:**
```yaml
Endpoint: `GET /job/{id}/info`
Description: Displays information about an existing job.
Response:
{
  "id": "<uuid>",
  "status": "<status>",
  "exitCode": "<exitCode>",
  "output": "<stdout>",
  "errors": "<stderr>"
}
```

### Architecture

![architecture](docs/architecture.drawio.svg)

### Folder Structure
```
- api/
  - v1/
    - worker.go: worker DTOs
    - job.go: job DTOs
- cmd/
  - client
    - main.go: decorates a client with a CLI interface
  - server
    - main.go: decorates a server (no CLI interface for now...)
- pkg/
  - client
    - client.go: provides application services. works like a proxy
    - endpoint.go: implements an anti-fragile layer, independent of transport
    - transport.go: encodes http requests, sends http requests, decodes http responses
  - server
    - server.go: decorates an `http.server` with routing (should I move routing to transport?)
    - transport.go: decodes http requests, encodes http responses
    - endpoint.go: implements an anti-fragile layer, independent of transport (somewhat like a controller)
    - logging.go: decorates services with custom logging (what is being called and when (I'm just playing with the middleware pattern here)
    - service.go: provides application-level features
  - worker
    - worker.go: manages jobs
  - job
    - job.go: manages a unix process
  - storage
    - store.go: in-memory repository
```
