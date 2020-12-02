package event

import (
	"hatoba_tsugu/pkg/kubernetes"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

func Init(stopChan chan struct{}) {
	controller := NewController(
		workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		cache.NewSharedIndexInformer(cache.NewListWatchFromClient(kubernetes.KubeClient.CoreV1().RESTClient(), "events", v1.NamespaceAll, fields.Everything()), &v1.Event{}, 0, cache.Indexers{}),
	)
	go controller.Run(1, stopChan)
}
