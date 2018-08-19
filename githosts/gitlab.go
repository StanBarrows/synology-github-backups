package githosts

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type gitlabHost struct {
	Provider string
	APIURL   string
}

func (provider gitlabHost) getAuthenticatedGitlabUserID(client http.Client) int {
	type gitLabNameResponse struct {
		ID int
	}
	// get user id
	getUserIDURL := provider.APIURL + string(os.PathSeparator) + "user"
	req, newReqErr := http.NewRequest(http.MethodGet, getUserIDURL, nil)
	if newReqErr != nil {
		logger.Fatal(newReqErr)
	}
	req.Header.Set("Private-Token", os.Getenv("GITLAB_TOKEN"))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	resp, reqErr := client.Do(req)
	if reqErr != nil {
		logger.Fatal(reqErr)
	}
	bodyB, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(bytes.Replace(bodyB, []byte("\r"), []byte("\r\n"), -1))
	var respObj gitLabNameResponse
	if err := json.Unmarshal([]byte(bodyStr), &respObj); err != nil {
		logger.Fatal(err)
	}
	return respObj.ID
}

type gitLabOwner struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type gitLabProject struct {
	Path              string      `json:"path"`
	PathWithNameSpace string      `json:"path_with_namespace"`
	HTTPSURL          string      `json:"http_url_to_repo"`
	SSHURL            string      `json:"ssh_url_to_repo"`
	Owner             gitLabOwner `json:"owner"`
}
type gitLabGetProjectsResponse []gitLabProject

func (provider gitlabHost) getProjectsByUserID(client http.Client, userID int) (repos []repository) {
	getUserIDURL := provider.APIURL + string(os.PathSeparator) + "users" + string(os.PathSeparator) + strconv.Itoa(userID) + string(os.PathSeparator) + "projects"
	req, _ := http.NewRequest(http.MethodGet, getUserIDURL, nil)
	req.Header.Set("Private-Token", os.Getenv("GITLAB_TOKEN"))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	resp, _ := client.Do(req)
	bodyB, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(bytes.Replace(bodyB, []byte("\r"), []byte("\r\n"), -1))
	var respObj gitLabGetProjectsResponse
	if err := json.Unmarshal([]byte(bodyStr), &respObj); err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}

	for _, project := range respObj {
		// gitlab replaces hyphens with spaces in owner names, so fix
		var owner string
		owner = strings.Replace(project.Owner.Name, " ", "-", -1)
		var repo = repository{
			Name:          project.Path,
			Owner:         owner,
			NameWithOwner: project.PathWithNameSpace,
			HTTPSUrl:      project.HTTPSURL,
			SSHUrl:        project.SSHURL,
			Domain:        "gitlab.com",
		}
		repos = append(repos, repo)
	}
	return repos
}

func (provider gitlabHost) describeRepos() describeReposOutput {
	logger.Println("listing GitLab repositories")
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	userID := provider.getAuthenticatedGitlabUserID(*client)

	result := describeReposOutput{
		Repos: provider.getProjectsByUserID(*client, userID),
	}
	return result
}

func (provider gitlabHost) getAPIURL() string {
	return provider.APIURL
}

func (provider gitlabHost) Backup(backupDIR string) {
	describeReposOutput := provider.describeRepos()

	for _, repo := range describeReposOutput.Repos {
		firstPos := strings.Index(repo.HTTPSUrl, "//")
		logger.Println("owner: ", repo.Owner)
		repo.URLWithToken = repo.HTTPSUrl[:firstPos+2] + repo.Owner + ":" + os.Getenv("GITLAB_TOKEN") + "@" + repo.HTTPSUrl[firstPos+2:]
		processBackup(repo, backupDIR)
	}
}
