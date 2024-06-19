package testutils

import (
	"std"
	"strings"
)

// TestAddress returns a fake std.Address that starts with "g1test", fills the middle with "0", and ends with the provided name.
// The provided name avoids using characters not in the charset by replacing them with valid characters.
// The address is limited to the charset: "qpzry9x8gf2tvdw0s3jn54khce6mua7l".
func TestAddress(name string) std.Address {
	const prefix = "g1"
	const addressLength = 31
	const charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

	if len(name) == 0 {
		panic("address name cannot be empty")
	}

	// Convert the provided name to lowercase.
	name = strings.ToLower(name)

	// Replace invalid characters in the provided name to match the charset.
	replacements := map[rune]rune{
		'1': 'l', 'b': '0', 'i': 'l', 'o': '0',
	}
	var builder strings.Builder
	for _, char := range name {
		if replacement, ok := replacements[char]; ok {
			builder.WriteRune(replacement)
		} else if strings.ContainsRune(charset, char) {
			builder.WriteRune(char)
		} else {
			panic("address name contains invalid characters")
		}
	}
	name = builder.String()

	// Calculate the number of zeros needed to fill the middle part.
	zerosLength := addressLength - len(prefix) - len(name) - 6
	if zerosLength < 0 {
		panic("address name is too long")
	}

	checksum := "asdfghf"

	// Create the fake address.
	fakeAddress := prefix + name + strings.Repeat("0", zerosLength) + checksum

	return std.Address(fakeAddress)
}
