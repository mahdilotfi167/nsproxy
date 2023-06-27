# NSProxy
![Build Release Status](https://github.com/mahdilotfi167/nsproxy/actions/workflows/ci.build.release.yml/badge.svg)  
Nowadays, with the significant increase in websites and people's usage of them, DNS servers have to handle a huge number of domain resolution requests. It is interesting to note that a typical home connected to the internet generates around 10,000 DNS queries per day!

DNS proxy acts as an intermediary between DNS clients and DNS servers, forwarding DNS requests and replies.
![proxy.png](docs/proxy.png)

## Installation
### Run using Docker (recommended)
- (Optional) Create a config file named `config.json` with the appropriate configuration values.
- Run a container using the following command (Modify the options as needed):  
`docker run -d -p 1053:53/udp --name=mynsproxy -v ./config.json:/etc/nsproxy.json mahdilotfi/nsproxy:latest`

### Docker Compose (persistent cache)
- (Optional) Edit the file `config.json` with the desired values.
- Run the required containers using the command:  
`docker-compose -f docker-compose.yml up -d`

### Install from source
- Execute `sudo make install` to build and install the project.
- Start the server by running `sudo systemctl start nsproxy`.

#### Uninstallation
- To stop the server and remove the installed files, run `sudo make uninstall`.

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
- [ ] Fast recovery on external cache failure
- [ ] Case-insensitive matching

## Additional Info & References

- **RFC 1034** DOMAIN NAMES - CONCEPTS AND FACILITIES
- **RFC 1035** DOMAIN NAMES - IMPLEMENTATION AND SPECIFICATION
- **RFC 3596** DNS Extensions to Support IP Version 6
