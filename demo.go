/*
Copyright 2016 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gorilla/mux"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func demo() {

	r := mux.NewRouter()
	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hello!"))
	})
	r.HandleFunc("/github-event/{repo}", GithubEvent)
	r.HandleFunc("/apollo-event", ApolloEvent)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalln("err", err)
	}

}
func hmacSha256(data []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func ApolloEvent(writer http.ResponseWriter, request *http.Request) {
	all, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(string(all))
}
func GithubEvent(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	repo := vars["repo"]
	all, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(string(all))

	if request.Header.Get("X-Hub-Signature-256")!="sha256="+hmacSha256(all,Md5(repo)){
		log.Println("invalid Secret")
		return
	}
	defer request.Body.Close()
	resp := GithubResp{}
	if err := json.Unmarshal(all, &resp); err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(resp.Repository.Name, resp.HeadCommit.Message)
	regex := regexp.MustCompile("^deploy:(v\\d{1,2}\\.\\d{1,2}\\.\\d{1,2})")
	submatch := regex.FindStringSubmatch(resp.HeadCommit.Message)
	if len(submatch) != 2 {
		log.Println("not a deploy request")
		return
	}
	if resp.Repository.Name != repo {
		log.Println("invalid servername:", resp.Repository.Name, repo)
		return
	}
	tag := submatch[1]
	image := fmt.Sprintf("ccr.ccs.tencentyun.com/comeonjy/%s:%s", repo, tag)
	if err := RestartDeploy(repo, image); err != nil {
		log.Println("RestartDeploy err", err.Error())
		_ = PostFieShu("RestartDeploy err: " + err.Error())
	} else {
		_ = PostFieShu(fmt.Sprintf("RestartDeploy success: %s %s", repo, image))
	}
}

func RestartDeploy(name string, image string) error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	log.Println("Updated deployment start...")
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := deploymentsClient.Get(context.TODO(), name, metav1.GetOptions{})
		if getErr != nil {
			return err
		}
		result.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().String()
		result.Spec.Template.Spec.Containers[0].Image = image
		_, updateErr := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		return retryErr
	}
	log.Println("Updated deployment end...")

	return nil
}

func Md5(str string) string  {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func PostFieShu(info string) error {
	msg := FeiShuMsg{
		MsgType: "text",
	}
	msg.Content.Text = fmt.Sprintf("working [ %s ] : %s", os.Getenv("APP_ENV"), info)
	marshal, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", "https://open.feishu.cn/open-apis/bot/v2/hook/67c37caa-a7c2-44b3-8726-b784081c2102", bytes.NewReader(marshal))
	if err != nil {
		return err
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return nil
}

type FeiShuMsg struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

type GithubResp struct {
	Ref        string `json:"ref"`
	Before     string `json:"before"`
	After      string `json:"after"`
	Repository struct {
		Id       int    `json:"id"`
		NodeId   string `json:"node_id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Private  bool   `json:"private"`
		Owner    struct {
			Name              string `json:"name"`
			Email             string `json:"email"`
			Login             string `json:"login"`
			Id                int    `json:"id"`
			NodeId            string `json:"node_id"`
			AvatarUrl         string `json:"avatar_url"`
			GravatarId        string `json:"gravatar_id"`
			Url               string `json:"url"`
			HtmlUrl           string `json:"html_url"`
			FollowersUrl      string `json:"followers_url"`
			FollowingUrl      string `json:"following_url"`
			GistsUrl          string `json:"gists_url"`
			StarredUrl        string `json:"starred_url"`
			SubscriptionsUrl  string `json:"subscriptions_url"`
			OrganizationsUrl  string `json:"organizations_url"`
			ReposUrl          string `json:"repos_url"`
			EventsUrl         string `json:"events_url"`
			ReceivedEventsUrl string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"owner"`
		HtmlUrl          string        `json:"html_url"`
		Description      string        `json:"description"`
		Fork             bool          `json:"fork"`
		Url              string        `json:"url"`
		ForksUrl         string        `json:"forks_url"`
		KeysUrl          string        `json:"keys_url"`
		CollaboratorsUrl string        `json:"collaborators_url"`
		TeamsUrl         string        `json:"teams_url"`
		HooksUrl         string        `json:"hooks_url"`
		IssueEventsUrl   string        `json:"issue_events_url"`
		EventsUrl        string        `json:"events_url"`
		AssigneesUrl     string        `json:"assignees_url"`
		BranchesUrl      string        `json:"branches_url"`
		TagsUrl          string        `json:"tags_url"`
		BlobsUrl         string        `json:"blobs_url"`
		GitTagsUrl       string        `json:"git_tags_url"`
		GitRefsUrl       string        `json:"git_refs_url"`
		TreesUrl         string        `json:"trees_url"`
		StatusesUrl      string        `json:"statuses_url"`
		LanguagesUrl     string        `json:"languages_url"`
		StargazersUrl    string        `json:"stargazers_url"`
		ContributorsUrl  string        `json:"contributors_url"`
		SubscribersUrl   string        `json:"subscribers_url"`
		SubscriptionUrl  string        `json:"subscription_url"`
		CommitsUrl       string        `json:"commits_url"`
		GitCommitsUrl    string        `json:"git_commits_url"`
		CommentsUrl      string        `json:"comments_url"`
		IssueCommentUrl  string        `json:"issue_comment_url"`
		ContentsUrl      string        `json:"contents_url"`
		CompareUrl       string        `json:"compare_url"`
		MergesUrl        string        `json:"merges_url"`
		ArchiveUrl       string        `json:"archive_url"`
		DownloadsUrl     string        `json:"downloads_url"`
		IssuesUrl        string        `json:"issues_url"`
		PullsUrl         string        `json:"pulls_url"`
		MilestonesUrl    string        `json:"milestones_url"`
		NotificationsUrl string        `json:"notifications_url"`
		LabelsUrl        string        `json:"labels_url"`
		ReleasesUrl      string        `json:"releases_url"`
		DeploymentsUrl   string        `json:"deployments_url"`
		CreatedAt        int           `json:"created_at"`
		UpdatedAt        time.Time     `json:"updated_at"`
		PushedAt         int           `json:"pushed_at"`
		GitUrl           string        `json:"git_url"`
		SshUrl           string        `json:"ssh_url"`
		CloneUrl         string        `json:"clone_url"`
		SvnUrl           string        `json:"svn_url"`
		Homepage         interface{}   `json:"homepage"`
		Size             int           `json:"size"`
		StargazersCount  int           `json:"stargazers_count"`
		WatchersCount    int           `json:"watchers_count"`
		Language         interface{}   `json:"language"`
		HasIssues        bool          `json:"has_issues"`
		HasProjects      bool          `json:"has_projects"`
		HasDownloads     bool          `json:"has_downloads"`
		HasWiki          bool          `json:"has_wiki"`
		HasPages         bool          `json:"has_pages"`
		ForksCount       int           `json:"forks_count"`
		MirrorUrl        interface{}   `json:"mirror_url"`
		Archived         bool          `json:"archived"`
		Disabled         bool          `json:"disabled"`
		OpenIssuesCount  int           `json:"open_issues_count"`
		License          interface{}   `json:"license"`
		AllowForking     bool          `json:"allow_forking"`
		IsTemplate       bool          `json:"is_template"`
		Topics           []interface{} `json:"topics"`
		Visibility       string        `json:"visibility"`
		Forks            int           `json:"forks"`
		OpenIssues       int           `json:"open_issues"`
		Watchers         int           `json:"watchers"`
		DefaultBranch    string        `json:"default_branch"`
		Stargazers       int           `json:"stargazers"`
		MasterBranch     string        `json:"master_branch"`
	} `json:"repository"`
	Pusher struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"pusher"`
	Sender struct {
		Login             string `json:"login"`
		Id                int    `json:"id"`
		NodeId            string `json:"node_id"`
		AvatarUrl         string `json:"avatar_url"`
		GravatarId        string `json:"gravatar_id"`
		Url               string `json:"url"`
		HtmlUrl           string `json:"html_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReposUrl          string `json:"repos_url"`
		EventsUrl         string `json:"events_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"sender"`
	Created bool        `json:"created"`
	Deleted bool        `json:"deleted"`
	Forced  bool        `json:"forced"`
	BaseRef interface{} `json:"base_ref"`
	Compare string      `json:"compare"`
	Commits []struct {
		Id        string    `json:"id"`
		TreeId    string    `json:"tree_id"`
		Distinct  bool      `json:"distinct"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
		Url       string    `json:"url"`
		Author    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
		Committer struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"committer"`
		Added    []string      `json:"added"`
		Removed  []interface{} `json:"removed"`
		Modified []interface{} `json:"modified"`
	} `json:"commits"`
	HeadCommit struct {
		Id        string    `json:"id"`
		TreeId    string    `json:"tree_id"`
		Distinct  bool      `json:"distinct"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
		Url       string    `json:"url"`
		Author    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
		Committer struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"committer"`
		Added    []string      `json:"added"`
		Removed  []interface{} `json:"removed"`
		Modified []interface{} `json:"modified"`
	} `json:"head_commit"`
}
