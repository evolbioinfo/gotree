package shell

import (
	"fmt"
	"os"

	"github.com/abiosoft/ishell"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Shell interface {
	AddCmd(*ishell.Cmd)
	Run()
	Println(...interface{})
}

func AddCommands(s Shell, rootcmd *cobra.Command, parent *ishell.Cmd, cmds ...*cobra.Command) {
	for _, cmd := range cmds {
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			// In console mode, we deactivate stdin input files
			if f.Value.String() == "stdin" || f.Value.String() == "-" ||
				f.DefValue == "stdin" || f.DefValue == "-" {
				f.Value.Set("none")
				f.DefValue = "none"
			}
		})
		if cmd.Name() != "help" {
			ishellcmd := &ishell.Cmd{
				Name:     cmd.Name(),
				Help:     cmd.Short,
				LongHelp: cmd.Long,
				Func: func(c *ishell.Context) {
					// We reinitialize all flags
					cobrac, _, err := rootcmd.Find(c.RawArgs)
					if err != nil {
						fmt.Fprintf(os.Stderr, "%s\n", err.Error())
					} else {
						cobrac.Flags().VisitAll(func(f *pflag.Flag) {
							// In console mode, we deactivate stdin input files
							if f.Value.String() == "stdin" || f.Value.String() == "-" ||
								f.DefValue == "stdin" || f.DefValue == "-" {
								f.Value.Set("none")
							} else {
								f.Value.Set(f.DefValue)
							}
						})
						// Then we execute the command using cobra
						rootcmd.SetArgs(c.RawArgs)
						rootcmd.Execute()
					}
				},
			}
			if parent == nil {
				s.AddCmd(ishellcmd)
			} else {
				parent.AddCmd(ishellcmd)
			}
			AddCommands(s, rootcmd, ishellcmd, cmd.Commands()...)
		}
	}
}

func New() Shell {
	return ishell.New()
}
