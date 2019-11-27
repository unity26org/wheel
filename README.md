## Overview

Wheel is a tool for creating and maintaining scalable and lightweight RESTful APIs. It runs through command line and generates codes (in Go Language), avoiding rework when designing the application architecture and maintenance.

## Features

- [MVC](http://wheel.unity26.org/features.html#mvc)
- [RESTful](http://wheel.unity26.org/features.html#restful)
- [JWT](http://wheel.unity26.org/features.html#jwt)
- [Session controller](http://wheel.unity26.org/features.html#session-controller)
- [Middleware](http://wheel.unity26.org/features.html#middleware)
- [Authorization](http://wheel.unity26.org/features.html#authorization)
- [Users management](http://wheel.unity26.org/features.html#users-management)
- [ORM](http://wheel.unity26.org/features.html#orm)
- [Migration](http://wheel.unity26.org/features.html#migration)
- [Search engine](http://wheel.unity26.org/features.html#search-engine)
- [Pagination](http://wheel.unity26.org/features.html#pagination)
- [Ordering](http://wheel.unity26.org/features.html#ordering)
- [Sends email](http://wheel.unity26.org/features.html#sends-email)
- [Internationalization (I18n)](http://wheel.unity26.org/features.html#i18n)
- [Log](http://wheel.unity26.org/features.html#log)


See full documentation of default features at http://wheel.unity26.org/features.html


## Install

### Go

[Install Golang](https://golang.org/doc/install)

### Dependences

```
$> go get github.com/iancoleman/strcase
$> go get github.com/jinzhu/inflection
```

### Wheel

```
$> go get github.com/unity26org/wheel
$> cd GOPATH/src/github.com/unity26org/wheel
$> go build -o wheel main.go 
$> sudo mv wheel /usr/bin
```

__GOPATH__ is where the Go packages and sources are installed

The example above, the executable file was moved to _/usr/bin_. But feel free to set it up to any directory you want. Just add the path to your _.profile_, as you see below:

```
export PATH=$PATH:YOUR_DESIRED_PATH
```


## Usage

Wheel has only two options: _new_ to create new APIs and _generate_ to add new functionalities to your API. 

Check _help_ for more details.


```
wheel --help
```

### New API 

Let's create an API for a Blog.

```
wheel new github.com/account_name/blog
```

It will output something like this:

```
"Go" seems installed
Checking dependences...
         ...
Generating new app...
         created: GOPATH/src/github.com/account_name/blog
         ...

Your RESTful API was successfully created!

Change to the root directory using the command line below: 
cd GOPATH/src/github.com/account_name/blog

Set up your database connection modifying the file config/database.yml

For more details call help:
go run main.go --help
```

Remember: __GOPATH__ is where Go packages and sources are installed


### Configure Your API

####  Database

Currently, Wheel has support only for Postgresql. Edit _config/database.yml_ and set up your database connection.

#### Email

To connect to your email provider edit _config/email.yml_ and set up with your send email account.

#### Application

Edit _config/app.yml_ and set the following options:


| Item | Definition |
| ------ | ----------- |
| _app_name_ | Your app name |
| _app_repository_ | Repository name |
| _frontend_base_url_ | URL to be used on your frontend |
| _secret_key_ | Key to encrypt passwords on database |
| _reset_password_expiration_seconds_ | After reset password, how long (in seconds) is it valid? |
| _token_expiration_seconds_ | After a JWT token is generated, how long (in seconds) is it valid? |
| _locales_ | List of available locales |

#### Locales

Words and phrases for internacionalization. You can add your own locales files, but remember to add to _config/app.yml_ configuration file first.


### Running

Before running you must be sure your database schema is up to date, just run the _migrate_ mode:

```
$> go run main.go -mode=migrate
```

Run:

```
$> go run main.go
```

Now go to http://localhost:8081 and you'll see:


```
{
  system_message: {
    type: "notice",
    content: "Yeah! Wheel is working!"
  }
}
```

See full documentation of default resources at http://wheel.unity26.org/default-resources.html


### New CRUD

Based on the Blog API above, let's create a new CRUD.

Don't forget to call the directory where the application were generated.
 
```
cd GOPATH/src/github.com/account_name/blog
wheel g scaffold post title:string description:text published:bool user:references
```

It will output something like this:

```
"Go" seems installed
Checking dependences...
         ...
Generating new CRUD...
         created: app/post/post_model.go
         created: app/post/post_view.go
         created: db/entities/post_entity.go
         created: app/handlers/post_handler.go
         updated: routes/routes.go
         updated: db/schema/migrate.go
         updated: routes/authorize.go
```

After any changing, don't forget to run the _migrate_ mode:

```
$> go run main.go -mode=migrate
$> go run main.go
```

## Full Documentation

See full documentation at http://wheel.unity26.org/


## License

Wheel is released under the [MIT License](https://opensource.org/licenses/MIT).