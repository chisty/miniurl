# miniurl
microservice to make short url

## Features
- Create unique short URL for each request URL/Link.
- Has cache support which will remain stored for 30 days. Cache will refresh for 30 days again if it is used before expiration.
- Basic authentication service implemented with jwt token.
- Unit test added.
- Has docker support for both local and prod environment with separate configuration.

## Tool & Reason
#### Primary database is Dyanmo DB
- It is a faster & simple key-value store. 
- For this project it is a good option, since in my design i have simple model to save. I am generating unique short id(key) for each requested url and save them in dynamo. It will be fast to fetch them by the id since it is primary key. 
- I am generating unique id for each request because it will keep the implementation simple & we have enough space to accommodate more than ~100 years data in 100k write scale.
#### Redis cache
-	docker for local, ElastiCache in AWS for PROD. 
-	Redis support TTL (time to live) for cache which is handy.
#### Snowfalke Id Generator
-	Snowflake is twitter's popular distributed id generator.
-	 I have used this as a package in this project because in real scenario this id generator should be separate microservice.
#### OAuth/JWT Token
-	I have created a dummy jwt token from jwt.io website. I have used HMAC algorithm & "my_super_secret_key" as a secret key to demonstrate that this app can parse jwt token from bearer token and authenticate request. 
-	In real scenario, token /refresh token generation will be handled by separate microservice.
#### docker-compose file
-	docker-compose.dev.yml file is to support local development whereas docker-compose.prod.yml will support prod deployment. 
-	I have listed all necessary environment variables in docker-compose file which should be saved in PROD env variable or CI/CD secret in real scenario.


##How to Run
###Local Environment

Running below command in console will fetch docker images and run miniurl app on port 9000 & 9001 in local environment:

    docker-compose   -f   docker-compose.dev.yml  up
    
Please notice that environment variables are mandatory in docker-compose file:

    AWS_REGION=ap-southeast-1
    AWS_ACCESS_KEY=XXXXXXXXXX
    AWS_SECRET_KEY=XXXXXXXXXXXXXXXXXXXXXXXXXXXXX
    REDIS_URL=redis:6379
    REDIS_TTL=10                        //TTL is saved in second(s). locally it is 10 second. 
    JWT_SECRET=my_super_secret_key      //secret key for jwt token
