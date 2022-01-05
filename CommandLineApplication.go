package go_cmdhelper

import (
	"bufio"
	"fmt"
	"os"
)

type TriggerTable map[Command]func(Command) bool

type CommandLineApplication struct {
	args                  []string
	triggers              TriggerTable // The return value determines if the program should end
	parser                func() Command
	prompt                string
	prefix                func()
	bgTasks               []func()
	delegate              *CommandLineApplicationDelegate
	unknownCommandHandler func(command Command)
	promptProvider        func(application *CommandLineApplication) string
}

func (table TriggerTable) getTrigger(c Command) func(Command) bool {
	for k, v := range table {
		if k.IsEqualTo(c) {
			return v
		}
	}
	return nil
}

func MakeCommandLineApplication() *CommandLineApplication {
	return &CommandLineApplication{
		os.Args,
		make(map[Command]func(Command) bool),
		nil,
		">> ",
		nil,
		make([]func(), 0),
		nil,
		func(command Command) {
			fmt.Printf("Command don't exists: %s\n", command.GetName())
		},
		func(app *CommandLineApplication) string {
			return app.prompt
		},
	}
}

func (app *CommandLineApplication) SetParser(parser func() Command) {
	if parser == nil || app == nil {
		return
	}
	app.parser = parser
}

func (app *CommandLineApplication) SetTrigger(cmd Command, handler func(Command) bool) {
	if cmd == nil || handler == nil || app == nil {
		return
	}
	app.triggers[cmd] = handler
}

func (app *CommandLineApplication) SetPrompt(prompt string) {
	if app == nil {
		return
	}
	app.prompt = prompt
}

func (app *CommandLineApplication) SetPrefix(p func()) {
	if p == nil || app == nil {
		return
	}
	app.prefix = p
}

func (app *CommandLineApplication) SetDelegate(delegate *CommandLineApplicationDelegate) {
	// Allows to pass a nil in to cancel the delegate
	if app == nil {
		return
	}
	app.delegate = delegate
}

func (app *CommandLineApplication) SetUnknownCommandHandler(handler func(cmd Command)) {
	if app == nil || handler == nil {
		return
	}
	app.unknownCommandHandler = handler
}

func (app *CommandLineApplication) SetPromptProvidingMethod(m func(application *CommandLineApplication) string) {
	if app == nil || m == nil {
		return
	}
	app.promptProvider = m
}

// RunProgram Start the application
func (app *CommandLineApplication) RunProgram() {
	if app == nil {
		return
	}
	// Execute the delegate if exists
	if app.delegate != nil {
		(*app.delegate).applicationWillStart(app)
	}
	// If prefix isn't nil
	if app.prefix != nil {
		app.prefix()
	}
	// Default Parser (If didn't exist, create one lazily)
	if app.parser == nil {
		app.parser = func() Command {
			l, _, _ := bufio.NewReader(os.Stdin).ReadLine()
			return MakeDefaultCommand(string(l))
		}
	}
	// Application Loop
	for {
		print(app.promptProvider(app))
		cmd := app.parser()
		// Check if the cmd exists
		trigger := app.triggers.getTrigger(cmd)
		if trigger == nil {
			// Error occurred
			app.unknownCommandHandler(cmd)
			continue
		}
		shouldExit := trigger(cmd)
		if shouldExit {
			// Execute the delegate if exists
			if app.delegate != nil {
				(*app.delegate).applicationWillEnd(app)
			}
			break
		}
	}
	// Execute the delegate if exists
	if app.delegate != nil {
		(*app.delegate).applicationDidEnd(app)
	}
}
