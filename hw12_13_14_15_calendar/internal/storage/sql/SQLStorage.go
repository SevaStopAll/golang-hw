package sqlstorage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/sevastopall/hw12_13_14_15_calendar/internal/storage/models"
	"strings"
	"time"
)

type SQLStorage struct {
	db           *sqlx.DB
	dbDriverName string
	dsn          string
}

func New(dbDriverName string, dsn string) *SQLStorage {
	dbDriverName = strings.ToLower(dbDriverName)
	dsn = strings.TrimSpace(dsn)
	return &SQLStorage{dbDriverName: dbDriverName, dsn: dsn}
}

func (r *SQLStorage) Connect(ctx context.Context) error {
	r.db = sqlx.MustConnect(r.dbDriverName, r.dsn)
	// Настройки ниже конфигурируют пулл подключений к базе данных. Их названия стандартны для большинства библиотек.
	// Ознакомиться с их описанием можно на примере документации Hikari pool:
	// https://github.com/brettwooldridge/HikariCP?tab=readme-ov-file#gear-configuration-knobs-baby
	r.db.SetMaxIdleConns(5)
	r.db.SetMaxOpenConns(20)
	r.db.SetConnMaxLifetime(1 * time.Minute)
	r.db.SetConnMaxIdleTime(10 * time.Minute)
	return nil
}

func (r *SQLStorage) Close(ctx context.Context) error {
	// TODO
	return nil
}

func (r *SQLStorage) Create(event models.Event) (id int64, err error) {
	result, err := r.db.Query(
		"INSERT INTO event (name) VALUES($1) RETURNING id", event.Title)
	if err != nil {
		return 0, err
	}
	if result.Next() {
		err = result.Scan(&id)
		return id, err
	}
	return id, err
}

func (r *SQLStorage) Update(event models.Event) {
	r.db.Query(
		"UPDATE event SET title =$1, date_time = $2 WHERE id = $3", event.Title, event.DateTime, event.ID)
}

func (r *SQLStorage) DeleteByID(eventID int64) (err error) {
	_, err = r.db.Query(
		"DELETE from event where id = $1", eventID)
	return err
}

func (r *SQLStorage) FindEventsByDay(date time.Time) (res []models.Event, err error) {
	err = r.db.Select(&res,
		"SELECT * from event WHERE date_time = $1", date.Format("2006-01-02"))
	return res, err
}

func (r *SQLStorage) FindEventsByWeek(date time.Time) (res []models.Event, err error) {
	err = r.db.Select(&res,
		"SELECT * from event WHERE date_time = $1 AND $2", date.Format("2006-01-02"), date.AddDate(0, 0, 7))
	return res, err
}

func (r *SQLStorage) FindEventsByMonth(date time.Time) (res []models.Event, err error) {
	err = r.db.Select(&res,
		"SELECT * from event WHERE date_time = $1 AND $2", date.Format("2006-01-02"), date.AddDate(0, 0, 30))
	return res, err
}

func (r *SQLStorage) FindByID(id int64) (res, err error) {
	err = r.db.Get(&res,
		"SELECT * from event where id = $1", id)
	return res, err
}

func (r *SQLStorage) FindAll() (res []models.Event, err error) {
	err = r.db.Select(&res,
		"SELECT * from event")
	return res, err
}

func (r *SQLStorage) FindByIDs(ids []int64) (res []models.Event, err error) {
	query, args, err := sqlx.In(
		"SELECT * from event where id IN(?)", ids)
	if err == nil {
		query = r.db.Rebind(query)
		var rows *sqlx.Rows
		rows, err = r.db.Queryx(query, args...)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var event models.Event
			err = rows.StructScan(&event)
			if err != nil {
				return nil, err
			}
			res = append(res, event)
		}
		return res, err
	}
	return make([]models.Event, 0), err
}

func (r *SQLStorage) DeleteByIDs(ids []int64) (err error) {
	query, args, err := sqlx.In(
		"DELETE from event where id IN(?)", ids)
	if err != nil {
		return err
	}
	query = r.db.Rebind(query)
	_, err = r.db.Query(query, args...)
	return err
}

func (r *SQLStorage) ExecuteQuery(query string) {
	_, err := r.db.Exec(query)
	if err != nil {
		return
	}
}
