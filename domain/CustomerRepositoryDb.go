package domain

import (
	"database/sql"
	"github.com/barnettt/banking/exceptions"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type CustomerRepositoryDb struct {
	client *sql.DB
}

func (repository CustomerRepositoryDb) FindAll() ([]Customer, *exceptions.AppError) {

	findAllCustomers := "select customer_id, name, city, postcode, date_of_birth, status from customers"
	rows, error := repository.client.Query(findAllCustomers)
	var customer Customer
	if error != nil {
		log.Println("Error accessing customer database ", error)
		return nil, exceptions.NewDatabaseError("Unexpected database error ")
	}
	return repository.getCustomersFromRows(rows, customer)
}

func (repository CustomerRepositoryDb) FindById(id string) (*Customer, *exceptions.AppError) {
	customerQuery := "select customer_id, name, city, postcode, date_of_birth, status from customers where customer_id = ?"
	var customer Customer
	row := repository.client.QueryRow(customerQuery, id)
	err := row.Scan(&customer.Id, &customer.Name, &customer.City, &customer.Postcode, &customer.dateOfBirth, &customer.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exceptions.NewNotFoundError("Customer not found")
		}
		log.Println("Error retrieving customer by id")
		return nil, exceptions.NewDatabaseError("Unexpected database error")
	}
	return &customer, nil
}

func (repository CustomerRepositoryDb) FindByStatus(status string) ([]Customer, *exceptions.AppError) {
	customerStatusQuery := "select customer_id, name, city, postcode, date_of_birth, status from customers where status = ?"

	rows, error := repository.client.Query(customerStatusQuery, status)
	var customer Customer
	if error != nil {
		log.Println("Error accessing customer database ", error)
		return nil, exceptions.NewDatabaseError("Unexpected database error ")
	}
	return repository.getCustomersFromRows(rows, customer)
}

func (repository CustomerRepositoryDb) getCustomersFromRows(rows *sql.Rows, customer Customer) ([]Customer, *exceptions.AppError) {
	customers := make([]Customer, 0)
	for rows.Next() {
		error := rows.Scan(&customer.Id, &customer.Name, &customer.City, &customer.Postcode, &customer.dateOfBirth, &customer.Status)
		if error != nil {
			log.Println("Error parsing customer row from database ")
			return nil, exceptions.NewDatabaseError("Unexpected database error ")
		}
		customers = append(customers, customer)
	}
	return customers, nil
}
func NewCustomerRepositoryDb() CustomerRepositoryDb {
	client, err := sql.Open("mysql", "banking:banking@tcp(localhost:3306)/banking")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	return CustomerRepositoryDb{client}
}
