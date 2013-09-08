package main

import (
    "net"
    "errors"
    "fmt"
    "bufio"
    "bytes"
    )

type Client struct {
  conn net.Conn
  reader *bufio.Reader
}

func (c *Client) check() error {
  if c.conn == nil {
    return errors.New("not connected yet!")
  }
  return nil
}

func (c *Client) readMulLines() ([]byte,error){
  if c.conn == nil {
    return nil, errors.New("not connected yet!")
  }

  result := new(bytes.Buffer)

  for {
    b,err := c.reader.ReadBytes(byte('\n'))
    if err != nil {
      return nil,err
    }

    result.Write(b)
    b = b[0:len(b)-2]
    if bytes.Compare(b,[]byte("END")) == 0{
      break
    }
  }

  return result.Bytes(),nil
}

func (c *Client) readLine() ([]byte,error) {
  if c.conn == nil {
    return nil, errors.New("not connected yet!")
  }
  b, err := c.reader.ReadBytes(byte('\n'))
  return b,err
}

func (c *Client) Connect(host string) error {
  var err error
  c.conn, err = net.Dial("tcp",host)
  if err != nil {
    return err
  }
  c.reader = bufio.NewReader(c.conn)
  return nil
}

func (c *Client) Stats() ([]byte,error){
  if c.conn == nil {
    return nil,errors.New("not connected!")
  }

  c.conn.Write([]byte("stats\r\n"))
  return c.readMulLines()

}

func (c *Client) Add(key,value []byte) ([]byte,error){
  err := c.check()
  if err != nil {
    return nil,err
  }

  s := fmt.Sprintf("add %s 0 0 %d \r\n",string(key),len(value))
  c.conn.Write([]byte(s))
  s = fmt.Sprintf("%s\r\n", string(value))
  c.conn.Write([]byte(s))
  return c.readLine()

}

func (c *Client) Get(key []byte)([]byte,error) {
  err := c.check()
  if err != nil {
    return nil,err
  }

  s := fmt.Sprintf("get %s\r\n", string(key))
  c.conn.Write([]byte(s))

  b,err := c.readMulLines()
  if err != nil{
    return nil,err
  }

  items := bytes.Split(b,[]byte("\r\n"))
  if len(items) == 2 {
    return nil,nil
  }

  return items[1],nil

}

func main(){
  var c Client
  err := c.Connect("localhost:11211")
  if err != nil{
    panic(err)
  }

  b,err := c.Stats()
  if err != nil {
    panic(err)
  }

  fmt.Print(string(b))

  b, err = c.Add([]byte("hello"),[]byte("world"))
  fmt.Print(string(b))
  b,err = c.Get([]byte("hello"))
  fmt.Println(string(b))

  b,err = c.Get([]byte("good"))
  if b == nil {
    fmt.Println("can`t find good")
  }
}
