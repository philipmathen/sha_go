package main

import (
	"fmt"
	"os"
	"strconv"
	// "error"
)

func main() {
	var arguments []string = os.Args[1:]
	var message string

	for i := range arguments {
		message += arguments[i]
	}
	//initialer State
	var h0 uint32 = 0x67452301
	var h1 uint32 = 0xEFCDAB89
	var h2 uint32 = 0x98BADCFE
	var h3 uint32 = 0x10325476
	var h4 uint32 = 0xC3D2E1F0

	message = convertStringToBin(message)
	message = padMessage(message)
	//aufteilen der message in 512bit Blöcke à 16*32bits
	numberOfBlocks := len(message) / 512
	var messageSlice [][]string
	for i := 0; i < numberOfBlocks; i++ {
		messageSlice = append(messageSlice, []string{})
		for j := 0; j < 512; j += 32 {
			messageSlice[i] = append(messageSlice[i], message[j:j+32])
		}
	}
	//für jeden 512bit block
	for i := range messageSlice {
		//erweitern von 16 auf 80 Wörter
		//fmt.Println(messageSlice)
		for j := 16; j < 80; j++ {
			word0, _ := strconv.ParseUint(messageSlice[i][j-3], 2, 32)
			word1, _ := strconv.ParseUint(messageSlice[i][j-8], 2, 32)
			word2, _ := strconv.ParseUint(messageSlice[i][j-14], 2, 32)
			word3, _ := strconv.ParseUint(messageSlice[i][j-16], 2, 32)

			var word uint32 = leftrotateInteger(uint32(word0)^uint32(word1)^uint32(word2)^uint32(word3), 1)
			messageSlice[i] = append(messageSlice[i], fmt.Sprintf("%.32b", word))
		}
	}

	for _, v := range messageSlice {
		var a uint32 = h0
		var b uint32 = h1
		var c uint32 = h2
		var d uint32 = h3
		var e uint32 = h4
		var f uint32
		var k uint32
		for i, _ := range v {

			if i >= 0 && i < 20 {
				f = (b & c) | (^b & d)
				k = 0x5A827999
			} else if i >= 20 && i < 40 {
				f = b ^ c ^ d
				k = 0x6ED9EBA1
			} else if i >= 40 && i < 60 {
				f = (b & c) | (b & d) | (c & d)
				k = 0x8F1BBCDC
			} else if i >= 60 && i < 80 {
				f = b ^ c ^ d
				k = 0xCA62C1D6
			}
			//fmt.Printf("%v | val -> %v\n", i, val)
			wordI, err := strconv.ParseUint(v[i], 2, 32)
			if err != nil {
				fmt.Println(err.Error())
			}
			temp := leftrotateInteger(a, 5) + f + e + k + uint32(wordI)

			e = d
			d = c
			c = leftrotateInteger(b, 30)
			b = a
			a = temp

		}

		h0 += a
		h1 += b
		h2 += c
		h3 += d
		h4 += e
	}

	fmt.Printf("%x", h0)
	fmt.Printf("%x", h1)
	fmt.Printf("%x", h2)
	fmt.Printf("%x", h3)
	fmt.Printf("%x", h4)

}

// convertStringToBin takes a String and returns the binary representation
func convertStringToBin(param string) (binRespresentation string) {
	for _, v := range param {
		binRespresentation += fmt.Sprintf("%.8b", v)
	}
	fmt.Printf("binary of message: %v\n", binRespresentation)
	return
}

// padMessage takes the binary string and returns it with the padding
func padMessage(param string) (paddedMessage string) {
	paddedMessage = param
	var messageLength int = len(paddedMessage)
	paddedMessage += "1"
	for len(paddedMessage)%512 != 448 {
		paddedMessage += "0"
	}
	paddedMessage += fmt.Sprintf("%.64b", messageLength)
	//fmt.Printf("padded msg length: %v\n", len(paddedMessage))
	return
}

func leftrotateArray[T any](param []T, n int) (resultArray []T) {
	resultArray = param
	var firstItem T
	var length = len(param)
	for i := 0; i < n; i++ {
		for j := 0; j < length; j++ {
			if j == 0 {
				firstItem = resultArray[0]
			} else if j == length-1 {
				resultArray[length-1] = firstItem
				break
			}
			resultArray[j] = resultArray[j+1]
		}
	}
	return
}

func leftrotateInteger(param uint32, n int) (result uint32) {
	paramBinaryRuneSlice := []rune(fmt.Sprintf("%.32b", param))
	rotatedSlice := leftrotateArray(paramBinaryRuneSlice, n)
	res, err := strconv.ParseUint(string(rotatedSlice), 2, 32)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
	}
	result = uint32(res)
	return
}
