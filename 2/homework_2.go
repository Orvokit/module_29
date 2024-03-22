/*
В работе часто возникает потребность правильно останавливать приложения.
Например, когда наш сервер обслуживает соединения, а нам хочется, чтобы все текущие соединения были обработаны
и лишь потом произошло выключение сервиса. Для этого существует паттерн graceful shutdown.

Напишите приложение, которое выводит квадраты натуральных чисел на экран, а после получения сигнала ^С обрабатывает этот сигнал,
пишет «выхожу из программы» и выходит.

Для реализации данного паттерна воспользуйтесь каналами и оператором select с default-кейсом
*/

package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	var result int

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	wg.Add(1)

	go func() {
		defer wg.Done()

		for i := 0; i < 100; i++ {
			select {
			case <-done:
				fmt.Println("выхожу из программы")
				return
			default:
				result = i * i
				fmt.Println(result)
				time.Sleep(time.Second)
			}
		}
	}()
	wg.Wait()
}
