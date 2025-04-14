const rpc = require('./rpc')

async function FunctionHandler(context) {
  return [context.Text.toString(), 200];
}

// Export the function
module.exports = { FunctionHandler };
