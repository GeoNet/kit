

# dl
`import "github.com/GeoNet/kit/seis/dl"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
The dl module has been writen as a lightweight replacement for the C libdali library.
It is aimed at clients that need to connect to a datalink server, either requesting inforamtion
or for uploading most likely miniseed records.




## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [func MarshalPreheader(r Preheader) (b [PreheaderSize]byte)](#MarshalPreheader)
* [type DLConn](#DLConn)
  * [func NewDLConn(service string, timeout time.Duration) (*DLConn, error)](#NewDLConn)
  * [func (d *DLConn) Id() string](#DLConn.Id)
  * [func (d *DLConn) SetId(program, username string) error](#DLConn.SetId)
  * [func (d *DLConn) Size() int](#DLConn.Size)
  * [func (d *DLConn) Writable() bool](#DLConn.Writable)
  * [func (d *DLConn) WriteMS(srcname string, start, end time.Time, data []byte) error](#DLConn.WriteMS)
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
* [type Packet](#Packet)
* [type Preheader](#Preheader)
  * [func UnmarshalPreheader(b [PreheaderSize]byte) (r Preheader)](#UnmarshalPreheader)


#### <a name="pkg-files">Package files</a>
[conn.go](/src/target/conn.go) [dlink.go](/src/target/dlink.go) [doc.go](/src/target/doc.go) [types.go](/src/target/types.go) 


## <a name="pkg-constants">Constants</a>
``` go
const (
    PreheaderSize = 3
)
```



## <a name="MarshalPreheader">func</a> [MarshalPreheader](/src/target/types.go?s=537:595#L37)
``` go
func MarshalPreheader(r Preheader) (b [PreheaderSize]byte)
```



## <a name="DLConn">type</a> [DLConn](/src/target/conn.go?s=198:300#L19)
``` go
type DLConn struct {
    net.Conn
    // contains filtered or unexported fields
}

```
DLConn provides connection information to a datalink service.







### <a name="NewDLConn">func</a> [NewDLConn](/src/target/conn.go?s=489:559#L31)
``` go
func NewDLConn(service string, timeout time.Duration) (*DLConn, error)
```
NewDLConn makes a connection to a datalink service, the function
SetId should be run after the connection is made and Close
should be called when the link is no longer required.





### <a name="DLConn.Id">func</a> (\*DLConn) [Id](/src/target/conn.go?s=872:900#L49)
``` go
func (d *DLConn) Id() string
```
Id returns the client identification sent to the server.




### <a name="DLConn.SetId">func</a> (\*DLConn) [SetId](/src/target/conn.go?s=1240:1294#L64)
``` go
func (d *DLConn) SetId(program, username string) error
```
SetId sens an ID message to the remote connection and decodes the connection capabilities.




### <a name="DLConn.Size">func</a> (\*DLConn) [Size](/src/target/conn.go?s=1098:1125#L59)
``` go
func (d *DLConn) Size() int
```
Size returns the current expected packet size.




### <a name="DLConn.Writable">func</a> (\*DLConn) [Writable](/src/target/conn.go?s=991:1023#L54)
``` go
func (d *DLConn) Writable() bool
```
Writable returns whether packets can be written over the connection.




### <a name="DLConn.WriteMS">func</a> (\*DLConn) [WriteMS](/src/target/conn.go?s=2214:2295#L100)
``` go
func (d *DLConn) WriteMS(srcname string, start, end time.Time, data []byte) error
```
WriteMS sends a raw miniseed packet to a datalink service, it will need to be decoded
prior to sending to allow for the required message identification values to added.




## <a name="DLink">type</a> [DLink](/src/target/dlink.go?s=82:177#L8)
``` go
type DLink struct {
    // contains filtered or unexported fields
}

```
DLink is a wrapper around a DLConn connection.







### <a name="NewDLink">func</a> [NewDLink](/src/target/dlink.go?s=863:916#L42)
``` go
func NewDLink(server string, opts ...DLinkOpt) *DLink
```
NewDLink returns a DLink pointer for the given server, optional settings can be passed
as DLinkOpt functions.





### <a name="DLink.Connect">func</a> (\*DLink) [Connect](/src/target/dlink.go?s=1543:1585#L71)
``` go
func (d *DLink) Connect() (*DLConn, error)
```
Connect returns a DLConn pointer on a successful connection to a datalink server.




### <a name="DLink.SetProgram">func</a> (\*DLink) [SetProgram](/src/target/dlink.go?s=1279:1315#L61)
``` go
func (d *DLink) SetProgram(s string)
```
SetProgram sets the program name used for connection requests.




### <a name="DLink.SetTimeout">func</a> (\*DLink) [SetTimeout](/src/target/dlink.go?s=1149:1192#L56)
``` go
func (d *DLink) SetTimeout(t time.Duration)
```
SetTimeout sets the timeout value used for connection requests.




### <a name="DLink.SetUsername">func</a> (\*DLink) [SetUsername](/src/target/dlink.go?s=1399:1436#L66)
``` go
func (d *DLink) SetUsername(s string)
```
SetUsername sets the username used for connection requests.




## <a name="DLinkOpt">type</a> [DLinkOpt](/src/target/dlink.go?s=244:270#L17)
``` go
type DLinkOpt func(*DLink)
```
DLinkOpt is a function for setting DLink internal parameters.







### <a name="SetDLProgram">func</a> [SetDLProgram](/src/target/dlink.go?s=512:548#L27)
``` go
func SetDLProgram(s string) DLinkOpt
```
SetDLProgram sets the program name for datalink connections.


### <a name="SetDLTimeout">func</a> [SetDLTimeout](/src/target/dlink.go?s=355:398#L20)
``` go
func SetDLTimeout(t time.Duration) DLinkOpt
```
SetDLTimeout sets the timeout for seedlink server commands and packet requests.


### <a name="SetDLUsername">func</a> [SetDLUsername](/src/target/dlink.go?s=659:696#L34)
``` go
func SetDLUsername(s string) DLinkOpt
```
SetDLUsername sets the username for datalink connections.





## <a name="Packet">type</a> [Packet](/src/target/types.go?s=69:148#L12)
``` go
type Packet struct {
    Preheader Preheader
    Header    []byte
    Body      []byte
}

```









## <a name="Preheader">type</a> [Preheader](/src/target/types.go?s=270:415#L26)
``` go
type Preheader struct {
    DL           [2]byte //ASCII String == "DL"
    HeaderLength uint8   //1 byte describing the length of rest of the header
}

```






### <a name="UnmarshalPreheader">func</a> [UnmarshalPreheader](/src/target/types.go?s=417:477#L31)
``` go
func UnmarshalPreheader(b [PreheaderSize]byte) (r Preheader)
```








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
