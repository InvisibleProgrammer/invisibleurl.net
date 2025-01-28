# invisibleurl.net

## Local dev:

If you follow the following steps, your instance will run on [localhost](https://localhost:3000)

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

### Add https on your local dev setup

Generate certificate:
```bash
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 36500 -nodes
```
Note: use localhost as fqdn
Note: your browser won't like that certificate. You have to add it to your browser's certificate store manually.

### Use air to run it locally

Install:
```bash
go install github.com/cosmtrek/air@latest
```

Run:
```bash
air .
```
