#Scheduler-extender <br>
###predicates.go <br>
在这里实现过滤方法，对于每个node，通过podFitsOnNode()函数判断pod能否在该node上运行。这里podFitsOnNode实现的逻辑是，在10以内取随机数，若随机数小于7，则通过，否则不通过。也就是说，每个node有70%的概率通过。
```
func podFitsOnNode(pod *v1.Pod, node v1.Node) (bool, []string, error) {
	var failReasons []string
	lucky := (rand.Intn(10) < 7)
	if !lucky{
		failReasons = append(failReasons, "Pod is unlucky.")
		log.Printf("pod %v/%v does not fit on node %v\n", pod.Name, pod.Namespace, node.Name)
		return false, failReasons, nil
	}
	log.Printf("pod %v/%v fit on node %v\n", pod.Name, pod.Namespace, node.Name)
	return true, failReasons, nil
}
```
###priorities.go
在这里对每个通过了过滤的node进行打分，分数是10以内的随机数。
```
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
```
