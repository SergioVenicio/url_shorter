package models

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/speps/go-hashids"
)

type Url struct {
	Id         string    `json:"id"`
	Shorted    string    `json:"shorted"`
	Url        string    `json:"url"`
	CreateDate time.Time `json:"create_date"`
}

func generateId() []int {
	hd := hashids.NewData()
	hd.Salt = os.Getenv("SHORTER_SALT")
	hashid, _ := hashids.NewWithData(hd)
	encodeData, _ := hashid.Encode([]int{int(time.Now().Unix())})
	id, _ := hashid.DecodeWithError(encodeData)
	return id
}

func intToString(numbers []int) string {
	var buffer bytes.Buffer

	for i := 0; i < len(numbers); i++ {
		buffer.WriteString(strconv.Itoa(numbers[i]))
	}

	return buffer.String()
}

func (u Url) GetId() string {
	return u.Id
}

func (u *Url) SetShorted() error {
	if u.Url == "" {
		return errors.New("url field is required!")
	}

	u.Id = intToString(generateId())
	u.Shorted = fmt.Sprintf("%s/%s", os.Getenv("BASE_URL"), u.Id)
	u.CreateDate = time.Now()

	return nil
}
