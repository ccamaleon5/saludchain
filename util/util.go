package util

import(
	"fmt"
	"time"
	"strconv"
	mgo "gopkg.in/mgo.v2"
)

//byteToHex ...
func byteToHex(input []byte) string {
	var hexValue string
	for _, v := range input {
		hexValue += fmt.Sprintf("%02x", v)
	}
	return hexValue
}

//FindTotalDocuments into mongodb
func FindTotalDocuments(db *mgo.Database) int64 {
	collections := [4]string{"medicalappointments", "patients", "doctors", "permission"}
	var sum int64

	for _, collection := range collections {
		count, _ := db.C(collection).Find(nil).Count()
		sum += int64(count)
	}

	return sum
}

// FindTimeFromObjectID ... Convert ObjectID string to Time
func FindTimeFromObjectID(id string) time.Time {
	ts, _ := strconv.ParseInt(id[0:8], 16, 64)
	return time.Unix(ts, 0)
}