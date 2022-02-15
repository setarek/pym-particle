## Pym Particle

Pym particle makes links shorter. it cant get simpler than that.
some pros of Pym particle:

- Has Microservice architecture
- Deletes links after 30 days
- Is super fast as a result of Redis caching
- Keeps track of visitors



### APIs

#### Post 

```
curl --location --request POST 'localhost:9006/api/v1/shorten' \
--header 'Content-Type: application/json' \
--data-raw '{
    "original_link":"http://google.com"
}'
```

#### Get

``` 
curl --location --request GET 'localhost:9005/TpXd4'
```



###  Roadmap

* [] Dockerise the project 
* [] Check inactivity on links before deleting them
* [] Adding cache rewrite mechanism