/*
Copyright © 2023 Eliott Teissonniere

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package cmd

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/eteissonniere/app-store-connect-cli/client"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	argKeyId      string
	argIssuerId   string
	argKeyPath    string
	argBundleId   string
	argUseSandbox bool

	apiClient *client.Client
)

var rootCmd = &cobra.Command{
	Use:   "app-store-connect-cli",
	Short: "A simple CLI tool to interact with the App Store Connect API",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		errMissingRequiredArgument := fmt.Errorf("missing required arguments: key-id, issuer-id, private-key, and bundle-id are required")
		if argKeyId == "" {
			return errMissingRequiredArgument
		}
		if argIssuerId == "" {
			return errMissingRequiredArgument
		}
		if argKeyPath == "" {
			return errMissingRequiredArgument
		}
		if argBundleId == "" {
			return errMissingRequiredArgument
		}

		p8Bytes, err := ioutil.ReadFile(argKeyPath)
		if err != nil {
			return fmt.Errorf("unable to read private key file: %w", err)
		}

		block, _ := pem.Decode(p8Bytes)
		if block == nil || block.Type != "PRIVATE KEY" {
			return fmt.Errorf("failed to decode PEM block containing private key")
		}

		parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse private key: %w", err)
		}

		ecdsaPrivateKey, ok := parsedKey.(*ecdsa.PrivateKey)
		if !ok {
			return fmt.Errorf("private key is not an ECDSA private key")
		}
		apiClient = client.New(client.ClientConfig{
			IssuerId:   argIssuerId,
			BundleId:   argBundleId,
			KeyId:      argKeyId,
			PrivateKey: ecdsaPrivateKey,
			UseSandbox: argUseSandbox,
		})

		log.Debug().
			Str("key-id", argKeyId).
			Str("issuer-id", argIssuerId).
			Str("bundle-id", argBundleId).
			Str("private-key", argKeyPath).
			Bool("use-sandbox", argUseSandbox).
			Msg("API client initialized")

		return nil
	},
}

func Execute() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&argKeyId, "key-id", "", "The key ID of the API key")
	rootCmd.PersistentFlags().StringVar(&argIssuerId, "issuer-id", "", "The issuer ID of the API key")
	rootCmd.PersistentFlags().StringVar(&argKeyPath, "private-key", "", "The path to the private key file")
	rootCmd.PersistentFlags().StringVar(&argBundleId, "bundle-id", "", "The bundle ID of the application")
	rootCmd.PersistentFlags().BoolVar(&argUseSandbox, "use-sandbox", false, "Use the sandbox environment")
}
