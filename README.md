# Ginney

A library which build to be a Gin best friend. Ginney comes with the abilities to
- Override the default log of Gin and make it becomes  Application log format.
- Manage the correlation id of Composite service and Microservice
- Make a http request with the GET, POST and PUT methods. If the correlation id is in the context, ginney will automatically add it to the header of the request.

## Ginney in  Application project

```go
// main.go in Composite Service
package main

import (
	...
	"github.com/chaiyawatkit/ginney"
)

...

func main() {
	ginEngine := gin.New()

	ginEngine.Use(gin.Recovery())
	ginEngine.Use(ginney.LogWithCorrelationIdMiddleware(gin.DefaultWriter, []string{"/health"}))
	ginEngine.Use(ginney.CompositeCorrelationIdMiddleware())
	ginEngine.Use(ginney.FromGinContextToContextMiddleware())

	// Init Boredom Log
	boredom.Init()
    
	// Init usecase and repository
	...

}
```

```go
// /app/modules/wallet_service/repository/microservice_repository.go in Composite Service
package repository

import ...

func (repository *repository) SugarDaddySign(ctx context.Context, xdrOps []string) (*string, error) {
	errContextMsg := "fail to sign sugar daddy xdr ops by requesting to stellar-account-service " + env.StellarServiceUrl

	preSignRequestJSON := struct {
		XdrOps []string `json:"xdrOps"`
	}{xdrOps}
	requestBody, _ := json.Marshal(preSignRequestJSON)

	url := fmt.Sprintf("%s/v1/sugar-daddies.sign", env.StellarServiceUrl)
	resp, _ := fin.Post(ctx, url, constants.ContentTypeJSON, bytes.NewBuffer(requestBody))

	responseBodyData, err := utils.VeloResponseParser(resp, errContextMsg)
	if err != nil {
		return nil, err
	}

	responseModel, err := new(models.SugarDaddyXdrResponse).Parse(responseBodyData)
	if err != nil {
		return nil, errors.Wrap(err, errContextMsg)
	}

	return responseModel.SignedXdr, nil
}
```

```go
// main.go in Microservice Service
package main

import (
	...
	"github.com/chaiyawatkit/ginney"
)

...

func main() {
	ginEngine := gin.New()

	ginEngine.Use(gin.Recovery())
	ginEngine.Use(ginney.LogWithCorrelationIdMiddleware(gin.DefaultWriter, []string{"/health"}))
	ginEngine.Use(ginney.MicroServiceCorrelationIdMiddleware())
	ginEngine.Use(ginney.FromGinContextToContextMiddleware())

	// Init Boredom Log
	boredom.Init()
    
	// Init usecase and repository
	...

}
```

## Available API
```go
// converting from gin context to context by placing the gin.Context to a very specific key of context
ginney.FromGinContextToContext(gc)
// converting from context to gin context by getting the gin.Context from a very specific key of context
ginney.FromContextToGinContext(ctx)
// a very specific key of context which hold gin.Context
ginney.GinContextKey
// a key of correlation id in the header
ginney.CorrelationIdHeaderKey
```