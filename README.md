# Notify

Notify pushes msgs to an HTTP enpoint in batches and within a specified interval.

Provides `pkg/http_pusher`, which can be used separately.

## Usage

```
go install github.com/go-task/task/v3/cmd/task@latest
task build
```

```
Usage of ./cmd/notify/notify:
  -batch-amount int
        Max number of notifications per a single request (default 50)
  -help
        Show the help screen
  -interval int
        Time between buffer flushes, seconds (default 5)
  -timeout int
        HTTP timeout, seconds (default 30)
  -url string
        Destination URL (default "http://localhost:5050")
```

## Taskfile

```
task: Available tasks for this project:
* build:
* start:
* start-echo-server:
* test:
* test-echo:
```

## Sample output (no errors)

Using 15 messages from `test/msgs` and following commands:
- `task start-echo-server`
- `task test-echo`

### Notify

```
task: [test-echo] go run cmd/notify/notify.go --batch-amount 2 --interval 1 < test/msgs
2022/02/28 13:08:35 Push
2022/02/28 13:08:35 Push
2022/02/28 13:08:35 Push
2022/02/28 13:08:35 Push
2022/02/28 13:08:35 Push
2022/02/28 13:08:35 Push
2022/02/28 13:08:35 Push
2022/02/28 13:08:35 Push
2022/02/28 13:08:35 Push
2022/02/28 13:08:35 Push
2022/02/28 13:08:35 Push
2022/02/28 13:08:35 Push
2022/02/28 13:08:35 Push
2022/02/28 13:08:35 Push
2022/02/28 13:08:35 Push
2022/02/28 13:08:36 Flush
2022/02/28 13:08:37 Flush
2022/02/28 13:08:38 Flush
2022/02/28 13:08:39 Flush
2022/02/28 13:08:40 Flush
2022/02/28 13:08:41 Flush
2022/02/28 13:08:42 Flush
2022/02/28 13:08:43 Flush
```

### Echo server

```
2022/02/28 13:08:36 Received 2 msgs
2022/02/28 13:08:37 Received 2 msgs
2022/02/28 13:08:38 Received 2 msgs
2022/02/28 13:08:39 Received 2 msgs
2022/02/28 13:08:40 Received 2 msgs
2022/02/28 13:08:41 Received 2 msgs
2022/02/28 13:08:42 Received 2 msgs
2022/02/28 13:08:43 Received 1 msgs
```

## Sample output (with errors)

Same as above, but with random errors generated (1 out of 3 reqs will fail).

### Notify

```
task: [test-echo] go run cmd/notify/notify.go --batch-amount 2 --interval 1 < test/msgs
2022/02/28 13:43:41 Push
2022/02/28 13:43:41 Push
2022/02/28 13:43:41 Push
2022/02/28 13:43:41 Push
2022/02/28 13:43:41 Push
2022/02/28 13:43:41 Push
2022/02/28 13:43:41 Push
2022/02/28 13:43:41 Push
2022/02/28 13:43:41 Push
2022/02/28 13:43:41 Push
2022/02/28 13:43:41 Push
2022/02/28 13:43:41 Push
2022/02/28 13:43:41 Push
2022/02/28 13:43:41 Push
2022/02/28 13:43:41 Push
2022/02/28 13:43:42 Flush
2022/02/28 13:43:43 Flush
2022/02/28 13:43:44 Flush
2022/02/28 13:43:45 Flush
2022/02/28 13:43:45 HTTP error 500
2022/02/28 13:43:45 2 msgs with errors, re-sending
2022/02/28 13:43:45 PushMany
2022/02/28 13:43:46 Flush
2022/02/28 13:43:47 Flush
2022/02/28 13:43:48 Flush
2022/02/28 13:43:49 Flush
2022/02/28 13:43:50 Flush
2022/02/28 13:43:50 HTTP error 500
2022/02/28 13:43:50 1 msgs with errors, re-sending
2022/02/28 13:43:50 PushMany
2022/02/28 13:43:51 Flush
```

### Echo server 

```
2022/02/28 13:43:42 Received 2 msgs
2022/02/28 13:43:43 Received 2 msgs
2022/02/28 13:43:44 Received 2 msgs
2022/02/28 13:43:45 Failing...
2022/02/28 13:43:46 Received 2 msgs
2022/02/28 13:43:47 Received 2 msgs
2022/02/28 13:43:48 Received 2 msgs
2022/02/28 13:43:49 Received 2 msgs
2022/02/28 13:43:50 Failing...
2022/02/28 13:43:51 Received 1 msgs
```