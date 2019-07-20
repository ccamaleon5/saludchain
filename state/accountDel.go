package state

import (
	"errors"

	"github.com/mitchellh/mapstructure"
	"github.com/ccamaleon5/saludchain/store"
	"github.com/ccamaleon5/saludchain/transaction"
)

func checkAccountDelTransaction(tx *transaction.Transaction, store *store.State) error {
	data := &transaction.AccountDelData{}
	if err := mapstructure.Decode(tx.Data, data); err != nil {
		return err
	}
	tx.Data = data
	if !store.HasAccount(data.ID) {
		return errors.New("account doesn't exists")
	}
	k, err := store.GetAccountPubKey(data.ID)
	if err != nil {
		return errors.New("pubkey can't be loaded: " + err.Error())
	}
	if err = tx.Verify(k); err != nil {
		return errors.New("tx can't be verified: " + err.Error())
	}
	if err := tx.VerifyProofOfWork(transaction.DefaultProofOfWorkCost); err != nil {
		return err
	}
	return nil
}

func deliverAccountDelTransaction(tx *transaction.Transaction, store *store.State) error {
	data := &transaction.AccountDelData{}
	if err := mapstructure.Decode(tx.Data, &data); err != nil {
		return err
	}
	return store.DeleteAccount(data.ID)
}