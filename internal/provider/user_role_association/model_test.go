package user_role_association

import (
	"encoding/json"
	"testing"
)

func TestCreateUserRoleAssociationAPIResponseUnmarshal(t *testing.T) {
	sample := `{
		"object_version": 1,
		"creation_user_id": "creator-123",
		"update_user_id": null,
		"delete_user_id": null,
		"creation_timestamp": "2025-11-03T12:00:00Z",
		"update_timestamp": null,
		"user_role_id": "role-abc",
		"user_id": "user-xyz",
		"expiration_timestamp": null
	}`

	var resp CreateUserRoleAssociationAPIResponse
	if err := json.Unmarshal([]byte(sample), &resp); err != nil {
		t.Fatalf("failed to unmarshal sample response: %v", err)
	}

	if resp.UserRoleId != "role-abc" {
		t.Fatalf("unexpected user_role_id: %s", resp.UserRoleId)
	}
	if resp.UserID != "user-xyz" {
		t.Fatalf("unexpected user_id: %s", resp.UserID)
	}
	if resp.CreationUserID != "creator-123" {
		t.Fatalf("unexpected creation_user_id: %s", resp.CreationUserID)
	}
	if resp.ObjectVersion != 1 {
		t.Fatalf("unexpected object_version: %d", resp.ObjectVersion)
	}
	if resp.CreationTimestamp != "2025-11-03T12:00:00Z" {
		t.Fatalf("unexpected creation_timestamp: %s", resp.CreationTimestamp)
	}
	if resp.UpdateUserID != nil {
		t.Fatalf("expected update_user_id to be nil, got: %v", resp.UpdateUserID)
	}
	if resp.UpdateTimestamp != nil {
		t.Fatalf("expected update_timestamp to be nil, got: %v", resp.UpdateTimestamp)
	}
	if resp.ExpirationTimestamp != nil {
		t.Fatalf("expected expiration_timestamp to be nil, got: %v", resp.ExpirationTimestamp)
	}
}
