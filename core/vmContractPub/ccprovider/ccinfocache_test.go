/*
	This file is part of go-palletone.
	go-palletone is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.
	go-palletone is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.
	You should have received a copy of the GNU General Public License
	along with go-palletone.  If not, see <http://www.gnu.org/licenses/>.
*/
/*
 * Copyright IBM Corp. All Rights Reserved.
 * @author PalletOne core developers <dev@pallet.one>
 * @date 2018
 */

package ccprovider

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/palletone/go-palletone/core/vmContractPub/protos/peer"
	"github.com/stretchr/testify/assert"
	"github.com/palletone/go-palletone/vm/common"
)

func getDepSpec(name string, path string, version string, initArgs [][]byte) (*peer.ChaincodeDeploymentSpec, error) {
	spec := &peer.ChaincodeSpec{Type: 1, ChaincodeId: &peer.ChaincodeID{Name: name, Path: path, Version: version}, Input: &peer.ChaincodeInput{Args: initArgs}}

	codePackageBytes := bytes.NewBuffer(nil)
	gz := gzip.NewWriter(codePackageBytes)
	tw := tar.NewWriter(gz)

	err := util.WriteBytesToPackage("src/garbage.go", []byte(name+path+version), tw)
	if err != nil {
		return nil, err
	}

	tw.Close()
	gz.Close()

	return &peer.ChaincodeDeploymentSpec{ChaincodeSpec: spec, CodePackage: codePackageBytes.Bytes()}, nil
}

func buildPackage(name string, path string, version string, initArgs [][]byte) (CCPackage, error) {
	depSpec, err := getDepSpec(name, path, version, initArgs)
	if err != nil {
		return nil, err
	}

	buf, err := proto.Marshal(depSpec)
	if err != nil {
		return nil, err
	}
	cccdspack := &CDSPackage{}
	if _, err := cccdspack.InitFromBuffer(buf); err != nil {
		return nil, err
	}

	return cccdspack, nil
}

type mockCCInfoFSStorageMgrImpl struct {
	CCMap map[string]CCPackage
}

func (m *mockCCInfoFSStorageMgrImpl) GetChaincode(ccname string, ccversion string) (CCPackage, error) {
	return m.CCMap[ccname+ccversion], nil
}

// here we test the cache implementation itself
func TestCCInfoCache(t *testing.T) {
	ccname := "foo"
	ccver := "1.0"
	ccpath := "github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02"

	ccinfoFs := &mockCCInfoFSStorageMgrImpl{CCMap: map[string]CCPackage{}}
	cccache := NewCCInfoCache(ccinfoFs)

	// test the get side

	// the cc data is not yet in the cache
	_, err := cccache.GetChaincodeData(ccname, ccver)
	assert.Error(t, err)

	// put it in the file system
	pack, err := buildPackage(ccname, ccpath, ccver, [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")})
	assert.NoError(t, err)
	ccinfoFs.CCMap[ccname+ccver] = pack

	// expect it to be in the cache now
	cd1, err := cccache.GetChaincodeData(ccname, ccver)
	assert.NoError(t, err)

	// it should still be in the cache
	cd2, err := cccache.GetChaincodeData(ccname, ccver)
	assert.NoError(t, err)

	// they are not null
	assert.NotNil(t, cd1)
	assert.NotNil(t, cd2)

	// test the put side now..
	ccver = "2.0"
	// put it in the file system
	pack, err = buildPackage(ccname, ccpath, ccver, [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")})
	assert.NoError(t, err)
	ccinfoFs.CCMap[ccname+ccver] = pack

	// create a dep spec to put
	_, err = getDepSpec(ccname, ccpath, ccver, [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")})
	assert.NoError(t, err)

	// expect it to be cached
	cd1, err = cccache.GetChaincodeData(ccname, ccver)
	assert.NoError(t, err)

	// it should still be in the cache
	cd2, err = cccache.GetChaincodeData(ccname, ccver)
	assert.NoError(t, err)

	// they are not null
	assert.NotNil(t, cd1)
	assert.NotNil(t, cd2)
}

func TestPutChaincode(t *testing.T) {
	ccname := ""
	ccver := "1.0"
	ccpath := "github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02"

	ccinfoFs := &mockCCInfoFSStorageMgrImpl{CCMap: map[string]CCPackage{}}
	NewCCInfoCache(ccinfoFs)

	// Error case 1: ccname is empty
	// create a dep spec to put
	_, err := getDepSpec(ccname, ccpath, ccver, [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")})
	assert.NoError(t, err)

	// Error case 2: ccver is empty
	ccname = "foo"
	ccver = ""
	_, err = getDepSpec(ccname, ccpath, ccver, [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")})
	assert.NoError(t, err)

	// Error case 3: ccfs.PutChainCode returns an error
	ccinfoFs = &mockCCInfoFSStorageMgrImpl{CCMap: map[string]CCPackage{}}
	NewCCInfoCache(ccinfoFs)

	ccname = "foo"
	ccver = "1.0"
	_, err = getDepSpec(ccname, ccpath, ccver, [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")})
	assert.NoError(t, err)
}

// here we test the peer's built-in cache after enabling it
func TestCCInfoFSPeerInstance(t *testing.T) {
	ccname := "bar"
	ccver := "1.0"
	ccpath := "github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02"

	// the cc data is not yet in the cache
	_, err := GetChaincodeFromFS(ccname, ccver)
	assert.Error(t, err)

	// create a dep spec to put
	ds, err := getDepSpec(ccname, ccpath, ccver, [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")})
	assert.NoError(t, err)

	// put it
	err = PutChaincodeIntoFS(ds)
	assert.NoError(t, err)

	// Get all installed chaincodes, it should not return 0 chaincodes
	resp, err := GetInstalledChaincodes()
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotZero(t, len(resp.Chaincodes), "GetInstalledChaincodes should not have returned 0 chaincodes")

	//get chaincode data
	_, err = GetChaincodeData(ccname, ccver)
	assert.NoError(t, err)
}

func TestGetInstalledChaincodesErrorPaths(t *testing.T) {
	// Get the existing chaincode install path value and set it
	// back after we are done with the test
	cip := chaincodeInstallPath
	defer SetChaincodesPath(cip)

	// Create a temp dir and remove it at the end
	dir, err := ioutil.TempDir(os.TempDir(), "chaincodes")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	// Set the above created directory as the chaincode install path
	SetChaincodesPath(dir)
	err = ioutil.WriteFile(filepath.Join(dir, "idontexist.1.0"), []byte("test"), 0777)
	assert.NoError(t, err)
	resp, err := GetInstalledChaincodes()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(resp.Chaincodes),
		"Expected 0 chaincodes but GetInstalledChaincodes returned %s chaincodes", len(resp.Chaincodes))
}

func TestNewCCContext(t *testing.T) {
	ccctx := NewCCContext("foo", "foo", "1.0", "", false, nil, nil)
	assert.NotNil(t, ccctx)
	canName := ccctx.GetCanonicalName()
	assert.NotEmpty(t, canName)

	assert.Panics(t, func() {
		NewCCContext("foo", "foo", "", "", false, nil, nil)
	}, "NewCCContext should have paniced if version is empty")

	ccctx = &CCContext{"foo", "foo", "1.0", "", false, nil, nil, "", nil}
	assert.Panics(t, func() {
		ccctx.GetCanonicalName()
	}, "GetConnonicalName should have paniced if cannonical name is empty")
}

func TestChaincodePackageExists(t *testing.T) {
	_, err := ChaincodePackageExists("foo1", "1.0")
	assert.Error(t, err)
}

func TestSetChaincodesPath(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "setchaincodes")
	if err != nil {
		assert.Fail(t, err.Error(), "Unable to create temp dir")
	}
	defer os.RemoveAll(dir)
	t.Logf("created temp dir %s", dir)

	// Get the existing chaincode install path value and set it
	// back after we are done with the test
	cip := chaincodeInstallPath
	defer SetChaincodesPath(cip)

	f, err := ioutil.TempFile(dir, "chaincodes")
	assert.NoError(t, err)
	assert.Panics(t, func() {
		SetChaincodesPath(f.Name())
	}, "SetChaincodesPath should have paniced if a file is passed to it")

	// Following code works on mac but does not work in CI
	// // Make the directory read only
	// err = os.Chmod(dir, 0444)
	// assert.NoError(t, err)
	// cdir := filepath.Join(dir, "chaincodesdir")
	// assert.Panics(t, func() {
	// 	SetChaincodesPath(cdir)
	// }, "SetChaincodesPath should have paniced if it is not able to stat the dir")

	// // Make the directory read and execute
	// err = os.Chmod(dir, 0555)
	// assert.NoError(t, err)
	// assert.Panics(t, func() {
	// 	SetChaincodesPath(cdir)
	// }, "SetChaincodesPath should have paniced if it is not able to create the dir")
}

var ccinfocachetestpath = "/tmp/ccinfocachetest"

func TestMain(m *testing.M) {
	os.RemoveAll(ccinfocachetestpath)

	SetChaincodesPath(ccinfocachetestpath)
	rc := m.Run()
	os.RemoveAll(ccinfocachetestpath)
	os.Exit(rc)
}
