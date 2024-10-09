package main

import (
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Job struct {
	Class       string `json:"_class"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Jobs        []Job  `json:"jobs"`
	description string
}

type JobsResponse struct {
	Jobs        []Job `json:"jobs"`
	description string
}

var res = make([]Result, 0)

// 从jenkins页面F12 Network复制
var cookie = ""

func main() {
	jenkinsURL := "https://jenkins.osinfra.cn/api/json"
	client := &http.Client{}
	req, err := http.NewRequest("GET", jenkinsURL, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("cookie", cookie)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var jobsResponse JobsResponse
	err = json.Unmarshal(body, &jobsResponse)
	if err != nil {
		fmt.Println(err)
		return
	}

	getAllNonFolderJobs(jobsResponse.Jobs, client, req)

	// 创建一个新的Excel文件
	f := excelize.NewFile()
	sheetName := "Sheet1"

	// 设置表头
	headers := []string{"product", "jobName", "description", "repoUrl", "branch", "maintainerEmail", "jobUrl"}
	for i, header := range headers {
		cellName, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cellName, header)
	}

	// 填充数据
	for i, ac := range res {
		cellName, _ := excelize.CoordinatesToCellName(1, i+2)
		f.SetCellValue(sheetName, cellName, ac.product)
		cellName, _ = excelize.CoordinatesToCellName(2, i+2)
		f.SetCellValue(sheetName, cellName, ac.jobName)
		cellName, _ = excelize.CoordinatesToCellName(3, i+2)
		f.SetCellValue(sheetName, cellName, ac.description)
		cellName, _ = excelize.CoordinatesToCellName(4, i+2)
		f.SetCellValue(sheetName, cellName, ac.repoUrl)
		cellName, _ = excelize.CoordinatesToCellName(5, i+2)
		f.SetCellValue(sheetName, cellName, ac.branch)
		cellName, _ = excelize.CoordinatesToCellName(6, i+2)
		f.SetCellValue(sheetName, cellName, ac.maintainerEmail)
		cellName, _ = excelize.CoordinatesToCellName(7, i+2)
		f.SetCellValue(sheetName, cellName, ac.jobUrl)
		cellName, _ = excelize.CoordinatesToCellName(8, i+2)
		f.SetCellValue(sheetName, cellName, ac.serviceName)
	}

	// 保存Excel文件
	if err := f.SaveAs("people" + strconv.FormatInt(time.Now().UnixMilli(), 10) + ".xlsx"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Excel文件已成功创建：people.xlsx")
	}
}

func getAllNonFolderJobs(jobs []Job, client *http.Client, req *http.Request) {
	for _, job := range jobs {
		if job.Class == "com.cloudbees.hudson.plugins.folder.Folder" {
			getFolderJobs(job.URL+"api/json", client, req)
		} else if job.Class == "org.jenkinsci.plugins.workflow.job.WorkflowJob" {
			continue
		} else {
			fmt.Println("Job Type:", job.Class)
			fmt.Println("Job Name:", job.Name)
			fmt.Println("Job URL:", job.URL)
			if job.Jobs != nil {
				getFolderJobs(job.URL+"api/json", client, req)
			} else {
				if job.Name == "mindspore-prod-usercenter" {
					fmt.Println(111)
				}
				getLastBuildInfo(job, client, req)
			}
		}
	}
}

func getLastBuildInfo(job Job, client *http.Client, req *http.Request) {
	req, err := http.NewRequest("GET", job.URL+"lastSuccessfulBuild/api/json", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("cookie", cookie)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 404 {
		return
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	var lastBuildInfo MatrixBuild
	err = json.Unmarshal(body, &lastBuildInfo)
	if err != nil {
		log.Println(err)
	}

	result := Result{}

	for _, action := range lastBuildInfo.Actions {
		if action.Class == "hudson.model.ParametersAction" {
			for _, parameter := range action.Parameters {
				if parameter.Name == "branch" {
					result.branch = parameter.Value
				}
				if parameter.Name == "release" {
					result.branch = parameter.Value
				}
				if parameter.Name == "REPOSITORY" || parameter.Name == "REPO" || parameter.Name == "CODE_REPOSITORY" || parameter.Name == "REPOTISOTY" {
					result.repoUrl = parameter.Value
				}
				if parameter.Name == "CODE_BRANCH" {
					result.branch = parameter.Value
				}
				if parameter.Name == "PROJECT_NAME" {
					result.jobName = parameter.Value
				}
				if parameter.Name == "EMAIL" {
					result.maintainerEmail = parameter.Value
				}
			}
		}
		if action.Class == "hudson.plugins.git.util.BuildData" && len(action.RemoteUrls) > 0 && len(action.LastBuiltRevision.Branch) > 0 {
			if result.repoUrl == "" {
				result.repoUrl = action.RemoteUrls[0]
			}
			result.branch = strings.ReplaceAll(action.LastBuiltRevision.Branch[0].Name, "refs/remotes/origin/", "")
		}
	}
	if lastBuildInfo.URL == "" {
		return
	}
	split := strings.Split(lastBuildInfo.URL, "/")
	result.description = lastBuildInfo.FullDisplayName
	result.product = split[4]
	result.jobUrl = lastBuildInfo.URL
	if result.jobName == "" {
		result.jobName = split[len(split)-3]
	}
	if result.repoUrl != "" {
		repoUrlSplit := strings.Split(result.repoUrl, "/")
		result.jobName = strings.Replace(repoUrlSplit[len(repoUrlSplit)-1], ".git", "", 1)
		result.serviceName = result.jobName
	}
	res = append(res, result)
}

func getFolderJobs(folderURL string, client *http.Client, req *http.Request) {
	req, err := http.NewRequest("GET", folderURL, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("cookie", cookie)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var folderResponse JobsResponse
	err = json.Unmarshal(body, &folderResponse)
	if err != nil {
		fmt.Println(err)
		return
	}

	getAllNonFolderJobs(folderResponse.Jobs, client, req)
}

type Result struct {
	product         string
	description     string
	jobName         string
	repoUrl         string
	branch          string
	maintainerEmail string
	jobUrl          string
	serviceName     string
}

type Cause struct {
	Class            string `json:"_class"`
	ShortDescription string `json:"shortDescription"`
	UpstreamBuild    int    `json:"upstreamBuild"`
	UpstreamProject  string `json:"upstreamProject"`
	UpstreamUrl      string `json:"upstreamUrl"`
}

type Parameter struct {
	Class string `json:"_class"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Action struct {
	Class              string      `json:"_class"`
	Causes             []Cause     `json:"causes,omitempty"`
	Parameters         []Parameter `json:"parameters,omitempty"`
	BuildsByBranchName BuildData   `json:"buildsByBranchName,omitempty"`
	RemoteUrls         []string    `json:"remoteUrls"`
	LastBuiltRevision  struct {
		SHA1   string   `json:"SHA1"`
		Branch []Branch `json:"branch"`
	} `json:"lastBuiltRevision"`
}

type Branch struct {
	SHA1 string `json:"SHA1"`
	Name string `json:"name"`
}

type Build struct {
	Class  string `json:"_class"`
	Number int    `json:"number"`
	Result string `json:"result"`
	Marked struct {
		SHA1   string   `json:"SHA1"`
		Branch []Branch `json:"branch"`
	} `json:"marked"`
	Revision struct {
		SHA1   string   `json:"SHA1"`
		Branch []Branch `json:"branch"`
	} `json:"revision"`
}

type BuildData struct {
	Class              string           `json:"_class"`
	BuildsByBranchName map[string]Build `json:"buildsByBranchName"`
}

type Run struct {
	Number int    `json:"number"`
	URL    string `json:"url"`
}

type MatrixBuild struct {
	Class             string   `json:"_class"`
	Actions           []Action `json:"actions"`
	Building          bool     `json:"building"`
	Description       string   `json:"description"`
	DisplayName       string   `json:"displayName"`
	Duration          int      `json:"duration"`
	EstimatedDuration int      `json:"estimatedDuration"`
	Executor          string   `json:"executor"`
	FullDisplayName   string   `json:"fullDisplayName"`
	ID                string   `json:"id"`
	InProgress        bool     `json:"inProgress"`
	KeepLog           bool     `json:"keepLog"`
	Number            int      `json:"number"`
	QueueId           int      `json:"queueId"`
	Result            string   `json:"result"`
	Timestamp         int64    `json:"timestamp"`
	URL               string   `json:"url"`
	BuiltOn           string   `json:"builtOn"`
	Runs              []Run    `json:"runs"`
}
