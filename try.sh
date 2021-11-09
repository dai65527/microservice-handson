#!/bin/bash
set -eux

curl -s -XPOST -d '{"name":"gopher"}' localhost:4000/auth/signup | jq .
TOKEN=$(curl -s -XPOST -d '{"name":"gopher"}' localhost:4000/auth/signin | jq .access_token -r)
curl -s -XGET -H "authorization: bearer $TOKEN" localhost:4000/catalog/items | jq .
curl -s -XPOST -d '{"title":"Keyboard","price":4000}' -H "authorization: bearer $TOKEN" localhost:4000/catalog/items | jq .
curl -s -XGET -H "authorization: bearer $TOKEN" localhost:4000/catalog/items | jq .

curl -s -XGET -H "authorization: bearer $TOKEN" localhost:4000/catalog/items/e0e58243-4138-48e5-8aba-448a8888e2ff | jq .

