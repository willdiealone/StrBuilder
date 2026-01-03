package main

import (
	"fmt"
	"unsafe"
)

type MyStringBuilder struct {
	buffer []byte
}

// NewMyStringBuilder Конструктор
func NewMyStringBuilder() MyStringBuilder { return MyStringBuilder{} }

// Write добавляет данные в буфер и возвращаем кол-во записанных байт
func (b *MyStringBuilder) Write(p []byte) (n int) {
	b.buffer = append(b.buffer, p...)
	return len(b.buffer)
}

// String возвращает итоговую неизменяемую строку.
func (b *MyStringBuilder) String() string {
	// Создаем новый header с длинной на исходый массив,
	// чтобы избежать копирования.
	return unsafe.String(unsafe.SliceData(b.buffer), len(b.buffer))
}

// Grow расширяет буффер
func (b *MyStringBuilder) Grow(n int) {
	if n < 0 {
		panic("negative buffer")
	}
	// Проверяем есть ли свободное место
	if cap(b.buffer)-len(b.buffer) < n {
		b.grow(n)
	}
}

// grow расширяет буффер
func (b *MyStringBuilder) grow(n int) {
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

	myBuilder := NewMyStringBuilder()
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
