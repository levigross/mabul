// Copyright © 2017 Levi Gross <levi@levigross.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"
	"time"

	"github.com/levigross/mabul/udp"
	"github.com/spf13/cobra"
)

var udpAttackConfig udp.AttackConfig

// udpCmd represents the udp command
var udpCmd = &cobra.Command{
	Use:   "udp",
	Short: "Launches UDP style attacks (stateless)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		udpAttacker := &udp.Attacker{
			Config: &udpAttackConfig,
			Target: target,
			Log:    setupLogging(),
		}
		if err := udpAttacker.Attack(&target); err != nil {
			udpAttacker.Log.Info("Unable to execute attack: ", err)
			os.Exit(-1)
		}
	},
}

func init() {
	RootCmd.AddCommand(udpCmd)
	addTargetFlags(udpCmd)

	udpCmd.Flags().DurationVar(&udpAttackConfig.AttackDuration, "attackDuration", time.Second*10, "Attack time duration")
	udpCmd.Flags().UintVar(&udpAttackConfig.NumThreads, "numThreads", 10, "Number of threads")
}
