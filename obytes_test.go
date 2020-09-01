// +build go1.12

package obytes

import (
	"crypto/rand"
	"testing"
)

//  0x55        0x55
// [0101 0101] [0101 0101]

type iotest struct {
	/*
	inb     []byte
	outb    []byte
	*/
	in      byte // input
	out     byte // output
	mask_lo byte
	mask_hi byte
}

var tests = []iotest{
	{byte('A'), 0, 0x55, 0x55},
	{byte('A'), 0, 0x05, 0xf5},
	{byte('F'), 0, 0xf5, 0x05},
	{byte('F'), 0, 0xff, 0x00},
}

func TestObfuscateByte(t *testing.T) {
	for i := 0; i < len(tests); i++ {
		obfOut, err := obfuscateByte(tests[i].in, tests[i].mask_hi, tests[i].mask_lo)
		if err != nil {
			t.Logf("ObfuscateByte error\n")
			t.Fail()
		}
		//t.Logf("ObfuscateByte in: %x out: %x\n", tests[i].in, obfOut)

		outIn, err := deobfuscateByte(obfOut, tests[i].mask_hi, tests[i].mask_lo)
		if err != nil {
			t.Logf("DeObfuscateByte error\n")
			t.Fail()
		}
		//t.Logf("DeObfuscateByte in: %x out: %x\n", obfOut, outIn)

		if tests[i].in != outIn {
			t.Logf("matching error %02x vs %02x\n", tests[i].in, outIn)
			t.Fail()
		}
	}

	for i := 0; i < 100000; i++ {
		obfIn := make([]byte, 1)
		rand.Read(obfIn)

		obfOut, err := obfuscateByte(obfIn[0], 0x55, 0x55)
		if err != nil {
			t.Logf("ObfuscateByte error\n")
			t.Fail()
		}

		outIn, err := deobfuscateByte(obfOut, 0x55, 0x55)
		if err != nil {
			t.Logf("DeObfuscateByte error\n")
			t.Fail()
		}

		if obfIn[0] != outIn {
			t.Logf("matching error %02x vs %02x\n", obfIn[0], outIn)
			t.Fail()
		}
	}
}

func TestObfuscateError(t *testing.T) {
	/*
	obfOut, err := obfuscateByte(0, 0x55, 0x55)
	if err == nil || obfOut != nil {
		t.Logf("ObfuscateByte Error Test\n")
		t.Fail()
	}
	*/

	obfOut, err := Obfuscate(nil)
	if err == nil || obfOut != nil {
		t.Logf("Obfuscate Error Test\n")
		t.Fail()
	}
}

func TestDeObfuscateError(t *testing.T) {
	obfOut, err := deobfuscateByte(nil, 0x55, 0x55)
	if err == nil || obfOut != 0 {
		t.Logf("DeobfuscateByte Error Test\n")
		t.Fail()
	}

	obfOut, err = deobfuscateByte([]byte("123"), 0x55, 0x55)
	if err == nil || obfOut != 0 {
		t.Logf("DeobfuscateByte Error Test\n")
		t.Fail()
	}

	obfOutByte, err := Deobfuscate(nil)
	if err == nil || obfOutByte != nil {
		t.Logf("Deobfuscate Error Test\n")
		t.Fail()
	}
}
