from parliament import Context
from flask import Request
import random
import string
import json
import uuid
import time
import os
from concurrent.futures import ThreadPoolExecutor
import requests

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
    headers = {'Content-Type': 'application/json'}
    result_txt = requests.post(os.environ["SN_COMPOSE_POST"], json=req["req"], headers=headers).text
    response = json.loads(result_txt)
    time.sleep(req["st"])
    return response

def main(context: Context):
    n_req = 11
    st = 1
    if 'request' in context.keys():
        req_j = json.loads(context.request.get_data(as_text=True))
        n_req = req_j["n"]
        st = req_j["st"]

    reqs = []
    for i in range(n_req):
        req_id = str(uuid.uuid4())
        req = {"req": generate_data(), "st": st}
        reqs.append(req)

    ex = ThreadPoolExecutor(max_workers=n_req)
    all_results_futures = ex.map(Entry, reqs)
    res = []
    for result in all_results_futures:
        res.append(result)

    return json.dumps(res), 200
