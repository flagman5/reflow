load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "alt_exit.go",
        "doc.go",
        "entry.go",
        "exported.go",
        "formatter.go",
        "hooks.go",
        "json_formatter.go",
        "logger.go",
        "logrus.go",
        "terminal_bsd.go",
        "terminal_linux.go",
        "terminal_notwindows.go",
        "terminal_solaris.go",
        "terminal_windows.go",
        "text_formatter.go",
        "writer.go",
    ],
    importmap = "grail/go/src/vendor/github.com/Sirupsen/logrus",
    importpath = "github.com/Sirupsen/logrus",
    visibility = ["//visibility:public"],
    deps = select({
        "@io_bazel_rules_go//go/platform:solaris": [
            "//go/src/vendor/golang.org/x/sys/unix:go_default_library",
        ],
        "//conditions:default": [],
    }),
)
