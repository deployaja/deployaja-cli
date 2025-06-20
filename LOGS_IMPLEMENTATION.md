# Logs Follow Implementation Guide

This document explains how to implement the follow functionality for logs in the DeployAja CLI and what changes are needed in the backend API.

## CLI Changes Made

### 1. Updated API Client (`internal/api/client.go`)

- **Modified `GetLogs` method**: Now only handles non-follow mode and returns an error if follow is requested
- **Added `GetLogsStream` method**: New method that handles streaming logs using Server-Sent Events (SSE)

### 2. Updated Logs Command (`cmd/logs.go`)

- **Split functionality**: Regular logs vs streaming logs
- **Added `streamLogs` function**: Handles real-time log streaming with graceful shutdown
- **Signal handling**: Supports Ctrl+C to stop following logs

## Backend API Changes Required

### 1. New Streaming Endpoint

The backend needs a new endpoint for streaming logs:

```
GET /api/v1/logs/{name}/stream?tail={number}
```

**Headers required:**
- `Accept: text/event-stream`
- `Cache-Control: no-cache`
- `Connection: keep-alive`

### 2. Server-Sent Events (SSE) Format

The streaming endpoint should return data in SSE format:

```
data: {"timestamp":"2024-01-01T12:00:00Z","level":"INFO","message":"Log message","source":"app"}

data: {"timestamp":"2024-01-01T12:00:01Z","level":"ERROR","message":"Error message","source":"app"}

: keepalive

data: [DONE]
```

### 3. Backend Implementation Requirements

#### Regular Logs Endpoint (existing)
```
GET /api/v1/logs/{name}?tail={number}
```
Returns JSON response:
```json
{
  "logs": [
    {
      "timestamp": "2024-01-01T12:00:00Z",
      "level": "INFO",
      "message": "Log message",
      "source": "app"
    }
  ]
}
```

#### Streaming Logs Endpoint (new)
```
GET /api/v1/logs/{name}/stream?tail={number}
```

**Implementation steps:**

1. **Set SSE headers:**
   ```go
   w.Header().Set("Content-Type", "text/event-stream")
   w.Header().Set("Cache-Control", "no-cache")
   w.Header().Set("Connection", "keep-alive")
   ```

2. **Send initial logs:**
   - Fetch the last `tail` number of logs
   - Send each log as an SSE event

3. **Stream new logs:**
   - Subscribe to new log events for the deployment
   - Send each new log as it arrives
   - Handle client disconnection gracefully

4. **Keepalive:**
   - Send `: keepalive` every 30 seconds to prevent connection timeout

5. **End stream:**
   - Send `data: [DONE]` when streaming should end

### 4. Log Storage Considerations

For real-time log streaming, consider:

- **Database triggers**: Use database triggers to notify when new logs are inserted
- **Message queues**: Use Redis pub/sub, RabbitMQ, or similar for real-time notifications
- **File watching**: Watch log files for changes (if using file-based logging)
- **WebSocket fallback**: Consider WebSocket as an alternative to SSE

### 5. Example Backend Implementation (Go)

```go
func handleStreamLogs(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    name := vars["name"]
    
    tailStr := r.URL.Query().Get("tail")
    tail := 100 // default
    if tailStr != "" {
        if t, err := strconv.Atoi(tailStr); err == nil {
            tail = t
        }
    }

    // Set SSE headers
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

    // Send initial logs
    initialLogs := getLogsFromStorage(name, tail)
    for _, logEntry := range initialLogs {
        data, _ := json.Marshal(logEntry)
        fmt.Fprintf(w, "data: %s\n\n", data)
        w.(http.Flusher).Flush()
    }

    // Subscribe to new logs
    logChan := subscribeToNewLogs(name)
    notify := w.(http.CloseNotifier).CloseNotify()
    
    for {
        select {
        case logEntry := <-logChan:
            data, _ := json.Marshal(logEntry)
            fmt.Fprintf(w, "data: %s\n\n", data)
            w.(http.Flusher).Flush()
            
        case <-notify:
            // Client disconnected
            unsubscribeFromLogs(name)
            return
            
        case <-time.After(30 * time.Second):
            // Send keepalive
            fmt.Fprintf(w, ": keepalive\n\n")
            w.(http.Flusher).Flush()
        }
    }
}
```

## Usage

### CLI Commands

```bash
# Regular logs (non-follow)
deployaja logs my-app --tail 50

# Follow logs in real-time
deployaja logs my-app --follow

# Follow logs with custom tail
deployaja logs my-app --follow --tail 200
```

### Features

- **Real-time streaming**: See logs as they happen
- **Graceful shutdown**: Press Ctrl+C to stop following
- **Initial logs**: Shows the last N logs before starting to follow
- **Error handling**: Proper error handling for network issues
- **Keepalive**: Prevents connection timeouts

## Testing

To test the implementation:

1. Start your backend server with the new streaming endpoint
2. Run `deployaja logs my-app --follow`
3. Generate some logs in your application
4. Verify that new logs appear in real-time
5. Press Ctrl+C to test graceful shutdown 