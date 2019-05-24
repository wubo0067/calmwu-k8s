/*
 * @Author: calm.wu
 * @Date: 2019-05-22 10:26:35
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-05-22 18:51:39
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

// admitFunc is the type we use for all of our validators and mutators
type admitFunc func(v1beta1.AdmissionReview) *v1beta1.AdmissionResponse

// serve handles the http portion of a request prior to handing to an admit
// function
// func serveAdmitFunc(w http.ResponseWriter, r *http.Request, admit admitFunc) {
// 	var body []byte
// 	if r.Body != nil {
// 		if data, err := ioutil.ReadAll(r.Body); err == nil {
// 			body = data
// 		}
// 	}

// 	// verify the content type is accurate
// 	contentType := r.Header.Get("Content-Type")
// 	if contentType != "application/json" {
// 		klog.Errorf("contentType=%s, expect application/json", contentType)
// 		return
// 	}

// 	klog.V(0).Info(fmt.Sprintf("handling request: %s", body))

// 	// The AdmissionReview that was sent to the webhook
// 	requestedAdmissionReview := v1beta1.AdmissionReview{}

// 	// The AdmissionReview that will be returned
// 	responseAdmissionReview := v1beta1.AdmissionReview{}

// 	deserializer := codecs.UniversalDeserializer()
// 	if _, _, err := deserializer.Decode(body, nil, &requestedAdmissionReview); err != nil {
// 		klog.Error(err)
// 		responseAdmissionReview.Response = toAdmissionResponse(err)
// 	} else {
// 		// pass to admitFunc
// 		responseAdmissionReview.Response = admit(requestedAdmissionReview)
// 	}

// 	// Return the same UID
// 	responseAdmissionReview.Response.UID = requestedAdmissionReview.Request.UID

// 	klog.V(0).Info(fmt.Sprintf("sending response: %v", responseAdmissionReview.Response))

// 	respBytes, err := json.Marshal(responseAdmissionReview)
// 	if err != nil {
// 		klog.Error(err)
// 	}
// 	if _, err := w.Write(respBytes); err != nil {
// 		klog.Error(err)
// 	}
// }

func serveHandler(w http.ResponseWriter, r *http.Request) {
	klog.V(0).Info("Handling webhook request...")

	// body数据
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

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		klog.Errorf("Content-Type=%s, expect application/json", contentType)
		http.Error(w, "invalid Content-Type, expect `application/json`", http.StatusUnsupportedMediaType)
		return
	}

	var admissionResponse *v1beta1.AdmissionResponse

	requestedAdmissionReview := v1beta1.AdmissionReview{}
	if _, _, err := deserializer.Decode(body, nil, &requestedAdmissionReview); err != nil {
		klog.Errorf("Can't decode body: %v", err)
		admissionResponse = &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	} else {
		klog.Info(r.URL.Path)
		if r.URL.Path == "/mutate" {
			admissionResponse = serveMutate(&requestedAdmissionReview)
		}
	}

	responseAdmissionReview := v1beta1.AdmissionReview{}
	if admissionResponse != nil {
		responseAdmissionReview.Response = admissionResponse
		if requestedAdmissionReview.Request != nil {
			responseAdmissionReview.Response.UID = requestedAdmissionReview.Request.UID
		}
	}
	
	resp, err := json.Marshal(responseAdmissionReview)
	if err != nil {
		klog.Errorf("Can't encode response: %v", err)
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
	}
	klog.Infof("Ready to write reponse ...")
	if _, err := w.Write(resp); err != nil {
		klog.Errorf("Can't write response: %v", err)
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}	
}

func main() {
	var config Config
	config.addFlags()
	flag.Parse()

	// 设置输出的log
	fLog, _ := os.OpenFile("webhook.log", os.O_RDWR|os.O_CREATE, 0666)
	defer fLog.Close()

	klog.SetOutputBySeverity("INFO", fLog)

	klog.V(0).Info("webhook start running.....")

	http.HandleFunc("/mutate", serveHandler)
	//http.HandleFunc("/pods", servePods)
	httpServer := &http.Server{
		Addr:      ":443",
		TLSConfig: configTLS(config),
	}

	httpServer.ListenAndServeTLS("", "")
}
