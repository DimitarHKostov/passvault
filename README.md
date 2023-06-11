# passvault

### docker network create my-network

### curl -X POST  localhost:80/api/v1/login  -H "Content-Type: application/json" -d '{"password":"Test"}' -vvv

### curl -X POST -b "passvault-cookie=value" -d '{"domain":"google.com", "username":"test", "password":"test"}' localhost:80/api/v1/save -vvv

### curl -b "passvault-cookie=value" -d '{"domain":"google.com"}' localhost:80/api/v1/retrieve -X GET -vvv

### curl -b "passvault-cookie=value" -d '{"domain":"google.com", "username":"test", "password":"test2"}' localhost:80/api/v1/update -X PUT -vvv


TODO:
## logs
## unit tests
## load from env