package logging

import (
	"fmt"
	"github.com/svenschaper/goproperties"
	"time"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"crypto/sha512"
	"io/ioutil"
	"strings"
	"regexp"
)

type Logger struct {
	Class    string
	Loglevel int
}

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[34m"
var White = "\033[97m"
var Secret = "\033[1;37m"

const (
	ERROR       = 1
	COMMUNICATE = 2
	REPORT      = 3
	INFO        = 4
	WARN        = 5
	DEBUG       = 6
)

var logger Logger
var prop *properties.Propertie
var certpath string

func init() {
	properties.SetPropertiePath("config.yml")
	prop2, _ := properties.LoadProperty()
	prop = prop2
}

func GeneralInitLogger(class string) Logger {
	var logger Logger
	s := prop.GetProperty("log.level." + class)
	if s != "" {
		return logger.initialize(class, s)
	}
	return logger.initialize(class, prop.GetProperty("log.level"))
}

func (l Logger) initialize(class string, level string) Logger {
	var loglevel int
	switch level {
	case "ERROR":
		loglevel = 1
	case "COMMUNICATE":
		loglevel = 2
	case "REPORT":
		loglevel = 3
	case "INFO":
		loglevel = 4
	case "WARN":
		loglevel = 5
	case "DEBUG":
		loglevel = 6
	}
	return Logger{Class: class, Loglevel: loglevel}
}

func (l Logger) Info(s string, args ...interface{} ) {

	if l.Loglevel >= INFO {
		fmt.Println(Green + "[" + time.Now().Format(time.RFC850) + "]" + "   " + "   INFO          [" + l.Class + "]   " + fmt.Sprintf(s, args...) + Reset)
	}
}

func (l Logger) Error(s string, args ...interface{}) {
	if l.Loglevel >= ERROR {
		fmt.Println(Red + "[" + time.Now().Format(time.RFC850) + "]" + "   " + "   ERROR         [" + l.Class + "]   " + fmt.Sprintf(s, args...) + Reset)
	}
}

func (l Logger) Warn(s string, args ...interface{}) {
	if l.Loglevel >= WARN {
		fmt.Println(Purple + "[" + time.Now().Format(time.RFC850) + "]" + "   " + "   WARN          [" + l.Class + "]   " + fmt.Sprintf(s, args...) + Reset)
	}
}

func (l Logger) Debug(s string, args ...interface{}) {
	if l.Loglevel >= DEBUG {
		fmt.Println(Blue + "[" + time.Now().Format(time.RFC850) + "]" + "   " + "   DEBUG         [" + l.Class + "]   " + fmt.Sprintf(s, args...) + Reset)
	}
}

func (l Logger) Report(s string, args ...interface{}) {
	if l.Loglevel >= REPORT {
		fmt.Println(Cyan + "[" + time.Now().Format(time.RFC850) + "]" + "   " + "   REPORT        [" + l.Class + "]   " + fmt.Sprintf(s, args...) + Reset)
	}
}

func (l Logger) Communicate(s string, args ...interface{}) {
	if l.Loglevel >= COMMUNICATE {
		fmt.Println(Yellow + "[" + time.Now().Format(time.RFC850) + "]" + "   " + "   COMMUNICATE   [" + l.Class + "]   " + fmt.Sprintf(s, args...) + Reset)
	}
}



func (l Logger) InfoWithC(s string, encrypted string, args ...interface{}) {

	if l.Loglevel >= INFO {
		fmt.Println(Green + "[" + time.Now().Format(time.RFC850) + "]" + "   " + "   INFO          [" + l.Class + "]   " + fmt.Sprintf(s, args...)  + Secret+ "          [Secured Information] " + EncryptMessage(encrypted) + Reset)
	}
}

func (l Logger) ErrorWithC(s string, encrypted string, args ...interface{}) {
	if l.Loglevel >= ERROR {
		fmt.Println(Red + "[" + time.Now().Format(time.RFC850) + "]" + "   " + "   ERROR         [" + l.Class + "]   " + fmt.Sprintf(s, args...) + Secret+ "          [Secured Information] " + EncryptMessage(encrypted) + Reset)
	}
}

func (l Logger) WarnWithC(s string, encrypted string, args ...interface{}) {
	if l.Loglevel >= WARN {
		fmt.Println(Purple + "[" + time.Now().Format(time.RFC850) + "]" + "   " + "   WARN          [" + l.Class + "]   " + fmt.Sprintf(s, args...) + Secret+ "          [Secured Information] " + EncryptMessage(encrypted) + Reset)
	}
}

func (l Logger) DebugWithC(s string, encrypted string, args ...interface{}) {
	if l.Loglevel >= DEBUG {
		fmt.Println(Blue + "[" + time.Now().Format(time.RFC850) + "]" + "   " + "   DEBUG         [" + l.Class + "]   " + fmt.Sprintf(s, args...) + Secret+ "          [Secured Information] " + EncryptMessage(encrypted) + Reset)
	}
}

func (l Logger) ReportWithC(s string, encrypted string, args ...interface{}) {
	if l.Loglevel >= REPORT {
		fmt.Println(Cyan + "[" + time.Now().Format(time.RFC850) + "]" + "   " + "   REPORT        [" + l.Class + "]   " + fmt.Sprintf(s, args...) + Secret+ "          [Secured Information] " + EncryptMessage(encrypted) + Reset)
	}
}

func (l Logger) CommunicateWithC(s string, encrypted string, args ...interface{}) {
	if l.Loglevel >= COMMUNICATE {
		fmt.Println(Yellow + "[" + time.Now().Format(time.RFC850) + "]" + "   " + "   COMMUNICATE   [" + l.Class + "]   " + fmt.Sprintf(s, args...) + Secret+ "          [Secured Information] " + EncryptMessage(encrypted) + Reset)
	}
}

func EncryptMessage(s string) string{
	x,e:= Encrypt(prop.GetProperty("cert.path"), s)
	if e  != nil {
		fmt.Println(e.Error())
	}
	return extractMessageFromBlock(x)
}

func extractMessageFromBlock(input string) string{
	s := strings.Split(strings.Split(input, " -----")[1], "-----END")[0]
	re := regexp.MustCompile(`\r?\n`)
	input = re.ReplaceAllString(s, " ")
	return input	
}


func Encrypt(publicKeyPath, plainText string) (string, error) {
	bytes, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return "", err
	}
 
	publicKey, err := convertBytesToPublicKey(bytes)
	if err != nil {
		return "", err
	}
 
	cipher, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, publicKey, []byte(plainText), nil)
	if err != nil {
		return "", err
	}
 
	return cipherToPemString(cipher), nil
}


func convertBytesToPublicKey(keyBytes []byte) (*rsa.PublicKey, error) {
	var err error
 
	block, _ := pem.Decode(keyBytes)
	blockBytes := block.Bytes
	ok := x509.IsEncryptedPEMBlock(block)
 
	if ok {
		blockBytes, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
	}
 
	publicKey, err := x509.ParsePKCS1PublicKey(blockBytes)
	if err != nil {
		return nil, err
	}
 
	return publicKey, nil
}
 
func cipherToPemString(cipher []byte) string {
	return string(
		pem.EncodeToMemory(
			&pem.Block{
				Bytes: cipher,
			},
		),
	)
}


