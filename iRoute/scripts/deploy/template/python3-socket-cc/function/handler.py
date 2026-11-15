# Replace the following with your code
workflow_name = "workflow"
writer_up = ("writer_user", "writer_password")
reader_up = ("reader_user", "reader_password")

def get_lat_dist(lats, dag):
    """
    Calculate the latency distribution from the list of latencies.
    Args:
        lats: [(message_data, timestamp)]
        - message_data.req_id: str, request id
        - message_data.req: str, request result of a final stage func
        - timestamp: float, timestamp of receiving the message, in seconds

        dag: list of stages in the workflow, each stage is a list of function names

    Please replace this with your actual implementation.
    """

    return ""
    


def handle(req):
    """handle a request to the function
    Args:
        req (str): request body
    """

    return req
