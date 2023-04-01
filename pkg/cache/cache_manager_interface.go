package cache

type CacheManagerInterface interface {
	Set(string, string)
	Get(string) (string, bool)
}
