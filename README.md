## POC of distributed file storage

The minimal implementation of distributed file storage.

## Requirements
- We have single incoming http server where file should be uploaded throught PUT request.
- The incoming file should be devided to several parts (chunks) and placed to all available storge servers evenly.
- We need 2 end-points:
  - PUT /upload
  - GET /download/{filename} - the end-point should join all chunks in realtime and return the file to a client.

## Run
1. We must compile service:
```
go mod tidy
go build
```
2. Run service in first terminal:
```
./server
```
3. Test using `curl` (You have to have `curl` utility in your system):
```
curl --verbose --request PUT -F file=@example.jpg http://0.0.0.0:8080/upload
curl -o new_example.jpg http://0.0.0.0:8080/download/example.jpg
```