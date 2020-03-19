package testing

import (
	"encoding/json"
	"fmt"
	fake "bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper/client"
	"net/http"
	"testing"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/identity/tokens"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/testhelper"
)

// TokenOutput is a sample response to a AccessToken call.
var TokenOutput = fmt.Sprintf(`
{
   "access": "%s",
   "refresh": "%s"
}`, fake.AccessToken,
	fake.RefreshToken,
)

var expectedToken = tokens.Token{
	Access:  fake.AccessToken,
	Refresh: fake.RefreshToken,
}

func getGetResult(t *testing.T) tokens.TokenResult {
	result := tokens.TokenResult{}
	result.Header = http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", fake.AccessToken)},
	}
	err := json.Unmarshal([]byte(TokenOutput), &result.Body)
	testhelper.AssertNoErr(t, err)
	return result
}
