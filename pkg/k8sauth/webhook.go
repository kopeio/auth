package k8sauth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang/glog"
	authenticationv1beta1 "k8s.io/client-go/pkg/apis/authentication/v1beta1"
	"kope.io/auth/pkg/tokenstore"
)

type Webhook struct {
	Tokenstore tokenstore.Interface
}

func (h *Webhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var review authenticationv1beta1.TokenReview
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		glog.Warningf("failed to decode body: %v", err)
		http.Error(w, fmt.Sprintf("failed to decode body: %v", err), http.StatusBadRequest)
		return
	}

	if review.APIVersion != authenticationv1beta1.SchemeGroupVersion.String() {
		glog.Warningf("unknown version: %v", review.APIVersion)
		http.Error(w, fmt.Sprintf("unknown version: %v", review.APIVersion), http.StatusBadRequest)
		return
	}

	resp := &authenticationv1beta1.TokenReview{}
	resp.APIVersion = authenticationv1beta1.SchemeGroupVersion.String()
	userInfo, err := h.Tokenstore.LookupToken(review.Spec.Token)
	if err != nil {
		glog.Warningf("error looking up token: %v", err)
		resp.Status.Authenticated = false
		resp.Status.Error = err.Error()
	} else if userInfo == nil {
		glog.V(2).Infof("Invalid token")
		resp.Status.Authenticated = false
	} else {
		glog.V(2).Infof("Valid token for user %q", userInfo.Username)
		resp.Status.Authenticated = true
		resp.Status.User = *userInfo
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		glog.Warningf("error writing response: %v", err)
	}
}
