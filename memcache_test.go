package go_memcached


import (
    "testing"
    "bytes"
    "strings"
    //"fmt"
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

  keys := []string{"one","two","three"}
  values:= []string{"v1","v2","v3"}

  // add and get keys one by one
  for i:=0 ; i<len(keys);i++ {
    c.Delete([]byte(keys[i]))
    b,err := c.Add([]byte(keys[i]),[]byte(values[i]))
    if err != nil {
      t.Errorf("add error %s",err)
    }
    if bytes.Compare(b,[]byte("STORED\r\n")) !=0 {
      t.Errorf("add error result is %s", string(b))
    }
    v,err := c.Get([]byte(keys[i]))
    if bytes.Compare(v.Value,[]byte(values[i])) != 0{
      t.Errorf("get key %s error, expect %s get %s",keys[i],values[i], string(v.Value))
    }

  }

  s := strings.Join(keys," ")

  // get multiple keys
  m,err:= c.Gets(bytes.Split([]byte(s),[]byte(" ")))
  if err != nil{
    t.Errorf("gets error")
  }

  for i:=0 ; i<len(keys); i++{
    v,e := m[keys[i]]
    if e == false{
      t.Errorf("key %s do not exist",keys[i])
    }

    if bytes.Compare(v.Value, []byte(values[i])) != 0{
      t.Errorf("get %s value error, expect %s, get %s", keys[i], values[i], string(v.Value))
    }
  }

  // add an existed key
  b,err := c.Add([]byte("one"),[]byte("hello"))
  if bytes.Compare(b,[]byte("NOT_STORED\r\n")) != 0 {
    t.Errorf("store an existed key error")
  }

}

func TestSet(t *testing.T){
  if initial == false{
    t.Errorf("client didn`t connent")
  }

  b,err := c.Set([]byte("set"),[]byte("good"))
  if err != nil {
    t.Errorf("set error")
  }

  if bytes.Compare(b,[]byte("STORED\r\n")) != 0{
    t.Errorf("set result error %s", string(b))
  }

  v,err := c.Get([]byte("set"))
    if bytes.Compare(v.Value,[]byte("good")) != 0{
      t.Errorf("get key %s error, expect %s get %s","set","good", string(v.Value))
    }

  // set to another value
  b,err = c.Set([]byte("set"),[]byte("good bye"))

  if bytes.Compare(b,[]byte("STORED\r\n")) != 0{
    t.Errorf("set result error %s", string(b))
  }

  v,err = c.Get([]byte("set"))
    if bytes.Compare(v.Value,[]byte("good bye")) != 0{
      t.Errorf("get key %s error, expect %s get %s","set","good bye", string(v.Value))
    }


}

func init(){
  err := c.Connect(host)
  if err == nil {
    initial = true
  }
}
