package user

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "strings"

    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/types"

)

// Ensure UserResource implements CRUD interfaces
var _ resource.Resource = &UserResource{}

// CREATE
func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan UserResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    payload := map[string]interface{}{
        "first_name": plan.FirstName.ValueString(),
        "last_name":  plan.LastName.ValueString(),
        "email":      plan.Email.ValueString(),
    }

    if !plan.GroupID.IsUnknown() && !plan.GroupID.IsNull() {
        payload["group_id"] = plan.GroupID.ValueString()
    }
    if !plan.CompanyName.IsUnknown() && !plan.CompanyName.IsNull() {
        payload["company_name"] = plan.CompanyName.ValueString()
    }
    if !plan.PhoneNumber.IsUnknown() && !plan.PhoneNumber.IsNull() {
        payload["phone_number"] = plan.PhoneNumber.ValueString()
    }
    if !plan.PrivacyAccepted.IsUnknown() && !plan.PrivacyAccepted.IsNull() {
        payload["privacy_accepted"] = plan.PrivacyAccepted.ValueBool()
    }
    if !plan.Locale.IsUnknown() && !plan.Locale.IsNull() {
        payload["locale"] = plan.Locale.ValueString()
    }
    if !plan.Title.IsUnknown() && !plan.Title.IsNull() {
        payload["title"] = plan.Title.ValueString()
    }
    if !plan.AgreeToContact.IsUnknown() && !plan.AgreeToContact.IsNull() {
        payload["agree_to_contact"] = plan.AgreeToContact.ValueBool()
    }

    body, err := json.Marshal(payload)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error marshalling payload: %s", err))
        return
    }

    url := fmt.Sprintf("%s/ident/v1/user", r.client.HostURL)
    reqHTTP, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating request: %s", err))
        return
    }
    reqHTTP.Header.Set("Content-Type", "application/json")

    resBody, err := r.client.DoRequest(reqHTTP, nil)
    if err != nil {
        errSummary := "API Error"
        if strings.Contains(err.Error(), "status: 403") {
            errSummary = "API Error: Forbidden (Check Authentication and Permissions/Token Scope)"
        }

        resp.Diagnostics.AddError(errSummary, fmt.Sprintf("POST %s failed: %s\nResponse body:\n%s", reqHTTP.URL.String(), err, string(resBody)))
        return
    }

    var apiResp CreateUserAPIResponse
    if err := json.Unmarshal(resBody, &apiResp); err != nil {
        resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding response: %s", err))
        return
    }

    // Build state
    state := UserResourceModel{
        ObjectVersion:     types.Int64Value(apiResp.ObjectVersion),
        CreationUserID:    types.StringValue(apiResp.CreationUserID),
        UpdateUserID:      types.StringPointerValue(apiResp.UpdateUserID),
        CreationTimestamp: types.StringValue(apiResp.CreationTimestamp),
        UpdateTimestamp:   types.StringPointerValue(apiResp.UpdateTimestamp),
        UserID:            types.StringValue(apiResp.UserID),
        GroupID:           types.StringPointerValue(apiResp.GroupID),
        FirstName:         types.StringValue(apiResp.FirstName),
        LastName:          types.StringValue(apiResp.LastName),
        Email:             types.StringValue(apiResp.Email),
        CompanyName:       types.StringPointerValue(apiResp.CompanyName),
        PhoneNumber:       types.StringPointerValue(apiResp.PhoneNumber),
        PrivacyAccepted:   types.BoolPointerValue(apiResp.PrivacyAccepted),
        Locale:            types.StringPointerValue(apiResp.Locale),
        LastLoginTimestamp: types.StringPointerValue(apiResp.LastLoginTimestamp),
        Title:             types.StringPointerValue(apiResp.Title),
        AgreeToContact:    types.BoolPointerValue(apiResp.AgreeToContact),
        Source:            types.StringValue(apiResp.Source),
    }

    // Keep the plan values for fields that are not returned by API (if any)
    // e.g. if client provided certain optional fields, preserve them in state
    if !plan.GroupID.IsUnknown() && !plan.GroupID.IsNull() {
        state.GroupID = plan.GroupID
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// READ
func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state UserResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    url := fmt.Sprintf("%s/ident/v1/user/%s", r.client.HostURL, state.UserID.ValueString())
    reqHTTP, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating read request: %s", err))
        return
    }

    resBody, err := r.client.DoRequest(reqHTTP, nil)
    if err != nil {
        if strings.Contains(err.Error(), "status: 404") {
            resp.State.RemoveResource(ctx)
            return
        }
        resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error reading user %s: %s\nResponse body:\n%s", state.UserID.ValueString(), err, string(resBody)))
        return
    }

    var apiResp UserAPIResponse
    if err := json.Unmarshal(resBody, &apiResp); err != nil {
        resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding response: %s", err))
        return
    }

    state.ObjectVersion = types.Int64Value(apiResp.ObjectVersion)
    state.CreationUserID = types.StringValue(apiResp.CreationUserID)
    state.UpdateUserID = types.StringPointerValue(apiResp.UpdateUserID)
    state.CreationTimestamp = types.StringValue(apiResp.CreationTimestamp)
    state.UpdateTimestamp = types.StringPointerValue(apiResp.UpdateTimestamp)
    state.UserID = types.StringValue(apiResp.UserID)
    state.GroupID = types.StringPointerValue(apiResp.GroupID)
    state.FirstName = types.StringValue(apiResp.FirstName)
    state.LastName = types.StringValue(apiResp.LastName)
    state.Email = types.StringValue(apiResp.Email)
    state.CompanyName = types.StringPointerValue(apiResp.CompanyName)
    state.PhoneNumber = types.StringPointerValue(apiResp.PhoneNumber)
    state.PrivacyAccepted = types.BoolPointerValue(apiResp.PrivacyAccepted)
    state.Locale = types.StringPointerValue(apiResp.Locale)
    state.LastLoginTimestamp = types.StringPointerValue(apiResp.LastLoginTimestamp)
    state.Title = types.StringPointerValue(apiResp.Title)
    state.AgreeToContact = types.BoolPointerValue(apiResp.AgreeToContact)
    state.Source = types.StringValue(apiResp.Source)

    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// UPDATE
func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan, state UserResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    payload := map[string]interface{}{
        "object_version": state.ObjectVersion.ValueInt64(),
    }

    if !plan.GroupID.Equal(state.GroupID) {
        if plan.GroupID.IsNull() {
            payload["group_id"] = ""
        } else {
            payload["group_id"] = plan.GroupID.ValueString()
        }
    }
    if !plan.FirstName.Equal(state.FirstName) {
        payload["first_name"] = plan.FirstName.ValueString()
    }
    if !plan.LastName.Equal(state.LastName) {
        payload["last_name"] = plan.LastName.ValueString()
    }
    if !plan.Email.Equal(state.Email) {
        payload["email"] = plan.Email.ValueString()
    }
    if !plan.CompanyName.Equal(state.CompanyName) {
        if plan.CompanyName.IsNull() {
            payload["company_name"] = ""
        } else {
            payload["company_name"] = plan.CompanyName.ValueString()
        }
    }
    if !plan.PhoneNumber.Equal(state.PhoneNumber) {
        if plan.PhoneNumber.IsNull() {
            payload["phone_number"] = ""
        } else {
            payload["phone_number"] = plan.PhoneNumber.ValueString()
        }
    }
    if !plan.PrivacyAccepted.Equal(state.PrivacyAccepted) {
        payload["privacy_accepted"] = plan.PrivacyAccepted.ValueBool()
    }
    if !plan.Locale.Equal(state.Locale) {
        if plan.Locale.IsNull() {
            payload["locale"] = ""
        } else {
            payload["locale"] = plan.Locale.ValueString()
        }
    }
    if !plan.Title.Equal(state.Title) {
        if plan.Title.IsNull() {
            payload["title"] = ""
        } else {
            payload["title"] = plan.Title.ValueString()
        }
    }
    if !plan.AgreeToContact.Equal(state.AgreeToContact) {
        payload["agree_to_contact"] = plan.AgreeToContact.ValueBool()
    }

    body, err := json.Marshal(payload)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error marshalling payload: %s", err))
        return
    }

    url := fmt.Sprintf("%s/ident/v1/user/%s", r.client.HostURL, state.UserID.ValueString())
    reqHTTP, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating update request: %s", err))
        return
    }
    reqHTTP.Header.Set("Content-Type", "application/json")

    resBody, err := r.client.DoRequest(reqHTTP, nil)
    if err != nil {
        resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error updating user: %s\nResponse body:\n%s", err, string(resBody)))
        return
    }

    var apiResp UserAPIResponse
    if err := json.Unmarshal(resBody, &apiResp); err != nil {
        resp.Diagnostics.AddError("Decode Error", fmt.Sprintf("Error decoding response: %s", err))
        return
    }

    // Preserve nulls
    if plan.GroupID.IsNull() {
        apiResp.GroupID = nil
    }

    // Update state
    state.ObjectVersion = types.Int64Value(apiResp.ObjectVersion)
    state.CreationUserID = types.StringValue(apiResp.CreationUserID)
    state.UpdateUserID = types.StringPointerValue(apiResp.UpdateUserID)
    state.CreationTimestamp = types.StringValue(apiResp.CreationTimestamp)
    state.UpdateTimestamp = types.StringPointerValue(apiResp.UpdateTimestamp)
    state.UserID = types.StringValue(apiResp.UserID)
    state.GroupID = types.StringPointerValue(apiResp.GroupID)
    state.FirstName = types.StringValue(apiResp.FirstName)
    state.LastName = types.StringValue(apiResp.LastName)
    state.Email = types.StringValue(apiResp.Email)
    state.CompanyName = types.StringPointerValue(apiResp.CompanyName)
    state.PhoneNumber = types.StringPointerValue(apiResp.PhoneNumber)
    state.PrivacyAccepted = types.BoolPointerValue(apiResp.PrivacyAccepted)
    state.Locale = types.StringPointerValue(apiResp.Locale)
    state.LastLoginTimestamp = types.StringPointerValue(apiResp.LastLoginTimestamp)
    state.Title = types.StringPointerValue(apiResp.Title)
    state.AgreeToContact = types.BoolPointerValue(apiResp.AgreeToContact)
    state.Source = types.StringValue(apiResp.Source)

    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// DELETE
func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state UserResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    url := fmt.Sprintf("%s/ident/v1/user/%s", r.client.HostURL, state.UserID.ValueString())
    reqHTTP, err := http.NewRequest(http.MethodDelete, url, nil)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Error creating delete request: %s", err))
        return
    }

    _, err = r.client.DoRequest(reqHTTP, nil)
    if err != nil {
        if strings.Contains(err.Error(), "status: 404") {
            // already gone
            return
        }
        resp.Diagnostics.AddError("API Error", fmt.Sprintf("Error deleting user: %s", err))
    }
}
