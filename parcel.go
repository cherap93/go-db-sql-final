package main

import (
	"database/sql"
	"errors"
)

// структуры для работы с посылками в БД
type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

// Добавление в БД строки с помощью INSERT INTO
func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at)VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// Получение одной строки из БД, возвращает структуру Parcel
func (s ParcelStore) Get(number int) (Parcel, error) {
	p := Parcel{}
	row := s.db.QueryRow("SELECT client, status, address, created_at FROM parcel WHERE number = :number",
		sql.Named("number", number))
	err := row.Scan(&p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}

	return p, nil
}

// Получение коллекции посылок определенного клиента с помощью Query (select)
func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	var res []Parcel

	rows, err := s.db.Query("SELECT number, status, address, created_at FROM parcel WHERE client = :client",
		sql.Named("client", client))
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := Parcel{}
		p.Client = client
		err := rows.Scan(&p.Number, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}

	return res, nil
}

// Обновление статуса в БД
func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		return err
	}
	return nil
}

// Обнолвение адреса, работает только если статус Registered
func (s ParcelStore) SetAddress(number int, address string) error {
	// менять адрес можно только если значение статуса registered
	status, err := s.getStatus(number)
	if err != nil {
		return err
	}
	if status != ParcelStatusRegistered {
		return errors.New("Parcel must be registered")
	}

	_, err = s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number",
		sql.Named("address", address),
		sql.Named("number", number))
	if err != nil {
		return err
	}
	return nil
}

// Удаление строки, работает только если статус Registered
func (s ParcelStore) Delete(number int) error {
	// удалять строку можно только если значение статуса registered
	status, err := s.getStatus(number)
	if err != nil {
		return err
	}
	if status != ParcelStatusRegistered {
		return errors.New("Parcel must be registered")
	}

	_, err = s.db.Exec("DELETE FROM parcel WHERE number = :number", sql.Named("number", number))
	if err != nil {
		return err
	}
	return nil
}

// Функция возвращает статус посылки из БД
func (s ParcelStore) getStatus(number int) (string, error) {
	var status string
	row := s.db.QueryRow("SELECT status FROM parcel WHERE number = :number", sql.Named("number", number))
	err := row.Scan(&status)
	if err != nil {
		return "", err
	}
	return status, nil
}
