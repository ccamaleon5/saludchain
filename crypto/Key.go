package crypto

import (
	"fmt"
	"golang.org/x/crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"errors"
)

//Key ...
type Key struct {
	pub  *ed25519.PublicKey
	priv *ed25519.PrivateKey
}

//CreateKeyPair ed25519
func CreateKeyPair() (*Key, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &Key{priv: &priv, pub: &pub}, nil
}

//NewFromStrings ...
func NewFromStrings(pub, priv string) (*Key, error) {
	k := &Key{}
	if pub == "" && priv == "" {
		return nil, errors.New("no key material supplied")
	}
	if pub != "" {
		if err := k.SetPubString(pub); err != nil {
			return nil, err
		}
	}
	if priv != "" && pub == "" {
		return nil, errors.New("no pubkey to privkey supplied")
	}
	if priv != "" {
		if err := k.SetPrivString(priv); err != nil {
			return nil, err
		}
	}
	return k, nil
}

//GetPubString ...
func (k *Key) GetPubString() string {
	return base64.StdEncoding.EncodeToString(*k.pub)
}

//GetPrivString ...
func (k *Key) GetPrivString() string {
	return base64.StdEncoding.EncodeToString(*k.priv)
}

//SetPubString ...
func (k *Key) SetPubString(pub string) error {
	bs, err := base64.StdEncoding.DecodeString(pub)
	if err != nil {
		return err
	}

	var pubKey ed25519.PublicKey = bs

	k.pub = &pubKey
	return nil
}

//SetPrivString ...
func (k *Key) SetPrivString(priv string) error {
	bs, err := base64.StdEncoding.DecodeString(priv)
	if err != nil {
		return err
	}
	
	var privKey ed25519.PrivateKey = bs

	k.priv = &privKey
	return nil
}

//Sign ...
func (k *Key) Sign(hash []byte) (string, error) {
	fmt.Println("privada",*k.priv)
	signature := ed25519.Sign(*k.priv, hash)
	
	sStr := base64.StdEncoding.EncodeToString(signature)
	return sStr, nil
}

//Verify ...
func (k *Key) Verify(hash []byte, signature string) error {
	if signature == "" {
		fmt.Println("firma malformada")
		return errors.New("malformed signature")
	}

	sBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		fmt.Println("error",err)
		return err
	}
	
	fmt.Println("publica",*k.pub)
	fmt.Println("hash",hash)
	fmt.Println("firma",sBytes)

	//if !ed25519.Verify(*k.pub, hash, sBytes) {
	//	fmt.Println("mala firma")
	//	return errors.New("bad signature")
	//}
	return nil
}