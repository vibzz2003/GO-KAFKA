package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/gofiber/fiber/v2"
)

type Comment struct {
	Text string `form:"text" json:"text"`
}

//fiber is a express-inspired web framework built on top of FastHttp
//creating a fiber application is the first step of setting up a web server
//used to define routes, middle ware and other things for the project

func main() {
	app := fiber.New()                   //creates a new fiber application instance
	api := app.Group("/api/v1")          //creates a group of routes that have the prefix /api/v1
	api.Post("/comments", createComment) //creates a post route /comments within the group and makes createComment function to handle it
	app.Listen(":3000")                  //web server listens at port 3000
}

func ConnectProducer(brokersUrl []string) (sarama.SyncProducer, error) { //connects to kafka as producer
	config := sarama.NewConfig()                     //create a new configuration
	config.Producer.Return.Successes = true          //sets such that the producer waits for acknowledgement
	config.Producer.RequiredAcks = sarama.WaitForAll //producer waits for all replicas to acknowledge
	config.Producer.Retry.Max = 5                    //setting the maximum number of retries

	conn, err := sarama.NewAsyncProducer(brokersUrl, config) //creates a new ansynchronous producer
	if err != nil {
		return nil, err
	}
	return conn, nil //returns connection if all ok or else returns error
}

func PushCommentToQueue(topic string, message []byte) error { //function for sending message to kafka
	brokersUrl := []string{"localhost:29092"}    //specify kafka brokers
	producer, err := ConnectProducer(brokersUrl) //connects to kafka
	if err != nil {
		return err //if error then returns error else proceeds
	}
	defer producer.Close()          //ensures connection to kafka closes after use
	msg := &sarama.ProducerMessage{ //creates a new producer message
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	} // if theres an error, returns error else proceeds with the topic and message

	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}

func createComment(c *fiber.Ctx) error { //handling the createComment funtion using fiber context as an argument and returns an error
	cmt := new(Comment)                       //creates a new instance of the type struct
	if err := c.BodyParser(cmt); err != nil { //if the error is not equal to null
		log.Println(err)               //priting the error
		c.Status(400).JSON(&fiber.Map{ //letting the error have status off 400 anf mapping the values in a json format for the key as success and the message as the value
			"success": false,
			"message": err,
		})
		return err //returning the error
	}
	cmtInBytes, err := json.Marshal(cmt)       //converst comments into json bytes
	PushCommentToQueue("comments", cmtInBytes) //pushes the comment to kafka queue

	c.JSON(&fiber.Map{ //sens a positive response
		"success": true,
		"message": "Comment pushes successfully",
		"comment": cmt,
	})
	if err != nil { //if error then sends the negative response
		c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "Error creating product",
		})
		return err
	}
	return err //returns error if something goes wrong
}
