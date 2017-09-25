// Code generated by protoc-gen-go. DO NOT EDIT.
// source: data/user.proto

package data

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type User struct {
	Id       uint64 `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	HomeCell uint64 `protobuf:"varint,2,opt,name=home_cell,json=homeCell" json:"home_cell,omitempty"`
	WorkCell uint64 `protobuf:"varint,3,opt,name=work_cell,json=workCell" json:"work_cell,omitempty"`
}

func (m *User) Reset()                    { *m = User{} }
func (m *User) String() string            { return proto.CompactTextString(m) }
func (*User) ProtoMessage()               {}
func (*User) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *User) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *User) GetHomeCell() uint64 {
	if m != nil {
		return m.HomeCell
	}
	return 0
}

func (m *User) GetWorkCell() uint64 {
	if m != nil {
		return m.WorkCell
	}
	return 0
}

func init() {
	proto.RegisterType((*User)(nil), "data.User")
}

func init() { proto.RegisterFile("data/user.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 112 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4f, 0x49, 0x2c, 0x49,
	0xd4, 0x2f, 0x2d, 0x4e, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x01, 0x09, 0x28,
	0x05, 0x70, 0xb1, 0x84, 0x16, 0xa7, 0x16, 0x09, 0xf1, 0x71, 0x31, 0x65, 0xa6, 0x48, 0x30, 0x2a,
	0x30, 0x6a, 0xb0, 0x04, 0x31, 0x65, 0xa6, 0x08, 0x49, 0x73, 0x71, 0x66, 0xe4, 0xe7, 0xa6, 0xc6,
	0x27, 0xa7, 0xe6, 0xe4, 0x48, 0x30, 0x81, 0x85, 0x39, 0x40, 0x02, 0xce, 0xa9, 0x39, 0x39, 0x20,
	0xc9, 0xf2, 0xfc, 0xa2, 0x6c, 0x88, 0x24, 0x33, 0x44, 0x12, 0x24, 0x00, 0x92, 0x4c, 0x62, 0x03,
	0x1b, 0x6f, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xe1, 0x6d, 0x74, 0x2e, 0x71, 0x00, 0x00, 0x00,
}
