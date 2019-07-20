package transaction

import(
	"time"
	"gopkg.in/mgo.v2/bson"
)

//Permission ..
type Permission struct {
	ID         	bson.ObjectId `bson:"_id" json:"_id"`
	Doctors    	[]bson.ObjectId `bson:"doctors" json:"doctors"`
	Secret     	string `bson:"secret" json:"secret"`
}

//Record ...
type Record struct {
	Date     time.Time `json:"date"`
	Comment  string `json:"comment"`
}
