class Kubernetes {
  static url(group, version) {
    var kubernetesUrl = window.AppSettings.kubernetesUrl || "https://userapimock.useast1.k8s.justinsb.com";
    if (!kubernetesUrl.endsWith("/")) {
      kubernetesUrl += "/";
    }
    return  kubernetesUrl + "apis/" + group + "/" + version + "/";
  }
};

export default Kubernetes;