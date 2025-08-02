# SentinelAI Architecture (Week 1)

               +---------------------+
               |     Frontend        |
               | React Dashboard     |
               +----------+----------+
                          |
                          v
+---------+    +-----------------+    +---------+
|  Redis  |<-->| Ingest Service  |--->| ML/AI   |
|  Queue  |    |   (Go)          |    | FastAPI |
+---------+    +-----------------+    +---------+
                          |
                          v
                  +---------------+
                  |  TimescaleDB  |
                  |  PostgreSQL   |
                  +---------------+
