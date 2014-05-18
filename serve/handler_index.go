package serve
import (
  "net/http"
  "github.com/hidu/goutils"
  "text/template"
  "strings"
  "github.com/googollee/go-socket.io"
  "log"
//  "fmt"
)
/**
*https://github.com/googollee/go-socket.io
*/
func new_req(ns *socketio.NameSpace, req_path string,host string) {
//  var name string
//  name = ns.Session.Values["name"].(string)
//  fmt.Printf("%s said in %s, title: %s, body: %s, article number: %i", name, ns.Endpoint(), title, body, article_num)
}
func send_req(ns *socketio.NameSpace,sid int64,host string,req_path string) {
  ns.Emit("req",sid,host,req_path)
}

func (ser *ProxyServe)initWs(){
	sock_config := &socketio.Config{HeartbeatTimeout:2,ClosingTimeout:4}
	ser.ws= socketio.NewSocketIOServer(sock_config)
	ser.wsClients=make(map[string]*wsClient)
	ser.ws.On("connect", func(ns *socketio.NameSpace){
	  log.Println("ws connected",ns.Id()," in channel ", ns.Endpoint())
	  ser.wsClients[ns.Id()]=&wsClient{ns:ns,user:"guest"}
	})
	ser.ws.On("disconnect", func(ns *socketio.NameSpace){
	  log.Println("ws disconnect",ns.Id()," in channel ", ns.Endpoint())
	  if _,has:=ser.wsClients[ns.Id()];has{
	    delete(ser.wsClients,ns.Id())
	  }
	})
//	ser.ws.On("req", new_req)
}

func (ser *ProxyServe)handleLocalReq(w http.ResponseWriter, req *http.Request){
   if(strings.HasPrefix(req.URL.Path,"/socket.io/1/")){
			ser.ws.ServeHTTP(w,req)
			return;
	 }
   if(req.Method=="GET"){
	    if(strings.HasPrefix(req.URL.Path,"/res/")){
	       goutils.DefaultResource.HandleStatic(w,req,req.URL.Path)
	    }else{
		   msg:=goutils.DefaultResource.Load("/res/tpl/index.html")
		   tpl,_:=template.New("page").Parse(string(msg))
		   values :=make(map[string]string)
		   values["host"]=req.Host
		   values["title"]=""
		   values["version"]="0.1"
		   tpl.Execute(w,values)
	    }
   }
}

