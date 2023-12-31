// Code generated by protoc-gen-go-pulsar. DO NOT EDIT.
package gashubv1beta1

import (
	fmt "fmt"
	runtime "github.com/cosmos/cosmos-proto/runtime"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoiface "google.golang.org/protobuf/runtime/protoiface"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	io "io"
	reflect "reflect"
	sync "sync"
)

var (
	md_EventUpdateMsgGasParams              protoreflect.MessageDescriptor
	fd_EventUpdateMsgGasParams_msg_type_url protoreflect.FieldDescriptor
	fd_EventUpdateMsgGasParams_from_value   protoreflect.FieldDescriptor
	fd_EventUpdateMsgGasParams_to_value     protoreflect.FieldDescriptor
)

func init() {
	file_cosmos_gashub_v1beta1_event_proto_init()
	md_EventUpdateMsgGasParams = File_cosmos_gashub_v1beta1_event_proto.Messages().ByName("EventUpdateMsgGasParams")
	fd_EventUpdateMsgGasParams_msg_type_url = md_EventUpdateMsgGasParams.Fields().ByName("msg_type_url")
	fd_EventUpdateMsgGasParams_from_value = md_EventUpdateMsgGasParams.Fields().ByName("from_value")
	fd_EventUpdateMsgGasParams_to_value = md_EventUpdateMsgGasParams.Fields().ByName("to_value")
}

var _ protoreflect.Message = (*fastReflection_EventUpdateMsgGasParams)(nil)

type fastReflection_EventUpdateMsgGasParams EventUpdateMsgGasParams

func (x *EventUpdateMsgGasParams) ProtoReflect() protoreflect.Message {
	return (*fastReflection_EventUpdateMsgGasParams)(x)
}

func (x *EventUpdateMsgGasParams) slowProtoReflect() protoreflect.Message {
	mi := &file_cosmos_gashub_v1beta1_event_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

var _fastReflection_EventUpdateMsgGasParams_messageType fastReflection_EventUpdateMsgGasParams_messageType
var _ protoreflect.MessageType = fastReflection_EventUpdateMsgGasParams_messageType{}

type fastReflection_EventUpdateMsgGasParams_messageType struct{}

func (x fastReflection_EventUpdateMsgGasParams_messageType) Zero() protoreflect.Message {
	return (*fastReflection_EventUpdateMsgGasParams)(nil)
}
func (x fastReflection_EventUpdateMsgGasParams_messageType) New() protoreflect.Message {
	return new(fastReflection_EventUpdateMsgGasParams)
}
func (x fastReflection_EventUpdateMsgGasParams_messageType) Descriptor() protoreflect.MessageDescriptor {
	return md_EventUpdateMsgGasParams
}

// Descriptor returns message descriptor, which contains only the protobuf
// type information for the message.
func (x *fastReflection_EventUpdateMsgGasParams) Descriptor() protoreflect.MessageDescriptor {
	return md_EventUpdateMsgGasParams
}

// Type returns the message type, which encapsulates both Go and protobuf
// type information. If the Go type information is not needed,
// it is recommended that the message descriptor be used instead.
func (x *fastReflection_EventUpdateMsgGasParams) Type() protoreflect.MessageType {
	return _fastReflection_EventUpdateMsgGasParams_messageType
}

// New returns a newly allocated and mutable empty message.
func (x *fastReflection_EventUpdateMsgGasParams) New() protoreflect.Message {
	return new(fastReflection_EventUpdateMsgGasParams)
}

// Interface unwraps the message reflection interface and
// returns the underlying ProtoMessage interface.
func (x *fastReflection_EventUpdateMsgGasParams) Interface() protoreflect.ProtoMessage {
	return (*EventUpdateMsgGasParams)(x)
}

// Range iterates over every populated field in an undefined order,
// calling f for each field descriptor and value encountered.
// Range returns immediately if f returns false.
// While iterating, mutating operations may only be performed
// on the current field descriptor.
func (x *fastReflection_EventUpdateMsgGasParams) Range(f func(protoreflect.FieldDescriptor, protoreflect.Value) bool) {
	if x.MsgTypeUrl != "" {
		value := protoreflect.ValueOfString(x.MsgTypeUrl)
		if !f(fd_EventUpdateMsgGasParams_msg_type_url, value) {
			return
		}
	}
	if x.FromValue != "" {
		value := protoreflect.ValueOfString(x.FromValue)
		if !f(fd_EventUpdateMsgGasParams_from_value, value) {
			return
		}
	}
	if x.ToValue != "" {
		value := protoreflect.ValueOfString(x.ToValue)
		if !f(fd_EventUpdateMsgGasParams_to_value, value) {
			return
		}
	}
}

// Has reports whether a field is populated.
//
// Some fields have the property of nullability where it is possible to
// distinguish between the default value of a field and whether the field
// was explicitly populated with the default value. Singular message fields,
// member fields of a oneof, and proto2 scalar fields are nullable. Such
// fields are populated only if explicitly set.
//
// In other cases (aside from the nullable cases above),
// a proto3 scalar field is populated if it contains a non-zero value, and
// a repeated field is populated if it is non-empty.
func (x *fastReflection_EventUpdateMsgGasParams) Has(fd protoreflect.FieldDescriptor) bool {
	switch fd.FullName() {
	case "cosmos.gashub.v1beta1.EventUpdateMsgGasParams.msg_type_url":
		return x.MsgTypeUrl != ""
	case "cosmos.gashub.v1beta1.EventUpdateMsgGasParams.from_value":
		return x.FromValue != ""
	case "cosmos.gashub.v1beta1.EventUpdateMsgGasParams.to_value":
		return x.ToValue != ""
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: cosmos.gashub.v1beta1.EventUpdateMsgGasParams"))
		}
		panic(fmt.Errorf("message cosmos.gashub.v1beta1.EventUpdateMsgGasParams does not contain field %s", fd.FullName()))
	}
}

// Clear clears the field such that a subsequent Has call reports false.
//
// Clearing an extension field clears both the extension type and value
// associated with the given field number.
//
// Clear is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_EventUpdateMsgGasParams) Clear(fd protoreflect.FieldDescriptor) {
	switch fd.FullName() {
	case "cosmos.gashub.v1beta1.EventUpdateMsgGasParams.msg_type_url":
		x.MsgTypeUrl = ""
	case "cosmos.gashub.v1beta1.EventUpdateMsgGasParams.from_value":
		x.FromValue = ""
	case "cosmos.gashub.v1beta1.EventUpdateMsgGasParams.to_value":
		x.ToValue = ""
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: cosmos.gashub.v1beta1.EventUpdateMsgGasParams"))
		}
		panic(fmt.Errorf("message cosmos.gashub.v1beta1.EventUpdateMsgGasParams does not contain field %s", fd.FullName()))
	}
}

// Get retrieves the value for a field.
//
// For unpopulated scalars, it returns the default value, where
// the default value of a bytes scalar is guaranteed to be a copy.
// For unpopulated composite types, it returns an empty, read-only view
// of the value; to obtain a mutable reference, use Mutable.
func (x *fastReflection_EventUpdateMsgGasParams) Get(descriptor protoreflect.FieldDescriptor) protoreflect.Value {
	switch descriptor.FullName() {
	case "cosmos.gashub.v1beta1.EventUpdateMsgGasParams.msg_type_url":
		value := x.MsgTypeUrl
		return protoreflect.ValueOfString(value)
	case "cosmos.gashub.v1beta1.EventUpdateMsgGasParams.from_value":
		value := x.FromValue
		return protoreflect.ValueOfString(value)
	case "cosmos.gashub.v1beta1.EventUpdateMsgGasParams.to_value":
		value := x.ToValue
		return protoreflect.ValueOfString(value)
	default:
		if descriptor.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: cosmos.gashub.v1beta1.EventUpdateMsgGasParams"))
		}
		panic(fmt.Errorf("message cosmos.gashub.v1beta1.EventUpdateMsgGasParams does not contain field %s", descriptor.FullName()))
	}
}

// Set stores the value for a field.
//
// For a field belonging to a oneof, it implicitly clears any other field
// that may be currently set within the same oneof.
// For extension fields, it implicitly stores the provided ExtensionType.
// When setting a composite type, it is unspecified whether the stored value
// aliases the source's memory in any way. If the composite value is an
// empty, read-only value, then it panics.
//
// Set is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_EventUpdateMsgGasParams) Set(fd protoreflect.FieldDescriptor, value protoreflect.Value) {
	switch fd.FullName() {
	case "cosmos.gashub.v1beta1.EventUpdateMsgGasParams.msg_type_url":
		x.MsgTypeUrl = value.Interface().(string)
	case "cosmos.gashub.v1beta1.EventUpdateMsgGasParams.from_value":
		x.FromValue = value.Interface().(string)
	case "cosmos.gashub.v1beta1.EventUpdateMsgGasParams.to_value":
		x.ToValue = value.Interface().(string)
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: cosmos.gashub.v1beta1.EventUpdateMsgGasParams"))
		}
		panic(fmt.Errorf("message cosmos.gashub.v1beta1.EventUpdateMsgGasParams does not contain field %s", fd.FullName()))
	}
}

// Mutable returns a mutable reference to a composite type.
//
// If the field is unpopulated, it may allocate a composite value.
// For a field belonging to a oneof, it implicitly clears any other field
// that may be currently set within the same oneof.
// For extension fields, it implicitly stores the provided ExtensionType
// if not already stored.
// It panics if the field does not contain a composite type.
//
// Mutable is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_EventUpdateMsgGasParams) Mutable(fd protoreflect.FieldDescriptor) protoreflect.Value {
	switch fd.FullName() {
	case "cosmos.gashub.v1beta1.EventUpdateMsgGasParams.msg_type_url":
		panic(fmt.Errorf("field msg_type_url of message cosmos.gashub.v1beta1.EventUpdateMsgGasParams is not mutable"))
	case "cosmos.gashub.v1beta1.EventUpdateMsgGasParams.from_value":
		panic(fmt.Errorf("field from_value of message cosmos.gashub.v1beta1.EventUpdateMsgGasParams is not mutable"))
	case "cosmos.gashub.v1beta1.EventUpdateMsgGasParams.to_value":
		panic(fmt.Errorf("field to_value of message cosmos.gashub.v1beta1.EventUpdateMsgGasParams is not mutable"))
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: cosmos.gashub.v1beta1.EventUpdateMsgGasParams"))
		}
		panic(fmt.Errorf("message cosmos.gashub.v1beta1.EventUpdateMsgGasParams does not contain field %s", fd.FullName()))
	}
}

// NewField returns a new value that is assignable to the field
// for the given descriptor. For scalars, this returns the default value.
// For lists, maps, and messages, this returns a new, empty, mutable value.
func (x *fastReflection_EventUpdateMsgGasParams) NewField(fd protoreflect.FieldDescriptor) protoreflect.Value {
	switch fd.FullName() {
	case "cosmos.gashub.v1beta1.EventUpdateMsgGasParams.msg_type_url":
		return protoreflect.ValueOfString("")
	case "cosmos.gashub.v1beta1.EventUpdateMsgGasParams.from_value":
		return protoreflect.ValueOfString("")
	case "cosmos.gashub.v1beta1.EventUpdateMsgGasParams.to_value":
		return protoreflect.ValueOfString("")
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: cosmos.gashub.v1beta1.EventUpdateMsgGasParams"))
		}
		panic(fmt.Errorf("message cosmos.gashub.v1beta1.EventUpdateMsgGasParams does not contain field %s", fd.FullName()))
	}
}

// WhichOneof reports which field within the oneof is populated,
// returning nil if none are populated.
// It panics if the oneof descriptor does not belong to this message.
func (x *fastReflection_EventUpdateMsgGasParams) WhichOneof(d protoreflect.OneofDescriptor) protoreflect.FieldDescriptor {
	switch d.FullName() {
	default:
		panic(fmt.Errorf("%s is not a oneof field in cosmos.gashub.v1beta1.EventUpdateMsgGasParams", d.FullName()))
	}
	panic("unreachable")
}

// GetUnknown retrieves the entire list of unknown fields.
// The caller may only mutate the contents of the RawFields
// if the mutated bytes are stored back into the message with SetUnknown.
func (x *fastReflection_EventUpdateMsgGasParams) GetUnknown() protoreflect.RawFields {
	return x.unknownFields
}

// SetUnknown stores an entire list of unknown fields.
// The raw fields must be syntactically valid according to the wire format.
// An implementation may panic if this is not the case.
// Once stored, the caller must not mutate the content of the RawFields.
// An empty RawFields may be passed to clear the fields.
//
// SetUnknown is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_EventUpdateMsgGasParams) SetUnknown(fields protoreflect.RawFields) {
	x.unknownFields = fields
}

// IsValid reports whether the message is valid.
//
// An invalid message is an empty, read-only value.
//
// An invalid message often corresponds to a nil pointer of the concrete
// message type, but the details are implementation dependent.
// Validity is not part of the protobuf data model, and may not
// be preserved in marshaling or other operations.
func (x *fastReflection_EventUpdateMsgGasParams) IsValid() bool {
	return x != nil
}

// ProtoMethods returns optional fastReflectionFeature-path implementations of various operations.
// This method may return nil.
//
// The returned methods type is identical to
// "google.golang.org/protobuf/runtime/protoiface".Methods.
// Consult the protoiface package documentation for details.
func (x *fastReflection_EventUpdateMsgGasParams) ProtoMethods() *protoiface.Methods {
	size := func(input protoiface.SizeInput) protoiface.SizeOutput {
		x := input.Message.Interface().(*EventUpdateMsgGasParams)
		if x == nil {
			return protoiface.SizeOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Size:              0,
			}
		}
		options := runtime.SizeInputToOptions(input)
		_ = options
		var n int
		var l int
		_ = l
		l = len(x.MsgTypeUrl)
		if l > 0 {
			n += 1 + l + runtime.Sov(uint64(l))
		}
		l = len(x.FromValue)
		if l > 0 {
			n += 1 + l + runtime.Sov(uint64(l))
		}
		l = len(x.ToValue)
		if l > 0 {
			n += 1 + l + runtime.Sov(uint64(l))
		}
		if x.unknownFields != nil {
			n += len(x.unknownFields)
		}
		return protoiface.SizeOutput{
			NoUnkeyedLiterals: input.NoUnkeyedLiterals,
			Size:              n,
		}
	}

	marshal := func(input protoiface.MarshalInput) (protoiface.MarshalOutput, error) {
		x := input.Message.Interface().(*EventUpdateMsgGasParams)
		if x == nil {
			return protoiface.MarshalOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Buf:               input.Buf,
			}, nil
		}
		options := runtime.MarshalInputToOptions(input)
		_ = options
		size := options.Size(x)
		dAtA := make([]byte, size)
		i := len(dAtA)
		_ = i
		var l int
		_ = l
		if x.unknownFields != nil {
			i -= len(x.unknownFields)
			copy(dAtA[i:], x.unknownFields)
		}
		if len(x.ToValue) > 0 {
			i -= len(x.ToValue)
			copy(dAtA[i:], x.ToValue)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.ToValue)))
			i--
			dAtA[i] = 0x1a
		}
		if len(x.FromValue) > 0 {
			i -= len(x.FromValue)
			copy(dAtA[i:], x.FromValue)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.FromValue)))
			i--
			dAtA[i] = 0x12
		}
		if len(x.MsgTypeUrl) > 0 {
			i -= len(x.MsgTypeUrl)
			copy(dAtA[i:], x.MsgTypeUrl)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.MsgTypeUrl)))
			i--
			dAtA[i] = 0xa
		}
		if input.Buf != nil {
			input.Buf = append(input.Buf, dAtA...)
		} else {
			input.Buf = dAtA
		}
		return protoiface.MarshalOutput{
			NoUnkeyedLiterals: input.NoUnkeyedLiterals,
			Buf:               input.Buf,
		}, nil
	}
	unmarshal := func(input protoiface.UnmarshalInput) (protoiface.UnmarshalOutput, error) {
		x := input.Message.Interface().(*EventUpdateMsgGasParams)
		if x == nil {
			return protoiface.UnmarshalOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Flags:             input.Flags,
			}, nil
		}
		options := runtime.UnmarshalInputToOptions(input)
		_ = options
		dAtA := input.Buf
		l := len(dAtA)
		iNdEx := 0
		for iNdEx < l {
			preIndex := iNdEx
			var wire uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
				}
				if iNdEx >= l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				wire |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			fieldNum := int32(wire >> 3)
			wireType := int(wire & 0x7)
			if wireType == 4 {
				return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: EventUpdateMsgGasParams: wiretype end group for non-group")
			}
			if fieldNum <= 0 {
				return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: EventUpdateMsgGasParams: illegal tag %d (wire type %d)", fieldNum, wire)
			}
			switch fieldNum {
			case 1:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field MsgTypeUrl", wireType)
				}
				var stringLen uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					stringLen |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				intStringLen := int(stringLen)
				if intStringLen < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + intStringLen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				x.MsgTypeUrl = string(dAtA[iNdEx:postIndex])
				iNdEx = postIndex
			case 2:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field FromValue", wireType)
				}
				var stringLen uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					stringLen |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				intStringLen := int(stringLen)
				if intStringLen < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + intStringLen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				x.FromValue = string(dAtA[iNdEx:postIndex])
				iNdEx = postIndex
			case 3:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field ToValue", wireType)
				}
				var stringLen uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					stringLen |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				intStringLen := int(stringLen)
				if intStringLen < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + intStringLen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				x.ToValue = string(dAtA[iNdEx:postIndex])
				iNdEx = postIndex
			default:
				iNdEx = preIndex
				skippy, err := runtime.Skip(dAtA[iNdEx:])
				if err != nil {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, err
				}
				if (skippy < 0) || (iNdEx+skippy) < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if (iNdEx + skippy) > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				if !options.DiscardUnknown {
					x.unknownFields = append(x.unknownFields, dAtA[iNdEx:iNdEx+skippy]...)
				}
				iNdEx += skippy
			}
		}

		if iNdEx > l {
			return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
		}
		return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, nil
	}
	return &protoiface.Methods{
		NoUnkeyedLiterals: struct{}{},
		Flags:             protoiface.SupportMarshalDeterministic | protoiface.SupportUnmarshalDiscardUnknown,
		Size:              size,
		Marshal:           marshal,
		Unmarshal:         unmarshal,
		Merge:             nil,
		CheckInitialized:  nil,
	}
}

// Since: cosmos-sdk 0.43

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.0
// 	protoc        (unknown)
// source: cosmos/gashub/v1beta1/event.proto

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// EventUpdateMsgGasParams is emitted when updating a message's gas params
type EventUpdateMsgGasParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// msg_type_url is the type url of the message
	MsgTypeUrl string `protobuf:"bytes,1,opt,name=msg_type_url,json=msgTypeUrl,proto3" json:"msg_type_url,omitempty"`
	// from_value is the previous gas params
	FromValue string `protobuf:"bytes,2,opt,name=from_value,json=fromValue,proto3" json:"from_value,omitempty"`
	// to_value is the new gas params
	ToValue string `protobuf:"bytes,3,opt,name=to_value,json=toValue,proto3" json:"to_value,omitempty"`
}

func (x *EventUpdateMsgGasParams) Reset() {
	*x = EventUpdateMsgGasParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cosmos_gashub_v1beta1_event_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventUpdateMsgGasParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventUpdateMsgGasParams) ProtoMessage() {}

// Deprecated: Use EventUpdateMsgGasParams.ProtoReflect.Descriptor instead.
func (*EventUpdateMsgGasParams) Descriptor() ([]byte, []int) {
	return file_cosmos_gashub_v1beta1_event_proto_rawDescGZIP(), []int{0}
}

func (x *EventUpdateMsgGasParams) GetMsgTypeUrl() string {
	if x != nil {
		return x.MsgTypeUrl
	}
	return ""
}

func (x *EventUpdateMsgGasParams) GetFromValue() string {
	if x != nil {
		return x.FromValue
	}
	return ""
}

func (x *EventUpdateMsgGasParams) GetToValue() string {
	if x != nil {
		return x.ToValue
	}
	return ""
}

var File_cosmos_gashub_v1beta1_event_proto protoreflect.FileDescriptor

var file_cosmos_gashub_v1beta1_event_proto_rawDesc = []byte{
	0x0a, 0x21, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2f, 0x67, 0x61, 0x73, 0x68, 0x75, 0x62, 0x2f,
	0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x15, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2e, 0x67, 0x61, 0x73, 0x68,
	0x75, 0x62, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x22, 0x75, 0x0a, 0x17, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x73, 0x67, 0x47, 0x61, 0x73, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x20, 0x0a, 0x0c, 0x6d, 0x73, 0x67, 0x5f, 0x74, 0x79, 0x70,
	0x65, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6d, 0x73, 0x67,
	0x54, 0x79, 0x70, 0x65, 0x55, 0x72, 0x6c, 0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x72, 0x6f, 0x6d, 0x5f,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x72, 0x6f,
	0x6d, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x74, 0x6f, 0x5f, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x74, 0x6f, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x42, 0xd3, 0x01, 0x0a, 0x19, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73,
	0x2e, 0x67, 0x61, 0x73, 0x68, 0x75, 0x62, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x42,
	0x0a, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x34, 0x63,
	0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x73, 0x64, 0x6b, 0x2e, 0x69, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2f, 0x67, 0x61, 0x73, 0x68, 0x75, 0x62, 0x2f, 0x76, 0x31,
	0x62, 0x65, 0x74, 0x61, 0x31, 0x3b, 0x67, 0x61, 0x73, 0x68, 0x75, 0x62, 0x76, 0x31, 0x62, 0x65,
	0x74, 0x61, 0x31, 0xa2, 0x02, 0x03, 0x43, 0x47, 0x58, 0xaa, 0x02, 0x15, 0x43, 0x6f, 0x73, 0x6d,
	0x6f, 0x73, 0x2e, 0x47, 0x61, 0x73, 0x68, 0x75, 0x62, 0x2e, 0x56, 0x31, 0x62, 0x65, 0x74, 0x61,
	0x31, 0xca, 0x02, 0x15, 0x43, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x5c, 0x47, 0x61, 0x73, 0x68, 0x75,
	0x62, 0x5c, 0x56, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0xe2, 0x02, 0x21, 0x43, 0x6f, 0x73, 0x6d,
	0x6f, 0x73, 0x5c, 0x47, 0x61, 0x73, 0x68, 0x75, 0x62, 0x5c, 0x56, 0x31, 0x62, 0x65, 0x74, 0x61,
	0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x17,
	0x43, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x3a, 0x3a, 0x47, 0x61, 0x73, 0x68, 0x75, 0x62, 0x3a, 0x3a,
	0x56, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cosmos_gashub_v1beta1_event_proto_rawDescOnce sync.Once
	file_cosmos_gashub_v1beta1_event_proto_rawDescData = file_cosmos_gashub_v1beta1_event_proto_rawDesc
)

func file_cosmos_gashub_v1beta1_event_proto_rawDescGZIP() []byte {
	file_cosmos_gashub_v1beta1_event_proto_rawDescOnce.Do(func() {
		file_cosmos_gashub_v1beta1_event_proto_rawDescData = protoimpl.X.CompressGZIP(file_cosmos_gashub_v1beta1_event_proto_rawDescData)
	})
	return file_cosmos_gashub_v1beta1_event_proto_rawDescData
}

var file_cosmos_gashub_v1beta1_event_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_cosmos_gashub_v1beta1_event_proto_goTypes = []interface{}{
	(*EventUpdateMsgGasParams)(nil), // 0: cosmos.gashub.v1beta1.EventUpdateMsgGasParams
}
var file_cosmos_gashub_v1beta1_event_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_cosmos_gashub_v1beta1_event_proto_init() }
func file_cosmos_gashub_v1beta1_event_proto_init() {
	if File_cosmos_gashub_v1beta1_event_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cosmos_gashub_v1beta1_event_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventUpdateMsgGasParams); i {
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
			RawDescriptor: file_cosmos_gashub_v1beta1_event_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_cosmos_gashub_v1beta1_event_proto_goTypes,
		DependencyIndexes: file_cosmos_gashub_v1beta1_event_proto_depIdxs,
		MessageInfos:      file_cosmos_gashub_v1beta1_event_proto_msgTypes,
	}.Build()
	File_cosmos_gashub_v1beta1_event_proto = out.File
	file_cosmos_gashub_v1beta1_event_proto_rawDesc = nil
	file_cosmos_gashub_v1beta1_event_proto_goTypes = nil
	file_cosmos_gashub_v1beta1_event_proto_depIdxs = nil
}
