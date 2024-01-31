package main

import (
	"context"
	"fmt"

	"github.com/yasir7ca/sui-go-sdk/constant"
	"github.com/yasir7ca/sui-go-sdk/models"
	"github.com/yasir7ca/sui-go-sdk/signer"
	"github.com/yasir7ca/sui-go-sdk/sui"
)

var ctx = context.Background()
var cli = sui.NewSuiClient(constant.SuiPublicTestnet)

func main() {

	kmsSigner, err := signer.GetAwsSigner("alias/autodeleverager-service", "ap-northeast-1")
	fmt.Println(kmsSigner.KmsService, err)
	fmt.Println(kmsSigner.Address)

	txnMetadata, err1 := cli.SplitCoinEqual(ctx, models.SplitCoinEqualRequest{
		Signer:       "0x2365a98757b84f3335f51914b9fc4303ada752416bc3057e0930c32d859c5bf7",
		CoinObjectId: "0xd97ce4857fc751a9d82ece38157b8ed171de0339d323f110fb4e6fdf386edca1",
		SplitCount:   "2",
		Gas:          "0x0042e988b1d76930fddc65b6401e7ca2f23875f2491a0086d3d39bf5feecd20c",
		GasBudget:    "100000000",
	})
	fmt.Println(txnMetadata, err1)
	executionResult, err2 := cli.SignAndExecuteTransactionBlockWithKMS(ctx, models.SignAndExecuteTransactionBlockRequestWithKMS{
		TxnMetaData: txnMetadata,
		KeyId:       "alias/autodeleverager-service",
		PublicKey:   kmsSigner.PublicKey,
		Kms:         kmsSigner.KmsService,
		Options: models.SuiTransactionBlockOptions{
			ShowObjectChanges: true,
			ShowInput:         true,
			ShowEvents:        true,
			ShowEffects:       true,
		},
		RequestType: "WaitForLocalExecution",
	})
	fmt.Println(executionResult, err2)

	/*
		signature, err := cli.SignAndExecuteTransactionBlockWithKMS(ctx, models.SignAndExecuteTransactionBlockRequestWithKMS{
			//nolint:all
			TxnMetaData: models.TxnMetaData{
				//Gas:          100000,
				//InputObjects: []
				TxBytes: "AAACACBTyRzDtCtw6jRv1XtMXizH7Y2Mw9zn1As6Sve96LfNTAAIgJaYAAAAAAACAgABAQEAAQECAAABAAAjZamHV7hPMzX1GRS5/EMDradSQWvDBX4JMMMthZxb9wHZfOSFf8dRqdguzjgVe47Rcd4DOdMj8RD7Tm/fOG7coZ05DgAAAAAAIFCA2DxgHXRnqGpQiaDWKawbu1HSyEw1BRe7htoPmdWmI2Wph1e4TzM19RkUufxDA62nUkFrwwV+CTDDLYWcW/foAwAAAAAAAICWmAAAAAAAAA==",
			},
			Kms:       kmsSigner.KmsService,
			KeyId:     "alias/autodeleverager-service",
			PublicKey: kmsSigner.PublicKey,
			// only fetch the effects field
			Options: models.SuiTransactionBlockOptions{
				ShowObjectChanges: true,
				ShowInput:         true,
				ShowEvents:        true,
				ShowEffects:       true,
			},
			RequestType: "WaitForLocalExecution",
		})

		normalSigner, _ := signer.NewSignertWithMnemonic("cloud athlete merge matter select rib message elephant announce creek border catalog")
		fmt.Println("Public key of normal", normalSigner.PubKey)
		fmt.Println(base64.StdEncoding.EncodeToString(normalSigner.PubKey))



	*/
	/*otherSignature, jh := cli.SignAndExecuteTransactionBlock(ctx, models.SignAndExecuteTransactionBlockRequest{
		//nolint:all
		TxnMetaData: models.TxnMetaData{
			//Gas:          100000,
			//InputObjects: []
			TxBytes: "AAACACBTyRzDtCtw6jRv1XtMXizH7Y2Mw9zn1As6Sve96LfNTAAIAOH1BQAAAAACAgABAQEAAQECAAABAAAjZamHV7hPMzX1GRS5/EMDradSQWvDBX4JMMMthZxb9wHZfOSFf8dRqdguzjgVe47Rcd4DOdMj8RD7Tm/fOG7coZw5DgAAAAAAIGJdpOcubHxHNL2gwBWUd+GSgCeH/demBlsXQVePm2BTI2Wph1e4TzM19RkUufxDA62nUkFrwwV+CTDDLYWcW/foAwAAAAAAAEBCDwAAAAAAAA==",
		},
		PriKey: normalSigner.PriKey,
		// only fetch the effects field
		Options: models.SuiTransactionBlockOptions{
			ShowObjectChanges: true,
			ShowInput:         true,
			ShowEvents:        true,
			ShowEffects:       true,
		},
		RequestType: "WaitForLocalExecution",
	})
	fmt.Println(otherSignature)*/

	//fmt.Println("___)__)_)__", jh)
	//fmt.Println(signature)
	//fmt.Println(err)

}
