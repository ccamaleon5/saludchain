package transaction

import (
	"fmt"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"sort"
	"time"

	"github.com/ccamaleon5/saludchain/crypto"

	"golang.org/x/crypto/sha3"
)

//Transaction ...
type Transaction struct {
	Type      TransactionType `json:"type"`
	Timestamp time.Time       `json:"timestamp"`
	Signature string          `json:"signature"`
	Nonce     uint32          `json:"nonce"`
	Data      interface{}     `json:"data"`
}

//Hashable ...
type Hashable interface {
	Hash() []byte
}

type TransactionType string

const (
	AccountAdd     			TransactionType = "add-account"
	AccountDel     			TransactionType = "del-account"
	AttendPatient		    TransactionType = "attend-patient"
	DoctorAdd 				TransactionType = "add-doctor"
	DoctorDel				TransactionType = "del-doctor"
	MedicalAppointmentAdd	TransactionType = "add-medical-appointment"
	MedicalAppointmentDel	TransactionType = "del-medical-appointment"
	MedicalRecordAdd		TransactionType = "add-medical-record"
	MedicalRecordDel		TransactionType = "del-medical-record"
	MedicalRecordUpdate		TransactionType = "update-medical-record"
	PatientAdd    			TransactionType = "add-patient"
	PatientDel				TransactionType = "del-patient"	
	ReputationGive 			TransactionType = "give-reputation"
	SecretAdd      			TransactionType = "secret-add"
	SecretUpdate   			TransactionType = "secret-update"
	SecretDel      			TransactionType = "secret-del"
	SecretShare    			TransactionType = "secret-share"
)

const DefaultProofOfWorkCost byte = 16

//FromBytes ...
func (t *Transaction) FromBytes(bs []byte) error {
	return json.Unmarshal(bs, t)
}

//ToBytes ...
func (t *Transaction) ToBytes() ([]byte, error) {
	return json.Marshal(t)
}

//Hash ...
func (t *Transaction) Hash() []byte {
	hash := sha3.New512()
	encoder := json.NewEncoder(hash)
	encoder.Encode(t.Type)
	encoder.Encode(t.Timestamp)
	if hashable, ok := t.Data.(Hashable); ok {
		hash.Write(hashable.Hash())
	} else {
		encoder.Encode(t.Data)
	}
	return hash.Sum(nil)
}

//Hash2 ...
func (t *Transaction) Hash2() []byte {
	hash := sha3.New512()
	fmt.Println("type:"+t.Type)
	hash.Write([]byte(t.Type))
	return hash.Sum(nil)
}

//Sign transaction
func (t *Transaction) Sign(key *crypto.Key) error {
	hash := t.Hash()
	signature, err := key.Sign(hash)
	if err != nil {
		return err
	}
	t.Signature = signature
	return nil
}

//Verify transaction
func (t *Transaction) Verify(key *crypto.Key) error {
	hash := t.Hash2()
	fmt.Println("hash 2:",hash)
	return key.Verify(hash, t.Signature)
}

//ProofOfWork ...
func (t *Transaction) ProofOfWork(cost byte) error {
	for round := 0; round < (1 << 32); round++ {
		t.Nonce = uint32(round)
		if err := t.VerifyProofOfWork(cost); err == nil {
			return nil
		}
	}
	return errors.New("can not find pow")
}

//VerifyProofOfWork ...
func (t *Transaction) VerifyProofOfWork(cost byte) error {
	hasher := sha3.New512()
	hasher.Write(t.Hash())
	binary.Write(hasher, binary.LittleEndian, t.Nonce)
	tip := uint64(0)
	buf := bytes.NewBuffer(hasher.Sum(nil))
	binary.Read(buf, binary.LittleEndian, &tip)
	if tip<<(64-cost) == 0 {
		return nil
	}
	return errors.New("failed to validate proof of work")
}

//New ...
func New(t TransactionType, data interface{}) *Transaction {
	fmt.Println("TIEMPO:",time.Now().String())
	return &Transaction{t, time.Now(), "", 0, data}
}

func hashStringMap(m map[string]interface{}) []byte {
	hash := sha3.New512()
	encoder := json.NewEncoder(hash)
	keys := make([]string, len(m))
	i := 0
	for id := range m {
		keys[i] = id
		i++
	}
	sort.Strings(keys)
	for _, key := range keys {
		encoder.Encode(key)
		encoder.Encode(m[key])
	}
	return hash.Sum(nil)
}