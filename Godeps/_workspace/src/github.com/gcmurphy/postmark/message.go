package postmark
import (
    "os"
    "fmt"
    "io/ioutil"
    "encoding/json"
    "encoding/base64"
    "path"
    "mime"
)

type Header struct {
    Name string
    Value string
}

type Attachment struct {
    Name string
    Content string // Base 64 encoded string
    ContentType string
}

type Response struct {
    ErrorCode int
    Message string
    MessageID string
    SubmittedAt string //Date 
    To string
}

type BatchResponse []Response

type Message struct {
    From string
    To string
    Cc string
    Bcc string
    Subject string
    Tag string
    HtmlBody string
    TextBody string
    ReplyTo string
    Headers []Header
    Attachments []Attachment
}

type BatchMessage []Message


func (p *Message) String() string{
    js, e := json.MarshalIndent(p, "", "")
    if e != nil {
        return ""
    }
    return string(js)
}

// Attach file to message (base64 encoded)
func (p *Message) Attach(file string)(error){

    finfo, e := os.Stat(file)
    if e != nil {
        return e
    }

    if finfo.Size() > int64(10e6){
        return fmt.Errorf("File size %d exceeds 10MB limit.", finfo.Size())
    }

    fh, e := os.Open(file)
    if e != nil {
        return e
    }

    // Even though we only have 10MB limit..
    // I probably shouldn't do this..
    cnt, e := ioutil.ReadAll(fh)
    if e != nil {
        return e
    }
    fh.Close()

    mimeType := mime.TypeByExtension(path.Ext(file))
    if len(mimeType) == 0 {
        return fmt.Errorf("Unknown mime type for attachment: %s", file)
    }

    attachment := Attachment{
        Name: finfo.Name(),
        Content: base64.StdEncoding.EncodeToString(cnt),
        ContentType: mimeType,
    }
    p.Attachments = append(p.Attachments, attachment)
    return nil
}

func unmarshal (msg []byte, i interface{})(error){
    e := json.Unmarshal(msg, i)
    if e != nil {
        return e
    }
    return nil
}

func (m *Message) Marshal()([]byte, error){
    return json.Marshal(*m)
}

func UnmarshalMessage(msg []byte)(*Message, error){
    var m Message
    e := unmarshal(msg, &m)
    return &m, e
}

func (r *Response) Marshal()([]byte, error){
    return json.Marshal(*r)
}

func UnmarshalResponse(rsp []byte)(*Response, error){
    var r Response
    e := unmarshal(rsp, &r)
    return &r, e
}

func (r *Response) String() string{
    js, e := json.MarshalIndent(r, "", "")
    if e != nil {
        return ""
    }
    return string(js)
}
