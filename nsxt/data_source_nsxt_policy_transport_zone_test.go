/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccDataSourceNsxtPolicyTransportZone_basic(t *testing.T) {
	transportZoneName := getVlanTransportZoneName()
	testResourceName := "data.nsxt_policy_transport_zone.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccNSXPolicyTransportZonePrecheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNSXPolicyTransportZoneReadTemplate(transportZoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", transportZoneName),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "is_default"),
					resource.TestCheckResourceAttr(testResourceName, "transport_type", "VLAN_BACKED"),
				),
			},
			{
				Config: testAccNSXPolicyTransportZoneWithTransportTypeTemplate(transportZoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "display_name", transportZoneName),
					resource.TestCheckResourceAttrSet(testResourceName, "id"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "is_default"),
					resource.TestCheckResourceAttr(testResourceName, "transport_type", "VLAN_BACKED"),
				),
			},
		},
	})
}

func testAccNSXPolicyTransportZoneReadTemplate(transportZoneName string) string {
	if testAccIsGlobalManager() {
		sitePath := getTestSitePath()
		return testAccNSXGlobalPolicyTransportZoneReadTemplate(transportZoneName, sitePath)
	}
	return fmt.Sprintf(`
data "nsxt_policy_transport_zone" "test" {
  display_name = "%s"
}`, transportZoneName)
}

func testAccNSXPolicyTransportZoneWithTransportTypeTemplate(transportZoneName string) string {
	if testAccIsGlobalManager() {
		sitePath := getTestSitePath()
		return testAccNSXGlobalPolicyTransportZoneWithTransportTypeTemplate(transportZoneName, sitePath)
	}
	return fmt.Sprintf(`
data "nsxt_policy_transport_zone" "test" {
  display_name   = "%s"
  transport_type = "VLAN_BACKED"
}`, transportZoneName)
}

func testAccNSXGlobalPolicyTransportZoneReadTemplate(transportZoneName string, sitePath string) string {
	return fmt.Sprintf(`
data "nsxt_policy_transport_zone" "test" {
  display_name = "%s"
  site_path = "%s"
}`, transportZoneName, sitePath)
}

func testAccNSXPolicyTransportZonePrecheck(t *testing.T) {
	testAccPreCheck(t)
	if testAccIsGlobalManager() && getTestSitePath() == "" {
		str := fmt.Sprintf("%s must be set for this acceptance test", "NSXT_TEST_SITE_PATH")
		t.Fatal(str)
	}
}

func testAccNSXGlobalPolicyTransportZoneWithTransportTypeTemplate(transportZoneName string, sitePath string) string {
	return fmt.Sprintf(`
data "nsxt_policy_transport_zone" "test" {
  display_name = "%s"
  site_path = "%s"
  transport_type = "VLAN_BACKED"
}`, transportZoneName, sitePath)
}
