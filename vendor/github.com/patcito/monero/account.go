package monero

import (
	"bytes"
	"errors"
	"io"

	"github.com/patcito/monero/crypto"
)

// Account contains public and private keys for the spend and view
// aspects of a Monero account.
type Account struct {
	spendP, spendS [32]byte
	viewP, viewS   [32]byte
}

// Address returns the address of a given account.
func (a *Account) Address() *Address {
	return &Address{spend: a.spendP, view: a.viewP}
}

func (a *Account) String() string {
	return a.Address().String()
}

// Secret returns the spend secret key that
// can be used to regenerate the account.
func (a *Account) Secret() [32]byte {
	return a.spendS
}

var NonDeterministic = errors.New("account is non-deterministic")

// Mnemonic returns an Electrum style mnemonic representation
// of the account spend secret key.
func (a *Account) Mnemonic() ([]string, error) {
	ss := a.spendS[:]

	var test [32]byte
	crypto.ViewFromSpend(&test, &a.spendS)

	if !bytes.Equal(a.viewS[:], test[:]) {
		return nil, NonDeterministic
	}

	return BytesToWords(ss)
}

// Recover recovers an account using a secret key.
func RecoverAccount(secret [32]byte) (*Account, error) {
	if !crypto.CheckSecret(&secret) {
		return nil, crypto.InvalidSecret
	}

	a := &Account{spendS: secret}
	crypto.PublicFromSecret(&a.spendP, &a.spendS)
	crypto.ViewFromSpend(&a.viewS, &a.spendS)
	crypto.PublicFromSecret(&a.viewP, &a.viewS)
	return a, nil
}

// RecoverAccountWithMnemonic recovers an account
// with an Electrum style word list.
func RecoverAccountWithMnemonic(words []string) (*Account, error) {
	var seed [32]byte
	if err := WordsToBytes(&seed, words); err != nil {
		return nil, err
	}
	return RecoverAccount(seed)
}

// GenerateAccountGenerates a new account.
func GenerateAccount(random io.Reader) (acc *Account, err error) {
	acc = new(Account)
	acc.spendS, err = crypto.GenerateSecret(random)
	if err != nil {
		return nil, err
	}

	crypto.PublicFromSecret(&acc.spendP, &acc.spendS)
	crypto.ViewFromSpend(&acc.viewS, &acc.spendS)
	crypto.PublicFromSecret(&acc.viewP, &acc.viewS)
	return acc, nil
}
