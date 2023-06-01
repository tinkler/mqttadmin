// Code generated by github.com/tinkler/mqttadmin; DO NOT EDIT.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: page/v1/page.proto

package page_v1

import (
	_ "github.com/tinkler/mqttadmin/user/v1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Page struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Page    int32 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	PerPage int32 `protobuf:"varint,2,opt,name=per_page,json=perPage,proto3" json:"per_page,omitempty"`
	Total   int32 `protobuf:"varint,3,opt,name=total,proto3" json:"total,omitempty"`
}

func (x *Page) Reset() {
	*x = Page{}
	if protoimpl.UnsafeEnabled {
		mi := &file_page_v1_page_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Page) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Page) ProtoMessage() {}

func (x *Page) ProtoReflect() protoreflect.Message {
	mi := &file_page_v1_page_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Page.ProtoReflect.Descriptor instead.
func (*Page) Descriptor() ([]byte, []int) {
	return file_page_v1_page_proto_rawDescGZIP(), []int{0}
}

func (x *Page) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *Page) GetPerPage() int32 {
	if x != nil {
		return x.PerPage
	}
	return 0
}

func (x *Page) GetTotal() int32 {
	if x != nil {
		return x.Total
	}
	return 0
}

type PageRow struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RowNo int32 `protobuf:"varint,1,opt,name=row_no,json=rowNo,proto3" json:"row_no,omitempty"`
}

func (x *PageRow) Reset() {
	*x = PageRow{}
	if protoimpl.UnsafeEnabled {
		mi := &file_page_v1_page_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PageRow) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PageRow) ProtoMessage() {}

func (x *PageRow) ProtoReflect() protoreflect.Message {
	mi := &file_page_v1_page_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PageRow.ProtoReflect.Descriptor instead.
func (*PageRow) Descriptor() ([]byte, []int) {
	return file_page_v1_page_proto_rawDescGZIP(), []int{1}
}

func (x *PageRow) GetRowNo() int32 {
	if x != nil {
		return x.RowNo
	}
	return 0
}

var File_page_v1_page_proto protoreflect.FileDescriptor

var file_page_v1_page_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x61, 0x67, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x61, 0x67, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x70, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x19, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61,
	0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x12, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x76,
	0x31, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4b, 0x0a, 0x04,
	0x50, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x70, 0x65, 0x72, 0x5f,
	0x70, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x70, 0x65, 0x72, 0x50,
	0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x22, 0x20, 0x0a, 0x07, 0x50, 0x61, 0x67,
	0x65, 0x52, 0x6f, 0x77, 0x12, 0x15, 0x0a, 0x06, 0x72, 0x6f, 0x77, 0x5f, 0x6e, 0x6f, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x72, 0x6f, 0x77, 0x4e, 0x6f, 0x32, 0x88, 0x01, 0x0a, 0x08,
	0x50, 0x61, 0x67, 0x65, 0x47, 0x73, 0x72, 0x76, 0x12, 0x3b, 0x0a, 0x0d, 0x50, 0x61, 0x67, 0x65,
	0x46, 0x65, 0x74, 0x63, 0x68, 0x55, 0x73, 0x65, 0x72, 0x12, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x1a,
	0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x41, 0x6e, 0x79, 0x12, 0x3f, 0x0a, 0x0d, 0x50, 0x61, 0x67, 0x65, 0x52, 0x6f, 0x77,
	0x47, 0x65, 0x6e, 0x52, 0x6f, 0x77, 0x12, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x1a, 0x14, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41,
	0x6e, 0x79, 0x28, 0x01, 0x30, 0x01, 0x42, 0x5e, 0x0a, 0x21, 0x69, 0x6e, 0x6b, 0x2e, 0x73, 0x66,
	0x73, 0x2e, 0x74, 0x69, 0x6e, 0x6b, 0x6c, 0x65, 0x72, 0x2e, 0x6d, 0x71, 0x74, 0x74, 0x61, 0x64,
	0x6d, 0x69, 0x6e, 0x2e, 0x70, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x42, 0x09, 0x70, 0x61, 0x67,
	0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x69, 0x6e, 0x6b, 0x6c, 0x65, 0x72, 0x2f, 0x6d, 0x71, 0x74,
	0x74, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x70, 0x61, 0x67, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x70,
	0x61, 0x67, 0x65, 0x5f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_page_v1_page_proto_rawDescOnce sync.Once
	file_page_v1_page_proto_rawDescData = file_page_v1_page_proto_rawDesc
)

func file_page_v1_page_proto_rawDescGZIP() []byte {
	file_page_v1_page_proto_rawDescOnce.Do(func() {
		file_page_v1_page_proto_rawDescData = protoimpl.X.CompressGZIP(file_page_v1_page_proto_rawDescData)
	})
	return file_page_v1_page_proto_rawDescData
}

var file_page_v1_page_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_page_v1_page_proto_goTypes = []interface{}{
	(*Page)(nil),      // 0: page.v1.Page
	(*PageRow)(nil),   // 1: page.v1.PageRow
	(*anypb.Any)(nil), // 2: google.protobuf.Any
}
var file_page_v1_page_proto_depIdxs = []int32{
	2, // 0: page.v1.PageGsrv.PageFetchUser:input_type -> google.protobuf.Any
	2, // 1: page.v1.PageGsrv.PageRowGenRow:input_type -> google.protobuf.Any
	2, // 2: page.v1.PageGsrv.PageFetchUser:output_type -> google.protobuf.Any
	2, // 3: page.v1.PageGsrv.PageRowGenRow:output_type -> google.protobuf.Any
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_page_v1_page_proto_init() }
func file_page_v1_page_proto_init() {
	if File_page_v1_page_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_page_v1_page_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Page); i {
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
		file_page_v1_page_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PageRow); i {
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
			RawDescriptor: file_page_v1_page_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_page_v1_page_proto_goTypes,
		DependencyIndexes: file_page_v1_page_proto_depIdxs,
		MessageInfos:      file_page_v1_page_proto_msgTypes,
	}.Build()
	File_page_v1_page_proto = out.File
	file_page_v1_page_proto_rawDesc = nil
	file_page_v1_page_proto_goTypes = nil
	file_page_v1_page_proto_depIdxs = nil
}
