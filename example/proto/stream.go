package proto

// copy from: https://github.com/secondtonone1/wentserver/protocol/stream.go

import (
	"encoding/binary"
	"fmt"
)

var ErrBuffOverflow = fmt.Errorf("剩余的缓存空间不足")

type IReadStream interface {
	Size() int   // 缓存空间大小
	Length() int // 已使用的长度
	Right() int  // 剩余空间大小
	Reset([]byte)
	Data() []byte
	ReadByte() (b byte, err error)
	ReadUint16() (b uint16, err error)
	ReadUint32() (b uint32, err error)
	ReadUint64() (b uint64, err error)
	ReadBuff(size int) (b []byte, err error)
	PeekBuff(size int) (buff []byte, err error)
	CopyBuff(b []byte) error
}

type IWriteStream interface {
	Size() int   // 缓存空间大小
	Length() int // 已使用的长度
	Right() int  // 剩余空间大小
	Reset([]byte)
	Data() []byte
	WriteByte(b byte) error
	WriteUint16(b uint16) error
	WriteUint32(b uint32) error
	WriteUint64(b uint64) error
	WriteBuff(b []byte) error
	CleanBuff() error // 清楚已经读取过的内容以便节省空间
}

type BigEndianStreamImpl struct {
	pos  int
	len  int
	buff []byte
}

func NewBigEndianStream(buff []byte) *BigEndianStreamImpl {
	return &BigEndianStreamImpl{
		buff: buff,
	}
}

func (impl *BigEndianStreamImpl) Size() int   { return len(impl.buff) }
func (impl *BigEndianStreamImpl) Length() int { return impl.len }

func (impl *BigEndianStreamImpl) Data() []byte { return impl.buff }

func (impl *BigEndianStreamImpl) Right() int { return len(impl.buff) - impl.pos }

func (impl *BigEndianStreamImpl) Reset(buff []byte) { impl.pos = 0; impl.len = 0; impl.buff = buff }

func (impl *BigEndianStreamImpl) ReadByte() (b byte, err error) {
	if impl.Right() < 1 {
		return 0, ErrBuffOverflow
	}
	b = impl.buff[impl.pos]
	impl.pos += 1
	return b, nil
}

func (impl *BigEndianStreamImpl) ReadUint16() (b uint16, err error) {
	if impl.Right() < 2 {
		return 0, ErrBuffOverflow
	}
	b = binary.BigEndian.Uint16(impl.buff[impl.pos:])
	impl.pos += 2
	return b, nil
}

func (impl *BigEndianStreamImpl) ReadUint32() (b uint32, err error) {
	if impl.Right() < 4 {
		return 0, ErrBuffOverflow
	}
	b = binary.BigEndian.Uint32(impl.buff[impl.pos:])
	impl.pos += 4
	return b, nil
}

func (impl *BigEndianStreamImpl) ReadUint64() (b uint64, err error) {
	if impl.Right() < 8 {
		return 0, ErrBuffOverflow
	}
	b = binary.BigEndian.Uint64(impl.buff[impl.pos:])
	impl.pos += 8
	return b, nil
}

func (impl *BigEndianStreamImpl) ReadBuff(size int) (buff []byte, err error) {
	if impl.Right() < size {
		return nil, ErrBuffOverflow
	}
	buff = make([]byte, size, size)
	copy(buff, impl.buff[impl.pos:impl.pos+size])
	impl.pos += size
	return buff, nil
}

// PeekBuff 从缓存中预读取指定字节数据不修改pos
func (impl *BigEndianStreamImpl) PeekBuff(size int) (buff []byte, err error) {
	if impl.Right() < size {
		return nil, ErrBuffOverflow
	}
	buff = make([]byte, size, size)
	copy(buff, impl.buff[impl.pos:impl.pos+size])

	return buff, nil
}

func (impl *BigEndianStreamImpl) CopyBuff(b []byte) error {
	if impl.Right() < len(b) {
		return ErrBuffOverflow
	}
	copy(b, impl.buff[impl.pos:impl.pos+len(b)])
	return nil
}

func (impl *BigEndianStreamImpl) WriteByte(b byte) error {
	if impl.Right() < 1 {
		return ErrBuffOverflow
	}
	impl.buff[impl.len] = b
	impl.len += 1
	return nil
}

func (impl *BigEndianStreamImpl) WriteUint16(b uint16) error {
	if impl.Right() < 2 {
		return ErrBuffOverflow
	}
	binary.BigEndian.PutUint16(impl.buff[impl.len:], b)
	impl.len += 2
	return nil
}

func (impl *BigEndianStreamImpl) WriteUint32(b uint32) error {
	if impl.Right() < 4 {
		return ErrBuffOverflow
	}
	binary.BigEndian.PutUint32(impl.buff[impl.len:], b)
	impl.len += 4
	return nil
}

func (impl *BigEndianStreamImpl) WriteUint64(b uint64) error {
	if impl.Right() < 8 {
		return ErrBuffOverflow
	}
	binary.BigEndian.PutUint64(impl.buff[impl.len:], b)
	impl.len += 8
	return nil
}

func (impl *BigEndianStreamImpl) WriteBuff(buff []byte) error {
	if impl.Right() < len(buff) {
		return ErrBuffOverflow
	}
	copy(impl.buff[impl.len:], buff)
	impl.len += len(buff)
	return nil
}

// CleanBuff 将已经读取过的数据从缓存中清除
func (impl *BigEndianStreamImpl) CleanBuff() error {
	tmpBuffer, err := impl.ReadBuff(impl.len - impl.pos)
	if nil != err {
		return err
	}

	impl.Reset(make([]byte, len(impl.buff), len(impl.buff)))
	return impl.WriteBuff(tmpBuffer)
}

type LittleEndianStreamImpl struct {
	pos  int
	len  int
	buff []byte
}

func NewLittleEndianStream(buff []byte) *LittleEndianStreamImpl {
	return &LittleEndianStreamImpl{
		buff: buff,
	}
}

func (impl *LittleEndianStreamImpl) Size() int   { return len(impl.buff) }
func (impl *LittleEndianStreamImpl) Length() int { return impl.len }

func (impl *LittleEndianStreamImpl) Data() []byte { return impl.buff }

// Right 剩余空间大小
func (impl *LittleEndianStreamImpl) Right() int { return len(impl.buff) - impl.pos }

func (impl *LittleEndianStreamImpl) Reset(buff []byte) { impl.pos = 0; impl.len = 0; impl.buff = buff }

func (impl *LittleEndianStreamImpl) ReadByte() (b byte, err error) {
	if impl.Right() < 1 {
		return 0, ErrBuffOverflow
	}
	b = impl.buff[impl.pos]
	impl.pos += 1
	return b, nil
}

func (impl *LittleEndianStreamImpl) ReadUint16() (b uint16, err error) {
	if impl.Right() < 2 {
		return 0, ErrBuffOverflow
	}
	b = binary.LittleEndian.Uint16(impl.buff[impl.pos:])
	impl.pos += 2
	return b, nil
}

func (impl *LittleEndianStreamImpl) ReadUint32() (b uint32, err error) {
	if impl.Right() < 4 {
		return 0, ErrBuffOverflow
	}
	b = binary.LittleEndian.Uint32(impl.buff[impl.pos:])
	impl.pos += 4
	return b, nil
}

func (impl *LittleEndianStreamImpl) ReadUint64() (b uint64, err error) {
	if impl.Right() < 8 {
		return 0, ErrBuffOverflow
	}
	b = binary.LittleEndian.Uint64(impl.buff[impl.pos:])
	impl.pos += 8
	return b, nil
}

func (impl *LittleEndianStreamImpl) ReadBuff(size int) (buff []byte, err error) {
	if impl.Right() < size {
		return nil, ErrBuffOverflow
	}
	buff = make([]byte, size, size)
	copy(buff, impl.buff[impl.pos:impl.pos+size])
	impl.pos += size
	return buff, nil
}

// PeekBuff 从缓存中预读取指定字节数据不修改pos
func (impl *LittleEndianStreamImpl) PeekBuff(size int) (buff []byte, err error) {
	if impl.Right() < size {
		return nil, ErrBuffOverflow
	}
	buff = make([]byte, size, size)
	copy(buff, impl.buff[impl.pos:impl.pos+size])

	return buff, nil
}

func (impl *LittleEndianStreamImpl) CopyBuff(b []byte) error {
	if impl.Right() < len(b) {
		return ErrBuffOverflow
	}
	copy(b, impl.buff[impl.pos:impl.pos+len(b)])
	return nil
}

func (impl *LittleEndianStreamImpl) WriteByte(b byte) error {
	if impl.Right() < 1 {
		return ErrBuffOverflow
	}
	impl.buff[impl.len] = b
	impl.len += 1
	return nil
}

func (impl *LittleEndianStreamImpl) WriteUint16(b uint16) error {
	if impl.Right() < 2 {
		return ErrBuffOverflow
	}
	binary.LittleEndian.PutUint16(impl.buff[impl.len:], b)
	impl.len += 2
	return nil
}

func (impl *LittleEndianStreamImpl) WriteUint32(b uint32) error {
	if impl.Right() < 4 {
		return ErrBuffOverflow
	}
	binary.LittleEndian.PutUint32(impl.buff[impl.len:], b)
	impl.len += 4
	return nil
}

func (impl *LittleEndianStreamImpl) WriteUint64(b uint64) error {
	if impl.Right() < 8 {
		return ErrBuffOverflow
	}
	binary.LittleEndian.PutUint64(impl.buff[impl.len:], b)
	impl.len += 8
	return nil
}

func (impl *LittleEndianStreamImpl) WriteBuff(buff []byte) error {
	if impl.Right() < len(buff) {
		return ErrBuffOverflow
	}
	copy(impl.buff[impl.len:], buff)
	impl.len += len(buff)
	return nil
}

// CleanBuff 将已经读取过的数据从缓存中清除
func (impl *LittleEndianStreamImpl) CleanBuff() error {
	tmpBuffer, err := impl.ReadBuff(impl.len - impl.pos)
	if nil != err {
		return err
	}

	impl.Reset(make([]byte, len(impl.buff), len(impl.buff)))
	return impl.WriteBuff(tmpBuffer)
}
