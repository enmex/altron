import re

def remove_spaces_until_first_char(s):
    index = 0
    while index < len(s) and s[index].isspace():
        index += 1
    return s[index:]

def generate_sploit(data):
    message = ''
    url = ''

    for item in data:
        for line in str(item).split('\n'):
            if 'import requests' in line:
                continue

            if re.search('requests.[\w]+\((.*)\)', line):
                line = line.replace('requests', 's')
                try:
                    line_url = re.search("s.[\w]+\(\.*('http[s]?:\/\/[a-z0-9A-Z.:]+).*\)", line).group(1)
                except:
                    line_url = re.search("s.[\w]+\(\.*('http[s]?:\/\/[\[a-z0-9A-Z.:\]]+).*\)", line).group(1)
                url = line_url
                line = line.replace(line_url, "url + '")

            line = remove_spaces_until_first_char(line)

            message += line + '\n'
    message += 's.close()'

    start = '''import sys
import requests
import sys
import random
import string

def generator(size=12, chars=string.digits + string.ascii_letters):
    return ''.join(random.choice(chars) for _ in range(size))

port = 9090
#url = "http://" + str(sys.argv[1]) + ":" + str(port)
url = {}'

s = requests.Session()'''.format(url)
    sploit = start + message
    return sploit