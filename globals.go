package bomcomp

import "archive/zip"
import "bytes"
import _ "code.google.com/p/odbc"
import "database/sql"
import "fmt"
import "io"
import "log"
import "os"
import "strings"

var S Settings
var U Uom
var Q SqlStatements
var Cnn, Dlcmd string
var ticks int
var RESET_TICKS = func() { ticks = 0 }
var TICK = func() {
  if ticks++; ticks%1000 == 0 {
    fmt.Print(".")
  }
}

func ClearWorkFiles() {
  fmt.Println("Deleting work files")
  db, err := sql.Open("odbc", Cnn)
  if err != nil {
    log.Fatalf("Open Database: %v\n", err)
  }
  defer db.Close()
  ClearTable(db, "_lst1_bomitm_lcl")
  ClearTable(db, "_lst2_bomitm_lcl")
  ClearTable(db, "bomlist1")
  ClearTable(db, "bomlist2")
  ClearTable(db, "bomhdr_match")
  ClearTable(db, "bomitm_work")
}

func ClearTable(db *sql.DB, t string) {
  if _, err := db.Exec(Dlcmd + t); err != nil {
    log.Fatalf("Delete sql table: %v\n", err)
  }
}

func UfmtMatnr(m string) string {
  r := strings.NewReplacer(".", "", "/", "", "-", "")
  return r.Replace(m)
}

func Zipgen(ifnam, ofnam string) {
  outf, err := os.Create(ofnam)
  if err != nil {
    log.Fatalf("Create: %v\n", err)
  }
  w := zip.NewWriter(outf)
  f, err := w.Create(ifnam)
  if err != nil {
    log.Fatal(err)
  }
  inf, err := os.Open(ifnam)
  if err != nil {
    log.Fatalf("Open: %v\n", err)
  }
  fs, _ := inf.Stat()
  ibuf := make([]byte, fs.Size())
  _, err = inf.Read(ibuf)
  if err != nil {
    log.Fatal(err)
  }
  inf.Close()
  _, err = f.Write(ibuf)
  if err != nil {
    log.Fatal(err)
  }
  err = w.Close()
  if err != nil {
    log.Fatal(err)
  }
}

func Zipdsp(fname string) {
  rc, err := zip.OpenReader(fname)
  if err != nil {
    log.Fatal(err)
  }
  defer rc.Close()
  for _, f := range rc.File {
    d, err := f.Open()
    if err != nil {
      log.Fatal(err)
    }
    defer d.Close()
    buf := new(bytes.Buffer)
    buf.ReadFrom(d)
    for iline, err := buf.ReadString(byte('\n')); err != io.EOF; iline, err = buf.ReadString(byte('\n')) {
      flds := strings.Split(string(iline), "\t")
      fmt.Println(flds[0] + "/" + flds[1])
    }
  }
}
