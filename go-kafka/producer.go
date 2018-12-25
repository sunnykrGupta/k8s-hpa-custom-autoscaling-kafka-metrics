// function-based Apache Kafka producer
package main

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"syscall"
	"os/signal"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/Pallinder/go-randomdata"
	"encoding/json"
)


//struct to keep MSG
type ProfilePara struct {
    Bio *randomdata.Profile
    P1 string
	P2 string
	P3 string
	P4 string
	P5 string
}

func profileProvider() ([]byte, error) {

	profilePara := ProfilePara{
			Bio: randomdata.GenerateProfile(randomdata.RandomGender),
			P1: randomdata.Paragraph() + randomdata.Paragraph() + randomdata.Paragraph() + randomdata.Paragraph(),
			P2: randomdata.Paragraph() + randomdata.Paragraph() + randomdata.Paragraph() + randomdata.Paragraph(),
			P3: randomdata.Paragraph() + randomdata.Paragraph() + randomdata.Paragraph() + randomdata.Paragraph(),
			P4: randomdata.Paragraph() + randomdata.Paragraph() + randomdata.Paragraph() + randomdata.Paragraph(),
			P5: randomdata.Paragraph() + randomdata.Paragraph() + randomdata.Paragraph() + randomdata.Paragraph(),
	}
	jsonBytes, err := json.Marshal(profilePara)

	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

//----------------

func routineStats(msgBurst int, statsDone chan bool, counterChannel chan int,
                timelast_Channel chan time.Time, timenow_Channel chan time.Time ){
    //infinite loop for channel read wait
	run := true
    for run == true {

		select {
		case <- statsDone:
			run = false
			fmt.Println("Closed Stats Goroutine")
			break

		case t_last := <- timelast_Channel:
	        t_now := <- timenow_Channel
	        counter := <- counterChannel
            // result will be duration
            t_diff := t_now.Sub(t_last)
            fmt.Println(t_last.Format("15:04:05.999"), t_diff.Seconds()*1000)
            fmt.Println(float64(msgBurst)/t_diff.Seconds(), " RPS Rate | ", "Produced :", counter)
            fmt.Println()
            //
		}
    }
}



func main() {

	if len(os.Args) < 9 {
		fmt.Fprintf(os.Stderr, "Usage: %s <broker:port>  <topic>  <msgBurst>  <total>  <none/gzip>  <lingerMs>  <waitMs>  <ssl/plaintext>  ...<ssl params-optional> \n",
			os.Args[0])
		os.Exit(1)
	}

	broker := os.Args[1]
	topic := os.Args[2]
	msgBurst, _ := strconv.Atoi(os.Args[3])
	total, _ := strconv.Atoi(os.Args[4])
	compression_codec := os.Args[5]
	lingerMs, _ := strconv.Atoi(os.Args[6])
	waitMs, _ := strconv.Atoi(os.Args[7])
	securityProtocol := os.Args[8]


	//linger.ms : Delay in milliseconds to wait for messages in the producer queue to accumulate before
	//"message.max.bytes": 1000,000, Maximum Kafka protocol request message(batch) size.
	//"queue.buffering.max.kbytes" 1048576, (1GB) Maximum total message size sum allowed on the producer queue.

	kafkaCMap := &kafka.ConfigMap{"bootstrap.servers": broker,
								"request.required.acks": 1,
								"queue.buffering.max.kbytes": 204800,
								"queue.buffering.max.messages": 100000,
								"linger.ms": lingerMs,
								"security.protocol": securityProtocol,
								"compression.codec": compression_codec,
	}

	if securityProtocol == "ssl" {
		if len(os.Args) != 12 {
			fmt.Fprintf(os.Stderr, "Usage: %s <broker:port> <topic> <msgBurst> <total> <none/gzip> <lingerMs> <waitMs> <ssl/plaintext/sasl_plaintext/sasl_ssl> \n Additional Params needed : <ssl.ca.location> <ssl.keystore.location> <ssl.keystore.password> \n",
				os.Args[0])
			os.Exit(1)
		}

		ca_location := os.Args[9]
		keystore_location := os.Args[10]
		keystore_password := os.Args[11]

		kafkaCMap.SetKey("ssl.ca.location", ca_location)
		kafkaCMap.SetKey("ssl.keystore.location", keystore_location)
		kafkaCMap.SetKey("ssl.keystore.password", keystore_password)

	}



	p, err := kafka.NewProducer(kafkaCMap)

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Producer %v\n", p)

	//Channel declaration
	statsDone := make(chan bool)
    counterChannel := make(chan int)
    timelast_Channel := make(chan time.Time)
    timenow_Channel := make(chan time.Time)

	// goroutine invoke which process stats
	go routineStats(msgBurst, statsDone, counterChannel, timelast_Channel, timenow_Channel)

	go func(){
		//defer close(doneChan)

		for e := range p.Events(){
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				}

			default :
				fmt.Printf("Ignored event: %s\n", ev)
			}
		}
	}()


	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	msgCnt := 0
	run := true
	t_last := time.Now()
	jsonBytes, _ := profileProvider()

	fmt.Println("\n\t Total packets to produce : %d", total)

	// To run indefinitely
	for run == true {

		select {
			case sig := <- sigchan:
				fmt.Printf("Caught Signal %v: terminating\n", sig)
				fmt.Printf("\n\tProduced : %d\n", msgCnt)
				run = false
				// terminate goroutine stats
				statsDone <- true

			default:
				//AVOID randomdata
				//jsonBytes, err := profileProvider()
				//fmt.Printf("%% %s\n", jsonBytes[:70])

				err = p.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
					Value:          []byte(jsonBytes),
				}, nil)

				if err == nil {
					msgCnt++
				}else {
					fmt.Println("Produce Err : ", err)
					fmt.Printf("Wait for delivery : %d ms\n\n", waitMs)
					time.Sleep(time.Millisecond * time.Duration(waitMs))
				}

				//Periodically invoke stats crunchers
				if msgCnt % msgBurst == 0  && msgCnt != 0 && run != false {
		            t_now :=  time.Now()
					//Write to channels
		            timelast_Channel <- t_last
		            timenow_Channel <- t_now
		            counterChannel <- msgCnt
		            t_last = t_now

		            // added sleep to avoid heavy disk filling in broker
		            time.Sleep(time.Millisecond * time.Duration(waitMs))
				}
		}

		if run == false{
			break
		}
	}


	n := p.Flush(5000)
	fmt.Println("Outstanding events unflushed : ", n)

	p.Close()
}
