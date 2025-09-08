# MacroPod Function Templates

These are the templates for all MacroPod functions. Please copy the directory for the language your function is built in. Please note that:

1. The template uses the MacroPod go, node and python base templates for the specific branch, so if you wish to change the image base please modify the appropriate line in the Dockerfile.
2. Each language has their own specific way of handling imports. If you need more imports than the standard requirements.txt, func.js or go.mod file import please add the appropriate lines to the Dockerfile.
3. Each language has their own format for RPCs for downstream. While the format is similar please consult the below examples.
4. All RPC calls fit the format of context (passed in parameter to function), name of downstream function (refer to workflow definition) and payload/payloads of specified type.

Go Example:

	package function

	import (

	)

	func FunctionHandler(context Context) (string, int) {
		ret_val := ""
		ret_code := 200
		ret_vals := []string{}
		ret_codes := []int32{}
		ret_val, ret_code = Invoke(context, "FOO", "txt")
		ret_val, ret_code = Invoke_JSON(context, "FOO", map[string]interface{"bar": "txt"})
		ret_val, ret_code = Invoke_Data(context, "FOO", []byte{0,0,0,0,0,0,0,0})
		ret_vals, ret_codes = Invoke_Multi(context, "FOO", []string{"txt"})
		ret_vals, ret_codes = Invoke_Multi_JSON(context, "FOO", []map[string]interface{}{map[string]interface{"bar": "txt"}})
		ret_vals, ret_codes = Invoke_Multi_Data(context, "FOO", [][]byte{[]byte{0,0,0,0,0,0,0,0}})
		return context.Text, 200
	}

NodeJs Example:

	const rpc = require('./rpc')

	async function FunctionHandler(context) {
	  let [res, code] = await rpc.Invoke(context, "FOO", "txt");
	  let [res, code] = await rpc.Invoke_JSON(context, "FOO", {"bar": "txt"});
	  let [res, code] = await rpc.Invoke_Data(context, "FOO", new ArrayBuffer(8));
	  let [results, codes] = await rpc.Invoke_Multi(context, "FOO", ["txt"]);
	  let [results, codes] = await rpc.Invoke_Multi_JSON(context, "FOO", [{"bar": "txt"}]);
	  let [results, codes] = await rpc.Invoke_Multi_Data(context, "FOO", [new ArrayBuffer(8)]);
	  return [context.Text, 200];
	}

	// Export the function
	module.exports = { FunctionHandler };

Python Example:

	from rpc import Invoke

	def FunctionHandler(context):
	    ret_val, ret_code = Invoke(context, "FOO", "txt")
	    ret_val, ret_code = Invoke_JSON(context, "FOO", {"bar": "txt"})
	    ret_val, ret_code = Invoke_Data(context, "FOO", b'txt')
	    ret_vals, ret_codes = Invoke_Multi(context, "FOO", ["txt"])
	    ret_vals, ret_codes = Invoke_Multi_JSON(context, "FOO", [{"bar": "txt"}])
	    ret_vals, ret_codes = Invoke_Multi_Data(context, "FOO", [b'txt'])
	    return context["Text"], 200