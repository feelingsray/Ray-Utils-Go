package encode

import (
    "bytes"
    "encoding/binary"
    "encoding/hex"
    "fmt"
)

func HexStringToByte(src string) ([]byte, error) {
    var dst []byte
    _, err := fmt.Sscanf(src, "%X", &dst)
    if err != nil {
        return nil, err
    }
    return dst, nil
}

func ByteToHexString(src []byte) string {
    dst := hex.EncodeToString(src)
    return dst
}

func BytesToInt(bys []byte) (int, error) {
    bytebuff := bytes.NewBuffer(bys)
    var data int64
    err := binary.Read(bytebuff, binary.BigEndian, &data)
    if err != nil {
        return 0, err
    }
    return int(data), nil
}
