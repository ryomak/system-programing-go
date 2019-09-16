 package main

 import(
   "encoding/binary"
  "fmt"
  "io"
  "os"
 )

 func dumpChunk(chunk io.Reader) {
   var length int32
   binary.Read(chunk,binary.BigEndian,&length)
   buffer := make([]byte ,4)
   chunk.Read(buffer)
   fmt.Printf("chunk %s (%d bytes)\n",string(buffer),length)
 }

 func readChunks(file *os.File) []io.Reader {
   var chunks []io.Reader
   var offset int64 = 8

   file.Seek(8,0)

   for {
     var length int32
     if err := binary.Read(file,binary.BigEndian,&length);err == io.EOF {
       break
     }
    chunks = append(chunks,io.NewSectionReader(file,offset,int64(length)+12))
    offset,_ = file.Seek(int64(length+8),1)
   }
   return chunks
 }


 func main() {
   file ,err := os.Open(os.Args[1]) 
   if err != nil {
     panic (err)
   }
   defer file.Close()
   chunks := readChunks(file)
   for _,v := range chunks {
     dumpChunk(v)
   }
 }
