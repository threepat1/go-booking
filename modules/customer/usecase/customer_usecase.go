package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"booking-app/modules/customer/domain"
	"booking-app/modules/customer/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/gomail.v2"
)

type CustomerUsecase struct {
	CustomerRepo repository.CustomerRepository
	EmailSender  *gomail.Dialer
	AppURL       string
}

func NewCustomerUsecase(cr repository.CustomerRepository, emailSender *gomail.Dialer, appURL string) *CustomerUsecase {
	return &CustomerUsecase{
		CustomerRepo: cr,
		EmailSender:  emailSender,
		AppURL:       appURL,
	}
}

func (cu *CustomerUsecase) GetCustomers(ctx context.Context) ([]domain.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return cu.CustomerRepo.FindAll(ctx)
}

func (cu *CustomerUsecase) CreateCustomer(ctx context.Context, customer *domain.Customer) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	customer.ID = primitive.NewObjectID()
	customer.VerificationToken = generateToken()
	customer.IsVerified = false
	if err := cu.CustomerRepo.Save(ctx, customer); err != nil {
		return err
	}
	return cu.sendVerificationEmail(customer)
}

func (cu *CustomerUsecase) UpdateCustomer(ctx context.Context, id string, customer *domain.Customer) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	customer.ID = oid
	return cu.CustomerRepo.Update(ctx, oid, customer)
}

func (cu *CustomerUsecase) VerifyEmail(ctx context.Context, token string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return cu.CustomerRepo.VerifyEmail(ctx, token)
}

func (cu *CustomerUsecase) sendVerificationEmail(customer *domain.Customer) error {
	verificationLink := fmt.Sprintf("%s/verify-email?token=%s", cu.AppURL, customer.VerificationToken)
	msg := gomail.NewMessage()
	msg.SetHeader("From", "no-reply@example.com")
	msg.SetHeader("To", customer.Email)
	msg.SetHeader("Subject", "Email Verification")
	msg.SetBody("text/plain", fmt.Sprintf("Please verify your email by clicking on the following link: %s", verificationLink))
	return cu.EmailSender.DialAndSend(msg)
}

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
