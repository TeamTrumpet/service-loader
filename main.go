package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/minio/minio-go"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "service-loader"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "bucket, b",
		},
		cli.StringFlag{
			Name: "access_key_id",
		},
		cli.StringFlag{
			Name: "aws_secret_access_key",
		},
	}
	app.Action = run
	app.Run(os.Args)
}

func generateSha256(b []byte) (string, error) {
	h := sha256.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil)), nil
}

func run(c *cli.Context) error {
	s3Client, err := minio.New("s3.amazonaws.com", c.String("access_key_id"), c.String("aws_secret_access_key"), true)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	bucketName := c.String("bucket")
	appName := c.Args()[0]
	tagName := c.Args()[1]

	goos := os.Getenv("GOOS")
	goarch := os.Getenv("GOARCH")
	if goos == "" || goarch == "" {
		goos = runtime.GOOS
		goarch = runtime.GOARCH
	}

	fileName := fmt.Sprintf("%s_%s_%s_%s", appName, tagName, goos, goarch)

	tarFileName := fileName + ".tar.gz"

	tarFile, err := s3Client.GetObject(bucketName, tarFileName)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	var localBuffer = bytes.NewBuffer(nil)

	if _, err = io.Copy(localBuffer, tarFile); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	shaFile, err := s3Client.GetObject(bucketName, fileName+".sha256")
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	remoteSha256, err := ioutil.ReadAll(shaFile)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	remoteChecksum := strings.Fields(string(remoteSha256))[0]

	localChecksum, err := generateSha256(localBuffer.Bytes())
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if remoteChecksum != localChecksum {
		err = errors.New("checksum mismatch")

		return cli.NewExitError(err.Error(), 1)
	}

	localFile, err := os.Create(tarFileName)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if _, err = io.Copy(localFile, localBuffer); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Printf("Downloaded the %s release of %s to %s.\n", tagName, appName, tarFileName)

	return nil
}
