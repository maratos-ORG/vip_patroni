# Development Guide

Welcome to vip_patroni app development guide!

## Dependencies

You have to install these tools:

- SCM: `git`
- Programming language: `golang 1.16`
- Build tools (depend on build method):
  - `make` - to use Makefile
  - `golangci-lint`- to use make linter
     <!-- INSTALL golangci-lint
     curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.43.0 
     -->

## Build and run

1. Clone repo:
  ```bash
  git clone git@github.com:maratos-ORG/vip_patroni.git
  cd vip_patroni
  ```

2. Make build and run app
  ```bash
  make build RELEASE=1.0 && ./bin/vip_patroni --version
  ```

Other options:

```bash
make help
```
