package ipmi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type marshalable interface {
	Marshal() ([]byte, error)
}

// MessageClass identifies the type of message in its lower five bits, and whether the message is an
// ACK in the most significant bit.
type MessageClass uint8

// IsACK checks the top bit of the message class; if it is set, returns true.
func (c MessageClass) IsACK() bool {
	return (c & 0x80) != 0
}

// Reserved returns the value of the two reserved bits.
func (c MessageClass) Reserved() uint8 {
	return uint8((c >> 5) & 0x03)
}

// Class returns the numeric message class value.
func (c MessageClass) Class() uint8 {
	return uint8(c & 0x1f)
}

// SameClass checks if this message has the same class as the other, ignoring ACK/reserved.
func (c MessageClass) SameClass(other MessageClass) bool {
	return c.Class() == other.Class()
}

// String returns a log-friendly representation of the message class.
func (c MessageClass) String() string {
	prefix := ""
	if c.IsACK() {
		prefix = " ACK;"
	}
	if r := c.Reserved(); r != 0 {
		prefix += fmt.Sprintf(" Reserved=%01x;", r)
	}
	return fmt.Sprintf("%02x:%s %s", uint8(c), prefix, c.Name())
}

// MarshalJSON returns the search-friendly JSON representation of the MessageClass.
func (c MessageClass) MarshalJSON() ([]byte, error) {
	type temp struct {
		Raw      uint8  `json:"raw"`
		ACK      bool   `json:"is_ack,omitempty"`
		Name     string `json:"name"`
		Reserved uint8  `json:"reserved,omitempty"`
		Class    uint8  `json:"class"`
	}
	return json.Marshal(&temp{
		Raw:      uint8(c),
		ACK:      c.IsACK(),
		Name:     c.Name(),
		Reserved: c.Reserved(),
		Class:    c.Class(),
	})
}

// IPMIAuthType provides an interface to access the individual fields packed into the IPMI
// authentication type information byte defined in section 13.6, table 13, and utilized further in
// the channel auth capabilities response of 22.13's table 22, where the values here number the bits
// there.
type IPMIAuthType uint8

// Reserved returns the four reserved bits.
func (v IPMIAuthType) Reserved() uint8 {
	return uint8(v >> 4)
}

// Get the lower four bits, which correspond to the actual authentication type.
func (v IPMIAuthType) Get() uint8 {
	return uint8(v & 0x0f)
}

// String returns a log-friendly string representation of the type.
func (v IPMIAuthType) String() string {
	return fmt.Sprintf("%02x: %s", uint8(v), v.Name())
}

// Validate that the auth type is good (i.e. that it doesn't have an unrecognized, but non-reserved
// type).
func (v IPMIAuthType) Validate() error {
	// strict validation could check that Reserved == 0, but for forward compatibility, we'll ignore that.
	_, ok := ipmiAuthTypeNames[v.Get()]
	if !ok {
		return fmt.Errorf("0x%02x is not a valid IPMI auth type", v)
	}
	return nil
}

// MarshalJSON returns the search-friendly JSON representation of the auth type.
func (v IPMIAuthType) MarshalJSON() ([]byte, error) {
	type temp struct {
		Raw      uint8  `json:"raw"`
		Type     uint8  `json:"type"`
		Name     string `json:"name"`
		Reserved uint8  `json:"reserved,omitempty"`
	}
	return json.Marshal(&temp{
		Raw:      uint8(v),
		Type:     v.Get(),
		Name:     v.Name(),
		Reserved: v.Reserved(),
	})
}

// IPMIPayloadType represents the type of a payload.
type IPMIPayloadType uint8

// Encrypted returns true if the high bit is set (see e.g. section 13.6).
func (v IPMIPayloadType) Encrypted() bool {
	return v&0x80 != 0
}

// Authenticated returns true if bit 6 is set (see e.g. section 13.6).
func (v IPMIPayloadType) Authenticated() bool {
	return v&0x40 != 0
}

// Get just the payload type identifier (mask off the top two bits).
func (v IPMIPayloadType) Get() uint8 {
	return uint8(v & 0x3F)
}

// String returns a log-friendly representation of the payload type.
func (v IPMIPayloadType) String() string {
	prefix := ""
	if v.Encrypted() {
		prefix = " encrypted;"
	}
	if v.Authenticated() {
		prefix += " authenticated;"
	}
	return fmt.Sprintf("%02x:%s %s", uint8(v), prefix, v.Name())
}

// LogicalUnitNumber is a two-bit field identifying the target of a request, usually packed in with
// the NetFn (stored here in a PackedNetFn).
type LogicalUnitNumber uint8

// LogicalUnitNumberJSON is the search-friendly version of LogicalUnitNumber
type LogicalUnitNumberJSON struct {
	Raw  uint8  `json:"raw"`
	Name string `json:"name,omitempty"`
}

// MarshalJSON returns the search-friendly JSON encoding of the LUN.
func (v LogicalUnitNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(&LogicalUnitNumberJSON{
		Raw:  uint8(v),
		Name: v.Name(),
	})
}

// NetFn is a 6-bit integer identifying a class of functions, with the least significant bit
// distinguishing responses (set) from requests (set).
type NetFn uint8

// Pack a LUN with this NetFn to get a PackedNetFn.
func (v NetFn) Pack(lun LogicalUnitNumber) PackedNetFn {
	return PackedNetFn(uint8(v)<<2 | uint8(lun))
}

// IsRequest returns true if this NetFn represents a request.
func (v NetFn) IsRequest() bool {
	return v&1 == 0
}

// IsResponse returns true if this NetFn represents a response.
func (v NetFn) IsResponse() bool {
	return !v.IsRequest()
}

// Get the standardized int value of the NetFn, with the least significant bit masked off (so
// NetFn.Get() is always even).
func (v NetFn) Get() uint8 {
	return uint8(v & 0xFE)
}

// NetFnJSON is the search-friendly JSON format of the NetFn.
type NetFnJSON struct {
	Raw        uint8  `json:"raw"`
	Value      uint8  `json:"value"`
	IsRequest  bool   `json:"is_request,omitempty"`
	IsResponse bool   `json:"is_response,omitempty"`
	Name       string `json:"name,omitempty"`
}

// MarshalJSON returns the search-friendly JSON representation of the NetFn.
func (v NetFn) MarshalJSON() ([]byte, error) {
	return json.Marshal(&NetFnJSON{
		Raw:        uint8(v),
		Value:      v.Get(),
		IsRequest:  v.IsRequest(),
		IsResponse: v.IsResponse(),
		Name:       v.Name(),
	})
}

// PackedNetFn represents a combined NetFn and LUN.
type PackedNetFn uint8

// NetFn gets the NetFn portion (the upper 6 bits).
func (v PackedNetFn) NetFn() NetFn {
	return NetFn(v >> 2)
}

// LUN gets the logical unit number portion (the lower two bits).
func (v PackedNetFn) LUN() LogicalUnitNumber {
	return LogicalUnitNumber(v & 0x03)
}

// SetNetFn part of the byte.
func (v *PackedNetFn) SetNetFn(netFn NetFn) *PackedNetFn {
	*v = PackedNetFn(uint8(*v&0x03) | (uint8(netFn) << 2))
	return v
}

// SetLUN part of the byte.
func (v *PackedNetFn) SetLUN(lun LogicalUnitNumber) *PackedNetFn {
	*v = PackedNetFn(uint8(*v&0xFC) | uint8(lun))
	return v
}

// PackedNetFnJSON is the search-friendly JSON format of the PackedNetFn.
type PackedNetFnJSON struct {
	Raw   uint8             `json:"raw"`
	NetFn NetFn             `json:"net_fn"`
	LUN   LogicalUnitNumber `json:"lun"`
}

// MarshalJSON returns the search-friendly JSON representation of the packed NetFn + LUN.
func (v PackedNetFn) MarshalJSON() ([]byte, error) {
	return json.Marshal(&PackedNetFnJSON{
		Raw:   uint8(v),
		NetFn: v.NetFn(),
		LUN:   v.LUN(),
	})
}

// RMCPHeader, defined in figure section 13.1.3.
type RMCPHeader struct {
	Version        byte         `json:"version"`
	Reserved       byte         `json:"reserved,omitempty"`
	SequenceNumber byte         `json:"sequence_number"`
	MessageClass   MessageClass `json:"message_class"`
}

// Write the RMCPHeader to a Writer, return the number of bytes written and/or any error.
func (header *RMCPHeader) Write(dst io.Writer) (int, error) {
	return writeAndCount(dst, header)
}

// Read this RMCPHeader's contents from a Reader, and return the number of bytes read and/or any
// error.
func (header *RMCPHeader) Read(src io.Reader) (int, error) {
	return readObject(src, header, 4)
}

// IPMICommandNumber identifies a command within the NetFn. These are defined in appendix G.
type IPMICommandNumber uint8

// IPMICommandNumberJSON is the search-friendly JSON format of the IPMICommandNumber.
// NOTE: IPMICommandNumber has no MarshalJSON(), since its name depends on the NetFn number.
type IPMICommandNumberJSON struct {
	Raw  uint8  `json:"raw"`
	Name string `json:"name,omitempty"`
}

// IPMICommandPayload is defined in section 13.8.
type IPMICommandPayload struct {
	// Responder address.
	RsAddr IPMIAddressByte `json:"rs_addr"`

	// Network function code with LUN.
	NetFn PackedNetFn `json:"net_fn"`

	// Chk1 = -((RsAddr + NetFn) & 0xff)
	Chk1 uint8 `json:"chk1"`

	// Requestor address.
	RqAddr IPMIAddressByte `json:"rq_addr"`

	// Requestor sequence number.
	RqSeq uint8 `json:"rq_seq"`

	// Command byte.
	// Note: the JSON field is overridden, since the name field requires the NetFn to be available.
	Cmd IPMICommandNumber `json:"cmd"`

	// Data to be sent.
	Data []byte `json:"data,omitempty"`

	// Chk2 = -((RqAddr + RqSeq + Cmd + [data]) & 0xff)
	Chk2 uint8 `json:"chk2"`
}

// CalculateChecksums returns the expected values for chk1, chk2, given the rest of the payload.
func (payload *IPMICommandPayload) CalculateChecksums() (uint8, uint8) {
	chk1 := -(uint8(payload.RsAddr) + uint8(payload.NetFn))
	temp := uint8(payload.RqAddr) + uint8(payload.RqSeq) + uint8(payload.Cmd)
	for _, v := range payload.Data {
		temp += uint8(v)
	}
	chk2 := -temp
	return chk1, chk2
}

// Validate that the checksums have their expected values.
func (payload *IPMICommandPayload) Validate() error {
	chk1, chk2 := payload.CalculateChecksums()
	if chk1 == payload.Chk1 && chk2 == payload.Chk2 {
		return nil
	}
	msg := "Validation error"
	if chk1 != payload.Chk1 {
		msg += fmt.Sprintf("; chk1 expected 0x%02x, got 0x%02x", chk1, payload.Chk1)
	}
	if chk2 != payload.Chk2 {
		msg += fmt.Sprintf("; chk2 expected 0x%02x, got 0x%02x", chk2, payload.Chk2)
	}
	return errors.New(msg)
}

// SetChecksums calculates and sets the checksums, given the present state of the payload.
func (payload *IPMICommandPayload) SetChecksums() *IPMICommandPayload {
	payload.Chk1, payload.Chk2 = payload.CalculateChecksums()
	return payload
}

// IPMICommandPayloadAlias is an alias for IPMICommandPayload, used to avoid circular references.
type IPMICommandPayloadAlias IPMICommandPayload

// IPMICommandPayloadJSON is the search-friendly JSON representation of the IPMI command payload.
type IPMICommandPayloadJSON struct {
	CmdOverride   *IPMICommandNumberJSON `json:"cmd"`
	ChecksumError bool                   `json:"checksum_error,omitempty"`
	*IPMICommandPayloadAlias
}

// MarshalJSON returns the search-friendly representation of the IPMI command payload.
func (payload *IPMICommandPayload) MarshalJSON() ([]byte, error) {
	return json.Marshal(&IPMICommandPayloadJSON{
		CmdOverride: &IPMICommandNumberJSON{
			Raw:  uint8(payload.Cmd),
			Name: payload.Cmd.Name(payload.NetFn.NetFn()),
		},
		ChecksumError:           payload.Validate() != nil,
		IPMICommandPayloadAlias: (*IPMICommandPayloadAlias)(payload),
	})
}

// Marshal the payload into a byte slice.
func (payload *IPMICommandPayload) Marshal() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := write(buf, []byte{
		byte(payload.RsAddr),
		byte(payload.NetFn),
		byte(payload.Chk1),
		byte(payload.RqAddr),
		byte(payload.RqSeq),
		byte(payload.Cmd),
	}); err != nil {
		return buf.Bytes(), err
	}
	if _, err := write(buf, payload.Data); err != nil {
		return buf.Bytes(), err
	}
	if _, err := write(buf, []byte{payload.Chk2}); err != nil {
		return buf.Bytes(), err
	}
	return buf.Bytes(), nil
}

// Unmarshal a byte slice into this payload.
func (payload *IPMICommandPayload) Unmarshal(data []byte) error {
	reader := bytes.NewReader(data)
	if _, err := readObject(reader, &payload.RsAddr, 1); err != nil {
		return err
	}
	if _, err := readObject(reader, &payload.NetFn, 1); err != nil {
		return err
	}
	if _, err := readObject(reader, &payload.Chk1, 1); err != nil {
		return err
	}
	if _, err := readObject(reader, &payload.RqAddr, 1); err != nil {
		return err
	}
	if _, err := readObject(reader, &payload.RqSeq, 1); err != nil {
		return err
	}
	if _, err := readObject(reader, &payload.Cmd, 1); err != nil {
		return err
	}
	if len(data) < 7 {
		return io.ErrShortBuffer
	}
	payload.Data = make([]byte, len(data)-7)
	copy(payload.Data, data[6:len(data)-1])
	payload.Chk2 = data[len(data)-1]
	return nil
}

// SetData and update the checksums.
func (payload *IPMICommandPayload) SetData(data marshalable) error {
	body, err := data.Marshal()
	if err != nil {
		return err
	}
	payload.Data = body
	payload.SetChecksums()
	return err
}

// Write the payload to the writer. Caller is responsible for invoking SetChecksums().
func (payload *IPMICommandPayload) Write(writer io.Writer) (int, error) {
	temp, err := payload.Marshal()
	if err != nil {
		return 0, err
	}
	return write(writer, temp)
}

// GetChannelAuthenticationCapabilitiesRequest is defined in section 22.13.
// This is used at the start of establishing connections, in cleartext prior to any authentication,
// and IP-based channels must support it -- so it is ideal for scanning purposes.
type GetChannelAuthenticationCapabilitiesRequest struct {
	ChannelNumber         uint8
	ChannelPrivilegeLevel uint8
}

// ChannelReserved gets the value of the reserved bits (6:4) in the ChannelNumber.
func (v *GetChannelAuthenticationCapabilitiesRequest) ChannelReserved() uint8 {
	return (v.ChannelNumber >> 4) & 0x07
}

// PrivilegeReserved gets the value of the reserved bits (7:4) in the ChannelPrivilegeLevel.
func (v *GetChannelAuthenticationCapabilitiesRequest) PrivilegeReserved() uint8 {
	return v.ChannelPrivilegeLevel >> 4
}

// Extended returns true if the top bit of the channel number is set, which indicates that IPMIv2.0+
// extended data should be returned.
func (v *GetChannelAuthenticationCapabilitiesRequest) Extended() bool {
	return v.ChannelNumber&0x80 == 0x80
}

// Channel returns the channel identifier. The special value 0x0E means "retrieve information for
// channel this request was issued on".
func (v *GetChannelAuthenticationCapabilitiesRequest) Channel() uint8 {
	return v.ChannelNumber & 0x0F
}

// PrivilegeLevel returns the requested privilege level (bits 3:0).
func (v *GetChannelAuthenticationCapabilitiesRequest) PrivilegeLevel() ChannelPrivilegeLevel {
	return ChannelPrivilegeLevel(v.ChannelPrivilegeLevel & 0x0F)
}

// SetChannel (the low four bits) on the request.
func (v *GetChannelAuthenticationCapabilitiesRequest) SetChannel(ch uint8) *GetChannelAuthenticationCapabilitiesRequest {
	v.ChannelPrivilegeLevel = (v.ChannelPrivilegeLevel & 0xf0) | (ch & 0x0f)
	return v
}

// SetExtended on the request (either set or clear the high bit).
func (v *GetChannelAuthenticationCapabilitiesRequest) SetExtended(ex bool) *GetChannelAuthenticationCapabilitiesRequest {
	if ex {
		v.ChannelPrivilegeLevel |= 0x80
	} else {
		v.ChannelPrivilegeLevel &= 0x7F
	}
	return v
}

// SetPrivilegeLevel on the request (the lower four bits)
func (v *GetChannelAuthenticationCapabilitiesRequest) SetPrivilegeLevel(lvl ChannelPrivilegeLevel) *GetChannelAuthenticationCapabilitiesRequest {
	v.ChannelPrivilegeLevel = (v.ChannelPrivilegeLevel & 0xf0) | uint8(lvl&0x0f)
	return v
}

// Set the extended bit, the channel, and the privilege level on the request.
func (v *GetChannelAuthenticationCapabilitiesRequest) Set(ex bool, ch uint8, lvl ChannelPrivilegeLevel) *GetChannelAuthenticationCapabilitiesRequest {
	b := uint8(0x00)
	if ex {
		b = 0x80
	}
	v.ChannelNumber = b | (ch & 0x0f)
	v.ChannelPrivilegeLevel = uint8(lvl & 0x0f)
	return v
}

// Marshal the request into a byte slice.
func (v *GetChannelAuthenticationCapabilitiesRequest) Marshal() ([]byte, error) {
	return []byte{
		v.ChannelNumber, v.ChannelPrivilegeLevel,
	}, nil
}

// Unmarshal a byte slice into this request.
func (v *GetChannelAuthenticationCapabilitiesRequest) Unmarshal(data []byte) error {
	if len(data) != 2 {
		return fmt.Errorf("bad data: GetChannelAuthenticationCapabilitiesRequest should be 2 bytes, got %d bytes", len(data))
	}
	v.ChannelNumber = data[0]
	v.ChannelPrivilegeLevel = data[1]
	return nil
}

// CompletionCode is defined in section 5.2, and is the first byte of the response to a command.
// Errors outside of the command level are handled with a different error reporting mechanism.
type CompletionCode uint8

// CompletionCodeJSON is the search-friendly version of CompletionCode.
type CompletionCodeJSON struct {
	Raw  uint8  `json:"raw"`
	Name string `json:"name,omitempty"`
}

// MarshalJSON returns the search-friendly JSON encoding of CompletionCodeJSON
func (v CompletionCode) MarshalJSON() ([]byte, error) {
	return json.Marshal(&CompletionCodeJSON{
		Raw:  uint8(v),
		Name: v.Name(),
	})
}

// ChannelPrivilegeLevel identifies the requested level of privilege for the session.
type ChannelPrivilegeLevel uint8

// SupportedAuthTypes is a bit mask of supported AuthTypes defined in e.g. the second part of table
// 22 in section 22.13. The bit indices (incidentally?) coincide with the scalar values from the
// AuthType enum (e.g. AuthTypeNone = 0, SupportedAuthTypesNone = 2^0), though bit 6 is not defined.
type SupportedAuthTypes uint8

// Extended indicates that extended auth options are used.
func (v SupportedAuthTypes) Extended() bool {
	return v&0x80 == 0x80
}

// Reserved returns the value of the reserved bits (bit 3 and bit 6) as a mask.
func (v SupportedAuthTypes) Reserved() uint8 {
	return uint8(v & 0x48)
}

// Supports checks of the given IPMIAuthType is supported by this mask.
func (v SupportedAuthTypes) Supports(which IPMIAuthType) bool {
	iType := which.Get()
	if iType > 5 {
		// only 5 are supported
		return false
	}
	return v&(1<<iType) != 0
}

// Get this supported types mask as a map[IPMIAuthType]bool.
func (v SupportedAuthTypes) Get() map[IPMIAuthType]bool {
	types := []IPMIAuthType{AuthTypeNone, AuthTypeMD2, AuthTypeMD5, AuthTypePassword, AuthTypeOEMProprietary}
	ret := make(map[IPMIAuthType]bool)
	for _, t := range types {
		ret[t] = v.Supports(t)
	}
	return ret
}

// SupportedAuthTypesJSON is the search-friendly representation of the SupportedAuthTypes.
type SupportedAuthTypesJSON struct {
	Raw            uint8 `json:"raw"`
	Extended       bool  `json:"extended,omitempty"`
	Reserved       uint8 `json:"reserved,omitempty"`
	None           bool  `json:"none,omitempty"`
	MD2            bool  `json:"md2,omitempty"`
	MD5            bool  `json:"md5,omitempty"`
	Password       bool  `json:"password,omitempty"`
	OEMProprietary bool  `json:"oem_proprietary,omitempty"`
}

// MarshalJSON returns the search-friendly representation of the SupportedAuthTypes.
func (v SupportedAuthTypes) MarshalJSON() ([]byte, error) {
	return json.Marshal(&SupportedAuthTypesJSON{
		Raw:            uint8(v),
		Extended:       v.Extended(),
		Reserved:       v.Reserved(),
		None:           v.Supports(AuthTypeNone),
		MD2:            v.Supports(AuthTypeMD2),
		MD5:            v.Supports(AuthTypeMD5),
		Password:       v.Supports(AuthTypePassword),
		OEMProprietary: v.Supports(AuthTypeOEMProprietary),
	})
}

// AuthStatus is defined in section 22.13.
type AuthStatus uint8

// Reserved is the value of the upper two reserved bits.
func (v AuthStatus) Reserved() uint8 {
	return uint8(v & 0xC0)
}

// KG returns whether the K_G bit (bit #5) is set to a non-zero value.
func (v AuthStatus) KG() bool {
	return v&0x20 != 0
}

// AuthEachMessage returns true if the per-message authentication bit (#4) is set (v1.5/v2.0 only).
func (v AuthStatus) AuthEachMessage() bool {
	return v&0x10 != 0
}

// UserAuthDisabled returns true if the user-level authentication status bit (#3) is set (and so
// user-level authentication is disabled).
func (v AuthStatus) UserAuthDisabled() bool {
	return v&0x08 != 0
}

// HasNamedUsers returns true if the "non-null usernames enabled" bit (#2) is set.
func (v AuthStatus) HasNamedUsers() bool {
	return v&0x04 != 0
}

// HasAnonymousUsers returns true if the "null usernames enabled" bit (#1) is set.
func (v AuthStatus) HasAnonymousUsers() bool {
	return v&0x02 != 0
}

// AnonymousLoginEnabled returns true if the "anonymous login enabled" bit (#0) is set.
func (v AuthStatus) AnonymousLoginEnabled() bool {
	return v&0x01 != 0
}

// AuthStatusJSON represents the search-friendly representation of the AuthStatus.
type AuthStatusJSON struct {
	Reserved              uint8 `json:"reserved,omitempty"`
	KG                    bool  `json:"two_key_login_required,omitempty"`
	AuthEachMessage       bool  `json:"auth_each_message,omitempty"`
	UserAuthDisabled      bool  `json:"user_auth_disabled,omitempty"`
	HasNamedUsers         bool  `json:"has_named_users,omitempty"`
	HasAnonymousUsers     bool  `json:"has_anonymous_users,omitempty"`
	AnonymousLoginEnabled bool  `json:"anonymous_login_enabled,omitempty"`
}

// MarshalJSON returns the search-friendly representation of the AuthStatus.
func (v AuthStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(&AuthStatusJSON{
		Reserved:              v.Reserved(),
		KG:                    v.KG(),
		AuthEachMessage:       v.AuthEachMessage(),
		UserAuthDisabled:      v.UserAuthDisabled(),
		HasNamedUsers:         v.HasNamedUsers(),
		HasAnonymousUsers:     v.HasAnonymousUsers(),
		AnonymousLoginEnabled: v.AnonymousLoginEnabled(),
	})
}

// ExtendedCapabilities are defined in table 22 of section 22.13, and are defined only in IPMIv2.0+.
type ExtendedCapabilities uint8

// Reserved returns the reserved bits (all but the bottom two).
func (v ExtendedCapabilities) Reserved() uint8 {
	return uint8(v & 0xFC)
}

// SupportsIPMIv2_0 returns true if the 2.0 bit (#1) is set.
func (v ExtendedCapabilities) SupportsIPMIv2_0() bool {
	return v&0x02 != 0
}

// SupportsIPMIv1_5 returns true if the 1.5 bit (#1) is set.
func (v ExtendedCapabilities) SupportsIPMIv1_5() bool {
	return v&0x01 != 0
}

// ExtendedCapabilitiesJSON represents the search-friendly version of ExtendedCapabilities.
type ExtendedCapabilitiesJSON struct {
	Reserved         uint8 `json:"reserved,omitempty"`
	SupportsIPMIv2_0 bool  `json:"supports_ipmi_v2_0"`
	SupportsIPMIv1_5 bool  `json:"supports_ipmi_v1_5"`
}

// MarshalJSON returns the search-friendly representation of the ExtendedCapabilities.
func (v ExtendedCapabilities) MarshalJSON() ([]byte, error) {
	return json.Marshal(&ExtendedCapabilitiesJSON{
		Reserved:         v.Reserved(),
		SupportsIPMIv1_5: v.SupportsIPMIv1_5(),
		SupportsIPMIv2_0: v.SupportsIPMIv2_0(),
	})
}

// GetChannelAuthenticationCapabilitiesResponse is defined in the second part of table 22 in
// section 22.13.
type GetChannelAuthenticationCapabilitiesResponse struct {
	CompletionCode       CompletionCode       `json:"completion_code"`
	ChannelNumber        uint8                `json:"channel_number"`
	SupportedAuthTypes   SupportedAuthTypes   `json:"supported_auth_types,omitempty"`
	AuthStatus           AuthStatus           `json:"auth_status,omitempty"`
	ExtendedCapabilities ExtendedCapabilities `json:"extended_capabilities,omitempty"`
	OEM_ID               [3]byte              `json:"oem_id,omitempty"`
	OEMData              uint8                `json:"oem_data,omitempty"`
}

// Read the payload's contents from the given Reader, return the number of bytes actually read
// and/or the error.
func (v *GetChannelAuthenticationCapabilitiesResponse) Read(src io.Reader) (int, error) {
	return readObject(src, v, 9)
}

// Unmarshal the data into this payload.
func (v *GetChannelAuthenticationCapabilitiesResponse) Unmarshal(data []byte) error {
	n, err := v.Read(bytes.NewReader(data))
	if err != nil {
		return err
	}
	if n < len(data) {
		return io.ErrShortBuffer
	}
	return nil
}

// MarshalJSON returns the search-friendly representation of the capabilities.
func (v *GetChannelAuthenticationCapabilitiesResponse) MarshalJSON() ([]byte, error) {
	// 	Hack around infinite MarshalJSON loop by aliasing parent type (http://choly.ca/post/go-json-marshalling/)
	type Alias GetChannelAuthenticationCapabilitiesResponse
	type temp struct {
		OEMIDOverride []byte `json:"oem_id,omitempty"`
		*Alias
	}
	ret := &temp{
		OEMIDOverride: v.OEM_ID[:],
		Alias:         (*Alias)(v),
	}
	if v.OEM_ID == [3]byte{0, 0, 0} {
		ret.OEMIDOverride = nil
	}
	return json.Marshal(ret)
}

// IPMISessionHeader_v1_5 is a header for a v1.5 IPMI packet.
type IPMISessionHeader_v1_5 struct {
	AuthType              IPMIAuthType `json:"auth_type"`
	SessionSequenceNumber uint32       `json:"sequence_number"`
	SessionID             uint32       `json:"session_id,omitempty"`
	// Absent if AuthType == None
	AuthCode [16]byte `json:"auth_code,omitempty"`
	Length   uint8    `json:"length,omitempty"`
}

var emptyAuthCode [16]byte

// MarshalJSON returns the search-friendly representation of the header.
func (header *IPMISessionHeader_v1_5) MarshalJSON() ([]byte, error) {
	// 	Hack around infinite MarshalJSON loop by aliasing parent type (http://choly.ca/post/go-json-marshalling/)
	type Alias IPMISessionHeader_v1_5
	type temp_noAuth struct {
		AuthCodeOmitted []byte `json:"auth_code,omitempty"`
		*Alias
	}
	type temp_auth struct {
		AuthCodeOverride []byte `json:"auth_code"`
		*Alias
	}
	if header.AuthType == AuthTypeNone && header.AuthCode == emptyAuthCode {
		return json.Marshal(&temp_noAuth{
			Alias: (*Alias)(header),
		})
	}
	return json.Marshal(&temp_auth{
		AuthCodeOverride: header.AuthCode[:],
		Alias:            (*Alias)(header),
	})
}

// Unmarshal the data into this header.
func (header *IPMISessionHeader_v1_5) Unmarshal(data []byte) error {
	buf := bytes.NewReader(data)
	_, err := header.Read(buf)
	return err
}

// Read the data from the reader into this header; return the number of bytes actually read and/or
// any error.
func (header *IPMISessionHeader_v1_5) Read(src_ io.Reader) (int, error) {
	src := countingReader(src_)
	if _, err := readObject(src, &header.AuthType, 1); err != nil {
		return src.numRead, err
	}
	if _, err := readObject(src, &header.SessionSequenceNumber, 4); err != nil {
		return src.numRead, err
	}
	if _, err := readObject(src, &header.SessionID, 4); err != nil {
		return src.numRead, err
	}

	if header.AuthType != AuthTypeNone {
		if _, err := read(src, header.AuthCode[:]); err != nil {
			return src.numRead, err
		}
	}
	if _, err := readObject(src, &header.Length, 1); err != nil {
		return src.numRead, err
	}
	return src.numRead, nil
}

// Write the header to the Writer, and return the number of bytes written and any error.
func (header *IPMISessionHeader_v1_5) Write(dst io.Writer) (int, error) {
	// If it's not None, we can just directly write it
	if header.AuthType != AuthTypeNone {
		return writeAndCount(dst, header)
	}
	// Otherwise, we must skip the AuthCode
	temp := bytes.NewBuffer(make([]byte, 0, 10))
	if _, err := write(temp, []byte{byte(header.AuthType)}); err != nil {
		return 0, err
	}
	if n, err := writeAndCount(temp, &header.SessionSequenceNumber); err != nil || n != 4 {
		if err == nil {
			err = io.ErrShortWrite
		}
		return 0, err
	}
	if n, err := writeAndCount(temp, &header.SessionID); err != nil || n != 4 {
		if err == nil {
			err = io.ErrShortWrite
		}
		return 0, err
	}
	if _, err := write(temp, []byte{byte(header.Length)}); err != nil {
		return 0, err
	}
	return write(dst, temp.Bytes())
}

// GetPacket returns an IPMISessionPacket with a copy of this header and the given body.
func (header *IPMISessionHeader_v1_5) GetPacket(body marshalable) (*IPMISessionPacket_v1_5, error) {
	ret := &IPMISessionPacket_v1_5{
		IPMISessionHeader_v1_5: *header,
	}
	if err := ret.SetBody(body); err != nil {
		return nil, err
	}
	return ret, nil
}

// IPMISessionPacket_v1_5 is a header + body for a v1.5 IPMI packet.
type IPMISessionPacket_v1_5 struct {
	IPMISessionHeader_v1_5
	Body []byte `json:"body,omitempty"`
}

// SetBody for this packet by marshaling it and setting the length.
func (p *IPMISessionPacket_v1_5) SetBody(body marshalable) error {
	data, err := body.Marshal()
	if err != nil {
		return err
	}
	if len(data) > 0xff {
		return fmt.Errorf("data too long: IPMISessionPacket_v1_5 body too large (%d > 255)", len(data))
	}
	p.Body = data
	p.Length = uint8(len(data))
	return nil
}

// Marshal the packet into a byte slice.
func (p *IPMISessionPacket_v1_5) Marshal() ([]byte, error) {
	temp := bytes.NewBuffer(make([]byte, 0, 16+len(p.Body)))
	if _, err := p.IPMISessionHeader_v1_5.Write(temp); err != nil {
		return nil, err
	}
	if _, err := write(temp, p.Body); err != nil {
		return nil, err
	}
	ret := temp.Bytes()
	return ret, nil
}

// Write the packet to the Writer, returning the number of bytes written and any error.
func (p *IPMISessionPacket_v1_5) Write(dst io.Writer) (int, error) {
	ret, err := p.Marshal()
	if err != nil {
		return 0, err
	}
	return write(dst, ret)
}

// Unmarshal this packet's contents from the given data or return an error.
func (p *IPMISessionPacket_v1_5) Unmarshal(data []byte) error {
	buf := bytes.NewReader(data)
	_, err := p.Read(buf)
	return err
}

// Read this packet's contents from the reader or return an error.
func (p *IPMISessionPacket_v1_5) Read(src io.Reader) (int, error) {
	n, err := p.IPMISessionHeader_v1_5.Read(src)
	if err != nil {
		return n, err
	}
	p.Body = make([]byte, p.Length)
	n2, err := read(src, p.Body)
	if n2 < int(p.Length) {
		p.Body = p.Body[0:n2]
	}
	return n + n2, err
}

// IPMIAddressByte represents either a slave address (lsb=0) or software ID (lsb=1).
// Upper 7 bits hold the actual address / ID.
// Always 0x20 for BMC.
type IPMIAddressByte uint8

// IsSlaveAddress returns true if this address represents a slave device (i.e. its lowest bit is not
// set).
func (v IPMIAddressByte) IsSlaveAddress() bool {
	return v&1 == 0
}

// IsSoftwareID returns true if this address represents a software identifier (i.e. its lowest bit
// is set).
func (v IPMIAddressByte) IsSoftwareID() bool {
	return v&1 == 1
}

// Addr returns this address byte's actual address (its upper 7 bits).
func (v IPMIAddressByte) Addr() uint8 {
	return uint8(v) >> 1
}

// Set the software ID/slave bit and address of this address byte.
func (v *IPMIAddressByte) Set(isSoftwareID bool, addr uint8) *IPMIAddressByte {
	lsb := uint8(0)
	if isSoftwareID {
		lsb = 1
	}
	*v = IPMIAddressByte(lsb | (addr << 1))
	return v
}

// IPMIAddressByteJSON represents the search-friendly version of the IPMIAddressByte.
type IPMIAddressByteJSON struct {
	Raw        uint8 `json:"raw"`
	Slave      bool  `json:"slave,omitempty"`
	SoftwareID bool  `json:"software_id,omitempty"`
	Address    uint8 `json:"address"`
}

// MarshalJSON returns the search-friendly representation of the IPMIAddressByte.
func (v IPMIAddressByte) MarshalJSON() ([]byte, error) {
	return json.Marshal(&IPMIAddressByteJSON{
		Raw:        uint8(v),
		Slave:      v.IsSlaveAddress(),
		SoftwareID: v.IsSoftwareID(),
		Address:    v.Addr(),
	})
}
