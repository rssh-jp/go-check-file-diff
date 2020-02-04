# go-check-file-modify

## usage
```
package main
import(
    "log"
    "bytes"

    "github.com/rssh-jp/go-check-file-diff"
)

func main(){
    file1 := bytes.NewBufferString("012345")
    file2 := bytes.NewBufferString("012345")
    issame, err := checkfilediff.IsSame(file1, file2)
    if err != nil{
        log.Fatal(err)
    }

    log.Println(issame)
}
```
