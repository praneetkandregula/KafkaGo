# KafkaGo
Using Kafka to connect microservice and keeping them highly available.
Even if they go down (for example) when they come back up they can start reading from where they stopped.

## Usage
* Start kafka and zookeeper with ```docker-compose up -d```
* Start the producer by running ```go run main.go``` in the root directory
* Start consumers:
  * In ./datateam run ```go run main.go```
  * In ./processor run ```go run main.go```
* Try stopping either/both consumers
* Bring them back up with the above commands, they'll start reading from the last read value and get all the missed values from the producer
