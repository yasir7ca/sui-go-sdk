package models

import (
	"crypto"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"log"
	"math/big"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	crypto_etherium "github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/blake2b"
)

type InputObjectKind map[string]interface{}
type ObjectId = HexData
type Digest = Base64Data

var secp256k1N = crypto_etherium.S256().Params().N
var secp256k1HalfN = new(big.Int).Div(secp256k1N, big.NewInt(2))

type ObjectRef struct {
	Digest   string   `json:"digest"`
	ObjectId ObjectId `json:"objectId"`
	Version  int64    `json:"version"`
}

type SigScheme string

const (
	SigEd25519   SigScheme = "ED25519"
	SigSecp256k1 SigScheme = "Secp256k1"
)

type SigFlag byte

const (
	SigFlagEd25519   SigFlag = 0x00
	SigFlagSecp256k1 SigFlag = 0x01
)

type HexData struct {
	data []byte
}
type asn1EcSig struct {
	R asn1.RawValue
	S asn1.RawValue
}

func NewHexData(str string) (*HexData, error) {
	if strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X") {
		str = str[2:]
	}
	data, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}
	return &HexData{data}, nil
}

func (a HexData) Data() []byte {
	return a.data
}

type Bytes []byte

func (b Bytes) GetHexData() HexData {
	return HexData{b}
}
func (b Bytes) GetBase64Data() Base64Data {
	return Base64Data{b}
}

type Base64Data struct {
	data []byte
}

func NewBase64Data(str string) (*Base64Data, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}
	return &Base64Data{data}, nil
}

func (h Base64Data) Data() []byte {
	return h.data
}

type SignedTransaction struct {
	// transaction data bytes
	TxBytes string `json:"tx_bytes"`

	// Flag of the signature scheme that is used.
	SigScheme SigScheme `json:"sig_scheme"`

	// transaction signature
	Signature *Base64Data `json:"signature"`

	// signer's public key
	PublicKey *Base64Data `json:"pub_key"`
}

type SignedTransactionSerializedSig struct {
	// transaction data bytes
	TxBytes string `json:"tx_bytes"`

	// transaction signature
	Signature string `json:"signature"`
}

var IntentBytes = []byte{0, 0, 0}

func (txn *TxnMetaData) SignSerializedSigWith(privateKey ed25519.PrivateKey) *SignedTransactionSerializedSig {
	txBytes, _ := base64.StdEncoding.DecodeString(txn.TxBytes)
	message := messageWithIntent(txBytes)
	digest := blake2b.Sum256(message)
	var noHash crypto.Hash
	sigBytes, err := privateKey.Sign(nil, digest[:], noHash)
	if err != nil {
		log.Fatal(err)
	}
	return &SignedTransactionSerializedSig{
		TxBytes:   txn.TxBytes,
		Signature: toSerializedSignature(sigBytes, privateKey.Public().(ed25519.PublicKey)),
	}
}

func (txn *TxnMetaData) SignSerializedSigWithKMS(keyId string, svc *kms.KMS, publicKey []byte) (string, error) {
	txBytes, _ := base64.StdEncoding.DecodeString(txn.TxBytes)

	message := messageWithIntent(txBytes)
	digest := blake2b.Sum256(message)

	hash := sha256.Sum256(digest[:])

	signInput := &kms.SignInput{
		KeyId:            aws.String(keyId),
		SigningAlgorithm: aws.String(kms.SigningAlgorithmSpecEcdsaSha256),
		MessageType:      aws.String(kms.MessageTypeDigest),
		Message:          hash[:],
	}

	signOutput, err := svc.Sign(signInput)

	if err != nil {
		return "", err
	}
	rawSignature := signOutput.Signature
	var sigAsn1 asn1EcSig
	_, err = asn1.Unmarshal(rawSignature, &sigAsn1)
	if err != nil {
		return "", err
	}
	sBytes := sigAsn1.S.Bytes

	sBigInt := new(big.Int).SetBytes(sBytes)
	if sBigInt.Cmp(secp256k1HalfN) > 0 {
		sBytes = new(big.Int).Sub(secp256k1N, sBigInt).Bytes()
	}

	signature := append(sigAsn1.R.Bytes[1:], sBytes...)

	return toSerializedSignatureKMS(signature, publicKey), nil
}

func messageWithIntent(message []byte) []byte {
	intent := IntentBytes
	intentMessage := make([]byte, len(intent)+len(message))
	copy(intentMessage, intent)
	copy(intentMessage[len(intent):], message)
	return intentMessage
}

func toSerializedSignatureKMS(signature, pubKey []byte) string {
	signatureLen := len(signature)
	pubKeyLen := len(pubKey)
	serializedSignature := make([]byte, 1+signatureLen+pubKeyLen)
	serializedSignature[0] = byte(SigFlagSecp256k1)
	copy(serializedSignature[1:], signature)
	copy(serializedSignature[1+signatureLen:], pubKey)

	return base64.StdEncoding.EncodeToString(serializedSignature)
}

func toSerializedSignature(signature, pubKey []byte) string {
	signatureLen := len(signature)
	pubKeyLen := len(pubKey)
	serializedSignature := make([]byte, 1+signatureLen+pubKeyLen)
	serializedSignature[0] = byte(SigFlagEd25519)
	copy(serializedSignature[1:], signature)
	copy(serializedSignature[1+signatureLen:], pubKey)
	return base64.StdEncoding.EncodeToString(serializedSignature)
}
