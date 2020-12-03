package channel

import (
	"fmt"
	"io"
	"os"
)

type StdoutChannel struct {
	Config
}

func (sc *StdoutChannel) Send(msg string) bool {
	_, err := fmt.Fprint(os.Stdout, msg)
	return err == io.EOF
}
