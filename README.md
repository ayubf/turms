### Pre-requisites

- [Install Nodejs](https://nodejs.org/en)

- [Install Go](https://go.dev/)

- [Install MongoDB](https://www.mongodb.com/docs/manual/administration/install-community/)

### Set Up 

1. Clone project files

```shell
    git clone https://github.com/ayubf/turms.git
```

2. Set up

Set up frontend
```shell 
    cd turms/frontend
    npm i
```

Make sure MongoDB is running after install
```shell
    mongosh 
```

3. Run 

In two seperate terminal windows, and make sure mongo:

```shell
    cd turms/frontend && npm run start
```

```shell
    cd turms/backend && go run main.go
```
Then visit https://localhost:3000