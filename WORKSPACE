git_repository(
    name = "io_bazel_rules_go",
    remote = "https://github.com/bazelbuild/rules_go.git",
    tag = "0.5.4",
)

load("@io_bazel_rules_go//go:def.bzl", "go_repositories", "go_repository")

go_repositories()

#============================================================================
git_repository(
    name = "io_bazel_rules_docker",
    remote = "https://github.com/bazelbuild/rules_docker.git",
    tag = "v0.0.1",
)

load(
  "@io_bazel_rules_docker//docker:docker.bzl",
  "docker_repositories"
)
docker_repositories()

#=============================================================================

git_repository(
    name = "org_pubref_rules_protobuf",
    remote = "https://github.com/pubref/rules_protobuf.git",
    tag = "v0.7.2",
)

load("@org_pubref_rules_protobuf//go:rules.bzl", "go_proto_repositories")

go_proto_repositories()

#=============================================================================

# for building docker base images
debs = (
    (
        "deb_busybox",
        "5f81f140777454e71b9e5bfdce9c89993de5ddf4a7295ea1cfda364f8f630947",
        "http://ftp.us.debian.org/debian/pool/main/b/busybox/busybox-static_1.22.0-19+b3_amd64.deb",
    ),
    (
        "deb_libc",
        "b3f7278d80d5d0dc428fe92309bbc0e0a1ed665548a9f660663c1e1151335ae9",
        "http://ftp.us.debian.org/debian/pool/main/g/glibc/libc6_2.24-11+deb9u1_amd64.deb",
    ),
)

[http_file(
    name = name,
    sha256 = sha256,
    url = url,
) for name, sha256, url in debs]

#=============================================================================
# client-go

go_repository(
    name = "io_k8s_client_go",
    commit = "0389c75147549613c776c45d2de9511339b0c072",
    importpath = "k8s.io/client-go",
)

go_repository(
    name = "io_k8s_apiserver",
    commit = "f0eaebe542e55a09c7301fb1da38694f433d9b72",
    importpath = "k8s.io/apiserver",
)

go_repository(
    name = "io_k8s_apimachinery",
    commit = "84c15da65eb86243c295d566203d7689cc6ac04b",
    importpath = "k8s.io/apimachinery",
)

go_repository(
    name = "com_github_PuerkitoBio_purell",
    commit = "8a290539e2e8629dbc4e6bad948158f790ec31f4",
    importpath = "github.com/PuerkitoBio/purell",
)

go_repository(
    name = "com_github_PuerkitoBio_urlesc",
    commit = "5bd2802263f21d8788851d5305584c82a5c75d7e",
    importpath = "github.com/PuerkitoBio/urlesc",
)

go_repository(
    name = "com_github_coreos_go_oidc",
    commit = "5644a2f50e2d2d5ba0b474bc5bc55fea1925936d",
    importpath = "github.com/coreos/go-oidc",
)

go_repository(
    name = "com_github_coreos_pkg",
    commit = "fa29b1d70f0beaddd4c7021607cc3c3be8ce94b8",
    importpath = "github.com/coreos/pkg",
)

go_repository(
    name = "com_github_davecgh_go_spew",
    commit = "5215b55f46b2b919f50a1df0eaa5886afe4e3b3d",
    importpath = "github.com/davecgh/go-spew",
)

go_repository(
    name = "com_github_docker_distribution",
    commit = "cd27f179f2c10c5d300e6d09025b538c475b0d51",
    importpath = "github.com/docker/distribution",
)

go_repository(
    name = "com_github_emicklei_go_restful",
    commit = "09691a3b6378b740595c1002f40c34dd5f218a22",
    importpath = "github.com/emicklei/go-restful",
)

go_repository(
    name = "com_github_ghodss_yaml",
    commit = "73d445a93680fa1a78ae23a5839bad48f32ba1ee",
    importpath = "github.com/ghodss/yaml",
)

go_repository(
    name = "com_github_go_openapi_jsonpointer",
    commit = "46af16f9f7b149af66e5d1bd010e3574dc06de98",
    importpath = "github.com/go-openapi/jsonpointer",
)

go_repository(
    name = "com_github_go_openapi_analysis",
    commit = "b44dc874b601d9e4e2f6e19140e794ba24bead3b",
    importpath = "github.com/go-openapi/analysis",
)

go_repository(
    name = "com_github_go_openapi_jsonreference",
    commit = "13c6e3589ad90f49bd3e3bbe2c2cb3d7a4142272",
    importpath = "github.com/go-openapi/jsonreference",
)

go_repository(
    name = "com_github_go_openapi_loads",
    commit = "18441dfa706d924a39a030ee2c3b1d8d81917b38",
    importpath = "github.com/go-openapi/loads",
)

go_repository(
    name = "com_github_go_openapi_spec",
    commit = "6aced65f8501fe1217321abf0749d354824ba2ff",
    importpath = "github.com/go-openapi/spec",
)

go_repository(
    name = "com_github_go_openapi_swag",
    commit = "1d0bd113de87027671077d3c71eb3ac5d7dbba72",
    importpath = "github.com/go-openapi/swag",
)

go_repository(
    name = "com_github_gogo_protobuf",
    commit = "c0656edd0d9eab7c66d1eb0c568f9039345796f7",
    importpath = "github.com/gogo/protobuf",
)

go_repository(
    name = "com_github_golang_glog",
    commit = "44145f04b68cf362d9c4df2182967c2275eaefed",
    importpath = "github.com/golang/glog",
)

go_repository(
    name = "com_github_golang_groupcache",
    commit = "02826c3e79038b59d737d3b1c0a1d937f71a4433",
    importpath = "github.com/golang/groupcache",
)

go_repository(
    name = "com_github_golang_protobuf",
    commit = "8616e8ee5e20a1704615e6c8d7afcdac06087a67",
    importpath = "github.com/golang/protobuf",
)

go_repository(
    name = "com_github_google_gofuzz",
    commit = "44d81051d367757e1c7c6a5a86423ece9afcf63c",
    importpath = "github.com/google/gofuzz",
)

go_repository(
    name = "com_github_howeyc_gopass",
    commit = "3ca23474a7c7203e0a0a070fd33508f6efdb9b3d",
    importpath = "github.com/howeyc/gopass",
)

go_repository(
    name = "com_github_imdario_mergo",
    commit = "6633656539c1639d9d78127b7d47c622b5d7b6dc",
    importpath = "github.com/imdario/mergo",
)

go_repository(
    name = "com_github_jonboulle_clockwork",
    commit = "72f9bd7c4e0c2a40055ab3d0f09654f730cce982",
    importpath = "github.com/jonboulle/clockwork",
)

go_repository(
    name = "com_github_juju_ratelimit",
    commit = "77ed1c8a01217656d2080ad51981f6e99adaa177",
    importpath = "github.com/juju/ratelimit",
)

go_repository(
    name = "com_github_mailru_easyjson",
    commit = "d5b7844b561a7bc640052f1b935f7b800330d7e0",
    importpath = "github.com/mailru/easyjson",
)

go_repository(
    name = "com_github_pmezard_go_difflib",
    commit = "d8ed2627bdf02c080bf22230dbb337003b7aba2d",
    importpath = "github.com/pmezard/go-difflib",
)

go_repository(
    name = "com_github_spf13_pflag",
    commit = "5ccb023bc27df288a957c5e994cd44fd19619465",
    importpath = "github.com/spf13/pflag",
)

go_repository(
    name = "com_github_stretchr_testify",
    commit = "e3a8ff8ce36581f87a15341206f205b1da467059",
    importpath = "github.com/stretchr/testify",
)

go_repository(
    name = "com_github_ugorji_go",
    commit = "ded73eae5db7e7a0ef6f55aace87a2873c5d2b74",
    importpath = "github.com/ugorji/go",
)

go_repository(
    name = "com_google_cloud_go_compute_metadata",
    commit = "3b1ae45394a234c385be014e9a488f2bb6eef821",
    importpath = "cloud.google.com/go/compute/metadata",
)

go_repository(
    name = "com_google_cloud_go_internal",
    commit = "3b1ae45394a234c385be014e9a488f2bb6eef821",
    importpath = "cloud.google.com/go/internal",
)

go_repository(
    name = "in_gopkg_inf_v0",
    commit = "3887ee99ecf07df5b447e9b00d9c0b2adaa9f3e4",
    importpath = "gopkg.in/inf.v0",
)

go_repository(
    name = "in_gopkg_yaml_v2",
    commit = "53feefa2559fb8dfa8d81baad31be332c97d6c77",
    importpath = "gopkg.in/yaml.v2",
)

go_repository(
    name = "org_golang_google_appengine",
    commit = "4f7eeb5305a4ba1966344836ba4af9996b7b4e05",
    importpath = "google.golang.org/appengine",
)

go_repository(
    name = "org_golang_google_appengine_internal",
    commit = "4f7eeb5305a4ba1966344836ba4af9996b7b4e05",
    importpath = "google.golang.org/appengine/internal",
)

go_repository(
    name = "org_golang_google_appengine_internal_app_identity",
    commit = "4f7eeb5305a4ba1966344836ba4af9996b7b4e05",
    importpath = "google.golang.org/appengine/internal/app_identity",
)

go_repository(
    name = "org_golang_google_appengine_internal_base",
    commit = "4f7eeb5305a4ba1966344836ba4af9996b7b4e05",
    importpath = "google.golang.org/appengine/internal/base",
)

go_repository(
    name = "org_golang_google_appengine_internal_datastore",
    commit = "4f7eeb5305a4ba1966344836ba4af9996b7b4e05",
    importpath = "google.golang.org/appengine/internal/datastore",
)

go_repository(
    name = "org_golang_google_appengine_internal_log",
    commit = "4f7eeb5305a4ba1966344836ba4af9996b7b4e05",
    importpath = "google.golang.org/appengine/internal/log",
)

go_repository(
    name = "org_golang_google_appengine_internal_modules",
    commit = "4f7eeb5305a4ba1966344836ba4af9996b7b4e05",
    importpath = "google.golang.org/appengine/internal/modules",
)

go_repository(
    name = "org_golang_google_appengine_internal_remote_api",
    commit = "4f7eeb5305a4ba1966344836ba4af9996b7b4e05",
    importpath = "google.golang.org/appengine/internal/remote_api",
)

go_repository(
    name = "org_golang_x_crypto",
    commit = "d172538b2cfce0c13cee31e647d0367aa8cd2486",
    importpath = "golang.org/x/crypto",
)

go_repository(
    name = "org_golang_x_net",
    commit = "e90d6d0afc4c315a0d87a568ae68577cc15149a0",
    importpath = "golang.org/x/net",
)

go_repository(
    name = "org_golang_x_oauth2",
    commit = "3c3a985cb79f52a3190fbc056984415ca6763d01",
    importpath = "golang.org/x/oauth2",
)

go_repository(
    name = "org_golang_x_sys",
    commit = "8f0908ab3b2457e2e15403d3697c9ef5cb4b57a9",
    importpath = "golang.org/x/sys",
)

go_repository(
    name = "org_golang_x_text",
    commit = "2910a502d2bf9e43193af9d68ca516529614eed3",
    importpath = "golang.org/x/text",
)

#=============================================================================
# other deps

go_repository(
    name = "com_github_18F_hmacauth",
    commit = "9232a6386b737d7d1e5c1c6e817aa48d5d8ee7cd",
    importpath = "github.com/18F/hmacauth",
)

go_repository(
    name = "org_golang_google_api",
    commit = "650535c7d6201e8304c92f38c922a9a3a36c6877",
    importpath = "google.golang.org/api",
)

go_repository(
    name = "com_google_cloud_go",
    commit = "dbe4740b523eecbc19b2050f0691772c312aa07b",
    importpath = "cloud.google.com/go",
)

go_repository(
    name = "com_github_googleapis_gax_go",
    commit = "8c5154c0fe5bf18cf649634d4c6df50897a32751",
    importpath = "github.com/googleapis/gax-go",
)

go_repository(
    name = "com_github_coreos_etcd",
    commit = "cc198e22d3b8fd7ec98304c95e68ee375be54589",
    importpath = "github.com/coreos/etcd",
)

go_repository(
    name = "com_github_pborman_uuid",
    commit = "ca53cad383cad2479bbba7f7a1a05797ec1386e4",
    importpath = "github.com/pborman/uuid",
)

go_repository(
    name = "com_github_prometheus_client_golang",
    commit = "e51041b3fa41cece0dca035740ba6411905be473",
    importpath = "github.com/prometheus/client_golang",
)

go_repository(
    name = "com_github_prometheus_client_model",
    commit = "fa8ad6fec33561be4280a8f0514318c79d7f6cb6",
    importpath = "github.com/prometheus/client_model",
)

go_repository(
    name = "com_github_prometheus_common",
    commit = "ffe929a3f4c4faeaa10f2b9535c2b1be3ad15650",
    importpath = "github.com/prometheus/common",
)

go_repository(
    name = "com_github_prometheus_procfs",
    commit = "454a56f35412459b5e684fd5ec0f9211b94f002a",
    importpath = "github.com/prometheus/procfs",
)

go_repository(
    name = "com_github_beorn7_perks",
    commit = "3ac7bf7a47d159a033b107610db8a1b6575507a4",
    importpath = "github.com/beorn7/perks",
)

go_repository(
    name = "org_bitbucket_ww_goautoneg",
    commit = "75cd24fc2f2c2a2088577d12123ddee5f54e0675",
    importpath = "bitbucket.org/ww/goautoneg",
)

go_repository(
    name = "com_github_matttproud_golang_protobuf_extensions",
    commit = "c12348ce28de40eed0136aa2b644d0ee0650e56c",
    importpath = "github.com/matttproud/golang_protobuf_extensions",
)

go_repository(
    name = "com_github_grpc_ecosystem_grpc_gateway",
    commit = "f52d055dc48aec25854ed7d31862f78913cf17d1",
    importpath = "github.com/grpc-ecosystem/grpc-gateway",
)

go_repository(
    name = "com_github_coreos_go_systemd",
    commit = "48702e0da86bd25e76cfef347e2adeb434a0d0a6",
    importpath = "github.com/coreos/go-systemd",
)

go_repository(
    name = "com_github_pkg_errors",
    commit = "a22138067af1c4942683050411a841ade67fe1eb",
    importpath = "github.com/pkg/errors",
)

go_repository(
    name = "in_gopkg_natefinch_lumberjack_v2",
    commit = "20b71e5b60d756d3d2f80def009790325acc2b23",
    importpath = "gopkg.in/natefinch/lumberjack.v2",
)

go_repository(
    name = "com_github_elazarl_go_bindata_assetfs",
    commit = "3dcc96556217539f50599357fb481ac0dc7439b9",
    importpath = "github.com/elazarl/go-bindata-assetfs",
)

go_repository(
    name = "com_github_evanphx_json_patch",
    commit = "ba18e35c5c1b36ef6334cad706eb681153d2d379",
    importpath = "github.com/evanphx/json-patch",
)

go_repository(
    name = "org_golang_google_grpc",
    commit = "231b4cfea0e79843053a33f5fe90bd4d84b23cd3",
    importpath = "google.golang.org/grpc",
)
