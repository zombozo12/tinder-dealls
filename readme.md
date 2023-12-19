# Technical Test
## Detail
Name: Wiguna Ramadhan

Applied Position: Senior Software Engineer

Cases: Tinder-like API

## How to run
1. Create postgres database
2. Run `misc/tinder.sql` file to create table
3. Update `config.json` file with your database and redis configuration
4. Run `go mod tidy` and `go mod vendor` to download all dependencies
5. Run `go run cmd/seeder/main.go` to seed the database
6. Run `go run main.go` to run the program

## Configuration
Configuration is stored in `config.json` file. You can change the configuration there.

## Folder Structure
```bash
tinder-dealls
├── cmd
│       ├── seeder
│       │       └── main.go
│       └── tinder-http
│           └── main.go
├── config.json
├── domain
│       ├── auth.go
│       ├── config.go
│       ├── helper.go
│       ├── inventory.go
│       ├── matched.go
│       ├── notification.go
│       ├── profile.go
│       └── user.go
├── go.mod
├── go.sum
├── handler
│       └── resthttp
│           ├── auth.go
│           ├── dependencies.go
│           ├── matcher.go
│           ├── middleware.go
│           ├── profile.go
│           ├── recommendation.go
│           ├── response.go
│           └── router.go
├── misc
│       └── tinder.sql
├── readme.md
├── repository
│       ├── auth
│       │       ├── auth.go
│       │       └── database.go
│       ├── inventory
│       │       ├── database.go
│       │       └── inventory.go
│       ├── matched
│       │       ├── database.go
│       │       └── matched.go
│       ├── notification
│       │       ├── database.go
│       │       └── notification.go
│       ├── profile
│       │       ├── database.go
│       │       └── profile.go
│       └── rds
│           └── redis.go
└── services
    ├── auth.go
    ├── auth_test.go
    ├── dependencies.go
    ├── dependencies_mock_test.go
    ├── matcher.go
    ├── matcher_test.go
    ├── profile.go
    └── recommendation.go
```

## API Documentation
### Postman Collection
You can import the postman collection from `misc/tinder.postman_collection.json` file.
### API Contract
#### Authentication
1. To sign in, call `POST /api/auth/in` with body:
    ```json
    {
        "email": "anything@mail.com",
        "password": "password"
    }
    ```
2. To sign up, call `POST /api/auth/up` with body:
    ```json
    {
        "phone": "+6212341234121",
        "code": 1234
    }
    ```
#### Profile
Authentication is required to access this endpoint. You can use `Authorization` header with value `Bearer <token>` to authenticate.
1. To create profile, call `POST /api/profile/create` with body:
    ```json
    {
        "name": "test",
        "gender": "male",
        "interest_in": "female"
    }
    ```
2. To update profile picture, call `POST /api/profile/pic` with body:
    ```json
    {
        "pic": "https://placehold.co/600x400/EEE/31343C"
    }
    ```
3. To update profile, call `POST /api/profile/update` with body:
    ```json
    {
        "name": "test",
        "gender": "male",
        "interest_in": "female"
    }
    ```
#### Recommendation
1. To get recommendation, call `GET /api/recommendation/get`

#### Matcher
1. To like someone, call `POST /api/matcher/like` with body:
    ```json
    {
        "target_user_id": 1
    }
    ```
2. To dislike someone, call `POST /api/matcher/dislike` with body:
    ```json
    {
        "target_user_id": 1
    }
    ```
3. To superlike someone, call `GET /api/matcher/superlike` with body:
    ```json
    {
        "target_user_id": 1
    }
    ```