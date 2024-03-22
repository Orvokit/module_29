/*
Реализуйте паттерн-конвейер:

	1. Программа принимает числа из стандартного ввода в бесконечном цикле и передаёт число в горутину.
	2. Квадрат: горутина высчитывает квадрат этого числа и передаёт в следующую горутину.
	3. Произведение: следующая горутина умножает квадрат числа на 2.
	4. При вводе «стоп» выполнение программы останавливается.

Воспользуйтесь небуферизированными каналами и waitgroup.
*/

package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

func main() {

	var wg sync.WaitGroup

	ic := input(&wg)
	sc := square(ic, &wg)
	mc := multiply(sc, &wg)
	fmt.Println("умножение квадрата на 2:", <-mc)
	wg.Wait()
}

func input(wg *sync.WaitGroup) chan int {
	inputChan := make(chan int)

	go func() {
		defer func() {
			wg.Done()
			fmt.Println("Закрываем канал input")
			close(inputChan)
		}()

		var x string
		for {
			fmt.Print("Введите число, если хотите завершить программу, введите 'стоп': ")

			_, err := fmt.Scan(&x)
			if err != nil {
				log.Println(err)
				continue
			}

			xInt, err := strconv.Atoi(x)
			if err != nil {
				if x == "стоп" {
					break
				}
				log.Println(err, "ошибка конвертации строки в число")
				continue
			}
			fmt.Printf("Кладем значение %v в канал inputChan\n", xInt)
			inputChan <- xInt
			time.Sleep(3 * time.Second) //нужно чтобы новое значение запрашивалось после отработки двух следующих горутин
		}
	}()
	return inputChan
}

func square(inputChan chan int, wg *sync.WaitGroup) chan int {
	squareChan := make(chan int)
	wg.Add(1)

	go func() {
		defer func() {
			wg.Done()
			fmt.Println("Закрываем канал square")
		}()

		for val := range inputChan {
			fmt.Printf("Забираем значение %v из канада input и кладем значение %v в square\n", val, val*val)
			squareChan <- val * val
		}
	}()
	return squareChan
}

func multiply(squareChan chan int, wg *sync.WaitGroup) chan int {
	multiplyChan := make(chan int)
	wg.Add(1)

	go func() {
		defer func() {
			wg.Done()
			fmt.Println("Закрываем канал multiply")
		}()

		for val := range squareChan {
			fmt.Println("возведение в квадрат:", val)
			fmt.Printf("Забираем значение %v из канала squarechan и кладем значение %v в multiplyChan\n", val, val*2)
			multiplyChan <- val * 2
		}
	}()
	return multiplyChan
}
