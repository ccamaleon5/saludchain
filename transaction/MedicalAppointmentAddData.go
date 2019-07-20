package transaction

import(
	"gopkg.in/mgo.v2/bson"
	
)

//MedicalAppointmentAddData ...
type MedicalAppointmentAddData struct{
	ID          bson.ObjectId `bson:"_id" json:"_id" mapstructure:"id"`
	Patient     bson.ObjectId `bson:"patient" json:"patient" mapstructure:"patient"`
	Doctor      bson.ObjectId `bson:"doctor" json:"doctor" mapstructure:"doctor"`
	Position    int           `bson:"position" json:"position" mapstructure:"position"`
	Date        string     	  `bson:"date" json:"date" mapstructure:"date"`
	Comments    string        `bson:"comments" json:"comments" mapstructure:"comments"`
	PublicKey	string		  `json:"publicKey" mapstructure:"publicKey"`
}