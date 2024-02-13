package function

import (
	"os"
	"strings"
)

func function_handler(context Context) (string, int) {
	if context.request_type != "GRPC" {
		payloads := []string{"TEST"}
		return "[" + strings.Join(RPC(os.Getenv("TEST"), payloads, context.workflow_id), ", ") + "]", 200
	}
	return context.request, 200
}

