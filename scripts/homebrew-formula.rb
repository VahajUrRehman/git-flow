# GitFlow TUI Homebrew Formula
# Usage: Copy this to a homebrew-tap repo and update the URLs/SHA256

class GitflowTui < Formula
  desc "Terminal-based Git management application with beautiful TUI"
  homepage "https://github.com/gitflow/tui"
  version "0.1.0"
  
  if OS.mac? && Hardware::CPU.intel?
    url "https://github.com/gitflow/tui/releases/download/v#{version}/gitflow-tui-darwin-amd64.tar.gz"
    sha256 "DARWIN_AMD64_SHA256"
  elsif OS.mac? && Hardware::CPU.arm?
    url "https://github.com/gitflow/tui/releases/download/v#{version}/gitflow-tui-darwin-arm64.tar.gz"
    sha256 "DARWIN_ARM64_SHA256"
  elsif OS.linux? && Hardware::CPU.intel?
    url "https://github.com/gitflow/tui/releases/download/v#{version}/gitflow-tui-linux-amd64.tar.gz"
    sha256 "LINUX_AMD64_SHA256"
  elsif OS.linux? && Hardware::CPU.arm?
    url "https://github.com/gitflow/tui/releases/download/v#{version}/gitflow-tui-linux-arm64.tar.gz"
    sha256 "LINUX_ARM64_SHA256"
  end

  def install
    bin.install "gitflow-tui"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/gitflow-tui --version")
  end
end
