# URL Shortener

## Overview
This is a URL Shortener application built using Golang with the Gin framework, MySQL as the database, and GORM as the ORM.

## How to Run the Application

### 1. Build and Run Using Docker
Ensure you have Docker and Docker Compose installed.

1. Clone the repository:
   ```sh
   git clone https://github.com/your-repo/urlshortener.git
   cd urlshortener
   ```

2. Build and start the application using Docker Compose:
   ```sh
   docker-compose up --build
   ```

3. To stop the application:
   ```sh
   docker-compose down
   ```

## API Endpoints with Curl Examples

### 1. Shorten a URL
**Endpoint:** `POST /api/v1/shorten`
```sh
curl -X POST http://localhost:8080/api/v1/shorten \
     -H "Content-Type: application/json" \
     -d '{"url": "https://example.com"}'
```
**Response:**
```json
{
  "short_url": "http://localhost:8080/abc123"
}
```

### 2. Retrieve Original URL
**Endpoint:** `GET /:shortCode`
```sh
curl -X GET http://localhost:8080/abc123 -v
```
**Response:** Redirects to `https://example.com`

### 3. Get Top Domains
**Endpoint:** `GET /api/v1/metrics/top-domains`
```sh
curl -X GET http://localhost:8080/api/v1/metrics/top-domains
```
**Response:**
```json
{
  "domains": [
    { "domain": "example.com", "count": 10 },
    { "domain": "test.com", "count": 5 }
  ]
}
```

## Design Decisions Explained

### 1. **Gin Framework for HTTP Handling**
**Why?** Gin is lightweight, fast, and provides built-in middleware for request handling. It suits high-performance applications like URL shorteners.

### 2. **MySQL for Data Storage**
**Why?**
- Relational integrity ensures no duplicate short codes.
- Transactions ensure atomic operations.
- Indexing on `short_code` allows fast lookups.
- Can be easily extended for analytics.

### 3. **GORM as ORM**
**Why?**
- Simplifies database interactions.
- Auto-migration for schema changes.
- Easier query building and struct mapping.

### 4. **Short Code Generation Using Random Characters**
**Why?**
- Provides a compact and unique identifier.
- Uses a predefined charset with alphanumeric values for easy readability.

## Possible Improvements

### 1. **Rate Limiting**
- Prevents abuse of the shortening service.
- Can be implemented using middleware or Redis.

### 2. **Expiration for Shortened URLs**
- Allows URLs to expire after a defined time.
- Can be stored as a TTL field in the database.

### 3. **Analytics & Click Tracking**
- Store click timestamps for tracking trends.
- Show analytics per URL (number of visits, user agents, locations).

### 4. **Caching with Redis**
- Cache frequent lookups for faster redirections.
- Reduce database load.

### 5. **Support for Custom Short Codes**
- Allow users to specify their own short URLs.
- Ensure uniqueness before assigning.

### 6. **Scale the system to support large number of concurrent users**

### 7. **Use NoSql Database**
- Provides high R/W throughput.
- Easily scalable in comparison to RDBMS.
