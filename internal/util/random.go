package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func RandomUsername() string {
	return RandomString(10)
}
func RandomEmail() string {
	return fmt.Sprintf("%s@%s.com", RandomString(6), RandomString(5))
}

var randGen *rand.Rand

func init() {
	randSrc := rand.NewSource(time.Now().UnixNano())
	randGen = rand.New(randSrc)
}

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomString(n int) string {
	strb := strings.Builder{}
	strb.Grow(n)

	for range n {
		err := strb.WriteByte(
			alphabet[randGen.Intn(len(alphabet))],
		)
		if err != nil {
			panic(fmt.Sprintf("failed to write byte: %v", err))
		}
	}
	return strb.String()
}

func RandomInt(min, max int64) int64 {
	return min + randGen.Int63n(max-min+1)
}

func RandomPhone() string {
	number := fmt.Sprintf("+%d%d%d%d%d%d%d%d%d%d",
		RandomInt(1, 9),
		RandomInt(0, 9), RandomInt(0, 9), RandomInt(0, 9),
		RandomInt(0, 9), RandomInt(0, 9), RandomInt(0, 9),
		RandomInt(0, 9), RandomInt(0, 9), RandomInt(0, 9),
	)
	return number
}
