package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", time.Second*10, "timeout of closing")
}

func main() {
	flag.Parse()

	if !hasCorrectArguments(os.Args) {
		log.Fatalln("Required arguments has missing (host or port)")
	}

	ctxNotify, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

	ctxTimeout, cancel := context.WithTimeout(ctxNotify, timeout)
	defer cancel()

	// Получаем из аргументов адрес хоста и порт.
	countArgs := len(os.Args)
	address := net.JoinHostPort(os.Args[countArgs-2], os.Args[countArgs-1])
	telnet := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	defer func() {
		err := telnet.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	err := telnet.Connect()
	if err != nil {
		log.Println(err)
		return
	}

	// Запускаем горутину для отправки данных в сокет.
	go func() {
		for {
			select {
			case <-ctxTimeout.Done():
				return
			default:
				err = telnet.Send()
				if hasBeenError(err) {
					log.Printf("Failed send data to socket, %v", err)
				}

				stopGoroutinesIfEndOfFile(err, stop)
			}
		}
	}()

	// Запускаем горутину для получения данных из сокета.
	go func() {
		for {
			select {
			case <-ctxTimeout.Done():
				return
			default:
				err = telnet.Receive()
				if hasBeenError(err) {
					log.Printf("Failed receive data from socket, %v", err)
				}

				stopGoroutinesIfEndOfFile(err, stop)
			}
		}
	}()

	<-ctxTimeout.Done()
	err = ctxTimeout.Err()
	if hasBeenError(err) {
		log.Println(err)
	}
}
