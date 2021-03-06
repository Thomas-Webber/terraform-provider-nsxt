/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause

   Generated by: https://github.com/swagger-api/swagger-codegen.git */

package administration

// SNMP Service properties
type SnmpServiceProperties struct {

	// SNMP v1, v2c community
	Communities []SnmpCommunity `json:"communities,omitempty"`

	// Start when system boots
	StartOnBoot bool `json:"start_on_boot"`
}
