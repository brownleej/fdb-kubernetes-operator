/*
 * replace_failed_process_groups_test.go
 *
 * This source file is part of the FoundationDB open source project
 *
 * Copyright 2024 Apple Inc. and the FoundationDB project authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package replacements

import (
	fdbv1beta2 "github.com/FoundationDB/fdb-kubernetes-operator/api/v1beta2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/utils/pointer"
)

var _ = Describe("replace_failed_process_groups", func() {
	DescribeTable("check if removal is allowed", func(cluster *fdbv1beta2.FoundationDBCluster, maxReplacements int, faultDomainsWithReplacements map[fdbv1beta2.FaultDomain]fdbv1beta2.None, faultDomain fdbv1beta2.FaultDomain, expected bool) {
		Expect(removalAllowed(cluster, maxReplacements, faultDomainsWithReplacements, faultDomain)).To(Equal(expected))
	},
		Entry("process group based replacement: with 1 replacement allowed",
			&fdbv1beta2.FoundationDBCluster{},
			1,
			nil,
			fdbv1beta2.FaultDomain(""),
			true,
		),
		Entry("process group based replacement: with 0 replacements allowed",
			&fdbv1beta2.FoundationDBCluster{},
			0,
			nil,
			fdbv1beta2.FaultDomain(""),
			false,
		),
		Entry("fault domain based replacement: no ongoing replacements",
			&fdbv1beta2.FoundationDBCluster{
				Spec: fdbv1beta2.FoundationDBClusterSpec{
					AutomationOptions: fdbv1beta2.FoundationDBClusterAutomationOptions{
						Replacements: fdbv1beta2.AutomaticReplacementOptions{
							FaultDomainBasedReplacements: pointer.Bool(true),
						},
					},
				},
			},
			0,
			nil,
			fdbv1beta2.FaultDomain("zone1"),
			true,
		),
		Entry("fault domain based replacement: ongoing replacements same fault domain",
			&fdbv1beta2.FoundationDBCluster{
				Spec: fdbv1beta2.FoundationDBClusterSpec{
					AutomationOptions: fdbv1beta2.FoundationDBClusterAutomationOptions{
						Replacements: fdbv1beta2.AutomaticReplacementOptions{
							FaultDomainBasedReplacements: pointer.Bool(true),
						},
					},
				},
			},
			0,
			map[fdbv1beta2.FaultDomain]fdbv1beta2.None{
				"zone1": {},
			},
			fdbv1beta2.FaultDomain("zone1"),
			true,
		),
		Entry("fault domain based replacement: ongoing replacements different fault domain",
			&fdbv1beta2.FoundationDBCluster{
				Spec: fdbv1beta2.FoundationDBClusterSpec{
					AutomationOptions: fdbv1beta2.FoundationDBClusterAutomationOptions{
						Replacements: fdbv1beta2.AutomaticReplacementOptions{
							FaultDomainBasedReplacements: pointer.Bool(true),
						},
					},
				},
			},
			0,
			map[fdbv1beta2.FaultDomain]fdbv1beta2.None{
				"zone1": {},
			},
			fdbv1beta2.FaultDomain("zone2"),
			false,
		),
		Entry("fault domain based replacement: too many ongoing replacements same fault domain",
			&fdbv1beta2.FoundationDBCluster{
				Spec: fdbv1beta2.FoundationDBClusterSpec{
					AutomationOptions: fdbv1beta2.FoundationDBClusterAutomationOptions{
						Replacements: fdbv1beta2.AutomaticReplacementOptions{
							FaultDomainBasedReplacements: pointer.Bool(true),
						},
					},
				},
			},
			0,
			map[fdbv1beta2.FaultDomain]fdbv1beta2.None{
				"zone1": {},
				"zone2": {},
			},
			fdbv1beta2.FaultDomain("zone1"),
			false,
		),
	)

})
