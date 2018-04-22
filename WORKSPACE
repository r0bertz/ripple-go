http_archive(
    name = "io_bazel_rules_go",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.9.0/rules_go-0.9.0.tar.gz",
    sha256 = "4d8d6244320dd751590f9100cf39fd7a4b75cd901e1f3ffdfd6f048328883695",
)
http_archive(
    name = "bazel_gazelle",
    url = "https://github.com/bazelbuild/bazel-gazelle/releases/download/0.8/bazel-gazelle-0.8.tar.gz",
    sha256 = "e3dadf036c769d1f40603b86ae1f0f90d11837116022d9b06e4cd88cae786676",
)
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains", "go_repository")
go_rules_dependencies()
go_register_toolchains()
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
gazelle_dependencies()

go_repository(
    name = "com_github_gorilla_websocket",
    tag = "v1.2.0",
    importpath = "github.com/gorilla/websocket",
)

go_repository(
    name = "com_github_mitchellh_mapstructure",
    commit = "b4575eea38cca1123ec2dc90c26529b5c5acfcff",
    importpath = "github.com/mitchellh/mapstructure",
)

go_repository(
    name = "com_github_rubblelabs_ripple",
    commit = "0e03ed41baca64e37507128bc593822292b3349b",
    importpath = "github.com/rubblelabs/ripple",
)

go_repository(
    name = "com_github_r0bertz_ripple",
    commit = "HEAD",
    importpath = "github.com/r0bertz/ripple",
)

go_repository(
    name = "com_github_willf_bitset",
    tag = "v1.1.3",
    importpath = "github.com/willf/bitset",
)

go_repository(
    name = "com_github_agl_ed25519",
    commit = "5312a61534124124185d41f09206b9fef1d88403",
    importpath = "github.com/agl/ed25519",
)

go_repository(
    name = "com_github_btcsuite_btcd",
    commit = "2be2f12b358dc57d70b8f501b00be450192efbc3",
    importpath = "github.com/btcsuite/btcd",
)

go_repository(
    name = "org_golang_x_crypto",
    importpath = "golang.org/x/crypto",
    strip_prefix="crypto-81e90905daefcd6fd217b62423c0908922eadb30",
    type="zip",
    urls=['https://codeload.github.com/golang/crypto/zip/81e90905daefcd6fd217b62423c0908922eadb30'],
)

go_repository(
    name = "com_github_golang_collections_collections",
    commit = "604e922904d35e97f98a774db7881f049cd8d970",
    importpath = "github.com/golang-collections/collections",
)
