# Ratelimiter

## Introduction

**Ratelimiter** is a simple rate limiter for `Gin` and `orbit`. It is based on [token bucket](https://en.wikipedia.org/wiki/Token_bucket) algorithm. It is designed to be used in a protected API endpoint.

`Ratelimiter` has two modes:

-   **Base on IP**: Each IP address has its own bucket
-   **All share one**: All requests share the same bucket

`Ratelimiter` based golang native package and other powerful packages:

-   [gin](https://github.com/gin-gonic/gin)
-   [orbit](https://github.com/shengyanli1982/orbit)
-   [xxhash](https://github.com/cespare/xxhash)
-   [golang.org/x/time/rate](https://pkg.go.dev/golang.org/x/time/rate)

So it is safe to use in a concurrent environment.

## Installation

```bash
go get github.com/shengyanli1982/orbit-contrib/pkg/ratelimiter
```

## Quick Start

`Ratelimiter` is designed to be used in a protected API endpoint. It is recommended to use it with `orbit` or `gin`.

### Config

`Ratelimiter` has a config object, which can be used to configure the batch process behavior. The config object can be used following methods to set.

-   `WithCallback` : set the callback function. The default is `&emptyCallback{}`.
-   `WithRate`: set the rate. The default is `float64(1)`.
-   `WithBurst` : set the burst. The default is `1`.
-   `WithMatchFunc` : set the match function. The default is `DefaultLimitMatchFunc`.
-   `WithWhitelist` : set the whitelist. The default is `DefaultIpWhitelist`.

### Components

#### 1. Ratelimiter

`Ratelimiter` is used the most component. It is used to limit the rate of requests.

**Methods**

-   `GetLimiter` : get the limiter.
-   `SetRate` : set the rate for limiter, thread safe.
-   `SetBurst` : set the burst for limiter, thread safe.
-   `HandlerFunc` : return a `gin.HandlerFunc` for `orbit` or `gin`.
-   `Stop` : stop the limiter. It is empty function, no need to call it.

**Example**

```go

```

#### 2. Ip Ratelimiter

**Methods**

-   `GetLimiter` : get the limiter.
-   `SetRate` : set the rate for limiter, thread safe.
-   `SetBurst` : set the burst for limiter, thread safe.
-   `HandlerFunc` : return a `gin.HandlerFunc` for `orbit` or `gin`.
-   `Stop` : stop the limiter. It is used to release the resources.

**Example**

```go

```
