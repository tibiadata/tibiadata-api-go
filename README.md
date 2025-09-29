# TibiaData API in Golang

[![GitHub CI](https://img.shields.io/github/actions/workflow/status/tibiadata/tibiadata-api-go/build.yml?branch=main&logo=github)](https://github.com/tibiadata/tibiadata-api-go/actions/workflows/build.yml)
[![Codecov](https://codecov.io/gh/TibiaData/tibiadata-api-go/branch/main/graph/badge.svg?token=PSBNLBI10C)](https://codecov.io/gh/TibiaData/tibiadata-api-go)
[![GitHub go.mod version](https://img.shields.io/github/go-mod/go-version/tibiadata/tibiadata-api-go?logo=go)](https://github.com/tibiadata/tibiadata-api-go/blob/main/go.mod)
[![GitHub release](https://img.shields.io/github/v/release/tibiadata/tibiadata-api-go?sort=semver&logo=github)](https://github.com/tibiadata/tibiadata-api-go/releases)
[![Docker image size (tag)](https://img.shields.io/docker/image-size/tibiadata/tibiadata-api-go/latest?logo=docker)](https://hub.docker.com/r/tibiadata/tibiadata-api-go)
[![GitHub license](https://img.shields.io/github/license/tibiadata/tibiadata-api-go)](https://github.com/tibiadata/tibiadata-api-go/blob/main/LICENSE)

TibiaData API written in Golang and deployed in container (version v3 and above).

Documentation of API endpoints can be found on [docs.tibiadata.com](https://docs.tibiadata.com).

### Table of Contents

- [API versions](#api-versions)
- [How to use](#how-to-use)
  - [Helm](#helm)
  - [Docker](#docker)
  - [Docker-compose](#docker-compose)
  - [Local development](#local-development)
  - [Environment variables](#environment-variables)
  - [Deployment note](#deployment-note)
- [API documentation](#api-documentation)
  - [Available endpoints](#available-endpoints)
  - [Deprecated endpoints](#deprecated-endpoints)
  - [Restricted endpoints](#restricted-endpoints)
- [General information](#general-information)
- [Credits](#credits)

## API versions

Here is a summary of the TibiaData API versions

> **v4** is released _(since 1st December 2024 )_\
> **v3** is deprecated _(since 31rd January 2024)_\
> **v2** is deprecated _(since 30rd April 2022)_\
> **v1** is deprecated _(since 30rd April 2018)_

## How to use

You can either use it in a Docker container or go download the code and deploy it yourself on any server.

Keep in mind that there are restrictions on tibia.com that might impact the usage of the application being hosted yourself.

### Helm

We have a Helm chart available for you to use to deploy your application to Kubernetes.

All our charts are available through [charts.tibiadata.com](https://charts.tibiadata.com).

Add the repository

```console
helm repo add tibiadata https://charts.tibiadata.com
helm repo update
```

Search for the chart

```console
helm repo search tibiadata
```

The charts-repository is located in [tibiadata-helm-charts](https://github.com/tibiadata/tibiadata-helm-charts).

### Docker

Our images are available on both [GitHub Container Registry](https://github.com/tibiadata/tibiadata-api-go/pkgs/container/tibiadata-api-go) and [Docker Hub](https://hub.docker.com/r/tibiadata/tibiadata-api-go).

This is how to pull and run the _latest_ release of TibiaData from [GHCR](https://github.com/tibiadata/tibiadata-api-go/pkgs/container/tibiadata-api-go):

```console
docker pull ghcr.io/tibiadata/tibiadata-api-go:latest
docker run -p 127.0.0.1:80:8080/tcp --rm -it ghcr.io/tibiadata/tibiadata-api-go:latest
```

You can also use [Docker Hub](https://hub.docker.com/r/tibiadata/tibiadata-api-go) to pull your images from.

If you want to run the latest code you can switch from _latest_ to _edge_.

### Docker-compose

This is a simple example on how you can get up and running with TibiaData in docker-compose, which will be running on port 8080 and be exposed locally.

```yaml
version: "3"

services:
  tibiadata:
    image: ghcr.io/tibiadata/tibiadata-api-go:latest
    restart: always
    environment:
      - TIBIADATA_HOST=tibiadata.example.com
    ports:
      - 8080:8080
```

### Local development

Build the code on your computer

```console
docker build -t tibiadata .
```

Run your build locally

```console
docker run -p 127.0.0.1:80:8080/tcp --rm -it tibiadata
```

### Environment variables

_Information will be added at a later stage._

### Deployment note

You should consider to add a layer in front of this application, so you can do caching of endpoints, access controll or what ever your needs are.

We do so at least by using [Kong](https://github.com/Kong/kong) API Gateway, which solves features like caching, rate-limiting, authentication and more.

## API documentation

The hosted API documentation for our [api.tibiadata.com](https://api.tibiadata.com) service can be viewd at [docs.tibiadata.com](https://docs.tibiadata.com).

There is a swagger-generated documentation available for download on the [GitHub Release](https://github.com/tibiadata/tibiadata-api-go/releases) of the version you are looking for.

### Available endpoints

Those are the current existing endpoints.

- GET `/`
- GET `/ping`
- GET `/healthz`
- GET `/readyz`
- GET `/v4/boostablebosses`
- GET `/v4/character/:name`
- GET `/v4/creature/:race`
- GET `/v4/creatures`
- GET `/v4/fansites`
- GET `/v4/guild/:name`
- GET `/v4/guilds/:world`
- GET `/v4/highscores/:world/:category/:vocation/:page`
- GET `/v4/house/:world/:house_id`
- GET `/v4/houses/:world/:town`
- GET `/v4/killstatistics/:world`
- GET `/v4/news/archive`
- GET `/v4/news/archive/:days`
- GET `/v4/news/id/:news_id`
- GET `/v4/news/latest`
- GET `/v4/news/newsticker`
- GET `/v4/spell/:spell_id`
- GET `/v4/spells`
- GET `/v4/world/:name`
- GET `/v4/worlds`
- GET `/versions`

### Deprecated Endpoints

In addition to the deprecated API versions like v1, v2 and v3, there are also some endpoints that are deprecated. As of now, those are:

- GET `/health`

### Restricted endpoints

There are some endpoints that can be deviant between the container documentation and the hosted version. This is due to restricted mode that restrict certain API actions due to high load on the tibia.com servers.

- `/v4/highscores`-filtering on vocation is removed, only the `all` category is valid.

## General information

Tibia is a registered trademark of [CipSoft GmbH](https://www.cipsoft.com/en/). Tibia and all products related to Tibia are copyright by [CipSoft GmbH](https://www.cipsoft.com/en/).

## Credits

- Authors: [Tobias Lindberg](https://github.com/tobiasehlert) – [List of contributors](https://github.com/tibiadata/tibiadata-api-go/graphs/contributors)
- Distributed under [MIT License](LICENSE)
