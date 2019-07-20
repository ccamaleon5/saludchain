package cmd

import (
    "fmt"
    "time"
    "golang.org/x/crypto/ed25519"
    "gopkg.in/mgo.v2/bson"
    "github.com/tendermint/tendermint/rpc/client"
    "github.com/tendermint/tendermint/types"
    cmn "github.com/tendermint/tendermint/libs/common"
    b64 "encoding/base64"
	"encoding/hex"
    "encoding/json"
    
    "github.com/ccamaleon5/saludchain/transaction"
    "github.com/ccamaleon5/saludchain/store"
    "github.com/ccamaleon5/saludchain/crypto"
)


//TxBody ...
type TxBody struct{
    Type      string         `json:"type"`
    Entity    interface{}    `json:"entity"`
}

//Key ...
type Key struct{
    PublicKey string `json:"publicKey"`
    PrivateKey string `json:"privateKey"`
}

//AttendPatientData ...
type AttendPatientData struct{
	ID          bson.ObjectId `bson:"_id" json:"_id" mapstructure:"id"`
	Record  	string		  `json:"record" mapstructure:"id"`
	PublicKey	string		  `json:"publicKey" mapstructure:"publicKey"`
}

//PatientAddData ...
type PatientAddData struct {
	ID        		bson.ObjectId 	`bson:"_id" json:"id"  mapstructure:"id"`
    IDAccount       string          `bson:"account" json:"account" mapstructure:"account"`
    Name      		string        	`bson:"name" json:"name" mapstructure:"name"`
	LastName  		string        	`bson:"lastName" json:"lastName" mapstructure:"lastName"`
    PublicKey 		string        	`bson:"publicKey" json:"publicKey" mapstructure:"publicKey"`
    MedicalRecord 	[10]string    	`bson:"medicalRecord" json:"medicalRecord" mapstructure:"medicalRecord"`
}

//DoctorAddData ...
type DoctorAddData struct{
    ID          bson.ObjectId   `bson:"id" json:"id"  mapstructure:"id"`
    IDAccount   string          `bson:"account" json:"account" mapstructure:"account"`
    Name        string          `bson:"name" json:"name" mapstructure:"name"`
    LastName    string          `bson:"lastName" json:"lastName" mapstructure:"lastName"`
    PublicKey   string          `bson:"publicKey" json:"publicKey" mapstructure:"publicKey"`
    Speciality  string          `bosn:"speciality" json:"speciality" mapstructure:"speciality"`  
}

//MedicalAppointmentAddData ...
type MedicalAppointmentAddData struct{
	ID          bson.ObjectId `bson:"_id" json:"_id" mapstructure:"id"`
	Patient     bson.ObjectId `bson:"patient" json:"patient" mapstructure:"patient"`
	Doctor      bson.ObjectId `bson:"doctor" json:"doctor" mapstructure:"doctor"`
	Position    int           `bson:"position" json:"position" mapstructure:"position"`
	Date        string        `bson:"date" json:"date" mapstructure:"date"`
	Comments    string        `bson:"comments" json:"comments" mapstructure:"comments"`
	PublicKey	string		  `json:"publicKey" mapstructure:"publicKey"`
}

//Transaction ...
type Transaction struct{
    Body       string       `json:"body"`
    Signature  string       `json:"signature"`
    PublicKey  string       `json:"publicKey"`
}

//CreateAccount ...
func CreateAccount(pubKey string){
    pub, _ := hex.DecodeString(pubKey)
    
    acc := &store.Account{ID: bson.NewObjectId().Hex(),PubKey: b64.StdEncoding.EncodeToString(pub)}
    fmt.Println("id",acc.ID)
    tx := transaction.New(transaction.AccountAdd,&transaction.AccountAddData{Account: acc})

    bs, _ := tx.ToBytes()
    
    client := client.NewHTTP("tcp://0.0.0.0:26657", "/websocket")
    result, err := client.BroadcastTxCommit(types.Tx(bs))
    if err != nil{
        panic(err)
    }

    fmt.Println(result)
}

func addSecret(accountID string, id string, value string){
    secret := &store.Secret{ID:id, Value:value, Shares: make(map[string]string),Owners: map[string]bool{accountID: true,}}
    aesKey, _ := secret.Encrypt()
    fmt.Println("aesKey",aesKey)
    tx := transaction.New(transaction.SecretAdd,&transaction.SecretAddData{Secret: secret})

    bs, _ := tx.ToBytes()
    
    client := client.NewHTTP("tcp://0.0.0.0:26657", "/websocket")
    result, err := client.BroadcastTxCommit(types.Tx(bs))
    if err != nil{
        panic(err)
    }

    fmt.Println(result)
}

//CreatePatient ...
func CreatePatient(id string, name string, lastName string, privKey string){
    var privateKey, _ = hex.DecodeString(privKey)
    publicKey := make([]byte, 32) 
    
    copy(publicKey, privateKey[32:])

    fmt.Println("privateKey:",privateKey)
    fmt.Println("publicKey:",publicKey)

    var patient PatientAddData
    patient.ID = bson.ObjectIdHex(id)
    patient.IDAccount = id
    patient.Name = name
    patient.LastName = lastName
    patient.PublicKey = b64.StdEncoding.EncodeToString(publicKey) 

    fmt.Println("patient",patient)

    tx := transaction.New(transaction.PatientAdd,&patient)

    fmt.Println("tx",tx)

    key, err1 := crypto.NewFromStrings(b64.StdEncoding.EncodeToString(publicKey),b64.StdEncoding.EncodeToString(privateKey))
    if err1 != nil{
        panic(err1)
    }
    if err := tx.Sign(key); err != nil {
		fmt.Println("no puede firmar")
    }

    bs, _ := tx.ToBytes()

    fmt.Println(patient)

    client := client.NewHTTP("tcp://0.0.0.0:26657", "/websocket")
    result, err := client.BroadcastTxCommit(types.Tx(bs))
    if err != nil{
        panic(err)
    }

    fmt.Println(result)
}

//CreateDoctor ...
func CreateDoctor(id string, name string, lastName string, speciality string, privKey string){
    var privateKey, _ = hex.DecodeString(privKey)
    publicKey := make([]byte, 32) 
    copy(publicKey, privateKey[32:])

    var doctor DoctorAddData
    doctor.ID = bson.ObjectIdHex(id)
    doctor.IDAccount = id
    doctor.Name = name
    doctor.LastName = lastName
    doctor.Speciality = speciality
    doctor.PublicKey = b64.StdEncoding.EncodeToString(publicKey) 

    fmt.Println("doctor",doctor)

    tx := transaction.New(transaction.DoctorAdd,&doctor)

    fmt.Println("tx",tx)
    fmt.Println("txBytes",tx.Data)

    key, err1 := crypto.NewFromStrings(b64.StdEncoding.EncodeToString(publicKey),b64.StdEncoding.EncodeToString(privateKey))
    if err1 != nil{
        panic(err1)
    }
    if err := tx.Sign(key); err != nil {
		fmt.Println("no puede firmar")
    }

    bs, _ := tx.ToBytes()

    fmt.Println(doctor)

    client := client.NewHTTP("tcp://0.0.0.0:26657", "/websocket")
    result, err := client.BroadcastTxCommit(types.Tx(bs))
    if err != nil{
        panic(err)
    }

    fmt.Println(result)
}

//GenerateKey ...
func GenerateKey(){
    pub, priv, err := ed25519.GenerateKey(nil)
    if err != nil {
        panic(err)
    }

    fmt.Println("public key:",string(byteToHex(pub)))
    fmt.Println("private key:",string(byteToHex(priv)))
}

func byteToHex(input []byte) string {
	var hexValue string
	for _, v := range input {
		hexValue += fmt.Sprintf("%02x", v)
	}
	return hexValue
}

//NewMedicalAppointment ...
func NewMedicalAppointment(patientID string, doctorID string, privKey string, comments string){
    var privateKey, _ = hex.DecodeString(privKey)
    publicKey := make([]byte, 32) 
    copy(publicKey, privateKey[32:])

    var medicalAppointment MedicalAppointmentAddData
    medicalAppointment.ID = bson.NewObjectId()
    medicalAppointment.Patient = bson.ObjectIdHex(patientID)
    medicalAppointment.Doctor = bson.ObjectIdHex(doctorID)
    t := time.Now()
    fmt.Println(t.Format("2006-01-02 15:04:05"))
    medicalAppointment.Date = t.Format("2006-01-02 15:04:05")
    medicalAppointment.Comments = comments
    medicalAppointment.PublicKey = b64.StdEncoding.EncodeToString(publicKey) 

    fmt.Println("patient",medicalAppointment)

    tx := transaction.New(transaction.MedicalAppointmentAdd,&medicalAppointment)

    fmt.Println("tx",tx)
    fmt.Println("txBytes",tx.Data)

    key, err1 := crypto.NewFromStrings(b64.StdEncoding.EncodeToString(publicKey),b64.StdEncoding.EncodeToString(privateKey))
    if err1 != nil{
        panic(err1)
    }
    if err := tx.Sign(key); err != nil {
		fmt.Println("no puede firmar")
    }

    bs, _ := tx.ToBytes()

    fmt.Println(medicalAppointment)

    client := client.NewHTTP("tcp://0.0.0.0:26657", "/websocket")
    result, err := client.BroadcastTxCommit(types.Tx(bs))
    if err != nil{
        panic(err)
    }

    fmt.Println(result)
}

func callRPC(txBody TxBody, pubKey string, privKey string){
    privKeyBytes, err := hex.DecodeString(privKey)
    if err != nil {
        panic(err)
    }
    
    txBodyJSON,_ := json.Marshal(txBody) 
    signature := ed25519.Sign(privKeyBytes, []byte(txBodyJSON))

    fmt.Println(string(byteToHex(signature)))

    var tx Transaction
    
    tx.Body = string(txBodyJSON)
    tx.Signature = string(byteToHex(signature))
    tx.PublicKey = pubKey
    
    fmt.Println(tx)

    message,_ := json.Marshal(tx)

    var checkTx types.Tx
    checkTx = []byte(message)  
    client := client.NewHTTP("tcp://0.0.0.0:26657", "/websocket")
    result, err := client.BroadcastTxCommit(checkTx)
    if err != nil{
        panic(err)
    }

    fmt.Println(result)
}

func getPatient(pubKey string){
    var query cmn.HexBytes
    query = []byte(pubKey)
    client := client.NewHTTP("tcp://0.0.0.0:26657","/websocket")
    result, err := client.ABCIQuery("patient",query)
    if err != nil {
        panic(err)
    }

    fmt.Println("result:",result)
}

func getDoctor(pubKey string){
    var query cmn.HexBytes
    query = []byte(pubKey)
    client := client.NewHTTP("tcp://0.0.0.0:26657","/websocket")
    result, err := client.ABCIQuery("doctor",query)
    if err != nil {
        panic(err)
    }

    fmt.Println("result:",result)
}

func getAccount(pubKey string){
    var query cmn.HexBytes
    query = []byte(pubKey)
    client := client.NewHTTP("tcp://0.0.0.0:26657","/websocket")
    result, err := client.ABCIQuery("account",query)
    if err != nil {
        panic(err)
    }

    fmt.Println("resultadooo",result)
}

//AttendPatient ...
func AttendPatient(privKey string, medicalAppointmentID string, recordMedical string){
    var privateKey, _ = hex.DecodeString(privKey)
    publicKey := make([]byte, 32) 
    copy(publicKey, privateKey[32:])

    var attendPatient AttendPatientData
    attendPatient.ID = bson.ObjectIdHex(medicalAppointmentID)
    attendPatient.PublicKey = b64.StdEncoding.EncodeToString(publicKey)
    attendPatient.Record = recordMedical 
    
    fmt.Println("attendPatient",attendPatient)

    tx := transaction.New(transaction.AttendPatient,&attendPatient)

    fmt.Println("tx",tx)

    key, err1 := crypto.NewFromStrings(b64.StdEncoding.EncodeToString(publicKey),b64.StdEncoding.EncodeToString(privateKey))
    if err1 != nil{
        panic(err1)
    }
    if err := tx.Sign(key); err != nil {
		fmt.Println("can't sign")
    }

    bs, _ := tx.ToBytes()

    client := client.NewHTTP("tcp://0.0.0.0:26657", "/websocket")
    result, err := client.BroadcastTxCommit(types.Tx(bs))
    if err != nil{
        panic(err)
    }

    fmt.Println(result)   
}