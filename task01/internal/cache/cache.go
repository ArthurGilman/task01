package cache

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"task01/internal/config"
	"task01/internal/db"
	"time"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()
var client *redis.Client

func Init() error {
	addr := config.Get().Redis.Host + ":" + config.Get().Redis.Port
	pass := config.Get().Redis.Password
	db, err := strconv.Atoi(config.Get().Redis.DB)
	if err != nil {
		return errors.New("Failed to convert str to int kafka database number from env")
	}

	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	_, err = client.Ping(Ctx).Result()
	if err != nil {
		return errors.New("Failed to connect to kafka")
	}

	return nil
}

func CachePerson(person *db.Person) error {

	data, err := json.Marshal(person)
	if err != nil {
		return err
	}

	err = client.Set(Ctx, strconv.Itoa(person.Id), data, 1*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func getCachedPerson(key string) (*db.Person, error) {
	p := &db.Person{}

	data, err := client.Get(Ctx, key).Bytes()
	if err != nil {
		return p, err
	}

	err = json.Unmarshal(data, &p)
	if err != nil {
		return p, err
	}

	return p, nil
}

func GetPersonsFromCacheOrDatabase(id string) (*db.Person, error) {
	cachedPersons, err := getCachedPerson(id)
	if err == nil {
		return cachedPersons, nil
	}

	i, _ := strconv.Atoi(id)

	p, err := db.GetPersonById(i)

	if err != nil {
		return p, err
	}

	err = CachePerson(p)

	if err != nil {
		return p, err
	}

	return p, nil
}

func DeleteFromCache(id string) error {
	err := client.Del(Ctx, id).Err()

	if err != nil {
		return errors.New("Failed to delete person from cache")
	}

	return nil
}
