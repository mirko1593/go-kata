package main

import (
	"bytes"
	"fmt"
)

// data represents a table of input and expected output.
var data = []struct {
	input  []byte
	output []byte
}{
	{[]byte("abc"), []byte("abc")},
	{[]byte("elvis"), []byte("Elvis")},
	{[]byte("aElvis"), []byte("aElvis")},
	{[]byte("abcelvis"), []byte("abcElvis")},
	{[]byte("eelvis"), []byte("eElvis")},
	{[]byte("aelvis"), []byte("aElvis")},
	{[]byte("aabeeeelvis"), []byte("aabeeeElvis")},
	{[]byte("e l v i s"), []byte("e l v i s")},
	{[]byte("aa bb e l v i saa"), []byte("aa bb e l v i saa")},
	{[]byte(" elvi s"), []byte(" elvi s")},
	{[]byte("elvielvis"), []byte("elviElvis")},
	{[]byte("elvielvielviselvi1"), []byte("elvielviElviselvi1")},
	{[]byte("elvielviselvis"), []byte("elviElvisElvis")},
}

func assembleInputStream() []byte {
	var in []byte
	for _, d := range data {
		in = append(in, d.input...)
	}

	return in
}

func assembleOutputStream() []byte {
	var out []byte
	for _, d := range data {
		out = append(out, d.output...)
	}

	return out
}

func main() {
	var output bytes.Buffer
	in := assembleInputStream()
	out := assembleOutputStream()

	find := []byte("elvis")
	repl := []byte("Elvis")

	fmt.Println("=======================================\nRunning Algorithm One")
	output.Reset()
	algOne(in, find, repl, &output)
	matched := bytes.Compare(out, output.Bytes())
	fmt.Printf("Matched: %v\nInp: [%s]\nExp: [%s]\nGot: [%s]\n", matched == 0, in, out, output.Bytes())

	fmt.Println("=======================================\nRunning Algorithm Two")
	output.Reset()
	algTwo(in, find, repl, &output)
	matched = bytes.Compare(out, output.Bytes())
	fmt.Printf("Matched: %v\nInp: [%s]\nExp: [%s]\nGot: [%s]\n", matched == 0, in, out, output.Bytes())
}

func algOne(data []byte, find []byte, repl []byte, output *bytes.Buffer) {
	input := bytes.NewBuffer(data)
	// input := &bytes.Buffer{buf: data}

	size := len(find)

	buf := make([]byte, 5)
	end := size - 1

	// read initial chunk
	if n, err := input.Read(buf[:end]); err != nil {
		// if n, err := io.ReadFull(input, buf[:end]); err != nil {
		output.Write(buf[:n])
		return
	}

	for {
		// slice one byte
		// if _, err := io.ReadFull(input, buf[end:]); err != nil {
		// 	output.Write(buf[:end])
		// 	return
		// }
		var err error
		buf[end:][0], err = input.ReadByte()
		if err != nil {
			output.Write(buf[:end])
			return
		}

		// replace if match
		if bytes.Equal(buf, find) {
			output.Write(repl)

			// read new initial chunk
			// if n, err := io.ReadFull(input, buf[:end]); err != nil {
			if n, err := input.Read(buf[:end]); err != nil {
				output.Write(buf[:n])
				return
			}

			continue
		}

		output.WriteByte(buf[0])

		// splice buf by one byte
		copy(buf, buf[1:])
	}
}

func algTwo(data []byte, find []byte, repl []byte, output *bytes.Buffer) {
	input := bytes.NewReader(data)

	size := len(find)

	idx := 0

	for {
		b, err := input.ReadByte()
		if err != nil {
			break
		}

		if b == find[idx] {
			idx++

			if idx == size {
				output.Write(repl)
				idx = 0
			}

			continue
		}

		if idx != 0 {
			output.Write(find[:idx])

			input.UnreadByte()

			idx = 0

			continue
		}

		output.WriteByte(b)
		idx = 0
	}
}
