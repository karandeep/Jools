package lib

import (
	"../config"
	"github.com/bradfitz/gomemcache/memcache"
)

var mc *memcache.Client

func GetMCConnection() *memcache.Client {
	if mc == nil {
		configData := config.GetConfig()
		mc = memcache.New(
			configData.MC_1,
		)
	}
	return mc
}

func GetPrefixedKey(key string) string {
	configData := config.GetConfig()
	return configData.MC_KEY_PREFIX + key
}
