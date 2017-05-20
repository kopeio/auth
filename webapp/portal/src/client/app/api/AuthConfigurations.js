import KubernetesTypeWrapper from "./KubernetesTypeWrapper";

class AuthConfigurations extends KubernetesTypeWrapper {
  constructor(namespace) {
    super("config.auth.kope.io", "v1alpha1", "authconfigurations", namespace);
  }

  static namespace(namespace) {
    return new AuthConfigurations(namespace);
  }
}

export default AuthConfigurations;