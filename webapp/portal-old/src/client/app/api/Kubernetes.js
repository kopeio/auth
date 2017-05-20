class Kubernetes {
  static url(group, version) {
      return "https://userapimock.useast1.k8s.justinsb.com/apis/" + group + "/" + version + "/";
  }
};

export default Kubernetes;