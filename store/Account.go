package store

import (
	"encoding/json"
	"errors"

	"github.com/ccamaleon5/saludchain/crypto"
)

//Account ...
type Account struct {
	ID         string         `json:"id" mapstructure:"id"`
	PubKey     string         `json:"pubkey" mapstructure:"pubkey"`
}

//AddAccount ...
func (s *State) AddAccount(account *Account) error {
	if s.HasAccount(account.ID) {
		return errors.New("account already exists")
	}
	return s.SetAccount(account)
}

//SetAccount ...
func (s *State) SetAccount(account *Account) error {
	bs, err := json.Marshal(account)
	if err != nil {
		return err
	}
	s.Tree.Set([]byte(accountPrefix+account.ID), bs)
	return nil
}

//HasAccount ...
func (s *State) HasAccount(id string) bool {
	return s.Tree.Has([]byte(accountPrefix + id))
}

//GetAccount ...
func (s *State) GetAccount(id string) (*Account, error) {
	_, bs := s.Tree.Get([]byte(accountPrefix + id))
	if bs == nil {
		return nil, errors.New("no such account")
	}
	acc := &Account{}
	return acc, json.Unmarshal(bs, acc)
}

//DeleteAccount ...
func (s *State) DeleteAccount(id string) error {
	_, removed := s.Tree.Remove([]byte(accountPrefix + id))
	if !removed {
		return errors.New("no such account")
	}
	return nil
}

//GetAccountPubKey ...
func (s *State) GetAccountPubKey(id string) (*crypto.Key, error) {
	acc, err := s.GetAccount(id)
	if err != nil {
		return nil, err
	}
	return crypto.NewFromStrings(acc.PubKey, "")
}

//ListAccounts ...
func (s *State) ListAccounts() (result []*Account, err error) {
	start := accountPrefix
	end := start[:len(start)-1]
	end = end + string(start[len(start)-1]+1)
	result = make([]*Account, 0)
	s.Tree.IterateRange([]byte(start), []byte(end), true, func(key []byte, value []byte) bool {
		acc := &Account{}
		err = json.Unmarshal(value, acc)
		if err != nil {
			return true
		}
		result = append(result, acc)
		return false
	})
	return
}