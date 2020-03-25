package bomcomp

import _ "code.google.com/p/odbc"
import "database/sql"
import "fmt"
import "log"
import "strconv"

type sExplN struct{
  bomid, systm, werks, matnr, maktx, stlal, mtart, beskz, sobsl, wtbom sql.NullString
}
type sList  struct{
  bomid, systm, werks, matnr, maktx, stlal, mtart, beskz, sobsl, wtbom string
  matnb string
}

func PrepMatchList(src string) {
  fmt.Print("Determining BOMs for comparison"); RESET_TICKS()
  if src == "auto" {
    prepList("1")
    prepList("2")
  }
  db, err := sql.Open("odbc", Cnn)
  if err != nil {
    log.Fatalf("Open Database: %v\n", err)
  }
  defer db.Close()
  ClearTable(db, "bomhdr_match")
  if src == "auto" {
    _, err = db.Exec(Q.Cmd["isrBomMtchi"])
  } else {
    _, err = db.Exec(Q.Cmd["isrBomMtche"])
  }
  if err != nil {
    log.Fatalf("Execute sql isrBomMtch_: %v\n", err)
  }
  fmt.Println("")
}

func prepList(i string) {
  db, err := sql.Open("odbc", Cnn)
  if err != nil {
    log.Fatalf("Open Database: %v\n", err)
  }
  defer db.Close()
  ClearTable(db, "bomlist"+i)
  ds, err := db.Query(Q.Cmd["selBomItem"+i])
  if err != nil {
    log.Fatalf("Execute sql selBomItem_: %v\n", err)
  }
  for ds.Next() {
    e := new(sExplN)
    l := new(sList)
    err = ds.Scan(&e.bomid, &e.systm, &e.werks, &e.matnr, &e.maktx, &e.stlal, &e.mtart, &e.beskz, &e.sobsl, &e.wtbom)
    if err != nil {
      log.Fatalf("Scan selBomItem_ row: %v\n", err)
    }
    mapBomItm(e, l)
    l.matnb = UfmtMatnr(l.matnr)
    l.bomid = l.systm + l.werks + l.matnb + l.stlal
    _, err = db.Exec(Q.Cmd["isrBomList"+i], l.bomid, l.systm, l.werks, l.matnr, l.maktx, l.stlal, l.mtart, l.beskz, l.sobsl, l.wtbom, l.matnb)
    if err != nil {
      log.Fatalf("Execute sql isrBomList: %v\n", err)
    }
    TICK()
  }
}

func mapBomItm(e *sExplN, l *sList) {
  if e.bomid.Valid { l.bomid = e.bomid.String }
  if e.systm.Valid { l.systm = e.systm.String }
  if e.werks.Valid { l.werks = e.werks.String }
  if e.matnr.Valid { l.matnr = e.matnr.String }
  if e.maktx.Valid { l.maktx = e.maktx.String }
  if e.stlal.Valid { l.stlal = e.stlal.String }
  if e.mtart.Valid { l.mtart = e.mtart.String }
  if e.beskz.Valid { l.beskz = e.beskz.String }
  if e.sobsl.Valid { l.sobsl = e.sobsl.String }
  if e.wtbom.Valid { l.wtbom = e.wtbom.String }
}

// Match list reader
type sBomidN struct {
  Systm, Werks, Matnr, Maktx, Stlal, Mtart, Beskz, Sobsl, Wtbom sql.NullString
}
type SMtchN struct {
  Matnb, Stlab, Idenb sql.NullString
  P1, P2 sBomidN
}
func (m *SMtchN) ScanBomMtch(ds *sql.Rows) {
  err := ds.Scan(&m.Matnb, &m.Stlab, &m.Idenb,
    &m.P1.Systm, &m.P1.Werks, &m.P1.Matnr, &m.P1.Maktx, &m.P1.Stlal, &m.P1.Mtart, &m.P1.Beskz, &m.P1.Sobsl, &m.P1.Wtbom,
    &m.P2.Systm, &m.P2.Werks, &m.P2.Matnr, &m.P2.Maktx, &m.P2.Stlal, &m.P2.Mtart, &m.P2.Beskz, &m.P2.Sobsl, &m.P2.Wtbom)
  if err != nil {
    log.Fatalf("Scan selBomMatch row: %v\n", err)
  }
}

// Header comparison writer
type sBomid struct {
  Systm, Werks, Matnr, Maktx, Stlal, Mtart, Beskz, Sobsl, Wtbom string
}
type SCmpHdr struct {
  Seqnh, Items, Commo, Exclu, Commp, Comeq, Comdf, Comep, Tfldd, Mrate int
  Matnb, Stlab, Idenb, Gbcod, Bucod, Bomdf string
  P1, P2 sBomid
}
func (h *SCmpHdr) Clear() {
  h.Items, h.Commo, h.Exclu, h.Commp, h.Comeq, h.Comdf, h.Comep, h.Tfldd, h.Mrate = 0, 0, 0, 0, 0, 0, 0, 0, 0
  h.Matnb, h.Stlab, h.Idenb, h.Gbcod, h.Bucod, h.Bomdf = "", "", "", "", "", ""
  h.P1.Systm, h.P1.Werks, h.P1.Matnr, h.P1.Maktx, h.P1.Stlal, h.P1.Mtart, h.P1.Beskz, h.P1.Sobsl = "", "", "", "", "", "", "", ""
  h.P1.Wtbom = ""
  h.P2.Systm, h.P2.Werks, h.P2.Matnr, h.P2.Maktx, h.P2.Stlal, h.P2.Mtart, h.P2.Beskz, h.P2.Sobsl = "", "", "", "", "", "", "", ""
  h.P2.Wtbom = ""
}
func (h *SCmpHdr) MapBomMtch(m *SMtchN) {
  h.Clear()
  if m.Matnb.Valid    { h.Matnb    = m.Matnb.String    }
  if m.Stlab.Valid    { h.Stlab    = m.Stlab.String    }
  if m.Idenb.Valid    { h.Idenb    = m.Idenb.String    }
  if m.P1.Systm.Valid { h.P1.Systm = m.P1.Systm.String }
  if m.P1.Werks.Valid { h.P1.Werks = m.P1.Werks.String }
  if m.P1.Matnr.Valid { h.P1.Matnr = m.P1.Matnr.String }
  if m.P1.Maktx.Valid { h.P1.Maktx = m.P1.Maktx.String }
  if m.P1.Stlal.Valid { h.P1.Stlal = m.P1.Stlal.String }
  if m.P1.Mtart.Valid { h.P1.Mtart = m.P1.Mtart.String }
  if m.P1.Beskz.Valid { h.P1.Beskz = m.P1.Beskz.String }
  if m.P1.Sobsl.Valid { h.P1.Sobsl = m.P1.Sobsl.String }
  if m.P1.Wtbom.Valid { h.P1.Wtbom = m.P1.Wtbom.String }
  if m.P2.Systm.Valid { h.P2.Systm = m.P2.Systm.String }
  if m.P2.Werks.Valid { h.P2.Werks = m.P2.Werks.String }
  if m.P2.Matnr.Valid { h.P2.Matnr = m.P2.Matnr.String }
  if m.P2.Maktx.Valid { h.P2.Maktx = m.P2.Maktx.String }
  if m.P2.Stlal.Valid { h.P2.Stlal = m.P2.Stlal.String }
  if m.P2.Mtart.Valid { h.P2.Mtart = m.P2.Mtart.String }
  if m.P2.Beskz.Valid { h.P2.Beskz = m.P2.Beskz.String }
  if m.P2.Sobsl.Valid { h.P2.Sobsl = m.P2.Sobsl.String }
  if m.P2.Wtbom.Valid { h.P2.Wtbom = m.P2.Wtbom.String }
}
func (h *SCmpHdr) CalcStatsBom() {
  h.Bomdf, h.Commp, h.Comep = "", 0, 0
  if h.Exclu != 0 || h.Tfldd != 0 {
    h.Bomdf = "X"
  }
  h.Commo = h.Items - h.Exclu
  if h.Items != 0 {
    h.Commp = int(float32(h.Commo)/float32(h.Items)*100 + 0.51)
  }
  h.Comeq = h.Commo - h.Comdf
  if h.Commo != 0 {
    h.Comep = int(float32(h.Comeq)/float32(h.Commo)*100 + 0.51)
  }
  if h.Commp < h.Comep {
    h.Mrate = h.Commp
  } else {
    h.Mrate = h.Comep
  }
}
func (h *SCmpHdr) SetExclusiveBom() {
  h.Bomdf, h.Items, h.Commo, h.Exclu, h.Commp, h.Comeq, h.Comdf, h.Comep, h.Tfldd, h.Mrate = "X", 0, 0, 0, 0, 0, 0, 0, 0, 0
}
func (h *SCmpHdr) IsrtCmpHdr(db *sql.DB) {
  _, err := db.Exec(Q.Cmd["isrBomHdrComp"], strconv.Itoa(h.Seqnh), S.Gbcod, S.Bucod, h.Idenb,
    h.P1.Systm, h.P1.Matnr, h.P1.Maktx, h.P1.Werks, h.P1.Stlal, h.P1.Mtart, h.P1.Beskz, h.P1.Sobsl, h.P1.Wtbom,
    h.P2.Systm, h.P2.Matnr, h.P2.Maktx, h.P2.Werks, h.P2.Stlal, h.P2.Mtart, h.P2.Beskz, h.P2.Sobsl, h.P2.Wtbom,
    h.Bomdf, strconv.Itoa(h.Items), strconv.Itoa(h.Commo), strconv.Itoa(h.Exclu), strconv.Itoa(h.Commp), strconv.Itoa(h.Comeq),
    strconv.Itoa(h.Comdf), strconv.Itoa(h.Comep), strconv.Itoa(h.Tfldd), strconv.Itoa(h.Mrate))
  if err != nil {
    log.Fatalf("Execute sql isrBomHdrComp: %v\n", err)
  }
}
