package main

import "fmt"
import "net"
import "log"
import 	"io/ioutil"
import "github.com/pkg/sftp"
import "golang.org/x/crypto/ssh"
import "os"
import  "io"
import "encoding/base64"

var conn *ssh.Client
var sftpClient *sftp.Client
var err error

func createConn() {

  addr := "10.4.1.142:22"
   
  config := &ssh.ClientConfig{
   User: "sftpuser",
   Auth: []ssh.AuthMethod{
     ssh.Password("sftpuser"),
   },
   HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
     fmt.Println("server host key's type :", key.Type())
     

     fmt.Println(ssh.FingerprintSHA256(key))

     str := base64.StdEncoding.EncodeToString(key.Marshal())
     fmt.Println(str)
     return nil
   },
  }

 conn, err = ssh.Dial("tcp", addr, config)
 if err != nil {
   panic("Failed to dial: " + err.Error())
 }
 log.Println("inside function : ",conn)
// return conn
}


func createSession() {

// open an SFTP session over an existing ssh connection.
sftpClient, err = sftp.NewClient(conn)
if err != nil {
   log.Fatal(err)    /* panic("Failed to create client: " + err.Error()) */
}
//return sftpClient
log.Println("inside function : ",sftpClient)

}

//get working directory
func getWorkingDir (){
  cwd, err := sftpClient.Getwd()
  log.Println("Current working directory:", cwd)
  if err != nil {
    log.Fatal(err)
  }
}

// walk a directory
func walkDir (){
  w := sftpClient.Walk("/sftpuser/sftp-test/finalDir")
  for w.Step() {
     if w.Err() != nil {
         continue
     }
     log.Println(w.Path())
  }
}

// create a file
func leaveMark (content []byte){
  
  localFile, err := os.Open("/home/nishadi/Desktop/Booklet.pdf")
    
    defer localFile.Close()
    remoteFile, err := sftpClient.Create("sftpuser/sftp-test/finalDir/x.pdf")
    println(err)
    if err != nil {
      log.Fatal(err)
      }
    _, err = io.Copy(remoteFile, localFile)
    log.Fatal(err)
  
/*
  f, err := sftpClient.Create("/sftpuser/sftp-test/finalDir/build.sbt")
  if err != nil {
     log.Fatal(err)
  }
  if _, err := f.Write(content); err != nil {
     log.Fatal(err)
  }*/
}


// check it's there
func isThere (){
  fi, err := sftpClient.Lstat("/sftpuser/sftp-test/finalDir/hello1.txt")
  if err != nil {
     log.Fatal(err)
  }
  log.Println(fi)
}

//rename a file
func renameFile (){
  err := sftpClient.Rename("/sftpuser/sftp-test/tempDir/build.sbt","/sftpuser/sftp-test/finalDir/build2.sbt")
  if err != nil {
    log.Fatal(err)
  }  
}

//remove a file
func remove (){
  err := sftpClient.Remove("/sftpuser/sftp-test/finalDir/hello1.txt")
  if err != nil {
    log.Fatal(err)
  }
}

//Close connection
func close(){
  err := sftpClient.Close()
  if err != nil {
    log.Fatal(err)
  }
  log.Println("Done!")
}

//read a file (into a byte array)
func readToByteArray() []byte{
  content, err := ioutil.ReadFile("testdata.txt")
  if err != nil {
    log.Fatal(err)
  }
  return content
}


func main() {
  fmt.Printf("hello, world\n")
 var content []byte = readToByteArray();
  createConn();
  createSession();
  getWorkingDir();
  walkDir();
 leaveMark(content);
  renameFile();
 isThere();
 remove();
  //close  connection
  defer close();
  
}
