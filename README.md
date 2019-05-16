# Naturae Server
This is the server implementation of Naturae, the client code can be found in this repository,
[Naturae UI](https://github.com/nanaeaubry/Naturae_UI) <br />
In order to receive responses from the client, it must be connected to the running server on localhost or some other hosting method.
In our case, we chose AWS, Amazon Web Services as our provider.

## Deployment
The server requires that you have Golang installed on your system. Refer to the repository
[Go]https://github.com/golang/go for more information on how to install it. You must also have configured credentials for MongoDB, Amazon S3, and AWS(or some other provider) in order to run the server. For obvious reasons, we cannot provide our own. 

Once everything is setup, navigate to the main directory and run the command
```
go run main.go
```

## Implementation
The server is implemented using gRPC framework. gRPC functions can be viewed in the included Naturae proto file. Subsequent updates to it will require the protobuf compiler and Go plugin. Although a pb.go file has already been included. The server code imports this generated stub and implements them accordingly.
```
protoc --go_out=plugins=grpc:. naturaeproto/Naturae.proto
```

## Built With
* [Go]https://github.com/golang/go 
* [AWS]https://aws.amazon.com/
* [S3]https://aws.amazon.com/s3/
* [MongoDB]https://www.mongodb.com/
* [gRPC-go]https://github.com/grpc/grpc-go
* [protobuf]https://github.com/protocolbuffers/protobuf

## Contributors

* **Visal Hok** -  [HappyLyfe123](https://github.com/HappyLyfe123)
* **Steven Lim** - [LimStevenLBW](https://github.com/LimStevenLBW)

## License

This project is licensed under the MIT License, see the [LICENSE.md](LICENSE.md) file for details
