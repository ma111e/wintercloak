package encryption

func Xor(key []byte, cleartext []byte) []byte {
	res := make([]byte, len(cleartext))

	for idx, b := range cleartext {
		res[idx] = b ^ key[idx%len(key)]
	}

	return res
}
