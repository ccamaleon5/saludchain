package state

import (
	"fmt"
	"errors"

	"github.com/mitchellh/mapstructure"
	"github.com/ccamaleon5/saludchain/store"
	"github.com/ccamaleon5/saludchain/transaction"
)

func checkSecretAddTransaction(tx *transaction.Transaction, store *store.State) error {
	fmt.Println("tx",tx)
	data := &transaction.SecretAddData{}
	if err := mapstructure.Decode(tx.Data, data); err != nil {
		return err
	}
	tx.Data = data
	if store.HasSecret(data.Secret.ID) {
		fmt.Println("secret exists")
		return errors.New("secret exists")
	}
	if len(data.Secret.Shares) == 0 {
		fmt.Println("no shares supplied")
		return errors.New("no shares supplied")
	}
	if len(data.Secret.Owners) == 0 {
		fmt.Println("no owners supplied")
		return errors.New("no owners supplied")
	}
	
	return nil
}

func deliverSecretAddTransaction(tx *transaction.Transaction, store *store.State) error {
	data := &transaction.SecretAddData{}
	if err := mapstructure.Decode(tx.Data, &data); err != nil {
		return err
	}
	return store.AddSecret(data.Secret)
}