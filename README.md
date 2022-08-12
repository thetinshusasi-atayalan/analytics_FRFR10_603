# Analytics

Analytics Service

## Local Development

### Dependencies

#### InfluxDB

- Update service.environment to point to development Influx DB

#### Env Variables

- export LOCAL_DEV_OR_CI_MODE=true
<!-- Add local influx db token in -->

### Build

<!-- - go mod tidy (if there are package updates/additional imports) -->
<!-- @Pranshanth , please check this -->

- go work sync (if there are package updates/additional imports)

- make build

### Run

- go mod tidy <!-- install go deps -->
- mkdir build (create build directory for server logs)
- make run

### Test

- make test
