package transaction

import(
	"gopkg.in/mgo.v2/bson"
)

//PatientAddData ...
type PatientAddData struct {
	ID        		bson.ObjectId 	`bson:"_id" json:"id"  mapstructure:"id"`
	IDAccount       string          `bson:"account" json:"account" mapstructure:"account"`
	Name      		string        	`bson:"name" json:"name" mapstructure:"name"`
	LastName  		string        	`bson:"lastName" json:"lastName" mapstructure:"lastName"`
	PublicKey 		string        	`bson:"publicKey" json:"publicKey" mapstructure:"publicKey"`
	HashCredential  string          `bson:"hashCredential" json:"hashCredential" mapstructure:"hashCredential"`
    MedicalRecord 	[10]string    	`bson:"medicalRecord" json:"medicalRecord" mapstructure:"medicalRecord"`
}