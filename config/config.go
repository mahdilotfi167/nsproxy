/*
Copyright © 2023 Mahdi Lotfi mahdilotfi167@gmail.com
*/

package config

type ServerConfig struct {
	CacheExpirationTime uint     `json:"cache-expiration-time"`
	ExternalDNSServers  []string `json:"external-dns-servers"`
}
