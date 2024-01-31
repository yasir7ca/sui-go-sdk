package signer

import (
	"crypto/ed25519"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/tyler-smith/go-bip39"
	"github.com/yasir7ca/sui-go-sdk/common/keypair"
	"golang.org/x/crypto/blake2b"
)

var secp256k1N = crypto.S256().Params().N
var secp256k1HalfN = new(big.Int).Div(secp256k1N, big.NewInt(2))

const awsKmsSignOperationMessageType = "DIGEST"
const awsKmsSignOperationSigningAlgorithm = "ECDSA_SHA_256"
const (
	SigntureFlagEd25519     = 0x0
	SigntureFlagSecp256k1   = 0x1
	AddressLength           = 64
	DerivationPathEd25519   = `m/44'/784'/0'/0'/0'`
	DerivationPathSecp256k1 = `m/54'/784'/0'/0/0`
)

type PublicKey []byte
type Signer struct {
	PriKey  ed25519.PrivateKey
	PubKey  ed25519.PublicKey
	Address string
}

type AwsSigner struct {
	KmsService      *kms.KMS
	PublicKey       []byte
	PublicKeyStr    string
	PublicKeyBase64 string
	Address         string
	KmsId           string
}

type ECDSAPubKeyAlgorithm struct {
	Algorithm  asn1.ObjectIdentifier
	Parameters asn1.ObjectIdentifier
}

type EcdsaPubKey struct {
	Algo   ECDSAPubKeyAlgorithm
	PubKey asn1.BitString
}
type asn1EcPublicKey struct {
	EcPublicKeyInfo asn1EcPublicKeyInfo
	PublicKey       asn1.BitString
}

type asn1EcPublicKeyInfo struct {
	Algorithm  asn1.ObjectIdentifier
	Parameters asn1.ObjectIdentifier
}

type asn1EcSig struct {
	R asn1.RawValue
	S asn1.RawValue
}

func NewSigner(seed []byte) *Signer {
	priKey := ed25519.NewKeyFromSeed(seed[:])
	pubKey := priKey.Public().(ed25519.PublicKey)

	tmp := []byte{byte(keypair.Ed25519Flag)}
	tmp = append(tmp, pubKey...)
	addrBytes := blake2b.Sum256(tmp)
	addr := "0x" + hex.EncodeToString(addrBytes[:])[:AddressLength]

	return &Signer{
		PriKey:  priKey,
		PubKey:  pubKey,
		Address: addr,
	}
}

func NewSignertWithMnemonic(mnemonic string) (*Signer, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}
	key, err := DeriveForPath("m/44'/784'/0'/0'/0'", seed)
	if err != nil {
		return nil, err
	}
	return NewSigner(key.Key), nil
}

func GetAwsSigner(KeyId string, Region string) (*AwsSigner, error) {

	config := aws.NewConfig()
	config.Region = &Region

	session, sessionErr := session.NewSession(config)

	if sessionErr != nil {
		if aerr, ok := sessionErr.(awserr.Error); ok {
			fmt.Println("Error", aerr.Code(), aerr.Error())
			return nil, errors.New(strings.Join([]string{aerr.Code(), aerr.Error()}, "-"))
		} else {
			return nil, errors.New(sessionErr.Error())
		}
	}

	kmsSvc := kms.New(session)

	input := &kms.GetPublicKeyInput{
		KeyId: aws.String(KeyId),
	}

	pubKey, err := kmsSvc.GetPublicKey(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println("Error", aerr.Error())
			return nil, errors.New(strings.Join([]string{aerr.Code(), aerr.Error()}, "-"))
		} else {
			return nil, errors.New(err.Error())
		}
	}

	rawPK := pubKey.PublicKey

	var asn1pubk asn1EcPublicKey
	_, err = asn1.Unmarshal(rawPK, &asn1pubk)
	if err != nil {
		return nil, error(err)
	}

	pubkey, err := crypto.UnmarshalPubkey(asn1pubk.PublicKey.Bytes)

	if err != nil {
		return nil, error(err)
	}

	compressedPubkey := secp256k1.CompressPubkey(pubkey.X, pubkey.Y)

	tmp := []byte{byte(keypair.Secp256k1Flag)}
	tmp = append(tmp, compressedPubkey...)
	addrBytes := blake2b.Sum256(tmp)
	addr := "0x" + hex.EncodeToString(addrBytes[:])[:64]
	publicKeyStr := hex.EncodeToString(compressedPubkey)
	publicKeyBase64 := base64.StdEncoding.EncodeToString(compressedPubkey)

	signer := AwsSigner{KmsService: kmsSvc, Address: addr, PublicKey: compressedPubkey, PublicKeyStr: publicKeyStr,
		PublicKeyBase64: publicKeyBase64, KmsId: KeyId,
	}

	return &signer, nil

}
