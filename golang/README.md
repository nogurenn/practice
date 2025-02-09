# golang

Each project has its own Makefile. You only need `make` and Docker + Compose install to run everything.

```bash
make runSimpleEventWorker

# ...
# go run cmd/main.go
# [main] 2025/01/26 01:24:06 550e8400-e29b-41d4-a716-446655440003-Alice: {Status: Recalled, Balance: 1000}
# ...

make runLeetcode
```
