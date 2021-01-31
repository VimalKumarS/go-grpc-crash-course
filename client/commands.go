package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	casinopb "gitlab.com/preslavmihaylov/go-grpc-exercise/gen/casino"
	commonpb "gitlab.com/preslavmihaylov/go-grpc-exercise/gen/common"
	"google.golang.org/grpc"
)

type command string

var errStopGambling = errors.New("user exits gambling session")

func setupClient() (casinopb.CasinoClient, *grpc.ClientConn) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithBlock())

	conn, err := grpc.Dial("localhost:10000", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}

	return casinopb.NewCasinoClient(conn), conn
}

func buyTokens(tokensCnt int) (string, error) {
	res, err := client.BuyTokens(context.Background(), &commonpb.Payment{
		User: &commonpb.User{
			Id: username,
		},
		Amount: int32(tokensCnt),
	})
	if err != nil {
		return "", fmt.Errorf("couldn't buy tokens: %w", err)
	}

	return fmt.Sprintf("Successfully bought %v tokens!", res.Count), nil
}

func withdraw(tokensCnt int) (string, error) {
	res, err := client.Withdraw(context.Background(), &casinopb.WithdrawRequest{
		User:      &commonpb.User{Id: username},
		TokensCnt: int32(tokensCnt),
	})
	if err != nil {
		return "", fmt.Errorf("couldn't withdraw tokens: %w", err)
	}

	return fmt.Sprintf("Successfully withdrew %d tokens and got %d$!", tokensCnt, res.GetAmount()), nil
}

func tokenBalance() (string, error) {
	res, err := client.GetTokenBalance(context.Background(), &commonpb.User{
		Id: username,
	})
	if err != nil {
		return "", fmt.Errorf("couldn't get token balance: %w", err)
	}

	return fmt.Sprintf("Your token balance is %v.", res.GetCount()), nil
}

func payments() (string, error) {
	stream, err := client.GetPayments(context.Background(), &commonpb.User{Id: username})
	if err != nil {
		return "", err
	}

	payments := []*commonpb.Payment{}
	for payment, err := stream.Recv(); err != io.EOF; {
		if err != nil {
			return "", err
		}

		payments = append(payments, payment)
	}

	return fmt.Sprintf("Here's your payment history:\n%v", paymentHistoryString(payments)), nil
}

func paymentStatement() (string, error) {
	panic("not implemented")
}

func gamble() (string, error) {
	panic("not implemented")
}
