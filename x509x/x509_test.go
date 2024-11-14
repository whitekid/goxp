package x509x

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseCSR(t *testing.T) {
	type args struct {
		commonName         string
		dnsNames           []string
		signatureAlgorithm x509.SignatureAlgorithm
		organization       string
		organizationUnit   string
		country            string
		province           string
		locality           string
	}
	tests := [...]struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{"cn", []string{"example.local"}, x509.ECDSAWithSHA384, "O", "OU", "KR", "ST", "LOC"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, pemBytes, err := CreateCertificateRequest(&x509.CertificateRequest{
				DNSNames:           tt.args.dnsNames,
				SignatureAlgorithm: tt.args.signatureAlgorithm,
				Subject: pkix.Name{
					CommonName:         tt.args.commonName,
					Organization:       []string{tt.args.organization},
					OrganizationalUnit: []string{tt.args.organizationUnit},
					Country:            []string{tt.args.country},
					Province:           []string{tt.args.province},
					Locality:           []string{tt.args.locality},
				},
			})
			require.NoError(t, err)

			got, err := ParseCSR(pemBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCSRPEM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, tt.args.commonName, got.Subject.CommonName)
			require.Equal(t, tt.args.dnsNames, got.DNSNames)
			require.Equal(t, tt.args.signatureAlgorithm, got.SignatureAlgorithm)
			require.Equal(t, tt.args.organization, got.Subject.Organization[0])
			require.Equal(t, tt.args.organizationUnit, got.Subject.OrganizationalUnit[0])
			require.Equal(t, tt.args.country, got.Subject.Country[0])
			require.Equal(t, tt.args.province, got.Subject.Province[0])
			require.Equal(t, tt.args.locality, got.Subject.Locality[0])
		})
	}
}

func TestParseRevocationList(t *testing.T) {
	type args struct {
		crlBytes []byte
	}
	tests := [...]struct {
		name       string
		args       args
		wantErr    bool
		wantOrg    string
		wantNumber *big.Int
	}{
		{"valid", args{MustReadFile("fixtures/DigiCertTLSHybridECCSHA3842020CA1-1.crl")}, false, "DigiCert Inc", big.NewInt(499)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := x509.ParseRevocationList(tt.args.crlBytes)
			if (err != nil) != tt.wantErr {
				require.Failf(t, "ParseRevocationList() failed", "error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.wantOrg, got.Issuer.Organization[0])
			require.Equal(t, tt.wantNumber, got.Number)
		})
	}
}

func TestGenerateKey(t *testing.T) {
	type args struct {
		algorithm x509.SignatureAlgorithm
	}
	tests := [...]struct {
		name    string
		args    args
		wantErr bool
	}{
		{"SHA256WithRSA", args{algorithm: x509.SHA256WithRSA}, false},
		{"ECDSAWithSHA256", args{algorithm: x509.ECDSAWithSHA256}, false},
		{"PureEd25519", args{algorithm: x509.PureEd25519}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateKey(tt.args.algorithm)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_, err = EncodePrivateKeyToPEM(got)
			require.NoError(t, err)

			pemBytes, err := x509.MarshalPKIXPublicKey(got.Public())
			require.NoError(t, err)
			require.NotNil(t, pemBytes)
		})
	}
}

func TestPrivateKeyAlgorithm(t *testing.T) {
	type args struct {
		algo x509.SignatureAlgorithm
	}
	tests := [...]struct {
		name string
		args args
	}{
		{"valid", args{x509.ECDSAWithSHA256}},
		{"valid", args{x509.ECDSAWithSHA384}},
		{"valid", args{x509.ECDSAWithSHA512}},
		{"valid", args{x509.SHA256WithRSA}},
		{"valid", args{x509.SHA384WithRSA}},
		{"valid", args{x509.SHA512WithRSA}},
		{"valid", args{x509.PureEd25519}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := GenerateKey(tt.args.algo)
			require.NoError(t, err)
			got := PrivateKeyAlgorithm(key)
			require.Equalf(t, tt.args.algo, got, "want %s but got %s", tt.args.algo.String(), got.String())
		})
	}
}

func MustReadFile(name string) []byte {
	data, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}

	return data
}
