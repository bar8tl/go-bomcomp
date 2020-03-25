/*** bomcomp.go : 2012-04-12 BAR8TL - Migration Tools: Simple BOM comparison Version 2.0.1 (mtrl&docm line items) ***/
package main

import rb "bar8tl/p/bomcomp"
import _  "code.google.com/p/odbc"
import "database/sql"
import "fmt"
import "log"
import "math"

func init() {
  rb.S.NewSettings("config.xml")
  rb.U.NewUom(rb.S.Uomfl)
  rb.Q.NewSqlStatements(rb.S.Sqlfl, rb.S.Dbase, rb.S.Cmode, rb.S.Rules)
}

func main() {
  if rb.S.Pstep[0] == 'L' {
    rb.LocalCopyBomItm()
  }
  if rb.S.Pstep[1] == 'M' {
    rb.PrepMatchList("auto")
  }
  if rb.S.Pstep[2] == 'C' {
    bwseBomsToCompare()
    if rb.S.Hkeep == "Y" {
      rb.ClearWorkFiles()
    }
  }
}

var h rb.SCmpHdr

func bwseBomsToCompare() {
  fmt.Print("Comparing BOMs"); rb.RESET_TICKS()
  db, err := sql.Open("odbc", rb.Cnn)
  if err != nil {
    log.Fatalf("Open Database: %v\n", err)
  }
  defer db.Close()
  rb.ClearTable(db, "bomhdr_comp")
  rb.ClearTable(db, "bomitm_work")
  var ds *sql.Rows
  ds, err = db.Query(rb.Q.Cmd["selBomMtch"])
  if err != nil {
    log.Fatalf("Execute sql selBomMtch: %v\n", err)
  }
  var m rb.SMtchN
  for ds.Next() {
    m.ScanBomMtch(ds); h.MapBomMtch(&m)
    h.Seqnh++
    if h.Idenb == "InList1_InList2" || rb.S.Fcomp == "Y" {
      compareBom(db)
      h.CalcStatsBom()
    } else {
      h.SetExclusiveBom()
    }
    h.IsrtCmpHdr(db)
  }
  fmt.Println("")
  rb.ReseqBomItems(db)
}

var d1, d2 rb.SExplN
var e1, e2 rb.SExpl
var w rb.SWrkItm

func compareBom(db *sql.DB) {
  h.Items, h.Exclu, h.Comdf, h.Tfldd = 0, 0, 0, 0
  var ds1, ds2 *sql.Rows
  if rb.S.Citms == "Mtrl" {
    ds1, _ = db.Query(rb.Q.Cmd["selMtrlItm1"], h.P1.Systm, h.P1.Werks, h.P1.Matnr, h.P1.Stlal)
    ds2, _ = db.Query(rb.Q.Cmd["selMtrlItm2"], h.P2.Systm, h.P2.Werks, h.P2.Matnr, h.P2.Stlal)
  } else if rb.S.Citms == "Docu" {
    ds1, _ = db.Query(rb.Q.Cmd["selDocuItm1"], h.P1.Systm, h.P1.Werks, h.P1.Matnr, h.P1.Stlal)
    ds2, _ = db.Query(rb.Q.Cmd["selDocuItm2"], h.P2.Systm, h.P2.Werks, h.P2.Matnr, h.P2.Stlal)
  } else if rb.S.Citms == "MtDc" {
    ds1, _ = db.Query(rb.Q.Cmd["selMtDcItm1"], h.P1.Systm, h.P1.Werks, h.P1.Matnr, h.P1.Stlal, h.P1.Systm, h.P1.Werks, h.P1.Matnr, h.P1.Stlal)
    ds2, _ = db.Query(rb.Q.Cmd["selMtDcItm2"], h.P2.Systm, h.P2.Werks, h.P2.Matnr, h.P2.Stlal, h.P2.Systm, h.P2.Werks, h.P2.Matnr, h.P2.Stlal)
  }
  readItem1(ds1)
  readItem2(ds2)
  for !e1.Eof && !e2.Eof {
    if e1.Key > e2.Key {
      notinBOM1inBOM2(db, ds2)
    } else if e1.Key < e2.Key {
      inBOM1notinBOM2(db, ds1)
    } else {
      inBOM1inBOM2(db, ds1, ds2)
    }
    rb.TICK()
  }
  for e1.Eof && !e2.Eof {
    notinBOM1inBOM2(db, ds2); rb.TICK()
  }
  for e2.Eof && !e1.Eof {
    inBOM1notinBOM2(db, ds1); rb.TICK()
  }
}

func notinBOM1inBOM2(db *sql.DB, ds2 *sql.Rows) {
  h.Items++; h.Exclu++
  w.Ident, w.Ibom1, w.Ibom2, w.Itmdf, w.Ifldd = "NotInBom1_InBom2", "", "X", "X", 0
  w.P1.Systm = h.P1.Systm; w.P1.Werks = h.P1.Werks; w.P1.Matnr = h.P1.Matnr; w.P1.Stlal = h.P1.Stlal
  w.HideItemBom1().ClearMarksItemBom1().ShowItemBom2(e2).ClearMarksItemBom2().IsrtWrkItm(db, h); readItem2(ds2)
}

func inBOM1notinBOM2(db *sql.DB, ds1 *sql.Rows) {
  h.Items++; h.Exclu++
  w.Ident, w.Ibom1, w.Ibom2, w.Itmdf, w.Ifldd = "InBom1_NotInBom2", "X", "", "X", 0
  w.P2.Systm = h.P2.Systm; w.P2.Werks = h.P2.Werks; w.P2.Matnr = h.P2.Matnr; w.P2.Stlal = h.P2.Stlal
  w.ShowItemBom1(e1).ClearMarksItemBom1().HideItemBom2().ClearMarksItemBom2().IsrtWrkItm(db, h); readItem1(ds1)
}

func inBOM1inBOM2(db *sql.DB, ds1, ds2 *sql.Rows) {
  h.Items++
  w.Ident, w.Ibom1, w.Ibom2, w.Itmdf, w.Ifldd = "InBom1_InBom2", "X", "X", "", 0
  w.ShowItemBom1(e1); w.ShowItemBom2(e2); w.ClearMarksItemBom1(); w.ClearMarksItemBom2()
  if (rb.S.Level == "lv1" || rb.S.Level == "") &&
     (rb.UfmtMatnr(e1.Idnrk) != rb.UfmtMatnr(e2.Idnrk) || e1.Menge != e2.Menge || e1.Meins != e2.Meins || e1.Dokar != e2.Dokar ||
      rb.UfmtMatnr(e1.Doknr) != rb.UfmtMatnr(e2.Doknr) || e1.Dokvr != e2.Dokvr || e1.Doktl != e2.Doktl) {
    checkMatchLev1()
  } else {
  if (rb.S.Level == "lv2" || rb.S.Level == "") &&
     (rb.UfmtMatnr(e1.Idnrk) != rb.UfmtMatnr(e2.Idnrk) || e1.Menge != e2.Menge || e1.Meins != e2.Meins || e1.Dokar != e2.Dokar ||
      rb.UfmtMatnr(e1.Doknr) != rb.UfmtMatnr(e2.Doknr) || e1.Dokvr != e2.Dokvr || e1.Doktl != e2.Doktl || e1.Postp != e2.Postp) {
    checkMatchLev1(); checkMatchLev2()
  } else {
  if (rb.S.Level == "lv3") &&
     (rb.UfmtMatnr(e1.Idnrk) != rb.UfmtMatnr(e2.Idnrk) || e1.Menge != e2.Menge || e1.Meins != e2.Meins || e1.Dokar != e2.Dokar ||
      rb.UfmtMatnr(e1.Doknr) != rb.UfmtMatnr(e2.Doknr) || e1.Dokvr != e2.Dokvr || e1.Doktl != e2.Doktl || e1.Postp != e2.Postp ||
      e1.Sanko != e2.Sanko || e1.Sanfe != e2.Sanfe || e1.Sanka != e2.Sanka) {
    checkMatchLev1(); checkMatchLev2(); checkMatchLev3()
  } else {
  if (rb.S.Level == "lv4") &&
     (rb.UfmtMatnr(e1.Idnrk) != rb.UfmtMatnr(e2.Idnrk) || e1.Menge != e2.Menge || e1.Meins != e2.Meins || e1.Dokar != e2.Dokar ||
      rb.UfmtMatnr(e1.Doknr) != rb.UfmtMatnr(e2.Doknr) || e1.Dokvr != e2.Dokvr || e1.Doktl != e2.Doktl || e1.Postp != e2.Postp ||
      e1.Alpos != e2.Alpos || e1.Alpgr != e2.Alpgr || e1.Alprf != e2.Alprf || e1.Alpst != e2.Alpst || e1.Ewahr != e2.Ewahr ||
      e1.Sanko != e2.Sanko || e1.Sanfe != e2.Sanfe || e1.Sanka != e2.Sanka) {
    checkMatchLev1(); checkMatchLev2(); checkMatchLev3(); checkMatchLev4() }}}
  }
  if w.Ifldd > 0 {
    w.Itmdf = "X"
    h.Comdf++
  }
  w.IsrtWrkItm(db, h); readItem1(ds1); readItem2(ds2)
}

func checkMatchLev1() {
  var factor, conver, difere, tolera float64
  if rb.UfmtMatnr(e1.Idnrk) != rb.UfmtMatnr(e2.Idnrk) { w.Ifldd++; h.Tfldd++; w.P1.Idnrkd, w.P2.Idnrkd = "*", "*" }
  if e1.Meins != e2.Meins {
    if factor = float64(rb.U.InUom(e1.Meins, e2.Meins)); factor > 0 {
    } else                { w.Ifldd++; h.Tfldd++; w.P1.Meinsd, w.P2.Meinsd = "*", "*" }}
  if factor > 0 && factor != 1 {
    conver = e1.Menge * factor; difere = e2.Menge - conver; tolera = float64(rb.S.Convt) * e2.Menge }
  if ((e1.Meins == e2.Meins || factor == 1) && e1.Menge == e2.Menge) ||
     ((e1.Meins != e2.Meins && factor != 1) && factor > 0 &&
     ((tolera > 0 && float64(math.Abs(float64(difere))) <= tolera) || (tolera == 0 && difere == 0))) {
  } else                  { w.Ifldd++; h.Tfldd++; w.P1.Menged, w.P2.Menged = "*", "*" }
  if e1.Dokar != e2.Dokar { w.Ifldd++; h.Tfldd++; w.P1.Dokard, w.P2.Dokard = "*", "*" }
  if rb.UfmtMatnr(e1.Doknr) != rb.UfmtMatnr(e2.Doknr) { w.Ifldd++; h.Tfldd++; w.P1.Doknrd, w.P2.Doknrd = "*", "*" }
  if e1.Dokvr != e2.Dokvr { w.Ifldd++; h.Tfldd++; w.P1.Dokvrd, w.P2.Dokvrd = "*", "*" }
  if e1.Doktl != e2.Doktl { w.Ifldd++; h.Tfldd++; w.P1.Doktld, w.P2.Doktld = "*", "*" }
}

func checkMatchLev2() {
  if e1.Postp != e2.Postp { w.Ifldd++; h.Tfldd++; w.P1.Postpd, w.P2.Postpd = "*", "*" }
}

func checkMatchLev3() {
  if e1.Sanko != e2.Sanko { w.Ifldd++; h.Tfldd++; w.P1.Sankod, w.P2.Sankod = "*", "*" }
  if e1.Sanfe != e2.Sanfe { w.Ifldd++; h.Tfldd++; w.P1.Sanfed, w.P2.Sanfed = "*", "*" }
  if e1.Sanka != e2.Sanka { w.Ifldd++; h.Tfldd++; w.P1.Sankad, w.P2.Sankad = "*", "*" }
}

func checkMatchLev4() {
  if e1.Alpos != e2.Alpos { w.Ifldd++; h.Tfldd++; w.P1.Alposd, w.P2.Alposd = "*", "*" }
  if e1.Alpgr != e2.Alpgr { w.Ifldd++; h.Tfldd++; w.P1.Alpgrd, w.P2.Alpgrd = "*", "*" }
  if e1.Alprf != e2.Alprf { w.Ifldd++; h.Tfldd++; w.P1.Alprfd, w.P2.Alprfd = "*", "*" }
  if e1.Alpst != e2.Alpst { w.Ifldd++; h.Tfldd++; w.P1.Alpstd, w.P2.Alpstd = "*", "*" }
  if e1.Ewahr != e2.Ewahr { w.Ifldd++; h.Tfldd++; w.P1.Ewahrd, w.P2.Ewahrd = "*", "*" }
}

func readItem1(ds1 *sql.Rows) {
  if ds1.Next() {
    d1.ScanExpl(ds1); e1.MapItem(d1)
    e1.Eof = false
  } else {
    e1.Eof = true
  }
}

func readItem2(ds2 *sql.Rows) {
  if ds2.Next() {
    d2.ScanExpl(ds2); e2.MapItem(d2)
    e2.Eof = false
  } else {
    e2.Eof = true
  }
}
