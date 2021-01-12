package main
import(
  // "fmt"
  // "net/http"
  "log"
  "html"
  "context"
  "crypto/tls"
  "flag"
  "os"
  "os/signal"
  "syscall"
  // "github.com/golang/glog"


  "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
	"k8s.io/api/admission/v1beta1"
	"k8s.io/api/core/v1"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	port = "8080"
)

var (
	tlscert, tlskey string
)

type ValidationAdmissionHandler struct {
}

func processRequests(w http.ResponseWriter, r *http.Request){
    curr_url := html.EscapeString(r.URL.Path)
    fmt.Println(curr_url)

    switch curr_url {
    case "/":
      fmt.Fprintf(w, "Welcome to the HomePage!")
      fmt.Println("Endpoint Hit: homePage")
    case "/validate":
      serveK8s()
    default:
      fmt.Fprintf(w, "Sorry, no junk in the server :| !")
    }
}

func handleRequests() {
    http.HandleFunc("/", processRequests)
    http.HandleFunc("/validate", processRequests)
    log.Fatal(http.ListenAndServe(":10000", nil))
}

func serveK8s() {
  flag.StringVar(&tlscert, "tlsCertFile", "certs/gandalf-crt.pem", "File containing the x509 Certificate for HTTPS.")
	flag.StringVar(&tlskey, "tlsKeyFile", "certs/gandalf-key.pem", "File containing the x509 private key to --tlsCertFile.")

	flag.Parse()

	certs, err := tls.LoadX509KeyPair(tlscert, tlskey)
	if err != nil {
		glog.Errorf("Filed to load key pair: %v", err)
	}

	server := &http.Server{
		Addr:      fmt.Sprintf(":%v", port),
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{certs}},
	}

	// define http server and server handler
	gs := ValidationAdmissionHandler{}
	mux := http.NewServeMux()
	mux.HandleFunc("/validate", gs.serve)
	server.Handler = mux

	// start webhook server in new rountine
	go func() {
		if err := server.ListenAndServeTLS("", ""); err != nil {
			glog.Errorf("Failed to listen and serve webhook server: %v", err)
		}
	}()

	glog.Infof("Server running listening in port: %s", port)

	// listening shutdown singal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	glog.Info("Got shutdown signal, shutting down webhook server gracefully...")
	server.Shutdown(context.Background())
}

func (gs *ValidationAdmissionHandler) serve(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	if len(body) == 0 {
		glog.Error("empty body")
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}
	glog.Info("Received request")

	if r.URL.Path != "/validate" {
		glog.Error("no validate")
		http.Error(w, "no validate", http.StatusBadRequest)
		return
	}

	arRequest := v1beta1.AdmissionReview{}
	if err := json.Unmarshal(body, &arRequest); err != nil {
		glog.Error("incorrect body")
		http.Error(w, "incorrect body", http.StatusBadRequest)
	}

	raw := arRequest.Request.Object.Raw
	pod := v1.Pod{}
	if err := json.Unmarshal(raw, &pod); err != nil {
		glog.Error("error deserializing pod")
		return
	}
	if pod.Name == "smooth-app" {
		return
	}

	arResponse := v1beta1.AdmissionReview{
		Response: &v1beta1.AdmissionResponse{
			Allowed: false,
			Result: &metav1.Status{
				Message: "Keep calm and not add more crap in the cluster!",
			},
		},
	}
	resp, err := json.Marshal(arResponse)
	if err != nil {
		glog.Errorf("Can't encode response: %v", err)
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
	}
	glog.Infof("Ready to write reponse ...")
	if _, err := w.Write(resp); err != nil {
		glog.Errorf("Can't write response: %v", err)
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
}

func main() {
    // handleRequests()
    serveK8s()
}
