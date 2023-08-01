package group

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetFunc func(string) ([]byte, error)

func (f GetFunc) Get(key string) ([]byte, error) {
	return f(key)
}
