load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "aes_cbc.go",
        "aes_cbc_content_cipher.go",
        "aes_cbc_padder.go",
        "aes_gcm.go",
        "aes_gcm_content_cipher.go",
        "cipher.go",
        "cipher_builder.go",
        "cipher_util.go",
        "decryption_client.go",
        "doc.go",
        "encryption_client.go",
        "envelope.go",
        "hash_io.go",
        "helper.go",
        "key_handler.go",
        "kms_key_handler.go",
        "mat_desc.go",
        "padder.go",
        "pkcs7_padder.go",
        "strategy.go",
    ],
    importmap = "grail/go/src/vendor/github.com/aws/aws-sdk-go/service/s3/s3crypto",
    importpath = "github.com/aws/aws-sdk-go/service/s3/s3crypto",
    visibility = ["//visibility:public"],
    deps = [
        "//go/src/vendor/github.com/aws/aws-sdk-go/aws:go_default_library",
        "//go/src/vendor/github.com/aws/aws-sdk-go/aws/awserr:go_default_library",
        "//go/src/vendor/github.com/aws/aws-sdk-go/aws/client:go_default_library",
        "//go/src/vendor/github.com/aws/aws-sdk-go/aws/request:go_default_library",
        "//go/src/vendor/github.com/aws/aws-sdk-go/internal/sdkio:go_default_library",
        "//go/src/vendor/github.com/aws/aws-sdk-go/service/kms:go_default_library",
        "//go/src/vendor/github.com/aws/aws-sdk-go/service/kms/kmsiface:go_default_library",
        "//go/src/vendor/github.com/aws/aws-sdk-go/service/s3:go_default_library",
        "//go/src/vendor/github.com/aws/aws-sdk-go/service/s3/s3iface:go_default_library",
    ],
)
