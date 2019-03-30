import sys
import requests

url = "http://localhost:8080/"

"""
python3 heraclt.py get key
python3 heraclt.py put key value
python3 heraclt.py delete key
"""

def getValue(key):
	path = url + "get" + key 
	response = requests.get(path)
	print(response.status_code)


def putPair(key, value):
	payload = {"key": key, "value": value}
	response = requests.put(url, json=payload)
	print(response.status_code)


def delete(key):
	payload = {"key": key}
	response = requests.delete(url, json=payload)
	print(response.status_code)


def main():
	length = len(sys.argv)
	if length < 3 or length > 4:
		print("Invalid Command")
		return

	op = sys.argv[1]
	if op == 'get' and length == 3:
		getValue(sys.argv[2])
	elif op == 'put' and length == 4:
		putPair(sys.argv[2], sys.argv[3])
	elif op == 'delete' and length == 3:
		delete(sys.argv[2])
	else:
		print("Invalid Command")

if __name__ == "__main__":
	main()




