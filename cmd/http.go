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

	"github.com/levigross/mabul/http"
	"github.com/spf13/cobra"
)

var httpAttack http.AttackConfig

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "This is designed to execute layer 7 attacks",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		attacker := &http.Attacker{
			Config: &httpAttack,
			Log:    setupLogging(),
		}
		if err := attacker.Attack(&httpAttack); err != nil {
			attacker.Log.Info("Error executing attack: ", err)
			os.Exit(-1)
		}
	},
}

func init() {
	RootCmd.AddCommand(httpCmd)

	httpCmd.Flags().DurationVar(&httpAttack.AttackDuration, "attackDuration", time.Second*10, "Attack time duration")
	httpCmd.Flags().StringVar(&httpAttack.HTTPClient, "httpClient", "fasthttp", "The HTTP client you wish to use (fasthttp or net/http)")
	httpCmd.Flags().UintVar(&httpAttack.NumThreads, "numThreads", 10, "Number of threads")
	httpCmd.Flags().UintVar(&httpAttack.NumConnections, "numConnections", 100, "Number of connections per thread")
	httpCmd.Flags().DurationVar(&httpAttack.Timeout, "requestTimeout", time.Second, "Request timeout per request")
	httpCmd.Flags().IntVar(&httpAttack.ErrorThreshold, "errorThreshold", 50,
		"The precentage of errors you are willing to enjoy - use -1 for unlimited")
	httpCmd.Flags().StringVar(&httpAttack.URL, "url", "http://localhost:8000/", "The URL you wish to attack")
}
