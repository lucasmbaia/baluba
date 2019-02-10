package core

/*import (
	"errors"
	"strconv"
	"net/http"
	"golang.org/x/net/http2"
)

const (
	defaultPort = 5200
)

type ServerH2 struct {
	server	    *http.Server
	certificate string
	key	    string
}

type ServerH2Config struct {
	Port	    int
	Certificate string
	Key	    string
}

func NewServerH2(config ServerH2Config) (ServerH2, error) {
	var (
		s   ServerH2
		err error
	)

	if config.Port == 0 {
		config.Port = defaultPort
	}

	if config.Certificate == "" {
		return s, errors.New("Certificate must be specified")
	}

	if config.Key == "" {
		return s, errors.New("Key must be specified")
	}

	s.server = &http.Server{
		Addr: fmt.Sprintf(":%s", strconv.Itoa(config.Port)),
	}

	http2.ConfigureServer(s.server, nil)
	http.HandleFunc("/upload", s.Upload)

	return s, nil
}

func (s *ServerH2) Listen() error {
	if err := s.server.ListenAndServeLTS(s.certificate, s.key); err != nil {
		return errors.Wrapf(err, "failed during server listen and serve")
	}

	return nil
}

func (s *ServerH2) Upload(w http.ResponseWriter, r *http.Request) {
	var (
		err	      error
		bytesReceived int64 = 0
		buf	      = new(bytes.Buffer)
	)

	if bytesReceived, err = io.Copy(buf, r.Body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "%+v", err)
		return
	}

	fmt.Println("bytes_received", bytes_received)
	return
}*/
