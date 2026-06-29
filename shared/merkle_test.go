package shared

import "testing"

func TestMerkleRootFromHashesIsOrderIndependent(t *testing.T) {
	hashes := []string{"c", "a", "b"}
	reordered := []string{"b", "c", "a"}

	got := MerkleRootFromHashes(hashes)
	gotReordered := MerkleRootFromHashes(reordered)

	if got != gotReordered {
		t.Fatalf("expected root to be independent of input order, got %q and %q", got, gotReordered)
	}
	if got == "" {
		t.Fatal("expected a non-empty root for non-empty input")
	}
}

func TestMerkleRootFromHashesDoesNotMutateInput(t *testing.T) {
	hashes := []string{"c", "a", "b"}
	_ = MerkleRootFromHashes(hashes)

	if hashes[0] != "c" || hashes[1] != "a" || hashes[2] != "b" {
		t.Fatalf("expected input slice to be left untouched, got %v", hashes)
	}
}

func TestMerkleRootFromHashesEmpty(t *testing.T) {
	if got := MerkleRootFromHashes(nil); got != "" {
		t.Fatalf("expected empty root for no hashes, got %q", got)
	}
}

func TestMerkleRootFromHashesDetectsTampering(t *testing.T) {
	original := MerkleRootFromHashes([]string{"a", "b", "c"})
	tampered := MerkleRootFromHashes([]string{"a", "b", "d"})

	if original == tampered {
		t.Fatal("expected changing a leaf hash to change the root")
	}
}
