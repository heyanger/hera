# Project Hera

## How to begin?
1. Clone this project in your GOPATH, i.e. if your GOPATH is `/opt/code/go`, this project should be located in 

```
/opt/code/go/src/github.com/funkytennisball/hera
```

2. install dep and link it to path

```
curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
```

3. run dep ensure to install updates

```
dep ensure
```

4. run the program

```
go run main.o
```

## Developing Q&A

Client-side code goes in folder `client`, Frontfacing service goes in `service`