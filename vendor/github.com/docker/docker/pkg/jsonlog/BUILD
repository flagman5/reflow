load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "jsonlog.go",
        "jsonlog_marshalling.go",
        "jsonlogbytes.go",
        "time_marshalling.go",
    ],
    importmap = "grail/go/src/vendor/github.com/docker/docker/pkg/jsonlog",
    importpath = "github.com/docker/docker/pkg/jsonlog",
    visibility = ["//visibility:public"],
)
