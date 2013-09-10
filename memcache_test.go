package go_memcached


import (
    "testing"
    "bytes"
    )

// change this config item if necessary
var host = "localhost:11211"
var c Client
var initial = false

func TestConnect(t *testing.T){
  var n Client
  err := n.Connect(host)
  if err != nil {
    t.Errorf("can`t connect to %s",host)
  }
  n.Close()
}

func TestAdd(t *testing.T){
  if initial == false{
    t.Errorf("client didn`t connent")
  }

  c.Delete([]byte("key"))
  b,err:=c.Add([]byte("key"),[]byte("value"))
  if err != nil{
    t.Errorf("add error %s",err)
  }

  if bytes.Compare(b,[]byte("STORED\r\n")) !=0 {
    t.Errorf("add error result is %s", string(b))
  }

  /*
  v,err := c.Get([]byte("key"))
  fmt.Println(string(v.Value))


  m,err := c.Gets(bytes.Split([]byte("hello key good day"),[]byte(" ")))
  for key,value := range m {
    fmt.Println("key " + key)
    fmt.Println("value " + string(value.Value))
  }
  */
}
func init(){
  err := c.Connect(host)
  if err == nil {
    initial = true
  }
}
