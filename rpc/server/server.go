package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc/v2"
	rpcjson "github.com/gorilla/rpc/v2/json2"

	"github.com/anyswap/CrossChain-Bridge/cmd/utils"
	"github.com/anyswap/CrossChain-Bridge/log"
	"github.com/anyswap/CrossChain-Bridge/params"
	"github.com/anyswap/CrossChain-Bridge/rpc/restapi"
	"github.com/anyswap/CrossChain-Bridge/rpc/rpcapi"
)

// StartAPIServer start api server
func StartAPIServer() {
	router := initRouter()

	apiPort := params.GetAPIPort()
	apiServer := params.GetConfig().APIServer
	allowedOrigins := apiServer.AllowedOrigins

	corsOptions := []handlers.CORSOption{
		handlers.AllowedMethods([]string{"GET", "POST"}),
	}
	if len(allowedOrigins) != 0 {
		corsOptions = append(corsOptions,
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"}),
			handlers.AllowedOrigins(allowedOrigins),
		)
	}

	log.Info("JSON RPC service listen and serving", "port", apiPort, "allowedOrigins", allowedOrigins)
	svr := http.Server{
		Addr:         fmt.Sprintf(":%v", apiPort),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		Handler:      handlers.CORS(corsOptions...)(router),
	}
	go func() {
		if err := svr.ListenAndServe(); err != nil {
			log.Error("ListenAndServe error", "err", err)
		}
	}()

	utils.TopWaitGroup.Add(1)
	go utils.WaitAndCleanup(func() { doCleanup(&svr) })
}

func doCleanup(svr *http.Server) {
	defer utils.TopWaitGroup.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := svr.Shutdown(ctx); err != nil {
		log.Error("Server Shutdown failed", "err", err)
	}
	log.Info("Close http server success")
}

// nolint:funlen // put together handle func
func initRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	rpcserver := rpc.NewServer()
	rpcserver.RegisterCodec(rpcjson.NewCodec(), "application/json")
	_ = rpcserver.RegisterService(new(rpcapi.RPCAPI), "swap")

	r.Handle("/rpc", rpcserver)
	r.HandleFunc("/serverinfo", restapi.ServerInfoHandler).Methods("GET")
	r.HandleFunc("/versioninfo", restapi.VersionInfoHandler).Methods("GET")
	r.HandleFunc("/nonceinfo", restapi.NonceInfoHandler).Methods("GET")
	r.HandleFunc("/pairinfo/{pairid}", restapi.TokenPairInfoHandler).Methods("GET")
	r.HandleFunc("/statistics/{pairid}", restapi.StatisticsHandler).Methods("GET")
	r.HandleFunc("/swapin/post/{pairid}/{txid}", restapi.PostSwapinHandler).Methods("POST")
	r.HandleFunc("/swapout/post/{pairid}/{txid}", restapi.PostSwapoutHandler).Methods("POST")
	r.HandleFunc("/swapin/p2sh/{txid}/{bind}", restapi.PostP2shSwapinHandler).Methods("POST")
	r.HandleFunc("/swapin/retry/{pairid}/{txid}", restapi.RetrySwapinHandler).Methods("POST")
	r.HandleFunc("/swapin/{pairid}/{txid}", restapi.GetSwapinHandler).Methods("GET")
	r.HandleFunc("/swapout/{pairid}/{txid}", restapi.GetSwapoutHandler).Methods("GET")
	r.HandleFunc("/swapin/{pairid}/{txid}/raw", restapi.GetRawSwapinHandler).Methods("GET")
	r.HandleFunc("/swapout/{pairid}/{txid}/raw", restapi.GetRawSwapoutHandler).Methods("GET")
	r.HandleFunc("/swapin/{pairid}/{txid}/rawresult", restapi.GetRawSwapinResultHandler).Methods("GET")
	r.HandleFunc("/swapout/{pairid}/{txid}/rawresult", restapi.GetRawSwapoutResultHandler).Methods("GET")
	r.HandleFunc("/swapin/history/{pairid}/{address}", restapi.SwapinHistoryHandler).Methods("GET")
	r.HandleFunc("/swapout/history/{pairid}/{address}", restapi.SwapoutHistoryHandler).Methods("GET")
	r.HandleFunc("/p2sh/{address}", restapi.GetP2shAddressInfo).Methods("GET", "POST")
	r.HandleFunc("/p2sh/bind/{address}", restapi.RegisterP2shAddress).Methods("GET", "POST")
	r.HandleFunc("/registered/{address}", restapi.GetRegisteredAddress).Methods("GET", "POST")
	r.HandleFunc("/register/{address}", restapi.RegisterAddress).Methods("GET", "POST")

	methodsExcluesGet := []string{"POST", "HEAD", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}
	methodsExcluesPost := []string{"GET", "HEAD", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}
	methodsExcluesGetAndPost := []string{"HEAD", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}

	r.HandleFunc("/serverinfo", warnHandler).Methods(methodsExcluesGet...)
	r.HandleFunc("/versioninfo", warnHandler).Methods(methodsExcluesGet...)
	r.HandleFunc("/nonceinfo", warnHandler).Methods(methodsExcluesGet...)
	r.HandleFunc("/pairinfo/{pairid}", warnHandler).Methods(methodsExcluesGet...)
	r.HandleFunc("/statistics/{pairid}", warnHandler).Methods(methodsExcluesGet...)
	r.HandleFunc("/swapin/post/{pairid}/{txid}", warnHandler).Methods(methodsExcluesPost...)
	r.HandleFunc("/swapout/post/{pairid}/{txid}", warnHandler).Methods(methodsExcluesPost...)
	r.HandleFunc("/swapin/p2sh/{txid}/{bind}", warnHandler).Methods(methodsExcluesPost...)
	r.HandleFunc("/swapin/retry/{pairid}/{txid}", warnHandler).Methods(methodsExcluesPost...)
	r.HandleFunc("/swapin/{pairid}/{txid}", warnHandler).Methods(methodsExcluesGet...)
	r.HandleFunc("/swapout/{pairid}/{txid}", warnHandler).Methods(methodsExcluesGet...)
	r.HandleFunc("/swapin/{pairid}/{txid}/raw", warnHandler).Methods(methodsExcluesGet...)
	r.HandleFunc("/swapout/{pairid}/{txid}/raw", warnHandler).Methods(methodsExcluesGet...)
	r.HandleFunc("/swapin/{pairid}/{txid}/rawresult", warnHandler).Methods(methodsExcluesGet...)
	r.HandleFunc("/swapout/{pairid}/{txid}/rawresult", warnHandler).Methods(methodsExcluesGet...)
	r.HandleFunc("/swapin/history/{pairid}/{address}", warnHandler).Methods(methodsExcluesGet...)
	r.HandleFunc("/swapout/history/{pairid}/{address}", warnHandler).Methods(methodsExcluesGet...)
	r.HandleFunc("/p2sh/{address}", warnHandler).Methods(methodsExcluesGetAndPost...)
	r.HandleFunc("/p2sh/bind/{address}", warnHandler).Methods(methodsExcluesGetAndPost...)
	r.HandleFunc("/registered/{address}", warnHandler).Methods(methodsExcluesGetAndPost...)
	r.HandleFunc("/register/{address}", warnHandler).Methods(methodsExcluesGetAndPost...)

	return r
}

func warnHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Forbid '%v' on '%v'\n", r.Method, r.RequestURI)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
