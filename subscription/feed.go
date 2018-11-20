package subscription

import (
	"bytes"
)

const atomHint = "http://www.w3.org/2005/Atom"
const detectionLimit = 1024

// Determines if the document is an Atom feed by searching for the
// format namespace. Search is given up after 1024 bytes.
func isAtomFeed(doc []byte) bool {
	upper := len(doc)
	if upper > detectionLimit {
		upper = detectionLimit
	}
	return bytes.Contains(doc[:upper], []byte(atomHint))
}
