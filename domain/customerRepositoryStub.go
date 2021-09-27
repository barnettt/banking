package domain

type CustomerRepositoryStub struct {
	customers []Customer
}

func (stub CustomerRepositoryStub) FindAll() ([]Customer, error) {
	return stub.customers, nil
}

func NewCustomerRepositoryStub() CustomerRepositoryStub {
	customers := []Customer{
		{"1001", "Ayyub", "Luton", "LT01 8BH", "24/05/1994", "active"}, {"1002", "Umayamah", "London", "SE6 7TH", "09/02/1997", "active"},
		{"1003", "Sumayah", "London", "BR3 2QQ", "10/02/1998", "active"},
	}
	return CustomerRepositoryStub{customers}
}
