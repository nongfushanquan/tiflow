// Copyright 2020 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package security

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/pingcap/tidb-tools/pkg/utils"
	cerror "github.com/pingcap/tiflow/pkg/errors"
	pd "github.com/tikv/pd/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Credential holds necessary path parameter to build a tls.Config
type Credential struct {
	CAPath        string   `toml:"ca-path" json:"ca-path"`
	CertPath      string   `toml:"cert-path" json:"cert-path"`
	KeyPath       string   `toml:"key-path" json:"key-path"`
	CertAllowedCN []string `toml:"cert-allowed-cn" json:"cert-allowed-cn"`
}

// IsTLSEnabled checks whether TLS is enabled or not.
func (s *Credential) IsTLSEnabled() bool {
	return len(s.CAPath) != 0
}

// PDSecurityOption creates a new pd SecurityOption from Security
func (s *Credential) PDSecurityOption() pd.SecurityOption {
	return pd.SecurityOption{
		CAPath:   s.CAPath,
		CertPath: s.CertPath,
		KeyPath:  s.KeyPath,
	}
}

// ToGRPCDialOption constructs a gRPC dial option.
func (s *Credential) ToGRPCDialOption() (grpc.DialOption, error) {
	tlsCfg, err := s.ToTLSConfig()
	if err != nil || tlsCfg == nil {
		return grpc.WithInsecure(), err
	}
	return grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg)), nil
}

// ToTLSConfig generates tls's config from *Security
func (s *Credential) ToTLSConfig() (*tls.Config, error) {
	cfg, err := utils.ToTLSConfig(s.CAPath, s.CertPath, s.KeyPath)
	return cfg, cerror.WrapError(cerror.ErrToTLSConfigFailed, err)
}

// ToTLSConfigWithVerify generates tls's config from *Security and requires
// the remote common name to be verified.
func (s *Credential) ToTLSConfigWithVerify() (*tls.Config, error) {
	cfg, err := utils.ToTLSConfigWithVerify(s.CAPath, s.CertPath, s.KeyPath, s.CertAllowedCN)
	return cfg, cerror.WrapError(cerror.ErrToTLSConfigFailed, err)
}

func (s *Credential) getSelfCommonName() (string, error) {
	if s.CertPath == "" {
		return "", nil
	}
	data, err := os.ReadFile(s.CertPath)
	if err != nil {
		return "", cerror.WrapError(cerror.ErrToTLSConfigFailed, err)
	}
	block, _ := pem.Decode(data)
	if block == nil || block.Type != "CERTIFICATE" {
		return "", cerror.ErrToTLSConfigFailed.GenWithStack("failed to decode PEM block to certificate")
	}
	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", cerror.WrapError(cerror.ErrToTLSConfigFailed, err)
	}
	return certificate.Subject.CommonName, nil
}

// AddSelfCommonName add Common Name in certificate that specified by s.CertPath
// to s.CertAllowedCN
func (s *Credential) AddSelfCommonName() error {
	cn, err := s.getSelfCommonName()
	if err != nil {
		return err
	}
	if cn == "" {
		return nil
	}
	s.CertAllowedCN = append(s.CertAllowedCN, cn)
	return nil
}
