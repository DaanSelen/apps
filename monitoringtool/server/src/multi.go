package main

import (
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

var (
	key *otp.Key
)

func initTOTP() {
	key, _ = totp.Generate(totp.GenerateOpts{
		Issuer:      "NMAS",
		AccountName: "dselen@nerthus.nl",
	})
}

func getCode() string {
	currentTime := time.Now()
	optCode, _ := totp.GenerateCode(key.Secret(), currentTime)
	return optCode
}
