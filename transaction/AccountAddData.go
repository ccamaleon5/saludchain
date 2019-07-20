package transaction

import (
	"encoding/json"
	"sort"

	"github.com/ccamaleon5/saludchain/store"
	"golang.org/x/crypto/sha3"
)

type AccountAddData struct {
	Account *store.Account
}

func (data *AccountAddData) Hash() []byte {
	hash := sha3.New512()
	encoder := json.NewEncoder(hash)
	encoder.Encode(data.Account.ID)
	encoder.Encode(data.Account.PubKey)
	return hash.Sum(nil)
}

func hashReputation(m map[string]int) []byte {
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