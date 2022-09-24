package common

import (
	"fmt"
	"math"
	"strings"

	"github.com/samber/do"
)

type HashCoder struct {
}

func ConstructHashCoder(injector *do.Injector) (*HashCoder, error) {
	encodingTable = make(map[string]int, len(decodingTable))

	for key, value := range decodingTable {
		encodingTable[value] = key
	}

	return &HashCoder{}, nil
}

func (t *HashCoder) Encode(target int) string {
	hash := ""

	denominator := len(decodingTable)
	quotient := target
	for {
		if quotient == 0 {
			break
		}

		remainder := quotient % denominator
		quotient = quotient / denominator

		hash = decodingTable[remainder] + hash
	}

	if MINIMAL_HASH_LENGTH <= len(hash) {
		return hash
	}

	return fmt.Sprintf("%0*s", MINIMAL_HASH_LENGTH, hash)
}

func (t *HashCoder) Decode(target string) int {
	hash := strings.TrimLeft(target, "0")
	if hash == "" {
		return 0
	}

	result := 0
	base := float64(len(encodingTable))
	for i := 0; i < len(hash); i++ {
		encoded := encodingTable[string(hash[len(hash)-i-1])]
		result += encoded * int(math.Pow(base, float64(i)))
	}

	return result
}

const MINIMAL_HASH_LENGTH int = 4

var encodingTable map[string]int
var decodingTable = map[int]string{
	0:  "0",
	1:  "1",
	2:  "2",
	3:  "3",
	4:  "4",
	5:  "5",
	6:  "6",
	7:  "7",
	8:  "8",
	9:  "9",
	10: "A",
	11: "B",
	12: "C",
	13: "D",
	14: "E",
	15: "F",
	16: "G",
	17: "H",
	18: "I",
	19: "J",
	20: "K",
	21: "L",
	22: "M",
	23: "N",
	24: "O",
	25: "P",
	26: "Q",
	27: "R",
	28: "S",
	29: "T",
	30: "U",
	31: "V",
	32: "W",
	33: "X",
	34: "Y",
	35: "Z",
}
