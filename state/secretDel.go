package state

import (
	"errors"

	"github.com/mitchellh/mapstructure"
	"github.com/ccamaleon5/saludchain/store"
	"github.com/ccamaleon5/saludchain/transaction"
)

func checkSecretDelTransaction(tx *transaction.Transaction, store *store.State) error {
	data := &transaction.SecretDelData{}
	if err := mapstructure.Decode(tx.Data, data); err != nil {
		return err
	}
	tx.Data = data
	if !store.HasSecret(data.ID) {
		return errors.New("secret doesn't exists")
	}
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

func deliverSecretDelTransaction(tx *transaction.Transaction, store *store.State) error {
	data := &transaction.SecretDelData{}
	if err := mapstructure.Decode(tx.Data, &data); err != nil {
		return err
	}
	return store.DeleteSecret(data.ID)
}