package crypto

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

func decodeScalar(s string) *[32]byte {
	sc := new([32]byte)
	hex.Decode(sc[:], []byte(s))
	return sc
}

func readTestLines(t *testing.T, name string) (lines [][]string) {
	file, err := os.Open("tests.txt")
	if err != nil {
		t.Fatal("Failed to open tests file,", err)
	}
	defer file.Close()

	r := bufio.NewReader(file)

	lines = make([][]string, 0, 32)

	var line string
	for {
		line, err = r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			t.Fatal("Failed to test from file,", err)
		}

		if strings.HasPrefix(line, name) {
			break
		}
	}

	l := len(name)

	for {
		lines = append(lines, strings.Fields(line[l:]))
		line, err = r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			t.Fatal("Failed to tests from file,", err)
		}

		if !strings.HasPrefix(line, name) {
			return
		}
	}
	return
}

func TestCheckScalar(t *testing.T) {
	lines := readTestLines(t, "check_scalar")
	for i, args := range lines {
		s := decodeScalar(args[0])
		ok := scCheck(s)

		if ok != (args[1] == "true") {
			t.Errorf("CheckScalar %d: %s, want %s, got %v",
				i, args[0], args[1], ok)
		}
	}
}

func TestHashToScalar(t *testing.T) {
	lines := readTestLines(t, "hash_to_scalar")
	var s [32]byte
	for i, args := range lines {
		in, _ := hex.DecodeString(args[0])
		hashToScalar(&s, in)
		out, _ := hex.DecodeString(args[1])

		if !bytes.Equal(s[:], out) {
			t.Errorf("HashToScalar %d: %s, want %s, got %x",
				i, args[0], args[1], s[:])
		}
	}
}

func TestCheckKey(t *testing.T) {
	lines := readTestLines(t, "check_key")
	for i, args := range lines {
		key, _ := hex.DecodeString(args[0])
		ok := checkKey(key)
		if ok != (args[1] == "true") {
			t.Errorf("check_key %d: %s, want %s, got %v",
				i, args[0], args[1], ok)
		}
	}
}

func TestSecretKeyToPublicKey(t *testing.T) {
	lines := readTestLines(t, "secret_key_to_public_key")
	publicTest := new([32]byte)
	for i, args := range lines {
		secret := decodeScalar(args[0])
		PublicFromSecret(publicTest, secret)

		valid := args[1] == "true"
		validTest := CheckSecret(secret)

		if valid != validTest {
			t.Errorf("secret_key_to_public_key %d: want error to be %v, error was %v", i, valid, validTest)
		}

		if valid {
			publicControl, _ := hex.DecodeString(args[2])
			if !bytes.Equal(publicTest[:], publicControl) {
				t.Errorf("secret_key_to_public_key %d: want %s, got %x",
					i, args[2], publicTest[:],
				)
			}
		}
	}
}

func TestGenerateKeyDerivation(t *testing.T) {
	lines := readTestLines(t, "generate_key_derivation")
	for i, args := range lines {
		public := decodeScalar(args[0])
		secret := decodeScalar(args[1])
		valid := (args[2] == "true")

		d, err := generateKeyDerivation(public, secret)
		if (err == nil) != valid {
			t.Errorf("generate_key_derivation %d: want error to be %v, error was %v", i, valid, err == nil)
		}

		if valid {
			control, _ := hex.DecodeString(args[3])

			if !bytes.Equal(d[:], control) {
				t.Errorf("generate_key_derivation %d: want %s, got %x", i, args[3], d)
			}
		}
	}
}

func TestDerivePublicKey(t *testing.T) {
	lines := readTestLines(t, "derive_public_key")
	var (
		derivation, control []byte
		test                *[32]byte

		outputIndex uint64
		base        *[32]byte
		valid       bool

		err error
	)

	for i, args := range lines {
		derivation, _ = hex.DecodeString(args[0])
		outputIndex, _ = strconv.ParseUint(args[1], 10, 64)
		base = decodeScalar(args[2])
		valid = args[3] == "true"
		if valid {
			control, _ = hex.DecodeString(args[4])
		}

		test, err = derivePublicKey(derivation, outputIndex, base)
		if (err == nil) != valid {
			t.Errorf("derive_public_key %d: want error to be %v, error was %v", i, valid, err == nil)
		}

		if valid {
			if !bytes.Equal(test[:], control) {
				t.Errorf("derive_public_key %d: want %s, got %x", i, args[4], test[:])
			}
		}
	}
}

func TestDeriveSecretKey(t *testing.T) {
	for i, args := range readTestLines(t, "derive_secret_key") {
		derivation, _ := hex.DecodeString(args[0])
		outputIndex, _ := strconv.ParseUint(args[1], 10, 64)
		base := decodeScalar(args[2])
		control, _ := hex.DecodeString(args[3])

		test, err := deriveSecretKey(derivation, outputIndex, base)
		if err != nil {
			t.Errorf("derive_secret_key %d error: %v", err)
		}

		if !bytes.Equal(test[:], control) {
			t.Errorf("derive_secret_key %d: want %s, got %x", i, args[3], test)
		}
	}
}

/*
func TestGenerateSignature(t *testing.T) {
	lines := readTestLines(t, "generate_signature")
	var (
		prefixHash []byte
		public     *Public
		secret     *Secret

		test, control []byte
	)

	for i, args := range lines {
		prefixHash, _ = hex.DecodeString(args[0])
		public = decodeScalar(args[1])
		secret = decodeScalar(args[2])
		control, _ = hex.DecodeString(args[3])

		test = generateSignature(prefixHash, public, secret)
		if !bytes.Equal(test, control) {
			t.Errorf("generate_signature %d: want %s, got %x",
				i, args[3], test)
		}
	}
}
*/

func TestCheckSignature(t *testing.T) {
	lines := readTestLines(t, "check_signature")
	for i, args := range lines {
		prefixHash, _ := hex.DecodeString(args[0])
		pub := decodeScalar(args[1])
		sig, _ := hex.DecodeString(args[2])
		control := args[3] == "true"

		test := checkSignature(prefixHash, pub, sig)
		if test && !control {
			t.Errorf("check_signature %d: should not be valid.", i)
		} else if !test && control {
			t.Errorf("check_signature %d: should be valid.", i)
		}
	}
}

func TestHashToPoint(t *testing.T) {
	lines := readTestLines(t, "hash_to_point")
	for i, args := range lines {
		hash, _ := hex.DecodeString(args[0])
		control, _ := hex.DecodeString(args[1])

		point := hashToPoint(hash)

		if !bytes.Equal(point[:], control) {
			t.Errorf("hash_to_point %d: want %s, got %x", i, args[1], point[:])
		}
	}
}

func TestHashToEC(t *testing.T) {
	lines := readTestLines(t, "hash_to_ec")
	for i, args := range lines {
		key := decodeScalar(args[0])
		control := args[1]

		test := hashToEC(key)

		if test.String() != control {
			t.Errorf("hash_to_point %d: want %s, got %s", i, control, test)
		}
	}
}

func TestGenerateKeyImage(t *testing.T) {
	lines := readTestLines(t, "generate_key_image")
	for i, args := range lines {
		public := decodeScalar(args[0])
		secret := decodeScalar(args[1])
		control, _ := hex.DecodeString(args[2])

		test := generateKeyImage(public, secret)
		if !bytes.Equal(test[:], control) {
			t.Errorf("generate_key_image %d: want %s, got %x", i, args[2], test)
		}
	}
}

func TestCheckRingSignature(t *testing.T) {
	for i, args := range readTestLines(t, "check_ring_signature") {
		prefixHash, _ := hex.DecodeString(args[0])
		image, _ := hex.DecodeString(args[1])
		pubCount, _ := strconv.Atoi(args[2])
		pubs := make([]*[32]byte, pubCount)
		for j := 0; j < pubCount; j++ {
			pubs[j] = decodeScalar(args[3+j])
		}

		sigs, _ := hex.DecodeString(args[3+pubCount])

		control := args[4+pubCount] == "true"
		test := checkRingSignature(prefixHash, image, pubs, sigs)
		if test != control {
			t.Fatalf("check_ring_signature %d: want %v, got %v", i, control, test)
		}
	}
}
