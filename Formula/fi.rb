class Fi < Formula
  desc "Display detailed file information with tree-like formatting"
  homepage "https://github.com/oh-tarnished/fi"
  url "https://github.com/oh-tarnished/fi/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "" # Will be filled after creating a release
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "-o", bin/"fi"
  end

  test do
    system "#{bin}/fi", "--help"
  end
end
