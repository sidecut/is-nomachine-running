# is-nomachine-running

An extremely basic web server that simply displays whether nomachine is running and, if so, whether a client is connected.

## Caddy reverse proxy

This repository contains a caddyfile for use by the [Caddy](https://caddyserver.com) [reverse proxy](https://en.wikipedia.org/wiki/Reverse_proxy) to serve the site via SSL.

The cert file and cert key in pem format should be named `certs/certfile.crt` and `certs/certfile.key`.

## Project setup

```bash
yarn install
```

### Compiles and hot-reloads for development

```bash
yarn run serve
```

### Compiles and minifies for production

```bash
yarn run build
```

### Run your tests

```bash
yarn run test
```

### Lints and fixes files

```bash
yarn run lint
```

### Customize configuration

See [Configuration Reference](https://cli.vuejs.org/config/).

## TO DO

- [x] Test on macOS
- [ ] Test on Linux
- [ ] Test on Windows
- [ ] Support virtual clients, multiple clients, etc.
- [ ] Upgrade to Vue 3
