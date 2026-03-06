# Perf Test Target

A lightweight Go web application designed as a target for performance testing. It exposes several HTTP endpoints and tracks per-endpoint request counts in real time via a built-in dashboard.

![Screenshot](screenshot.png)

## Endpoints

| Method | Path          | Description                  |
|--------|---------------|------------------------------|
| GET    | `/`           | Dashboard (HTML)             |
| GET    | `/users`      | List users                   |
| POST   | `/users`      | Create a user                |
| GET    | `/tasks`      | List tasks                   |
| GET    | `/documents`  | List documents               |
| POST   | `/documents`  | Create a document            |
| GET    | `/health`     | Health check                 |
| GET    | `/api/stats`  | Per-endpoint request counts  |

## Prerequisites

- Go 1.25+
- Docker (optional, for containerised builds)
- [Vegeta](https://github.com/tsenart/vegeta) (optional, for load testing)

## Configuration

The app reads environment variables from a `.env` file (via [godotenv](https://github.com/joho/godotenv)):

| Variable         | Description                          | Example          |
|------------------|--------------------------------------|------------------|
| `HTTP_BIND_ADDR` | Address the server listens on        | `0.0.0.0:3000`   |
| `LOGLEVEL`       | Log level (`debug`,`info`,`warn`,`error`) | `info`       |

## Running locally

```bash
# create a .env file with at least:
echo 'HTTP_BIND_ADDR="0.0.0.0:3000"' > .env
echo 'LOGLEVEL="debug"' >> .env

# run directly
./run-dev.sh
```

## Docker

```bash
# build
docker build -t rogierlommers/perftest .

# run
docker run -p 3000:3000 \
  -e HTTP_BIND_ADDR="0.0.0.0:3000" \
  -e LOGLEVEL="info" \
  rogierlommers/perftest
```

Or use the convenience script:

```bash
./build-and-push.sh
```

## Load testing

A [Vegeta](https://github.com/tsenart/vegeta)-based load test script is included:

```bash
# install vegeta (macOS)
brew install vegeta

# run all endpoint attacks in parallel for 30 s
./loadtest.sh
```

Results (text reports and latency histograms) are written to `vegeta-results/`.

## CI/CD

- **GitHub Actions** — `.github/workflows/build-and-push.yaml` builds and pushes the Docker image to Docker Hub on every push to `main`.
- **Azure DevOps** — `devops-pipeline/pipeline.yaml` builds, pushes, and deploys to Azure Container Apps.
