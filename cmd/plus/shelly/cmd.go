package shelly

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/jodydadescott/shelly-go-sdk/plus/shelly"
)

type callback interface {
	WriteStdout(any) error
	WriteStderr(s string)
	ReadInput() ([]byte, error)
	Shelly() (*shelly.Client, error)
	RebootDevice(ctx context.Context) error
}

func NewCmd(callback callback) *cobra.Command {

	var stageArg string
	var urlArg string
	var appendArg string
	var autorebootArg bool

	var ignoreAuth bool
	var ignoreNetArg bool
	var ignoreBluetoothArg bool
	var ignoreCloudArg bool
	var ignoreInputArg bool
	var ignoreLightArg bool
	var ignoreMqttArg bool
	var ignoreSwitchArg bool
	var ignoreSystemArg bool
	var ignoreWebsocketArg bool
	var markupArg bool

	rootCmd := &cobra.Command{
		Use:   "shelly",
		Short: "Shelly Component",
	}

	getConfigCmd := &cobra.Command{
		Use:   "get-config",
		Short: "Returns config",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := callback.Shelly()
			if err != nil {
				return err
			}

			config, err := client.GetConfig(cmd.Context(), markupArg)
			if err != nil {
				return err
			}

			return callback.WriteStdout(config)
		},
	}

	getConfigCmd.PersistentFlags().BoolVar(&markupArg, "markup", false, "returns config that can be used as a template")

	getStatusCmd := &cobra.Command{
		Use:   "get-status",
		Short: "Returns status",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := callback.Shelly()
			if err != nil {
				return err
			}

			result, err := client.GetStatus(cmd.Context())
			if err != nil {
				return err
			}

			return callback.WriteStdout(result)
		},
	}

	getInfoCmd := &cobra.Command{
		Use:   "get-info",
		Short: "Returns device info",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := callback.Shelly()
			if err != nil {
				return err
			}

			result, err := client.GetDeviceInfo(cmd.Context())
			if err != nil {
				return err
			}

			return callback.WriteStdout(result)
		},
	}

	getMethodsCmd := &cobra.Command{
		Use:   "get-methods",
		Short: "Returns all available RPC methods for device",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := callback.Shelly()
			if err != nil {
				return err
			}

			result, err := client.ListMethods(cmd.Context())
			if err != nil {
				return err
			}

			return callback.WriteStdout(result)
		},
	}

	getUpdatesCmd := &cobra.Command{
		Use:   "get-updates",
		Short: "Returns available update info",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := callback.Shelly()
			if err != nil {
				return err
			}

			result, err := client.CheckForUpdate(cmd.Context())
			if err != nil {
				return err
			}

			return callback.WriteStdout(result)
		},
	}

	rebootCmd := &cobra.Command{
		Use:   "reboot",
		Short: "Executes device reboot",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := callback.Shelly()
			if err != nil {
				return err
			}

			err = client.Reboot(cmd.Context())
			if err != nil {
				return err
			}

			return nil
		},
	}

	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Returns available update info",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := callback.Shelly()
			if err != nil {
				return err
			}

			params := &shelly.ShellyUpdateConfig{}

			if stageArg != "" {
				params.Stage = &stageArg
			}

			if urlArg != "" {
				params.Url = &urlArg
			}

			return client.Update(cmd.Context(), params)
		},
	}

	updateCmd.PersistentFlags().StringVar(&stageArg, "stage", "", "The type of the new version - either stable or beta. By default updates to stable version. Optional")
	updateCmd.PersistentFlags().StringVar(&urlArg, "url", "", "Url address of the update. Optional")

	factoryResetCmd := &cobra.Command{
		Use:   "factory-reset",
		Short: "Executes factory reset",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := callback.Shelly()
			if err != nil {
				return err
			}

			return client.FactoryReset(cmd.Context())
		},
	}

	resetWifiConfigCmd := &cobra.Command{
		Use:   "reset-wifi-config",
		Short: "Executes Wifi config reset",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := callback.Shelly()
			if err != nil {
				return err
			}

			return client.ResetWiFiConfig(cmd.Context())
		},
	}

	setConfigCmd := &cobra.Command{
		Use:   "set-config",
		Short: "Sets config",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := callback.Shelly()
			if err != nil {
				return err
			}

			b, err := callback.ReadInput()
			if err != nil {
				return err
			}

			var config *shelly.ShellyConfig

			var errors *multierror.Error

			err = json.Unmarshal(b, &config)
			if err != nil {
				errors = multierror.Append(errors, err)
				err = yaml.Unmarshal(b, &config)

				if err != nil {
					errors = multierror.Append(errors, err)
					errors = multierror.Append(errors, fmt.Errorf("Invalid format. Expect JSON or YAML"))
					return errors.ErrorOrNil()
				}
			}

			if ignoreAuth {
				config.Auth = nil
			}

			if ignoreNetArg {
				config.Ethernet = nil
				config.Wifi = nil
			}

			if ignoreBluetoothArg {
				config.Bluetooth = nil
			}

			if ignoreCloudArg {
				config.Cloud = nil
			}

			if ignoreInputArg {
				config.Input = nil
			}

			if ignoreLightArg {
				config.Light = nil
			}

			if ignoreMqttArg {
				config.Mqtt = nil
			}

			if ignoreSwitchArg {
				config.Switch = nil
			}

			if ignoreSystemArg {
				config.System = nil
			}

			if ignoreWebsocketArg {
				config.Websocket = nil
			}

			report := client.SetConfig(cmd.Context(), config)

			err = report.Error()

			if err != nil {
				return err
			}

			if report.RebootRequired() {
				if autorebootArg {
					callback.WriteStderr("reboot is required; rebooting ...")
					return callback.RebootDevice(cmd.Context())

				}

				callback.WriteStderr("reboot is required")
			}

			return callback.WriteStdout(report)
		},
	}

	setConfigCmd.PersistentFlags().BoolVar(&autorebootArg, "autoreboot", false, "automatically reboot device is necessary")

	setConfigCmd.PersistentFlags().BoolVar(&ignoreAuth, "ignore-auth", false, "ignore Auth config")
	setConfigCmd.PersistentFlags().BoolVar(&ignoreNetArg, "ignore-net", false, "ignore net (Ethernet/Wifi) config")

	setConfigCmd.PersistentFlags().BoolVar(&ignoreBluetoothArg, "ignore-bluetooth", false, "ignore Bluetooth config")
	setConfigCmd.PersistentFlags().BoolVar(&ignoreCloudArg, "ignore-cloud", false, "ignore Cloud config")
	setConfigCmd.PersistentFlags().BoolVar(&ignoreInputArg, "ignore-input", false, "ignore Input config")
	setConfigCmd.PersistentFlags().BoolVar(&ignoreLightArg, "ignore-light", false, "ignore Light config")
	setConfigCmd.PersistentFlags().BoolVar(&ignoreMqttArg, "ignore-mqtt", false, "ignore Mqtt config")
	setConfigCmd.PersistentFlags().BoolVar(&ignoreSwitchArg, "ignore-switch", false, "ignore Switch config")
	setConfigCmd.PersistentFlags().BoolVar(&ignoreSystemArg, "ignore-system", false, "ignore System config")
	setConfigCmd.PersistentFlags().BoolVar(&ignoreWebsocketArg, "ignore-websocket", false, "ignore Websocket config")

	getAppend := func() (bool, error) {

		switch strings.ToLower(appendArg) {

		case "true":
			return true, nil

		case "false":
			return false, nil

		default:
			return false, fmt.Errorf("append must be set to true or false")

		}

	}

	putTlsClientCertCmd := &cobra.Command{
		Use:   "put-tls-client-cert",
		Short: "Sets TLS Client Cert",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := callback.Shelly()
			if err != nil {
				return err
			}

			append, err := getAppend()
			if err != nil {
				return err
			}

			b, err := callback.ReadInput()
			if err != nil {
				return err
			}

			data := string(b)

			return client.PutTLSClientCert(cmd.Context(), &shelly.ShellyTLSConfig{
				Data:   &data,
				Append: &append,
			})
		},
	}

	appendMsg := "true if more data will be appended afterwards, default false"

	putTlsClientCertCmd.PersistentFlags().StringVar(&appendArg, "append", "", appendMsg)

	putTlsClientKeyCmd := &cobra.Command{
		Use:   "put-tls-client-key",
		Short: "Sets TLS Client Key",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := callback.Shelly()
			if err != nil {
				return err
			}

			append, err := getAppend()
			if err != nil {
				return err
			}

			b, err := callback.ReadInput()
			if err != nil {
				return err
			}

			data := string(b)

			return client.PutTLSClientKey(cmd.Context(), &shelly.ShellyTLSConfig{
				Data:   &data,
				Append: &append,
			})
		},
	}

	putTlsClientKeyCmd.PersistentFlags().StringVar(&appendArg, "append", "", appendMsg)

	putUserCACmd := &cobra.Command{
		Use:   "put-user-ca",
		Short: "Sets Users CA",
		RunE: func(cmd *cobra.Command, args []string) error {

			client, err := callback.Shelly()
			if err != nil {
				return err
			}

			append, err := getAppend()
			if err != nil {
				return err
			}

			b, err := callback.ReadInput()
			if err != nil {
				return err
			}

			data := string(b)

			return client.PutUserCA(cmd.Context(), &shelly.ShellyTLSConfig{
				Data:   &data,
				Append: &append,
			})
		},
	}

	putUserCACmd.PersistentFlags().StringVar(&appendArg, "append", "", appendMsg)

	rootCmd.AddCommand(getConfigCmd, getStatusCmd, getInfoCmd, getMethodsCmd,
		getUpdatesCmd, rebootCmd, updateCmd,
		factoryResetCmd, resetWifiConfigCmd, setConfigCmd,
		putTlsClientCertCmd, putTlsClientKeyCmd, putUserCACmd)
	return rootCmd
}
