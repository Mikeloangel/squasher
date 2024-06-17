package interfaces

type Storager interface {
	Get(short string) (string, error)
	Set(url string) string
}
