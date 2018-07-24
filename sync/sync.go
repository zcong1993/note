package sync

import (
	"context"
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"github.com/spf13/viper"
	"github.com/zcong1993/note/utils"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	// SYNC_CONFIG_NAME if sync client config filename
	SYNC_CONFIG_NAME = "note_sync_config"
	// DB_Name is bolt db filename
	DB_Name = ".note.bolt.db"
	// FILE_NAME is qiniu s3 file key
	FILE_NAME = "note-nolt.db"
)

// Client is sync client of qiniu
type Client struct {
	Bucket      string
	Domain      string
	QnAccessKey string
	QnSecretKey string
}

// NewClient is constructor for sync client
func NewClient() *Client {
	viper.AutomaticEnv()
	viper.SetConfigName(SYNC_CONFIG_NAME)
	viper.AddConfigPath("$HOME")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	Bucket := viper.GetString("Bucket")
	if Bucket == "" {
		log.Fatal("Bucket should not be nil")
	}
	Domain := viper.GetString("Domain")
	if Domain == "" {
		log.Fatal("Domain should not be nil")
	}
	QnAccessKey := viper.GetString("QnAccessKey")
	if QnAccessKey == "" {
		log.Fatal("QnAccessKey should not be nil")
	}
	QnSecretKey := viper.GetString("QnSecretKey")
	if QnSecretKey == "" {
		log.Fatal("QnSecretKey should not be nil")
	}

	return &Client{
		Bucket:      Bucket,
		QnAccessKey: QnAccessKey,
		QnSecretKey: QnSecretKey,
		Domain:      Domain,
	}
}

// Upload can backup db to qiniu s3
func (c *Client) Upload() {
	file := utils.MustGetDb(DB_Name)

	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", c.Bucket, FILE_NAME),
	}
	mac := qbox.NewMac(c.QnAccessKey, c.QnSecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	cfg.Zone = &storage.ZoneHuanan
	cfg.UseHTTPS = false
	cfg.UseCdnDomains = false

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}

	err := formUploader.PutFile(context.Background(), &ret, upToken, FILE_NAME, file, nil)
	if err != nil {
		log.Fatalf("upload error, %+v\n", err)
	}
	fmt.Printf("backup success. key: %s, hash: %s \n", ret.Key, ret.Hash)
}

func (c *Client) getUrl() string {
	mac := qbox.NewMac(c.QnAccessKey, c.QnSecretKey)
	deadline := time.Now().Add(time.Second * 3600).Unix()
	privateAccessURL := storage.MakePrivateURL(mac, c.Domain, FILE_NAME, deadline)
	return privateAccessURL
}

// Download can load the db from
func (c *Client) Download() {
	ff := utils.MustGetDb(DB_Name)
	if utils.IsFileExists(ff) {
		fmt.Printf("db file exists, delete it if you want to force load. %s \n", ff)
		return
	}
	url := c.getUrl()
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	dst, err := os.OpenFile(ff, os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(dst, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	dst.Close()
	fmt.Println("sync download done!")
}
