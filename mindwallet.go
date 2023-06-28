package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/btcsuite/btcutil/base58"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/patcito/monero"
	"github.com/patcito/monero/crypto"
	"github.com/vsergeev/btckeygenie/btckey"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/pbkdf2"
)

func usage() {
	fmt.Printf("Usage: %s [bitcoin|monero|litecoin|ethereum] [Passphrase] [Salt] [normal (default)|strong|super_strong|ridiculously_strong]\n\n", os.Args[0])
	fmt.Printf("Example: ./mindwallet ethereum MyPassword MyPassphrase %s \n\n")
	os.Exit(1)
}

func fail(err error) {
	fmt.Printf("%s\n", err)
	os.Exit(1)
}

func main() {
	var passphrase string
	var salt string
	var currency string
	var strong string
	var priv btckey.PrivateKey
	var result [32]byte
	var hash_pad byte
	var key []byte

	if len(os.Args) < 4 || len(os.Args) > 5 {
		usage()
	}

	currency = os.Args[1]
	passphrase = os.Args[2]
	salt = os.Args[3]
	if len(os.Args) > 4 {
		strong = os.Args[4]
	}
	if currency == "bitcoin" {
		hash_pad = 1
	} else if currency == "litecoin" {
		hash_pad = 2
	} else if currency == "monero" {
		hash_pad = 3
	} else if currency == "ethereum" {
		hash_pad = 4
	} else {
		usage()
	}
	if len(os.Args) == 6 {
		argval, err := strconv.Atoi(os.Args[6])
		if err != nil {
			usage()
		} else {
			hash_pad = byte(argval)
		}
	}
	start := time.Now()

	fmt.Printf("Passphrase: %s\nSalt: %s\nHashPad: %d\n", passphrase, salt, hash_pad)
	_passphrase := fmt.Sprint(passphrase, string(rune(hash_pad)))
	_salt := fmt.Sprint(salt, string(rune(hash_pad)))
	if strong == "strong" {
		fmt.Printf("\n\n Depending on your CPU, this make take up to 4 minute or more\n\n")
		key = argon2.IDKey([]byte(_passphrase), []byte(_salt), 256, 10*256*1024, 32, 32)
	} else if strong == "super_strong" {
		fmt.Printf("\n\n Depending on your CPU, this make take up to 8 minutes or more\n\n")
		key = argon2.IDKey([]byte(_passphrase), []byte(_salt), 512, 15*256*1024, 32, 32)
	} else if strong == "ridiculously_strong" {
		fmt.Printf("\n\n Depending on your CPU, this make take up to 20 minutes or more, you've been warned!\n\n")
		key = argon2.IDKey([]byte(_passphrase), []byte(_salt), 1024, 20*256*1024, 32, 32)
	} else {
		fmt.Printf("\n\n Depending on your CPU, this make take up to 10 seconds, you've been warned!\n\n")
		key = argon2.IDKey([]byte(_passphrase), []byte(_salt), 8, 10*256*1024, 32, 32)
	}

	//	key := argon2.Key([]byte("password"), []byte("somesalt"), 10, 1024, 1, 32)
	//fmt.Printf("salt: %x\n _passphrase: %x\n _salt: %x\n key: %x\n BYTEPASS: %x\n", salt, _passphrase, _salt, key, []byte("password"))

	_passphrase = fmt.Sprint(passphrase, string(rune(hash_pad+1)))
	_salt = fmt.Sprint(salt, string(rune(hash_pad+1)))
	fmt.Printf("\n\n Now generating pbkdf2 key\n\n")

	key2 := pbkdf2.Key([]byte(_passphrase), []byte(_salt), 65536, 32, sha256.New)

	for i := 0; i < len(key); i++ {
		result[i] = key[i] ^ key2[i]
	}
	fmt.Printf("Seed: %x\n", result)

	if currency == "bitcoin" {
		priv.FromBytes(result[:])
		privkey := priv.ToWIF()
		address_uncompressed := priv.ToAddressUncompressed()
		fmt.Printf("Bitcoin Address: %s\n", address_uncompressed)
		fmt.Printf("Private Key: %s\n", privkey)

	} else if currency == "litecoin" {
		var version byte = 48
		priv.FromBytes(result[:])
		privkey := base58.CheckEncode(result[:], version+0x80)
		addr_bytes, ver, err := base58.CheckDecode(priv.ToAddressUncompressed())
		if err != nil || ver != 0 {
			fail(err)
		}
		address_uncompressed := base58.CheckEncode(addr_bytes, version)
		fmt.Printf("Litecoin Address: %s\n", address_uncompressed)
		fmt.Printf("Private Key: %s\n", privkey)

	} else if currency == "ethereum" {
		priv.FromBytes(result[:])
		privkey := priv.ToBytesUncompressed()
		address := ethcrypto.Keccak256(privkey[1:])[12:]

		fmt.Printf("Ethereum Address: 0x%x\n", address)
		fmt.Printf("Private Key: %x\n", result)

	} else if currency == "monero" {
		var secret [32]byte
		var view_secret [32]byte
		var spend_secret [32]byte
		crypto.SecretFromSeed(&secret, &result)
		account, err := monero.RecoverAccount(secret)
		if err != nil {
			fail(err)
		}
		spend_secret = account.Secret()
		crypto.ViewFromSpend(&view_secret, &spend_secret)
		fmt.Printf("Address: %s\n", account.Address().String())
		fmt.Printf("Private Spend Key: %x\n", spend_secret)
		fmt.Printf("Private View Key: %x\n", view_secret)
		mnemonic, err := account.Mnemonic()
		if err != nil {
			fail(err)
		}
		fmt.Printf("Mnemonic: %s\n", mnemonic)
	}
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Printf("Elapsed time: %s\n", elapsed)
}
