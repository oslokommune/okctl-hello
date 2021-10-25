package communicationtest

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const URLToTest = "http://tempo.monitoring:3100"

// HandlerFunc returns the health handler
func HandlerFunc() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := handle(writer)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write([]byte(fmt.Sprintf("Error while doing communication test: %s", err.Error())))
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func handle(writer io.Writer) error {
	response := strings.Builder{}

	client := &http.Client{
		Timeout: 5 * time.Second, // nolint: gomnd
	}

	req, err := http.NewRequest(http.MethodGet, URLToTest, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	response.WriteString(fmt.Sprintf("Doing GET to %s\n", URLToTest))

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("doing request: %w", err)
	}

	response.WriteString(fmt.Sprintf("Response status code: %d\n", resp.StatusCode))
	response.WriteString(fmt.Sprintf("Response status: %s\n", resp.Status))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	response.WriteString(fmt.Sprintf("Response body:\n%s\n", body))

	_, err = writer.Write([]byte(response.String()))
	if err != nil {
		return fmt.Errorf("creating response string: %w", err)
	}

	return nil
}
