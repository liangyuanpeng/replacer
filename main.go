package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	admissionv1 "k8s.io/api/admission/v1"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	v1 "k8s.io/kubernetes/pkg/apis/core/v1"
)

const (
	NA = "not_available"
)

var (
	runtimeScheme = runtime.NewScheme()
	// codecs        = serializer.NewCodecFactory(runtimeScheme)
	// deserializer  = codecs.UniversalDeserializer()
	// (https://github.com/kubernetes/kubernetes/issues/57982)
	// defaulter = runtime.ObjectDefaulter(runtimeScheme)
)

func init() {
	_ = corev1.AddToScheme(runtimeScheme)
	_ = admissionregistrationv1.AddToScheme(runtimeScheme)
	// defaulting with webhooks:
	// https://github.com/kubernetes/kubernetes/issues/57982
	_ = v1.AddToScheme(runtimeScheme)
}

func main() {
	router := gin.Default()
	router.POST("/mutate", mutateHandler)

	cert := flag.String("cert", "", "cert file ")
	key := flag.String("key", "", "key file")
	flag.Parse()

	log.Printf("tls file:%s|%s", *cert, *key)

	if *cert != "" && *key != "" {
		log.Print("run tls")
		router.RunTLS(":8181", *cert, *key)
	} else {
		log.Printf("run not tls")
		router.Run(":8181")
	}

}

func mutateHandler(c *gin.Context) {
	fmt.Println("begin mutateHandler")
	ar := admissionv1.AdmissionReview{}
	err := c.BindJSON(&ar)
	if err != nil {
		log.Println("parse.request.failed!", err)
	}
	b, err := json.Marshal(ar)
	if err != nil {
		fmt.Println("failed.to.Marshal:", err)
	} else {
		fmt.Println("data.size:", len(b))
		// fmt.Println("s:", string(b))
	}
	req := ar.Request
	var (
		resourceName string
	)
	fmt.Printf("AdmissionReview for Kind=%s, Namespace=%s Name=%s (%s) UID=%s patchOperation=%s UserInfo=%s \n",
		req.Kind, req.Namespace, req.Name, resourceName, req.UID, req.Operation, req.UserInfo)

	containers := make(map[int]string, 1)
	admissionReview := admissionv1.AdmissionReview{}

	switch req.Kind.Kind {
	case "Deployment":
		var deployment appsv1.Deployment
		if err := json.Unmarshal(req.Object.Raw, &deployment); err != nil {
			fmt.Printf("Could not unmarshal raw object: %s \n", err.Error())
			admissionReview.Response = &admissionv1.AdmissionResponse{
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
			admissionReview.Response.UID = ar.Request.UID

			c.JSON(200, &admissionReview)
			return
		}

		for i, c := range deployment.Spec.Template.Spec.Containers {
			fmt.Println("container.image:", c.Image)
			if strings.HasPrefix(c.Image, "lank8s.cn") || strings.HasPrefix(c.Image, "k8s.lank8s.cn") || strings.HasPrefix(c.Image, "gcr.lank8s.cn") {
				continue
			}
			if strings.HasPrefix(c.Image, "k8s.gcr.io") {
				containers[i] = strings.ReplaceAll(c.Image, "k8s.gcr.io", "lank8s.cn")
			} else if strings.HasPrefix(c.Image, "gcr.io") {
				containers[i] = strings.ReplaceAll(c.Image, "gcr.io", "gcr.lank8s.cn")
			}
		}
	case "StatefulSet":
		var statefulSet appsv1.StatefulSet
		if err := json.Unmarshal(req.Object.Raw, &statefulSet); err != nil {
			fmt.Printf("Could not unmarshal raw object: %s \n", err.Error())
			admissionReview.Response = &admissionv1.AdmissionResponse{
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
			admissionReview.Response.UID = ar.Request.UID

			c.JSON(200, &admissionReview)
			return
		}
		for i, c := range statefulSet.Spec.Template.Spec.Containers {
			fmt.Println("container.image:", c.Image)
			if strings.HasPrefix(c.Image, "lank8s.cn") || strings.HasPrefix(c.Image, "k8s.lank8s.cn") || strings.HasPrefix(c.Image, "gcr.lank8s.cn") {
				continue
			}
			if strings.HasPrefix(c.Image, "k8s.gcr.io") {
				containers[i] = strings.ReplaceAll(c.Image, "k8s.gcr.io", "lank8s.cn")
			} else if strings.HasPrefix(c.Image, "gcr.io") {
				containers[i] = strings.ReplaceAll(c.Image, "gcr.io", "gcr.lank8s.cn")
			}
		}
	case "DaemonSet":
		var daemonSet appsv1.DaemonSet
		if err := json.Unmarshal(req.Object.Raw, &daemonSet); err != nil {
			fmt.Printf("Could not unmarshal raw object: %s \n", err.Error())
			admissionReview.Response = &admissionv1.AdmissionResponse{
				Result: &metav1.Status{
					Message: err.Error(),
				},
			}
			admissionReview.Response.UID = ar.Request.UID

			c.JSON(200, &admissionReview)
			return
		}
		for i, c := range daemonSet.Spec.Template.Spec.Containers {
			fmt.Println("container.image:", c.Image)
			if strings.HasPrefix(c.Image, "lank8s.cn") || strings.HasPrefix(c.Image, "k8s.lank8s.cn") || strings.HasPrefix(c.Image, "gcr.lank8s.cn") {
				continue
			}
			if strings.HasPrefix(c.Image, "k8s.gcr.io") {
				containers[i] = strings.ReplaceAll(c.Image, "k8s.gcr.io", "lank8s.cn")
			} else if strings.HasPrefix(c.Image, "gcr.io") {
				containers[i] = strings.ReplaceAll(c.Image, "gcr.io", "gcr.lank8s.cn")
			}
		}
	case "Pod":
	}
	var patch []patchOperation
	fmt.Println("begin replace image", containers)
	patch = replaceImage(containers, patch)
	patchBytes, err := json.Marshal(patch)
	if err != nil {
		admissionReview.Response = &admissionv1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
		admissionReview.Response.UID = ar.Request.UID

		c.JSON(200, &admissionReview)
		return
	}

	fmt.Printf("AdmissionResponse: patch=%s\n", string(patchBytes))

	admissionReview.Response = &admissionv1.AdmissionResponse{
		Allowed: true,
		Patch:   patchBytes,
		PatchType: func() *admissionv1.PatchType {
			pt := admissionv1.PatchTypeJSONPatch
			return &pt
		}(),
	}
	admissionReview.APIVersion = ar.APIVersion
	admissionReview.Kind = ar.Kind
	admissionReview.Response.UID = ar.Request.UID

	c.JSON(200, &admissionReview)
}

func replaceImage(containers map[int]string, patch []patchOperation) []patchOperation {

	for i, c := range containers {
		str := strconv.Itoa(i)
		p := patchOperation{
			Op:    "replace",
			Path:  "/spec/template/spec/containers/" + str + "/image",
			Value: c,
		}
		fmt.Printf("add.patch:%s|%s\n", p.Path, p.Value)
		patch = append(patch, p)
	}
	return patch
}

type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}
