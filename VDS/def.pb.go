// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: def.proto

package vdsdata

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type VDSOrdType int32

const (
	VDSOrdType_ORD_MARKET VDSOrdType = 0 //市价
	VDSOrdType_ORD_LIMIT  VDSOrdType = 1 //限价
)

// Enum value maps for VDSOrdType.
var (
	VDSOrdType_name = map[int32]string{
		0: "ORD_MARKET",
		1: "ORD_LIMIT",
	}
	VDSOrdType_value = map[string]int32{
		"ORD_MARKET": 0,
		"ORD_LIMIT":  1,
	}
)

func (x VDSOrdType) Enum() *VDSOrdType {
	p := new(VDSOrdType)
	*p = x
	return p
}

func (x VDSOrdType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (VDSOrdType) Descriptor() protoreflect.EnumDescriptor {
	return file_def_proto_enumTypes[0].Descriptor()
}

func (VDSOrdType) Type() protoreflect.EnumType {
	return &file_def_proto_enumTypes[0]
}

func (x VDSOrdType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use VDSOrdType.Descriptor instead.
func (VDSOrdType) EnumDescriptor() ([]byte, []int) {
	return file_def_proto_rawDescGZIP(), []int{0}
}

type VDSTrdType int32

const (
	VDSTrdType_TRD_BUY    VDSTrdType = 0 //买
	VDSTrdType_TRD_SELL   VDSTrdType = 1 //卖
	VDSTrdType_TRD_CANCEL VDSTrdType = 2 //撤单
)

// Enum value maps for VDSTrdType.
var (
	VDSTrdType_name = map[int32]string{
		0: "TRD_BUY",
		1: "TRD_SELL",
		2: "TRD_CANCEL",
	}
	VDSTrdType_value = map[string]int32{
		"TRD_BUY":    0,
		"TRD_SELL":   1,
		"TRD_CANCEL": 2,
	}
)

func (x VDSTrdType) Enum() *VDSTrdType {
	p := new(VDSTrdType)
	*p = x
	return p
}

func (x VDSTrdType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (VDSTrdType) Descriptor() protoreflect.EnumDescriptor {
	return file_def_proto_enumTypes[1].Descriptor()
}

func (VDSTrdType) Type() protoreflect.EnumType {
	return &file_def_proto_enumTypes[1]
}

func (x VDSTrdType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use VDSTrdType.Descriptor instead.
func (VDSTrdType) EnumDescriptor() ([]byte, []int) {
	return file_def_proto_rawDescGZIP(), []int{1}
}

type VDSOrdSide int32

const (
	VDSOrdSide_OS_B VDSOrdSide = 0
	VDSOrdSide_OS_S VDSOrdSide = 1
)

// Enum value maps for VDSOrdSide.
var (
	VDSOrdSide_name = map[int32]string{
		0: "OS_B",
		1: "OS_S",
	}
	VDSOrdSide_value = map[string]int32{
		"OS_B": 0,
		"OS_S": 1,
	}
)

func (x VDSOrdSide) Enum() *VDSOrdSide {
	p := new(VDSOrdSide)
	*p = x
	return p
}

func (x VDSOrdSide) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (VDSOrdSide) Descriptor() protoreflect.EnumDescriptor {
	return file_def_proto_enumTypes[2].Descriptor()
}

func (VDSOrdSide) Type() protoreflect.EnumType {
	return &file_def_proto_enumTypes[2]
}

func (x VDSOrdSide) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use VDSOrdSide.Descriptor instead.
func (VDSOrdSide) EnumDescriptor() ([]byte, []int) {
	return file_def_proto_rawDescGZIP(), []int{2}
}

type VDSInterfaceType int32

const (
	VDSInterfaceType_ReqLogon       VDSInterfaceType = 0   //登录
	VDSInterfaceType_ReqStatCode    VDSInterfaceType = 1   //请求码表
	VDSInterfaceType_ReqKLine       VDSInterfaceType = 2   //请求K线
	VDSInterfaceType_ReqIndicator   VDSInterfaceType = 3   //请求指标
	VDSInterfaceType_SubSnapshot    VDSInterfaceType = 100 //订阅快照
	VDSInterfaceType_SubOrder       VDSInterfaceType = 101
	VDSInterfaceType_SubTransaction VDSInterfaceType = 102
)

// Enum value maps for VDSInterfaceType.
var (
	VDSInterfaceType_name = map[int32]string{
		0:   "ReqLogon",
		1:   "ReqStatCode",
		2:   "ReqKLine",
		3:   "ReqIndicator",
		100: "SubSnapshot",
		101: "SubOrder",
		102: "SubTransaction",
	}
	VDSInterfaceType_value = map[string]int32{
		"ReqLogon":       0,
		"ReqStatCode":    1,
		"ReqKLine":       2,
		"ReqIndicator":   3,
		"SubSnapshot":    100,
		"SubOrder":       101,
		"SubTransaction": 102,
	}
)

func (x VDSInterfaceType) Enum() *VDSInterfaceType {
	p := new(VDSInterfaceType)
	*p = x
	return p
}

func (x VDSInterfaceType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (VDSInterfaceType) Descriptor() protoreflect.EnumDescriptor {
	return file_def_proto_enumTypes[3].Descriptor()
}

func (VDSInterfaceType) Type() protoreflect.EnumType {
	return &file_def_proto_enumTypes[3]
}

func (x VDSInterfaceType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use VDSInterfaceType.Descriptor instead.
func (VDSInterfaceType) EnumDescriptor() ([]byte, []int) {
	return file_def_proto_rawDescGZIP(), []int{3}
}

type Exch int32

const (
	Exch_SZ Exch = 0
	Exch_SH Exch = 1
	Exch_BJ Exch = 2
)

// Enum value maps for Exch.
var (
	Exch_name = map[int32]string{
		0: "SZ",
		1: "SH",
		2: "BJ",
	}
	Exch_value = map[string]int32{
		"SZ": 0,
		"SH": 1,
		"BJ": 2,
	}
)

func (x Exch) Enum() *Exch {
	p := new(Exch)
	*p = x
	return p
}

func (x Exch) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Exch) Descriptor() protoreflect.EnumDescriptor {
	return file_def_proto_enumTypes[4].Descriptor()
}

func (Exch) Type() protoreflect.EnumType {
	return &file_def_proto_enumTypes[4]
}

func (x Exch) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Exch.Descriptor instead.
func (Exch) EnumDescriptor() ([]byte, []int) {
	return file_def_proto_rawDescGZIP(), []int{4}
}

type VDSSecType int32

const (
	VDSSecType_STOCK VDSSecType = 0
	VDSSecType_FUND  VDSSecType = 1
	VDSSecType_BOND  VDSSecType = 2
)

// Enum value maps for VDSSecType.
var (
	VDSSecType_name = map[int32]string{
		0: "STOCK",
		1: "FUND",
		2: "BOND",
	}
	VDSSecType_value = map[string]int32{
		"STOCK": 0,
		"FUND":  1,
		"BOND":  2,
	}
)

func (x VDSSecType) Enum() *VDSSecType {
	p := new(VDSSecType)
	*p = x
	return p
}

func (x VDSSecType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (VDSSecType) Descriptor() protoreflect.EnumDescriptor {
	return file_def_proto_enumTypes[5].Descriptor()
}

func (VDSSecType) Type() protoreflect.EnumType {
	return &file_def_proto_enumTypes[5]
}

func (x VDSSecType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use VDSSecType.Descriptor instead.
func (VDSSecType) EnumDescriptor() ([]byte, []int) {
	return file_def_proto_rawDescGZIP(), []int{5}
}

type VDSKLineType int32

const (
	VDSKLineType_KLINE_MIN1  VDSKLineType = 0
	VDSKLineType_KLINE_MIN10 VDSKLineType = 1
	VDSKLineType_KLINE_MIN30 VDSKLineType = 2
	VDSKLineType_KLINE_MIN60 VDSKLineType = 3
	VDSKLineType_KLINE_DAY   VDSKLineType = 10
	VDSKLineType_KLINE_WEEK  VDSKLineType = 20
	VDSKLineType_KLINE_MONTH VDSKLineType = 30
	VDSKLineType_KLINE_YEAR  VDSKLineType = 40
)

// Enum value maps for VDSKLineType.
var (
	VDSKLineType_name = map[int32]string{
		0:  "KLINE_MIN1",
		1:  "KLINE_MIN10",
		2:  "KLINE_MIN30",
		3:  "KLINE_MIN60",
		10: "KLINE_DAY",
		20: "KLINE_WEEK",
		30: "KLINE_MONTH",
		40: "KLINE_YEAR",
	}
	VDSKLineType_value = map[string]int32{
		"KLINE_MIN1":  0,
		"KLINE_MIN10": 1,
		"KLINE_MIN30": 2,
		"KLINE_MIN60": 3,
		"KLINE_DAY":   10,
		"KLINE_WEEK":  20,
		"KLINE_MONTH": 30,
		"KLINE_YEAR":  40,
	}
)

func (x VDSKLineType) Enum() *VDSKLineType {
	p := new(VDSKLineType)
	*p = x
	return p
}

func (x VDSKLineType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (VDSKLineType) Descriptor() protoreflect.EnumDescriptor {
	return file_def_proto_enumTypes[6].Descriptor()
}

func (VDSKLineType) Type() protoreflect.EnumType {
	return &file_def_proto_enumTypes[6]
}

func (x VDSKLineType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use VDSKLineType.Descriptor instead.
func (VDSKLineType) EnumDescriptor() ([]byte, []int) {
	return file_def_proto_rawDescGZIP(), []int{6}
}

var File_def_proto protoreflect.FileDescriptor

var file_def_proto_rawDesc = []byte{
	0x0a, 0x09, 0x64, 0x65, 0x66, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2a, 0x2b, 0x0a, 0x0a, 0x56,
	0x44, 0x53, 0x4f, 0x72, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0e, 0x0a, 0x0a, 0x4f, 0x52, 0x44,
	0x5f, 0x4d, 0x41, 0x52, 0x4b, 0x45, 0x54, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x4f, 0x52, 0x44,
	0x5f, 0x4c, 0x49, 0x4d, 0x49, 0x54, 0x10, 0x01, 0x2a, 0x37, 0x0a, 0x0a, 0x56, 0x44, 0x53, 0x54,
	0x72, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x54, 0x52, 0x44, 0x5f, 0x42, 0x55,
	0x59, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x54, 0x52, 0x44, 0x5f, 0x53, 0x45, 0x4c, 0x4c, 0x10,
	0x01, 0x12, 0x0e, 0x0a, 0x0a, 0x54, 0x52, 0x44, 0x5f, 0x43, 0x41, 0x4e, 0x43, 0x45, 0x4c, 0x10,
	0x02, 0x2a, 0x20, 0x0a, 0x0a, 0x56, 0x44, 0x53, 0x4f, 0x72, 0x64, 0x53, 0x69, 0x64, 0x65, 0x12,
	0x08, 0x0a, 0x04, 0x4f, 0x53, 0x5f, 0x42, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x4f, 0x53, 0x5f,
	0x53, 0x10, 0x01, 0x2a, 0x84, 0x01, 0x0a, 0x10, 0x56, 0x44, 0x53, 0x49, 0x6e, 0x74, 0x65, 0x72,
	0x66, 0x61, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0c, 0x0a, 0x08, 0x52, 0x65, 0x71, 0x4c,
	0x6f, 0x67, 0x6f, 0x6e, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x52, 0x65, 0x71, 0x53, 0x74, 0x61,
	0x74, 0x43, 0x6f, 0x64, 0x65, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x52, 0x65, 0x71, 0x4b, 0x4c,
	0x69, 0x6e, 0x65, 0x10, 0x02, 0x12, 0x10, 0x0a, 0x0c, 0x52, 0x65, 0x71, 0x49, 0x6e, 0x64, 0x69,
	0x63, 0x61, 0x74, 0x6f, 0x72, 0x10, 0x03, 0x12, 0x0f, 0x0a, 0x0b, 0x53, 0x75, 0x62, 0x53, 0x6e,
	0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x10, 0x64, 0x12, 0x0c, 0x0a, 0x08, 0x53, 0x75, 0x62, 0x4f,
	0x72, 0x64, 0x65, 0x72, 0x10, 0x65, 0x12, 0x12, 0x0a, 0x0e, 0x53, 0x75, 0x62, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x10, 0x66, 0x2a, 0x1e, 0x0a, 0x04, 0x45, 0x78,
	0x63, 0x68, 0x12, 0x06, 0x0a, 0x02, 0x53, 0x5a, 0x10, 0x00, 0x12, 0x06, 0x0a, 0x02, 0x53, 0x48,
	0x10, 0x01, 0x12, 0x06, 0x0a, 0x02, 0x42, 0x4a, 0x10, 0x02, 0x2a, 0x2b, 0x0a, 0x0a, 0x56, 0x44,
	0x53, 0x53, 0x65, 0x63, 0x54, 0x79, 0x70, 0x65, 0x12, 0x09, 0x0a, 0x05, 0x53, 0x54, 0x4f, 0x43,
	0x4b, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x46, 0x55, 0x4e, 0x44, 0x10, 0x01, 0x12, 0x08, 0x0a,
	0x04, 0x42, 0x4f, 0x4e, 0x44, 0x10, 0x02, 0x2a, 0x91, 0x01, 0x0a, 0x0c, 0x56, 0x44, 0x53, 0x4b,
	0x4c, 0x69, 0x6e, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0e, 0x0a, 0x0a, 0x4b, 0x4c, 0x49, 0x4e,
	0x45, 0x5f, 0x4d, 0x49, 0x4e, 0x31, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x4b, 0x4c, 0x49, 0x4e,
	0x45, 0x5f, 0x4d, 0x49, 0x4e, 0x31, 0x30, 0x10, 0x01, 0x12, 0x0f, 0x0a, 0x0b, 0x4b, 0x4c, 0x49,
	0x4e, 0x45, 0x5f, 0x4d, 0x49, 0x4e, 0x33, 0x30, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x4b, 0x4c,
	0x49, 0x4e, 0x45, 0x5f, 0x4d, 0x49, 0x4e, 0x36, 0x30, 0x10, 0x03, 0x12, 0x0d, 0x0a, 0x09, 0x4b,
	0x4c, 0x49, 0x4e, 0x45, 0x5f, 0x44, 0x41, 0x59, 0x10, 0x0a, 0x12, 0x0e, 0x0a, 0x0a, 0x4b, 0x4c,
	0x49, 0x4e, 0x45, 0x5f, 0x57, 0x45, 0x45, 0x4b, 0x10, 0x14, 0x12, 0x0f, 0x0a, 0x0b, 0x4b, 0x4c,
	0x49, 0x4e, 0x45, 0x5f, 0x4d, 0x4f, 0x4e, 0x54, 0x48, 0x10, 0x1e, 0x12, 0x0e, 0x0a, 0x0a, 0x4b,
	0x4c, 0x49, 0x4e, 0x45, 0x5f, 0x59, 0x45, 0x41, 0x52, 0x10, 0x28, 0x42, 0x0c, 0x5a, 0x0a, 0x2e,
	0x2f, 0x3b, 0x76, 0x64, 0x73, 0x64, 0x61, 0x74, 0x61, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_def_proto_rawDescOnce sync.Once
	file_def_proto_rawDescData = file_def_proto_rawDesc
)

func file_def_proto_rawDescGZIP() []byte {
	file_def_proto_rawDescOnce.Do(func() {
		file_def_proto_rawDescData = protoimpl.X.CompressGZIP(file_def_proto_rawDescData)
	})
	return file_def_proto_rawDescData
}

var file_def_proto_enumTypes = make([]protoimpl.EnumInfo, 7)
var file_def_proto_goTypes = []interface{}{
	(VDSOrdType)(0),       // 0: VDSOrdType
	(VDSTrdType)(0),       // 1: VDSTrdType
	(VDSOrdSide)(0),       // 2: VDSOrdSide
	(VDSInterfaceType)(0), // 3: VDSInterfaceType
	(Exch)(0),             // 4: Exch
	(VDSSecType)(0),       // 5: VDSSecType
	(VDSKLineType)(0),     // 6: VDSKLineType
}
var file_def_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_def_proto_init() }
func file_def_proto_init() {
	if File_def_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_def_proto_rawDesc,
			NumEnums:      7,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_def_proto_goTypes,
		DependencyIndexes: file_def_proto_depIdxs,
		EnumInfos:         file_def_proto_enumTypes,
	}.Build()
	File_def_proto = out.File
	file_def_proto_rawDesc = nil
	file_def_proto_goTypes = nil
	file_def_proto_depIdxs = nil
}
