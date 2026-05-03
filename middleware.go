package main
import(
	"net/http"
	"fmt"
)

func middlewareLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Println(req.Method, req.URL.Path)

		next.ServeHTTP(w, req)
	})
}