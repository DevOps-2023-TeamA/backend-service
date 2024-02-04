# TSAO Backend
## contents

## setup
1. install Go (https://go.dev/doc/install)
2. create an .env file with the secret key (instructions below in bash). 
```bash
echo "SECRET_KEY=THE SECRET KEY" > .env
```
> This is what the `.env` file should look like.
> ```env
> SECRET_KEY=THE SECRET KEY
> ```

> [!Warning]
> Do not commit the `.env` file!
> 
> It should be automatically ignored as part of the `.gitignore`.

3. create your database string in the format below
```
username:password@tcp(127.0.0.1:3306)
```

## serving
> [!note]
> Make sure the database is set up and being served on port 3306.
>
> For setup instructions, refer to [tsao-db](https://github.com/DevOps-2023-TeamA/tsao-db).

### on macOS/linux
```bash
# serve auth microservice
go run microservices/auth/*.go -sql "DATABASE STRING"

# serve accounts microservice
go run microservices/accounts/*.go -sql "DATABASE STRING"

# serve records microservice
go run microservices/records/*.go -sql "DATABASE STRING"
```
> [!Tip]
> Add a `&` after the command to run it in the background
> 
> ```
> go run microservices/auth/*.go -sql "DATABASE STRING" &
> ```
>
> To kill the service,
> 1. Retrieve the PID: Run `sudo lsof -i :PORT_NUMBER` replacing `PORT_NUMBER` with the port number of the service (see [here](#port-references)).
> 2. Run `sudo kill -9 PID`, replacing PID with the service's PID.

### on windows
```bash
# serve auth microservice
go run microservices/auth/. -sql "DATABASE STRING"

# serve accounts microservice
go run microservices/accounts/. -sql "DATABASE STRING"

# serve records microservice
go run microservices/records/. -sql "DATABASE STRING"
```

## port references
| port | what's running?       |
|------|-----------------------|
| 8000 | auth microservice     |
| 8001 | records microservice  |
| 8002 | accounts microservice |

## maintainers
- [Yee Jia Chen](https://github.com/jiachenyee) S10219344C
- [Isabelle Pak Yi Shan](https://github.com/isabellepakyishan) S10222456J
- [Ho Kuan Zher](https://github.com/Kuan-Zher) S10223870D
- [Cheah Seng Jun](https://github.com/DanielCheahSJ) S10227333K
- [Chua Guo Jun](https://github.com/GuojunLoser) s10227743H
