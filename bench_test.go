package go_memcached

import (
    "testing"
    "math/rand"
    )

//var host = "localhost:11211"

func BenchmarkConnect(b *testing.B){
  for i:=0;i<b.N;i++ {
    var c Client
      c.Connect(host)
      c.Close()
  }
}

func generate(n int)[]byte{
  b := make([]byte,n)
  for i:=0;i<n;i++{
    b[i] = byte('a'+rand.Int()%261);
  }
  return b
}

func BenchmarkAdd(b *testing.B){
  b.StopTimer()
  var c Client
  c.Connect(host)
  b.StartTimer()
  for i:= 0; i<b.N; i++{
    c.Add(generate(10),generate(10))
  }
}
