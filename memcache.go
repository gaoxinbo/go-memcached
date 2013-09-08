package go_memcached

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

func (c *Client) Close(){
  c.conn.Close()
  c.conn = nil
}

func (c *Client) Stats() ([]byte,error){
  if c.conn == nil {
    return nil,errors.New("not connected!")
  }

  c.conn.Write([]byte("stats\r\n"))
  return c.readMulLines()

}

func (c *Client) deletionCommand(key []byte)([]byte,error){
  err := c.check()
  if err != nil {
    return nil,err
  }

  c.send(getDeletionCommond(key))
  return c.readLine()
}

func (c *Client) send(cmd string) error{
  c.conn.Write([]byte(cmd))
  // TODO handle error
  return nil
}

func (c *Client) storageCommand(key,value []byte, flag, expire int, op string) ([]byte,error){
  err := c.check()
  if err != nil {
    return nil,err
  }

  err = c.send(getStorageCommond(key,value,flag,expire,op))
  if err != nil {
    return nil,err
  }

  return c.readLine()
}

func (c *Client) Replace(key,value []byte) ([]byte,error){
  return c.ReplaceWithExpire(key,value,0,0)
}

func (c *Client) ReplaceWithExpire(key,value []byte,flag, expire int) ([]byte,error){
  return c.storageCommand(key,value,flag,expire,"replace")
}

func (c *Client) Add(key,value []byte) ([]byte,error){
  return c.AddWithExpire(key,value,0,0)
}

func (c *Client) AddWithExpire(key,value []byte, flag,expire int)([]byte, error){
  return c.storageCommand(key,value,flag,expire,"add")
}

func (c *Client) Append(key,value []byte) ([]byte,error){
  return c.AppendWithExpire(key,value,0,0)
}

func (c *Client) AppendWithExpire(key,value []byte, flag,expire int)([]byte, error){
  return c.storageCommand(key,value,flag,expire,"append")
}

func (c *Client) Prepend(key,value []byte) ([]byte,error){
  return c.PrependWithExpire(key,value,0,0)
}

func (c *Client) PrependWithExpire(key,value []byte, flag,expire int)([]byte, error){
  return c.storageCommand(key,value,flag,expire,"prepend")
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

func (c *Client) Delete(key []byte)([]byte,error){
  return c.deletionCommand(key)
}

/*
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
*/
