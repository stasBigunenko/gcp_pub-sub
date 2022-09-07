# GCP PUB-SUB pet-project

## Introduction
This is a simple stream microservice example project, written in Golang and using Google Cloud Platform's service/API:

Cloud Pub/Sub for messaging
Gin-gonic/gin as HTTP web framework for data processing
PostgreSQL for storage

This was built by me, as part of my exploratory work with GCP Pub/Sub.

Running this requires a Google account and enabling GCP services, which will incur certain costs.


## Overview
This service receive data in JSON format from the front-end service, like marketplace. Data is the ID of the action and ID of the product and sends to the specified topic.
This data forwarded to the GCP.
Microservice receive this data and save it to the table ``` user_activities```.
Microservice has 4 different endpoints:
- ``` GET /bucket"```
- ``` GET /outofbucket"```
- ``` GET /description```
- ``` GET /descriptionandbucket```

You will receive the information in JSON format according the products and action what have been made with the product and between the dates you specified.
Example of request:
```
{
    "actionID":"00000000-0000-1000-0000-000000000000",
    "fromDateYear":"2022",
    "fromDateMonth":"08",
    "fromDateDay":"24",
    "toDateYear":"2022",
    "toDateMonth":"09",
    "toDateDay":"30"
} 
```
And example of response:
```
{"products":[{"actionID":"00000000-0000-1000-0000-000000000000","createdAt":"0001-01-01T00:00:00Z","productID":"00000000-0000-0000-0000-000000000001","name":"Shampoo","description":"Gel","price":100,"category":"00000000-0000-0000-1000-000000000000"}]}
```
## Getting started
First, make sure you have docker and docker-compose installed on your system. After cloning this repo to run the product with default configuration, open the root folder in bash and type:

```docker-compose up --build```

You can specify all the data of the services (ports, usernames, passwords etc) in 

```docker-compose.yaml```