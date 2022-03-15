## Sharsec

This repository will host a couple toy implementation of interesting secret-sharing algorithms.

You will find Shamir's secret-sharing protocol in `shamir.go`.

An implementation of Schoenmakers' PVSS can be found in `pvss.go`

All implementations should implement the `sharsec.SSS` interface

Please do not use this in production. This is only a personal academic exercise.

#### Build:

`go mod tidy`

`go build`
