# invisibleurl.net

## Local dev:

### Use air to run it locally

Install:
```bash
go install github.com/cosmtrek/air@latest
```

Run:
```bash
air .
```

### To be able to send emails locally, use mailhog

Install:
```bash
brew install mailhog
```

Run:
```bash
brew services run mailhog
```
**Note**: Web UI: http://localhost:8025/

### To run the tests:

```bash
go test ./...
```
