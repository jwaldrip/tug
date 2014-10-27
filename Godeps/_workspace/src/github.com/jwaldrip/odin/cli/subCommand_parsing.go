package cli

func (cmd *CLI) parseSubCommands(args []string) ([]string, bool) {
	if len(args) == 0 || len(cmd.subCommands) == 0 {
		return args, false
	}
	name := args[0]
	subcmd, ok := cmd.subCommands[name]
	if !ok {
		cmd.errf("invalid command: %s", name)
		return args, false
	}

	// Inherit Outputs
	if subcmd.errOutput == nil {
		subcmd.errOutput = cmd.errOutput
	}
	if subcmd.stdOutput == nil {
		subcmd.stdOutput = cmd.stdOutput
	}

	subcmd.Start(args...)

	return []string{}, true
}
