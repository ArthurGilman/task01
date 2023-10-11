package person

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInfo(t *testing.T) {
	a := assert.New(t)

	exp := &Info{
		Age:         56,
		Gender:      "male",
		Nationality: "Bulgaria",
	}

	i, err := GetInfo("Nikolay")
	if a.NoError(err) {
		a.Equal(exp, i)
	}
}

func TestGetInfoError(t *testing.T) {
	a := assert.New(t)

	expErr := errors.New("Incorrect format of FIO")
	_, err := GetInfo("Ar.tur")

	if a.Error(err) {
		a.Equal(expErr, err)
	}
}
