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