[![GoDoc](https://godoc.org/github.com/KarpelesLab/rc?status.svg)](https://godoc.org/github.com/KarpelesLab/rc)

# rc

A simple object to be used when implementing `io.ReaderFrom` to keep track of how many bytes
have been read from a given stream.

## Example

```go
func (obj *myObject) ReadFrom(r io.Reader) (int64, error) {
    rc := rc.New(r)

    err := binary.Read(rc, binary.BigEndian, &obj.Value)
    if err != nil {
        // return the error
        return rc.Error64(err)
    }
    // more reading

    return rc.Ret64()
}
