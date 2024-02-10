package cli

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/inner-daydream/kvz/internal/kv"
)

type KvCmd struct {
	s kv.KvService
}

type SetCmd struct {
	KvCmd
	Key   string `arg:""`
	Value string `arg:""`
}

type GetCmd struct {
	KvCmd
	Key string `arg:""`
}

type Cli struct {
	Set SetCmd `cmd:"" help:"Set a kv pair"`
	Get GetCmd `cmd:"" help:"Get the value of a key"`
}

func (c *SetCmd) Run() error {
	err := c.s.Set(c.Key, c.Value)
	if err != nil {
		return err
	}
	return nil
}

func (c *GetCmd) Run() error {
	val, err := c.s.Get(c.Key)
	if err != nil {
		return err
	}
	fmt.Println(val)
	return nil
}

func NewCli(kvService *kv.KvService) *Cli {
	kvCmd := KvCmd{
		s: *kvService,
	}
	return &Cli{
		Set: SetCmd{
			KvCmd: kvCmd,
		},
		Get: GetCmd{
			KvCmd: kvCmd,
		},
	}
}

func ParseAndExecute(cli *Cli) {
	ctx := kong.Parse(cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
