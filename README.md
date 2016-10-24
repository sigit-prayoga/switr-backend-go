# switr-backend-go
**[Switr](http://jlp.community/switr)** is a really simple web application that made for a demo of some recent technologies.
It's been dedicated to anyone, community or individual, who wants to share the particular tech.

**Frontend Tech Stack:**
* [AngularJS](https://github.com/sigit-prayoga/switr-web)
* [Angular 2.0](https://github.com/sigit-prayoga/switr-web-angular2)
* ReactJS
* Bootstrap
* Socket IO Client
* Ionic in iOS and Android (coming soon)

**Backend Tech Stack:**
* [NodeJS](https://github.com/sigit-prayoga/switr-backend)
* GO

**Database**
* MongoDB
* PostgreSQL (coming soon)

**This Repo:**
* GO
* MongoDB

[Go](https://golang.org/) is an open source programming language that makes it easy to build simple, reliable, and efficient software.
[MongoDB](https://www.mongodb.com/) Enterprise Advanced features MongoDB Enterprise Server and a finely-tuned package of advanced software, support, certifications, and other services. More than one-third of the Fortune 100 rely on MongoDB Enterprise Advanced to help run their mission critical applications.

***
[See GO installation here](https://golang.org/doc/install)

### Package to import
* [Goji](https://goji.io/) is a HTTP request multiplexer
```sh
$ go get goji.io
```

* [mgo](https://labix.org/mgo) is a MongoDB driver for Go
```sh
$ go get gopkg.in/mgo.v2
```

* [cors](go get github.com/rs/cors) Go net/http configurable handler to handle CORS requests
```sh
$ go get github.com/rs/cors
```

### Navigate to the project root folder
```sh
$ cd [projectRoot]
```

### Run GO main app
```sh
$ go run main.go
```
> *I should quit and rerun whenever create changes in a code.*
>> **In this case, you might need [```gin```, the autoreload](https://github.com/codegangsta/gin)**

### To install ```gin```
```sh
$ go get github.com/codegangsta/gin
```

### Run ```gin``` in your project root
```sh
$ gin
```

> Use ```-i``` flag to immediately run after build succeed.

***

## **Other version of the backend**
In case you wonder to learn this backend using [NodeJS](http://nodejs.org), you can get the Node version [here.](https://github.com/sigit-prayoga/switr-backend)
Also, we have already had the frontend part using AngularJS 1.5, take 'em out [here](https://github.com/sigit-prayoga/switr-web)

> *'I like Angular 2.0 more than 1.5'*
>> Alright, alright! [check this out!](https://github.com/sigit-prayoga/switr-web-angular2)

***