package shared

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/cyberphone/json-canonicalization/go/src/webpki.org/jsoncanonicalizer"
)

// CanonicalHash returns a SHA256 hash of the canonical JSON representation of v.
func CanonicalHash(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	canonical, err := jsoncanonicalizer.Transform(data)

	if err != nil {
		return "", err
	}

	h := sha256.Sum256(canonical)
	return hex.EncodeToString(h[:]), nil
}

// MerkleTree represents a simple merkle tree.
type MerkleTree struct {
	root *merkleNode
}

type merkleNode struct {
	hash  string
	left  *merkleNode
	right *merkleNode
}

// NewMerkleTree builds a merkle tree from the given hashes.
func NewMerkleTree(hashes []string) *MerkleTree {
	if len(hashes) == 0 {
		return &MerkleTree{}
	}

	// Create leaf nodes
	nodes := make([]*merkleNode, len(hashes))
	for i, hash := range hashes {
		nodes[i] = &merkleNode{hash: hash}
	}

	// Build tree bottom-up
	for len(nodes) > 1 {
		var parentNodes []*merkleNode
		for i := 0; i < len(nodes); i += 2 {
			if i+1 < len(nodes) {
				parentHash := combinedHash(nodes[i].hash, nodes[i+1].hash)
				parentNodes = append(parentNodes, &merkleNode{
					hash:  parentHash,
					left:  nodes[i],
					right: nodes[i+1],
				})
			} else {
				// Odd number of nodes, promote the last node
				parentNodes = append(parentNodes, nodes[i])
			}
		}
		nodes = parentNodes
	}

	return &MerkleTree{root: nodes[0]}
}

// Root returns the root hash of the merkle tree.
func (mt *MerkleTree) Root() string {
	if mt.root == nil {
		return ""
	}
	return mt.root.hash
}

// combinedHash returns the hash of two concatenated hashes.
func combinedHash(left, right string) string {
	combined := left + right
	h := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(h[:])
}
