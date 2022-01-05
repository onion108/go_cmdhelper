package go_cmdhelper

type CommandLineApplicationDelegate interface {
	applicationWillStart(application *CommandLineApplication)
	applicationWillEnd(application *CommandLineApplication)
	applicationDidEnd(application *CommandLineApplication)
}
