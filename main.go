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
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

type T struct {
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

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		log.Println("/")
	})
	http.HandleFunc("/github-event/account", func(writer http.ResponseWriter, request *http.Request) {
		log.Println("/github-event/account")
		all, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return
		}
		defer request.Body.Close()
		log.Println(string(all))
	})
	if err:=http.ListenAndServe(":80", nil);err!=nil{
		log.Fatalln("err",err)
	}
}

func devops() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		// get pods in all the namespaces by omitting namespace
		// Or specify namespace to get pods in particular namespace
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		// Examples for error handling:
		// - Use helper functions e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		_, err = clientset.CoreV1().Pods("default").Get(context.TODO(), "example-xxxxx", metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod example-xxxxx not found in default namespace\n")
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found example-xxxxx pod in default namespace\n")
		}

		deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			// Retrieve the latest version of Deployment before attempting update
			// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
			result, getErr := deploymentsClient.Get(context.TODO(), "account", metav1.GetOptions{})
			if getErr != nil {
				panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
			}
			result.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().String()
			result.Spec.Template.Spec.Containers[0].Image = "ccr.ccs.tencentyun.com/comeonjy/account:v0.0.1"
			_, updateErr := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
			return updateErr
		})
		if retryErr != nil {
			panic(fmt.Errorf("Update failed: %v", retryErr))
		}
		fmt.Println("Updated deployment...")

		time.Sleep(20 * time.Second)
	}
}
func int32Ptr(i int32) *int32 { return &i }
