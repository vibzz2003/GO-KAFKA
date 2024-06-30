package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

func main() {
	topic := "comments"                                         //defining the kafka topics to take messages from
	worker, err := connectConsumer([]string{"localhost:29092"}) //connects to kafka brokers
	if err != nil {
		panic(err)
	}

	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest) //starts to consume messages from the specified partition of the topic ( from the oldest partition)
	if err != nil {
		panic(err)
	}

	fmt.Println("consumer started")
	sigchan := make(chan os.Signal, 1)                      //creates a channel to recieve os messages
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM) //register the channel to recieve interuppt and terminate signals

	msgCount := 0                 //initialises the counter number of messages processed
	doneCh := make(chan struct{}) //creates a channel to signal when to stop the consumer

	go func() { //enters the go routine
		for {
			select { //listens on multiple channels
			case err := <-consumer.Errors(): //listens for errors from the consumer
				fmt.Println(err)
			case msg := <-consumer.Messages(): //listems for new messages from the consumer
				msgCount++                                                                                                   //incremrnts the message count
				fmt.Printf("Received message count: %d | Topic(%s) | Message(%s)\n", msgCount, msg.Topic, string(msg.Value)) // prints the recieved message in correct format
			case <-sigchan: //listens for os signals
				fmt.Println("Interruption detected")
				doneCh <- struct{}{}
				return // Exit the goroutine
			}
		}
	}()

	<-doneCh //waits for the signal to stop
	fmt.Println("Processed", msgCount, "messages")
	if err := worker.Close(); err != nil {
		panic(err)
	}
}

func connectConsumer(brokersUrl []string) (sarama.Consumer, error) { //configures and connects the kafka brokers
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true //to enable returning errors
	conn, err := sarama.NewConsumer(brokersUrl, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
