//Package mitm is the core of this project and is responsible for creating a proxy that
//intercepts all the HTTP/HTTPS traffic going through it.
//The SSL bumping is currently made by using a fake CA whose keys are stored in
//path.Join(os.Getenv("HOME"), ".mitm")
//This is a fork of github.com/kr/mitm
package mitm

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"os"
	"path"

	"github.com/empijei/wapty/config"
)

var (
	localhostname, _ = os.Hostname()
	keyFile          = path.Join(config.ConfDir, "ca-key.pem")
	certFile         = path.Join(config.ConfDir, "ca-cert.crt")
)

// LoadCA loads the ca from "HOME/dir"
func LoadCA() (cert tls.Certificate, err error) {
	// TODO(kr): check file permissions
	cert, err = tls.LoadX509KeyPair(certFile, keyFile)
	if os.IsNotExist(err) {
		cert, err = genCA()
	}
	if err == nil {
		cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	}
	return
}

func genCA() (cert tls.Certificate, err error) {
	certPEM, keyPEM, err := GenerateCA(localhostname)
	if err != nil {
		return
	}
	cert, _ = tls.X509KeyPair(certPEM, keyPEM)
	err = ioutil.WriteFile(certFile, certPEM, 0400)
	if err == nil {
		err = ioutil.WriteFile(keyFile, keyPEM, 0400)
	}
	return cert, err
}
