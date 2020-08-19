class Earthly < Formula
  desc "Build automation tool for the post-container era"
  homepage "https://earthly.dev/"
  url ":URL:"
  sha256 ":SHA256:"
  license "MPL-2.0"
  head "https://github.com/earthly/earthly.git"

  bottle do
    cellar :any_skip_relocation
    sha256 "5929dc77833f8c3df5401ac67fb3cd6f519b9f24152d6c8e52d8675919826934" => :catalina
    sha256 "b6d2e459b9d8d12033899bb16acb33b5db19119f8bb198983ecb051bb367ab2a" => :mojave
    sha256 "5d6412cc573aebf07cf98b90d755964747d4868135f2ca721761d30e7133d53d" => :high_sierra
  end

  depends_on "go" => :build

  def install
    system "go", "build",
        "-tags", ":TAGS:",
        "-ldflags", ":LDFLAGS:",
        *std_go_args,
        "-o", bin/"earth",
        "./cmd/earth/main.go"
  end

  test do
    (testpath/"build.earth").write <<~EOS

      default:
      \tRUN echo Homebrew
    EOS

    output = shell_output("#{bin}/earth --buildkit-host 127.0.0.1 +default 2>&1", 1).strip
    assert_match "Error while dialing invalid address 127.0.0.1", output
  end
end
