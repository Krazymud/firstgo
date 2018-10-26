package des

//Decrypt using des algo.
func Decrypt(cipherText, key []byte) []byte {
	if len(key) != 8 {
		//
	}

	plainText := make([]byte, len(cipherText))
	var subKeys [16][]byte
	leftKey := permute(key, permutedChoice1)[0:4]
	rightKey := permute(key, permutedChoice1)[4:8]
	for i := 0; i < 16; i++ {
		subKeys[i] = generateSubkey(leftKey, rightKey, i+1)
	}

	//handle the plainText block by block
	for i := 0; i < len(cipherText); i += 8 {
		block := permute(cipherText[i:i+8], initialPermutation)
		left := block[0:4]
		right := block[4:8]
		//iteration
		for j := 0; j < 16; j++ {
			newRight := feistel(left, right, subKeys[15-j])
			left = right
			right = newRight
		}
		finalBlock := make([]byte, 8)
		for i := 0; i < 4; i++ {
			finalBlock[i] = right[i]
			finalBlock[i+4] = left[i]
		}
		finalBlock = permute(finalBlock, finalPermutation)
		for j := 0; j < 8; j++ {
			plainText[i+j] = finalBlock[j]
		}
	}
	count, val := 1, 0
	for i := len(plainText) - 1; i >= len(plainText)-7; i-- {
		val = int(plainText[i])
		if plainText[i-1] == plainText[i] {
			count++
		} else {
			break
		}
	}
	if count == val {
		plainText = plainText[:len(plainText)-count]
	}
	return plainText
}
