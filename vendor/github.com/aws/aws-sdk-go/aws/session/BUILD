load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "doc.go",
        "env_config.go",
        "session.go",
        "shared_config.go",
    ],
    importmap = "grail/go/src/vendor/github.com/aws/aws-sdk-go/aws/session",
    importpath = "github.com/aws/aws-sdk-go/aws/session",
    visibility = ["//visibility:public"],
    deps = [
        "//go/src/vendor/github.com/aws/aws-sdk-go/aws:go_default_library",
        "//go/src/vendor/github.com/aws/aws-sdk-go/aws/awserr:go_default_library",
        "//go/src/vendor/github.com/aws/aws-sdk-go/aws/client:go_default_library",
        "//go/src/vendor/github.com/aws/aws-sdk-go/aws/corehandlers:go_default_library",
        "//go/src/vendor/github.com/aws/aws-sdk-go/aws/credentials:go_default_library",
        "//go/src/vendor/github.com/aws/aws-sdk-go/aws/credentials/stscreds:go_default_library",
        "//go/src/vendor/github.com/aws/aws-sdk-go/aws/csm:go_default_library",
        "//go/src/vendor/github.com/aws/aws-sdk-go/aws/defaults:go_default_library",
        "//go/src/vendor/github.com/aws/aws-sdk-go/aws/endpoints:go_default_library",
        "//go/src/vendor/github.com/aws/aws-sdk-go/aws/request:go_default_library",
        "//go/src/vendor/github.com/go-ini/ini:go_default_library",
    ],
)
