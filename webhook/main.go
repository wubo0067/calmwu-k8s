/*
 * @Author: calm.wu
 * @Date: 2019-05-22 10:26:35
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-05-24 14:57:14
 */
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

// toAdmissionResponse is a helper function to create an AdmissionResponse
// with an embedded error
func toAdmissionResponse(err error) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}

//
func isKubeNamespace(ns string) bool {
	return ns == metav1.NamespacePublic || ns == metav1.NamespaceSystem
}

func webHookServerHandler(w http.ResponseWriter, r *http.Request) {
	klog.V(0).Info("Handling webhook request...")

	// step 1 请求校验，只能是Post请求，json content type
	var errStr string
	// 必须是post方法
	if r.Method != http.MethodPost {
		errStr = fmt.Sprintf("invalid method %s, only POST requests are allowed", r.Method)
		http.Error(w, errStr, http.StatusMethodNotAllowed)
		klog.Error(errStr)
		return
	}

	// 读取body数据
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	if len(body) == 0 {
		klog.Error("empty body")
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}

	// Content-Type 必须是json
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		errStr = fmt.Sprintf("Content-Type=%s, expect application/json", contentType)
		klog.Errorf(errStr)
		http.Error(w, errStr, http.StatusUnsupportedMediaType)
		return
	}

	// step 2 解析AdmissionReview request，并处理

	// 回应对象
	var admissionResponse *v1beta1.AdmissionResponse

	// request解包
	admissionReviewReq := v1beta1.AdmissionReview{}
	if _, _, err := universalDeserializer.Decode(body, nil, &admissionReviewReq); err != nil {
		errStr = fmt.Sprintf("could not deserialize request: %v", err)
		klog.Error(errStr)
		admissionResponse = &v1beta1.AdmissionResponse{
			Allowed: false,
			Result: &metav1.Status{
				Message: errStr,
			},
		}
	} else if admissionReviewReq.Request == nil {
		errStr = "malformed admission review: request is nil"
		klog.Error(errStr)
		admissionResponse = &v1beta1.AdmissionResponse{
			Allowed: false,
			Result: &metav1.Status{
				Message: errStr,
			},
		}
	} else {
		klog.Infof("request path:%s", r.URL.Path)
		// only for non-Kubernetes namespaces. For objects in Kubernetes namespaces
		if !isKubeNamespace(admissionReviewReq.Request.Namespace) {
			// 处理请求
			admissionResponse = serveMutate(&admissionReviewReq)
		} else {
			klog.Infof("no-Kubernetes namespaces req:%#v", admissionReviewReq)
		}
	}

	admissionReviewRes := v1beta1.AdmissionReview{}
	if admissionResponse != nil {
		klog.Infof("admissionResponse is not nil, use admissionResponse")
		admissionReviewRes.Response = admissionResponse
		admissionReviewRes.Response.UID = admissionReviewReq.Request.UID
	} else {
		klog.Infof("admissionResponse is nil, default allow passed!")
		admissionReviewRes.Response = &v1beta1.AdmissionResponse{
			UID:     admissionReviewReq.Request.UID,
			Allowed: true,
		}
	}

	klog.Infof("Response:%#v", admissionReviewRes.Response)

	resp, err := json.Marshal(admissionReviewRes)
	if err != nil {
		errStr = fmt.Sprintf("could not encode response: %v", err)
		klog.Error(errStr)
		http.Error(w, errStr, http.StatusInternalServerError)
	}
	klog.Info("------------ Ready to write reponse ... ------------")
	if _, err := w.Write(resp); err != nil {
		errStr = fmt.Sprintf("Can't write response: %v", err)
		klog.Error(errStr)
		http.Error(w, errStr, http.StatusInternalServerError)
	}
}

func main() {
	var config Config
	config.addFlags()
	flag.Parse()

	// 设置输出的log
	//fLog, _ := os.OpenFile("webhook.log", os.O_RDWR|os.O_CREATE, 0666)
	//defer fLog.Close()
	klog.SetOutputBySeverity("INFO", os.Stderr)
	klog.V(0).Info("webhook start running.....")

	http.HandleFunc("/mutate", webHookServerHandler)
	//http.HandleFunc("/pods", servePods)
	httpServer := &http.Server{
		Addr:      ":443",
		TLSConfig: configTLS(config),
	}

	httpServer.ListenAndServeTLS("", "")
}
