package bomcomp

import "encoding/xml"
import "flag"
import "io/ioutil"
import "log"
import "os"
import "strconv"

type Settings struct {
  Dbase string `xml:"database,attr"`
  Gbcod string `xml:"gb,attr"`
  Bucod string `xml:"bu,attr"`
  Citms string `xml:"cmpItems,attr"`
  Level string `xml:"level,attr"`
  Cmode string `xml:"cmpMode,attr"`
  Fcomp string `xml:"fullCmp,attr"`
  Rules string `xml:"rulesSet,attr"`
  Ptole string `xml:"convTol,attr"`
  Hkeep string `xml:"hkeeping,attr"`
  Dstep string `xml:"dfltSteps,attr"`
  Uomfl string `xml:"uomFile,attr"`
  Sqlfl string `xml:"sqlCmdFile,attr"`
  Bitm1 string `xml:"bomItm1,attr"`
  Bitm2 string `xml:"bomItm2,attr"`
  Pstep string
  Convt float32
}

func (s *Settings) NewSettings(fname string) {
  f, err := os.Open(fname)
  if err != nil {
    log.Fatalf("Open Config file: %v\n", err)
  }
  defer f.Close()
  xmlv, _ := ioutil.ReadAll(f)
  err = xml.Unmarshal(xmlv, &s)
  if err != nil {
    log.Fatalf("Unmarshall Config file: %v\n", err)
  }
  temp, _ := strconv.ParseFloat(s.Ptole, 32)
  s.Convt = float32(temp)
  pstep := flag.String("steps", s.Dstep, "Processing steps: [LMC|other combination]. L=Local copy XLS files, M=Match list prep, C=Compare BOMs")
  flag.Parse()
  s.Pstep = *pstep
}
