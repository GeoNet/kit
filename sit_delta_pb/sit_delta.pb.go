// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1-devel
// 	protoc        v3.17.3
// source: sit_delta.proto

// protoc --proto_path=protobuf/sit_delta --go_out=./ protobuf/sit_delta/sit_delta.proto

package sit_delta_pb

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// A geographical point on NZGD2000
type Point struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Latitude - geographical latitude of the point.
	Latitude float64 `protobuf:"fixed64,1,opt,name=latitude,proto3" json:"latitude,omitempty"`
	// Longitude - geographical longitude of the point.
	Longitude float64 `protobuf:"fixed64,2,opt,name=longitude,proto3" json:"longitude,omitempty"`
	// Elevation - geographical height of the point.
	Elevation float64 `protobuf:"fixed64,3,opt,name=elevation,proto3" json:"elevation,omitempty"`
	// Datum
	Datum string `protobuf:"bytes,4,opt,name=datum,proto3" json:"datum,omitempty"`
}

func (x *Point) Reset() {
	*x = Point{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sit_delta_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Point) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Point) ProtoMessage() {}

func (x *Point) ProtoReflect() protoreflect.Message {
	mi := &file_sit_delta_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Point.ProtoReflect.Descriptor instead.
func (*Point) Descriptor() ([]byte, []int) {
	return file_sit_delta_proto_rawDescGZIP(), []int{0}
}

func (x *Point) GetLatitude() float64 {
	if x != nil {
		return x.Latitude
	}
	return 0
}

func (x *Point) GetLongitude() float64 {
	if x != nil {
		return x.Longitude
	}
	return 0
}

func (x *Point) GetElevation() float64 {
	if x != nil {
		return x.Elevation
	}
	return 0
}

func (x *Point) GetDatum() string {
	if x != nil {
		return x.Datum
	}
	return ""
}

// A time span that has a start and an end.
type Span struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Start - time in Unix seconds.
	Start int64 `protobuf:"varint,1,opt,name=start,proto3" json:"start,omitempty"`
	// End - time in Unix seconds.  A future date of 9999-01-01T00:00:00Z is used to indicate still open.
	End int64 `protobuf:"varint,2,opt,name=end,proto3" json:"end,omitempty"`
}

func (x *Span) Reset() {
	*x = Span{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sit_delta_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Span) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Span) ProtoMessage() {}

func (x *Span) ProtoReflect() protoreflect.Message {
	mi := &file_sit_delta_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Span.ProtoReflect.Descriptor instead.
func (*Span) Descriptor() ([]byte, []int) {
	return file_sit_delta_proto_rawDescGZIP(), []int{1}
}

func (x *Span) GetStart() int64 {
	if x != nil {
		return x.Start
	}
	return 0
}

func (x *Span) GetEnd() int64 {
	if x != nil {
		return x.End
	}
	return 0
}

// A site record (represents a seismic site OR gps mark (OR tsunami station?))
type Site struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//The site code or station ID of the site
	Code string `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	//The location of the site
	Point *Point `protobuf:"bytes,2,opt,name=point,proto3" json:"point,omitempty"`
	//The ground relationship
	GroundRelationship float64 `protobuf:"fixed64,3,opt,name=ground_relationship,json=groundRelationship,proto3" json:"ground_relationship,omitempty"`
	//The network code
	Network string `protobuf:"bytes,4,opt,name=network,proto3" json:"network,omitempty"`
	//The date the site was established
	Span *Span `protobuf:"bytes,5,opt,name=span,proto3" json:"span,omitempty"`
	//Information for a 'Mark' site (will only exist if site is a mark)
	Mark *Mark `protobuf:"bytes,6,opt,name=mark,proto3" json:"mark,omitempty"`
	//List of 'Location' (for seismic + tsunami)
	Locations         []*Location          `protobuf:"bytes,7,rep,name=locations,proto3" json:"locations,omitempty"`
	EquipmentInstalls []*Equipment_Install `protobuf:"bytes,8,rep,name=equipment_installs,json=equipmentInstalls,proto3" json:"equipment_installs,omitempty"`
}

func (x *Site) Reset() {
	*x = Site{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sit_delta_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Site) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Site) ProtoMessage() {}

func (x *Site) ProtoReflect() protoreflect.Message {
	mi := &file_sit_delta_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Site.ProtoReflect.Descriptor instead.
func (*Site) Descriptor() ([]byte, []int) {
	return file_sit_delta_proto_rawDescGZIP(), []int{2}
}

func (x *Site) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *Site) GetPoint() *Point {
	if x != nil {
		return x.Point
	}
	return nil
}

func (x *Site) GetGroundRelationship() float64 {
	if x != nil {
		return x.GroundRelationship
	}
	return 0
}

func (x *Site) GetNetwork() string {
	if x != nil {
		return x.Network
	}
	return ""
}

func (x *Site) GetSpan() *Span {
	if x != nil {
		return x.Span
	}
	return nil
}

func (x *Site) GetMark() *Mark {
	if x != nil {
		return x.Mark
	}
	return nil
}

func (x *Site) GetLocations() []*Location {
	if x != nil {
		return x.Locations
	}
	return nil
}

func (x *Site) GetEquipmentInstalls() []*Equipment_Install {
	if x != nil {
		return x.EquipmentInstalls
	}
	return nil
}

// A location record (seismic site, tsunami guage, etc.)
type Location struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//The actual 'location' value (i.e. 40)
	Location string `protobuf:"bytes,1,opt,name=location,proto3" json:"location,omitempty"`
	//The geographical location
	Point  *Point `protobuf:"bytes,2,opt,name=point,proto3" json:"point,omitempty"`
	Status string `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
	//depth? maybe?
	GroundRelationship float64 `protobuf:"fixed64,4,opt,name=ground_relationship,json=groundRelationship,proto3" json:"ground_relationship,omitempty"`
	Notes              string  `protobuf:"bytes,5,opt,name=notes,proto3" json:"notes,omitempty"`
	Span               *Span   `protobuf:"bytes,6,opt,name=span,proto3" json:"span,omitempty"`
}

func (x *Location) Reset() {
	*x = Location{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sit_delta_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Location) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Location) ProtoMessage() {}

func (x *Location) ProtoReflect() protoreflect.Message {
	mi := &file_sit_delta_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Location.ProtoReflect.Descriptor instead.
func (*Location) Descriptor() ([]byte, []int) {
	return file_sit_delta_proto_rawDescGZIP(), []int{3}
}

func (x *Location) GetLocation() string {
	if x != nil {
		return x.Location
	}
	return ""
}

func (x *Location) GetPoint() *Point {
	if x != nil {
		return x.Point
	}
	return nil
}

func (x *Location) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Location) GetGroundRelationship() float64 {
	if x != nil {
		return x.GroundRelationship
	}
	return 0
}

func (x *Location) GetNotes() string {
	if x != nil {
		return x.Notes
	}
	return ""
}

func (x *Location) GetSpan() *Span {
	if x != nil {
		return x.Span
	}
	return nil
}

// An equipment record
type Equipment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AssetNumber  string  `protobuf:"bytes,1,opt,name=asset_number,json=assetNumber,proto3" json:"asset_number,omitempty"`
	SerialNumber string  `protobuf:"bytes,2,opt,name=serial_number,json=serialNumber,proto3" json:"serial_number,omitempty"`
	Manufacturer string  `protobuf:"bytes,3,opt,name=manufacturer,proto3" json:"manufacturer,omitempty"`
	Model        string  `protobuf:"bytes,4,opt,name=model,proto3" json:"model,omitempty"`
	Type         string  `protobuf:"bytes,5,opt,name=type,proto3" json:"type,omitempty"`
	Owner        string  `protobuf:"bytes,6,opt,name=owner,proto3" json:"owner,omitempty"`
	Height       float64 `protobuf:"fixed64,9,opt,name=height,proto3" json:"height,omitempty"`
	Location     string  `protobuf:"bytes,10,opt,name=location,proto3" json:"location,omitempty"`
	Orientation  float64 `protobuf:"fixed64,11,opt,name=orientation,proto3" json:"orientation,omitempty"`
}

func (x *Equipment) Reset() {
	*x = Equipment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sit_delta_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Equipment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Equipment) ProtoMessage() {}

func (x *Equipment) ProtoReflect() protoreflect.Message {
	mi := &file_sit_delta_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Equipment.ProtoReflect.Descriptor instead.
func (*Equipment) Descriptor() ([]byte, []int) {
	return file_sit_delta_proto_rawDescGZIP(), []int{4}
}

func (x *Equipment) GetAssetNumber() string {
	if x != nil {
		return x.AssetNumber
	}
	return ""
}

func (x *Equipment) GetSerialNumber() string {
	if x != nil {
		return x.SerialNumber
	}
	return ""
}

func (x *Equipment) GetManufacturer() string {
	if x != nil {
		return x.Manufacturer
	}
	return ""
}

func (x *Equipment) GetModel() string {
	if x != nil {
		return x.Model
	}
	return ""
}

func (x *Equipment) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Equipment) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

func (x *Equipment) GetHeight() float64 {
	if x != nil {
		return x.Height
	}
	return 0
}

func (x *Equipment) GetLocation() string {
	if x != nil {
		return x.Location
	}
	return ""
}

func (x *Equipment) GetOrientation() float64 {
	if x != nil {
		return x.Orientation
	}
	return 0
}

// An equipment install record
type Equipment_Install struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Equipment *Equipment `protobuf:"bytes,1,opt,name=equipment,proto3" json:"equipment,omitempty"`
	Installed *Span      `protobuf:"bytes,2,opt,name=installed,proto3" json:"installed,omitempty"`
}

func (x *Equipment_Install) Reset() {
	*x = Equipment_Install{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sit_delta_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Equipment_Install) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Equipment_Install) ProtoMessage() {}

func (x *Equipment_Install) ProtoReflect() protoreflect.Message {
	mi := &file_sit_delta_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Equipment_Install.ProtoReflect.Descriptor instead.
func (*Equipment_Install) Descriptor() ([]byte, []int) {
	return file_sit_delta_proto_rawDescGZIP(), []int{5}
}

func (x *Equipment_Install) GetEquipment() *Equipment {
	if x != nil {
		return x.Equipment
	}
	return nil
}

func (x *Equipment_Install) GetInstalled() *Span {
	if x != nil {
		return x.Installed
	}
	return nil
}

// A GNSS Mark.
type Mark struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InstalledMonument []*InstalledMonument `protobuf:"bytes,1,rep,name=installed_monument,json=installedMonument,proto3" json:"installed_monument,omitempty"`
	Point             *Point               `protobuf:"bytes,2,opt,name=point,proto3" json:"point,omitempty"`
}

func (x *Mark) Reset() {
	*x = Mark{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sit_delta_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Mark) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Mark) ProtoMessage() {}

func (x *Mark) ProtoReflect() protoreflect.Message {
	mi := &file_sit_delta_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Mark.ProtoReflect.Descriptor instead.
func (*Mark) Descriptor() ([]byte, []int) {
	return file_sit_delta_proto_rawDescGZIP(), []int{6}
}

func (x *Mark) GetInstalledMonument() []*InstalledMonument {
	if x != nil {
		return x.InstalledMonument
	}
	return nil
}

func (x *Mark) GetPoint() *Point {
	if x != nil {
		return x.Point
	}
	return nil
}

type InstalledMonument struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Span     *Span     `protobuf:"bytes,1,opt,name=span,proto3" json:"span,omitempty"`
	Monument *Monument `protobuf:"bytes,2,opt,name=monument,proto3" json:"monument,omitempty"`
}

func (x *InstalledMonument) Reset() {
	*x = InstalledMonument{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sit_delta_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InstalledMonument) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstalledMonument) ProtoMessage() {}

func (x *InstalledMonument) ProtoReflect() protoreflect.Message {
	mi := &file_sit_delta_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstalledMonument.ProtoReflect.Descriptor instead.
func (*InstalledMonument) Descriptor() ([]byte, []int) {
	return file_sit_delta_proto_rawDescGZIP(), []int{7}
}

func (x *InstalledMonument) GetSpan() *Span {
	if x != nil {
		return x.Span
	}
	return nil
}

func (x *InstalledMonument) GetMonument() *Monument {
	if x != nil {
		return x.Monument
	}
	return nil
}

// A monument for a Mark
type Monument struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DomesNumber string  `protobuf:"bytes,1,opt,name=domes_number,json=domesNumber,proto3" json:"domes_number,omitempty"`
	Height      float64 `protobuf:"fixed64,3,opt,name=height,proto3" json:"height,omitempty"`
}

func (x *Monument) Reset() {
	*x = Monument{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sit_delta_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Monument) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Monument) ProtoMessage() {}

func (x *Monument) ProtoReflect() protoreflect.Message {
	mi := &file_sit_delta_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Monument.ProtoReflect.Descriptor instead.
func (*Monument) Descriptor() ([]byte, []int) {
	return file_sit_delta_proto_rawDescGZIP(), []int{8}
}

func (x *Monument) GetDomesNumber() string {
	if x != nil {
		return x.DomesNumber
	}
	return ""
}

func (x *Monument) GetHeight() float64 {
	if x != nil {
		return x.Height
	}
	return 0
}

var File_sit_delta_proto protoreflect.FileDescriptor

var file_sit_delta_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x73, 0x69, 0x74, 0x5f, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x09, 0x73, 0x69, 0x74, 0x5f, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x22, 0x75, 0x0a, 0x05,
	0x50, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x6c, 0x61, 0x74, 0x69, 0x74, 0x75, 0x64,
	0x65, 0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x09, 0x6c, 0x6f, 0x6e, 0x67, 0x69, 0x74, 0x75, 0x64, 0x65, 0x12,
	0x1c, 0x0a, 0x09, 0x65, 0x6c, 0x65, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x09, 0x65, 0x6c, 0x65, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a,
	0x05, 0x64, 0x61, 0x74, 0x75, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x64, 0x61,
	0x74, 0x75, 0x6d, 0x22, 0x2e, 0x0a, 0x04, 0x53, 0x70, 0x61, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03,
	0x65, 0x6e, 0x64, 0x22, 0xd7, 0x02, 0x0a, 0x04, 0x73, 0x69, 0x74, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65,
	0x12, 0x26, 0x0a, 0x05, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x10, 0x2e, 0x73, 0x69, 0x74, 0x5f, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x2e, 0x50, 0x6f, 0x69, 0x6e,
	0x74, 0x52, 0x05, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x2f, 0x0a, 0x13, 0x67, 0x72, 0x6f, 0x75,
	0x6e, 0x64, 0x5f, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x68, 0x69, 0x70, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x12, 0x67, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x6c,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x68, 0x69, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x6e, 0x65, 0x74,
	0x77, 0x6f, 0x72, 0x6b, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6e, 0x65, 0x74, 0x77,
	0x6f, 0x72, 0x6b, 0x12, 0x23, 0x0a, 0x04, 0x73, 0x70, 0x61, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0f, 0x2e, 0x73, 0x69, 0x74, 0x5f, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x2e, 0x53, 0x70,
	0x61, 0x6e, 0x52, 0x04, 0x73, 0x70, 0x61, 0x6e, 0x12, 0x23, 0x0a, 0x04, 0x6d, 0x61, 0x72, 0x6b,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x73, 0x69, 0x74, 0x5f, 0x64, 0x65, 0x6c,
	0x74, 0x61, 0x2e, 0x4d, 0x61, 0x72, 0x6b, 0x52, 0x04, 0x6d, 0x61, 0x72, 0x6b, 0x12, 0x31, 0x0a,
	0x09, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x13, 0x2e, 0x73, 0x69, 0x74, 0x5f, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x2e, 0x4c, 0x6f, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x09, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x12, 0x4b, 0x0a, 0x12, 0x65, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x6e,
	0x73, 0x74, 0x61, 0x6c, 0x6c, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x73,
	0x69, 0x74, 0x5f, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x2e, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65,
	0x6e, 0x74, 0x5f, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x52, 0x11, 0x65, 0x71, 0x75, 0x69,
	0x70, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x73, 0x22, 0xd2, 0x01,
	0x0a, 0x08, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x6f,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x6f,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x26, 0x0a, 0x05, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x73, 0x69, 0x74, 0x5f, 0x64, 0x65, 0x6c, 0x74,
	0x61, 0x2e, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x05, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x2f, 0x0a, 0x13, 0x67, 0x72, 0x6f, 0x75, 0x6e, 0x64,
	0x5f, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x68, 0x69, 0x70, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x12, 0x67, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x52, 0x65, 0x6c, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x68, 0x69, 0x70, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x6f, 0x74, 0x65, 0x73,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6e, 0x6f, 0x74, 0x65, 0x73, 0x12, 0x23, 0x0a,
	0x04, 0x73, 0x70, 0x61, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x73, 0x69,
	0x74, 0x5f, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x2e, 0x53, 0x70, 0x61, 0x6e, 0x52, 0x04, 0x73, 0x70,
	0x61, 0x6e, 0x22, 0x8d, 0x02, 0x0a, 0x09, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74,
	0x12, 0x21, 0x0a, 0x0c, 0x61, 0x73, 0x73, 0x65, 0x74, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x73, 0x73, 0x65, 0x74, 0x4e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x5f, 0x6e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x65, 0x72, 0x69,
	0x61, 0x6c, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x22, 0x0a, 0x0c, 0x6d, 0x61, 0x6e, 0x75,
	0x66, 0x61, 0x63, 0x74, 0x75, 0x72, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x6d, 0x61, 0x6e, 0x75, 0x66, 0x61, 0x63, 0x74, 0x75, 0x72, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6d, 0x6f, 0x64,
	0x65, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06,
	0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06, 0x68, 0x65,
	0x69, 0x67, 0x68, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x20, 0x0a, 0x0b, 0x6f, 0x72, 0x69, 0x65, 0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x0b, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0b, 0x6f, 0x72, 0x69, 0x65, 0x6e, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x22, 0x76, 0x0a, 0x11, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x5f,
	0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x12, 0x32, 0x0a, 0x09, 0x65, 0x71, 0x75, 0x69, 0x70,
	0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x73, 0x69, 0x74,
	0x5f, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x2e, 0x45, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74,
	0x52, 0x09, 0x65, 0x71, 0x75, 0x69, 0x70, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x2d, 0x0a, 0x09, 0x69,
	0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f,
	0x2e, 0x73, 0x69, 0x74, 0x5f, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x2e, 0x53, 0x70, 0x61, 0x6e, 0x52,
	0x09, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x65, 0x64, 0x22, 0x7b, 0x0a, 0x04, 0x4d, 0x61,
	0x72, 0x6b, 0x12, 0x4b, 0x0a, 0x12, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x65, 0x64, 0x5f,
	0x6d, 0x6f, 0x6e, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c,
	0x2e, 0x73, 0x69, 0x74, 0x5f, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61,
	0x6c, 0x6c, 0x65, 0x64, 0x4d, 0x6f, 0x6e, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x11, 0x69, 0x6e,
	0x73, 0x74, 0x61, 0x6c, 0x6c, 0x65, 0x64, 0x4d, 0x6f, 0x6e, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x12,
	0x26, 0x0a, 0x05, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10,
	0x2e, 0x73, 0x69, 0x74, 0x5f, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x2e, 0x50, 0x6f, 0x69, 0x6e, 0x74,
	0x52, 0x05, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x22, 0x69, 0x0a, 0x11, 0x49, 0x6e, 0x73, 0x74, 0x61,
	0x6c, 0x6c, 0x65, 0x64, 0x4d, 0x6f, 0x6e, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x23, 0x0a, 0x04,
	0x73, 0x70, 0x61, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x73, 0x69, 0x74,
	0x5f, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x2e, 0x53, 0x70, 0x61, 0x6e, 0x52, 0x04, 0x73, 0x70, 0x61,
	0x6e, 0x12, 0x2f, 0x0a, 0x08, 0x6d, 0x6f, 0x6e, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x73, 0x69, 0x74, 0x5f, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x2e,
	0x4d, 0x6f, 0x6e, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x08, 0x6d, 0x6f, 0x6e, 0x75, 0x6d, 0x65,
	0x6e, 0x74, 0x22, 0x45, 0x0a, 0x08, 0x4d, 0x6f, 0x6e, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x21,
	0x0a, 0x0c, 0x64, 0x6f, 0x6d, 0x65, 0x73, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x6f, 0x6d, 0x65, 0x73, 0x4e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x06, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x42, 0x10, 0x5a, 0x0e, 0x2e, 0x2f, 0x73,
	0x69, 0x74, 0x5f, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x5f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_sit_delta_proto_rawDescOnce sync.Once
	file_sit_delta_proto_rawDescData = file_sit_delta_proto_rawDesc
)

func file_sit_delta_proto_rawDescGZIP() []byte {
	file_sit_delta_proto_rawDescOnce.Do(func() {
		file_sit_delta_proto_rawDescData = protoimpl.X.CompressGZIP(file_sit_delta_proto_rawDescData)
	})
	return file_sit_delta_proto_rawDescData
}

var file_sit_delta_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_sit_delta_proto_goTypes = []interface{}{
	(*Point)(nil),             // 0: sit_delta.Point
	(*Span)(nil),              // 1: sit_delta.Span
	(*Site)(nil),              // 2: sit_delta.site
	(*Location)(nil),          // 3: sit_delta.Location
	(*Equipment)(nil),         // 4: sit_delta.Equipment
	(*Equipment_Install)(nil), // 5: sit_delta.Equipment_Install
	(*Mark)(nil),              // 6: sit_delta.Mark
	(*InstalledMonument)(nil), // 7: sit_delta.InstalledMonument
	(*Monument)(nil),          // 8: sit_delta.Monument
}
var file_sit_delta_proto_depIdxs = []int32{
	0,  // 0: sit_delta.site.point:type_name -> sit_delta.Point
	1,  // 1: sit_delta.site.span:type_name -> sit_delta.Span
	6,  // 2: sit_delta.site.mark:type_name -> sit_delta.Mark
	3,  // 3: sit_delta.site.locations:type_name -> sit_delta.Location
	5,  // 4: sit_delta.site.equipment_installs:type_name -> sit_delta.Equipment_Install
	0,  // 5: sit_delta.Location.point:type_name -> sit_delta.Point
	1,  // 6: sit_delta.Location.span:type_name -> sit_delta.Span
	4,  // 7: sit_delta.Equipment_Install.equipment:type_name -> sit_delta.Equipment
	1,  // 8: sit_delta.Equipment_Install.installed:type_name -> sit_delta.Span
	7,  // 9: sit_delta.Mark.installed_monument:type_name -> sit_delta.InstalledMonument
	0,  // 10: sit_delta.Mark.point:type_name -> sit_delta.Point
	1,  // 11: sit_delta.InstalledMonument.span:type_name -> sit_delta.Span
	8,  // 12: sit_delta.InstalledMonument.monument:type_name -> sit_delta.Monument
	13, // [13:13] is the sub-list for method output_type
	13, // [13:13] is the sub-list for method input_type
	13, // [13:13] is the sub-list for extension type_name
	13, // [13:13] is the sub-list for extension extendee
	0,  // [0:13] is the sub-list for field type_name
}

func init() { file_sit_delta_proto_init() }
func file_sit_delta_proto_init() {
	if File_sit_delta_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_sit_delta_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Point); i {
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
		file_sit_delta_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Span); i {
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
		file_sit_delta_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Site); i {
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
		file_sit_delta_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Location); i {
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
		file_sit_delta_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Equipment); i {
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
		file_sit_delta_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Equipment_Install); i {
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
		file_sit_delta_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Mark); i {
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
		file_sit_delta_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InstalledMonument); i {
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
		file_sit_delta_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Monument); i {
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
			RawDescriptor: file_sit_delta_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_sit_delta_proto_goTypes,
		DependencyIndexes: file_sit_delta_proto_depIdxs,
		MessageInfos:      file_sit_delta_proto_msgTypes,
	}.Build()
	File_sit_delta_proto = out.File
	file_sit_delta_proto_rawDesc = nil
	file_sit_delta_proto_goTypes = nil
	file_sit_delta_proto_depIdxs = nil
}
