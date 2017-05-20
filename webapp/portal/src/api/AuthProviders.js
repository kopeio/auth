// @flow
import KubernetesTypeWrapper from "./KubernetesTypeWrapper";

class AuthProviders extends KubernetesTypeWrapper {
  constructor(namespace: string) {
    super("config.auth.kope.io", "v1alpha1", "authproviders", namespace);
  }

  static namespace(namespace: string) : AuthProviders {
    return new AuthProviders(namespace);
  }
}

export default AuthProviders;