// +build go1.12

package obytes

import (
	"crypto/rand"

	chacha "golang.org/x/crypto/chacha20poly1305"
)

// this just an obfuscation mechanism to hide "structure" in a flow of bytes
// with an overhead, it is in ANY WAY a ways to SECURE DATA.
// obfuscation.
// XXX need to add an offset prefix of real random data
func NewBlob(in []byte, mask_hi, mask_lo byte, prefixlen int) (out []byte, err error) {
	randpfx := make([]byte, prefixlen)
	_, err = rand.Read(randpfx)
	if err != nil {
		err = ErrCrypto
		return
	}

	xkey := make([]byte, chacha.KeySize)
	_, err = rand.Read(xkey)
	if err != nil {
		err = ErrCrypto
		return
	}

	nonce := make([]byte, chacha.NonceSizeX)
	_, err = rand.Read(nonce)
	if err != nil {
		err = ErrCrypto
		return
	}

	oxkey, err := ObfuscateWithMask(xkey, mask_hi, mask_lo)
	if err != nil {
		err = ErrUnexpected
		return
	}

	ononce, err := ObfuscateWithMask(nonce, mask_hi, mask_lo)
	if err != nil {
		err = ErrUnexpected
		return
	}

	aead, err := chacha.NewX(xkey)
	if err != nil {
		err = ErrCrypto
		return
	}

	ct := aead.Seal(nil, nonce, in, nil)

	out = append(out, randpfx...)
	out = append(out, ononce...)
	out = append(out, oxkey...)
	out = append(out, ct...)

	return
}

func OpenBlob(in []byte, mask_hi, mask_lo byte, prefixlen int) (out []byte, err error) {
	ononcelen := chacha.NonceSizeX * 2
	okeylen := chacha.KeySize * 2
	minlen := prefixlen + ononcelen + okeylen + 16

	if len(in) < minlen {
		err = ErrInvalidInput
		return
	}

	ononce := in[prefixlen : prefixlen+ononcelen]
	nonce, err := DeobfuscateWithMask(ononce, mask_hi, mask_lo)
	if err != nil {
		//err = derr
		err = ErrUnexpected
		return
	}

	oxkey := in[prefixlen+ononcelen : prefixlen+ononcelen+okeylen]
	xkey, err := DeobfuscateWithMask(oxkey, mask_hi, mask_lo)
	if err != nil {
		//err = xerr
		err = ErrUnexpected
		return
	}

	aead, err := chacha.NewX(xkey)
	if err != nil {
		//err = cerr
		err = ErrCrypto
		return
	}

	ct := in[prefixlen+ononcelen+okeylen:]
	out, err = aead.Open(nil, nonce, ct, nil)
	if err != nil {
		//err = fmt.Errorf("invalid plural box")
		err = ErrCrypto
		return
	}

	return
}
