package main

import "bar8tl/p/bomcomp"
import "flag"
import "fmt"

func main() {
  proc := flag.String("p", "na", "Procedure: [uomgen|uomdsp|zipgen|zipdsp]")
  ifnm := flag.String("i", "na", "Input file name")
  ofnm := flag.String("o", "na", "Output file name")
  flag.Parse()
  var t bomcomp.Umline
  switch *proc {
    case "uomgen" :  t.Uomgen(*ifnm, *ofnm)
    case "uomdsp" :  t.Uomdsp(*ifnm)
    case "zipgen" : bomcomp.Zipgen(*ifnm, *ofnm)
    case "zipdsp" : bomcomp.Zipdsp(*ifnm)
    case "na" : fmt.Println("Option invalid")
  }
}
