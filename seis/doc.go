// The seis module has been writen as a lightweight replacement for the C libraries libmseed and libslink.
// It is aimed at clients that need to decode miniseed data either directly or collected from a seedlink
// server.
//
// The seedlink code is not a direct replacement for libslink. It can run in two modes, either as a
// raw connection to the client connection (SLConn) which allows mechanisms to monitor or have a finer
// control of the SeedLink connection, or in the collection mode (SLink) where a connection is established
// and received miniseed blocks can be processed with a call back function. A context can be passed into
// the collection loop to allow interuption or as a shutdown mechanism. It is not passed to the underlying
// seedlink connection messaging which is managed via a deadline mechanism, e.g. the `SetTimeout` option.
//
// An example Seedlink application can be as simple as:
//
//  slink := seis.NewSLink("localhost:18000")
//
//  if err := slink.Collect(func(seq string, data []byte) (bool, error) {
//          if ms, err := seis.NewMSRecord(data); err == nil {
//              log.Println(ms.SrcName(false), time.Since(ms.EndTime()))
//         }
//         return false, nil
//  }); err != nil {
//          log.Fatal(err)
//  }
//
// The conversion to miniseed can be coupled together, is in this example, but it is not required. Either
// the raw packet can be managed as a whole unit or it can be unpacked using another mechanism.
//
package seis
