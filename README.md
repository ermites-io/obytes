obytes
======

an easy to use Obfuscated Bytes package.


Description
===========

This simple package provide simple bytes obfuscations mechanisms.
* merge data with of a set of random bytes according to a user defined bitmask, a single byte becomes two obfuscated bytes.
* random obfuscated key, random obfuscated nonce, xchacha20+poly1305.


How the obfuscation works
=========================

"Masks" defines how random bits are distributed in each obfuscated byte.
The default mask is 0x55 : [0101 0101]

The byte obfuscation transform one data byte nibble into one obfuscated byte, such that
* bits that are 0 become random bits.
* bits that are 1 become data bits.

You need a mask byte for one byte nibble of data, hence user provided masks requires a nibble high mask and nibble low mask.
!IMPORTANT! you should ALWAYS make sure to have evenly distributed random/data bits in each mask.

Carefully chosen masks allows to recover your byte out of 2 obfuscated bytes.


Examples
========

Obfuscate bytes using default mask:   
```go
	obfuscated, err := Obfuscate(input)
	if err != nil {
		// handle error
	}
	// obfuscated is len(input)*2 and is obfuscated

	deobfuscated_input, err := Deobfuscate(obfuscated)
	if err != nil {
		// handle error
	}

```


Obfuscate bytes using a user defined mask:
```go
	obfuscated, err := ObfuscateWithMask(input, 0xf0, 0x0f)
	if err != nil {
		// handle error
	}
	// obfuscated is len(input)*2 and is obfuscated

	deobfuscated_input, err := DeobfuscateWithMask(obfuscated, 0xf0, 0x0f)
	if err != nil {
		// handle error
	}
```


Obfuscate using an encrypted blob:
```go
	obfuscated, err := NewBlob(input, 0x55, 0x55, 234)
	if err != nil {
		// handle error
	}
	// obfuscated is a blob of len(input)+(keysize*2)+(noncesize*2)+(aead overhead)+234
	// obfuscated is prefixed by 234 random bytes.
```
