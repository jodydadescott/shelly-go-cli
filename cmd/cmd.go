package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hokaccha/go-prettyjson"
	shelly "github.com/jodydadescott/shelly-go-sdk"

	"github.com/jodydadescott/shelly-go-sdk/plus"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	pluscmd "github.com/jodydadescott/shelly-go-cli/cmd/plus"
	"github.com/jodydadescott/shelly-go-cli/logging"
	"github.com/jodydadescott/shelly-go-cli/types"
)

type Cmd struct {
	initalized bool
	*cobra.Command
	_client         *shelly.Client
	_plusClient     *plus.Client
	hostnameArg     string
	passwordArg     string
	outputArg       string
	filenameArg     string
	debugEnabledArg bool
}

func NewCmd() *Cmd {

	t := &Cmd{}

	command := &cobra.Command{

		Use: BinaryName,

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if t.debugEnabledArg {
				zap.ReplaceGlobals(logging.GetDebugZapLogger())
				zap.L().Debug("debug is enabled")
			} else {
				zap.ReplaceGlobals(logging.GetDefaultZapLogger())
			}
		},

		SilenceUsage: true,

		PersistentPostRun: func(cmd *cobra.Command, args []string) {

			if t._client != nil {
				t._client.Close()
			}

		},
	}

	t.Command = command

	t.PersistentFlags().StringVarP(&t.hostnameArg, "hostname", "H", "", fmt.Sprintf("Hostname; optionally use env var '%s'", ShellyHostnameEnvVar))
	t.PersistentFlags().StringVarP(&t.hostnameArg, "password", "p", "", fmt.Sprintf("Password; optionally use env var '%s'", ShellyPasswordEnvVar))
	t.PersistentFlags().StringVarP(&t.outputArg, "output", "o", ShellyOutputDefault, fmt.Sprintf("Output format. One of: prettyjson | json | jsonpath | yaml ; Optionally use env var '%s'", ShellyOutputEnvVar))
	t.PersistentFlags().StringVarP(&t.filenameArg, "filename", "f", "", "Filename or Dirname")
	t.PersistentFlags().BoolVarP(&t.debugEnabledArg, "debug", "d", false, "debug to STDERR")
	t.AddCommand(pluscmd.NewCmd(t))

	return t
}

func (t *Cmd) client() *shelly.Client {

	if t._client != nil {
		return t._client
	}

	config := &shelly.Config{
		DebugEnabled: t.debugEnabledArg,
		Hostname:     t.hostnameArg,
		Password:     t.passwordArg,
	}

	if config.Hostname == "" {
		config.Hostname = os.Getenv(ShellyHostnameEnvVar)
	}

	if config.Password == "" {
		config.Password = os.Getenv(ShellyPasswordEnvVar)
	}

	t._client = shelly.New(config)
	return t._client
}

func (t *Cmd) PlusClient() (*plus.Client, error) {

	if t._plusClient != nil {
		return t._plusClient, nil
	}

	return t.client().PlusClient()
}

// WriteObject writes object in desired format to STDOUT
func (t *Cmd) WriteStdout(input any) error {

	switch input.(type) {

	case nil:
		return nil

	case string:
		fmt.Println(input.(string))
		return nil

	case *string:
		fmt.Println(*input.(*string))
		return nil

	}

	// Arg can be in the format of 'format' or 'format=...'.
	// For example jsonpath expects a second arg such as jsonpath=$.

	outputArgSplit := strings.Split(t.outputArg, "=")

	switch strings.ToLower(outputArgSplit[0]) {

	case "prettyjson":
		data, err := prettyjson.Marshal(input)
		if err != nil {
			return err
		}
		fmt.Println(strings.TrimSpace(string(data)))
		return nil

	case "jsonpath":

		if len(outputArgSplit) > 1 {

			v := interface{}(nil)
			data, err := json.Marshal(input)
			if err != nil {
				return err
			}

			err = json.Unmarshal(data, &v)
			if err != nil {
				return err
			}

			data2, err := jsonpath.Get(outputArgSplit[1], v)
			if err != nil {
				return err
			}

			switch data2.(type) {
			case string:
				fmt.Println(data2)

			case []interface{}:
				s := reflect.ValueOf(data2)
				for i := 0; i < s.Len(); i++ {
					fmt.Println(s.Index(i))
				}
			}

			return nil
		}

		return fmt.Errorf("Missing jsonpath arg. Expect jsonpath=...")

	case "json":
		data, err := json.Marshal(input)
		if err != nil {
			return err
		}
		fmt.Println(strings.TrimSpace(string(data)))
		return nil

	case "yaml":
		data, err := yaml.Marshal(input)
		if err != nil {
			return err
		}
		fmt.Println(strings.TrimSpace(string(data)))
		return nil

	}

	return fmt.Errorf("format type %s is unknown", t.outputArg)
}

func (t *Cmd) WriteStderr(s string) {
	fmt.Fprintln(os.Stderr, s)
}

// // ReadInput reads input from file if specified. If no file is specified STDIN will be checked and
// // if present will be returned. If filename is not specified and no STDIN data is present then an error
// // will be returned.
// func (t *Cmd) ReadInput() ([]byte, error) {

// 	if t.filenameArg != "" {
// 		return os.ReadFile(t.filenameArg)
// 	}

// 	fi, err := os.Stdin.Stat()
// 	if err != nil {
// 		return nil, err
// 	}

// 	if (fi.Mode() & os.ModeCharDevice) == 0 {
// 		return io.ReadAll(os.Stdin)
// 	}

// 	return nil, fmt.Errorf("Data input is required. Use filename or pipe to STDIN")
// }

func (t *Cmd) GetFiles() (*types.Files, error) {
	return types.NewFiles(t.filenameArg)
}
