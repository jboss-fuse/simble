/*
 * Copyright (C) 2018 Red Hat, Inc.
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
package http

import (
	"github.com/jboss-fuse/simble/v1/pkg/simble"
	"github.com/jboss-fuse/simble/v1/pkg/simble/echo"
	"github.com/jboss-fuse/simble/v1/pkg/simble/static"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {

	// Define the command options.
	addIntFlag("port", 3000, "The port the server accepts connections on.")
	addStringFlag("path", "", "The directory path to serve files from.")
	addStringFlag("prefix", "/", "The URL prefix to serve content from.")
	addBoolFlag("spa", false, "Run in Single Page App mode.")
	addBoolFlag("etags", false, "Calculate ETag headers to assist caching.")
	addIntFlag("tls-port", 3443, "The port the server accepts TLS connections on.")
	addStringFlag("tls-lets-encrypt", "", "The directory used to cache letsencrypt.org certificates. Enables TLS.")
	addStringFlag("tls-key", "", "The TLS private key file. Enables TLS.")
	addStringFlag("tls-cert", "", "The TLS certificate file. Enables TLS.")

	// Support using env vars to configure command options.
	viper.SetEnvPrefix(Command.Use)
	viper.AutomaticEnv()
}
var (
	Command = &cobra.Command{
		Use: "http",
		Short: "Starts an http server for serving static content",
		RunE: func(cmd *cobra.Command, args []string) error {
			server := simble.New()
			server.AddContext(&echo.EchoContext{
				Port:              viper.GetInt("port"),
				TLSPort:           viper.GetInt("tls-port"),
				TLSLetsEncryptDir: viper.GetString("tls-lets-encrypt"),
				TLSKeyFile:        viper.GetString("tls-key"),
				TLSCertFile:       viper.GetString("tls-cert"),
			})
			server.AddContext(&static.StaticContext{
				URLPath: viper.GetString("prefix"),
				DirPath: viper.GetString("path"),
				SinglePageAppMode: viper.GetBool("spa"),
				ETags: viper.GetBool("etags"),
			})
			return server.Run();
		},
	}
)


func addStringFlag(flagName string, defaultValue string, description string) {
	Command.Flags().String(flagName, defaultValue, description)
	viper.BindPFlag(flagName, Command.Flags().Lookup(flagName))
}
func addIntFlag(flagName string, defaultValue int, description string) {
	Command.Flags().Int(flagName, defaultValue, description)
	viper.BindPFlag(flagName, Command.Flags().Lookup(flagName))
}
func addBoolFlag(flagName string, defaultValue bool, description string) {
	Command.Flags().Bool(flagName, defaultValue, description)
	viper.BindPFlag(flagName, Command.Flags().Lookup(flagName))
}