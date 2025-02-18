# musicx-spotify-api

## Requirements
1. `docker` and `docker-compose`
2. [Optional] `direnv` for autoloading envvars.
    - I've  also included a script for loading the envvars manually. See `scripts/loadenv.sh`.
3. [Optional] `Postman` for testing the endpoints.
    - I have provided a quick and dirty Postman collection you could import for comfort.

Please execute all commands from project root, at least.

## Environment setup
```sh
make env

direnv allow # auto-loads envvars into your shell session upon entering the project directory
```

## Running the tests
```sh
make up; make migrate; make test
```

## Running the app
I suggest running `make down` for a bit of testdata cleanup. See design write-ups for context.
```sh
make down; make up; make migrate; make app
```

See also below for spotify artist IDs used in testing.
```
7tNO3vJC9zlHy2IJOx34ga
0XATRDCYuuGhk0oE7C0o5G
0BqRGrwqndrtNkojXiqIzL
46UMQ0cW8ToR8egkBRwAxZ
3WwGRA2o4Ux1RRMYaYDh7N
```

## Resetting the db
```sh
make down; make up; make migrate
```

## Exam Decisions
1. A few unit and integration tests were missed. Also purely because of my time constraints. Coverage was done on a best-effort and practicality basis.


## Design Decisions
1. 3 Musketeers is implemented for reproducibility of development, build, and runtime environments.
1. This server's entrypoint is located in `cmd/server/main.go`, in simulation of future requirements needing multiple binaries within one repository. Maybe `workers` in the future?
1. There was some difficulty implementing the correct code for proper separation of concerns among `Transaction`, `Repository`, and `sqlx.DB. As such, some of the unit and integration tests have hacks and shortcuts to get around time sink problems such as db-level uniqueness constraints. Solving this would have been a priority enhancement because it would cause problems for everyone managing the repository.
1. An injected logger was one of the planned enhancements. It did not happen due to time constraints.
1. The app's Dockerfile is composed of `build` and `run` stages, in anticipation of more decoupled CICD operations.
1. `Flyway` is my personal weapon of choice for database migrations.
    - Separating db migrations from the application provides safer testing and development of the db. Having the application layer handle migrations is dangerous on many levels.
    - I enforce "`forward` migrations only". In practice, maintaining `down` migrations only does more harm than good in production. Coding and testing `down` migrations is a time sink with no considerable long-term benefits.
    - One of the highlighted downsides is coordinating migrations with app deployments in production.
1. Deleting a `genre` record triggers a `cascade delete` on the artist-genre association table. It is a business decision that could go either way -- the current behavior was just chosen for unambiguity. See `V0003__create_artists_table.sql`.
1. A mock generator was used for faster iterations of mocks. I found writing your own mocks is best paired with IntelliJ's GoLand IDE because it can edit and replace function signatures and interfaces globally. As always, concerns like this depends on the team's values and priorities.
1. Default pagination values are enforced when consumer is providing invalid values such as negative numbers. Limit of `0` would evaluate to `limit: 30` through the service layer.
1. Handlers have been implemented such that any business rule and/or default values are managed by the service layer.

## Known Flaws
1. `GET /artists/{spotifyArtistId}` works more naively than desired. It handles the artist and genre relationship a bit carelessly. This was not intended.
1. The interface design among `Transaction`, `Repository`, and `sqlx.DB` is not as decoupled as I had hoped. The tests inside `{artist|genre}.repo_test.go` prove this.
1. No tests have been written to evaluate app behavior with Unicode values like non-English artists and genres.
1. A few failure cases of handler responses are unexpected. Aside from misleading responses for "already created" objects, nothing really major.

## Sample responses
### GET /artists/7tNO3vJC9zlHy2IJOx34ga
```json
{
    "artist": {
        "id": "3f268431-ed72-4704-8fe4-e2ab9ed00f43",
        "name": "BINI",
        "spotify_artist_id": "7tNO3vJC9zlHy2IJOx34ga",
        "spotify_uri": "spotify:artist:7tNO3vJC9zlHy2IJOx34ga",
        "created_at": "2024-12-17T10:53:41.356948Z",
        "updated_at": "2024-12-17T10:53:41.356948Z"
    },
    "genres": [
        "opm",
        "pinoy idol pop",
        "p-pop"
    ]
}
```

### GET /genres?page=2&limit=3
```json
{
    "genres": [
        {
            "id": "ccd2b49f-9f34-43c7-819b-5e5a1ca7fea7",
            "name": "japanese singer-songwriter",
            "created_at": "2024-12-17T10:53:41.356948Z",
            "updated_at": "2024-12-17T10:53:41.356948Z"
        },
        // .....
    ],
    "metadata": {
        "page": 2,
        "limit": 3,
        "total": 11,
        "total_pages": 4
    }
}
```
