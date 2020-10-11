# webapp
Testing auto build - test
#### To run the go server:
```
docker build -t webapp .
docker run -p 8080:8080 webapp
```

#### To run the go server and mysql:
```
docker pull mysql
docker-compose up -d
```

#### To terminate go server and mysql:
```
docker-compose down
```


#### Example create a user:
```
curl -v -X POST "http://localhost:8080/v1/user" -H  "accept: application/json" -H  "Content-Type: application/json" -d "{\"first_name\":\"Jane\",\"last_name\":\"Doe\",\"password\":\"1*Skdjfhskdfjhg\",\"username\":\"jane.doe@example.com\"}"
```
It returns Auth token and the created user.

#### Example get self (\<token> is obtained from the previous POST):
```
curl -v -X GET "http://localhost:8080/v1/user/self" -H "accept: application/json" -H  "Content-Type: application/json" -H "Authorization: Bearer <token>"
```
It returns current user's data.

#### Example update self (\<token> is obtained from the previous POST): 
```
curl -v -X PUT "http://localhost:8080/v1/user/self" -H "accept: application/json" -H  "Content-Type: application/json" -H "Authorization: Bearer <token>" -d "{\"first_name\":\"Boran\",\"last_name\":\"Yildirim\",\"password\":\"1*Skdjfhskdfjhg\",\"username\":\"jane.doe@example.com\"}"
```

#### Example get a user with id:
```
curl -v -X GET "http://localhost:8080/v1/user/<user-id>"
```
It returns user data with \<user-id>.
