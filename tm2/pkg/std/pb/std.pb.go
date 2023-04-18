// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: std.proto

package pb

import (
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

// messages
type BaseAccount struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address       string     `protobuf:"bytes,1,opt,name=Address,proto3" json:"Address,omitempty"`
	Coins         string     `protobuf:"bytes,2,opt,name=Coins,proto3" json:"Coins,omitempty"`
	PubKey        *anypb.Any `protobuf:"bytes,3,opt,name=PubKey,proto3" json:"PubKey,omitempty"`
	AccountNumber uint64     `protobuf:"varint,4,opt,name=AccountNumber,proto3" json:"AccountNumber,omitempty"`
	Sequence      uint64     `protobuf:"varint,5,opt,name=Sequence,proto3" json:"Sequence,omitempty"`
}

func (x *BaseAccount) Reset() {
	*x = BaseAccount{}
	if protoimpl.UnsafeEnabled {
		mi := &file_std_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BaseAccount) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BaseAccount) ProtoMessage() {}

func (x *BaseAccount) ProtoReflect() protoreflect.Message {
	mi := &file_std_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BaseAccount.ProtoReflect.Descriptor instead.
func (*BaseAccount) Descriptor() ([]byte, []int) {
	return file_std_proto_rawDescGZIP(), []int{0}
}

func (x *BaseAccount) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *BaseAccount) GetCoins() string {
	if x != nil {
		return x.Coins
	}
	return ""
}

func (x *BaseAccount) GetPubKey() *anypb.Any {
	if x != nil {
		return x.PubKey
	}
	return nil
}

func (x *BaseAccount) GetAccountNumber() uint64 {
	if x != nil {
		return x.AccountNumber
	}
	return 0
}

func (x *BaseAccount) GetSequence() uint64 {
	if x != nil {
		return x.Sequence
	}
	return 0
}

type MemFile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	Body string `protobuf:"bytes,2,opt,name=Body,proto3" json:"Body,omitempty"`
}

func (x *MemFile) Reset() {
	*x = MemFile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_std_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MemFile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MemFile) ProtoMessage() {}

func (x *MemFile) ProtoReflect() protoreflect.Message {
	mi := &file_std_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MemFile.ProtoReflect.Descriptor instead.
func (*MemFile) Descriptor() ([]byte, []int) {
	return file_std_proto_rawDescGZIP(), []int{1}
}

func (x *MemFile) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *MemFile) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

type MemPackage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name  string     `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	Path  string     `protobuf:"bytes,2,opt,name=Path,proto3" json:"Path,omitempty"`
	Files []*MemFile `protobuf:"bytes,3,rep,name=Files,proto3" json:"Files,omitempty"`
}

func (x *MemPackage) Reset() {
	*x = MemPackage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_std_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MemPackage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MemPackage) ProtoMessage() {}

func (x *MemPackage) ProtoReflect() protoreflect.Message {
	mi := &file_std_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MemPackage.ProtoReflect.Descriptor instead.
func (*MemPackage) Descriptor() ([]byte, []int) {
	return file_std_proto_rawDescGZIP(), []int{2}
}

func (x *MemPackage) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *MemPackage) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *MemPackage) GetFiles() []*MemFile {
	if x != nil {
		return x.Files
	}
	return nil
}

type InternalError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *InternalError) Reset() {
	*x = InternalError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_std_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InternalError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InternalError) ProtoMessage() {}

func (x *InternalError) ProtoReflect() protoreflect.Message {
	mi := &file_std_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InternalError.ProtoReflect.Descriptor instead.
func (*InternalError) Descriptor() ([]byte, []int) {
	return file_std_proto_rawDescGZIP(), []int{3}
}

type TxDecodeError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *TxDecodeError) Reset() {
	*x = TxDecodeError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_std_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TxDecodeError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TxDecodeError) ProtoMessage() {}

func (x *TxDecodeError) ProtoReflect() protoreflect.Message {
	mi := &file_std_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TxDecodeError.ProtoReflect.Descriptor instead.
func (*TxDecodeError) Descriptor() ([]byte, []int) {
	return file_std_proto_rawDescGZIP(), []int{4}
}

type InvalidSequenceError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *InvalidSequenceError) Reset() {
	*x = InvalidSequenceError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_std_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InvalidSequenceError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InvalidSequenceError) ProtoMessage() {}

func (x *InvalidSequenceError) ProtoReflect() protoreflect.Message {
	mi := &file_std_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InvalidSequenceError.ProtoReflect.Descriptor instead.
func (*InvalidSequenceError) Descriptor() ([]byte, []int) {
	return file_std_proto_rawDescGZIP(), []int{5}
}

type UnauthorizedError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UnauthorizedError) Reset() {
	*x = UnauthorizedError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_std_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnauthorizedError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnauthorizedError) ProtoMessage() {}

func (x *UnauthorizedError) ProtoReflect() protoreflect.Message {
	mi := &file_std_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnauthorizedError.ProtoReflect.Descriptor instead.
func (*UnauthorizedError) Descriptor() ([]byte, []int) {
	return file_std_proto_rawDescGZIP(), []int{6}
}

type InsufficientFundsError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *InsufficientFundsError) Reset() {
	*x = InsufficientFundsError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_std_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InsufficientFundsError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InsufficientFundsError) ProtoMessage() {}

func (x *InsufficientFundsError) ProtoReflect() protoreflect.Message {
	mi := &file_std_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InsufficientFundsError.ProtoReflect.Descriptor instead.
func (*InsufficientFundsError) Descriptor() ([]byte, []int) {
	return file_std_proto_rawDescGZIP(), []int{7}
}

type UnknownRequestError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UnknownRequestError) Reset() {
	*x = UnknownRequestError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_std_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnknownRequestError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnknownRequestError) ProtoMessage() {}

func (x *UnknownRequestError) ProtoReflect() protoreflect.Message {
	mi := &file_std_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnknownRequestError.ProtoReflect.Descriptor instead.
func (*UnknownRequestError) Descriptor() ([]byte, []int) {
	return file_std_proto_rawDescGZIP(), []int{8}
}

type InvalidAddressError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *InvalidAddressError) Reset() {
	*x = InvalidAddressError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_std_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InvalidAddressError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InvalidAddressError) ProtoMessage() {}

func (x *InvalidAddressError) ProtoReflect() protoreflect.Message {
	mi := &file_std_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InvalidAddressError.ProtoReflect.Descriptor instead.
func (*InvalidAddressError) Descriptor() ([]byte, []int) {
	return file_std_proto_rawDescGZIP(), []int{9}
}

type UnknownAddressError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UnknownAddressError) Reset() {
	*x = UnknownAddressError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_std_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnknownAddressError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnknownAddressError) ProtoMessage() {}

func (x *UnknownAddressError) ProtoReflect() protoreflect.Message {
	mi := &file_std_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnknownAddressError.ProtoReflect.Descriptor instead.
func (*UnknownAddressError) Descriptor() ([]byte, []int) {
	return file_std_proto_rawDescGZIP(), []int{10}
}

type InvalidPubKeyError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *InvalidPubKeyError) Reset() {
	*x = InvalidPubKeyError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_std_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InvalidPubKeyError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InvalidPubKeyError) ProtoMessage() {}

func (x *InvalidPubKeyError) ProtoReflect() protoreflect.Message {
	mi := &file_std_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InvalidPubKeyError.ProtoReflect.Descriptor instead.
func (*InvalidPubKeyError) Descriptor() ([]byte, []int) {
	return file_std_proto_rawDescGZIP(), []int{11}
}

type InsufficientCoinsError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *InsufficientCoinsError) Reset() {
	*x = InsufficientCoinsError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_std_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InsufficientCoinsError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InsufficientCoinsError) ProtoMessage() {}

func (x *InsufficientCoinsError) ProtoReflect() protoreflect.Message {
	mi := &file_std_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InsufficientCoinsError.ProtoReflect.Descriptor instead.
func (*InsufficientCoinsError) Descriptor() ([]byte, []int) {
	return file_std_proto_rawDescGZIP(), []int{12}
}

type InvalidCoinsError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *InvalidCoinsError) Reset() {
	*x = InvalidCoinsError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_std_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InvalidCoinsError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InvalidCoinsError) ProtoMessage() {}

func (x *InvalidCoinsError) ProtoReflect() protoreflect.Message {
	mi := &file_std_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InvalidCoinsError.ProtoReflect.Descriptor instead.
func (*InvalidCoinsError) Descriptor() ([]byte, []int) {
	return file_std_proto_rawDescGZIP(), []int{13}
}

type InvalidGasWantedError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *InvalidGasWantedError) Reset() {
	*x = InvalidGasWantedError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_std_proto_msgTypes[14]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InvalidGasWantedError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InvalidGasWantedError) ProtoMessage() {}

func (x *InvalidGasWantedError) ProtoReflect() protoreflect.Message {
	mi := &file_std_proto_msgTypes[14]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InvalidGasWantedError.ProtoReflect.Descriptor instead.
func (*InvalidGasWantedError) Descriptor() ([]byte, []int) {
	return file_std_proto_rawDescGZIP(), []int{14}
}

type OutOfGasError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *OutOfGasError) Reset() {
	*x = OutOfGasError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_std_proto_msgTypes[15]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OutOfGasError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OutOfGasError) ProtoMessage() {}

func (x *OutOfGasError) ProtoReflect() protoreflect.Message {
	mi := &file_std_proto_msgTypes[15]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OutOfGasError.ProtoReflect.Descriptor instead.
func (*OutOfGasError) Descriptor() ([]byte, []int) {
	return file_std_proto_rawDescGZIP(), []int{15}
}

type MemoTooLargeError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *MemoTooLargeError) Reset() {
	*x = MemoTooLargeError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_std_proto_msgTypes[16]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MemoTooLargeError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MemoTooLargeError) ProtoMessage() {}

func (x *MemoTooLargeError) ProtoReflect() protoreflect.Message {
	mi := &file_std_proto_msgTypes[16]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MemoTooLargeError.ProtoReflect.Descriptor instead.
func (*MemoTooLargeError) Descriptor() ([]byte, []int) {
	return file_std_proto_rawDescGZIP(), []int{16}
}

type InsufficientFeeError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *InsufficientFeeError) Reset() {
	*x = InsufficientFeeError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_std_proto_msgTypes[17]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InsufficientFeeError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InsufficientFeeError) ProtoMessage() {}

func (x *InsufficientFeeError) ProtoReflect() protoreflect.Message {
	mi := &file_std_proto_msgTypes[17]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InsufficientFeeError.ProtoReflect.Descriptor instead.
func (*InsufficientFeeError) Descriptor() ([]byte, []int) {
	return file_std_proto_rawDescGZIP(), []int{17}
}

type TooManySignaturesError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *TooManySignaturesError) Reset() {
	*x = TooManySignaturesError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_std_proto_msgTypes[18]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TooManySignaturesError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TooManySignaturesError) ProtoMessage() {}

func (x *TooManySignaturesError) ProtoReflect() protoreflect.Message {
	mi := &file_std_proto_msgTypes[18]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TooManySignaturesError.ProtoReflect.Descriptor instead.
func (*TooManySignaturesError) Descriptor() ([]byte, []int) {
	return file_std_proto_rawDescGZIP(), []int{18}
}

type NoSignaturesError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *NoSignaturesError) Reset() {
	*x = NoSignaturesError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_std_proto_msgTypes[19]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NoSignaturesError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NoSignaturesError) ProtoMessage() {}

func (x *NoSignaturesError) ProtoReflect() protoreflect.Message {
	mi := &file_std_proto_msgTypes[19]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NoSignaturesError.ProtoReflect.Descriptor instead.
func (*NoSignaturesError) Descriptor() ([]byte, []int) {
	return file_std_proto_rawDescGZIP(), []int{19}
}

type GasOverflowError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GasOverflowError) Reset() {
	*x = GasOverflowError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_std_proto_msgTypes[20]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GasOverflowError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GasOverflowError) ProtoMessage() {}

func (x *GasOverflowError) ProtoReflect() protoreflect.Message {
	mi := &file_std_proto_msgTypes[20]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GasOverflowError.ProtoReflect.Descriptor instead.
func (*GasOverflowError) Descriptor() ([]byte, []int) {
	return file_std_proto_rawDescGZIP(), []int{20}
}

var File_std_proto protoreflect.FileDescriptor

var file_std_proto_rawDesc = []byte{
	0x0a, 0x09, 0x73, 0x74, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x73, 0x74, 0x64,
	0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xad, 0x01, 0x0a, 0x0b,
	0x42, 0x61, 0x73, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x41,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x43, 0x6f, 0x69, 0x6e, 0x73, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x43, 0x6f, 0x69, 0x6e, 0x73, 0x12, 0x2c, 0x0a, 0x06, 0x50,
	0x75, 0x62, 0x4b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e,
	0x79, 0x52, 0x06, 0x50, 0x75, 0x62, 0x4b, 0x65, 0x79, 0x12, 0x24, 0x0a, 0x0d, 0x41, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x0d, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12,
	0x1a, 0x0a, 0x08, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x08, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x22, 0x31, 0x0a, 0x07, 0x4d,
	0x65, 0x6d, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x42, 0x6f,
	0x64, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x42, 0x6f, 0x64, 0x79, 0x22, 0x58,
	0x0a, 0x0a, 0x4d, 0x65, 0x6d, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x50, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x50, 0x61, 0x74, 0x68, 0x12, 0x22, 0x0a, 0x05, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x73, 0x74, 0x64, 0x2e, 0x4d, 0x65, 0x6d, 0x46, 0x69, 0x6c,
	0x65, 0x52, 0x05, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x22, 0x0f, 0x0a, 0x0d, 0x49, 0x6e, 0x74, 0x65,
	0x72, 0x6e, 0x61, 0x6c, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x0f, 0x0a, 0x0d, 0x54, 0x78, 0x44,
	0x65, 0x63, 0x6f, 0x64, 0x65, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x16, 0x0a, 0x14, 0x49, 0x6e,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x22, 0x13, 0x0a, 0x11, 0x55, 0x6e, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a,
	0x65, 0x64, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x18, 0x0a, 0x16, 0x49, 0x6e, 0x73, 0x75, 0x66,
	0x66, 0x69, 0x63, 0x69, 0x65, 0x6e, 0x74, 0x46, 0x75, 0x6e, 0x64, 0x73, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x22, 0x15, 0x0a, 0x13, 0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x15, 0x0a, 0x13, 0x49, 0x6e, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22,
	0x15, 0x0a, 0x13, 0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x14, 0x0a, 0x12, 0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x50, 0x75, 0x62, 0x4b, 0x65, 0x79, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x18, 0x0a, 0x16,
	0x49, 0x6e, 0x73, 0x75, 0x66, 0x66, 0x69, 0x63, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x69, 0x6e,
	0x73, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x13, 0x0a, 0x11, 0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x43, 0x6f, 0x69, 0x6e, 0x73, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x17, 0x0a, 0x15, 0x49,
	0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x47, 0x61, 0x73, 0x57, 0x61, 0x6e, 0x74, 0x65, 0x64, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x22, 0x0f, 0x0a, 0x0d, 0x4f, 0x75, 0x74, 0x4f, 0x66, 0x47, 0x61, 0x73,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x13, 0x0a, 0x11, 0x4d, 0x65, 0x6d, 0x6f, 0x54, 0x6f, 0x6f,
	0x4c, 0x61, 0x72, 0x67, 0x65, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x16, 0x0a, 0x14, 0x49, 0x6e,
	0x73, 0x75, 0x66, 0x66, 0x69, 0x63, 0x69, 0x65, 0x6e, 0x74, 0x46, 0x65, 0x65, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x22, 0x18, 0x0a, 0x16, 0x54, 0x6f, 0x6f, 0x4d, 0x61, 0x6e, 0x79, 0x53, 0x69, 0x67,
	0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x13, 0x0a, 0x11,
	0x4e, 0x6f, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x22, 0x12, 0x0a, 0x10, 0x47, 0x61, 0x73, 0x4f, 0x76, 0x65, 0x72, 0x66, 0x6c, 0x6f, 0x77,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x42, 0x27, 0x5a, 0x25, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6e, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2f, 0x67, 0x6e, 0x6f, 0x2f,
	0x74, 0x6d, 0x32, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x73, 0x74, 0x64, 0x2f, 0x70, 0x62, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_std_proto_rawDescOnce sync.Once
	file_std_proto_rawDescData = file_std_proto_rawDesc
)

func file_std_proto_rawDescGZIP() []byte {
	file_std_proto_rawDescOnce.Do(func() {
		file_std_proto_rawDescData = protoimpl.X.CompressGZIP(file_std_proto_rawDescData)
	})
	return file_std_proto_rawDescData
}

var file_std_proto_msgTypes = make([]protoimpl.MessageInfo, 21)
var file_std_proto_goTypes = []interface{}{
	(*BaseAccount)(nil),            // 0: std.BaseAccount
	(*MemFile)(nil),                // 1: std.MemFile
	(*MemPackage)(nil),             // 2: std.MemPackage
	(*InternalError)(nil),          // 3: std.InternalError
	(*TxDecodeError)(nil),          // 4: std.TxDecodeError
	(*InvalidSequenceError)(nil),   // 5: std.InvalidSequenceError
	(*UnauthorizedError)(nil),      // 6: std.UnauthorizedError
	(*InsufficientFundsError)(nil), // 7: std.InsufficientFundsError
	(*UnknownRequestError)(nil),    // 8: std.UnknownRequestError
	(*InvalidAddressError)(nil),    // 9: std.InvalidAddressError
	(*UnknownAddressError)(nil),    // 10: std.UnknownAddressError
	(*InvalidPubKeyError)(nil),     // 11: std.InvalidPubKeyError
	(*InsufficientCoinsError)(nil), // 12: std.InsufficientCoinsError
	(*InvalidCoinsError)(nil),      // 13: std.InvalidCoinsError
	(*InvalidGasWantedError)(nil),  // 14: std.InvalidGasWantedError
	(*OutOfGasError)(nil),          // 15: std.OutOfGasError
	(*MemoTooLargeError)(nil),      // 16: std.MemoTooLargeError
	(*InsufficientFeeError)(nil),   // 17: std.InsufficientFeeError
	(*TooManySignaturesError)(nil), // 18: std.TooManySignaturesError
	(*NoSignaturesError)(nil),      // 19: std.NoSignaturesError
	(*GasOverflowError)(nil),       // 20: std.GasOverflowError
	(*anypb.Any)(nil),              // 21: google.protobuf.Any
}
var file_std_proto_depIdxs = []int32{
	21, // 0: std.BaseAccount.PubKey:type_name -> google.protobuf.Any
	1,  // 1: std.MemPackage.Files:type_name -> std.MemFile
	2,  // [2:2] is the sub-list for method output_type
	2,  // [2:2] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_std_proto_init() }
func file_std_proto_init() {
	if File_std_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_std_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BaseAccount); i {
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
		file_std_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MemFile); i {
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
		file_std_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MemPackage); i {
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
		file_std_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InternalError); i {
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
		file_std_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TxDecodeError); i {
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
		file_std_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InvalidSequenceError); i {
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
		file_std_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnauthorizedError); i {
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
		file_std_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InsufficientFundsError); i {
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
		file_std_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnknownRequestError); i {
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
		file_std_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InvalidAddressError); i {
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
		file_std_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnknownAddressError); i {
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
		file_std_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InvalidPubKeyError); i {
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
		file_std_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InsufficientCoinsError); i {
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
		file_std_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InvalidCoinsError); i {
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
		file_std_proto_msgTypes[14].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InvalidGasWantedError); i {
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
		file_std_proto_msgTypes[15].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OutOfGasError); i {
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
		file_std_proto_msgTypes[16].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MemoTooLargeError); i {
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
		file_std_proto_msgTypes[17].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InsufficientFeeError); i {
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
		file_std_proto_msgTypes[18].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TooManySignaturesError); i {
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
		file_std_proto_msgTypes[19].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NoSignaturesError); i {
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
		file_std_proto_msgTypes[20].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GasOverflowError); i {
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
			RawDescriptor: file_std_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   21,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_std_proto_goTypes,
		DependencyIndexes: file_std_proto_depIdxs,
		MessageInfos:      file_std_proto_msgTypes,
	}.Build()
	File_std_proto = out.File
	file_std_proto_rawDesc = nil
	file_std_proto_goTypes = nil
	file_std_proto_depIdxs = nil
}
