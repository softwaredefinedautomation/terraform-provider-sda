package role

import (
    "encoding/json"
    "testing"
)

func TestCreateRoleAPIResponseUnmarshal(t *testing.T) {
    sample := `{
        "object_version": 1,
        "creation_user_id": "creator-123",
        "update_user_id": null,
        "creation_timestamp": "2025-11-03T12:00:00Z",
        "update_timestamp": null,
        "user_role_id": "role-abc",
        "name": "editor",
        "group_id": "group-1",
        "description": "Editor role",
        "policies": [],
        "is_system_role": false,
        "sso_group_mapping": ["sso-group-1"]
    }`

    var resp CreateRoleAPIResponse
    if err := json.Unmarshal([]byte(sample), &resp); err != nil {
        t.Fatalf("failed to unmarshal sample response: %v", err)
    }

    if resp.UserRoleID != "role-abc" {
        t.Fatalf("unexpected user_role_id: %s", resp.UserRoleID)
    }
    if resp.Name != "editor" {
        t.Fatalf("unexpected name: %s", resp.Name)
    }
}
