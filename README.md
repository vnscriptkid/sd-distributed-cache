# sd-distributed-cache

- low level (1 machine)
  - storage
  - cache full (how do we know)
  - eviction: when to use which policies
  - communication
    - HTTP
    - raw TCP
    - GRPC 
  - concurrency: when to use which
    - pesstimistic: how to not lock the whole map
    - optimistic
  - TTL
    - pro-active: random sampling (20)
    - lazy deletion (at the time of calling GET key)

- high level (distributed)
  -  storage
  -  routing
  -  availability
  -  reliability
  -  scaling
