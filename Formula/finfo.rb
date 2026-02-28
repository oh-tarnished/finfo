# typed: false
# frozen_string_literal: true

class Finfo < Formula
  desc "Display detailed file information with tree-like formatting"
  homepage "https://github.com/oh-tarnished/finfo"
  url "https://github.com/oh-tarnished/finfo/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "0019dfc4b32d63c1392aa264aed2253c1e0c2fb09216f8e2cc269bbfb8bb49b5"
  license "MIT"
  head "https://github.com/oh-tarnished/finfo.git", branch: "main"

  livecheck do
    url :stable
    regex(/^v?(\d+(?:\.\d+)+)$/i)
  end

  depends_on "go" => [:build, :test]

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w")
  end

  test do
    (testpath/"testfile").write "hello"
    assert_match "Path", shell_output("#{bin}/finfo #{testpath}/testfile")
  end
end
