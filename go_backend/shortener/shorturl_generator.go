package shortener

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"
	"github.com/itchyny/base58-go"
)

// generates and return hash of input string
func sha256Of(input string) []byte {
	algo := sha256.New()
	algo.Write([]byte(input))
	return algo.Sum(nil)
}

func base58Encoded(bytes []byte) string {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(encoded)
}

// GenerateShortLink is the production entry point. It reads the clock to build
// a salt, then delegates to the pure generateShortLink. Because it touches the
// real clock it is non-deterministic and is NOT the thing we unit-test.
func GenerateShortLink(initialLink string) string {
	salt := strconv.FormatInt(time.Now().UnixNano(), 10)
	return generateShortLink(initialLink, salt)
}

// generateShortLink is pure: the same initialLink + salt always produce the
// same short link. Tests call this directly with a fixed salt so results are
// deterministic and can be asserted against known values.
func generateShortLink(initialLink string, salt string) string {
	urlHashBytes := sha256Of(initialLink + salt)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	return finalString[:8]
}