package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"hash"
	"io"

	"io/ioutil"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/warlockz/ase-service/db"
	routes "github.com/warlockz/ase-service/router"
	keys "github.com/warlockz/ase-service/utils"
	"github.com/warlockz/ase-service/utils/hedera"
)

func main() {
	router := gin.Default()
	fmt.Println("Hello world")
	db.ConnectDB()
	keys.ReadUserKeys()
	hedera.ConnectHedera()
	routes.RouteIndex(router)

	router.Run("localhost:3000")

	// createUserKeys("owner")
	// createUserKeys("user")
	/**
	1. Upload File
	2. Extract Keywords
	3. Encrypt File
	4.
	*/
	// ownerPrivateKey, ownerPublicKey, userPrivateKey, userPublicKey, err := readUserKeys()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// //	fmt.Println(ownerPrivateKey, ownerPublicKey, userPrivateKey, userPublicKey)
	// // 1. Upload File Along with keywords

	// //content, err := ioutil.ReadFile("thermopylae.txt")
	// key := make([]byte, 32)

	// keywordArray := [3]string{"sample", "Thermopylae", "Greek"}
	// keywordHashes := generateKeywordHash(keywordArray[:])

	// b, errs := json.Marshal(keywordHashes)
	// if errs != nil {
	// 	log.Fatal(errs)
	// }
	// fmt.Println(string(b))

	// fileHash, filePath, fileName := encryptFile(key)

	// // finalFileName :=
	// //fmt.Println(string(fileHash[:]), string(ciphertext), filePath)

	// label := []byte("")
	// hash := sha256.New()

	// // Encrypt the key
	//keyCiphertext, err := encryptKey(key, userPublicKey, ownerPrivateKey, hash, label)

	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// KeywordArray := types.ParametersType{
	// 	Datatype: "stringArray",
	// 	Array:    keywordArray[:],
	// }

	// FileHash := types.ParametersType{
	// 	Datatype: "string",
	// 	Value:    hex.EncodeToString(fileHash[:]),
	// }
	// fmt.Println(hex.EncodeToString((keyCiphertext[:])))
	// KeyCipherText := types.ParametersType{
	// 	Datatype: "string",
	// 	Value:    hex.EncodeToString((keyCiphertext[:])),
	// 	// Value: "sample",
	// }

	// parametersArray := []types.ParametersType{KeywordArray, FileHash, KeyCipherText}

	// db.UploadFile(filePath, string(fileName)+".bin")
	// hedera.SetData(parametersArray)
	// configOperatorID := configs.EnvHederaOperatorID()
	// // Keyword := types.ParametersType{
	// // 	Datatype: "string",
	// // 	Value:    "Greek",
	// // }

	// // OwnerAddress := types.ParametersType{
	// // 	Datatype: "address",
	// // 	Value:    configOperatorID,
	// // }
	// // viewParametersArray := []types.ParametersType{OwnerAddress, Keyword}

	// // utils.ViewData(viewParametersArray)

	// // Generate Trapdoor data
	// searchString := []string{"Greek"}
	// SearchArray := types.ParametersType{
	// 	Datatype: "stringArray",
	// 	Array:    searchString,
	// }

	// DataOwner := types.ParametersType{
	// 	Datatype: "address",
	// 	Value:    configOperatorID,
	// }

	// EndTimeStamp := types.ParametersType{
	// 	Datatype: "uint256",
	// 	Value:    "20000",
	// }
	// trapdoorData := []types.ParametersType{SearchArray, DataOwner, EndTimeStamp}

	// retrievedKeyCiphertext := hedera.GenerateTrapdoor(trapdoorData)
	// key, errn := decryptKey(userPrivateKey, ownerPublicKey, retrievedKeyCiphertext, hash, label)

	// if errn != nil {
	// 	fmt.Println(errn)
	// 	os.Exit(1)
	// }
	// db.DownloadFile(fileName + ".bin")
	// decryptFile(key, fileName)

}

func encryptFile(key []byte) ([32]byte, string, string) {
	log.Print("File encryption example")

	plaintext, err := ioutil.ReadFile("sample/number1.txt")
	if err != nil {
		log.Fatal(err)
	}

	fileHash := createFileHash(plaintext)

	// The key should be 16 bytes (AES-128), 24 bytes (AES-192) or
	// 32 bytes (AES-256)
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panic(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Panic(err)
	}

	// Never use more than 2^32 random nonces with a given key
	// because of the risk of repeat.
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	fileName := hex.EncodeToString(fileHash[:])
	// fmt.Println(fileName)
	filePath := "upload/" + string(fileName) + ".bin"
	// Save back to file
	err = ioutil.WriteFile(filePath, ciphertext, 0777)
	if err != nil {
		log.Panic(err)
	}

	return fileHash, filePath, fileName
}

func decryptFile(key []byte, fileName string) {
	fmt.Println("In Decrypt")
	ciphertext, err := ioutil.ReadFile("download/" + fileName + ".bin")
	if err != nil {
		log.Fatal(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panic(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Panic(err)
	}

	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile("plaintext_new.txt", plaintext, 0777)
	if err != nil {
		log.Panic(err)
	}
}

func encryptKey(key []byte, userPublicKey *rsa.PublicKey, ownerPrivateKey *rsa.PrivateKey, hash hash.Hash, label []byte) ([]byte, error) {

	// Encrypt the key
	ciphertext, err := rsa.EncryptOAEP(
		hash,
		rand.Reader,
		userPublicKey,
		key,
		label,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("OAEP encrypted  to \n[%x]\n", ciphertext)

	return ciphertext, err

}

func decryptKey(userPrivateKey *rsa.PrivateKey, ownerPublicKey *rsa.PublicKey, ciphertext []byte, hash hash.Hash, label []byte) ([]byte, error) {
	key, err := rsa.DecryptOAEP(
		hash,
		rand.Reader,
		userPrivateKey,
		ciphertext,
		label,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// fmt.Printf("OAEP decrypted [%x] to \n[%s]\n", ciphertext, key)

	return key, err

}

func generateRandomHash() []byte {
	key := make([]byte, 32)

	_, err := rand.Read(key)
	if err != nil {
		// handle error here
	}

	return key
}

func createFileHash(fileContent []byte) [32]byte {
	fileHash := sha256.Sum256(fileContent)
	return fileHash
}

func generateKeywordHash(wordArray []string) []string {
	wordHashArray := make([]string, len(wordArray))

	for j := 0; j < len(wordArray); j++ {
		hashData := []byte(wordArray[j])
		wordHash := sha256.Sum256(hashData)
		//fmt.Println(hex.EncodeToString(wordHash[:]))
		wordHashArray[j] = hex.EncodeToString(wordHash[:])
	}

	return wordHashArray
}

func createEncIndexMap() {

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

func readUserKeys() (*rsa.PrivateKey, *rsa.PublicKey, *rsa.PrivateKey, *rsa.PublicKey, error) {
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

	ownerPrivateKey, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ownerPublicKey := &ownerPrivateKey.PublicKey

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

	userPrivateKey, err := x509.ParsePKCS1PrivateKey(userdata.Bytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	userPublicKey := &userPrivateKey.PublicKey

	return ownerPrivateKey, ownerPublicKey, userPrivateKey, userPublicKey, nil

}
