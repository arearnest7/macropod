from rpc import Invoke_JSON
import random
import string
import json
import uuid
import time
from concurrent.futures import ThreadPoolExecutor

def gen_random_string(i):
    choices = string.ascii_letters + string.digits
    return "".join([choices[random.randint(0, len(choices)-1)] for j in range(i)])

def generate_data():
    # func for generate req for message
    user_index = random.randint(1, 1000)

    # num_user_mentions = random.randint(1, 5)
    num_user_mentions = 1
    num_urls = random.randint(1, 5)
    num_media = random.randint(1, 4)
    text = ""

    user_mentions = []
    for i in range(num_user_mentions):
        while True:
            user_mention_id = random.randint(1, 1000)
            if user_mention_id != user_index and user_mention_id not in user_mentions:
                user_mentions.append(user_mention_id)
                break

    text += " ".join(["@username_"+str(i) for i in user_mentions])
    text += " "

    for i in range(num_urls):
        text += ("http://"+gen_random_string(64)+" ")

    media_ids = []
    media_types = ["png"] * num_media
    for i in range(num_media):
        media_ids.append(gen_random_string(5))

    data = {
        "username": "username_" + str(user_index),
        "user_id": user_index,
        "post_type": 0,
        "text": text,
        "media_ids": ",".join(media_ids),
        "media_types": ",".join(media_types)
    }

    return data

def Entry(req):
    result_txt, code = Invoke_JSON(req["context"], "SN_COMPOSE_POST", req["req"])
    response = json.loads(result_txt)
    time.sleep(req["st"])
    return response

def FunctionHandler(context):
    n_req = 11
    st = 1
    if context["Text"] != "":
        req_j = json.loads(context["Text"])
        n_req = req_j["n"]
        st = req_j["st"]

    reqs = []
    for i in range(n_req):
        req_id = str(uuid.uuid4())
        req = generate_data()
        reqs.append({"req": req, "context": context, "st": st})

    ex = ThreadPoolExecutor(max_workers=n_req)
    all_results_futures = ex.map(Entry, reqs)
    res = []
    for result in all_results_futures:
        res.append(result)

    return json.dumps(res), 200
