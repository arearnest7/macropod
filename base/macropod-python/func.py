from rpc import Invoke

def FunctionHandler(context):
    return context["Text"], 200
