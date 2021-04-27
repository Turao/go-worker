# go-worker

## Features
### Dispatch a new Job
Allows users to create and execute processes.

:memo: This is a not a blocking operation (i.e. does not wait for the process to start).

```
POST /job
```

Receives:
```json
{
  "name": "<your-command-here>",
  "args": [
    "<your-commands-first-argument>",
    "<your-commands-second-argument>",
    "<...>"
  ]
}
```

Returns:
```json
{
  "id": "<uuid>"
}
```

### Forcefully Stop a Job
Allows users to stop a running job.

:memo: This is a not a blocking operation (i.e. does not wait for the process to stop).

```
POST /job/{id}/stop
```

Returns:
- `HTTP Status 200 (Sucessful)`
- `HTTP Status 404 (Not Found)`
- `HTTP Status 500 (Internal Error)`

### Query Job Info
Displays information about an existing job.

```
GET /job/{id}/info
```

Returns:
```json
{
  "id": "<uuid>",
  "status": "<status>",
  "exitCode": "<exitCode>",
  "output": "<stdout>",
  "errors": "<stderr>"
}
```

### Architecture

> I'm not really sure how to best structure the client layers yet...

- api/
  - v1/
    - worker.go: worker endpoints (client/server) API
- cmd/
  - client
    - main.go: decorates a client with a nice CLI interface
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
    - endpoint.go: implements an anti-fragile layer, independent of transport (just like a controller)
    - logging.go: decorates services with custom logging (what is being called and when)
    - service.go: provides application-level features
  - worker
    - worker.go: manages jobs
  - job
    - job.go: manages a unix process
  - storage
    - store.go: in-memory repository