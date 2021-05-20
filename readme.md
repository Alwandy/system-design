# System Design By Julian Alwandy

![Image of Infrastructure idea](https://github.com/Alwandy/system-design/blob/master/infrastructure.png?raw=true)


The design is based on using Lambda with DynamoDB to store the keys. 

With autoscaling group and Lambda, we are able to define in the specifications when to scale up and down based on the requirements 
- Scaling target: 1000+ req/s, after scaling-up/out without major code change (https://docs.aws.amazon.com/lambda/latest/dg/invocation-scaling.html) 

Why have we gone with DynamoDB? 
Reason being is that with this application, we only needed to store the shortenUrl with the long url and get the item based on shortenUrl ID. Using NoSQL is faster rather than relationship for these type of queries. On top of that, AWS offers autoscaling with dynamodb that we can configure based on requirements stated in project, https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/AutoScaling.html 

In term the infrastructure is very simple and requires no maintenance as AWS DynamoDB is a managed service and the lambda function is only invoked on call. This will keep cost down and allow scaling based on invocation and minimize any impact when scaling up. 

NOTE: 
To run the Dockerfile, you need to run with ENV variables:
```
export AWS_ACCESS_KEY_ID=
export AWS_SECRET_ACCESS_KEY=
export AWS_DEFAULT_REGION=eu-west-1
``` 
As the code will initially on every run to try create table if required. If more time, just do a check if table already exists before trying to create or go with assumption it does exist already through infrastructure as code. 


Some credits to: https://github.com/donnemartin/system-design-primer/tree/master/solutions/system_design/pastebin 