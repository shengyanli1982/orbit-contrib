<div>
	<h1>orbit-contrib</h1>
    <p>Collection of middlewares created by the community</p>
	<img src="assets/logo.png" alt="logo" width="450px">
    
</div>

## Middlewares

The following middlewares can be used with [`orbit`](https://github.com/shengyanli1982/orbit) and [`gin`](https://github.com/gin-gonic/gin).

All middlewares are located in the [pkg](./pkg/) directory.

### 1. RateLimiter

The [**RateLimiter**](./pkg/ratelimiter/) middleware is used to limit the rate of requests. It is based on the [token bucket](https://en.wikipedia.org/wiki/Token_bucket) algorithm.

### 2. Compressor

The [**Compressor**](./pkg/compressor/) middleware is used to compress the response body. It supports `gzip` and `deflate` algorithms.

### 3. Rewriter

The [**Rewriter**](./pkg/rewriter/) middleware is used to rewrite the request path. It supports rewriting of `url.URL` related elements.
