package bomcomp

import "bufio"
import "bytes"
import "encoding/binary"
import "encoding/csv"
import "fmt"
import "io"
import "log"
import "os"
import "strconv"
import "strings"
import "unsafe"

type Umline struct {
  Fuom, Tuom [3]byte
  Den, Num   float32
}

type Uom struct {
  Size int
  Line []Umline
}

func (u *Uom) NewUom(fname string) {
  u.Line = make([]Umline, 100)
  rlen := len(u.Line[0].Fuom) + len(u.Line[0].Tuom) + int(unsafe.Sizeof(u.Line[0].Den)) + int(unsafe.Sizeof(u.Line[0].Num))
  f, err := os.Open(fname)
  if err != nil {
    log.Fatalf("Open UOM file: %v\n", err)
  }
  defer f.Close()
  stat, err := f.Stat()
  if err != nil {
    log.Fatalf("Get statistics of UOM file: %v\n", err)
  }
  var size int64 = stat.Size()
  tbuf := make([]byte, size)
  rdr := bufio.NewReader(f)
  err = binary.Read(rdr, binary.LittleEndian, tbuf)
  if err != nil {
    log.Fatalf("Read UOM file: %v\n", err)
  }
  u.Size = int(size / int64(rlen))
  for r := 0; r < u.Size; r++ {
    row := tbuf[r*rlen : (r+1)*rlen]
    copy(u.Line[r].Fuom[:], row[0:3])
    copy(u.Line[r].Tuom[:], row[3:6])
    temp := bytes.NewBuffer(row[6:10])
    binary.Read(temp, binary.LittleEndian, &u.Line[r].Den)
    temp = bytes.NewBuffer(row[10:14])
    binary.Read(temp, binary.LittleEndian, &u.Line[r].Num)
  }
}

func (u *Uom) InUom(fUOM, tUOM string) float32 {
  var factor float32 = 0
  for i := 0; i < u.Size && factor == 0; i++ {
    sFuom := u.CToGoString(u.Line[i].Fuom[:])
    sTuom := u.CToGoString(u.Line[i].Tuom[:])
    if fUOM == strings.TrimSpace(sFuom) && tUOM == strings.TrimSpace(sTuom) {
      factor = u.Line[i].Num / u.Line[i].Den
    }
  }
  return factor
}

func (u *Uom) CToGoString(c []byte) string {
  n := -1
  for i, b := range c {
    if b == 0 {
      break
    }
    n = i
  }
  return string(c[:n+1])
}

func (t *Umline) Uomgen(ifnam, ofnam string) {
  ofil, _ := os.Create(ofnam)
  defer ofil.Close()
  ifil, _ := os.Open(ifnam)
  defer ifil.Close()
  rdr := csv.NewReader(ifil)
  flds, err := rdr.Read()
  buf := new(bytes.Buffer)
  for ; err != io.EOF; flds, err = rdr.Read() {
    for i := 0; i < 3; i++ {
      t.Fuom[i], t.Tuom[i] = 0x00, 0x00
    }
    copy(t.Fuom[:], flds[0])
    copy(t.Tuom[:], flds[1])
    temp, _ := strconv.ParseFloat(flds[2], 32)
    t.Den = float32(temp)
    temp, _ = strconv.ParseFloat(flds[3], 32)
    t.Num = float32(temp)
    err = binary.Write(buf, binary.LittleEndian, t)
    if err != nil {
      log.Fatalf("binary.Write failed: %v\n", err)
    }
  }
  out := buf.Bytes()
  ofil.Write(out)
  ofil.Sync()
}

func (t *Umline) Uomdsp(fname string) {
  ifil, _ := os.Open(fname)
  defer ifil.Close()
  ista, _ := ifil.Stat()
  var isze int64 = ista.Size()
  tmp := make([]byte, isze)
  irdr := bufio.NewReader(ifil)
  err := binary.Read(irdr, binary.LittleEndian, tmp)
  if err != nil {
    log.Fatalf("binary.Read failed: %v\n", err)
  }
  lrow := len(t.Fuom) + len(t.Tuom) + int(unsafe.Sizeof(t.Den)) + int(unsafe.Sizeof(t.Num))
  for r := 0; r < int(isze/int64(lrow)); r++ {
    row := tmp[r*lrow : (r+1)*lrow]
    copy(t.Fuom[:], row[0:3])
    copy(t.Tuom[:], row[3:6])
    dbuf := bytes.NewBuffer(row[6:10])
    binary.Read(dbuf, binary.LittleEndian, &t.Den)
    nbuf := bytes.NewBuffer(row[10:14])
    binary.Read(nbuf, binary.LittleEndian, &t.Num)
    fmt.Printf("%s %s %f %f\n", t.Fuom, t.Tuom, t.Den, t.Num)
  }
}
