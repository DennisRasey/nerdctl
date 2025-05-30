//go:build !linux

/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package container

import (
	"github.com/containerd/nerdctl/v2/pkg/inspecttypes/native"
	"github.com/containerd/nerdctl/v2/pkg/statsutil"
)

func setContainerStatsAndRenderStatsEntry(previousStats *statsutil.ContainerStats, firstSet bool, anydata interface{}, pid int, interfaces []native.NetInterface, systemInfo statsutil.SystemInfo) (statsutil.StatsEntry, error) {
	return statsutil.StatsEntry{}, nil
}

// getSystemCPUUsage reads the system's CPU usage from /proc/stat and returns
// the total CPU usage in nanoseconds and the number of CPUs.
func getSystemCPUUsage() (uint64, uint32, error) {
	return 0, 0, nil
}
