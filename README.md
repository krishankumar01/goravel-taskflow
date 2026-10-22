# TaskFlow

A task management REST API built with [Goravel](https://www.goravel.dev) v1.18 and Go 1.25: accounts, roles, tasks, file uploads, background jobs and email.

This is the companion code for the [Goravel series on DevChalk](https://devchalk.com/goravel). Each part of the series ends at a tag in this repo.

## Follow along

- **Start a part:** check out the previous part's tag, for example `git checkout part-01` before Part 2.
- **See what a part changed:** compare two tags on GitHub, for example `/compare/part-01...part-02`.
- **Skip ahead:** check out any tag; each one builds and runs on its own.

## Run it

You need Go 1.25 or newer and Docker.

```bash
cp .env.example .env
go run . artisan key:generate
docker compose up -d
go run .
```

The API answers on http://localhost:3000. Postgres runs on port 5433 and Redis on 6380, as set in `docker-compose.yml` and `.env.example`.

## Found a problem?

Open an issue and say which part it is about.
