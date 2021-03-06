package go_memcached

import (
    "net"
    "errors"
    "bufio"
    "bytes"
    "strconv"
    )

type Client struct {
  conn net.Conn
  reader *bufio.Reader
}

type Value struct{
  Value []byte
  Flag int
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
func (c *Client) send(cmd string) error{
  c.conn.Write([]byte(cmd))
  // TODO handle error
  return nil
}

func (c *Client) singleLineCommand(cmd string)([]byte,error){
  err := c.check()
  if err != nil {
    return nil,err
  }
  c.send(cmd)
  return c.readLine()
}

func (c *Client) multipleLinesCommand(cmd string)([]byte, error){
  err := c.check()
  if err != nil {
    return nil,err
  }
  c.send(cmd)
  return c.readMulLines()
}

func (c *Client) Replace(key,value []byte) ([]byte,error){
  return c.ReplaceWithExpire(key,value,0,0)
}

func (c *Client) ReplaceWithExpire(key,value []byte,flag, expire int) ([]byte,error){
  return c.singleLineCommand(getStorageCommond(key,value,flag,expire,"replace"))
}

func (c *Client) Add(key,value []byte) ([]byte,error){
  return c.AddWithExpire(key,value,0,0)
}

func (c *Client) AddWithExpire(key,value []byte, flag,expire int)([]byte, error){
  return c.singleLineCommand(getStorageCommond(key,value,flag,expire,"add"))
}

func (c *Client) Set(key,value []byte) ([]byte,error){
  return c.SetWithExpire(key,value,0,0)
}

func (c *Client) SetWithExpire(key,value []byte, flag,expire int)([]byte, error){
  return c.singleLineCommand(getStorageCommond(key,value,flag,expire,"set"))
}

func (c *Client) Append(key,value []byte) ([]byte,error){
  return c.AppendWithExpire(key,value,0,0)
}

func (c *Client) AppendWithExpire(key,value []byte, flag,expire int)([]byte, error){
  return c.singleLineCommand(getStorageCommond(key,value,flag,expire,"append"))
}

func (c *Client) Prepend(key,value []byte) ([]byte,error){
  return c.PrependWithExpire(key,value,0,0)
}

func (c *Client) PrependWithExpire(key,value []byte, flag,expire int)([]byte, error){
  return c.singleLineCommand(getStorageCommond(key,value,flag,expire,"prepend"))
}

// retrive one key
func (c *Client) Get(key []byte)(*Value,error) {
  var k = make([][]byte,1)
  k = append(k,key)

  m,err := c.Gets(k)
  if err != nil {
    return nil,err
  }
  v,b:= m[string(key)]
  if b == false {
    return nil,nil
  }

  return &v,nil
}

// get multiple keys
func (c *Client) Gets(key [][]byte) (map[string]Value,error){
  b,err := c.multipleLinesCommand(getGetCommand(key))
  //b,err := c.getCommand(key)
  if err != nil{
    return nil,err
  }
  m := parseGet(b)
  return m,nil
}

func parseGet(b []byte) (map[string] Value){
  m := make(map[string]Value)
  var key string
  var flag int
  var length int
  results := bytes.Split(b,[]byte("\r\n"))

  // last one is "END\r\n"
  for i:=0;i<len(results)-1;i++ {
    if i % 2 == 0 {
      items := bytes.Split(results[i],[]byte(" "))
      // invalid 
      if len(items) < 4 || bytes.Compare(items[0], []byte("VALUE"))!=0 {
        key = ""
      }else{
        key = string(items[1])
        flag,_ = strconv.Atoi(string(items[2]))
        length,_ = strconv.Atoi(string(items[3]))
      }
    } else {
      if len(key) == 0 {
        continue
      } else{
        t := results[i][0:length]
        m[key] = Value{Value:t, Flag:flag}
      }
    }
  }
  return m
}

func (c *Client) Delete(key []byte)([]byte,error){
  return c.singleLineCommand(getDeletionCommond(key))
}

func (c *Client) Incr(key,value []byte) ([]byte,error){
  return c.singleLineCommand(getIncrOrDecrCommand(key,value,"incr"))
}

func (c *Client) Decr(key, value[]byte) ([]byte, error){
  return c.singleLineCommand(getIncrOrDecrCommand(key,value,"decr"))
}

func (c *Client) Touch(key []byte, expire int)([]byte,error){
  return c.singleLineCommand(getTouchCommand(key,expire))
}

