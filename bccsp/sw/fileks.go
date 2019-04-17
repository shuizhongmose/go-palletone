/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package sw

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/btcsuite/btcutil/base58"
	"github.com/palletone/go-palletone/bccsp"
	"github.com/palletone/go-palletone/bccsp/utils"
	"github.com/palletone/go-palletone/common/log"
    "github.com/tjfoc/gmsm/sm2"
)

// NewFileBasedKeyStore instantiated a file-based key store at a given position.
// The key store can be encrypted if a non-empty password is specified.
// It can be also be set as read only. In this case, any store operation
// will be forbidden
func NewFileBasedKeyStore(pwd []byte, path string, readOnly bool) (bccsp.KeyStore, error) {
	ks := &fileBasedKeyStore{}
	return ks, ks.Init(pwd, path, readOnly)
}

// fileBasedKeyStore is a folder-based KeyStore.
// Each key is stored in a separated file whose name contains the key's SKI
// and flags to identity the key's type. All the keys are stored in
// a folder whose path is provided at initialization time.
// The KeyStore can be initialized with a password, this password
// is used to encrypt and decrypt the files storing the keys.
// A KeyStore can be read only to avoid the overwriting of keys.
type fileBasedKeyStore struct {
	path string

	readOnly bool
	isOpen   bool

	pwd []byte

	// Sync
	m sync.Mutex
}

// Init initializes this KeyStore with a password, a path to a folder
// where the keys are stored and a read only flag.
// Each key is stored in a separated file whose name contains the key's SKI
// and flags to identity the key's type.
// If the KeyStore is initialized with a password, this password
// is used to encrypt and decrypt the files storing the keys.
// The pwd can be nil for non-encrypted KeyStores. If an encrypted
// key-store is initialized without a password, then retrieving keys from the
// KeyStore will fail.
// A KeyStore can be read only to avoid the overwriting of keys.
func (ks *fileBasedKeyStore) Init(pwd []byte, path string, readOnly bool) error {
	// Validate inputs
	// pwd can be nil

	if len(path) == 0 {
		return errors.New("An invalid KeyStore path provided. Path cannot be an empty string.")
	}

	ks.m.Lock()
	defer ks.m.Unlock()

	if ks.isOpen {
		return errors.New("KeyStore already initilized.")
	}

	ks.path = path
	ks.pwd = utils.Clone(pwd)

	err := ks.createKeyStoreIfNotExists()
	if err != nil {
		return err
	}

	err = ks.openKeyStore()
	if err != nil {
		return err
	}

	ks.readOnly = readOnly

	return nil
}

// ReadOnly returns true if this KeyStore is read only, false otherwise.
// If ReadOnly is true then StoreKey will fail.
func (ks *fileBasedKeyStore) ReadOnly() bool {
	return ks.readOnly
}

// GetKey returns a key object whose SKI is the one passed.
func (ks *fileBasedKeyStore) GetKey(ski []byte) (bccsp.Key, error) {
	// Validate arguments
	if len(ski) == 0 {
		return nil, errors.New("Invalid SKI. Cannot be of zero length.")
	}

	suffix := ks.getSuffix(ski2Address(ski))

	switch suffix {
	case "key":
		// Load the key
		key, err := ks.loadKey(ski2Address(ski))
		if err != nil {
			return nil, fmt.Errorf("Failed loading key [%x] [%s]", ski, err)
		}

		return &aesPrivateKey{key, false}, nil
	case "sk":
		// Load the private key
		key, err := ks.loadPrivateKey(ski2Address(ski))
		if err != nil {
			return nil, fmt.Errorf("Failed loading secret key [%x] [%s]", ski, err)
		}

		switch key.(type) {
		case *ecdsa.PrivateKey:
			return &ecdsaPrivateKey{key.(*ecdsa.PrivateKey)}, nil
		case *rsa.PrivateKey:
			return &rsaPrivateKey{key.(*rsa.PrivateKey)}, nil
		default:
			return nil, errors.New("Secret key type not recognized")
		}
	case "pk":
		// Load the public key
		key, err := ks.loadPublicKey(ski2Address(ski))
		if err != nil {
			return nil, fmt.Errorf("Failed loading public key [%x] [%s]", ski, err)
		}

		switch key.(type) {
		case *ecdsa.PublicKey:
			return &ecdsaPublicKey{key.(*ecdsa.PublicKey)}, nil
		case *rsa.PublicKey:
			return &rsaPublicKey{key.(*rsa.PublicKey)}, nil
		default:
			return nil, errors.New("Public key type not recognized")
		}
	default:
		return ks.searchKeystoreForSKI(ski)
	}
}

// StoreKey stores the key k in this KeyStore.
// If this KeyStore is read only then the method will fail.
func (ks *fileBasedKeyStore) StoreKey(k bccsp.Key) (err error) {
	if ks.readOnly {
		return errors.New("Read only KeyStore.")
	}

	if k == nil {
		return errors.New("Invalid key. It must be different from nil.")
	}
	switch k.(type) {
	case *ecdsaPrivateKey:
		kk := k.(*ecdsaPrivateKey)

		err = ks.storePrivateKey(ski2Address(k.SKI()), kk.privKey)
		if err != nil {
			return fmt.Errorf("Failed storing ECDSA private key [%s]", err)
		}

	case *ecdsaPublicKey:
		kk := k.(*ecdsaPublicKey)

		err = ks.storePublicKey(ski2Address(k.SKI()), kk.pubKey)
		if err != nil {
			return fmt.Errorf("Failed storing ECDSA public key [%s]", err)
		}

	case *rsaPrivateKey:
		kk := k.(*rsaPrivateKey)

		err = ks.storePrivateKey(ski2Address(k.SKI()), kk.privKey)
		if err != nil {
			return fmt.Errorf("Failed storing RSA private key [%s]", err)
		}

	case *rsaPublicKey:
		kk := k.(*rsaPublicKey)

		err = ks.storePublicKey(ski2Address(k.SKI()), kk.pubKey)
		if err != nil {
			return fmt.Errorf("Failed storing RSA public key [%s]", err)
		}

	case *aesPrivateKey:
		kk := k.(*aesPrivateKey)

		err = ks.storeKey(ski2Address(k.SKI()), kk.privKey)
		if err != nil {
			return fmt.Errorf("Failed storing AES key [%s]", err)
		}

	default:
		return fmt.Errorf("Key type not reconigned [%s]", k)
	}

	return
}
func ski2Address(ski []byte) string {
	if len(ski) < 20 {
		return ""
	}
	return "P" + base58.CheckEncode(ski[0:20], byte(0))
}
func (ks *fileBasedKeyStore) searchKeystoreForSKI(ski []byte) (k bccsp.Key, err error) {

	files, _ := ioutil.ReadDir(ks.path)
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if f.Size() > (1 << 16) { //64k, somewhat arbitrary limit, considering even large RSA keys
			continue
		}

		raw, err := ioutil.ReadFile(filepath.Join(ks.path, f.Name()))
		if err != nil {
			continue
		}

		key, err := utils.PEMtoPrivateKey(raw, ks.pwd)
		if err != nil {
			continue
		}

		switch key.(type) {
		case *ecdsa.PrivateKey:
			k = &ecdsaPrivateKey{key.(*ecdsa.PrivateKey)}
		case *rsa.PrivateKey:
			k = &rsaPrivateKey{key.(*rsa.PrivateKey)}
		default:
			continue
		}

		if !bytes.Equal(k.SKI(), ski) {
			continue
		}

		return k, nil
	}
	return nil, fmt.Errorf("Key with SKI %s not found in %s", ski2Address(ski), ks.path)
}

func (ks *fileBasedKeyStore) getSuffix(alias string) string {
	files, _ := ioutil.ReadDir(ks.path)
	for _, f := range files {
		if strings.HasPrefix(f.Name(), alias) {
			if strings.HasSuffix(f.Name(), "sk") {
				return "sk"
			}
			if strings.HasSuffix(f.Name(), "pk") {
				return "pk"
			}
			if strings.HasSuffix(f.Name(), "key") {
				return "key"
			}
			break
		}
	}
	return ""
}

func (ks *fileBasedKeyStore) storePrivateKey(alias string, privateKey interface{}) error {
	log.Debugf("Try store private key to file.%s",alias)
	rawKey, err := utils.PrivateKeyToPEM(privateKey, ks.pwd)
	if err != nil {
		log.Errorf("Failed converting private key to PEM [%s]: [%s]", alias, err)
		return err
	}

	err = ioutil.WriteFile(ks.getPathForAlias(alias, "sk"), rawKey, 0600)
	if err != nil {
		log.Errorf("Failed storing private key [%s]: [%s]", alias, err)
		return err
	}

	return nil
}

func (ks *fileBasedKeyStore) storePublicKey(alias string, publicKey interface{}) error {
	rawKey, err := utils.PublicKeyToPEM(publicKey, ks.pwd)
	if err != nil {
		log.Errorf("Failed converting public key to PEM [%s]: [%s]", alias, err)
		return err
	}

	err = ioutil.WriteFile(ks.getPathForAlias(alias, "pk"), rawKey, 0600)
	if err != nil {
		log.Errorf("Failed storing private key [%s]: [%s]", alias, err)
		return err
	}

	return nil
}

func (ks *fileBasedKeyStore) storeKey(alias string, key []byte) error {
	pem, err := utils.AEStoEncryptedPEM(key, ks.pwd)
	if err != nil {
		log.Errorf("Failed converting key to PEM [%s]: [%s]", alias, err)
		return err
	}

	err = ioutil.WriteFile(ks.getPathForAlias(alias, "key"), pem, 0600)
	if err != nil {
		log.Errorf("Failed storing key [%s]: [%s]", alias, err)
		return err
	}

	return nil
}

func (ks *fileBasedKeyStore) loadPrivateKey(alias string) (interface{}, error) {
	path := ks.getPathForAlias(alias, "sk")
	log.Debugf("Loading private key [%s] at [%s]...", alias, path)

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("Failed loading private key [%s]: [%s].", alias, err.Error())

		return nil, err
	}

	privateKey, err := utils.PEMtoPrivateKey(raw, ks.pwd)
	if err != nil {
		log.Errorf("Failed parsing private key [%s]: [%s].", alias, err.Error())

		return nil, err
	}

	return privateKey, nil
}

func (ks *fileBasedKeyStore) loadPublicKey(alias string) (interface{}, error) {
	path := ks.getPathForAlias(alias, "pk")
	log.Debugf("Loading public key [%s] at [%s]...", alias, path)

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("Failed loading public key [%s]: [%s].", alias, err.Error())

		return nil, err
	}

	privateKey, err := utils.PEMtoPublicKey(raw, ks.pwd)
	if err != nil {
		log.Errorf("Failed parsing private key [%s]: [%s].", alias, err.Error())

		return nil, err
	}

	return privateKey, nil
}

func (ks *fileBasedKeyStore) loadKey(alias string) ([]byte, error) {
	path := ks.getPathForAlias(alias, "key")
	log.Debugf("Loading key [%s] at [%s]...", alias, path)

	pem, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("Failed loading key [%s]: [%s].", alias, err.Error())

		return nil, err
	}

	key, err := utils.PEMtoAES(pem, ks.pwd)
	if err != nil {
		log.Errorf("Failed parsing key [%s]: [%s]", alias, err)

		return nil, err
	}

	return key, nil
}

func (ks *fileBasedKeyStore) createKeyStoreIfNotExists() error {
	// Check keystore directory
	ksPath := ks.path
	missing, err := utils.DirMissingOrEmpty(ksPath)

	if missing {
		log.Debugf("KeyStore path [%s] missing [%t]: [%s]", ksPath, missing, utils.ErrToString(err))

		err := ks.createKeyStore()
		if err != nil {
			log.Errorf("Failed creating KeyStore At [%s]: [%s]", ksPath, err.Error())
			return nil
		}
	}

	return nil
}

func (ks *fileBasedKeyStore) createKeyStore() error {
	// Create keystore directory root if it doesn't exist yet
	ksPath := ks.path
	log.Debugf("Creating KeyStore at [%s]...", ksPath)

	os.MkdirAll(ksPath, 0755)

	log.Debugf("KeyStore created at [%s].", ksPath)
	return nil
}

func (ks *fileBasedKeyStore) openKeyStore() error {
	if ks.isOpen {
		return nil
	}
	ks.isOpen = true
	log.Debugf("KeyStore opened at [%s]...done", ks.path)

	return nil
}

func (ks *fileBasedKeyStore) getPathForAlias(alias, suffix string) string {
	return filepath.Join(ks.path, alias+"_"+suffix)
}

func readKeyFromFile(privKeyPath , pubKeyPath string , pass []byte)(*sm2.PrivateKey,*sm2.PublicKey,  bool){
    fmt.Println("begin to read！！！",privKeyPath)
    fmt.Println("pass is ！！！",pass)
	privateKey, e := sm2.ReadPrivateKeyFromPem(privKeyPath, pass)
	if e != nil{
		fmt.Println("failed to read privateKey ！！！")
		return nil,nil,false
	}
	fmt.Printf("privateKey is %+v\n",privateKey)
	publicKey, i := sm2.ReadPublicKeyFromPem(pubKeyPath, pass)
	if i!=nil{
		fmt.Println("failed to read publicKey ！！！")
		return nil,nil,false
	}
	fmt.Printf("publicKey is %+v\n",publicKey)
	return privateKey,publicKey,true
}
