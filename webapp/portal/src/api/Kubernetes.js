class Kubernetes {
  static url(group, version) {
    var kubernetesUrl = window.AppSettings.kubernetesUrl || ("https://" + window.host);
    if (!kubernetesUrl.endsWith("/")) {
      kubernetesUrl += "/";
    }
    return kubernetesUrl + "apis/" + group + "/" + version + "/";
  }
};

export default Kubernetes;