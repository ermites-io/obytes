// +build go1.12

package obytes

import (
	"bytes"
	"crypto/rand"
	"testing"
)

//  0x55        0x55
// [0101 0101] [0101 0101]

// super incomplete but it's a start
type ioBlobTestEntry struct {
	in       []byte // input
	plen     int    // prefix len
	mask_lo  byte
	mask_hi  byte
	inerr    error // expected error
	outerr   error // expected error
	equality bool
}

var blobTestVector = []ioBlobTestEntry{
	{[]byte("A"), 0, 0x55, 0x55, nil, nil, true},
	{[]byte("A"), 0, 0x05, 0xf5, nil, nil, true},
	{[]byte("F"), 0, 0xf5, 0x05, nil, nil, true},
	{[]byte("F"), 0, 0xff, 0x00, nil, nil, true},
}

func TestBlob(t *testing.T) {
	for i := 0; i < len(blobTestVector); i++ {
		blob, err := NewBlob(blobTestVector[i].in, tests[i].mask_hi, tests[i].mask_lo, blobTestVector[i].plen)
		if err != blobTestVector[i].inerr {
			t.Logf("NewBlob error\n")
			t.Fail()
		}

		plain, err := OpenBlob(blob, blobTestVector[i].mask_hi, tests[i].mask_lo, blobTestVector[i].plen)
		if err != blobTestVector[i].inerr {
			t.Logf("DeObfuscateByte error\n")
			t.Fail()
		}

		if blobTestVector[i].equality && !bytes.Equal(blobTestVector[i].in, plain) {
			t.Logf("matching error %02x vs %02x\n", blobTestVector[i].in, plain)
			t.Fail()
		}
	}
}

func TestRandBlob(t *testing.T) {
	for i := 0; i < 1000; i++ {
		input := make([]byte, 1024)
		_, err := rand.Read(input)
		if err != nil {
			t.Logf("rand error\n")
			t.Fail()
		}

		blob, err := NewBlob(input, 0x55, 0x55, i)
		if err != nil {
			t.Logf("NewBlob error\n")
			t.Fail()
		}

		output, err := OpenBlob(blob, 0x55, 0x55, i)
		if err != nil {
			t.Logf("OpenBlob error\n")
			t.Fail()
		}

		if !bytes.Equal(input, output) {
			t.Logf("matching error %.16x vs %.16x\n", input, output)
			t.Fail()
		}
	}
}
