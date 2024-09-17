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
    - multiper servers (shared nothing arch) - mutually exclusive
  -  routing: proxy (has small db)
    - hash-based
      - hash(key) % numOfServers
      - consistent hashing 
    - range-based
  -  availability
    - mornitoring (orchestrator)
      - healthcheck -> inform proxy when node's down (quickly asap) -> update routing logic
      - prometheus
    - replicas, standby nodes
  -  reliability
    - WAL  
  -  scaling
    - minimize data transfer when node added, removed -> consistent hashing
    - practicals
      - node added
      - node removed gracefully
      - node shutdown-ed abruptly 


# Misc
- hash table
  - hash function
  - collision
    - chaining: linked list
    - linear probing: next empty slot
  - resizing
