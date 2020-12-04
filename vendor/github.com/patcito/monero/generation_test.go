package monero

import (
	crand "crypto/rand"
	"math/rand"
	"strings"
	"testing"
	"time"
	//"github.com/patcito/monero/crypto"
)

func TestWordToBytes(t *testing.T) {
	rand.Seed(time.Now().Unix())

	words := make([]string, 24)

	for i := 0; i < 24; i++ {
		words[i] = wordsArray[rand.Intn(numwords)]
	}

	var b [32]byte
	WordsToBytes(&b, words)
	result, err := BytesToWords(b[:])
	if err != nil {
		t.Fatalf("error recovering word list, %s", err)
	}
	if len(result) != 24 {
		t.Fatalf("recovered word list was %d, not 24 long", len(result))
	}

	for i := 0; i < 24; i++ {
		if words[i] != result[i] {
			t.Errorf("Word mismatch: %02d - %q != %q", i, words[i], result[i])
		}
	}
}

type test struct {
	addr          string
	words         []string
	deterministic bool
}

var tests = []*test{
	{
		"43GVEVSitCqRxtuRXWbpsu6trCmHXhqNM4myNZka86JMMtw75eWVduKRJ2rz3yjTUoCPf9mkJHjfC9JRZUM7f3fSM52vYej",
		strings.Fields("happen former recall kill tonight magic mercy threw somehow arrive meant sheet charm victim once indeed hug bubble wash storm hill bid respect excuse"),
		true,
	},
	{
		"4A4ZprYVAdC8iKm1bVwPmQi2Q7KtEZL94d7tCHEaQV9FBa4UgeUCDaRAeqvSRgwbeQ67xrSmABVQyMZX2KuuNAV3Bk8cLW1",
		strings.Fields("empty favorite good iron spend memory grand dark direction brain out pleasure climb hardly out claim neither lick hidden button aim shiver gently treat"),
		true,
	},
	{
		"468vZRyTA7F5mjPFXTFcwp3bsTBHbMLztKD9B3FbHZewiyHXBWKmAgFcQChApA5gM34eAyX6siSpy2vwifZ8Cd6nSuq5Dau",
		strings.Fields("handle brother whistle realize money upon leg level doubt shove count memory wonderful pop clear huge tap less age circle slowly weather gasp grief"),
		true,
	},
	{
		"455MJ7FvGZL8rmxHSjE4z7AFsCJTmF9L2U3nLw4fD7zDRT1K5xjYUadVEcuekSbDereYgAQcWrJGyd42K4L9bTgb7WKJiFV",
		strings.Fields("claim pride forward strain piece group torture stream balance unknown lick common useless empty prayer good sunlight trouble return snap gone focus measure scale"),
		true,
	},
	{
		"48ZgBHE2G4oYyinm59WWEmeG2BLvHCx96WMQbzEvE8TDiGx7WUsjSpMgBNghuVYscY8VVYqzpydvSSL5BNTUifVPAyc1k6y",
		strings.Fields("rhyme remain sigh leg pray freedom nine around table planet down connect apart inhale daughter defense spin mind especially easy grow quiet coward belly"),
		false,
	},
	{
		"49ifQoT1Sn8B3Mbn4RAfo2N3hdwr26txzWGP83JbJYSe86HQLi2dteAaviHk7rFg4gMh1Qjo6XL4kRLuzc9FUaaE69mPCiR",
		strings.Fields("flood wolf lot swear orange act marry tap steel scream finish calm river friendship pulse fact storm ugly fill creak tickle stress guitar hurt"),
		false,
	},
	{
		"47DASo17ysKUy6K14tqXhzBbY8VvNRVXDU2SiNSEWefvMeVfdAbYhPSZvSDyRVJHhpWKxcpVAbHWLfdSGZHrkjaN6rXGwJe",
		strings.Fields("depth frame youth learn journey strip rude prefer finish edge slice verse split prefer release shower language burst bounce slowly flower smoke spring seem"),
		false,
	},
	{
		"42roLXCWzE4RaEpZwY9r6YfQfacKzMya2hMJa6tVN3L4RkkMbQCHBX6ZFtMg8AmAh7F91MdtTWwzBCt5WLwithYpBwzV3Qu",
		strings.Fields("yellow nature circle rush wet iron angel plate much said heel flood deeply twirl teeth future crack ashamed salty describe hidden idiot pants liquid"),
		false,
	},
}

func TestRecovery(t *testing.T) {
	for _, test := range tests {
		if !test.deterministic {
			// This doesn't work,
			continue
		}
		account, err := RecoverAccountWithMnemonic(test.words)
		if err != nil {
			t.Fatal("mnemonic recovery failed,", err)
		}
		if test.addr != account.String() {
			t.Fatalf("mnemonic recovery failed,\nwanted %q\ngot    %q", test.addr, account)
		}

		words, err := account.Mnemonic()
		if test.deterministic == (err != nil) {
			t.Errorf("Mnemonic incorrectly determined determinism for %s", test.addr)
		}
		for i := 0; i < len(words); i++ {
			if test.words[i] != words[i] {
				t.Errorf("Mnemonic() failed %d, wanted %q got %q", i, test.words[i], words[i])
			}
		}
	}

	for i := 0; i < 32; i++ {
		first, err := GenerateAccount(crand.Reader)
		if err != nil {
			t.Fatal("GenerateAccount:", err)
		}

		words, err := first.Mnemonic()
		if err != nil {
			t.Fatal("Mnemonic:", err)
		}

		second, err := RecoverAccountWithMnemonic(words)
		if err != nil {
			t.Fatal("RecoverAccountWithMnemonic:", err)
		}

		if first.String() != second.String() {
			t.Errorf("account generation and recovery failed\n %s != %s", first, second)
		}
	}
}
