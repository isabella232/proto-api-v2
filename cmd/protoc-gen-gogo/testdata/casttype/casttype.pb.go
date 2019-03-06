// Code generated by protoc-gen-go. DO NOT EDIT.
// source: casttype.proto

package casttype

import (
	_ "github.com/gogo/protobuf/gogoproto"
	casttype "github.com/gogo/protobuf/test/casttype"
	proto "github.com/golang/protobuf/proto"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Castaway struct {
	Int32Ptr             *int32                  `protobuf:"varint,1,opt,name=Int32Ptr,casttype=int32" json:"Int32Ptr,omitempty"`
	Int32                int32                   `protobuf:"varint,2,opt,name=Int32,casttype=int32" json:"Int32,omitempty"`
	MyUint64Ptr          *casttype.MyUint64Type  `protobuf:"varint,3,opt,name=MyUint64Ptr,casttype=github.com/gogo/protobuf/test/casttype.MyUint64Type" json:"MyUint64Ptr,omitempty"`
	MyUint64             casttype.MyUint64Type   `protobuf:"varint,4,opt,name=MyUint64,casttype=github.com/gogo/protobuf/test/casttype.MyUint64Type" json:"MyUint64,omitempty"`
	MyFloat32Ptr         *casttype.MyFloat32Type `protobuf:"fixed32,5,opt,name=MyFloat32Ptr,casttype=github.com/gogo/protobuf/test/casttype.MyFloat32Type" json:"MyFloat32Ptr,omitempty"`
	MyFloat32            casttype.MyFloat32Type  `protobuf:"fixed32,6,opt,name=MyFloat32,casttype=github.com/gogo/protobuf/test/casttype.MyFloat32Type" json:"MyFloat32,omitempty"`
	MyFloat64Ptr         *casttype.MyFloat64Type `protobuf:"fixed64,7,opt,name=MyFloat64Ptr,casttype=github.com/gogo/protobuf/test/casttype.MyFloat64Type" json:"MyFloat64Ptr,omitempty"`
	MyFloat64            casttype.MyFloat64Type  `protobuf:"fixed64,8,opt,name=MyFloat64,casttype=github.com/gogo/protobuf/test/casttype.MyFloat64Type" json:"MyFloat64,omitempty"`
	MyBytes              *casttype.Bytes         `protobuf:"bytes,9,opt,name=MyBytes,casttype=github.com/gogo/protobuf/test/casttype.Bytes" json:"MyBytes,omitempty"`
	NormalBytes          []byte                  `protobuf:"bytes,10,opt,name=NormalBytes'" json:"NormalBytes,omitempty"`
	MyUint64S            *casttype.MyUint64Type  `protobuf:"varint,11,rep,name=MyUint64s,casttype=github.com/gogo/protobuf/test/casttype.MyUint64Type" json:"MyUint64s,omitempty"`
	MyMap                map[string]uint64       `protobuf:"bytes,12,rep,name=MyMap,casttype=github.com/gogo/protobuf/test/casttype.MyMapType" json:"MyMap,omitempty" protobuf_key:"bytes,1,opt,name=key'" protobuf_val:"varint,2,opt,name=value'"`
	MyCustomMap          map[string]uint64       `protobuf:"bytes,13,rep,name=MyCustomMap'" json:"MyCustomMap,omitempty" protobuf_key:"bytes,1,opt,name=key'" protobuf_val:"varint,2,opt,name=value'"`
	MyNullableMap        map[int32]*Wilson       `protobuf:"bytes,14,rep,name=MyNullableMap'" json:"MyNullableMap,omitempty" protobuf_key:"varint,1,opt,name=key'" protobuf_val:"bytes,2,opt,name=value'"`
	MyEmbeddedMap        map[int32]*Wilson       `protobuf:"bytes,15,rep,name=MyEmbeddedMap'" json:"MyEmbeddedMap,omitempty" protobuf_key:"varint,1,opt,name=key'" protobuf_val:"bytes,2,opt,name=value'"`
	String_              *casttype.MyStringType  `protobuf:"bytes,16,opt,name=String,casttype=github.com/gogo/protobuf/test/casttype.MyStringType" json:"String,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *Castaway) Reset()         { *m = Castaway{} }
func (m *Castaway) String() string { return proto.CompactTextString(m) }
func (*Castaway) ProtoMessage()    {}
func (*Castaway) Descriptor() ([]byte, []int) {
	return fileDescriptor_28473e82a641ed5b, []int{0}
}

func (m *Castaway) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Castaway.Unmarshal(m, b)
}
func (m *Castaway) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Castaway.Marshal(b, m, deterministic)
}
func (m *Castaway) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Castaway.Merge(m, src)
}
func (m *Castaway) XXX_Size() int {
	return xxx_messageInfo_Castaway.Size(m)
}
func (m *Castaway) XXX_DiscardUnknown() {
	xxx_messageInfo_Castaway.DiscardUnknown(m)
}

var xxx_messageInfo_Castaway proto.InternalMessageInfo

func (m *Castaway) GetInt32Ptr() int32 {
	if m != nil && m.Int32Ptr != nil {
		return *m.Int32Ptr
	}
	return 0
}

func (m *Castaway) GetInt32() int32 {
	if m != nil && m.Int32 != nil {
		return m.Int32
	}
	return 0
}

func (m *Castaway) GetMyUint64Ptr() casttype.MyUint64Type {
	if m != nil && m.MyUint64Ptr != nil {
		return *m.MyUint64Ptr
	}
	return 0
}

func (m *Castaway) GetMyUint64() casttype.MyUint64Type {
	if m != nil && m.MyUint64 != nil {
		return m.MyUint64
	}
	return 0
}

func (m *Castaway) GetMyFloat32Ptr() casttype.MyFloat32Type {
	if m != nil && m.MyFloat32Ptr != nil {
		return *m.MyFloat32Ptr
	}
	return 0
}

func (m *Castaway) GetMyFloat32() casttype.MyFloat32Type {
	if m != nil && m.MyFloat32 != nil {
		return m.MyFloat32
	}
	return 0
}

func (m *Castaway) GetMyFloat64Ptr() casttype.MyFloat64Type {
	if m != nil && m.MyFloat64Ptr != nil {
		return *m.MyFloat64Ptr
	}
	return 0
}

func (m *Castaway) GetMyFloat64() casttype.MyFloat64Type {
	if m != nil && m.MyFloat64 != nil {
		return m.MyFloat64
	}
	return 0
}

func (m *Castaway) GetMyBytes() casttype.Bytes {
	if m != nil {
		return *m.MyBytes
	}
	return nil
}

func (m *Castaway) GetNormalBytes() []byte {
	if m != nil {
		return m.NormalBytes
	}
	return nil
}

func (m *Castaway) GetMyUint64S() casttype.MyUint64Type {
	if m != nil {
		return *m.MyUint64S
	}
	return nil
}

func (m *Castaway) GetMyMap() map[string]uint64 {
	if m != nil {
		return m.MyMap
	}
	return nil
}

func (m *Castaway) GetMyCustomMap() map[string]uint64 {
	if m != nil {
		return m.MyCustomMap
	}
	return nil
}

func (m *Castaway) GetMyNullableMap() map[int32]*Wilson {
	if m != nil {
		return m.MyNullableMap
	}
	return nil
}

func (m *Castaway) GetMyEmbeddedMap() map[int32]*Wilson {
	if m != nil {
		return m.MyEmbeddedMap
	}
	return nil
}

func (m *Castaway) GetString_() casttype.MyStringType {
	if m != nil && m.String_ != nil {
		return *m.String_
	}
	return ""
}

type Wilson struct {
	Int64                *int64   `protobuf:"varint,1,opt,name=Int64'" json:"Int64,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Wilson) Reset()         { *m = Wilson{} }
func (m *Wilson) String() string { return proto.CompactTextString(m) }
func (*Wilson) ProtoMessage()    {}
func (*Wilson) Descriptor() ([]byte, []int) {
	return fileDescriptor_28473e82a641ed5b, []int{1}
}

func (m *Wilson) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Wilson.Unmarshal(m, b)
}
func (m *Wilson) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Wilson.Marshal(b, m, deterministic)
}
func (m *Wilson) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Wilson.Merge(m, src)
}
func (m *Wilson) XXX_Size() int {
	return xxx_messageInfo_Wilson.Size(m)
}
func (m *Wilson) XXX_DiscardUnknown() {
	xxx_messageInfo_Wilson.DiscardUnknown(m)
}

var xxx_messageInfo_Wilson proto.InternalMessageInfo

func (m *Wilson) GetInt64() int64 {
	if m != nil && m.Int64 != nil {
		return *m.Int64
	}
	return 0
}

func init() {
	proto.RegisterFile("casttype.proto", fileDescriptor_28473e82a641ed5b)
	proto.RegisterType((*Castaway)(nil), "casttype.Castaway")
	proto.RegisterMapType((map[string]uint64)(nil), "casttype.Castaway.MyCustomMapEntry")
	proto.RegisterMapType((map[int32]*Wilson)(nil), "casttype.Castaway.MyEmbeddedMapEntry")
	proto.RegisterMapType((map[string]uint64)(nil), "casttype.Castaway.MyMapEntry")
	proto.RegisterMapType((map[int32]*Wilson)(nil), "casttype.Castaway.MyNullableMapEntry")
	proto.RegisterType((*Wilson)(nil), "casttype.Wilson")
}

var fileDescriptor_28473e82a641ed5b = []byte{
	// 685 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x95, 0xbf, 0x6f, 0xd3, 0x4c,
	0x18, 0xc7, 0x7d, 0x4d, 0xd3, 0x26, 0x97, 0xa6, 0x6f, 0x74, 0x7a, 0x07, 0xab, 0x12, 0x67, 0xab,
	0x55, 0x91, 0x07, 0x48, 0xab, 0x34, 0x2a, 0x55, 0x41, 0x0c, 0xae, 0x8a, 0x54, 0x84, 0x0b, 0x32,
	0x54, 0x15, 0x88, 0xc5, 0x69, 0x4d, 0xb0, 0x70, 0xec, 0x28, 0xbe, 0x80, 0xbc, 0x55, 0x65, 0x40,
	0xe2, 0x2f, 0x61, 0x64, 0x41, 0x62, 0x64, 0xec, 0xd8, 0x91, 0xc9, 0xad, 0xcd, 0x52, 0xb6, 0x8e,
	0x55, 0x26, 0x74, 0x77, 0x8e, 0xed, 0xfe, 0x00, 0xa5, 0xee, 0x76, 0xcf, 0xdd, 0xf3, 0x7c, 0x9e,
	0xef, 0x3d, 0xf7, 0xdc, 0x1d, 0x9c, 0xde, 0x31, 0x3c, 0x42, 0xfc, 0xae, 0x59, 0xef, 0xf6, 0x5c,
	0xe2, 0xa2, 0xd2, 0xd0, 0x9e, 0xb9, 0xdb, 0xb6, 0xc8, 0xdb, 0x7e, 0xab, 0xbe, 0xe3, 0x76, 0x16,
	0xda, 0x6e, 0xdb, 0x5d, 0x60, 0x0e, 0xad, 0xfe, 0x1b, 0x66, 0x31, 0x83, 0x8d, 0x78, 0xe0, 0xec,
	0xef, 0x2a, 0x2c, 0xad, 0x19, 0x1e, 0x31, 0x3e, 0x18, 0x3e, 0x9a, 0x87, 0xa5, 0x0d, 0x87, 0x2c,
	0x35, 0x9e, 0x91, 0x9e, 0x08, 0x64, 0xa0, 0x14, 0xd4, 0xf2, 0x20, 0x90, 0x8a, 0x16, 0x9d, 0xd3,
	0x93, 0x25, 0x34, 0x07, 0x8b, 0x6c, 0x2c, 0x8e, 0x31, 0x9f, 0xea, 0x41, 0x20, 0x09, 0xa9, 0x1f,
	0x5f, 0x43, 0x2f, 0x61, 0x45, 0xf3, 0xb7, 0x2c, 0x87, 0x2c, 0x37, 0x29, 0xae, 0x20, 0x03, 0x65,
	0x5c, 0xbd, 0x37, 0x08, 0xa4, 0xa5, 0xbf, 0x0a, 0x24, 0xa6, 0x47, 0x16, 0x92, 0x8d, 0x0d, 0xa3,
	0x5f, 0xf8, 0x5d, 0x53, 0xcf, 0xb2, 0xd0, 0x36, 0x2c, 0x0d, 0x4d, 0x71, 0x9c, 0x71, 0xef, 0xc7,
	0x12, 0x72, 0xb1, 0x13, 0x18, 0x7a, 0x0d, 0xa7, 0x34, 0xff, 0x91, 0xed, 0x1a, 0x71, 0x0d, 0x8a,
	0x32, 0x50, 0xc6, 0xd4, 0x95, 0x41, 0x20, 0x35, 0x47, 0x06, 0xc7, 0xe1, 0x8c, 0x7c, 0x8e, 0x86,
	0x5e, 0xc1, 0x72, 0x62, 0x8b, 0x13, 0x0c, 0xfd, 0x20, 0xd6, 0x9d, 0x0f, 0x9f, 0xe2, 0x32, 0xca,
	0x79, 0xb9, 0x27, 0x65, 0xa0, 0x80, 0x3c, 0xca, 0xe3, 0x9a, 0x9c, 0xa3, 0x65, 0x94, 0x2f, 0x37,
	0xc5, 0x12, 0x43, 0xe7, 0x54, 0x1e, 0xe3, 0x53, 0x1c, 0x7a, 0x0c, 0x27, 0x35, 0x5f, 0xf5, 0x89,
	0xe9, 0x89, 0x65, 0x19, 0x28, 0x53, 0xea, 0xe2, 0x20, 0x90, 0xee, 0x8c, 0x48, 0x65, 0x71, 0xfa,
	0x10, 0x80, 0x64, 0x58, 0xd9, 0x74, 0x7b, 0x1d, 0xc3, 0xe6, 0x3c, 0x48, 0x79, 0x7a, 0x76, 0x0a,
	0x6d, 0xd1, 0x9d, 0xf0, 0xd3, 0xf6, 0xc4, 0x8a, 0x5c, 0xb8, 0x49, 0x4f, 0xa6, 0x24, 0x64, 0xc1,
	0xa2, 0xe6, 0x6b, 0x46, 0x57, 0x9c, 0x92, 0x0b, 0x4a, 0xa5, 0x71, 0xab, 0x9e, 0x44, 0x0c, 0xef,
	0x56, 0x9d, 0xad, 0xaf, 0x3b, 0xa4, 0xe7, 0xab, 0xcd, 0x41, 0x20, 0x2d, 0x8e, 0x9c, 0x51, 0x33,
	0xba, 0x2c, 0x1d, 0xcf, 0x80, 0xbe, 0x01, 0x7a, 0xb1, 0xd6, 0xfa, 0x1e, 0x71, 0x3b, 0x34, 0x63,
	0x95, 0x65, 0x9c, 0xbb, 0x32, 0x63, 0xe2, 0xc5, 0xf3, 0x3a, 0xfb, 0x47, 0xd7, 0xd8, 0xe9, 0x73,
	0xd2, 0xb3, 0x9c, 0x36, 0x4d, 0xfd, 0xf9, 0x28, 0xf7, 0xa5, 0x4d, 0x14, 0xa0, 0x8f, 0x00, 0x56,
	0x35, 0x7f, 0xb3, 0x6f, 0xdb, 0x46, 0xcb, 0x36, 0xa9, 0xf2, 0x69, 0xa6, 0x7c, 0xfe, 0x4a, 0xe5,
	0x19, 0x3f, 0xae, 0x7d, 0x79, 0xff, 0x48, 0x6a, 0x8c, 0x2c, 0x82, 0x3d, 0x41, 0x4c, 0xc3, 0xf9,
	0x9c, 0xe8, 0x13, 0x53, 0xb1, 0xde, 0x69, 0x99, 0xbb, 0xbb, 0xe6, 0x2e, 0x55, 0xf1, 0xdf, 0x3f,
	0x54, 0x64, 0xfc, 0xb8, 0x8a, 0x55, 0xda, 0xf5, 0xf9, 0x95, 0x64, 0x78, 0xe8, 0x29, 0x9c, 0xe0,
	0x15, 0x16, 0x6b, 0x32, 0x50, 0xca, 0xd7, 0x6c, 0xc3, 0xf4, 0x70, 0xf4, 0x18, 0x33, 0xb3, 0x02,
	0x61, 0xda, 0x63, 0xa8, 0x06, 0x0b, 0xef, 0x4c, 0x9f, 0xbd, 0xe2, 0x65, 0x9d, 0x0e, 0xd1, 0xff,
	0xb0, 0xf8, 0xde, 0xb0, 0xfb, 0x26, 0x7b, 0xb5, 0xc7, 0x75, 0x6e, 0xac, 0x8e, 0xad, 0x80, 0x99,
	0x87, 0xb0, 0x76, 0xb1, 0x57, 0xae, 0x15, 0xaf, 0x43, 0x74, 0xf9, 0xc4, 0xb2, 0x84, 0x22, 0x27,
	0xdc, 0xce, 0x12, 0x2a, 0x8d, 0x5a, 0x5a, 0xf3, 0x6d, 0xcb, 0xf6, 0x5c, 0xe7, 0x12, 0xf3, 0x62,
	0xfd, 0x6f, 0xc6, 0x9c, 0xc5, 0x70, 0x82, 0x4f, 0xd2, 0xbd, 0x6c, 0xb0, 0xef, 0x83, 0xfd, 0x72,
	0x3a, 0x37, 0xd4, 0x27, 0x07, 0x21, 0x16, 0x0e, 0x43, 0x2c, 0xfc, 0x0c, 0xb1, 0x70, 0x1c, 0x62,
	0x70, 0x12, 0x62, 0x70, 0x1a, 0x62, 0x70, 0x16, 0x62, 0xb0, 0x17, 0x61, 0xf0, 0x25, 0xc2, 0xe0,
	0x6b, 0x84, 0xc1, 0xf7, 0x08, 0x83, 0x1f, 0x11, 0x06, 0x07, 0x11, 0x16, 0x0e, 0x23, 0x2c, 0x1c,
	0x47, 0x18, 0x9c, 0x44, 0x58, 0x38, 0x8d, 0x30, 0x38, 0x8b, 0xb0, 0xb0, 0xf7, 0x0b, 0x0b, 0x7f,
	0x02, 0x00, 0x00, 0xff, 0xff, 0xdd, 0xa0, 0x64, 0x14, 0xa3, 0x07, 0x00, 0x00,
}
