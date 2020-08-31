package main

import (
    "encoding/json"
    "fmt"
    "wxbizjsonmsgcrypt"
    "net/http"
    "log"
    "strings"
    "net/url"
    "io/ioutil"
)

const token = "gY1AGR3mjBhzy"
const receiverId = "wwabfd0cec7171e769"
const encodingAeskey = "g8VGfQEqluUhoKOlyjmmll8Q9C5tVFUTX5T2qkmI9Sv"

func getString(str, endstr string,start int, msg * string) int{
    end := strings.Index(str,endstr)
    *msg = str[start:end]
    return end + len(endstr)
}

func VerifyURL(w http.ResponseWriter, r *http.Request) {
    //httpstr := `&{GET /?msg_signature=825075c093249d5a60967fe4a613cae93146636b&timestamp=1597998748&nonce=1597483820&echostr=neLB8CftccHiz19tluVb%2BUBnUVMT3xpUMZU8qvDdD17eH8XfEsbPYC%2FkJyPsZOOc6GdsCeu8jSIa2noSJ%2Fez2w%3D%3D HTTP/1.1 1 1 map[Cache-Control:[no-cache] Accept:[*/*] Pragma:[no-cache] User-Agent:[Mozilla/4.0]] 0x86c180 0 [] false 100.108.211.112:8893 map[] map[] <nil> map[] 100.108.79.233:59663 /?msg_signature=825075c093249d5a60967fe4a613cae93146636b&timestamp=1597998748&nonce=1597483820&echostr=neLB8CftccHiz19tluVb%2BUBnUVMT3xpUMZU8qvDdD17eH8XfEsbPYC%2FkJyPsZOOc6GdsCeu8jSIa2noSJ%2Fez2w%3D%3D <nil>}`
    fmt.Println(r,r.Body)
    httpstr := r.URL.RawQuery
    start := strings.Index(httpstr,"msg_signature=")
    start += len("msg_signature=")
    
    var msg_signature string
    next := getString(httpstr,"&timestamp=",start, &msg_signature)

    var timestamp string
    next = getString(httpstr,"&nonce=",next, &timestamp)

    var nonce string
    next = getString(httpstr,"&echostr=",next, &nonce)

    echostr := httpstr[next:len(httpstr)]

    echostr,_ = url.QueryUnescape(echostr)
    fmt.Println(msg_signature,timestamp,nonce,echostr,next)

    wxcpt := wxbizjsonmsgcrypt.NewWXBizMsgCrypt(token, encodingAeskey, receiverId, wxbizjsonmsgcrypt.JsonType)
    echoStr, cryptErr := wxcpt.VerifyURL(msg_signature, timestamp, nonce, echostr)
    if nil != cryptErr {
        fmt.Println("verifyUrl fail", cryptErr)
    }
    fmt.Println("verifyUrl success echoStr", string(echoStr))
    fmt.Fprintf(w,string(echoStr))


}


type MsgContent struct {
    ToUsername   string `json:"ToUserName"`
    FromUsername string `json:"FromUserName"`
    CreateTime   uint32 `json:"CreateTime"`
    MsgType      string `json:"MsgType"`
    Content      string `json:"Content"`
    Msgid        uint64 `json:"MsgId"`
    Agentid      uint32 `json:"AgentId"`
}

func MsgHandler(w http.ResponseWriter, r *http.Request){
    httpstr := r.URL.RawQuery
    start := strings.Index(httpstr,"msg_signature=")
    start += len("msg_signature=")

    var msg_signature string
    next := getString(httpstr,"&timestamp=",start, &msg_signature)

    var timestamp string
    next = getString(httpstr,"&nonce=",next, &timestamp)

    nonce := httpstr[next:len(httpstr)]
    fmt.Println(msg_signature,timestamp,nonce)
  
    body, err := ioutil.ReadAll(r.Body)
    fmt.Println(string(body),err)
    wxcpt := wxbizjsonmsgcrypt.NewWXBizMsgCrypt(token, encodingAeskey, receiverId, wxbizjsonmsgcrypt.JsonType)
 
    msg, err_ := wxcpt.DecryptMsg(msg_signature,timestamp,nonce,body)
    fmt.Println(string(msg),err_)
    var msgContent MsgContent
    err = json.Unmarshal(msg, &msgContent)
    if nil != err {
        fmt.Println("Unmarshal fail", err)
    } else {
        fmt.Println("struct", msgContent)
    }
    

    fmt.Println(msgContent,err)
    ToUsername := msgContent.ToUsername
    msgContent.ToUsername = msgContent.FromUsername
    msgContent.FromUsername = ToUsername
    fmt.Println("replaymsg", msgContent)
    replayJson,err := json.Marshal(&msgContent)
   
    encryptMsg, cryptErr := wxcpt.EncryptMsg(string(replayJson), "1409659589", "1409659589")
    if nil != cryptErr {
        fmt.Println("DecryptMsg fail", cryptErr)
    }

    sEncryptMsg := string(encryptMsg)
    
    fmt.Println("after encrypt sEncryptMsg: ", sEncryptMsg)
    fmt.Fprintf(w,sEncryptMsg);
}


func CallbackHandler(w http.ResponseWriter, r *http.Request) {
    httpstr := r.URL.RawQuery
    echo := strings.Index(httpstr,"echostr")
    if (echo != -1){
        VerifyURL(w,r)
    }else{
        MsgHandler(w,r)
    }
    
    fmt.Println("finished CallbackHandler",httpstr)
}

func main() {
    http.HandleFunc("/", CallbackHandler)               //      设置访问路由
    log.Fatal(http.ListenAndServe(":8893", nil))
}