package dto

type CustomerResponse struct {
	// set the json output tags in struct
	Id          string `json:"customerId" xml:"customerId"`
	Name        string `json:"fullName" xml:"name"`
	City        string `json:"city" xml:"city"`
	Postcode    string `json:"postcode"  xml:"postcode"`
	DateOfBirth string `json:"dateOfBirth"  xml:"dateOfBirth"`
	Status      string `json:"status"  xml:"status"`
}
