package database

type SavebleItem interface {
	GetId() string
}

type IDatabase interface {
	Save(keyPattern string, item SavebleItem) error
	Get(keyPattern string, id string) (interface{}, error)
	List(keyPattern string) ([]interface{}, error)
}
