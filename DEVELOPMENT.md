# Development Guide

Welcome to vip_patroni app development guide!

## Dependencies

Make sure you have the following tools installed:

| Tool           | Purpose                             |
|----------------|-------------------------------------|
| git            | Version control (SCM)               |
| golang >= 1.16 | Main programming language           |
| make           | Build automation                    |
| golangci-lint  | Code linting (used in Makefile)     |

- SCM: `git`
- Programming language: `golang 1.16`
- Build tools (depend on build method):
  - `make` - to use Makefile
  - `golangci-lint`- to use make linter  
    ```bash 
    💡 Install golangci-lint
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.43.0
    ```
## Build & Run

1.	Clone the repository:
```bash
git clone git@github.com:maratos-ORG/vip_patroni.git
cd vip_patroni
```

2. Build and run the app:
```bash
make build RELEASE=1.0 && ./bin/vip_patroni --version
```

3. More Options
List available make targets:
```bash
make help
```
