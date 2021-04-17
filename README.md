# kami-go

## Features
### Dispatch Job
Allows users to create and execute processes.

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
  "exitCode": "<exitCode>"
}
```

### Query Job Logs
Display the logs of an existing job.

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
