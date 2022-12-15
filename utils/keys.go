package keys

import (
	"bufio"
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
