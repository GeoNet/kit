

# seis
`import "github.com/GeoNet/kit/seis"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
The seis module has been writen as a lightweight replacement for the C libraries libmseed and libslink.
It is aimed at clients that need to decode miniseed data either directly or collected from a seedlink
server.

The seedlink code is not a direct replacement for libslink. It can run in two modes, either as a
raw connection to the client connection (SLConn) which allows mechanisms to monitor or have a finer
control of the SeedLink connection, or in the collection mode (SLink) where a connection is established
and received miniseed blocks can be processed with a call back function. A context can be passed into
the collection loop to allow interuption or as a shutdown mechanism. It is not passed to the underlying
seedlink connection messaging which is managed via a deadline mechanism, e.g. the `SetTimeout` option.

An example Seedlink application can be as simple as:


	slink := seis.NewSLink("localhost:18000")
	
	if err := slink.Collect(func(seq string, data []byte) (bool, error) {
	        if ms, err := seis.NewMSRecord(data); err == nil {
	            log.Println(ms.SrcName(false), time.Since(ms.EndTime()))
	       }
	       return false, nil
	}); err != nil {
	        log.Fatal(err)
	}

The conversion to miniseed can be coupled together, is in this example, but it is not required. Either
the raw packet can be managed as a whole unit or it can be unpacked using another mechanism.




## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [func MarshalBTime(r BTime) (b [BTimeSize]byte)](#MarshalBTime)
* [func MarshalBlockette1000(r Blockette1000) (b [Blockette1000Size]byte)](#MarshalBlockette1000)
* [func MarshalBlockette1001(r Blockette1001) (b [Blockette1001Size]byte)](#MarshalBlockette1001)
* [func MarshalBlocketteHeader(r BlocketteHeader) (b [BlocketteHeaderSize]byte)](#MarshalBlocketteHeader)
* [func MarshalDLPreheader(r DLPreheader) (b [DLPreheaderSize]byte)](#MarshalDLPreheader)
* [func MarshalMSDataHeader(r MSDataHeader) (b [MSDataHeaderSize]byte)](#MarshalMSDataHeader)
* [type BTime](#BTime)
  * [func NewBTime(t time.Time) BTime](#NewBTime)
  * [func UnmarshalBTime(b [BTimeSize]byte) (r BTime)](#UnmarshalBTime)
  * [func (b BTime) Time() time.Time](#BTime.Time)
* [type Blockette1000](#Blockette1000)
  * [func UnmarshalBlockette1000(b [Blockette1000Size]byte) (r Blockette1000)](#UnmarshalBlockette1000)
* [type Blockette1001](#Blockette1001)
  * [func UnmarshalBlockette1001(b [Blockette1001Size]byte) (r Blockette1001)](#UnmarshalBlockette1001)
* [type BlocketteHeader](#BlocketteHeader)
  * [func UnmarshalBlocketteHeader(b [BlocketteHeaderSize]byte) (r BlocketteHeader)](#UnmarshalBlocketteHeader)
* [type DLConn](#DLConn)
  * [func NewDLConn(service string, timeout time.Duration) (*DLConn, error)](#NewDLConn)
  * [func (d *DLConn) Id() string](#DLConn.Id)
  * [func (d *DLConn) SetId(program, username string) error](#DLConn.SetId)
  * [func (d *DLConn) Size() int](#DLConn.Size)
  * [func (d *DLConn) Writable() bool](#DLConn.Writable)
  * [func (d *DLConn) WriteMS(data []byte) error](#DLConn.WriteMS)
* [type DLPacket](#DLPacket)
* [type DLPreheader](#DLPreheader)
  * [func UnmarshalDLPreheader(b [DLPreheaderSize]byte) (r DLPreheader)](#UnmarshalDLPreheader)
* [type DLink](#DLink)
  * [func NewDLink(server string, opts ...DLinkOpt) *DLink](#NewDLink)
  * [func (d *DLink) Connect() (*DLConn, error)](#DLink.Connect)
  * [func (d *DLink) SetProgram(s string)](#DLink.SetProgram)
  * [func (d *DLink) SetTimeout(t time.Duration)](#DLink.SetTimeout)
  * [func (d *DLink) SetUsername(s string)](#DLink.SetUsername)
* [type DLinkOpt](#DLinkOpt)
  * [func SetDLProgram(s string) DLinkOpt](#SetDLProgram)
  * [func SetDLTimeout(t time.Duration) DLinkOpt](#SetDLTimeout)
  * [func SetDLUsername(s string) DLinkOpt](#SetDLUsername)
* [type MSDataHeader](#MSDataHeader)
  * [func UnmarshalMSDataHeader(b [MSDataHeaderSize]byte) (r MSDataHeader)](#UnmarshalMSDataHeader)
* [type MSEncoding](#MSEncoding)
* [type MSRecord](#MSRecord)
  * [func NewMSRecord(buf []byte) (*MSRecord, error)](#NewMSRecord)
  * [func (m *MSRecord) BlockSize() int](#MSRecord.BlockSize)
  * [func (m *MSRecord) ByteOrder() WordOrder](#MSRecord.ByteOrder)
  * [func (m *MSRecord) Channel() string](#MSRecord.Channel)
  * [func (m *MSRecord) DataQuality() byte](#MSRecord.DataQuality)
  * [func (m *MSRecord) Encoding() MSEncoding](#MSRecord.Encoding)
  * [func (m *MSRecord) EndTime() time.Time](#MSRecord.EndTime)
  * [func (m *MSRecord) Location() string](#MSRecord.Location)
  * [func (m *MSRecord) Network() string](#MSRecord.Network)
  * [func (m *MSRecord) NumSamples() int](#MSRecord.NumSamples)
  * [func (m *MSRecord) SampleCount() int](#MSRecord.SampleCount)
  * [func (m *MSRecord) SamplePeriod() time.Duration](#MSRecord.SamplePeriod)
  * [func (m *MSRecord) SampleRate() float64](#MSRecord.SampleRate)
  * [func (m *MSRecord) SampleType() MSSampleType](#MSRecord.SampleType)
  * [func (m *MSRecord) SequenceNumber() int](#MSRecord.SequenceNumber)
  * [func (m *MSRecord) SrcName(quality bool) string](#MSRecord.SrcName)
  * [func (m *MSRecord) StartTime() time.Time](#MSRecord.StartTime)
  * [func (m *MSRecord) Station() string](#MSRecord.Station)
  * [func (m *MSRecord) String() string](#MSRecord.String)
  * [func (m *MSRecord) ToBytes() []byte](#MSRecord.ToBytes)
  * [func (m *MSRecord) ToFloat64s() []float64](#MSRecord.ToFloat64s)
  * [func (m *MSRecord) ToInt32s() []int32](#MSRecord.ToInt32s)
  * [func (m *MSRecord) ToInts() []int](#MSRecord.ToInts)
  * [func (m *MSRecord) ToStrings() []string](#MSRecord.ToStrings)
  * [func (m *MSRecord) Unpack(buf []byte, dataflag bool) error](#MSRecord.Unpack)
* [type MSRecordFunc](#MSRecordFunc)
* [type MSRecordStream](#MSRecordStream)
  * [func NewMSRecordStream(msr *MSRecord) *MSRecordStream](#NewMSRecordStream)
  * [func (m *MSRecordStream) PackInt32(start time.Time, quality int, locked bool, data []int32, fn MSRecordFunc) error](#MSRecordStream.PackInt32)
* [type MSSampleType](#MSSampleType)
* [type SLCollectFunc](#SLCollectFunc)
* [type SLConn](#SLConn)
  * [func NewSLConn(service string, timeout time.Duration) (*SLConn, error)](#NewSLConn)
  * [func (c *SLConn) Collect() (*SLPacket, error)](#SLConn.Collect)
  * [func (c *SLConn) CommandCat() ([]byte, error)](#SLConn.CommandCat)
  * [func (c *SLConn) CommandClose() ([]byte, error)](#SLConn.CommandClose)
  * [func (c *SLConn) CommandData(sequence string, starttime time.Time) error](#SLConn.CommandData)
  * [func (c *SLConn) CommandEnd() error](#SLConn.CommandEnd)
  * [func (c *SLConn) CommandHello() ([]byte, error)](#SLConn.CommandHello)
  * [func (c *SLConn) CommandId() ([]byte, error)](#SLConn.CommandId)
  * [func (c *SLConn) CommandSelect(selection string) error](#SLConn.CommandSelect)
  * [func (c *SLConn) CommandStation(station, network string) error](#SLConn.CommandStation)
  * [func (c *SLConn) CommandTime(starttime, endtime time.Time) error](#SLConn.CommandTime)
  * [func (c *SLConn) GetInfo(level string) ([]byte, error)](#SLConn.GetInfo)
  * [func (c *SLConn) GetSLInfo(level string) (*SLInfo, error)](#SLConn.GetSLInfo)
* [type SLInfo](#SLInfo)
  * [func (s *SLInfo) Unmarshal(data []byte) error](#SLInfo.Unmarshal)
* [type SLPacket](#SLPacket)
  * [func NewSLPacket(data []byte) (*SLPacket, error)](#NewSLPacket)
* [type SLState](#SLState)
  * [func (s *SLState) Add(station SLStation)](#SLState.Add)
  * [func (s *SLState) Find(stn SLStation) *SLStation](#SLState.Find)
  * [func (s *SLState) Marshal() ([]byte, error)](#SLState.Marshal)
  * [func (s *SLState) ReadFile(path string) error](#SLState.ReadFile)
  * [func (s *SLState) Stations() []SLStation](#SLState.Stations)
  * [func (s *SLState) Unmarshal(data []byte) error](#SLState.Unmarshal)
  * [func (s *SLState) WriteFile(path string) error](#SLState.WriteFile)
* [type SLStation](#SLStation)
  * [func (s SLStation) Key() SLStation](#SLStation.Key)
* [type SLink](#SLink)
  * [func NewSLink(server string, opts ...SLinkOpt) *SLink](#NewSLink)
  * [func (s *SLink) Collect(fn SLCollectFunc) error](#SLink.Collect)
  * [func (s *SLink) CollectWithContext(ctx context.Context, fn SLCollectFunc) error](#SLink.CollectWithContext)
  * [func (s *SLink) SetEndTime(t time.Time)](#SLink.SetEndTime)
  * [func (s *SLink) SetKeepAlive(d time.Duration)](#SLink.SetKeepAlive)
  * [func (s *SLink) SetNetTo(d time.Duration)](#SLink.SetNetTo)
  * [func (s *SLink) SetRefresh(d time.Duration)](#SLink.SetRefresh)
  * [func (s *SLink) SetSelectors(selectors string)](#SLink.SetSelectors)
  * [func (s *SLink) SetSequence(sequence int)](#SLink.SetSequence)
  * [func (s *SLink) SetStartTime(t time.Time)](#SLink.SetStartTime)
  * [func (s *SLink) SetState(stations ...SLStation)](#SLink.SetState)
  * [func (s *SLink) SetStateFile(f string)](#SLink.SetStateFile)
  * [func (s *SLink) SetStreams(streams string)](#SLink.SetStreams)
  * [func (s *SLink) SetTimeout(d time.Duration)](#SLink.SetTimeout)
* [type SLinkOpt](#SLinkOpt)
  * [func SetSLEndTime(t time.Time) SLinkOpt](#SetSLEndTime)
  * [func SetSLKeepAlive(d time.Duration) SLinkOpt](#SetSLKeepAlive)
  * [func SetSLNetTo(d time.Duration) SLinkOpt](#SetSLNetTo)
  * [func SetSLRefresh(d time.Duration) SLinkOpt](#SetSLRefresh)
  * [func SetSLSelectors(selectors string) SLinkOpt](#SetSLSelectors)
  * [func SetSLSequence(sequence int) SLinkOpt](#SetSLSequence)
  * [func SetSLStartTime(t time.Time) SLinkOpt](#SetSLStartTime)
  * [func SetSLState(stations ...SLStation) SLinkOpt](#SetSLState)
  * [func SetSLStateFile(f string) SLinkOpt](#SetSLStateFile)
  * [func SetSLStreams(streams string) SLinkOpt](#SetSLStreams)
  * [func SetSLTimeout(d time.Duration) SLinkOpt](#SetSLTimeout)
* [type WordOrder](#WordOrder)


#### <a name="pkg-files">Package files</a>
[dl_conn.go](/src/target/dl_conn.go) [dl_dlink.go](/src/target/dl_dlink.go) [dl_types.go](/src/target/dl_types.go) [doc.go](/src/target/doc.go) [ms_decode.go](/src/target/ms_decode.go) [ms_encode.go](/src/target/ms_encode.go) [ms_pack.go](/src/target/ms_pack.go) [ms_record.go](/src/target/ms_record.go) [ms_steim.go](/src/target/ms_steim.go) [ms_types.go](/src/target/ms_types.go) [sl_conn.go](/src/target/sl_conn.go) [sl_info.go](/src/target/sl_info.go) [sl_packet.go](/src/target/sl_packet.go) [sl_slink.go](/src/target/sl_slink.go) [sl_state.go](/src/target/sl_state.go) [sl_stream.go](/src/target/sl_stream.go) 


## <a name="pkg-constants">Constants</a>
``` go
const (
    MSDataHeaderSize    = 48
    BTimeSize           = 10
    BlocketteHeaderSize = 4
    Blockette1000Size   = 4
    Blockette1001Size   = 4
)
```
``` go
const (
    DLPreheaderSize = 3
)
```
``` go
const (
    SLPacketSize = 8 + 512
)
```



## <a name="MarshalBTime">func</a> [MarshalBTime](/src/target/ms_types.go?s=1095:1141#L62)
``` go
func MarshalBTime(r BTime) (b [BTimeSize]byte)
```


## <a name="MarshalBlockette1000">func</a> [MarshalBlockette1000](/src/target/ms_types.go?s=4632:4702#L193)
``` go
func MarshalBlockette1000(r Blockette1000) (b [Blockette1000Size]byte)
```


## <a name="MarshalBlockette1001">func</a> [MarshalBlockette1001](/src/target/ms_types.go?s=5199:5269#L218)
``` go
func MarshalBlockette1001(r Blockette1001) (b [Blockette1001Size]byte)
```


## <a name="MarshalBlocketteHeader">func</a> [MarshalBlocketteHeader](/src/target/ms_types.go?s=4108:4184#L170)
``` go
func MarshalBlocketteHeader(r BlocketteHeader) (b [BlocketteHeaderSize]byte)
```


## <a name="MarshalDLPreheader">func</a> [MarshalDLPreheader](/src/target/dl_types.go?s=557:621#L37)
``` go
func MarshalDLPreheader(r DLPreheader) (b [DLPreheaderSize]byte)
```


## <a name="MarshalMSDataHeader">func</a> [MarshalMSDataHeader](/src/target/ms_types.go?s=2978:3045#L130)
``` go
func MarshalMSDataHeader(r MSDataHeader) (b [MSDataHeaderSize]byte)
```



## <a name="BTime">type</a> [BTime](/src/target/ms_types.go?s=191:398#L16)
``` go
type BTime struct {
    Year   uint16
    Doy    uint16
    Hour   uint8
    Minute uint8
    Second uint8
    Unused byte   //Required for "alignment"
    S0001  uint16 //.0001 of a second 0-9999
}

```






### <a name="NewBTime">func</a> [NewBTime](/src/target/ms_types.go?s=596:628#L39)
``` go
func NewBTime(t time.Time) BTime
```

### <a name="UnmarshalBTime">func</a> [UnmarshalBTime](/src/target/ms_types.go?s=839:887#L50)
``` go
func UnmarshalBTime(b [BTimeSize]byte) (r BTime)
```




### <a name="BTime.Time">func</a> (BTime) [Time](/src/target/ms_types.go?s=400:431#L26)
``` go
func (b BTime) Time() time.Time
```



## <a name="Blockette1000">type</a> [Blockette1000](/src/target/ms_types.go?s=4305:4462#L177)
``` go
type Blockette1000 struct {
    Encoding     uint8
    WordOrder    uint8
    RecordLength uint8
    Reserved     uint8
}

```






### <a name="UnmarshalBlockette1000">func</a> [UnmarshalBlockette1000](/src/target/ms_types.go?s=4464:4536#L184)
``` go
func UnmarshalBlockette1000(b [Blockette1000Size]byte) (r Blockette1000)
```




## <a name="Blockette1001">type</a> [Blockette1001](/src/target/ms_types.go?s=4826:5021#L202)
``` go
type Blockette1001 struct {
    TimingQuality uint8
    MicroSec      int8 //Increased accuracy for starttime
    Reserved      uint8
    FrameCount    uint8
}

```






### <a name="UnmarshalBlockette1001">func</a> [UnmarshalBlockette1001](/src/target/ms_types.go?s=5023:5095#L209)
``` go
func UnmarshalBlockette1001(b [Blockette1001Size]byte) (r Blockette1001)
```




## <a name="BlocketteHeader">type</a> [BlocketteHeader](/src/target/ms_types.go?s=3790:3911#L158)
``` go
type BlocketteHeader struct {
    BlocketteType uint16
    NextBlockette uint16 //Byte of next blockette, 0 if last blockette
}

```






### <a name="UnmarshalBlocketteHeader">func</a> [UnmarshalBlocketteHeader](/src/target/ms_types.go?s=3913:3991#L163)
``` go
func UnmarshalBlocketteHeader(b [BlocketteHeaderSize]byte) (r BlocketteHeader)
```




## <a name="DLConn">type</a> [DLConn](/src/target/dl_conn.go?s=204:306#L19)
``` go
type DLConn struct {
    net.Conn
    // contains filtered or unexported fields
}

```
DLConn provides connection information to a datalink service.







### <a name="NewDLConn">func</a> [NewDLConn](/src/target/dl_conn.go?s=495:565#L31)
``` go
func NewDLConn(service string, timeout time.Duration) (*DLConn, error)
```
NewDLConn makes a connection to a datalink service, the function
SetId should be run after the connection is made and Close
should be called when the link is no longer required.





### <a name="DLConn.Id">func</a> (\*DLConn) [Id](/src/target/dl_conn.go?s=878:906#L49)
``` go
func (d *DLConn) Id() string
```
Id returns the client identification sent to the server.




### <a name="DLConn.SetId">func</a> (\*DLConn) [SetId](/src/target/dl_conn.go?s=1246:1300#L64)
``` go
func (d *DLConn) SetId(program, username string) error
```
SetId sens an ID message to the remote connection and decodes the connection capabilities.




### <a name="DLConn.Size">func</a> (\*DLConn) [Size](/src/target/dl_conn.go?s=1104:1131#L59)
``` go
func (d *DLConn) Size() int
```
Size returns the current expected packet size.




### <a name="DLConn.Writable">func</a> (\*DLConn) [Writable](/src/target/dl_conn.go?s=997:1029#L54)
``` go
func (d *DLConn) Writable() bool
```
Writable returns whether packets can be written over the connection.




### <a name="DLConn.WriteMS">func</a> (\*DLConn) [WriteMS](/src/target/dl_conn.go?s=2196:2239#L100)
``` go
func (d *DLConn) WriteMS(data []byte) error
```
WriteMS sends a raw miniseed packet to a datalink service, it will be decoded to update the
required message identification values required.




## <a name="DLPacket">type</a> [DLPacket](/src/target/dl_types.go?s=73:156#L12)
``` go
type DLPacket struct {
    Preheader DLPreheader
    Header    []byte
    Body      []byte
}

```









## <a name="DLPreheader">type</a> [DLPreheader](/src/target/dl_types.go?s=282:429#L26)
``` go
type DLPreheader struct {
    DL           [2]byte //ASCII String == "DL"
    HeaderLength uint8   //1 byte describing the length of rest of the header
}

```






### <a name="UnmarshalDLPreheader">func</a> [UnmarshalDLPreheader](/src/target/dl_types.go?s=431:497#L31)
``` go
func UnmarshalDLPreheader(b [DLPreheaderSize]byte) (r DLPreheader)
```




## <a name="DLink">type</a> [DLink](/src/target/dl_dlink.go?s=84:179#L8)
``` go
type DLink struct {
    // contains filtered or unexported fields
}

```
DLink is a wrapper around a DLConn connection.







### <a name="NewDLink">func</a> [NewDLink](/src/target/dl_dlink.go?s=865:918#L42)
``` go
func NewDLink(server string, opts ...DLinkOpt) *DLink
```
NewDLink returns a DLink pointer for the given server, optional settings can be passed
as DLinkOpt functions.





### <a name="DLink.Connect">func</a> (\*DLink) [Connect](/src/target/dl_dlink.go?s=1545:1587#L71)
``` go
func (d *DLink) Connect() (*DLConn, error)
```
Connect returns a DLConn pointer on a successful connection to a datalink server.




### <a name="DLink.SetProgram">func</a> (\*DLink) [SetProgram](/src/target/dl_dlink.go?s=1281:1317#L61)
``` go
func (d *DLink) SetProgram(s string)
```
SetProgram sets the program name used for connection requests.




### <a name="DLink.SetTimeout">func</a> (\*DLink) [SetTimeout](/src/target/dl_dlink.go?s=1151:1194#L56)
``` go
func (d *DLink) SetTimeout(t time.Duration)
```
SetTimeout sets the timeout value used for connection requests.




### <a name="DLink.SetUsername">func</a> (\*DLink) [SetUsername](/src/target/dl_dlink.go?s=1401:1438#L66)
``` go
func (d *DLink) SetUsername(s string)
```
SetUsername sets the username used for connection requests.




## <a name="DLinkOpt">type</a> [DLinkOpt](/src/target/dl_dlink.go?s=246:272#L17)
``` go
type DLinkOpt func(*DLink)
```
DLinkOpt is a function for setting DLink internal parameters.







### <a name="SetDLProgram">func</a> [SetDLProgram](/src/target/dl_dlink.go?s=514:550#L27)
``` go
func SetDLProgram(s string) DLinkOpt
```
SetDLProgram sets the program name for datalink connections.


### <a name="SetDLTimeout">func</a> [SetDLTimeout](/src/target/dl_dlink.go?s=357:400#L20)
``` go
func SetDLTimeout(t time.Duration) DLinkOpt
```
SetDLTimeout sets the timeout for seedlink server commands and packet requests.


### <a name="SetDLUsername">func</a> [SetDLUsername](/src/target/dl_dlink.go?s=661:698#L34)
``` go
func SetDLUsername(s string) DLinkOpt
```
SetDLUsername sets the username for datalink connections.





## <a name="MSDataHeader">type</a> [MSDataHeader](/src/target/ms_types.go?s=1355:2156#L74)
``` go
type MSDataHeader struct {
    SeqNumber   [6]byte //ASCII String representing a 7 digit number
    DataQuality byte    //ASCI: D, R, Q or M
    Reserved    byte

    //These are ascii strings
    Station  [5]byte
    Location [2]byte
    Channel  [3]byte
    Network  [2]byte

    StartTime            BTime
    SampleCount          uint16 // Number of Samples in the data block which may or may not be unpacked.
    SampleRateFactor     int16  // >0: Samples/Second <0: Second/Samples 0: Seconds/Sample, ASCII/OPAQUE DATA records
    SampleRateMultiplier int16  // >0: Multiplication Factor <0: Division Factor

    //flags are bit flags
    ActivityFlags    byte
    IOClockFlags     byte
    DataQualityFlags byte

    BlockettesToFollow uint8
    TimeCorrection     int32 // 0.0001 second units
    BeginningOfData    uint16
    FirstBlockette     uint16
}

```






### <a name="UnmarshalMSDataHeader">func</a> [UnmarshalMSDataHeader](/src/target/ms_types.go?s=2158:2227#L101)
``` go
func UnmarshalMSDataHeader(b [MSDataHeaderSize]byte) (r MSDataHeader)
```




## <a name="MSEncoding">type</a> [MSEncoding](/src/target/ms_record.go?s=92:113#L13)
``` go
type MSEncoding uint8
```

``` go
const (
    EncodingASCII      MSEncoding = 0
    EncodingInt16      MSEncoding = 1
    EncodingInt24      MSEncoding = 2
    EncodingInt32      MSEncoding = 3
    EncodingIEEEFloat  MSEncoding = 4
    EncodingIEEEDouble MSEncoding = 5

    EncodingSTEIM1            MSEncoding = 10
    EncodingSTEIM2            MSEncoding = 11
    EncodingGEOSCOPE24bit     MSEncoding = 12
    EncodingGEOSCOPE16bit3exp MSEncoding = 13
    EncodingGEOSCOPE16bit4exp MSEncoding = 14
    EncodingUSNN              MSEncoding = 15
    EncodingCDSN              MSEncoding = 16
    EncodingGraefenberg       MSEncoding = 17
    EncodingIPG               MSEncoding = 18
    EncodingSTEIM3            MSEncoding = 19

    EncodingSRO    MSEncoding = 30
    EncodingHGLP   MSEncoding = 31
    EncodingDWWSSN MSEncoding = 32
    EncodingRSTN   MSEncoding = 33
)
```









## <a name="MSRecord">type</a> [MSRecord](/src/target/ms_record.go?s=1226:1405#L58)
``` go
type MSRecord struct {
    Header MSDataHeader
    B1000  Blockette1000 //If Present
    B1001  Blockette1001 //If Present
    Data   []byte
    // contains filtered or unexported fields
}

```
MSRecord represents the raw miniseed record.







### <a name="NewMSRecord">func</a> [NewMSRecord](/src/target/ms_record.go?s=1939:1986#L87)
``` go
func NewMSRecord(buf []byte) (*MSRecord, error)
```
NewMSRecord decodes and unpacks the record samples from a byte slice and
returns an MSRecord pointer, or an empty pointer and an error if it
could not be decoded.





### <a name="MSRecord.BlockSize">func</a> (\*MSRecord) [BlockSize](/src/target/ms_record.go?s=16762:16796#L539)
``` go
func (m *MSRecord) BlockSize() int
```
PacketSize returns the length of the packet




### <a name="MSRecord.ByteOrder">func</a> (\*MSRecord) [ByteOrder](/src/target/ms_record.go?s=15092:15132#L479)
``` go
func (m *MSRecord) ByteOrder() WordOrder
```
ByteOrder returns the miniseed data byte order.




### <a name="MSRecord.Channel">func</a> (\*MSRecord) [Channel](/src/target/ms_record.go?s=12618:12653#L387)
``` go
func (m *MSRecord) Channel() string
```
Channel returns the cleaned miniseed header channel string.




### <a name="MSRecord.DataQuality">func</a> (\*MSRecord) [DataQuality](/src/target/ms_record.go?s=12776:12813#L392)
``` go
func (m *MSRecord) DataQuality() byte
```
DataQuality returns the miniseed header data quality flag.




### <a name="MSRecord.Encoding">func</a> (\*MSRecord) [Encoding](/src/target/ms_record.go?s=14958:14998#L474)
``` go
func (m *MSRecord) Encoding() MSEncoding
```
Encoding returns the miniseed data format encoding.




### <a name="MSRecord.EndTime">func</a> (\*MSRecord) [EndTime](/src/target/ms_record.go?s=15979:16017#L518)
``` go
func (m *MSRecord) EndTime() time.Time
```
EndTime returns the calculated time of the last sample.




### <a name="MSRecord.Location">func</a> (\*MSRecord) [Location](/src/target/ms_record.go?s=12457:12493#L382)
``` go
func (m *MSRecord) Location() string
```
Location returns the cleaned miniseed header location string.




### <a name="MSRecord.Network">func</a> (\*MSRecord) [Network](/src/target/ms_record.go?s=12137:12172#L372)
``` go
func (m *MSRecord) Network() string
```
Network returns the cleaned miniseed header network string.




### <a name="MSRecord.NumSamples">func</a> (\*MSRecord) [NumSamples](/src/target/ms_record.go?s=15385:15420#L493)
``` go
func (m *MSRecord) NumSamples() int
```
NumSamples returns the number of samples decoded from the miniseed record.




### <a name="MSRecord.SampleCount">func</a> (\*MSRecord) [SampleCount](/src/target/ms_record.go?s=14827:14863#L469)
``` go
func (m *MSRecord) SampleCount() int
```
SampleCount returns the number of samples in the record, independent of whether they are decoded or not.




### <a name="MSRecord.SamplePeriod">func</a> (\*MSRecord) [SamplePeriod](/src/target/ms_record.go?s=16573:16620#L534)
``` go
func (m *MSRecord) SamplePeriod() time.Duration
```
SamplePeriod converts the sample rate into a time interval.
For invalid sampling rates a zero duration is returned.




### <a name="MSRecord.SampleRate">func</a> (\*MSRecord) [SampleRate](/src/target/ms_record.go?s=14587:14626#L464)
``` go
func (m *MSRecord) SampleRate() float64
```
SampleRate returns the decoded header sampling rate in samples per second.




### <a name="MSRecord.SampleType">func</a> (\*MSRecord) [SampleType](/src/target/ms_record.go?s=15777:15821#L510)
``` go
func (m *MSRecord) SampleType() MSSampleType
```
SampleType returns the type of samples decoded, or UnknownType if no data has been decoded.




### <a name="MSRecord.SequenceNumber">func</a> (\*MSRecord) [SequenceNumber](/src/target/ms_record.go?s=11774:11813#L362)
``` go
func (m *MSRecord) SequenceNumber() int
```
SequenceNumber returns the record sequence number.




### <a name="MSRecord.SrcName">func</a> (\*MSRecord) [SrcName](/src/target/ms_record.go?s=19295:19342#L633)
``` go
func (m *MSRecord) SrcName(quality bool) string
```
SrcName returns the standard stream representation of the packet header.




### <a name="MSRecord.StartTime">func</a> (\*MSRecord) [StartTime](/src/target/ms_record.go?s=12902:12942#L397)
``` go
func (m *MSRecord) StartTime() time.Time
```
StartTime returns the corrected header start time.




### <a name="MSRecord.Station">func</a> (\*MSRecord) [Station](/src/target/ms_record.go?s=12296:12331#L377)
``` go
func (m *MSRecord) Station() string
```
Station returns the cleaned miniseed header station string.




### <a name="MSRecord.String">func</a> (\*MSRecord) [String](/src/target/ms_record.go?s=2202:2236#L99)
``` go
func (m *MSRecord) String() string
```
String implements the Stringer interface and provides a short summary of the
miniseed record header.




### <a name="MSRecord.ToBytes">func</a> (\*MSRecord) [ToBytes](/src/target/ms_record.go?s=17445:17480#L566)
``` go
func (m *MSRecord) ToBytes() []byte
```
ToBytes returns a cleaned byte slice of the expected ASCII record.




### <a name="MSRecord.ToFloat64s">func</a> (\*MSRecord) [ToFloat64s](/src/target/ms_record.go?s=18592:18633#L611)
``` go
func (m *MSRecord) ToFloat64s() []float64
```
ToFloat64s returns the record numerical data converted to an float64 slice.




### <a name="MSRecord.ToInt32s">func</a> (\*MSRecord) [ToInt32s](/src/target/ms_record.go?s=18372:18409#L602)
``` go
func (m *MSRecord) ToInt32s() []int32
```



### <a name="MSRecord.ToInts">func</a> (\*MSRecord) [ToInts](/src/target/ms_record.go?s=17719:17752#L576)
``` go
func (m *MSRecord) ToInts() []int
```
ToInts returns the record numerical data converted to an int slice.




### <a name="MSRecord.ToStrings">func</a> (\*MSRecord) [ToStrings](/src/target/ms_record.go?s=16949:16988#L547)
``` go
func (m *MSRecord) ToStrings() []string
```
ToStrings returns an ASCII encodied record as a slice of line strings.




### <a name="MSRecord.Unpack">func</a> (\*MSRecord) [Unpack](/src/target/ms_record.go?s=7117:7175#L234)
``` go
func (m *MSRecord) Unpack(buf []byte, dataflag bool) error
```
Unpack the record form a byte slice, the dataflag can be used to suppress decoding the waveform data
for efficency if only the header information is required.




## <a name="MSRecordFunc">type</a> [MSRecordFunc](/src/target/ms_pack.go?s=1794:1847#L68)
``` go
type MSRecordFunc func(*MSRecord, []byte, bool) error
```
MSRecordFunc is used to process packed miniseed blocks.










## <a name="MSRecordStream">type</a> [MSRecordStream](/src/target/ms_pack.go?s=144:439#L11)
``` go
type MSRecordStream struct {
    // stream details
    Station     string
    Location    string
    Channel     string
    Network     string
    DataQuality byte

    // sampling rate
    Factor     int
    Multiplier int

    // block details
    WordOrder    WordOrder
    RecordLength int

    // running details
    SeqNumber int
}

```
MSRecordStream is used as a template to encode miniseed blocks.







### <a name="NewMSRecordStream">func</a> [NewMSRecordStream](/src/target/ms_pack.go?s=1226:1279#L52)
``` go
func NewMSRecordStream(msr *MSRecord) *MSRecordStream
```
NewMSRecordStream builds a MSRecordStream pointer from an MSRecord pointer.





### <a name="MSRecordStream.PackInt32">func</a> (\*MSRecordStream) [PackInt32](/src/target/ms_pack.go?s=2565:2679#L101)
``` go
func (m *MSRecordStream) PackInt32(start time.Time, quality int, locked bool, data []int32, fn MSRecordFunc) error
```



## <a name="MSSampleType">type</a> [MSSampleType](/src/target/ms_record.go?s=985:1007#L47)
``` go
type MSSampleType byte
```

``` go
const (
    UnknownType MSSampleType = 0
    ByteType    MSSampleType = 'a'
    IntegerType MSSampleType = 'i'
    FloatType   MSSampleType = 'f'
    DoubleType  MSSampleType = 'd'
)
```









## <a name="SLCollectFunc">type</a> [SLCollectFunc](/src/target/sl_slink.go?s=4708:4761#L190)
``` go
type SLCollectFunc func(string, []byte) (bool, error)
```
SLCollectFunc is a function run on each returned seedlink packet. It should return a true value
to stop collecting data without an error message. A non-nil returned error will also stop
collection but with an assumed errored state.










## <a name="SLConn">type</a> [SLConn](/src/target/sl_conn.go?s=1494:1643#L57)
``` go
type SLConn struct {
    net.Conn
    // contains filtered or unexported fields
}

```






### <a name="NewSLConn">func</a> [NewSLConn](/src/target/sl_conn.go?s=1837:1907#L71)
``` go
func NewSLConn(service string, timeout time.Duration) (*SLConn, error)
```
NewSLConn returns a new connection to the named seedlink server with a given command timeout. It is expected that the
Close function be called when the connection is no longer required.





### <a name="SLConn.Collect">func</a> (\*SLConn) [Collect](/src/target/sl_conn.go?s=9519:9564#L384)
``` go
func (c *SLConn) Collect() (*SLPacket, error)
```
Collect returns a seedlink packet if available within the optional timout. Any error returned should be
checked that it isn't a timeout, this should be handled as appropriate for the request.




### <a name="SLConn.CommandCat">func</a> (\*SLConn) [CommandCat](/src/target/sl_conn.go?s=6761:6806#L296)
``` go
func (c *SLConn) CommandCat() ([]byte, error)
```
CommandStationList sends a CAT command to the seedlink server.




### <a name="SLConn.CommandClose">func</a> (\*SLConn) [CommandClose](/src/target/sl_conn.go?s=6608:6655#L291)
``` go
func (c *SLConn) CommandClose() ([]byte, error)
```
CommandClose sends a BYE command to the seedlink server.




### <a name="SLConn.CommandData">func</a> (\*SLConn) [CommandData](/src/target/sl_conn.go?s=8142:8214#L333)
``` go
func (c *SLConn) CommandData(sequence string, starttime time.Time) error
```
CommandData sends a DATA command to the seedlink server.




### <a name="SLConn.CommandEnd">func</a> (\*SLConn) [CommandEnd](/src/target/sl_conn.go?s=9164:9199#L375)
``` go
func (c *SLConn) CommandEnd() error
```
CommandEnd sends an END command to the seedlink server.




### <a name="SLConn.CommandHello">func</a> (\*SLConn) [CommandHello](/src/target/sl_conn.go?s=6461:6508#L286)
``` go
func (c *SLConn) CommandHello() ([]byte, error)
```
CommandHello sends a HELLO command to the seedlink server.




### <a name="SLConn.CommandId">func</a> (\*SLConn) [CommandId](/src/target/sl_conn.go?s=6314:6358#L281)
``` go
func (c *SLConn) CommandId() ([]byte, error)
```
CommandId sends an INFO ID command to the seedlink server.




### <a name="SLConn.CommandSelect">func</a> (\*SLConn) [CommandSelect](/src/target/sl_conn.go?s=7849:7903#L323)
``` go
func (c *SLConn) CommandSelect(selection string) error
```
CommandSelect sends a SELECT command to the seedlink server.




### <a name="SLConn.CommandStation">func</a> (\*SLConn) [CommandStation](/src/target/sl_conn.go?s=6910:6972#L301)
``` go
func (c *SLConn) CommandStation(station, network string) error
```
CommandStation sends a STATION command to the seedlink server.




### <a name="SLConn.CommandTime">func</a> (\*SLConn) [CommandTime](/src/target/sl_conn.go?s=8636:8700#L353)
``` go
func (c *SLConn) CommandTime(starttime, endtime time.Time) error
```
CommandTime sends a TIME command to the seedlink server.




### <a name="SLConn.GetInfo">func</a> (\*SLConn) [GetInfo](/src/target/sl_conn.go?s=5521:5575#L252)
``` go
func (c *SLConn) GetInfo(level string) ([]byte, error)
```
GetInfoLevel requests the seedlink server return an INFO request for the given level.




### <a name="SLConn.GetSLInfo">func</a> (\*SLConn) [GetSLInfo](/src/target/sl_conn.go?s=6015:6072#L266)
``` go
func (c *SLConn) GetSLInfo(level string) (*SLInfo, error)
```
GetInfo requests the seedlink server return an INFO request for the given level. The results
are returned as a decoded SLInfo pointer, or an error otherwise.




## <a name="SLInfo">type</a> [SLInfo](/src/target/sl_info.go?s=42:862#L7)
``` go
type SLInfo struct {
    XMLName xml.Name `xml:"seedlink"`

    Software     string `xml:"software,attr"`
    Organization string `xml:"organization,attr"`
    Started      string `xml:"started,attr"`
    Capability   []struct {
        Name string `xml:"name,attr"`
    } `xml:"capability"`
    Station []struct {
        Name        string `xml:"name,attr"`
        Network     string `xml:"network,attr"`
        Description string `xml:"description,attr"`
        BeginSeq    string `xml:"begin_seq,attr"`
        EndSeq      string `xml:"end_seq,attr"`
        StreamCheck string `xml:"stream_check,attr"`
        Stream      []struct {
            Location  string `xml:"location,attr"`
            Seedname  string `xml:"seedname,attr"`
            Type      string `xml:"type,attr"`
            BeginTime string `xml:"begin_time,attr"`
            EndTime   string `xml:"end_time,attr"`
        } `xml:"stream"`
    } `xml:"station"`
}

```









### <a name="SLInfo.Unmarshal">func</a> (\*SLInfo) [Unmarshal](/src/target/sl_info.go?s=864:909#L33)
``` go
func (s *SLInfo) Unmarshal(data []byte) error
```



## <a name="SLPacket">type</a> [SLPacket](/src/target/sl_packet.go?s=68:211#L11)
``` go
type SLPacket struct {
    SL   [2]byte   // ASCII String == "SL"
    Seq  [6]byte   // ASCII sequence number
    Data [512]byte // Fixed size payload
}

```






### <a name="NewSLPacket">func</a> [NewSLPacket](/src/target/sl_packet.go?s=213:261#L17)
``` go
func NewSLPacket(data []byte) (*SLPacket, error)
```




## <a name="SLState">type</a> [SLState](/src/target/sl_state.go?s=653:741#L29)
``` go
type SLState struct {
    // contains filtered or unexported fields
}

```
SLState maintains the current state information for a seedlink connection.










### <a name="SLState.Add">func</a> (\*SLState) [Add](/src/target/sl_state.go?s=1376:1416#L62)
``` go
func (s *SLState) Add(station SLStation)
```
Add inserts or updates the station collection details into the connection state.




### <a name="SLState.Find">func</a> (\*SLState) [Find](/src/target/sl_state.go?s=1691:1739#L75)
``` go
func (s *SLState) Find(stn SLStation) *SLStation
```



### <a name="SLState.Marshal">func</a> (\*SLState) [Marshal](/src/target/sl_state.go?s=2223:2266#L106)
``` go
func (s *SLState) Marshal() ([]byte, error)
```



### <a name="SLState.ReadFile">func</a> (\*SLState) [ReadFile](/src/target/sl_state.go?s=2387:2432#L116)
``` go
func (s *SLState) ReadFile(path string) error
```



### <a name="SLState.Stations">func</a> (\*SLState) [Stations](/src/target/sl_state.go?s=816:856#L37)
``` go
func (s *SLState) Stations() []SLStation
```
Stations returns a sorted slice of current station state information.




### <a name="SLState.Unmarshal">func</a> (\*SLState) [Unmarshal](/src/target/sl_state.go?s=2013:2059#L92)
``` go
func (s *SLState) Unmarshal(data []byte) error
```



### <a name="SLState.WriteFile">func</a> (\*SLState) [WriteFile](/src/target/sl_state.go?s=2615:2661#L134)
``` go
func (s *SLState) WriteFile(path string) error
```



## <a name="SLStation">type</a> [SLStation](/src/target/sl_state.go?s=184:364#L13)
``` go
type SLStation struct {
    Network   string    `json:"network"`
    Station   string    `json:"station"`
    Sequence  int       `json:"sequence"`
    Timestamp time.Time `json:"timestamp"`
}

```
SLStation stores the latest state information for the given network and station combination.










### <a name="SLStation.Key">func</a> (SLStation) [Key](/src/target/sl_state.go?s=469:503#L21)
``` go
func (s SLStation) Key() SLStation
```
Key returns a blank SLStation except for the Network and Station entries, this useful as a map key.




## <a name="SLink">type</a> [SLink](/src/target/sl_slink.go?s=169:442#L13)
``` go
type SLink struct {
    // contains filtered or unexported fields
}

```
SLink is a wrapper around an SLConn to provide
handling of timeouts and keep alive messages.







### <a name="NewSLink">func</a> [NewSLink](/src/target/sl_slink.go?s=2612:2665#L114)
``` go
func NewSLink(server string, opts ...SLinkOpt) *SLink
```
NewSlink returns a SLink pointer for the given server, optional settings can be passed
as SLinkOpt functions.





### <a name="SLink.Collect">func</a> (\*SLink) [Collect](/src/target/sl_slink.go?s=8759:8806#L337)
``` go
func (s *SLink) Collect(fn SLCollectFunc) error
```
Collect calls CollectWithContext with a background Context and a handler function.




### <a name="SLink.CollectWithContext">func</a> (\*SLink) [CollectWithContext](/src/target/sl_slink.go?s=5602:5681#L201)
``` go
func (s *SLink) CollectWithContext(ctx context.Context, fn SLCollectFunc) error
```
CollectWithContext makes a connection to the seedlink server, recovers initial client information and
the sets the connection into streaming mode. Recovered packets are passed to a given function
to process, if this function returns a true value or a non-nil error value the collection will
stop and the function will return.
If a call returns with a timeout error a check is made whether a keepalive is needed or whether
the function should return as no data has been received for an extended period of time. It is
assumed the calling function will attempt a reconnection with an updated set of options, specifically
any start or end time parameters. The Context parameter can be used to to cancel the data collection
independent of the function as this may never be called if no appropriate has been received.




### <a name="SLink.SetEndTime">func</a> (\*SLink) [SetEndTime](/src/target/sl_slink.go?s=3954:3993#L168)
``` go
func (s *SLink) SetEndTime(t time.Time)
```
SetEndTime sets the initial end time of the request.




### <a name="SLink.SetKeepAlive">func</a> (\*SLink) [SetKeepAlive](/src/target/sl_slink.go?s=3563:3608#L153)
``` go
func (s *SLink) SetKeepAlive(d time.Duration)
```
SetKeepAlive sets the time interval needed without any packets for
a check message is sent.




### <a name="SLink.SetNetTo">func</a> (\*SLink) [SetNetTo](/src/target/sl_slink.go?s=3405:3446#L147)
``` go
func (s *SLink) SetNetTo(d time.Duration)
```
SetNetTo sets the overall timeout after which a reconnection is tried.




### <a name="SLink.SetRefresh">func</a> (\*SLink) [SetRefresh](/src/target/sl_slink.go?s=3145:3188#L137)
``` go
func (s *SLink) SetRefresh(d time.Duration)
```
SetRefresh sets the interval used for state refreshes if enabled.




### <a name="SLink.SetSelectors">func</a> (\*SLink) [SetSelectors](/src/target/sl_slink.go?s=4233:4279#L178)
``` go
func (s *SLink) SetSelectors(selectors string)
```
SetSelectors sets the channel selectors used for seedlink connections.




### <a name="SLink.SetSequence">func</a> (\*SLink) [SetSequence](/src/target/sl_slink.go?s=3695:3736#L158)
``` go
func (s *SLink) SetSequence(sequence int)
```
SetSequence sets the start sequence for the initial request.




### <a name="SLink.SetStartTime">func</a> (\*SLink) [SetStartTime](/src/target/sl_slink.go?s=3828:3869#L163)
``` go
func (s *SLink) SetStartTime(t time.Time)
```
SetStartTime sets the initial starting time of the request.




### <a name="SLink.SetState">func</a> (\*SLink) [SetState](/src/target/sl_slink.go?s=4374:4421#L183)
``` go
func (s *SLink) SetState(stations ...SLStation)
```
SetState sets the initial list of station state information.




### <a name="SLink.SetStateFile">func</a> (\*SLink) [SetStateFile](/src/target/sl_slink.go?s=3270:3308#L142)
``` go
func (s *SLink) SetStateFile(f string)
```
SetStateFile sets the file for storing state information.




### <a name="SLink.SetStreams">func</a> (\*SLink) [SetStreams](/src/target/sl_slink.go?s=4090:4132#L173)
``` go
func (s *SLink) SetStreams(streams string)
```
SetStreams sets the channel streams used for seedlink connections.




### <a name="SLink.SetTimeout">func</a> (\*SLink) [SetTimeout](/src/target/sl_slink.go?s=3012:3055#L132)
``` go
func (s *SLink) SetTimeout(d time.Duration)
```
SetTimeout sets the timeout value used for connection requests.




## <a name="SLinkOpt">type</a> [SLinkOpt](/src/target/sl_slink.go?s=506:532#L33)
``` go
type SLinkOpt func(*SLink)
```
SLink is a function for setting SLink internal parameters.







### <a name="SetSLEndTime">func</a> [SetSLEndTime](/src/target/sl_slink.go?s=1836:1875#L85)
``` go
func SetSLEndTime(t time.Time) SLinkOpt
```
SetSLEndTime sets the end of the initial request from the seedlink server.


### <a name="SetSLKeepAlive">func</a> [SetSLKeepAlive](/src/target/sl_slink.go?s=1315:1360#L64)
``` go
func SetSLKeepAlive(d time.Duration) SLinkOpt
```
SetSLKeepAlive sets the time to send an ID message to server if no packets have been received.


### <a name="SetSLNetTo">func</a> [SetSLNetTo](/src/target/sl_slink.go?s=1128:1169#L57)
``` go
func SetSLNetTo(d time.Duration) SLinkOpt
```
SetSLNetTo sets the time to after which the connection is closed after no packets have been received.


### <a name="SetSLRefresh">func</a> [SetSLRefresh](/src/target/sl_slink.go?s=777:820#L43)
``` go
func SetSLRefresh(d time.Duration) SLinkOpt
```
SetSLRefresh sets the interval for state refreshing if enabled.


### <a name="SetSLSelectors">func</a> [SetSLSelectors](/src/target/sl_slink.go?s=2202:2248#L99)
``` go
func SetSLSelectors(selectors string) SLinkOpt
```
SetSLSelectors sets the default list of selectors to use for seedlink stream requests.


### <a name="SetSLSequence">func</a> [SetSLSequence](/src/target/sl_slink.go?s=1478:1519#L71)
``` go
func SetSLSequence(sequence int) SLinkOpt
```
SetSLSequence sets the start sequence for the initial request.


### <a name="SetSLStartTime">func</a> [SetSLStartTime](/src/target/sl_slink.go?s=1659:1700#L78)
``` go
func SetSLStartTime(t time.Time) SLinkOpt
```
SetSLStartTime sets the start of the initial request from the seedlink server.


### <a name="SetSLState">func</a> [SetSLState](/src/target/sl_slink.go?s=2374:2421#L106)
``` go
func SetSLState(stations ...SLStation) SLinkOpt
```
SetSLState sets the default list of station state information.


### <a name="SetSLStateFile">func</a> [SetSLStateFile](/src/target/sl_slink.go?s=933:971#L50)
``` go
func SetSLStateFile(f string) SLinkOpt
```
SetSLStateFile sets the file for storing state information.


### <a name="SetSLStreams">func</a> [SetSLStreams](/src/target/sl_slink.go?s=2014:2056#L92)
``` go
func SetSLStreams(streams string) SLinkOpt
```
SetSLStreams sets the list of stations and streams to from the seedlink server.


### <a name="SetSLTimeout">func</a> [SetSLTimeout](/src/target/sl_slink.go?s=617:660#L36)
``` go
func SetSLTimeout(d time.Duration) SLinkOpt
```
SetSLTimeout sets the timeout for seedlink server commands and packet requests.





## <a name="WordOrder">type</a> [WordOrder](/src/target/ms_record.go?s=896:916#L40)
``` go
type WordOrder uint8
```

``` go
const (
    LittleEndian WordOrder = 0
    BigEndian    WordOrder = 1
)
```













- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
