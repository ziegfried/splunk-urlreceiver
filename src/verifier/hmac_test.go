package verifier

import "testing"
import "crypto/sha1"

func TestVerifyHmac(t *testing.T) {
	if verifyHmac([]byte("hello world"), "", "", sha1.New) {
		t.Error()
	}
	if !verifyHmac([]byte("hello world!"), "test123", "89e296a8504189d677de2eaebbf87babf137fd1a", sha1.New) {
		t.Error()
	}
}

func TestExtractHash(t *testing.T) {
	if h, str := extractHash("sha1=cafebabecafebabe"); h == nil || str != "cafebabecafebabe" {
		t.Error()
	}
	if h, str := extractHash("foobar"); h != nil || str != "" {
		t.Error()
	}
}
