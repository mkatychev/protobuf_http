#!/usr/bin/env bash


NOTE_NAME=${NOTE_NAME:-"Test_Note"}
NTBK_PORT=${NTBK_PORT:-8080}

# r.HandleFunc("/notebook", repo.CreateNotebook).Methods("POST")
curl -d '{"name":"'$NOTE_NAME'"}' localhost:$NTBK_PORT/notebook
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
' localhost:$NTBK_PORT/note | jq -r .id)

echo $NOTE_ID

# r.HandleFunc("/note", repo.UpdateNote).Methods("UPDATE")
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
# r.HandleFunc("/note", repo.UpdateNote).Methods("UPDATE")
curl -X UPDATE -d  "$UPDATE_BODY" localhost:$NTBK_PORT/note


# r.HandleFunc("/notebook", repo.GetNotebook).Methods("GET")
curl -X GET -d "{\"name\": \"$NOTE_NAME\"}" localhost:$NTBK_PORT/notebook
