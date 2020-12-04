package monero

import "testing"

const control = "42uWJwLRQRmSLm6DntgH3h2BvdLA3xqo1amRwPjQysCiii56jQE2uyG7vmQgZzCRpZarxg5LCUhPFRGE4VtHK5oqG1uvTnZ"

func TestDecodeEncodeAddress(t *testing.T) {
	addr, err := DecodeAddress(control)
	if err != nil {
		t.Fatal("Error decoding address,", err)
	}
	if control != addr.String() {
		t.Errorf("Decoding and encoding failed,\nwanted %s,\ngot    %s", control, addr)
	}
}
