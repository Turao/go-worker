# kami-go

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
  "command": "<your-command-here>"
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
This is a blocking operation.

```
GET /job/{id}/info
```

Returns:
```json
{
  "id": "<uuid>",
  "status": "<status>",
  "exitCode": "<exitCode>"
}
```

### Query Job Logs
Display the logs of an existing job.
This is a blocking operation.

```
GET /job/{id}/logs
```

Returns:
```json
{
  "output": "<stdout>",
  "errors": "<stderr>"
}
```


### Architecture

- apiserver
  - server.go: decorates an `http.server` with `worker` handlers
  - transport.go: decodes http requests, encodes http responses
  - endpoint.go: implements an anti-fragile layer, independent of transport (just like a controller)
  - logging.go: decorates services with custom logging (what is being called and when)
  - service.go: provides application-level features 
- worker
  - worker.go: provides appliciation-level features
- job
  - job.go: provides core business features such as start, stop, info, ...
    - (this is the main domain entity)
- storage
  - store.go: in-memory repository