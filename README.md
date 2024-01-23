# orbit-contrib

Collection of middlewares created by the community

## Middlewares

Following middlewares can be used with [`orbit`](https://github.com/shengyanli1982/orbit) and [`gin`](https://github.com/gin-gonic/gin).

All middlewares in [pkg](./pkg/) directory.

### 1. Rate Limiter

[**Ratelimiter**](./pkg/ratelimiter/) is used to limit the rate of requests. Based on [token bucket](https://en.wikipedia.org/wiki/Token_bucket) algorithm.

### 2. Compressor

[**Compressor**](./pkg/compressor/) is used to compress the response body. Supports `gzip` and `deflate` algorithms.
