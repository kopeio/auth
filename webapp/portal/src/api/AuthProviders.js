// @flow
import KubernetesTypeWrapper from "./KubernetesTypeWrapper";

class AuthProviders extends KubernetesTypeWrapper {
  constructor() {
    super("config.auth.kope.io", "v1alpha1", "authproviders");
  }

  static build() : AuthProviders {
    return new AuthProviders();
  }
}

export default AuthProviders;