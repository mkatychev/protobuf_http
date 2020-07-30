#!/usr/bin/env bash

NOTE_NAME=$(gofaker name first)
curl -d '{"name":"'$NOTE_NAME'"}' localhost:8080/notebook
NOTE_ID=$(curl -d '
{
  "notebook_name": "'$NOTE_NAME'",
  "title": "newNote",
  "body": "shoddy",
  "tags": [
    "homie",
    "growie"
  ]
}
' localhost:8080/note | jq -r .id)
echo $NOTE_ID
UPDATE_BODY=$(
cat <<EOF 
{
  "id": "$NOTE_ID",
  "notebook_name": "$NOTE_NAME",
  "title": "newerNote",
  "body": "clean",
  "tags": [
    "crawl",
    "growl"
  ]
}
EOF
)
echo "$UPDATE_BODY"
curl -X UPDATE -d  "$UPDATE_BODY" localhost:8080/note

