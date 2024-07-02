package urlgenerator

import (
	"fmt"
	"hash/fnv"
)

// HashURL generates a shortened version of the given URL
// using the FNV-1a hash function.
func HashURL(url string) string {
	hasher := fnv.New32a()
	hasher.Write([]byte(url))
	return fmt.Sprintf("%x", hasher.Sum32())
}
