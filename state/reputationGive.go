package state

import (
	"errors"

	"github.com/mitchellh/mapstructure"
	"github.com/ccamaleon5/saludchain/store"
	"github.com/ccamaleon5/saludchain/transaction"
)

func checkReputationGiveTransaction(tx *transaction.Transaction, store *store.State) error {
	data := &transaction.ReputationGiveData{}
	if err := mapstructure.Decode(tx.Data, data); err != nil {
		return err
	}
	tx.Data = data
	if !store.HasAccount(data.From) {
		return errors.New("reject give-rep because id doesnt exist: " + data.From)
	}
	if !store.HasAccount(data.To) {
		return errors.New("reject give-rep because id doesnt exist: " + data.From)
	}
	if data.Value < -3 || data.Value > 3 {
		return errors.New("reject give-rep because bad value")
	}
	k, err := store.GetAccountPubKey(data.From)
	if err != nil {
		return errors.New("reject give-rep because pubkey cant be loaded")
	}
	if err = tx.Verify(k); err != nil {
		return errors.New("reject give-rep because signature cant be verified")
	}
	if err := tx.VerifyProofOfWork(transaction.DefaultProofOfWorkCost); err != nil {
		return err
	}
	return nil
}

func deliverReputationGiveTransaction(tx *transaction.Transaction, store *store.State) error {
	data := &transaction.ReputationGiveData{}
	if err := mapstructure.Decode(tx.Data, &data); err != nil {
		return err
	}
	acc, err := store.GetAccount(data.To)
	if err != nil {
		return err
	}
	
	err = store.SetAccount(acc)
	if err != nil {
		return err
	}
	return nil
}