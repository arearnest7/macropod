package function

import (
	"os"
	"strings"
)

func FunctionHandler(context Context) (string, int) {
	if os.Getenv("TEST") != "" {
		payloads := [][]byte{[]byte(strings.Repeat("a", 10000000))}
		return "[" + strings.Join(RPC(context, os.Getenv("TEST"), payloads), ",") + "]", 200
	}
	return string(context.Request), 200
}

