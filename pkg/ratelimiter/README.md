# Ratelimiter

## Introduction

**Ratelimiter** is a simple rate limiter for `Gin` and `orbit`. It is based on [token bucket](https://en.wikipedia.org/wiki/Token_bucket) algorithm. It is designed to be used in a protected API endpoint.

`Ratelimiter` has two modes:

-   **Base on IP**: Each IP address has its own bucket
-   **All share one**: All requests share the same bucket

`Ratelimiter` based golang native package and other powerful packages:

-   [Gin](https://github.com/gin-gonic/gin)
-   [orbit](https://github.com/shengyanli1982/orbit)
-   [xxhash](https://github.com/cespare/xxhash)
-   [golang.org/x/time/rate](https://pkg.go.dev/golang.org/x/time/rate)

So it is safe to use in a concurrent environment.

## Installation

```bash
go get github.com/shengyanli1982/orbit-contrib/pkg/ratelimiter
```

## Quick Start
