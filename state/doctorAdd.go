package state

import (
	"fmt"
	"errors"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/ccamaleon5/saludchain/store"
	"github.com/ccamaleon5/saludchain/transaction"

	"gopkg.in/mgo.v2/bson"

	mgo "gopkg.in/mgo.v2"
)

func checkDoctorAddTransaction(tx *transaction.Transaction, store *store.State) error {
	fmt.Println("tx",tx)
	data := &transaction.DoctorAddData{}
	if err := mapstructure.Decode(tx.Data, data); err != nil {
		return err
	}
	tx.Data = data
	if !store.HasAccount(data.IDAccount) {
		return errors.New("account doesn't exists")
	}
	data.ID = bson.ObjectIdHex(data.IDAccount)
	k, err := store.GetAccountPubKey(data.ID.Hex())
	if err != nil {
		return errors.New("pubkey can't be loaded: " + err.Error())
	}
	data.PublicKey = k.GetPubString()
	if err = tx.Verify(k); err != nil {
		return errors.New("tx can't be verified: " + err.Error())
	}
	if bson.IsObjectIdHex(data.ID.Hex()) != true {
		return errors.New("id bad data type")
	}

	if strings.TrimSpace(data.Name) == "" {
		return errors.New("name bad data type")
	}

	if strings.TrimSpace(data.LastName) == "" {
		return errors.New("lastName bad data type")
	}

	if strings.TrimSpace(data.Speciality) == "" {
		return errors.New("lastName bad data type")
	}
	return nil
}

func deliverDoctorAddTransaction(tx *transaction.Transaction, store *store.State, db *mgo.Database) error {
	data := &transaction.DoctorAddData{}
	fmt.Println("Deliver data",tx)
	if err := mapstructure.Decode(tx.Data, &data); err != nil {
		return err
	}
	return createDoctor(db, data)
}

func createDoctor(db *mgo.Database, doctor *transaction.DoctorAddData) error{
	doctor.ID = bson.ObjectIdHex(doctor.IDAccount)
	doctor.PublicKey = strings.ToUpper(doctor.PublicKey)

	dbErr := db.C("doctors").Insert(doctor)

	if dbErr != nil {
		panic(dbErr)
	}	

	return nil
}