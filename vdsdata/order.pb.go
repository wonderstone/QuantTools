// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: order.proto

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

//逐笔委托
type VDSOrder struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Symbol string `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`        //代码
	Exch   Exch   `protobuf:"varint,2,opt,name=exch,proto3,enum=Exch" json:"exch,omitempty"` //市场
	///消息序号
	//SH:逐笔成交单独序号， 在同一个 ChannelNo 内唯一,从 1 开始连续.
	//SZ:逐笔成交与委托统一序号， 在同一个 ChannelNo 内唯一， 从 1 开始连续
	Seqno      int64   `protobuf:"varint,3,opt,name=seqno,proto3" json:"seqno,omitempty"`
	Channalno  int32   `protobuf:"varint,4,opt,name=channalno,proto3" json:"channalno,omitempty"`   //频道
	Tradedate  int32   `protobuf:"varint,5,opt,name=tradedate,proto3" json:"tradedate,omitempty"`   //交易日期
	Updatetime int32   `protobuf:"varint,6,opt,name=updatetime,proto3" json:"updatetime,omitempty"` //更新时间
	Price      float64 `protobuf:"fixed64,7,opt,name=price,proto3" json:"price,omitempty"`
	Trdvol     float64 `protobuf:"fixed64,8,opt,name=trdvol,proto3" json:"trdvol,omitempty"`
	Trdmoney   float64 `protobuf:"fixed64,9,opt,name=trdmoney,proto3" json:"trdmoney,omitempty"`
	//int64   trdbuyno =10;
	//int64   trdsellno =11;
	Ordtype  VDSOrdType `protobuf:"varint,12,opt,name=ordtype,proto3,enum=VDSOrdType" json:"ordtype,omitempty"`
	Ordside  VDSOrdSide `protobuf:"varint,13,opt,name=ordside,proto3,enum=VDSOrdSide" json:"ordside,omitempty"`
	Bizindex int64      `protobuf:"varint,20,opt,name=bizindex,proto3" json:"bizindex,omitempty"`
}

func (x *VDSOrder) Reset() {
	*x = VDSOrder{}
	if protoimpl.UnsafeEnabled {
		mi := &file_order_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VDSOrder) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VDSOrder) ProtoMessage() {}

func (x *VDSOrder) ProtoReflect() protoreflect.Message {
	mi := &file_order_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VDSOrder.ProtoReflect.Descriptor instead.
func (*VDSOrder) Descriptor() ([]byte, []int) {
	return file_order_proto_rawDescGZIP(), []int{0}
}

func (x *VDSOrder) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *VDSOrder) GetExch() Exch {
	if x != nil {
		return x.Exch
	}
	return Exch_SZ
}

func (x *VDSOrder) GetSeqno() int64 {
	if x != nil {
		return x.Seqno
	}
	return 0
}

func (x *VDSOrder) GetChannalno() int32 {
	if x != nil {
		return x.Channalno
	}
	return 0
}

func (x *VDSOrder) GetTradedate() int32 {
	if x != nil {
		return x.Tradedate
	}
	return 0
}

func (x *VDSOrder) GetUpdatetime() int32 {
	if x != nil {
		return x.Updatetime
	}
	return 0
}

func (x *VDSOrder) GetPrice() float64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *VDSOrder) GetTrdvol() float64 {
	if x != nil {
		return x.Trdvol
	}
	return 0
}

func (x *VDSOrder) GetTrdmoney() float64 {
	if x != nil {
		return x.Trdmoney
	}
	return 0
}

func (x *VDSOrder) GetOrdtype() VDSOrdType {
	if x != nil {
		return x.Ordtype
	}
	return VDSOrdType_ORD_MARKET
}

func (x *VDSOrder) GetOrdside() VDSOrdSide {
	if x != nil {
		return x.Ordside
	}
	return VDSOrdSide_OS_B
}

func (x *VDSOrder) GetBizindex() int64 {
	if x != nil {
		return x.Bizindex
	}
	return 0
}

var File_order_proto protoreflect.FileDescriptor

var file_order_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x09, 0x64,
	0x65, 0x66, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe3, 0x02, 0x0a, 0x08, 0x56, 0x44, 0x53,
	0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x12, 0x19, 0x0a,
	0x04, 0x65, 0x78, 0x63, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x05, 0x2e, 0x45, 0x78,
	0x63, 0x68, 0x52, 0x04, 0x65, 0x78, 0x63, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x65, 0x71, 0x6e,
	0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x73, 0x65, 0x71, 0x6e, 0x6f, 0x12, 0x1c,
	0x0a, 0x09, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x61, 0x6c, 0x6e, 0x6f, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x09, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x61, 0x6c, 0x6e, 0x6f, 0x12, 0x1c, 0x0a, 0x09,
	0x74, 0x72, 0x61, 0x64, 0x65, 0x64, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x09, 0x74, 0x72, 0x61, 0x64, 0x65, 0x64, 0x61, 0x74, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72,
	0x69, 0x63, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x74, 0x72, 0x64, 0x76, 0x6f, 0x6c, 0x18, 0x08, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x06, 0x74, 0x72, 0x64, 0x76, 0x6f, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x74, 0x72, 0x64, 0x6d,
	0x6f, 0x6e, 0x65, 0x79, 0x18, 0x09, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x74, 0x72, 0x64, 0x6d,
	0x6f, 0x6e, 0x65, 0x79, 0x12, 0x25, 0x0a, 0x07, 0x6f, 0x72, 0x64, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x0c, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0b, 0x2e, 0x56, 0x44, 0x53, 0x4f, 0x72, 0x64, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x74, 0x79, 0x70, 0x65, 0x12, 0x25, 0x0a, 0x07, 0x6f,
	0x72, 0x64, 0x73, 0x69, 0x64, 0x65, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0b, 0x2e, 0x56,
	0x44, 0x53, 0x4f, 0x72, 0x64, 0x53, 0x69, 0x64, 0x65, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x73, 0x69,
	0x64, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x62, 0x69, 0x7a, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x14,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x62, 0x69, 0x7a, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x42, 0x0c,
	0x5a, 0x0a, 0x2e, 0x2f, 0x3b, 0x76, 0x64, 0x73, 0x64, 0x61, 0x74, 0x61, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_order_proto_rawDescOnce sync.Once
	file_order_proto_rawDescData = file_order_proto_rawDesc
)

func file_order_proto_rawDescGZIP() []byte {
	file_order_proto_rawDescOnce.Do(func() {
		file_order_proto_rawDescData = protoimpl.X.CompressGZIP(file_order_proto_rawDescData)
	})
	return file_order_proto_rawDescData
}

var file_order_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_order_proto_goTypes = []interface{}{
	(*VDSOrder)(nil), // 0: VDSOrder
	(Exch)(0),        // 1: Exch
	(VDSOrdType)(0),  // 2: VDSOrdType
	(VDSOrdSide)(0),  // 3: VDSOrdSide
}
var file_order_proto_depIdxs = []int32{
	1, // 0: VDSOrder.exch:type_name -> Exch
	2, // 1: VDSOrder.ordtype:type_name -> VDSOrdType
	3, // 2: VDSOrder.ordside:type_name -> VDSOrdSide
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_order_proto_init() }
func file_order_proto_init() {
	if File_order_proto != nil {
		return
	}
	file_def_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_order_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VDSOrder); i {
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
			RawDescriptor: file_order_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_order_proto_goTypes,
		DependencyIndexes: file_order_proto_depIdxs,
		MessageInfos:      file_order_proto_msgTypes,
	}.Build()
	File_order_proto = out.File
	file_order_proto_rawDesc = nil
	file_order_proto_goTypes = nil
	file_order_proto_depIdxs = nil
}