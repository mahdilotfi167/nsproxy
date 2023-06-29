/*
Copyright © 2023 Mahdi Lotfi mahdilotfi167@gmail.com
*/

package config

type ServerConfig struct {
	ExternalDNSServers  []string `mapstructure:"external-dns-servers"`
	ExternalDNSTimeout  uint     `mapstructure:"external-dns-timeout"`
	CacheExpirationTime uint     `mapstructure:"cache-expiration-time"`
}

type CacheConfig struct {
	CacheExpirationTime uint   `mapstructure:"cache-expiration-time"`
	CacheURL            string `mapstructure:"cache-url"`
}
