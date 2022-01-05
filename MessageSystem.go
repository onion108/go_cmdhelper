package go_cmdhelper

import "log"

type MessagingDelegate interface {
	messageReceived(message *Message)
	messageWillExecute(message *Message)
	messageDidExecute(message *Message)
}

/* Message */

const (
	MSGTYPE_BUILTIN_TERMINATE = "[terminate_t]"
	MSGTYPE_BUILTIN_SKIP      = "[skip_t]"
	MSGTYPE_BUILTIN_EXECNOW   = "[execute_only]"
)

type Message struct {
	executable func()
	data       interface{}
	msgtype    string
	delegate   *MessagingDelegate
}

func MakeMessage(msgtype string, data interface{}) Message {
	return Message{nil, data, msgtype, nil}
}

func (msg *Message) BindDelegate(delegate *MessagingDelegate) {
	if msg == nil || delegate == nil {
		return
	}
	msg.delegate = delegate
}

func (msg *Message) SetExecutable(executable func()) {
	if msg == nil || executable == nil {
		return
	}
	msg.executable = executable
}

func (msg *Message) GetType() string {
	if msg != nil {
		return msg.msgtype
	}
	return "FUCK YOUR MOM THAT YOU CALL THIS METHOD ON NIL!!!?"
}

func (msg *Message) GetData() interface{} {
	if msg != nil {
		return msg.data
	}
	return nil
}

/* Message Thread */

type MessageListener struct {
	onHandleMessage func(message Message, output func(Message))
}

func MakeMessageListener() MessageListener {
	return MessageListener{nil}
}

func (msgListener *MessageListener) SetMessageListener(l func(message Message, output func(Message))) {
	if msgListener == nil || l == nil {
		return
	}
	msgListener.onHandleMessage = l
}

func (msgListener *MessageListener) ListenMessage(inputChan <-chan Message, outputChan chan<- Message) {
	// Arguments checking
	if msgListener == nil || inputChan == nil || outputChan == nil {
		return
	}
	go func() {
		for {
			nextMsg := <-inputChan
			// Builtin message handling
			if nextMsg.msgtype == MSGTYPE_BUILTIN_TERMINATE {
				break
			}
			if nextMsg.msgtype == MSGTYPE_BUILTIN_SKIP {
				continue
			}
			if nextMsg.msgtype == MSGTYPE_BUILTIN_EXECNOW {
				nextMsg.executable()
				continue
			}
			// Otherwise, call the default callback
			if msgListener.onHandleMessage == nil {
				log.Fatalf("Unhandled message:\n\ttype:%s\n\t", nextMsg.msgtype)
			}
			msgListener.onHandleMessage(nextMsg, func(m Message) {
				go func() { outputChan <- m }() // Avoid blocking
			})
		}
	}()
}
