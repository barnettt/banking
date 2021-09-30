package domain

import (
	"database/sql"
	_ "errors"
	"fmt"
	"github.com/barnettt/banking/exceptions"
	"github.com/barnettt/banking/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"os"
	"time"
)

type CustomerRepositoryDb struct {
	client *sqlx.DB
}

func (repository CustomerRepositoryDb) FindAll() ([]Customer, *exceptions.AppError) {
	findAllCustomers := "select customer_id, name, city, postcode, date_of_birth, status from customers"
	customers := make([]Customer, 0)
	error := repository.client.Select(&customers, findAllCustomers)
	if error != nil {
		logger.Error("Error accessing customer database " + error.Error())
		return nil, exceptions.NewDatabaseError("Unexpected database error ")
	}
	return customers, nil
}

func (repository CustomerRepositoryDb) FindById(id string) (*Customer, *exceptions.AppError) {
	customerQuery := "select customer_id, name, city, postcode, date_of_birth, status from customers where customer_id = ?"
	var customer Customer
	error := repository.client.Get(&customer, customerQuery, id)
	if error != nil {
		if error == sql.ErrNoRows {
			return nil, exceptions.NewNotFoundError("Customer not found")
		}
		logger.Error("Error retrieving customer by id " + error.Error())
		return nil, exceptions.NewDatabaseError("Unexpected database error")
	}
	return &customer, nil
}

func (repository CustomerRepositoryDb) FindByStatus(status string) ([]Customer, *exceptions.AppError) {
	customerStatusQuery := "select customer_id, name, city, postcode, date_of_birth, status from customers where status = ?"
	customers := make([]Customer, 0)
	error := repository.client.Select(&customers, customerStatusQuery, status)
	if error != nil {
		logger.Error("Error accessing customer database " + error.Error())
		return nil, exceptions.NewDatabaseError("Unexpected database error ")
	}
	return customers, nil
}
func NewCustomerRepositoryDb() CustomerRepositoryDb {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWD")
	dbName := os.Getenv("DB_NAME")
	dbProtocol := os.Getenv("DB_PROTOCOL")
	dbDrivername := os.Getenv("DB_DRIVER_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	client, err := sqlx.Open(fmt.Sprintf("%s", dbDrivername), fmt.Sprintf("%s:%s@%s(%s:%s)/%s", user, password, dbProtocol, dbHost, dbPort, dbName))
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	return CustomerRepositoryDb{client}
}
