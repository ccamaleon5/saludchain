package state

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/ccamaleon5/saludchain/code"
	"github.com/ccamaleon5/saludchain/store"
	"github.com/ccamaleon5/saludchain/transaction"
	"github.com/ccamaleon5/saludchain/util"

	"github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/iavl"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ types.Application = (*JSONStateApplication)(nil)

// JSONStateApplication ...
type JSONStateApplication struct {
	types.BaseApplication
	store 					*store.State
	db						*mgo.Database
}

// NewJSONStateApplication ...
func NewJSONStateApplication(dbCopy *mgo.Database) *JSONStateApplication {
	tree := iavl.NewMutableTree(db.NewMemDB(),0)
	return &JSONStateApplication{store: store.NewStateFromTree(*tree),db: dbCopy}
}

// Info ...
func (app *JSONStateApplication) Info(req types.RequestInfo) (resInfo types.ResponseInfo) {
	return types.ResponseInfo{Data: fmt.Sprintf("{\"size\":%v}", app.store.Tree.Size())}
}

// DeliverTx ... Update the global state
func (app *JSONStateApplication) DeliverTx(txBytes []byte) types.ResponseDeliverTx {
	tx := &transaction.Transaction{}
	if err := tx.FromBytes(txBytes); err != nil {
		return types.ResponseDeliverTx{Code: code.CodeTypeBadData, Tags: nil}
	}
	switch tx.Type {
	case transaction.AccountAdd:
		{
			if err := deliverAccountAddTransaction(tx, app.store); err != nil {
				return types.ResponseDeliverTx{Code: code.CodeTypeBadData, Tags: nil, Log:err.Error()}
			}
		}

	case transaction.AccountDel:
		{
			if err := deliverAccountDelTransaction(tx, app.store); err != nil {
				return types.ResponseDeliverTx{Code: code.CodeTypeBadData, Tags: nil, Log:err.Error()}
			}
		}
	case transaction.AttendPatient:
		{
			if err := deliverAttendPatientTransaction(tx, app.store, app.db); err != nil {
				return types.ResponseDeliverTx{Code: code.CodeTypeBadData, Tags: nil, Log:err.Error()}
			}
		}	
	case transaction.SecretAdd:
		{
			if err := deliverSecretAddTransaction(tx, app.store); err != nil {
				return types.ResponseDeliverTx{Code: code.CodeTypeBadData, Tags: nil, Log:err.Error()}
			}
		}
	case transaction.SecretUpdate:
		{
			if err := deliverSecretUpdateTransaction(tx, app.store); err != nil {
				return types.ResponseDeliverTx{Code: code.CodeTypeBadData, Tags: nil, Log:err.Error()}
			}
		}
	case transaction.SecretDel:
		{
			if err := deliverSecretDelTransaction(tx, app.store); err != nil {
				return types.ResponseDeliverTx{Code: code.CodeTypeBadData, Tags: nil, Log:err.Error()}
			}
		}
	case transaction.SecretShare:
		{
			if err := deliverSecretShareTransaction(tx, app.store); err != nil {
				return types.ResponseDeliverTx{Code: code.CodeTypeBadData, Tags: nil, Log:err.Error()}
			}
		}
	case transaction.PatientAdd:
		{
			if err := deliverPatientAddTransaction(tx, app.store, app.db); err != nil{
				return types.ResponseDeliverTx{Code: code.CodeTypeBadData, Tags: nil, Log:err.Error()}
			}
		}
	case transaction.DoctorAdd:
		{
			if err := deliverDoctorAddTransaction(tx, app.store, app.db); err != nil{
				return types.ResponseDeliverTx{Code: code.CodeTypeBadData, Tags: nil, Log:err.Error()}
			}
		}
	case transaction.MedicalAppointmentAdd:
		{
			if err := deliverMedicalAppointmentAddTransaction(tx, app.store, app.db); err != nil{
				return types.ResponseDeliverTx{Code: code.CodeTypeBadData, Tags: nil, Log:err.Error()}
			}
		}	
	}
	
	return types.ResponseDeliverTx{Code: code.CodeTypeOK, Tags: nil}
}

// CheckTx ... Verify the transaction
func (app *JSONStateApplication) CheckTx(txBytes []byte) types.ResponseCheckTx {
	tx := &transaction.Transaction{}
	if err := tx.FromBytes(txBytes); err != nil{
		fmt.Println(err)
		return types.ResponseCheckTx{Code: code.CodeTypeBadData}
	}
	
	switch tx.Type{
		case transaction.AccountAdd:
			{
				if err := checkAccountAddTransaction(tx, app.store); err != nil {
					return types.ResponseCheckTx{Code: code.CodeTypeBadData, Log:err.Error()}
				}
			}

		case transaction.AccountDel:
			{
				if err := checkAccountDelTransaction(tx, app.store); err != nil {
					return types.ResponseCheckTx{Code: code.CodeTypeBadData, Log:err.Error()}
				}
			}
		case transaction.AttendPatient:
			{
				if err := checkAttendPatientTransaction(tx, app.db); err != nil {
					return types.ResponseCheckTx{Code: code.CodeTypeBadData, Log:err.Error()}
				}
			}
		case transaction.SecretAdd:
			{
				if err := checkSecretAddTransaction(tx, app.store); err != nil {
					return types.ResponseCheckTx{Code: code.CodeTypeBadData, Log:err.Error()}
				}
			}
		case transaction.SecretUpdate:
			{
			if err := checkSecretUpdateTransaction(tx, app.store); err != nil {
					return types.ResponseCheckTx{Code: code.CodeTypeBadData, Log:err.Error()}
				}
			}
		case transaction.SecretDel:
			{
				if err := checkSecretDelTransaction(tx, app.store); err != nil {
					return types.ResponseCheckTx{Code: code.CodeTypeBadData, Log:err.Error()}
				}
			}
		case transaction.SecretShare:
			{
				if err := checkSecretShareTransaction(tx, app.store); err != nil {
					return types.ResponseCheckTx{Code: code.CodeTypeBadData, Log:err.Error()}
				}
			}
		case transaction.PatientAdd:
			{
				if err := checkPatientAddTransaction(tx, app.store); err != nil {
					return types.ResponseCheckTx{Code: code.CodeTypeBadData, Log:err.Error()}
				}
			}
		case transaction.DoctorAdd:
			{
				if err := checkDoctorAddTransaction(tx, app.store); err != nil {
					return types.ResponseCheckTx{Code: code.CodeTypeBadData, Log:err.Error()}
				}
			}
		case transaction.MedicalAppointmentAdd:
			{
				if err := checkMedicalAppointmentAddTransaction(tx, app.db); err != nil {
					return types.ResponseCheckTx{Code: code.CodeTypeBadData, Log:err.Error()}
				}
			}
	}
	codeType := code.CodeTypeOK

	return types.ResponseCheckTx{Code: codeType}
}

// Commit ...Commit the block. Calculate the appHash
func (app *JSONStateApplication) Commit() types.ResponseCommit {
	fmt.Println("BLOCK COMMMITED")
	appHash := make([]byte, 8)

	count := util.FindTotalDocuments(app.db)

	binary.PutVarint(appHash, count)

	return types.ResponseCommit{Data: appHash}
}

// Query ... Query the blockchain.
func (app *JSONStateApplication) Query(reqQuery types.RequestQuery) (resQuery types.ResponseQuery) {
	switch reqQuery.Path{
	case "patient":
		return getPatient(reqQuery.Data, app.db)
	case "doctor":
		return getDoctor(reqQuery.Data, app.db)
	case "account":
		{
			var (
				result interface{}
				err    error
			)
			if reqQuery.Data == nil {
				result, err = app.store.ListAccounts()
				log.Printf("got account list: %+v", result)
			} else {
				result, err = app.store.GetAccount(string(reqQuery.Data))
				log.Printf("got account: %+v", result)
			}
			if err != nil {
				resQuery.Code = code.CodeTypeBadData
				resQuery.Log = err.Error()
				return
			}
			bs, _ := json.Marshal(result)
			resQuery.Value = bs
		}
	case "secret":
		{
			var (
				result interface{}
				err    error
			)
			if reqQuery.Data == nil {
				result, err = app.store.ListSecrets()
				log.Printf("got secret list: %+v", result)
			} else {
				result, err = app.store.GetSecret(string(reqQuery.Data))
				log.Printf("got secret: %+v", result)
			}
			if err != nil {
				resQuery.Code = code.CodeTypeBadData
				resQuery.Log = err.Error()
				return
			}
			bs, _ := json.Marshal(result)
			resQuery.Value = bs
		}
	default:
		{
			resQuery.Code = code.CodeTypeBadData
			resQuery.Log = "wrong path"
			return
		}
	}

	return
}

func getPatient(data []byte, db *mgo.Database) types.ResponseQuery {
	id := strings.ToUpper(string(data))
	patient := transaction.PatientAddData{} 
	error := db.C("patients").Find(bson.M{"publicKey": id}).One(&patient)
	if error != nil{
		panic(error)
	}

	response, _ := json.Marshal(patient)

	return types.ResponseQuery{Value: response}
}

func getDoctor(data []byte, db *mgo.Database) types.ResponseQuery {
	id := strings.ToUpper(string(data))
	doctor := transaction.DoctorAddData{} 
	error := db.C("doctors").Find(bson.M{"publicKey": id}).One(&doctor)
	if error != nil{
		panic(error)
	}
	
	response, _ := json.Marshal(doctor)

	return types.ResponseQuery{Value: response}
}