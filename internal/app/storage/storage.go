package storage

type S interface {
	Add(url string) (string, error)
	Get(hash string) (string, error)
}
