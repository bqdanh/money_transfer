package jwt

import (
	"fmt"
	"testing"
)

func TestJwtRSADecoder_GenerateToken(t *testing.T) {
	privateKey := GetPrivateKey("./cert/id_rsa_test")
	jwtRSADecoder := JwtRSADecoder{
		privateKey: privateKey,
		publicKey:  nil,
	}

	m := map[string]interface{}{
		"username": "admin",
		"exp":      10000000000,
	}

	token, err := jwtRSADecoder.GenerateToken(m)
	if err != nil {
		t.Errorf("GenerateToken() error = %v", err)
		return
	}
	fmt.Println("token: ", token)
}

func TestJwtRSADecoder_ValidateToken(t *testing.T) {
	jwtToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjEwMDAwMDAwMDAwLCJ1c2VybmFtZSI6ImFkbWluIn0.FXb7EKu7ti5B1best-E6WfqCvj8lW_2Arnzmb_D2xC17zxIsoM0TFGqntTQ7Mem2Bems1raZpnEcrnySOWe4ghH7SFo4KV2-MInEqYMQTI-U-ITy5u-lrqRubKFqpuJ0u2g4e8kLX-n6DCixWtO4qFPXmXITAabSfx0JUYvvUwiW6H6sy9cy0mOIhb0BS4gyAXC6voy9AWS-mLAsGy7oCHXweBDb4wq48dvg0YjSPPjgd7KtFJJs6NqZ-TMz8aGjpNbhlY2gEKb5T3dyo-oqieJjuPqfDfjiZi91b834WnAlTDHH4MMvXFLltMPDGFOJmDcTlTGKouKvseGERLr_yJEQGCwhB7NFND2KBhqR1nJi9r4axZ21LOmNJhdCzQ-yKlj37zCy_VuMpP_oqV61nE8Ei2S6R2j74WIghDoCUyDNOD4ZdX_aaPKx_YUcmjZMooTvxX_g9FI7DnaipTYaqttj1iUl2D3ShNH_lPcL4g8rv6CqMDn535ZTp0SU4APgBWyZBjkhFwngBgou0pJkV4wypBpFunPU6UPN-vAst8_-acboCKOXiBhwwwitqIP1FtEaJ6DBMmb8KaP9UpNW4WWApkxpp0kuzVqLIDH11QriLrWVVCFpBnwfKeD93TGL5zF1XPAFp1hscOKKawnhvx2HI-lP_LkI2ycgoMJuHlM\n"
	publicKey := GetPublicKey("./cert/id_rsa_test.pub")
	jwtRSADecoder := JwtRSADecoder{
		privateKey: nil,
		publicKey:  publicKey,
	}
	data, isVerify, err := jwtRSADecoder.ValidateToken(jwtToken)
	if err != nil {
		t.Errorf("GenerateToken() error = %v", err)
		return
	}
	if !isVerify {
		t.Errorf("GenerateToken() not verify")
		return
	}

	fmt.Println("verify: ", data)
}
