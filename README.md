# User Analytics Service by Bagus

Counts **Daily Unique Users** and **Monthly Unique Users** based on user login events.
---

## Overview

* **Raw events:** store all the raw user login events in `user_logins`
* **Rollups:** for a performant and uniqueness queries, using rollup table `daily_unique_users` and `monthly_unique_users` with idempotency.
* **Queries:** expose two API to retrieve Daily Unique Users and Monthly Unique Users

---

## Quick Start (using Docker Compose)

### Prerequisites

* Docker & Docker Compose

### 1) Environment

```bash
cp .env.example .env
```

### 2) Run

```bash
docker compose up --build
```

### 3) Stop & clean

```bash
docker compose down
# ALTERNATIVE to also remove db data:
docker compose down -v
```

---

## API

### Ingest login

```
POST /v1/user/login
Content-Type: application/json
```

Body:

```json
{
  "user_id": "05378ca8-961d-49e7-a903-8026dad78bb7",
  "login_time": "2025-08-10T02:15:00Z"
}
```

Response: `202 Accepted`

### Daily uniques

```
GET /v1/user/uniques/daily?date=YYYY-MM-DD
```

Example:

```bash
curl "http://localhost:8080/v1/user/uniques/daily?date=2025-08-10"
```

Response:

```json
{"date":"2025-08-10","unique_users":123}
```

### Monthly uniques

```
GET /v1/user/uniques/monthly?month=YYYY-MM
```

Example:

```bash
curl "http://localhost:8080/v1/user/uniques/monthly?month=2025-08"
```

Response:

```json
{"month":"2025-08","unique_users":987}
```

---

## Sample Flow (Retrieve daily unique user)

Add two login as seed data, and retrieve the daily unique:

```bash
curl -s -XPOST localhost:8080/v1/user/login -H 'content-type: application/json' -d '{
  "user_id":"05378ca8-961d-49e7-a903-8026dad78bb7",
  "login_time":"2025-08-09T16:30:00Z"
}'

curl -s -XPOST localhost:8080/v1/user/login -H 'content-type: application/json' -d '{
  "user_id":"05378ca8-961d-49e7-a903-8026dad78bb7",
  "login_time":"2025-08-09T17:30:00Z"
}'

curl "http://localhost:8080/v1/user/uniques/daily?date=2025-08-09"
# => {"date":"2025-08-09","unique_users":1}
```

---

## Design Decisions & Assumptions

1. **[Database design decision] rollup tables (daily and monthly) to deduplicate user and faster query**
    * alternatives:
      * just have raw table `user_logins` and query the raw table directly using `SELECT COUNT(DISTINCT user_id)` over a time window for day and month
      * have additional tables `daily_unique_user_count` and `monthly_unique_user_count` that just store day/month with the count
    * decision:

      | Approach                                        | Write cost                         | Read cost                | Freshness | Complexity |
      | ----------------------------------------------- | ---------------------------------- |--------------------------| --------- | ---------- |
      | **just user_logins table**                      | Low                                | **High** `DISTINCT`      | Real‑time | Simple     |
      | **add unique_users table**                      | Moderate                | **Low** (`COUNT(*)`)     | Real‑time | Moderate   |
      | **add unique_users and count table**            | Higher                 | **Lowest** (Primary Key) | Real‑time | Highest    |

      * chose the middle complexity, that has low read cost and a moderate write cost. 
      * in addition, add the `ON CONFLICT DO NOTHING` on the unique table to ensure uniqueness during insertion
      * use transaction for all of these three tables insertion, to ensure consistency
      * the count table can be additional improvement once query latency became an issue/bottleneck

2. **Store everything in UTC**
    * alternatives:
      * store the timezone as separate column
      * query using timezone
    * decisions:
      * Store `login_time` in UTC
      * this will make us to have 1 single source of truth, in one timezone (UTC)
      * one of the concern is, if user login in Singapore with day 2025-08-12 but in UTC its still 2025-08-11, this might not fully represent the fact from user perspective
      * so it really depends on the purpose of querying this data for

3. **REST API for User login event ingestion and no queue**
    * alternatives:
      * the event source produce event directly to some queue, and the analytics service directly consuming from queue
      * have a queue in between the API and writing to database
    * decision:
      * Simple and straightforward for our first implementation of the service.
      * It is also easier to evaluate since it has strong consistency.
      * For backpressure, we can always return 429 to ask caller to call again later, and since we have idempotency, it should be fine even if the call was made more than one.
      * Queue adds complexity, and make it eventual consistent. Only add this once we see performance degradation, especially the DB write capacity.