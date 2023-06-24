# NSProxy

A simple and lightweight DNS proxy

## Installation
### Run using docker (recommended)

### Install from source

## Configuration

| Key                   | Description                                                                       | Default Value              |
|-----------------------|-----------------------------------------------------------------------------------|----------------------------|
| cache-expiration-time | The duration for which cache entries will remain active                           | TTL of the RRs             |
| cache-url             | The URL of the cache utilized for caching resource records (RRs)                  | "" (means in memory cache) |
| external-dns-servers  | A list of external servers utilized to resolve unresolved DNS requests            | ["8.8.8.8:53"]             |
| external-dns-timeout  | The expiration time for an external DNS request before attempting the next server | 0 (means no timeout)       |

## Upcoming

- [ ] Additional info from cache
- [ ] Generic wildcard matching

## Additional Info & References

- **RFC 1034** DOMAIN NAMES - CONCEPTS AND FACILITIES
- **RFC 1035** DOMAIN NAMES - IMPLEMENTATION AND SPECIFICATION
- **RFC 3596** DNS Extensions to Support IP Version 6
