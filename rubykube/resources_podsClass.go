package rubykube

import (
	mruby "github.com/mitchellh/go-mruby"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	kapi "k8s.io/client-go/pkg/api/v1"
)

type podListTypeAlias kapi.PodList

//go:generate gotemplate "./templates/resource" "podsClass(\"Pods\", pods, podListTypeAlias)"

func (c *podsClass) getList(ns string, listOptions meta.ListOptions) (*kapi.PodList, error) {
	return c.rk.clientset.Core().Pods(ns).List(listOptions)
}

func (c *podsClass) getItem(pods podListTypeAlias, index int) (*podClassInstance, error) {
	newPodObj, err := c.rk.classes.Pod.New()
	if err != nil {
		return nil, err
	}
	pod := pods.Items[index]
	newPodObj.vars.pod = podTypeAlias(pod)
	return newPodObj, nil
}

//go:generate gotemplate "./templates/resource/list" "podsListModule(podsClass, \"Pods\", pods, podListTypeAlias)"

func (c *podsClass) defineOwnMethods() {
	c.defineListMethods()
	c.rk.appendMethods(c.class, map[string]methodDefintion{
		"logs": {
			mruby.ArgsNone(), func(m *mruby.Mrb, self *mruby.MrbValue) (mruby.Value, mruby.Value) {
				vars, err := c.LookupVars(self)
				if err != nil {
					return nil, createException(m, err.Error())
				}

				newPodLogsObj, err := c.rk.classes.PodLogs.New()
				if err != nil {
					return nil, createException(m, err.Error())
				}
				newPodLogsObj.vars.pods = vars.pods.Items
				return callWithException(m, newPodLogsObj.self, "get!")
			},
			instanceMethod,
		},
	})
}

func (o *podsClassInstance) Update(args ...*mruby.MrbValue) (mruby.Value, error) {
	return call(o.self, "get!", args...)
}
