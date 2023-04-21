// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: statcode.proto

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

//静态码表结构
type VDSStatCode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Symbol    string     `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`                           //代码
	Name      string     `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`                               //名称
	Exch      Exch       `protobuf:"varint,3,opt,name=exch,proto3,enum=Exch" json:"exch,omitempty"`                    //市场
	Sectype   VDSSecType `protobuf:"varint,4,opt,name=sectype,proto3,enum=VDSSecType" json:"sectype,omitempty"`        //类别
	PreClose  float64    `protobuf:"fixed64,11,opt,name=pre_close,json=preClose,proto3" json:"pre_close,omitempty"`    //昨收价
	PriceUp   float64    `protobuf:"fixed64,12,opt,name=price_up,json=priceUp,proto3" json:"price_up,omitempty"`       //涨停价
	PriceDown float64    `protobuf:"fixed64,13,opt,name=price_down,json=priceDown,proto3" json:"price_down,omitempty"` //跌停价
	Buyunit   int32      `protobuf:"varint,14,opt,name=buyunit,proto3" json:"buyunit,omitempty"`                       //最小买
	Sellunit  int32      `protobuf:"varint,15,opt,name=sellunit,proto3" json:"sellunit,omitempty"`                     //最小卖
	Date      int32      `protobuf:"varint,16,opt,name=date,proto3" json:"date,omitempty"`                             //交易日期
}

func (x *VDSStatCode) Reset() {
	*x = VDSStatCode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_statcode_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VDSStatCode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VDSStatCode) ProtoMessage() {}

func (x *VDSStatCode) ProtoReflect() protoreflect.Message {
	mi := &file_statcode_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VDSStatCode.ProtoReflect.Descriptor instead.
func (*VDSStatCode) Descriptor() ([]byte, []int) {
	return file_statcode_proto_rawDescGZIP(), []int{0}
}

func (x *VDSStatCode) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *VDSStatCode) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *VDSStatCode) GetExch() Exch {
	if x != nil {
		return x.Exch
	}
	return Exch_SZ
}

func (x *VDSStatCode) GetSectype() VDSSecType {
	if x != nil {
		return x.Sectype
	}
	return VDSSecType_STOCK
}

func (x *VDSStatCode) GetPreClose() float64 {
	if x != nil {
		return x.PreClose
	}
	return 0
}

func (x *VDSStatCode) GetPriceUp() float64 {
	if x != nil {
		return x.PriceUp
	}
	return 0
}

func (x *VDSStatCode) GetPriceDown() float64 {
	if x != nil {
		return x.PriceDown
	}
	return 0
}

func (x *VDSStatCode) GetBuyunit() int32 {
	if x != nil {
		return x.Buyunit
	}
	return 0
}

func (x *VDSStatCode) GetSellunit() int32 {
	if x != nil {
		return x.Sellunit
	}
	return 0
}

func (x *VDSStatCode) GetDate() int32 {
	if x != nil {
		return x.Date
	}
	return 0
}

var File_statcode_proto protoreflect.FileDescriptor

var file_statcode_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x73, 0x74, 0x61, 0x74, 0x63, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x09, 0x64, 0x65, 0x66, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9c, 0x02, 0x0a, 0x0b,
	0x56, 0x44, 0x53, 0x53, 0x74, 0x61, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73,
	0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x6d,
	0x62, 0x6f, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x04, 0x65, 0x78, 0x63, 0x68, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x05, 0x2e, 0x45, 0x78, 0x63, 0x68, 0x52, 0x04, 0x65, 0x78,
	0x63, 0x68, 0x12, 0x25, 0x0a, 0x07, 0x73, 0x65, 0x63, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x0b, 0x2e, 0x56, 0x44, 0x53, 0x53, 0x65, 0x63, 0x54, 0x79, 0x70, 0x65,
	0x52, 0x07, 0x73, 0x65, 0x63, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x72, 0x65,
	0x5f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x70, 0x72,
	0x65, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x70, 0x72, 0x69, 0x63, 0x65, 0x5f,
	0x75, 0x70, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x01, 0x52, 0x07, 0x70, 0x72, 0x69, 0x63, 0x65, 0x55,
	0x70, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x69, 0x63, 0x65, 0x5f, 0x64, 0x6f, 0x77, 0x6e, 0x18,
	0x0d, 0x20, 0x01, 0x28, 0x01, 0x52, 0x09, 0x70, 0x72, 0x69, 0x63, 0x65, 0x44, 0x6f, 0x77, 0x6e,
	0x12, 0x18, 0x0a, 0x07, 0x62, 0x75, 0x79, 0x75, 0x6e, 0x69, 0x74, 0x18, 0x0e, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x07, 0x62, 0x75, 0x79, 0x75, 0x6e, 0x69, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65,
	0x6c, 0x6c, 0x75, 0x6e, 0x69, 0x74, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x73, 0x65,
	0x6c, 0x6c, 0x75, 0x6e, 0x69, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x10,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x64, 0x61, 0x74, 0x65, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f,
	0x3b, 0x76, 0x64, 0x73, 0x64, 0x61, 0x74, 0x61, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_statcode_proto_rawDescOnce sync.Once
	file_statcode_proto_rawDescData = file_statcode_proto_rawDesc
)

func file_statcode_proto_rawDescGZIP() []byte {
	file_statcode_proto_rawDescOnce.Do(func() {
		file_statcode_proto_rawDescData = protoimpl.X.CompressGZIP(file_statcode_proto_rawDescData)
	})
	return file_statcode_proto_rawDescData
}

var file_statcode_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_statcode_proto_goTypes = []interface{}{
	(*VDSStatCode)(nil), // 0: VDSStatCode
	(Exch)(0),           // 1: Exch
	(VDSSecType)(0),     // 2: VDSSecType
}
var file_statcode_proto_depIdxs = []int32{
	1, // 0: VDSStatCode.exch:type_name -> Exch
	2, // 1: VDSStatCode.sectype:type_name -> VDSSecType
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_statcode_proto_init() }
func file_statcode_proto_init() {
	if File_statcode_proto != nil {
		return
	}
	file_def_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_statcode_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VDSStatCode); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_statcode_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_statcode_proto_goTypes,
		DependencyIndexes: file_statcode_proto_depIdxs,
		MessageInfos:      file_statcode_proto_msgTypes,
	}.Build()
	File_statcode_proto = out.File
	file_statcode_proto_rawDesc = nil
	file_statcode_proto_goTypes = nil
	file_statcode_proto_depIdxs = nil
}
