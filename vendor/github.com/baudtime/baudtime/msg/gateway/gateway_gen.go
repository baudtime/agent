package gateway

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/baudtime/baudtime/msg"
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *AddRequest) DecodeMsg(dc *msgp.Reader) (err error) {
	var zb0001 uint32
	zb0001, err = dc.ReadArrayHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 1 {
		err = msgp.ArrayError{Wanted: 1, Got: zb0001}
		return
	}
	var zb0002 uint32
	zb0002, err = dc.ReadArrayHeader()
	if err != nil {
		err = msgp.WrapError(err, "Series")
		return
	}
	if cap(z.Series) >= int(zb0002) {
		z.Series = (z.Series)[:zb0002]
	} else {
		z.Series = make([]*msg.Series, zb0002)
	}
	for za0001 := range z.Series {
		if dc.IsNil() {
			err = dc.ReadNil()
			if err != nil {
				err = msgp.WrapError(err, "Series", za0001)
				return
			}
			z.Series[za0001] = nil
		} else {
			if z.Series[za0001] == nil {
				z.Series[za0001] = new(msg.Series)
			}
			err = z.Series[za0001].DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "Series", za0001)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *AddRequest) EncodeMsg(en *msgp.Writer) (err error) {
	// array header, size 1
	err = en.Append(0x91)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Series)))
	if err != nil {
		err = msgp.WrapError(err, "Series")
		return
	}
	for za0001 := range z.Series {
		if z.Series[za0001] == nil {
			err = en.WriteNil()
			if err != nil {
				return
			}
		} else {
			err = z.Series[za0001].EncodeMsg(en)
			if err != nil {
				err = msgp.WrapError(err, "Series", za0001)
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *AddRequest) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// array header, size 1
	o = append(o, 0x91)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Series)))
	for za0001 := range z.Series {
		if z.Series[za0001] == nil {
			o = msgp.AppendNil(o)
		} else {
			o, err = z.Series[za0001].MarshalMsg(o)
			if err != nil {
				err = msgp.WrapError(err, "Series", za0001)
				return
			}
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *AddRequest) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	if zb0001 != 1 {
		err = msgp.ArrayError{Wanted: 1, Got: zb0001}
		return
	}
	var zb0002 uint32
	zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err, "Series")
		return
	}
	if cap(z.Series) >= int(zb0002) {
		z.Series = (z.Series)[:zb0002]
	} else {
		z.Series = make([]*msg.Series, zb0002)
	}
	for za0001 := range z.Series {
		if msgp.IsNil(bts) {
			bts, err = msgp.ReadNilBytes(bts)
			if err != nil {
				return
			}
			z.Series[za0001] = nil
		} else {
			if z.Series[za0001] == nil {
				z.Series[za0001] = new(msg.Series)
			}
			bts, err = z.Series[za0001].UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Series", za0001)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *AddRequest) Msgsize() (s int) {
	s = 1 + msgp.ArrayHeaderSize
	for za0001 := range z.Series {
		if z.Series[za0001] == nil {
			s += msgp.NilSize
		} else {
			s += z.Series[za0001].Msgsize()
		}
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *InstantQueryRequest) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "time":
			z.Time, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Time")
				return
			}
		case "timeout":
			z.Timeout, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Timeout")
				return
			}
		case "query":
			z.Query, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Query")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z InstantQueryRequest) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "time"
	err = en.Append(0x83, 0xa4, 0x74, 0x69, 0x6d, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.Time)
	if err != nil {
		err = msgp.WrapError(err, "Time")
		return
	}
	// write "timeout"
	err = en.Append(0xa7, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74)
	if err != nil {
		return
	}
	err = en.WriteString(z.Timeout)
	if err != nil {
		err = msgp.WrapError(err, "Timeout")
		return
	}
	// write "query"
	err = en.Append(0xa5, 0x71, 0x75, 0x65, 0x72, 0x79)
	if err != nil {
		return
	}
	err = en.WriteString(z.Query)
	if err != nil {
		err = msgp.WrapError(err, "Query")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z InstantQueryRequest) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "time"
	o = append(o, 0x83, 0xa4, 0x74, 0x69, 0x6d, 0x65)
	o = msgp.AppendString(o, z.Time)
	// string "timeout"
	o = append(o, 0xa7, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74)
	o = msgp.AppendString(o, z.Timeout)
	// string "query"
	o = append(o, 0xa5, 0x71, 0x75, 0x65, 0x72, 0x79)
	o = msgp.AppendString(o, z.Query)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *InstantQueryRequest) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "time":
			z.Time, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Time")
				return
			}
		case "timeout":
			z.Timeout, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Timeout")
				return
			}
		case "query":
			z.Query, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Query")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z InstantQueryRequest) Msgsize() (s int) {
	s = 1 + 5 + msgp.StringPrefixSize + len(z.Time) + 8 + msgp.StringPrefixSize + len(z.Timeout) + 6 + msgp.StringPrefixSize + len(z.Query)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *LabelValuesRequest) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "name":
			z.Name, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Name")
				return
			}
		case "constraint":
			z.Constraint, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Constraint")
				return
			}
		case "timeout":
			z.Timeout, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Timeout")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z LabelValuesRequest) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "name"
	err = en.Append(0x83, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.Name)
	if err != nil {
		err = msgp.WrapError(err, "Name")
		return
	}
	// write "constraint"
	err = en.Append(0xaa, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x74)
	if err != nil {
		return
	}
	err = en.WriteString(z.Constraint)
	if err != nil {
		err = msgp.WrapError(err, "Constraint")
		return
	}
	// write "timeout"
	err = en.Append(0xa7, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74)
	if err != nil {
		return
	}
	err = en.WriteString(z.Timeout)
	if err != nil {
		err = msgp.WrapError(err, "Timeout")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z LabelValuesRequest) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "name"
	o = append(o, 0x83, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.Name)
	// string "constraint"
	o = append(o, 0xaa, 0x63, 0x6f, 0x6e, 0x73, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x74)
	o = msgp.AppendString(o, z.Constraint)
	// string "timeout"
	o = append(o, 0xa7, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74)
	o = msgp.AppendString(o, z.Timeout)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *LabelValuesRequest) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "name":
			z.Name, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Name")
				return
			}
		case "constraint":
			z.Constraint, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Constraint")
				return
			}
		case "timeout":
			z.Timeout, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Timeout")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z LabelValuesRequest) Msgsize() (s int) {
	s = 1 + 5 + msgp.StringPrefixSize + len(z.Name) + 11 + msgp.StringPrefixSize + len(z.Constraint) + 8 + msgp.StringPrefixSize + len(z.Timeout)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *QueryResponse) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "result":
			z.Result, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Result")
				return
			}
		case "status":
			err = z.Status.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "Status")
				return
			}
		case "errorMsg":
			z.ErrorMsg, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "ErrorMsg")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *QueryResponse) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "result"
	err = en.Append(0x83, 0xa6, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74)
	if err != nil {
		return
	}
	err = en.WriteString(z.Result)
	if err != nil {
		err = msgp.WrapError(err, "Result")
		return
	}
	// write "status"
	err = en.Append(0xa6, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73)
	if err != nil {
		return
	}
	err = z.Status.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "Status")
		return
	}
	// write "errorMsg"
	err = en.Append(0xa8, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x73, 0x67)
	if err != nil {
		return
	}
	err = en.WriteString(z.ErrorMsg)
	if err != nil {
		err = msgp.WrapError(err, "ErrorMsg")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *QueryResponse) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "result"
	o = append(o, 0x83, 0xa6, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74)
	o = msgp.AppendString(o, z.Result)
	// string "status"
	o = append(o, 0xa6, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73)
	o, err = z.Status.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Status")
		return
	}
	// string "errorMsg"
	o = append(o, 0xa8, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x73, 0x67)
	o = msgp.AppendString(o, z.ErrorMsg)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *QueryResponse) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "result":
			z.Result, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Result")
				return
			}
		case "status":
			bts, err = z.Status.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Status")
				return
			}
		case "errorMsg":
			z.ErrorMsg, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "ErrorMsg")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *QueryResponse) Msgsize() (s int) {
	s = 1 + 7 + msgp.StringPrefixSize + len(z.Result) + 7 + z.Status.Msgsize() + 9 + msgp.StringPrefixSize + len(z.ErrorMsg)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *RangeQueryRequest) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "start":
			z.Start, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Start")
				return
			}
		case "end":
			z.End, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "End")
				return
			}
		case "step":
			z.Step, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Step")
				return
			}
		case "timeout":
			z.Timeout, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Timeout")
				return
			}
		case "query":
			z.Query, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Query")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *RangeQueryRequest) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 5
	// write "start"
	err = en.Append(0x85, 0xa5, 0x73, 0x74, 0x61, 0x72, 0x74)
	if err != nil {
		return
	}
	err = en.WriteString(z.Start)
	if err != nil {
		err = msgp.WrapError(err, "Start")
		return
	}
	// write "end"
	err = en.Append(0xa3, 0x65, 0x6e, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.End)
	if err != nil {
		err = msgp.WrapError(err, "End")
		return
	}
	// write "step"
	err = en.Append(0xa4, 0x73, 0x74, 0x65, 0x70)
	if err != nil {
		return
	}
	err = en.WriteString(z.Step)
	if err != nil {
		err = msgp.WrapError(err, "Step")
		return
	}
	// write "timeout"
	err = en.Append(0xa7, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74)
	if err != nil {
		return
	}
	err = en.WriteString(z.Timeout)
	if err != nil {
		err = msgp.WrapError(err, "Timeout")
		return
	}
	// write "query"
	err = en.Append(0xa5, 0x71, 0x75, 0x65, 0x72, 0x79)
	if err != nil {
		return
	}
	err = en.WriteString(z.Query)
	if err != nil {
		err = msgp.WrapError(err, "Query")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *RangeQueryRequest) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 5
	// string "start"
	o = append(o, 0x85, 0xa5, 0x73, 0x74, 0x61, 0x72, 0x74)
	o = msgp.AppendString(o, z.Start)
	// string "end"
	o = append(o, 0xa3, 0x65, 0x6e, 0x64)
	o = msgp.AppendString(o, z.End)
	// string "step"
	o = append(o, 0xa4, 0x73, 0x74, 0x65, 0x70)
	o = msgp.AppendString(o, z.Step)
	// string "timeout"
	o = append(o, 0xa7, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74)
	o = msgp.AppendString(o, z.Timeout)
	// string "query"
	o = append(o, 0xa5, 0x71, 0x75, 0x65, 0x72, 0x79)
	o = msgp.AppendString(o, z.Query)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *RangeQueryRequest) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "start":
			z.Start, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Start")
				return
			}
		case "end":
			z.End, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "End")
				return
			}
		case "step":
			z.Step, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Step")
				return
			}
		case "timeout":
			z.Timeout, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Timeout")
				return
			}
		case "query":
			z.Query, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Query")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *RangeQueryRequest) Msgsize() (s int) {
	s = 1 + 6 + msgp.StringPrefixSize + len(z.Start) + 4 + msgp.StringPrefixSize + len(z.End) + 5 + msgp.StringPrefixSize + len(z.Step) + 8 + msgp.StringPrefixSize + len(z.Timeout) + 6 + msgp.StringPrefixSize + len(z.Query)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *SeriesLabelsRequest) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "matches":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				err = msgp.WrapError(err, "Matches")
				return
			}
			if cap(z.Matches) >= int(zb0002) {
				z.Matches = (z.Matches)[:zb0002]
			} else {
				z.Matches = make([]string, zb0002)
			}
			for za0001 := range z.Matches {
				z.Matches[za0001], err = dc.ReadString()
				if err != nil {
					err = msgp.WrapError(err, "Matches", za0001)
					return
				}
			}
		case "start":
			z.Start, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Start")
				return
			}
		case "end":
			z.End, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "End")
				return
			}
		case "timeout":
			z.Timeout, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Timeout")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *SeriesLabelsRequest) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 4
	// write "matches"
	err = en.Append(0x84, 0xa7, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Matches)))
	if err != nil {
		err = msgp.WrapError(err, "Matches")
		return
	}
	for za0001 := range z.Matches {
		err = en.WriteString(z.Matches[za0001])
		if err != nil {
			err = msgp.WrapError(err, "Matches", za0001)
			return
		}
	}
	// write "start"
	err = en.Append(0xa5, 0x73, 0x74, 0x61, 0x72, 0x74)
	if err != nil {
		return
	}
	err = en.WriteString(z.Start)
	if err != nil {
		err = msgp.WrapError(err, "Start")
		return
	}
	// write "end"
	err = en.Append(0xa3, 0x65, 0x6e, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.End)
	if err != nil {
		err = msgp.WrapError(err, "End")
		return
	}
	// write "timeout"
	err = en.Append(0xa7, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74)
	if err != nil {
		return
	}
	err = en.WriteString(z.Timeout)
	if err != nil {
		err = msgp.WrapError(err, "Timeout")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *SeriesLabelsRequest) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "matches"
	o = append(o, 0x84, 0xa7, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Matches)))
	for za0001 := range z.Matches {
		o = msgp.AppendString(o, z.Matches[za0001])
	}
	// string "start"
	o = append(o, 0xa5, 0x73, 0x74, 0x61, 0x72, 0x74)
	o = msgp.AppendString(o, z.Start)
	// string "end"
	o = append(o, 0xa3, 0x65, 0x6e, 0x64)
	o = msgp.AppendString(o, z.End)
	// string "timeout"
	o = append(o, 0xa7, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74)
	o = msgp.AppendString(o, z.Timeout)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *SeriesLabelsRequest) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "matches":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Matches")
				return
			}
			if cap(z.Matches) >= int(zb0002) {
				z.Matches = (z.Matches)[:zb0002]
			} else {
				z.Matches = make([]string, zb0002)
			}
			for za0001 := range z.Matches {
				z.Matches[za0001], bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "Matches", za0001)
					return
				}
			}
		case "start":
			z.Start, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Start")
				return
			}
		case "end":
			z.End, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "End")
				return
			}
		case "timeout":
			z.Timeout, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Timeout")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *SeriesLabelsRequest) Msgsize() (s int) {
	s = 1 + 8 + msgp.ArrayHeaderSize
	for za0001 := range z.Matches {
		s += msgp.StringPrefixSize + len(z.Matches[za0001])
	}
	s += 6 + msgp.StringPrefixSize + len(z.Start) + 4 + msgp.StringPrefixSize + len(z.End) + 8 + msgp.StringPrefixSize + len(z.Timeout)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *SeriesLabelsResponse) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "labels":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				err = msgp.WrapError(err, "Labels")
				return
			}
			if cap(z.Labels) >= int(zb0002) {
				z.Labels = (z.Labels)[:zb0002]
			} else {
				z.Labels = make([][]msg.Label, zb0002)
			}
			for za0001 := range z.Labels {
				var zb0003 uint32
				zb0003, err = dc.ReadArrayHeader()
				if err != nil {
					err = msgp.WrapError(err, "Labels", za0001)
					return
				}
				if cap(z.Labels[za0001]) >= int(zb0003) {
					z.Labels[za0001] = (z.Labels[za0001])[:zb0003]
				} else {
					z.Labels[za0001] = make([]msg.Label, zb0003)
				}
				for za0002 := range z.Labels[za0001] {
					err = z.Labels[za0001][za0002].DecodeMsg(dc)
					if err != nil {
						err = msgp.WrapError(err, "Labels", za0001, za0002)
						return
					}
				}
			}
		case "status":
			err = z.Status.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "Status")
				return
			}
		case "errorMsg":
			z.ErrorMsg, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "ErrorMsg")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *SeriesLabelsResponse) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "labels"
	err = en.Append(0x83, 0xa6, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Labels)))
	if err != nil {
		err = msgp.WrapError(err, "Labels")
		return
	}
	for za0001 := range z.Labels {
		err = en.WriteArrayHeader(uint32(len(z.Labels[za0001])))
		if err != nil {
			err = msgp.WrapError(err, "Labels", za0001)
			return
		}
		for za0002 := range z.Labels[za0001] {
			err = z.Labels[za0001][za0002].EncodeMsg(en)
			if err != nil {
				err = msgp.WrapError(err, "Labels", za0001, za0002)
				return
			}
		}
	}
	// write "status"
	err = en.Append(0xa6, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73)
	if err != nil {
		return
	}
	err = z.Status.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "Status")
		return
	}
	// write "errorMsg"
	err = en.Append(0xa8, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x73, 0x67)
	if err != nil {
		return
	}
	err = en.WriteString(z.ErrorMsg)
	if err != nil {
		err = msgp.WrapError(err, "ErrorMsg")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *SeriesLabelsResponse) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "labels"
	o = append(o, 0x83, 0xa6, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Labels)))
	for za0001 := range z.Labels {
		o = msgp.AppendArrayHeader(o, uint32(len(z.Labels[za0001])))
		for za0002 := range z.Labels[za0001] {
			o, err = z.Labels[za0001][za0002].MarshalMsg(o)
			if err != nil {
				err = msgp.WrapError(err, "Labels", za0001, za0002)
				return
			}
		}
	}
	// string "status"
	o = append(o, 0xa6, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73)
	o, err = z.Status.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Status")
		return
	}
	// string "errorMsg"
	o = append(o, 0xa8, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x73, 0x67)
	o = msgp.AppendString(o, z.ErrorMsg)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *SeriesLabelsResponse) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "labels":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Labels")
				return
			}
			if cap(z.Labels) >= int(zb0002) {
				z.Labels = (z.Labels)[:zb0002]
			} else {
				z.Labels = make([][]msg.Label, zb0002)
			}
			for za0001 := range z.Labels {
				var zb0003 uint32
				zb0003, bts, err = msgp.ReadArrayHeaderBytes(bts)
				if err != nil {
					err = msgp.WrapError(err, "Labels", za0001)
					return
				}
				if cap(z.Labels[za0001]) >= int(zb0003) {
					z.Labels[za0001] = (z.Labels[za0001])[:zb0003]
				} else {
					z.Labels[za0001] = make([]msg.Label, zb0003)
				}
				for za0002 := range z.Labels[za0001] {
					bts, err = z.Labels[za0001][za0002].UnmarshalMsg(bts)
					if err != nil {
						err = msgp.WrapError(err, "Labels", za0001, za0002)
						return
					}
				}
			}
		case "status":
			bts, err = z.Status.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Status")
				return
			}
		case "errorMsg":
			z.ErrorMsg, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "ErrorMsg")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *SeriesLabelsResponse) Msgsize() (s int) {
	s = 1 + 7 + msgp.ArrayHeaderSize
	for za0001 := range z.Labels {
		s += msgp.ArrayHeaderSize
		for za0002 := range z.Labels[za0001] {
			s += z.Labels[za0001][za0002].Msgsize()
		}
	}
	s += 7 + z.Status.Msgsize() + 9 + msgp.StringPrefixSize + len(z.ErrorMsg)
	return
}
