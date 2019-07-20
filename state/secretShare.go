package state

import (
	"errors"

	"github.com/mitchellh/mapstructure"
	"github.com/ccamaleon5/saludchain/store"
	"github.com/ccamaleon5/saludchain/transaction"
)

func checkSecretShareTransaction(tx *transaction.Transaction, store *store.State) error {
	data := &transaction.SecretShareData{}
	if err := mapstructure.Decode(tx.Data, data); err != nil {
		return err
	}
	tx.Data = data
	secret, err := store.GetSecret(data.ID)
	if err != nil {
		return err
	}
	if _, ok := secret.Shares[data.SenderID]; !ok {
		return errors.New("sender has no share on this secret")
	}
	if _, ok := secret.Owners[data.SenderID]; !ok {
		return errors.New("sender is not owner of this secret")
	}
	if _, ok := secret.Shares[data.AccountID]; ok {
		return errors.New("share receiver already has a share on this secret")
	}
	k, err := store.GetAccountPubKey(data.SenderID)
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

func deliverSecretShareTransaction(tx *transaction.Transaction, store *store.State) error {
	data := &transaction.SecretShareData{}
	if err := mapstructure.Decode(tx.Data, &data); err != nil {
		return err
	}
	secret, err := store.GetSecret(data.ID)
	if err != nil {
		return err
	}
	secret.Shares[data.AccountID] = data.Key
	return store.SetSecret(secret)
}