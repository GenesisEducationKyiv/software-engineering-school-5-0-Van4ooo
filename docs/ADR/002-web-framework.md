# ADR-002: Choosing Web Framework

---

**Status:** Accepted   
**Author:** Ivan Rybachyk   
**Date:** 08.06.2025  

---

## Context
We require a lightweight, high-performance HTTP framework for:
- Serving the Weather App's REST API endpoints (subscribe, confirm, update)
- Middleware support (CORS, logging)
- Easy integration with PostgreSQL via GORM
- Good swagger documentation support

## Solution
### Options Considered
- **Gin**
- **Echo**
- **Fiber**
- **Go stdlib ```net/http```**

### Evaluation

**Gin**:
- Pros:
  - High performance, minimal memory footprint.
  - Middleware support.
  - Mature ecosystem with good documentation.
  - Simple routing and error handling.
  - Good support for Swagger documentation.
- Cons:
  - Slightly more complex than the standard library.
  - Slightly larger binary size compared to stdlib.

**Echo**:
- Pros:
  - High performance, similar to Gin.
  - Rich middleware ecosystem.
  - Good documentation and community support.
- Cons:
  - Less community momentum than Gin

**Go stdlib `net/http`**:
- Pros:
  - No external dependencies, minimal binary size.
  - Familiar and well-documented.
  - Full control over request handling.
- Cons:
  - More boilerplate code required for routing.
  - Less convenient for complex applications.

**Fiber**:
- Pros:
  - Very high performance.
  - Middleware support and easy routing.
  - Good documentation.
- Cons:
  - Newer, less mature than Gin or Echo.
  - Larger binary size due to additional features.
  - Less idiomatic Go style.

## Decision
After evaluating the options, we chose **Gin**.

## Consequences
**Positive**:
- High performance and low memory usage.
- Rich middleware ecosystem for CORS, logging, etc.
- Simple and intuitive routing.
- Easy Swagger/OpenAPI integration with ```swaggo/gin-swagger```

**Negative**:
- Slightly larger binary size compared to the standard library.
- Slightly more complex than using the standard library, but manageable.
- Dependency on an external library, which may require updates and maintenance.
