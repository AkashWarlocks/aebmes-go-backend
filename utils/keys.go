package keys

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

var ownerPrivateKey *rsa.PrivateKey
var ownerPublicKey *rsa.PublicKey
var userPrivateKey *rsa.PrivateKey
var userPublicKey *rsa.PublicKey

func ReadUserKeys() (*rsa.PrivateKey, *rsa.PublicKey, *rsa.PrivateKey, *rsa.PublicKey, error) {
	ownerprivateKeyFile, err := os.Open("owner_private_key.pem")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ownerpemfileinfo, _ := ownerprivateKeyFile.Stat()
	var size int64 = ownerpemfileinfo.Size()
	pembytes := make([]byte, size)
	buffer := bufio.NewReader(ownerprivateKeyFile)
	_, err = buffer.Read(pembytes)
	data, _ := pem.Decode([]byte(pembytes))
	ownerprivateKeyFile.Close()

	ownerPrivateKey, err = x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ownerPublicKey = &ownerPrivateKey.PublicKey

	userprivateKeyFile, err := os.Open("user_private_key.pem")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	userpemfileinfo, _ := userprivateKeyFile.Stat()
	var usersize int64 = userpemfileinfo.Size()
	userpembytes := make([]byte, usersize)
	userbuffer := bufio.NewReader(userprivateKeyFile)
	_, err = userbuffer.Read(userpembytes)
	userdata, _ := pem.Decode([]byte(pembytes))
	ownerprivateKeyFile.Close()

	userPrivateKey, err = x509.ParsePKCS1PrivateKey(userdata.Bytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	userPublicKey = &userPrivateKey.PublicKey

	return ownerPrivateKey, ownerPublicKey, userPrivateKey, userPublicKey, nil

}

func createUserKeys(userType string) {

	userPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err.Error)
		os.Exit(1)
	}
	// userPublicKey := &userPrivateKey.PublicKey

	privKeyBytes := x509.MarshalPKCS1PrivateKey(userPrivateKey)
	fmt.Println(privKeyBytes)
	pemPrivateFile, err := os.Create(userType + "_private_key.pem")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = pem.Encode(pemPrivateFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privKeyBytes,
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// fmt.Println("Private Key :", string(priv_pem), "end")
	// fmt.Println("Public key ", userPublicKey)
	pemPrivateFile.Close()
}

func OwnerPrivateKey() *rsa.PrivateKey {
	return ownerPrivateKey
}

func OwnerPublicKey() *rsa.PublicKey {
	return ownerPublicKey
}

func UserPrivateKey() *rsa.PrivateKey {
	return userPrivateKey
}

func UserPublicKey() *rsa.PublicKey {
	return userPublicKey
}
