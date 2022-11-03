package commands

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"embed"
	"encoding/pem"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"
)

func setConfigFile(path, data string, t *testing.T) {
	f, err := os.Create(filepath.Join(path))
	if err != nil {
		t.Error(err)
		return
	}

	defer f.Close()
	_, err = f.WriteString(strings.Trim(data, "\n"))
	if err != nil {
		t.Error(err)
		return
	}
	t.Setenv("IZE_CONFIG_FILE", filepath.Join(path))
}

func makeSSHKeyPair(pubKeyPath, privateKeyPath string) (string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return "", err
	}

	// generate and write private key as PEM
	privateKeyFile, err := os.Create(privateKeyPath)
	defer privateKeyFile.Close()
	if err != nil {
		return "", err
	}
	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return "", err
	}

	// generate and write public key
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(pubKeyPath, ssh.MarshalAuthorizedKey(pub), 0655)
	if err != nil {
		return "", err
	}

	return string(ssh.MarshalAuthorizedKey(pub)), nil
}

func resetEnv(environ []string) {
	for _, s := range environ {
		kv := strings.Split(s, "=")
		os.Setenv(kv[0], kv[1])
	}
}

func copyEmbedData(fsys embed.FS, sourceDir string, targetDir string) error {
	subdirs, err := fsys.ReadDir(sourceDir)
	if err != nil {
		return err
	}
	for _, d := range subdirs {
		sourcePath := path.Join(sourceDir, d.Name())
		if d.IsDir() {
			err = copyEmbedData(fsys, path.Join(sourceDir, d.Name()), path.Join(targetDir, d.Name()))
			if err != nil {
				return err
			}
		} else {
			localPath := filepath.Join(targetDir, d.Name())

			content, err := fsys.ReadFile(sourcePath)
			if err != nil {
				return err
			}
			err = os.MkdirAll(filepath.Dir(localPath), 0755)
			if err != nil {
				return err
			}
			err = os.WriteFile(localPath, content, 0755)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
