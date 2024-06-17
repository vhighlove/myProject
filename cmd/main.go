package main

import (
	"context"
	"fmt"
	"net/http"
)

// -http методы, url
// -443, 80
// -какие бывают body (postman)
// -headers, headercontenttype
// -Context в go
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/road/ty145", GetMiddleware(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "key1", "qwerty")
		ToPrintContext(ctx)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`Hello World GetMethod`))

	}))
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
	// server := new(myproject.Server)
	// if err := server.Run("8080"); err != nil {
	// 	log.Fatalf("error with running server: %s", err.Error())
	// }
}

func GetMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`not allowed method`))
			return
		}
		ctx := context.WithValue(r.Context(), "key", "Hello World")
		r = r.WithContext(ctx)
		next(w, r)
	}
}

func ToPrintContext(ctx context.Context) {
	fmt.Println(ctx.Value("key1"))
	fmt.Println(ctx.Value("key"))
	fmt.Println(ctx.Value("ke"))
}
