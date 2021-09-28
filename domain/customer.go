package domain

type Customer struct {
	// set the json output tags in struct
	Id          string
	Name        string `json:"full_name" xml:"name"`
	City        string `json:"city" xml:"city"`
	Postcode    string `json:"post_code"  xml:"postcode"`
	dateOfBirth string
	Status      string
}

type CustomerRepository interface {
	FindAll() ([]Customer, error)
}
