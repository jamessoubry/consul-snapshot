package backup

import (
	//"crypto/sha256"
	//"encoding/hex"
	"encoding/json"
	"fmt"
	//"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/pshima/consul-snapshot/consul"
)

const (
	writeBackupChecksum = "cd77faecc2e0dd4eb6e03e87354837763d4f6eb0742446ee6670d5944982c4e1"
	calc256teststring   = "Look Mom I'm Testing My SHA256!"
	tmpDir              = "/tmp"
)

var (
	// KV test data
	testkvpair1 = consulapi.KVPair{Key: "test1", Value: []byte("yes")}
	testkvpair2 = consulapi.KVPair{Key: "test2", Value: []byte("no")}
	kvpairlist  = []*consulapi.KVPair{&testkvpair1, &testkvpair2}

	// PQ test data
	testpqdata1 = &consulapi.PreparedQueryDefinition{ID: "99", Name: "pqtest1"}
	pqtestlist  = []*consulapi.PreparedQueryDefinition{testpqdata1}

	// ACL test data
	acltestdata1 = &consulapi.ACLEntry{ID: "98", Name: "acltest1"}
	acltestlist  = []*consulapi.ACLEntry{acltestdata1}
)

// Setup some basic structs we can use across tests
func testingStructs() *Backup {
	consulClient := &consul.Consul{}
	consulClient.KeyData = kvpairlist
	consulClient.PQData = pqtestlist
	consulClient.ACLData = acltestlist
	backup := &Backup{
		Client: consulClient,
	}
	backup.StartTime = time.Now().Unix()

	return backup
}

func writeTestFile(filename string, contents []byte) error {
	writepath := filepath.Join(tmpDir, filename)

	handle, err := os.OpenFile(writepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("Could not open local file for writing!: %v", err)
	}
	defer handle.Close()

	_, err = handle.Write(contents)
	if err != nil {
		return fmt.Errorf("Could not write data to file!: %v", err)
	}
	return nil
}

func TestCalc256(t *testing.T) {
	_, err := calcSha256("/temper/should/not/exist")
	if err == nil {
		t.Error("Can't read from a file that doesnt exist!!")
	}

	writeTestFile("testcalc256", []byte(calc256teststring))
	sha, err := calcSha256(filepath.Join(tmpDir, "testcalc256"))
	if err != nil {
		t.Errorf("Failed to call sha256: %v", err)
	}
	if sha != "79200c6c6b00cc4c2b7bd7e83fc54cfd2b9e8930127aeffe307fdd7631d9a8a0" {
		t.Errorf("Bad SHA: %v", sha)
	}
}

func TestKeysToJSON(t *testing.T) {
	backup := testingStructs()
	backup.KeysToJSON()

	marshallSouce, err := json.Marshal(kvpairlist)
	if err != nil {
		t.Errorf("Unable to marshall source testing data: %v", err)
	}

	reflecttest := reflect.DeepEqual(backup.KVJSONData, marshallSouce)

	if reflecttest != true {
		t.Errorf("JSON marshall did not equal. Got %v, expected %v", marshallSouce, backup.KVJSONData)
	}
}

func TestPQsToJSON(t *testing.T) {
	backup := testingStructs()
	backup.PQsToJSON()

	marshallSouce, err := json.Marshal(pqtestlist)
	if err != nil {
		t.Errorf("Unable to marshall source testing data: %v", err)
	}

	reflecttest := reflect.DeepEqual(backup.PQJSONData, marshallSouce)

	if reflecttest != true {
		t.Errorf("JSON marshall did not equal. Got %v, expected %v", marshallSouce, backup.PQJSONData)
	}
}

func TestACLsToJSON(t *testing.T) {
	backup := testingStructs()
	backup.ACLsToJSON()

	marshallSouce, err := json.Marshal(acltestlist)
	if err != nil {
		t.Errorf("Unable to marshall source testing data: %v", err)
	}

	reflecttest := reflect.DeepEqual(backup.ACLJSONData, marshallSouce)

	if reflecttest != true {
		t.Errorf("JSON marshall did not equal. Got %v, expected %v", marshallSouce, backup.ACLJSONData)
	}
}

/*
// Write the file locally and then check that we get the same checksum back
func TestWriteBackupLocal(t *testing.T) {
	backup := testingStructs()
	backup.KeysToJSON()

	backup.writeBackupLocal()
	shacheck := sha256.New()

	filepath := fmt.Sprintf("%v/%v", backup.LocalFilePath, backup.LocalFileName)
	f, err := ioutil.ReadFile(filepath)
	shacheck.Write(f)
	if err != nil {
		t.Errorf("Unable to read local backup file: %v", err)
	}

	expectedSum := writeBackupChecksum
	receivedSum := hex.EncodeToString(shacheck.Sum(nil))

	if expectedSum != receivedSum {
		t.Errorf("Expected checksum %s, recieved %s", expectedSum, receivedSum)
	}

	err = os.Remove(filepath)
	if err != nil {
		t.Errorf("Unable to remove temporary backup file: %v", err)
	}

}
*/
