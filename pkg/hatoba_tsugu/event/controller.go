package event

// code base:
// https://github.com/kubernetes/client-go/blob/master/examples/workqueue/main.go
// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-api-machinery/controllers.md

import (
	"fmt"
	go_tool "github.com/maxiloEmmmm/go-tool"
	"hatoba_tsugu/pkg/channel"
	"html/template"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"strings"

	"time"
)

// Controller demonstrates how to implement a controller with client-go.
type Controller struct {
	queue    workqueue.RateLimitingInterface
	informer cache.SharedIndexInformer
}

// NewController creates a new Controller.
func NewController(queue workqueue.RateLimitingInterface, informer cache.SharedIndexInformer) *Controller {
	c := &Controller{
		informer: informer,
		queue:    queue,
	}
	c.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.onAdd,
		UpdateFunc: c.onUpdate,
		DeleteFunc: c.onDelete,
	})
	return c
}

func (c *Controller) onAdd(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err == nil {
		c.queue.Add(key)
	}
}

func (c *Controller) onUpdate(oldObj interface{}, newObj interface{}) {
	//key, err := cache.MetaNamespaceKeyFunc(newObj)
	//if err == nil {
	//	c.queue.Add(key)
	//}
}

func (c *Controller) onDelete(obj interface{}) {
	//key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	//if err == nil {
	//	c.queue.Add(key)
	//}
}

func (c *Controller) processNextItem() bool {
	// Wait until there is a new item in the working queue
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	// Tell the queue that we are done with processing this key. This unblocks the key for other workers
	// This allows safe parallel processing because two pods with the same key are never processed in
	// parallel.
	defer c.queue.Done(key)

	err := c.processItem(key.(string))
	if err == nil {
		c.queue.Forget(key)
	} else if c.queue.NumRequeues(key) < 3 {
		c.queue.AddRateLimited(key)
	} else {
		c.queue.Forget(key)
		runtime.HandleError(err)
	}
	return true
}

func (c *Controller) processItem(key string) error {
	obj, _, err := c.informer.GetIndexer().GetByKey(key)
	if err != nil {
		return err
	}

	evt := obj.(*v1.Event)
	for _, notify := range notifications {
		tpl, err := template.New("tpl").Parse(notify.Spec.Tpl)
		if err != nil {
			continue
		}

		// 挨个过滤
		var result = false
		for _, filter := range notify.Spec.Filter {
			val, _ := go_tool.Get(evt, filter.Key)
			if str, ok := val.(string); ok {
				// 具体规则
				switch filter.Type {
				case EqFilterType:
					result = str == filter.Val
				case InFilterType:
					result = go_tool.InArray(strings.Split(filter.Val, ","), str)
				}

				if !result {
					break
				}
			}
		}

		if result {
			builder := &strings.Builder{}
			err := tpl.Execute(builder, evt)
			if err != nil {
				continue
			}

			if c := channel.NewChannel(notify.Spec.Engine); c != nil {
				c.Send(builder.String())
			}
			break
		}
	}
	return nil
}

// handleErr checks if an error happened and makes sure we will retry later.
func (c *Controller) handleErr(err error, key interface{}) {
	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization.
		// This ensures that future processing of updates for this key is not delayed because of
		// an outdated error history.
		c.queue.Forget(key)
		return
	}

	// This controller retries 5 times if something goes wrong. After that, it stops trying.
	if c.queue.NumRequeues(key) < 5 {
		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	// Report to an external entity that, even after several retries, we could not successfully process this key
	runtime.HandleError(err)
}

// Run begins watching and syncing.
func (c *Controller) Run(threadiness int, stopCh chan struct{}) {
	defer runtime.HandleCrash()

	// Let the workers stop when we are done
	defer c.queue.ShutDown()

	go c.informer.Run(stopCh)

	// Wait for all involved caches to be synced, before processing items from the queue is started
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}
