package state

import (
	"fmt"
	"errors"
	"strings"
	
	b64 "encoding/base64"
    "gopkg.in/mgo.v2/bson"
	mgo "gopkg.in/mgo.v2"

	"github.com/ccamaleon5/saludchain/store"
	"github.com/ccamaleon5/saludchain/transaction"
	"github.com/ccamaleon5/saludchain/crypto"
)

func checkAttendPatientTransaction(tx *transaction.Transaction, db *mgo.Database) error {
	fmt.Println("tx",tx)
	var temp transaction.AttendPatientData
	var message map[string]interface{}

	message = tx.Data.(map[string]interface{})
	
	temp.ID = bson.ObjectIdHex(message["_id"].(string))
	temp.PublicKey = message["publicKey"].(string) 
	temp.Record = message["record"].(string)

	tx.Data = temp
	
	pub, _:= b64.StdEncoding.DecodeString(message["publicKey"].(string))
	
	key, err1 := crypto.NewFromStrings(b64.StdEncoding.EncodeToString(pub),"")
    if err1 != nil{
        panic(err1)
    }
	if err := tx.Verify(key); err != nil {
		return errors.New("tx can't be verified: " + err.Error())
	}
	countDoctor, _ := db.C("doctors").Find(bson.M{"publicKey": strings.ToUpper(message["publicKey"].(string))}).Count()
		
	if countDoctor == 0 {
		return errors.New("publicKey doesn't exist: ")
	}

	if strings.TrimSpace(message["record"].(string)) == "" {
		return errors.New("record bad data type")
	}

	if countMedicalRecord, _ := db.C("medicalappointments").Find(bson.M{"_id":temp.ID}).Count(); countMedicalRecord == 0{
		return errors.New("MedicalAppointment doesn't exist: ")
	}

	return nil
}

func deliverAttendPatientTransaction(tx *transaction.Transaction, store *store.State, db *mgo.Database) error {
	fmt.Println("Deliver data",tx)
	data := &transaction.AttendPatientData{}
	var message map[string]interface{}

	message = tx.Data.(map[string]interface{})

	data.ID = bson.ObjectIdHex(message["_id"].(string))
	data.PublicKey = message["publicKey"].(string) 
	data.Record = message["record"].(string)

	return attendPatient(db, data)
}

func attendPatient(db *mgo.Database, attendPatient *transaction.AttendPatientData) error{
	var medicalAppointment transaction.MedicalAppointmentAddData
	
	if err := db.C("medicalappointments").Find(bson.M{"_id": attendPatient.ID}).One(&medicalAppointment); err != nil{
		panic(err)
	}

	var patient = transaction.PatientAddData{}
	err1 := db.C("patients").Find(bson.M{"_id": medicalAppointment.Patient}).One(&patient)
	if err1 != nil {
		panic(err1)
	}

	slice := patient.MedicalRecord[0:9]

	for  i := range slice {
		if slice[i] == ""{
			slice[i] = attendPatient.Record
			break
		}	
	}

	dbErr := db.C("patients").Update(bson.M{"_id": medicalAppointment.Patient}, patient)

	if dbErr != nil {
		panic(dbErr)
	}

	return nil
}