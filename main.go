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
	// Создаем новуй header с длинной на исходый массив,
	// чтобы избежать копирования.
	return unsafe.String(unsafe.SliceData(b.buffer), len(b.buffer))
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
}
