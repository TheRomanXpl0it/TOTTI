import requests
import random
import string
import time

random.seed(time.time())

FLAG_LEN = 32

def random_flag(length:int=FLAG_LEN-1):
	return "".join(random.choices(string.ascii_uppercase + string.digits, k=length)) + "="

flags = [
	{
		"flag": random_flag(),
		"username": "tofu",
		"exploit": "sploit",
		"team_ip": "10.60.1.1"
	}
	for _ in range(20)
]


r = requests.post("http://localhost:5000/api/flags", json=flags)
print(r)
