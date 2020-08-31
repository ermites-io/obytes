// +build go1.12

package obytes

import (
	"crypto/rand"
)

const (
	MASK_HI = 1
	MASK_LO = 0
)

/*
I NEED TO DEFINE MASKs as a pair of byte constants..
const (
    MASK0 = []byte { 0x55, 0x55 }
    MASK1 = []byte { 0xf0, 0x0f }
    MASK2 = []byte { 0x0f, 0xf0 }
)
*/

/* simple obfuscation keep you busy and suck your time. */
func obfuscateByte(in byte, mask_hi, mask_lo byte) (out []byte, err error) {
	rnd := make([]byte, 2)
	out = make([]byte, 2)

	// get random bytes & apply the mask
	_, err = rand.Read(rnd)
	if err != nil {
		return nil, ErrUnexpected //fmt.Errorf("rand error")
	}
	rnd[0] &= ^(mask_lo)
	rnd[1] &= ^(mask_hi)

	/* lo nibble */
	for i := 0; i < 8; i++ {
		if (mask_lo>>uint(i))&0x01 == 0x01 {
			out[0] |= (in & 0x01) << uint(i)
			in >>= 0x01
		} else {
			out[0] |= ((rnd[0] >> uint(i) & 0x01) << uint(i))
		}
	}

	/* hi nibble */
	for i := 0; i < 8; i++ {
		if (mask_hi>>uint(i))&0x01 == 0x01 {
			out[1] |= (in & 0x01) << uint(i)
			in >>= 0x01
		} else {
			out[1] |= ((rnd[1] >> uint(i) & 0x01) << uint(i))
		}
	}

	/*
	   fmt.Printf("MASKED RND[]: 0x%02x 0x%02x\n", rnd[1], rnd[0])
	   fmt.Printf("OBF IN[]: %02x\n", in[0])
	   fmt.Printf("OBF OUT[]: 0x%02x 0x%02x\n", out[1], out[0])
	*/
	return
}


func deobfuscateByte(in []byte, mask_hi, mask_lo byte) (out byte, err error) {
	var cnt int

	if len(in) != 2 {
		err = ErrInvalidInput
		return
	}

	for i := 7; i >= 0; i-- {
		if (mask_hi>>uint(i))&0x01 == 0x01 {
			out |= (in[1] >> uint(i)) & 0x01
			if cnt != 7 {
				out <<= 0x01
			}
			cnt++
		}
	}

	for i := 7; i >= 0; i-- {
		if (mask_lo>>uint(i))&0x01 == 0x01 {
			out |= (in[0] >> uint(i)) & 0x01
			//if i != 0 { // NO NEED TO SHIFT ON THE LAST BIT
			if cnt != 7 { // THE EIGHT BIT MASK CAN BE DISTRIBUTED HOW YOU WANT INSIDE MASK_HI MASK_LO
				out <<= 0x01
			}
			cnt++
		}
	}
	return
}

/* can process it with any number of bytes */
func Obfuscate(in []byte) (out []byte, err error) {
	if len(in) == 0 {
		err = ErrInvalidInput
		return
	}

	//out = make([]byte, len(in)*2)
	for i := range in {
		/* mask have to be provided out */
		tmp, derr := obfuscateByte(in[i], 0x55, 0x55)
		if derr != nil {
			err = derr
			return
		}
		out = append(out, tmp...)
	}

	return
}

/* can process it with any number of bytes */
func Deobfuscate(in []byte) (out []byte, err error) {
	if len(in) == 0 || len(in) %2 != 0 {
		err = ErrInvalidInput
		return
	}

	out = make([]byte, len(in)/2)
	for i := range out {
		tmp, derr := deobfuscateByte(in[i*2:(i*2)+2], 0x55, 0x55)
		if derr != nil {
			err = derr
			return
		}
		out[i] = tmp
	}
	return
}
