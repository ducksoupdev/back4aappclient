# back4app

Go library for accessing the [Back4App API](https://www.back4app.com/)

## Installation ##

Back4app is compatible with modern Go releases in module mode, with Go installed:

```bash
go get github.com/ducksoupdev/back4app/user
go get github.com/ducksoupdev/back4app/object
go get github.com/ducksoupdev/back4app/util
```

will resolve and add the package to the current development module, along with its dependencies.

Alternatively the same can be achieved if you use import in a package:

```go
import "github.com/ducksoupdev/back4app/user"
import "github.com/ducksoupdev/back4app/object"
import "github.com/ducksoupdev/back4app/util"
```

and run `go get` without parameters.

## Usage

### User

Construct a new user, then use the methods on the user to
login, sign up and request password reset. For example:

```go
u := user.NewUser("applicationId", "restApiKey")

// login user
sessionToken, err := user.Login("username", "password")

// sign up user
var data = make(map[string]interface{})
data["username"] = "username"
data["password"] = "password"
sessionToken, _ := u.SignUp(data)

// if user records are protected
// sign up user with session token in the body
var data = make(map[string]interface{})
data["username"] = "username"
data["password"] = "password"
data["sessionToken"] = "sessionToken"
sessionToken, _ := u.SignUp(data)

// current user
user, err := u.CurrentUser("sessionToken")

// request password reset
err := u.RequestPasswordReset("email")

// request email verification
err := u.VerificationEmailRequest("email")
```

### Object

Construct a new object, then use the methods on the object to
create, update, delete and query objects. For example:

```go
o := object.NewObject("applicationId", "restApiKey", "sessionToken")

// create object
var data = make(map[string]interface{})
data["name"] = "name"
data["description"] = "description"
object, err := o.Create("className", data)

// update object
var data = make(map[string]interface{})
data["name"] = "name"
data["description"] = "description"
isUpdated, err := o.Update("className", "objectId", data)

// delete object
isDeleted, err := o.Delete("className", "objectId")

// read object
object, err := o.Read("className", "objectId")

// list objects
objects, err := o.List("className")

// count objects
objects, err := o.List("className", WithCount(10))

// limit objects
objects, err := o.List("className", WithLimit(10))

// skip objects
objects, err := o.List("className", WithSkip(10))

// order objects
objects, err := o.List("className", WithOrder("name"))

// objects with constraints
objects, err := o.List("className", WithConstraints("{"title": "My post title", "likes": { "$gt": 100 }}"))
```

### Utility functions

The util package contains some useful functions. For example:

```go
// generate a back4app date object
date := utility.ToBack4AppDate('2020-01-01T00:00:00.000Z')

// parse a back4app date object from map[string]interface{}
date := utility.ParseBack4AppDate(map[string]interface{}{"__type": "Date", "iso": "2020-01-01T00:00:00.000Z"})

// convert a back4app date object to a time.Time object
time, err := utility.Back4AppDateToTime(date)

// convert a back4app date object to a string
str := utility.Back4AppDateToIsoString(date)
```

## License

This project is licensed under the MIT License - see the [`LICENSE`](LICENSE) file for details.

## Contributing

Any kind of positive contribution is welcome! Please help us to grow by contributing to the project.

If you wish to contribute, you can work on any features you think would enhance the library. After adding your code, please send us a Pull Request.

> Please read [CONTRIBUTING](CONTRIBUTING.md) for details on our [CODE OF CONDUCT](CODE_OF_CONDUCT.md), and the process for submitting pull requests to us.

## Support

We all need support and motivation. Please give this project a ⭐️ to encourage and show that you liked it. Don't forget to leave a star ⭐️ before you move away.
