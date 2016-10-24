package value_type

type baseOperation struct {
	err error
	commandName string
}

func (o *baseOperation) GetError() error {
	return o.err
}

func (o *baseOperation) GetCommandName() string {
	return o.commandName
}
