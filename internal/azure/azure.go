package azure

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

func DeleteSpec(id string) error {
	var client = &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, "https://hka-lab-api-management.management.azure-api.net/subscriptions/1ed07f18-4c8e-424d-bd80-19e1247574ad/resourceGroups/RG-APIGenerator/providers/Microsoft.ApiManagement/service/hka-lab-api-management/apis/"+id+"?api-version=2021-08-01", nil)
	if err != nil {
		return err
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "SharedAccessSignature "+os.Getenv("AZURE_AUTH_KEY"))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	if resp.StatusCode != 200 && resp.StatusCode != 204 && resp.StatusCode != 202 {
		return fmt.Errorf(string(body))
	}
	return nil
}

func UploadSpec(id string, bodyString string) error {
	var client = &http.Client{}
	req, err := http.NewRequest(http.MethodPut, "https://hka-lab-api-management.management.azure-api.net/subscriptions/1ed07f18-4c8e-424d-bd80-19e1247574ad/resourceGroups/RG-APIGenerator/providers/Microsoft.ApiManagement/service/hka-lab-api-management/apis/"+id+"?api-version=2021-08-01", bytes.NewBuffer([]byte(bodyString)))
	if err != nil {
		return err
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "SharedAccessSignature "+os.Getenv("AZURE_AUTH_KEY"))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	if resp.StatusCode != 200 && resp.StatusCode != 204 && resp.StatusCode != 202 {
		return fmt.Errorf(string(body))
	}
	return nil
}
