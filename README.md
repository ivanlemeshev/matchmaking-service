# Matchmaking Service

![Build](https://github.com/ivanlemeshev/matchmaking-service/actions/workflows/build.yml/badge.svg)
![Tests](https://github.com/ivanlemeshev/matchmaking-service/actions/workflows/test.yml/badge.svg)
![Linter](https://github.com/ivanlemeshev/matchmaking-service/actions/workflows/lint.yml/badge.svg)

## Description

This is a simple matchmaking service that allows players to find matches based
on their level.

It is possible to configure the number of players in a match and the
matchmaking timeout.

It uses Protobuf and gRPC for communication between clients and the server.

Buf is used for generating Protobuf files.

There are configured GitHub Actions for building, testing and linting.

It is not production-ready and does not have a lot of nessesary features like:

- Persistence
- Graceful shutdown
- Health checks
- Logging
- Monitoring
- Security
- Deployment
- etc.

I have implemented only the core functionality and matchmaking logic. It still
requires a lot of work to be production-ready and some refactoring.

I used GitHub Copilot for some autocomplete suggestions, tests and comments.

I also used Gemini Advanced 2.0 Flash to resolve some issues and find some
matchmaking algorithms I can use.

## Requirements

- [Go 1.24.0](https://go.dev/doc/install)
- [Buf 1.50.0](https://buf.build/docs/cli/installation/)
- [grpcurl](https://github.com/fullstorydev/grpcurl?tab=readme-ov-file#installation)

## Run and make requests

```bash
make run
```

Request example:

```bash
grpcurl -plaintext -d '{"player_id": "123", "player_level": 1}' localhost:8080 matchmaking.v1.MatchmakingService/FindMatch
ERROR:
  Code: NotFound
  Message: Failed to find match
```

```bash
grpcurl -plaintext -d '{"player_id": "123", "player_level": 1}' localhost:8080 matchmaking.v1.MatchmakingService/FindMatch
{
  "matchId": "1"
}
```
