package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/inner-daydream/kvz/internal/kv"
)

type BaseKvCmd struct {
	s kv.KvService
}

type SetCmd struct {
	*BaseKvCmd
	Key       string `arg:""`
	Value     string `arg:""`
	ShowHooks bool   `help:"Show the stdout of the hooks that ran."`
	Verbose   bool   `help:"Display additional informations about the hooks that ran (status, name)"`
}

type GetCmd struct {
	*BaseKvCmd
	Key string `arg:""`
}

type AddHookCmd struct {
	*BaseKvCmd
	Name   string `arg:""`
	Script string `arg:""`
	File   bool   `help:"Takes in a file as a script. it will be saved in the store and changes to the file won't be picked up automatically"`
	Link   bool   `help:"Takes in a path to a a file, the path will be called on update"`
}

type AttachHookCmd struct {
	*BaseKvCmd
	Key      string `arg:""`
	HookName string `arg:""`
}

type LsKeysCmd struct {
	*BaseKvCmd
}

type RmKeyCmd struct {
	*BaseKvCmd
	Key string `arg:""`
}

type LsHooksCmd struct {
	*BaseKvCmd
}

type RmHookCmd struct {
	*BaseKvCmd
	Name string `arg:""`
}

type KvSubCmd struct {
	Set SetCmd    `cmd:"" help:"Set a kv pair"`
	Get GetCmd    `cmd:"" help:"Get the value of a key"`
	Ls  LsKeysCmd `cmd:"" help:"List the keys"`
	Rm  RmKeyCmd  `cmd:"" help:"Remove a key"`
}
type HookSubCmd struct {
	Add    AddHookCmd    `cmd:"" help:"create a hook, when attached to a key, the provided script will run whenever the value of the key is changed"`
	Attach AttachHookCmd `cmd:"" help:"attach a hook to a key, it will run whenever the value of the key is changed"`
	Ls     LsHooksCmd    `cmd:"" help:"List the hook names"`
	Rm     RmHookCmd     `cmd:"" help:"Remove a hook"`
}

type Cli struct {
	Kv   KvSubCmd   `cmd:""`
	Hook HookSubCmd `cmd:""`
}

func (c *SetCmd) Run() error {
	err := c.s.Set(c.Key, c.Value)
	if err != nil {
		return err
	}
	hooks, err := c.s.GetAttachedHooks(c.Key)
	if err != nil {
		return err
	}
	if len(hooks) == 0 {
		return nil
	}
	cmds, err := c.s.ExecHooks(hooks, c.Value)
	if err != nil {
		return err
	}
	if len(cmds) == 0 {
		return nil
	}

	for _, cmd := range cmds {
		if c.Verbose {
			if cmd.Error != nil {
				fmt.Printf("Hook %s execution failed: %s\n the following error occurred: %s\n", cmd.Caller, cmd.Error, cmd.Stderr)
			} else {
				fmt.Printf("Hook %s success.\nResult:\n%s", cmd.Caller, cmd.Stdout)
			}
		}
		if cmd.Error != nil {
			fmt.Printf("%s\n%s", cmd.Error, cmd.Stderr)
		}
		if c.ShowHooks {
			fmt.Printf(cmd.Stdout)
		}

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

func (c *AddHookCmd) Run() error {
	if !c.File && !c.Link {
		return c.s.AddScriptHook(c.Name, c.Script)
	}
	_, err := os.Stat(c.Script)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("the specified file does not exist")
	}
	filePath, err := filepath.Abs(c.Script)
	if err != nil {
		return fmt.Errorf("could not determine full path of the script")
	}
	if c.Link {
		return c.s.AddFilePathHook(c.Name, filePath)
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("could not read the script: %w", err)
	}
	return c.s.AddFileHook(c.Name, string(content))
}

func (c *AttachHookCmd) Run() error {
	return c.s.AttachHook(c.Key, c.HookName)
}

func (c *LsKeysCmd) Run() error {
	keys, err := c.s.ListKeys()
	if err != nil {
		return err
	}
	for _, key := range keys {
		fmt.Println(key)
	}
	return nil
}

func (c *LsHooksCmd) Run() error {
	hookNames, err := c.s.ListHooks()
	if err != nil {
		return err
	}
	for _, hookName := range hookNames {
		fmt.Println(hookName)
	}
	return nil
}

func (c *RmHookCmd) Run() error {
	return c.s.DeleteHook(c.Name)
}

func (c *RmKeyCmd) Run() error {
	return c.s.Delete(c.Key)
}

func NewCli(kvService kv.KvService) *Cli {
	baseKvCmd := BaseKvCmd{
		s: kvService,
	}
	return &Cli{
		Kv: KvSubCmd{
			Set: SetCmd{
				BaseKvCmd: &baseKvCmd,
			},
			Get: GetCmd{
				BaseKvCmd: &baseKvCmd,
			},
			Ls: LsKeysCmd{
				BaseKvCmd: &baseKvCmd,
			},
			Rm: RmKeyCmd{
				BaseKvCmd: &baseKvCmd,
			},
		},
		Hook: HookSubCmd{
			Add: AddHookCmd{
				BaseKvCmd: &baseKvCmd,
			},
			Attach: AttachHookCmd{
				BaseKvCmd: &baseKvCmd,
			},
			Ls: LsHooksCmd{
				BaseKvCmd: &baseKvCmd,
			},
			Rm: RmHookCmd{
				BaseKvCmd: &baseKvCmd,
			},
		},
	}
}

func ParseAndExecute(cli *Cli) {
	ctx := kong.Parse(cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
