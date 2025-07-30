# invisibleurl.net

## Local dev:

If you follow the following steps, your instance will run on [localhost](https://localhost:3000)

### To be able to send emails locally, use mailhog

Ask for SMTP settings and fill your .env file

OR you can use mailhog for testing purposes, like this:

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

### SSL section

At localhost, you don't need a certificate.
On the production environment, there will be NGinX that handles SSL.
Be careful, there is a chance that you will have to fight with your browser on local dev setup to be able to reach the site without SSL.
Or, if you want, you can still set up NGinX locally as well and generate a certificate for yourself:

Generate certificate:
```bash
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 36500 -nodes
```
Note: use localhost as fqdn
Note: your browser won't like that certificate. You have to add it to your browser's certificate store manually.


### Use air to run it locally

Install:
```bash
go install github.com/air-verse/air@latest
```

Run:
```bash
air .
```
