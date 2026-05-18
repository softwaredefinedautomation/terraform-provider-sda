package clients

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// SelectTenant - Select tenant for the current session, return 204 if successful
func (c *Client) SelectTenant(tenantId string, authToken *string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/ident/v1/tenant/select/%s", c.HostURL, tenantId), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", *authToken)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return nil
}
