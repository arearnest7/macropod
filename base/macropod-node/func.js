const rpc = require('./rpc')

async function FunctionHandler(context) {
  return [context.Text, 200];
}

// Export the function
module.exports = { FunctionHandler };
