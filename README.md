# passvault

### curl -X POST  localhost:80/api/v1/login  -H "Content-Type: application/json" -d '{"password":"Test"}' -vvv

### curl -b "passvault-cookie=cookie_value"  localhost:80/api/v1/logout -vvv

### curl -X POST -b "passvault-cookie=value" -d '{"domain":"google.com", "username":"test", "password":"test"}' localhost:80/api/v1/save -vvv

### curl -b "passvault-cookie=cookie_value"  localhost:80/api/v1/retrieve -vvv