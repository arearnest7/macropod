import redis

state_list = ['AK', 'AL', 'AR', 'AZ', 'CA', 'CO', 'CT', 'DC', 'DE', 'FL', 'GA', 'HI', 'IA', 'ID', 'IL', 'IN', 'KS', 'KY', 'LA', 'MA', 'MD', 'ME', 'MI', 'MN', 'MO', 'MS', 'MT', 'NC', 'ND', 'NE', 'NH', 'NJ', 'NM', 'NV', 'NY', 'OH', 'OK', 'OR', 'PA', 'RI', 'SC', 'SD', 'TN', 'TX', 'U']
candidates = ['John_Smith', 'Adam_Carter', 'Jane_Langley']

def main():
	r = redis.Redis(host='10.125.189.107', port=6379, password='password1234')
	for i in range(2000):
		r.set("voter-" + str(i), "")
main()
