# NSProxy

A simple and lightweight DNS proxy

## Installation
### Run using docker (recommended)
- (Optional) Create config file `config.json` with the appropriate config values
- Run a container using the following command (Change options as desired)  
` docker run -d -p 1053:53/udp --name=mynsproxy -v ./config.json:/etc/nsproxy.json mahdilotfi/nsproxy:latest`

### Docker compose (persistent cache)
- (Optional) Edit file `config.json` with the appropriate values
- Run required containers using `docker compose up -f docker-compose.yml up -d`

### Install from source

## Configuration
You can override the following configurations in the `/etc/nsproxy.json` file.

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
