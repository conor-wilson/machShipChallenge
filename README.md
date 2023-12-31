# Git Retriever

Git Retriever is an API that was designed for a Coding Challenge as part of MachShip's hiring process. The API contains an endpoint called retrieveUsers that can be used to look up the basic information of specified GitHub users. 

## Prerequisites

To run this application you must the following installed on your machine: 
- A up-to-date Golang compiler
- The `curl` command-line tool

## Setup

After downloading the application source code, you can start the API using your favorite command terminal from inside the `machShipChallenge` directory.

```bash
go build .
./gitRetriever
```

## Usage

To use the API, execute the following in a separate terminal from the API's executable:

```bash
curl "http://localhost:8080/retrieveUsers?users=[username1]&users=[username2]&users=[username3]"
```

You must specify the usernames of the GitHub users you wish to obtain the information of using the `users` input array (`[username1]`, `[username2]` and `[username3]` in the above). There is no limit to the number of users you can include. 
