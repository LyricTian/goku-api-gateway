package backend

import (
	"context"
	"github.com/eolinker/goku-api-gateway/config"
	"github.com/eolinker/goku-api-gateway/goku-node/common"
	"github.com/eolinker/goku-api-gateway/goku-service/application"
	"github.com/eolinker/goku-api-gateway/goku-service/balance"
	"github.com/eolinker/goku-api-gateway/node/gateway/application/action"
	"github.com/eolinker/goku-api-gateway/node/gateway/application/interpreter"
	"github.com/eolinker/goku-api-gateway/node/gateway/response"
	"io/ioutil"
	"strings"
	"time"
)

type Layer struct {
	BalanceName string
	Balance application.IHttpApplication
	HasBalance bool
	Protocol  string


	Filter action.Filter
	Method  string
	Path    interpreter.Interpreter
	Decode  response.DecodeHandle

	Body interpreter.Interpreter
	Encode string
	Target  string
	Group[]   string
	Retry   int
	TimeOut time.Duration

}

func (b *Layer) Send(ctx *common.Context,variables *interpreter.Variables,deadline context.Context) (*BackendResponse, error) {
	path:= b.Path.Execution(variables)
	body:= b.Body.Execution(variables)
	method:= b.Method

	r, finalTargetServer, retryTargetServers, err := b.Balance.Send(b.Protocol,method, path, ctx.ProxyRequest.Querys(), ctx.ProxyRequest.Headers(),[]byte(body), b.TimeOut, b.Retry)

	if err!=nil{
		return nil,err
	}
	backendResponse := &BackendResponse{
		Method:            strings.ToUpper(method),
		Protocol:          b.Protocol,
		//Response:           r,
		TargetUrl:          path,
		FinalTargetServer:  finalTargetServer,
		RetryTargetServers: retryTargetServers,
		Header:r.Header,
	}

	defer r.Body.Close()
	backendResponse.BodyOrg, err = ioutil.ReadAll(r.Body)
	if err!= nil{
		return backendResponse,nil
	}


	rp,e:=response.Decode(backendResponse.BodyOrg,b.Decode)
	if e!= nil{
		backendResponse.Body = nil
		return nil,e
	}



	b.Filter.Do(rp)

	if b.Target!= ""{
		rp.ReTarget(b.Target)
	}
	if len(b.Group)>0{
		rp.Group(b.Group)
	}
	backendResponse.Body = rp.Data
	return backendResponse,nil
}



func NewLayer(step *config.APIStepConfig) *Layer {
	var b = &Layer{
		BalanceName: step.Balance,
		Balance:     nil,
		HasBalance:  false,
		Protocol:    step.Proto,
		Filter:      genFilter(step.BlackList,step.WhiteList,step.Actions),
		Method:      strings.ToUpper(step.Method),
		Path:       interpreter.GenPath( step.Path),
		Decode:      response.GetDecoder(step.Decode),
		Encode:      step.Encode,
		Target:      step.Target,
		Group:      nil,
		TimeOut:time.Duration(step.TimeOut)*time.Millisecond,
		Body:interpreter.Gen(step.Body),
		Retry:     step.Retry,
	}
	if step.Group != ""{
		b.Group =  strings.Split(step.Group,".")
	}

	b.Balance, b.HasBalance = balance.GetByName(b.BalanceName)

	return b
}
