
# idcra-telegram-reminder

This is the codebase for Indonesia Digital Caries Risk Assessment Telegram Reminder, or IDCRA Telegram Reminder for short.

# Original Readme

#### Structure
    go-graphql-starer
    │   README.md
    │   server.go           --- the entry file
    │   .env                --- the configuration file
    └───app                 --- application app like db configuration
    └───data                --- storing the sql data patch for different version
    │   └───1.0             --- storing sql data patch for version 1.0
    │      └───...          --- sql files
    └───controller          --- the controller to define method, controller will call service
    └───model               --- the folder putting struct file
    └───service             --- services for users, authorization etc.
    └───util                --- utilities

#### Requirement:

1. Mysql database
2. Golang

Remark: If you want to use other databases, please feel free to change the driver in `app/database.go`

#### Usage(Without docker):

1. Run the sql scripts under `data/1.0` folder inside Mysql database console

2. Sync Dependency

3. Setup GOPATH (Optional if already set)

   For example: MacOS
    ```
    export GOPATH=/Users/${username}/go
    export PATH=$PATH:$GOPATH/bin
    ```

4. Start the server (Ensure your mysql database is live and its setting in .env is correct)
    ```
    go build server.go
    ```