package main

import (
	"context"
	"log"
	"math/rand"
	"time"
	yellowcard "yellowcard-go"
)

const apiKey = "c5315180696a51ab885023bdc1ae3c0e"
const secret = "81473a4ca26a28203aa4a8e26afe571bf097ec5a2a5c2acd41c83a7968c4cf3b"

func main() {
	client := yellowcard.New(apiKey, secret, yellowcard.WithEnvironment(yellowcard.EnvironmentSandbox))
	ctx := context.Background()

	//channels, err := client.GetChannels(ctx, "")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//log.Printf("%+v", channels[0])

	//networks, err := client.GetNetworks(ctx, yellowcard.CountryCodeKE)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//log.Printf("%+v", networks)

	//rates, err := client.GetRates(ctx, "")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//log.Printf("%+v", rates)

	//accountDetails, err := client.ResolveBankAccount(ctx, &yellowcard.ResolveBankAccountRequest{
	//	AccountNumber: "589000",
	//	NetworkID:     "41109c18-9604-4389-8472-44ff4378c6cb",
	//})
	//
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//log.Printf("%+v",accountDetails)

	// {
	//   "channelId":"81018280-e320-4c81-9b2f-6f636c2239d8",
	//   "sequenceId":"1234-432121231a",
	//   "amount":7491.65,
	//   "reason":"entertainment",
	//   "destination":{
	//      "accountNumber":"+12222222222",
	//      "accountType":"momo",
	//      "networkId":"41109c18-9604-4389-8472-44ff4378c6cb",
	//      "country":"ZA",
	//      "accountBank":"589000",
	//      "accountName":"Ken Adams"
	//   },
	//   "sender":{
	//      "name":"Sample Name",
	//      "country":"US",
	//      "phone":"+12222222222",
	//      "address":"Sample Address",
	//      "dob":"10/10/1950",
	//      "email":"email@domain.com",
	//      "idNumber":"0123456789",
	//      "idType":"license"
	//   }
	//}

	paymentRequest := yellowcard.PaymentRequest{
		Amount:    7491.65,
		ChannelID: "81018280-e320-4c81-9b2f-6f636c2239d8",
		Destination: yellowcard.Destination{
			AccountBank:   "589000",
			AccountName:   "Ken Adams",
			AccountNumber: "+12222222222",
			AccountType:   "momo",
			Country:       "ZA",
			NetworkID:     "41109c18-9604-4389-8472-44ff4378c6cb",
		},
		Reason: "entertainment",
		Sender: yellowcard.Sender{
			Address:  "Sample Address",
			Country:  "US",
			Dob:      "10/10/1950",
			Email:    "email@domain.com",
			IdNumber: "0123456789",
			IDType:   "license",
			Name:     "Sample Name",
			Phone:    "+12222222222",
		},
		SequenceID: randSeq(10),
	}

	_ = paymentRequest

	paymentRequestResp, err := client.MakePayment(ctx, &paymentRequest)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("[paymentRequestResp] %+v", paymentRequestResp.ID)

	time.Sleep(10 * time.Second)
	acceptPaymentResp, err := client.AcceptPaymentRequest(ctx, "8b094b7b-5538-5645-a6ba-21038c2ca0d2")
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%+v", acceptPaymentResp)

	denyPaymentResp, err := client.DenyPaymentRequest(ctx, "970636ef-7a95-5fd9-875a-3103dd427e76")
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%+v", denyPaymentResp)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
