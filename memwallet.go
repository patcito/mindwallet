package main

import (
    "github.com/vsergeev/btckeygenie/btckey"
    "github.com/btcsuite/btcutil/base58"
    "golang.org/x/crypto/pbkdf2"
    "golang.org/x/crypto/scrypt"
    "github.com/ehmry/monero"
    "github.com/ehmry/monero/crypto"
    ethcrypto "github.com/ethereum/go-ethereum/crypto"
    "crypto/sha256"
    "strconv"
    "fmt"
    "os"
)

func usage () {
  fmt.Printf("Usage: %s [bitcoin|monero|litecoin|ethereum] [Passphrase] [Salt]\n\n", os.Args[0])
  os.Exit(1)
}

func fail (err error) {
  fmt.Printf("%s\n", err)
  os.Exit(1)
}

func main () {
    var passphrase string
    var salt string
    var currency  string
    var priv btckey.PrivateKey
    var result [32]byte
    var hash_pad byte

    if len(os.Args) < 4 || len(os.Args) > 5 {
      usage()
    }

    currency = os.Args[1]
    passphrase = os.Args[2]
    salt = os.Args[3]


    if currency == "bitcoin" {
      hash_pad = 1
    } else if  currency == "litecoin" {
      hash_pad = 2
    } else if  currency == "monero" {
      hash_pad = 3
    } else if  currency == "ethereum" {
      hash_pad = 4
    } else {
      usage()
    }
    if len(os.Args) == 5 {
      argval, err := strconv.Atoi(os.Args[5])
      if err != nil {
        usage()
      } else {
        hash_pad = byte(argval)
      }
    }

    fmt.Printf("Passphrase: %s\nSalt: %s\nHashPad: %d\n", passphrase, salt, hash_pad)
    _passphrase := fmt.Sprint(passphrase, string(rune(hash_pad)))
    _salt := fmt.Sprint(salt, string(rune(hash_pad)))
    key, _ := scrypt.Key([]byte(_passphrase), []byte(_salt), 262144, 8, 1, 32)

    _passphrase = fmt.Sprint(passphrase, string(rune(hash_pad + 1)))
    _salt = fmt.Sprint(salt, string(rune(hash_pad + 1)))
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
      privkey := base58.CheckEncode(result[:], version + 0x80)
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
      var secret [32]byte;
      var view_secret [32]byte;
      var spend_secret [32]byte;
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
}
