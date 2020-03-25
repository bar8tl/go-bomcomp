package bomcomp

import _ "code.google.com/p/odbc"
import "database/sql"
import "fmt"
import "log"
import "strings"
import "strconv"

func LocalCopyBomItm() {
  fmt.Print("Start local copy of BOM files to compare"); RESET_TICKS()
  copyBomItm("1")
  copyBomItm("2")
  fmt.Println("")
}
func copyBomItm(bitm string) {
  isrLclBomItm := strings.Replace(Q.Cmd["isrLclBomItm"], "#", bitm, -1)
  selLclBomItm := strings.Replace(Q.Cmd["selLclBomItm"], "#", bitm, -1)
  updLclBomItm := strings.Replace(Q.Cmd["updLclBomItm"], "#", bitm, -1)
  db, err := sql.Open("odbc", Cnn)
  if err != nil {
    log.Fatalf("Open Database: %v\n", err)
  }
  defer db.Close()
  ClearTable(db, "_lst"+bitm+"_bomitm_lcl")
  _, err = db.Exec(isrLclBomItm)
  if err != nil {
    log.Fatalf("Execute sql isrLclBomItm: %v\n", err)
  }
  ds, err := db.Query(selLclBomItm)
  if err != nil {
    log.Fatalf("Execute sql selLclBomItm: %v\n", err)
  }
  for ds.Next() {
    var seqno float32
    var matnr, match string
    var nmtch sql.NullString
    ds.Scan(&seqno, &matnr, &nmtch)
    var seqnr int = int(seqno)
    if nmtch.Valid {
      match = nmtch.String
    }
    matnr = UfmtMatnr(matnr)
    match = UfmtMatnr(match)
    db.Exec(updLclBomItm, match, strconv.Itoa(seqnr))
    TICK()
  }
}

func ReseqBomItems(db *sql.DB) {
  fmt.Print("Preparing comparison report"); RESET_TICKS()
  ClearTable(db,"bomitm_comp")
  ds, err := db.Query(Q.Cmd["selBomItmWork"])
  if err != nil {
    log.Fatalf("Execute sql selBomItmWork: %v\n", err)
  }
  var j SWrkItmN
  var i SCmpItm
  var bomno, wbomn int
  for ds.Next() {
    j.ScanWrkItm(ds); i.MapItmCmp(j)
    bomno = i.Seqnh
    if bomno == wbomn {
      i.Nline++;
    } else {
      wbomn = bomno
      i.Nline = 1
    }
    i.Seqni++
    i.IsrtCmpItm(db)
    TICK()
  }
  fmt.Println("")
}

// Explosion file reader
type SExplN struct {
  Systm, Werks, Matnr, Maktx, Stlal, Mtart, Beskz, Sobsl sql.NullString
  Posnr, Postp, Idnrk, Cmakt, Cmtar, Cbesk, Csobs, Meins, Alpos, Alpgr, Alprf, Alpst, Sanko, Sanfe, Sanka, Dokar, Doknr, Dokvr,
    Doktl sql.NullString
  Menge, Ewahr sql.NullFloat64
  Match, Exstr sql.NullString
}
func (d *SExplN) ScanExpl(ds *sql.Rows) {
  err := ds.Scan(&d.Systm, &d.Werks, &d.Matnr, &d.Maktx, &d.Stlal, &d.Mtart, &d.Beskz, &d.Sobsl, &d.Posnr, &d.Postp, &d.Idnrk,
    &d.Cmakt, &d.Cmtar, &d.Cbesk, &d.Csobs, &d.Menge, &d.Meins, &d.Alpos, &d.Alpgr, &d.Alprf, &d.Alpst, &d.Ewahr, &d.Sanko,
    &d.Sanfe, &d.Sanka, &d.Dokar, &d.Doknr, &d.Dokvr, &d.Doktl, &d.Match, &d.Exstr)
  if err != nil {
    log.Fatalf("Scan sql ScanExpl: %v\n", err)
  }
}

//
type SExpl struct {
  Systm, Werks, Matnr, Maktx, Stlal, Mtart, Beskz, Sobsl string
  Posnr, Postp, Idnrk, Cmakt, Cmtar, Cbesk, Csobs, Meins, Alpos, Alpgr, Alprf, Alpst, Sanko, Sanfe, Sanka, Dokar, Doknr, Dokvr,
    Doktl string
  Menge, Ewahr float64
  Match, Exstr, Key string
  Eof bool
}
func (e *SExpl) Clear() {
  e.Systm, e.Werks, e.Matnr, e.Maktx, e.Stlal, e.Mtart, e.Beskz, e.Sobsl = "", "", "", "", "", "", "", ""
  e.Posnr, e.Postp, e.Idnrk, e.Cmakt, e.Cmtar, e.Cbesk, e.Csobs, e.Meins, e.Alpos, e.Alpgr = "", "", "", "", "", "", "", "", "", ""
  e.Alprf, e.Alpst, e.Sanko, e.Sanfe, e.Sanka, e.Dokar, e.Doknr, e.Dokvr, e.Doktl = "", "", "", "", "", "", "", "", ""
  e.Menge, e.Ewahr = 0, 0
}
func (e *SExpl) MapItem (d SExplN) {
  e.Clear()
  if d.Systm.Valid { e.Systm = d.Systm.String  }
  if d.Werks.Valid { e.Werks = d.Werks.String  }
  if d.Matnr.Valid { e.Matnr = d.Matnr.String  }
  if d.Maktx.Valid { e.Maktx = d.Maktx.String  }
  if d.Stlal.Valid { e.Stlal = d.Stlal.String  }
  if d.Mtart.Valid { e.Mtart = d.Mtart.String  }
  if d.Beskz.Valid { e.Beskz = d.Beskz.String  }
  if d.Sobsl.Valid { e.Sobsl = d.Sobsl.String  }
  if d.Posnr.Valid { e.Posnr = d.Posnr.String  }
  if d.Postp.Valid { e.Postp = d.Postp.String  }
  if d.Idnrk.Valid { e.Idnrk = d.Idnrk.String  }
  if d.Cmakt.Valid { e.Cmakt = d.Cmakt.String  }
  if d.Cmtar.Valid { e.Cmtar = d.Cmtar.String  }
  if d.Cbesk.Valid { e.Cbesk = d.Cbesk.String  }
  if d.Csobs.Valid { e.Csobs = d.Csobs.String  }
  if d.Menge.Valid { e.Menge = d.Menge.Float64 }
  if d.Meins.Valid { e.Meins = d.Meins.String  }
  if d.Alpos.Valid { e.Alpos = d.Alpos.String  }
  if d.Alpgr.Valid { e.Alpgr = d.Alpgr.String  }
  if d.Alprf.Valid { e.Alprf = d.Alprf.String  }
  if d.Alpst.Valid { e.Alpst = d.Alpst.String  }
  if d.Ewahr.Valid { e.Ewahr = d.Ewahr.Float64 }
  if d.Sanko.Valid { e.Sanko = d.Sanko.String  }
  if d.Sanfe.Valid { e.Sanfe = d.Sanfe.String  }
  if d.Sanka.Valid { e.Sanka = d.Sanka.String  }
  if d.Dokar.Valid { e.Dokar = d.Dokar.String  }
  if d.Doknr.Valid { e.Doknr = d.Doknr.String  }
  if d.Dokvr.Valid { e.Dokvr = d.Dokvr.String  }
  if d.Doktl.Valid { e.Doktl = d.Doktl.String  }
  if d.Match.Valid { e.Match = d.Match.String  }
  if d.Exstr.Valid { e.Exstr = d.Exstr.String  }
  e.Key = UfmtMatnr(e.Idnrk)
}

// Work comparison file writer
type SItmdt struct {
  Systm, Werks, Matnr, Maktx, Stlal, Mtart, Beskz, Sobsl string
  Posnr, Postp, Idnrk, Cmakt, Cmtar, Cbesk, Csobs, Meins, Alpos, Alpgr, Alprf, Alpst, Sanko, Sanfe, Sanka, Dokar, Doknr, Dokvr,
    Doktl  string
  Menge, Ewahr float64
  Match, Exstr  string
  Posnrd,Postpd,Idnrkd,Cmaktd, Cmtard,Cbeskd,Csobsd,Meinsd,Alposd,Alpgrd,Alprfd,Alpstd,Sankod,Sanfed,Sankad,Dokard,Doknrd,Dokvrd,
    Doktld string
  Menged,Ewahrd string
}
type SWrkItm struct {
  Seqnw, Seqni, Ifldd, Nline int
  Ident, Ibom1, Ibom2, Itmdf, Posnr string
  P1, P2 SItmdt
}
func (w *SWrkItm) HideItemBom1() *SWrkItm {
  w.P1.Maktx, w.P1.Mtart, w.P1.Beskz, w.P1.Sobsl, w.P1.Posnr, w.P1.Postp, w.P1.Idnrk, w.P1.Cmakt = "", "", "", "", "", "", "", ""
  w.P1.Cmtar, w.P1.Cbesk, w.P1.Csobs, w.P1.Meins, w.P1.Alpos, w.P1.Alpgr, w.P1.Alprf, w.P1.Alpst = "", "", "", "", "", "", "", ""
  w.P1.Sanko, w.P1.Sanfe, w.P1.Sanka, w.P1.Dokar, w.P1.Doknr, w.P1.Dokvr, w.P1.Doktl             = "", "", "", "", "", "", ""
  w.P1.Menge, w.P1.Ewahr = 0, 0
  return w
}
func (w *SWrkItm) HideItemBom2() *SWrkItm {
  w.P2.Maktx, w.P2.Mtart, w.P2.Beskz, w.P2.Sobsl, w.P2.Posnr, w.P2.Postp, w.P2.Idnrk, w.P2.Cmakt = "", "", "", "", "", "", "", ""
  w.P2.Cmtar, w.P2.Cbesk, w.P2.Csobs, w.P2.Meins, w.P2.Alpos, w.P2.Alpgr, w.P2.Alprf, w.P2.Alpst = "", "", "", "", "", "", "", ""
  w.P2.Sanko, w.P2.Sanfe, w.P2.Sanka, w.P2.Dokar, w.P2.Doknr, w.P2.Dokvr, w.P2.Doktl             = "", "", "", "", "", "", ""
  w.P2.Menge, w.P2.Ewahr = 0, 0
  return w
}
func (w *SWrkItm) ShowItemBom1(e1 SExpl) *SWrkItm {
  w.Posnr    = e1.Posnr; w.P1.Systm = e1.Systm; w.P1.Werks = e1.Werks; w.P1.Matnr = e1.Matnr; w.P1.Maktx = e1.Maktx
  w.P1.Stlal = e1.Stlal; w.P1.Mtart = e1.Mtart; w.P1.Beskz = e1.Beskz; w.P1.Sobsl = e1.Sobsl; w.P1.Posnr = e1.Posnr
  w.P1.Postp = e1.Postp; w.P1.Idnrk = e1.Idnrk; w.P1.Cmakt = e1.Cmakt; w.P1.Cmtar = e1.Cmtar; w.P1.Cbesk = e1.Cbesk
  w.P1.Csobs = e1.Csobs; w.P1.Menge = e1.Menge; w.P1.Meins = e1.Meins; w.P1.Dokar = e1.Dokar; w.P1.Doknr = e1.Doknr
  w.P1.Dokvr = e1.Dokvr; w.P1.Doktl = e1.Doktl; w.P1.Alpos = e1.Alpos; w.P1.Alpgr = e1.Alpgr; w.P1.Alprf = e1.Alprf
  w.P1.Alpst = e1.Alpst; w.P1.Ewahr = e1.Ewahr; w.P1.Sanko = e1.Sanko; w.P1.Sanfe = e1.Sanfe; w.P1.Sanka = e1.Sanka
  return w
}
func (w *SWrkItm) ShowItemBom2(e2 SExpl) *SWrkItm {
  w.Posnr    = e2.Posnr; w.P2.Systm = e2.Systm; w.P2.Werks = e2.Werks; w.P2.Matnr = e2.Matnr; w.P2.Maktx = e2.Maktx
  w.P2.Stlal = e2.Stlal; w.P2.Mtart = e2.Mtart; w.P2.Beskz = e2.Beskz; w.P2.Sobsl = e2.Sobsl; w.P2.Posnr = e2.Posnr
  w.P2.Postp = e2.Postp; w.P2.Idnrk = e2.Idnrk; w.P2.Cmakt = e2.Cmakt; w.P2.Cmtar = e2.Cmtar; w.P2.Cbesk = e2.Cbesk
  w.P2.Csobs = e2.Csobs; w.P2.Menge = e2.Menge; w.P2.Meins = e2.Meins; w.P2.Dokar = e2.Dokar; w.P2.Doknr = e2.Doknr
  w.P2.Dokvr = e2.Dokvr; w.P2.Doktl = e2.Doktl; w.P2.Alpos = e2.Alpos; w.P2.Alpgr = e2.Alpgr; w.P2.Alprf = e2.Alprf
  w.P2.Alpst = e2.Alpst; w.P2.Ewahr = e2.Ewahr; w.P2.Sanko = e2.Sanko; w.P2.Sanfe = e2.Sanfe; w.P2.Sanka = e2.Sanka
  return w
}
func (w *SWrkItm) ClearMarksItemBom1() *SWrkItm {
  w.P1.Posnrd, w.P1.Postpd, w.P1.Idnrkd, w.P1.Cmaktd, w.P1.Cmtard, w.P1.Cbeskd, w.P1.Csobsd = "", "", "", "", "", "", ""
  w.P1.Menged, w.P1.Meinsd, w.P1.Alposd, w.P1.Alpgrd, w.P1.Alprfd, w.P1.Alpstd, w.P1.Ewahrd = "", "", "", "", "", "", ""
  w.P1.Sankod, w.P1.Sanfed, w.P1.Sankad, w.P1.Dokard, w.P1.Doknrd, w.P1.Dokvrd, w.P1.Doktld = "", "", "", "", "", "", ""
  return w
}
func (w *SWrkItm) ClearMarksItemBom2() *SWrkItm {
  w.P2.Posnrd, w.P2.Postpd, w.P2.Idnrkd, w.P2.Cmaktd, w.P2.Cmtard, w.P2.Cbeskd, w.P2.Csobsd = "", "", "", "", "", "", ""
  w.P2.Menged, w.P2.Meinsd, w.P2.Alposd, w.P2.Alpgrd, w.P2.Alprfd, w.P2.Alpstd, w.P2.Ewahrd = "", "", "", "", "", "", ""
  w.P2.Sankod, w.P2.Sanfed, w.P2.Sankad, w.P2.Dokard, w.P2.Doknrd, w.P2.Dokvrd, w.P2.Doktld = "", "", "", "", "", "", ""
  return w
}
func (w *SWrkItm) IsrtWrkItm(db *sql.DB, h SCmpHdr) {
  w.Seqnw++
  sSeqnw := strconv.Itoa(w.Seqnw)
  sIfldd := strconv.Itoa(w.Ifldd)
  sItems := strconv.Itoa(h.Items)
  sSeqnh := strconv.Itoa(h.Seqnh)
  smeng1 := strconv.FormatFloat(float64(w.P1.Menge), 'f', -1, 32)
  sewah1 := strconv.FormatFloat(float64(w.P1.Ewahr), 'f', -1, 32)
  smeng2 := strconv.FormatFloat(float64(w.P2.Menge), 'f', -1, 32)
  sewah2 := strconv.FormatFloat(float64(w.P2.Ewahr), 'f', -1, 32)
  _, err := db.Exec(Q.Cmd["isrBomItmWork"], sSeqnw, S.Gbcod, S.Bucod, w.Ident,   w.Ibom1, w.Ibom2, sIfldd, w.Itmdf, sItems,
    w.P1.Systm,  w.P1.Matnr, w.P1.Maktx,  w.P1.Werks, w.P1.Stlal,  w.P1.Mtart, w.P1.Beskz,  w.P1.Sobsl, w.P1.Posnrd, w.P1.Posnr,
    w.P1.Postpd, w.P1.Postp, w.P1.Idnrkd, w.P1.Idnrk, w.P1.Cmaktd, w.P1.Cmakt, w.P1.Cmtard, w.P1.Cmtar, w.P1.Cbeskd, w.P1.Cbesk,
    w.P1.Csobsd, w.P1.Csobs, w.P1.Menged, smeng1,     w.P1.Meinsd, w.P1.Meins, w.P1.Alposd, w.P1.Alpos, w.P1.Alpgrd, w.P1.Alpgr,
    w.P1.Alprfd, w.P1.Alprf, w.P1.Alpstd, w.P1.Alpst, w.P1.Ewahrd, sewah1,     w.P1.Sankod, w.P1.Sanko, w.P1.Sanfed, w.P1.Sanfe,
    w.P1.Sankad, w.P1.Sanka, w.P1.Dokard, w.P1.Dokar, w.P1.Doknrd, w.P1.Doknr, w.P1.Dokvrd, w.P1.Dokvr, w.P1.Doktld, w.P1.Doktl,
    w.P2.Systm,  w.P2.Matnr, w.P2.Maktx,  w.P2.Werks, w.P2.Stlal,  w.P2.Mtart, w.P2.Beskz,  w.P2.Sobsl, w.P2.Posnrd, w.P2.Posnr,
    w.P2.Postpd, w.P2.Postp, w.P2.Idnrkd, w.P2.Idnrk, w.P2.Cmaktd, w.P2.Cmakt, w.P2.Cmtard, w.P2.Cmtar, w.P2.Cbeskd, w.P2.Cbesk,
    w.P2.Csobsd, w.P2.Csobs, w.P2.Menged, smeng2,     w.P2.Meinsd, w.P2.Meins, w.P2.Alposd, w.P2.Alpos, w.P2.Alpgrd, w.P2.Alpgr,
    w.P2.Alprfd, w.P2.Alprf, w.P2.Alpstd, w.P2.Alpst, w.P2.Ewahrd, sewah2,     w.P2.Sankod, w.P2.Sanko, w.P2.Sanfed, w.P2.Sanfe,
    w.P2.Sankad, w.P2.Sanka, w.P2.Dokard, w.P2.Dokar, w.P2.Doknrd, w.P2.Doknr, w.P2.Dokvrd, w.P2.Dokvr, w.P2.Doktld, w.P2.Doktl,
    sSeqnh, w.Posnr)
  if err != nil {
    log.Fatalf("Execute sql isrBomItmWork: %v\n", err)
  }
}

// Work comparison file reader
type SItmdtN struct {
  Systm, Werks, Matnr, Maktx, Stlal, Mtart, Beskz, Sobsl sql.NullString
  Posnr, Postp, Idnrk, Cmakt, Cmtar, Cbesk, Csobs, Meins, Alpos, Alpgr, Alprf, Alpst, Sanko, Sanfe, Sanka, Dokar, Doknr, Dokvr,
    Doktl  sql.NullString
  Menge, Ewahr  sql.NullFloat64
  Posnrd,Postpd,Idnrkd,Cmaktd,Cmtard,Cbeskd,Csobsd,Meinsd,Alposd,Alpgrd,Alprfd,Alpstd,Sankod,Sanfed,Sankad,Dokard,Doknrd,Dokvrd,
    Doktld sql.NullString
  Menged,Ewahrd sql.NullString
}
type SWrkItmN struct {
  Seqnw, Ifldd, Items, Seqnh sql.NullInt64
  Gbcod, Bucod, Ident, Ibom1, Ibom2, Itmdf, Posnr sql.NullString
  P1, P2 SItmdtN
}
func (j *SWrkItmN) ScanWrkItm(ds *sql.Rows) {
  err := ds.Scan(&j.Seqnw, &j.Gbcod, &j.Bucod, &j.Ident, &j.Ibom1, &j.Ibom2, &j.Ifldd, &j.Itmdf, &j.Items,
    &j.P1.Systm,  &j.P1.Matnr, &j.P1.Maktx,  &j.P1.Werks, &j.P1.Stlal,  &j.P1.Mtart, &j.P1.Beskz,  &j.P1.Sobsl,
    &j.P1.Posnrd, &j.P1.Posnr, &j.P1.Postpd, &j.P1.Postp, &j.P1.Idnrkd, &j.P1.Idnrk, &j.P1.Cmaktd, &j.P1.Cmakt,
    &j.P1.Cmtard, &j.P1.Cmtar, &j.P1.Cbeskd, &j.P1.Cbesk, &j.P1.Csobsd, &j.P1.Csobs, &j.P1.Menged, &j.P1.Menge,
    &j.P1.Meinsd, &j.P1.Meins, &j.P1.Alposd, &j.P1.Alpos, &j.P1.Alpgrd, &j.P1.Alpgr, &j.P1.Alprfd, &j.P1.Alprf,
    &j.P1.Alpstd, &j.P1.Alpst, &j.P1.Ewahrd, &j.P1.Ewahr, &j.P1.Sankod, &j.P1.Sanko, &j.P1.Sanfed, &j.P1.Sanfe,
    &j.P1.Sankad, &j.P1.Sanka, &j.P1.Dokard, &j.P1.Dokar, &j.P1.Doknrd, &j.P1.Doknr, &j.P1.Dokvrd, &j.P1.Dokvr,
    &j.P1.Doktld, &j.P1.Doktl,
    &j.P2.Systm,  &j.P2.Matnr, &j.P2.Maktx,  &j.P2.Werks, &j.P2.Stlal,  &j.P2.Mtart, &j.P2.Beskz,  &j.P2.Sobsl,
    &j.P2.Posnrd, &j.P2.Posnr, &j.P2.Postpd, &j.P2.Postp, &j.P2.Idnrkd, &j.P2.Idnrk, &j.P2.Cmaktd, &j.P2.Cmakt,
    &j.P2.Cmtard, &j.P2.Cmtar, &j.P2.Cbeskd, &j.P2.Cbesk, &j.P2.Csobsd, &j.P2.Csobs, &j.P2.Menged, &j.P2.Menge,
    &j.P2.Meinsd, &j.P2.Meins, &j.P2.Alposd, &j.P2.Alpos, &j.P2.Alpgrd, &j.P2.Alpgr, &j.P2.Alprfd, &j.P2.Alprf,
    &j.P2.Alpstd, &j.P2.Alpst, &j.P2.Ewahrd, &j.P2.Ewahr, &j.P2.Sankod, &j.P2.Sanko, &j.P2.Sanfed, &j.P2.Sanfe,
    &j.P2.Sankad, &j.P2.Sanka, &j.P2.Dokard, &j.P2.Dokar, &j.P2.Doknrd, &j.P2.Doknr, &j.P2.Dokvrd, &j.P2.Dokvr,
    &j.P2.Doktld, &j.P2.Doktl,
    &j.Seqnh, &j.Posnr)
  if err != nil {
    log.Fatalf("Scan sql ScanWorkItm: %v\n", err)
  }
}

// Item comparison file writer
type SCmpdt struct {
  Systm, Werks, Matnr, Maktx, Stlal, Mtart, Beskz, Sobsl string
  Posnr, Postp, Idnrk, Cmakt, Cmtar, Cbesk, Csobs, Meins, Alpos, Alpgr, Alprf, Alpst, Sanko, Sanfe, Sanka, Dokar, Doknr, Dokvr,
    Doktl  string
  Menge, Ewahr  float64
  Posnrd,Postpd,Idnrkd,Cmaktd,Cmtard,Cbeskd,Csobsd,Meinsd,Alposd,Alpgrd,Alprfd,Alpstd,Sankod,Sanfed,Sankad,Dokard,Doknrd,Dokvrd,
    Doktld string
  Menged,Ewahrd string
}
type SCmpItm struct {
  Seqni, Ifldd, Items, Nline, Seqnh, Seqnw int
  Gbcod, Bucod, Ident, Ibom1, Ibom2, Itmdf, Posnr string
  P1, P2 SCmpdt
}
func (i *SCmpItm) Clear() {
  i.Ifldd, i.Items, i.Nline, i.Seqnh, i.Seqnw = 0, 0, 0, 0, 0
  i.Gbcod, i.Bucod, i.Ident, i.Ibom1, i.Ibom2, i.Itmdf, i.Posnr = "", "", "", "", "", "", ""
  i.P1.Systm, i.P1.Werks, i.P1.Matnr, i.P1.Maktx, i.P1.Stlal, i.P1.Mtart, i.P1.Beskz, i.P1.Sobsl = "", "", "", "", "", "", "", ""
  i.P1.Posnr, i.P1.Postp, i.P1.Idnrk, i.P1.Cmakt, i.P1.Cmtar, i.P1.Cbesk, i.P1.Csobs, i.P1.Meins = "", "", "", "", "", "", "", ""
  i.P1.Alpos, i.P1.Alpgr, i.P1.Alprf, i.P1.Alpst, i.P1.Sanko, i.P1.Sanfe, i.P1.Sanka, i.P1.Dokar = "", "", "", "", "", "", "", ""
  i.P1.Doknr, i.P1.Dokvr, i.P1.Doktl = "", "", ""
  i.P1.Menge, i.P1.Ewahr = 0, 0
  i.P1.Posnrd, i.P1.Postpd, i.P1.Idnrkd, i.P1.Cmaktd, i.P1.Cmtard, i.P1.Cbeskd, i.P1.Csobsd = "", "", "", "", "", "", ""
  i.P1.Meinsd, i.P1.Alposd, i.P1.Alpgrd, i.P1.Alprfd, i.P1.Alpstd, i.P1.Sankod, i.P1.Sanfed = "", "", "", "", "", "", ""
  i.P1.Sankad, i.P1.Dokard, i.P1.Doknrd, i.P1.Dokvrd, i.P1.Doktld = "", "", "", "", ""
  i.P1.Menged, i.P1.Ewahrd = "", ""
  i.P2.Systm, i.P2.Werks, i.P2.Matnr, i.P2.Maktx, i.P2.Stlal, i.P2.Mtart, i.P2.Beskz, i.P2.Sobsl = "", "", "", "", "", "", "", ""
  i.P2.Posnr, i.P2.Postp, i.P2.Idnrk, i.P2.Cmakt, i.P2.Cmtar, i.P2.Cbesk, i.P2.Csobs, i.P2.Meins = "", "", "", "", "", "", "", ""
  i.P2.Alpos, i.P2.Alpgr, i.P2.Alprf, i.P2.Alpst, i.P2.Sanko, i.P2.Sanfe, i.P2.Sanka, i.P2.Dokar = "", "", "", "", "", "", "", ""
  i.P2.Doknr, i.P2.Dokvr, i.P2.Doktl = "", "", ""
  i.P2.Menge, i.P2.Ewahr = 0, 0
  i.P2.Posnrd, i.P2.Postpd, i.P2.Idnrkd, i.P2.Cmaktd, i.P2.Cmtard, i.P2.Cbeskd, i.P2.Csobsd = "", "", "", "", "", "", ""
  i.P2.Meinsd, i.P2.Alposd, i.P2.Alpgrd, i.P2.Alprfd, i.P2.Alpstd, i.P2.Sankod, i.P2.Sanfed = "", "", "", "", "", "", ""
  i.P2.Sankad, i.P2.Dokard, i.P2.Doknrd, i.P2.Dokvrd, i.P2.Doktld = "", "", "", "", ""
  i.P2.Menged, i.P2.Ewahrd = "", ""
}
func (i *SCmpItm) MapItmCmp(j SWrkItmN) {
  i.Clear()
  if j.Seqnw.Valid { i.Seqnw = int(j.Seqnw.Int64) }
  if j.Gbcod.Valid { i.Gbcod = j.Gbcod.String }
  if j.Bucod.Valid { i.Bucod = j.Bucod.String }
  if j.Ident.Valid { i.Ident = j.Ident.String }
  if j.Ibom1.Valid { i.Ibom1 = j.Ibom1.String }
  if j.Ibom2.Valid { i.Ibom2 = j.Ibom2.String }
  if j.Ifldd.Valid { i.Ifldd = int(j.Ifldd.Int64) }
  if j.Itmdf.Valid { i.Itmdf = j.Itmdf.String }
  if j.Items.Valid { i.Items = int(j.Items.Int64) }
  if j.Seqnh.Valid { i.Seqnh = int(j.Seqnh.Int64) }
  if j.Posnr.Valid { i.Posnr = j.Posnr.String }
  if j.P1.Systm.Valid  { i.P1.Systm  = j.P1.Systm.String  }
  if j.P1.Matnr.Valid  { i.P1.Matnr  = j.P1.Matnr.String  }
  if j.P1.Maktx.Valid  { i.P1.Maktx  = j.P1.Maktx.String  }
  if j.P1.Werks.Valid  { i.P1.Werks  = j.P1.Werks.String  }
  if j.P1.Stlal.Valid  { i.P1.Stlal  = j.P1.Stlal.String  }
  if j.P1.Mtart.Valid  { i.P1.Mtart  = j.P1.Mtart.String  }
  if j.P1.Beskz.Valid  { i.P1.Beskz  = j.P1.Beskz.String  }
  if j.P1.Sobsl.Valid  { i.P1.Sobsl  = j.P1.Sobsl.String  }
  if j.P1.Posnrd.Valid { i.P1.Posnrd = j.P1.Posnrd.String }
  if j.P1.Posnr.Valid  { i.P1.Posnr  = j.P1.Posnr.String  }
  if j.P1.Postpd.Valid { i.P1.Postpd = j.P1.Postpd.String }
  if j.P1.Postp.Valid  { i.P1.Postp  = j.P1.Postp.String  }
  if j.P1.Idnrkd.Valid { i.P1.Idnrkd = j.P1.Idnrkd.String }
  if j.P1.Idnrk.Valid  { i.P1.Idnrk  = j.P1.Idnrk.String  }
  if j.P1.Cmaktd.Valid { i.P1.Cmaktd = j.P1.Cmaktd.String }
  if j.P1.Cmakt.Valid  { i.P1.Cmakt  = j.P1.Cmakt.String  }
  if j.P1.Cmtard.Valid { i.P1.Cmtard = j.P1.Cmtard.String }
  if j.P1.Cmtar.Valid  { i.P1.Cmtar  = j.P1.Cmtar.String  }
  if j.P1.Cbeskd.Valid { i.P1.Cbeskd = j.P1.Cbeskd.String }
  if j.P1.Cbesk.Valid  { i.P1.Cbesk  = j.P1.Cbesk.String  }
  if j.P1.Csobsd.Valid { i.P1.Csobsd = j.P1.Csobsd.String }
  if j.P1.Csobs.Valid  { i.P1.Csobs  = j.P1.Csobs.String  }
  if j.P1.Menged.Valid { i.P1.Menged = j.P1.Menged.String }
  if j.P1.Menge.Valid  { i.P1.Menge  = j.P1.Menge.Float64 }
  if j.P1.Meinsd.Valid { i.P1.Meinsd = j.P1.Meinsd.String }
  if j.P1.Meins.Valid  { i.P1.Meins  = j.P1.Meins.String  }
  if j.P1.Alposd.Valid { i.P1.Alposd = j.P1.Alposd.String }
  if j.P1.Alpos.Valid  { i.P1.Alpos  = j.P1.Alpos.String  }
  if j.P1.Alpgrd.Valid { i.P1.Alpgrd = j.P1.Alpgrd.String }
  if j.P1.Alpgr.Valid  { i.P1.Alpgr  = j.P1.Alpgr.String  }
  if j.P1.Alprfd.Valid { i.P1.Alprfd = j.P1.Alprfd.String }
  if j.P1.Alprf.Valid  { i.P1.Alprf  = j.P1.Alprf.String  }
  if j.P1.Alpstd.Valid { i.P1.Alpstd = j.P1.Alpstd.String }
  if j.P1.Alpst.Valid  { i.P1.Alpst  = j.P1.Alpst.String  }
  if j.P1.Ewahrd.Valid { i.P1.Ewahrd = j.P1.Ewahrd.String }
  if j.P1.Ewahr.Valid  { i.P1.Ewahr  = j.P1.Ewahr.Float64 }
  if j.P1.Sankod.Valid { i.P1.Sankod = j.P1.Sankod.String }
  if j.P1.Sanko.Valid  { i.P1.Sanko  = j.P1.Sanko.String  }
  if j.P1.Sanfed.Valid { i.P1.Sanfed = j.P1.Sanfed.String }
  if j.P1.Sanfe.Valid  { i.P1.Sanfe  = j.P1.Sanfe.String  }
  if j.P1.Sankad.Valid { i.P1.Sankad = j.P1.Sankad.String }
  if j.P1.Sanka.Valid  { i.P1.Sanka  = j.P1.Sanka.String  }
  if j.P1.Dokard.Valid { i.P1.Dokard = j.P1.Dokard.String }
  if j.P1.Dokar.Valid  { i.P1.Dokar  = j.P1.Dokar.String  }
  if j.P1.Doknrd.Valid { i.P1.Doknrd = j.P1.Doknrd.String }
  if j.P1.Doknr.Valid  { i.P1.Doknr  = j.P1.Doknr.String  }
  if j.P1.Dokvrd.Valid { i.P1.Dokvrd = j.P1.Dokvrd.String }
  if j.P1.Dokvr.Valid  { i.P1.Dokvr  = j.P1.Dokvr.String  }
  if j.P1.Doktld.Valid { i.P1.Doktld = j.P1.Doktld.String }
  if j.P1.Doktl.Valid  { i.P1.Doktl  = j.P1.Doktl.String  }
  if j.P2.Systm.Valid  { i.P2.Systm  = j.P2.Systm.String  }
  if j.P2.Matnr.Valid  { i.P2.Matnr  = j.P2.Matnr.String  }
  if j.P2.Maktx.Valid  { i.P2.Maktx  = j.P2.Maktx.String  }
  if j.P2.Werks.Valid  { i.P2.Werks  = j.P2.Werks.String  }
  if j.P2.Stlal.Valid  { i.P2.Stlal  = j.P2.Stlal.String  }
  if j.P2.Mtart.Valid  { i.P2.Mtart  = j.P2.Mtart.String  }
  if j.P2.Beskz.Valid  { i.P2.Beskz  = j.P2.Beskz.String  }
  if j.P2.Sobsl.Valid  { i.P2.Sobsl  = j.P2.Sobsl.String  }
  if j.P2.Posnrd.Valid { i.P2.Posnrd = j.P2.Posnrd.String }
  if j.P2.Posnr.Valid  { i.P2.Posnr  = j.P2.Posnr.String  }
  if j.P2.Postpd.Valid { i.P2.Postpd = j.P2.Postpd.String }
  if j.P2.Postp.Valid  { i.P2.Postp  = j.P2.Postp.String  }
  if j.P2.Idnrkd.Valid { i.P2.Idnrkd = j.P2.Idnrkd.String }
  if j.P2.Idnrk.Valid  { i.P2.Idnrk  = j.P2.Idnrk.String  }
  if j.P2.Cmaktd.Valid { i.P2.Cmaktd = j.P2.Cmaktd.String }
  if j.P2.Cmakt.Valid  { i.P2.Cmakt  = j.P2.Cmakt.String  }
  if j.P2.Cmtard.Valid { i.P2.Cmtard = j.P2.Cmtard.String }
  if j.P2.Cmtar.Valid  { i.P2.Cmtar  = j.P2.Cmtar.String  }
  if j.P2.Cbeskd.Valid { i.P2.Cbeskd = j.P2.Cbeskd.String }
  if j.P2.Cbesk.Valid  { i.P2.Cbesk  = j.P2.Cbesk.String  }
  if j.P2.Csobsd.Valid { i.P2.Csobsd = j.P2.Csobsd.String }
  if j.P2.Csobs.Valid  { i.P2.Csobs  = j.P2.Csobs.String  }
  if j.P2.Menged.Valid { i.P2.Menged = j.P2.Menged.String }
  if j.P2.Menge.Valid  { i.P2.Menge  = j.P2.Menge.Float64 }
  if j.P2.Meinsd.Valid { i.P2.Meinsd = j.P2.Meinsd.String }
  if j.P2.Meins.Valid  { i.P2.Meins  = j.P2.Meins.String  }
  if j.P2.Alposd.Valid { i.P2.Alposd = j.P2.Alposd.String }
  if j.P2.Alpos.Valid  { i.P2.Alpos  = j.P2.Alpos.String  }
  if j.P2.Alpgrd.Valid { i.P2.Alpgrd = j.P2.Alpgrd.String }
  if j.P2.Alpgr.Valid  { i.P2.Alpgr  = j.P2.Alpgr.String  }
  if j.P2.Alprfd.Valid { i.P2.Alprfd = j.P2.Alprfd.String }
  if j.P2.Alprf.Valid  { i.P2.Alprf  = j.P2.Alprf.String  }
  if j.P2.Alpstd.Valid { i.P2.Alpstd = j.P2.Alpstd.String }
  if j.P2.Alpst.Valid  { i.P2.Alpst  = j.P2.Alpst.String  }
  if j.P2.Ewahrd.Valid { i.P2.Ewahrd = j.P2.Ewahrd.String }
  if j.P2.Ewahr.Valid  { i.P2.Ewahr  = j.P2.Ewahr.Float64 }
  if j.P2.Sankod.Valid { i.P2.Sankod = j.P2.Sankod.String }
  if j.P2.Sanko.Valid  { i.P2.Sanko  = j.P2.Sanko.String  }
  if j.P2.Sanfed.Valid { i.P2.Sanfed = j.P2.Sanfed.String }
  if j.P2.Sanfe.Valid  { i.P2.Sanfe  = j.P2.Sanfe.String  }
  if j.P2.Sankad.Valid { i.P2.Sankad = j.P2.Sankad.String }
  if j.P2.Sanka.Valid  { i.P2.Sanka  = j.P2.Sanka.String  }
  if j.P2.Dokard.Valid { i.P2.Dokard = j.P2.Dokard.String }
  if j.P2.Dokar.Valid  { i.P2.Dokar  = j.P2.Dokar.String  }
  if j.P2.Doknrd.Valid { i.P2.Doknrd = j.P2.Doknrd.String }
  if j.P2.Doknr.Valid  { i.P2.Doknr  = j.P2.Doknr.String  }
  if j.P2.Dokvrd.Valid { i.P2.Dokvrd = j.P2.Dokvrd.String }
  if j.P2.Dokvr.Valid  { i.P2.Dokvr  = j.P2.Dokvr.String  }
  if j.P2.Doktld.Valid { i.P2.Doktld = j.P2.Doktld.String }
  if j.P2.Doktl.Valid  { i.P2.Doktl  = j.P2.Doktl.String  }
}
func (i *SCmpItm) IsrtCmpItm(db *sql.DB) {
  sSeqni := strconv.Itoa(i.Seqni)
  sIfldd := strconv.Itoa(i.Ifldd)
  sNline := strconv.Itoa(i.Nline)
  sSeqnh := strconv.Itoa(i.Seqnh)
  sSeqnw := strconv.Itoa(i.Seqnw)
  smeng1 := strconv.FormatFloat(float64(i.P1.Menge), 'f', -1, 32)
  sewah1 := strconv.FormatFloat(float64(i.P1.Ewahr), 'f', -1, 32)
  smeng2 := strconv.FormatFloat(float64(i.P2.Menge), 'f', -1, 32)
  sewah2 := strconv.FormatFloat(float64(i.P2.Ewahr), 'f', -1, 32)
  _, err := db.Exec(Q.Cmd["isrBomItmComp"], sSeqni, i.Gbcod, i.Bucod, i.Ident, i.Ibom1, i.Ibom2, sIfldd, i.Itmdf, sNline,
    i.P1.Systm,  i.P1.Matnr, i.P1.Maktx,  i.P1.Werks, i.P1.Stlal,  i.P1.Mtart, i.P1.Beskz,  i.P1.Sobsl, i.P1.Posnrd, i.P1.Posnr,
    i.P1.Postpd, i.P1.Postp, i.P1.Idnrkd, i.P1.Idnrk, i.P1.Cmaktd, i.P1.Cmakt, i.P1.Cmtard, i.P1.Cmtar, i.P1.Cbeskd, i.P1.Cbesk,
    i.P1.Csobsd, i.P1.Csobs, i.P1.Menged, smeng1,     i.P1.Meinsd, i.P1.Meins, i.P1.Alposd, i.P1.Alpos, i.P1.Alpgrd, i.P1.Alpgr,
    i.P1.Alprfd, i.P1.Alprf, i.P1.Alpstd, i.P1.Alpst, i.P1.Ewahrd, sewah1,     i.P1.Sankod, i.P1.Sanko, i.P1.Sanfed, i.P1.Sanfe,
    i.P1.Sankad, i.P1.Sanka, i.P1.Dokard, i.P1.Dokar, i.P1.Doknrd, i.P1.Doknr, i.P1.Dokvrd, i.P1.Dokvr, i.P1.Doktld, i.P1.Doktl,
    i.P2.Systm,  i.P2.Matnr, i.P2.Maktx,  i.P2.Werks, i.P2.Stlal,  i.P2.Mtart, i.P2.Beskz,  i.P2.Sobsl, i.P2.Posnrd, i.P2.Posnr,
    i.P2.Postpd, i.P2.Postp, i.P2.Idnrkd, i.P2.Idnrk, i.P2.Cmaktd, i.P2.Cmakt, i.P2.Cmtard, i.P2.Cmtar, i.P2.Cbeskd, i.P2.Cbesk,
    i.P2.Csobsd, i.P2.Csobs, i.P2.Menged, smeng2,     i.P2.Meinsd, i.P2.Meins, i.P2.Alposd, i.P2.Alpos, i.P2.Alpgrd, i.P2.Alpgr,
    i.P2.Alprfd, i.P2.Alprf, i.P2.Alpstd, i.P2.Alpst, i.P2.Ewahrd, sewah2,     i.P2.Sankod, i.P2.Sanko, i.P2.Sanfed, i.P2.Sanfe,
    i.P2.Sankad, i.P2.Sanka, i.P2.Dokard, i.P2.Dokar, i.P2.Doknrd, i.P2.Doknr, i.P2.Dokvrd, i.P2.Dokvr, i.P2.Doktld, i.P2.Doktl,
    sSeqnh, i.Posnr, sSeqnw)
  if err != nil {
    log.Fatalf("Execute sql isrBomItmComp: %v\n", err)
  }
}
