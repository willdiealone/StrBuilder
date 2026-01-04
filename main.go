package main

import (
	"fmt"
	"unicode/utf8"
	"unsafe"
)

type MyStringBuilder struct {
	addr   *MyStringBuilder
	buffer []byte
}

// NewMyStringBuilder Конструктор
func NewMyStringBuilder() MyStringBuilder { return MyStringBuilder{} }

// copyCheck разделяемый буфер: Если вы скопируете Builder,
// оба экземпляра (старый и новый) будут указывать на один и тот же массив байт в памяти
// (так как слайс внутри — это просто указатель на массив, длина и емкость).
// Если вы вызвали .String() у первого билдера, вы получили строку,
// которая «смотрит» на его внутренний буфер. Если вы затем скопируете этот билдер во второй и
// продолжите писать во второй, он может изменить данные в том же самом буфере.
// Значение «неизменяемой» строки в первом билдере внезапно изменится.
// Это приведет к непредсказуемому поведению программы и нарушению безопасности памяти.
func (b *MyStringBuilder) copyCheck() {
	if b.addr == nil {
		// Скрываем указатель от escape анализа, предотвращая попадание в кучу
		// Не разрешено, по этому придется иметь дело с heap
		//b.addr = (*MyStringBuilder)(abi.NoEscape(unsafe.Pointer(b)))
		b.addr = b
	} else if b.addr != b {
		panic("bad address")
	}
}

// Write добавляет данные в буфер и возвращаем кол-во записанных байт
func (b *MyStringBuilder) Write(p []byte) (n int) {
	b.copyCheck()
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
	b.copyCheck()
	// Создаем новый слайс
	// Можно было бы и так, но мы будем занулять память
	//buf := bytealg.MakeNoZero(2*cap(b.buf) + n)[:len(b.buf)]
	buf := make([]byte, len(b.buffer), 2*cap(b.buffer)+n)
	// Копируем данные из старого
	copy(buf, b.buffer)
	// Переопределяем
	b.buffer = buf
}

// Len возващает длину буфера
func (b *MyStringBuilder) Len() int { return len(b.buffer) }

// Cap возващает емкость буфера
func (b *MyStringBuilder) Cap() int { return cap(b.buffer) }

// Reset сбрасываем билдер
func (b *MyStringBuilder) Reset() {
	b.addr = nil
	b.buffer = nil
}

// WriteByte записываем байт в буфер
func (b *MyStringBuilder) WriteByte(c byte) error {
	b.copyCheck()
	b.buffer = append(b.buffer, c)
	return nil
}

// WriteString записываем строку в буфер
func (b *MyStringBuilder) WriteString(s string) (int, error) {
	b.copyCheck()
	b.buffer = append(b.buffer, s...)
	return len(s), nil
}

// WriteRune Записываем руну в буфер
func (b *MyStringBuilder) WriteRune(r rune) (int, error) {
	b.copyCheck()
	n := len(b.buffer)
	b.buffer = utf8.AppendRune(b.buffer, r)
	return len(b.buffer) - n, nil
}

func main() {

	myBuilder := NewMyStringBuilder()
	myBuilder.Write(([]byte)("Hello❤️👍"))
	result := myBuilder.String()
	pointerRes := unsafe.StringData(result)
	slice := unsafe.Slice(pointerRes, len(result))
	fmt.Printf("%08b", slice)
	//     1         1        1       1         1       3                           3                         4
	// [01001000 01100101 01101100 01101100 01101111 11100010 10011101 10100100 11101111 10111000 10001111 11110000 10011111 10010001 10001101]
	// LittleEndian
	fmt.Println(len(myBuilder.buffer))
	myBuilder.Grow(len(result))
	myBuilder.Write(([]byte)("Hello"))
	result = myBuilder.String()
	fmt.Println(result)
	fmt.Println(len(myBuilder.buffer))
	// Hello❤️👍Hello
	// 20
	myBuilder.Reset()
	_, _ = myBuilder.WriteString("Monkey👍")
	result = myBuilder.String()
	fmt.Println(result)
	// Monkey👍

}
