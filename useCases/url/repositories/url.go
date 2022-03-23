package repositories

import (
	"fmt"

	"github.com/SergioVenicio/url_shorter/shared/database"
	"github.com/SergioVenicio/url_shorter/useCases/url/models"
)

func Save(u *models.Url) error {
	return database.Save[*models.Url]("urls", u)
}

func Get(id string) (models.Url, error) {
	return database.Get[models.Url]("urls", id)
}

func List(offset int64, limit int64) ([]models.Url, error) {
	return database.List[models.Url]("urls:*", offset, limit)
}

func GetAccess(id string) (models.Access, error) {
	var access models.Access
	counter := database.GetCounter(fmt.Sprintf("access:%s", id))
	access.Visits = int(counter)
	access.Url, _ = Get(id)
	return access, nil
}

func IncrementAccess(id string) {
	database.Increment(fmt.Sprintf("access:%s", id))
}
