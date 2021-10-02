package service

import (
	"fmt"
	"github.com/barnettt/banking/domain"
	"github.com/barnettt/banking/dto"
	"github.com/barnettt/banking/exceptions"
	"github.com/barnettt/banking/logger"
	"github.com/golang-sql/civil"
	"time"
)

type AccountService interface {
	// NewAccount(dto.AccountRequest) (*dto.NewAccountResponse, err *exceptions.AppError)
	Save(account *dto.AccountRequest) (*dto.NewAccountResponse, *exceptions.AppError)
}

type DefaultAccountService struct {
	repository domain.AccountRepository
}

func (defaultAccountService DefaultAccountService) Save(account *dto.AccountRequest) (*dto.NewAccountResponse, *exceptions.AppError) {
	dtoAcc := getDomainAccount(*account)
	accResponse, err := defaultAccountService.repository.Save(*dtoAcc)
	if err != nil {
		logger.Error(err.Message)
		return nil, err
	}
	return accResponse, nil
}

func getDomainAccount(request dto.AccountRequest) *domain.Account {
	amount := fmt.Sprintf("%f", request.Amount)
	request.Validate()
	// using time.Now().Format("2021-10-02T11:26:20")
	// kept failing as the date produced seemed to be mangled
	// using civil datetime seems to be ok with mysql db
	dateTime := civil.DateTimeOf(time.Now())
	domAcc := domain.Account{CustomerId: request.CustomerId,
		Amount:      amount,
		OpeningDate: dateTime.String(),
		AccountType: request.AccountType,
		Status:      "1",
	}
	return &domAcc
}

func NewAccountService(repo domain.AccountRepository) DefaultAccountService {
	return DefaultAccountService{repository: repo}
}
