package go_memcached


import (
    "testing"
    "bytes"
    )

// change this config item if necessary
var host = "localhost:11211"

func TestConnect(t *testing.T){
  var c Client
  err := c.Connect(host)
  if err != nil {
    t.Errorf("can`t connect to %s",host)
  }
  c.Close()
}

func TestAdd(t *testing.T){
  var c Client
  err := c.Connect(host)
  if err != nil {
    t.Errorf("can`t connect to %s",host)
  }

  c.Delete([]byte("key"))
  b,err:=c.Add([]byte("key"),[]byte("value"))
  if err != nil{
    t.Errorf("add error %s",err)
  }

  if bytes.Compare(b,[]byte("STORED\r\n")) !=0 {
    t.Errorf("add error result is %s", string(b))
  }
}
