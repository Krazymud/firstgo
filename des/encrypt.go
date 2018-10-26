package des

//Encrypt using des algo.
func Encrypt(plainText, key []byte) []byte {
	if len(key) != 8 {
		//
	}
	addByte := 8 - len(plainText)%8
	for i := 0; i < addByte; i++ {
		plainText = append(plainText, byte(addByte))
	}

	cipherText := make([]byte, len(plainText))

	leftKey := permute(key, permutedChoice1)[0:4]
	rightKey := permute(key, permutedChoice1)[4:8]

	//handle the plainText block by block
	for i := 0; i < len(plainText); i += 8 {
		block := permute(plainText[i:i+8], initialPermutation)
		left := block[0:4]
		right := block[4:8]
		//iteration
		for j := 0; j < 16; j++ {
			subkey := generateSubkey(leftKey, rightKey, j+1)
			newRight := feistel(left, right, subkey)
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
			cipherText[i+j] = finalBlock[j]
		}
	}
	return cipherText
}

//feistel function
func feistel(left []byte, right []byte, key []byte) []byte {
	expanded := permute(right, expFunction)
	for i := range expanded {
		expanded[i] ^= key[i]
	}
	substitute := make([]byte, 4)

	substitute[0] = (sBoxes[0][(expanded[0]>>2)&0x01|(expanded[0]>>7)][(expanded[0]>>3)&0x0F] << 4) | (sBoxes[1][expanded[0]&0x02|(expanded[1]<<3>>7)][((expanded[0]&0x01)<<3)|(expanded[1]>>5)])
	substitute[1] = (sBoxes[2][(expanded[1]>>2)&0x02|(expanded[2]<<1>>7)][(expanded[1]<<1)&0x0E|(expanded[2]>>7)]) | (sBoxes[3][(expanded[2]&0x01)|(expanded[2]>>4&0x02)][(expanded[2]>>1)&0x0F])
	substitute[2] = (sBoxes[4][(expanded[3]>>2)&0x01|(expanded[3]>>7)][(expanded[3]>>3)&0x0F] << 4) | (sBoxes[5][expanded[3]&0x02|(expanded[4]<<3>>7)][((expanded[3]&0x01)<<3)|(expanded[4]>>5)])
	substitute[3] = (sBoxes[6][(expanded[4]>>2)&0x02|(expanded[5]<<1>>7)][(expanded[4]<<1)&0x0E|(expanded[5]>>7)]) | (sBoxes[7][(expanded[5]&0x01)|(expanded[5]>>4&0x02)][(expanded[5]>>1)&0x0F])

	permuted := permute(substitute, permutation)
	for i := range permuted {
		permuted[i] ^= left[i]
	}
	return permuted
}

//generateSubkey during iteration
func generateSubkey(leftKey, rightKey []byte, turn int) []byte {
	var loop int
	if turn == 1 || turn == 2 || turn == 9 || turn == 16 {
		loop = 1
	} else {
		loop = 2
	}
	for l := 0; l < loop; l++ {
		lHead := leftKey[0] >> 7 << 4
		rHead := (rightKey[0] & 0x0F) >> 3
		for i := 0; i < 3; i++ {
			lShift := leftKey[i+1] >> 7
			rShift := rightKey[i+1] >> 7
			leftKey[i] <<= 1
			rightKey[i] <<= 1
			rightKey[i] |= rShift
			leftKey[i] |= lShift
		}
		leftKey[3] = (leftKey[3]<<1)&0xE0 | lHead
		rightKey[0] &= 0x0F
		rightKey[3] = (rightKey[3] << 1) | rHead
	}
	subKey := make([]byte, 7)
	for i := 0; i < 3; i++ {
		subKey[i] = leftKey[i]
		subKey[i+4] = rightKey[i+1]
	}
	subKey[3] = leftKey[3] | rightKey[0]
	return permute(subKey, permutedChoice2)
}

//permute is used to perform permutations due to table
func permute(bt, table []byte) []byte {
	res := make([]byte, len(table)/8)
	for i, v := range table {
		btIndex := (v - 1) / 8
		btOffset := (v - 1) % 8
		tbIndex := i / 8
		tbOffset := byte(i % 8)
		res[tbIndex] |= ((bt[btIndex] << btOffset) & 0x80) >> tbOffset
	}
	return res
}
