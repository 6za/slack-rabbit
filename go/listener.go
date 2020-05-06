package main
import (
	"flag"
	"log"
	"net/http"
	"github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
  "fmt"
  "os/exec"
  "os"
  "strings"
  "app/slackevents"
  //"io/ioutil"
  "github.com/gin-gonic/gin"
  "github.com/streadway/amqp"
)
import "encoding/json"
//"io/ioutil"
import _ "time"

func initPrometheus(){
	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var bot_id = os.Getenv("BOT_USER_ID")
var queue = os.Getenv("OUT_QUEUE")
var queueHostname = os.Getenv("QUEUE_HOSTNAME")
var queueUser = os.Getenv("QUEUE_USER")
var queuePassword = os.Getenv("QUEUE_PASSWORD")

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.(Prometheus)")
var (
	sampleCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "listener_sample", Help: "Sample echo listener"	})
)

func init() {
	// Register the summary and the histogram with Prometheus's default registry.
	prometheus.MustRegister(sampleCounter)
}

type Payload struct {
  Data string `json:"data"  binding:"exists"`
}

type Sender struct {
  Platform       string  `json:"platform"`
  SenderId         int64   `json:"senderId"`
} 

func messageIn(c *gin.Context) {
  //- Setup 
    sampleCounter.Inc()    
    conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:5672/", queueUser,queuePassword,queueHostname))
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close()

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()
    failOnError(err, "Failed to declare a queue")
    q, err := ch.QueueDeclare(
      queue, // name
      false,   // durable
      false,   // delete when unused
      false,   // exclusive
      false,   // no-wait
      nil,     // arguments
    )
    failOnError(err, "Failed to declare a queue")



  //- Payload tranformation
    var data slackevents.SlackEvent
    c.BindJSON(&data)
    fmt.Printf("Data to store: %v\n", data)
    //log.Printf(" [jsonPackage] Sent x %v", string(data)) 
    //Trying to remove echo
    messageDebug, err := json.Marshal(data)
    fmt.Printf("-->> Bot ID/UserID: %v - %v - %s\n", bot_id, data.Event.User, string(messageDebug) )
    if data.Event.User != bot_id && data.Event.User != "" {

      message := slackevents.Message{User: data.Event.User, 
        Text: data.Event.Text,
        Source: "slack"}
      message.ReplyTo.Channel =  data.Event.Channel
        //payload transformation    
        //- Publish
      jsonPackage, _ := json.Marshal(message)
      err = ch.Publish(
        "",     // exchange
        q.Name, // routing key
        false,  // mandatory
        false,  // immediate
        amqp.Publishing{
          ContentType: "application/json",
          Body:        jsonPackage,
        })
        log.Printf(" [jsonPackage] Sent x %v", string(jsonPackage))    
      failOnError(err, "Failed to publish a message")
    }  else {
      log.Printf(" [jsonPackage-data] Echo from %v", bot_id)   
    } 
  //- Response
    buf := make([]byte, 1024)
    num, _ := c.Request.Body.Read(buf)
    reqBody := string(buf[0:num])
    c.JSON(http.StatusOK, reqBody)

}


//for slack hand-shake
func slack(c *gin.Context) {
  buf := make([]byte, 1024)
	num, _ := c.Request.Body.Read(buf)
  reqBody := string(buf[0:num])
  log.Printf(" [jsonPackage-data] Sent x %v", string(reqBody))   
  c.JSON(http.StatusOK, reqBody)
}

func main() {
	 go initPrometheus()
   r := gin.Default()
   r.POST("/message", messageIn)
   r.POST("/slack", messageIn)
   r.Run(":9090")
}


//- Support
  func checkError(e error){
      if e != nil {
          panic(e)
      }
  }
  func printCommand(cmd *exec.Cmd) {
    fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
  }
  func printError(err error) {
    if err != nil {
      os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
    }
  }
  func printOutput(outs []byte) {
    if len(outs) > 0 {
      fmt.Printf("==> Output: %s\n", string(outs))
    }
  }
  func failOnError(err error, msg string) {
    if err != nil {
      log.Fatalf("%s: %s", msg, err)
    }
  }