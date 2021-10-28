package domain

import (
	"database/sql"
	_ "errors"
	"github.com/barnettt/banking/exceptions"
	"github.com/barnettt/banking/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type CustomerRepositoryDb struct {
	client *sqlx.DB
}

func (repository CustomerRepositoryDb) FindAll() ([]Customer, *exceptions.AppError) {
	findAllCustomers := "select customer_id, name, city, postcode, date_of_birth, status from customers"
	customers := make([]Customer, 0)
	err := repository.client.Select(&customers, findAllCustomers)
	if err != nil {
		logger.Error("Error accessing customer database " + err.Error())
		return nil, exceptions.NewDatabaseError("Unexpected database error ")
	}
	return customers, nil
}

func (repository CustomerRepositoryDb) FindById(id string) (*Customer, *exceptions.AppError) {
	customerQuery := "select customer_id, name, city, postcode, date_of_birth, status from customers where customer_id = ?"
	var customer Customer
	err := repository.client.Get(&customer, customerQuery, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exceptions.NewNotFoundError("Customer not found")
		}
		logger.Error("Error retrieving customer by id " + err.Error())
		return nil, exceptions.NewDatabaseError("Unexpected database error")
	}
	return &customer, nil
}

func (repository CustomerRepositoryDb) FindByStatus(status string) ([]Customer, *exceptions.AppError) {
	customerStatusQuery := "select customer_id, name, city, postcode, date_of_birth, status from customers where status = ?"
	customers := make([]Customer, 0)
	err := repository.client.Select(&customers, customerStatusQuery, status)
	if err != nil {
		logger.Error("Error accessing customer database " + err.Error())
		return nil, exceptions.NewDatabaseError("Unexpected database error ")
	}
	return customers, nil
}
func NewCustomerRepositoryDb(dbClient *sqlx.DB) CustomerRepositoryDb {
	return CustomerRepositoryDb{dbClient}
}
