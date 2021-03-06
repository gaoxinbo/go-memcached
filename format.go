package go_memcached

import (
    "fmt"
    "bytes"
    )

func getDeletionCommond(key []byte) string{
    return fmt.Sprintf("delete %s\r\n",string(key))
}

func getStorageCommond(key,value []byte,flag,expire int, op string) string{
    return fmt.Sprintf("%s %s %d %d %d\r\n%s\r\n",op,string(key),flag,expire,len(value),string(value))
}

func getGetCommand(key [][]byte) string {
  b := bytes.Join(key,[]byte(" "))
  return fmt.Sprintf("get %s\r\n", string(b))
}

func getIncrOrDecrCommand(key,value []byte,op string) string{
  return fmt.Sprintf("%s %s %s\r\n" ,op, string(key), string(value))
}

func getTouchCommand(key[]byte,expire int)string{
  return fmt.Sprintf("touch %s %d\r\n", string(key),expire);
}
