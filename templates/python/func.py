from rpc import RPC

def FunctionHandler(context):
    return str(context["Text"]), 200
