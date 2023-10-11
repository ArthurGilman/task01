package person

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"task01/internal/logger"

	"github.com/biter777/countries"
)

type Age struct {
	Age int `json:"age"`
}

type Gender struct {
	Gender string `json:"gender"`
}

type Country struct {
	CountryId   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type Nationality struct {
	Countries []Country `json:"country"`
}

type Info struct {
	Age         int
	Gender      string
	Nationality string
}

func FillField(formatURL string, name string, info interface{}) {
	url := fmt.Sprintf(formatURL, name)

	resp, err := http.Get(url)

	if err != nil {
		logger.InfoLog.Println(err)
		return
	}

	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)

	if err != nil {
		logger.InfoLog.Println(err)
		return
	}

	if err = json.Unmarshal(buf, info); err != nil {
		logger.InfoLog.Println(err)
		return
	}
}

func GetInfo(name string) (*Info, error) {
	i := &Info{}
	a := &Age{}
	g := &Gender{}
	n := &Nationality{}

	FillField("https://api.agify.io/?name=%s", name, a)
	FillField("https://api.genderize.io/?name=%s", name, g)
	FillField("https://api.nationalize.io/?name=%s", name, n)

	if a.Age == 0 || g.Gender == "" || len(n.Countries) == 0 {
		return i, errors.New("Incorrect format of FIO")
	}

	c := countries.ByName(n.Countries[0].CountryId).String()

	i.Age = a.Age
	i.Gender = g.Gender
	i.Nationality = c

	return i, nil
}
