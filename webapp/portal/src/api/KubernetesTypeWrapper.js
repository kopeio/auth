// @flow
import Kubernetes from "./Kubernetes";

class KubernetesTypeWrapper {
  constructor(group: string, version: string, kind: string, namespace: string) {
    this.group = group;
    this.version = version;
    this.kind = kind;
    this.namespace = namespace;
  }

  _url(name: string) : string {
    var u = Kubernetes.url(this.group, this.version) + "namespaces/" + this.namespace + "/" + this.kind + "/";
    if (name) {
      u += name;
    }
    return u;
  }

  list() {
    return fetch(this._url())
      .then(response => response.json());
  };

  get(name: string) {
    return fetch(this._url(name))
      .then(response => response.json());
  };

  put(name: string, data) {
    return fetch(this._url(name), {
      method: 'PUT',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data)
    });
  };
}

export default KubernetesTypeWrapper;