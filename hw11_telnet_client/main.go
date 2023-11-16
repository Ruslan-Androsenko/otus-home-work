package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", time.Second*10, "timeout of closing")
}

func main() {
	flag.Parse()

	if len(os.Args) < 3 {
		log.Fatalln("Required arguments has missing (host or port)")
	}

	ctxNotify, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

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
			case <-ctxNotify.Done():
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
			case <-ctxNotify.Done():
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

	<-ctxNotify.Done()
	errNotify := ctxNotify.Err()
	if hasBeenError(errNotify) {
		log.Println(errNotify)
	}
}

// Проверяем что это пришла ошибка а не конец файла.
func hasBeenError(err error) bool {
	return err != nil && !errors.Is(err, io.EOF)
}

// Остановить горутины если ошибка является концом файла.
func stopGoroutinesIfEndOfFile(err error, stop context.CancelFunc) {
	if errors.Is(err, io.EOF) {
		stop()
	}
}
