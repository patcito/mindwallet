package monero

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// WordsToBytes converts 12 or 24 Electrum words into a key buffer.
func WordsToBytes(dst *[32]byte, words []string) error {
	if len(words) != 12 && len(words) != 24 {
		return fmt.Errorf("word list must be 12 or 24 words long, not %d", len(words))
	}

	buf := bytes.NewBuffer(dst[:0])

	var (
		val, w1, w2, w3 uint32
		word            string
		ok              bool
	)

	for i := 0; i < len(words)/3; i++ {

		word = words[i*3]
		w1, ok = wordsMap[word]
		if ok {
			word = words[i*3+1]
			w2, ok = wordsMap[word]
			if ok {
				word = words[i*3+2]
				w3, ok = wordsMap[word]
			}
		}
		if !ok {
			return fmt.Errorf("%q not in wordlist", word)
		}

		val = w1 + numwords*(((numwords-w1)+w2)%numwords) + numwords*numwords*(((numwords-w2)+w3)%numwords)
		if !(val%numwords == w1) {
			return fmt.Errorf("%d %% numwords == %d", val, w1)
		}

		if err := binary.Write(buf, binary.LittleEndian, val); err != nil {
			return err
		}
	}

	if len(words) == 12 {
		copy(dst[16:], dst[:16])
	}
	return nil
}

// WordsToBytes converts bytes into Electrum words.
func BytesToWords(b []byte) (words []string, err error) {
	buf := bytes.NewReader(b)
	var w1, w2, w3, val uint32

	if len(b)%4 != 0 {
		err = fmt.Errorf("BytesToWords called on a slice not divisible by 4")
		return
	}

	words = make([]string, len(b)/4*3)

	i := 0
	for {
		if err := binary.Read(buf, binary.LittleEndian, &val); err != nil {
			break
		}
		w1 = val % numwords
		w2 = ((val / numwords) + w1) % numwords
		w3 = (((val / numwords) / numwords) + w2) % numwords

		words[i] = wordsArray[w1]
		i++
		words[i] = wordsArray[w2]
		i++
		words[i] = wordsArray[w3]
		i++
	}
	return
}
