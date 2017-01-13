# gonfig
a config parser for GO with support for path selectors

##package documentation

https://godoc.org/github.com/creamdog/gonfig

##overview

a loaded config gonfig.FromJson("path/to/config.json") exposes the following methods

- GetString("path/to/key", "default value / can be nil") (string, error)
- GetInt("path/to/key", 0) (string, error)
- GetFloat("path/to/key", nil) (string, error)
- GetBool("path/to/key", false) (string, error)
- GetAs("path/to/key", &myCustomStruct) error

##usage example

###example config file

```json
{
  "services" : {
    "facebook" : {
      "host" : "facebook.com",
      "port" : 443,
      "message" : {
        "template" : {
          "subject" : "this is my default subject",
          "message" : "whoops, looks like I forgot to write a body!"
        }
      }
    }
  },
  "credentials" : {
    "facebook" : {
      "username" : "jane32",
      "password" : "supersecret"
    }
  }
}
```

###example main.go

```go
import(
  "github.com/creamdog/gonfig"
  "os"
)

type ExampleMessageStruct struct {
  Message string
  Subject string
}

func main() {
  f, err := os.Open("myconfig.json")
  if err != nil {
    // TODO: error handling
  }
  defer f.Close();
  config, err := gonfig.FromJson(f)
  if err != nil {
    // TODO: error handling
  }
  
  
  username, err := config.GetString("credentials/facebook/username", "scooby")
  if err != nil {
    // TODO: error handling
  }
  password, err := config.GetString("credentials/facebook/password", "123456")
  if err != nil {
    // TODO: error handling
  }
  host, err := config.GetString("services/facebook/host", "localhost")
  if err != nil {
    // TODO: error handling
  }
  port, err := config.GetInt("services/facebook/port", 80)
  if err != nil {
    // TODO: error handling
  }
  
  // TODO: example something
  // login(host, port, username, password)
  
  var template ExampleMessageStruct
  if err := config.GetAs("services/facebook/message/template", &template); err != nil {
    // TODO: error handling
  }
  
  /// TODO: example something
  // template.Message = "I just want to say, oh my!"
  // sendMessage(host, port, template)
  
}
```
