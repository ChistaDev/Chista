package helpers

import (
	"fmt"
	"testing"

	"github.com/Chista-Framework/Chista/logger"
	"github.com/stretchr/testify/assert"
)

func TestVerifyCensysCredentials(t *testing.T) {
	api_id := "test"
	secret := "test"
	InitiliazeWebSocketConnection()
	VerifyCensysCredentials(api_id, secret)
}

// Tests curl -X 'GET' 'https://search.censys.io/api/v2/certificates/search?q=xn--p-bga9c0ugb.com&per_page=100' -H 'accept: application/json' -H 'Authorization: Basic <Token>'
func TestApiRequester(t *testing.T) {
	a := assert.New(t)
	input := []string{
		"paypal.com",
		"sibersaldirilar.com",
	}
	var extracted_domains_from_ct []string

	for i := 0; i < len(input); i++ {
		/*url := "https://search.censys.io/api/v2/certificates/search?per_page=100&q=" + input[i]
		auth_key := "Authorization"
		auth_value := "Basic " + base64.StdEncoding.EncodeToString([]byte(GoDotEnvVariable("CENSYS_API_ID")+":"+GoDotEnvVariable("CENSYS_API_SECRET")))
		*/
		logger.Log.Infof("Looking for CT Transparency logs for %s...", input[i])

		/*err, censys_resp_model := ApiRequester(url, "GET", auth_key, auth_value, nil)
		if err != nil {
			logger.Log.Errorf("APIRequster error  %v...", err)
		}*/

		// Extract the domain name
		/*hits := censys_resp_model
		for _, hit := range hits {
			extracted_domains_from_ct = append(extracted_domains_from_ct, hit.Names...)
		}
		*/
		// Add cursor if there is any new page, (collect data for all pages)
		/*cursor := censys_resp_model.Result.Links.Next
		for cursor != "" {
			url = "https://search.censys.io/api/v2/certificates/search?per_page=100&q=" + input[i] + "&cursor=" + cursor
			err, censys_resp_model = ApiRequester(url, "GET", auth_key, auth_value, nil)
			if err != nil {
				logger.Log.Errorf("APIRequster error  %v...", err)
			}

			// Extract the domain name
			hits := censys_resp_model.Result.Hits
			for _, hit := range hits {
				extracted_domains_from_ct = append(extracted_domains_from_ct, hit.Names...)
			}
		}
		*/
		logger.Log.Infof("CT Transparency  logs collected for %s", input[i])

	}
	fmt.Printf("\n\nEXTRACTED DOMAINS: %s\n\n", extracted_domains_from_ct)
	a.NotNil(extracted_domains_from_ct)

	// Make unique the array
	uniqued_domains_ct_logs := UniqueStrArray(extracted_domains_from_ct)
	fmt.Printf("\n\n UNIQUE DOMAINS: %s\n\n", uniqued_domains_ct_logs)
}

func TestConvertToPunnyCodeDomain(t *testing.T) {
	a := assert.New(t)
	input := "süper.com"

	converted_domain, err := ConvertToPunnyCodeDomain(input)
	if err != nil {
		a.Error(err)
	}

	fmt.Printf("\t[T] Converted domain: %v\n", converted_domain)
	a.NotNil(converted_domain, "Converted domain: %v", converted_domain)

}

func TestGenerateDomainsWithUnsupportedChars(t *testing.T) {
	input := "super"
	t.Run(input, func(t *testing.T) {
		results := GenerateDomainsWithUnsupportedChars(input)
		t.Logf("Result strings: %v", results)
		if len(results) == 0 {
			t.Errorf("Expected a non-empty list, but got an empty list")
		}
	})

}

func TestPsuedoString(t *testing.T) {
	testCases := []struct {
		input string
	}{
		{
			input: "microsoft",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			results := PsuedoString(tc.input)
			t.Logf("Result strings: %v", results)
			if len(results) == 0 {
				t.Errorf("Expected a non-empty list, but got an empty list")
			}

		})
	}

}
