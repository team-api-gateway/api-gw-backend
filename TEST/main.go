package main

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement"
)

func main() {
	credential, err := azidentity.NewClientSecretCredential("a8cea0e8-2d5f-4772-9440-2820c0ef44b9", "612291ac-20cc-43c9-bf94-3bc78c928fcf", "ubf7Q~QAq6WYS6HiX2pR7J.orjlWhoTlExFH4", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	/*credential, err := azidentity.NewManagedIdentityCredential(&azidentity.ManagedIdentityCredentialOptions{ID: azidentity.ClientID("08ccc1fd-d017-4a1c-aa8d-3264f10d13b9")})
	if err != nil {
		fmt.Println(err)
		return
	}*/
	apiClient := armapimanagement.NewAPIClient("1ed07f18-4c8e-424d-bd80-19e1247574ad", credential, nil)
	pager := apiClient.ListByService("RG-APIGenerator", "hka-lab-api-management", nil)
	fmt.Printf("%#v\n", pager)
	fmt.Println(pager.Err())
	pager.NextPage(context.Background())
	fmt.Println(pager.Err())

}
