package scheduler

import (
	"testing"

	"github.com/andrewrutherfoord/fed-bill-poc/shared"
	sharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models"
)

func TestVerifyProviderLineItemMerkleRoot(t *testing.T) {
	batches := []sharedmodels.ChargeBatch{
		{MerkleRoot: "batch-root-2"},
		{MerkleRoot: "batch-root-1"},
		{MerkleRoot: "batch-root-3"},
	}
	claimedRoot := shared.MerkleRootFromHashes([]string{"batch-root-1", "batch-root-2", "batch-root-3"})

	if !verifyProviderLineItemMerkleRoot(batches, claimedRoot) {
		t.Fatal("expected line item root built from the same batch roots to verify, regardless of fetch order")
	}
}

func TestVerifyProviderLineItemMerkleRootDetectsMismatch(t *testing.T) {
	batches := []sharedmodels.ChargeBatch{
		{MerkleRoot: "batch-root-1"},
		{MerkleRoot: "batch-root-2"},
	}

	if verifyProviderLineItemMerkleRoot(batches, "some-other-claimed-root") {
		t.Fatal("expected a line item root that doesn't match any combination of the batch roots to fail verification")
	}
}

func TestVerifyProviderLineItemMerkleRootDetectsMissingBatch(t *testing.T) {
	allBatches := []sharedmodels.ChargeBatch{
		{MerkleRoot: "batch-root-1"},
		{MerkleRoot: "batch-root-2"},
		{MerkleRoot: "batch-root-3"},
	}
	claimedRoot := shared.MerkleRootFromHashes([]string{"batch-root-1", "batch-root-2", "batch-root-3"})

	// Simulate a batch the billing provider didn't report (or that was dropped/tampered with).
	incompleteBatches := allBatches[:2]

	if verifyProviderLineItemMerkleRoot(incompleteBatches, claimedRoot) {
		t.Fatal("expected verification to fail when a batch backing the claimed root is missing")
	}
}
