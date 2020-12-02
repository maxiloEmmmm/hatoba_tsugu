package event

import (
	"hatoba_tsugu/pkg/kubernetes"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"log"
)

var notifications []Notification

func FetchEventNotification() {
	response, err := kubernetes.Client.R().SetResult(&NotificationListCrd{}).Get(kubernetes.HatobaTsuguEventNotificationPath.MultiPath())
	if err != nil {
		log.Fatal(err)
	} else if err = kubernetes.ResponseOk(response); err != nil {
		log.Fatal(err)
	}

	notifications = response.Result().(*NotificationListCrd).Items
}

func Init(stopChan chan struct{}) {
	FetchEventNotification()

	controller := NewController(
		workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		cache.NewSharedIndexInformer(cache.NewListWatchFromClient(kubernetes.KubeClient.CoreV1().RESTClient(), "events", v1.NamespaceAll, fields.Everything()), &v1.Event{}, 0, cache.Indexers{}),
	)
	go controller.Run(1, stopChan)
}

type NotificationCrd struct {
	metaV1.TypeMeta   `json:",omitempty"`
	metaV1.ObjectMeta `json:"metadata,omitempty"`
	Spec              Notification `json:"spec,omitempty"`
}

type NotificationListCrd struct {
	v1.ServiceList
	metaV1.TypeMeta `json:",inline"`
	metaV1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []Notification `json:"spec,omitempty"`
}

type Notification struct {
	Name   string               `json:"name"`
	Tpl    string               `json:"tpl"`
	Engine string               `json:"engine"`
	Filter []NotificationFilter `json:"notify_filter"`
}

type NotificationFilter struct {
	Type string `json:"type"`
	Key  string `json:"key"`
	Val  string `json:"val"`
}
