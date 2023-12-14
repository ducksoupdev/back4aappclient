# back4app

Go library for accessing the [Back4App API](https://www.back4app.com/)

## Installation ##

Back4app is compatible with modern Go releases in module mode, with Go installed:

```bash
go get github.com/ducksoupdev/back4app
```

will resolve and add the package to the current development module, along with its dependencies.

Alternatively the same can be achieved if you use import in a package:

```go
import "github.com/ducksoupdev/back4app/user"
import "github.com/ducksoupdev/back4app/object"
```

and run `go get` without parameters.

## Usage

Construct a new user, then use the methods on the user to
login, sign up and request password reset. For example:

```go
u := user.NewUser("applicationId", "restApiKey")

// login user
sessionToken, err := user.login("username", "password")

// sign up user
var data = make(map[string]interface{})
data["username"] = "username"
data["password"] = "password"
sessionToken, _ := u.signUp(data)

// request password reset
err := u.requestPasswordReset("email")
```

Construct a new object, then use the methods on the object to
create, update, delete and query objects. For example:

```go
o := object.NewObject("applicationId", "restApiKey", "sessionToken")

// create object
var data = make(map[string]interface{})
data["name"] = "name"
data["description"] = "description"
object, err := o.create("className", data)

// update object
var data = make(map[string]interface{})
data["name"] = "name"
data["description"] = "description"
isUpdated, err := o.update("className", "objectId", data)

// delete object
isDeleted, err := o.delete("className", "objectId")

// read object
object, err := o.read("className", "objectId")

// list objects
objects, err := o.list("className")
```

## Contributing

TODO
