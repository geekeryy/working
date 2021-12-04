// Package service @Description  TODO
// @Author  	 jiangyang
// @Created  	 2021/11/30 11:43 下午
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/comeonjy/go-kit/grpc/reloadconfig"
	"github.com/comeonjy/go-kit/pkg/util"
	"github.com/comeonjy/go-kit/pkg/xenv"
	"github.com/comeonjy/working/pkg/consts"
	"github.com/comeonjy/working/pkg/notify"
	"google.golang.org/grpc"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

func (svc *WorkingService) GithubEvent(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	all, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(string(all))

	defer r.Body.Close()
	resp := GithubResp{}
	if err := json.Unmarshal(all, &resp); err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(resp.Repository.Name, resp.HeadCommit.Message)

	if r.Header.Get("X-Hub-Signature-256") != "sha256="+util.HmacSha256(all, util.Md5(resp.Repository.Name)) {
		log.Println("invalid Secret")
		return
	}

	regex := regexp.MustCompile("^deploy:(v\\d{1,2}\\.\\d{1,2}\\.\\d{1,2})")
	submatch := regex.FindStringSubmatch(resp.HeadCommit.Message)
	if len(submatch) != 2 {
		log.Println("not a deploy request")
		return
	}

	tag := submatch[1]
	image := fmt.Sprintf("%s/%s:%s", consts.EnvMap["images_repo"], resp.Repository.Name, tag)
	if err := svc.restartDeploy(resp.Repository.Name, image); err != nil {
		log.Println("RestartDeploy err", err.Error())
		_ = notify.PostFieShu("RestartDeploy err: " + err.Error())
	} else {
		_ = notify.PostFieShu(fmt.Sprintf("RestartDeploy success: %s %s", resp.Repository.Name, image))
	}
}

func (svc *WorkingService) restartDeploy(name string, image string) error {
	log.Println("Updated deployment start...")
	deployments, err := svc.k8sClient.AppsV1().Deployments(apiv1.NamespaceDefault).List(context.Background(), metav1.ListOptions{
		LabelSelector: "githubRepoName=" + name,
	})
	if err != nil {
		return err
	}
	for _, result := range deployments.Items {
		if retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			result.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().String()
			result.Spec.Template.Spec.Containers[0].Image = image
			_, updateErr := svc.k8sClient.AppsV1().Deployments(apiv1.NamespaceDefault).Update(context.TODO(), &result, metav1.UpdateOptions{})
			return updateErr
		}); retryErr != nil {
			log.Println(retryErr)
			continue
		}
		log.Println("restartDeploy:",result.Name)
	}
	log.Println("Updated deployment end...")

	return nil
}

func (svc *WorkingService) ApolloEvent(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	env := r.URL.Query().Get("env")
	all, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(string(all))
	defer r.Body.Close()
	resp := ApolloEventReq{}
	if err := json.Unmarshal(all, &resp); err != nil {
		log.Println(err.Error())
		return
	}
	appId := resp.AppId

	log.Println(env, appId)

	if xenv.GetApolloCluster("") != resp.ClusterName {
		log.Println("ClusterName:", resp.ClusterName, "!=", xenv.GetApolloCluster(""))
		return
	}

	list, err := svc.k8sClient.CoreV1().Pods(apiv1.NamespaceDefault).List(context.Background(), metav1.ListOptions{
		LabelSelector: "apolloAppId=" + appId,
	})
	if err != nil {
		log.Println("K8s List err:", err)
		return
	}

	ipMap := make(map[string]struct{})
	for _, v := range list.Items {
		ipMap[v.Status.PodIP] = struct{}{}
	}

	log.Println("ipMap:", ipMap)

	for ip := range ipMap {
		dial, err := grpc.Dial(ip+":"+xenv.GetEnv(xenv.GrpcPort), grpc.WithInsecure())
		if err != nil {
			log.Println("grpc.Dial err", ip, err)
			continue
		}
		_, err = reloadconfig.NewReloadConfigClient(dial).ReloadConfig(context.Background(), &reloadconfig.Empty{})
		if err != nil {
			log.Println("NewReloadConfigClient.ReloadConfig err", ip, err)
			return
		}
	}

	if err != nil {
		log.Println(fmt.Sprintf("%s:%s err:%s", env, appId, err.Error()))
	} else {
		log.Println(fmt.Sprintf("%s:%s success", env, appId))
	}
}

type ApolloEventReq struct {
	Id                   int    `json:"id"`
	AppId                string `json:"appId"`
	ClusterName          string `json:"clusterName"`
	NamespaceName        string `json:"namespaceName"`
	BranchName           string `json:"branchName"`
	Operator             string `json:"operator"`
	ReleaseId            int    `json:"releaseId"`
	ReleaseTitle         string `json:"releaseTitle"`
	ReleaseComment       string `json:"releaseComment"`
	ReleaseTime          string `json:"releaseTime"`
	ReleaseTimeFormatted string `json:"releaseTimeFormatted"`
	Configuration        []struct {
		FirstEntity  string `json:"firstEntity"`
		SecondEntity string `json:"secondEntity"`
	} `json:"configuration"`
	IsReleaseAbandoned bool `json:"isReleaseAbandoned"`
	PreviousReleaseId  int  `json:"previousReleaseId"`
	Operation          int  `json:"operation"`
	OperationContext   struct {
		IsEmergencyPublish bool `json:"isEmergencyPublish"`
	} `json:"operationContext"`
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
