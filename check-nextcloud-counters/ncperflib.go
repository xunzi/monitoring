package ncperflib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/kr/pretty"
)

var Debug bool = false

func Debugprint(msg string) {
	if Debug {
		fmt.Printf("DEBUG: %s\n", msg)
	}
}

type ncPerfData struct {
	Ocs struct {
		Meta struct {
			Status     string `json:"status"`
			Statuscode int    `json:"statuscode"`
			Message    string `json:"message"`
		} `json:"meta"`
		Data struct {
			Nextcloud struct {
				System struct {
					Version             string    `json:"version"`
					Theme               string    `json:"theme"`
					EnableAvatars       string    `json:"enable_avatars"`
					EnablePreviews      string    `json:"enable_previews"`
					MemcacheLocal       string    `json:"memcache.local"`
					MemcacheDistributed string    `json:"memcache.distributed"`
					FilelockingEnabled  string    `json:"filelocking.enabled"`
					MemcacheLocking     string    `json:"memcache.locking"`
					Debug               string    `json:"debug"`
					Freespace           int64     `json:"freespace"`
					Cpuload             []float64 `json:"cpuload"`
					MemTotal            int       `json:"mem_total"`
					MemFree             int       `json:"mem_free"`
					SwapTotal           int       `json:"swap_total"`
					SwapFree            int       `json:"swap_free"`
					Apps                struct {
						NumInstalled        int           `json:"num_installed"`
						NumUpdatesAvailable int           `json:"num_updates_available"`
						AppUpdates          []interface{} `json:"app_updates"`
					} `json:"apps"`
				} `json:"system"`
				Storage struct {
					NumUsers         int `json:"num_users"`
					NumFiles         int `json:"num_files"`
					NumStorages      int `json:"num_storages"`
					NumStoragesLocal int `json:"num_storages_local"`
					NumStoragesHome  int `json:"num_storages_home"`
					NumStoragesOther int `json:"num_storages_other"`
				} `json:"storage"`
				Shares struct {
					NumShares               int `json:"num_shares"`
					NumSharesUser           int `json:"num_shares_user"`
					NumSharesGroups         int `json:"num_shares_groups"`
					NumSharesLink           int `json:"num_shares_link"`
					NumSharesMail           int `json:"num_shares_mail"`
					NumSharesRoom           int `json:"num_shares_room"`
					NumSharesLinkNoPassword int `json:"num_shares_link_no_password"`
					NumFedSharesSent        int `json:"num_fed_shares_sent"`
					NumFedSharesReceived    int `json:"num_fed_shares_received"`
					Permissions31           int `json:"permissions_3_1"`
				} `json:"shares"`
			} `json:"nextcloud"`
			Server struct {
				Webserver string `json:"webserver"`
				Php       struct {
					Version           string `json:"version"`
					MemoryLimit       int    `json:"memory_limit"`
					MaxExecutionTime  int    `json:"max_execution_time"`
					UploadMaxFilesize int    `json:"upload_max_filesize"`
				} `json:"php"`
				Database struct {
					Type    string `json:"type"`
					Version string `json:"version"`
					Size    int    `json:"size"`
				} `json:"database"`
			} `json:"server"`
			ActiveUsers struct {
				Last5Minutes int `json:"last5minutes"`
				Last1Hour    int `json:"last1hour"`
				Last24Hours  int `json:"last24hours"`
			} `json:"activeUsers"`
		} `json:"data"`
	} `json:"ocs"`
}

func CheckArguments(counter string, warning int64, critical int64) {
	if warning >= critical {
		NagiosResult(3, "Warning must be smaller than Critical")
	}
	allowedCounters := []string{"AppUdatesAvailable", "FreeSpace", "NumShares", "ActiveUsers5Min", "DbSize"}
	for c := range allowedCounters {
		if allowedCounters[c] == counter {
			return
		}
	}
	NagiosResult(3, "Unknown Counter")
}

func NagiosResult(ret int, message string) {
	switch ret {
	case 0:
		fmt.Printf("OK: %s\n", message)

		os.Exit(ret)
	case 1:
		fmt.Printf("WARNING: %s \n", message)
		os.Exit(ret)
	case 2:
		fmt.Printf("CRITICAL: %s\n", message)
		os.Exit(ret)
	default:
		fmt.Printf("UNKNOWN: %s\n", message)
		os.Exit(3)
	}
}

func FetchPerformanceInfo(counter string, hostname string, uri string, username string, password string) []byte {
	perfURL := fmt.Sprintf("https://%s/%s?format=json", hostname, uri)
	Debugprint(fmt.Sprintf("initiating GET request to %s", perfURL))
	req, err := http.NewRequest("GET", perfURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(username, password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode > 299 {
		NagiosResult(3, fmt.Sprintf("Http request returned %d (%s)", resp.StatusCode, http.StatusText(resp.StatusCode)))
	} else {
		Debugprint(fmt.Sprintf("Status %s:", resp.Status))
		Debugprint(fmt.Sprintf("%# v", pretty.Formatter(resp.Header)))
	}

	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	Debugprint(string(b))

	return b
}

func ParsePerfData(pdata []byte, counter string) int64 {
	var m ncPerfData
	json.Unmarshal(pdata, &m)

	switch counter {
	case "AppUdatesAvailable":
		return int64(m.Ocs.Data.Nextcloud.System.Apps.NumUpdatesAvailable)
	case "FreeSpace":
		return m.Ocs.Data.Nextcloud.System.Freespace
	case "NumShares":
		return int64(m.Ocs.Data.Nextcloud.Shares.NumShares)
	case "ActiveUsers5Min":
		return int64(m.Ocs.Data.ActiveUsers.Last5Minutes)
	case "DbSize":
		return int64(m.Ocs.Data.Server.Database.Size)
	default:
		return -1
	}

}
