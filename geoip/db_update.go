package geoip

import (
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	conf "github.com/spf13/viper"
)

const zeroMd5 = "00000000000000000000000000000000"

func periodicUpdate(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for t := range ticker.C {
		log.Debug("Periodic update at: ", t)
		tryGetNewDB()
	}
}

func DBupdate() {
	updateInterval := conf.GetInt("interval")
	go periodicUpdate(time.Duration(updateInterval) * time.Hour)

	if stats := tryGetNewDB(); stats == false {
		md5Old, _ := getMD5(conf.GetString("database"))
		if md5Old != zeroMd5 {
			OpenDatabase()
		} else {
			log.Fatalf("[DBupdate] Can't read database file")
		}
	}
}

func tryGetNewDB() bool {
	for i := 1; i <= 3; i++ {
		if stats, err := downloadNewDB(); err != nil {
			log.Errorf("[downloadNewDB] %s", err)
			time.Sleep(5 * time.Second)
		} else {
			return stats
		}
	}

	return false
}

// https://github.com/maxmind/geoipupdate/blob/4f30969a125b55f986382d1d3b4f9534929fd5b5/pkg/geoipupdate/database/http_reader.go#L43
func downloadNewDB() (bool, error) {
	log.Debug("Trying update database.")

	md5Old, err := getMD5(conf.GetString("database"))
	if err != nil {
		log.Errorf("[getMD5] %s", err)
	}

	file_name := conf.GetString("database") + ".gz"

	maxMindURL := fmt.Sprintf(
		"https://updates.maxmind.com/geoip/databases/%s/update?db_md5=%s",
		conf.GetString("product_id"),
		url.QueryEscape(md5Old),
	)

	req, err := http.NewRequest(http.MethodGet, maxMindURL, nil)
	if err != nil {
		return false, err
	}

	req.SetBasicAuth(conf.GetString("uid"), conf.GetString("license"))

	log.Debug("[downloadNewDB] Calling get request to ", maxMindURL)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	if header := resp.Header.Get("X-Database-Md5"); header == "" {
		return false, errors.New("Don't have X-Database-Md5 header")
	}

	md5Response := resp.Header.Get("X-Database-Md5")
	log.Debugf("[downloadNewDB] X-database-md5: %s", md5Response)

	if md5Response != md5Old {
		if err = downloadToFile(resp, file_name); err != nil {
			return false, err
		}
		if err = gunzipAndReplace(file_name, md5Response); err != nil {
			return false, err
		}
		OpenDatabase()
		return true, nil
	}

	log.Debug("No file to download")
	return false, nil
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
		err := fmt.Errorf("bad md5 from file %s != %s", md5New, md5Response)
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

func getMD5(filePath string) (string, error) {
	var returnMD5 string

	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		log.Debugf("Can't open database: %s", err)
		return zeroMd5, nil
	}

	hash := md5.New()

	_, err = io.Copy(hash, file)
	if err != nil {
		return zeroMd5, err
	}

	hashInBytes := hash.Sum(nil)
	returnMD5 = hex.EncodeToString(hashInBytes)

	log.Debugf("[getMD5] MD5 of %s is %s", filePath, returnMD5)

	return returnMD5, nil
}
