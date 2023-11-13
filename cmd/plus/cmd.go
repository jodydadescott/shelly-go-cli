package plus

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/jodydadescott/shelly-go-sdk/plus"
	"github.com/jodydadescott/shelly-go-sdk/plus/input"
	"github.com/jodydadescott/shelly-go-sdk/plus/light"
	"github.com/jodydadescott/shelly-go-sdk/plus/shelly"
	"github.com/jodydadescott/shelly-go-sdk/plus/switchx"
	"github.com/jodydadescott/shelly-go-sdk/plus/system"
	"github.com/jodydadescott/shelly-go-sdk/plus/wifi"

	lightcmd "github.com/jodydadescott/shelly-go-cli/cmd/plus/light"
	shellycmd "github.com/jodydadescott/shelly-go-cli/cmd/plus/shelly"
	switchxcmd "github.com/jodydadescott/shelly-go-cli/cmd/plus/switchx"
	wificmd "github.com/jodydadescott/shelly-go-cli/cmd/plus/wifi"
	"github.com/jodydadescott/shelly-go-cli/types"
)

type callback interface {
	PlusClient() (*plus.Client, error)
	WriteStdout(any) error
	WriteStderr(string)
	GetFiles() (*types.Files, error)
}

type Cmd struct {
	_client *plus.Client
	*cobra.Command
	callback
}

func NewCmd(callback callback) *cobra.Command {

	t := &Cmd{}

	t.Command = &cobra.Command{
		Use:   "plus",
		Short: "Shelly Plus",
	}

	t.callback = callback

	t.AddCommand(shellycmd.NewCmd(t), wificmd.NewCmd(t), switchxcmd.NewCmd(t), lightcmd.NewCmd(t))
	return t.Command
}

func (t *Cmd) LogDebug(s string) {
	t.WriteStderr(s)
}

func (t *Cmd) System() (*system.Client, error) {
	client, err := t.callback.PlusClient()
	if err != nil {
		return nil, err
	}
	return client.System(), nil
}

func (t *Cmd) Shelly() (*shelly.Client, error) {
	client, err := t.callback.PlusClient()
	if err != nil {
		return nil, err
	}
	return client.Shelly(), nil
}

func (t *Cmd) Wifi() (*wifi.Client, error) {
	client, err := t.callback.PlusClient()
	if err != nil {
		return nil, err
	}
	return client.Wifi(), nil
}

func (t *Cmd) Switch() (*switchx.Client, error) {
	client, err := t.callback.PlusClient()
	if err != nil {
		return nil, err
	}
	return client.Switch(), nil
}

func (t *Cmd) Light() (*light.Client, error) {
	client, err := t.callback.PlusClient()
	if err != nil {
		return nil, err
	}
	return client.Light(), nil
}

func (t *Cmd) Input() (*input.Client, error) {
	client, err := t.callback.PlusClient()
	if err != nil {
		return nil, err
	}
	return client.Input(), nil
}

func (t *Cmd) RebootDevice(ctx context.Context) error {
	shelly, err := t.Shelly()
	if err != nil {
		return err
	}
	return shelly.Reboot(ctx)
}
