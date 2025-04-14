package function

import (

)

func FunctionHandler(context Context) (string, int) {
	return string(context.Text), 200
}

