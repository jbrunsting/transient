#!/bin/sh

SESSION_ID=$(curl -v 'http://localhost:443/api/user/login' -H 'Origin: http://localhost:443' -H 'Accept-Encoding: gzip, deflate, br' -H 'Accept-Language: en-US,en;q=0.9' -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.108 Safari/537.36' -H 'Content-Type: application/json;charset=UTF-8' -H 'Accept: application/json, text/plain, */*' -H 'Referer: http://localhost:443/' -H 'Connection: keep-alive' --data-binary '{"username":"ff","password":"ff"}' --compressed 2>&1 | grep -Po "(?<=sessionId=).*(?=; Path=)")

TITLE=$(shuf -n 1 titles.txt)
SENTANCE=$(shuf -n 1 sentances.txt)
echo "${TITLE}"
echo "${SENTANCE}"
curl -v 'http://localhost:443/api/post' -H 'Origin: http://localhost:443' -H 'Accept-Encoding: gzip, deflate, br' -H 'Accept-Language: en-US,en;q=0.9' -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.108 Safari/537.36' -H 'Content-Type: application/json;charset=UTF-8' -H 'Accept: application/json, text/plain, */*' -H 'Referer: http://localhost:443/profile' -H "Cookie: sessionId=${SESSION_ID}" -H 'Connection: keep-alive' --data-binary "{\"title\":\"${TITLE}\",\"content\":\"${SENTANCE}\"}" --compressed
