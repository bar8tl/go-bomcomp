package bomcomp

import "archive/zip"
import "bytes"
import "io"
import "log"
import "strings"

type SqlStatements struct {
  Cmd map[string]string
}

func (q *SqlStatements) NewSqlStatements(fname, dbase, cmode, rules string) {
  q.Cmd = make(map[string]string)
  rc, err := zip.OpenReader(fname)
  if err != nil {
    log.Fatalf("Open Sql Statements Archive file: %v\n", err)
  }
  defer rc.Close()
  for _, f := range rc.File {
    var d io.ReadCloser
    d, err = f.Open()
    if err != nil {
      log.Fatalf("Open SQL Statements archived file: %v\n", err)
    }
    defer d.Close()
    buf := new(bytes.Buffer)
    buf.ReadFrom(d)
    for line, err := buf.ReadString(byte('\n')); err != io.EOF; line, err = buf.ReadString(byte('\n')) {
      flds := strings.Split(string(line), "\t")
      if cmode == "lcl" && strings.Contains(flds[1], "_xls") {
        flds[1] = strings.Replace(flds[1], "_xls", "_lcl", -1)
      }
      q.Cmd[flds[0]] = flds[1]
    }
  }
  if rules == "TlP" {
    q.Cmd["selBomItem1"] = q.Cmd["selBomItem1tr"]
    q.Cmd["selBomItem2"] = q.Cmd["selBomItem2tr"]
    q.Cmd["selMtrlItm1"] = q.Cmd["selMtrlItm1tr"]
    q.Cmd["selMtrlItm2"] = q.Cmd["selMtrlItm2tr"]
    q.Cmd["selMtDcItm1"] = q.Cmd["selMtDcItm1tr"]
    q.Cmd["selMtDcItm2"] = q.Cmd["selMtDcItm2tr"]
  } else {
    q.Cmd["selBomItem1"] = q.Cmd["selBomItem1nr"]
    q.Cmd["selBomItem2"] = q.Cmd["selBomItem2nr"]
    q.Cmd["selMtrlItm1"] = q.Cmd["selMtrlItm1nr"]
    q.Cmd["selMtrlItm2"] = q.Cmd["selMtrlItm2nr"]
    q.Cmd["selMtDcItm1"] = q.Cmd["selMtDcItm1nr"]
    q.Cmd["selMtDcItm2"] = q.Cmd["selMtDcItm2nr"]
  }
  Cnn = q.Cmd["sCnnODBC"] + dbase + ";"
  Dlcmd = q.Cmd["delTabl"]
}
