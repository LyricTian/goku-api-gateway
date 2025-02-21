package backend

import (
	"fmt"
	"github.com/eolinker/goku-api-gateway/config"
	"github.com/eolinker/goku-api-gateway/goku-node/common"
	"github.com/eolinker/goku-api-gateway/goku-service/application"
	"github.com/eolinker/goku-api-gateway/goku-service/balance"
	"io/ioutil"
	"strings"

	"github.com/eolinker/goku-api-gateway/node/gateway/application/interpreter"
	"github.com/eolinker/goku-api-gateway/node/gateway/response"
	"time"
)

type Proxy struct {
	//step * config.APIStepConfig
	BalanceName string
	Balance application.IHttpApplication
	HasBalance bool
	Protocol  string


	Method  string
	Path    interpreter.Interpreter
	OrgPath string
	Decode  response.DecodeHandle

	RequestPath string

	Retry   int
	TimeOut time.Duration

}
func NewProxyBackendTarget(step *config.APIStepConfig,requestPath string,balanceTarget string) *Proxy {
	b:= &Proxy{
		BalanceName:balanceTarget,
		Protocol:step.Proto,
		Method: strings.ToUpper(step.Method),
		Path: interpreter.GenPath( step.Path),

		RequestPath: requestPath,
		Decode: response.GetDecoder(step.Decode),


		TimeOut:time.Duration(step.TimeOut)*time.Millisecond,
		Retry: step.Retry,

	}


	b.Balance, b.HasBalance = balance.GetByName(balanceTarget)

	return b
}


func (b *Proxy) Send(ctx *common.Context,variables *interpreter.Variables)( *BackendResponse,error) {


	if !b.HasBalance{
		err := fmt.Errorf("get balance error:%s", b.BalanceName)
		return nil, err
	}





	path:= b.Path.Execution(variables)

	// 不是restful时，将匹配路由之后对url拼接到path之后
	if len(variables.Restful)==0{
		orgRequestUrl := ctx.RequestOrg.URL().RawPath
		lessPath := strings.TrimPrefix(orgRequestUrl,b.RequestPath)
		lessPath = strings.TrimPrefix(lessPath,"/")
		path = strings.TrimSuffix(path,"/")
		path = fmt.Sprint(path,"/",lessPath)
	}

	method:= b.Method
	if method == "FOLLOW"{
		method = ctx.ProxyRequest.Method
	}
	r, finalTargetServer, retryTargetServers, err := b.Balance.Send(b.Protocol,method , path, ctx.ProxyRequest.Querys(), ctx.ProxyRequest.Headers(),variables.Org, b.TimeOut, b.Retry)


	backendResponse := &BackendResponse{
		Method:            method,
		Protocol:           b.Protocol,
		StatusCode:200,
		Status:"200" ,
		//Response:           r,
		TargetUrl:          path,
		FinalTargetServer:  finalTargetServer,
		RetryTargetServers: retryTargetServers,
		Header:r.Header,
		//Cookies:r.Cookies(),
	}
	if err!=nil{
		backendResponse.StatusCode,backendResponse.Status = 503,"503"
		return backendResponse,err
	}
	defer r.Body.Close()
	backendResponse.BodyOrg, err = ioutil.ReadAll(r.Body)
	if err!= nil{
		return backendResponse,nil
	}

	if b.Decode!= nil{
		rp,e:=response.Decode(backendResponse.BodyOrg,b.Decode)
		if e!= nil{
			backendResponse.Body = nil
		}else {
			backendResponse.Body = rp.Data
		}
	}

	return backendResponse,nil

}
