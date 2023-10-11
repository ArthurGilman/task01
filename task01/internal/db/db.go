package db

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"task01/internal/config"

	"github.com/jackc/pgx/v4"
)

type Person struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patrynymic"`
	Age        int    `json:"age"`
	Gender     string `json:"gender"`
	Country    string `json:"country"`
}
type Filter struct {
	Start   string `json:"startage"`
	End     string `json:"endage"`
	Orderby string `json:"orderby"`
	Limit   string `json:"limit"`
	Offset  string `json:"offset"`
}

var (
	conn  *pgx.Conn
	Ctx   = context.Background()
	table = ""
)

func ConnectionPsql() error {
	dsn := config.Get().PostgreSQL.DSN
	table = config.Get().PostgreSQL.Talbe
	maxAttempts := 5

	var err error

	for attempt := 1; attempt < maxAttempts; attempt++ {
		conn, err = pgx.Connect(Ctx, dsn)
		if err == nil {
			break
		} else {
			time.Sleep(1 * time.Second)
		}
	}

	if err != nil {
		return errors.New("Failed to connect postgres")
	}

	return nil
}

func InitStruct() error {
	file, err := os.Open("struct.sql")

	if err != nil {
		return err
	}
	defer file.Close()

	scan := bufio.NewScanner(file)

	for scan.Scan() {
		query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s", table, scan.Text())

		_, err := conn.Exec(Ctx, query)

		if err != nil {
			return err
		}
	}

	return nil
}

func AddPerson(p *Person) (int, error) {
	type Id struct {
		id int
	}

	i := &Id{}

	query := fmt.Sprintf("INSERT INTO %s (name, surname, patronymic, age, gender, country) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", table)

	row := conn.QueryRow(Ctx, query, p.Name, p.Surname, p.Patronymic, p.Age, p.Gender, p.Country)

	if err := row.Scan(i); err != nil {
		return 0, err
	}

	return i.id, nil
}

func DeletePerson(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = %d", table, id)

	_, err := conn.Exec(Ctx, query)

	return err
}

func UpdatePerson(p *Person) error {
	query := fmt.Sprintf("UPDATE %s SET (name, surname, patronymic, age, gender, country) = ($1, $2, $3, $4, $5, $6) WHERE id = %d", table, p.Id)

	_, err := conn.Exec(Ctx, query, p.Name, p.Surname, p.Patronymic, p.Age, p.Gender, p.Country)

	return err
}

func GetFiltredAge(startAge string, endAge string, orderBy string, limit string, offset string) ([]Person, error) {
	pers := make([]Person, 0, 10)

	query := fmt.Sprintf("SELECT * FROM %s WHERE age BETWEEN %s AND %s ORDER BY %s LIMIT %s OFFSET %s",
		table, startAge, endAge, orderBy, limit, offset)

	rows, err := conn.Query(Ctx, query)
	defer rows.Close()

	if err != nil {
		return pers, err
	}

	for rows.Next() {
		p := Person{}

		err := rows.Scan(&p.Id, &p.Name, &p.Surname, &p.Patronymic, &p.Age, &p.Gender, &p.Country)

		if err != nil {
			return pers, err
		}
		pers = append(pers, p)
	}

	return pers, err
}

func GetPersonById(id int) (*Person, error) {
	per := &Person{}

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=%d", table, id)

	row := conn.QueryRow(Ctx, query)

	err := row.Scan(per)

	if err != nil {
		return per, err
	}

	return per, err
}
