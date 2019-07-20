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

func checkPatientAddTransaction(tx *transaction.Transaction, store *store.State) error {
	fmt.Println("check-tx",tx)
	data := &transaction.PatientAddData{}
	if err := mapstructure.Decode(tx.Data, data); err != nil {
		return err
	}
	tx.Data = data
	if !store.HasAccount(data.IDAccount) {
		fmt.Println("no tiene cuenta")
		return errors.New("account doesn't exists")
	}
	data.ID = bson.ObjectIdHex(data.IDAccount)
	k, err := store.GetAccountPubKey(data.ID.Hex())
	if err != nil {
		fmt.Println("llave publica no puede ser cargada")
		return errors.New("pubkey can't be loaded: " + err.Error())
	}
	data.PublicKey = k.GetPubString()
	if err = tx.Verify(k); err != nil {
		fmt.Println("tx no puede ser verificada")
		return errors.New("tx can't be verified: " + err.Error())
	}
	if bson.IsObjectIdHex(data.ID.Hex()) != true {
		fmt.Println("mal tipo de id")
		return errors.New("id bad data type")
	}

	if strings.TrimSpace(data.Name) == "" {
		fmt.Println("mal tipo de nombre")
		return errors.New("name bad data type")
	}

	if strings.TrimSpace(data.LastName) == "" {
		fmt.Println("mal tipo de apellido")
		return errors.New("lastName bad data type")
	}
return nil
}

func deliverPatientAddTransaction(tx *transaction.Transaction, store *store.State, db *mgo.Database) error {
	data := &transaction.PatientAddData{}
	fmt.Println("Deliver-tx",tx)
	if err := mapstructure.Decode(tx.Data, &data); err != nil {
		return err
	}
	return createPatient(db, data)
}

func createPatient(db *mgo.Database, patient *transaction.PatientAddData) error{
	patient.ID = bson.ObjectIdHex(patient.IDAccount)
	patient.PublicKey = strings.ToUpper(patient.PublicKey)
	
	dbErr := db.C("patients").Insert(patient)

	if dbErr != nil {
		panic(dbErr)
	}	

	return nil
}