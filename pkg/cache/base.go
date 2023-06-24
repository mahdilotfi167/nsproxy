package cache

type Manager interface {
	Get(key interface{}) interface{}
	Put(key interface{}, value interface{})
}

type InMemory struct {
}

type Redis struct {
	Url string
}
