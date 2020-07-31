package nexus

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	//myJson "github.com/siangyeh8818/gdeyamlOperator/internal/json"
	//IO "github.com/siangyeh8818/gdeyamlOperator/internal/myIo"
)

type Nexus struct {
	NexusApiMethod          string
	NexusReqBody            string
	NexusOutputPattern      string
	NexusPromoteType        string
	NexusPromoteDestination string
	NexusPromoteUrl         string
	NexusPromoteSource      string
}

func PostNesusAPI(nexusurl string, nexus_user string, nexus_password string, request_body string) {

	fmt.Printf("your request url : %s\n", nexusurl)
	req, err := http.NewRequest("POST", nexusurl, strings.NewReader(request_body))
	if err != nil {
		fmt.Println(err)
		// handle err
	}
	if nexus_user != "" && nexus_password != "" {
		req.SetBasicAuth(nexus_user, nexus_password)
	}
	fmt.Println("---------------34---------------")
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "multipart/form-data")

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Println("request failed")
	}
	//responseData, err := ioutil.ReadAll(resp.Body)
	//log.Println(string(responseData))
	//log.Println(resp.Status)
	//log.Println(resp)

}

/*
func PutNesusAPI(nexusurl string, nexus_user string, nexus_password string, request_body string) {

	// curl -X POST "https://package.pentium.network/service/rest/v1/staging/move/events-preview?repository=events&name=siang-test%2F01%2Fevent.yml" -H "accept: application/json"

	//fmt.Printf("your request url : %s\n", nexusurl)
	req, err := http.NewRequest("PUT", nexusurl, strings.NewReader(request_body))
	if err != nil {
		// handle err
	}
	if nexus_user != "" && nexus_password != "" {
		req.SetBasicAuth(nexus_user, nexus_password)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Println("request failed")
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	//log.Println(string(responseData))
	//log.Println(resp.Status)
	//log.Println(resp)
}

func DeleteNesusAPI(nexusurl string, nexus_user string, nexus_password string, request_body string) {

	// curl -X POST "https://package.pentium.network/service/rest/v1/staging/move/events-preview?repository=events&name=siang-test%2F01%2Fevent.yml" -H "accept: application/json"

	fmt.Printf("your request url : %s\n", nexusurl)
	req, err := http.NewRequest("DELETE", nexusurl, strings.NewReader(request_body))
	if err != nil {
		// handle err
	}
	if nexus_user != "" && nexus_password != "" {
		req.SetBasicAuth(nexus_user, nexus_password)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Println("request failed")
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	log.Println(string(responseData))
	log.Println(resp.Status)
	log.Println(resp)
}
*/

func POSTForm_NesusAPI(filename string) {
	log.Println("POSTForm_NesusAPI")
	NexusServer := os.Getenv("NEXUS_SERVER")
	NexusUser := os.Getenv("NEXUS_USER")
	NexusPassword := os.Getenv("NEXUS_PASSWORD")
	NexusRepository := os.Getenv("NEXUS_REPOSITORY")
	nTime := time.Now()
	local1, _ := time.LoadLocation("Asia/Taipei") //等同于"CST"

	logDay := nTime.In(local1).Format("20060102")
	//curl -X POST "https://package.pentium.network/service/rest/v1/components?repository=scripts-qa"
	//-H "accept: application/json" -H "Content-Type: multipart/form-data"
	//-F "raw.directory=/byos-host-bootstrap-uninstall/0.3/dist" -F "raw.asset1=@script.zip;type=application/zip" -F "raw.asset1.filename=script.zip"
	path, _ := os.Getwd()
	path += "/" + filename
	//fmt.Printf("your request url : %s\n", nexusurl)
	extraParams := map[string]string{
		"raw.directory":       logDay,
		"raw.asset1":          "@" + filename + ";type=text/csv",
		"raw.asset1.filename": filename,
	}
	log.Println("start to newfileUploadRequest")
	req, err := newfileUploadRequest(NexusServer+"/service/rest/v1/components?repository="+NexusRepository, extraParams, filename, filename)
	//req, err := http.NewRequest("POST", nexusurl, strings.NewReader(request_body))
	if err != nil {
		// handle err
	}
	req.SetBasicAuth(NexusUser, NexusPassword)
	req.Header.Set("Accept", "application/json")

	//req.Header.Set("Content-Type", "multipart/form-data")
	//log.Println(req.Header)

	log.Println("------------------")
	log.Println(req)
	resp, err := http.DefaultClient.Do(req)

	defer resp.Body.Close()
	if err != nil {
		log.Println("request failed")
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	log.Println(string(responseData))
	log.Println(resp.Status)
	//log.Println(resp)

}

func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	log.Println("newfileUploadRequest")
	log.Println(path)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	//log.Println(body)
	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func ExecShell(s_command string) (string, string) {
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := exec.Command("/bin/bash", "-c", s_command)
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	var errStdout, errStderr error
	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	/*
		go func() {
			_, errStdout = io.Copy(stdout, stdoutIn)
		}()
		go func() {
			_, errStderr = io.Copy(stderr, stderrIn)
		}()
	*/
	_, errStdout = io.Copy(stdout, stdoutIn)
	_, errStderr = io.Copy(stderr, stderrIn)
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	if errStdout != nil || errStderr != nil {
		fmt.Printf("stdout: %v, stderr: %v\n", errStdout, errStderr)
		log.Fatal("failed to capture stdout or stderr\n")
	}
	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	return outStr, errStr
}
func RunCommand(commandStr string) string {
	cmdstr := commandStr
	out, _ := exec.Command("sh", "-c", cmdstr).Output()
	strout := string(out)

	return strout
}
