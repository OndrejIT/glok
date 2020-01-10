package geoip

import (
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	conf "github.com/spf13/viper"
)

func periodicUpdate(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for t := range ticker.C {
		log.Debug("Periodic update at: ", t)
		tryGetNewDB()
	}
}

func DBupdate() {
	update_interval := conf.GetInt("interval")
	go periodicUpdate(time.Duration(update_interval) * time.Hour)
	tryGetNewDB()
}

func tryGetNewDB() {
	for i := 1; i <= 3; i++ {
		if err := downloadNewDB(); err != nil {
			log.Errorf("[downloadNewDB] %s", err)
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
}

func downloadNewDB() error {
	log.Debug("Trying update database.")

	md5Old, err := getMD5(conf.GetString("database"))
	if err != nil {
		log.Errorf("[getMD5] %s", err)
	}

	ipAddr, err := getClientIp()
	if err != nil {
		return err
	}

	challengeMD5 := getChallengeMD5(conf.GetString("license"), ipAddr) // challengemd5 == md5sum(license+clientipaddr)
	file_name := conf.GetString("database") + ".gz"
	dl_url := fmt.Sprintf("https://updates.maxmind.com/app/update_secure?db_md5=%s&challenge_md5=%s&user_id=%s&edition_id=%s",
		md5Old, challengeMD5, conf.GetString("uid"), conf.GetString("product_id"))

	log.Debug("[downloadNewDB] Calling get request to ", dl_url)
	resp, err := http.Get(dl_url)

	if err != nil {
		return err
	}

	if header := resp.Header.Get("X-Database-Md5"); header == "" {
		return errors.New("Don't have X-Database-Md5 header.")
	}

	md5Response := resp.Header.Get("X-Database-Md5")
	log.Debugf("[downloadNewDB] X-database-md5: %s", md5Response)

	if md5Response != md5Old {
		if err = downloadToFile(resp, file_name); err != nil {
			return err
		}
		if err = gunzipAndReplace(file_name, md5Response); err != nil {
			return err
		}
	} else {
		log.Debug("No file to download")
	}

	return nil
}

func gunzipAndReplace(file_name string, md5Response string) error {
	log.Debug("[gunzipAndReplace] unziping and replacing ", file_name)
	file, err := os.Open(file_name)
	defer file.Close()
	if err != nil {
		os.Remove(file_name)
		return err
	}

	archive, err := gzip.NewReader(file)
	defer archive.Close()
	if err != nil {
		os.Remove(file_name)
		return err
	}

	writer, err := os.Create(file_name[:len(file_name)-3])
	defer writer.Close()
	if err != nil {
		os.Remove(file_name)
		return err
	}

	_, err = io.Copy(writer, archive)
	if err != nil {
		log.Error("[gunzipAndReplace] ", err)
		return err
	}
	os.Remove(file_name)

	md5New, err := getMD5(conf.GetString("database"))
	if err != nil {
		log.Errorf("[getMD5] %s", err)
		return err
	}

	if md5New != md5Response {
		err := errors.New(fmt.Sprintf("bad md5 from file %s != %s", md5New, md5Response))
		return err
	}

	return nil
}

func downloadToFile(resp *http.Response, file_name string) error {
	log.Debug("[downloadToFile] creating file ", file_name)
	file, err := os.Create(file_name)
	defer file.Close()
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	_, err = io.Copy(file, resp.Body)

	if err != nil {
		return err
	}
	return nil
}

func getChallengeMD5(license string, ipAddr string) string {
	Hash := md5.Sum([]byte(license + ipAddr))
	encodedHash := hex.EncodeToString(Hash[:16])
	log.Debug("[getChallengeMD5] License is ", license)
	log.Debug("[getChallengeMD5] Challenge md5 is ", encodedHash)
	return encodedHash
}

func getMD5(filePath string) (string, error) {
	var returnMD5 string

	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		log.Debugf("Can't open database: %s", err)
		return "00000000000000000000000000000000", nil
	}

	hash := md5.New()

	_, err = io.Copy(hash, file)
	if err != nil {
		return returnMD5, err
	}

	hashInBytes := hash.Sum(nil)
	returnMD5 = hex.EncodeToString(hashInBytes)

	log.Debugf("[getMD5] MD5 of %s is %s", filePath, returnMD5)

	return returnMD5, nil
}

func getClientIp() (string, error) {
	url := "https://updates.maxmind.com/app/update_getipaddr"

	log.Debug("[getClientIp] get request to ", url)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		log.Error("[getClientIp] Can't get client IP")
		return "", err
	}

	defer resp.Body.Close()
	log.Debug("[getClientIp] http status ", resp.Status)

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Errorf("[getClientIp] Can't get client IP")
		return "", err
	}

	bodyString := string(bodyBytes)
	log.Debug("[getClientIp] client IP is ", bodyString)

	return bodyString, nil
}
