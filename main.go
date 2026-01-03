package main

import (
	"fmt"
	"unsafe"
)

type MyBuilder struct {
	buffer []byte
}

// NewMyBuilder Конструктор
func NewMyBuilder() MyBuilder { return MyBuilder{} }

// Write, который добавляет данные в буфер и возвращаем кол-во записанных байт
func (b *MyBuilder) Write(p []byte) (n int) {
	// Просто добавляем последовательность байт
	b.buffer = append(b.buffer, p...)
	return len(b.buffer)
}

// String, который возвращает итоговую неизменяемую строку.
func (b *MyBuilder) String() string {
	// Создаем новый header с длинной на исходый массив,
	// чтобы избежать копирования.
	return unsafe.String(unsafe.SliceData(b.buffer), len(b.buffer))
}

// Grow расширяет буффер
func (b *MyBuilder) Grow(n int) {
	if n < 0 {
		panic("negative buffer")
	}
	// Проверяем есть ли свободное место
	if cap(b.buffer)-len(b.buffer) < n {
		b.grow(n)
	}
}

func (b *MyBuilder) grow(n int) {
	// Создаем новый слайс
	// Можно было бы и так, но мы будем занулять память
	//buf := bytealg.MakeNoZero(2*cap(b.buf) + n)[:len(b.buf)]
	buf := make([]byte, len(b.buffer), 2*cap(b.buffer)+n)
	// Копируем данные из старого
	copy(buf, b.buffer)
	// Переопределяем
	b.buffer = buf
}

func main() {

	myBuilder := NewMyBuilder()
	myBuilder.Write(([]byte)("Hello"))
	myBuilder.Write(([]byte)(" World"))
	myBuilder.Write(([]byte)("!!!"))
	myBuilder.Write(([]byte)("!!!"))
	myBuilder.Write(([]byte)("!!!"))
	myBuilder.Write(([]byte)("!!!"))
	myBuilder.Write(([]byte)("!!!"))
	myBuilder.Write(([]byte)("!!!"))
	result := myBuilder.String()
	fmt.Println(result)
	fmt.Println(len(myBuilder.buffer))
	myBuilder.Grow(len(result))
	myBuilder.Write(([]byte)("Hello"))
	result = myBuilder.String()
	fmt.Println(result)
	fmt.Println(len(myBuilder.buffer))
}

// Hello World!!!!!!!!!!!!!!!!!!
// 29
// Hello World!!!!!!!!!!!!!!!!!!Hello
// 34
