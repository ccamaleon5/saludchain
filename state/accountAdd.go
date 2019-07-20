package state

import (
	"errors"
    "fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/ccamaleon5/saludchain/crypto"
	"github.com/ccamaleon5/saludchain/store"
	"github.com/ccamaleon5/saludchain/transaction"
)

func checkAccountAddTransaction(tx *transaction.Transaction, store *store.State) error {
	fmt.Println("tx check:",tx)
	data := &transaction.AccountAddData{}
	if err := mapstructure.Decode(tx.Data, data); err != nil {
		return err
	}
	
	tx.Data = data
	if store.HasAccount(data.Account.ID) {
		return errors.New("account exists")
	}
	if _, err := crypto.NewFromStrings(data.Account.PubKey, ""); err != nil {
		return err
	}
	
	return nil
}

func deliverAccountAddTransaction(tx *transaction.Transaction, store *store.State) error {
	fmt.Println("tx-deliver:",tx)
	data := &transaction.AccountAddData{}
	if err := mapstructure.Decode(tx.Data, &data); err != nil {
		return err
	}
	return store.AddAccount(data.Account)
}