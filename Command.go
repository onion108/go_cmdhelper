package go_cmdhelper

type Command interface {
	GetArgs() []string
	GetName() string
	IsEqualTo(another Command) bool
}

// The default implementation

type DefaultCommand struct {
	args []string
}

func MakeDefaultCommand(cmd string) *DefaultCommand {
	return &DefaultCommand{[]string{cmd}}
}

func (cmd *DefaultCommand) GetArgs() []string {
	return cmd.args
}

func (cmd *DefaultCommand) GetName() string {
	return cmd.args[0]
}

func (cmd *DefaultCommand) IsEqualTo(another Command) bool {
	return another.GetName() == cmd.GetName()
}
