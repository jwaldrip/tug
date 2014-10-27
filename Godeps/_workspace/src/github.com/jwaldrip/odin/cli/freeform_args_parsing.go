package cli

import "github.com/nitrous-io/tug/Godeps/_workspace/src/github.com/jwaldrip/odin/cli/values"

func (cmd *CLI) assignUnparsedArgs(args []string) {
	for _, arg := range args {
		str := ""
		cmd.unparsedArgs = append(cmd.unparsedArgs, values.NewString(arg, &str))
	}
}
