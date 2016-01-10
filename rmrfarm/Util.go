package rmrfarm

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"strconv"
	"strings"
)

/*
Binary Reader
*/
type Reader interface {
	ReadByte() byte
	ReadInt8() int8
	ReadInt16() int16
	ReadInt32() int32
	ReadInt64() int64
	ReadFloat() float32
	ReadDouble() float64
	ReadBool() bool
	ReadUtfString() string
	ReadBytesArray() []byte
	GetValue() []byte
}

type binaryreader struct {
	cursor int16
	byte   []byte
}

func CreateBinaryReader(bytearray []byte) Reader {
	bytearraydest := make([]byte, len(bytearray))
	copy(bytearraydest, bytearray)
	return &binaryreader{0, bytearraydest}
}

func (binaryreader *binaryreader) GetValue() []byte{
	t := make([]byte, len(binaryreader.byte))
	copy(t, binaryreader.byte)
	return t
}

func (binaryreader *binaryreader) ReadByte() byte {
	result := byte(binaryreader.byte[binaryreader.cursor])
	binaryreader.cursor += 1
	return byte(result)
}

func (binaryreader *binaryreader) ReadInt8() int8 {
	result := byte(binaryreader.byte[binaryreader.cursor])
	binaryreader.cursor += 1
	return int8(result)
}

func (binaryreader *binaryreader) ReadInt16() int16 {
	result := binary.LittleEndian.Uint16(binaryreader.byte[binaryreader.cursor : binaryreader.cursor+2])
	binaryreader.cursor += 2
	return int16(result)
}

func (binaryreader *binaryreader) ReadInt32() int32 {
	result := binary.LittleEndian.Uint32(binaryreader.byte[binaryreader.cursor : binaryreader.cursor+4])
	binaryreader.cursor += 4
	return int32(result)
}

func (binaryreader *binaryreader) ReadInt64() int64 {
	result := binary.LittleEndian.Uint64(binaryreader.byte[binaryreader.cursor : binaryreader.cursor+8])
	binaryreader.cursor += 8
	return int64(result)
}

func (binaryreader *binaryreader) ReadFloat() (result float32) {
	binary.Read(bytes.NewReader(binaryreader.byte[binaryreader.cursor:binaryreader.cursor+4]), binary.LittleEndian, &result)
	binaryreader.cursor += 4
	return result
}

func (binaryreader *binaryreader) ReadDouble() (result float64) {
	binary.Read(bytes.NewReader(binaryreader.byte[binaryreader.cursor:binaryreader.cursor+8]), binary.LittleEndian, &result)
	binaryreader.cursor += 8
	return result
}

func (binaryreader *binaryreader) ReadBool() bool {
	return binaryreader.ReadInt8() == 1
}

func (binaryreader *binaryreader) ReadLast() []byte {
	length := int(len(binaryreader.byte))
	return binaryreader.ReadBytes(length - int(binaryreader.cursor))
}

func (binaryreader *binaryreader) ReadUtfString() string {
	str := ""
	for {
		char := binaryreader.ReadByte()
		if char == 0 {
			break
		}
		str += string(char)
	}
	return str
}

func (binaryreader *binaryreader) ReadBytes(i int) []byte {
	byte := make([]byte, i)
	for x := 0; x < i; x++ {
		byte[x] = binaryreader.ReadByte()
	}
	return byte
}

func (binaryreader *binaryreader) ReadBytesArray() []byte {
	size := binaryreader.ReadByte()
	bytearray := make([]byte, size)
	for x := byte(0); x < size; x++ {
		bytearray[x] = binaryreader.ReadByte()
	}
	return bytearray
}

/*
Binary Writer
*/
type Writer interface {
	WriteInt8(byte)
	WriteInt16(int16)
	WriteInt32(int32)
	WriteInt64(int64)
	WriteFloat(float32)
	WriteDouble(float64)
	WriteBool(bool)
	WriteUtfString(string)
	WriteByteArray([]byte)
	WriteUtfStringArray([]string)
}

type binarywriter struct {
	byte []byte
}

func CreateBinaryWriter() *binarywriter {
	return &binarywriter{make([]byte, 0)}
}

func (binarywriter *binarywriter) WriteInt8(data byte) {
	binarywriter.byte = append(binarywriter.byte, data)
}

func (binarywriter *binarywriter) WriteInt16(data int16) {
	bytearray := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytearray, uint16(data))
	binarywriter.byte = append(binarywriter.byte, bytearray...)
}

func (binarywriter *binarywriter) WriteInt32(data int32) {
	bytearray := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytearray, uint32(data))
	binarywriter.byte = append(binarywriter.byte, bytearray...)
}

func (binarywriter *binarywriter) WriteInt64(data int64) {
	bytearray := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytearray, uint64(data))
	binarywriter.byte = append(binarywriter.byte, bytearray...)
}

func (binarywriter *binarywriter) WriteFloat(data float32) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, data)
	binarywriter.byte = append(binarywriter.byte, buf.Bytes()...)
}

func (binarywriter *binarywriter) WriteDouble(data float64) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, data)
	binarywriter.byte = append(binarywriter.byte, buf.Bytes()...)
}

func (binarywriter *binarywriter) WriteBool(data bool) {
	if data{
		binarywriter.WriteInt8(1)
	}else{
		binarywriter.WriteInt8(0)
	}
}

func (binarywriter *binarywriter) WriteUtfString(data string) {

	binarywriter.byte = append(binarywriter.byte, []byte(data)...)
	binarywriter.WriteInt8(byte(0))

}

func (binarywriter *binarywriter) WriteByteArray(data []byte) {
	binarywriter.WriteInt32(int32(len(data)))
	for _, b := range data {
		binarywriter.WriteInt8(b)
	}
}

func (binarywriter *binarywriter) WriteUtfStringArray(data []string) {
	binarywriter.WriteInt32(int32(len(data)))
	for _, s := range data {
		binarywriter.WriteUtfString(s)
	}
}

func (binarywriter *binarywriter) Bytes() []byte {
	return binarywriter.byte
}

/*
Util Function
*/
func GenerateToken(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func NewIdPool() func() int32 {
	i := int32(0)
	return func() int32 {
		i += 1
		return i
	}
}

func ToStrings(a ...interface{}) string {
	stringList := make([]string, len(a))
	for argNum := 0; argNum < len(a); argNum++ {
		arg := a[argNum]
		switch f := arg.(type) {
		case bool:
			if f {
				stringList[argNum] = "true"
			} else {
				stringList[argNum] = "false"
			}
		case float32:
			stringList[argNum] = strconv.FormatFloat(float64(f), byte('g'), 1, 64)
		case float64:
			stringList[argNum] = strconv.FormatFloat(float64(f), byte('g'), 1, 64)
		case complex64:
			stringList[argNum] = "no complex Parser"
		case complex128:
			stringList[argNum] = "no complex Parser"
		case int:
			stringList[argNum] = strconv.FormatInt(int64(f), 10)
		case int8:
			stringList[argNum] = strconv.FormatInt(int64(f), 10)
		case int16:
			stringList[argNum] = strconv.FormatInt(int64(f), 10)
		case int32:
			stringList[argNum] = strconv.FormatInt(int64(f), 10)
		case int64:
			stringList[argNum] = strconv.FormatInt(int64(f), 10)
		case uint:
			stringList[argNum] = strconv.FormatUint(uint64(f), 10)
		case uint8:
			stringList[argNum] = strconv.FormatUint(uint64(f), 10)
		case uint16:
			stringList[argNum] = strconv.FormatUint(uint64(f), 10)
		case uint32:
			stringList[argNum] = strconv.FormatUint(uint64(f), 10)
		case uint64:
			stringList[argNum] = strconv.FormatUint(uint64(f), 10)
		case uintptr:
			stringList[argNum] = strconv.FormatUint(uint64(f), 10)
		case string:
			stringList[argNum] = f
		case []byte:
			stringList[argNum] = "no byte Parser"
		default:
			stringList[argNum] = "no default Parser"
		}
	}
	return strings.Join(stringList, " ")
}
