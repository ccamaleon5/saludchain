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

func checkMedicalAppointmentAddTransaction(tx *transaction.Transaction, db *mgo.Database) error {
	fmt.Println("tx",tx)
	var temp transaction.MedicalAppointmentAddData
	var message map[string]interface{}

	message = tx.Data.(map[string]interface{})
	
	temp.ID = bson.ObjectIdHex(message["_id"].(string))
	temp.Patient = bson.ObjectIdHex(message["patient"].(string))
	temp.Doctor = bson.ObjectIdHex(message["doctor"].(string))
	temp.Position = int(message["position"].(float64))
	temp.Date = message["date"].(string)
	temp.PublicKey = message["publicKey"].(string)
	temp.Comments = message["comments"].(string)

	tx.Data = temp
	
	pub, _:= b64.StdEncoding.DecodeString(message["publicKey"].(string))
	
	key, err1 := crypto.NewFromStrings(b64.StdEncoding.EncodeToString(pub),"")
    if err1 != nil{
        panic(err1)
    }
	if err := tx.Verify(key); err != nil {
		return errors.New("tx can't be verified: " + err.Error())
	}
	countPatient, _ := db.C("patients").Find(bson.M{"publicKey": strings.ToUpper(message["publicKey"].(string))}).Count()
	countDoctor, _ := db.C("doctors").Find(bson.M{"publicKey": strings.ToUpper(message["publicKey"].(string))}).Count()
		
	if countPatient == 0 && countDoctor == 0 {
		return errors.New("publicKey doesn't exist: ")
	}
	return nil
}

func deliverMedicalAppointmentAddTransaction(tx *transaction.Transaction, store *store.State, db *mgo.Database) error {
	fmt.Println("Deliver data",tx)
	data := &transaction.MedicalAppointmentAddData{}
	var message map[string]interface{}

	message = tx.Data.(map[string]interface{})

	data.ID = bson.ObjectIdHex(message["_id"].(string))
	data.Patient = bson.ObjectIdHex(message["patient"].(string))
	data.Doctor = bson.ObjectIdHex(message["doctor"].(string))
	data.Position = int(message["position"].(float64))
	data.Date = message["date"].(string)
	data.PublicKey = message["publicKey"].(string)
	data.Comments = message["comments"].(string)

	return createMedicalAppointment(db, data)
}

func createMedicalAppointment(db *mgo.Database, medicalAppointment *transaction.MedicalAppointmentAddData) error{
	var patient = transaction.PatientAddData{}
	err := db.C("patients").Find(bson.M{"publicKey": strings.ToUpper(medicalAppointment.PublicKey)}).One(&patient)
	if err != nil {
		panic(err)
	}
	medicalAppointment.Patient = patient.ID

	var doctor = transaction.DoctorAddData{}
	err = db.C("doctors").Find(bson.M{"_id": medicalAppointment.Doctor}).One(&doctor)
	if err != nil {
		panic(err)
	}

	fmt.Println("doctor.ID:",doctor.ID)

	count, _ := db.C("medicalappointments").Find(bson.M{"doctor": doctor.ID}).Count()

	fmt.Println("count:",count)

	if count > 0 {
		medicalAppointment.Position = count + 1
	} else {
		medicalAppointment.Position = 1
	}

	dbErr := db.C("medicalappointments").Insert(medicalAppointment)

	if dbErr != nil {
		panic(dbErr)
	}

	return nil
}