package models

type Access struct {
	Url    Url
	Visits int
}

func (a Access) GetId() string {
	return a.Url.Id
}
