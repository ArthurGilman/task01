package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"task01/internal/config"
	"task01/internal/db"
	"task01/internal/logger"
	"task01/internal/person"

	"github.com/segmentio/kafka-go"
)

var (
	r   *kafka.Reader
	w   *kafka.Writer
	Ctx = context.Background()
)

type FIO struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

func Init() error {
	conf := config.Get()
	host := conf.Kafka.Host
	port := conf.Kafka.Port

	addr := host + ":" + port

	r = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{addr},
		Topic:    "FIO",
		MaxBytes: 5e6,
	})

	w = &kafka.Writer{
		Addr:     kafka.TCP(addr),
		Topic:    "FIO_FAILED",
		Balancer: &kafka.LeastBytes{},
	}

	stats := r.Stats()

	if stats.Dials == 0 {
		return errors.New("Failed to connect kafka")
	}

	return nil
}

func ReadMessage() (*db.Person, error) {
	fio := &FIO{}

	m, err := r.ReadMessage(Ctx)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(m.Value, fio)

	if err != nil {
		err = errors.New("Invalid json data")
		return nil, err
	}

	if fio.Name == "" || fio.Surname == "" {
		err = errors.New("No required field")
		return nil, err
	}

	info, err := person.GetInfo(fio.Name)

	p := &db.Person{}

	p.Name = fio.Name
	p.Surname = fio.Surname
	p.Patronymic = fio.Patronymic
	p.Age = info.Age
	p.Gender = info.Gender
	p.Country = info.Nationality

	return p, err
}

func ReportFail(key string, errMsg string) error {
	err := w.WriteMessages(Ctx,
		kafka.Message{
			Key:   []byte(key),
			Value: []byte(errMsg),
		})

	return err
}

func StartReading() {
	id := 0

	for {
		p, err := ReadMessage()
		if err != nil {
			id++
			logger.InfoLog.Println(err)
			ReportFail(strconv.Itoa(id), err.Error())
		}

		db.AddPerson(p)
	}

}
