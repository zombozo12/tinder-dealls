# Tinder-like API

## Document Signature

|                       |                                            |
|-----------------------|--------------------------------------------|
| Creator               | Wiguna Ramadhan                            |
| Github                | https://github.com/zombozo12/tinder-dealls |
| LinkedIn              | https://linkedin.com/in/wigunarrr          |

## Background / Problem Statement
Tinder-like API is a technical test from Dealls. 
This API is used to find a match between two users. 
This API is built using Go and PostgreSQL.

## Goals
1. Create a Tinder-like API
2. Create a documentation for the API
3. Create a test for the API

## Non-Goals
1. Create a frontend for the API

## Functional Requirements
1. Golang 1.21

    A great language for backend development. I choose this language because it's fast, easy to use, and has a great community.
    For this project, I use Go 1.21 because it's the latest version of Go. Fiber is a great framework for Go, it's fast and easy to use.
    I use GORM as the ORM for this project because it's easy to use and has a great community.
2. PostgreSQL

    PostgreSQL is a great database for relational data. I choose this database because it's easy to use and has a great community.
    I use pgx as the PostgreSQL driver for this project because it's fast and easy to use.
3. Redis

    Redis is a great database for caching. I choose this database because it's fast and easy to use.
    I use redigo as the Redis driver for this project because it's fast and easy to use.

### Non-Functional Requirements
1. Sanity

    I use `golangci-lint` to check the sanity of the code.
2. Faith

    I use `go test` to check the faith of the code.

## Design
### ERD
![image](https://github.com/meshery/meshery/assets/21243980/c6673507-ceba-497e-b3de-71f9b391f679)

## Test Cases
### Unit Test
1. Test for authentication services