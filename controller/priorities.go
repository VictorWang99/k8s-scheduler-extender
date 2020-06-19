package controller

import (
	"log"
	"math/rand"

	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
)

// It'd better to only define one custom priority per extender
// as current extender interface only supports one single weight mapped to one extender
// and also it returns HostPriorityList, rather than []HostPriorityList


// it's webhooked to pkg/scheduler/core/generic_scheduler.go#PrioritizeNodes()
// you can't see existing scores calculated so far by default scheduler
// instead, scores output by this function will be added back to default scheduler
func prioritize(args schedulerapi.ExtenderArgs) *schedulerapi.HostPriorityList {
	pod := args.Pod
	nodes := args.Nodes.Items

	hostPriorityList := make(schedulerapi.HostPriorityList, len(nodes))
	for i, node := range nodes {
		score := rand.Intn(i*10)
		if score > schedulerapi.MaxPriority{
			score = schedulerapi.MaxPriority
		}
		log.Printf("pod %v/%v is lucky to get score %v\n", pod.Name, pod.Namespace, score)
		hostPriorityList[i] = schedulerapi.HostPriority{
			Host:  node.Name,
			Score: score,
		}
	}

	return &hostPriorityList
}
