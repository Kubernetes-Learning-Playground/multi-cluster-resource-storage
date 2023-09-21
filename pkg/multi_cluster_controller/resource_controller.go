package multi_cluster_controller

import (
	"context"
	"fmt"
	"github.com/practice/multi_resource/pkg/apis/resource/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

func (mc *MultiClusterHandler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	// 获取 Resource
	rr := &v1alpha1.MultiClusterResource{}
	err := mc.Get(ctx, req.NamespacedName, rr)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// 删除状态，会等到 Finalizer 字段清空后才会真正删除
	// 1、删除所有集群资源
	// 2、清空 Finalizer，更新状态
	if !rr.DeletionTimestamp.IsZero() {
		err = mc.resourceDelete(rr)
		if err != nil {
			mc.EventRecorder.Event(rr, corev1.EventTypeNormal, "Delete", fmt.Sprintf("delete %s fail", rr.Name))
			return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 60}, err
		}

		err = mc.Client.Update(ctx, rr)
		if err != nil {
			mc.EventRecorder.Event(rr, corev1.EventTypeWarning, "UpdateFailed", fmt.Sprintf("update %s fail", rr.Name))
			return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 60}, err
		}

		return reconcile.Result{}, nil
	}

	// 设置 crd 对象的 Finalizer 字段，并判断是否改变
	forDelete, finalizer, isChange := mc.setResourceFinalizer(rr)

	// 如果 Finalizer 字段改变，
	// 代表可能是需要进行特定集群的删除资源操作
	if isChange {
		err = mc.resourceDeleteBySlice(rr, forDelete)
		if err != nil {
			mc.EventRecorder.Event(rr, corev1.EventTypeWarning, "DeleteFailed", fmt.Sprintf("resourceDeleteBySlice %s fail", rr.Name))
			return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 60}, err
		}
		// 删除后覆盖
		rr.Finalizers = finalizer
		err = mc.Client.Update(ctx, rr)
		if err != nil {
			mc.EventRecorder.Event(rr, corev1.EventTypeWarning, "UpdateFailed", fmt.Sprintf("update %s fail", rr.Name))
			return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 60}, err
		}
	}

	// apply 操作
	err = mc.resourceApply(rr)
	if err != nil {
		mc.EventRecorder.Event(rr, corev1.EventTypeWarning, "ApplyFailed", fmt.Sprintf("resourceApply %s fail", rr.Name))
		return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 60}, err
	}

	return reconcile.Result{}, nil
}