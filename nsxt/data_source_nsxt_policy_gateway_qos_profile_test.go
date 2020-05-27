/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	gm_infra "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra"
	gm_model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
	"testing"
)

func TestAccDataSourceNsxtPolicyGatewayQosProfile_basic(t *testing.T) {
	name := "terraform_test"
	testResourceName := "data.nsxt_policy_gateway_qos_profile.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "3.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccDataSourceNsxtPolicyGatewayQosProfileDeleteByName(name)
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if err := testAccDataSourceNsxtPolicyGatewayQosProfileCreate(name); err != nil {
						panic(err)
					}
				},
				Config: testAccNsxtPolicyGatewayQosProfileReadTemplate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", name),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
				),
			},
			{
				Config: testAccNsxtPolicyEmptyTemplate(),
			},
		},
	})
}

func testAccDataSourceNsxtPolicyGatewayQosProfileCreate(name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}

	displayName := name
	description := name
	obj := model.GatewayQosProfile{
		Description: &description,
		DisplayName: &displayName,
	}

	// Generate a random ID for the resource
	id := newUUID()

	converter := bindings.NewTypeConverter()
	converter.SetMode(bindings.REST)
	if testAccIsGlobalManager() {
		dataValue, err1 := converter.ConvertToVapi(obj, model.GatewayQosProfileBindingType())
		if err1 != nil {
			return err1[0]
		}
		gmObj, err2 := converter.ConvertToGolang(dataValue, gm_model.GatewayQosProfileBindingType())
		if err2 != nil {
			return err2[0]
		}
		gmProfile := gmObj.(gm_model.GatewayQosProfile)
		client := gm_infra.NewDefaultGatewayQosProfilesClient(connector)
		err = client.Patch(id, gmProfile)
	} else {
		client := infra.NewDefaultGatewayQosProfilesClient(connector)
		err = client.Patch(id, obj)
	}

	if err != nil {
		return handleCreateError("GatewayQosProfile", id, err)
	}
	return nil
}

func testAccDataSourceNsxtPolicyGatewayQosProfileDeleteByName(name string) error {
	connector, err := testAccGetPolicyConnector()
	if err != nil {
		return fmt.Errorf("Error during test client initialization: %v", err)
	}

	// Find the object by name and delete it
	if testAccIsGlobalManager() {
		objID, err := testGetObjIDByName(name, "GatewayQosProfile")
		if err == nil {
			client := gm_infra.NewDefaultQosProfilesClient(connector)
			err := client.Delete(objID)
			if err != nil {
				return handleDeleteError("GatewayQosProfile", objID, err)
			}
			return nil
		}
	} else {
		client := infra.NewDefaultGatewayQosProfilesClient(connector)
		// Find the object by name
		objList, err := client.List(nil, nil, nil, nil, nil, nil)
		if err != nil {
			return handleListError("GatewayQosProfile", err)
		}
		for _, objInList := range objList.Results {
			if *objInList.DisplayName == name {
				err := client.Delete(*objInList.Id)
				if err != nil {
					return handleDeleteError("GatewayQosProfile", *objInList.Id, err)
				}
				return nil
			}
		}
	}
	return fmt.Errorf("Error while deleting GatewayQosProfile '%s': resource not found", name)
}

func testAccNsxtPolicyGatewayQosProfileReadTemplate(name string) string {
	return fmt.Sprintf(`
data "nsxt_policy_gateway_qos_profile" "test" {
  display_name = "%s"
}`, name)
}
