package aebmes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/warlockz/ase-service/db"
	"github.com/warlockz/ase-service/types"
	keys "github.com/warlockz/ase-service/utils"
	"github.com/warlockz/ase-service/utils/hedera"
)

var hashValue hash.Hash = sha256.New()

func EncryptFile(keywordArray []string, file string) {
	start := time.Now()

	ownerPrivateKey := keys.OwnerPrivateKey()
	userPublicKey := keys.UserPublicKey()
	// fmt.Println(ownerPrivateKey, userPublicKey)
	key := make([]byte, 32)
	// Encrypt File using key
	fileHash, filePath, fileName := _encryptFile(key, file)
	fmt.Println("File Encrypted AES")

	label := []byte("")
	// Encrypt the key

	keyCiphertext, err := _encryptKey(key, userPublicKey, ownerPrivateKey, hashValue, label)
	fmt.Println("Key Encrypted RSA")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	KeywordArray := types.ParametersType{
		Datatype: "stringArray",
		Array:    keywordArray[:],
	}

	FileHash := types.ParametersType{
		Datatype: "string",
		Value:    hex.EncodeToString(fileHash[:]),
	}

	KeyCipherText := types.ParametersType{
		Datatype: "string",
		Value:    hex.EncodeToString((keyCiphertext[:])),
		// Value: "sample",
	}

	parametersArray := []types.ParametersType{KeywordArray, FileHash, KeyCipherText}

	db.UploadFile(filePath, string(fileName)+".bin")

	hedera.SetData(parametersArray)
	fmt.Println("Data stored in blockchain")

	elapsed := time.Since(start)
	fmt.Println("Total Time Taken to Encrypt: ", +elapsed)
}

func GenerateTrapdoor(keywordArray []string, ownerId string) (string, error) {
	start := time.Now()

	ownerPublicKey := keys.OwnerPublicKey()
	userPrivateKey := keys.UserPrivateKey()
	SearchArray := types.ParametersType{
		Datatype: "stringArray",
		Array:    keywordArray,
	}

	DataOwner := types.ParametersType{
		Datatype: "address",
		Value:    ownerId,
	}

	EndTimeStamp := types.ParametersType{
		Datatype: "uint256",
		Value:    "20000",
	}

	trapdoorData := []types.ParametersType{SearchArray, DataOwner, EndTimeStamp}

	// Search smart contract uisng trapdoor
	retrievedKeyCiphertext, fileName, trapdoorHash := hedera.GenerateTrapdoor(trapdoorData)
	label := []byte("")
	fmt.Println("Trapdoor stored in Blockchain")
	// Verify Trapdoore
	TrapdoorHash := types.ParametersType{
		Datatype: "bytes32",
		Bytes32:  trapdoorHash,
	}

	DataUser := types.ParametersType{
		Datatype: "address",
		Value:    ownerId,
	}

	verifyTrapdoorData := []types.ParametersType{TrapdoorHash, DataUser}

	verifyTrapdoor := hedera.VerifyTrapdoor(verifyTrapdoorData)
	if !verifyTrapdoor {
		return "nil", errors.New("Trapdoor not verified")
	}

	fmt.Println("Trapdoor Verified")

	key, errn := _decryptKey(userPrivateKey, ownerPublicKey, retrievedKeyCiphertext, hashValue, label)
	fmt.Println("Key Decrypted RSA ")

	if errn != nil {
		fmt.Println(errn)
		os.Exit(1)
	}
	db.DownloadFile(fileName + ".bin")
	newFileHash, newFileName := _decryptFile(key, fileName)
	fmt.Println("File Decrypted AES ")

	//Verify Search
	NewFileHash := types.ParametersType{
		Datatype: "string",
		Value:    hex.EncodeToString(newFileHash[:]),
	}
	verifySearchData := []types.ParametersType{TrapdoorHash, NewFileHash, DataUser}
	verifyResult := hedera.VerifyResult(verifySearchData)

	if !verifyResult {
		return "", errors.New("Result not verified")
	}

	fmt.Println("Result Verified")
	elapsed := time.Since(start)
	fmt.Println("Total Time Taken to Decrypt: ", +elapsed)
	return newFileName, nil
}

func _encryptFile(key []byte, file string) ([32]byte, string, string) {
	log.Print("File encryption example")

	plaintext, err := ioutil.ReadFile("upload/" + file)
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

func _decryptFile(key []byte, fileName string) ([32]byte, string) {
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
	newFileName := uuid.New().String()
	fileHash := createFileHash(plaintext)

	err = ioutil.WriteFile("output/"+newFileName, plaintext, 0777)
	if err != nil {
		log.Panic(err)
	}
	return fileHash, "output/" + newFileName
}

func _decryptKey(userPrivateKey *rsa.PrivateKey, ownerPublicKey *rsa.PublicKey, ciphertext []byte, hash hash.Hash, label []byte) ([]byte, error) {
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

func _encryptKey(key []byte, userPublicKey *rsa.PublicKey, ownerPrivateKey *rsa.PrivateKey, hash hash.Hash, label []byte) ([]byte, error) {
	fmt.Println(key, userPublicKey, ownerPrivateKey, hash, label)
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
