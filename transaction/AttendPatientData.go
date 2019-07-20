package transaction

import(
	"gopkg.in/mgo.v2/bson"
	
)

//AttendPatientData ...
type AttendPatientData struct{
	ID          bson.ObjectId `bson:"_id" json:"_id" mapstructure:"id"`
	Record  	string		  `json:"record" mapstructure:"id"`
	PublicKey	string		  `json:"publicKey" mapstructure:"publicKey"`
}