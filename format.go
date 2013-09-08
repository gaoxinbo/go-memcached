package go_memcached

import (
    "fmt"
    )

func getDeletionCommond(key []byte) string{
    return fmt.Sprintf("delete %s\r\n",string(key))
}

func getStorageCommond(key,value []byte,flag,expire int, op string) string{
    return fmt.Sprintf("%s %s %d %d %d\r\n%s\r\n",op,string(key),flag,expire,len(value),string(value))
}
