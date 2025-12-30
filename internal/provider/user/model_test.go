package user

import (
    "encoding/json"
    "testing"
)

func TestCreateUserAPIResponseUnmarshal(t *testing.T) {
    sample := `{
        "object_version": 1,
        "creation_user_id": "creator-123",
        "update_user_id": null,
        "creation_timestamp": "2025-11-03T12:00:00Z",
        "update_timestamp": null,
        "user_id": "user-abc",
        "group_id": "group-1",
        "first_name": "John",
        "last_name": "Doe",
        "email": "john.doe@example.com",
        "company_name": "Acme Corp",
        "phone_number": "+441234567890",
        "privacy_accepted": true,
        "locale": "en-GB",
        "last_login_timestamp": null,
        "title": "Engineer",
        "agree_to_contact": true,
        "source": "SDA"
    }`

    var resp CreateUserAPIResponse
    if err := json.Unmarshal([]byte(sample), &resp); err != nil {
        t.Fatalf("failed to unmarshal sample response: %v", err)
    }

    if resp.UserID != "user-abc" {
        t.Fatalf("unexpected user_id: %s", resp.UserID)
    }
    if resp.FirstName != "John" || resp.LastName != "Doe" {
        t.Fatalf("unexpected name: %s %s", resp.FirstName, resp.LastName)
    }
    if resp.Source != "SDA" {
        t.Fatalf("unexpected source: %s", resp.Source)
    }
}
