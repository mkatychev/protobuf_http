Example Protobuf implementation scoped to HTTP development

## Notebook
A `Notebook` is an object containing several `Notes`:

```proto
message Note {
  string                    id            = 1;
  string                    title         = 2;
  string                    body          = 3;
  repeated string           tags          = 4;
  google.protobuf.Timestamp created       = 5;
  google.protobuf.Timestamp last_modified = 6;
}
```

* Notebooks can be retrieved along with the metadata of the notes within, if a `tags` is provided, only the relevant notes will be returned
* Notes can be created, updated, deleted, and retrieved

## Setup

* `brew install clang-format protoc-gen-go protobuf` to setup on macOS for development
* `./proto-format.sh` will autoformat the proto upon modification
* to generate/regenerate `*.pb.go` files: `protoc notebook.proto --go_out=.`

## Running
- `go run .` will run the service on `localhost:8080`
- `NTBK_PORT="9001" go run .` will run the service on `localhost:9001`

## Runnig using `docker`
- To build: `docker build --tag notebook:1.0 .`
- And run: `docker run -it --publish 8080:8080 notebook:1.0`

## Example Commands
* Use `curl.sh` to see example interactions:

```sh
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
```

* **NOTE**: `curl.sh` port and notebook name can be overridden: `NTBK_PORT="9001"  NOTE_NAME="Not_Test_Note" go run .`
