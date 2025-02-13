# Matchmaking Service

![Build](https://github.com/ivanlemeshev/matchmaking-service/actions/workflows/build.yml/badge.svg)
![Tests](https://github.com/ivanlemeshev/matchmaking-service/actions/workflows/test.yml/badge.svg)
![Linter](https://github.com/ivanlemeshev/matchmaking-service/actions/workflows/lint.yml/badge.svg)

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
grpcurl -plaintext -d '{"player_id": "123", "player_level": 1}' localhost:8080 matchmaking.v1.MatchmakingService/Enqueue
{
  "matchId": "1"
}
```
