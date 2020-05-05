// Copyright 2018 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package security

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/coreos/ignition/v2/tests/register"
	"github.com/coreos/ignition/v2/tests/servers"
	"github.com/coreos/ignition/v2/tests/types"

	"github.com/vincent-petithory/dataurl"
)

func init() {
	cer, err := tls.X509KeyPair(publicKey, privateKey)
	if err != nil {
		panic(fmt.Sprintf("error loading x509 keypair: %v", err))
	}
	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	customCAServer.TLS = config
	customCAServer.StartTLS()

	cer2, err := tls.X509KeyPair(publicKey2, privateKey2)
	if err != nil {
		panic(fmt.Sprintf("error loading x509 keypair2: %v", err))
	}
	config2 := &tls.Config{Certificates: []tls.Certificate{cer2}}
	customCAServer2.TLS = config2
	customCAServer2.StartTLS()

	register.Register(register.PositiveTest, AppendConfigCustomCert())
	register.Register(register.PositiveTest, FetchFileCustomCert())
	register.Register(register.PositiveTest, FetchFileCABundleCert())
	register.Register(register.PositiveTest, FetchFileCustomCertHTTP())
	register.Register(register.PositiveTest, FetchFileCABundleCertHTTP())
	register.Register(register.PositiveTest, FetchFileCustomCertHTTPCompressed())
	register.Register(register.PositiveTest, FetchFileCustomCertHTTPUsingHeaders())
	register.Register(register.PositiveTest, FetchFileCustomCertHTTPUsingHeadersWithRedirect())
	register.Register(register.PositiveTest, FetchFileCustomCertHTTPUsingOverwrittenHeaders())
}

var (
	// generated via:
	// openssl ecparam -genkey -name secp384r1 -out server.key
	privateKey = []byte(`-----BEGIN EC PARAMETERS-----
BgUrgQQAIg==
-----END EC PARAMETERS-----
-----BEGIN EC PRIVATE KEY-----
MIGkAgEBBDB6yW6RIYfTXdYVuPY0V0L6EtZ6vZD86vgbsw52Y3/U5nZ2JE++JrKu
tt2Xt/NMzG6gBwYFK4EEACKhZANiAAQDEhfHEulYKlANw9eR5l455gwzAIQuraa0
49RhvM7PPywaiD8DobteQmE8wn7cJSzOYw6GLvrL4Q1BO5EFUXknkW50t8lfnUeH
veCNsqvm82F1NVevVoExAUhDYmMREa4=
-----END EC PRIVATE KEY-----`)

	// generated via:
	// openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
	publicKey = []byte(`-----BEGIN CERTIFICATE-----
MIICzTCCAlKgAwIBAgIJALTP0pfNBMzGMAoGCCqGSM49BAMCMIGZMQswCQYDVQQG
EwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNj
bzETMBEGA1UECgwKQ29yZU9TIEluYzEUMBIGA1UECwwLRW5naW5lZXJpbmcxEzAR
BgNVBAMMCmNvcmVvcy5jb20xHTAbBgkqhkiG9w0BCQEWDm9lbUBjb3Jlb3MuY29t
MB4XDTE4MDEyNTAwMDczOVoXDTI4MDEyMzAwMDczOVowgZkxCzAJBgNVBAYTAlVT
MRMwEQYDVQQIDApDYWxpZm9ybmlhMRYwFAYDVQQHDA1TYW4gRnJhbmNpc2NvMRMw
EQYDVQQKDApDb3JlT1MgSW5jMRQwEgYDVQQLDAtFbmdpbmVlcmluZzETMBEGA1UE
AwwKY29yZW9zLmNvbTEdMBsGCSqGSIb3DQEJARYOb2VtQGNvcmVvcy5jb20wdjAQ
BgcqhkjOPQIBBgUrgQQAIgNiAAQDEhfHEulYKlANw9eR5l455gwzAIQuraa049Rh
vM7PPywaiD8DobteQmE8wn7cJSzOYw6GLvrL4Q1BO5EFUXknkW50t8lfnUeHveCN
sqvm82F1NVevVoExAUhDYmMREa6jZDBiMA8GA1UdEQQIMAaHBH8AAAEwHQYDVR0O
BBYEFEbFy0SPiF1YXt+9T3Jig2rNmBtpMB8GA1UdIwQYMBaAFEbFy0SPiF1YXt+9
T3Jig2rNmBtpMA8GA1UdEwEB/wQFMAMBAf8wCgYIKoZIzj0EAwIDaQAwZgIxAOul
t3MhI02IONjTDusl2YuCxMgpy2uy0MPkEGUHnUOsxmPSG0gEBCNHyeKVeTaPUwIx
AKbyaAqbChEy9CvDgyv6qxTYU+eeBImLKS3PH2uW5etc/69V/sDojqpH3hEffsOt
9g==
-----END CERTIFICATE-----`)

	// generated via
	// openssl ecparam -genkey -name secp384r1 -out server.key
	privateKey2 = []byte(`-----BEGIN EC PARAMETERS-----
BgUrgQQAIg==
-----END EC PARAMETERS-----
-----BEGIN EC PRIVATE KEY-----
MIGkAgEBBDCfXncsl/kqihUWRHThBdGEDpv/bavwHYEi2tjrHiRkm+b7zhFlup8o
aP1l1zP1LhKgBwYFK4EEACKhZANiAAQ/J0D0C3h2a55JU3/EANe1d3e2/mfcoXGq
P8soiFdYntRIC4+V4dnRJuHRR+FHR/3531EIf2WXsoIJr/IRhR/j0tAeXpZ++G+E
vaooXf7gShnhRYKM4viPx4+DhSPjmqw=
-----END EC PRIVATE KEY-----`)

	// generate csr:
	// openssl req -new -key server.key -out server.csr
	// generate certificate:
	// openssl x509 -req -days 3650 -in server.csr -signkey server.key -out
	// server.crt -extensions v3_req -extfile extfile.conf
	// where extfile.conf has the following details:
	// $ cat extfile.conf
	// [ v3_req ]
	// subjectAltName = IP:127.0.0.1
	// subjectKeyIdentifier=hash
	// authorityKeyIdentifier=keyid
	// basicConstraints = critical,CA:TRUE
	publicKey2 = []byte(`-----BEGIN CERTIFICATE-----
MIICrDCCAjOgAwIBAgIUbFS1ugcEYYGQoTiV6O//r3wdO58wCgYIKoZIzj0EAwIw
gYQxCzAJBgNVBAYTAlVTMQswCQYDVQQIDAJOQzEQMA4GA1UEBwwHUmFsZWlnaDEQ
MA4GA1UECgwHUmVkIEhhdDEUMBIGA1UECwwLRW5naW5lZXJpbmcxDzANBgNVBAMM
BkNvcmVPUzEdMBsGCSqGSIb3DQEJARYOb2VtQGNvcmVvcy5jb20wHhcNMjAwNTA3
MjIzMzA3WhcNMzAwNTA1MjIzMzA3WjCBhDELMAkGA1UEBhMCVVMxCzAJBgNVBAgM
Ak5DMRAwDgYDVQQHDAdSYWxlaWdoMRAwDgYDVQQKDAdSZWQgSGF0MRQwEgYDVQQL
DAtFbmdpbmVlcmluZzEPMA0GA1UEAwwGQ29yZU9TMR0wGwYJKoZIhvcNAQkBFg5v
ZW1AY29yZW9zLmNvbTB2MBAGByqGSM49AgEGBSuBBAAiA2IABD8nQPQLeHZrnklT
f8QA17V3d7b+Z9yhcao/yyiIV1ie1EgLj5Xh2dEm4dFH4UdH/fnfUQh/ZZeyggmv
8hGFH+PS0B5eln74b4S9qihd/uBKGeFFgozi+I/Hj4OFI+OarKNkMGIwDwYDVR0R
BAgwBocEfwAAATAdBgNVHQ4EFgQUovVgWNFFPhrF7XzaRteDnpfPXxowHwYDVR0j
BBgwFoAUovVgWNFFPhrF7XzaRteDnpfPXxowDwYDVR0TAQH/BAUwAwEB/zAKBggq
hkjOPQQDAgNnADBkAjBvCIr9k43oR18Z4HLTzaRfzacFzo75Lt5n0pk3PA5CrUg3
sXU6o4IxyLNFHzIJn7cCMGTMVKEzoSZDclxkEgu53WM7PQljHgL9FJScEt4hzO2u
FFNjhq0ODV1LNc1i8pQCAg==
-----END CERTIFICATE-----`)
	// catting publicKey and publicKey2
	caBundle = []byte(`-----BEGIN CERTIFICATE-----
MIICzTCCAlKgAwIBAgIJALTP0pfNBMzGMAoGCCqGSM49BAMCMIGZMQswCQYDVQQG
EwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwNU2FuIEZyYW5jaXNj
bzETMBEGA1UECgwKQ29yZU9TIEluYzEUMBIGA1UECwwLRW5naW5lZXJpbmcxEzAR
BgNVBAMMCmNvcmVvcy5jb20xHTAbBgkqhkiG9w0BCQEWDm9lbUBjb3Jlb3MuY29t
MB4XDTE4MDEyNTAwMDczOVoXDTI4MDEyMzAwMDczOVowgZkxCzAJBgNVBAYTAlVT
MRMwEQYDVQQIDApDYWxpZm9ybmlhMRYwFAYDVQQHDA1TYW4gRnJhbmNpc2NvMRMw
EQYDVQQKDApDb3JlT1MgSW5jMRQwEgYDVQQLDAtFbmdpbmVlcmluZzETMBEGA1UE
AwwKY29yZW9zLmNvbTEdMBsGCSqGSIb3DQEJARYOb2VtQGNvcmVvcy5jb20wdjAQ
BgcqhkjOPQIBBgUrgQQAIgNiAAQDEhfHEulYKlANw9eR5l455gwzAIQuraa049Rh
vM7PPywaiD8DobteQmE8wn7cJSzOYw6GLvrL4Q1BO5EFUXknkW50t8lfnUeHveCN
sqvm82F1NVevVoExAUhDYmMREa6jZDBiMA8GA1UdEQQIMAaHBH8AAAEwHQYDVR0O
BBYEFEbFy0SPiF1YXt+9T3Jig2rNmBtpMB8GA1UdIwQYMBaAFEbFy0SPiF1YXt+9
T3Jig2rNmBtpMA8GA1UdEwEB/wQFMAMBAf8wCgYIKoZIzj0EAwIDaQAwZgIxAOul
t3MhI02IONjTDusl2YuCxMgpy2uy0MPkEGUHnUOsxmPSG0gEBCNHyeKVeTaPUwIx
AKbyaAqbChEy9CvDgyv6qxTYU+eeBImLKS3PH2uW5etc/69V/sDojqpH3hEffsOt
9g==
-----END CERTIFICATE-----
# CustomCAServer1 certificate
-----BEGIN CERTIFICATE-----
MIICrDCCAjOgAwIBAgIUbFS1ugcEYYGQoTiV6O//r3wdO58wCgYIKoZIzj0EAwIw
gYQxCzAJBgNVBAYTAlVTMQswCQYDVQQIDAJOQzEQMA4GA1UEBwwHUmFsZWlnaDEQ
MA4GA1UECgwHUmVkIEhhdDEUMBIGA1UECwwLRW5naW5lZXJpbmcxDzANBgNVBAMM
BkNvcmVPUzEdMBsGCSqGSIb3DQEJARYOb2VtQGNvcmVvcy5jb20wHhcNMjAwNTA3
MjIzMzA3WhcNMzAwNTA1MjIzMzA3WjCBhDELMAkGA1UEBhMCVVMxCzAJBgNVBAgM
Ak5DMRAwDgYDVQQHDAdSYWxlaWdoMRAwDgYDVQQKDAdSZWQgSGF0MRQwEgYDVQQL
DAtFbmdpbmVlcmluZzEPMA0GA1UEAwwGQ29yZU9TMR0wGwYJKoZIhvcNAQkBFg5v
ZW1AY29yZW9zLmNvbTB2MBAGByqGSM49AgEGBSuBBAAiA2IABD8nQPQLeHZrnklT
f8QA17V3d7b+Z9yhcao/yyiIV1ie1EgLj5Xh2dEm4dFH4UdH/fnfUQh/ZZeyggmv
8hGFH+PS0B5eln74b4S9qihd/uBKGeFFgozi+I/Hj4OFI+OarKNkMGIwDwYDVR0R
BAgwBocEfwAAATAdBgNVHQ4EFgQUovVgWNFFPhrF7XzaRteDnpfPXxowHwYDVR0j
BBgwFoAUovVgWNFFPhrF7XzaRteDnpfPXxowDwYDVR0TAQH/BAUwAwEB/zAKBggq
hkjOPQQDAgNnADBkAjBvCIr9k43oR18Z4HLTzaRfzacFzo75Lt5n0pk3PA5CrUg3
sXU6o4IxyLNFHzIJn7cCMGTMVKEzoSZDclxkEgu53WM7PQljHgL9FJScEt4hzO2u
FFNjhq0ODV1LNc1i8pQCAg==
-----END CERTIFICATE-----`)

	customCAServerFile = []byte(`{
			"ignition": { "version": "3.0.0" },
			"storage": {
				"files": [{
					"path": "/foo/bar",
					"contents": { "source": "data:,example%20file%0A" }
				}]
			}
		}`)

	customCAServerFile2 = []byte(`{
			"ignition": { "version": "3.0.0" },
			"storage": {
				"files": [{
					"path": "/foo/bar2",
					"contents": { "source": "data:,example%20file2%0A" }
				}]
			}
		}`)

	customCAServer = httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(customCAServerFile)
	}))
	customCAServer2 = httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(customCAServerFile2)
	}))
)

func AppendConfigCustomCert() types.Test {
	name := "tls.appendconfig"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	config := fmt.Sprintf(`{
		"ignition": {
			"version": "$version",
			"config": {
			  "merge": [{
				"source": %q
			  }]
			},
			"security": {
				"tls": {
					"certificateAuthorities": [{
						"source": %q
					}]
				}
			}
		}
	}`, customCAServer.URL, dataurl.EncodeBytes(publicKey))
	configMinVersion := "3.0.0"

	out[0].Partitions.AddFiles("ROOT", []types.File{
		{
			Node: types.Node{
				Name:      "bar",
				Directory: "foo",
			},
			Contents: "example file\n",
		},
	})

	return types.Test{
		Name:             name,
		In:               in,
		Out:              out,
		Config:           config,
		ConfigMinVersion: configMinVersion,
	}
}

func FetchFileCustomCert() types.Test {
	name := "tls.fetchfile"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	config := fmt.Sprintf(`{
		"ignition": {
			"version": "$version",
			"security": {
				"tls": {
					"certificateAuthorities": [{
						"source": %q
					}]
				}
			}
		},
		"storage": {
			"files": [{
				"path": "/foo/bar",
				"contents": {
					"source": %q
				}
			}]
		}
	}`, dataurl.EncodeBytes(publicKey), customCAServer.URL)
	configMinVersion := "3.0.0"

	out[0].Partitions.AddFiles("ROOT", []types.File{
		{
			Node: types.Node{
				Directory: "foo",
				Name:      "bar",
			},
			Contents: string(customCAServerFile),
		},
	})

	return types.Test{
		Name:             name,
		In:               in,
		Out:              out,
		Config:           config,
		ConfigMinVersion: configMinVersion,
	}
}

func FetchFileCABundleCert() types.Test {
	name := "tls.fetchfile.cabundle"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	config := fmt.Sprintf(`{
		"ignition": {
			"version": "$version",
			"security": {
				"tls": {
					"certificateAuthorities": [{
						"source": %q
					}]
				}
			}
		},
		"storage": {
			"files": [{
				"path": "/foo/bar",
				"contents": {
					"source": %q
				}
			},
			{
				"path": "/foo/bar2",
				"contents": {
					"source": %q
				}
			}]
		}
	}`, dataurl.EncodeBytes(caBundle), customCAServer.URL, customCAServer2.URL)
	configMinVersion := "3.0.0"

	out[0].Partitions.AddFiles("ROOT", []types.File{
		{
			Node: types.Node{
				Directory: "foo",
				Name:      "bar",
			},
			Contents: string(customCAServerFile),
		},
		{
			Node: types.Node{
				Directory: "foo",
				Name:      "bar2",
			},
			Contents: string(customCAServerFile2),
		},
	})

	return types.Test{
		Name:             name,
		In:               in,
		Out:              out,
		Config:           config,
		ConfigMinVersion: configMinVersion,
	}
}

func FetchFileCustomCertHTTP() types.Test {
	name := "tls.fetchfile.http"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	config := fmt.Sprintf(`{
		"ignition": {
			"version": "$version",
			"security": {
				"tls": {
					"certificateAuthorities": [{
						"source": "http://127.0.0.1:8080/certificates"
					}]
				}
			}
		},
		"storage": {
			"files": [{
				"path": "/foo/bar",
				"contents": {
					"source": %q
				}
			}]
		}
	}`, customCAServer.URL)
	configMinVersion := "3.0.0"

	out[0].Partitions.AddFiles("ROOT", []types.File{
		{
			Node: types.Node{
				Directory: "foo",
				Name:      "bar",
			},
			Contents: string(customCAServerFile),
		},
	})

	return types.Test{
		Name:             name,
		In:               in,
		Out:              out,
		Config:           config,
		ConfigMinVersion: configMinVersion,
	}
}

func FetchFileCABundleCertHTTP() types.Test {
	name := "tls.fetchfile.http.cabundle"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	config := fmt.Sprintf(`{
		"ignition": {
			"version": "$version",
			"security": {
				"tls": {
					"certificateAuthorities": [{
						"source": "http://127.0.0.1:8080/caBundle"
					}]
				}
			}
		},
		"storage": {
			"files": [{
				"path": "/foo/bar",
				"contents": {
					"source": %q
				}
			},
			{
				"path": "/foo/bar2",
				"contents": {
					"source": %q
				}
			}]
		}
	}`, customCAServer.URL, customCAServer2.URL)
	configMinVersion := "3.0.0"

	out[0].Partitions.AddFiles("ROOT", []types.File{
		{
			Node: types.Node{
				Directory: "foo",
				Name:      "bar",
			},
			Contents: string(customCAServerFile),
		},
		{
			Node: types.Node{
				Directory: "foo",
				Name:      "bar2",
			},
			Contents: string(customCAServerFile2),
		},
	})

	return types.Test{
		Name:             name,
		In:               in,
		Out:              out,
		Config:           config,
		ConfigMinVersion: configMinVersion,
	}
}

func FetchFileCustomCertHTTPCompressed() types.Test {
	name := "tls.fetchfile.http.compressed"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	config := fmt.Sprintf(`{
		"ignition": {
			"version": "$version",
			"security": {
				"tls": {
					"certificateAuthorities": [{
						"compression": "gzip",
						"source": "http://127.0.0.1:8080/certificates_compressed",
						"verification": {
							"hash": "sha512-%v"
						}
					}]
				}
			}
		},
		"storage": {
			"files": [{
				"path": "/foo/bar",
				"contents": {
					"source": %q
				}
			}]
		}
	}`, servers.PublicKeyHash, customCAServer.URL)
	configMinVersion := "3.1.0"

	out[0].Partitions.AddFiles("ROOT", []types.File{
		{
			Node: types.Node{
				Directory: "foo",
				Name:      "bar",
			},
			Contents: string(customCAServerFile),
		},
	})

	return types.Test{
		Name:             name,
		In:               in,
		Out:              out,
		Config:           config,
		ConfigMinVersion: configMinVersion,
	}
}

func FetchFileCustomCertHTTPUsingHeaders() types.Test {
	name := "tls.fetchfile.http.headers"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	config := fmt.Sprintf(`{
		"ignition": {
			"version": "$version",
			"security": {
				"tls": {
					"certificateAuthorities": [{
						"httpHeaders": [{"name": "X-Auth", "value": "r8ewap98gfh4d8"}, {"name": "Keep-Alive", "value": "300"}],
						"source": "http://127.0.0.1:8080/certificates_headers"
					}]
				}
			}
		},
		"storage": {
			"files": [{
				"path": "/foo/bar",
				"contents": {
					"source": %q
				}
			}]
		}
	}`, customCAServer.URL)
	configMinVersion := "3.1.0"

	out[0].Partitions.AddFiles("ROOT", []types.File{
		{
			Node: types.Node{
				Directory: "foo",
				Name:      "bar",
			},
			Contents: string(customCAServerFile),
		},
	})

	return types.Test{
		Name:             name,
		In:               in,
		Out:              out,
		Config:           config,
		ConfigMinVersion: configMinVersion,
	}
}

func FetchFileCustomCertHTTPUsingHeadersWithRedirect() types.Test {
	name := "tls.fetchfile.http.headers.redirect"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	config := fmt.Sprintf(`{
		"ignition": {
			"version": "$version",
			"security": {
				"tls": {
					"certificateAuthorities": [{
						"httpHeaders": [{"name": "X-Auth", "value": "r8ewap98gfh4d8"}, {"name": "Keep-Alive", "value": "300"}],
						"source": "http://127.0.0.1:8080/certificates_headers_redirect"
					}]
				}
			}
		},
		"storage": {
			"files": [{
				"path": "/foo/bar",
				"contents": {
					"source": %q
				}
			}]
		}
	}`, customCAServer.URL)
	configMinVersion := "3.1.0"

	out[0].Partitions.AddFiles("ROOT", []types.File{
		{
			Node: types.Node{
				Directory: "foo",
				Name:      "bar",
			},
			Contents: string(customCAServerFile),
		},
	})

	return types.Test{
		Name:             name,
		In:               in,
		Out:              out,
		Config:           config,
		ConfigMinVersion: configMinVersion,
	}
}

func FetchFileCustomCertHTTPUsingOverwrittenHeaders() types.Test {
	name := "tls.fetchfile.http.headers.overwrite"
	in := types.GetBaseDisk()
	out := types.GetBaseDisk()
	config := fmt.Sprintf(`{
		"ignition": {
			"version": "$version",
			"security": {
				"tls": {
					"certificateAuthorities": [{
						"httpHeaders": [
							{"name": "Keep-Alive", "value": "1000"},
							{"name": "Accept", "value": "application/json"},
							{"name": "Accept-Encoding", "value": "identity, compress"},
							{"name": "User-Agent", "value": "MyUA"}
						],
						"source": "http://127.0.0.1:8080/certificates_headers_overwrite"
					}]
				}
			}
		},
		"storage": {
			"files": [{
				"path": "/foo/bar",
				"contents": {
					"source": %q
				}
			}]
		}
	}`, customCAServer.URL)
	configMinVersion := "3.1.0"

	out[0].Partitions.AddFiles("ROOT", []types.File{
		{
			Node: types.Node{
				Directory: "foo",
				Name:      "bar",
			},
			Contents: string(customCAServerFile),
		},
	})

	return types.Test{
		Name:             name,
		In:               in,
		Out:              out,
		Config:           config,
		ConfigMinVersion: configMinVersion,
	}
}
