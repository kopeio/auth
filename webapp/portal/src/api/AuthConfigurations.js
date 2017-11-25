// @flow
import KubernetesTypeWrapper from "./KubernetesTypeWrapper";

class AuthConfigurations extends KubernetesTypeWrapper {
  constructor() {
    super("config.auth.kope.io", "v1alpha1", "authconfigurations");
  }

  static build() : AuthConfigurations {
    return new AuthConfigurations();
  }
}

export default AuthConfigurations;