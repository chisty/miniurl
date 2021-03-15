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


## How to Run
### Local Environment

Running below command in console will fetch docker images and run miniurl app on port 9000 & 9001 in local environment:

    docker-compose   -f   docker-compose.dev.yml  up
    
Please notice that environment variables are mandatory in docker-compose file:

    AWS_REGION=ap-southeast-1
    AWS_ACCESS_KEY=XXXXXXXXXX
    AWS_SECRET_KEY=XXXXXXXXXXXXXXXXXXXXXXXXXXXXX
    REDIS_URL=redis:6379
    REDIS_TTL=10                        //TTL is saved in second(s). locally it is 10 second. 
    JWT_SECRET=my_super_secret_key      //secret key for jwt token
    
    
Local environment is using AWS dynamoDB as main db, redis docker as cache db. To connect with the dynamo db, we need the Access & Secret Keys.
After executing docker-compose, go to localhost: 9001 or 9002.

![2021-03-15_04h26_48](https://user-images.githubusercontent.com/2306114/111114068-50f84000-8595-11eb-9bda-f5a9a0951f35.png)


### POST Call/ Save URL:

We can now save data using Post call to our app. Using Postman, it will be like below:

![2021-03-15_04h30_49](https://user-images.githubusercontent.com/2306114/111114181-75541c80-8595-11eb-9ead-73de0c63d7b5.png)

However, before that we have to set our bearer token here:

![2021-03-15_04h31_55](https://user-images.githubusercontent.com/2306114/111114237-8a30b000-8595-11eb-9b3c-8d5e3e80813f.png)

For now we can use the defined token below as bearer token

    eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHAiOiJtaW5pdXJsIn0.lYWHpDPJYf5NjJQlYCLgzmiOeqfKLa7V5qEheyTmOUc
    


JWT token can be generated from jwt.io webpage. miniurl app only checks if it has any claim with id: ***“app” & value: “miniurl”***. Currently it does not check any expiration data. Notice, the secret key is defined as “**my_super_secret_key**”

![2021-03-15_04h38_43](https://user-images.githubusercontent.com/2306114/111114509-ebf11a00-8595-11eb-9b12-835948c5ccc2.png)

Copying the token from Encoded section (3) and putting it as bearer token will be sufficient for now. Both **GET & POST** is validating our basic JWT token.

### GET Call/ Fetch URL:

Setting the jwt token & putting **{id}** in the route will fetch the stored long URL

![2021-03-15_04h47_45](https://user-images.githubusercontent.com/2306114/111114638-14791400-8596-11eb-807e-531b7ae3d230.png)

For now, app is returning 307 (temporary redirect) with full response. 


## PROD Environment

In production, miniurl is running on AWS free tier only. 
- Deployed on a docker container inside free ec2 instance. 
- PROD and Local app using same dynamo db table.
- PROD using AWS ElastiCache (redis) free version.

Testing PROD instance is same like local instance. Both needs same bearer token now. PROD web URL:

    http://ec2-54-251-165-137.ap-southeast-1.compute.amazonaws.com/
    
![2021-03-15_04h59_22](https://user-images.githubusercontent.com/2306114/111114743-41c5c200-8596-11eb-986c-d3a1246b86a6.png)

### POST data in PROD:

![2021-03-15_04h59_55](https://user-images.githubusercontent.com/2306114/111114794-55712880-8596-11eb-9c89-214511daf69e.png)

### GET data in PROD:

![2021-03-15_05h00_41](https://user-images.githubusercontent.com/2306114/111114824-615cea80-8596-11eb-9f28-a8e02161007d.png)


    
