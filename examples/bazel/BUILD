load("@io_bazel_rules_docker//java:image.bzl", "java_image")
load("@rules_java//java:defs.bzl", "java_binary")

java_binary(
    name = "ProjectRunner",
    srcs = glob(["src/main/java/com/example/*.java"]),
)

java_image(
    name = "ProjectRunnerImage",
    srcs = glob(["src/main/java/com/example/*.java"]),
    main_class = "com.example.ProjectRunner",
)
